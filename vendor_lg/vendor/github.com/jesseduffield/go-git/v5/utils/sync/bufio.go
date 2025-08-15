package sync

import (
	"bufio"
	"io"
	"sync"
)

var bufioReader = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

// GetBufioReader returns a *bufio.Reader that is managed by a sync.Pool.
// Returns a bufio.Reader that is reset with reader and ready for use.
//
// After use, the *bufio.Reader should be put back into the sync.Pool
// by calling PutBufioReader.
func GetBufioReader(reader io.Reader) *bufio.Reader {
	r := bufioReader.Get().(*bufio.Reader)
	r.Reset(reader)
	return r
}

// PutBufioReader puts reader back into its sync.Pool.
func PutBufioReader(reader *bufio.Reader) {
	bufioReader.Put(reader)
}
