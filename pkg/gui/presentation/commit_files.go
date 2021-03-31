package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetCommitFileLine(name string, diffName string, commitFile *models.CommitFile, patchManager *patch.PatchManager, parent string) string {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	defaultColor := color.New(theme.DefaultTextColor)
	diffTerminalColor := color.New(theme.DiffTerminalColor)

	if commitFile == nil {
		return name
	}

	colour := defaultColor
	if diffName == name {
		colour = diffTerminalColor
	} else if commitFile != nil {
		status := patchManager.GetFileStatus(commitFile.Name, parent)
		switch status {
		case patch.UNSELECTED:
			colour = defaultColor
		case patch.WHOLE:
			colour = green
		case patch.PART:
			colour = yellow
		}
	}

	return utils.ColoredString(commitFile.ChangeStatus, getColorForChangeStatus(commitFile.ChangeStatus)) + " " + colour.Sprint(name)
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
