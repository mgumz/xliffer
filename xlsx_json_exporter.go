// This file is part of *xliffer*
//
// Copyright (C) 2017, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"os"
	"path"
)

type xlsxJsonExporter struct {
	file   *os.File
	pretty bool
}

func (exp *xlsxJsonExporter) Open(folder, base string) error {
	name := path.Join(folder, base) + ".json"
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	exp.file = file
	return nil
}

func (exp *xlsxJsonExporter) Filename() string {
	return exp.file.Name()
}

func (exp *xlsxJsonExporter) Close() error {
	return exp.file.Close()
}

func (exp *xlsxJsonExporter) Write(units []xlsxTransUnit) error {

	var (
		err     error
		buf     []byte
		keyVals = make(map[string]string)
	)

	for _, unit := range units {
		keyVals[unit.Id] = unit.Target
	}

	if exp.pretty {
		buf, err = json.MarshalIndent(&keyVals, "", "\t")
	} else {
		buf, err = json.Marshal(&keyVals)
	}

	if err == nil {
		_, err = exp.file.Write(buf)
	}

	return err
}
