package gocui

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTextArea(content string, cursor int) *TextArea {
	textarea := &TextArea{}
	textarea.TypeString(content)
	textarea.cursor = cursor
	return textarea
}

func TestVimWordMotions(t *testing.T) {
	tests := []struct {
		content  string
		cursor   int
		motion   func(*TextArea, int, int) int
		count    int
		expected int
	}{
		{"foo bar baz", 0, (*TextArea).vimWordForward, 1, 4},
		{"foo bar baz", 0, (*TextArea).vimWordForward, 2, 8},
		{"foo bar baz", 4, (*TextArea).vimWordForward, 5, 11},
		{"foo.bar", 0, (*TextArea).vimWordForward, 1, 3},
		{"foo.bar", 3, (*TextArea).vimWordForward, 1, 4},
		{"foo\nbar", 0, (*TextArea).vimWordForward, 1, 4},
		{"foo  \n  bar", 0, (*TextArea).vimWordForward, 1, 8},
		{"字字 foo", 0, (*TextArea).vimWordForward, 1, 7},
		{"", 0, (*TextArea).vimWordForward, 1, 0},

		{"foo bar baz", 0, (*TextArea).vimWordEnd, 1, 2},
		{"foo bar baz", 2, (*TextArea).vimWordEnd, 1, 6},
		{"foo bar baz", 0, (*TextArea).vimWordEnd, 2, 6},
		{"foo.bar", 0, (*TextArea).vimWordEnd, 2, 3},
		{"foo\nbar", 2, (*TextArea).vimWordEnd, 1, 6},
		{"foo", 2, (*TextArea).vimWordEnd, 1, 2},

		{"foo bar baz", 8, (*TextArea).vimWordBack, 1, 4},
		{"foo bar baz", 8, (*TextArea).vimWordBack, 2, 0},
		{"foo bar", 6, (*TextArea).vimWordBack, 1, 4},
		{"foo.bar", 4, (*TextArea).vimWordBack, 1, 3},
		{"foo\nbar", 4, (*TextArea).vimWordBack, 1, 0},
		{"foo bar", 0, (*TextArea).vimWordBack, 1, 0},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			textarea := makeTextArea(test.content, test.cursor)
			assert.Equal(t, test.expected, test.motion(textarea, test.cursor, test.count))
		})
	}
}

func TestVimLineHelpers(t *testing.T) {
	textarea := makeTextArea("foo\n  bar baz\n\nqux", 0)

	assert.Equal(t, 0, textarea.hardLineStart(2))
	assert.Equal(t, 3, textarea.hardLineEnd(2))
	assert.Equal(t, 4, textarea.hardLineStart(8))
	assert.Equal(t, 13, textarea.hardLineEnd(8))
	assert.Equal(t, 14, textarea.hardLineStart(14))
	assert.Equal(t, 14, textarea.hardLineEnd(14))

	assert.Equal(t, 6, textarea.vimFirstNonBlank(4))
	assert.Equal(t, 0, textarea.vimFirstNonBlank(1))
	assert.Equal(t, 14, textarea.vimFirstNonBlank(14))

	assert.Equal(t, 12, textarea.vimLastCharOfLine(4))
	assert.Equal(t, 14, textarea.vimLastCharOfLine(14))

	assert.Equal(t, 0, textarea.vimLineNumber(2))
	assert.Equal(t, 1, textarea.vimLineNumber(8))
	assert.Equal(t, 4, textarea.vimLineCount())
	assert.Equal(t, 4, textarea.vimGotoLine(1))
	assert.Equal(t, 14, textarea.vimGotoLine(2))
	assert.Equal(t, 15, textarea.vimGotoLine(3))
	assert.Equal(t, 15, textarea.vimGotoLine(99))
}

