package main

import (
	"fmt"
	"go/token"
	"go/types"
	"sort"
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
		switch obj.(type) {
		case *types.Const:
			p.Declare(NewConstDecl(p, obj))
		case *types.TypeName:
			switch t := obj.Type().Underlying().(type) {
			case *types.Array:
				p.Array(obj, t)
			case *types.Basic:
				p.Declare(NewTypedefDecl(p, obj, t))
			case *types.Interface:
				p.Interface(obj, t)
			case *types.Struct:
				p.Struct(obj, t)
			case *types.Slice:
				p.Slice(obj, t)
			case *types.Map:
				p.Map(obj, t)
			}
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

// Map adds a map declaration to the reciever.
func (p *Package) Map(obj types.Object, mapType *types.Map) {
	keyType := mapType.Key()
	valType := mapType.Elem()
	p.Declare(NewMapDecl(p, obj, keyType.(types.Object), valType.(types.Object)))
}

// Slice adds a slice declaration to the receiver.
func (p *Package) Slice(obj types.Object, sliceType *types.Slice) {
	elType := sliceType.Elem()
	p.Declare(NewArrayDecl(p, obj, elType, sizer.Sizeof(sliceType.Elem())))
}

// Array adds an array declaration to the receiver.
func (p *Package) Array(obj types.Object, arrayType *types.Array) {
	elType := arrayType.Elem().Underlying()
	size := sizer.Sizeof(arrayType)
	p.Declare(NewArrayDecl(p, obj, elType, size))
}

// Struct adds a struct declaration to the receiver.
func (p *Package) Struct(obj types.Object, structType *types.Struct) {
	decl := NewStructDecl(p, obj, sizer.Sizeof(structType))
	fields := make([]*types.Var, structType.NumFields())
	for i := 0; i < structType.NumFields(); i++ {
		fields[i] = structType.Field(i)
	}
	offsets := sizer.Offsetsof(fields)
	for i := 0; i < structType.NumFields(); i++ {
		field := fields[i]
		if field.Anonymous() {
			// TODO: allow embedding
			continue
		}
		fieldType := field.Type()
		f := NewStructField(p, field, sizer.Sizeof(fieldType), offsets[i], sizer.Alignof(fieldType)) // XXX check pos
		decl.AddField(f)
	}
	p.Declare(decl)
}

func makeMethodArgs(pkg *Package, obj types.Object, args *types.Tuple, prefix string) []*MethodArg {
	ma := make([]*MethodArg, 0)
	for i := 0; i < args.Len(); i++ {
		arg := args.At(i)
		name := arg.Name()
		if name == "" {
			name = fmt.Sprintf("%s%d", prefix, i+1)
		}
		ma = append(ma, NewMethodArg(pkg, arg, name))
	}
	return ma
}

// Interface adds an interface declaration to the receiver.
func (p *Package) Interface(obj types.Object, interfaceType *types.Interface) {
	xi := NewInterfaceDecl(p, obj)
	for i := 0; i < interfaceType.NumMethods(); i++ {
		fn := interfaceType.Method(i)
		sig := fn.Type().(*types.Signature)
		args := makeMethodArgs(p, fn, sig.Params(), "arg")
		results := makeMethodArgs(p, fn, sig.Results(), "res")
		xm := NewMethod(p, fn, args, results)
		xi.Declare(xm)
	}
	p.Declare(xi)
}

// Position returns the token.Position given a declaration's token.Pos
func (p *Package) Position(pos token.Pos) token.Position {
	return p.fset.Position(pos)
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
