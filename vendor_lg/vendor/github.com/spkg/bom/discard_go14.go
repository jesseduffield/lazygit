// +build !go1.5

package bom

import "bufio"

func discardBytes(buf *bufio.Reader, n int) {
	// cannot use the buf.Discard method as it was introduced in Go 1.5
	for i := 0; i < n; i++ {
		buf.ReadByte()
	}
}
