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

// setLang is simple converter to set the "lang=XX" attribute for
// each translation unit.
type setLang struct {
	inFile     string
	sourceLang string
	targetLang string
}

const _KEEP = "keep"

func init() {
	registeredConverters["set-lang"] = new(setLang)
}

func (s *setLang) Description() string {
	return "Sets the \"lang\" attribute of all translation units in a XLIFF"
}

func (s *setLang) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" set-lang", flag.ExitOnError)
	fs.StringVar(&s.inFile, "in", "", "infile")
	fs.StringVar(&s.targetLang, "target", _KEEP, "target language")
	fs.StringVar(&s.sourceLang, "source", _KEEP, "source language")
	return fs.Parse(args)
}

func (s *setLang) Convert(w io.Writer) error {

	var doc, err = xliffFromFile(s.inFile)
	if err != nil {
		return err
	}

	for i := range doc.File {

		setOrKeep(&doc.File[i].SourceLang, s.sourceLang)
		setOrKeep(&doc.File[i].SourceLang, s.sourceLang)

		for j := range doc.File[i].Body.TransUnit {
			setOrKeep(&doc.File[i].Body.TransUnit[j].Source.Lang, s.sourceLang)
			setOrKeep(&doc.File[i].Body.TransUnit[j].Target.Lang, s.targetLang)
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

func setOrKeep(to *string, from string) {
	if from != _KEEP {
		*to = from
	}
}
