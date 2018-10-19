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
	log.SetFlags(0)
	log.SetPrefix(myname + ": ")

	templateNames := NewStringSlice()

	outputFilename := flag.String("o", "", "write output to `filename`")
	flag.Var(templateNames, "t", "generate output using `template`")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", myname, "[options] [path]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		path, err := os.Getwd()
		if err == nil {
			err = ridl(path, *outputFilename, templateNames)
		}
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for _, path := range flag.Args() {
			err := ridl(path, filepath.Join(path, *outputFilename), templateNames)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
