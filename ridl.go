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
		logdebug("ParseFile %q", filename)
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
	return NewPackage(pkg, fset), nil
}

func generateOutput(pkg *Package, directory string, filenames []string, templateNames []string) error {
	templateContext := NewContext(directory, filenames, pkg)
	for _, templateName := range templateNames {
		templateFilename := FindTemplate(templateName)
		if templateFilename == "" {
			return fmt.Errorf("%q: template file not found", templateName)
		}
		outputFilename, err := getOutputFilename(templateFilename, templateName, directory, pkg.PackageName)
		if err != nil {
			return err
		}
		output, err := getOutputWriter(outputFilename)
		if err != nil {
			return err
		}
		err1 := ExpandTemplate(templateFilename, templateName, templateContext, output)
		err2 := output.Close()
		if err1 == nil {
			err1 = err2
		}
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func getOutputWriter(filename string) (io.WriteCloser, error) {
	if *dryRunFlag {
		return NopWriteCloser(io.Discard), nil
	}
	if filename == StdoutFilename || filename == "" {
		return NopWriteCloser(os.Stdout), nil
	}
	w, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", outputFilename, err)
	}
	return w, nil
}

func getOutputFilename(templateFilename, templateName, directory, pkgname string) (string, error) {
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
