package patch_exploring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrigin(t *testing.T) {
	type scenario struct {
		name            string
		origin          int
		bufferHeight    int
		firstLineIdx    int
		lastLineIdx     int
		selectedLineIdx int
		selectMode      selectMode
		expected        int
	}

	scenarios := []scenario{
		{
			name:            "selection above scroll window",
			origin:          50,
			bufferHeight:    100,
			firstLineIdx:    10,
			lastLineIdx:     10,
			selectedLineIdx: 10,
			selectMode:      LINE,
			expected:        10,
		},
		{
			name:            "selection below scroll window",
			origin:          0,
			bufferHeight:    100,
			firstLineIdx:    150,
			lastLineIdx:     150,
			selectedLineIdx: 150,
			selectMode:      LINE,
			expected:        50,
		},
		{
			name:            "selection within scroll window",
			origin:          0,
			bufferHeight:    100,
			firstLineIdx:    50,
			lastLineIdx:     50,
			selectedLineIdx: 50,
			selectMode:      LINE,
			expected:        0,
		},
		{
			name:            "range ending below scroll window with selection at end of range",
			origin:          0,
			bufferHeight:    100,
			firstLineIdx:    40,
			lastLineIdx:     150,
			selectedLineIdx: 150,
			selectMode:      RANGE,
			expected:        50,
		},
		{
			name:            "range ending below scroll window with selection at beginning of range",
			origin:          0,
			bufferHeight:    100,
			firstLineIdx:    40,
			lastLineIdx:     150,
			selectedLineIdx: 40,
			selectMode:      RANGE,
			expected:        40,
		},
		{
			name:            "range starting above scroll window with selection at beginning of range",
			origin:          50,
			bufferHeight:    100,
			firstLineIdx:    40,
			lastLineIdx:     150,
			selectedLineIdx: 40,
			selectMode:      RANGE,
			expected:        40,
		},
		{
			name:            "hunk extending beyond both bounds of scroll window",
			origin:          50,
			bufferHeight:    100,
			firstLineIdx:    40,
			lastLineIdx:     200,
			selectedLineIdx: 70,
			selectMode:      HUNK,
			expected:        40,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			assert.EqualValues(t, s.expected, calculateOrigin(s.origin, s.bufferHeight, s.firstLineIdx, s.lastLineIdx, s.selectedLineIdx, s.selectMode))
		})
	}
}
