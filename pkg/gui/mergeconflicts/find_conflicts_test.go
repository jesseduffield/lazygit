package mergeconflicts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineLineType(t *testing.T) {
	type scenario struct {
		line     string
		expected LineType
	}

	scenarios := []scenario{
		{
			line:     "",
			expected: NOT_A_MARKER,
		},
		{
			line:     "blah",
			expected: NOT_A_MARKER,
		},
		{
			line:     "<<<<<<< HEAD",
			expected: START,
		},
		{
			line:     "<<<<<<< HEAD:my_branch",
			expected: START,
		},
		{
			line:     "<<<<<<< MERGE_HEAD:my_branch",
			expected: START,
		},
		{
			line:     "<<<<<<< Updated upstream:my_branch",
			expected: START,
		},
		{
			line:     "<<<<<<< ours:my_branch",
			expected: START,
		},
		{
			line:     "=======",
			expected: MIDDLE,
		},
		{
			line:     ">>>>>>> blah",
			expected: END,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, determineLineType(s.line))
	}
}
