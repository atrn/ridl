// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"fmt"
	"go/token"
	"go/types"
)

type DeclKind int

const (
	DeclKindConst = iota
	DeclKindTypedef
	DeclKindArray
	DeclKindStruct
	DeclKindStructField
	DeclKindMap
	DeclKindInterface
	DeclKindMethod
	DeclKindMethodArg
)

func (d DeclKind) String() string {
	switch d {
	case DeclKindConst:
		return "DeclKindConst"
	case DeclKindTypedef:
		return "DeclKindTypedef"
	case DeclKindArray:
		return "DeclKindArray"
	case DeclKindStruct:
		return "DeclKindStruct"
	case DeclKindStructField:
		return "DeclKindStructField"
	case DeclKindMap:
		return "DeclKindMap"
	case DeclKindInterface:
		return "DeclKindInterface"
	case DeclKindMethod:
		return "DeclKindMethod"
	case DeclKindMethodArg:
		return "DeclKindMethodArg"
	}
	panic(fmt.Errorf("bad DeclKind value == %d", int(d)))
}

//  ================================================================

// The Decl interface is used to retrieve information common
// to all declarations, their name and type.
type Decl interface {
	// The Name method returns the name of the receiver, an
	// unqualified Go identifer.
	//
	Name() string

	// The Typename method returns the receiver's Go type as a string.
	//
	Typename() string

	// The kind of declaration.
	//
	Kind() DeclKind

	// The Position method returns the location of the declaration.
	//
	Position() token.Position
}

type SizedDecl interface {
	Decl
	Size() int64
}

//  ================================================================

// The decl struct holds data common to all declarations.  The
// decl type is embedded in declaration types.
type decl struct {
	pkg  *Package
	obj  types.Object
	kind DeclKind
}

// Name returns the receiver's name, a Go identifier.
func (d *decl) Name() string {
	return d.obj.Id()
}

func (d *decl) Typename() string {
	return TrimUntyped(d.obj.Type().String())
}

// Kind returns the kind of declaration.
func (d *decl) Kind() DeclKind {
	return d.kind
}

// Position returns the declaration's token.Position.
func (d *decl) Position() token.Position {
	return d.pkg.Position(d.obj.Pos())
}

// IsConst returns true if the declaration is a const
func (d *decl) IsConst() bool {
	return d.kind == DeclKindConst
}

// IsTypedef returns true if the declaration is a typedef
func (d *decl) IsTypedef() bool {
	return d.kind == DeclKindTypedef
}

// IsArray returns true if the declaration is an array
func (d *decl) IsArray() bool {
	return d.kind == DeclKindArray
}

// IsStruct returns true if the declaration is a struct
func (d *decl) IsStruct() bool {
	return d.kind == DeclKindStruct
}

// IsStructField returns true if the declaration is a struct field
func (d *decl) IsStructField() bool {
	return d.kind == DeclKindStructField
}

// IsMap returns true if the declaration is a map
func (d *decl) IsMap() bool {
	return d.kind == DeclKindMap
}

// IsInterface returns true if the declaration is an interface
func (d *decl) IsInterface() bool {
	return d.kind == DeclKindInterface
}

// IsMethod returns true if the declaration is a method
func (d *decl) IsMethod() bool {
	return d.kind == DeclKindMethod
}

// IsMethodArg returns true if the declaration is a method argument
func (d *decl) IsMethodArg() bool {
	return d.kind == DeclKindMethodArg
}

//  ================================================================

// The sizedDecl struct is a decl that has a size, in bytes.
// It implements the SizedDecl interface and is embedded in
// declaration types that have a size.
type sizedDecl struct {
	decl
	size int64
}

func (d *sizedDecl) Size() int64 {
	return d.size
}

//  ================================================================

// The ConstDecl type represents a constant.
type ConstDecl struct {
	decl
	IsEnumerator bool
}

