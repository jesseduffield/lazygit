package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/theme"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name           string
	Recency        string
	Pushables      string
	Pullables      string
	Selected       bool
	encodedStrings *utils.EncodedStrings
}

// GetDisplayStrings returns the display string of branch
func (b *Branch) GetDisplayStrings(isFocused bool) []string {
	displayName := utils.ColoredString(b.Name, GetBranchColor(b.Name))
	if isFocused && b.Selected && b.Pushables != "" && b.Pullables != "" {
		displayName = fmt.Sprintf("%s %s%s%s%s", displayName, b.encodedStrings.UpArrow, b.Pushables, b.encodedStrings.DownArrow, b.Pullables)
	}

	return []string{b.Recency, displayName}
}

// GetBranchColor branch color
func GetBranchColor(name string) color.Attribute {
	branchType := strings.Split(name, "/")[0]

	switch branchType {
	case "feature":
		return color.FgGreen
	case "bugfix":
		return color.FgYellow
	case "hotfix":
		return color.FgRed
	default:
		return theme.DefaultTextColor
	}
}
