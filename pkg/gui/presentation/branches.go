package presentation

import (
	"fmt"
	"strings"
	"time"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/logs"
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
	commitStore *models.CommitStore,
) [][]string {
	isContainedInMainBranch := isContainedInMainBranchFn(
		userConfig.Git.MainBranches,
		branches,
		commitStore,
	)

	if commitStore.Size() > 10 {
		logs.Global.Warnf("commitStore.Slice()[0:10]: %v", commitStore.Slice()[0:10])
	}

	return slices.Map(branches, func(branch *models.Branch) []string {
		diffed := branch.Name == diffName
		return getBranchDisplayStrings(branch, fullDescription, diffed, tr, userConfig, isContainedInMainBranch)
	})
}

// getBranchDisplayStrings returns the display string of branch
func getBranchDisplayStrings(
	b *models.Branch,
	fullDescription bool,
	diffed bool,
	tr *i18n.TranslationSet,
	userConfig *config.UserConfig,
	isContainedInMainBranch func(string) bool,
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
	coloredName = fmt.Sprintf("%s %s", coloredName, branchStatus)

	recencyColor := style.FgCyan
	if b.Recency == "  *" {
		recencyColor = style.FgGreen
	}

	res := make([]string, 0, 6)
	res = append(res, recencyColor.Sprint(b.Recency))
	if icons.IsIconEnabled() {
		res = append(res, nameTextStyle.Sprint(icons.IconForBranch(b)))
	}

	if fullDescription || userConfig.Gui.ShowBranchCommitHash {
		var hashStyle style.TextStyle
		if isContainedInMainBranch(b.CommitHash) {
			hashStyle = style.FgGreen
		} else {
			hashStyle = style.FgYellow
		}
		coloredHash := hashStyle.Sprint(utils.ShortSha(b.CommitHash))

		res = append(res, coloredHash)
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

// returns a function that tells us if a given branch's commit is contained in any of the main branches
func isContainedInMainBranchFn(mainBranches []string, branches []*models.Branch, commitStore *models.CommitStore) func(string) bool {
	mainBranchHashes := []string{}
	for _, branch := range branches {
		if lo.Contains(mainBranches, branch.Name) {
			mainBranchHashes = append(mainBranchHashes, branch.CommitHash)
		}
	}

	logs.Global.Warnf("mainBranchHashes: %v", mainBranchHashes)

	t := time.Now()

	ancestorSlices := lo.Map(mainBranchHashes, func(hash string, _ int) map[string]bool {
		return commitStore.FindAncestors(
			hash,
			lo.Map(branches, func(branch *models.Branch, _ int) string {
				return branch.CommitHash
			}),
		)
	})

	ancestors := lo.Reduce(ancestorSlices, mergeMaps, map[string]bool{})

	logs.Global.Warnf("isContainedInMainBranchFn took %v", time.Since(t))

	return func(commitHash string) bool {
		return ancestors[commitHash]
	}
}

func mergeMaps(a, b map[string]bool, _ int) map[string]bool {
	res := map[string]bool{}
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}

// it's faster to first check if master is a proper ancestor of my commit, because it likely is and it will be faster to find out. If I go the other way around, I need to traverse from master to the root for every branch potentially. Is there a way to speed that up?
