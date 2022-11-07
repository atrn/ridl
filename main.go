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
			outputSpec := filepath.Join(path, *outputFilename)
			err = ridlDir(path, outputSpec, templateNames.Slice())
		} else {
			outputSpec := filepath.Join(filepath.Dir(path), *outputFilename)
			err = ridlFile(path, outputSpec, templateNames.Slice())
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getcwd() string {
	path, err := os.Getwd()
	if err != nil { // unexpected but possible
		log.Fatal(err)
	}
	return path
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
