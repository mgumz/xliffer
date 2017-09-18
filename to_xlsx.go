// This file is part of *xliffer*
//
// Copyright (C) 2017, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/tealeg/xlsx"
)

// toXLSX creates a .xlsx spreadsheet with the following table content
//   key | note            | source | ... | <target>
//   a.b | first entry     | hi     | ... | hey
//   b.c | important entry | cya    | ....| bye
type toXLSX struct {
	inFile       string
	appendFile   string
	appendSheet  string
	keyMatch     string
	keyTo        string
	keyColumn    int
	headRow      int
	targetColumn int

	xlFile  *xlsx.File
	xlSheet *xlsx.Sheet
}

func init() {
	registeredConverters["to-xlsx"] = new(toXLSX)
}

func (plugin *toXLSX) Description() string {
	return "Converts XLIFF to XLSX"
}

func (conv *toXLSX) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" to-xlsx", flag.ExitOnError)
	fs.StringVar(&conv.inFile, "in", "", "infile")
	fs.StringVar(&conv.appendFile, "append", "", ".xlsx file to append to")
	fs.StringVar(&conv.appendSheet, "sheet", "", "sheet of .xlsx to append to")
	fs.IntVar(&conv.headRow, "head-row", 0, "row which holds the header")
	fs.IntVar(&conv.keyColumn, "key-column", 0, "column holding the key / msgid")
	fs.StringVar(&conv.keyMatch, "key-match", "", "translate chars in key (regexp)")
	fs.StringVar(&conv.keyTo, "key-to", "", "chars of key gets translated to (string)")
	fs.IntVar(&conv.targetColumn, "target-column", -1, "column which will hold the translated text")
	return fs.Parse(args)
}

// in case of an "append to .xlsx" situation we need to open the .xlsx file
// before any convertion happens: the -o parameter from xliffer might point
// to the exact same .xlsx file which "truncates" the .xlsx file in order to
// provide the 'w' parameter of the .Convert() function. thus, we read in
// the .xlsx completely before and then it does not matter what happens to
// to the file
func (conv *toXLSX) Prepare() error {

	var (
		file  *xlsx.File
		sheet *xlsx.Sheet
		err   error
	)

	if file, sheet, err = conv.newOrAppend(); err != nil {
		return err
	}

	conv.xlFile = file
	conv.xlSheet = sheet

	return nil
}

func (conv *toXLSX) Convert(w io.Writer) error {

	var (
		err error
		doc *xliffDoc

		keyTrans = func(in string) string { return in }
	)

	if doc, err = xliffFromFile(conv.inFile); err != nil {
		return err
	}

	if conv.keyMatch != "" {
		rx, err := regexp.CompilePOSIX(conv.keyMatch)
		if err != nil {
			return err
		}
		keyTrans = func(in string) string {
			return rx.ReplaceAllString(in, conv.keyTo)
		}
	}

	// an empty sheet has no header-row
	conv.ensureRowExists(conv.xlSheet, conv.headRow)

	targetColumn := conv.getTargetColumn(conv.xlSheet, conv.headRow)

	existingKeys := conv.keyRowMap(conv.xlSheet, conv.headRow, conv.keyColumn)
	//fmt.Println("existing keys:", len(existingKeys))

	conv.createSheetHeader(conv.xlSheet, targetColumn)
	// 2 phases:
	//
	// phase-a - find existing keys and write the translations unit.Target at
	// the last column of the found row. if the given unit.ID is NOT found in
	// existingKeys container, the unit will be added to the end of the sheet
	// in phase-b

	type xlEntry struct {
		Key    string `xlsx:"0"`
		Note   string `xlsx:"1"`
		Source string `xlsx:"2"`
		Target string `xlsx:"3"`
	}
	appendix := make([]xlEntry, 0)

	for _, file := range doc.File {
		for _, unit := range file.Body.TransUnit {

			key := keyTrans(unit.ID)
			entry := xlEntry{key, unit.Note, unit.Source.Inner, unit.Target.Inner}
			row, exists := existingKeys[key]
			if !exists {
				appendix = append(appendix, entry)
				continue
			}

			conv.setCell(conv.xlSheet, row, targetColumn, entry.Target)
		}
	}

	// phase-b - write all entries in appendix to the end of the current sheet
	//
	if len(appendix) > 0 {

		cell := conv.xlSheet.Cell(len(conv.xlSheet.Rows)+1, conv.keyColumn+XLSX_NOTE_COLUMN)
		cell.SetStyle(conv.boldStyle())
		cell.SetString("added from " + conv.inFile)

		for _, entry := range appendix {
			row := len(conv.xlSheet.Rows)
			conv.setCell(conv.xlSheet, row, conv.keyColumn, entry.Key)
			conv.setCell(conv.xlSheet, row, conv.keyColumn+XLSX_NOTE_COLUMN, entry.Note)
			conv.setCell(conv.xlSheet, row, conv.keyColumn+XLSX_SOURCE_COLUMN, entry.Source)
			conv.setCell(conv.xlSheet, row, targetColumn, entry.Target)
		}
	}

	conv.xlFile.Write(w)

	return err
}

