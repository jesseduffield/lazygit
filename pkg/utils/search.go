package utils

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

func FilterStrings(needle string, haystack []string, useFuzzySearch bool) []string {
	if needle == "" {
		return []string{}
	}

	matches := Find(needle, haystack, useFuzzySearch)

	return lo.Map(matches, func(match fuzzy.Match, _ int) string {
		return match.Str
	})
}

// Duplicated from the fuzzy package because it's private there
type stringSource []string

func (ss stringSource) String(i int) string {
	return ss[i]
}

func (ss stringSource) Len() int { return len(ss) }

// Drop-in replacement for fuzzy.Find (except that it doesn't fill out
// MatchedIndexes or Score, but we are not using these)
func FindSubstrings(pattern string, data []string) fuzzy.Matches {
	return FindSubstringsFrom(pattern, stringSource(data))
}

// Drop-in replacement for fuzzy.FindFrom (except that it doesn't fill out
// MatchedIndexes or Score, but we are not using these)
func FindSubstringsFrom(pattern string, data fuzzy.Source) fuzzy.Matches {
	substrings := strings.Fields(pattern)
	result := fuzzy.Matches{}

outer:
	for i := 0; i < data.Len(); i++ {
		s := data.String(i)
		for _, sub := range substrings {
			if !CaseAwareContains(s, sub) {
				continue outer
			}
		}
		result = append(result, fuzzy.Match{Str: s, Index: i})
	}

	return result
}

func Find(pattern string, data []string, useFuzzySearch bool) fuzzy.Matches {
	if useFuzzySearch {
		return fuzzy.Find(pattern, data)
	}
	return FindSubstrings(pattern, data)
}

func FindFrom(pattern string, data fuzzy.Source, useFuzzySearch bool) fuzzy.Matches {
	if useFuzzySearch {
		return fuzzy.FindFrom(pattern, data)
	}
	return FindSubstringsFrom(pattern, data)
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
