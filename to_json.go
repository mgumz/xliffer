// This file is part of *xliffer*
//
// Copyright (C) 2015, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"regexp"
)

// toJSON converts the target translation to a simple key:value
// structered JSON file
type toJSON struct {
	inFile   string
	pretty   bool
	keyMatch string
	keyTo    string
}

func init() {
	registeredConverters["to-json"] = new(toJSON)
}

func (tj *toJSON) Description() string {
	return "Converts XLIFF to JSON (key,value)"
}

func (tj *toJSON) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" to-json", flag.ExitOnError)
	fs.StringVar(&tj.inFile, "in", "", "infile")
	fs.StringVar(&tj.keyMatch, "key-match", "", "translate chars in key (regexp)")
	fs.StringVar(&tj.keyTo, "key-to", "", "chars of key gets translated to (string)")
	fs.BoolVar(&tj.pretty, "pretty", tj.pretty, "pretty print the resulting json")
	return fs.Parse(args)
}

func (tj *toJSON) Prepare() error {
	return nil
}

func (tj *toJSON) Convert(w io.Writer) error {

	var doc, err = xliffFromFile(tj.inFile)
	if err != nil {
		return err
	}

	var keyTrans = func(in string) string { return in }
	if tj.keyMatch != "" {
		rx, err := regexp.CompilePOSIX(tj.keyMatch)
		if err != nil {
			return err
		}
		keyTrans = func(in string) string {
			return rx.ReplaceAllString(in, tj.keyTo)
		}
	}

	var mappings = map[string]string{}
	for _, file := range doc.File {

		// note: no support for "groups" yet.
		for _, unit := range file.Body.TransUnit {

			unitID := keyTrans(unit.ID)

			if _, exist := mappings[unitID]; exist {
				log.Printf("warning: double entry for key %q", unitID)
			}

			mappings[unitID] = unit.Target.Inner
		}
	}

	var out []byte
	if tj.pretty {
		out, err = json.MarshalIndent(&mappings, "", "\t")
	} else {
		out, err = json.Marshal(&mappings)
	}

	if err == nil {
		w.Write(out)
	}

	return err
}
