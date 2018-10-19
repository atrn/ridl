package main

import "fmt"

// A StringSlice is slice of strings used to collect command line arguments.
//
type StringSlice []string

// NewStringSlice returns a new, empty, StringSlice.
//
func NewStringSlice() *StringSlice {
	s := make(StringSlice, 0, 8)
	return &s
}

// String implements the Stringer interface for StringSlice.
//
func (s *StringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set implements the flag.Value interface
//
func (s *StringSlice) Set(val string) error {
	*s = append(*s, val)
	return nil
}

// Slice returns the collection of strings.
//
func (s *StringSlice) Slice() []string {
	return *s
}
