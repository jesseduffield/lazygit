package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetCommitFileListDisplayStrings(commitFiles []*models.CommitFile, diffName string) [][]string {
	if len(commitFiles) == 0 {
		return [][]string{{utils.ColoredString("(none)", color.FgRed)}}
	}

	lines := make([][]string, len(commitFiles))

	for i := range commitFiles {
		diffed := commitFiles[i].Name == diffName
		lines[i] = getCommitFileDisplayStrings(commitFiles[i], diffed)
	}

	return lines
}

// getCommitFileDisplayStrings returns the display string of branch
func getCommitFileDisplayStrings(f *models.CommitFile, diffed bool) []string {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	defaultColor := color.New(theme.DefaultTextColor)
	diffTerminalColor := color.New(theme.DiffTerminalColor)

	var colour *color.Color
	switch f.PatchStatus {
	case patch.UNSELECTED:
		colour = defaultColor
	case patch.WHOLE:
		colour = green
	case patch.PART:
		colour = yellow
	}
	if diffed {
		colour = diffTerminalColor
	}
	return []string{utils.ColoredString(f.ChangeStatus, getColorForChangeStatus(f.ChangeStatus)), colour.Sprint(f.Name)}
}

func getColorForChangeStatus(changeStatus string) color.Attribute {
	switch changeStatus {
	case "A":
		return color.FgGreen
	case "M", "R":
		return color.FgYellow
	case "D":
		return color.FgRed
	case "C":
		return color.FgCyan
	case "T":
		return color.FgMagenta
	default:
		return theme.DefaultTextColor
	}
}
