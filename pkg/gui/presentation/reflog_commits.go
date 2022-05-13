package presentation

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
)

func GetReflogCommitListDisplayStrings(commits []*models.Commit, fullDescription bool, cherryPickedCommitShaSet *set.Set[string], diffName string, timeFormat string, parseEmoji bool) [][]string {
	var displayFunc func(*models.Commit, string, bool, bool, bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	return slices.Map(commits, func(commit *models.Commit) []string {
		diffed := commit.Sha == diffName
		cherryPicked := cherryPickedCommitShaSet.Includes(commit.Sha)
		return displayFunc(commit, timeFormat, cherryPicked, diffed, parseEmoji)
	})
}

func reflogShaColor(cherryPicked, diffed bool) style.TextStyle {
	if diffed {
		return theme.DiffTerminalColor
	}

	shaColor := style.FgBlue
	if cherryPicked {
		shaColor = theme.CherryPickedCommitTextStyle
	}

	return shaColor
}

func getFullDescriptionDisplayStringsForReflogCommit(c *models.Commit, timeFormat string, cherryPicked, diffed, parseEmoji bool) []string {
	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		reflogShaColor(cherryPicked, diffed).Sprint(c.ShortSha()),
		style.FgMagenta.Sprint(utils.UnixToDate(c.UnixTimestamp, timeFormat)),
		theme.DefaultTextColor.Sprint(name),
	}
}

func getDisplayStringsForReflogCommit(c *models.Commit, timeFormat string, cherryPicked, diffed, parseEmoji bool) []string {
	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		reflogShaColor(cherryPicked, diffed).Sprint(c.ShortSha()),
		theme.DefaultTextColor.Sprint(name),
	}
}
