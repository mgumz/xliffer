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
)

type dumpXLIFF struct {
	inFile string
}

func init() {
	registeredConverters["dump"] = new(dumpXLIFF)
}

func (d *dumpXLIFF) Description() string {
	return "Dumps XLIFF as parsed"
}

func (d *dumpXLIFF) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" dump", flag.ExitOnError)
	fs.StringVar(&d.inFile, "in", "", "infile")
	return fs.Parse(args)
}

func (d *dumpXLIFF) Prepare() error {
	return nil
}

func (d *dumpXLIFF) Convert(w io.Writer) error {

	var doc, err = xliffFromFile(d.inFile)
	if err != nil {
		return err
	}

	for _, file := range doc.File {

		for _, unit := range file.Body.TransUnit {

			fmt.Fprintf(w, "unit %s\n", unit.ID)
			fmt.Fprintf(w, " source: %v\n", unit.Source)
			fmt.Fprintf(w, " target: %v\n", unit.Target)

		}
	}

	return err
}
