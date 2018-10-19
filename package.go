package main

import (
	"fmt"
	"go/types"
	"log"
	"strings"
)

// The Package type represents a single package, a named collection of
// declarations, constants, types, and associated imported packages.
//
// Decls and Imports are in declaration order.
//
type Package struct {
	PackageName string
	Decls       []Decl
	Imports     []string
	importIndex map[string]struct{} // aka set[string]
}

// NewPackage creates a new Package that has the given name.  The
// Package is created with a nil, as opposed to empty, Decls and
// Imports slices.
//
func NewPackage(pkg *types.Package) *Package {
	p := &Package{
		PackageName: pkg.Name(),
		importIndex: make(map[string]struct{}),
	}

	for _, imported := range pkg.Imports() {
		p.Import(imported.Path())
	}

	scope := pkg.Scope()
	names := scope.Names()
	for i := 0; i < scope.Len(); i++ {
		obj := scope.Lookup(names[i])
		switch actual := obj.(type) {
		case *types.Const:
			p.Const(actual)
		case *types.TypeName:
			p.TypeName(actual)
		default:
			log.Printf("X1:  %T  ->  %#v\n", actual, actual)
		}
	}

	return p
}

// Declare appends a Decl to the receiver's collection of declarations.
//
func (p *Package) Declare(decl Decl) {
	p.Decls = append(p.Decls, decl)
}

// Import appends the name of an imported package to the receiver's
// collection of imports.
//
func (p *Package) Import(path string) {
	if _, exists := p.importIndex[path]; !exists {
		p.Imports = append(p.Imports, path)
		p.importIndex[path] = struct{}{}
	}
}

// Const adds a declaration of a constant to the receiver.
//
func (p *Package) Const(obj *types.Const) {
	p.Declare(NewConstDecl(obj.Name(), cleanTypename(obj.Type()), obj.Val().ExactString()))
}

// TypeName adds a type declaration to the receiver.
//
func (p *Package) TypeName(obj *types.TypeName) {
	switch t := obj.Type().Underlying().(type) {
	case *types.Array:
		p.Array(obj.Name(), t)
	case *types.Basic:
		p.Typedef(obj.Name(), t)
	case *types.Interface:
		p.Interface(obj.Name(), t)
	case *types.Struct:
		p.Struct(obj.Name(), t)
	case *types.Slice:
		p.Slice(obj.Name(), t)
	case *types.Map:
		p.Map(obj.Name(), t)
	default:
		log.Printf("X2 %T %#v\n", t, t)
	}
}

// Map adds a map declaration to the reciever.
//
func (p *Package) Map(name string, obj *types.Map) {
	keytyp := obj.Key().String()
	valtyp := obj.Elem().String()
	p.Declare(NewMapDecl(name, keytyp, valtyp))
}

// Slice adds a slice declaration to the receiver.
//
func (p *Package) Slice(name string, obj *types.Slice) {
	typ := obj.Elem().String()
	p.Declare(NewArrayDecl(name, typ, 0))
}

// Array adds an array declaration to the receiver.
//
func (p *Package) Array(name string, obj *types.Array) {
	size := obj.Len()
	typ := obj.Elem().String()
	p.Declare(NewArrayDecl(name, typ, int(size)))
}

// Struct adds a struct declaration to the receiver.
//
func (p *Package) Struct(name string, obj *types.Struct) {
	astruct := NewStructDecl(name)
	for i := 0; i < obj.NumFields(); i++ {
		avar := obj.Field(i)
		if avar.Anonymous() {
			// fixme - should ridl allow embedding?
			//
			continue
		}
		f := NewStructField(avar.Name(), avar.Type().String())
		astruct.AddField(f)
	}
	p.Declare(astruct)
}

func makeMethodArgs(args *types.Tuple, prefix string) []*MethodArg {
	ma := make([]*MethodArg, 0)
	for i := 0; i < args.Len(); i++ {
		arg := args.At(i)
		name := arg.Name()
		if name == "" {
			name = fmt.Sprintf("%s%d", prefix, i+1)
		}
		ma = append(ma, NewMethodArg(name, cleanTypename(arg.Type())))
	}
	return ma
}

// Interface adds an interface declaration to the receiver.
//
func (p *Package) Interface(name string, obj *types.Interface) {
	xi := NewInterfaceDecl(name)
	for i := 0; i < obj.NumMethods(); i++ {
		fn := obj.Method(i)
		sig := fn.Type().(*types.Signature)
		args := makeMethodArgs(sig.Params(), "arg")
		results := makeMethodArgs(sig.Results(), "res")
		xm := NewMethod(fn.Name(), args, results)
		xi.Declare(xm)
	}
	p.Declare(xi)
}

// Typedef adds a type alias declaration to the receiver.
//
func (p *Package) Typedef(name string, obj *types.Basic) {
	p.Declare(NewTypedefDecl(name, obj.String()))
}

func cleanTypename(t types.Type) string {
	return strings.TrimPrefix(t.String(), "untyped ")
}
