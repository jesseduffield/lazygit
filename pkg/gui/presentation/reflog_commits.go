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

	var displayFunc func(*models.Commit, map[string]bool, bool, bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	for i := range commits {
		diffed := commits[i].Sha == diffName
		lines[i] = displayFunc(commits[i], cherryPickedCommitShaMap, diffed, parseEmoji)
	}

	return lines
}

func coloredReflogSha(c *models.Commit, cherryPickedCommitShaMap map[string]bool) string {
	shaColor := style.FgBlue
	if cherryPickedCommitShaMap[c.Sha] {
		shaColor = style.FgCyan.SetColor(style.BgBlue)
	}

	return shaColor.Sprint(c.ShortSha())
}

func getFullDescriptionDisplayStringsForReflogCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed, parseEmoji bool) []string {
	colorAttr := theme.DefaultTextColor
	if diffed {
		colorAttr = theme.DiffTerminalColor
	}

	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		coloredReflogSha(c, cherryPickedCommitShaMap),
		style.FgMagenta.Sprint(utils.UnixToDate(c.UnixTimestamp)),
		colorAttr.Sprint(name),
	}
}

func getDisplayStringsForReflogCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed, parseEmoji bool) []string {
	name := c.Name
	if parseEmoji {
		name = emoji.Sprint(name)
	}

	return []string{
		coloredReflogSha(c, cherryPickedCommitShaMap),
		theme.DefaultTextColor.Sprint(name),
	}
}
