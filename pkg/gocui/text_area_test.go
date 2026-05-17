package gocui

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextArea(t *testing.T) {
	tests := []struct {
		actions           func(*TextArea)
		wrapWidth         int
		expectedContent   string
		expectedCursor    int
		expectedClipboard string
	}{
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("c")
			},
			expectedContent:   "abc",
			expectedCursor:    3,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("\n")
				textarea.TypeCharacter("c")
			},
			expectedContent:   "a\nc",
			expectedCursor:    3,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abcd")
			},
			expectedContent:   "abcd",
			expectedCursor:    4,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("a字cd")
			},
			expectedContent:   "a字cd",
			expectedCursor:    6,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.BackSpaceChar()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.BackSpaceChar()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.BackSpaceChar()
			},
			expectedContent:   "a",
			expectedCursor:    1,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.DeleteChar()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.DeleteChar()
			},
			expectedContent:   "a",
			expectedCursor:    1,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.MoveCursorLeft()
				textarea.DeleteChar()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("c")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.DeleteChar()
			},
			expectedContent:   "ac",
			expectedCursor:    1,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.MoveCursorLeft()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.MoveCursorLeft()
			},
			expectedContent:   "a",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.MoveCursorLeft()
			},
			expectedContent:   "ab",
			expectedCursor:    1,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.MoveCursorRight()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.MoveCursorRight()
			},
			expectedContent:   "a",
			expectedCursor:    1,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.MoveCursorLeft()
				textarea.MoveCursorRight()
			},
			expectedContent:   "ab",
			expectedCursor:    2,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("漢")
				textarea.TypeCharacter("字")
				textarea.MoveCursorLeft()
			},
			expectedContent:   "漢字",
			expectedCursor:    3,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.ToggleOverwrite()
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
			},
			expectedContent:   "ab",
			expectedCursor:    2,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("c")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.ToggleOverwrite()
				textarea.TypeCharacter("d")
			},
			expectedContent:   "adc",
			expectedCursor:    2,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb")
				textarea.MoveLeftWord()
			},
			expectedContent: "aaa bbb",
			expectedCursor:  4,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa\nbbb")
				textarea.MoveLeftWord()
			},
			expectedContent: "aaa\nbbb",
			expectedCursor:  4,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb")
				textarea.GoToStartOfLine()
				textarea.MoveLeftWord()
			},
			wrapWidth:       4,
			expectedContent: "aaa bbb",
			expectedCursor:  0,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb\n")
				textarea.MoveLeftWord()
			},
			expectedContent: "aaa bbb\n",
			expectedCursor:  7,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb")
				textarea.MoveLeftWord()
				textarea.MoveLeftWord()
			},
			expectedContent: "aaa bbb",
			expectedCursor:  0,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa")
				textarea.GoToStartOfLine()
				textarea.MoveLeftWord()
			},
			expectedContent: "aaa",
			expectedCursor:  0,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb")
				textarea.MoveRightWord()
			},
			expectedContent: "aaa bbb",
			expectedCursor:  7,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa\nbbb")
				textarea.GoToStartOfLine()
				textarea.MoveCursorLeft()
				textarea.MoveRightWord()
			},
			expectedContent: "aaa\nbbb",
			expectedCursor:  4,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb")
				textarea.GoToStartOfLine()
				textarea.MoveCursorLeft()
				textarea.MoveRightWord()
			},
			wrapWidth:       4,
			expectedContent: "aaa bbb",
			expectedCursor:  7,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb")
				textarea.GoToStartOfLine()
				textarea.MoveRightWord()
			},
			expectedContent: "aaa bbb",
			expectedCursor:  3,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("aaa bbb\n")
				textarea.MoveCursorLeft()
				textarea.GoToStartOfLine()
				textarea.MoveRightWord()
				textarea.MoveRightWord()
			},
			expectedContent: "aaa bbb\n",
			expectedCursor:  7,
		},
		{
			actions: func(textarea *TextArea) {
				// overwrite mode acts same as normal mode when cursor is at the end
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("c")
				textarea.ToggleOverwrite()
				textarea.TypeCharacter("d")
			},
			expectedContent:   "abcd",
			expectedCursor:    4,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.DeleteToStartOfLine()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.DeleteToStartOfLine()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "ab",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.DeleteToStartOfLine()
			},
			expectedContent:   "ab",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("\n")
				textarea.DeleteToStartOfLine()
			},
			expectedContent:   "ab",
			expectedCursor:    2,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("\n")
				textarea.TypeCharacter("c")
				textarea.TypeCharacter("d")
				textarea.DeleteToStartOfLine()
			},
			expectedContent:   "ab\n",
			expectedCursor:    3,
			expectedClipboard: "cd",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.GoToStartOfLine()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.MoveCursorLeft()
				textarea.GoToStartOfLine()
			},
			expectedContent:   "a",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("\n")
				textarea.TypeCharacter("c")
				textarea.TypeCharacter("d")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.GoToStartOfLine()
			},
			expectedContent:   "ab\ncd",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("\n")
				textarea.TypeCharacter("c")
				textarea.TypeCharacter("d")
				textarea.GoToStartOfLine()
			},
			expectedContent:   "ab\ncd",
			expectedCursor:    3,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("\n")
				textarea.TypeCharacter("c")
				textarea.TypeCharacter("d")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.GoToStartOfLine()
			},
			expectedContent:   "ab\ncd",
			expectedCursor:    3,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.GoToEndOfLine()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("a")
				textarea.TypeCharacter("b")
				textarea.TypeCharacter("\n")
				textarea.TypeCharacter("c")
				textarea.TypeCharacter("d")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.GoToEndOfLine()
			},
			expectedContent:   "ab\ncd",
			expectedCursor:    5,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.SetCursor2D(10, 10)
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.SetCursor2D(-1, -1)
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\ncd")
				textarea.SetCursor2D(0, 0)
			},
			expectedContent:   "ab\ncd",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\ncd")
				textarea.SetCursor2D(2, 0)
			},
			expectedContent:   "ab\ncd",
			expectedCursor:    2,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\ncd\nef")
				textarea.SetCursor2D(2, 1)
			},
			expectedContent:   "ab\ncd\nef",
			expectedCursor:    5,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abcd\n\nijkl")
				textarea.MoveCursorUp()
			},
			expectedContent:   "abcd\n\nijkl",
			expectedCursor:    5,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abcdef\n老老老")
				textarea.MoveCursorLeft()
				textarea.MoveCursorUp()
			},
			expectedContent:   "abcdef\n老老老",
			expectedCursor:    4,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abcdef\n老老老")
				textarea.MoveCursorUp()
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.MoveCursorDown()
			},
			expectedContent:   "abcdef\n老老老",
			expectedCursor:    13,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abcd\nef")
				textarea.MoveCursorUp()
				textarea.GoToEndOfLine()
				textarea.MoveCursorDown()
			},
			expectedContent:   "abcd\nef",
			expectedCursor:    7,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abcd")
				textarea.MoveCursorUp()
			},
			expectedContent:   "abcd",
			expectedCursor:    4,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abcdefg`)
				textarea.Clear()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abcdefg`)
				textarea.Clear()
			},
			expectedContent:   "",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc def`)
				textarea.MoveCursorLeft()
				textarea.BackSpaceWord()
			},
			expectedContent:   "abc f",
			expectedCursor:    4,
			expectedClipboard: "de",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc  def   `)
				textarea.BackSpaceWord()
			},
			expectedContent:   "abc  ",
			expectedCursor:    5,
			expectedClipboard: "def   ",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abc def\nghi")
				textarea.BackSpaceWord()
			},
			expectedContent:   "abc def\n",
			expectedCursor:    8,
			expectedClipboard: "ghi",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abc def\nghi")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.BackSpaceWord()
			},
			expectedContent:   "abc defghi",
			expectedCursor:    7,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc(def)`)
				textarea.BackSpaceWord()
			},
			expectedContent:   "abc(def",
			expectedCursor:    7,
			expectedClipboard: ")",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc(def`)
				textarea.BackSpaceWord()
			},
			expectedContent:   "abc(",
			expectedCursor:    4,
			expectedClipboard: "def",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc`)
				textarea.GoToStartOfLine()
				textarea.BackSpaceWord()
			},
			expectedContent:   "abc",
			expectedCursor:    0,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc`)
				textarea.Yank()
			},
			expectedContent:   "abc",
			expectedCursor:    3,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc def`)
				textarea.DeleteToStartOfLine()
				textarea.Yank()
				textarea.Yank()
			},
			expectedContent:   "abc defabc def",
			expectedCursor:    14,
			expectedClipboard: "abc def",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abc\ndef")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.MoveCursorUp()
				textarea.DeleteToEndOfLine()
			},
			expectedContent:   "a\ndef",
			expectedCursor:    1,
			expectedClipboard: "bc",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abc\ndef")
				textarea.MoveCursorUp()
				textarea.DeleteToEndOfLine()
			},
			expectedContent:   "abcdef",
			expectedCursor:    3,
			expectedClipboard: "",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc def`)
				textarea.BackSpaceWord()
				textarea.Yank()
				textarea.Yank()
			},
			expectedContent:   "abc defdef",
			expectedCursor:    10,
			expectedClipboard: "def",
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString(`abc def`)
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.DeleteToEndOfLine()
				textarea.Yank()
				textarea.Yank()
			},
			expectedContent:   "abc defef",
			expectedCursor:    9,
			expectedClipboard: "ef",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			textarea := &TextArea{}
			if test.wrapWidth > 0 {
				textarea.AutoWrap = true
				textarea.AutoWrapWidth = test.wrapWidth
			}
			test.actions(textarea)
			assert.EqualValues(t, test.expectedContent, textarea.GetUnwrappedContent())
			assert.EqualValues(t, test.expectedCursor, textarea.cursor)
			assert.EqualValues(t, test.expectedClipboard, textarea.clipboard)
		})
	}
}

