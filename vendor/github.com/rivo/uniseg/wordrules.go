package uniseg

import "unicode/utf8"

// The states of the word break parser.
const (
	wbAny = iota
	wbCR
	wbLF
	wbNewline
	wbWSegSpace
	wbHebrewLetter
	wbALetter
	wbWB7
	wbWB7c
	wbNumeric
	wbWB11
	wbKatakana
	wbExtendNumLet
	wbOddRI
	wbEvenRI
	wbZWJBit = 16 // This bit is set for any states followed by at least one zero-width joiner (see WB4 and WB3c).
)

// The word break parser's breaking instructions.
const (
	wbDontBreak = iota
	wbBreak
)

// The word break parser's state transitions. It's anologous to grTransitions,
// see comments there for details. Unicode version 14.0.0.
var wbTransitions = map[[2]int][3]int{
	// WB3b.
	{wbAny, prNewline}: {wbNewline, wbBreak, 32},
	{wbAny, prCR}:      {wbCR, wbBreak, 32},
	{wbAny, prLF}:      {wbLF, wbBreak, 32},

	// WB3a.
	{wbNewline, prAny}: {wbAny, wbBreak, 31},
	{wbCR, prAny}:      {wbAny, wbBreak, 31},
	{wbLF, prAny}:      {wbAny, wbBreak, 31},

	// WB3.
	{wbCR, prLF}: {wbLF, wbDontBreak, 30},

	// WB3d.
	{wbAny, prWSegSpace}:       {wbWSegSpace, wbBreak, 9990},
	{wbWSegSpace, prWSegSpace}: {wbWSegSpace, wbDontBreak, 34},

	// WB5.
	{wbAny, prALetter}:               {wbALetter, wbBreak, 9990},
	{wbAny, prHebrewLetter}:          {wbHebrewLetter, wbBreak, 9990},
	{wbALetter, prALetter}:           {wbALetter, wbDontBreak, 50},
	{wbALetter, prHebrewLetter}:      {wbHebrewLetter, wbDontBreak, 50},
	{wbHebrewLetter, prALetter}:      {wbALetter, wbDontBreak, 50},
	{wbHebrewLetter, prHebrewLetter}: {wbHebrewLetter, wbDontBreak, 50},

	// WB7. Transitions to wbWB7 handled by transitionWordBreakState().
	{wbWB7, prALetter}:      {wbALetter, wbDontBreak, 70},
	{wbWB7, prHebrewLetter}: {wbHebrewLetter, wbDontBreak, 70},

	// WB7a.
	{wbHebrewLetter, prSingleQuote}: {wbAny, wbDontBreak, 71},

	// WB7c. Transitions to wbWB7c handled by transitionWordBreakState().
	{wbWB7c, prHebrewLetter}: {wbHebrewLetter, wbDontBreak, 73},

	// WB8.
	{wbAny, prNumeric}:     {wbNumeric, wbBreak, 9990},
	{wbNumeric, prNumeric}: {wbNumeric, wbDontBreak, 80},

	// WB9.
	{wbALetter, prNumeric}:      {wbNumeric, wbDontBreak, 90},
	{wbHebrewLetter, prNumeric}: {wbNumeric, wbDontBreak, 90},

	// WB10.
	{wbNumeric, prALetter}:      {wbALetter, wbDontBreak, 100},
	{wbNumeric, prHebrewLetter}: {wbHebrewLetter, wbDontBreak, 100},

	// WB11. Transitions to wbWB11 handled by transitionWordBreakState().
	{wbWB11, prNumeric}: {wbNumeric, wbDontBreak, 110},

	// WB13.
	{wbAny, prKatakana}:      {wbKatakana, wbBreak, 9990},
	{wbKatakana, prKatakana}: {wbKatakana, wbDontBreak, 130},

	// WB13a.
	{wbAny, prExtendNumLet}:          {wbExtendNumLet, wbBreak, 9990},
	{wbALetter, prExtendNumLet}:      {wbExtendNumLet, wbDontBreak, 131},
	{wbHebrewLetter, prExtendNumLet}: {wbExtendNumLet, wbDontBreak, 131},
	{wbNumeric, prExtendNumLet}:      {wbExtendNumLet, wbDontBreak, 131},
	{wbKatakana, prExtendNumLet}:     {wbExtendNumLet, wbDontBreak, 131},
	{wbExtendNumLet, prExtendNumLet}: {wbExtendNumLet, wbDontBreak, 131},

	// WB13b.
	{wbExtendNumLet, prALetter}:      {wbALetter, wbDontBreak, 132},
	{wbExtendNumLet, prHebrewLetter}: {wbHebrewLetter, wbDontBreak, 132},
	{wbExtendNumLet, prNumeric}:      {wbNumeric, wbDontBreak, 132},
	{wbExtendNumLet, prKatakana}:     {prKatakana, wbDontBreak, 132},
}

