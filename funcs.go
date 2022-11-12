// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"unicode"
)

func cpptype(fullType string, asArg bool) string {
	result := func(t string) string {
		logdebug("cpptype %q -> %q", fullType, t)
		return t
	}
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
	} else if strings.HasPrefix(fullType, "map[") {
		// FIXME: this is simplistic and fails for things like map[[2]int]some_type
		end := strings.Index(fullType, "]")
		if end == -1 {
			panic("malformed type: " + fullType)
		}
		gokey := fullType[4:end]
		goval := fullType[end+1:]
		ckey := cpptype(gokey, false)
		if goval == "struct{}" {
			return result(fmt.Sprintf("std::set<%s>", ckey))
		}
		cval := cpptype(goval, false)
		return result(fmt.Sprintf("std::map<%s, %s>", ckey, cval))
	}
	switch goType {
	case "byte":
		return result("std::byte" + arraySuffix)
	case "error":
		if asArg {
			return result(constref(rterror) + arraySuffix)
		}
		return result(rterror + arraySuffix)
	case "string":
		if asArg {
			return result(constref(stdstr) + arraySuffix)
		}
		return result(stdstr + arraySuffix)
	case "float32":
		return result("float" + arraySuffix)
	case "float64":
		return result("double" + arraySuffix)
	case "rune":
		return result("uint32_t" + arraySuffix)
	case "bool":
		return result(goType + arraySuffix)
	case "int":
		return result(goType + arraySuffix)
	case "uint":
		return result("unsigned int" + arraySuffix)
	case "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
		return result(goType + "_t" + arraySuffix)
	}
	if asArg {
		return result(constref(goType) + arraySuffix)
	}
	return result(goType + arraySuffix)
}

func cppType(t string) string {
	c := cpptype(t, false)
	logdebug("cpptype %q -> %q", t, c)
	return c
}

func argType(t string) string {
	c := cpptype(t, true)
	logdebug("argtype %q -> %q", t, c)
	return c
}

func resType(t string) string {
	c := cpptype(t, false)
	if strings.HasSuffix(c, " *") {
		c = strings.TrimSuffix(c, " *")
		c = fmt.Sprintf("std::vector<%s>", c)
	}
	logdebug("restype %q -> %q", t, c)
	return c
}

func basename(path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	logdebug("basename %q -> %q", path, base)
	return base
}

func tolower(s string) string {
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
	et := t[end+1:]
	logdebug("eltype %q -> %q", t, et)
	return et
}

func dims(t string) string {
	if t[0] != '[' {
		return ""
	}
	end := strings.Index(t, "]")
	if end == -1 {
		panic("malformed type: " + t)
	}
	d := t[0 : end+1]
	logdebug("dims %q -> %q", t, d)
	return d
}

func isslice(t string) bool {
	return strings.HasPrefix(t, "[]")
}

func decap(s string) string {
	if s == "" {
		return ""
	}

	r := strings.NewReader(s)
	ch, _, err := r.ReadRune()
	if err != nil {
		log.Printf("WARNING: bad Unicode string detected %q", s)
		return s
	}

	if unicode.IsUpper(ch) {
		bytes, _ := io.ReadAll(r)
		return string(unicode.ToLower(ch)) + string(bytes)
	}

	return s
}

var cppTemplateFuncs = map[string]interface{}{
	"argtype":  argType,
	"basename": basename,
	"cpptype":  cppType,
	"dims":     dims,
	"eltype":   eltype,
	"isslice":  isslice,
	"plus":     plus,
	"restype":  resType,
	"tolower":  tolower,
	"decap":    decap,
}
