package presentation

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

var branchPrefixColorCache = make(map[string]style.TextStyle)

func GetBranchListDisplayStrings(
	branches []*models.Branch,
	fullDescription bool,
	diffName string,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
	worktrees []*models.Worktree,
) [][]string {
	return lo.Map(branches, func(branch *models.Branch, _ int) []string {
		diffed := branch.Name == diffName
		return getBranchDisplayStrings(branch, fullDescription, diffed, tr, userConfig, worktrees)
	})
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(
	b *models.Branch,
	fullDescription bool,
	diffed bool,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
	worktrees []*models.Worktree,
) []string {
	displayName := b.Name
	if b.DisplayName != "" {
		displayName = b.DisplayName
	}

	nameTextStyle := GetBranchTextStyle(b.Name)
	if diffed {
		nameTextStyle = theme.DiffTerminalColor
	}

	coloredName := nameTextStyle.Sprint(displayName)
	branchStatus := utils.WithPadding(ColoredBranchStatus(b, tr), 2, utils.AlignLeft)
	if git_commands.CheckedOutByOtherWorktree(b, worktrees) {
		worktreeIcon := lo.Ternary(icons.IsIconEnabled(), icons.LINKED_WORKTREE_ICON.Icon, fmt.Sprintf("(%s)", tr.LcWorktree))
    
		coloredName = fmt.Sprintf("%s %s", coloredName, style.FgDefault.Sprint(worktreeIcon))
	}
	coloredName = fmt.Sprintf("%s %s", coloredName, branchStatus)

	recencyColor := style.FgCyan
	if b.Recency == "  *" {
		recencyColor = style.FgGreen
	}

	res := make([]string, 0, 6)
	res = append(res, recencyColor.Sprint(b.Recency))

	if icons.IsIconEnabled() {
		res = append(res, nameTextStyle.Sprint(icons.IconForBranch(b).Icon))
	}

	if fullDescription || userConfig.Gui.ShowBranchCommitHash {
		res = append(res, utils.ShortSha(b.CommitHash))
	}

	res = append(res, coloredName)

	if fullDescription {
		res = append(
			res,
			fmt.Sprintf("%s %s",
				style.FgYellow.Sprint(b.UpstreamRemote),
				style.FgYellow.Sprint(b.UpstreamBranch),
			),
			utils.TruncateWithEllipsis(b.Subject, 60),
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

func ColoredBranchStatus(branch *models.Branch, tr *i18n.TranslationSet) string {
	colour := style.FgYellow
	if branch.UpstreamGone {
		colour = style.FgRed
	} else if branch.MatchesUpstream() {
		colour = style.FgGreen
	} else if branch.RemoteBranchNotStoredLocally() {
		colour = style.FgMagenta
	}

	return colour.Sprint(BranchStatus(branch, tr))
}

func BranchStatus(branch *models.Branch, tr *i18n.TranslationSet) string {
	if !branch.IsTrackingRemote() {
		return ""
	}

	if branch.UpstreamGone {
		return tr.UpstreamGone
	}

	if branch.MatchesUpstream() {
		return "✓"
	}
	if branch.RemoteBranchNotStoredLocally() {
		return "?"
	}

	result := ""
	if branch.HasCommitsToPush() {
		result = fmt.Sprintf("↑%s", branch.Pushables)
	}
	if branch.HasCommitsToPull() {
		result = fmt.Sprintf("%s↓%s", result, branch.Pullables)
	}

	return result
}

func SetCustomBranches(customBranchColors map[string]string) {
	branchPrefixColorCache = utils.SetCustomColors(customBranchColors)
}
