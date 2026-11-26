package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminalDrawsBranchDrawingSymbols(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{name: "kitty", version: "0.36.1", expected: false},
		{name: "kitty", version: "0.36.2", expected: true},
		{name: "kitty", version: "0.44.0", expected: true},
		{name: "ghostty", version: "0.9.0", expected: false},
		{name: "ghostty", version: "1.0.0", expected: true},
		{name: "ghostty", version: "1.3.0-main+0123abcd", expected: true},
		{name: "WezTerm", version: "20250601-102030-89abcdef", expected: false},
		{name: "tmux", version: "3.5a", expected: false},
		{name: "iTerm2", version: "3.6.4", expected: false},
		{name: "kitty", version: "", expected: false},
		{name: "", version: "", expected: false},
	}

	for _, test := range tests {
		t.Run(test.name+" "+test.version, func(t *testing.T) {
			assert.Equal(t, test.expected, TerminalDrawsBranchDrawingSymbols(test.name, test.version))
		})
	}
}
