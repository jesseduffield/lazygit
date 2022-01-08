package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetCommitFileLine(name string, diffName string, commitFile *models.CommitFile, status patch.PatchStatus) string {
	var colour style.TextStyle
	if diffName == name {
		colour = theme.DiffTerminalColor
	} else {
		switch status {
		case patch.WHOLE:
			colour = style.FgGreen
		case patch.PART:
			colour = style.FgYellow
		case patch.UNSELECTED:
			colour = theme.DefaultTextColor
		}
	}

	name = utils.EscapeSpecialChars(name)
	if commitFile == nil {
		return colour.Sprint(name)
	}

	return getColorForChangeStatus(commitFile.ChangeStatus).Sprint(commitFile.ChangeStatus) + " " + colour.Sprint(name)
}

func getColorForChangeStatus(changeStatus string) style.TextStyle {
	switch changeStatus {
	case "A":
		return style.FgGreen
	case "M", "R":
		return style.FgYellow
	case "D":
		return style.FgRed
	case "C":
		return style.FgCyan
	case "T":
		return style.FgMagenta
	default:
		return theme.DefaultTextColor
	}
}
