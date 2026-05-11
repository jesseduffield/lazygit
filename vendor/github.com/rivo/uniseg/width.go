package uniseg

// EastAsianAmbiguousWidth specifies the monospace width for East Asian
// characters classified as Ambiguous. The default is 1 but some rare fonts
// render them with a width of 2.
var EastAsianAmbiguousWidth = 1

// runeWidth returns the monospace width for the given rune. The provided
// grapheme property is a value mapped by the [graphemeCodePoints] table.
//
// Every rune has a width of 1, except for runes with the following properties
// (evaluated in this order):
//
//   - Control, CR, LF, Extend, ZWJ: Width of 0
//   - \u2e3a, TWO-EM DASH: Width of 3
//   - \u2e3b, THREE-EM DASH: Width of 4
//   - East-Asian width Fullwidth and Wide: Width of 2 (Ambiguous and Neutral
//     have a width of 1)
//   - Regional Indicator: Width of 2
//   - Extended Pictographic: Width of 2, unless Emoji Presentation is "No".
func runeWidth(r rune, graphemeProperty int) int {
	switch graphemeProperty {
	case prControl, prCR, prLF, prExtend, prZWJ:
		return 0
	case prRegionalIndicator:
		return 2
	case prExtendedPictographic:
		if property(emojiPresentation, r) == prEmojiPresentation {
			return 2
		}
		return 1
	}

	switch r {
	case 0x2e3a:
		return 3
	case 0x2e3b:
		return 4
	}

	switch propertyEastAsianWidth(r) {
	case prW, prF:
		return 2
	case prA:
		return EastAsianAmbiguousWidth
	}

	return 1
}

// StringWidth returns the monospace width for the given string, that is, the
// number of same-size cells to be occupied by the string.
func StringWidth(s string) (width int) {
	state := -1
	for len(s) > 0 {
		var w int
		_, s, w, state = FirstGraphemeClusterInString(s, state)
		width += w
	}
	return
}
