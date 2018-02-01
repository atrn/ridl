# `ridl` - a retargetable interface compiler

`ridl` (pronounced _riddle_) is an _IDL compiler_. Of sorts. When
compared with other tools that call themselves IDL compilers `ridl`
stands out. First because `ridl` does not have an interface definition
language. There is no IDL. And second `ridl` doesn't actual _compile_
anything.  So, _ridl_ is an an _IDL compiler_ with no IDL and that
doesn't compile. Sound good?

## What ridl does

Instead of defining an IDL `ridl` uses an existing programming
language, Go, and interprets its type and constant declarations
as definitions of _interfaces_. `ridl` parses Go code, using
standard Go packages, and builds a data structure to represent
the various entities defined by interface definitions. `ridl`
then uses Go's standard template processing package to _execute_
user-defined template files using the data structure as the
so-called _template context_. The result is a transformation
of the declarations in the input to whatever output produced
by the supplied template files. What passes for _compilation_.

- `ridl` has no IDL, it just uses a subset of Go.
- `ridl` doesn't compile, it expands templates

An example template converts Go to C++.

## What ridl doesn't do

Other languages are more expressive when it comes to data
structuring and other details. Being based on Go and its
syntax ridl has no unions nor does it provide any support
for the versioning issue.

### Missing features list

- no _union_ types
- versioning is not expressed


## So what is an interface?

An _interface_ is a collection of constants, data structures and
_messages_. Users use these to interact with some sort of _service_.
The details of how this is done vary considerably. Numerous frameworks,
systems and libraries exist to tackle this problem. Far too many to
list here. But regardless of the _encoding_ used and the _transport_
and _addressing_ requirements they are have similar base concerns
regarding defining values, types and messages.

## Using go as an IDL

`ridl` reads the Go constant and type declarations, using standard Go
pacakges, and creates a data structure representing the entities. This
data structure is then used to generate output using macro processing,
which, these days, we call _template expansion_.

Standard Go packages are used to read the code and perform the
template processing.

## Defining Interfaces in Go

### Basic types

Go's basic types - boolean, integer, floating point, string - provide
the starting point.

### Aggregates

Arrays, slices and maps of all supported types.

### const

Constants of all supported types. As `ridl` uses standard Go tools and
semantics features such as `iota` are fully supported in interface
definitions.

### structures

Structures contains fields of all supported types.

Embedding is supported.

### interface types

`interface` types are supported. Code generation uses interfaces to
define _messages_.


# Definitions

## Interface Descriptions

A ridl interface description is a Go package only containing constant
and type declarations, i.e. no functions or methods. By convention
these are named with a `.ridl` extension but the file is a `.go` file
and can be processed by standard Go tools.

Like Go, or more correctly because of it, `ridl` has a module system
based on _packages_.  Each interface belongs to a package and an
interface must declare its pacakge as its first statement.

Packages may be used in code generation. For instance, C++ generators
typically map the Go package to a C++ namespace.

So to repeat the point, ridl interfaces are just Go declarations.

But...

## Interface Types

`ridl` interprets a Go `interface` type as the definition of a set of
messages where each method defined in the `interface` represents a
distinct _message_. The method name identifies the message within the
particular interface and any arguments define the message's _payload_,
the data communicated along with the message identification.  If the
method has results, they define the payload of a _result_ message.

## Arrays, Slices and Maps

Arrays, slice and map types may be used. They are easily mapped to
some appropriate construct in most target languages.

## `error`

Go's `error` type is permitted and is mapped to whatever
representation of errors is used with the target language.

The error type's interface is limited and only permits
extracting a string that communicates the error. We use
this to communicate errors as strings (along with a flag
to indicate the string is valid, empty error strings are
valid so may not be used to indicate a lack of error).

## Restrictions

Function, channel and pointer types are not permitted.

Function declarations are by default errors. A _permissive_
mode can be used to parse Go files and ignore constructs
not permitted in ridl files.

## Assumed Communications Semantics

The actual semantics of the messaging are not defined by `ridl` itself
but depend on the underlying code generator and target environment.

## Code Generation

### Context Data Structure

### Example

C++ with zeromq

This template generates C++ code for RIDL interfaces with ZeroMQ used
as the underlying transport. Message passing is implemented atop
ZeroMQ's C API. Messages are expressed as C++ PODs held in ZeroMQ
message types, i.e. ZeroMQ _owns_ the memory (avoiding copies).
## Implementation

`ridl` parses Go using the standard `go/types` package (and the other
packages it uses). Code-generation is then done using Go's
`text/template` package.  `ridl` parses the code using `go/types` and
converts the resultant AST into a _context_ data structure for use in
template expansion.

So `ridl` does no actual compilation but instead transforms the
`go/types` generated AST into a form more amenable for templates
and expands one or more templates.

