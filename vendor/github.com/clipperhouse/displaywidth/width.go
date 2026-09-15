package displaywidth

import (
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
)

// String calculates the display width of a string,
// by iterating over grapheme clusters in the string
// and summing their widths.
func String(s string) int {
	return DefaultOptions.String(s)
}

// String calculates the display width of a string, for the given options, by
// iterating over grapheme clusters in the string and summing their widths.
func (options Options) String(s string) int {
	width := 0
	pos := 0

	for pos < len(s) {
		// Try ASCII optimization
		asciiLen := printableASCIILength(s[pos:])
		if asciiLen > 0 {
			width += asciiLen
			pos += asciiLen
			continue
		}

		// Not ASCII, use grapheme parsing
		g := graphemes.FromString(s[pos:])
		g.AnsiEscapeSequences = options.ControlSequences
		g.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

		start := pos

		for g.Next() {
			v := g.Value()
			width += graphemeWidth(v, options)
			pos += len(v)

			// Quick check: if remaining might have printable ASCII, break to outer loop
			if pos < len(s) && s[pos] >= 0x20 && s[pos] <= 0x7E {
				break
			}
		}

		// Defensive, should not happen: if no progress was made,
		// skip a byte to prevent infinite loop. Only applies if
		// the grapheme parser misbehaves.
		if pos == start {
			pos++
		}
	}

	return width
}

// Bytes calculates the display width of a []byte,
// by iterating over grapheme clusters in the byte slice
// and summing their widths.
func Bytes(s []byte) int {
	return DefaultOptions.Bytes(s)
}

// Bytes calculates the display width of a []byte, for the given options, by
// iterating over grapheme clusters in the slice and summing their widths.
func (options Options) Bytes(s []byte) int {
	width := 0
	pos := 0

	for pos < len(s) {
		// Try ASCII optimization
		asciiLen := printableASCIILength(s[pos:])
		if asciiLen > 0 {
			width += asciiLen
			pos += asciiLen
			continue
		}

		// Not ASCII, use grapheme parsing
		g := graphemes.FromBytes(s[pos:])
		g.AnsiEscapeSequences = options.ControlSequences
		g.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

		start := pos

		for g.Next() {
			v := g.Value()
			width += graphemeWidth(v, options)
			pos += len(v)

			// Quick check: if remaining might have printable ASCII, break to outer loop
			if pos < len(s) && s[pos] >= 0x20 && s[pos] <= 0x7E {
				break
			}
		}

		// Defensive, should not happen: if no progress was made,
		// skip a byte to prevent infinite loop. Only applies if
		// the grapheme parser misbehaves.
		if pos == start {
			pos++
		}
	}

	return width
}

// Rune calculates the display width of a rune. You
// should almost certainly use [String] or [Bytes] for
// most purposes.
//
// The smallest unit of display width is a grapheme
// cluster, not a rune. Iterating over runes to measure
// width is incorrect in many cases.
func Rune(r rune) int {
	return DefaultOptions.Rune(r)
}

// Rune calculates the display width of a rune, for the given options.
//
// You should almost certainly use [String] or [Bytes] for most purposes.
//
// The smallest unit of display width is a grapheme cluster, not a rune.
// Iterating over runes to measure width is incorrect in many cases.
func (options Options) Rune(r rune) int {
	if r < utf8.RuneSelf {
		return asciiWidth(byte(r))
	}

	// Surrogates (U+D800-U+DFFF) are invalid UTF-8.
	if r >= 0xD800 && r <= 0xDFFF {
		return 0
	}

	var buf [4]byte
	n := utf8.EncodeRune(buf[:], r)

	// Skip the grapheme iterator
	return graphemeWidth(buf[:n], options)
}

const _Default property = 0

// graphemeWidth returns the display width of a grapheme cluster.
// The passed string must be a single grapheme cluster.
func graphemeWidth[T ~string | []byte](s T, options Options) int {
	if len(s) == 0 {
		return 0
	}

	// C1 controls (0x80-0x9F) are zero-width when 8-bit control sequences
	// are enabled. This must be checked before the single-byte optimization
	// below, which would otherwise return width 1 for these bytes.
	if options.ControlSequences8Bit && s[0] >= 0x80 && s[0] <= 0x9F {
		return 0
	}

	// Optimization: single-byte graphemes need no property lookup
	if len(s) == 1 {
		return asciiWidth(s[0])
	}

	// Multi-byte grapheme clusters led by a C0 control (0x00-0x1F)
	if s[0] <= 0x1F {
		return 0
	}

	p, sz := lookup(s)
	prop := property(p)

	// Variation Selector 16 (VS16) requests emoji presentation
	if prop != _Wide && sz > 0 && len(s) >= sz+3 {
		vs := s[sz : sz+3]
		if isVS16(vs) {
			prop = _Wide
		}
		// VS15 (0x8E) requests text presentation but does not affect width,
		// in my reading of Unicode TR51. Falls through to return the base
		// character's property.
	}

	if options.EastAsianWidth && prop == _East_Asian_Ambiguous {
		prop = _Wide
	}

	if prop > upperBound {
		prop = _Default
	}

	return propertyWidths[prop]
}

func asciiWidth(b byte) int {
	if b <= 0x1F || b == 0x7F {
		return 0
	}
	return 1
}

// printableASCIILength returns the length of consecutive printable ASCII bytes
// starting at the beginning of s.
func printableASCIILength[T string | []byte](s T) int {
	i := 0
	for ; i < len(s); i++ {
		b := s[i]
		// Printable ASCII is 0x20-0x7E (space through tilde)
		if b < 0x20 || b > 0x7E {
			break
		}
	}

	// If the next byte is non-ASCII (>= 0x80), back off by 1. The grapheme
	// parser may group the last ASCII byte with subsequent non-ASCII bytes,
	// such as combining marks.
	if i > 0 && i < len(s) && s[i] >= 0x80 {
		i--
	}

	return i
}

// isVS16 checks if the slice matches VS16 (U+FE0F) UTF-8 encoding
// (EF B8 8F). It assumes len(s) >= 3.
func isVS16[T ~string | []byte](s T) bool {
	return s[0] == 0xEF && s[1] == 0xB8 && s[2] == 0x8F
}

// propertyWidths is a jump table of sorts, instead of a switch
var propertyWidths = [4]int{
	_Default:              1,
	_Zero_Width:           0,
	_Wide:                 2,
	_East_Asian_Ambiguous: 1,
}

const upperBound = property(len(propertyWidths) - 1)
