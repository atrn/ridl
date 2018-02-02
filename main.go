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

var outputFilename = ""

func main() {
	myname := filepath.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(myname + ": ")
	templates := NewStrings()
	flag.Var(templates, "t", "generate output using `template`")
	flag.StringVar(&outputFilename, "o", "", "output `filename`")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", myname, "[options] [path]")
		flag.PrintDefaults()
	}
	flag.Parse()
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	switch flag.NArg() {
	case 0:
	case 1:
		path = flag.Arg(0)
	default:
		flag.Usage()
		os.Exit(1)
	}
	pkg, err := ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}
	err = GenerateOutput(pkg, path, templates.Slice())
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateOutput performs the ridl "code-generation" step on a Package.
//
func GenerateOutput(pkg *Package, path string, templates []string) error {
	templateContext := NewContext(path, pkg)
	for _, templateName := range templates {
		out := NopWriteCloser(os.Stdout)
		if outputFilename != "" {
			file, err := MakeOutputFile(path, templateName)
			if err != nil {
				return err
			}
			out = file
		}
		err := ExpandTemplate(templateName, templateContext, out)
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

// ParseFiles parses all the .ridl files in the given directory
// and returns a new Package representing their content.
//
func ParseFiles(path string) (*Package, error) {
	var err error
	fset := token.NewFileSet()
	filenames, err := filepath.Glob(filepath.Join(path, "*.ridl"))
	if err != nil {
		return nil, err
	}
	files := make([]*ast.File, 0, len(filenames))
	for _, filename := range filenames {
		var file *ast.File
		file, err = parser.ParseFile(fset, filename, nil, 0)
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

// CleanTypename removes the "untypeed" suffix appended to
// the type names of untyped numeric constants.
//
func CleanTypename(t types.Type) string {
	return strings.TrimPrefix(t.String(), "untyped ")
}

// MakeOutputFile creates a file for storing output.
//
func MakeOutputFile(path, templateName string) (io.WriteCloser, error) {
	if outputFilename != "" {
		return os.Create(outputFilename)
	}
	sansExt := func(s string) string {
		return strings.TrimSuffix(s, filepath.Ext(s))
	}
	name := sansExt(path) + "-" + sansExt(filepath.Base(templateName)) + ".out"
	return os.Create(name)
}
