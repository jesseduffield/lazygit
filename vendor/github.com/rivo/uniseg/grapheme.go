package uniseg

import "unicode/utf8"

// The states of the grapheme cluster parser.
const (
	grAny = iota
	grCR
	grControlLF
	grL
	grLVV
	grLVTT
	grPrepend
	grExtendedPictographic
	grExtendedPictographicZWJ
	grRIOdd
	grRIEven
)

// The grapheme cluster parser's breaking instructions.
const (
	grNoBoundary = iota
	grBoundary
)

// The grapheme cluster parser's state transitions. Maps (state, property) to
// (new state, breaking instruction, rule number). The breaking instruction
// always refers to the boundary between the last and next code point.
//
// This map is queried as follows:
//
//   1. Find specific state + specific property. Stop if found.
//   2. Find specific state + any property.
//   3. Find any state + specific property.
//   4. If only (2) or (3) (but not both) was found, stop.
//   5. If both (2) and (3) were found, use state and breaking instruction from
//      the transition with the lower rule number, prefer (3) if rule numbers
//      are equal. Stop.
//   6. Assume grAny and grBoundary.
var grTransitions = map[[2]int][3]int{
	// GB5
	{grAny, prCR}:      {grCR, grBoundary, 50},
	{grAny, prLF}:      {grControlLF, grBoundary, 50},
	{grAny, prControl}: {grControlLF, grBoundary, 50},

	// GB4
	{grCR, prAny}:        {grAny, grBoundary, 40},
	{grControlLF, prAny}: {grAny, grBoundary, 40},

	// GB3.
	{grCR, prLF}: {grAny, grNoBoundary, 30},

	// GB6.
	{grAny, prL}: {grL, grBoundary, 9990},
	{grL, prL}:   {grL, grNoBoundary, 60},
	{grL, prV}:   {grLVV, grNoBoundary, 60},
	{grL, prLV}:  {grLVV, grNoBoundary, 60},
	{grL, prLVT}: {grLVTT, grNoBoundary, 60},

	// GB7.
	{grAny, prLV}: {grLVV, grBoundary, 9990},
	{grAny, prV}:  {grLVV, grBoundary, 9990},
	{grLVV, prV}:  {grLVV, grNoBoundary, 70},
	{grLVV, prT}:  {grLVTT, grNoBoundary, 70},

	// GB8.
	{grAny, prLVT}: {grLVTT, grBoundary, 9990},
	{grAny, prT}:   {grLVTT, grBoundary, 9990},
	{grLVTT, prT}:  {grLVTT, grNoBoundary, 80},

	// GB9.
	{grAny, prExtend}: {grAny, grNoBoundary, 90},
	{grAny, prZWJ}:    {grAny, grNoBoundary, 90},

	// GB9a.
	{grAny, prSpacingMark}: {grAny, grNoBoundary, 91},

	// GB9b.
	{grAny, prPreprend}: {grPrepend, grBoundary, 9990},
	{grPrepend, prAny}:  {grAny, grNoBoundary, 92},

	// GB11.
	{grAny, prExtendedPictographic}:                     {grExtendedPictographic, grBoundary, 9990},
	{grExtendedPictographic, prExtend}:                  {grExtendedPictographic, grNoBoundary, 110},
	{grExtendedPictographic, prZWJ}:                     {grExtendedPictographicZWJ, grNoBoundary, 110},
	{grExtendedPictographicZWJ, prExtendedPictographic}: {grExtendedPictographic, grNoBoundary, 110},

	// GB12 / GB13.
	{grAny, prRegionalIndicator}:    {grRIOdd, grBoundary, 9990},
	{grRIOdd, prRegionalIndicator}:  {grRIEven, grNoBoundary, 120},
	{grRIEven, prRegionalIndicator}: {grRIOdd, grBoundary, 120},
}

// Graphemes implements an iterator over Unicode extended grapheme clusters,
// specified in the Unicode Standard Annex #29. Grapheme clusters correspond to
// "user-perceived characters". These characters often consist of multiple
// code points (e.g. the "woman kissing woman" emoji consists of 8 code points:
// woman + ZWJ + heavy black heart (2 code points) + ZWJ + kiss mark + ZWJ +
// woman) and the rules described in Annex #29 must be applied to group those
// code points into clusters perceived by the user as one character.
type Graphemes struct {
	// The code points over which this class iterates.
	codePoints []rune

	// The (byte-based) indices of the code points into the original string plus
	// len(original string). Thus, len(indices) = len(codePoints) + 1.
	indices []int

	// The current grapheme cluster to be returned. These are indices into
	// codePoints/indices. If start == end, we either haven't started iterating
	// yet (0) or the iteration has already completed (1).
	start, end int

	// The index of the next code point to be parsed.
	pos int

	// The current state of the code point parser.
	state int
}

