package presentation

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

var branchPrefixColorCache = make(map[string]style.TextStyle)

func GetBranchListDisplayStrings(branches []*models.Branch, fullDescription bool, diffName string) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		diffed := branches[i].Name == diffName
		lines[i] = getBranchDisplayStrings(branches[i], prs, fullDescription, diffed)
	}

	return lines
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(
	b *models.Branch,
	prs map[*models.Branch]*models.GithubPullRequest,
	fullDescription bool,
	diffed bool) []string {
	displayName := b.Name
	if b.DisplayName != "" {
		displayName = b.DisplayName
	}

	nameTextStyle := GetBranchTextStyle(b.Name)
	if diffed {
		nameTextStyle = theme.DiffTerminalColor
	}
	coloredName := nameTextStyle.Sprint(displayName)
	if b.IsTrackingRemote() {
		coloredName = fmt.Sprintf("%s %s", coloredName, ColoredBranchStatus(b))
	}

	recencyColor := style.FgCyan
	if b.Recency == "  *" {
		recencyColor = style.FgGreen
	}

	res := []string{recencyColor.Sprint(b.Recency)}
	pr, hasPr := prs[b]
	res = append(res, coloredPrNumber(pr, hasPr), coloredName)

	if fullDescription {
		return append(
			res,
			fmt.Sprintf("%s %s",
				style.FgYellow.Sprint(b.UpstreamRemote),
				style.FgYellow.Sprint(b.UpstreamBranch),
			),
		)
	}
	return res
}

// GetBranchTextStyle branch color
func GetBranchTextStyle(name string) style.TextStyle {
	branchType := strings.Split(name, "/")[0]

	if value, ok := branchPrefixColorCache[branchType]; ok {
		return value
	}

	switch branchType {
	case "feature":
		return style.FgGreen
	case "bugfix":
		return style.FgYellow
	case "hotfix":
		return style.FgRed
	default:
		return theme.DefaultTextColor
	}
}

func ColoredBranchStatus(branch *models.Branch) string {
	colour := style.FgYellow
	if branch.MatchesUpstream() {
		colour = style.FgGreen
	} else if !branch.IsTrackingRemote() {
		colour = style.FgRed
	}

	return colour.Sprint(BranchStatus(branch))
}

func BranchStatus(branch *models.Branch) string {
	return fmt.Sprintf("↑%s↓%s", branch.Pushables, branch.Pullables)
}

func SetCustomBranches(customBranchColors map[string]string) {
	branchPrefixColorCache = utils.SetCustomColors(customBranchColors)
}
