package gui

import (
	"unicode"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
)

func (gui *Gui) handleEditorKeypress(textArea *gocui.TextArea, key gocui.Key, ch rune, mod gocui.Modifier, allowMultiline bool) bool {
	newlineKey, ok := keybindings.GetKey(gui.c.UserConfig.Keybinding.Universal.AppendNewline).(gocui.Key)
	if !ok {
		newlineKey = gocui.KeyAltEnter
	}

	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		textArea.BackSpaceChar()
	case key == gocui.KeyCtrlD || key == gocui.KeyDelete:
		textArea.DeleteChar()
	case key == gocui.KeyArrowDown:
		textArea.MoveCursorDown()
	case key == gocui.KeyArrowUp:
		textArea.MoveCursorUp()
	case key == gocui.KeyArrowLeft && (mod&gocui.ModAlt) != 0:
		textArea.MoveLeftWord()
	case key == gocui.KeyArrowLeft || key == gocui.KeyCtrlB:
		textArea.MoveCursorLeft()
	case key == gocui.KeyArrowRight && (mod&gocui.ModAlt) != 0:
		textArea.MoveRightWord()
	case key == gocui.KeyArrowRight || key == gocui.KeyCtrlF:
		textArea.MoveCursorRight()
	case key == newlineKey:
		if allowMultiline {
			textArea.TypeRune('\n')
		} else {
			return false
		}
	case key == gocui.KeySpace:
		textArea.TypeRune(' ')
	case key == gocui.KeyInsert:
		textArea.ToggleOverwrite()
	case key == gocui.KeyCtrlU:
		textArea.DeleteToStartOfLine()
	case key == gocui.KeyCtrlK:
		textArea.DeleteToEndOfLine()
	case key == gocui.KeyCtrlA || key == gocui.KeyHome:
		textArea.GoToStartOfLine()
	case key == gocui.KeyCtrlE || key == gocui.KeyEnd:
		textArea.GoToEndOfLine()
	case key == gocui.KeyCtrlW:
		textArea.BackSpaceWord()
	case key == gocui.KeyCtrlY:
		textArea.Yank()

		// TODO: see if we need all three of these conditions: maybe the final one is sufficient
	case ch != 0 && mod == 0 && unicode.IsPrint(ch):
		textArea.TypeRune(ch)
	default:
		return false
	}

	return true
}

// we've just copy+pasted the editor from gocui to here so that we can also re-
// render the commit message length on each keypress
func (gui *Gui) commitMessageEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, true)

	// This function is called again on refresh as part of the general resize popup call,
	// but we need to call it here so that when we go to render the text area it's not
	// considered out of bounds to add a newline, meaning we can avoid unnecessary scrolling.
	err := gui.resizePopupPanel(v, v.TextArea.GetContent())
	if err != nil {
		gui.c.Log.Error(err)
	}
	v.RenderTextArea()
	gui.RenderCommitLength()

	return matched
}

func (gui *Gui) defaultEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, false)

	v.RenderTextArea()

	if gui.findSuggestions != nil {
		input := v.TextArea.GetContent()
		gui.suggestionsAsyncHandler.Do(func() func() {
			suggestions := gui.findSuggestions(input)
			return func() { gui.setSuggestions(suggestions) }
		})
	}

	return matched
}
