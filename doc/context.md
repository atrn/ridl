# ridl Template Context

## Context

| Variable    | Type            | Description                                      |
|:------------|:----------------|:-------------------------------------------------|
| PackageName | string          | Name of the package.                             |
| Decls       | []Decl          | Array of all declarations in the package.        |
| Imports     | []string        | Names of all imported packages.                  |
| RidlVersion | string          | Version of ridl being used.                      |
| Directory   | string          | Name of the directory being processed.           |
| Filenames   | []string        | Names of all .ridl files being processed.        |
| BuildTime   | time.Time       | Time of processing.                              |
| Username    | string          | Name of user running ridl.                       |
| Hostname    | string          | Name of host on which ridl is being run.         |
| Typedefs    | []TypedefDecl   | All type/alias declarations.                     |
| ArrayTypes  | []ArrayDecl     | All array type declarations.                     |
| MapTypes    | []MapDecl       | All map type declarations.                       |
| StructTypes | []StructDecl    | All struct type declarations.                    |
| Interfaces  | []InterfaceDecl | All interface declarations.                      |
| Constants   | []ConstDecl     | All constant declarations.                       |
| Enums       | []Enum          | All enum-like constant declarations.             |
| NotEnums    | []ConstDecl     | All constant declarations that are not enum-like |

## Decl

| Variable | Type           | Description                                  |
|:---------|:---------------|:---------------------------------------------|
| Name     | string         | The declarations's identifier.               |
| TypeName | string         | The name of the declaration's type.          |
| Kind     | DeclKind       | The kind of declaration (see below).         |
| Position | token.Position | The source of the declaration.               |

### DeclKind

| Value     | Description                          |
|:----------|:-------------------------------------|
| const     | A constant                           |
| type      | A type (aka typedef)                 |
| array     | An array                             |
| struct    | A struct                             |
| field     | A struct field                       |
| map       | A map                                |
| interface | An interface                         |
| method    | An interface method (function)       |
| argument  | An argument to or result of a method |

### Delcaration Predicates

| Predicate     | Description                                                  |
|:--------------|:-------------------------------------------------------------|
| IsConst       | True if the declaration declares a constant                  |
| IsTypedef     | True if the declaration declares a type                      |
| IsArray       | True if the declaration declares an array type               |
| IsStruct      | True if the declaration declares an struct type              |
| IsStructField | True if the declaration declares a field of a struct         |
| IsMap         | True if the declaration declares a map type                  |
| IsInterface   | True if the declaration declares an interface                |
| IsMethod      | True if the declaration declares a method                    |
| IsMethodArg   | True if the declaration declares a method argument or result |


## ConstDecl

| Variable     | Type           | Description                                                 |
|:-------------|:---------------|:------------------------------------------------------------|
| Value        | constant.Value | The constant's computed value.                              |
| ExactValue   | string         | The constant's exact value as per go/types                  |
| IsEnumerator | bool           | True if this constant is an enumerator of an enum-like type |

## TypedefDecl

| Variable | Type   | Description                                                             |
|:---------|:-------|:------------------------------------------------------------------------|
| TypeName | string | The name of the underlying type                                         |
| IsEnum   | bool   | True of the type is used as the underlying type for enum-like constants |


## ArrayDecl

| Variable         | Type   | Description                                             |
|:-----------------|:-------|:--------------------------------------------------------|
| Length           | int    | Number of elements in the array or 0 if variably sized. |
| ElTypeName       | string | Name of the element type.                               |
| TypeName         | string | Go representation of the array type.                    |
| IsVariableLength | bool   | True if the array has variable size.                    |
| IsFixedLength    | bool   | True if the array has a fixed size.                     |

## StructDecl

| Variable         | Type   | Description                                             |
|:-----------------|:-------|:--------------------------------------------------------|
| TypeName | string | Name of the type |

## StructField

| Variable  | Type   | Description                       |
|:----------|:-------|:----------------------------------|
| Name      | string | Name of the field                 |
| Offset    | int    | Offset, in bytes, of the field    |
| Alignment | int    | ALignment, in bytes, of the field |
| Tags      | []Tag  | Tags associated with the field    |

### Tag

| Variable | Type   | Description              |
|:---------|:-------|:-------------------------|
| Key      | string | Name of the field tag    |
| Value    | string | Value of the field's tag |

### MapDecl

| Variable | Type       | Description                      |
|:---------|:-----------|:---------------------------------|
| Key      | types.Type | The type of the map's key values |
| Value    | types.Type | The type of the map's values     |

### InterfaceDecl

| Variable | Type         | Description              |
|:---------|:-------------|:-------------------------|
| Methods  | []MethodDecl | Methods of the interface |

### MethodDecl

| Variable | Type            | Description |
|:---------|:----------------|:------------|
| TypeName | string          |             |
| Args     | []MethodArgDecl |             |
| Results  | []MethodArgDecl |             |


### Enum

| Variable    | Type        | Description                                    |
|:------------|:------------|:-----------------------------------------------|
| Type        | TypedefDecl | Type of enumerators.                           |
| Enumerators | []ConstDecl | The enumerators                                |
| IsDense     | bool        | True if the enumerator values form a dense set |


## Template Functions
