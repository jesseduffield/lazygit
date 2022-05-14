package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuzzySearch(t *testing.T) {
	type scenario struct {
		needle   string
		haystack []string
		expected []string
	}

	scenarios := []scenario{
		{
			needle:   "",
			haystack: []string{"test"},
			expected: []string{},
		},
		{
			needle:   "test",
			haystack: []string{"test"},
			expected: []string{"test"},
		},
		{
			needle:   "o",
			haystack: []string{"a", "o", "e"},
			expected: []string{"o"},
		},
		{
			needle:   "mybranch",
			haystack: []string{"my_branch", "mybranch", "branch", "this is my branch"},
			expected: []string{"mybranch", "my_branch", "this is my branch"},
		},
		{
			needle:   "test",
			haystack: []string{"not a good match", "this 'test' is a good match", "test"},
			expected: []string{"test", "this 'test' is a good match"},
		},
		{
			needle:   "test",
			haystack: []string{"Test"},
			expected: []string{"Test"},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, FuzzySearch(s.needle, s.haystack))
	}
}

type node struct {
	value string
}

func TestFuzzySearchItems(t *testing.T) {
	type scenario struct {
		needle   string
		haystack []node
		expected []node
	}

	scenarios := []scenario{
		{
			needle:   "",
			haystack: []node{{"test"}},
			expected: []node{},
		},
		{
			needle:   "test",
			haystack: []node{{"test"}},
			expected: []node{{"test"}},
		},
		{
			needle:   "o",
			haystack: []node{{"a"}, {"o"}, {"e"}},
			expected: []node{{"o"}},
		},
		{
			needle:   "mybranch",
			haystack: []node{{"my_branch"}, {"mybranch"}, {"branch"}, {"this is my branch"}},
			expected: []node{{"mybranch"}, {"my_branch"}, {"this is my branch"}},
		},
		{
			needle:   "test",
			haystack: []node{{"not a good match"}, {"this 'test' is a good match"}, {"test"}},
			expected: []node{{"test"}, {"this 'test' is a good match"}},
		},
		{
			needle:   "test",
			haystack: []node{{"Test"}},
			expected: []node{{"Test"}},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, FuzzySearchItems(s.needle, s.haystack, func(n node) string { return n.value }))
	}
}
