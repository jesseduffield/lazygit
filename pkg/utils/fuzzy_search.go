package utils

import (
	"sort"

	"github.com/jesseduffield/generics/slices"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
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

type fuzzySource[T any] struct {
	items    []T
	toString func(T) string
}

func (self *fuzzySource[T]) String(i int) string {
	return self.toString(self.items[i])
}

func (self *fuzzySource[T]) Len() int {
	return len(self.items)
}

var _ fuzzy.Source = &fuzzySource[any]{}

func FuzzySearchItems[T any](needle string, items []T, toString func(T) string) []T {
	source := &fuzzySource[T]{
		items:    items,
		toString: toString,
	}

	matches := fuzzy.FindFrom(needle, source)
	return lo.Map(matches, func(match fuzzy.Match, _ int) T {
		return items[match.Index]
	})
}
