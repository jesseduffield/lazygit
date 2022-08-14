package uniseg

import "unicode/utf8"

// Graphemes implements an iterator over Unicode grapheme clusters, or
// user-perceived characters. While iterating, it also provides information
// about word boundaries, sentence boundaries, and line breaks.
//
// After constructing the class via [NewGraphemes] for a given string "str",
// [Next] is called for every grapheme cluster in a loop until it returns false.
// Inside the loop, information about the grapheme cluster as well as boundary
// information is available via the various methods (see examples below).
//
// Using this class to iterate over a string is convenient but it is much slower
// than using this package's [Step] or [StepString] functions or any of the
// other specialized functions starting with "First".
type Graphemes struct {
	// The original string.
	original string

	// The remaining string to be parsed.
	remaining string

	// The current grapheme cluster.
	cluster string

	// The byte offset of the current grapheme cluster relative to the original
	// string.
	offset int

	// The current boundary information of the Step() parser.
	boundaries int

	// The current state of the Step() parser.
	state int
}

// NewGraphemes returns a new grapheme cluster iterator.
func NewGraphemes(s string) *Graphemes {
	return &Graphemes{
		original:  s,
		remaining: s,
		state:     -1,
	}
}

// Next advances the iterator by one grapheme cluster and returns false if no
// clusters are left. This function must be called before the first cluster is
// accessed.
func (g *Graphemes) Next() bool {
	if len(g.remaining) == 0 {
		// We're already past the end.
		g.state = -2
		g.cluster = ""
		return false
	}
	g.offset += len(g.cluster)
	g.cluster, g.remaining, g.boundaries, g.state = StepString(g.remaining, g.state)
	return true
}

// Runes returns a slice of runes (code points) which corresponds to the current
// grapheme cluster. If the iterator is already past the end or [Next] has not
// yet been called, nil is returned.
func (g *Graphemes) Runes() []rune {
	if g.state < 0 {
		return nil
	}
	return []rune(g.cluster)
}

// Str returns a substring of the original string which corresponds to the
// current grapheme cluster. If the iterator is already past the end or [Next]
// has not yet been called, an empty string is returned.
func (g *Graphemes) Str() string {
	return g.cluster
}

// Bytes returns a byte slice which corresponds to the current grapheme cluster.
// If the iterator is already past the end or [Next] has not yet been called,
// nil is returned.
func (g *Graphemes) Bytes() []byte {
	if g.state < 0 {
		return nil
	}
	return []byte(g.cluster)
}

// Positions returns the interval of the current grapheme cluster as byte
// positions into the original string. The first returned value "from" indexes
// the first byte and the second returned value "to" indexes the first byte that
// is not included anymore, i.e. str[from:to] is the current grapheme cluster of
// the original string "str". If [Next] has not yet been called, both values are
// 0. If the iterator is already past the end, both values are 1.
func (g *Graphemes) Positions() (int, int) {
	if g.state == -1 {
		return 0, 0
	} else if g.state == -2 {
		return 1, 1
	}
	return g.offset, g.offset + len(g.cluster)
}

// IsWordBoundary returns true if a word ends after the current grapheme
// cluster.
func (g *Graphemes) IsWordBoundary() bool {
	if g.state < 0 {
		return true
	}
	return g.boundaries&MaskWord != 0
}

// IsSentenceBoundary returns true if a sentence ends after the current
// grapheme cluster.
func (g *Graphemes) IsSentenceBoundary() bool {
	if g.state < 0 {
		return true
	}
	return g.boundaries&MaskSentence != 0
}

// LineBreak returns whether the line can be broken after the current grapheme
// cluster. A value of [LineDontBreak] means the line may not be broken, a value
// of [LineMustBreak] means the line must be broken, and a value of
// [LineCanBreak] means the line may or may not be broken.
func (g *Graphemes) LineBreak() int {
	if g.state == -1 {
		return LineDontBreak
	}
	if g.state == -2 {
		return LineMustBreak
	}
	return g.boundaries & MaskLine
}

// Reset puts the iterator into its initial state such that the next call to
// [Next] sets it to the first grapheme cluster again.
func (g *Graphemes) Reset() {
	g.state = -1
	g.offset = 0
	g.cluster = ""
	g.remaining = g.original
}

// GraphemeClusterCount returns the number of user-perceived characters
// (grapheme clusters) for the given string.
func GraphemeClusterCount(s string) (n int) {
	state := -1
	for len(s) > 0 {
		_, s, _, state = FirstGraphemeClusterInString(s, state)
		n++
	}
	return
}

// FirstGraphemeCluster returns the first grapheme cluster found in the given
// byte slice according to the rules of Unicode Standard Annex #29, Grapheme
// Cluster Boundaries. This function can be called continuously to extract all
// grapheme clusters from a byte slice, as illustrated in the example below.
//
// If you don't know the current state, for example when calling the function
// for the first time, you must pass -1. For consecutive calls, pass the state
// and rest slice returned by the previous call.
//
// The "rest" slice is the sub-slice of the original byte slice "b" starting
// after the last byte of the identified grapheme cluster. If the length of the
// "rest" slice is 0, the entire byte slice "b" has been processed. The
// "cluster" byte slice is the sub-slice of the input slice containing the
// identified grapheme cluster.
//
// Given an empty byte slice "b", the function returns nil values.
//
// While slightly less convenient than using the Graphemes class, this function
// has much better performance and makes no allocations. It lends itself well to
// large byte slices.
//
// The "reserved" return value is a placeholder for future functionality and may
// be ignored for the time being.
func FirstGraphemeCluster(b []byte, state int) (cluster, rest []byte, reserved, newState int) {
	// An empty byte slice returns nothing.
	if len(b) == 0 {
		return
	}

	// Extract the first rune.
	r, length := utf8.DecodeRune(b)
	if len(b) <= length { // If we're already past the end, there is nothing else to parse.
		return b, nil, 0, grAny
	}

	// If we don't know the state, determine it now.
	if state < 0 {
		state, _ = transitionGraphemeState(state, r)
	}

	// Transition until we find a boundary.
	var boundary bool
	for {
		r, l := utf8.DecodeRune(b[length:])
		state, boundary = transitionGraphemeState(state, r)

		if boundary {
			return b[:length], b[length:], 0, state
		}

		length += l
		if len(b) <= length {
			return b, nil, 0, grAny
		}
	}
}

// FirstGraphemeClusterInString is like [FirstGraphemeCluster] but its input and
// outputs are strings.
func FirstGraphemeClusterInString(str string, state int) (cluster, rest string, reserved, newState int) {
	// An empty string returns nothing.
	if len(str) == 0 {
		return
	}

	// Extract the first rune.
	r, length := utf8.DecodeRuneInString(str)
	if len(str) <= length { // If we're already past the end, there is nothing else to parse.
		return str, "", 0, grAny
	}

	// If we don't know the state, determine it now.
	if state < 0 {
		state, _ = transitionGraphemeState(state, r)
	}

	// Transition until we find a boundary.
	var boundary bool
	for {
		r, l := utf8.DecodeRuneInString(str[length:])
		state, boundary = transitionGraphemeState(state, r)

		if boundary {
			return str[:length], str[length:], 0, state
		}

		length += l
		if len(str) <= length {
			return str, "", 0, grAny
		}
	}
}
