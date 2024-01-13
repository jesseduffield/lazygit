//go:build go1.12
// +build go1.12

package bluemonday

import "io"

type stringWriterWriter interface {
	io.Writer
	io.StringWriter
}
