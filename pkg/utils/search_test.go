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

func TestViewFilterPattern(t *testing.T) {
	const def = "re:"
	pat, re := ViewFilterPattern("substring", "re:^main", def)
	assert.True(t, re)
	assert.Equal(t, "^main", pat)

	pat, re = ViewFilterPattern("substring", "plain", def)
	assert.False(t, re)
	assert.Equal(t, "plain", pat)

	pat, re = ViewFilterPattern("regexp", "^main", def)
	assert.True(t, re)
	assert.Equal(t, "^main", pat)

	pat, re = ViewFilterPattern("fuzzy", "re:a.c", def)
	assert.True(t, re)
	assert.Equal(t, "a.c", pat)

	pat, re = ViewFilterPattern("substring", "rx:^x", "rx:")
	assert.True(t, re)
	assert.Equal(t, "^x", pat)

	pat, re = ViewFilterPattern("substring", "re:^main", "rx:")
	assert.False(t, re)
	assert.Equal(t, "re:^main", pat)
}

func TestFindFromRegexp(t *testing.T) {
	haystack := []string{"main", "amain", "xmain"}
	src := stringSource(haystack)

	got := FindFrom("^main", src, false, true)
	assert.Len(t, got, 1)
	assert.Equal(t, 0, got[0].Index)
	assert.Equal(t, "main", got[0].Str)

	// '.' is regexp metacharacter: 'foo.go' matches 'fooXgo'; literal dot needs '\.'
	dotHay := stringSource([]string{"fooXgo", "foo.go"})
	got = FindFrom("foo.go", dotHay, false, true)
	assert.Len(t, got, 2)
	got = FindFrom(`foo\.go`, dotHay, false, true)
	assert.Len(t, got, 1)
	assert.Equal(t, "foo.go", got[0].Str)

	// invalid pattern => no matches
	got = FindFrom("(", src, false, true)
	assert.Empty(t, got)

	// case: lowercase pattern matches uppercase (implicit (?i))
	got = FindFrom("main", stringSource([]string{"Main"}), false, true)
	assert.Len(t, got, 1)

	// uppercase in pattern => case-sensitive
	got = FindFrom("Main", stringSource([]string{"main", "Main"}), false, true)
	assert.Len(t, got, 1)
	assert.Equal(t, "Main", got[0].Str)
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
