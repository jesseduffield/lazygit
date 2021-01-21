package presentation

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetBranchListDisplayStrings(branches []*models.Branch, fullDescription bool, diffName string) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		diffed := branches[i].Name == diffName
		lines[i] = getBranchDisplayStrings(branches[i], fullDescription, diffed)
	}

	return lines
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(b *models.Branch, fullDescription bool, diffed bool) []string {
	displayName := b.Name
	if b.DisplayName != "" {
		displayName = b.DisplayName
	}

	nameColorAttr := GetBranchColor(b.Name)
	if diffed {
		nameColorAttr = theme.DiffTerminalColor
	}
	coloredName := utils.ColoredString(displayName, nameColorAttr)
	if b.Pushables != "" && b.Pullables != "" && b.Pushables != "?" && b.Pullables != "?" {
		trackColor := color.FgYellow
		if b.Pushables == "0" && b.Pullables == "0" {
			trackColor = color.FgGreen
		}
		track := utils.ColoredString(fmt.Sprintf("↑%s↓%s", b.Pushables, b.Pullables), trackColor)
		coloredName = fmt.Sprintf("%s %s", coloredName, track)
	}

	if !b.Merged {
		coloredName = coloredName + utils.ColoredString("Δ", color.Bold, color.FgMagenta)
	}

	recencyColor := color.FgCyan
	if b.Recency == "  *" {
		recencyColor = color.FgGreen
	}

	if fullDescription {
		return []string{utils.ColoredString(b.Recency, recencyColor), coloredName, utils.ColoredString(b.UpstreamName, color.FgYellow)}
	}

	return []string{utils.ColoredString(b.Recency, recencyColor), coloredName}
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
