package presentation

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func makeAtomic(v int32) *atomic.Int32 {
	var result atomic.Int32
	result.Store(v)
	return &result
}

func TestFormatPullRequestHeader(t *testing.T) {
	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)
	icons.SetNerdFontsVersion("")

	pr := &models.GithubPullRequest{
		Title:  "Improve checks",
		Number: 5871,
		State:  "OPEN",
		Url:    "https://github.com/jesseduffield/lazygit/pull/5871",
	}
	numberText := style.FgCyan.Sprint("#5871")

	actual := FormatPullRequestHeader(pr)

	expected := style.PrintHyperlink("Open  Improve checks  "+numberText+"\n", pr.Url)
	assert.Equal(t, expected, actual)
}

func Test_getBranchDisplayStrings(t *testing.T) {
	scenarios := []struct {
		branch               *models.Branch
		itemOperation        types.ItemOperation
		fullDescription      bool
		viewWidth            int
		useIcons             bool
		checkedOutByWorktree bool
		showDivergenceCfg    string
		expected             []string
	}{
		// First some tests for when the view is wide enough so that everything fits:
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_name"},
		},
		{
			branch:               &models.Branch{Name: "🍉_special_char", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            19,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "🍉_special_char"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_name (worktree other-worktree)"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             true,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_name (󰌹 other-worktree)"},
		},
		{
			branch: &models.Branch{
				Name:           "branch_name",
				Recency:        "1m",
				UpstreamRemote: "origin",
				AheadForPull:   "0",
				BehindForPull:  "0",
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_name ✓"},
		},
		{
			branch: &models.Branch{
				Name:           "branch_name",
				Recency:        "1m",
				UpstreamRemote: "origin",
				AheadForPull:   "3",
				BehindForPull:  "5",
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_name (worktree other-worktree) ↓5↑3"},
		},
		{
			branch: &models.Branch{
				Name:             "branch_name",
				Recency:          "1m",
				BehindBaseBranch: *makeAtomic(2),
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            20,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "onlyArrow",
			expected:             []string{"1m", "", "branch_name    ↓"},
		},
		{
			branch: &models.Branch{
				Name:             "branch_name",
				Recency:          "1m",
				UpstreamRemote:   "origin",
				AheadForPull:     "0",
				BehindForPull:    "0",
				BehindBaseBranch: *makeAtomic(2),
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            22,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "arrowAndNumber",
			expected:             []string{"1m", "", "branch_name ✓   ↓2"},
		},
		{
			branch: &models.Branch{
				Name:             "branch_name",
				Recency:          "1m",
				UpstreamRemote:   "origin",
				AheadForPull:     "3",
				BehindForPull:    "5",
				BehindBaseBranch: *makeAtomic(2),
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            26,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "arrowAndNumber",
			expected:             []string{"1m", "", "branch_name ↓5↑3    ↓2"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_name Pushing ●∙∙"},
		},
		{
			branch: &models.Branch{
				Name:           "branch_name",
				Recency:        "1m",
				CommitHash:     "1234567890",
				UpstreamRemote: "origin",
				UpstreamBranch: "branch_name",
				AheadForPull:   "0",
				BehindForPull:  "0",
				Subject:        "commit title",
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      true,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "12345678", "branch_name ✓", "origin branch_name", "commit title"},
		},

		// Now tests for how we truncate the branch name when there's not enough room:
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            14,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_na…"},
		},
		{
			branch:               &models.Branch{Name: "🍉_special_char", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            18,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "🍉_special_ch…"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            14,
			useIcons:             false,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "bra… (worktree)"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            12,
			useIcons:             true,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branc… 󰌹"},
		},
		{
			branch: &models.Branch{
				Name:           "branch_name",
				Recency:        "1m",
				UpstreamRemote: "origin",
				AheadForPull:   "0",
				BehindForPull:  "0",
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            14,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_… ✓"},
		},
		{
			branch: &models.Branch{
				Name:             "branch_name",
				Recency:          "1m",
				UpstreamRemote:   "origin",
				AheadForPull:     "3",
				BehindForPull:    "5",
				BehindBaseBranch: *makeAtomic(4),
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            21,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "arrowAndNumber",
			expected:             []string{"1m", "", "branch_n… ↓5↑3 ↓4"},
		},
		{
			branch: &models.Branch{
				Name:           "branch_name",
				Recency:        "1m",
				UpstreamRemote: "origin",
				AheadForPull:   "3",
				BehindForPull:  "5",
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            30,
			useIcons:             false,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "branch_na… (worktree) ↓5↑3"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            20,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "bra… Pushing ●∙∙"},
		},
		{
			branch:               &models.Branch{Name: "abc", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "abc Pushing ●∙∙"},
		},
		{
			branch:               &models.Branch{Name: "ab", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "ab Pushing ●∙∙"},
		},
		{
			branch:               &models.Branch{Name: "a", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "a Pushing ●∙∙"},
		},
		{
			branch: &models.Branch{
				Name:           "branch_name",
				Recency:        "1m",
				CommitHash:     "1234567890",
				UpstreamRemote: "origin",
				UpstreamBranch: "branch_name",
				AheadForPull:   "0",
				BehindForPull:  "0",
				Subject:        "commit title",
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      true,
			viewWidth:            20,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "", "12345678", "bran… ✓", "origin branch_name", "commit title"},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	c := common.NewDummyCommon()
	SetCustomBranches(c.UserConfig().Gui.BranchColorPatterns, true)

	for i, s := range scenarios {
		icons.SetNerdFontsVersion(lo.Ternary(s.useIcons, "3", ""))
		c.UserConfig().Gui.ShowDivergenceFromBaseBranch = s.showDivergenceCfg

		worktrees := []*models.Worktree{}
		if s.checkedOutByWorktree {
			worktrees = append(worktrees, &models.Worktree{Branch: s.branch.Name, Name: "other-worktree"})
		}

		t.Run(fmt.Sprintf("getBranchDisplayStrings_%d", i), func(t *testing.T) {
			strings := getBranchDisplayStrings(s.branch, s.itemOperation, s.fullDescription, false, s.viewWidth, c.Tr, c.UserConfig(), worktrees, time.Time{}, map[string]*models.GithubPullRequest{})
			assert.Equal(t, s.expected, strings)
		})
	}
}
