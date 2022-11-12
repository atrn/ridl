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
	"regexp"
	"strings"
	"text/template"
)

// FindTemplate searches for a template file with the given name
// and returns its full path.
func FindTemplate(name string) string {
	logdebug("looking for template %q", name)
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
	if path, exists := fileexists(name); exists {
		logdebug("template %q found - %q", name, path)
		return path
	}
	for _, directory := range templateDirs.Slice() {
		filename := filepath.Join(directory, name)
		if path, exists := fileexists(filename); exists {
			logdebug("template %q found - %q", name, path)
			return path
		}
		filename += ".template"
		if path, exists := fileexists(filename); exists {
			logdebug("template %q found - %q", name, path)
			return path
		}
	}
	logdebug("template %q not found", name)
	return ""
}

// ExpandTemplate executes the template in the given file, using
// the supplied context and writing output to the given io.Writer.
func ExpandTemplate(filename string, name string, context *Context, w io.Writer) error {
	t := template.New(name).Funcs(cppTemplateFuncs)
	if t, err := ParseTemplate(t, filename); err != nil {
		return fmt.Errorf("ExpandTemplate %q: %w", filename, err)
	} else if err = t.Execute(w, context); err != nil {
		return err
	} else {
		return nil
	}
}

var ridlCommentPattern = regexp.MustCompile("^//\\s*ridl:\\s*(.*)\\s*$")

// ParseTemplate parses the template file and adds it to the given
// template. Special comment lines in the file are skipped.
func ParseTemplate(t *template.Template, filename string) (*template.Template, error) {
	logdebug("parse template %q", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ParseTemplate %q: %w", filename, err)
	}
	defer file.Close()
	in := bufio.NewScanner(file)
	var data []byte
	for in.Scan() {
		line := in.Text()
		if !ridlCommentPattern.MatchString(line) {
			data = append(data, []byte(line)...)
			data = append(data, '\n')
		}
	}
	return t.Parse(string(data))
}

// GetEmbeddedOutputFilename looks for a special comment that defines
// a ridl template's output filename spec.
func GetEmbeddedOutputFilename(templateFilename string) (string, error) {
	comments, err := ParseComments(templateFilename)
	if err != nil {
		return "", err
	}
	outputSpec := ""
	for _, comment := range comments {
		fields := strings.Fields(comment)
		switch fields[0] {
		case "output":
			if outputSpec != "" {
				return outputSpec, fmt.Errorf("%q: template file contains multiple output specifications", templateFilename)
			}
			outputSpec := strings.TrimSpace(comment[len(fields[0]):])
			logdebug("template %q defines output spec %q", templateFilename, outputSpec)
		}
	}
	return outputSpec, nil
}

// ParseComments return the text of all of the special ridl comments in a file.
//
// Ridl's "special comments" are single line comments starting at the
// beginning of a line and who's first, whitespace-separated, field is
// "ridl:". The "ridl:" prefix is followed by the comment's "text",
// one or more whitespace-separated "words".
//
// This function returns a slice of strings containing all of the
// comment texts defined by the special "// ridl:" comments.
//
//
func ParseComments(filename string) ([]string, error) {
	logdebug("ParseComments %q", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ParseComments %q: %w", filename, err)
	}
	defer file.Close()
	var texts []string
	for r := bufio.NewScanner(file); r.Scan(); {
		parts := ridlCommentPattern.FindStringSubmatch(r.Text())
		if parts != nil && parts[1] != "" {
			texts = append(texts, parts[1])
		}
	}
	return texts, nil
}
