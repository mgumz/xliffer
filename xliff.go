// This file is part of *xliffer*
//
// Copyright (C) 2015, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
)

type xliffSource struct {
	XMLName xml.Name
	Inner   string `xml:",chardata"`
	Lang    string `xml:"lang,attr"`
	Space   string `xml:"space,attr,omitempty"`
	State   string `xml:"state,attr,omitempty"`
}

// xliffTarget might containt <mrk> tags which are leftovers from
// translation tools. as a result, "chardata" of a <target> node
// might be empty because all the translations are contained inside
// several <mrk>tags</mrk>. this is why we treat <target> similar to
// <source> but not equally.
type xliffTarget xliffSource

type xliffTransUnit struct {
	ID     string       `xml:"id,attr"`
	Source xliffSource  `xml:"source"`
	Target *xliffTarget `xml:"target,omitempty"`
	Note   string       `xml:"note,omitempty"`
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

	return xliffFromReader(f)
}

func xliffFromReader(r io.Reader) (*xliffDoc, error) {

	doc := new(xliffDoc)
	dec := xml.NewDecoder(r)

	if err := dec.Decode(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (to *xliffTarget) Copy(from *xliffSource) {
	to.XMLName = from.XMLName
	to.Inner = from.Inner
	to.Lang = from.Lang
	to.Space = from.Space
	to.State = from.State
}

// extract all chardata from a <target>-node, including all subnodes
func (target *xliffTarget) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	buf := bytes.NewBuffer(nil)

	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "target" {
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "lang":
						target.Lang = attr.Value
					case "space":
						target.Space = attr.Value
					case "state":
						target.State = attr.Value
					}
				}
			}
		case xml.CharData:
			buf.Write(t)
		}
	}

	target.Inner = buf.String()

	return nil
}
