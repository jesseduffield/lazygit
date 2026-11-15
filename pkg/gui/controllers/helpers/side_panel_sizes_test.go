package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestSidePanelSizeForScreen(t *testing.T) {
	laptopAndMonitor := []config.SidePanelSize{
		{ScreenSize: 160, Ratio: 0.3125}, // 50 columns
		{ScreenSize: 380, Ratio: 0.2},    // 76 columns
	}

	scenarios := []struct {
		name         string
		sizes        []config.SidePanelSize
		screenSize   int
		expectedSize int
		expectedOk   bool
	}{
		{
			name:       "no sizes",
			sizes:      nil,
			screenSize: 100,
			expectedOk: false,
		},
		{
			name:         "a single size keeps its ratio at any screen size",
			sizes:        []config.SidePanelSize{{ScreenSize: 160, Ratio: 0.25}},
			screenSize:   80,
			expectedSize: 20,
			expectedOk:   true,
		},
		{
			name:         "at a screen size that the user set it for",
			sizes:        laptopAndMonitor,
			screenSize:   380,
			expectedSize: 76,
			expectedOk:   true,
		},
		{
			name:         "below the smallest screen size, keeps its ratio",
			sizes:        laptopAndMonitor,
			screenSize:   80,
			expectedSize: 25,
			expectedOk:   true,
		},
		{
			name:         "above the largest screen size, keeps its ratio",
			sizes:        laptopAndMonitor,
			screenSize:   500,
			expectedSize: 100,
			expectedOk:   true,
		},
		{
			name:         "in between, interpolates the size",
			sizes:        laptopAndMonitor,
			screenSize:   270,
			expectedSize: 63,
			expectedOk:   true,
		},
		{
			name: "in between, interpolates the size rather than the ratio",
			sizes: []config.SidePanelSize{
				{ScreenSize: 100, Ratio: 0.5}, // 50 columns
				{ScreenSize: 400, Ratio: 0.1}, // 40 columns
			},
			screenSize:   250,
			expectedSize: 45,
			expectedOk:   true,
		},
		{
			name: "picks the two neighboring sizes",
			sizes: []config.SidePanelSize{
				{ScreenSize: 100, Ratio: 0.3},
				{ScreenSize: 200, Ratio: 0.3},
				{ScreenSize: 300, Ratio: 0.2},
			},
			screenSize:   250,
			expectedSize: 60,
			expectedOk:   true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			size, ok := sidePanelSizeForScreen(s.sizes, s.screenSize)
			assert.Equal(t, s.expectedOk, ok)
			assert.Equal(t, s.expectedSize, size)
		})
	}
}

func TestRecordSidePanelSize(t *testing.T) {
	scenarios := []struct {
		name       string
		sizes      []config.SidePanelSize
		screenSize int
		size       int
		expected   []config.SidePanelSize
	}{
		{
			name:       "first size",
			sizes:      nil,
			screenSize: 160,
			size:       50,
			expected:   []config.SidePanelSize{{ScreenSize: 160, Ratio: 0.3125}},
		},
		{
			name:       "a clearly different screen size adds a size",
			sizes:      []config.SidePanelSize{{ScreenSize: 160, Ratio: 0.3125}},
			screenSize: 380,
			size:       80,
			expected: []config.SidePanelSize{
				{ScreenSize: 160, Ratio: 0.3125},
				{ScreenSize: 380, Ratio: 0.2105},
			},
		},
		{
			name:       "a smaller screen size is sorted first",
			sizes:      []config.SidePanelSize{{ScreenSize: 160, Ratio: 0.3125}},
			screenSize: 100,
			size:       40,
			expected: []config.SidePanelSize{
				{ScreenSize: 100, Ratio: 0.4},
				{ScreenSize: 160, Ratio: 0.3125},
			},
		},
		{
			name: "a similar screen size replaces the size",
			sizes: []config.SidePanelSize{
				{ScreenSize: 160, Ratio: 0.3125},
				{ScreenSize: 380, Ratio: 0.2105},
			},
			screenSize: 400,
			size:       100,
			expected: []config.SidePanelSize{
				{ScreenSize: 160, Ratio: 0.3125},
				{ScreenSize: 400, Ratio: 0.25},
			},
		},
		{
			name: "replaces all similar screen sizes",
			sizes: []config.SidePanelSize{
				{ScreenSize: 100, Ratio: 0.4},
				{ScreenSize: 110, Ratio: 0.35},
				{ScreenSize: 120, Ratio: 0.3},
				{ScreenSize: 200, Ratio: 0.2},
			},
			screenSize: 110,
			size:       33,
			expected: []config.SidePanelSize{
				{ScreenSize: 110, Ratio: 0.3},
				{ScreenSize: 200, Ratio: 0.2},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, recordSidePanelSize(s.sizes, s.screenSize, s.size))
		})
	}
}

func TestRecordedSidePanelSizeComesBackUnchanged(t *testing.T) {
	for screenSize := 20; screenSize <= 2000; screenSize++ {
		for _, size := range []int{3, screenSize / 3, screenSize - 3} {
			sizes := recordSidePanelSize(nil, screenSize, size)
			result, ok := sidePanelSizeForScreen(sizes, screenSize)
			assert.True(t, ok)
			if result != size {
				t.Fatalf("size %d at screen size %d came back as %d", size, screenSize, result)
			}
		}
	}
}

func TestValidSidePanelSizes(t *testing.T) {
	input := map[config.SidePanelLayout][]config.SidePanelSize{
		config.SidePanelLayoutNormal: {
			{ScreenSize: 300, Ratio: 0.2},
			{ScreenSize: 0, Ratio: 0.3},
			{ScreenSize: 100, Ratio: 0.4},
			{ScreenSize: 200, Ratio: 0},
			{ScreenSize: 200, Ratio: 1},
			{ScreenSize: 300, Ratio: 0.25},
		},
		config.SidePanelLayoutHalfScreen: {
			{ScreenSize: -1, Ratio: 0.5},
		},
	}

	assert.Equal(t,
		map[config.SidePanelLayout][]config.SidePanelSize{
			config.SidePanelLayoutNormal: {
				{ScreenSize: 100, Ratio: 0.4},
				{ScreenSize: 300, Ratio: 0.25},
			},
		},
		ValidSidePanelSizes(input),
	)
}
