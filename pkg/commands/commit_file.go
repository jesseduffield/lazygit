package commands

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

// CommitFile : A git commit file
type CommitFile struct {
	Sha           string
	Name          string
	DisplayString string
	Status        int // one of 'WHOLE' 'PART' 'NONE'
}

const (
	// UNSELECTED is for when the commit file has not been added to the patch in any way
	UNSELECTED = iota
	// WHOLE is for when you want to add the whole diff of a file to the patch,
	// including e.g. if it was deleted
	WHOLE = iota
	// PART is for when you're only talking about specific lines that have been modified
	PART
)

// GetDisplayStrings is a function.
func (f *CommitFile) GetDisplayStrings(isFocused bool) []string {
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	defaultColor := color.New(theme.DefaultTextColor)

	var colour *color.Color
	switch f.Status {
	case UNSELECTED:
		colour = defaultColor
	case WHOLE:
		colour = green
	case PART:
		colour = yellow
	}
	return []string{colour.Sprint(f.DisplayString)}
}
