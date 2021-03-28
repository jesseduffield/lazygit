package gui

import (
	"strings"
	"unicode"

	"github.com/jesseduffield/gocui"
)

// we've just copy+pasted the editor from gocui to here so that we can also re-
// render the commit message length on each keypress
func (gui *Gui) commitMessageEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	newlineKey, ok := gui.getKey(gui.Config.GetUserConfig().Keybinding.Universal.AppendNewline).(gocui.Key)
	if !ok {
		newlineKey = gocui.KeyTab
	}

	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyArrowDown:
		_, cursorY := v.Cursor()
		_, viewY := v.Size()
		if cursorY+1 == viewY {
			increaseViewHeight(v, cursorY)
		} else {
			v.MoveCursor(0, 1, false)
		}
	case key == gocui.KeyArrowUp:
		cursorX, cursorY := v.Cursor()
		_, viewY := v.Size()
		if viewY > 5 && cursorY+1 == viewY {
			decreaseViewHeightWhenLastLineIsBlank(v, cursorY)
			v.SetCursor(cursorX, cursorY-1)
		} else {
			v.MoveCursor(0, -1, false)
		}
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == newlineKey:
		v.EditNewLine()
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyCtrlU:
		v.EditDeleteToStartOfLine()
	case key == gocui.KeyCtrlA:
		v.EditGotoToStartOfLine()
	case key == gocui.KeyCtrlE:
		v.EditGotoToEndOfLine()
	case unicode.IsPrint(ch):
		v.EditWrite(ch)
	}

	gui.RenderCommitLength()
}

func (gui *Gui) defaultEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyCtrlU:
		v.EditDeleteToStartOfLine()
	case key == gocui.KeyCtrlA:
		v.EditGotoToStartOfLine()
	case key == gocui.KeyCtrlE:
		v.EditGotoToEndOfLine()
	case unicode.IsPrint(ch):
		v.EditWrite(ch)
	}

	if gui.findSuggestions != nil {
		input := v.Buffer()
		suggestions := gui.findSuggestions(input)
		gui.setSuggestions(suggestions)
	}
}

func decreaseViewHeightWhenLastLineIsBlank(v *gocui.View, currentLineNumber int) {
	line, err := v.Line(currentLineNumber)
	if err != nil {
		return
	}

	if len(strings.TrimSpace(line)) > 0 {
		return
	}

	lineLength := len(line)
	v.SetCursor(0, currentLineNumber)
	for i := 0; i < lineLength; i++ {
		v.EditDelete(false)
	}
	v.EditDelete(true)
}

func increaseViewHeight(v *gocui.View, currentLineNumber int) {
	line, err := v.Line(currentLineNumber)
	if err != nil {
		return
	}

	lineLength := len(line)
	v.SetCursor(lineLength, currentLineNumber)
	v.EditNewLine()
}
