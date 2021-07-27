package presentation

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
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
		shaColor = style.FgCyan.SetColor(style.BgBlue)
	}

	tagString := ""
	secondColumnString := style.FgBlue.Sprint(utils.UnixToDate(c.UnixTimestamp))
	if c.Action != "" {
		secondColumnString = actionColorMap(c.Action).Sprint(c.Action)
	} else if c.ExtraInfo != "" {
		tagString = theme.DiffTerminalColor.SetBold(true).Sprint(c.ExtraInfo) + " "
	}

	truncatedAuthor := utils.TruncateWithEllipsis(c.Author, 17)

	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		shaColor.Sprint(c.ShortSha()),
		secondColumnString,
		style.FgYellow.Sprint(truncatedAuthor),
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
		shaColor = style.FgCyan.SetColor(style.BgBlue)
	}

	actionString := ""
	tagString := ""
	if c.Action != "" {
		actionString = actionColorMap(c.Action).Sprint(utils.WithPadding(c.Action, 7)) + " "
	} else if len(c.Tags) > 0 {
		tagString = theme.DiffTerminalColor.SetBold(true).Sprint(strings.Join(c.Tags, " ")) + " "
	}

	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		shaColor.Sprint(c.ShortSha()),
		actionString + tagString + theme.DefaultTextColor.Sprint(name),
	}
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
