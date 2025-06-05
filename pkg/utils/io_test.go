package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_forEachLineInStream(t *testing.T) {
	scenarios := []struct {
		name          string
		input         string
		expectedLines []string
	}{
		{
			name:          "empty input",
			input:         "",
			expectedLines: []string{},
		},
		{
			name:          "single line",
			input:         "abc\n",
			expectedLines: []string{"abc\n"},
		},
		{
			name:          "single line without line feed",
			input:         "abc",
			expectedLines: []string{"abc"},
		},
		{
			name:          "multiple lines",
			input:         "abc\ndef\n",
			expectedLines: []string{"abc\n", "def\n"},
		},
		{
			name:          "multiple lines including empty lines",
			input:         "abc\n\ndef\n",
			expectedLines: []string{"abc\n", "\n", "def\n"},
		},
		{
			name:          "multiple lines without linefeed at end of file",
			input:         "abc\ndef\nghi",
			expectedLines: []string{"abc\n", "def\n", "ghi"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			lines := []string{}
			forEachLineInStream(strings.NewReader(s.input), func(line string, i int) {
				lines = append(lines, line)
			})
			assert.EqualValues(t, s.expectedLines, lines)
		})
	}
}