// NewConstDecl returns a new ConstDecl with the given name, type and value.
func NewConstDecl(pkg *Package, obj types.Object) *ConstDecl {
	return &ConstDecl{decl{pkg, obj, DeclKindConst}, false}
}

func (decl *ConstDecl) Value() string {
	return decl.obj.(*types.Const).Val().String()
}

func (decl *ConstDecl) ExactValue() string {
	return decl.obj.(*types.Const).Val().ExactString()
}

//  ================================================================

// A TypedefDecl records a type alias formed by a type declaration
// of the form "type <identifier> <identifier>".
//
// If IsEnum is true the type is used as a Go-style enum and
// appears in the Enums slice.
//
type TypedefDecl struct {
	decl
	typedef *types.Basic
	IsEnum  bool
}

// NewTypedefDecl returns a new TypdefDecl with the given
// name and aliased type.
func NewTypedefDecl(pkg *Package, obj types.Object, typedef *types.Basic) *TypedefDecl {
	return &TypedefDecl{decl{pkg, obj, DeclKindTypedef}, typedef, false}
}

// Type returns the receiver's type, the alias part of
// the type declaration.
func (decl *TypedefDecl) Typename() string {
	return decl.typedef.String()
}

//  ================================================================

// ArrayDecl rerpresents an array or slice declaration (a slice
// being interpreted as an unbounded array).
type ArrayDecl struct {
	sizedDecl
	elType types.Type
}

// NewArrayDecl returns a new ArrayDecl with the supplied name,
// element type and size. A size of 0 implies an unbounded
// array, or vector, type.
func NewArrayDecl(pkg *Package, obj types.Object, elType types.Type, size int64) *ArrayDecl {
	return &ArrayDecl{sizedDecl{decl{pkg, obj, DeclKindArray}, size}, elType}
}

func (a *ArrayDecl) asArray() *types.Array {
	return a.obj.Type().Underlying().(*types.Array)
}

// Length returns the number of elements in the receiver.
func (a *ArrayDecl) Length() int {
	switch t := a.obj.Type().Underlying().(type) {
	case *types.Array:
		return int(t.Len())
	case *types.Slice:
		return 0
	}
	panic(fmt.Errorf("unexpected underlying type %T", a.obj.Type().Underlying()))
}

// ElTypename returns the type of the elements.
func (a *ArrayDecl) ElTypename() string {
	return a.elType.String()
}

// Typename returns the name of the receiver's type.
func (a *ArrayDecl) Typename() string {
	if a.IsVariableLength() {
		return fmt.Sprintf("[]%s", a.ElTypename())
	}
	return fmt.Sprintf("[%d]%s", a.Length(), a.ElTypename())
}

func (a *ArrayDecl) IsVariableLength() bool {
	return a.Length() == 0
}

func (a *ArrayDecl) IsFixedLength() bool {
	return a.Length() != 0
}

// ================================================================

// The StructDecl type represents a struct type declaration. A struct
// type has a name and zero or more fields, represented by StructField
// values.
type StructDecl struct {
	sizedDecl
	Fields []*StructField
}

// NewStructDecl returns a new, empty, StructDecl with the
// given name.
func NewStructDecl(pkg *Package, obj types.Object, size int64) *StructDecl {
	return &StructDecl{sizedDecl{decl{pkg, obj, DeclKindStruct}, size}, nil}
}

// AddField appends a field to the receiver's collection of fields.
func (decl *StructDecl) AddField(f *StructField) {
	decl.Fields = append(decl.Fields, f)
}

// Type returns the receiver's type.
func (decl *StructDecl) Typename() string {
	return decl.Name()
}

//  ================================================================

// The StructField type represents a field within a structure.  Each
// field has a name and a type. Embedded types are represented by
// fields with an empty Name.
type StructField struct {
	sizedDecl
	Tag       []StructFieldTag
	offset    int64
	alignment int64
}

