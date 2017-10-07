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

// copyUnits copies the source translation units onto
// the target translation units
type copyUnits struct {
	inFile string
}

func init() {
	registeredConverters["copy"] = new(copyUnits)
}

func (c *copyUnits) Description() string {
	return "Copies SOURCE to TARGET units in a XLIFF"
}

func (c *copyUnits) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" copy", flag.ExitOnError)
	fs.StringVar(&c.inFile, "in", "", "infile")
	return fs.Parse(args)
}

func (c *copyUnits) Prepare() error {
	return nil
}

func (c *copyUnits) Convert(w io.Writer) error {

	var doc, err = xliffFromFile(c.inFile)
	if err != nil {
		return err
	}

	for i := range doc.File {

		doc.File[i].TargetLang = doc.File[i].SourceLang
		for j := range doc.File[i].Body.TransUnit {
			name := doc.File[i].Body.TransUnit[j].Target.XMLName
			doc.File[i].Body.TransUnit[j].Target = new(xliffTarget)
			doc.File[i].Body.TransUnit[j].Target.Copy(&doc.File[i].Body.TransUnit[j].Source)
			doc.File[i].Body.TransUnit[j].Target.XMLName = name
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
