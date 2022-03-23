package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/kyokomi/emoji/v2"
)

func GetReflogCommitListDisplayStrings(commits []*models.Commit, fullDescription bool, cherryPickedCommitShaMap map[string]bool, diffName string, parseEmoji bool) [][]string {
	lines := make([][]string, len(commits))

	var displayFunc func(*models.Commit, bool, bool, bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	for i := range commits {
		diffed := commits[i].Sha == diffName
		cherryPicked := cherryPickedCommitShaMap[commits[i].Sha]
		lines[i] = displayFunc(commits[i], cherryPicked, diffed, parseEmoji)
	}

	return lines
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

func getFullDescriptionDisplayStringsForReflogCommit(c *models.Commit, cherryPicked, diffed, parseEmoji bool) []string {
	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		reflogShaColor(cherryPicked, diffed).Sprint(c.ShortSha()),
		style.FgMagenta.Sprint(utils.UnixToDate(c.UnixTimestamp)),
		theme.DefaultTextColor.Sprint(name),
	}
}

func getDisplayStringsForReflogCommit(c *models.Commit, cherryPicked, diffed, parseEmoji bool) []string {
	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		reflogShaColor(cherryPicked, diffed).Sprint(c.ShortSha()),
		theme.DefaultTextColor.Sprint(name),
	}
}
