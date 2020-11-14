package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetFileListDisplayStrings(files []*models.File, diffName string, submoduleConfigs []*models.SubmoduleConfig) [][]string {
	lines := make([][]string, len(files))

	for i := range files {
		lines[i] = getFileDisplayStrings(files[i], diffName, submoduleConfigs)
	}

	return lines
}

// getFileDisplayStrings returns the display string of branch
func getFileDisplayStrings(f *models.File, diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	output := GetStatusNodeLine(f.HasUnstagedChanges, f.ShortStatus, f.Name, diffName, submoduleConfigs, f)

	return []string{output}
}

func GetStatusNodeLine(hasUnstagedChanges bool, shortStatus string, name string, diffName string, submoduleConfigs []*models.SubmoduleConfig, file *models.File) string {
	// potentially inefficient to be instantiating these color
	// objects with each render
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	diffColor := color.New(theme.DiffTerminalColor)

	var restColor *color.Color
	if name == diffName {
		restColor = diffColor
	} else if hasUnstagedChanges {
		restColor = red
	} else {
		restColor = green
	}

	// this is just making things look nice when the background attribute is 'reverse'
	firstChar := shortStatus[0:1]
	firstCharCl := green
	if firstChar == "?" {
		firstCharCl = red
	} else if firstChar == " " {
		firstCharCl = restColor
	}

	secondChar := shortStatus[1:2]
	secondCharCl := red
	if secondChar == " " {
		secondCharCl = restColor
	}

	output := firstCharCl.Sprint(firstChar)
	output += secondCharCl.Sprint(secondChar)
	output += restColor.Sprintf(" %s", name)

	if file != nil && file.IsSubmodule(submoduleConfigs) {
		output += utils.ColoredString(" (submodule)", theme.DefaultTextColor)
	}

	return output
}
