package theme

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/stretchr/testify/assert"
)

func TestGetGocuiStyle(t *testing.T) {
	scenarios := []struct {
		name     string
		keys     []string
		expected gocui.Attribute
	}{
		{
			name:     "named color",
			keys:     []string{"red"},
			expected: gocui.ColorRed,
		},
		{
			name:     "named color with modifiers",
			keys:     []string{"red", "bold", "underline"},
			expected: gocui.ColorRed | gocui.AttrBold | gocui.AttrUnderline,
		},
		{
			name:     "strikethrough",
			keys:     []string{"red", "strikethrough"},
			expected: gocui.ColorRed | gocui.AttrStrikeThrough,
		},
		{
			name:     "dim",
			keys:     []string{"default", "dim"},
			expected: gocui.ColorDefault | gocui.AttrDim,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			assert.Equal(t, scenario.expected, GetGocuiStyle(scenario.keys))
		})
	}
}
