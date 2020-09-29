package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetReflogCommitListDisplayStrings(commits []*models.Commit, fullDescription bool, cherryPickedCommitShaMap map[string]bool, diffName string) [][]string {
	lines := make([][]string, len(commits))

	var displayFunc func(*models.Commit, map[string]bool, bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	for i := range commits {
		diffed := commits[i].Sha == diffName
		lines[i] = displayFunc(commits[i], cherryPickedCommitShaMap, diffed)
	}

	return lines
}

func coloredReflogSha(c *models.Commit, cherryPickedCommitShaMap map[string]bool) string {
	var shaColor *color.Color
	if cherryPickedCommitShaMap[c.Sha] {
		shaColor = color.New(color.FgCyan, color.BgBlue)
	} else {
		shaColor = color.New(color.FgBlue)
	}

	return shaColor.Sprint(c.ShortSha())
}

func getFullDescriptionDisplayStringsForReflogCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed bool) []string {
	colorAttr := theme.DefaultTextColor
	if diffed {
		colorAttr = theme.DiffTerminalColor
	}

	return []string{
		coloredReflogSha(c, cherryPickedCommitShaMap),
		utils.ColoredString(utils.UnixToDate(c.UnixTimestamp), color.FgMagenta),
		utils.ColoredString(c.Name, colorAttr),
	}
}

func getDisplayStringsForReflogCommit(c *models.Commit, cherryPickedCommitShaMap map[string]bool, diffed bool) []string {
	defaultColor := color.New(theme.DefaultTextColor)

	return []string{
		coloredReflogSha(c, cherryPickedCommitShaMap),
		defaultColor.Sprint(c.Name),
	}
}
