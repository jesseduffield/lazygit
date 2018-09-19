package commands

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Branch : A git branch
// duplicating this for now
type Branch struct {
	Name    string
	Recency string
}

// GetDisplayStrings returns the dispaly string of branch
func (b *Branch) GetDisplayStrings() []string {
	return []string{b.Recency, utils.ColoredString(b.Name, b.GetColor())}
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
