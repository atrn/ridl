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
	"go/token"
	"go/types"
	"sort"
	"strings"
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
		return "const"
	case DeclKindTypedef:
		return "type"
	case DeclKindArray:
		return "array"
	case DeclKindStruct:
		return "struct"
	case DeclKindStructField:
		return "field"
	case DeclKindMap:
		return "map"
	case DeclKindInterface:
		return "interface"
	case DeclKindMethod:
		return "method"
	case DeclKindMethodArg:
		return "argument"
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

	// The TypeName method returns the receiver's Go type as a string.
	//
	TypeName() string

	// The kind of declaration.
	//
	Kind() DeclKind

	// The Position method returns the location of the declaration.
	//
	Position() token.Position
}

//  ================================================================

// The decl struct holds data common to all declarations.  The
// decl type is embedded in declaration types.
type decl struct {
	pkg    *Package
	Object types.Object
	kind   DeclKind
}

// Name returns the receiver's name, a Go identifier.
func (d *decl) Name() string {
	return d.Object.Id()
}

func (d *decl) TypeName() string {
	return TrimUntyped(d.Object.Type().String())
}

func (d *decl) IsUntyped() bool {
	return strings.HasPrefix(d.Object.Type().String(), "untyped ")
}

// Kind returns the kind of declaration.
func (d *decl) Kind() DeclKind {
	return d.kind
}

