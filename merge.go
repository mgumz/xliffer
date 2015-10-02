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
)

// mergeConv is converter which merges 2 .xliff files
type mergeConv struct {
	aFile string
	bFile string
}

func init() {
	registeredConverters["merge"] = new(mergeConv)
}

func (m *mergeConv) Description() string {
	return "Merge two XLIFF files"
}

func (m *mergeConv) ParseArgs(base string, args []string) error {
	var fs = flag.NewFlagSet(base+" merge", flag.ExitOnError)
	fs.StringVar(&m.aFile, "a", "", "a file")
	fs.StringVar(&m.bFile, "b", "", "b file")
	return fs.Parse(args)
}

func (m *mergeConv) Convert() error {

	var (
		err  error
		aDoc *xliffDoc
		bDoc *xliffDoc
		out  []byte
	)

	if aDoc, err = xliffFromFile(m.aFile); err != nil {
		return fmt.Errorf("%s: %s", m.aFile, err)
	}
	if bDoc, err = xliffFromFile(m.bFile); err != nil {
		return fmt.Errorf("%s: %s", m.bFile, err)
	}

	// TODO: check uniqueness of keys|translation unit ids?

	aDoc.File = append(aDoc.File, bDoc.File...)

	if out, err = xml.MarshalIndent(aDoc, "", "  "); err != nil {
		return err
	}

	os.Stdout.WriteString(xml.Header)
	os.Stdout.Write(out)

	return nil
}
