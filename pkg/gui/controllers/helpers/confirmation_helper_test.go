package helpers

import (
	"testing"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func Test_underlineLinks(t *testing.T) {
	scenarios := []struct {
		name           string
		text           string
		expectedResult string
	}{
		{
			name:           "empty string",
			text:           "",
			expectedResult: "",
		},
		{
			name:           "no links",
			text:           "abc",
			expectedResult: "abc",
		},
		{
			name:           "entire string is a link",
			text:           "https://example.com",
			expectedResult: "\x1b]8;;https://example.com\x1b\\https://example.com\x1b]8;;\x1b\\",
		},
		{
			name:           "link preceded and followed by text",
			text:           "bla https://example.com xyz",
			expectedResult: "bla \x1b]8;;https://example.com\x1b\\https://example.com\x1b]8;;\x1b\\ xyz",
		},
		{
			name:           "more than one link",
			text:           "bla https://link1 blubb https://link2 xyz",
			expectedResult: "bla \x1b]8;;https://link1\x1b\\https://link1\x1b]8;;\x1b\\ blubb \x1b]8;;https://link2\x1b\\https://link2\x1b]8;;\x1b\\ xyz",
		},
		{
			name:           "link in angle brackets",
			text:           "See <https://example.com> for details",
			expectedResult: "See <\x1b]8;;https://example.com\x1b\\https://example.com\x1b]8;;\x1b\\> for details",
		},
		{
			name:           "link followed by newline",
			text:           "URL: https://example.com\nNext line",
			expectedResult: "URL: \x1b]8;;https://example.com\x1b\\https://example.com\x1b]8;;\x1b\\\nNext line",
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result := underlineLinks(s.text)
			assert.Equal(t, s.expectedResult, result)
		})
	}
}
