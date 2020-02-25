package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func GetFileListDisplayStrings(files []*commands.File) [][]string {
	lines := make([][]string, len(files))

	for i := range files {
		lines[i] = getFileDisplayStrings(files[i])
	}

	return lines
}

// getFileDisplayStrings returns the display string of branch
func getFileDisplayStrings(f *commands.File) []string {
	// potentially inefficient to be instantiating these color
	// objects with each render
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	if !f.Tracked && !f.HasStagedChanges {
		return []string{red.Sprint(f.DisplayString)}
	}

	output := green.Sprint(f.DisplayString[0:1])
	output += red.Sprint(f.DisplayString[1:3])
	if f.HasUnstagedChanges {
		output += red.Sprint(f.Name)
	} else {
		output += green.Sprint(f.Name)
	}
	return []string{output}
}
