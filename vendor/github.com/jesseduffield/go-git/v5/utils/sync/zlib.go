package sync

import (
	"bytes"
	"compress/zlib"
	"io"
	"sync"
)

var (
	zlibInitBytes = []byte{0x78, 0x9c, 0x01, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x01}
	zlibReader    = sync.Pool{
		New: func() interface{} {
			r, _ := zlib.NewReader(bytes.NewReader(zlibInitBytes))
			return ZLibReader{
				Reader: r.(zlibReadCloser),
			}
		},
	}
	zlibWriter = sync.Pool{
		New: func() interface{} {
			return zlib.NewWriter(nil)
		},
	}
)

type zlibReadCloser interface {
	io.ReadCloser
	zlib.Resetter
}

type ZLibReader struct {
	dict   *[]byte
	Reader zlibReadCloser
}

// GetZlibReader returns a ZLibReader that is managed by a sync.Pool.
// Returns a ZLibReader that is reset using a dictionary that is
// also managed by a sync.Pool.
//
// After use, the ZLibReader should be put back into the sync.Pool
// by calling PutZlibReader.
func GetZlibReader(r io.Reader) (ZLibReader, error) {
	z := zlibReader.Get().(ZLibReader)
	z.dict = GetByteSlice()

	err := z.Reader.Reset(r, *z.dict)

	return z, err
}

// PutZlibReader puts z back into its sync.Pool, first closing the reader.
// The Byte slice dictionary is also put back into its sync.Pool.
func PutZlibReader(z ZLibReader) {
	z.Reader.Close()
	PutByteSlice(z.dict)
	zlibReader.Put(z)
}

// GetZlibWriter returns a *zlib.Writer that is managed by a sync.Pool.
// Returns a writer that is reset with w and ready for use.
//
// After use, the *zlib.Writer should be put back into the sync.Pool
// by calling PutZlibWriter.
func GetZlibWriter(w io.Writer) *zlib.Writer {
	z := zlibWriter.Get().(*zlib.Writer)
	z.Reset(w)
	return z
}

// PutZlibWriter puts w back into its sync.Pool.
func PutZlibWriter(w *zlib.Writer) {
	zlibWriter.Put(w)
}
