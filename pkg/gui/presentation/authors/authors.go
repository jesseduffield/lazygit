package authors

import (
	"crypto/md5"
	"strings"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/uniseg"
)

type authorNameCacheKey struct {
	authorName string
	truncateTo int
}

// if these being global variables causes trouble we can wrap them in a struct
// attached to the gui state.
var (
	authorInitialCache = make(map[string]string)
	authorNameCache    = make(map[authorNameCacheKey]string)

	// The styles from gui.authorColors
	customAuthorStyles = make(map[string]*style.TextStyle)
	// The styles derived from the names of the other authors
	authorStyleCache = make(map[string]*style.TextStyle)

	// Whether the terminal has a light background, for the derived styles to
	// stand out against
	lightBackground bool

	colorsVersion int
)

const authorNameWildcard = "*"

func ShortAuthor(authorName string) string {
	if value, ok := authorInitialCache[authorName]; ok {
		return value
	}

	initials := getInitials(authorName)
	if initials == "" {
		return ""
	}

	value := AuthorStyle(authorName).Sprint(initials)
	authorInitialCache[authorName] = value

	return value
}

func LongAuthor(authorName string, length int) string {
	cacheKey := authorNameCacheKey{authorName: authorName, truncateTo: length}
	if value, ok := authorNameCache[cacheKey]; ok {
		return value
	}

	paddedAuthorName := utils.WithPadding(authorName, length, utils.AlignLeft)
	truncatedName := utils.TruncateWithEllipsis(paddedAuthorName, length)
	value := AuthorStyle(authorName).Sprint(truncatedName)
	authorNameCache[cacheKey] = value

	return value
}

// AuthorWithLength returns a representation of the author that fits into a
// given maximum length:
// - if the length is less than 2, it returns an empty string
// - if the length is 2, it returns the initials
// - otherwise, it returns the author name truncated to the maximum length
func AuthorWithLength(authorName string, length int) string {
	if length < 2 {
		return ""
	}

	if length == 2 {
		return ShortAuthor(authorName)
	}

	return LongAuthor(authorName, length)
}

func AuthorStyle(authorName string) *style.TextStyle {
	if value, ok := customAuthorStyles[authorName]; ok {
		return value
	}

	// use the unified style whatever the author name is
	if value, ok := customAuthorStyles[authorNameWildcard]; ok {
		return value
	}

	if value, ok := authorStyleCache[authorName]; ok {
		return value
	}

	value := trueColorStyle(authorName)

	authorStyleCache[authorName] = &value

	return &value
}

func trueColorStyle(str string) style.TextStyle {
	c := colorAtPosition(ColorPosition(str))

	return style.New().SetFg(style.NewRGBColor(color.RGB(uint8(c.R*255), uint8(c.G*255), uint8(c.B*255))))
}

// To check the colors at the edges of the ranges below, run
// `go run ./cmd/author_colors_repo <path>` and open the repository it creates.
func colorAtPosition(hue, saturation, lightness float64) colorful.Color {
	// The lightness of an HSLuv color is how bright it looks, so every author
	// comes out about equally readable whichever hue their name lands on. Plain
	// HSL spreads them instead. At one and the same lightness, it gives a
	// glaring yellow and a blue that all but disappears.
	//
	// There is one lightness range for a dark background and one for a light
	// background. Each keeps every author above a contrast ratio of 4.5:1
	// against common backgrounds of its kind, such as #1e1e1e and #fdf6e3.
	//
	// Saturation in HSLuv is a fraction of the most colorful a hue can get at
	// that lightness, and pale colors are hard to tell apart, so keep it near
	// the top of its range.
	minLightness := 0.57
	if lightBackground {
		minLightness = 0.31
	}

	return colorful.HSLuv(hue*360.0, 0.8+0.2*saturation, minLightness+0.15*lightness)
}

// ColorPosition says where an author's color lies within the range of hues,
// saturations and lightnesses that colorAtPosition picks from. Each is a
// fraction from 0 up to 1, derived from a hash of the author's name.
func ColorPosition(authorName string) (hue, saturation, lightness float64) {
	hash := md5.Sum([]byte(authorName))
	return randFloat(hash[0:4]), randFloat(hash[4:8]), randFloat(hash[8:12])
}

func randFloat(hash []byte) float64 {
	return float64(randInt(hash, 100)) / 100
}

func randInt(hash []byte, max int) int {
	sum := 0
	for _, b := range hash {
		sum = (sum + int(b)) % max
	}
	return sum
}

func getInitials(authorName string) string {
	if authorName == "" {
		return authorName
	}

	firstChar, _, width, _ := uniseg.FirstGraphemeClusterInString(authorName, -1)
	if width > 1 {
		return firstChar
	}

	split := strings.Split(authorName, " ")
	if len(split) == 1 {
		return utils.LimitStr(authorName, 2)
	}

	return utils.LimitStr(split[0], 1) + utils.LimitStr(split[1], 1)
}

func SetCustomAuthors(customAuthorColors map[string]string) {
	customAuthorStyles = utils.SetCustomColors(customAuthorColors)
	colorsChanged()
}

// SetLightBackground says whether the terminal has a light background, for the
// colors of authors to stand out against.
func SetLightBackground(light bool) {
	if light == lightBackground {
		return
	}

	lightBackground = light
	authorStyleCache = make(map[string]*style.TextStyle)
	colorsChanged()
}

// colorsChanged drops what was rendered with the previous colors of authors.
func colorsChanged() {
	authorInitialCache = make(map[string]string)
	authorNameCache = make(map[authorNameCacheKey]string)
	colorsVersion++
}

// ColorsVersion changes whenever the colors of authors change, so that
// whatever keeps the styles of authors around can tell when they are out of
// date.
func ColorsVersion() int {
	return colorsVersion
}
