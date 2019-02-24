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
	Copied        bool   // to know if this commit is ready to be cherry-picked somewhere
}

// GetDisplayStrings is a function.
func (c *Commit) GetDisplayStrings(isFocused bool) []string {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)
	white := color.New(color.FgWhite)

	// for some reason, setting the background to blue pads out the other commits
	// horizontally. For the sake of accessibility I'm considering this a feature,
	// not a bug
	copied := color.New(color.FgCyan, color.BgBlue)

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

	if c.Copied {
		shaColor = copied
	}

	actionString := ""
	if c.Action != "" {
		actionString = cyan.Sprint(utils.WithPadding(c.Action, 7)) + " "
	}

	return []string{shaColor.Sprint(c.Sha), actionString + white.Sprint(c.Name)}
}
