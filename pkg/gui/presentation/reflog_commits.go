package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetReflogCommitListDisplayStrings(commits []*commands.Commit, fullDescription bool, diffName string) [][]string {
	lines := make([][]string, len(commits))

	var displayFunc func(*commands.Commit, bool) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	for i := range commits {
		diffed := commits[i].Sha == diffName
		lines[i] = displayFunc(commits[i], diffed)
	}

	return lines
}

func getFullDescriptionDisplayStringsForReflogCommit(c *commands.Commit, diffed bool) []string {
	colorAttr := theme.DefaultTextColor
	if diffed {
		colorAttr = theme.DiffTerminalColor
	}

	return []string{
		utils.ColoredString(c.ShortSha(), color.FgBlue),
		utils.ColoredString(utils.UnixToDate(c.UnixTimestamp), color.FgMagenta),
		utils.ColoredString(c.Name, colorAttr),
	}
}

func getDisplayStringsForReflogCommit(c *commands.Commit, diffed bool) []string {
	defaultColor := color.New(theme.DefaultTextColor)

	return []string{utils.ColoredString(c.ShortSha(), color.FgBlue), defaultColor.Sprint(c.Name)}
}
