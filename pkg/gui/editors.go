package gui

import (
	"strings"

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
	gui.onPromptContentChanged(v)

	return matched
}

func (gui *Gui) onPromptContentChanged(v *gocui.View) {
	suggestionsContext := gui.State.Contexts.Suggestions
	// Capture the suggestions function and the input here, on the UI thread; the
	// main thread rewrites State.FindSuggestions when it (re)creates a prompt
	// panel, so reading it from the worker below would race that write.
	if findSuggestions := suggestionsContext.State.FindSuggestions; findSuggestions != nil {
		input := v.TextArea.GetContent()
		suggestionsContext.State.AsyncHandler.Do(func() func() {
			suggestions := findSuggestions(input)
			return func() { suggestionsContext.SetSuggestions(suggestions) }
		})
	}
}

func (gui *Gui) menuFilterEditor(v *gocui.View, key gocui.Key) bool {
	contentBefore := v.TextArea.GetContent()

	matched := gui.handleEditorKeypress(v, key, false)
	if !matched {
		// Give the global keybindings a chance at the key, e.g. so that ctrl-c
		// still quits while a menu is open.
		return false
	}

	v.RenderTextArea()

	content := v.TextArea.GetContent()
	if content == contentBefore {
		// The key just moved the cursor around within the filter; refiltering would
		// throw away the menu's selection for nothing.
		return true
	}

	gui.onMenuFilterContentChanged(content)

	return true
}

func (gui *Gui) pasteIntoMenuFilter(text string) {
	v := gui.Views.MenuFilter
	if typePastedTextIntoSingleLineView(v, text) {
		gui.onMenuFilterContentChanged(v.TextArea.GetContent())
	}
}

// typePastedTextIntoSingleLineView returns false if nothing is left of the text
// to type.
func typePastedTextIntoSingleLineView(v *gocui.View, text string) bool {
	text = pastedTextAsSingleLine(text)
	if text == "" {
		return false
	}

	v.TextArea.TypeString(text)
	v.RenderTextArea()
	return true
}

// Views that take a single line of text get pasted text without the newlines at
// its start and end, which are often part of a copied line, and with spaces in
// place of the newlines and tabs inside it.
func pastedTextAsSingleLine(text string) string {
	text = strings.Trim(text, "\n")
	return strings.NewReplacer("\n", " ", "\t", " ").Replace(text)
}

func (gui *Gui) onMenuFilterContentChanged(content string) {
	menuContext := gui.State.Contexts.Menu
	menuContext.SetFilterStarted(true)
	gui.helpers.Search.ApplyFilter(menuContext, content)
}

func (gui *Gui) searchEditor(v *gocui.View, key gocui.Key) bool {
	matched := gui.handleEditorKeypress(v, key, false)
	v.RenderTextArea()

	searchString := v.TextArea.GetContent()

	gui.helpers.Search.OnPromptContentChanged(searchString)

	return matched
}
