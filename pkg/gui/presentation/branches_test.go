package presentation

import (
	"fmt"
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

func Test_getBranchDisplayStrings(t *testing.T) {
	scenarios := []struct {
		branch               *models.Branch
		itemOperation        types.ItemOperation
		fullDescription      bool
		viewWidth            int
		useIcons             bool
		checkedOutByWorktree bool
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
			expected:             []string{"1m", "branch_name"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: true,
			expected:             []string{"1m", "branch_name (worktree)"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             true,
			checkedOutByWorktree: true,
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
			expected:             []string{"1m", "branch_name (worktree) ↑3↓5"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            100,
			useIcons:             false,
			checkedOutByWorktree: false,
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
			expected:             []string{"1m", "branch_na…"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            14,
			useIcons:             false,
			checkedOutByWorktree: true,
			expected:             []string{"1m", "bra… (worktree)"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationNone,
			fullDescription:      false,
			viewWidth:            14,
			useIcons:             true,
			checkedOutByWorktree: true,
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
			expected:             []string{"1m", "branch_na… (worktree) ↑3↓5"},
		},
		{
			branch:               &models.Branch{Name: "branch_name", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            20,
			useIcons:             false,
			checkedOutByWorktree: false,
			expected:             []string{"1m", "branc… Pushing |"},
		},
		{
			branch:               &models.Branch{Name: "abc", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			expected:             []string{"1m", "abc Pushing |"},
		},
		{
			branch:               &models.Branch{Name: "ab", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
			expected:             []string{"1m", "ab Pushing |"},
		},
		{
			branch:               &models.Branch{Name: "a", Recency: "1m"},
			itemOperation:        types.ItemOperationPushing,
			fullDescription:      false,
			viewWidth:            -1,
			useIcons:             false,
			checkedOutByWorktree: false,
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
			expected:             []string{"1m", "12345678", "bran… ✓", "origin branch_name", "commit title"},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	c := utils.NewDummyCommon()

	for i, s := range scenarios {
		icons.SetNerdFontsVersion(lo.Ternary(s.useIcons, "3", ""))

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
