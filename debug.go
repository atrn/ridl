package main

import (
	"fmt"
)

// Write a package as C++
//
func (p *Package) Dump() {
	for _, path := range p.Imports {
		fmt.Printf("#include %q\n", path)
	}
	fmt.Println("namespace", p.PackageName, "{\n")
	for _, decl := range p.Decls {
		switch actual := decl.(type) {
		case *TypedefDecl:
			fmt.Print("typedef ", decl.Type(), " ", decl.Name(), ";\n")
			fmt.Println("")
		case *StructDecl:
			fmt.Print(decl.Type(), " {\n")
			for _, field := range actual.Fields {
				fmt.Print("    ", cpptype(field.Type(), false), " _", field.Name(), ";\n")
			}
			fmt.Println("};")
			fmt.Println("")
		case *ConstDecl:
			fmt.Print("extern const ", cpptype(decl.Type(), false), " ", decl.Name(), ";\n")
			fmt.Println("")
		case *InterfaceDecl:
			fmt.Println("class ", decl.Name(), "{")
			fmt.Println("public:")
			fmt.Print("    virtual ~", decl.Name(), "();\n")
			for _, method := range actual.Methods {
				fmt.Print("    virtual void ", method.Name(), "(")
				argsep := ""
				for _, arg := range method.Args {
					fmt.Print(argsep, cpptype(arg.Type(), true), " ", arg.Name())
					argsep = ", "
				}
				for index, res := range method.Results {
					name := res.Name()
					if name == "" {
						name = fmt.Sprintf("out%d", index+1)
					}
					fmt.Print(argsep, cpptype(res.Type(), false), " *", name)
				}
				fmt.Println(") = 0;")
			}
			fmt.Println("};")
			fmt.Println("")
		default:
			fmt.Print(cpptype(decl.Type(), false), " ", decl.Name(), ";\n")
		}
	}
	fmt.Println("} // namespace", p.PackageName)
}
