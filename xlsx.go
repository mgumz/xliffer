// This file is part of *xliffer*
//
// Copyright (C) 2015, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

// converts a single master.xlsx into many .xliff files, ready
// to be shipped to translators.
//
// structure of the master.xlsx:
//
// name of sheet: PROJECT_NAME
//
// col0 - keys
// col1 - master-translation
// col2 - COUNTRY-LANG (iso 2 letter code)
//
// example master.xlsx:
//
// keys | note-column        | DE-de  | DE-en | EN-en | EN-x
// x    | greeting           | hallo  | hello | hello | hello
// y    | saying goodbye     | tschüß | bye   | bye   | bye
//
//
// from-xlsx creates translation .xliff files in the given directory:
//
// master.xlsx ->
//                 DE-de.xliff
//                 DE-en.xliff
//                 EN-en.xliff
//                 EN-x.xliff
//

const (
	XLSX_KEY_COLUMN    = 0
	XLSX_NOTE_COLUMN   = 1
	XLSX_SOURCE_COLUMN = 2
	XLSX_TARGET_COLUMN = 3

	XLSX_MAX_SHEETNAME = 30 // xlsx-limit
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

	destDir  string
	exporter xlsxExporter
}

type xlsxTransUnit struct {
	Id         string
	Source     string
	SourceLang string
	Target     string
	TargetLang string
	Note       string
}

type xlsxExporter interface {
	Open(folder, base string) error
	Write(units []xlsxTransUnit) error
	Close() error
	Filename() string
}

func init() {
	registeredConverters["from-xlsx"] = new(xlsxConverter)
}

func (x *xlsxConverter) Description() string {
	return "Converts an Excel sheet to XLIFF, JSON"
}

func (x *xlsxConverter) ParseArgs(base string, args []string) error {

	fs := flag.NewFlagSet(base+" from-xlsx", flag.ExitOnError)
	destType := "xliff"
	pretty := false

	fs.StringVar(&destType, "to", "json", "output format (xliff, json)")
	fs.BoolVar(&pretty, "pretty", pretty, "pretty print the output files")
	fs.StringVar(&x.fileName, "in", "", "infile")
	fs.IntVar(&x.skipRows, "skip-rows", 0, "number of rows to skip")
	fs.IntVar(&x.sheetNumber, "sheet", 1, "number of the sheet containing the translations")
	fs.IntVar(&x.keyColumn, "key-column", 0, "column holding the key / msgid")
	fs.IntVar(&x.sourceColumn, "source-col", -1, "column holding the source for the translation")
	fs.IntVar(&x.noteColumn, "note-col", -1, "column holding notes (0 - not used)")
	fs.IntVar(&x.targetColumn, "target-col", -1, "column holding the target translation")
	fs.StringVar(&x.sourceLang, "source-lang", "en", "source language")
	fs.StringVar(&x.targetLang, "target-lang", "en", "target language")
	fs.StringVar(&x.destDir, "dir", "", "output directory")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if x.destDir == "" {
		x.destDir, _ = os.Getwd()
	}
	if err = os.MkdirAll(x.destDir, 0777); err != nil {
		return err
	}

	switch destType {
	case "json":
		x.exporter = &xlsxJsonExporter{pretty: pretty}
	case "xliff":
		x.exporter = &xlsxXliffExporter{pretty: pretty}
	default:
		return fmt.Errorf("unsupported 'type': %q", destType)
	}

	return nil
}

func (conv *xlsxConverter) Convert(w io.Writer) error {

	var xlFile, err = xlsx.OpenFile(conv.fileName)
	if err != nil {
		return err
	}

	var sheet *xlsx.Sheet

	for s := range xlFile.Sheets {
		if s == (conv.sheetNumber - 1) {
			sheet = xlFile.Sheets[s]
			break
		}
	}

	if sheet == nil {
		return fmt.Errorf("did not find sheet %d in %s\n",
			conv.sheetNumber, conv.fileName)
	}

	keyRow, keyCol := conv.detectBounds(sheet)
	bodyRow := keyRow + 1

	srcCol := conv.sourceColumn
	if srcCol == -1 {
		// key | note | src | target-1 | target-2
		srcCol = keyCol + XLSX_SOURCE_COLUMN // TODO: detect "source" if not set
	}

	// the target column defines the column where the translated languages
	// start. if not defined, it defaults to the source column. why?
	// because that way the "initial" language also gets an export, either
	// to .xliff (which can then be handed to the translator team) or to
	// to .json and pretend it's already translated.
	targetCol := conv.targetColumn
	if targetCol == -1 {
		targetCol = srcCol // TODO: detect "target" automatically
	}

	head := sheet.Rows[keyRow].Cells

	//for x := range head[keyCol:] {
	//	fmt.Println(x, x+keyCol, head[x].String(), len(head), srcCol)
	//}

	srcLang := langFromCCLang(head[srcCol].String())
	rows := sheet.Rows

	for x := targetCol; x < len(rows[keyRow].Cells); x++ {

		cc_lang := head[x].String()
		lang := langFromCCLang(cc_lang)
		units := make([]xlsxTransUnit, 0, len(rows)-keyRow)

		for y := bodyRow; y < len(rows); y++ {
			row := rows[y]
			if (keyCol >= len(row.Cells)) || row.Cells[keyCol].String() == "" {
				continue
			}
			if (srcCol >= len(row.Cells)) || row.Cells[srcCol].String() == "" {
				continue
			}
			if (x >= len(row.Cells)) || row.Cells[x].String() == "" {
				continue
			}
			unit := conv.rowToTransUnit(y, keyCol, srcCol, x, srcLang, lang, rows)
			units = append(units, unit)
		}
		name := sheet.Name + "-" + cc_lang
		conv.exportUnits(conv.destDir, name, conv.exporter, units)
	}

	return nil
}

// extracts a row of the following form
//   key | comment | note | source | target-1 | target-2 | ...
// into a xlsxTransUnit
func (conv *xlsxConverter) rowToTransUnit(row, keyCol, sourceCol, targetCol int, sourceLang, targetLang string, rows []*xlsx.Row) xlsxTransUnit {
	unit := xlsxTransUnit{
		Id:         rows[row].Cells[keyCol].String(),
		Source:     rows[row].Cells[sourceCol].String(),
		SourceLang: sourceLang,
		Target:     rows[row].Cells[targetCol].String(),
		TargetLang: targetLang,
	}
	return unit
}

func (conv *xlsxConverter) exportUnits(dir, name string, x xlsxExporter, entries []xlsxTransUnit) {
	if err := x.Open(dir, name); err != nil {
		log.Printf("err: opening export file %q: %s", x.Filename(), err)
		return
	}
	defer x.Close()

	if err := x.Write(entries); err != nil {
		log.Printf("err: write entries to %q: %s", x.Filename(), err)
		return
	}

	log.Printf("written %d entries to %q: ok.", len(entries), x.Filename())
}

func langFromCCLang(cc_lang string) string {
	parts := strings.Split(cc_lang, "-")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// finds the key-column after x.skipRows, where the content starts
// aka "the header".
func (x *xlsxConverter) detectBounds(sheet *xlsx.Sheet) (int, int) {

	var cell *xlsx.Cell
	i, j := 0, 0

	for i = range sheet.Rows {
		if i < (x.skipRows) {
			continue
		}

		for j, cell = range sheet.Rows[i].Cells {
			if key := cell.String(); key != "" {
				return i, j
			}
		}
	}

	return -1, -1
}
