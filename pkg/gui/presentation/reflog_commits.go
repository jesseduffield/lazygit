package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetReflogCommitListDisplayStrings(commits []*commands.Commit, fullDescription bool) [][]string {
	lines := make([][]string, len(commits))

	var displayFunc func(*commands.Commit) []string
	if fullDescription {
		displayFunc = getFullDescriptionDisplayStringsForReflogCommit
	} else {
		displayFunc = getDisplayStringsForReflogCommit
	}

	for i := range commits {
		lines[i] = displayFunc(commits[i])
	}

	return lines
}

func getFullDescriptionDisplayStringsForReflogCommit(c *commands.Commit) []string {
	defaultColor := color.New(theme.DefaultTextColor)

	return []string{utils.ColoredString(c.ShortSha(), color.FgBlue), utils.ColoredString(c.Date, color.FgMagenta), defaultColor.Sprint(c.Name)}
}

func getDisplayStringsForReflogCommit(c *commands.Commit) []string {
	defaultColor := color.New(theme.DefaultTextColor)

	return []string{utils.ColoredString(c.ShortSha(), color.FgBlue), defaultColor.Sprint(c.Name)}
}
