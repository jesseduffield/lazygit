//go:build go1.1 && !go1.12
// +build go1.1,!go1.12

package bluemonday

import "io"

type stringWriterWriter interface {
	io.Writer
	StringWriter
}

type StringWriter interface {
	WriteString(s string) (n int, err error)
}
