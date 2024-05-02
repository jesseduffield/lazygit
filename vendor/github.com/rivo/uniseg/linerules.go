package uniseg

import "unicode/utf8"

// The states of the line break parser.
const (
	lbAny = iota
	lbBK
	lbCR
	lbLF
	lbNL
	lbSP
	lbZW
	lbWJ
	lbGL
	lbBA
	lbHY
	lbCL
	lbCP
	lbEX
	lbIS
	lbSY
	lbOP
	lbQU
	lbQUSP
	lbNS
	lbCLCPSP
	lbB2
	lbB2SP
	lbCB
	lbBB
	lbLB21a
	lbHL
	lbAL
	lbNU
	lbPR
	lbEB
	lbIDEM
	lbNUNU
	lbNUSY
	lbNUIS
	lbNUCL
	lbNUCP
	lbPO
	lbJL
	lbJV
	lbJT
	lbH2
	lbH3
	lbOddRI
	lbEvenRI
	lbExtPicCn
	lbZWJBit     = 64
	lbCPeaFWHBit = 128
)

// These constants define whether a given text may be broken into the next line.
// If the break is optional (LineCanBreak), you may choose to break or not based
// on your own criteria, for example, if the text has reached the available
// width.
const (
	LineDontBreak = iota // You may not break the line here.
	LineCanBreak         // You may or may not break the line here.
	LineMustBreak        // You must break the line here.
)

