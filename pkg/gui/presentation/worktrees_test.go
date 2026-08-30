package presentation

import (
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func Test_GetWorktreeDisplayString(t *testing.T) {
	tr := &i18n.TranslationSet{
		MainWorktree:    "(main worktree)",
		MissingWorktree: "(missing)",
	}

	scenarios := []struct {
		testName string
		worktree *models.Worktree
		expected []string
	}{
		{
			testName: "Current worktree on a branch",
			worktree: &models.Worktree{
				Name:      "my-worktree",
				Branch:    "my-branch",
				IsCurrent: true,
			},
			expected: []string{"  *", "my-worktree", "my-branch"},
		},
		{
			testName: "Non-current worktree on a branch",
			worktree: &models.Worktree{
				Name:   "my-worktree",
				Branch: "feature-branch",
			},
			expected: []string{"", "my-worktree", "feature-branch"},
		},
		{
			testName: "Main worktree on a branch",
			worktree: &models.Worktree{
				Name:   "repo",
				Branch: "main",
				IsMain: true,
			},
			expected: []string{"", "repo", "main (main worktree)"},
		},
		{
			testName: "Detached HEAD worktree",
			worktree: &models.Worktree{
				Name: "my-worktree",
				Head: "d85cc9d281fa6ae1665c68365fc70e75e82a042d",
			},
			expected: []string{"", "my-worktree", "HEAD detached at d85cc9d2"},
		},
		{
			testName: "Detached HEAD on main worktree",
			worktree: &models.Worktree{
				Name:   "repo",
				IsMain: true,
				Head:   "d85cc9d281fa6ae1665c68365fc70e75e82a042d",
			},
			expected: []string{"", "repo", "HEAD detached at d85cc9d2 (main worktree)"},
		},
		{
			testName: "Worktree with branch takes precedence over head",
			worktree: &models.Worktree{
				Name:   "my-worktree",
				Branch: "my-branch",
				Head:   "d85cc9d281fa6ae1665c68365fc70e75e82a042d",
			},
			expected: []string{"", "my-worktree", "my-branch"},
		},
		{
			testName: "Missing worktree",
			worktree: &models.Worktree{
				Name:          "my-worktree",
				IsPathMissing: true,
				Head:          "d85cc9d281fa6ae1665c68365fc70e75e82a042d",
			},
			expected: []string{"", "my-worktree (missing)", "HEAD detached at d85cc9d2"},
		},
		{
			testName: "Worktree with no branch and no head",
			worktree: &models.Worktree{
				Name: "my-worktree",
			},
			expected: []string{"", "my-worktree", ""},
		},
		{
			testName: "Main worktree with no branch and no head",
			worktree: &models.Worktree{
				Name:   "my-worktree",
				IsMain: true,
			},
			expected: []string{"", "my-worktree", " (main worktree)"},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := GetWorktreeDisplayString(tr, s.worktree)
			assert.Equal(t, s.expected, result)
		})
	}
}
