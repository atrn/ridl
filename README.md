# Retargetable interface compiler

`ridl` (pronounced _riddle_) is an _Interface Description Language_
compiler. Of sorts. Compared to other tools that call themselves IDL
compilers `ridl` stands out. Firstly because it does not use an
interface definition language. And secondly because it doesn't
actually _compile_ anything. ridl is an an _IDL compiler_ without
an IDL, that it doesn't compile.

## What ridl does

Instead of defining an IDL `ridl` uses an existing language, Go, and
re-interprets its constant and type declarations as definitions of
_interfaces_. `ridl` parses the Go code, using standard Go packages,
builds a data structure to represent the constants and types defined
by the code, and then uses Go's standard template package to _execute_
user-supplied templates using this data structure as the so-called
_template context_. The result is a transformation of the Go
declarations to whatever is produced by the templates.

ridl comes with some templates to demonstrate the approach. The
primary template is used to convert the Go declarations to C++.  This
is the first usage and ridl's initial features support this use.

## What ridl doesn't do

Other IDLs and RPC or protocol description languages are more
expressive. Being based on Go ridl has no unions and offers no direct
support for interfsce versioning.

### Missing features

- no _union_ types
- no _versioning_

## What is an interface?

An _interface_ is a collection of constants, data structures and
_messages_. Users use these to interact with some sort of _service_,
typically sending and receiving the messages defined by the interfaces
with the associated types and constant values.

The specific manner in how this messaging occurs varies
considerably. There are, and have been, numerous frameworks, systems
and libraries to tackle this problem. Far too many to list
here. However, regardless of the _encoding_ used by the underling
system, and the message _transport_ (with whatever _addressing_ it
uses) all have similar basic concerns about defining the constant
values, types and messages used in the inteface.

## Interface Descriptions

A ridl interface description is a Go package containing only constant
and type declarations, i.e. there are no functions. By convention
these are named with a `.ridl` extension but the file is a `.go` file
and can be processed by standard Go tools.

Like Go, or more correctly because of it, `ridl` has a module system
based upon _packages_.  Each interface belongs to a package and an
interface must declare its pacakge as its first statement (ridl *is*
go).

Packages may be used in template expansion. The C++ template uses the
Go package to name a C++ namespace.

## Interface Types

`ridl`, or more correctly the supplied C++ template file, interprets a
Go `interface` type as the definition of a set of messages. Each
method represents a distinct _message_ with the interface where a
method's name identifies the message within the interface.  Any
arguments to the method define the message _payload_, the data
communicated to the receiver along with the message identity.  If the
method has results they define the payload of a _result_ message
(who's identity is derived from the method name).

## Arrays, Slices and Maps

Arrays, slice and map types may be used. They are easily mapped to
some appropriate construct in most target languages.

## `error`

Go's `error` type is permitted and is mapped to whatever
representation of errors is used with the target language.

The error type's interface is limited, only permitting extracting the
error as a string, so ridl defines errors as string.

## Restrictions

Function, channel and pointer types are not permitted.

Function declarations are by default errors. A _permissive_ mode could
be used to parse Go files and ignore constructs not permitted in ridl
files.

## Assumed Communications Semantics

The actual semantics of the messaging are not defined by `ridl` itself
but depend on the underlying code generator and target environment.

## Code Generation

### Context Data Structure

... define Context here so people can write templates!

### Template functions

#### cpptype
#### argtype
#### basename
