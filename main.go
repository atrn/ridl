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

func main() {
	myname := filepath.Base(os.Args[0])

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", myname, "[options] [path]")
		flag.PrintDefaults()
	}

	log.SetFlags(0)
	log.SetPrefix(myname + ": ")

	templateNames := NewStringSlice()

	outputFilename := flag.String("o", "", "write output to `filename`")
	flag.Var(templateNames, "t", "generate output using `template`")

	flag.Parse()

	process := func(path, output string) {
		err := ridl(path, output, templateNames.Slice())
		if err != nil {
			log.Fatal(err)
		}
	}

	if flag.NArg() == 0 {
		process(getcwd(), *outputFilename)
	} else {
		for _, path := range flag.Args() {
			process(path, filepath.Join(path, *outputFilename))
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
