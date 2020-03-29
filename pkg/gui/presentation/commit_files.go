package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetCommitFileListDisplayStrings(commitFiles []*commands.CommitFile, diffName string) [][]string {
	lines := make([][]string, len(commitFiles))

	for i := range commitFiles {
		diffed := commitFiles[i].Name == diffName
		lines[i] = getCommitFileDisplayStrings(commitFiles[i], diffed)
	}

	return lines
}

// getCommitFileDisplayStrings returns the display string of branch
func getCommitFileDisplayStrings(f *commands.CommitFile, diffed bool) []string {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	defaultColor := color.New(theme.DefaultTextColor)
	diffTerminalColor := color.New(theme.DiffTerminalColor)

	var colour *color.Color
	switch f.Status {
	case commands.UNSELECTED:
		colour = defaultColor
	case commands.WHOLE:
		colour = green
	case commands.PART:
		colour = yellow
	}
	if diffed {
		colour = diffTerminalColor
	}
	return []string{colour.Sprint(f.DisplayString)}
}
