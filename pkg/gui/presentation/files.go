package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetFileListDisplayStrings(files []*commands.File, diffName string) [][]string {
	lines := make([][]string, len(files))

	for i := range files {
		diffed := files[i].Name == diffName
		lines[i] = getFileDisplayStrings(files[i], diffed)
	}

	return lines
}

// getFileDisplayStrings returns the display string of branch
func getFileDisplayStrings(f *commands.File, diffed bool) []string {
	// potentially inefficient to be instantiating these color
	// objects with each render
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	diffColor := color.New(theme.DiffTerminalColor)
	if !f.Tracked && !f.HasStagedChanges {
		return []string{red.Sprint(f.DisplayString)}
	}

	var restColor *color.Color
	if diffed {
		restColor = diffColor
	} else if f.HasUnstagedChanges {
		restColor = red
	} else {
		restColor = green
	}

	// this is just making things look nice when the background attribute is 'reverse'
	firstChar := f.DisplayString[0:1]
	firstCharCl := green
	if firstChar == " " {
		firstCharCl = restColor
	}

	secondChar := f.DisplayString[1:2]
	secondCharCl := red
	if secondChar == " " {
		secondCharCl = restColor
	}

	output := firstCharCl.Sprint(firstChar)
	output += secondCharCl.Sprint(secondChar)
	output += restColor.Sprintf(" %s", f.Name)

	if f.IsSubmodule {
		output += utils.ColoredString(" (submodule)", theme.DefaultTextColor)
	}

	return []string{output}
}
