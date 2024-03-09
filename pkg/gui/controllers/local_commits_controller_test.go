package controllers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

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

func Test_isFixupCommit(t *testing.T) {
	scenarios := []struct {
		subject                string
		expectedTrimmedSubject string
		expectedIsFixup        bool
	}{
		{
			subject:                "Bla",
			expectedTrimmedSubject: "Bla",
			expectedIsFixup:        false,
		},
		{
			subject:                "fixup Bla",
			expectedTrimmedSubject: "fixup Bla",
			expectedIsFixup:        false,
		},
		{
			subject:                "fixup! Bla",
			expectedTrimmedSubject: "Bla",
			expectedIsFixup:        true,
		},
		{
			subject:                "fixup! fixup! Bla",
			expectedTrimmedSubject: "Bla",
			expectedIsFixup:        true,
		},
		{
			subject:                "amend! squash! Bla",
			expectedTrimmedSubject: "Bla",
			expectedIsFixup:        true,
		},
		{
			subject:                "fixup!",
			expectedTrimmedSubject: "fixup!",
			expectedIsFixup:        false,
		},
	}
	for _, s := range scenarios {
		t.Run(s.subject, func(t *testing.T) {
			trimmedSubject, isFixupCommit := isFixupCommit(s.subject)
			assert.Equal(t, s.expectedTrimmedSubject, trimmedSubject)
			assert.Equal(t, s.expectedIsFixup, isFixupCommit)
		})
	}
}
