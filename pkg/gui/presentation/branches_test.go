package presentation

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func makeAtomic(v int32) (result atomic.Int32) {
	result.Store(v)
	return //nolint: nakedret
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
			expected:             []string{"1m", "branch_name"},
		},
		{
			branch:               &models.Branch{Name: "🍉_special_char", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            19,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "🍉_special_char"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "branch_name (worktree)"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             true,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "󰘬", "branch_name 󰌹"},
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
			expected:             []string{"1m", "branch_name ✓"},
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
			expected:             []string{"1m", "branch_name (worktree) ↓5↑3"},
		},
		{
			branch: &models.Branch{
				Name:             "branch_name",
				Recency:          "1m",
				BehindBaseBranch: makeAtomic(2),
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "onlyArrow",
			expected:             []string{"1m", "branch_name ↓"},
		},
		{
			branch: &models.Branch{
				Name:             "branch_name",
				Recency:          "1m",
				UpstreamRemote:   "origin",
				AheadForPull:     "0",
				BehindForPull:    "0",
				BehindBaseBranch: makeAtomic(2),
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "arrowAndNumber",
			expected:             []string{"1m", "branch_name ✓ ↓2"},
		},
		{
			branch: &models.Branch{
				Name:             "branch_name",
				Recency:          "1m",
				UpstreamRemote:   "origin",
				AheadForPull:     "3",
				BehindForPull:    "5",
				BehindBaseBranch: makeAtomic(2),
			},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "arrowAndNumber",
			expected:             []string{"1m", "branch_name ↓5↑3 ↓2"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "branch_name Pushing |"},
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
			expected:             []string{"1m", "12345678", "branch_name ✓", "origin branch_name", "commit title"},
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
			expected:             []string{"1m", "branch_na…"},
		},
		{
			branch:               &models.Branch{Name: "🍉_special_char", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            18,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "🍉_special_ch…"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            14,
			useIcons:             false,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "bra… (worktree)"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            14,
			useIcons:             true,
			checkedOutByWorktree: true,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "󰘬", "branc… 󰌹"},
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
			expected:             []string{"1m", "branch_… ✓"},
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
			expected:             []string{"1m", "branch_na… (worktree) ↓5↑3"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            20,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "branc… Pushing |"},
		},
		{
			branch:               &models.Branch{Name: "abc", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "abc Pushing |"},
		},
		{
			branch:               &models.Branch{Name: "ab", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "ab Pushing |"},
		},
		{
			branch:               &models.Branch{Name: "a", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			showDivergenceCfg:    "none",
			expected:             []string{"1m", "a Pushing |"},
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
			expected:             []string{"1m", "12345678", "bran… ✓", "origin branch_name", "commit title"},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	c := utils.NewDummyCommon()

	for i, s := range scenarios {
		icons.SetNerdFontsVersion(lo.Ternary(s.useIcons, "3", ""))
		c.UserConfig.Gui.ShowDivergenceFromBaseBranch = s.showDivergenceCfg

		worktrees := []*models.Worktree{}
		if s.checkedOutByWorktree {
			worktrees = append(worktrees, &models.Worktree{Branch: s.branch.Name})
		}

		t.Run(fmt.Sprintf("getBranchDisplayStrings_%d", i), func(t *testing.T) {
			strings := getBranchDisplayStrings(s.branch, s.itemOperation, s.fullDescription, false, s.viewWidth, c.Tr, c.UserConfig, worktrees, time.Time{})
			assert.Equal(t, s.expected, strings)
		})
	}
}
