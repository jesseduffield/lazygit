package git

import (
	"strings"

	"github.com/fatih/color"
)

// GetDisplayString returns the dispaly string of branch
// func (b *Branch) GetDisplayString() string {
// 	return gui.withPadding(b.Recency, 4) + gui.coloredString(b.Name, b.getColor())
// }

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
