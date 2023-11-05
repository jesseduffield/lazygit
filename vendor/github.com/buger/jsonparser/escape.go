package jsonparser

import (
	"bytes"
	"unicode/utf8"
)

// JSON Unicode stuff: see https://tools.ietf.org/html/rfc7159#section-7

const supplementalPlanesOffset = 0x10000
const highSurrogateOffset = 0xD800
const lowSurrogateOffset = 0xDC00

const basicMultilingualPlaneReservedOffset = 0xDFFF
const basicMultilingualPlaneOffset = 0xFFFF

func combineUTF16Surrogates(high, low rune) rune {
	return supplementalPlanesOffset + (high-highSurrogateOffset)<<10 + (low - lowSurrogateOffset)
}

const badHex = -1

func h2I(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'A' && c <= 'F':
		return int(c - 'A' + 10)
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	}
	return badHex
}

// decodeSingleUnicodeEscape decodes a single \uXXXX escape sequence. The prefix \u is assumed to be present and
// is not checked.
// In JSON, these escapes can either come alone or as part of "UTF16 surrogate pairs" that must be handled together.
// This function only handles one; decodeUnicodeEscape handles this more complex case.
func decodeSingleUnicodeEscape(in []byte) (rune, bool) {
	// We need at least 6 characters total
	if len(in) < 6 {
		return utf8.RuneError, false
	}

	// Convert hex to decimal
	h1, h2, h3, h4 := h2I(in[2]), h2I(in[3]), h2I(in[4]), h2I(in[5])
	if h1 == badHex || h2 == badHex || h3 == badHex || h4 == badHex {
		return utf8.RuneError, false
	}

	// Compose the hex digits
	return rune(h1<<12 + h2<<8 + h3<<4 + h4), true
}

// isUTF16EncodedRune checks if a rune is in the range for non-BMP characters,
// which is used to describe UTF16 chars.
// Source: https://en.wikipedia.org/wiki/Plane_(Unicode)#Basic_Multilingual_Plane
func isUTF16EncodedRune(r rune) bool {
	return highSurrogateOffset <= r && r <= basicMultilingualPlaneReservedOffset
}

func decodeUnicodeEscape(in []byte) (rune, int) {
	if r, ok := decodeSingleUnicodeEscape(in); !ok {
		// Invalid Unicode escape
		return utf8.RuneError, -1
	} else if r <= basicMultilingualPlaneOffset && !isUTF16EncodedRune(r) {
		// Valid Unicode escape in Basic Multilingual Plane
		return r, 6
	} else if r2, ok := decodeSingleUnicodeEscape(in[6:]); !ok { // Note: previous decodeSingleUnicodeEscape success guarantees at least 6 bytes remain
		// UTF16 "high surrogate" without manditory valid following Unicode escape for the "low surrogate"
		return utf8.RuneError, -1
	} else if r2 < lowSurrogateOffset {
		// Invalid UTF16 "low surrogate"
		return utf8.RuneError, -1
	} else {
		// Valid UTF16 surrogate pair
		return combineUTF16Surrogates(r, r2), 12
	}
}

// backslashCharEscapeTable: when '\X' is found for some byte X, it is to be replaced with backslashCharEscapeTable[X]
var backslashCharEscapeTable = [...]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

// unescapeToUTF8 unescapes the single escape sequence starting at 'in' into 'out' and returns
// how many characters were consumed from 'in' and emitted into 'out'.
// If a valid escape sequence does not appear as a prefix of 'in', (-1, -1) to signal the error.
func unescapeToUTF8(in, out []byte) (inLen int, outLen int) {
	if len(in) < 2 || in[0] != '\\' {
		// Invalid escape due to insufficient characters for any escape or no initial backslash
		return -1, -1
	}

	// https://tools.ietf.org/html/rfc7159#section-7
	switch e := in[1]; e {
	case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
		// Valid basic 2-character escapes (use lookup table)
		out[0] = backslashCharEscapeTable[e]
		return 2, 1
	case 'u':
		// Unicode escape
		if r, inLen := decodeUnicodeEscape(in); inLen == -1 {
			// Invalid Unicode escape
			return -1, -1
		} else {
			// Valid Unicode escape; re-encode as UTF8
			outLen := utf8.EncodeRune(out, r)
			return inLen, outLen
		}
	}

	return -1, -1
}

// unescape unescapes the string contained in 'in' and returns it as a slice.
// If 'in' contains no escaped characters:
//   Returns 'in'.
// Else, if 'out' is of sufficient capacity (guaranteed if cap(out) >= len(in)):
//   'out' is used to build the unescaped string and is returned with no extra allocation
// Else:
//   A new slice is allocated and returned.
func Unescape(in, out []byte) ([]byte, error) {
	firstBackslash := bytes.IndexByte(in, '\\')
	if firstBackslash == -1 {
		return in, nil
	}

	// Get a buffer of sufficient size (allocate if needed)
	if cap(out) < len(in) {
		out = make([]byte, len(in))
	} else {
		out = out[0:len(in)]
	}

	// Copy the first sequence of unescaped bytes to the output and obtain a buffer pointer (subslice)
	copy(out, in[:firstBackslash])
	in = in[firstBackslash:]
	buf := out[firstBackslash:]

	for len(in) > 0 {
		// Unescape the next escaped character
		inLen, bufLen := unescapeToUTF8(in, buf)
		if inLen == -1 {
			return nil, MalformedStringEscapeError
		}

		in = in[inLen:]
		buf = buf[bufLen:]

		// Copy everything up until the next backslash
		nextBackslash := bytes.IndexByte(in, '\\')
		if nextBackslash == -1 {
			copy(buf, in)
			buf = buf[len(in):]
			break
		} else {
			copy(buf, in[:nextBackslash])
			buf = buf[nextBackslash:]
			in = in[nextBackslash:]
		}
	}

	// Trim the out buffer to the amount that was actually emitted
	return out[:len(out)-len(buf)], nil
}
