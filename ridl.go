// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
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

func ridlDir(path string, outputSpec string, templateNames []string) error {
	if *debugFlag {
		log.Printf("Parsing .ridl files from directory %q", path)
	}
	paths, err := filepath.Glob(filepath.Join(path, "*.ridl"))
	if err != nil {
		return err
	}
	pkg, err := parseFiles(paths)
	if err == nil {
		err = generateOutput(pkg, path, templateNames, outputSpec)
	}
	return err
}

func ridlFile(path string, outputSpec string, templateNames []string) error {
	if *debugFlag {
		log.Printf("Parsing .ridl file %q", path)
	}
	pkg, err := parseFiles([]string{path})
	if err == nil {
		err = generateOutput(pkg, path, templateNames, outputSpec)
	}
	return err
}

func parseFiles(paths []string) (*Package, error) {
	fset := token.NewFileSet()
	files := make([]*ast.File, 0, len(paths))
	for _, path := range paths {
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", path, err)
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
		return nil, fmt.Errorf("type check %v: %w", paths, err)
	}

	if *debugFlag {
		log.Printf("parsed files %q ok", paths)
	}

	return NewPackage(pkg), nil
}

func generateOutput(pkg *Package, path string, templateNames []string, outputFilename string) error {
	if *debugFlag {
		log.Printf("generating output with templates %q, output spec %q", templateNames, outputFilename)
	}
	templateContext := NewContext(path, pkg)
	for _, templateName := range templateNames {
		w, err := getOutput(path, templateName, outputFilename)
		if err != nil {
			return fmt.Errorf("failed to determine output path for template %q: %w", templateName, err)
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
