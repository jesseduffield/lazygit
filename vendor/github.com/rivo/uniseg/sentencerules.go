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

// The sentence break parser's breaking instructions.
const (
	sbDontBreak = iota
	sbBreak
)

// The sentence break parser's state transitions. It's anologous to
// grTransitions, see comments there for details. Unicode version 14.0.0.
var sbTransitions = map[[2]int][3]int{
	// SB3.
	{sbAny, prCR}: {sbCR, sbDontBreak, 9990},
	{sbCR, prLF}:  {sbParaSep, sbDontBreak, 30},

	// SB4.
	{sbAny, prSep}:     {sbParaSep, sbDontBreak, 9990},
	{sbAny, prLF}:      {sbParaSep, sbDontBreak, 9990},
	{sbParaSep, prAny}: {sbAny, sbBreak, 40},
	{sbCR, prAny}:      {sbAny, sbBreak, 40},

	// SB6.
	{sbAny, prATerm}:     {sbATerm, sbDontBreak, 9990},
	{sbATerm, prNumeric}: {sbAny, sbDontBreak, 60},
	{sbSB7, prNumeric}:   {sbAny, sbDontBreak, 60}, // Because ATerm also appears in SB7.

	// SB7.
	{sbAny, prUpper}:   {sbUpper, sbDontBreak, 9990},
	{sbAny, prLower}:   {sbLower, sbDontBreak, 9990},
	{sbUpper, prATerm}: {sbSB7, sbDontBreak, 70},
	{sbLower, prATerm}: {sbSB7, sbDontBreak, 70},
	{sbSB7, prUpper}:   {sbUpper, sbDontBreak, 70},

	// SB8a.
	{sbAny, prSTerm}:           {sbSTerm, sbDontBreak, 9990},
	{sbATerm, prSContinue}:     {sbAny, sbDontBreak, 81},
	{sbATerm, prATerm}:         {sbATerm, sbDontBreak, 81},
	{sbATerm, prSTerm}:         {sbSTerm, sbDontBreak, 81},
	{sbSB7, prSContinue}:       {sbAny, sbDontBreak, 81},
	{sbSB7, prATerm}:           {sbATerm, sbDontBreak, 81},
	{sbSB7, prSTerm}:           {sbSTerm, sbDontBreak, 81},
	{sbSB8Close, prSContinue}:  {sbAny, sbDontBreak, 81},
	{sbSB8Close, prATerm}:      {sbATerm, sbDontBreak, 81},
	{sbSB8Close, prSTerm}:      {sbSTerm, sbDontBreak, 81},
	{sbSB8Sp, prSContinue}:     {sbAny, sbDontBreak, 81},
	{sbSB8Sp, prATerm}:         {sbATerm, sbDontBreak, 81},
	{sbSB8Sp, prSTerm}:         {sbSTerm, sbDontBreak, 81},
	{sbSTerm, prSContinue}:     {sbAny, sbDontBreak, 81},
	{sbSTerm, prATerm}:         {sbATerm, sbDontBreak, 81},
	{sbSTerm, prSTerm}:         {sbSTerm, sbDontBreak, 81},
	{sbSB8aClose, prSContinue}: {sbAny, sbDontBreak, 81},
	{sbSB8aClose, prATerm}:     {sbATerm, sbDontBreak, 81},
	{sbSB8aClose, prSTerm}:     {sbSTerm, sbDontBreak, 81},
	{sbSB8aSp, prSContinue}:    {sbAny, sbDontBreak, 81},
	{sbSB8aSp, prATerm}:        {sbATerm, sbDontBreak, 81},
	{sbSB8aSp, prSTerm}:        {sbSTerm, sbDontBreak, 81},

	// SB9.
	{sbATerm, prClose}:     {sbSB8Close, sbDontBreak, 90},
	{sbSB7, prClose}:       {sbSB8Close, sbDontBreak, 90},
	{sbSB8Close, prClose}:  {sbSB8Close, sbDontBreak, 90},
	{sbATerm, prSp}:        {sbSB8Sp, sbDontBreak, 90},
	{sbSB7, prSp}:          {sbSB8Sp, sbDontBreak, 90},
	{sbSB8Close, prSp}:     {sbSB8Sp, sbDontBreak, 90},
	{sbSTerm, prClose}:     {sbSB8aClose, sbDontBreak, 90},
	{sbSB8aClose, prClose}: {sbSB8aClose, sbDontBreak, 90},
	{sbSTerm, prSp}:        {sbSB8aSp, sbDontBreak, 90},
	{sbSB8aClose, prSp}:    {sbSB8aSp, sbDontBreak, 90},
	{sbATerm, prSep}:       {sbParaSep, sbDontBreak, 90},
	{sbATerm, prCR}:        {sbParaSep, sbDontBreak, 90},
	{sbATerm, prLF}:        {sbParaSep, sbDontBreak, 90},
	{sbSB7, prSep}:         {sbParaSep, sbDontBreak, 90},
	{sbSB7, prCR}:          {sbParaSep, sbDontBreak, 90},
	{sbSB7, prLF}:          {sbParaSep, sbDontBreak, 90},
	{sbSB8Close, prSep}:    {sbParaSep, sbDontBreak, 90},
	{sbSB8Close, prCR}:     {sbParaSep, sbDontBreak, 90},
	{sbSB8Close, prLF}:     {sbParaSep, sbDontBreak, 90},
	{sbSTerm, prSep}:       {sbParaSep, sbDontBreak, 90},
	{sbSTerm, prCR}:        {sbParaSep, sbDontBreak, 90},
	{sbSTerm, prLF}:        {sbParaSep, sbDontBreak, 90},
	{sbSB8aClose, prSep}:   {sbParaSep, sbDontBreak, 90},
	{sbSB8aClose, prCR}:    {sbParaSep, sbDontBreak, 90},
	{sbSB8aClose, prLF}:    {sbParaSep, sbDontBreak, 90},

	// SB10.
	{sbSB8Sp, prSp}:  {sbSB8Sp, sbDontBreak, 100},
	{sbSB8aSp, prSp}: {sbSB8aSp, sbDontBreak, 100},
	{sbSB8Sp, prSep}: {sbParaSep, sbDontBreak, 100},
	{sbSB8Sp, prCR}:  {sbParaSep, sbDontBreak, 100},
	{sbSB8Sp, prLF}:  {sbParaSep, sbDontBreak, 100},

	// SB11.
	{sbATerm, prAny}:     {sbAny, sbBreak, 110},
	{sbSB7, prAny}:       {sbAny, sbBreak, 110},
	{sbSB8Close, prAny}:  {sbAny, sbBreak, 110},
	{sbSB8Sp, prAny}:     {sbAny, sbBreak, 110},
	{sbSTerm, prAny}:     {sbAny, sbBreak, 110},
	{sbSB8aClose, prAny}: {sbAny, sbBreak, 110},
	{sbSB8aSp, prAny}:    {sbAny, sbBreak, 110},
	// We'll always break after ParaSep due to SB4.
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
	transition, ok := sbTransitions[[2]int{state, nextProperty}]
	if ok {
		// We have a specific transition. We'll use it.
		newState, sentenceBreak, rule = transition[0], transition[1] == sbBreak, transition[2]
	} else {
		// No specific transition found. Try the less specific ones.
		transAnyProp, okAnyProp := sbTransitions[[2]int{state, prAny}]
		transAnyState, okAnyState := sbTransitions[[2]int{sbAny, nextProperty}]
		if okAnyProp && okAnyState {
			// Both apply. We'll use a mix (see comments for grTransitions).
			newState, sentenceBreak, rule = transAnyState[0], transAnyState[1] == sbBreak, transAnyState[2]
			if transAnyProp[2] < transAnyState[2] {
				sentenceBreak, rule = transAnyProp[1] == sbBreak, transAnyProp[2]
			}
		} else if okAnyProp {
			// We only have a specific state.
			newState, sentenceBreak, rule = transAnyProp[0], transAnyProp[1] == sbBreak, transAnyProp[2]
			// This branch will probably never be reached because okAnyState will
			// always be true given the current transition map. But we keep it here
			// for future modifications to the transition map where this may not be
			// true anymore.
		} else if okAnyState {
			// We only have a specific property.
			newState, sentenceBreak, rule = transAnyState[0], transAnyState[1] == sbBreak, transAnyState[2]
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
