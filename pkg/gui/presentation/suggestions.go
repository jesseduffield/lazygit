package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func GetSuggestionListDisplayStrings(suggestions []*types.Suggestion) [][]string {
	lines := make([][]string, len(suggestions))

	for i := range suggestions {
		lines[i] = getSuggestionDisplayStrings(suggestions[i])
	}

	return lines
}

func getSuggestionDisplayStrings(suggestion *types.Suggestion) []string {
	return []string{suggestion.Label}
}