func (conv *toXLSX) newOrAppend() (*xlsx.File, *xlsx.Sheet, error) {

	var (
		file  *xlsx.File
		sheet *xlsx.Sheet
		err   error
	)

	if conv.appendFile == "" {

		file = xlsx.NewFile()
		sheetName := path.Base(conv.inFile)
		sheetName = sheetName[:XLSX_MAX_SHEETNAME]
		sheet, err = file.AddSheet(sheetName)
		if err != nil {
			return nil, nil, err
		}

	} else {

		f, err := xlsx.OpenFile(conv.appendFile)
		if err != nil {
			return nil, nil, err
		}
		file = f
		sheet = file.Sheets[0]
		if conv.appendSheet != "" {
			for i := range file.Sheets {
				if conv.appendSheet == file.Sheets[i].Name {
					sheet = file.Sheets[i]
				}
			}
		}
	}

	return file, sheet, nil
}

func (conv *toXLSX) keyRowMap(sheet *xlsx.Sheet, headRow, keyCol int) map[string]int {
	key2row := make(map[string]int)
	for i, row := range sheet.Rows {
		if i <= headRow {
			continue
		}

		cells := row.Cells
		if keyCol >= len(cells) {
			continue
		}
		key := strings.TrimSpace(cells[keyCol].String())

		if _, exists := key2row[key]; exists {
			fmt.Println("duplicate key detected:", key, ", ignoring.")
			continue
		}

		if key != "" {
			key2row[key] = i
		}
	}
	return key2row
}

func (conv *toXLSX) boldStyle() *xlsx.Style {
	font := xlsx.Font{}
	font.Name = xlsx.DefaultFont().Name
	font.Size = 2 * xlsx.DefaultFont().Size
	font.Bold = true
	style := xlsx.NewStyle()
	style.Font = font
	return style
}

func (conv *toXLSX) createSheetHeader(sheet *xlsx.Sheet, targetColumn int) {

	cell := sheet.Cell(conv.headRow, conv.keyColumn)
	keyStyle := cell.GetStyle()
	if cell.String() == "" {
		cell.SetString("key")
	}

	cell = sheet.Cell(conv.headRow, conv.keyColumn+XLSX_NOTE_COLUMN)
	if cell.String() == "" {
		cell.SetString("note")
	}

	cell = sheet.Cell(conv.headRow, conv.keyColumn+XLSX_SOURCE_COLUMN)
	if cell.String() == "" {
		cell.SetString("source")
	}

	cell = sheet.Cell(conv.headRow, targetColumn)
	cell.SetStyle(keyStyle)
	cell.SetString(path.Base(conv.inFile))
}

func (conv *toXLSX) setCell(sheet *xlsx.Sheet, row, col int, text string) {
	cell := sheet.Cell(row, col)
	cell.SetString(text)
}

func (conv *toXLSX) ensureRowExists(sheet *xlsx.Sheet, headRow int) {
	sheet.Cell(headRow, 0)
}

func (conv *toXLSX) getTargetColumn(sheet *xlsx.Sheet, headRow int) int {

	targetColumn := len(sheet.Rows[headRow].Cells)

	// an empty sheet has no cells ...
	if targetColumn < conv.keyColumn+XLSX_TARGET_COLUMN {
		targetColumn = conv.keyColumn + XLSX_TARGET_COLUMN
	}

	// the user called to-xlsx with a specific targetColumn
	if conv.targetColumn > -1 {
		targetColumn = conv.targetColumn
	}

	return targetColumn
}
