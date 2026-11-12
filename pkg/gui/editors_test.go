package gui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPastedTextAsSingleLine(t *testing.T) {
	scenarios := []struct {
		name     string
		text     string
		expected string
	}{
		{"single line", "foo bar", "foo bar"},
		{"newlines at start and end", "\n\nfoo\n", "foo"},
		{"newlines inside", "foo\nbar\n\nbaz", "foo bar  baz"},
		{"tabs", "\tfoo\tbar", " foo bar"},
		{"only newlines", "\n\n", ""},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, pastedTextAsSingleLine(s.text))
		})
	}
}
