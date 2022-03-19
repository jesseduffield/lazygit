package presentation

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func GetSuggestionListDisplayStrings(suggestions []*types.Suggestion) [][]string {
	return slices.Map(suggestions, func(suggestion *types.Suggestion) []string {
		return getSuggestionDisplayStrings(suggestion)
	})
}

func getSuggestionDisplayStrings(suggestion *types.Suggestion) []string {
	return []string{suggestion.Label}
}
