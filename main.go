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

var (
	templateNames  = NewStringSlice()
	outputFilename = flag.String("o", "", "write output to `filename` ('-' means stdout)")
	templatesDir   = flag.String("D", "", "search for templates in `dir`")
	debugFlag      = flag.Bool("debug", false, "enable debug output")
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

	flag.Parse()

	if *versionFlag {
		fmt.Print(versionNumber)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if *templatesDir == "" {
		if *templatesDir = os.Getenv("RIDL_TEMPLATES_DIR"); *templatesDir == "" {
			*templatesDir = filepath.Clean(filepath.Join(filepath.Dir(os.Args[0]), "../lib/ridl"))
		}
	}

	if templateNames.Len() > 0 && !isDir(*templatesDir) {
		log.Fatalf("%s: Not found or not a directory", *templatesDir)
	}

	for _, path := range flag.Args() {
		var err error
		if isDir(path) {
			err = ridlDir(path, makeOutputSpec(path), templateNames.Slice())
		} else {
			err = ridlFile(path, makeOutputSpec(filepath.Dir(path)), templateNames.Slice())
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func makeOutputSpec(dir string) string {
	if *outputFilename == "-" {
		return "-"
	}
	return filepath.Join(dir, *outputFilename)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
