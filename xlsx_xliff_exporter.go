// This file is part of *xliffer*
//
// Copyright (C) 2017, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/xml"
	"os"
	"path"
)

type xlsxXliffExporter struct {
	file   *os.File
	pretty bool
}

func (exp *xlsxXliffExporter) Open(folder, base string) error {
	name := path.Join(folder, base) + ".xliff"
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	exp.file = file
	return nil
}

func (exp *xlsxXliffExporter) Filename() string {
	return exp.file.Name()
}

func (exp *xlsxXliffExporter) Close() error {
	return exp.file.Close()
}

func (exp *xlsxXliffExporter) Write(units []xlsxTransUnit) error {

	var (
		doc    = newXliffDoc("", units[0].SourceLang)
		indent = ""
		buf    []byte
		err    error
	)

	if exp.pretty {
		indent = "  "
	}

	body := &doc.File[0].Body
	for _, unit := range units {
		xliffUnit := xliffTransUnit{
			ID: unit.Id,
			Source: xliffTransUnitInner{
				Lang:  unit.SourceLang,
				Inner: unit.Source,
				Space: "preserve",
			},
			Target: xliffTransUnitInner{
				Lang:  unit.TargetLang,
				Inner: unit.Target,
				Space: "preserve",
			},
			Note: unit.Note,
		}
		body.TransUnit = append(body.TransUnit, xliffUnit)
	}

	if buf, err = xml.MarshalIndent(doc, "", indent); err != nil {
		return err
	}

	if err == nil {
		if _, err = exp.file.WriteString(xml.Header); err != nil {
			return err
		}
		_, err = exp.file.Write(buf)
	}

	return err
}
