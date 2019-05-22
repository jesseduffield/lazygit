package commands

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name      string
	Recency   string
	Pushables string
	Pullables string
	Selected  bool
}

// GetDisplayStrings returns the display string of branch
func (b *Branch) GetDisplayStrings(isFocused bool) []string {
	displayName := utils.ColoredString(b.Name, b.GetColor())
	if isFocused && b.Selected && b.Pushables != "" && b.Pullables != "" {
		displayName = fmt.Sprintf("%s ↑%s↓%s", displayName, b.Pushables, b.Pullables)
	}

	return []string{b.Recency, displayName}
}

// GetColor branch color
func (b *Branch) GetColor() color.Attribute {
	switch b.getType() {
	case "feature":
		return color.FgGreen
	case "bugfix":
		return color.FgYellow
	case "hotfix":
		return color.FgRed
	default:
		return color.FgWhite
	}
}

// expected to return feature/bugfix/hotfix or blank string
func (b *Branch) getType() string {
	return strings.Split(b.Name, "/")[0]
}
