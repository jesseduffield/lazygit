package displaywidth

import (
	"github.com/clipperhouse/uax29/v2/graphemes"
)

// Graphemes is an iterator over grapheme clusters.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
type Graphemes[T ~string | []byte] struct {
	iter    *graphemes.Iterator[T]
	options Options
}

// Next advances the iterator to the next grapheme cluster.
func (g *Graphemes[T]) Next() bool {
	return g.iter.Next()
}

// Value returns the current grapheme cluster.
func (g *Graphemes[T]) Value() T {
	return g.iter.Value()
}

// Width returns the display width of the current grapheme cluster.
func (g *Graphemes[T]) Width() int {
	return graphemeWidth(g.Value(), g.options)
}

// StringGraphemes returns an iterator over grapheme clusters for the given
// string.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func StringGraphemes(s string) Graphemes[string] {
	return DefaultOptions.StringGraphemes(s)
}

// StringGraphemes returns an iterator over grapheme clusters for the given
// string, with the given options.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func (options Options) StringGraphemes(s string) Graphemes[string] {
	g := graphemes.FromString(s)
	g.AnsiEscapeSequences = options.ControlSequences
	g.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

	return Graphemes[string]{iter: g, options: options}
}

// BytesGraphemes returns an iterator over grapheme clusters for the given
// []byte.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func BytesGraphemes(s []byte) Graphemes[[]byte] {
	return DefaultOptions.BytesGraphemes(s)
}

// BytesGraphemes returns an iterator over grapheme clusters for the given
// []byte, with the given options.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func (options Options) BytesGraphemes(s []byte) Graphemes[[]byte] {
	g := graphemes.FromBytes(s)
	g.AnsiEscapeSequences = options.ControlSequences
	g.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

	return Graphemes[[]byte]{iter: g, options: options}
}
