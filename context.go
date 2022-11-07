// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"fmt"
	"os"
	"os/user"
	"time"
)

// The Context type defines the interface from RIDL templates to the
// interface being processed. A Context embeds a Package (see model.go)
// and acts as a Package, i.e. it has a name, contains collections of
// declarations of different types and defines a collection of imported
// package names.
//
// Context adds meta-data and indices to simplify template organization.
//
type Context struct {
	// Pointer to our Package
	*Package
	// The name of the file being processed.
	Filename string
	// The time at which processing is occurring.
	BuildTime time.Time
	// The name of the user doing the processing.
	Username string
	// The name of the host where processing is running.
	Hostname string
	// All decls in order
	Decls []Decl
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
//
func NewContext(filename string, pkg *Package) *Context {
	username := "unknown"
	hostname := "localhost"
	if u, err := user.Current(); err == nil {
		username = u.Username
	}
	if name, err := os.Hostname(); err == nil {
		hostname = name
	}
	context := &Context{
		Package:     pkg,
		Filename:    filename,
		BuildTime:   time.Now(),
		Username:    username,
		Hostname:    hostname,
		Decls:       pkg.Decls,
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

func isInteger(typ string) bool {
	integralTypes := map[string]struct{}{
		"byte":   struct{}{},
		"uint8":  struct{}{},
		"uint16": struct{}{},
		"uint32": struct{}{},
		"uint64": struct{}{},
		"int8":   struct{}{},
		"int16":  struct{}{},
		"int32":  struct{}{},
		"int64":  struct{}{},
		"int":    struct{}{},
	}
	_, isInt := integralTypes[typ]
	return isInt
}

func (c *Context) findEnums() {
	typedefs := make(map[string]*TypedefDecl, len(c.Typedefs))
	for _, t := range c.Typedefs {
		if isInteger(t.Alias) {
			typedefs[t.Name()] = t
		}
	}
	m := make(map[*TypedefDecl][]*ConstDecl)
	for _, constant := range c.Constants {
		t, found := typedefs[constant.Type()]
		if found {
			m[t] = append(m[t], constant)
		} else {
			c.NotEnums = append(c.NotEnums, constant)
		}
	}
	for typedef, constants := range m {
		c.Enums = append(c.Enums, &Enum{typedef, constants})
	}
}
