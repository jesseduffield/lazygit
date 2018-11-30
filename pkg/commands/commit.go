package commands

import (
	"github.com/fatih/color"
)

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Pushed        bool
	Merged        bool
	DisplayString string
}

// GetDisplayStrings is a function.
func (c *Commit) GetDisplayStrings() []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgGreen)
	green := color.New(color.FgYellow)
	white := color.New(color.FgWhite)

	shaColor := yellow
	if c.Pushed {
		shaColor = red
	} else if !c.Merged {
		shaColor = green
	}

	return []string{shaColor.Sprint(c.Sha), white.Sprint(c.Name)}
}
