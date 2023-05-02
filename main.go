// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	StdoutFilename = "-"
)

var (
	templateNames  = NewStringSlice()
	templateDirs   = NewStringSlice()
	outputFilename = flag.String("o", "", "write output to `filename` (use '-' for stdout)")
	debugFlag      = flag.Bool("debug", false, "enable debug output")
	dryRunFlag     = flag.Bool("n", false, "do not generate output, only parse files")
)

func main() {
	myname := filepath.Base(os.Args[0])

	log.SetFlags(0)
	log.SetPrefix(myname + ": ")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", myname, "[options] path")
		flag.PrintDefaults()
	}

	versionFlag := flag.Bool("version", false, "output version and exit")
	flag.Var(templateNames, "t", "generate output using `template`")
	flag.Var(templateDirs, "T", "search for templates in `dir`")
	typeMapFlag := flag.String("typemap", "", "type mapping `filename`")
	writeTypeMapFlag := flag.Bool("write-typemap", false, "output type mapping JSON and exit")

	if s := os.Getenv("RIDLPATH"); s != "" {
		*templateDirs = filepath.SplitList(s)
	}

	flag.Parse()

	if *versionFlag {
		fmt.Print(versionNumber)
		os.Exit(0)
	}

	initTypeMap()

	if *typeMapFlag != "" {
		if err := readTypeMap(*typeMapFlag); err != nil {
			log.Fatal(err)
		}
	}

	if *writeTypeMapFlag {
		writeTypeMap(os.Stdout)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for _, path := range flag.Args() {
		var err error
		if isDir(path) {
			err = ridlDir(path, templateNames.Slice())
		} else {
			err = ridlFile(path, templateNames.Slice())
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
