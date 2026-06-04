package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
)

func (gui *Gui) handleEditorKeypress(v *gocui.View, key gocui.Key, allowMultiline bool) bool {
	if key.Equals(gocui.NewKeyName(gocui.KeyEnter)) && allowMultiline {
		v.TextArea.TypeCharacter("\n")
		v.RenderTextArea()
		return true
	}

	return gocui.DefaultEditor.Edit(v, key)
}

// we've just copy+pasted the editor from gocui to here so that we can also re-
// render the commit message length on each keypress
func (gui *Gui) commitMessageEditor(v *gocui.View, key gocui.Key) bool {
	matched := gui.handleEditorKeypress(v, key, false)
	v.RenderTextArea()
	gui.c.Contexts().CommitMessage.RenderSubtitle()
	return matched
}

func (gui *Gui) commitDescriptionEditor(v *gocui.View, key gocui.Key) bool {
	matched := gui.handleEditorKeypress(v, key, true)
	v.RenderTextArea()
	return matched
}

func (gui *Gui) promptEditor(v *gocui.View, key gocui.Key) bool {
	matched := gui.handleEditorKeypress(v, key, false)

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

func (gui *Gui) searchEditor(v *gocui.View, key gocui.Key) bool {
	matched := gui.handleEditorKeypress(v, key, false)
	v.RenderTextArea()

	searchString := v.TextArea.GetContent()

	gui.helpers.Search.OnPromptContentChanged(searchString)

	return matched
}
