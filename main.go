// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var outputFilename = ""

func main() {
	myname := filepath.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(myname + ": ")
	templates := NewStrings()
	flag.Var(templates, "t", "generate output using `template`")
	flag.StringVar(&outputFilename, "o", "", "output `filename`")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", myname, "[options]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	pkg, err := ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}
	err = Ridl(pkg, path, templates.Slice())
	if err != nil {
		log.Fatal(err)
	}
}

type nopWriteCloser struct {
	w io.Writer
}

func (n *nopWriteCloser) Write(data []byte) (int, error) {
	return n.w.Write(data)
}

func (*nopWriteCloser) Close() error {
	return nil
}

// NopWriteCloser adds an empty Close() implementation to an
// io.Writer to transform it to an io.WriteCloser.
//
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}

// Ridl performs the ridl "code-generation" step on a Package.
//
func Ridl(pkg *Package, path string, templates []string) error {
	templateContext := NewContext(path, pkg)
	for _, templateName := range templates {
		out := NopWriteCloser(os.Stdout)
		if outputFilename != "" {
			file, err := MakeOutputFile(path, templateName)
			if err != nil {
				return err
			}
			out = file
		}
		err := ExpandTemplate(templateName, templateContext, out)
		err2 := out.Close()
		if err == nil {
			err = err2
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseFiles parses all the .ridl files in the given directory
// and returns a new Package representing their content.
//
func ParseFiles(path string) (*Package, error) {
	var err error
	fset := token.NewFileSet()
	filenames, err := filepath.Glob(filepath.Join(path, "*.ridl"))
	if err != nil {
		return nil, err
	}
	files := make([]*ast.File, 0, len(filenames))
	for _, filename := range filenames {
		var file *ast.File
		file, err = parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	conf := types.Config{
		IgnoreFuncBodies:         true,
		Importer:                 importer.Default(),
		DisableUnusedImportCheck: true,
	}
	pkg, err := conf.Check("", fset, files, nil)
	if err != nil {
		return nil, err
	}

	p := NewPackage(pkg.Name())
	for _, imported := range pkg.Imports() {
		p.Import(imported.Path())
	}
	scope := pkg.Scope()
	for i := 0; i < scope.Len(); i++ {
		obj := scope.Lookup(scope.Names()[i])
		switch actual := obj.(type) {
		case *types.Const:
			p.Const(actual)
		case *types.TypeName:
			p.TypeName(actual)
		default:
			log.Printf("X1:  %T  ->  %#v\n", actual, actual)
		}
	}
	return p, nil
}

// CleanTypename removes the "untypeed" suffix appended to
// the type names of untyped numeric constants.
//
func CleanTypename(t types.Type) string {
	return strings.TrimPrefix(t.String(), "untyped ")
}

// MakeMethodArgs builds a slice of new MethodArg values for
// the arguments in the supplied types.Tuple.
//
func MakeMethodArgs(args *types.Tuple, prefix string) []*MethodArg {
	ma := make([]*MethodArg, 0)
	for i := 0; i < args.Len(); i++ {
		arg := args.At(i)
		name := arg.Name()
		if name == "" {
			name = fmt.Sprintf("%s%d", prefix, i+1)
		}
		ma = append(ma, NewMethodArg(name, CleanTypename(arg.Type())))
	}
	return ma
}

// Const adds a declaration of a constant to the receiver.
//
func (p *Package) Const(obj *types.Const) {
	p.Declare(NewConstDecl(obj.Name(), CleanTypename(obj.Type()), obj.Val().ExactString()))
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
			continue
		}
		f := NewStructField(avar.Name(), avar.Type().String())
		astruct.AddField(f)
	}
	p.Declare(astruct)
}

// Interface adds an interface declaration to the receiver.
//
func (p *Package) Interface(name string, obj *types.Interface) {
	xi := NewInterfaceDecl(name)
	for i := 0; i < obj.NumMethods(); i++ {
		fn := obj.Method(i)
		sig := fn.Type().(*types.Signature)
		args := MakeMethodArgs(sig.Params(), "arg")
		results := MakeMethodArgs(sig.Results(), "res")
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

// MakeOutputFile creates a file for storing output.
//
func MakeOutputFile(path, templateName string) (io.WriteCloser, error) {
	if outputFilename != "" {
		return os.Create(outputFilename)
	}
	sansExt := func(s string) string {
		return strings.TrimSuffix(s, filepath.Ext(s))
	}
	name := sansExt(path) + "-" + sansExt(filepath.Base(templateName)) + ".out"
	return os.Create(name)
}
