package authors

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetInitials(t *testing.T) {
	for input, expectedOutput := range map[string]string{
		"Jesse Duffield":     "JD",
		"Jesse Duffield Man": "JD",
		"JesseDuffield":      "Je",
		"J":                  "J",
		"六书六書":               "六",
		"書":                  "書",
		"":                   "",
	} {
		output := getInitials(input)
		if output != expectedOutput {
			t.Errorf("Expected %s to be %s", output, expectedOutput)
		}
	}
}

func TestAuthorWithLength(t *testing.T) {
	scenarios := []struct {
		authorName     string
		length         int
		expectedOutput string
	}{
		{"Jesse Duffield", 0, ""},
		{"Jesse Duffield", 1, ""},
		{"Jesse Duffield", 2, "JD"},
		{"Jesse Duffield", 3, "Je…"},
		{"Jesse Duffield", 10, "Jesse Duf…"},
		{"Jesse Duffield", 14, "Jesse Duffield"},
	}
	for _, s := range scenarios {
		assert.Equal(t, s.expectedOutput, utils.Decolorise(AuthorWithLength(s.authorName, s.length)))
	}
}
