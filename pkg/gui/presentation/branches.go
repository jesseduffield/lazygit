package presentation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetBranchListDisplayStrings(branches []*models.Branch, fullDescription bool, diffName string, showGithub bool) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		diffed := branches[i].Name == diffName
		lines[i] = getBranchDisplayStrings(branches[i], fullDescription, diffed, showGithub)
	}

	return lines
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(b *models.Branch, fullDescription bool, diffed, showGithub bool) []string {
	displayName := b.Name
	if b.DisplayName != "" {
		displayName = b.DisplayName
	}

	nameColorAttr := GetBranchColor(b.Name)
	if diffed {
		nameColorAttr = theme.DiffTerminalColor
	}
	coloredName := utils.ColoredString(displayName, nameColorAttr)
	if b.IsTrackingRemote() {
		coloredName = fmt.Sprintf("%s %s", coloredName, ColoredBranchStatus(b))
	}

	recencyColor := color.FgCyan
	if b.Recency == "  *" {
		recencyColor = color.FgGreen
	}

	res := []string{utils.ColoredString(b.Recency, recencyColor)}
	if showGithub {
		if b.PR != nil {
			colour := color.FgMagenta // = state MERGED
			switch b.PR.State {
			case "OPEN":
				colour = color.FgGreen
			case "CLOSED":
				colour = color.FgRed
			}
			res = append(res, utils.ColoredString("#"+strconv.Itoa(b.PR.Number), colour))
		} else {
			res = append(res, "")
		}
	}

	if fullDescription {
		return append(res, coloredName, utils.ColoredString(b.UpstreamName, color.FgYellow))
	}
	return append(res, coloredName)
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

func ColoredBranchStatus(branch *models.Branch) string {
	colour := color.FgYellow
	if branch.MatchesUpstream() {
		colour = color.FgGreen
	} else if !branch.IsTrackingRemote() {
		colour = color.FgRed
	}

	return utils.ColoredString(BranchStatus(branch), colour)
}

func BranchStatus(branch *models.Branch) string {
	return fmt.Sprintf("↑%s↓%s", branch.Pushables, branch.Pullables)
}
