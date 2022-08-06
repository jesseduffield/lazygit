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

// The line break parser's state transitions. It's anologous to grTransitions,
// see comments there for details. Unicode version 14.0.0.
var lbTransitions = map[[2]int][3]int{
	// LB4.
	{lbAny, prBK}: {lbBK, LineCanBreak, 310},
	{lbBK, prAny}: {lbAny, LineMustBreak, 40},

	// LB5.
	{lbAny, prCR}: {lbCR, LineCanBreak, 310},
	{lbAny, prLF}: {lbLF, LineCanBreak, 310},
	{lbAny, prNL}: {lbNL, LineCanBreak, 310},
	{lbCR, prLF}:  {lbLF, LineDontBreak, 50},
	{lbCR, prAny}: {lbAny, LineMustBreak, 50},
	{lbLF, prAny}: {lbAny, LineMustBreak, 50},
	{lbNL, prAny}: {lbAny, LineMustBreak, 50},

	// LB6.
	{lbAny, prBK}: {lbBK, LineDontBreak, 60},
	{lbAny, prCR}: {lbCR, LineDontBreak, 60},
	{lbAny, prLF}: {lbLF, LineDontBreak, 60},
	{lbAny, prNL}: {lbNL, LineDontBreak, 60},

	// LB7.
	{lbAny, prSP}: {lbSP, LineDontBreak, 70},
	{lbAny, prZW}: {lbZW, LineDontBreak, 70},

	// LB8.
	{lbZW, prSP}:  {lbZW, LineDontBreak, 70},
	{lbZW, prAny}: {lbAny, LineCanBreak, 80},

	// LB11.
	{lbAny, prWJ}: {lbWJ, LineDontBreak, 110},
	{lbWJ, prAny}: {lbAny, LineDontBreak, 110},

	// LB12.
	{lbAny, prGL}: {lbGL, LineCanBreak, 310},
	{lbGL, prAny}: {lbAny, LineDontBreak, 120},

	// LB13 (simple transitions).
	{lbAny, prCL}: {lbCL, LineCanBreak, 310},
	{lbAny, prCP}: {lbCP, LineCanBreak, 310},
	{lbAny, prEX}: {lbEX, LineDontBreak, 130},
	{lbAny, prIS}: {lbIS, LineCanBreak, 310},
	{lbAny, prSY}: {lbSY, LineCanBreak, 310},

	// LB14.
	{lbAny, prOP}: {lbOP, LineCanBreak, 310},
	{lbOP, prSP}:  {lbOP, LineDontBreak, 70},
	{lbOP, prAny}: {lbAny, LineDontBreak, 140},

	// LB15.
	{lbQU, prSP}:   {lbQUSP, LineDontBreak, 70},
	{lbQU, prOP}:   {lbOP, LineDontBreak, 150},
	{lbQUSP, prOP}: {lbOP, LineDontBreak, 150},

	// LB16.
	{lbCL, prSP}:     {lbCLCPSP, LineDontBreak, 70},
	{lbNUCL, prSP}:   {lbCLCPSP, LineDontBreak, 70},
	{lbCP, prSP}:     {lbCLCPSP, LineDontBreak, 70},
	{lbNUCP, prSP}:   {lbCLCPSP, LineDontBreak, 70},
	{lbCL, prNS}:     {lbNS, LineDontBreak, 160},
	{lbNUCL, prNS}:   {lbNS, LineDontBreak, 160},
	{lbCP, prNS}:     {lbNS, LineDontBreak, 160},
	{lbNUCP, prNS}:   {lbNS, LineDontBreak, 160},
	{lbCLCPSP, prNS}: {lbNS, LineDontBreak, 160},

	// LB17.
	{lbAny, prB2}:  {lbB2, LineCanBreak, 310},
	{lbB2, prSP}:   {lbB2SP, LineDontBreak, 70},
	{lbB2, prB2}:   {lbB2, LineDontBreak, 170},
	{lbB2SP, prB2}: {lbB2, LineDontBreak, 170},

	// LB18.
	{lbSP, prAny}:     {lbAny, LineCanBreak, 180},
	{lbQUSP, prAny}:   {lbAny, LineCanBreak, 180},
	{lbCLCPSP, prAny}: {lbAny, LineCanBreak, 180},
	{lbB2SP, prAny}:   {lbAny, LineCanBreak, 180},

	// LB19.
	{lbAny, prQU}: {lbQU, LineDontBreak, 190},
	{lbQU, prAny}: {lbAny, LineDontBreak, 190},

	// LB20.
	{lbAny, prCB}: {lbCB, LineCanBreak, 200},
	{lbCB, prAny}: {lbAny, LineCanBreak, 200},

	// LB21.
	{lbAny, prBA}: {lbBA, LineDontBreak, 210},
	{lbAny, prHY}: {lbHY, LineDontBreak, 210},
	{lbAny, prNS}: {lbNS, LineDontBreak, 210},
	{lbAny, prBB}: {lbBB, LineCanBreak, 310},
	{lbBB, prAny}: {lbAny, LineDontBreak, 210},

	// LB21a.
	{lbAny, prHL}:    {lbHL, LineCanBreak, 310},
	{lbHL, prHY}:     {lbLB21a, LineDontBreak, 210},
	{lbHL, prBA}:     {lbLB21a, LineDontBreak, 210},
	{lbLB21a, prAny}: {lbAny, LineDontBreak, 211},

	// LB21b.
	{lbSY, prHL}:   {lbHL, LineDontBreak, 212},
	{lbNUSY, prHL}: {lbHL, LineDontBreak, 212},

	// LB22.
	{lbAny, prIN}: {lbAny, LineDontBreak, 220},

	// LB23.
	{lbAny, prAL}:  {lbAL, LineCanBreak, 310},
	{lbAny, prNU}:  {lbNU, LineCanBreak, 310},
	{lbAL, prNU}:   {lbNU, LineDontBreak, 230},
	{lbHL, prNU}:   {lbNU, LineDontBreak, 230},
	{lbNU, prAL}:   {lbAL, LineDontBreak, 230},
	{lbNU, prHL}:   {lbHL, LineDontBreak, 230},
	{lbNUNU, prAL}: {lbAL, LineDontBreak, 230},
	{lbNUNU, prHL}: {lbHL, LineDontBreak, 230},

	// LB23a.
	{lbAny, prPR}:  {lbPR, LineCanBreak, 310},
	{lbAny, prID}:  {lbIDEM, LineCanBreak, 310},
	{lbAny, prEB}:  {lbEB, LineCanBreak, 310},
	{lbAny, prEM}:  {lbIDEM, LineCanBreak, 310},
	{lbPR, prID}:   {lbIDEM, LineDontBreak, 231},
	{lbPR, prEB}:   {lbEB, LineDontBreak, 231},
	{lbPR, prEM}:   {lbIDEM, LineDontBreak, 231},
	{lbIDEM, prPO}: {lbPO, LineDontBreak, 231},
	{lbEB, prPO}:   {lbPO, LineDontBreak, 231},

	// LB24.
	{lbAny, prPO}: {lbPO, LineCanBreak, 310},
	{lbPR, prAL}:  {lbAL, LineDontBreak, 240},
	{lbPR, prHL}:  {lbHL, LineDontBreak, 240},
	{lbPO, prAL}:  {lbAL, LineDontBreak, 240},
	{lbPO, prHL}:  {lbHL, LineDontBreak, 240},
	{lbAL, prPR}:  {lbPR, LineDontBreak, 240},
	{lbAL, prPO}:  {lbPO, LineDontBreak, 240},
	{lbHL, prPR}:  {lbPR, LineDontBreak, 240},
	{lbHL, prPO}:  {lbPO, LineDontBreak, 240},

	// LB25 (simple transitions).
	{lbPR, prNU}:   {lbNU, LineDontBreak, 250},
	{lbPO, prNU}:   {lbNU, LineDontBreak, 250},
	{lbOP, prNU}:   {lbNU, LineDontBreak, 250},
	{lbHY, prNU}:   {lbNU, LineDontBreak, 250},
	{lbNU, prNU}:   {lbNUNU, LineDontBreak, 250},
	{lbNU, prSY}:   {lbNUSY, LineDontBreak, 250},
	{lbNU, prIS}:   {lbNUIS, LineDontBreak, 250},
	{lbNUNU, prNU}: {lbNUNU, LineDontBreak, 250},
	{lbNUNU, prSY}: {lbNUSY, LineDontBreak, 250},
	{lbNUNU, prIS}: {lbNUIS, LineDontBreak, 250},
	{lbNUSY, prNU}: {lbNUNU, LineDontBreak, 250},
	{lbNUSY, prSY}: {lbNUSY, LineDontBreak, 250},
	{lbNUSY, prIS}: {lbNUIS, LineDontBreak, 250},
	{lbNUIS, prNU}: {lbNUNU, LineDontBreak, 250},
	{lbNUIS, prSY}: {lbNUSY, LineDontBreak, 250},
	{lbNUIS, prIS}: {lbNUIS, LineDontBreak, 250},
	{lbNU, prCL}:   {lbNUCL, LineDontBreak, 250},
	{lbNU, prCP}:   {lbNUCP, LineDontBreak, 250},
	{lbNUNU, prCL}: {lbNUCL, LineDontBreak, 250},
	{lbNUNU, prCP}: {lbNUCP, LineDontBreak, 250},
	{lbNUSY, prCL}: {lbNUCL, LineDontBreak, 250},
	{lbNUSY, prCP}: {lbNUCP, LineDontBreak, 250},
	{lbNUIS, prCL}: {lbNUCL, LineDontBreak, 250},
	{lbNUIS, prCP}: {lbNUCP, LineDontBreak, 250},
	{lbNU, prPO}:   {lbPO, LineDontBreak, 250},
	{lbNUNU, prPO}: {lbPO, LineDontBreak, 250},
	{lbNUSY, prPO}: {lbPO, LineDontBreak, 250},
	{lbNUIS, prPO}: {lbPO, LineDontBreak, 250},
	{lbNUCL, prPO}: {lbPO, LineDontBreak, 250},
	{lbNUCP, prPO}: {lbPO, LineDontBreak, 250},
	{lbNU, prPR}:   {lbPR, LineDontBreak, 250},
	{lbNUNU, prPR}: {lbPR, LineDontBreak, 250},
	{lbNUSY, prPR}: {lbPR, LineDontBreak, 250},
	{lbNUIS, prPR}: {lbPR, LineDontBreak, 250},
	{lbNUCL, prPR}: {lbPR, LineDontBreak, 250},
	{lbNUCP, prPR}: {lbPR, LineDontBreak, 250},

	// LB26.
	{lbAny, prJL}: {lbJL, LineCanBreak, 310},
	{lbAny, prJV}: {lbJV, LineCanBreak, 310},
	{lbAny, prJT}: {lbJT, LineCanBreak, 310},
	{lbAny, prH2}: {lbH2, LineCanBreak, 310},
	{lbAny, prH3}: {lbH3, LineCanBreak, 310},
	{lbJL, prJL}:  {lbJL, LineDontBreak, 260},
	{lbJL, prJV}:  {lbJV, LineDontBreak, 260},
	{lbJL, prH2}:  {lbH2, LineDontBreak, 260},
	{lbJL, prH3}:  {lbH3, LineDontBreak, 260},
	{lbJV, prJV}:  {lbJV, LineDontBreak, 260},
	{lbJV, prJT}:  {lbJT, LineDontBreak, 260},
	{lbH2, prJV}:  {lbJV, LineDontBreak, 260},
	{lbH2, prJT}:  {lbJT, LineDontBreak, 260},
	{lbJT, prJT}:  {lbJT, LineDontBreak, 260},
	{lbH3, prJT}:  {lbJT, LineDontBreak, 260},

	// LB27.
	{lbJL, prPO}: {lbPO, LineDontBreak, 270},
	{lbJV, prPO}: {lbPO, LineDontBreak, 270},
	{lbJT, prPO}: {lbPO, LineDontBreak, 270},
	{lbH2, prPO}: {lbPO, LineDontBreak, 270},
	{lbH3, prPO}: {lbPO, LineDontBreak, 270},
	{lbPR, prJL}: {lbJL, LineDontBreak, 270},
	{lbPR, prJV}: {lbJV, LineDontBreak, 270},
	{lbPR, prJT}: {lbJT, LineDontBreak, 270},
	{lbPR, prH2}: {lbH2, LineDontBreak, 270},
	{lbPR, prH3}: {lbH3, LineDontBreak, 270},

	// LB28.
	{lbAL, prAL}: {lbAL, LineDontBreak, 280},
	{lbAL, prHL}: {lbHL, LineDontBreak, 280},
	{lbHL, prAL}: {lbAL, LineDontBreak, 280},
	{lbHL, prHL}: {lbHL, LineDontBreak, 280},

	// LB29.
	{lbIS, prAL}:   {lbAL, LineDontBreak, 290},
	{lbIS, prHL}:   {lbHL, LineDontBreak, 290},
	{lbNUIS, prAL}: {lbAL, LineDontBreak, 290},
	{lbNUIS, prHL}: {lbHL, LineDontBreak, 290},
}

