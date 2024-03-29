// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"fmt"
	"go/types"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	arrayDimensionsPattern = regexp.MustCompile("^\\[(.*)\\](.*)")
	mapKeyValuePattern     = regexp.MustCompile("^map\\[(.*)\\](.*)")
)

func cpptype(fullType string, asArg bool) string {
	result := func(t string, byRef bool) string {
		if asArg && byRef {
			return fmt.Sprintf("const %s &", t)
		}
		return t
	}

	goType := fullType

	parts := arrayDimensionsPattern.FindStringSubmatch(fullType)
	if parts != nil {
		if len(parts) != 3 {
			panic(fmt.Errorf("unexpected FindStringSubmatch extracting array info: %q", parts))
		}

		dim := strings.TrimSpace(parts[1])
		goType = strings.TrimSpace(parts[2])
		ctype, _ := mapGoToCpp(goType)

		if dim == "" {
			ctype = fmt.Sprintf("std::vector<%s>", ctype)
		} else {
			ctype = fmt.Sprintf("std::array<%s, %s>", ctype, dim)
		}

		return result(ctype, true)
	}

	parts = mapKeyValuePattern.FindStringSubmatch(fullType)
	if parts != nil {
		if len(parts) != 3 {
			panic(fmt.Errorf("unexpected FindStringSubmatch extracting map info: %q", parts))
		}
		ctype := ""
		gokey := strings.TrimSpace(parts[1])
		goval := strings.TrimSpace(parts[2])
		ckey := cpptype(gokey, false)
		if goval == "struct{}" {
			ctype = fmt.Sprintf("std::set<%s>", ckey)
		} else {
			cval := cpptype(goval, false)
			ctype = fmt.Sprintf("std::map<%s, %s>", ckey, cval)
		}
		return result(ctype, true)
	}

	return result(mapGoToCpp(goType))
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

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func multiply(a, b int) int {
	return a * b
}

func divide(a, b int) int {
	return a / b
}

func sizeof(t types.Type) int {
	return int(Sizer.Sizeof(t))
}

func trimprefix(s, p string) string {
	return strings.TrimPrefix(s, p)
}

func trimsuffix(s, f string) string {
	return strings.TrimSuffix(s, f)
}

var cppTemplateFuncs = map[string]interface{}{
	"argtype":    argType,
	"basename":   basename,
	"cpptype":    cppType,
	"dims":       dims,
	"eltype":     eltype,
	"isslice":    isslice,
	"add":        add,
	"subtract":   subtract,
	"multiply":   multiply,
	"divide":     divide,
	"restype":    resType,
	"tolower":    tolower,
	"decap":      decap,
	"sizeof":     sizeof,
	"trimprefix": trimprefix,
	"trimsuffix": trimsuffix,
}
