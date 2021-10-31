package presentation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetBranchListDisplayStrings(
	branches []*models.Branch,
	prs map[*models.Branch]*models.GithubPullRequest,
	fullDescription bool,
	diffName string,
	showGithub bool) [][]string {
	lines := make([][]string, len(branches))

	for i := range branches {
		diffed := branches[i].Name == diffName
		lines[i] = getBranchDisplayStrings(branches[i], prs, fullDescription, diffed, showGithub)
	}

	return lines
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(
	b *models.Branch,
	prs map[*models.Branch]*models.GithubPullRequest,
	fullDescription bool,
	diffed,
	showGithub bool) []string {
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

	res := []string{recencyColor.Sprint(b.Recency), coloredName}
	if showGithub {
		pr, hasPr := prs[b]
		res = append(res, coloredPrNumber(pr, hasPr))
	}

	if fullDescription {
		return append(res, style.FgYellow.Sprint(b.UpstreamName))
	}
	return res
}

// GetBranchTextStyle branch color
func GetBranchTextStyle(name string) style.TextStyle {
	branchType := strings.Split(name, "/")[0]

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

func coloredPrNumber(pr *models.GithubPullRequest, hasPr bool) string {
	if hasPr {
		colour := style.FgMagenta // = state MERGED
		switch pr.State {
		case "OPEN":
			colour = style.FgGreen
		case "CLOSED":
			colour = style.FgRed
		}
		return colour.Sprint("#" + strconv.Itoa(pr.Number))
	}

	return ""
}
