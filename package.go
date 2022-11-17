package main

import (
	"go/token"
	"go/types"
	"sort"
)

var Sizer = types.SizesFor("gc", "amd64")

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

	scope := pkg.Scope()
	names := scope.Names()
	objects := make([]types.Object, 0, len(names))
	for _, name := range names {
		objects = append(objects, scope.Lookup(name))
	}
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Pos() < objects[j].Pos()
	})

	for _, obj := range objects {
		switch t := obj.(type) {
		case *types.Const:
			p.Declare(NewConstDecl(p, obj))
		case *types.TypeName:
			p.Declare(makeDecl(p, t))
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

// Position returns the token.Position given a declaration's token.Pos
func (p *Package) Position(pos token.Pos) token.Position {
	return p.fset.Position(pos)
}
