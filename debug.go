package main

import (
	"fmt"
)

// Dump outputs a package in a somewhat readable format.
//
func (p *Package) Dump() {
	for _, path := range p.Imports {
		fmt.Printf("#include %q\n", path)
	}
	fmt.Println("namespace", p.PackageName, "{")
	fmt.Println("")
	for _, decl := range p.Decls {
		switch actual := decl.(type) {
		case *TypedefDecl:
			fmt.Print("typedef ", decl.Typename(), " ", decl.Name(), ";\n")
			fmt.Println("")
		case *StructDecl:
			fmt.Print(decl.Typename(), " {\n")
			for _, field := range actual.Fields {
				fmt.Print("    ", cpptype(field.Typename(), false), " _", field.Name(), ";\n")
			}
			fmt.Println("};")
			fmt.Println("")
		case *ConstDecl:
			fmt.Print("extern const ", cpptype(decl.Typename(), false), " ", decl.Name(), ";\n")
			fmt.Println("")
		case *InterfaceDecl:
			fmt.Println("class ", decl.Name(), "{")
			fmt.Println("public:")
			fmt.Print("    virtual ~", decl.Name(), "();\n")
			for _, method := range actual.Methods {
				fmt.Print("    virtual void ", method.Name(), "(")
				argsep := ""
				for _, arg := range method.Args {
					fmt.Print(argsep, cpptype(arg.Typename(), true), " ", arg.Name())
					argsep = ", "
				}
				for index, res := range method.Results {
					name := res.Name()
					if name == "" {
						name = fmt.Sprintf("out%d", index+1)
					}
					fmt.Print(argsep, cpptype(res.Typename(), false), " *", name)
				}
				fmt.Println(") = 0;")
			}
			fmt.Println("};")
			fmt.Println("")
		default:
			fmt.Print(cpptype(decl.Typename(), false), " ", decl.Name(), ";\n")
		}
	}
	fmt.Println("} // namespace", p.PackageName)
}
