package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFuzzySearch is a function.
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
