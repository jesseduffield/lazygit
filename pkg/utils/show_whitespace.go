package utils

import (
	"strings"
	"unicode/utf8"
)

const (
	spaceMarker = "\u00B7"
	tabMarker   = "\u2192"
	tabFillChar = "-"

	faintOn  = "\x1b[2m"
	faintOff = "\x1b[22m"
)

// ShowWhitespaceCharacters replaces spaces and tabs in the given content with
// visual indicators. ANSI escape sequences are preserved and whitespace within
// them is not replaced. Spaces are replaced with middle dots (·) and tabs are
// replaced with a right arrow (→) followed by horizontal lines to fill the
// tab stop, so the layout stays stable when toggling.
func ShowWhitespaceCharacters(content string, tabWidth int) string {
	var result strings.Builder
	result.Grow(len(content))

	if tabWidth < 1 {
		tabWidth = 4
	}

	inEscape := false
	inEscapeCSI := false
	col := 0

	for i := 0; i < len(content); {
		if inEscape {
			result.WriteByte(content[i])
			if inEscapeCSI {
				if (content[i] >= 'a' && content[i] <= 'z') || (content[i] >= 'A' && content[i] <= 'Z') {
					inEscape = false
					inEscapeCSI = false
				}
			} else {
				if content[i] == '[' {
					inEscapeCSI = true
				} else {
					inEscape = false
				}
			}
			i++
			continue
		}

		if content[i] == '\x1b' {
			result.WriteByte(content[i])
			inEscape = true
			i++
			continue
		}

		if content[i] == '\n' {
			result.WriteByte(content[i])
			col = 0
			i++
			continue
		}

		if content[i] == '\t' {
			numFill := tabWidth - (col % tabWidth)
			result.WriteString(faintOn)
			result.WriteString(tabMarker)
			col++
			for f := 1; f < numFill; f++ {
				result.WriteString(tabFillChar)
				col++
			}
			result.WriteString(faintOff)
			i++
			continue
		}

		if content[i] == ' ' {
			result.WriteString(faintOn)
			result.WriteString(spaceMarker)
			result.WriteString(faintOff)
			col++
			i++
			continue
		}

		r, size := utf8.DecodeRuneInString(content[i:])
		result.WriteRune(r)
		col++
		i += size
	}

	return result.String()
}
