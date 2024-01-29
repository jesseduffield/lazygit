package uniseg

import "unicode/utf8"

// The states of the sentence break parser.
const (
	sbAny = iota
	sbCR
	sbParaSep
	sbATerm
	sbUpper
	sbLower
	sbSB7
	sbSB8Close
	sbSB8Sp
	sbSTerm
	sbSB8aClose
	sbSB8aSp
)

// sbTransitions implements the sentence break parser's state transitions. It's
// anologous to [grTransitions], see comments there for details.
//
// Unicode version 15.0.0.
func sbTransitions(state, prop int) (newState int, sentenceBreak bool, rule int) {
	switch uint64(state) | uint64(prop)<<32 {
	// SB3.
	case sbAny | prCR<<32:
		return sbCR, false, 9990
	case sbCR | prLF<<32:
		return sbParaSep, false, 30

	// SB4.
	case sbAny | prSep<<32:
		return sbParaSep, false, 9990
	case sbAny | prLF<<32:
		return sbParaSep, false, 9990
	case sbParaSep | prAny<<32:
		return sbAny, true, 40
	case sbCR | prAny<<32:
		return sbAny, true, 40

	// SB6.
	case sbAny | prATerm<<32:
		return sbATerm, false, 9990
	case sbATerm | prNumeric<<32:
		return sbAny, false, 60
	case sbSB7 | prNumeric<<32:
		return sbAny, false, 60 // Because ATerm also appears in SB7.

	// SB7.
	case sbAny | prUpper<<32:
		return sbUpper, false, 9990
	case sbAny | prLower<<32:
		return sbLower, false, 9990
	case sbUpper | prATerm<<32:
		return sbSB7, false, 70
	case sbLower | prATerm<<32:
		return sbSB7, false, 70
	case sbSB7 | prUpper<<32:
		return sbUpper, false, 70

	// SB8a.
	case sbAny | prSTerm<<32:
		return sbSTerm, false, 9990
	case sbATerm | prSContinue<<32:
		return sbAny, false, 81
	case sbATerm | prATerm<<32:
		return sbATerm, false, 81
	case sbATerm | prSTerm<<32:
		return sbSTerm, false, 81
	case sbSB7 | prSContinue<<32:
		return sbAny, false, 81
	case sbSB7 | prATerm<<32:
		return sbATerm, false, 81
	case sbSB7 | prSTerm<<32:
		return sbSTerm, false, 81
	case sbSB8Close | prSContinue<<32:
		return sbAny, false, 81
	case sbSB8Close | prATerm<<32:
		return sbATerm, false, 81
	case sbSB8Close | prSTerm<<32:
		return sbSTerm, false, 81
	case sbSB8Sp | prSContinue<<32:
		return sbAny, false, 81
	case sbSB8Sp | prATerm<<32:
		return sbATerm, false, 81
	case sbSB8Sp | prSTerm<<32:
		return sbSTerm, false, 81
	case sbSTerm | prSContinue<<32:
		return sbAny, false, 81
	case sbSTerm | prATerm<<32:
		return sbATerm, false, 81
	case sbSTerm | prSTerm<<32:
		return sbSTerm, false, 81
	case sbSB8aClose | prSContinue<<32:
		return sbAny, false, 81
	case sbSB8aClose | prATerm<<32:
		return sbATerm, false, 81
	case sbSB8aClose | prSTerm<<32:
		return sbSTerm, false, 81
	case sbSB8aSp | prSContinue<<32:
		return sbAny, false, 81
	case sbSB8aSp | prATerm<<32:
		return sbATerm, false, 81
	case sbSB8aSp | prSTerm<<32:
		return sbSTerm, false, 81

	// SB9.
	case sbATerm | prClose<<32:
		return sbSB8Close, false, 90
	case sbSB7 | prClose<<32:
		return sbSB8Close, false, 90
	case sbSB8Close | prClose<<32:
		return sbSB8Close, false, 90
	case sbATerm | prSp<<32:
		return sbSB8Sp, false, 90
	case sbSB7 | prSp<<32:
		return sbSB8Sp, false, 90
	case sbSB8Close | prSp<<32:
		return sbSB8Sp, false, 90
	case sbSTerm | prClose<<32:
		return sbSB8aClose, false, 90
	case sbSB8aClose | prClose<<32:
		return sbSB8aClose, false, 90
	case sbSTerm | prSp<<32:
		return sbSB8aSp, false, 90
	case sbSB8aClose | prSp<<32:
		return sbSB8aSp, false, 90
	case sbATerm | prSep<<32:
		return sbParaSep, false, 90
	case sbATerm | prCR<<32:
		return sbParaSep, false, 90
	case sbATerm | prLF<<32:
		return sbParaSep, false, 90
	case sbSB7 | prSep<<32:
		return sbParaSep, false, 90
	case sbSB7 | prCR<<32:
		return sbParaSep, false, 90
	case sbSB7 | prLF<<32:
		return sbParaSep, false, 90
	case sbSB8Close | prSep<<32:
		return sbParaSep, false, 90
	case sbSB8Close | prCR<<32:
		return sbParaSep, false, 90
	case sbSB8Close | prLF<<32:
		return sbParaSep, false, 90
	case sbSTerm | prSep<<32:
		return sbParaSep, false, 90
	case sbSTerm | prCR<<32:
		return sbParaSep, false, 90
	case sbSTerm | prLF<<32:
		return sbParaSep, false, 90
	case sbSB8aClose | prSep<<32:
		return sbParaSep, false, 90
	case sbSB8aClose | prCR<<32:
		return sbParaSep, false, 90
	case sbSB8aClose | prLF<<32:
		return sbParaSep, false, 90

	// SB10.
	case sbSB8Sp | prSp<<32:
		return sbSB8Sp, false, 100
	case sbSB8aSp | prSp<<32:
		return sbSB8aSp, false, 100
	case sbSB8Sp | prSep<<32:
		return sbParaSep, false, 100
	case sbSB8Sp | prCR<<32:
		return sbParaSep, false, 100
	case sbSB8Sp | prLF<<32:
		return sbParaSep, false, 100

	// SB11.
	case sbATerm | prAny<<32:
		return sbAny, true, 110
	case sbSB7 | prAny<<32:
		return sbAny, true, 110
	case sbSB8Close | prAny<<32:
		return sbAny, true, 110
	case sbSB8Sp | prAny<<32:
		return sbAny, true, 110
	case sbSTerm | prAny<<32:
		return sbAny, true, 110
	case sbSB8aClose | prAny<<32:
		return sbAny, true, 110
	case sbSB8aSp | prAny<<32:
		return sbAny, true, 110
	// We'll always break after ParaSep due to SB4.

	default:
		return -1, false, -1
	}
}

