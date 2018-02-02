// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import "fmt"

// The Package type represents a single package, a named collection of
// declarations, constants and types, and imported packages.
//
// A Package's Decls and Imports are in declaration order.
//
type Package struct {
	PackageName string
	Decls       []Decl
	Imports     []string
}

// NewPackage creates a new Package that has the given name.
// The Package is created with a nil, not empty, Decls and
// Imports arrays.
//
func NewPackage(name string) *Package {
	return &Package{
		PackageName: name,
	}
}

// Declare appends a Decl to the receiver's collection of declarations.
//
func (pkg *Package) Declare(decl Decl) {
	pkg.Decls = append(pkg.Decls, decl)
}

// Import appends the name of an imported package to the receiver's
// collection of imports.
//
func (pkg *Package) Import(path string) {
	pkg.Imports = append(pkg.Imports, path)
}

//================================================================

// The Decl interface is used to access declarationss, either types
// or constants.
//
type Decl interface {
	// The Name method retrieves the name of the receiver,
	// an unqualified Go identifer.
	//
	Name() string

	// The Type method returns the reciever's, Go, type as a string.
	//
	Type() string
}

//================================================================

// The basicDecl struct holds data common to all declarations.
// The basicDecl type is intended to be embedded in more concrete
// declarations.
//
type basicDecl struct {
	name string
}

// Name returns the receiver's name, a Go identifier.
//
func (b *basicDecl) Name() string {
	return b.name
}

//================================================================

// The ConstDecl type represents a constant declration. A constant has
// a name, type and value, all represented as strings.
//
type ConstDecl struct {
	basicDecl
	typ   string
	Value string
}

// NewConstDecl returns a new ConstDecl withthe given name, type and value.
//
func NewConstDecl(name, typ, value string) *ConstDecl {
	return &ConstDecl{
		basicDecl: basicDecl{name},
		typ:       typ,
		Value:     value,
	}
}

// Type returns the receiver's type.
//
func (decl *ConstDecl) Type() string {
	return decl.typ
}

//================================================================

// A TypedefDecl records a type alias formed by a type declaration
// of the form "type <identifier> <identifier>".
//
type TypedefDecl struct {
	basicDecl
	Alias string
}

// NewTypedefDecl returns a new TypdefDecl with the given
// name and aliased type.
//
func NewTypedefDecl(name, alias string) *TypedefDecl {
	return &TypedefDecl{
		basicDecl{name},
		alias,
	}
}

// Type returns the receiver's type, the alias part of
// the type declaration.
//
func (t *TypedefDecl) Type() string {
	return t.Alias
}

//================================================================

// ArrayDecl rerpresents an array or slice declaration (a slice
// being interpreted as an unbounded array).
//
type ArrayDecl struct {
	basicDecl
	typ  string
	size int // 0 means variable, i.e. a slice
}

// NewArrayDecl returns a new ArrayDecl with the supplied name,
// element type and size. A size of 0 implies an unbounded
// array, or vector, type.
//
func NewArrayDecl(name, typ string, size int) *ArrayDecl {
	return &ArrayDecl{
		basicDecl: basicDecl{name},
		typ:       typ,
		size:      size,
	}
}

func (a *ArrayDecl) Size() int {
	return a.size
}

func (a *ArrayDecl) ElemType() string {
	return a.typ
}

// Type returns the receiver's type.
//
func (a *ArrayDecl) Type() string {
	if a.size == 0 {
		return "[]" + a.typ
	}
	return fmt.Sprintf("[%d]%s", a.size, a.typ)
}

// The StructDecl type represents a struct type declaration. A struct
// type has a name and zero or more fields, represented by StructField
// values.
//
type StructDecl struct {
	basicDecl
	Fields []*StructField
}

// NewStructDecl returns a new, empty, StructDecl with the
// given name.
//
func NewStructDecl(name string) *StructDecl {
	return &StructDecl{
		basicDecl: basicDecl{name},
		Fields:    nil,
	}
}

// AddField appends a field to the receiver's collection of fields.
//
func (decl *StructDecl) AddField(f *StructField) {
	decl.Fields = append(decl.Fields, f)
}

// Type returns the receiver's type.
//
func (decl *StructDecl) Type() string {
	return "struct " + decl.Name()
}

//================================================================

// The StructField type represents a field within a structure.  Each
// field has a name and a type. Embedded types are represented by
// fields with an empty Name.
//
type StructField struct {
	basicDecl
	typ string
	Tag []StructFieldTag
}

// NewStructField returns a new StructField with the given name and
// type.
//
func NewStructField(name, typ string) *StructField {
	return &StructField{
		basicDecl: basicDecl{name},
		typ:       typ,
	}
}

// Type returns the receiver's type.
//
func (sf *StructField) Type() string {
	return sf.typ
}

//================================================================

// The StructFieldTag type represents a single tag applied to
// the field of a struct.
//
type StructFieldTag struct {
	Key   string
	Value string
}

//================================================================

type MapDecl struct {
	basicDecl
	keytyp string
	valtyp string
}

func NewMapDecl(name, keytyp, valtyp string) *MapDecl {
	return &MapDecl{
		basicDecl: basicDecl{name},
		keytyp:       keytyp,
		valtyp: valtyp,
	}
}

func (decl *MapDecl) KeyType() string {
	return decl.keytyp
}

func (decl *MapDecl) Type() string {
	return decl.valtyp
}

//================================================================

// The InterfaceDecl type represents an interface type. An interface
// is a, named, collection of zero or more Methods.
//
type InterfaceDecl struct {
	basicDecl
	Methods []*Method
	Embeds  []string
}

// NewInterfaceDecl returns a new, empty, InterfaceDecl with the
// given name.
//
func NewInterfaceDecl(name string) *InterfaceDecl {
	return &InterfaceDecl{
		basicDecl: basicDecl{name},
		Methods:   nil,
	}
}

// Type returns the receiver's type.
//
func (decl *InterfaceDecl) Type() string {
	return "interface " + decl.Name()
}

// Declare appends a method declaration to the interface.
//
func (decl *InterfaceDecl) Declare(method *Method) {
	decl.Methods = append(decl.Methods, method)
}

// Embed appends an embedded interface to the interface.
//
func (decl *InterfaceDecl) Embed(n string) {
	decl.Embeds = append(decl.Embeds, n)
}

//================================================================

// The Method type represents a method declared within an interface.
// A method has a name and zero or more arguments, represented by
// MethodArg values, and zero of results, also represented by
// MethodArg values.
//
type Method struct {
	basicDecl
	Args    []*MethodArg
	Results []*MethodArg
}

// NewMethodName returns a new Method with the given name, arguments
// and results.
//
func NewMethod(name string, args []*MethodArg, results []*MethodArg) *Method {
	return &Method{
		basicDecl: basicDecl{name},
		Args:      args,
		Results:   results,
	}
}

// Type returns the receiver's type.
//
func (decl *Method) Type() string {
	return "func " + decl.Name()
}

//================================================================

// The MethodArg  type represents an  argument to  or a result  from a
// Method. A MethodArg has a name and a type. The name may be empty.
//
type MethodArg struct {
	basicDecl
	typ string
}

// NewMethodArg retusn a new MethodArg with the given name and type.
//
func NewMethodArg(name, typ string) *MethodArg {
	return &MethodArg{
		basicDecl{name},
		typ,
	}
}

// Type returns the receiver's type.
//
func (decl *MethodArg) Type() string {
	return decl.typ
}
