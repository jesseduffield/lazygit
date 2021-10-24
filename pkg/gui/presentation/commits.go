package presentation

import (
	"crypto/md5"
	"strings"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/mattn/go-runewidth"
)

func GetCommitListDisplayStrings(commits []*models.Commit, fullDescription bool, cherryPickedCommitShaMap map[string]bool, diffName string, parseEmoji bool) [][]string {
	lines := make([][]string, len(commits))

	var displayFunc func(*models.Commit, map[string]bool, bool, bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForCommit
	} else {
		displayFunc = getDisplayStringsForCommit
	}

	for i := range commits {
		diffed := commits[i].Sha == diffName
		lines[i] = displayFunc(commits[i], cherryPickedCommitShaMap, diffed, parseEmoji)
	}

	return lines
}

func getFullDescriptionDisplayStringsForCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed, parseEmoji bool) []string {
	shaColor := theme.DefaultTextColor
	switch c.Status {
	case "unpushed":
		shaColor = style.FgRed
	case "pushed":
		shaColor = style.FgYellow
	case "merged":
		shaColor = style.FgGreen
	case "rebasing":
		shaColor = style.FgBlue
	case "reflog":
		shaColor = style.FgBlue
	}

	if diffed {
		shaColor = theme.DiffTerminalColor
	} else if cherryPickedCommitShaMap[c.Sha] {
		// for some reason, setting the background to blue pads out the other commits
		// horizontally. For the sake of accessibility I'm considering this a feature,
		// not a bug
		shaColor = theme.CherryPickedCommitTextStyle
	}

	tagString := ""
	secondColumnString := style.FgBlue.Sprint(utils.UnixToDate(c.UnixTimestamp))
	if c.Action != "" {
		secondColumnString = actionColorMap(c.Action).Sprint(c.Action)
	} else if c.ExtraInfo != "" {
		tagString = style.FgMagenta.SetBold().Sprint(c.ExtraInfo) + " "
	}

	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		shaColor.Sprint(c.ShortSha()),
		secondColumnString,
		longAuthor(c.Author),
		tagString + theme.DefaultTextColor.Sprint(name),
	}
}

func getDisplayStringsForCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed, parseEmoji bool) []string {
	shaColor := theme.DefaultTextColor
	switch c.Status {
	case "unpushed":
		shaColor = style.FgRed
	case "pushed":
		shaColor = style.FgYellow
	case "merged":
		shaColor = style.FgGreen
	case "rebasing":
		shaColor = style.FgBlue
	case "reflog":
		shaColor = style.FgBlue
	}

	if diffed {
		shaColor = theme.DiffTerminalColor
	} else if cherryPickedCommitShaMap[c.Sha] {
		// for some reason, setting the background to blue pads out the other commits
		// horizontally. For the sake of accessibility I'm considering this a feature,
		// not a bug
		shaColor = theme.CherryPickedCommitTextStyle
	}

	actionString := ""
	tagString := ""
	if c.Action != "" {
		actionString = actionColorMap(c.Action).Sprint(utils.WithPadding(c.Action, 7)) + " "
	} else if len(c.Tags) > 0 {
		tagString = theme.DiffTerminalColor.SetBold().Sprint(strings.Join(c.Tags, " ")) + " "
	}

	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		shaColor.Sprint(c.ShortSha()),
		shortAuthor(c.Author),
		actionString + tagString + theme.DefaultTextColor.Sprint(name),
	}
}

var authorInitialCache = make(map[string]string)
var authorNameCache = make(map[string]string)

func shortAuthor(authorName string) string {
	if value, ok := authorInitialCache[authorName]; ok {
		return value
	}

	initials := getInitials(authorName)
	if initials == "" {
		return ""
	}

	value := authorColor(authorName).Sprint(initials)
	authorInitialCache[authorName] = value

	return value
}

func longAuthor(authorName string) string {
	if value, ok := authorNameCache[authorName]; ok {
		return value
	}

	truncatedName := utils.TruncateWithEllipsis(authorName, 17)
	value := authorColor(authorName).Sprint(truncatedName)
	authorNameCache[authorName] = value

	return value
}

func authorColor(authorName string) style.TextStyle {
	hash := md5.Sum([]byte(authorName))
	c := colorful.Hsl(randFloat(hash[0:4])*360.0, 0.6+0.4*randFloat(hash[4:8]), 0.4+randFloat(hash[8:12])*0.2)

	return style.New().SetFg(style.NewRGBColor(color.RGB(uint8(c.R*255), uint8(c.G*255), uint8(c.B*255))))
}

func randFloat(hash []byte) float64 {
	sum := 0
	for _, b := range hash {
		sum = (sum + int(b)) % 100
	}
	return float64(sum) / 100
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

func actionColorMap(str string) style.TextStyle {
	switch str {
	case "pick":
		return style.FgCyan
	case "drop":
		return style.FgRed
	case "edit":
		return style.FgGreen
	case "fixup":
		return style.FgMagenta
	default:
		return style.FgYellow
	}
}