func TestVimFindOnLine(t *testing.T) {
	tests := []struct {
		content  string
		cursor   int
		target   string
		forward  bool
		till     bool
		count    int
		expected int
	}{
		{"abcabc", 0, "c", true, false, 1, 2},
		{"abcabc", 0, "c", true, false, 2, 5},
		{"abcabc", 0, "c", true, true, 1, 1},
		{"abcabc", 5, "a", false, false, 1, 3},
		{"abcabc", 5, "a", false, true, 1, 4},
		{"abcabc", 0, "z", true, false, 1, -1},
		{"abc\ncde", 0, "d", true, false, 1, -1},
		{"abc\ncde", 5, "a", false, false, 1, -1},
		{"abcabc", 2, "c", true, false, 1, 5},
		{"a字c", 0, "c", true, false, 1, 4},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			textarea := makeTextArea(test.content, test.cursor)
			actual := textarea.vimFindOnLine(test.cursor, test.target, test.forward, test.till, test.count)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestVimWordBounds(t *testing.T) {
	tests := []struct {
		content       string
		cursor        int
		around        bool
		expectedStart int
		expectedEnd   int
	}{
		{"foo bar baz", 5, false, 4, 7},
		{"foo bar baz", 5, true, 4, 8},
		{"foo bar", 5, true, 3, 7},
		{"foo bar baz", 3, false, 3, 4},
		{"foo.bar", 1, false, 0, 3},
		{"foo.bar", 3, false, 3, 4},
		{"foo.bar", 3, true, 3, 4},
		{"one\ntwo\nthree", 5, false, 4, 7},
		{"one\ntwo", 3, false, 3, 3},
		{"字字 foo", 1, false, 0, 6},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			textarea := makeTextArea(test.content, test.cursor)
			start, end := textarea.vimWordBounds(test.cursor, test.around)
			assert.Equal(t, test.expectedStart, start)
			assert.Equal(t, test.expectedEnd, end)
		})
	}
}

func TestVimRangeEdits(t *testing.T) {
	textarea := makeTextArea("foo bar baz", 0)
	assert.Equal(t, "bar", textarea.getRange(4, 7))
	assert.Equal(t, "", textarea.getRange(7, 4))
	assert.Equal(t, "baz", textarea.getRange(8, 99))

	textarea.deleteRange(4, 8)
	assert.Equal(t, "foo baz", textarea.GetUnwrappedContent())
	assert.Equal(t, 4, textarea.cursor)

	textarea.insertAt(4, "bar ")
	assert.Equal(t, "foo bar baz", textarea.GetUnwrappedContent())
	assert.Equal(t, 8, textarea.cursor)

	textarea.replaceGraphemeAt(0, "g")
	assert.Equal(t, "goo bar baz", textarea.GetUnwrappedContent())
	assert.Equal(t, 0, textarea.cursor)

	textarea = makeTextArea("a\nb", 1)
	textarea.replaceGraphemeAt(1, "x")
	assert.Equal(t, "a\nb", textarea.GetUnwrappedContent())

	textarea = makeTextArea("字c", 0)
	textarea.replaceGraphemeAt(0, "a")
	assert.Equal(t, "ac", textarea.GetUnwrappedContent())
}

func TestVimGraphemeSteps(t *testing.T) {
	textarea := makeTextArea("a字c", 0)
	assert.Equal(t, 1, textarea.nextGraphemeStart(0))
	assert.Equal(t, 4, textarea.nextGraphemeStart(1))
	assert.Equal(t, 5, textarea.nextGraphemeStart(5))
	assert.Equal(t, 4, textarea.prevGraphemeStart(5))
	assert.Equal(t, 1, textarea.prevGraphemeStart(4))
	assert.Equal(t, 0, textarea.prevGraphemeStart(1))
	assert.Equal(t, "字", textarea.graphemeAt(1))
	assert.Equal(t, "", textarea.graphemeAt(5))
}

func TestVimWordMotionsWithSoftWrap(t *testing.T) {
	textarea := &TextArea{AutoWrap: true, AutoWrapWidth: 5}
	textarea.TypeString("aaa bbb ccc")

	assert.Equal(t, 4, textarea.vimWordForward(0, 1))
	assert.Equal(t, 8, textarea.vimWordForward(4, 1))
	assert.Equal(t, 4, textarea.vimWordBack(8, 1))
	assert.Equal(t, 6, textarea.vimWordEnd(4, 1))
	assert.Equal(t, 0, textarea.hardLineStart(8))
	assert.Equal(t, 11, textarea.hardLineEnd(3))
}
