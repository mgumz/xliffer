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
	"fmt"
	"os"

	"github.com/tealeg/xlsx"
)

type xlsxConverter struct {
	fileName     string
	skipRows     int
	sheetNumber  int
	keyColumn    int
	noteColumn   int
	sourceColumn int
	targetColumn int
	sourceLang   string
	targetLang   string
}

func init() {
	registeredConverters["from-xlsx"] = new(xlsxConverter)
}

func (x *xlsxConverter) Description() string {
	return "Converts an Excel sheet to XLIFF"
}

func (x *xlsxConverter) ParseArgs(base string, args []string) error {

	var fs = flag.NewFlagSet(base+" from-xlsx", flag.ExitOnError)

	fs.StringVar(&x.fileName, "in", "", "infile")
	fs.IntVar(&x.skipRows, "skipRows", 2, "number of rows to skip")
	fs.IntVar(&x.sheetNumber, "sheet", 1, "number of the sheet containing the translations")
	fs.IntVar(&x.keyColumn, "key-column", 3, "column holding the key / msgid")
	fs.IntVar(&x.sourceColumn, "source-col", 4, "column holding the source for the translation")
	fs.IntVar(&x.targetColumn, "target-col", 5, "column holding the target translation")
	fs.StringVar(&x.sourceLang, "source-lang", "en", "source language")
	fs.StringVar(&x.targetLang, "target-lang", "en", "target language")
	fs.IntVar(&x.noteColumn, "note-col", 0, "column holding notes (0 - not used)")

	return fs.Parse(args)
}

func (x *xlsxConverter) Convert() error {

	var xlFile, err = xlsx.OpenFile(x.fileName)
	if err != nil {
		return err
	}

	var sheet *xlsx.Sheet

	for s := range xlFile.Sheets {
		if s == (x.sheetNumber - 1) {
			sheet = xlFile.Sheets[s]
			break
		}
	}

	if sheet == nil {
		return fmt.Errorf("did not find sheet %d in %s\n",
			x.sheetNumber, x.fileName)
	}

	var doc = newXliffDoc(x.fileName, x.sourceLang)
	x.SheetToDoc(doc, sheet)

	var out []byte
	if out, err = xml.MarshalIndent(doc, "", "  "); err != nil {
		return err
	}

	os.Stdout.WriteString(xml.Header)
	os.Stdout.Write(out)

	return nil
}

func (x *xlsxConverter) SheetToDoc(doc *xliffDoc, sheet *xlsx.Sheet) {

	var sLang, tLang = x.sourceLang, x.targetLang
	var kCol, sCol, tCol = x.keyColumn - 1, x.sourceColumn - 1, x.targetColumn - 1
	var key, target, source string

	for r := range sheet.Rows {
		if r < (x.skipRows) {
			continue
		}

		var row = sheet.Rows[r]

		if len(row.Cells)-1 < kCol {
			continue
		}
		if key = row.Cells[kCol].String(); key == "" {
			continue
		}

		source, target = "", ""
		if sCol < len(row.Cells) {
			source = row.Cells[sCol].String()
		}
		if tCol < len(row.Cells) {
			target = row.Cells[tCol].String()
		}

		var unit = xliffTransUnit{
			ID:     key,
			Source: xliffTransUnitInner{source, sLang, "preserve", ""},
			Target: xliffTransUnitInner{target, tLang, "preserve", ""},
		}
		if (x.noteColumn > 0) && (x.noteColumn <= len(row.Cells)) {
			unit.Note = row.Cells[x.noteColumn-1].String()
		}

		doc.File[0].Body.TransUnit = append(doc.File[0].Body.TransUnit, unit)
	}
}
