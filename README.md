# ridl - retargetable IDL compiler

`ridl` (pronounced _riddle_) is an _Interface Description Language_
compiler. Of sorts. Compared with other tools that call themselves IDL
compilers `ridl` stands out. Firstly because it does not use an
interface definition language. And secondly because it doesn't
actually _compile_ anything. ridl is an an _IDL compiler_. With no
IDL. That it doesn't compile.

## What ridl does

Instead of defining an IDL `ridl` uses an existing language, Go, and
re-interprets its constant and type declarations to be definitions of
_interfaces_ (which, of course, they already are). `ridl` parses Go
code, using standard Go packages, builds a data structure to represent
the constants and types the code defines, and then uses Go's standard
template package to _execute_ user-supplied templates with this data
structure as the so-called _template context_. The result permits the
Go declarations to transformed to whatever is produced by the
templates.

ridl comes with example templates to convert Go declarations to C++.
These will be developed to allow for network messaging, using
different transports, with Go interface types defining sets of
messages.


## What ridl doesn't do

Being based on Go means ridl has no unions and offers no direct
support for interface versioning.

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

## Invoking `ridl`

Basic usage is,

`ridl` _options_ _[_ _filename... _]_

If no input files are named on the command line `ridl` reads all files
with the extension `.ridl` in the current directory.  If one, or more,
files are named on the command line only those files are processed.

### Options
- -t _template_  
Use _template_ to generate output.  More than one template may be used
in which case 
- -o _filename_  
Write output to _filename_.  The _filename_ `-` represents the
standard output.
- -D _directory_  
Read template files from _directory_.
- -version  
Output ridl version and exit.
- -debug  
Enable debug output.

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

## Templates

### Template Context

TBD.

### Template Functions

#### argtype
Returns a string with the C++ type corresponding to the given Go
type when a value of that type is used an argument to a function.
#### basename
Returns the base filename for a given pathname **without** any
extension.
#### cpptype
Returns a string with the C++ type corresponding to the given
Go type.
#### isslice
Determines if a Go type is a slice, i.e. if the Go type name
starts with the sequence `[]`.
#### tolower
Returns the lower case version of a string.
#### plus
Returns the sum of two integers.
#### eltype
Returns the element type of an array or slice type.
#### restype
Returns the C++ type to be used as a function result.
#### dims
TBD.
