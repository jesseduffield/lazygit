package context

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestIndexForGotoBottomInCommits(t *testing.T) {
	hashPool := &utils.StringPool{}
	commit := func(hash string, status models.CommitStatus) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{Hash: hash, Status: status})
	}

	branchCommits := []*models.Commit{
		commit("c3", models.StatusUnpushed),
		commit("c2", models.StatusUnpushed),
		commit("c1", models.StatusPushed),
		commit("m2", models.StatusMerged),
		commit("m1", models.StatusMerged),
	}

	scenarios := []struct {
		testName    string
		selectedIdx int
		commits     []*models.Commit
		expected    int
	}{
		{
			testName:    "goes to the beginning of the branch first",
			selectedIdx: 0,
			commits:     branchCommits,
			expected:    2,
		},
		{
			testName:    "goes from the beginning of the branch to the bottom",
			selectedIdx: 2,
			commits:     branchCommits,
			expected:    4,
		},
		{
			testName:    "goes from a merged commit to the bottom",
			selectedIdx: 3,
			commits:     branchCommits,
			expected:    4,
		},
		{
			testName:    "goes to the bottom if no commit is merged",
			selectedIdx: 0,
			commits: []*models.Commit{
				commit("c2", models.StatusUnpushed),
				commit("c1", models.StatusUnpushed),
			},
			expected: 1,
		},
		{
			testName:    "goes to the bottom if the branch has no commits of its own",
			selectedIdx: 0,
			commits: []*models.Commit{
				commit("m2", models.StatusMerged),
				commit("m1", models.StatusMerged),
			},
			expected: 1,
		},
		{
			testName:    "goes down from a merged commit that an unmerged one follows",
			selectedIdx: 1,
			commits: []*models.Commit{
				// c3 merges m2 of master into the branch
				commit("c3", models.StatusUnpushed),
				commit("m2", models.StatusMerged),
				commit("c2", models.StatusUnpushed),
				commit("c1", models.StatusUnpushed),
				commit("m1", models.StatusMerged),
			},
			expected: 4,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.Equal(t, s.expected, indexForGotoBottomInCommits(s.commits, s.selectedIdx))
		})
	}
}
