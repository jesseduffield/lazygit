package graphemes

import "github.com/clipperhouse/uax29/v2/internal/iterators"

type Iterator[T iterators.Stringish] struct {
	*iterators.Iterator[T]
}

var (
	splitFuncString = splitFunc[string]
	splitFuncBytes  = splitFunc[[]byte]
)

// FromString returns an iterator for the grapheme clusters in the input string.
// Iterate while Next() is true, and access the grapheme via Value().
func FromString(s string) Iterator[string] {
	return Iterator[string]{
		iterators.New(splitFuncString, s),
	}
}

// FromBytes returns an iterator for the grapheme clusters in the input bytes.
// Iterate while Next() is true, and access the grapheme via Value().
func FromBytes(b []byte) Iterator[[]byte] {
	return Iterator[[]byte]{
		iterators.New(splitFuncBytes, b),
	}
}
