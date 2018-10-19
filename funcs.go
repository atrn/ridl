// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"fmt"
	"path/filepath"
	"strings"
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
		}
		return rterror + arraySuffix
	case "string":
		if asArg {
			return constref(stdstr) + arraySuffix
		}
		return stdstr + arraySuffix
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
	}
	return goType + arraySuffix
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

var cppTemplateFuncs = map[string]interface{}{
	"cpptype":  cppType,
	"argtype":  argType,
	"restype":  resType,
	"basename": basename,
}
