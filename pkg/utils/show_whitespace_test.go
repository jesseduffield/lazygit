package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowWhitespaceCharacters(t *testing.T) {
	dot := faintOn + "\u00B7" + faintOff
	arr := faintOn + "\u2192"
	fill := "-"

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no whitespace",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "spaces only",
			input:    "hello world",
			expected: "hello" + dot + "world",
		},
		{
			name:     "tab at column 0",
			input:    "\thello",
			expected: arr + fill + fill + fill + faintOff + "hello",
		},
		{
			name:     "tab at column 1 (after char)",
			input:    "a\tb",
			expected: "a" + arr + fill + fill + faintOff + "b",
		},
		{
			name:     "multiple spaces",
			input:    "a  b",
			expected: "a" + dot + dot + "b",
		},
		{
			name:     "newlines preserved",
			input:    "line1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "ANSI escape preserved",
			input:    "\x1b[31mred text\x1b[0m",
			expected: "\x1b[31mred" + dot + "text\x1b[0m",
		},
		{
			name:     "trailing spaces",
			input:    "trailing  \n",
			expected: "trailing" + dot + dot + "\n",
		},
		{
			name:     "leading spaces",
			input:    "  leading",
			expected: dot + dot + "leading",
		},
		{
			name:     "tab followed by spaces",
			input:    "\t  hello",
			expected: arr + fill + fill + fill + faintOff + dot + dot + "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShowWhitespaceCharacters(tt.input, 4)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShowWhitespaceCharacters_DifferentTabWidths(t *testing.T) {
	arr := faintOn + "\u2192"
	fill := "-"

	result := ShowWhitespaceCharacters("\t", 2)
	assert.Equal(t, arr+fill+faintOff, result)

	result = ShowWhitespaceCharacters("\t", 8)
	assert.Equal(t, arr+fill+fill+fill+fill+fill+fill+fill+faintOff, result)
}
