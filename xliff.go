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
	"os"
)

type xliffTransUnitInner struct {
	Inner string `xml:",chardata"`
	Lang  string `xml:"lang,attr"`
	Space string `xml:"space,attr,omitempty"`
}

type xliffTransUnit struct {
	ID     string              `xml:"id,attr"`
	Source xliffTransUnitInner `xml:"source"`
	Target xliffTransUnitInner `xml:"target"`
	Note   string              `xml:"note,omitempty"`
}

type xliffBody struct {
	XMLName   xml.Name         `xml:"body"`
	TransUnit []xliffTransUnit `xml:"trans-unit"`
}

type xliffFile struct {
	Original   string    `xml:"original,attr"`
	SourceLang string    `xml:"source-language,attr,omitempty"`
	TargetLang string    `xml:"target-language,attr,omitempty"`
	DataType   string    `xml:"datatype,attr,omitempty"`
	Body       xliffBody `xml:"body"`
}

type xliffDoc struct {
	XMLName xml.Name    `xml:"xliff"`
	Version string      `xml:"version,attr"`
	Xmlns   string      `xml:"xmlns,attr"`
	File    []xliffFile `xml:"file"`
}

func newXliffDoc(original, origLang string) *xliffDoc {

	var doc = new(xliffDoc)

	doc.Version = "1.2"
	doc.Xmlns = "urn:oasis:names:tc:xliff:document:1.2"
	doc.File = make([]xliffFile, 1)
	doc.File[0].Original = original
	doc.File[0].Body.TransUnit = []xliffTransUnit{}
	doc.File[0].SourceLang = origLang
	doc.File[0].DataType = "html"

	return doc
}

func xliffFromFile(fileName string) (*xliffDoc, error) {
	var f, err = os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var doc = new(xliffDoc)
	var dec = xml.NewDecoder(f)
	if err = dec.Decode(doc); err != nil {
		return nil, err
	}
	return doc, err
}
