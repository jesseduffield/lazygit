package utils

import (
	"sort"
	"strings"

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

func CaseAwareContains(haystack, needle string) bool {
	// if needle contains an uppercase letter, we'll do a case sensitive search
	if ContainsUppercase(needle) {
		return strings.Contains(haystack, needle)
	}

	return CaseInsensitiveContains(haystack, needle)
}

func ContainsUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}

	return false
}

func CaseInsensitiveContains(haystack, needle string) bool {
	return strings.Contains(
		strings.ToLower(haystack),
		strings.ToLower(needle),
	)
}
