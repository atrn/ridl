// ridl - re-targetable IDL compiler
// Copyright © 2016 A.Newman.
//
// This file is licensed using the GNU Public License, version 2.
// See the file LICENSE for details.
//

package main

// Enum represents a C/C++ enumerated type that has been emulated
// using the Go idiom of defining a type and a series of constants of
// that type.
type Enum struct {
	Typedef   *TypedefDecl
	Constants []*ConstDecl
}
