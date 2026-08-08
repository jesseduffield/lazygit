package gocui

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const esc = "\x1b"

func runVimKeys(editor *VimEditor, ta *TextArea, keys string) {
	for _, r := range keys {
		ch := string(r)
		if ch == esc {
			editor.escape(ta)
			continue
		}
		if editor.mode == VimModeInsert {
			ta.TypeCharacter(ch)
			continue
		}
		editor.normalChar(ta, ch)
	}
}

func newVimFixture(content string, cursor int, multiline bool) (*VimEditor, *TextArea) {
	ta := makeTextArea(content, cursor)
	editor := NewVimEditor(nil, multiline)
	editor.mode = VimModeNormal
	return editor, ta
}

func TestVimEditorCommands(t *testing.T) {
	tests := []struct {
		content         string
		cursor          int
		keys            string
		expectedContent string
		expectedCursor  int
	}{
		{"foo bar baz", 0, "dw", "bar baz", 0},
		{"foo bar baz", 0, "d2w", "baz", 0},
		{"foo bar baz", 0, "2dw", "baz", 0},
		{"foo bar baz", 4, "diw", "foo  baz", 4},
		{"foo bar baz", 4, "daw", "foo baz", 4},
		{"foo bar baz", 4, "de", "foo  baz", 4},
		{"foo bar baz", 4, "d$", "foo ", 3},
		{"foo bar baz", 4, "D", "foo ", 3},
		{"foo bar baz", 4, "d0", "bar baz", 0},
		{"foo bar", 0, "dfa", "r", 0},
		{"foo bar", 0, "dta", "ar", 0},
		{"foo bar baz", 8, "db", "foo baz", 4},
		{"foo\nbar", 0, "dw", "\nbar", 0},
		{"foo\nbar", 0, "dd", "bar", 0},
		{"foo\nbar", 4, "dd", "foo", 0},
		{"a\nb\nc", 0, "2dd", "c", 0},
		{"a\nb\nc", 0, "dG", "", 0},
		{"a\nb\nc", 4, "dgg", "", 0},
		{"a\nb\nc", 0, "dj", "c", 0},
		{"a\nb\nc", 4, "dk", "a", 0},
		{"foo", 0, "x", "oo", 0},
		{"foo", 0, "3x", "", 0},
		{"foo", 0, "5x", "", 0},
		{"foo", 2, "X", "fo", 1},
		{"foo", 2, "2X", "o", 0},
		{"foo", 0, "rz", "zoo", 0},
		{"foo", 0, "3rz", "zzz", 2},
		{"foo", 0, "5rz", "zzz", 2},
		{"foo bar", 0, "cwme" + esc, "me bar", 1},
		{"foo bar", 0, "ciwme" + esc, "me bar", 1},
		{"foo bar", 3, "cwx" + esc, "fooxbar", 3},
		{"foo", 0, "cc" + "bar" + esc, "bar", 2},
		{"a\nb\nc", 0, "2ccx" + esc, "x\nc", 0},
		{"foo", 1, "C" + "x" + esc, "fx", 1},
		{"foo", 1, "S" + "x" + esc, "x", 0},
		{"foo", 1, "s" + "x" + esc, "fxo", 1},
		{"foo", 0, "ibar" + esc, "barfoo", 2},
		{"foo", 0, "abar" + esc, "fbaroo", 3},
		{"foo bar", 4, "Ix" + esc, "xfoo bar", 0},
		{"foo", 0, "Ax" + esc, "foox", 3},
		{"foo", 0, "ox" + esc, "foo\nx", 4},
		{"foo", 0, "Ox" + esc, "x\nfoo", 0},
		{"foo bar", 0, "w", "foo bar", 4},
		{"foo bar", 0, "ww", "foo bar", 6},
		{"foo bar baz", 0, "2w", "foo bar baz", 8},
		{"foo bar", 6, "bb", "foo bar", 0},
		{"foo bar", 0, "$", "foo bar", 6},
		{"foo bar", 6, "0", "foo bar", 0},
		{"  foo", 4, "^", "  foo", 2},
		{"foo bar", 0, "e", "foo bar", 2},
		{"abcabc", 0, "2fc", "abcabc", 5},
		{"abcabc", 0, "fc;", "abcabc", 5},
		{"abcabc", 0, "2fc,", "abcabc", 2},
		{"a\nb\nc", 0, "G", "a\nb\nc", 4},
		{"a\nb\nc", 4, "gg", "a\nb\nc", 0},
		{"a\nb\nc", 0, "2G", "a\nb\nc", 2},
		{"foo", 2, "hh", "foo", 0},
		{"foo", 0, "ll", "foo", 2},
		{"foo", 0, "3l", "foo", 2},
		{"foo bar", 0, "ywP", "foo foo bar", 3},
		{"foo", 0, "xp", "ofo", 1},
		{"foo", 0, "yiwp", "ffoooo", 3},
		{"a\nb", 0, "yyp", "a\na\nb", 2},
		{"a\nb", 2, "yyp", "a\nb\nb", 4},
		{"a\nb", 2, "yyP", "a\nb\nb", 2},
		{"foo bar", 0, "dwu", "foo bar", 0},
		{"foo bar", 0, "dwdwuu", "foo bar", 0},
		{"foo", 0, "ibar" + esc + "u", "foo", 0},
		{"foo", 0, "xxu", "oo", 0},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			editor, ta := newVimFixture(test.content, test.cursor, true)
			runVimKeys(editor, ta, test.keys)
			assert.Equal(t, test.expectedContent, ta.GetUnwrappedContent())
			assert.Equal(t, test.expectedCursor, ta.cursor)
		})
	}
}

