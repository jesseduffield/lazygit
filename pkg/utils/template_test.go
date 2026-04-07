package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResolvePlaceholderString is a function.
func TestResolvePlaceholderString(t *testing.T) {
	type scenario struct {
		templateString string
		arguments      map[string]string
		expected       string
	}

	scenarios := []scenario{
		{
			"",
			map[string]string{},
			"",
		},
		{
			"hello",
			map[string]string{},
			"hello",
		},
		{
			"hello {{arg}}",
			map[string]string{},
			"hello {{arg}}",
		},
		{
			"hello {{arg}}",
			map[string]string{"arg": "there"},
			"hello there",
		},
		{
			"hello",
			map[string]string{"arg": "there"},
			"hello",
		},
		{
			"{{nothing}}",
			map[string]string{"nothing": ""},
			"",
		},
		{
			"{{}} {{ this }} { should not throw}} an {{{{}}}} error",
			map[string]string{
				"blah": "blah",
				"this": "won't match",
			},
			"{{}} {{ this }} { should not throw}} an {{{{}}}} error",
		},
		{
			"{{a}}",
			map[string]string{
				"a": "X{{.a}}X",
			},
			"X{{.a}}X",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, ResolvePlaceholderString(s.templateString, s.arguments))
	}
}

func TestSanitizeTerminalTitle(t *testing.T) {
	scenarios := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "lazygit::myproject",
			expected: "lazygit::myproject",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string with spaces",
			input:    "lazygit :: my project",
			expected: "lazygit :: my project",
		},
		{
			name:     "string with unicode",
			input:    "lazygit::项目",
			expected: "lazygit::项目",
		},
		{
			name:     "removes null byte",
			input:    "lazy\x00git",
			expected: "lazygit",
		},
		{
			name:     "removes escape sequence",
			input:    "lazy\x1b[31mgit",
			expected: "lazy[31mgit",
		},
		{
			name:     "removes newline and tab",
			input:    "lazy\ngit\ttitle",
			expected: "lazygittitle",
		},
		{
			name:     "removes carriage return",
			input:    "lazy\rgit",
			expected: "lazygit",
		},
		{
			name:     "removes DEL character",
			input:    "lazy\x7fgit",
			expected: "lazygit",
		},
		{
			name:     "removes bell character",
			input:    "lazy\x07git",
			expected: "lazygit",
		},
		{
			name:     "preserves printable special chars",
			input:    "lazy&git|test<>",
			expected: "lazy&git|test<>",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.EqualValues(t, s.expected, SanitizeTerminalTitle(s.input))
		})
	}
}
