package main

import "fmt"

// A StringSlice is slice of strings used to collect command line arguments.
//
type StringSlice []string

// NewStringSlice returns a new, empty, StringSlice.
//
func NewStringSlice() *StringSlice {
	s := make(StringSlice, 0, 4) // total guess at an inital capacity
	return &s
}

// String implements the Stringer interface for StringSlice.  It is
// primarily used for debugging purposes so uses the fmt package to
// return the receiver's strings as a Go-literal.
//
func (s *StringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set implements the flag.Value interface and accumulates the
// supplied value by appending it to the slice.
//
func (s *StringSlice) Set(val string) error {
	*s = append(*s, val)
	return nil
}

// Slice returns the receiver's collection of strings.
//
func (s *StringSlice) Slice() []string {
	return *s
}

// Len returns the number of strings in the receiver
//
func (s *StringSlice) Len() int {
	return len(*s)
}