// transitionWordBreakState determines the new state of the word break parser
// given the current state and the next code point. It also returns whether a
// word boundary was detected. If more than one code point is needed to
// determine the new state, the byte slice or the string starting after rune "r"
// can be used (whichever is not nil or empty) for further lookups.
func transitionWordBreakState(state int, r rune, b []byte, str string) (newState int, wordBreak bool) {
	// Determine the property of the next character.
	nextProperty := property(workBreakCodePoints, r)

	// "Replacing Ignore Rules".
	if nextProperty == prZWJ {
		// WB4 (for zero-width joiners).
		if state == wbNewline || state == wbCR || state == wbLF {
			return wbAny | wbZWJBit, true // Make sure we don't apply WB4 to WB3a.
		}
		if state < 0 {
			return wbAny | wbZWJBit, false
		}
		return state | wbZWJBit, false
	} else if nextProperty == prExtend || nextProperty == prFormat {
		// WB4 (for Extend and Format).
		if state == wbNewline || state == wbCR || state == wbLF {
			return wbAny, true // Make sure we don't apply WB4 to WB3a.
		}
		if state == wbWSegSpace || state == wbAny|wbZWJBit {
			return wbAny, false // We don't break but this is also not WB3d or WB3c.
		}
		if state < 0 {
			return wbAny, false
		}
		return state, false
	} else if nextProperty == prExtendedPictographic && state >= 0 && state&wbZWJBit != 0 {
		// WB3c.
		return wbAny, false
	}
	if state >= 0 {
		state = state &^ wbZWJBit
	}

	// Find the applicable transition in the table.
	var rule int
	transition, ok := wbTransitions[[2]int{state, nextProperty}]
	if ok {
		// We have a specific transition. We'll use it.
		newState, wordBreak, rule = transition[0], transition[1] == wbBreak, transition[2]
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := wbTransitions[[2]int{state, prAny}]
		transAnyState, okAnyState := wbTransitions[[2]int{wbAny, nextProperty}]
		if okAnyProp && okAnyState {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, wordBreak, rule = transAnyState[0], transAnyState[1] == wbBreak, transAnyState[2]
			if transAnyProp[2] < transAnyState[2] {
				wordBreak, rule = transAnyProp[1] == wbBreak, transAnyProp[2]
			}
		} else if okAnyProp {
			// We only have a specific state.
			newState, wordBreak, rule = transAnyProp[0], transAnyProp[1] == wbBreak, transAnyProp[2]
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if okAnyState {
			// We only have a specific property.
			newState, wordBreak, rule = transAnyState[0], transAnyState[1] == wbBreak, transAnyState[2]
		} else {
			// No known transition. WB999: Any ÷ Any.
			newState, wordBreak, rule = wbAny, true, 9990
		}
	}

	// For those rules that need to look up runes further in the string, we
	// determine the property after nextProperty, skipping over Format, Extend,
	// and ZWJ (according to WB4). It's -1 if not needed, if such a rune cannot
	// be determined (because the text ends or the rune is faulty).
	farProperty := -1
	if rule > 60 &&
		(state == wbALetter || state == wbHebrewLetter || state == wbNumeric) &&
		(nextProperty == prMidLetter || nextProperty == prMidNumLet || nextProperty == prSingleQuote || // WB6.
			nextProperty == prDoubleQuote || // WB7b.
			nextProperty == prMidNum) { // WB12.
		for {
			var (
				r      rune
				length int
			)
			if b != nil { // Byte slice version.
				r, length = utf8.DecodeRune(b)
				b = b[length:]
			} else { // String version.
				r, length = utf8.DecodeRuneInString(str)
				str = str[length:]
			}
			if r == utf8.RuneError {
				break
			}
			prop := property(workBreakCodePoints, r)
			if prop == prExtend || prop == prFormat || prop == prZWJ {
				continue
			}
			farProperty = prop
			break
		}
	}

	// WB6.
	if rule > 60 &&
		(state == wbALetter || state == wbHebrewLetter) &&
		(nextProperty == prMidLetter || nextProperty == prMidNumLet || nextProperty == prSingleQuote) &&
		(farProperty == prALetter || farProperty == prHebrewLetter) {
		return wbWB7, false
	}

	// WB7b.
	if rule > 72 &&
		state == wbHebrewLetter &&
		nextProperty == prDoubleQuote &&
		farProperty == prHebrewLetter {
		return wbWB7c, false
	}

	// WB12.
	if rule > 120 &&
		state == wbNumeric &&
		(nextProperty == prMidNum || nextProperty == prMidNumLet || nextProperty == prSingleQuote) &&
		farProperty == prNumeric {
		return wbWB11, false
	}

	// WB15 and WB16.
	if newState == wbAny && nextProperty == prRegionalIndicator {
		if state != wbOddRI && state != wbEvenRI { // Includes state == -1.
			// Transition into the first RI.
			return wbOddRI, true
		}
		if state == wbOddRI {
			// Don't break pairs of Regional Indicators.
			return wbEvenRI, false
		}
		return wbOddRI, true // We can break after a pair.
	}

	return
}
