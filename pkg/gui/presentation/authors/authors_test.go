package authors

import (
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/lucasb-eyer/go-colorful"
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
	assert.Equal(t, style.FgBlue.Sprint("JD"), ShortAuthor("Jane Doe"))
	assert.Equal(t, style.FgBlue.Sprint("Jane Doe"), LongAuthor("Jane Doe", 8))
}

func TestSetLightBackground(t *testing.T) {
	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)
	t.Cleanup(func() {
		SetLightBackground(false)
		SetCustomAuthors(nil)
	})

	SetCustomAuthors(map[string]string{"Jane Doe": "red"})
	onDarkBackground := ShortAuthor("John Smith")

	SetLightBackground(true)
	assert.NotEqual(t, onDarkBackground, ShortAuthor("John Smith"))
	assert.Equal(t, style.FgRed.Sprint("JD"), ShortAuthor("Jane Doe"))
}

func TestAuthorColorsStandOutAgainstTheBackground(t *testing.T) {
	t.Cleanup(func() { SetLightBackground(false) })

	scenarios := []struct {
		name            string
		lightBackground bool
		backgrounds     []string
	}{
		{name: "dark", lightBackground: false, backgrounds: []string{"#000000", "#1e1e1e"}},
		{name: "light", lightBackground: true, backgrounds: []string{"#ffffff", "#fdf6e3"}},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			SetLightBackground(s.lightBackground)

			for _, backgroundHex := range s.backgrounds {
				background, err := colorful.Hex(backgroundHex)
				assert.NoError(t, err)

				// The edges of the range, which is where the contrast is lowest
				for hue := range 100 {
					for _, saturation := range []float64{0, 0.99} {
						for _, lightness := range []float64{0, 0.99} {
							c := colorAtPosition(float64(hue)/100, saturation, lightness)
							assert.GreaterOrEqual(t, contrastRatio(c, background), 4.5,
								"%s on %s", c.Hex(), backgroundHex)
						}
					}
				}
			}
		})
	}
}

// contrastRatio is as defined by the Web Content Accessibility Guidelines
func contrastRatio(a colorful.Color, b colorful.Color) float64 {
	luminance := func(c colorful.Color) float64 {
		r, g, b := c.LinearRgb()
		return 0.2126*r + 0.7152*g + 0.0722*b
	}
	lighter := max(luminance(a), luminance(b))
	darker := min(luminance(a), luminance(b))
	return (lighter + 0.05) / (darker + 0.05)
}
