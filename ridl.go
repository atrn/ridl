// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

func ridlDir(directoryPath string, templateNames []string) error {
	logdebug("Parsing all .ridl files from directory %q", directoryPath)
	if filenames, _ := filepath.Glob(filepath.Join(directoryPath, "*.ridl")); filenames == nil {
		return fmt.Errorf("%s: No .ridl files found in directory", directoryPath)
	} else {
		return ridlFiles(directoryPath, filenames, templateNames)
	}
}

func ridlFile(filename string, templateNames []string) error {
	logdebug("Parsing file %q", filename)
	if absPath, err := filepath.Abs(filename); err != nil {
		return err
	} else {
		return ridlFiles(filepath.Dir(absPath), []string{filename}, templateNames)
	}
}

func ridlFiles(directory string, filenames []string, templateNames []string) error {
	if pkg, err := parseFiles(filenames); err == nil {
		return generateOutput(pkg, directory, filenames, templateNames)
	} else {
		return err
	}
}

func parseFiles(filenames []string) (*Package, error) {
	fset := token.NewFileSet()
	files := make([]*ast.File, 0, len(filenames))
	for _, filename := range filenames {
		logdebug("parsing %q", filename)
		file, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", filename, err)
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
		return nil, fmt.Errorf("type check %q: %w", filenames, err)
	}
	return NewPackage(pkg), nil
}

func generateOutput(pkg *Package, directory string, filenames []string, templateNames []string) error {
	templateContext := NewContext(directory, filenames, pkg)
	for _, templateName := range templateNames {
		templateFilename := FindTemplate(templateName)
		if templateFilename == "" {
			return fmt.Errorf("%q: template file not found", templateName)
		}
		outputFilename, err := makeOutputFilename(templateFilename, templateName, directory, pkg.PackageName)
		if err != nil {
			return err
		}
		var w io.WriteCloser
		if *dryRunFlag {
			w = NopWriteCloser(io.Discard)
		} else if outputFilename == StdoutFilename {
			w = NopWriteCloser(os.Stdout)
		} else if outputFilename == "" {
			w = NopWriteCloser(io.Discard)
		} else if w, err = os.Create(outputFilename); err != nil {
			return fmt.Errorf("%q: %w", outputFilename, err)
		}
		err1 := ExpandTemplate(templateFilename, templateName, templateContext, w)
		err2 := w.Close()
		if err1 == nil {
			err1 = err2
		}
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func makeOutputFilename(templateFilename, templateName, directory, pkgname string) (string, error) {
	if *outputFilename != "" {
		return *outputFilename, nil
	}
	spec, err := GetEmbeddedOutputFilename(templateFilename)
	if err != nil {
		return "", err
	}
	t, err := template.New("output").Parse(spec)
	if err != nil {
		return "", err
	}
	type Context struct {
		Template  string
		Package   string
		Directory string
		Time      time.Time
		Username  string
		Hostname  string
	}
	context := Context{
		Template:  templateName,
		Package:   pkgname,
		Directory: directory,
		Time:      time.Now(),
		Username:  MustGetUsername(),
		Hostname:  MustGetHostname(),
	}
	var buffer bytes.Buffer
	if err = t.Execute(&buffer, context); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
