// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// FindTemplate searches for a template file with the given name
// and returns its full path.
//
func FindTemplate(filename string) string {
	fileexists := func(path string) (string, bool) {
		try := func(path string) (string, bool) {
			if info, err := os.Stat(path); err == nil {
				return path, !info.IsDir()
			}
			return "", false
		}
		p, ok := try(path)
		if !ok {
			p, ok = try(path + ".template")
		}
		return p, ok
	}
	if path, exists := fileexists(filename); exists {
		return path
	}
	if dir := os.Getenv("RIDLDIR"); dir != "" {
		name := filepath.Join(dir, filename)
		if path, exists := fileexists(name); exists {
			return path
		}
	}
	if *templatesDir != "" {
		name := filepath.Join(*templatesDir, filename)
		if path, exists := fileexists(name); exists {
			return path
		}
		name += ".template"
		if path, exists := fileexists(name); exists {
			return path
		}
	}
	return ""
}

// ExpandTemplate executes the template in the given file, using
// the supplied context and writing output to the given io.Writer.
//
func ExpandTemplate(filename string, context *Context, w io.Writer) error {
	if tfilename := FindTemplate(filename); tfilename != "" {
		filename = tfilename
	} else {
		return fmt.Errorf("%s: %s", filename, os.ErrNotExist)
	}
	t := template.New(filepath.Base(filename)).Funcs(cppTemplateFuncs)
	if t, err := ParseTemplates(t, filename); err != nil {
		return err
	} else if err = t.Execute(w, context); err != nil {
		return err
	} else {
		return nil
	}
}

// ParseTemplates parses the template file and adds it to
// the given template. Special comment lines in the file
// are skipped.
//
func ParseTemplates(t *template.Template, filename string) (*template.Template, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	in := bufio.NewScanner(file)
	skip := func() bool {
		return in.Scan() && strings.HasPrefix(in.Text(), "//")
	}
	if skip() {
		for skip() {
		}
		data := []byte(in.Text())
		for in.Scan() {
			data = append(data, []byte(in.Text())...)
			data = append(data, '\n')
		}
		// log.Print(string(data))
		return t.Parse(string(data))
	}
	return t.ParseFiles(filename)
}

// GetEmbeddedOutputFilename looks for a special comment that defines
// a (ridl) template's suggested output filename.
//
func GetEmbeddedOutputFilename(context *Context, filename string) (string, error) {
	comments, err := ParseComments(filename)
	if err != nil {
		return "", err
	}
	hadOutputSpec := false
	fname := ""
	for _, comment := range comments {
		fields := strings.Fields(comment)
		switch fields[0] {
		case "output":
			if hadOutputSpec {
				return fname, fmt.Errorf("%s: multiple output specifications", filename)
			}
			hadOutputSpec = true
			text := strings.TrimSpace(comment[len(fields[0]):])
			fname = strings.Replace(text, "{{.Filename}}", strings.TrimSuffix(context.Filename, filepath.Ext(context.Filename)), -1)
		}
	}
	return fname, nil
}

// ParseComments return all of the special ridl comments in a file.
//
func ParseComments(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	for r := bufio.NewScanner(file); r.Scan(); {
		const prefix = "// ridl:"
		if strings.HasPrefix(r.Text(), prefix) {
			lines = append(lines, strings.TrimSpace(r.Text()[len(prefix):]))
		}
	}
	return lines, nil
}
