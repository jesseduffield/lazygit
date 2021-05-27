package utils

import (
	"sort"

	"github.com/sahilm/fuzzy"
)

func FuzzySearch(needle string, haystack []string) []string {
	if needle == "" {
		return []string{}
	}

	matches := fuzzy.Find(needle, haystack)
	sort.Sort(matches)

	result := make([]string, len(matches))
	for i, match := range matches {
		result[i] = match.Str
	}

	return result
}