// transitionLineBreakState determines the new state of the line break parser
// given the current state and the next code point. It also returns the type of
// line break: LineDontBreak, LineCanBreak, or LineMustBreak. If more than one
// code point is needed to determine the new state, the byte slice or the string
// starting after rune "r" can be used (whichever is not nil or empty) for
// further lookups.
func transitionLineBreakState(state int, r rune, b []byte, str string) (newState int, lineBreak int) {
	// Determine the property of the next character.
	nextProperty, generalCategory := propertyWithGenCat(lineBreakCodePoints, r)

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
			ea := property(eastAsianWidth, r)
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
	transition, ok := lbTransitions[[2]int{state, nextProperty}]
	if ok {
		// We have a specific transition. We'll use it.
		newState, lineBreak, rule = transition[0], transition[1], transition[2]
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := lbTransitions[[2]int{state, prAny}]
		transAnyState, okAnyState := lbTransitions[[2]int{lbAny, nextProperty}]
		if okAnyProp && okAnyState {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, lineBreak, rule = transAnyState[0], transAnyState[1], transAnyState[2]
			if transAnyProp[2] < transAnyState[2] {
				lineBreak, rule = transAnyProp[1], transAnyProp[2]
			}
		} else if okAnyProp {
			// We only have a specific state.
			newState, lineBreak, rule = transAnyProp[0], transAnyProp[1], transAnyProp[2]
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if okAnyState {
			// We only have a specific property.
			newState, lineBreak, rule = transAnyState[0], transAnyState[1], transAnyState[2]
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
			pr, _ := propertyWithGenCat(lineBreakCodePoints, r)
			if pr == prNU {
				return lbNU, LineDontBreak
			}
		}
	}

	// LB30 (part one).
	if rule > 300 {
		if (state == lbAL || state == lbHL || state == lbNU || state == lbNUNU) && nextProperty == prOP {
			ea := property(eastAsianWidth, r)
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
		graphemeProperty := property(graphemeCodePoints, r)
		if graphemeProperty == prExtendedPictographic && generalCategory == gcCn {
			return lbExtPicCn, LineCanBreak
		}
	}

	return
}
