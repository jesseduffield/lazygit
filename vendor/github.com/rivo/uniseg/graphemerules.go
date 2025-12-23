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

// grTransitions implements the grapheme cluster parser's state transitions.
// Maps state and property to a new state, a breaking instruction, and rule
// number. The breaking instruction always refers to the boundary between the
// last and next code point. Returns negative values if no transition is found.
//
// This function is used as follows:
//
//  1. Find specific state + specific property. Stop if found.
//  2. Find specific state + any property.
//  3. Find any state + specific property.
//  4. If only (2) or (3) (but not both) was found, stop.
//  5. If both (2) and (3) were found, use state from (3) and breaking instruction
//     from the transition with the lower rule number, prefer (3) if rule numbers
//     are equal. Stop.
//  6. Assume grAny and grBoundary.
//
// Unicode version 15.0.0.
func grTransitions(state, prop int) (newState int, newProp int, boundary int) {
	// It turns out that using a big switch statement is much faster than using
	// a map.

	switch uint64(state) | uint64(prop)<<32 {
	// GB5
	case grAny | prCR<<32:
		return grCR, grBoundary, 50
	case grAny | prLF<<32:
		return grControlLF, grBoundary, 50
	case grAny | prControl<<32:
		return grControlLF, grBoundary, 50

	// GB4
	case grCR | prAny<<32:
		return grAny, grBoundary, 40
	case grControlLF | prAny<<32:
		return grAny, grBoundary, 40

	// GB3
	case grCR | prLF<<32:
		return grControlLF, grNoBoundary, 30

	// GB6
	case grAny | prL<<32:
		return grL, grBoundary, 9990
	case grL | prL<<32:
		return grL, grNoBoundary, 60
	case grL | prV<<32:
		return grLVV, grNoBoundary, 60
	case grL | prLV<<32:
		return grLVV, grNoBoundary, 60
	case grL | prLVT<<32:
		return grLVTT, grNoBoundary, 60

	// GB7
	case grAny | prLV<<32:
		return grLVV, grBoundary, 9990
	case grAny | prV<<32:
		return grLVV, grBoundary, 9990
	case grLVV | prV<<32:
		return grLVV, grNoBoundary, 70
	case grLVV | prT<<32:
		return grLVTT, grNoBoundary, 70

	// GB8
	case grAny | prLVT<<32:
		return grLVTT, grBoundary, 9990
	case grAny | prT<<32:
		return grLVTT, grBoundary, 9990
	case grLVTT | prT<<32:
		return grLVTT, grNoBoundary, 80

	// GB9
	case grAny | prExtend<<32:
		return grAny, grNoBoundary, 90
	case grAny | prZWJ<<32:
		return grAny, grNoBoundary, 90

	// GB9a
	case grAny | prSpacingMark<<32:
		return grAny, grNoBoundary, 91

	// GB9b
	case grAny | prPrepend<<32:
		return grPrepend, grBoundary, 9990
	case grPrepend | prAny<<32:
		return grAny, grNoBoundary, 92

	// GB11
	case grAny | prExtendedPictographic<<32:
		return grExtendedPictographic, grBoundary, 9990
	case grExtendedPictographic | prExtend<<32:
		return grExtendedPictographic, grNoBoundary, 110
	case grExtendedPictographic | prZWJ<<32:
		return grExtendedPictographicZWJ, grNoBoundary, 110
	case grExtendedPictographicZWJ | prExtendedPictographic<<32:
		return grExtendedPictographic, grNoBoundary, 110

	// GB12 / GB13
	case grAny | prRegionalIndicator<<32:
		return grRIOdd, grBoundary, 9990
	case grRIOdd | prRegionalIndicator<<32:
		return grRIEven, grNoBoundary, 120
	case grRIEven | prRegionalIndicator<<32:
		return grRIOdd, grBoundary, 120
	default:
		return -1, -1, -1
	}
}

// transitionGraphemeState determines the new state of the grapheme cluster
// parser given the current state and the next code point. It also returns the
// code point's grapheme property (the value mapped by the [graphemeCodePoints]
// table) and whether a cluster boundary was detected.
func transitionGraphemeState(state int, r rune) (newState, prop int, boundary bool) {
	// Determine the property of the next character.
	prop = propertyGraphemes(r)

	// Find the applicable transition.
	nextState, nextProp, _ := grTransitions(state, prop)
	if nextState >= 0 {
		// We have a specific transition. We'll use it.
		return nextState, prop, nextProp == grBoundary
	}

	// No specific transition found. Try the less specific ones.
	anyPropState, anyPropProp, anyPropRule := grTransitions(state, prAny)
	anyStateState, anyStateProp, anyStateRule := grTransitions(grAny, prop)
	if anyPropState >= 0 && anyStateState >= 0 {
		// Both apply. We'll use a mix (see comments for grTransitions).
		newState = anyStateState
		boundary = anyStateProp == grBoundary
		if anyPropRule < anyStateRule {
			boundary = anyPropProp == grBoundary
		}
		return
	}

	if anyPropState >= 0 {
		// We only have a specific state.
		return anyPropState, prop, anyPropProp == grBoundary
		// This branch will probably never be reached because okAnyState will
		// always be true given the current transition map. But we keep it here
		// for future modifications to the transition map where this may not be
		// true anymore.
	}

	if anyStateState >= 0 {
		// We only have a specific property.
		return anyStateState, prop, anyStateProp == grBoundary
	}

	// No known transition. GB999: Any ÷ Any.
	return grAny, prop, true
}
