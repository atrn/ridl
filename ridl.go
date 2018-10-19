// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ridl(path string, outputSpec string, templateNames *StringSlice) error {
	if pkg, err := parseFiles(path); err != nil {
		return err
	} else {
		return generateOutput(pkg, path, templateNames, outputSpec)
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

	return NewPackage(pkg), nil
}

func generateOutput(pkg *Package, path string, templateNames *StringSlice, outputFilename string) error {
	templateContext := NewContext(path, pkg)
	for _, templateName := range templateNames.Slice() {
		w, err := getOutput(path, templateName, outputFilename)
		if err != nil {
			return err
		}
		err = ExpandTemplate(templateName, templateContext, w)
		err2 := w.Close()
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
