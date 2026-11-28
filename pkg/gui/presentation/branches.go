package presentation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type branchColorPattern struct {
	pattern string
	style   style.TextStyle
}

var branchColorPatterns []branchColorPattern

var dimYellow = style.FgYellow.SetDim()

func GetBranchListDisplayStrings(
	branches []*models.Branch,
	getItemOperation func(item types.HasUrn) types.ItemOperation,
	prs map[string]*models.GithubPullRequest,
	fullDescription bool,
	diffName string,
	viewWidth int,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
	worktrees []*models.Worktree,
) [][]string {
	return lo.Map(branches, func(branch *models.Branch, _ int) []string {
		diffed := branch.Name == diffName
		return getBranchDisplayStrings(branch, getItemOperation(branch), fullDescription, diffed, viewWidth, tr, userConfig, worktrees, time.Now(), prs)
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
	prs map[string]*models.GithubPullRequest,
) []string {
	parts := getBranchRowParts(b, itemOperation, fullDescription, tr, userConfig, worktrees, now, prs)
	displayName := parts.displayName

	// If the worktree name doesn't fit in the available width, omit it
	worktreeIcon := parts.worktreeIcon
	if viewWidth-parts.widthExceptName(worktreeIcon) < utils.StringWidth(displayName) {
		worktreeIcon = parts.shortWorktreeIcon
	}
	availableWidth := viewWidth - parts.widthExceptName(worktreeIcon)

	nameTextStyle := GetBranchTextStyle(b.Name)
	if diffed {
		nameTextStyle = theme.DiffTerminalColor
	}

	// Don't bother shortening branch names that are already 3 characters or less
	if utils.StringWidth(displayName) > max(availableWidth, 3) {
		// Never shorten the branch name to less then 3 characters
		len := max(availableWidth, 4)
		displayName = utils.TruncateWithEllipsis(displayName, len)
	}
	coloredName := nameTextStyle.Sprint(displayName)
	if worktreeIcon != "" {
		coloredName = fmt.Sprintf("%s %s", coloredName, style.FgDefault.Sprint(worktreeIcon))
	}
	if len(parts.branchStatus) > 0 {
		coloredName = fmt.Sprintf("%s %s", coloredName, parts.branchStatus)
	}

	recencyColor := style.FgCyan
	if b.Recency == "  *" {
		recencyColor = style.FgGreen
	}

	res := make([]string, 0, 6)
	res = append(res, recencyColor.Sprint(b.Recency))

	var coloredPrIcon string
	pr, hasPr := prs[b.Name]
	if hasPr && ShouldShowPrForBranch(pr, b.Name, userConfig) {
		var prIcon string
		if icons.IsIconEnabled() {
			prIcon = icons.IconForRemoteUrl(pr.Url)
		} else {
			prIcon = "●"
		}
		coloredPrIcon = WithPrColor(pr.State, prIcon, false)
		if pr.State == "OPEN" {
			icon, _, textStyle := checksStatePresentation(pr.ChecksState, tr)
			if icon != "" {
				coloredPrIcon = textStyle.Sprint(icon)
			}
		}
	}
	res = append(res, coloredPrIcon)

	if parts.showCommitHash {
		// The CommitHash of a detached HEAD is empty, but its Name is the hash
		res = append(res, utils.ShortHash(lo.Ternary(b.DetachedHead, b.Name, b.CommitHash)))
	}

	if parts.divergence != "" {
		var paddingNeededForDivergence int
		if fullDescription {
			// don't right-align the divergence in half or full screen mode, since other fields
			// follow in that case and we don't know the width of our column
			paddingNeededForDivergence = 1
		} else {
			// The room that the name leaves, plus the space before the divergence
			paddingNeededForDivergence = availableWidth - utils.StringWidth(displayName) + 1
		}
		if paddingNeededForDivergence > 0 {
			coloredName += strings.Repeat(" ", paddingNeededForDivergence)
			coloredName += style.FgCyan.Sprint(parts.divergence)
		}
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

// branchRowParts holds the name to show in a branch's row, and everything that shares
// the row with it.
type branchRowParts struct {
	displayName    string
	showCommitHash bool
	// Whether the rows have a column for pull request icons
	hasPrColumn  bool
	branchStatus string
	divergence   string
	// The icon for a branch that another worktree has checked out, with the
	// worktree's name if it differs from the branch's, and the icon without the
	// name for a row that is too narrow for it. Both are empty if no other worktree
	// has the branch checked out.
	worktreeIcon      string
	shortWorktreeIcon string
}

func getBranchRowParts(
	b *models.Branch,
	itemOperation types.ItemOperation,
	fullDescription bool,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
	worktrees []*models.Worktree,
	now time.Time,
	prs map[string]*models.GithubPullRequest,
) branchRowParts {
	parts := branchRowParts{
		displayName:    b.Name,
		showCommitHash: fullDescription || userConfig.Gui.ShowBranchCommitHash,
		// if we have PRs then we assume that at least one branch in the list has one
		hasPrColumn:  len(prs) > 0,
		branchStatus: BranchStatus(b, itemOperation, tr, now, userConfig),
		divergence:   divergenceStr(b, itemOperation, tr, userConfig),
	}
	if b.DisplayName != "" {
		parts.displayName = b.DisplayName
	}

	if git_commands.CheckedOutByOtherWorktree(b, worktrees) {
		if icons.IsIconEnabled() {
			parts.shortWorktreeIcon = icons.LINKED_WORKTREE_ICON
		} else {
			parts.shortWorktreeIcon = fmt.Sprintf("(%s)", tr.LcWorktree)
		}
		parts.worktreeIcon = parts.shortWorktreeIcon

		if wt, ok := git_commands.WorktreeForBranch(b, worktrees); ok && wt.Name != b.Name {
			if icons.IsIconEnabled() {
				parts.worktreeIcon = fmt.Sprintf("(%s %s)", icons.LINKED_WORKTREE_ICON, wt.Name)
			} else {
				parts.worktreeIcon = fmt.Sprintf("(%s %s)", tr.LcWorktree, wt.Name)
			}
		}
	}

	return parts
}

// widthExceptName returns the width of everything in the row except the name, with
// the given worktree icon, including the spaces that separate the parts.
func (self branchRowParts) widthExceptName(worktreeIcon string) int {
	// Recency is always three characters, plus one for the space
	width := 4
	if len(self.divergence) > 0 {
		width += utils.StringWidth(self.divergence) + 1
	}
	if self.showCommitHash {
		width += utils.COMMIT_HASH_SHORT_SIZE + 1
	}
	if self.hasPrColumn {
		width += 2
	}
	if len(self.branchStatus) > 0 {
		width += utils.StringWidth(utils.Decolorise(self.branchStatus)) + 1
	}
	if worktreeIcon != "" {
		width += utils.StringWidth(worktreeIcon) + 1
	}
	return width
}

// GetBranchTextStyle branch color
func GetBranchTextStyle(name string) style.TextStyle {
	for _, p := range branchColorPatterns {
		if matched, _ := regexp.MatchString(p.pattern, name); matched {
			return p.style
		}
	}

	return theme.DefaultTextColor
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
		return style.FgCyan.Sprintf("%s %s", itemOperationStr, Loader(now, userConfig.Gui.Spinner))
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
			// A branch that diverged only because its upstream was rewritten
			// has no commits of its own, and fast-forwarding it resolves the
			// divergence. Dim it to set it apart from a branch whose
			// divergence needs a decision.
			divergenceStyle := lo.Ternary(branch.UpstreamRewritten.Load(), dimYellow, style.FgYellow)
			result = divergenceStyle.Sprintf("↓%s↑%s", branch.BehindForPull, branch.AheadForPull)
		} else if branch.IsBehindForPull() {
			result = style.FgYellow.Sprintf("↓%s", branch.BehindForPull)
		} else if branch.IsAheadForPull() {
			result = style.FgYellow.Sprintf("↑%s", branch.AheadForPull)
		}
	}

	return result
}

func divergenceStr(
	branch *models.Branch,
	itemOperation types.ItemOperation,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
) string {
	result := ""
	if ItemOperationToString(itemOperation, tr) == "" && userConfig.Gui.ShowDivergenceFromBaseBranch != "none" {
		behind := branch.BehindBaseBranch.Load()
		if behind != 0 {
			if userConfig.Gui.ShowDivergenceFromBaseBranch == "arrowAndNumber" {
				result += fmt.Sprintf("↓%d", behind)
			} else {
				result += "↓"
			}
		}
	}

	return result
}

func SetCustomBranches(patterns config.ColorPatterns) {
	branchColorPatterns = lo.Map(patterns, func(p config.ColorPattern, _ int) branchColorPattern {
		return branchColorPattern{pattern: p.Pattern, style: utils.CustomColorStyle(p.Color)}
	})
}

func WithPrColor(state string, text string, isBg bool) string {
	switch state {
	case "OPEN":
		return color.RGB(0x43, 0x84, 0x40, isBg).Sprint(text)
	case "CLOSED":
		return color.RGB(0xC9, 0x45, 0x3C, isBg).Sprint(text)
	case "MERGED":
		return color.RGB(0x82, 0x59, 0xDD, isBg).Sprint(text)
	case "DRAFT":
		return color.RGB(0x67, 0x6C, 0x75, isBg).Sprint(text)
	default:
		return lo.Ternary(isBg, style.BgDefault, style.FgDefault).Sprint(text)
	}
}

func FormatPullRequestHeader(pr *models.GithubPullRequest, tr *i18n.TranslationSet) string {
	icon := lo.Ternary(icons.IsIconEnabled(), icons.IconForRemoteUrl(pr.Url)+"  ", "")
	stateText := coloredPullRequestStateText(pr.State)
	checksStateText := coloredChecksStateText(pr.ChecksState, tr)
	numberText := style.FgCyan.Sprintf("#%d", pr.Number)

	// The checks status links to the checks tab, so it needs to be its own
	// hyperlink separate from the rest of the header.
	parts := []string{style.PrintHyperlink(icon+stateText, pr.Url)}
	if checksStateText != "" {
		parts = append(parts, style.PrintHyperlink(checksStateText, strings.TrimSuffix(pr.Url, "/")+"/checks"))
	}
	parts = append(parts, style.PrintHyperlink(fmt.Sprintf("%s  %s\n", pr.Title, numberText), pr.Url))

	return strings.Join(parts, "  ")
}

func pullRequestStateText(state string) string {
	var icon, label string
	switch state {
	case "OPEN":
		icon, label = " ", "Open"
	case "CLOSED":
		icon, label = " ", "Closed"
	case "MERGED":
		icon, label = " ", "Merged"
	case "DRAFT":
		icon, label = " ", "Draft"
	default:
		return ""
	}
	if icons.IsIconEnabled() {
		return icon + label
	}
	return label
}

func coloredPullRequestStateText(state string) string {
	if icons.IsIconEnabled() {
		return fmt.Sprintf("%s%s%s",
			WithPrColor(state, "", false),
			WithPrColor(state, color.RGB(0xFF, 0xFF, 0xFF, false).Sprint(pullRequestStateText(state)), true),
			WithPrColor(state, "", false))
	}

	return WithPrColor(state, pullRequestStateText(state), false)
}

func checksStatePresentation(state string, tr *i18n.TranslationSet) (string, string, style.TextStyle) {
	switch state {
	case "SUCCESS":
		return "✓", tr.PullRequestChecksPassing, style.FgGreen
	case "PENDING":
		return "●", tr.PullRequestChecksPending, style.FgYellow
	case "FAILURE":
		return "✗", tr.PullRequestChecksFailing, style.FgRed
	case "ERROR":
		return "!", tr.PullRequestChecksError, style.FgRed
	case "EXPECTED":
		return "○", tr.PullRequestChecksExpected, style.FgDefault
	default:
		return "", "", style.Nothing
	}
}

func coloredChecksStateText(state string, tr *i18n.TranslationSet) string {
	icon, text, textStyle := checksStatePresentation(state, tr)
	if text != "" {
		return textStyle.Sprintf("%s %s", icon, text)
	}
	return ""
}

func ShouldShowPrForBranch(pr *models.GithubPullRequest, branchName string, userConfig *config.UserConfig) bool {
	if !lo.Contains(userConfig.Git.MainBranches, branchName) {
		return true
	}

	// For main branches we only want to show the PR if it's open (or draft), on the assumption that a
	// closed PR for a main branch is always a mistake.
	return pr.State != "CLOSED" && pr.State != "MERGED"
}
