// +build go1.5

package bom

import "bufio"

func discardBytes(buf *bufio.Reader, n int) {
	// the Discard method was introduced in Go 1.5
	buf.Discard(n)
}
