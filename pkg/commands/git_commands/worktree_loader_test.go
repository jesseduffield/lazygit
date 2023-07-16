package git_commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUniqueNamesFromPaths(t *testing.T) {
	for _, scenario := range []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{},
			expected: []string{},
		},
		{
			input: []string{
				"/my/path/feature/one",
			},
			expected: []string{
				"one",
			},
		},
		{
			input: []string{
				"/my/path/feature/one/",
			},
			expected: []string{
				"one",
			},
		},
		{
			input: []string{
				"/a/b/c/d",
				"/a/b/c/e",
				"/a/b/f/d",
				"/a/e/c/d",
			},
			expected: []string{
				"b/c/d",
				"e",
				"f/d",
				"e/c/d",
			},
		},
	} {
		actual := getUniqueNamesFromPaths(scenario.input)
		assert.EqualValues(t, scenario.expected, actual)
	}
}
