package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

var (
	Version   = "0.5"
	GitHash   = ""
	BuildDate = ""
)

// print the version
func printVersion() {

	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)

	fmt.Fprintln(tw, "xliffer:\t"+Version)
	if GitHash != "" {
		fmt.Fprintln(tw, "git:\t"+GitHash)
	}
	if BuildDate != "" {
		fmt.Fprintln(tw, "build-date:\t"+BuildDate)
	}

	tw.Flush()
}
