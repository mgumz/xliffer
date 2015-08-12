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
	"log"
	"os"
)

type toJson struct {
	inFile string
	pretty bool
}

func init() {
	registeredConverters["to-json"] = new(toJson)
}

func (tj *toJson) Description() string {
	return "Converts XLIFF to JSON (key,value)"
}

func (tj *toJson) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" to-json", flag.ExitOnError)
	fs.StringVar(&tj.inFile, "in", "", "infile")
	fs.BoolVar(&tj.pretty, "pretty", tj.pretty, "pretty print the resulting json")
	return fs.Parse(args)
}

func (tj *toJson) Convert() error {

	var doc, err = xliffFromFile(tj.inFile)
	if err != nil {
		return err
	}

	var mappings = map[string]string{}
	for _, file := range doc.File {

		// note: no support for "groups" yet.
		for _, unit := range file.Body.TransUnit {

			if _, exist := mappings[unit.ID]; exist {
				log.Printf("warning: double entry for key %q", unit.ID)
			}

			mappings[unit.ID] = unit.Target.Inner
		}
	}

	var out []byte
	if tj.pretty {
		out, err = json.MarshalIndent(&mappings, "", "\t")
	} else {
		out, err = json.Marshal(&mappings)
	}

	if err == nil {
		os.Stdout.Write(out)
	}

	return err
}
