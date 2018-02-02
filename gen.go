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

func cpptype(fullType string, asArg bool) string {
	constref := func(t string) string {
		return "const " + t + " &"
	}
	const (
		stdstr  = "std::string"
		rterror = "std::runtime_error"
	)
	arraySuffix := ""
	var goType string
	if fullType[0] == '[' {
		// Array or slice
		end := strings.Index(fullType, "]")
		if end == -1 {
			panic("malformed type")
		}
		arraySuffix, goType = fullType[0:end+1], fullType[end+1:]
		if asArg || arraySuffix == "[]" {
			// slice -> pointer to T
			arraySuffix = " *"
		}
	} else {
		goType = fullType
	}
	switch goType {
	case "byte":
		return "uint8_t" + arraySuffix
	case "error":
		if asArg {
			return constref(rterror) + arraySuffix
		} else {
			return rterror + arraySuffix
		}
	case "string":
		if asArg {
			return constref(stdstr) + arraySuffix
		} else {
			return stdstr + arraySuffix
		}
	case "float32":
		return "float" + arraySuffix
	case "float64":
		return "double" + arraySuffix
	case "rune":
		return "uint32_t" + arraySuffix
	case "bool":
		return goType + arraySuffix
	case "int":
		return goType + arraySuffix
	case "uint":
		return "unsigned int" + arraySuffix
	case "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
		return goType + "_t" + arraySuffix
	}
	if asArg {
		return constref(goType) + arraySuffix
	} else {
		return goType + arraySuffix
	}
}

func cppType(t string) string {
	return cpptype(t, false)
}

func argType(t string) string {
	return cpptype(t, true)
}

func resType(t string) string {
	t = cpptype(t, false)
	if strings.HasSuffix(t, " *") {
		t = strings.TrimSuffix(t, " *")
		t = fmt.Sprintf("std::vector<%s>", t)
	}
	return t
}

func basename(path string) string {
	path = filepath.Base(path)
	path = strings.TrimSuffix(path, filepath.Ext(path))
	return path
}

var templatesDir = ""

func FindTemplate(filename string) string {
	if templatesDir == "" {
		templatesDir = filepath.Clean(filepath.Join(filepath.Dir(os.Args[0]), "../lib/ridl"))
	}
	fileexists := func(path string) (string, bool) {
		info, err := os.Stat(path)
		if err == nil {
			return path, !info.IsDir()
		}
		if !os.IsNotExist(err) {
			return "", false
		}
		path = path + ".template"
		info, err = os.Stat(path)
		if err == nil {
			return path, !info.IsDir()
		}
		return "", false
	}
	if path, exists := fileexists(filename); exists {
		return path
	}
	if dir := os.Getenv("RIDL"); dir != "" {
		name := filepath.Join(dir, filename)
		if path, exists := fileexists(name); exists {
			return path
		}
	}
	if templatesDir != "" {
		name := filepath.Join(templatesDir, filename)
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

func ExpandTemplate(filename string, context *Context, w io.Writer) error {
	if tfilename := FindTemplate(filename); tfilename != "" {
		filename = tfilename
	} else {
		return fmt.Errorf("%s: %s", filename, os.ErrNotExist)
	}
	funcMap := map[string]interface{}{
		"cpptype":  cppType,
		"argtype":  argType,
		"restype":  resType,
		"basename": basename,
	}
	t := template.New(filepath.Base(filename)).Funcs(funcMap)
	if t, err := ParseFiles(t, filename); err != nil {
		return err
	} else if err = t.Execute(w, context); err != nil {
		return err
	} else {
		return nil
	}
}

func ParseFiles(t *template.Template, filename string) (*template.Template, error) {
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
