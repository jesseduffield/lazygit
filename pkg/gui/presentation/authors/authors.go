package authors

import (
	"crypto/md5"
	"strings"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mattn/go-runewidth"
)

// if these being global variables causes trouble we can wrap them in a struct
// attached to the gui state.
var (
	authorInitialCache = make(map[string]string)
	authorNameCache    = make(map[string]string)
	authorStyleCache   = make(map[string]style.TextStyle)
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

func LongAuthor(authorName string) string {
	if value, ok := authorNameCache[authorName]; ok {
		return value
	}

	paddedAuthorName := utils.WithPadding(authorName, 17)
	truncatedName := utils.TruncateWithEllipsis(paddedAuthorName, 17)
	value := AuthorStyle(authorName).Sprint(truncatedName)
	authorNameCache[authorName] = value

	return value
}

func AuthorStyle(authorName string) style.TextStyle {
	if value, ok := authorStyleCache[authorName]; ok {
		return value
	}

	// use the unified style whatever the author name is
	if value, ok := authorStyleCache[authorNameWildcard]; ok {
		return value
	}

	value := trueColorStyle(authorName)

	authorStyleCache[authorName] = value

	return value
}

func trueColorStyle(str string) style.TextStyle {
	hash := md5.Sum([]byte(str))
	c := colorful.Hsl(randFloat(hash[0:4])*360.0, 0.6+0.4*randFloat(hash[4:8]), 0.4+randFloat(hash[8:12])*0.2)

	return style.New().SetFg(style.NewRGBColor(color.RGB(uint8(c.R*255), uint8(c.G*255), uint8(c.B*255))))
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

	firstRune := getFirstRune(authorName)
	if runewidth.RuneWidth(firstRune) > 1 {
		return string(firstRune)
	}

	split := strings.Split(authorName, " ")
	if len(split) == 1 {
		return utils.LimitStr(authorName, 2)
	}

	return utils.LimitStr(split[0], 1) + utils.LimitStr(split[1], 1)
}

func getFirstRune(str string) rune {
	// just using the loop for the sake of getting the first rune
	for _, r := range str {
		return r
	}
	// should never land here
	return 0
}

func SetCustomAuthors(customAuthorColors map[string]string) {
	authorStyleCache = utils.SetCustomColors(customAuthorColors)
}
