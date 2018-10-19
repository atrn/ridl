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
	goType := fullType
	if fullType[0] == '[' {
		end := strings.Index(fullType, "]")
		if end == -1 {
			panic("malformed type: " + fullType)
		}
		arraySuffix = fullType[0 : end+1]
		goType = fullType[end+1:]
		if asArg || arraySuffix == "[]" {
			// slice -> pointer to T
			arraySuffix = " *"
		}
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
	c := cpptype(t, false)
	if strings.HasSuffix(c, " *") {
		c = strings.TrimSuffix(c, " *")
		c = fmt.Sprintf("std::vector<%s>", c)
	}
	return c
}

func basename(path string) string {
	path = filepath.Base(path)
	path = strings.TrimSuffix(path, filepath.Ext(path))
	return path
}

func lc(s string) string {
	return strings.ToLower(s)
}

func plus(a, b int) int {
	return a + b
}

func eltype(t string) string {
	if t[0] != '[' {
		return t
	}
	end := strings.Index(t, "]")
	if end == -1 {
		panic("malformed type: " + t)
	}
	return t[end+1:]
}

func dims(t string) string {
	if t[0] != '[' {
		return ""
	}
	end := strings.Index(t, "]")
	if end == -1 {
		panic("malformed type: " + t)
	}
	return t[0 : end+1]
}

var cppTemplateFuncs = map[string]interface{}{
	"argtype":  argType,
	"basename": basename,
	"cpptype":  cppType,
	"lc":       lc,
	"plus":     plus,
	"eltype":   eltype,
	"restype":  resType,
	"dims":     dims,
}