// lbTransitions implements the line break parser's state transitions. It's
// anologous to [grTransitions], see comments there for details.
//
// Unicode version 15.0.0.
func lbTransitions(state, prop int) (newState, lineBreak, rule int) {
	switch uint64(state) | uint64(prop)<<32 {
	// LB4.
	case lbBK | prAny<<32:
		return lbAny, LineMustBreak, 40

	// LB5.
	case lbCR | prLF<<32:
		return lbLF, LineDontBreak, 50
	case lbCR | prAny<<32:
		return lbAny, LineMustBreak, 50
	case lbLF | prAny<<32:
		return lbAny, LineMustBreak, 50
	case lbNL | prAny<<32:
		return lbAny, LineMustBreak, 50

	// LB6.
	case lbAny | prBK<<32:
		return lbBK, LineDontBreak, 60
	case lbAny | prCR<<32:
		return lbCR, LineDontBreak, 60
	case lbAny | prLF<<32:
		return lbLF, LineDontBreak, 60
	case lbAny | prNL<<32:
		return lbNL, LineDontBreak, 60

	// LB7.
	case lbAny | prSP<<32:
		return lbSP, LineDontBreak, 70
	case lbAny | prZW<<32:
		return lbZW, LineDontBreak, 70

	// LB8.
	case lbZW | prSP<<32:
		return lbZW, LineDontBreak, 70
	case lbZW | prAny<<32:
		return lbAny, LineCanBreak, 80

	// LB11.
	case lbAny | prWJ<<32:
		return lbWJ, LineDontBreak, 110
	case lbWJ | prAny<<32:
		return lbAny, LineDontBreak, 110

	// LB12.
	case lbAny | prGL<<32:
		return lbGL, LineCanBreak, 310
	case lbGL | prAny<<32:
		return lbAny, LineDontBreak, 120

	// LB13 (simple transitions).
	case lbAny | prCL<<32:
		return lbCL, LineCanBreak, 310
	case lbAny | prCP<<32:
		return lbCP, LineCanBreak, 310
	case lbAny | prEX<<32:
		return lbEX, LineDontBreak, 130
	case lbAny | prIS<<32:
		return lbIS, LineCanBreak, 310
	case lbAny | prSY<<32:
		return lbSY, LineCanBreak, 310

	// LB14.
	case lbAny | prOP<<32:
		return lbOP, LineCanBreak, 310
	case lbOP | prSP<<32:
		return lbOP, LineDontBreak, 70
	case lbOP | prAny<<32:
		return lbAny, LineDontBreak, 140

	// LB15.
	case lbQU | prSP<<32:
		return lbQUSP, LineDontBreak, 70
	case lbQU | prOP<<32:
		return lbOP, LineDontBreak, 150
	case lbQUSP | prOP<<32:
		return lbOP, LineDontBreak, 150

	// LB16.
	case lbCL | prSP<<32:
		return lbCLCPSP, LineDontBreak, 70
	case lbNUCL | prSP<<32:
		return lbCLCPSP, LineDontBreak, 70
	case lbCP | prSP<<32:
		return lbCLCPSP, LineDontBreak, 70
	case lbNUCP | prSP<<32:
		return lbCLCPSP, LineDontBreak, 70
	case lbCL | prNS<<32:
		return lbNS, LineDontBreak, 160
	case lbNUCL | prNS<<32:
		return lbNS, LineDontBreak, 160
	case lbCP | prNS<<32:
		return lbNS, LineDontBreak, 160
	case lbNUCP | prNS<<32:
		return lbNS, LineDontBreak, 160
	case lbCLCPSP | prNS<<32:
		return lbNS, LineDontBreak, 160

	// LB17.
	case lbAny | prB2<<32:
		return lbB2, LineCanBreak, 310
	case lbB2 | prSP<<32:
		return lbB2SP, LineDontBreak, 70
	case lbB2 | prB2<<32:
		return lbB2, LineDontBreak, 170
	case lbB2SP | prB2<<32:
		return lbB2, LineDontBreak, 170

	// LB18.
	case lbSP | prAny<<32:
		return lbAny, LineCanBreak, 180
	case lbQUSP | prAny<<32:
		return lbAny, LineCanBreak, 180
	case lbCLCPSP | prAny<<32:
		return lbAny, LineCanBreak, 180
	case lbB2SP | prAny<<32:
		return lbAny, LineCanBreak, 180

	// LB19.
	case lbAny | prQU<<32:
		return lbQU, LineDontBreak, 190
	case lbQU | prAny<<32:
		return lbAny, LineDontBreak, 190

	// LB20.
	case lbAny | prCB<<32:
		return lbCB, LineCanBreak, 200
	case lbCB | prAny<<32:
		return lbAny, LineCanBreak, 200

	// LB21.
	case lbAny | prBA<<32:
		return lbBA, LineDontBreak, 210
	case lbAny | prHY<<32:
		return lbHY, LineDontBreak, 210
	case lbAny | prNS<<32:
		return lbNS, LineDontBreak, 210
	case lbAny | prBB<<32:
		return lbBB, LineCanBreak, 310
	case lbBB | prAny<<32:
		return lbAny, LineDontBreak, 210

	// LB21a.
	case lbAny | prHL<<32:
		return lbHL, LineCanBreak, 310
	case lbHL | prHY<<32:
		return lbLB21a, LineDontBreak, 210
	case lbHL | prBA<<32:
		return lbLB21a, LineDontBreak, 210
	case lbLB21a | prAny<<32:
		return lbAny, LineDontBreak, 211

	// LB21b.
	case lbSY | prHL<<32:
		return lbHL, LineDontBreak, 212
	case lbNUSY | prHL<<32:
		return lbHL, LineDontBreak, 212

	// LB22.
	case lbAny | prIN<<32:
		return lbAny, LineDontBreak, 220

	// LB23.
	case lbAny | prAL<<32:
		return lbAL, LineCanBreak, 310
	case lbAny | prNU<<32:
		return lbNU, LineCanBreak, 310
	case lbAL | prNU<<32:
		return lbNU, LineDontBreak, 230
	case lbHL | prNU<<32:
		return lbNU, LineDontBreak, 230
	case lbNU | prAL<<32:
		return lbAL, LineDontBreak, 230
	case lbNU | prHL<<32:
		return lbHL, LineDontBreak, 230
	case lbNUNU | prAL<<32:
		return lbAL, LineDontBreak, 230
	case lbNUNU | prHL<<32:
		return lbHL, LineDontBreak, 230

	// LB23a.
	case lbAny | prPR<<32:
		return lbPR, LineCanBreak, 310
	case lbAny | prID<<32:
		return lbIDEM, LineCanBreak, 310
	case lbAny | prEB<<32:
		return lbEB, LineCanBreak, 310
	case lbAny | prEM<<32:
		return lbIDEM, LineCanBreak, 310
	case lbPR | prID<<32:
		return lbIDEM, LineDontBreak, 231
	case lbPR | prEB<<32:
		return lbEB, LineDontBreak, 231
	case lbPR | prEM<<32:
		return lbIDEM, LineDontBreak, 231
	case lbIDEM | prPO<<32:
		return lbPO, LineDontBreak, 231
	case lbEB | prPO<<32:
		return lbPO, LineDontBreak, 231

	// LB24.
	case lbAny | prPO<<32:
		return lbPO, LineCanBreak, 310
	case lbPR | prAL<<32:
		return lbAL, LineDontBreak, 240
	case lbPR | prHL<<32:
		return lbHL, LineDontBreak, 240
	case lbPO | prAL<<32:
		return lbAL, LineDontBreak, 240
	case lbPO | prHL<<32:
		return lbHL, LineDontBreak, 240
	case lbAL | prPR<<32:
		return lbPR, LineDontBreak, 240
	case lbAL | prPO<<32:
		return lbPO, LineDontBreak, 240
	case lbHL | prPR<<32:
		return lbPR, LineDontBreak, 240
	case lbHL | prPO<<32:
		return lbPO, LineDontBreak, 240

	// LB25 (simple transitions).
	case lbPR | prNU<<32:
		return lbNU, LineDontBreak, 250
	case lbPO | prNU<<32:
		return lbNU, LineDontBreak, 250
	case lbOP | prNU<<32:
		return lbNU, LineDontBreak, 250
	case lbHY | prNU<<32:
		return lbNU, LineDontBreak, 250
	case lbNU | prNU<<32:
		return lbNUNU, LineDontBreak, 250
	case lbNU | prSY<<32:
		return lbNUSY, LineDontBreak, 250
	case lbNU | prIS<<32:
		return lbNUIS, LineDontBreak, 250
	case lbNUNU | prNU<<32:
		return lbNUNU, LineDontBreak, 250
	case lbNUNU | prSY<<32:
		return lbNUSY, LineDontBreak, 250
	case lbNUNU | prIS<<32:
		return lbNUIS, LineDontBreak, 250
	case lbNUSY | prNU<<32:
		return lbNUNU, LineDontBreak, 250
	case lbNUSY | prSY<<32:
		return lbNUSY, LineDontBreak, 250
	case lbNUSY | prIS<<32:
		return lbNUIS, LineDontBreak, 250
	case lbNUIS | prNU<<32:
		return lbNUNU, LineDontBreak, 250
	case lbNUIS | prSY<<32:
		return lbNUSY, LineDontBreak, 250
	case lbNUIS | prIS<<32:
		return lbNUIS, LineDontBreak, 250
	case lbNU | prCL<<32:
		return lbNUCL, LineDontBreak, 250
	case lbNU | prCP<<32:
		return lbNUCP, LineDontBreak, 250
	case lbNUNU | prCL<<32:
		return lbNUCL, LineDontBreak, 250
	case lbNUNU | prCP<<32:
		return lbNUCP, LineDontBreak, 250
	case lbNUSY | prCL<<32:
		return lbNUCL, LineDontBreak, 250
	case lbNUSY | prCP<<32:
		return lbNUCP, LineDontBreak, 250
	case lbNUIS | prCL<<32:
		return lbNUCL, LineDontBreak, 250
	case lbNUIS | prCP<<32:
		return lbNUCP, LineDontBreak, 250
	case lbNU | prPO<<32:
		return lbPO, LineDontBreak, 250
	case lbNUNU | prPO<<32:
		return lbPO, LineDontBreak, 250
	case lbNUSY | prPO<<32:
		return lbPO, LineDontBreak, 250
	case lbNUIS | prPO<<32:
		return lbPO, LineDontBreak, 250
	case lbNUCL | prPO<<32:
		return lbPO, LineDontBreak, 250
	case lbNUCP | prPO<<32:
		return lbPO, LineDontBreak, 250
	case lbNU | prPR<<32:
		return lbPR, LineDontBreak, 250
	case lbNUNU | prPR<<32:
		return lbPR, LineDontBreak, 250
	case lbNUSY | prPR<<32:
		return lbPR, LineDontBreak, 250
	case lbNUIS | prPR<<32:
		return lbPR, LineDontBreak, 250
	case lbNUCL | prPR<<32:
		return lbPR, LineDontBreak, 250
	case lbNUCP | prPR<<32:
		return lbPR, LineDontBreak, 250

	// LB26.
	case lbAny | prJL<<32:
		return lbJL, LineCanBreak, 310
	case lbAny | prJV<<32:
		return lbJV, LineCanBreak, 310
	case lbAny | prJT<<32:
		return lbJT, LineCanBreak, 310
	case lbAny | prH2<<32:
		return lbH2, LineCanBreak, 310
	case lbAny | prH3<<32:
		return lbH3, LineCanBreak, 310
	case lbJL | prJL<<32:
		return lbJL, LineDontBreak, 260
	case lbJL | prJV<<32:
		return lbJV, LineDontBreak, 260
	case lbJL | prH2<<32:
		return lbH2, LineDontBreak, 260
	case lbJL | prH3<<32:
		return lbH3, LineDontBreak, 260
	case lbJV | prJV<<32:
		return lbJV, LineDontBreak, 260
	case lbJV | prJT<<32:
		return lbJT, LineDontBreak, 260
	case lbH2 | prJV<<32:
		return lbJV, LineDontBreak, 260
	case lbH2 | prJT<<32:
		return lbJT, LineDontBreak, 260
	case lbJT | prJT<<32:
		return lbJT, LineDontBreak, 260
	case lbH3 | prJT<<32:
		return lbJT, LineDontBreak, 260

	// LB27.
	case lbJL | prPO<<32:
		return lbPO, LineDontBreak, 270
	case lbJV | prPO<<32:
		return lbPO, LineDontBreak, 270
	case lbJT | prPO<<32:
		return lbPO, LineDontBreak, 270
	case lbH2 | prPO<<32:
		return lbPO, LineDontBreak, 270
	case lbH3 | prPO<<32:
		return lbPO, LineDontBreak, 270
	case lbPR | prJL<<32:
		return lbJL, LineDontBreak, 270
	case lbPR | prJV<<32:
		return lbJV, LineDontBreak, 270
	case lbPR | prJT<<32:
		return lbJT, LineDontBreak, 270
	case lbPR | prH2<<32:
		return lbH2, LineDontBreak, 270
	case lbPR | prH3<<32:
		return lbH3, LineDontBreak, 270

	// LB28.
	case lbAL | prAL<<32:
		return lbAL, LineDontBreak, 280
	case lbAL | prHL<<32:
		return lbHL, LineDontBreak, 280
	case lbHL | prAL<<32:
		return lbAL, LineDontBreak, 280
	case lbHL | prHL<<32:
		return lbHL, LineDontBreak, 280

	// LB29.
	case lbIS | prAL<<32:
		return lbAL, LineDontBreak, 290
	case lbIS | prHL<<32:
		return lbHL, LineDontBreak, 290
	case lbNUIS | prAL<<32:
		return lbAL, LineDontBreak, 290
	case lbNUIS | prHL<<32:
		return lbHL, LineDontBreak, 290

	default:
		return -1, -1, -1
	}
}

