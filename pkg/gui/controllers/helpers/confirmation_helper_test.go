package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenuFilterPrompt(t *testing.T) {
	longPrompt := "Filter ('@' for keybindings): "
	shortPrompt := "Filter: "

	tests := []struct {
		name         string
		candidates   []string
		contentWidth int
		expected     string
	}{
		{name: "room for four characters", candidates: []string{shortPrompt}, contentWidth: 12, expected: shortPrompt},
		{name: "room for three characters", candidates: []string{shortPrompt}, contentWidth: 11, expected: ""},
		{name: "prefers the first candidate", candidates: []string{longPrompt, shortPrompt}, contentWidth: 34, expected: longPrompt},
		{name: "falls back to the next one", candidates: []string{longPrompt, shortPrompt}, contentWidth: 33, expected: shortPrompt},
		{name: "measures display width", candidates: []string{"篩選: "}, contentWidth: 9, expected: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, menuFilterPrompt(test.candidates, test.contentWidth))
		})
	}
}
