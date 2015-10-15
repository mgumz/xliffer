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
	"os"
	"path"
	"sort"
)

func main() {

	outFileName := flag.String("o", "-", "output file (default: \"-\"|stdout)")
	flag.Usage = usage
	flag.Parse()

	var args = flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(0)
	}

	var (
		conv  converter
		err   error
		exist bool
	)

	if conv, exist = registeredConverters[args[0]]; !exist {
		fmt.Fprintf(os.Stderr, "error: unknown command %v\n", args[0])
		flag.Usage()
		os.Exit(1)
	}

	if err = conv.ParseArgs(os.Args[0], args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: parsing %v\n", err)
		os.Exit(1)
	}

	var outWriter = os.Stdout
	if *outFileName != "" && *outFileName != "-" {
		outFile, err2 := os.Create(*outFileName)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "error: can't create %q: %v\n",
				*outFileName, err2)
			os.Exit(1)
		}
		defer outFile.Close()
		outWriter = outFile
	}

	if err := conv.Convert(outWriter); err != nil {
		fmt.Fprintf(os.Stderr, "error: converting %v\n", err)
	}
}

func usage() {
	fmt.Printf("%s converts to and from XLIFF files\n\n", path.Base(os.Args[0]))
	fmt.Printf("Usage: %s [-ho] <converter> [cflags]\n\n", os.Args[0])

	var converters = []string{}
	for c, _ := range registeredConverters {
		converters = append(converters, c)
	}
	sort.Strings(converters)

	fmt.Println("Available converters:\n")
	for _, c := range converters {
		fmt.Println("", c, "\t-", registeredConverters[c].Description())
	}
	fmt.Println()
	fmt.Println("Use <converter> -h to get the flags specific for the relevant converter")
	fmt.Println()

	flag.PrintDefaults()
}
