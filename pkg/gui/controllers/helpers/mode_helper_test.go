package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestScreenModeAfterClearingFiltering(t *testing.T) {
	tests := []struct {
		name             string
		current          types.ScreenMode
		forcedHalfScreen bool
		expected         types.ScreenMode
	}{
		{
			name:             "restore normal after filtering promoted screen mode",
			current:          types.SCREEN_HALF,
			forcedHalfScreen: true,
			expected:         types.SCREEN_NORMAL,
		},
		{
			name:             "preserve configured half screen mode",
			current:          types.SCREEN_HALF,
			forcedHalfScreen: false,
			expected:         types.SCREEN_HALF,
		},
		{
			name:             "preserve full screen mode chosen while filtering",
			current:          types.SCREEN_FULL,
			forcedHalfScreen: true,
			expected:         types.SCREEN_FULL,
		},
		{
			name:             "preserve normal when user already switched back",
			current:          types.SCREEN_NORMAL,
			forcedHalfScreen: true,
			expected:         types.SCREEN_NORMAL,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, screenModeAfterClearingFiltering(test.current, test.forcedHalfScreen))
		})
	}
}
