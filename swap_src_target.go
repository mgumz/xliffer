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

// swapSourceTarget is a converter which swaps source and target attributes.
// usefull when one has to work with .xliff files coming from sources not
// savy in using their xliff-editors correctly.
type swapSourceTarget struct {
	inFile string
}

func init() {
	registeredConverters["swap-source-target"] = new(swapSourceTarget)
}

func (s *swapSourceTarget) Description() string {
	return "Blanks the target-attribute of a XLIFF"
}

func (s *swapSourceTarget) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" swap-source-target", flag.ExitOnError)
	fs.StringVar(&s.inFile, "in", "", "infile")
	return fs.Parse(args)
}

func (s *swapSourceTarget) Convert(w io.Writer) error {

	var doc, err = xliffFromFile(s.inFile)
	if err != nil {
		return err
	}

	for i := range doc.File {
		doc.File[i].TargetLang = ""
		for j := range doc.File[i].Body.TransUnit {
			src := doc.File[i].Body.TransUnit[j].Source.Inner
			target := doc.File[i].Body.TransUnit[j].Target.Inner

			sname := doc.File[i].Body.TransUnit[j].Source.XMLName
			tname := doc.File[i].Body.TransUnit[j].Target.XMLName

			doc.File[i].Body.TransUnit[j].Source.Inner = target
			doc.File[i].Body.TransUnit[j].Target.Inner = src
			doc.File[i].Body.TransUnit[j].Source.XMLName = sname
			doc.File[i].Body.TransUnit[j].Target.XMLName = tname
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
