package commands

import (
	"github.com/fatih/color"
)

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Status        string // one of "unpushed", "pushed", "merged", or "rebasing"
	DisplayString string
}

// GetDisplayStrings is a function.
func (c *Commit) GetDisplayStrings(isFocused bool) []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	white := color.New(color.FgWhite)
	blue := color.New(color.FgBlue)

	var shaColor *color.Color
	switch c.Status {
	case "unpushed":
		shaColor = red
	case "pushed":
		shaColor = yellow
	case "merged":
		shaColor = green
	case "rebasing":
		shaColor = blue
	default:
		shaColor = white
	}

	return []string{shaColor.Sprint(c.Sha), white.Sprint(c.Name)}
}
