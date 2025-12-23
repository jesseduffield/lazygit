package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) handleEditorKeypress(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier, allowMultiline bool) bool {
	if key == gocui.KeyEnter && allowMultiline {
		v.TextArea.TypeCharacter("\n")
		v.RenderTextArea()
		return true
	}

	return gocui.DefaultEditor.Edit(v, key, ch, mod)
}

// we've just copy+pasted the editor from gocui to here so that we can also re-
// render the commit message length on each keypress
func (gui *Gui) commitMessageEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v, key, ch, mod, false)
	v.RenderTextArea()
	gui.c.Contexts().CommitMessage.RenderSubtitle()
	return matched
}

func (gui *Gui) commitDescriptionEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v, key, ch, mod, true)
	v.RenderTextArea()
	return matched
}

func (gui *Gui) promptEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	matched := gui.handleEditorKeypress(v, key, ch, mod, false)

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
	matched := gui.handleEditorKeypress(v, key, ch, mod, false)
	v.RenderTextArea()

	searchString := v.TextArea.GetContent()

	gui.helpers.Search.OnPromptContentChanged(searchString)

	return matched
}
