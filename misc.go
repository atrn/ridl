package main

import "io"

type nopWriteCloser struct {
	w io.Writer
}

func (n *nopWriteCloser) Write(data []byte) (int, error) {
	return n.w.Write(data)
}

func (*nopWriteCloser) Close() error {
	return nil
}

// NopWriteCloser adds an empty Close() implementation to an
// io.Writer to make it an io.WriteCloser.
//
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}
