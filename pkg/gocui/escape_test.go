package gocui

import (
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

func parseEscRunes(t *testing.T, ei *escapeInterpreter, runes string) {
	t.Helper()
	for _, b := range []byte(runes) {
		isEscape, err := ei.parseOne([]byte{b})
		assert.Equal(t, true, isEscape)
		assert.NoError(t, err)
	}
}
