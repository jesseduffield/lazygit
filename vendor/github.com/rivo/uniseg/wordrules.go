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

// wbTransitions implements the word break parser's state transitions. It's
// anologous to [grTransitions], see comments there for details.
//
// Unicode version 15.0.0.
func wbTransitions(state, prop int) (newState int, wordBreak bool, rule int) {
	switch uint64(state) | uint64(prop)<<32 {
	// WB3b.
	case wbAny | prNewline<<32:
		return wbNewline, true, 32
	case wbAny | prCR<<32:
		return wbCR, true, 32
	case wbAny | prLF<<32:
		return wbLF, true, 32

	// WB3a.
	case wbNewline | prAny<<32:
		return wbAny, true, 31
	case wbCR | prAny<<32:
		return wbAny, true, 31
	case wbLF | prAny<<32:
		return wbAny, true, 31

	// WB3.
	case wbCR | prLF<<32:
		return wbLF, false, 30

	// WB3d.
	case wbAny | prWSegSpace<<32:
		return wbWSegSpace, true, 9990
	case wbWSegSpace | prWSegSpace<<32:
		return wbWSegSpace, false, 34

	// WB5.
	case wbAny | prALetter<<32:
		return wbALetter, true, 9990
	case wbAny | prHebrewLetter<<32:
		return wbHebrewLetter, true, 9990
	case wbALetter | prALetter<<32:
		return wbALetter, false, 50
	case wbALetter | prHebrewLetter<<32:
		return wbHebrewLetter, false, 50
	case wbHebrewLetter | prALetter<<32:
		return wbALetter, false, 50
	case wbHebrewLetter | prHebrewLetter<<32:
		return wbHebrewLetter, false, 50

	// WB7. Transitions to wbWB7 handled by transitionWordBreakState().
	case wbWB7 | prALetter<<32:
		return wbALetter, false, 70
	case wbWB7 | prHebrewLetter<<32:
		return wbHebrewLetter, false, 70

	// WB7a.
	case wbHebrewLetter | prSingleQuote<<32:
		return wbAny, false, 71

	// WB7c. Transitions to wbWB7c handled by transitionWordBreakState().
	case wbWB7c | prHebrewLetter<<32:
		return wbHebrewLetter, false, 73

	// WB8.
	case wbAny | prNumeric<<32:
		return wbNumeric, true, 9990
	case wbNumeric | prNumeric<<32:
		return wbNumeric, false, 80

	// WB9.
	case wbALetter | prNumeric<<32:
		return wbNumeric, false, 90
	case wbHebrewLetter | prNumeric<<32:
		return wbNumeric, false, 90

	// WB10.
	case wbNumeric | prALetter<<32:
		return wbALetter, false, 100
	case wbNumeric | prHebrewLetter<<32:
		return wbHebrewLetter, false, 100

	// WB11. Transitions to wbWB11 handled by transitionWordBreakState().
	case wbWB11 | prNumeric<<32:
		return wbNumeric, false, 110

	// WB13.
	case wbAny | prKatakana<<32:
		return wbKatakana, true, 9990
	case wbKatakana | prKatakana<<32:
		return wbKatakana, false, 130

	// WB13a.
	case wbAny | prExtendNumLet<<32:
		return wbExtendNumLet, true, 9990
	case wbALetter | prExtendNumLet<<32:
		return wbExtendNumLet, false, 131
	case wbHebrewLetter | prExtendNumLet<<32:
		return wbExtendNumLet, false, 131
	case wbNumeric | prExtendNumLet<<32:
		return wbExtendNumLet, false, 131
	case wbKatakana | prExtendNumLet<<32:
		return wbExtendNumLet, false, 131
	case wbExtendNumLet | prExtendNumLet<<32:
		return wbExtendNumLet, false, 131

	// WB13b.
	case wbExtendNumLet | prALetter<<32:
		return wbALetter, false, 132
	case wbExtendNumLet | prHebrewLetter<<32:
		return wbHebrewLetter, false, 132
	case wbExtendNumLet | prNumeric<<32:
		return wbNumeric, false, 132
	case wbExtendNumLet | prKatakana<<32:
		return wbKatakana, false, 132

	default:
		return -1, false, -1
	}
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
	newState, wordBreak, rule = wbTransitions(state, nextProperty)
	if newState < 0 {
		// No specific transition found. Try the less specific ones.
		anyPropState, anyPropWordBreak, anyPropRule := wbTransitions(state, prAny)
		anyStateState, anyStateWordBreak, anyStateRule := wbTransitions(wbAny, nextProperty)
		if anyPropState >= 0 && anyStateState >= 0 {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, wordBreak, rule = anyStateState, anyStateWordBreak, anyStateRule
			if anyPropRule < anyStateRule {
				wordBreak, rule = anyPropWordBreak, anyPropRule
			}
		} else if anyPropState >= 0 {
			// We only have a specific state.
			newState, wordBreak, rule = anyPropState, anyPropWordBreak, anyPropRule
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if anyStateState >= 0 {
			// We only have a specific property.
			newState, wordBreak, rule = anyStateState, anyStateWordBreak, anyStateRule
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
