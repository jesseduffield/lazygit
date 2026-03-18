package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAIOutput(t *testing.T) {
	scenarios := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text unchanged",
			input:    "feat: add login button",
			expected: "feat: add login button",
		},
		{
			name:     "strips surrounding whitespace",
			input:    "  feat: add login button  ",
			expected: "feat: add login button",
		},
		{
			name:     "strips plain code fence",
			input:    "```\nfeat: add login button\n```",
			expected: "feat: add login button",
		},
		{
			name:     "strips fenced block with language tag",
			input:    "```text\nfeat: add login button\n```",
			expected: "feat: add login button",
		},
		{
			name:     "strips fenced block with surrounding text",
			input:    "Here is your commit:\n```\nfeat: add login button\n```\nDone.",
			expected: "feat: add login button",
		},
		{
			name:     "multiline commit message in fence",
			input:    "```\nfeat: add login button\n\nCloses #123\n```",
			expected: "feat: add login button\n\nCloses #123",
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, parseAIOutput(s.input))
		})
	}
}

func TestTryRemoveHardLineBreaks(t *testing.T) {
	scenarios := []struct {
		name           string
		message        string
		autoWrapWidth  int
		expectedResult string
	}{
		{
			name:           "empty",
			message:        "",
			autoWrapWidth:  7,
			expectedResult: "",
		},
		{
			name:           "all line breaks are needed",
			message:        "abc\ndef\n\nxyz",
			autoWrapWidth:  7,
			expectedResult: "abc\ndef\n\nxyz",
		},
		{
			name:           "some can be unwrapped",
			message:        "123\nabc def\nghi jkl\nmno\n456\n",
			autoWrapWidth:  7,
			expectedResult: "123\nabc def ghi jkl mno\n456\n",
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			actualResult := TryRemoveHardLineBreaks(s.message, s.autoWrapWidth)
			assert.Equal(t, s.expectedResult, actualResult)
		})
	}
}
