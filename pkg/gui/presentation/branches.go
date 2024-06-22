package presentation

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
)

var branchPrefixColorCache = make(map[string]style.TextStyle)

func GetBranchListDisplayStrings(
	branches []*models.Branch,
	getItemOperation func(item types.HasUrn) types.ItemOperation,
	fullDescription bool,
	diffName string,
	viewWidth int,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
	worktrees []*models.Worktree,
) [][]string {
	return lo.Map(branches, func(branch *models.Branch, _ int) []string {
		diffed := branch.Name == diffName
		return getBranchDisplayStrings(branch, getItemOperation(branch), fullDescription, diffed, viewWidth, tr, userConfig, worktrees, time.Now())
	})
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(
	b *models.Branch,
	itemOperation types.ItemOperation,
	fullDescription bool,
	diffed bool,
	viewWidth int,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
	worktrees []*models.Worktree,
	now time.Time,
) []string {
	checkedOutByWorkTree := git_commands.CheckedOutByOtherWorktree(b, worktrees)
	showCommitHash := fullDescription || userConfig.Gui.ShowBranchCommitHash
	branchStatus := BranchStatus(b, itemOperation, tr, now, userConfig)
	worktreeIcon := lo.Ternary(icons.IsIconEnabled(), icons.LINKED_WORKTREE_ICON, fmt.Sprintf("(%s)", tr.LcWorktree))

	// Recency is always three characters, plus one for the space
	availableWidth := viewWidth - 4
	if len(branchStatus) > 0 {
		availableWidth -= utils.StringWidth(utils.Decolorise(branchStatus)) + 1
	}
	if icons.IsIconEnabled() {
		availableWidth -= 2 // one for the icon, one for the space
	}
	if showCommitHash {
		availableWidth -= utils.COMMIT_HASH_SHORT_SIZE + 1
	}
	if checkedOutByWorkTree {
		availableWidth -= utils.StringWidth(worktreeIcon) + 1
	}

	displayName := b.Name
	if b.DisplayName != "" {
		displayName = b.DisplayName
	}

	nameTextStyle := GetBranchTextStyle(b.Name)
	if diffed {
		nameTextStyle = theme.DiffTerminalColor
	}

	// Don't bother shortening branch names that are already 3 characters or less
	if utils.StringWidth(displayName) > max(availableWidth, 3) {
		// Never shorten the branch name to less then 3 characters
		len := max(availableWidth, 4)
		displayName = runewidth.Truncate(displayName, len, "…")
	}
	coloredName := nameTextStyle.Sprint(displayName)
	if checkedOutByWorkTree {
		coloredName = fmt.Sprintf("%s %s", coloredName, style.FgDefault.Sprint(worktreeIcon))
	}
	if len(branchStatus) > 0 {
		coloredName = fmt.Sprintf("%s %s", coloredName, branchStatus)
	}

	recencyColor := style.FgCyan
	if b.Recency == "  *" {
		recencyColor = style.FgGreen
	}

	res := make([]string, 0, 6)
	res = append(res, recencyColor.Sprint(b.Recency))

	if icons.IsIconEnabled() {
		res = append(res, nameTextStyle.Sprint(icons.IconForBranch(b)))
	}

	if showCommitHash {
		res = append(res, utils.ShortHash(b.CommitHash))
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

func BranchStatus(
	branch *models.Branch,
	itemOperation types.ItemOperation,
	tr *i18n.TranslationSet,
	now time.Time,
	userConfig *config.UserConfig,
) string {
	itemOperationStr := ItemOperationToString(itemOperation, tr)
	if itemOperationStr != "" {
		return style.FgCyan.Sprintf("%s %s", itemOperationStr, utils.Loader(now, userConfig.Gui.Spinner))
	}

	result := ""
	if branch.IsTrackingRemote() {
		if branch.UpstreamGone {
			result = style.FgRed.Sprint(tr.UpstreamGone)
		} else if branch.MatchesUpstream() {
			result = style.FgGreen.Sprint("✓")
		} else if branch.RemoteBranchNotStoredLocally() {
			result = style.FgMagenta.Sprint("?")
		} else if branch.IsBehindForPull() && branch.IsAheadForPull() {
			result = style.FgYellow.Sprintf("↓%s↑%s", branch.BehindForPull, branch.AheadForPull)
		} else if branch.IsBehindForPull() {
			result = style.FgYellow.Sprintf("↓%s", branch.BehindForPull)
		} else if branch.IsAheadForPull() {
			result = style.FgYellow.Sprintf("↑%s", branch.AheadForPull)
		}
	}

	if userConfig.Gui.ShowDivergenceFromBaseBranch != "none" {
		behind := branch.BehindBaseBranch.Load()
		if behind != 0 {
			if result != "" {
				result += " "
			}
			if userConfig.Gui.ShowDivergenceFromBaseBranch == "arrowAndNumber" {
				result += style.FgCyan.Sprintf("↓%d", behind)
			} else {
				result += style.FgCyan.Sprintf("↓")
			}
		}
	}

	return result
}

func SetCustomBranches(customBranchColors map[string]string) {
	branchPrefixColorCache = utils.SetCustomColors(customBranchColors)
}
