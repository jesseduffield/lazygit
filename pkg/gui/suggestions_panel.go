package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) getSelectedSuggestionValue() string {
	selectedSuggestion := gui.getSelectedSuggestion()

	if selectedSuggestion != nil {
		return selectedSuggestion.Value
	}

	return ""
}

func (gui *Gui) getSelectedSuggestion() *types.Suggestion {
	return gui.State.Contexts.Suggestions.GetSelected()
}

func (gui *Gui) setSuggestions(suggestions []*types.Suggestion) {
	gui.State.Suggestions = suggestions
	gui.State.Contexts.Suggestions.SetSelectedLineIdx(0)
	_ = gui.resetOrigin(gui.Views.Suggestions)
	_ = gui.State.Contexts.Suggestions.HandleRender()
}