func TestGetCursorXY(t *testing.T) {
	tests := []struct {
		actions   func(*TextArea)
		wrapWidth int
		expectedX int
		expectedY int
	}{
		{
			actions: func(textarea *TextArea) {
				// do nothing
			},
			expectedX: 0,
			expectedY: 0,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("\n")
			},
			expectedX: 0,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("\na")
			},
			expectedX: 1,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("\n")
				textarea.MoveCursorUp()
			},
			expectedX: 0,
			expectedY: 0,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("\n\n")
				textarea.MoveCursorUp()
			},
			expectedX: 0,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\ncd")
			},
			expectedX: 2,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\n")
			},
			expectedX: 0,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\n")
				textarea.MoveCursorLeft()
			},
			expectedX: 2,
			expectedY: 0,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\n\n")
			},
			expectedX: 0,
			expectedY: 2,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("ab\n\n")
				textarea.MoveCursorLeft()
			},
			expectedX: 0,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeCharacter("漢")
				textarea.TypeCharacter("字")
			},
			expectedX: 4,
			expectedY: 0,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abc de")
				textarea.MoveCursorLeft()
			},
			wrapWidth: 4,
			expectedX: 1,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abc de")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
			},
			wrapWidth: 4,
			expectedX: 0,
			expectedY: 1,
		},
		{
			actions: func(textarea *TextArea) {
				textarea.TypeString("abc de")
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
				textarea.MoveCursorLeft()
			},
			wrapWidth: 4,
			expectedX: 3,
			expectedY: 0,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			textarea := &TextArea{}
			if test.wrapWidth > 0 {
				textarea.AutoWrap = true
				textarea.AutoWrapWidth = test.wrapWidth
			}
			test.actions(textarea)
			x, y := textarea.GetCursorXY()
			assert.EqualValues(t, test.expectedX, x)
			assert.EqualValues(t, test.expectedY, y)

			// As a sanity check, test that setting the cursor back to (x, y) results in the same cursor position:
			cursor := textarea.cursor
			textarea.SetCursor2D(x, y)
			assert.EqualValues(t, cursor, textarea.cursor)
		})
	}
}

