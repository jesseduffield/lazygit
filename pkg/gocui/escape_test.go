package gocui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOne(t *testing.T) {
	var ei *escapeInterpreter

	ei = newEscapeInterpreter(OutputNormal)
	isEscape, err := ei.parseOne([]byte{'a'})
	assert.Equal(t, false, isEscape)
	assert.NoError(t, err)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, "\x1b[0K")
	_, ok := ei.instruction.(eraseInLineFromCursor)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, "\x1b[K")
	_, ok = ei.instruction.(eraseInLineFromCursor)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, "\x1b[1K")
	_, ok = ei.instruction.(noInstruction)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, "\x1b(B")
	_, ok = ei.instruction.(noInstruction)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, "\x1b)0")
	_, ok = ei.instruction.(noInstruction)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, "\x1b*A")
	_, ok = ei.instruction.(noInstruction)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, "\x1b+K")
	_, ok = ei.instruction.(noInstruction)
	assert.Equal(t, true, ok)
}

func TestParseOneColours(t *testing.T) {
	scenarios := []struct {
		outputMode OutputMode
		input      string
		expectedFg Attribute
		expectedBg Attribute
	}{
		{OutputNormal, "\x1b[30m", ColorBlack, ColorDefault},
		{OutputNormal, "\x1b[31m", ColorRed, ColorDefault},
		{OutputNormal, "\x1b[32m", ColorGreen, ColorDefault},
		{OutputNormal, "\x1b[33m", ColorYellow, ColorDefault},
		{OutputNormal, "\x1b[34m", ColorBlue, ColorDefault},
		{OutputNormal, "\x1b[35m", ColorMagenta, ColorDefault},
		{OutputNormal, "\x1b[36m", ColorCyan, ColorDefault},
		{OutputNormal, "\x1b[37m", ColorWhite, ColorDefault},
		{OutputNormal, "\x1b[40m", ColorDefault, ColorBlack},
		{OutputNormal, "\x1b[41m", ColorDefault, ColorRed},
		{OutputNormal, "\x1b[42m", ColorDefault, ColorGreen},
		{OutputNormal, "\x1b[43m", ColorDefault, ColorYellow},
		{OutputNormal, "\x1b[44m", ColorDefault, ColorBlue},
		{OutputNormal, "\x1b[45m", ColorDefault, ColorMagenta},
		{OutputNormal, "\x1b[46m", ColorDefault, ColorCyan},
		{OutputNormal, "\x1b[47m", ColorDefault, ColorWhite},
		{OutputNormal, "\x1b[47;31m", ColorRed, ColorWhite},
		{OutputNormal, "\x1b[90m", Get256Color(8), ColorDefault},
		{OutputNormal, "\x1b[91m", Get256Color(9), ColorDefault},
		{OutputNormal, "\x1b[92m", Get256Color(10), ColorDefault},
		{OutputNormal, "\x1b[93m", Get256Color(11), ColorDefault},
		{OutputNormal, "\x1b[94m", Get256Color(12), ColorDefault},
		{OutputNormal, "\x1b[95m", Get256Color(13), ColorDefault},
		{OutputNormal, "\x1b[96m", Get256Color(14), ColorDefault},
		{OutputNormal, "\x1b[97m", Get256Color(15), ColorDefault},
		{OutputNormal, "\x1b[100m", ColorDefault, Get256Color(8)},
		{OutputNormal, "\x1b[101m", ColorDefault, Get256Color(9)},
		{OutputNormal, "\x1b[102m", ColorDefault, Get256Color(10)},
		{OutputNormal, "\x1b[103m", ColorDefault, Get256Color(11)},
		{OutputNormal, "\x1b[104m", ColorDefault, Get256Color(12)},
		{OutputNormal, "\x1b[105m", ColorDefault, Get256Color(13)},
		{OutputNormal, "\x1b[106m", ColorDefault, Get256Color(14)},
		{OutputNormal, "\x1b[107m", ColorDefault, Get256Color(15)},
		{Output256, "\x1b[38;5;32m", Get256Color(32), ColorDefault},
		{OutputTrue, "\x1b[38;5;32m", Get256Color(32), ColorDefault},
		{OutputTrue, "\x1b[38;2;50;103;205m", NewRGBColor(50, 103, 205), ColorDefault},
		{Output256, "\x1b[48;5;32m", ColorDefault, Get256Color(32)},
		{OutputTrue, "\x1b[48;5;32m", ColorDefault, Get256Color(32)},
		{OutputTrue, "\x1b[48;2;50;103;205m", ColorDefault, NewRGBColor(50, 103, 205)},
		{OutputTrue, "\x1b[1;95;48;2;255;224;224m", Get256Color(13), NewRGBColor(255, 224, 224)},
	}

	for _, scenario := range scenarios {
		ei := newEscapeInterpreter(scenario.outputMode)
		parseEscRunes(t, ei, scenario.input)
		assert.Equal(t, scenario.expectedFg, ei.curFgColor&AttrColorBits)
		assert.Equal(t, scenario.expectedBg, ei.curBgColor)
	}

	// resetting colours
	scenarios = []struct {
		outputMode OutputMode
		input      string
		expectedFg Attribute
		expectedBg Attribute
	}{
		{OutputNormal, "\x1b[39m", ColorDefault, ColorRed},
		{OutputNormal, "\x1b[49m", ColorRed, ColorDefault},
		{OutputNormal, "\x1b[0m", ColorDefault, ColorDefault},
	}

	for _, scenario := range scenarios {
		ei := newEscapeInterpreter(scenario.outputMode)
		ei.curFgColor = ColorRed
		ei.curBgColor = ColorRed
		parseEscRunes(t, ei, scenario.input)
		assert.Equal(t, scenario.expectedFg, ei.curFgColor)
		assert.Equal(t, scenario.expectedBg, ei.curBgColor)
	}

	// setting attributes
	attrScenarios := []struct {
		outputMode   OutputMode
		input        string
		expectedAttr Attribute
	}{
		{OutputNormal, "\x1b[1m", AttrBold},
		{OutputNormal, "\x1b[2m", AttrDim},
		{OutputNormal, "\x1b[3m", AttrItalic},
		{OutputNormal, "\x1b[4m", AttrUnderline},
		{OutputNormal, "\x1b[5m", AttrBlink},
		{OutputNormal, "\x1b[7m", AttrReverse},
		{OutputNormal, "\x1b[9m", AttrStrikeThrough},
	}

	for _, scenario := range attrScenarios {
		ei := newEscapeInterpreter(scenario.outputMode)
		parseEscRunes(t, ei, scenario.input)
		style := ei.curFgColor & AttrStyleBits
		assert.Equal(t, scenario.expectedAttr, style)
	}
}