// transitionLineBreakState determines the new state of the line break parser
// given the current state and the next code point. It also returns the type of
// line break: LineDontBreak, LineCanBreak, or LineMustBreak. If more than one
// code point is needed to determine the new state, the byte slice or the string
// starting after rune "r" can be used (whichever is not nil or empty) for
// further lookups.
func transitionLineBreakState(state int, r rune, b []byte, str string) (newState int, lineBreak int) {
	// Determine the property of the next character.
	nextProperty, generalCategory := propertyLineBreak(r)

	// Prepare.
	var forceNoBreak, isCPeaFWH bool
	if state >= 0 && state&lbCPeaFWHBit != 0 {
		isCPeaFWH = true // LB30: CP but ea is not F, W, or H.
		state = state &^ lbCPeaFWHBit
	}
	if state >= 0 && state&lbZWJBit != 0 {
		state = state &^ lbZWJBit // Extract zero-width joiner bit.
		forceNoBreak = true       // LB8a.
	}

	defer func() {
		// Transition into LB30.
		if newState == lbCP || newState == lbNUCP {
			ea := propertyEastAsianWidth(r)
			if ea != prF && ea != prW && ea != prH {
				newState |= lbCPeaFWHBit
			}
		}

		// Override break.
		if forceNoBreak {
			lineBreak = LineDontBreak
		}
	}()

	// LB1.
	if nextProperty == prAI || nextProperty == prSG || nextProperty == prXX {
		nextProperty = prAL
	} else if nextProperty == prSA {
		if generalCategory == gcMn || generalCategory == gcMc {
			nextProperty = prCM
		} else {
			nextProperty = prAL
		}
	} else if nextProperty == prCJ {
		nextProperty = prNS
	}

	// Combining marks.
	if nextProperty == prZWJ || nextProperty == prCM {
		var bit int
		if nextProperty == prZWJ {
			bit = lbZWJBit
		}
		mustBreakState := state < 0 || state == lbBK || state == lbCR || state == lbLF || state == lbNL
		if !mustBreakState && state != lbSP && state != lbZW && state != lbQUSP && state != lbCLCPSP && state != lbB2SP {
			// LB9.
			return state | bit, LineDontBreak
		} else {
			// LB10.
			if mustBreakState {
				return lbAL | bit, LineMustBreak
			}
			return lbAL | bit, LineCanBreak
		}
	}

	// Find the applicable transition in the table.
	var rule int
	newState, lineBreak, rule = lbTransitions(state, nextProperty)
	if newState < 0 {
		// No specific transition found. Try the less specific ones.
		anyPropProp, anyPropLineBreak, anyPropRule := lbTransitions(state, prAny)
		anyStateProp, anyStateLineBreak, anyStateRule := lbTransitions(lbAny, nextProperty)
		if anyPropProp >= 0 && anyStateProp >= 0 {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, lineBreak, rule = anyStateProp, anyStateLineBreak, anyStateRule
			if anyPropRule < anyStateRule {
				lineBreak, rule = anyPropLineBreak, anyPropRule
			}
		} else if anyPropProp >= 0 {
			// We only have a specific state.
			newState, lineBreak, rule = anyPropProp, anyPropLineBreak, anyPropRule
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if anyStateProp >= 0 {
			// We only have a specific property.
			newState, lineBreak, rule = anyStateProp, anyStateLineBreak, anyStateRule
		} else {
			// No known transition. LB31: ALL ÷ ALL.
			newState, lineBreak, rule = lbAny, LineCanBreak, 310
		}
	}

	// LB12a.
	if rule > 121 &&
		nextProperty == prGL &&
		(state != lbSP && state != lbBA && state != lbHY && state != lbLB21a && state != lbQUSP && state != lbCLCPSP && state != lbB2SP) {
		return lbGL, LineDontBreak
	}

	// LB13.
	if rule > 130 && state != lbNU && state != lbNUNU {
		switch nextProperty {
		case prCL:
			return lbCL, LineDontBreak
		case prCP:
			return lbCP, LineDontBreak
		case prIS:
			return lbIS, LineDontBreak
		case prSY:
			return lbSY, LineDontBreak
		}
	}

	// LB25 (look ahead).
	if rule > 250 &&
		(state == lbPR || state == lbPO) &&
		nextProperty == prOP || nextProperty == prHY {
		var r rune
		if b != nil { // Byte slice version.
			r, _ = utf8.DecodeRune(b)
		} else { // String version.
			r, _ = utf8.DecodeRuneInString(str)
		}
		if r != utf8.RuneError {
			pr, _ := propertyLineBreak(r)
			if pr == prNU {
				return lbNU, LineDontBreak
			}
		}
	}

	// LB30 (part one).
	if rule > 300 {
		if (state == lbAL || state == lbHL || state == lbNU || state == lbNUNU) && nextProperty == prOP {
			ea := propertyEastAsianWidth(r)
			if ea != prF && ea != prW && ea != prH {
				return lbOP, LineDontBreak
			}
		} else if isCPeaFWH {
			switch nextProperty {
			case prAL:
				return lbAL, LineDontBreak
			case prHL:
				return lbHL, LineDontBreak
			case prNU:
				return lbNU, LineDontBreak
			}
		}
	}

	// LB30a.
	if newState == lbAny && nextProperty == prRI {
		if state != lbOddRI && state != lbEvenRI { // Includes state == -1.
			// Transition into the first RI.
			return lbOddRI, lineBreak
		}
		if state == lbOddRI {
			// Don't break pairs of Regional Indicators.
			return lbEvenRI, LineDontBreak
		}
		return lbOddRI, lineBreak
	}

	// LB30b.
	if rule > 302 {
		if nextProperty == prEM {
			if state == lbEB || state == lbExtPicCn {
				return prAny, LineDontBreak
			}
		}
		graphemeProperty := propertyGraphemes(r)
		if graphemeProperty == prExtendedPictographic && generalCategory == gcCn {
			return lbExtPicCn, LineCanBreak
		}
	}

	return
}
