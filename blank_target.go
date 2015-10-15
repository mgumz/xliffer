// This file is part of *xliffer*
//
// Copyright (C) 2015, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/xml"
	"flag"
	"io"
)

// blankTarget is converter which blanks the target translation. it's purpose
// is to take a fully translated .xliff and create new templates for other
// languages.
type blankTarget struct {
	inFile string
}

func init() {
	registeredConverters["blank-target"] = new(blankTarget)
}

func (b *blankTarget) Description() string {
	return "Blanks the target-attribute of a XLIFF"
}

func (b *blankTarget) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" blank-target", flag.ExitOnError)
	fs.StringVar(&b.inFile, "in", "", "infile")
	return fs.Parse(args)
}

func (b *blankTarget) Convert(w io.Writer) error {

	var doc, err = xliffFromFile(b.inFile)
	if err != nil {
		return err
	}

	for i := range doc.File {
		doc.File[i].TargetLang = ""
		for j := range doc.File[i].Body.TransUnit {
			doc.File[i].Body.TransUnit[j].Target.Lang = ""
			doc.File[i].Body.TransUnit[j].Target.Inner = ""
		}
	}

	var out []byte
	if out, err = xml.MarshalIndent(doc, "", "  "); err != nil {
		return err
	}

	io.WriteString(w, xml.Header)
	w.Write(out)

	return nil
}
