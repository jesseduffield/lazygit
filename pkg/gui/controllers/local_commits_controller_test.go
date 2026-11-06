package controllers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestFindCommitDragBlock(t *testing.T) {
	commit := func(hash string) *models.Commit {
		return models.NewCommit(&utils.StringPool{}, models.NewCommitOpts{Hash: hash})
	}
	identities := []commitDragIdentity{
		commitDragIdentityForCommit(commit("b")),
		commitDragIdentityForCommit(commit("c")),
	}

	t.Run("finds the original block after selection changes", func(t *testing.T) {
		commits := []*models.Commit{commit("a"), commit("b"), commit("c"), commit("d")}

		actual, startIndex, endIndex, found := findCommitDragBlock(commits, identities)

		assert.True(t, found)
		assert.Equal(t, commits[1:3], actual)
		assert.Equal(t, 1, startIndex)
		assert.Equal(t, 2, endIndex)
	})

	t.Run("rejects a block that is no longer contiguous", func(t *testing.T) {
		_, _, _, found := findCommitDragBlock(
			[]*models.Commit{commit("a"), commit("b"), commit("d"), commit("c")}, identities,
		)

		assert.False(t, found)
	})

	t.Run("rejects an ambiguous block", func(t *testing.T) {
		_, _, _, found := findCommitDragBlock(
			[]*models.Commit{commit("b"), commit("c"), commit("b"), commit("c")}, identities,
		)

		assert.False(t, found)
	})
}

func Test_countSquashableCommitsAbove(t *testing.T) {
	scenarios := []struct {
		name           string
		commits        []*models.Commit
		selectedIdx    int
		rebaseStartIdx int
		expectedResult int
	}{
		{
			name: "no squashable commits",
			commits: []*models.Commit{
				{Name: "abc"},
				{Name: "def"},
				{Name: "ghi"},
			},
			selectedIdx:    2,
			rebaseStartIdx: 2,
			expectedResult: 0,
		},
		{
			name: "some squashable commits, including for the selected commit",
			commits: []*models.Commit{
				{Name: "fixup! def"},
				{Name: "fixup! ghi"},
				{Name: "abc"},
				{Name: "def"},
				{Name: "ghi"},
			},
			selectedIdx:    4,
			rebaseStartIdx: 4,
			expectedResult: 2,
		},
		{
			name: "base commit is below rebase start",
			commits: []*models.Commit{
				{Name: "fixup! def"},
				{Name: "abc"},
				{Name: "def"},
			},
			selectedIdx:    1,
			rebaseStartIdx: 1,
			expectedResult: 0,
		},
		{
			name: "base commit does not exist at all",
			commits: []*models.Commit{
				{Name: "fixup! xyz"},
				{Name: "abc"},
				{Name: "def"},
			},
			selectedIdx:    2,
			rebaseStartIdx: 2,
			expectedResult: 0,
		},
		{
			name: "selected commit is in the middle of fixups",
			commits: []*models.Commit{
				{Name: "fixup! def"},
				{Name: "abc"},
				{Name: "fixup! ghi"},
				{Name: "def"},
				{Name: "ghi"},
			},
			selectedIdx:    1,
			rebaseStartIdx: 4,
			expectedResult: 1,
		},
		{
			name: "selected commit is after rebase start",
			commits: []*models.Commit{
				{Name: "fixup! def"},
				{Name: "abc"},
				{Name: "def"},
				{Name: "ghi"},
			},
			selectedIdx:    3,
			rebaseStartIdx: 2,
			expectedResult: 1,
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expectedResult, countSquashableCommitsAbove(s.commits, s.selectedIdx, s.rebaseStartIdx))
		})
	}
}

func TestMoveCommitsToArgs(t *testing.T) {
	// The commits at indices 2 and 3 are moved. The topmost of them is the
	// head of branch "top"; the commit below them is the head of "x1" and
	// "x2", and the commit at index 0 is the head of "y".
	updateRefBranches := [][]string{{"y"}, nil, {"top"}, nil, {"x1", "x2"}, nil}

	scenarios := []struct {
		name                   string
		destination            commitsDestination
		expectedOffset         int
		expectedDestinationIdx int
		expectedBranchesAbove  []string
	}{
		{
			name:                   "past the commit below",
			destination:            commitsDestination{index: 5},
			expectedOffset:         1,
			expectedDestinationIdx: 5,
			expectedBranchesAbove:  []string{},
		},
		{
			name:                   "below the update-ref branches of the commit below",
			destination:            commitsDestination{index: 4, updateRefsAbove: 2},
			expectedOffset:         0,
			expectedDestinationIdx: 4,
			expectedBranchesAbove:  []string{"x1", "x2", "top"},
		},
		{
			name:                   "between the update-ref branches of the commit below",
			destination:            commitsDestination{index: 4, updateRefsAbove: 1},
			expectedOffset:         0,
			expectedDestinationIdx: 4,
			expectedBranchesAbove:  []string{"x1", "top"},
		},
		{
			name:                   "to the bottom",
			destination:            commitsDestination{index: 6},
			expectedOffset:         2,
			expectedDestinationIdx: 6,
			expectedBranchesAbove:  []string{},
		},
		{
			name:                   "above the update-ref branches of the topmost moved commit",
			destination:            commitsDestination{index: 2},
			expectedOffset:         0,
			expectedDestinationIdx: 4,
			expectedBranchesAbove:  []string{},
		},
		{
			name:                   "below the update-ref branches of the commit above",
			destination:            commitsDestination{index: 1},
			expectedOffset:         -1,
			expectedDestinationIdx: 1,
			expectedBranchesAbove:  []string{},
		},
		{
			name:                   "to the top, below the update-ref branches of the commit there",
			destination:            commitsDestination{index: 0, updateRefsAbove: 1},
			expectedOffset:         -2,
			expectedDestinationIdx: 0,
			expectedBranchesAbove:  []string{"y"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expectedOffset, s.destination.offset(2, 3))
			destinationIdx, branchesAbove := moveCommitsToArgs(updateRefBranches, 2, 3, s.destination)
			assert.Equal(t, s.expectedDestinationIdx, destinationIdx)
			assert.ElementsMatch(t, s.expectedBranchesAbove, branchesAbove)
		})
	}

	t.Run("without update-ref branches", func(t *testing.T) {
		destinationIdx, branchesAbove := moveCommitsToArgs(nil, 2, 3, commitsDestination{index: 5})
		assert.Equal(t, 5, destinationIdx)
		assert.Empty(t, branchesAbove)
	})
}