func TestVimEditorRedo(t *testing.T) {
	editor, ta := newVimFixture("foo bar", 0, true)
	runVimKeys(editor, ta, "dwu")
	assert.Equal(t, "foo bar", ta.GetUnwrappedContent())
	editor.normalKey(ta, NewKeyStrMod("r", ModCtrl))
	assert.Equal(t, "bar", ta.GetUnwrappedContent())
	runVimKeys(editor, ta, "u")
	assert.Equal(t, "foo bar", ta.GetUnwrappedContent())
}

func TestVimEditorRegisters(t *testing.T) {
	editor, ta := newVimFixture("foo bar", 0, true)
	runVimKeys(editor, ta, "yw")
	assert.Equal(t, "foo ", editor.register.text)
	assert.False(t, editor.register.linewise)

	runVimKeys(editor, ta, "yy")
	assert.Equal(t, "foo bar", editor.register.text)
	assert.True(t, editor.register.linewise)

	shared := &VimRegister{}
	editorA := NewVimEditor(shared, true)
	editorA.mode = VimModeNormal
	taA := makeTextArea("hello", 0)
	editorA.normalChar(taA, "y")
	editorA.normalChar(taA, "w")

	editorB := NewVimEditor(shared, true)
	editorB.mode = VimModeNormal
	taB := makeTextArea("x", 0)
	editorB.normalChar(taB, "P")
	assert.Equal(t, "hellox", taB.GetUnwrappedContent())
}

func TestVimEditorSingleLine(t *testing.T) {
	editor, ta := newVimFixture("foo", 0, false)
	runVimKeys(editor, ta, "o")
	assert.Equal(t, "foo", ta.GetUnwrappedContent())
	runVimKeys(editor, ta, "O")
	assert.Equal(t, "foo", ta.GetUnwrappedContent())

	editor.register.set("a\nb", true)
	runVimKeys(editor, ta, "P")
	assert.Equal(t, "a bfoo", ta.GetUnwrappedContent())
}

func TestVimEditorEscape(t *testing.T) {
	editor, ta := newVimFixture("foo", 0, true)

	editor.mode = VimModeInsert
	assert.True(t, editor.escape(ta))
	assert.Equal(t, VimModeNormal, editor.mode)

	editor.normalChar(ta, "d")
	assert.True(t, editor.escape(ta))
	assert.Equal(t, "", editor.op)

	editor.normalChar(ta, "3")
	assert.True(t, editor.escape(ta))
	assert.Equal(t, 0, editor.count)

	assert.False(t, editor.escape(ta))
}

func TestVimEditorEscapeCursor(t *testing.T) {
	editor, ta := newVimFixture("foo", 0, true)
	runVimKeys(editor, ta, "A")
	assert.Equal(t, 3, ta.cursor)
	runVimKeys(editor, ta, esc)
	assert.Equal(t, 2, ta.cursor)
}

func TestVimEditorReset(t *testing.T) {
	editor, ta := newVimFixture("foo", 0, true)
	runVimKeys(editor, ta, "x")
	assert.NotEmpty(t, editor.undo)

	editor.Reset(VimModeInsert)
	assert.Equal(t, VimModeInsert, editor.Mode())
	assert.Empty(t, editor.undo)
	_ = ta
}
