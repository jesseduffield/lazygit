package mergeconflicts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineLineType(t *testing.T) {
	// A markerSize of 0 means the file has no conflict-marker-size gitattribute,
	// so git's default size applies.
	type scenario struct {
		line       string
		markerSize int
		expected   LineType
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
		{
			line:     "<<<<<<<<",
			expected: NOT_A_MARKER,
		},
		// Markers without a label
		{
			line:     "<<<<<<<",
			expected: START,
		},
		{
			line:     "|||||||",
			expected: ANCESTOR,
		},
		{
			line:     ">>>>>>>",
			expected: END,
		},
		{
			line:       strings.Repeat("<", 32) + " HEAD",
			markerSize: 32,
			expected:   START,
		},
		{
			line:       strings.Repeat("|", 32) + " adf33b9",
			markerSize: 32,
			expected:   ANCESTOR,
		},
		{
			line:       strings.Repeat("=", 32),
			markerSize: 32,
			expected:   TARGET,
		},
		{
			line:       strings.Repeat(">", 32) + " blah",
			markerSize: 32,
			expected:   END,
		},
		// A file gets a bigger marker size precisely because its regular content
		// tends to contain marker-looking lines, so lines with the default size
		// must not be mistaken for markers
		{
			line:       "<<<<<<< HEAD",
			markerSize: 32,
			expected:   NOT_A_MARKER,
		},
		{
			line:       "=======",
			markerSize: 32,
			expected:   NOT_A_MARKER,
		},
		{
			line:       strings.Repeat("=", 33),
			markerSize: 32,
			expected:   NOT_A_MARKER,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, determineLineType(s.line, s.markerSize), s.line)
	}
}

func TestFindConflictsAux(t *testing.T) {
	// A markerSize of 0 means the file has no conflict-marker-size gitattribute,
	// so git's default size applies.
	type scenario struct {
		content    string
		markerSize int
		expected   bool
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
			content:  ">>>>>>>",
			expected: true,
		},
		{
			content:  "a\nb\nc\n<<<<<<< ",
			expected: true,
		},
		{
			content:    "a\nb\nc\n" + strings.Repeat("<", 32) + " HEAD",
			markerSize: 32,
			expected:   true,
		},
		{
			content:    "a\nb\nc\n" + strings.Repeat(">", 32) + " blah",
			markerSize: 32,
			expected:   true,
		},
		// Marker-looking lines of the default size are the file's regular content
		{
			content:    "a\nb\nc\n<<<<<<< HEAD\n=======\n>>>>>>> blah",
			markerSize: 32,
			expected:   false,
		},
	}

	for _, s := range scenarios {
		reader := strings.NewReader(s.content)
		result, err := fileHasConflictMarkersAux(reader, s.markerSize)
		assert.NoError(t, err)
		assert.EqualValues(t, s.expected, result, s.content)
	}
}