// transitionSentenceBreakState determines the new state of the sentence break
// parser given the current state and the next code point. It also returns
// whether a sentence boundary was detected. If more than one code point is
// needed to determine the new state, the byte slice or the string starting
// after rune "r" can be used (whichever is not nil or empty) for further
// lookups.
func transitionSentenceBreakState(state int, r rune, b []byte, str string) (newState int, sentenceBreak bool) {
	// Determine the property of the next character.
	nextProperty := property(sentenceBreakCodePoints, r)

	// SB5 (Replacing Ignore Rules).
	if nextProperty == prExtend || nextProperty == prFormat {
		if state == sbParaSep || state == sbCR {
			return sbAny, true // Make sure we don't apply SB5 to SB3 or SB4.
		}
		if state < 0 {
			return sbAny, true // SB1.
		}
		return state, false
	}

	// Find the applicable transition in the table.
	var rule int
	newState, sentenceBreak, rule = sbTransitions(state, nextProperty)
	if newState < 0 {
		// No specific transition found. Try the less specific ones.
		anyPropState, anyPropProp, anyPropRule := sbTransitions(state, prAny)
		anyStateState, anyStateProp, anyStateRule := sbTransitions(sbAny, nextProperty)
		if anyPropState >= 0 && anyStateState >= 0 {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, sentenceBreak, rule = anyStateState, anyStateProp, anyStateRule
			if anyPropRule < anyStateRule {
				sentenceBreak, rule = anyPropProp, anyPropRule
			}
		} else if anyPropState >= 0 {
			// We only have a specific state.
			newState, sentenceBreak, rule = anyPropState, anyPropProp, anyPropRule
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if anyStateState >= 0 {
			// We only have a specific property.
			newState, sentenceBreak, rule = anyStateState, anyStateProp, anyStateRule
		} else {
			// No known transition. SB999: Any × Any.
			newState, sentenceBreak, rule = sbAny, false, 9990
		}
	}

	// SB8.
	if rule > 80 && (state == sbATerm || state == sbSB8Close || state == sbSB8Sp || state == sbSB7) {
		// Check the right side of the rule.
		var length int
		for nextProperty != prOLetter &&
			nextProperty != prUpper &&
			nextProperty != prLower &&
			nextProperty != prSep &&
			nextProperty != prCR &&
			nextProperty != prLF &&
			nextProperty != prATerm &&
			nextProperty != prSTerm {
			// Move on to the next rune.
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
			nextProperty = property(sentenceBreakCodePoints, r)
		}
		if nextProperty == prLower {
			return sbLower, false
		}
	}

	return
}