// NewStructField returns a new StructField with the given name and
// type.
func NewStructField(pkg *Package, obj types.Object, size, offset, alignment int64) *StructField {
	return &StructField{sizedDecl{decl{pkg, obj, DeclKindStructField}, size}, nil, offset, alignment}
}

func (sf *StructField) Name() string {
	return sf.obj.Name()
}

func (sf *StructField) Offset() int64 {
	return sf.offset
}

func (sf *StructField) Alignment() int64 {
	return sf.alignment
}

//  ================================================================

// The StructFieldTag type represents a single tag applied to
// the field of a struct.
type StructFieldTag struct {
	Key   string
	Value string
}

//  ================================================================

// MapDecl represents a map declaration.
type MapDecl struct {
	decl
	keyType types.Object
	valType types.Object
}

// NewMapDecl returns a new MapDecl with the given name,
// and key and value types.
func NewMapDecl(pkg *Package, obj, keyType, valType types.Object) *MapDecl {
	return &MapDecl{decl{pkg, obj, DeclKindMap}, keyType, valType}
}

// KeyType returns the type, as a string, of the receiver's keys.
func (decl *MapDecl) KeyTypename() string {
	return decl.keyType.Name()
}

// Type returns the type of the receiver's values.
func (decl *MapDecl) Typename() string {
	return decl.valType.Name()
}

//  ================================================================

// The InterfaceDecl type represents an interface type. An interface
// is a, named, collection of zero or more Methods.
type InterfaceDecl struct {
	decl
	Methods  []*MethodDecl
	embedded []*InterfaceDecl
}

// NewInterfaceDecl returns a new, empty, InterfaceDecl with the
// given name.
func NewInterfaceDecl(pkg *Package, obj types.Object) *InterfaceDecl {
	return &InterfaceDecl{decl{pkg, obj, DeclKindInterface}, nil, nil}
}

// Type returns the receiver's type.
func (decl *InterfaceDecl) Typename() string {
	return "interface " + decl.obj.Name()
}

// Declare appends a method declaration to the interface.
func (decl *InterfaceDecl) Declare(method *MethodDecl) {
	decl.Methods = append(decl.Methods, method)
}

// Embed appends an embedded interface to the interface.
func (decl *InterfaceDecl) Embed(intf *InterfaceDecl) {
	decl.embedded = append(decl.embedded, intf)
}

//  ================================================================

// The MethodDecl type represents a method declared within an interface.
// A method has a name, zero or more arguments and zero or more results.
// Both arguments and results are represented by MethodArg values.
type MethodDecl struct {
	decl
	Args    []*MethodArg
	Results []*MethodArg
}

// NewMethod returns a new Method with the given name, arguments
// and results.
func NewMethod(pkg *Package, obj types.Object, args []*MethodArg, results []*MethodArg) *MethodDecl {
	return &MethodDecl{decl{pkg, obj, DeclKindMethod}, args, results}
}

// Type returns the receiver's type.
func (decl *MethodDecl) Typename() string {
	return decl.Name()
}

//  ================================================================

// The MethodArg  type represents an  argument to  or a result  from a
// Method. A MethodArg has a name and a type. The name may be empty.
type MethodArg struct {
	decl
	name string
}

// NewMethodArg retusn a new MethodArg with the given name and type.
func NewMethodArg(pkg *Package, obj types.Object, name string) *MethodArg {
	return &MethodArg{decl{pkg, obj, DeclKindMethodArg}, name}
}

// Type returns the receiver's type.
func (decl *MethodArg) Typename() string {
	return decl.obj.Type().String()
}

func (decl *MethodArg) Name() string {
	return decl.name
}

//  ================================================================

// Enum represents a C/C++ enumerated type that has been emulated
// using the Go idiom of defining a type and a series of constants of
// that type.
type Enum struct {
	Type        *TypedefDecl
	Enumerators []*ConstDecl
}
