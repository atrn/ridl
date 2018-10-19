# TODOs

## Support more types

## Messaging

Turn the methods of interface types into actual message structures and
assign message codes (ints, not strings) for de-mux. Generate code to
service requests and code to send them.

## Enums

The Go way of defining enumerations is a little verbose. A more direct
representation a la my go generate 'enums' tool would be nicer. Not a
huge problem but will annoy people used to DCE/CORBA IDL and C/C++.
