package runewidth

import (
	"os"

	"github.com/rivo/uniseg"
)

//go:generate go run script/generate.go

var (
	// EastAsianWidth will be set true if the current locale is CJK
	EastAsianWidth bool

	// StrictEmojiNeutral should be set false if handle broken fonts
	StrictEmojiNeutral bool = true

	// DefaultCondition is a condition in current locale
	DefaultCondition = &Condition{
		EastAsianWidth:     false,
		StrictEmojiNeutral: true,
	}
)

func init() {
	handleEnv()
}

func handleEnv() {
	env := os.Getenv("RUNEWIDTH_EASTASIAN")
	if env == "" {
		EastAsianWidth = IsEastAsian()
	} else {
		EastAsianWidth = env == "1"
	}
	// update DefaultCondition
	DefaultCondition.EastAsianWidth = EastAsianWidth
}

type interval struct {
	first rune
	last  rune
}

type table []interval

func inTables(r rune, ts ...table) bool {
	for _, t := range ts {
		if inTable(r, t) {
			return true
		}
	}
	return false
}

func inTable(r rune, t table) bool {
	if r < t[0].first {
		return false
	}

	bot := 0
	top := len(t) - 1
	for top >= bot {
		mid := (bot + top) >> 1

		switch {
		case t[mid].last < r:
			bot = mid + 1
		case t[mid].first > r:
			top = mid - 1
		default:
			return true
		}
	}

	return false
}

var private = table{
	{0x00E000, 0x00F8FF}, {0x0F0000, 0x0FFFFD}, {0x100000, 0x10FFFD},
}

var nonprint = table{
	{0x0000, 0x001F}, {0x007F, 0x009F}, {0x00AD, 0x00AD},
	{0x070F, 0x070F}, {0x180B, 0x180E}, {0x200B, 0x200F},
	{0x2028, 0x202E}, {0x206A, 0x206F}, {0xD800, 0xDFFF},
	{0xFEFF, 0xFEFF}, {0xFFF9, 0xFFFB}, {0xFFFE, 0xFFFF},
}

// Condition have flag EastAsianWidth whether the current locale is CJK or not.
type Condition struct {
	EastAsianWidth     bool
	StrictEmojiNeutral bool
}

// NewCondition return new instance of Condition which is current locale.
func NewCondition() *Condition {
	return &Condition{
		EastAsianWidth:     EastAsianWidth,
		StrictEmojiNeutral: StrictEmojiNeutral,
	}
}

// RuneWidth returns the number of cells in r.
// See http://www.unicode.org/reports/tr11/
func (c *Condition) RuneWidth(r rune) int {
	// optimized version, verified by TestRuneWidthChecksums()
	if !c.EastAsianWidth {
		switch {
		case r < 0x20 || r > 0x10FFFF:
			return 0
		case (r >= 0x7F && r <= 0x9F) || r == 0xAD: // nonprint
			return 0
		case r < 0x300:
			return 1
		case inTable(r, narrow):
			return 1
		case inTables(r, nonprint, combining):
			return 0
		case inTable(r, doublewidth):
			return 2
		default:
			return 1
		}
	} else {
		switch {
		case r < 0 || r > 0x10FFFF || inTables(r, nonprint, combining):
			return 0
		case inTable(r, narrow):
			return 1
		case inTables(r, ambiguous, doublewidth):
			return 2
		case !c.StrictEmojiNeutral && inTables(r, ambiguous, emoji, narrow):
			return 2
		default:
			return 1
		}
	}
}

// StringWidth return width as you can see
func (c *Condition) StringWidth(s string) (width int) {
	g := uniseg.NewGraphemes(s)
	for g.Next() {
		var chWidth int
		for _, r := range g.Runes() {
			chWidth = c.RuneWidth(r)
			if chWidth > 0 {
				break // Our best guess at this point is to use the width of the first non-zero-width rune.
			}
		}
		width += chWidth
	}
	return
}

// Truncate return string truncated with w cells
func (c *Condition) Truncate(s string, w int, tail string) string {
	if c.StringWidth(s) <= w {
		return s
	}
	w -= c.StringWidth(tail)
	var width int
	pos := len(s)
	g := uniseg.NewGraphemes(s)
	for g.Next() {
		var chWidth int
		for _, r := range g.Runes() {
			chWidth = c.RuneWidth(r)
			if chWidth > 0 {
				break // See StringWidth() for details.
			}
		}
		if width+chWidth > w {
			pos, _ = g.Positions()
			break
		}
		width += chWidth
	}
	return s[:pos] + tail
}

// Wrap return string wrapped with w cells
func (c *Condition) Wrap(s string, w int) string {
	width := 0
	out := ""
	for _, r := range []rune(s) {
		cw := c.RuneWidth(r)
		if r == '\n' {
			out += string(r)
			width = 0
			continue
		} else if width+cw > w {
			out += "\n"
			width = 0
			out += string(r)
			width += cw
			continue
		}
		out += string(r)
		width += cw
	}
	return out
}

// FillLeft return string filled in left by spaces in w cells
func (c *Condition) FillLeft(s string, w int) string {
	width := c.StringWidth(s)
	count := w - width
	if count > 0 {
		b := make([]byte, count)
		for i := range b {
			b[i] = ' '
		}
		return string(b) + s
	}
	return s
}

// FillRight return string filled in left by spaces in w cells
func (c *Condition) FillRight(s string, w int) string {
	width := c.StringWidth(s)
	count := w - width
	if count > 0 {
		b := make([]byte, count)
		for i := range b {
			b[i] = ' '
		}
		return s + string(b)
	}
	return s
}

// RuneWidth returns the number of cells in r.
// See http://www.unicode.org/reports/tr11/
func RuneWidth(r rune) int {
	return DefaultCondition.RuneWidth(r)
}

// IsAmbiguousWidth returns whether is ambiguous width or not.
func IsAmbiguousWidth(r rune) bool {
	return inTables(r, private, ambiguous)
}

// IsNeutralWidth returns whether is neutral width or not.
func IsNeutralWidth(r rune) bool {
	return inTable(r, neutral)
}

// StringWidth return width as you can see
func StringWidth(s string) (width int) {
	return DefaultCondition.StringWidth(s)
}

// Truncate return string truncated with w cells
func Truncate(s string, w int, tail string) string {
	return DefaultCondition.Truncate(s, w, tail)
}

// Wrap return string wrapped with w cells
func Wrap(s string, w int) string {
	return DefaultCondition.Wrap(s, w)
}

// FillLeft return string filled in left by spaces in w cells
func FillLeft(s string, w int) string {
	return DefaultCondition.FillLeft(s, w)
}

// FillRight return string filled in left by spaces in w cells
func FillRight(s string, w int) string {
	return DefaultCondition.FillRight(s, w)
}
