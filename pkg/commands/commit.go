package commands

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Status        string // one of "unpushed", "pushed", "merged", or "rebasing"
	DisplayString string
	Action        string // one of "", "pick", "edit", "squash", "reword", "drop", "fixup"
}

// GetDisplayStrings is a function.
func (c *Commit) GetDisplayStrings(isFocused bool) []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	white := color.New(color.FgWhite)
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)

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

	actionString := ""
	if c.Action != "" {
		actionString = cyan.Sprint(utils.WithPadding(c.Action, 7)) + " "
	}

	return []string{shaColor.Sprint(c.Sha), actionString + white.Sprint(c.Name)}
}