// NewGraphemes returns a new grapheme cluster iterator.
func NewGraphemes(s string) *Graphemes {
	l := utf8.RuneCountInString(s)
	codePoints := make([]rune, l)
	indices := make([]int, l+1)
	i := 0
	for pos, r := range s {
		codePoints[i] = r
		indices[i] = pos
		i++
	}
	indices[l] = len(s)
	g := &Graphemes{
		codePoints: codePoints,
		indices:    indices,
	}
	g.Next() // Parse ahead.
	return g
}

// Next advances the iterator by one grapheme cluster and returns false if no
// clusters are left. This function must be called before the first cluster is
// accessed.
func (g *Graphemes) Next() bool {
	g.start = g.end

	// The state transition gives us a boundary instruction BEFORE the next code
	// point so we always need to stay ahead by one code point.

	// Parse the next code point.
	for g.pos <= len(g.codePoints) {
		// GB2.
		if g.pos == len(g.codePoints) {
			g.end = g.pos
			g.pos++
			break
		}

		// Determine the property of the next character.
		nextProperty := property(g.codePoints[g.pos])
		g.pos++

		// Find the applicable transition.
		var boundary bool
		transition, ok := grTransitions[[2]int{g.state, nextProperty}]
		if ok {
			// We have a specific transition. We'll use it.
			g.state = transition[0]
			boundary = transition[1] == grBoundary
		} else {
			// No specific transition found. Try the less specific ones.
			transAnyProp, okAnyProp := grTransitions[[2]int{g.state, prAny}]
			transAnyState, okAnyState := grTransitions[[2]int{grAny, nextProperty}]
			if okAnyProp && okAnyState {
				// Both apply. We'll use a mix (see comments for grTransitions).
				g.state = transAnyState[0]
				boundary = transAnyState[1] == grBoundary
				if transAnyProp[2] < transAnyState[2] {
					g.state = transAnyProp[0]
					boundary = transAnyProp[1] == grBoundary
				}
			} else if okAnyProp {
				// We only have a specific state.
				g.state = transAnyProp[0]
				boundary = transAnyProp[1] == grBoundary
				// This branch will probably never be reached because okAnyState will
				// always be true given the current transition map. But we keep it here
				// for future modifications to the transition map where this may not be
				// true anymore.
			} else if okAnyState {
				// We only have a specific property.
				g.state = transAnyState[0]
				boundary = transAnyState[1] == grBoundary
			} else {
				// No known transition. GB999: Any x Any.
				g.state = grAny
				boundary = true
			}
		}

		// If we found a cluster boundary, let's stop here. The current cluster will
		// be the one that just ended.
		if g.pos-1 == 0 /* GB1 */ || boundary {
			g.end = g.pos - 1
			break
		}
	}

	return g.start != g.end
}

// Runes returns a slice of runes (code points) which corresponds to the current
// grapheme cluster. If the iterator is already past the end or Next() has not
// yet been called, nil is returned.
func (g *Graphemes) Runes() []rune {
	if g.start == g.end {
		return nil
	}
	return g.codePoints[g.start:g.end]
}

// Str returns a substring of the original string which corresponds to the
// current grapheme cluster. If the iterator is already past the end or Next()
// has not yet been called, an empty string is returned.
func (g *Graphemes) Str() string {
	if g.start == g.end {
		return ""
	}
	return string(g.codePoints[g.start:g.end])
}

// Bytes returns a byte slice which corresponds to the current grapheme cluster.
// If the iterator is already past the end or Next() has not yet been called,
// nil is returned.
func (g *Graphemes) Bytes() []byte {
	if g.start == g.end {
		return nil
	}
	return []byte(string(g.codePoints[g.start:g.end]))
}

// Positions returns the interval of the current grapheme cluster as byte
// positions into the original string. The first returned value "from" indexes
// the first byte and the second returned value "to" indexes the first byte that
// is not included anymore, i.e. str[from:to] is the current grapheme cluster of
// the original string "str". If Next() has not yet been called, both values are
// 0. If the iterator is already past the end, both values are 1.
func (g *Graphemes) Positions() (int, int) {
	return g.indices[g.start], g.indices[g.end]
}

// Reset puts the iterator into its initial state such that the next call to
// Next() sets it to the first grapheme cluster again.
func (g *Graphemes) Reset() {
	g.start, g.end, g.pos, g.state = 0, 0, 0, grAny
	g.Next() // Parse ahead again.
}

// GraphemeClusterCount returns the number of user-perceived characters
// (grapheme clusters) for the given string. To calculate this number, it
// iterates through the string using the Graphemes iterator.
func GraphemeClusterCount(s string) (n int) {
	g := NewGraphemes(s)
	for g.Next() {
		n++
	}
	return
}