func Test_AutoWrapContent(t *testing.T) {
	tests := []struct {
		name                   string
		content                string
		autoWrapWidth          int
		expectedWrappedContent string
		expectedSoftLineBreaks []int
	}{
		{
			name:                   "empty content",
			content:                "",
			autoWrapWidth:          7,
			expectedWrappedContent: "",
			expectedSoftLineBreaks: []int{},
		},
		{
			name:                   "no wrapping necessary",
			content:                "abcde",
			autoWrapWidth:          7,
			expectedWrappedContent: "abcde",
			expectedSoftLineBreaks: []int{},
		},
		{
			name:                   "wrap at whitespace",
			content:                "abcde xyz",
			autoWrapWidth:          7,
			expectedWrappedContent: "abcde \nxyz",
			expectedSoftLineBreaks: []int{6},
		},
		{
			name:                   "take wide characters into account",
			content:                "🏴󠁧󠁢󠁥󠁮󠁧󠁿 🏴󠁧󠁢󠁥󠁮󠁧󠁿 x y", // the flag has a width of 2
			autoWrapWidth:          7,
			expectedWrappedContent: "🏴󠁧󠁢󠁥󠁮󠁧󠁿 🏴󠁧󠁢󠁥󠁮󠁧󠁿 x \ny",
			expectedSoftLineBreaks: []int{60},
		},
		{
			name:                   "take wide characters into account at line end",
			content:                "🏴󠁧󠁢󠁥󠁮󠁧󠁿 🏴󠁧󠁢󠁥󠁮󠁧󠁿 🏴󠁧󠁢󠁥󠁮󠁧󠁿", // the flag has a width of 2
			autoWrapWidth:          7,
			expectedWrappedContent: "🏴󠁧󠁢󠁥󠁮󠁧󠁿 🏴󠁧󠁢󠁥󠁮󠁧󠁿 \n🏴󠁧󠁢󠁥󠁮󠁧󠁿",
			expectedSoftLineBreaks: []int{58},
		},
		{
			name:                   "lots of whitespace is preserved at end of line",
			content:                "abcde      xyz",
			autoWrapWidth:          7,
			expectedWrappedContent: "abcde      \nxyz",
			expectedSoftLineBreaks: []int{11},
		},
		{
			name:                   "don't wrap inside long word when there's no whitespace",
			content:                "abc defghijklmn opq",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc \ndefghijklmn \nopq",
			expectedSoftLineBreaks: []int{4, 16},
		},
		{
			name:                   "don't break at space after footnote symbol",
			content:                "abc\n[1]: https://long/link\ndef",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc\n[1]: https://long/link\ndef",
			expectedSoftLineBreaks: []int{},
		},
		{
			name:                   "don't break at space after footnote symbol at soft line start",
			content:                "abc def [1]: https://long/link\nghi",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc def \n[1]: https://long/link\nghi",
			expectedSoftLineBreaks: []int{8},
		},
		{
			name:                   "do break at subsequent space after footnote symbol",
			content:                "abc\n[1]: normal text follows\ndef",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc\n[1]: normal \ntext \nfollows\ndef",
			expectedSoftLineBreaks: []int{16, 21},
		},
		{
			name:                   "don't break at space after trailer",
			content:                "abc\nSigned-off-by: John Doe <john@doe.com>\nCo-authored-by: Jane Smith <jane@smith.com>\n",
			autoWrapWidth:          10,
			expectedWrappedContent: "abc\nSigned-off-by: John Doe <john@doe.com>\nCo-authored-by: Jane Smith <jane@smith.com>\n",
			expectedSoftLineBreaks: []int{},
		},
		{
			name:                   "do break at space after trailer if there is no space after the colon",
			content:                "abc\nSigned-off-by:John Doe <john@doe.com>\n",
			autoWrapWidth:          10,
			expectedWrappedContent: "abc\nSigned-off-by:John \nDoe \n<john@doe.com>\n",
			expectedSoftLineBreaks: []int{23, 27},
		},
		{
			name:                   "hard line breaks",
			content:                "abc\ndef\n",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc\ndef\n",
			expectedSoftLineBreaks: []int{},
		},
		{
			name:                   "mixture of hard and soft line breaks",
			content:                "abc def ghi jkl mno\npqr stu vwx yz\n",
			autoWrapWidth:          7,
			expectedWrappedContent: "abc def \nghi jkl \nmno\npqr stu \nvwx yz\n",
			expectedSoftLineBreaks: []int{8, 16, 28},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			textArea := &TextArea{content: tt.content, AutoWrapWidth: tt.autoWrapWidth, AutoWrap: true}
			cells, softLineBreakIndices := contentToCells(tt.content, tt.autoWrapWidth)
			textArea.cells = cells
			if !reflect.DeepEqual(textArea.GetContent(), tt.expectedWrappedContent) {
				t.Errorf("autoWrapContentImpl() wrappedContent = %v, expected %v", textArea.GetContent(), tt.expectedWrappedContent)
			}
			if !reflect.DeepEqual(softLineBreakIndices, tt.expectedSoftLineBreaks) {
				t.Errorf("autoWrapContentImpl() softLineBreakIndices = %v, expected %v", softLineBreakIndices, tt.expectedSoftLineBreaks)
			}

			// As a sanity check, run through all characters of the original content,
			// convert the cursor to the wrapped cursor, and check that the character
			// in the wrapped content at that position is the same:
			origCursor := 0
			for _, chr := range stringToGraphemes(tt.content) {
				wrappedIndex := textArea.contentCursorToCellCursor(origCursor)
				if chr != textArea.cells[wrappedIndex].char {
					t.Errorf("Runes in orig content and wrapped content don't match at %d: expected %v, got %v", origCursor, chr, textArea.cells[wrappedIndex].char)
				}

				// Also, check that converting the wrapped position back to the
				// orig position yields the original value again:
				origIndexAgain := textArea.cellCursorToContentCursor(wrappedIndex)
				if origCursor != origIndexAgain {
					t.Errorf("wrappedCursorToOrigCursor doesn't yield original position: expected %d, got %d", origCursor, origIndexAgain)
				}

				origCursor += len(chr)
			}
		})
	}
}

var testContent = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Quisque vehicula mi at elit pellentesque, eu pulvinar ligula molestie.
In vitae orci vitae elit fermentum lobortis sed in nisi.
Nam non odio nisi.
Donec vitae elit enim.
Pellentesque faucibus dolor at metus elementum sollicitudin.
Mauris eu orci vel odio ornare feugiat eget ac nisl.
Nam at dolor erat.
Integer sit amet rutrum lectus, mollis pretium sapien.
Maecenas ligula ipsum, congue vitae rhoncus eget, volutpat at quam.
Donec ac ultricies tortor, sit amet sollicitudin urna.
Integer porta ornare diam a imperdiet.
Praesent vulputate mi turpis, in porttitor diam commodo a.
Donec ut enim ligula.

[Thïs-is-not-à-fôôtnöte]: https://example.com/footnote

Sïgned-öff-by: This is not a trailer

[1]: This is a footnote

Signed-off-by: John Doe <john@doe.com>
`

func BenchmarkTypeCharacter(b *testing.B) {
	textArea := &TextArea{content: testContent, AutoWrapWidth: 72, AutoWrap: true}
	textArea.SetCursor2D(0, 0)

	b.ResetTimer()
	for b.Loop() {
		textArea.TypeCharacter("a")
	}
}
