package mergeconflicts

import (
	"strings"
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
			expected: TARGET,
		},
		{
			line:     ">>>>>>> blah",
			expected: END,
		},
		{
			line:     "||||||| adf33b9",
			expected: ANCESTOR,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, determineLineType(s.line))
	}
}

func TestFindConflictsAux(t *testing.T) {
	type scenario struct {
		content  string
		expected bool
	}

	scenarios := []scenario{
		{
			content:  "",
			expected: false,
		},
		{
			content:  "blah",
			expected: false,
		},
		{
			content:  ">>>>>>> ",
			expected: true,
		},
		{
			content:  "<<<<<<< ",
			expected: true,
		},
		{
			content:  " <<<<<<< ",
			expected: false,
		},
		{
			content:  "a\nb\nc\n<<<<<<< ",
			expected: true,
		},
	}

	for _, s := range scenarios {
		reader := strings.NewReader(s.content)
		assert.EqualValues(t, s.expected, fileHasConflictMarkersAux(reader))
	}
}
