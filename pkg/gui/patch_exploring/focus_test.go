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
		numLines        int
		firstLineIdx    int
		lastLineIdx     int
		selectedLineIdx int
		selectMode      selectMode
		expected        int
	}

	scenarios := []scenario{
		{
			name:            "selection above scroll window, enough room to put it in the middle",
			origin:          250,
			bufferHeight:    100,
			numLines:        500,
			firstLineIdx:    210,
			lastLineIdx:     210,
			selectedLineIdx: 210,
			selectMode:      LINE,
			expected:        160,
		},
		{
			name:            "selection above scroll window, not enough room to put it in the middle",
			origin:          50,
			bufferHeight:    100,
			numLines:        500,
			firstLineIdx:    10,
			lastLineIdx:     10,
			selectedLineIdx: 10,
			selectMode:      LINE,
			expected:        0,
		},
		{
			name:            "selection below scroll window, enough room to put it in the middle",
			origin:          0,
			bufferHeight:    100,
			numLines:        500,
			firstLineIdx:    150,
			lastLineIdx:     150,
			selectedLineIdx: 150,
			selectMode:      LINE,
			expected:        100,
		},
		{
			name:            "selection below scroll window, not enough room to put it in the middle",
			origin:          0,
			bufferHeight:    100,
			numLines:        200,
			firstLineIdx:    199,
			lastLineIdx:     199,
			selectedLineIdx: 199,
			selectMode:      LINE,
			expected:        99,
		},
		{
			name:            "selection within scroll window",
			origin:          0,
			bufferHeight:    100,
			numLines:        500,
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
			numLines:        500,
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
			numLines:        500,
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
			numLines:        500,
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
			numLines:        500,
			firstLineIdx:    40,
			lastLineIdx:     200,
			selectedLineIdx: 70,
			selectMode:      HUNK,
			expected:        40,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.EqualValues(t, s.expected, calculateOrigin(s.origin, s.bufferHeight, s.numLines, s.firstLineIdx, s.lastLineIdx, s.selectedLineIdx, s.selectMode))
		})
	}
}
