package ansi

import (
	"bytes"

	"github.com/mattn/go-runewidth"
)

// Buffer is a buffer aware of ANSI escape sequences.
type Buffer struct {
	bytes.Buffer
}

// PrintableRuneWidth returns the cell width of all printable runes in the
// buffer.
func (w Buffer) PrintableRuneWidth() int {
	return PrintableRuneWidth(w.String())
}

// PrintableRuneWidth returns the cell width of the given string.
func PrintableRuneWidth(s string) int {
	var n int
	var ansi bool

	for _, c := range s {
		if c == Marker {
			// ANSI escape sequence
			ansi = true
		} else if ansi {
			if IsTerminator(c) {
				// ANSI sequence terminated
				ansi = false
			}
		} else {
			n += runewidth.RuneWidth(c)
		}
	}

	return n
}
