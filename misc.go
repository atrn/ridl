package main

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func logdebug(format string, args ...interface{}) {
	if *debugFlag {
		log.Printf("DEBUG: "+format, args...)
	}
}

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

func MustGetUsername() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return u.Username
}

func MustGetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return hostname
}

func TrimExtension(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}

func TrimUntyped(t string) string {
	return strings.TrimPrefix(t, "untyped ")
}
