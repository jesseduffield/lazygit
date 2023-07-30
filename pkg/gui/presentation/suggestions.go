package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

func GetSuggestionListDisplayStrings(suggestions []*types.Suggestion) [][]string {
	return lo.Map(suggestions, func(suggestion *types.Suggestion, _ int) []string {
		return getSuggestionDisplayStrings(suggestion)
	})
}

func getSuggestionDisplayStrings(suggestion *types.Suggestion) []string {
	return []string{suggestion.Label}
}
