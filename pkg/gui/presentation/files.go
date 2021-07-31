package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetFileLine(hasUnstagedChanges bool, hasStagedChanges bool, name string, diffName string, submoduleConfigs []*models.SubmoduleConfig, file *models.File) string {
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
