package wordwrap

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/muesli/reflow/ansi"
)

var (
	defaultBreakpoints = []rune{'-'}
	defaultNewline     = []rune{'\n'}
)

// WordWrap contains settings and state for customisable text reflowing with
// support for ANSI escape sequences. This means you can style your terminal
// output without affecting the word wrapping algorithm.
type WordWrap struct {
	Limit        int
	Breakpoints  []rune
	Newline      []rune
	KeepNewlines bool

	buf   bytes.Buffer
	space bytes.Buffer
	word  ansi.Buffer

	lineLen int
	ansi    bool
}

// NewWriter returns a new instance of a word-wrapping writer, initialized with
// default settings.
func NewWriter(limit int) *WordWrap {
	return &WordWrap{
		Limit:        limit,
		Breakpoints:  defaultBreakpoints,
		Newline:      defaultNewline,
		KeepNewlines: true,
	}
}

// Bytes is shorthand for declaring a new default WordWrap instance,
// used to immediately word-wrap a byte slice.
func Bytes(b []byte, limit int) []byte {
	f := NewWriter(limit)
	_, _ = f.Write(b)
	_ = f.Close()

	return f.Bytes()
}

// String is shorthand for declaring a new default WordWrap instance,
// used to immediately word-wrap a string.
func String(s string, limit int) string {
	return string(Bytes([]byte(s), limit))
}

func (w *WordWrap) addSpace() {
	w.lineLen += w.space.Len()
	_, _ = w.buf.Write(w.space.Bytes())
	w.space.Reset()
}

func (w *WordWrap) addWord() {
	if w.word.Len() > 0 {
		w.addSpace()
		w.lineLen += w.word.PrintableRuneWidth()
		_, _ = w.buf.Write(w.word.Bytes())
		w.word.Reset()
	}
}

func (w *WordWrap) addNewLine() {
	_, _ = w.buf.WriteRune('\n')
	w.lineLen = 0
	w.space.Reset()
}

func inGroup(a []rune, c rune) bool {
	for _, v := range a {
		if v == c {
			return true
		}
	}
	return false
}

// Write is used to write more content to the word-wrap buffer.
func (w *WordWrap) Write(b []byte) (int, error) {
	if w.Limit == 0 {
		return w.buf.Write(b)
	}

	s := string(b)
	if !w.KeepNewlines {
		s = strings.Replace(strings.TrimSpace(s), "\n", " ", -1)
	}

	for _, c := range s {
		if c == '\x1B' {
			// ANSI escape sequence
			_, _ = w.word.WriteRune(c)
			w.ansi = true
		} else if w.ansi {
			_, _ = w.word.WriteRune(c)
			if (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) {
				// ANSI sequence terminated
				w.ansi = false
			}
		} else if inGroup(w.Newline, c) {
			// end of current line
			// see if we can add the content of the space buffer to the current line
			if w.word.Len() == 0 {
				if w.lineLen+w.space.Len() > w.Limit {
					w.lineLen = 0
				} else {
					// preserve whitespace
					_, _ = w.buf.Write(w.space.Bytes())
				}
				w.space.Reset()
			}

			w.addWord()
			w.addNewLine()
		} else if unicode.IsSpace(c) {
			// end of current word
			w.addWord()
			_, _ = w.space.WriteRune(c)
		} else if inGroup(w.Breakpoints, c) {
			// valid breakpoint
			w.addSpace()
			w.addWord()
			_, _ = w.buf.WriteRune(c)
		} else {
			// any other character
			_, _ = w.word.WriteRune(c)

			// add a line break if the current word would exceed the line's
			// character limit
			if w.lineLen+w.space.Len()+w.word.PrintableRuneWidth() > w.Limit &&
				w.word.PrintableRuneWidth() < w.Limit {
				w.addNewLine()
			}
		}
	}

	return len(b), nil
}

// Close will finish the word-wrap operation. Always call it before trying to
// retrieve the final result.
func (w *WordWrap) Close() error {
	w.addWord()
	return nil
}

// Bytes returns the word-wrapped result as a byte slice.
func (w *WordWrap) Bytes() []byte {
	return w.buf.Bytes()
}

// String returns the word-wrapped result as a string.
func (w *WordWrap) String() string {
	return w.buf.String()
}
