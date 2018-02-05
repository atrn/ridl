package main

import "fmt"

// A Strings is slice of strings used to collect command line arguments.
//
type Strings []string

// NewStrings returns a new, empty, Strings.
//
func NewStrings() *Strings {
	s := make(Strings, 0, 8)
	return &s
}

func (s *Strings) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set implements the flag.Value interface
func (s *Strings) Set(val string) error {
	*s = append(*s, val)
	return nil
}

// Slice returns the collection of strings.
func (s *Strings) Slice() []string {
	return *s
}
