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
