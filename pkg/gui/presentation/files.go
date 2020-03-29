package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
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

	output := green.Sprint(f.DisplayString[0:1])
	output += red.Sprint(f.DisplayString[1:3])

	var restColor *color.Color
	if diffed {
		restColor = diffColor
	} else if f.HasUnstagedChanges {
		restColor = red
	} else {
		restColor = green
	}
	output += restColor.Sprint(f.Name)
	return []string{output}
}
