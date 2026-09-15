package displaywidth

import (
	"strings"

	"github.com/clipperhouse/uax29/v2/graphemes"
)

// TruncateString truncates a string to the given maxWidth, and appends the
// given tail if the string is truncated.
//
// It ensures the visible width, including the width of the tail, is less than or
// equal to maxWidth.
//
// When [Options.ControlSequences] is true, 7-bit ANSI escape sequences that
// appear after the truncation point are preserved in the output. This ensures
// that escape sequences such as SGR resets are not lost, preventing color
// bleed in terminal output.
//
// [Options.ControlSequences8Bit] is ignored by truncation. 8-bit C1 byte values
// (0x80-0x9F) overlap with UTF-8 multi-byte encoding, so manipulating them
// during truncation can shift byte boundaries and form unintended visible
// characters. Use [Options.String] or [Options.Bytes] for 8-bit-aware width
// measurement.
func (options Options) TruncateString(s string, maxWidth int, tail string) string {
	// We deliberately ignore ControlSequences8Bit for truncation, see above.
	options.ControlSequences8Bit = false

	maxWidthWithoutTail := maxWidth - options.String(tail)

	var pos, total int
	g := graphemes.FromString(s)
	g.AnsiEscapeSequences = options.ControlSequences

	for g.Next() {
		gw := graphemeWidth(g.Value(), options)
		if total+gw <= maxWidthWithoutTail {
			pos = g.End()
		}
		total += gw
		if total > maxWidth {
			if options.ControlSequences {
				// Build result with trailing 7-bit ANSI escape sequences preserved
				var b strings.Builder
				b.Grow(len(s) + len(tail)) // at most original + tail
				b.WriteString(s[:pos])
				b.WriteString(tail)

				rem := graphemes.FromString(s[pos:])
				rem.AnsiEscapeSequences = options.ControlSequences

				for rem.Next() {
					v := rem.Value()
					// Only preserve 7-bit escapes (ESC = 0x1B) that measure
					// as zero-width on their own; some sequences (e.g. SOS)
					// are only valid in their original context.
					if len(v) > 0 && v[0] == 0x1B && options.String(v) == 0 {
						b.WriteString(v)
					}
				}
				return b.String()
			}
			return s[:pos] + tail
		}
	}
	// No truncation
	return s
}

// TruncateString truncates a string to the given maxWidth, and appends the
// given tail if the string is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func TruncateString(s string, maxWidth int, tail string) string {
	return DefaultOptions.TruncateString(s, maxWidth, tail)
}

// TruncateBytes truncates a []byte to the given maxWidth, and appends the
// given tail if the []byte is truncated.
//
// It ensures the visible width, including the width of the tail, is less than or
// equal to maxWidth.
//
// When [Options.ControlSequences] is true, 7-bit ANSI escape sequences that
// appear after the truncation point are preserved in the output. This ensures
// that escape sequences such as SGR resets are not lost, preventing color
// bleed in terminal output.
//
// [Options.ControlSequences8Bit] is ignored by truncation. 8-bit C1 byte values
// (0x80-0x9F) overlap with UTF-8 multi-byte encoding, so manipulating them
// during truncation can shift byte boundaries and form unintended visible
// characters. Use [Options.String] or [Options.Bytes] for 8-bit-aware width
// measurement.
func (options Options) TruncateBytes(s []byte, maxWidth int, tail []byte) []byte {
	// We deliberately ignore ControlSequences8Bit for truncation, see above.
	options.ControlSequences8Bit = false

	maxWidthWithoutTail := maxWidth - options.Bytes(tail)

	var pos, total int
	g := graphemes.FromBytes(s)
	g.AnsiEscapeSequences = options.ControlSequences

	for g.Next() {
		gw := graphemeWidth(g.Value(), options)
		if total+gw <= maxWidthWithoutTail {
			pos = g.End()
		}
		total += gw
		if total > maxWidth {
			if options.ControlSequences {
				// Build result with trailing 7-bit ANSI escape sequences preserved
				result := make([]byte, 0, len(s)+len(tail)) // at most original + tail
				result = append(result, s[:pos]...)
				result = append(result, tail...)

				rem := graphemes.FromBytes(s[pos:])
				rem.AnsiEscapeSequences = options.ControlSequences

				for rem.Next() {
					v := rem.Value()
					// Only preserve 7-bit escapes (ESC = 0x1B) that measure
					// as zero-width on their own; some sequences (e.g. SOS)
					// are only valid in their original context.
					if len(v) > 0 && v[0] == 0x1B && options.Bytes(v) == 0 {
						result = append(result, v...)
					}
				}
				return result
			}
			result := make([]byte, 0, pos+len(tail))
			result = append(result, s[:pos]...)
			result = append(result, tail...)
			return result
		}
	}
	// No truncation
	return s
}

// TruncateBytes truncates a []byte to the given maxWidth, and appends the
// given tail if the []byte is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func TruncateBytes(s []byte, maxWidth int, tail []byte) []byte {
	return DefaultOptions.TruncateBytes(s, maxWidth, tail)
}
