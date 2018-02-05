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
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	myname := filepath.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(myname + ": ")

	templates := NewStrings()
	flag.Var(templates, "t", "generate output using `template`")
	outputFilename := flag.String("o", "", "output `filename`")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", myname, "[options] [path]")
		flag.PrintDefaults()
	}
	flag.Parse()

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case flag.NArg() == 0:
		break
	case flag.NArg() == 1:
		path = flag.Arg(0)
	default:
		flag.Usage()
		os.Exit(1)
	}
	pkg, err := parseFiles(path)
	if err != nil {
		log.Fatal(err)
	}
	err = generateOutput(pkg, path, templates.Slice(), *outputFilename)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFiles(path string) (*Package, error) {
	fset := token.NewFileSet()
	paths, err := filepath.Glob(filepath.Join(path, "*.ridl"))
	if err != nil {
		return nil, err
	}
	files := make([]*ast.File, 0, len(paths))
	for _, path := range paths {
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	conf := types.Config{
		IgnoreFuncBodies:         true,
		Importer:                 importer.Default(),
		DisableUnusedImportCheck: true,
	}
	pkg, err := conf.Check("", fset, files, nil)
	if err != nil {
		return nil, err
	}

	p := NewPackage(pkg.Name())
	for _, imported := range pkg.Imports() {
		p.Import(imported.Path())
	}
	scope := pkg.Scope()
	for i := 0; i < scope.Len(); i++ {
		obj := scope.Lookup(scope.Names()[i])
		switch actual := obj.(type) {
		case *types.Const:
			p.Const(actual)
		case *types.TypeName:
			p.TypeName(actual)
		default:
			log.Printf("X1:  %T  ->  %#v\n", actual, actual)
		}
	}
	return p, nil
}

func generateOutput(pkg *Package, path string, templates []string, outputFilename string) error {
	templateContext := NewContext(path, pkg)
	for _, templateName := range templates {
		out, err := getOutput(path, templateName, outputFilename)
		if err != nil {
			return err
		}
		err = ExpandTemplate(templateName, templateContext, out)
		err2 := out.Close()
		if err == nil {
			err = err2
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func getOutput(path, templateName, outputFilename string) (io.WriteCloser, error) {
	if outputFilename == "-" {
		return NopWriteCloser(os.Stdout), nil
	}
	if outputFilename == "" {
		sansExt := func(s string) string {
			return strings.TrimSuffix(s, filepath.Ext(s))
		}
		outputFilename = sansExt(path) + "." + sansExt(filepath.Base(templateName))
	}
	return os.Create(outputFilename)
}
