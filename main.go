// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

import (
	"encoding/xml"
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
	log.SetPrefix(myname + " debug: ") // we only log debug output
	templates := NewStrings()
	flag.Var(templates, "t", "generate output using `template`")
	flag.StringVar(&outputFilename, "o", "", "output `filename`")
	writeXML := flag.Bool("write-xml", false, "Output context as XML")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", myname, "[options] filename...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	fail := func(err error) {
		fmt.Fprint(os.Stderr, myname, ": ", err.Error(), "\n")
		os.Exit(1)
	}
	for _, path := range flag.Args() {
		pkg, err := Parse(path)
		if err != nil {
			fail(err)
		}
		if *writeXML {
			WriteXML(pkg)
		} else {
			if err := Ridl(pkg, path, templates.Slice()); err != nil {
				fail(err)
			}
		}
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

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}

func Ridl(pkg *Package, path string, templates []string) error {
	templateContext := NewContext(path, pkg)
	for _, templateName := range templates {
		var out io.WriteCloser = NopWriteCloser(os.Stdout)
		if outputFilename != "" {
			if file, err := MakeOutputFile(path, templateName); err != nil {
				return err
			} else {
			    out = file
			}
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

func Parse(path string) (*Package, error) {
	var err error
	fset := token.NewFileSet()
	filenames, err := filepath.Glob(filepath.Join(filepath.Dir(path), "*.ridl"))
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

func CleanTypename(t types.Type) string {
	return strings.TrimPrefix(t.String(), "untyped ")
}

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

func (p *Package) Const(obj *types.Const) {
	p.Declare(NewConstDecl(obj.Name(), CleanTypename(obj.Type()), obj.Val().ExactString()))
}

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

func (p *Package) Map(name string, obj *types.Map) {
	keytyp := obj.Key().String()
	valtyp := obj.Elem().String()
	p.Declare(NewMapDecl(name, keytyp, valtyp))
}

func (p *Package) Slice(name string, obj *types.Slice) {
	typ := obj.Elem().String()
	p.Declare(NewArrayDecl(name, typ, 0))
}

func (p *Package) Array(name string, obj *types.Array) {
	size := obj.Len()
	typ := obj.Elem().String()
	p.Declare(NewArrayDecl(name, typ, int(size)))
}

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

func (p *Package) Typedef(name string, obj *types.Basic) {
	p.Declare(NewTypedefDecl(name, obj.String()))
}

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

func WriteXML(pkg *Package) {
	e := xml.NewEncoder(os.Stdout)
	e.Encode(pkg)
}
