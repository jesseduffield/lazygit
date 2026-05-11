package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
