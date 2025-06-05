package gui

import (
	"unicode"

	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleEditorKeypress(textArea *gocui.TextArea, key gocui.Key, ch rune, mod gocui.Modifier, allowMultiline bool) bool {
	switch {
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		textArea.BackSpaceChar()
	case key == gocui.KeyCtrlD || key == gocui.KeyDelete:
		textArea.DeleteChar()
	case key == gocui.KeyArrowDown:
		textArea.MoveCursorDown()
	case key == gocui.KeyArrowUp:
		textArea.MoveCursorUp()
	case (key == gocui.KeyArrowLeft || ch == 'b') && (mod&gocui.ModAlt) != 0:
		textArea.MoveLeftWord()
	case key == gocui.KeyArrowLeft || key == gocui.KeyCtrlB:
		textArea.MoveCursorLeft()
	case (key == gocui.KeyArrowRight || ch == 'f') && (mod&gocui.ModAlt) != 0:
		textArea.MoveRightWord()
	case key == gocui.KeyArrowRight || key == gocui.KeyCtrlF:
		textArea.MoveCursorRight()
	case key == gocui.KeyEnter:
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

	case unicode.IsPrint(ch):
		textArea.TypeRune(ch)
	default:
		return false
	}

	return true
}

// we've just copy+pasted the editor from gocui to here so that we can also re-
// render the commit message length on each keypress
func (gui *Gui) commitMessageEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, false)
	v.RenderTextArea()
	gui.c.Contexts().CommitMessage.RenderSubtitle()
	return matched
}

func (gui *Gui) commitDescriptionEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, true)
	v.RenderTextArea()
	return matched
}

func (gui *Gui) promptEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, false)

	v.RenderTextArea()

	suggestionsContext := gui.State.Contexts.Suggestions
	if suggestionsContext.State.FindSuggestions != nil {
		input := v.TextArea.GetContent()
		suggestionsContext.State.AsyncHandler.Do(func() func() {
			suggestions := suggestionsContext.State.FindSuggestions(input)
			return func() { suggestionsContext.SetSuggestions(suggestions) }
		})
	}

	return matched
}

func (gui *Gui) searchEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v.TextArea, key, ch, mod, false)
	v.RenderTextArea()

	searchString := v.TextArea.GetContent()

	gui.helpers.Search.OnPromptContentChanged(searchString)

	return matched
}