func TestParseOneIgnoresUnknownSequences(t *testing.T) {
	// Escape sequences the interpreter doesn't implement -- whether well-formed-but-unsupported
	// (private modes, DECSCUSR, …) or outright malformed -- must be silently
	// consumed rather than leaked into the view as literal text.
	scenarios := []string{
		"\x1b[?9001h",                            // DEC private-mode set (?-prefix)
		"\x1b[?25l",                              // hide cursor
		"\x1b[?25h",                              // show cursor
		"\x1b[2;J",                               // erase display (unusual 2;J variant)
		"\x1b[H",                                 // cursor home — re-anchors to row 1 (no-op when already there)
		"\x1bc",                                  // RIS — single-char ESC sequence
		"\x1b[;5H",                               // empty first param — defaults to row 1, no-op
		"\x1b[ q",                                // intermediate byte with no params (DECSCUSR family)
		"\x1b[0 q",                               // intermediate byte after a param
		"\x1b[1;;m",                              // malformed SGR: empty middle param
		"\x1b]8bogus\x07",                        // OSC 8 missing ';'
		"\x1b[" + strings.Repeat("0", 300) + "m", // single param overflows length cap
		"\x1b[" + strings.Repeat("1;", 25) + "1m", // too many params
	}

	for _, input := range scenarios {
		ei := newEscapeInterpreter(OutputNormal)
		parseEscRunes(t, ei, input)
		// An unimplemented/malformed sequence must leave no trace: no
		// pending instruction, no color change.
		_, noop := ei.instruction.(noInstruction)
		assert.True(t, noop, "input %q left a pending instruction", input)
		assert.Equal(t, ColorDefault, ei.curFgColor, "input %q mutated fg color", input)
		assert.Equal(t, ColorDefault, ei.curBgColor, "input %q mutated bg color", input)
	}
}

