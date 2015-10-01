// This file is part of *xliffer*
//
// Copyright (C) 2015, Travelping GmbH <copyright@travelping.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

// a converter transforms an input file. multiple converters
// register themself to a registry and are then exposed as
// xliffer commands
type converter interface {
	// the Description is shown while printing the usage
	Description() string

	// parse the flags specific to the converter
	ParseArgs(base string, args []string) error

	// converts the specified input file
	Convert() error
}

var registeredConverters = make(map[string]converter)
