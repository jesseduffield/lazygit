package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// TODO: move this back into presentation package and fix the import cycle

func getFileLine(hasUnstagedChanges bool, hasStagedChanges bool, name string, diffName string, submoduleConfigs []*models.SubmoduleConfig, file *models.File) string {
	// potentially inefficient to be instantiating these color
	// objects with each render
	partiallyModifiedColor := style.FgYellow

	restColor := style.FgGreen
	if name == diffName {
		restColor = theme.DiffTerminalColor
	} else if file == nil && hasStagedChanges && hasUnstagedChanges {
		restColor = partiallyModifiedColor
	} else if hasUnstagedChanges {
		restColor = style.FgRed
	}

	output := ""
	if file != nil {
		// this is just making things look nice when the background attribute is 'reverse'
		firstChar := file.ShortStatus[0:1]
		firstCharCl := style.FgGreen
		if firstChar == "?" {
			firstCharCl = style.FgRed
		} else if firstChar == " " {
			firstCharCl = restColor
		}

		secondChar := file.ShortStatus[1:2]
		secondCharCl := style.FgRed
		if secondChar == " " {
			secondCharCl = restColor
		}

		output = firstCharCl.Sprint(firstChar)
		output += secondCharCl.Sprint(secondChar)
		output += restColor.Sprint(" ")
	}

	output += restColor.Sprint(utils.EscapeSpecialChars(name))

	if file != nil && file.IsSubmodule(submoduleConfigs) {
		output += theme.DefaultTextColor.Sprint(" (submodule)")
	}

	return output
}

func getCommitFileLine(name string, diffName string, commitFile *models.CommitFile, status patch.PatchStatus) string {
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
