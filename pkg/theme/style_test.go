package theme

import (
	"reflect"
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

func TestGetTextStyle(t *testing.T) {
	scenarios := []struct {
		name       string
		keys       []string
		background bool
		expected   style.TextStyle
	}{
		{
			name:       "empty",
			keys:       []string{""},
			background: true,
			expected:   style.New(),
		},
		{
			name:       "named color, fg",
			keys:       []string{"blue"},
			background: false,
			expected:   style.New().SetFg(style.NewBasicColor(color.FgBlue)),
		},
		{
			name:       "named color, bg",
			keys:       []string{"blue"},
			background: true,
			expected:   style.New().SetBg(style.NewBasicColor(color.BgBlue)),
		},
		{
			name:       "hex color, fg",
			keys:       []string{"#123456"},
			background: false,
			expected:   style.New().SetFg(style.NewRGBColor(color.RGBColor{0x12, 0x34, 0x56, 0})),
		},
		{
			name:       "hex color, bg",
			keys:       []string{"#abcdef"},
			background: true,
			expected:   style.New().SetBg(style.NewRGBColor(color.RGBColor{0xab, 0xcd, 0xef, 1})),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			if actual := GetTextStyle(scenario.keys, scenario.background); !reflect.DeepEqual(actual, scenario.expected) {
				t.Errorf("GetTextStyle() = %v, expected %v", actual, scenario.expected)
			}
		})
	}
}
