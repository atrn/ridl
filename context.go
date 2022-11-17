// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"fmt"
	"go/constant"
	"go/types"
	"sort"
	"strings"
	"time"
)

// The Context type defines the interface from RIDL templates to the
// interface being processed. A Context embeds a Package (see model.go)
// and acts as a Package, i.e. it has a name, contains collections of
// declarations of different types and defines a collection of imported
// package names.
//
// Context adds meta-data and indices to simplify template organization.
type Context struct {
	// Pointer to our Package
	*Package
	// Ridl's version "number"
	RidlVersion string
	// The directory being processed
	Directory string
	// The names of all .ridl files used when parsing
	Filenames []string
	// The time at which processing is occurring.
	BuildTime time.Time
	// The name of the user doing the processing.
	Username string
	// The name of the host where processing is running.
	Hostname string
	// The basic alias types, "typedefs", "type <ident> <type>"...
	Typedefs []*TypedefDecl
	// The array and slice types.
	ArrayTypes []*ArrayDecl
	// The map types.
	MapTypes []*MapDecl
	// The struct types.
	StructTypes []*StructDecl
	// Interfaces.
	Interfaces []*InterfaceDecl
	// Constants.
	Constants []*ConstDecl
	// Enums.
	Enums []*Enum
	// NotEnums - constants that are not in Enums.
	NotEnums []*ConstDecl
}

// NewContext returns a new Context for the given file and Package.
func NewContext(directory string, filenames []string, pkg *Package) *Context {
	context := &Context{
		RidlVersion: strings.TrimSpace(versionNumber),
		Package:     pkg,
		Directory:   directory,
		Filenames:   filenames,
		BuildTime:   time.Now(),
		Username:    MustGetUsername(),
		Hostname:    MustGetHostname(),
		Typedefs:    make([]*TypedefDecl, 0),
		ArrayTypes:  make([]*ArrayDecl, 0),
		MapTypes:    make([]*MapDecl, 0),
		StructTypes: make([]*StructDecl, 0),
		Interfaces:  make([]*InterfaceDecl, 0),
		Constants:   make([]*ConstDecl, 0),
		Enums:       make([]*Enum, 0),
	}
	for _, decl := range pkg.Decls {
		switch d := decl.(type) {
		case *ConstDecl:
			context.Constants = append(context.Constants, d)
		case *TypedefDecl:
			context.Typedefs = append(context.Typedefs, d)
		case *ArrayDecl:
			context.ArrayTypes = append(context.ArrayTypes, d)
		case *MapDecl:
			context.MapTypes = append(context.MapTypes, d)
		case *StructDecl:
			context.StructTypes = append(context.StructTypes, d)
		case *InterfaceDecl:
			context.Interfaces = append(context.Interfaces, d)
		default:
			panic(fmt.Sprintf("unexpected Decl type: %T", d))
		}
	}
	context.findEnums()
	return context
}

func (c *Context) findEnums() {
	typedefs := make(map[string]*TypedefDecl, len(c.Typedefs))
	for _, t := range c.Typedefs {
		if (t.typedef.Info() & types.IsInteger) != 0 {
			typedefs[t.Name()] = t
		}
	}
	m := make(map[*TypedefDecl][]*ConstDecl)
	for _, constant := range c.Constants {
		t, found := typedefs[constant.TypeName()]
		if found {
			m[t] = append(m[t], constant)
			t.IsEnum = true
			constant.IsEnumerator = true
		} else {
			c.NotEnums = append(c.NotEnums, constant)
		}
	}
	for typedef, constants := range m {
		e := &Enum{typedef, constants, enumIsDense(constants)}
		c.Enums = append(c.Enums, e)
	}
}

func enumIsDense(enumerators []*ConstDecl) bool {
	sz := len(enumerators)
	if sz < 2 {
		return true
	}
	values := make([]int, sz)
	for i := range enumerators {
		n, _ := constant.Int64Val(enumerators[i].Value())
		values[i] = int(n)
	}
	sort.Ints(values)
	if values[0] != 0 {
		return false
	}
	for i := 1; i < len(values); i++ {
		if values[i]-values[i-1] != 1 {
			return false
		}
	}
	return true
}
