package commands

import "github.com/fatih/color"

// File : A file from git status
// duplicating this for now
type File struct {
	Name                    string
	HasStagedChanges        bool
	HasUnstagedChanges      bool
	Tracked                 bool
	Deleted                 bool
	HasMergeConflicts       bool
	HasInlineMergeConflicts bool
	DisplayString           string
	Type                    string // one of 'file', 'directory', and 'other'
}

// GetDisplayStrings returns the display string of a file
func (f *File) GetDisplayStrings(isFocused bool) []string {
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
