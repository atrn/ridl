package main

import (
	"fmt"
	"go/token"
	"go/types"
	"sort"
	"strings"
)

var sizer = types.SizesFor("gc", "amd64")

// The Package type represents a single package, a named collection of
// declarations, constants, types, and associated imported packages.
//
// Decls and Imports are in declaration order.
type Package struct {
	PackageName string
	Decls       []Decl
	Imports     []string
	importIndex map[string]struct{} // aka set[string]
	fset        *token.FileSet
}

// NewPackage creates a new Package that has the given name.  The
// Package is created with a nil, as opposed to empty, Decls and
// Imports slices.
func NewPackage(pkg *types.Package, fset *token.FileSet) *Package {
	p := &Package{
		PackageName: pkg.Name(),
		importIndex: make(map[string]struct{}),
		fset:        fset,
	}

	for _, imported := range pkg.Imports() {
		p.Import(imported.Path())
	}

	for _, obj := range objectsInDeclarationOrder(pkg) {
		switch actual := obj.(type) {
		case *types.Const:
			p.Const(actual)
		case *types.TypeName:
			p.TypeName(actual)
		}
	}

	return p
}

// Declare appends a Decl to the receiver's collection of declarations.
func (p *Package) Declare(decl Decl) {
	p.Decls = append(p.Decls, decl)
}

// Import appends the name of an imported package to the receiver's
// collection of imports.
func (p *Package) Import(path string) {
	if _, exists := p.importIndex[path]; !exists {
		p.Imports = append(p.Imports, path)
		p.importIndex[path] = struct{}{}
	}
}

// Const adds a declaration of a constant to the receiver.
func (p *Package) Const(obj *types.Const) {
	typ := cleanTypename(obj.Type())
	val := obj.Val().String()
	exact := obj.Val().ExactString()
	d := NewConstDecl(p, obj.Pos(), obj.Name(), typ, val, exact)
	p.Declare(d)
}

// TypeName adds a type declaration to the receiver.
func (p *Package) TypeName(obj *types.TypeName) {
	switch t := obj.Type().Underlying().(type) {
	case *types.Array:
		p.Array(obj.Pos(), obj.Name(), t)
	case *types.Basic:
		p.Typedef(obj.Pos(), obj.Name(), t)
	case *types.Interface:
		p.Interface(obj.Pos(), obj.Name(), t)
	case *types.Struct:
		p.Struct(obj.Pos(), obj.Name(), t)
	case *types.Slice:
		p.Slice(obj.Pos(), obj.Name(), t)
	case *types.Map:
		p.Map(obj.Pos(), obj.Name(), t)
	}
}

// Map adds a map declaration to the reciever.
func (p *Package) Map(pos token.Pos, name string, obj *types.Map) {
	keytyp := obj.Key().String()
	valtyp := obj.Elem().String()
	p.Declare(NewMapDecl(p, pos, name, keytyp, valtyp))
}

// Slice adds a slice declaration to the receiver.
func (p *Package) Slice(pos token.Pos, name string, obj *types.Slice) {
	typ := obj.Elem().String()
	p.Declare(NewArrayDecl(p, pos, name, typ, 0, sizer.Sizeof(obj.Elem())))
}

// Array adds an array declaration to the receiver.
func (p *Package) Array(pos token.Pos, name string, obj *types.Array) {
	length := obj.Len()
	typ := obj.Elem().String()
	size := sizer.Sizeof(obj)
	p.Declare(NewArrayDecl(p, pos, name, typ, int(length), size))
}

// Struct adds a struct declaration to the receiver.
func (p *Package) Struct(pos token.Pos, name string, obj *types.Struct) {
	decl := NewStructDecl(p, pos, name, sizer.Sizeof(obj))
	fields := make([]*types.Var, obj.NumFields())
	for i := 0; i < obj.NumFields(); i++ {
		fields[i] = obj.Field(i)
	}
	offsets := sizer.Offsetsof(fields)
	for i := 0; i < obj.NumFields(); i++ {
		field := fields[i]
		if field.Anonymous() {
			// TODO: allow embedding
			continue
		}
		typ := field.Type()
		f := NewStructField(p, pos, field.Name(), typ.String(), sizer.Sizeof(typ), offsets[i], sizer.Alignof(typ)) // XXX check pos
		decl.AddField(f)
	}
	p.Declare(decl)
}

func makeMethodArgs(pkg *Package, pos token.Pos, args *types.Tuple, prefix string) []*MethodArg {
	ma := make([]*MethodArg, 0)
	for i := 0; i < args.Len(); i++ {
		arg := args.At(i)
		name := arg.Name()
		if name == "" {
			name = fmt.Sprintf("%s%d", prefix, i+1)
		}
		ma = append(ma, NewMethodArg(pkg, pos, name, cleanTypename(arg.Type())))
	}
	return ma
}

// Interface adds an interface declaration to the receiver.
func (p *Package) Interface(pos token.Pos, name string, obj *types.Interface) {
	xi := NewInterfaceDecl(p, pos, name)
	for i := 0; i < obj.NumMethods(); i++ {
		fn := obj.Method(i)
		sig := fn.Type().(*types.Signature)
		args := makeMethodArgs(p, pos, sig.Params(), "arg")
		results := makeMethodArgs(p, pos, sig.Results(), "res")
		xm := NewMethod(p, pos, fn.Name(), args, results)
		xi.Declare(xm)
	}
	p.Declare(xi)
}

// Typedef adds a type alias declaration to the receiver.
func (p *Package) Typedef(pos token.Pos, name string, obj *types.Basic) {
	p.Declare(NewTypedefDecl(p, pos, name, obj.String()))
}

// SourceFile returns the name of the source file where the given declaration was defined.
func (p *Package) SourceFile(pos token.Pos) string {
	return p.fset.Position(pos).Filename
}

// SourceLine returns the line number corresponding to a declaration's token.Pos
func (p *Package) SourceLine(pos token.Pos) int {
	return p.fset.Position(pos).Line
}

// SourceColumn returns the column number corresponding to a declaration's token.Pos
func (p *Package) SourceColumn(pos token.Pos) int {
	return p.fset.Position(pos).Column
}

func cleanTypename(t types.Type) string {
	return strings.TrimPrefix(t.String(), "untyped ")
}

func objectsInDeclarationOrder(pkg *types.Package) []types.Object {
	scope := pkg.Scope()
	names := scope.Names()
	objects := make([]types.Object, 0, len(names))
	for _, name := range names {
		objects = append(objects, scope.Lookup(name))
	}
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Pos() < objects[j].Pos()
	})
	return objects
}
