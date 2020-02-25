package presentation

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetBranchListDisplayStrings(branches []*commands.Branch, isFocused bool, selectedLine int) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		showUpstreamDifferences := isFocused && i == selectedLine
		lines[i] = getBranchDisplayStrings(branches[i], showUpstreamDifferences)
	}

	return lines
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(b *commands.Branch, showUpstreamDifferences bool) []string {
	displayName := utils.ColoredString(b.Name, GetBranchColor(b.Name))
	if showUpstreamDifferences && b.Pushables != "" && b.Pullables != "" {
		displayName = fmt.Sprintf("%s ↑%s↓%s", displayName, b.Pushables, b.Pullables)
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
