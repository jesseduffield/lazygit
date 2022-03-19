package utils

import (
	"sort"

	"github.com/jesseduffield/generics/slices"
	"github.com/sahilm/fuzzy"
)

func FuzzySearch(needle string, haystack []string) []string {
	if needle == "" {
		return []string{}
	}

	matches := fuzzy.Find(needle, haystack)
	sort.Sort(matches)

	return slices.Map(matches, func(match fuzzy.Match) string {
		return match.Str
	})
}