// Position returns the declaration's token.Position.
func (d *decl) Position() token.Position {
	return d.pkg.Position(d.Object.Pos())
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

// The ConstDecl type represents a constant.
type ConstDecl struct {
	decl
	IsEnumerator bool
	EnumType     Decl
}

// NewConstDecl returns a new ConstDecl with the given name, type and value.
func NewConstDecl(pkg *Package, obj types.Object) *ConstDecl {
	return &ConstDecl{decl{pkg, obj, DeclKindConst}, false, nil}
}

func (decl *ConstDecl) Value() constant.Value {
	return decl.Object.(*types.Const).Val()
}

func (decl *ConstDecl) ExactValue() string {
	return decl.Object.(*types.Const).Val().ExactString()
}

//  ================================================================

// A TypedefDecl records a type definition formed by a type declaration
// of the form "type <identifier> <identifier>".
//
// If IsEnum is true the type is used as a Go-style enum and
// appears in the Enums slice.
type TypedefDecl struct {
	decl
	typedef   *types.Basic
	IsEnum    bool
	IsPointer bool
}

// NewTypedefDecl returns a new TypdefDecl with the given name and aliased
// type.
func NewTypedefDecl(pkg *Package, obj types.Object, typedef *types.Basic) *TypedefDecl {
	return &TypedefDecl{decl{pkg, obj, DeclKindTypedef}, typedef, false, false}
}

// Helper method to make newly created TypedefDecl values as pointers.
func (decl *TypedefDecl) MarkedAsPointer() *TypedefDecl {
	decl.IsPointer = true
	return decl
}

// Type returns the receiver's type, the alias part of
// the type declaration.
func (decl *TypedefDecl) TypeName() string {
	prefix := ""
	if decl.IsPointer {
		prefix = "*"
	}
	return prefix + decl.typedef.String()
}

//  ================================================================

// ArrayDecl rerpresents an array or slice declaration (a slice
// being interpreted as an unbounded array).
type ArrayDecl struct {
	decl
	elTypename string
	elType     types.Type
}

// NewArrayDecl returns a new ArrayDecl with the supplied name,
// element type and size. A size of 0 implies an unbounded
// array, or vector, type.
func NewArrayDecl(pkg *Package, obj types.Object, typename string, elType types.Type) *ArrayDecl {
	return &ArrayDecl{decl{pkg, obj, DeclKindArray}, typename, elType}
}

// Length returns the number of elements in the receiver.
func (a *ArrayDecl) Length() int {
	switch t := a.Object.Type().Underlying().(type) {
	case *types.Array:
		return int(t.Len())
	case *types.Slice:
		return 0
	}
	panic(fmt.Errorf("unexpected underlying type %T", a.Object.Type().Underlying()))
}

// ElTypeName returns the type of the elements of the array.
func (a *ArrayDecl) ElTypeName() string {
	return a.elTypename
}

// TypeName returns the name of the receiver's type.
func (a *ArrayDecl) TypeName() string {
	return a.ElTypeName()
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
	decl
	Fields []*StructField
}

// NewStructDecl returns a new, empty, StructDecl with the
// given name.
func NewStructDecl(pkg *Package, obj types.Object) *StructDecl {
	return &StructDecl{decl{pkg, obj, DeclKindStruct}, nil}
}

// AddField appends a field to the receiver's collection of fields.
func (decl *StructDecl) AddField(f *StructField) {
	decl.Fields = append(decl.Fields, f)
}

// Type returns the receiver's type.
func (decl *StructDecl) TypeName() string {
	return decl.Name()
}

//  ================================================================

// The StructField type represents a field within a structure.  Each
// field has a name and a type. Embedded types are represented by
// fields with an empty Name.
type StructField struct {
	decl
	Tags      []Tag
	offset    int
	alignment int
}

// NewStructField returns a new StructField
func NewStructField(pkg *Package, obj types.Object, offset, alignment int64) *StructField {
	return &StructField{decl{pkg, obj, DeclKindStructField}, nil, int(offset), int(alignment)}
}

func (sf *StructField) Name() string {
	return sf.Object.Name()
}

func (sf *StructField) Type() types.Type {
	return sf.Object.Type()
}

func (sf *StructField) Offset() int {
	return sf.offset
}

func (sf *StructField) Alignment() int {
	return sf.alignment
}

func (sf *StructField) HasTag(key string) bool {
	for _, tag := range sf.Tags {
		if tag.Key == key {
			return true
		}
	}
	return false
}

func (sf *StructField) TagValue(key string) string {
	for _, tag := range sf.Tags {
		if tag.Key == key {
			return tag.Value
		}
	}
	return ""
}

//  ================================================================

// Tag represents a single struct field tag, a key/value pair of strings.
type Tag struct {
	Key   string
	Value string
}

//  ================================================================

// MapDecl represents a map declaration.
type MapDecl struct {
	decl
	keyType types.Type
	valType types.Type
}

// NewMapDecl returns a new MapDecl with the given name,
// and key and value types.
func NewMapDecl(pkg *Package, obj types.Object, keyType, valType types.Type) *MapDecl {
	return &MapDecl{decl{pkg, obj, DeclKindMap}, keyType, valType}
}

func (decl *MapDecl) asMap() *types.Map {
	return decl.Object.Type().Underlying().(*types.Map)
}

// Type returns the type of the receiver's values.
func (decl *MapDecl) TypeName() string {
	return decl.asMap().String()
}

func (decl *MapDecl) Key() types.Type {
	return decl.keyType
}

func (decl *MapDecl) Value() types.Type {
	return decl.valType
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
func (decl *InterfaceDecl) TypeName() string {
	return "interface " + decl.Object.Name()
}

// Declare appends a method declaration to the interface.
func (decl *InterfaceDecl) addMethod(method *MethodDecl) {
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
func (decl *MethodDecl) TypeName() string {
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
func (decl *MethodArg) TypeName() string {
	return decl.Object.Type().String()
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
	IsDense     bool
}

//  ================================================================

func makeDecl(pkg *Package, obj *types.TypeName) Decl {
	switch t := obj.Type().Underlying().(type) {
	case *types.Array:
		return NewArrayDecl(pkg, obj, getTypeName(t.Elem()), t.Elem().Underlying())
	case *types.Basic:
		return NewTypedefDecl(pkg, obj, t)
	case *types.Interface:
		return makeInterface(pkg, obj, t)
	case *types.Struct:
		return makeStruct(pkg, obj, t)
	case *types.Slice:
		return NewArrayDecl(pkg, obj, getTypeName(t.Elem()), t.Elem().Underlying())
	case *types.Map:
		return NewMapDecl(pkg, obj, t.Key(), t.Elem())
	case *types.Pointer:
		return NewTypedefDecl(pkg, obj, t.Elem().(*types.Basic)).MarkedAsPointer()
	default:
		panic(fmt.Errorf("%T: not handling in makeDecl", t))
	}
}

func getTypeName(t types.Type) string {
	switch actual := t.(type) {
	case *types.Basic:
		return actual.Name()
	case *types.Named:
		return actual.Obj().Name()
	case *types.Pointer:
		return actual.String()
	default:
		panic(fmt.Errorf("getTypeName: %T", t))
	}
}

func makeStruct(pkg *Package, obj types.Object, structType *types.Struct) Decl {
	decl := NewStructDecl(pkg, obj)
	fields := make([]*types.Var, structType.NumFields())
	for i := 0; i < structType.NumFields(); i++ {
		fields[i] = structType.Field(i)
	}
	offsets := Sizer.Offsetsof(fields)
	for i := 0; i < structType.NumFields(); i++ {
		field := fields[i]
		if field.Anonymous() {
			// TODO: allow embedding
			continue
		}
		fieldType := field.Type()
		f := NewStructField(pkg, field, offsets[i], Sizer.Alignof(fieldType)) // XXX check pos
		decl.AddField(f)
	}
	return decl
}

func makeInterface(pkg *Package, obj types.Object, interfaceType *types.Interface) Decl {
	intf := NewInterfaceDecl(pkg, obj)
	methodsInOrder := make([]*types.Func, interfaceType.NumMethods())
	for i := 0; i < interfaceType.NumMethods(); i++ {
		methodsInOrder[i] = interfaceType.Method(i)
	}
	sort.Slice(methodsInOrder, func(i, j int) bool {
		return methodsInOrder[i].Pos() < methodsInOrder[j].Pos()
	})
	for i := 0; i < interfaceType.NumMethods(); i++ {
		method := methodsInOrder[i]
		signature := method.Type().(*types.Signature)
		args := makeArgs(pkg, method, signature.Params(), "arg")
		results := makeArgs(pkg, method, signature.Results(), "res")
		intf.addMethod(NewMethod(pkg, method, args, results))
	}
	return intf
}

func makeArgs(pkg *Package, obj types.Object, args *types.Tuple, prefix string) []*MethodArg {
	methodArgs := make([]*MethodArg, 0)
	for i := 0; i < args.Len(); i++ {
		arg := args.At(i)
		name := arg.Name()
		if name == "" {
			name = fmt.Sprintf("%s%d", prefix, i+1)
		}
		methodArgs = append(methodArgs, NewMethodArg(pkg, arg, name))
	}
	return methodArgs
}
