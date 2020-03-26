package presentation

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetCommitListDisplayStrings(commits []*commands.Commit, fullDescription bool, cherryPickedCommitShaMap map[string]bool) [][]string {
	lines := make([][]string, len(commits))

	var displayFunc func(*commands.Commit, map[string]bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForCommit
	} else {
		displayFunc = getDisplayStringsForCommit
	}

	for i := range commits {
		lines[i] = displayFunc(commits[i], cherryPickedCommitShaMap)
	}

	return lines
}

func getFullDescriptionDisplayStringsForCommit(c *commands.Commit, cherryPickedCommitShaMap map[string]bool) []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)
	defaultColor := color.New(theme.DefaultTextColor)
	magenta := color.New(color.FgMagenta)

	// for some reason, setting the background to blue pads out the other commits
	// horizontally. For the sake of accessibility I'm considering this a feature,
	// not a bug
	copied := color.New(color.FgCyan, color.BgBlue)

	var shaColor *color.Color
	switch c.Status {
	case "unpushed":
		shaColor = red
	case "pushed":
		shaColor = yellow
	case "merged":
		shaColor = green
	case "rebasing":
		shaColor = blue
	case "reflog":
		shaColor = blue
	case "selected":
		shaColor = magenta
	default:
		shaColor = defaultColor
	}

	if cherryPickedCommitShaMap[c.Sha] {
		shaColor = copied
	}

	tagString := ""
	truncatedDate := utils.TruncateWithEllipsis(c.Date, 15)
	secondColumnString := blue.Sprint(truncatedDate)
	if c.Action != "" {
		secondColumnString = cyan.Sprint(c.Action)
	} else if c.ExtraInfo != "" {
		tagColor := color.New(color.FgMagenta, color.Bold)
		tagString = utils.ColoredStringDirect(c.ExtraInfo, tagColor) + " "
	}

	truncatedAuthor := utils.TruncateWithEllipsis(c.Author, 17)

	return []string{shaColor.Sprint(c.ShortSha()), secondColumnString, yellow.Sprint(truncatedAuthor), tagString + defaultColor.Sprint(c.Name)}
}

func getDisplayStringsForCommit(c *commands.Commit, cherryPickedCommitShaMap map[string]bool) []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)
	defaultColor := color.New(theme.DefaultTextColor)
	magenta := color.New(color.FgMagenta)

	// for some reason, setting the background to blue pads out the other commits
	// horizontally. For the sake of accessibility I'm considering this a feature,
	// not a bug
	copied := color.New(color.FgCyan, color.BgBlue)

	var shaColor *color.Color
	switch c.Status {
	case "unpushed":
		shaColor = red
	case "pushed":
		shaColor = yellow
	case "merged":
		shaColor = green
	case "rebasing":
		shaColor = blue
	case "reflog":
		shaColor = blue
	case "selected":
		shaColor = magenta
	default:
		shaColor = defaultColor
	}

	if cherryPickedCommitShaMap[c.Sha] {
		shaColor = copied
	}

	actionString := ""
	tagString := ""
	if c.Action != "" {
		actionString = cyan.Sprint(utils.WithPadding(c.Action, 7)) + " "
	} else if len(c.Tags) > 0 {
		tagColor := color.New(color.FgMagenta, color.Bold)
		tagString = utils.ColoredStringDirect(strings.Join(c.Tags, " "), tagColor) + " "
	}

	return []string{shaColor.Sprint(c.ShortSha()), actionString + tagString + defaultColor.Sprint(c.Name)}
}