func TestParseOneCursorPositioning(t *testing.T) {
	// Cursor-positioning escapes that advance the row forward emit a
	// cursorDown instruction; backward / same-row moves are ignored
	// because the view's buffer is line-based.
	scenarios := []struct {
		input       string
		startRow    int // parser's screenRow before parsing
		wantAdvance int // 0 means "no instruction emitted"
	}{
		{"\x1b[5;1H", 1, 4}, // CUP — absolute row 5 from row 1
		{"\x1b[5H", 1, 4},   // CUP with only the row param
		{"\x1b[5;1H", 5, 0}, // CUP to the same row we're on — no-op
		{"\x1b[2;1H", 5, 0}, // CUP backward — ignored
		{"\x1b[5;1f", 1, 4}, // HVP alias for CUP
		{"\x1b[5d", 1, 4},   // VPA — absolute row
		{"\x1b[2d", 5, 0},   // VPA backward — ignored
		{"\x1b[3B", 1, 3},   // CUD — relative
		{"\x1b[B", 1, 1},    // CUD with default param of 1
		{"\x1b[2E", 1, 2},   // CNL — relative
	}

	for _, s := range scenarios {
		ei := newEscapeInterpreter(OutputNormal)
		ei.screenRow = s.startRow
		parseEscRunes(t, ei, s.input)
		if s.wantAdvance == 0 {
			_, noop := ei.instruction.(noInstruction)
			assert.True(t, noop, "input %q at row %d should be a no-op", s.input, s.startRow)
		} else {
			cd, ok := ei.instruction.(cursorDown)
			if assert.True(t, ok, "input %q at row %d should emit cursorDown", s.input, s.startRow) {
				assert.Equal(t, s.wantAdvance, cd.n, "input %q at row %d", s.input, s.startRow)
			}
		}
	}
}

func TestParseOneCursorHomeReanchors(t *testing.T) {
	// ConPTY emits cursor-home ([H) after [2J at the start of every screen.
	// In a view that isn't rewound in lockstep with ConPTY (the command log)
	// screenRow has drifted, so home must re-anchor it to the current write
	// position rather than be dropped as a backward move — otherwise the
	// absolute CUPs that follow compute negative, dropped advances and the
	// rows ConPTY positioned with collapse together.
	ei := newEscapeInterpreter(OutputNormal)
	ei.screenRow = 12 // accumulated drift from earlier command-log output

	parseEscRunes(t, ei, "\x1b[H")
	assert.Equal(t, 1, ei.screenRow, "home should re-anchor screenRow")
	_, noop := ei.instruction.(noInstruction)
	assert.True(t, noop, "home should not emit an instruction")

	// A subsequent CUP now advances relative to the re-anchored origin.
	parseEscRunes(t, ei, "\x1b[3;1H")
	cd, ok := ei.instruction.(cursorDown)
	if assert.True(t, ok, "CUP after home should emit cursorDown") {
		assert.Equal(t, 2, cd.n)
	}
}

func TestParseOneCursorForward(t *testing.T) {
	// CUF (\x1b[NC) emits a cursorForward instruction so the view can
	// materialize the N-cell gap as spaces. ConPTY uses this (often
	// paired with ECH) to encode runs of default-colored spaces.
	scenarios := []struct {
		input string
		wantN int
	}{
		{"\x1b[5C", 5},
		{"\x1b[1C", 1},
		{"\x1b[C", 1}, // no param defaults to 1
	}

	for _, s := range scenarios {
		ei := newEscapeInterpreter(OutputNormal)
		parseEscRunes(t, ei, s.input)
		cf, ok := ei.instruction.(cursorForward)
		if assert.True(t, ok, "input %q should emit cursorForward", s.input) {
			assert.Equal(t, s.wantN, cf.n, "input %q", s.input)
		}
	}
}

func parseEscRunes(t *testing.T, ei *escapeInterpreter, runes string) {
	t.Helper()
	for _, b := range []byte(runes) {
		isEscape, err := ei.parseOne([]byte{b})
		assert.Equal(t, true, isEscape)
		assert.NoError(t, err)
	}
}
