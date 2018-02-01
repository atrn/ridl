package main

import "fmt"

// A Strings is slice of strings used to collect command line arguments.
//
type Strings struct {
	col []string
}

func NewStrings() *Strings {
	return &Strings{make([]string, 0)}
}

func (c *Strings) String() string {
	return fmt.Sprintf("%v", c.col)
}

func (c *Strings) Set(val string) error {
	c.col = append(c.col, val)
	return nil
}

func (c *Strings) Slice() []string {
	return c.col
}
