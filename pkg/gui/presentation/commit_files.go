package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetCommitFileListDisplayStrings(branches []*commands.CommitFile) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		lines[i] = getCommitFileDisplayStrings(branches[i])
	}

	return lines
}

// getCommitFileDisplayStrings returns the display string of branch
func getCommitFileDisplayStrings(f *commands.CommitFile) []string {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	defaultColor := color.New(theme.DefaultTextColor)

	var colour *color.Color
	switch f.Status {
	case commands.UNSELECTED:
		colour = defaultColor
	case commands.WHOLE:
		colour = green
	case commands.PART:
		colour = yellow
	}
	return []string{colour.Sprint(f.DisplayString)}
}
