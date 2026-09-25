package authors

import (
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
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

func TestAuthorColorsFollowTheConfig(t *testing.T) {
	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)
	t.Cleanup(func() { SetCustomAuthors(nil) })

	SetCustomAuthors(map[string]string{"Jane Doe": "red"})
	assert.Equal(t, style.FgRed.Sprint("JD"), ShortAuthor("Jane Doe"))
	assert.Equal(t, style.FgRed.Sprint("Jane Doe"), LongAuthor("Jane Doe", 8))

	SetCustomAuthors(map[string]string{"Jane Doe": "blue"})
	/* EXPECTED:
	assert.Equal(t, style.FgBlue.Sprint("JD"), ShortAuthor("Jane Doe"))
	ACTUAL: */
	assert.Equal(t, style.FgRed.Sprint("JD"), ShortAuthor("Jane Doe"))
	/* EXPECTED:
	assert.Equal(t, style.FgBlue.Sprint("Jane Doe"), LongAuthor("Jane Doe", 8))
	ACTUAL: */
	assert.Equal(t, style.FgRed.Sprint("Jane Doe"), LongAuthor("Jane Doe", 8))
}
