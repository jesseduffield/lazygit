package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterStrings(t *testing.T) {
	type scenario struct {
		needle         string
		haystack       []string
		useFuzzySearch bool
		expected       []string
	}

	scenarios := []scenario{
		{
			needle:         "",
			haystack:       []string{"test"},
			useFuzzySearch: true,
			expected:       []string{},
		},
		{
			needle:         "test",
			haystack:       []string{"test"},
			useFuzzySearch: true,
			expected:       []string{"test"},
		},
		{
			needle:         "o",
			haystack:       []string{"a", "o", "e"},
			useFuzzySearch: true,
			expected:       []string{"o"},
		},
		{
			needle:         "mybranch",
			haystack:       []string{"my_branch", "mybranch", "branch", "this is my branch"},
			useFuzzySearch: true,
			expected:       []string{"mybranch", "my_branch", "this is my branch"},
		},
		{
			needle:         "test",
			haystack:       []string{"not a good match", "this 'test' is a good match", "test"},
			useFuzzySearch: true,
			expected:       []string{"test", "this 'test' is a good match"},
		},
		{
			needle:         "test",
			haystack:       []string{"Test"},
			useFuzzySearch: true,
			expected:       []string{"Test"},
		},
		{
			needle:         "test",
			haystack:       []string{"integration-testing", "t_e_s_t"},
			useFuzzySearch: false,
			expected:       []string{"integration-testing"},
		},
		{
			needle:         "integr test",
			haystack:       []string{"integration-testing", "testing-integration"},
			useFuzzySearch: false,
			expected:       []string{"integration-testing", "testing-integration"},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, FilterStrings(s.needle, s.haystack, s.useFuzzySearch))
	}
}

func TestCaseInsensitiveContains(t *testing.T) {
	testCases := []struct {
		haystack string
		needle   string
		expected bool
	}{
		{"Hello, World!", "world", true},           // Case-insensitive match
		{"Hello, World!", "WORLD", true},           // Case-insensitive match
		{"Hello, World!", "orl", true},             // Case-insensitive match
		{"Hello, World!", "o, W", true},            // Case-insensitive match
		{"Hello, World!", "hello", true},           // Case-insensitive match
		{"Hello, World!", "Foo", false},            // No match
		{"Hello, World!", "Hello, World!!", false}, // No match
		{"Hello, World!", "", true},                // Empty needle matches
		{"", "Hello", false},                       // Empty haystack doesn't match
		{"", "", true},                             // Empty strings match
		{"", " ", false},                           // Empty haystack, non-empty needle
		{" ", "", true},                            // Non-empty haystack, empty needle
	}

	for i, testCase := range testCases {
		result := CaseInsensitiveContains(testCase.haystack, testCase.needle)
		assert.Equal(t, testCase.expected, result, fmt.Sprintf("Test case %d failed. Expected '%v', got '%v' for '%s' in '%s'", i, testCase.expected, result, testCase.needle, testCase.haystack))
	}
}
