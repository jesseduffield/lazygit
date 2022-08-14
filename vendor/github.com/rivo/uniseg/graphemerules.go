package uniseg

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
//   5. If both (2) and (3) were found, use state from (3) and breaking instruction
//      from the transition with the lower rule number, prefer (3) if rule numbers
//      are equal. Stop.
//   6. Assume grAny and grBoundary.
//
// Unicode version 14.0.0.
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
	{grAny, prPrepend}: {grPrepend, grBoundary, 9990},
	{grPrepend, prAny}: {grAny, grNoBoundary, 92},

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

// transitionGraphemeState determines the new state of the grapheme cluster
// parser given the current state and the next code point. It also returns
// whether a cluster boundary was detected.
func transitionGraphemeState(state int, r rune) (newState int, boundary bool) {
	// Determine the property of the next character.
	nextProperty := property(graphemeCodePoints, r)

	// Find the applicable transition.
	transition, ok := grTransitions[[2]int{state, nextProperty}]
	if ok {
		// We have a specific transition. We'll use it.
		return transition[0], transition[1] == grBoundary
	}

	// No specific transition found. Try the less specific ones.
	transAnyProp, okAnyProp := grTransitions[[2]int{state, prAny}]
	transAnyState, okAnyState := grTransitions[[2]int{grAny, nextProperty}]
	if okAnyProp && okAnyState {
		// Both apply. We'll use a mix (see comments for grTransitions).
		newState = transAnyState[0]
		boundary = transAnyState[1] == grBoundary
		if transAnyProp[2] < transAnyState[2] {
			boundary = transAnyProp[1] == grBoundary
		}
		return
	}

	if okAnyProp {
		// We only have a specific state.
		return transAnyProp[0], transAnyProp[1] == grBoundary
		// This branch will probably never be reached because okAnyState will
		// always be true given the current transition map. But we keep it here
		// for future modifications to the transition map where this may not be
		// true anymore.
	}

	if okAnyState {
		// We only have a specific property.
		return transAnyState[0], transAnyState[1] == grBoundary
	}

	// No known transition. GB999: Any ÷ Any.
	return grAny, true
}
