package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFixupHelper_parseDiff(t *testing.T) {
	scenarios := []struct {
		name                     string
		diff                     string
		expectedDeletedLineHunks []*hunk
		expectedAddedLineHunks   []*hunk
	}{
		{
			name:                     "no diff",
			diff:                     "",
			expectedDeletedLineHunks: []*hunk{},
			expectedAddedLineHunks:   []*hunk{},
		},
		{
			name: "hunk with only deleted lines",
			diff: `
diff --git a/file1.txt b/file1.txt
index 9ce8efb33..aaf2a4666 100644
--- a/file1.txt
+++ b/file1.txt
@@ -3 +2,0 @@ bbb
-xxx
`,
			expectedDeletedLineHunks: []*hunk{
				{
					filename:     "file1.txt",
					startLineIdx: 3,
					numLines:     1,
				},
			},
			expectedAddedLineHunks: []*hunk{},
		},
		{
			name: "hunk with deleted and added lines",
			diff: `
diff --git a/file1.txt b/file1.txt
index 9ce8efb33..eb246cf98 100644
--- a/file1.txt
+++ b/file1.txt
@@ -3 +3 @@ bbb
-xxx
+yyy
`,
			expectedDeletedLineHunks: []*hunk{
				{
					filename:     "file1.txt",
					startLineIdx: 3,
					numLines:     1,
				},
			},
			expectedAddedLineHunks: []*hunk{},
		},
		{
			name: "hunk with only added lines",
			diff: `
diff --git a/file1.txt b/file1.txt
index 9ce8efb33..fb5e469e7 100644
--- a/file1.txt
+++ b/file1.txt
@@ -4,0 +5,2 @@ ddd
+xxx
+yyy
`,
			expectedDeletedLineHunks: []*hunk{},
			expectedAddedLineHunks: []*hunk{
				{
					filename:     "file1.txt",
					startLineIdx: 4,
					numLines:     2,
				},
			},
		},
		{
			name: "hunk with dashed lines",
			diff: `
diff --git a/file1.txt b/file1.txt
index 9ce8efb33..fb5e469e7 100644
--- a/file1.txt
+++ b/file1.txt
@@ -3,1 +3,1 @@
--- xxx
+-- yyy
`,
			expectedDeletedLineHunks: []*hunk{
				{
					filename:     "file1.txt",
					startLineIdx: 3,
					numLines:     1,
				},
			},
			expectedAddedLineHunks: []*hunk{},
		},
		{
			name: "several hunks in different files",
			diff: `
diff --git a/file1.txt b/file1.txt
index 9ce8efb33..0632e41b0 100644
--- a/file1.txt
+++ b/file1.txt
@@ -2 +1,0 @@ aaa
-bbb
@@ -4 +3 @@ ccc
-ddd
+xxx
@@ -6,0 +6 @@ fff
+zzz
diff --git a/file2.txt b/file2.txt
index 9ce8efb33..0632e41b0 100644
--- a/file2.txt
+++ b/file2.txt
@@ -0,3 +1,0 @@ aaa
-aaa
-bbb
-ccc
`,
			expectedDeletedLineHunks: []*hunk{
				{
					filename:     "file1.txt",
					startLineIdx: 2,
					numLines:     1,
				},
				{
					filename:     "file1.txt",
					startLineIdx: 4,
					numLines:     1,
				},
				{
					filename:     "file2.txt",
					startLineIdx: 0,
					numLines:     3,
				},
			},
			expectedAddedLineHunks: []*hunk{
				{
					filename:     "file1.txt",
					startLineIdx: 6,
					numLines:     1,
				},
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			deletedLineHunks, addedLineHunks := parseDiff(s.diff)
			assert.Equal(t, s.expectedDeletedLineHunks, deletedLineHunks)
			assert.Equal(t, s.expectedAddedLineHunks, addedLineHunks)
		})
	}
}

func TestFixupHelper_IsFixupCommit(t *testing.T) {
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
			trimmedSubject, isFixupCommit := IsFixupCommit(s.subject)
			assert.Equal(t, s.expectedTrimmedSubject, trimmedSubject)
			assert.Equal(t, s.expectedIsFixup, isFixupCommit)
		})
	}
}

func TestFixupHelper_findCommitWithPreviousMessage(t *testing.T) {
	hashPool := &utils.StringPool{}

	type commitDesc struct {
		Hash string
		Name string
	}

	scenarios := []struct {
		name         string
		commits      []commitDesc
		index        int
		expectedHash string
	}{
		{
			name: "not an amend commit",
			commits: []commitDesc{
				{"abc123", "Some feature"},
			},
			index:        0,
			expectedHash: "",
		},
		{
			name: "fixup commits don't change the message",
			commits: []commitDesc{
				{"abc123", "fixup! Some feature"},
				{"def456", "Some feature"},
			},
			index:        0,
			expectedHash: "",
		},
		{
			name: "a fixup of an amend commit doesn't change the message either",
			commits: []commitDesc{
				{"abc123", "fixup! amend! Some feature"},
				{"def456", "amend! Some feature"},
				{"ghi789", "Some feature"},
			},
			index:        0,
			expectedHash: "",
		},
		{
			name: "base commit right below the amend commit",
			commits: []commitDesc{
				{"abc123", "amend! Some feature"},
				{"def456", "Some feature"},
			},
			index:        0,
			expectedHash: "def456",
		},
		{
			name: "base commit further down the list",
			commits: []commitDesc{
				{"abc123", "amend! Some feature"},
				{"def456", "Unrelated commit"},
				{"ghi789", "Some feature"},
			},
			index:        0,
			expectedHash: "ghi789",
		},
		{
			name: "amend commit in the middle of the list",
			commits: []commitDesc{
				{"abc123", "Unrelated commit"},
				{"def456", "amend! Some feature"},
				{"ghi789", "Some feature"},
			},
			index:        1,
			expectedHash: "ghi789",
		},
		{
			name: "the nearest earlier amend commit holds the message",
			commits: []commitDesc{
				{"abc123", "amend! Some feature"},
				{"def456", "amend! Some feature"},
				{"ghi789", "Some feature"},
			},
			index:        0,
			expectedHash: "def456",
		},
		{
			name: "fixup and squash commits in between are skipped",
			commits: []commitDesc{
				{"abc123", "amend! Some feature"},
				{"def456", "fixup! Some feature"},
				{"ghi789", "squash! Some feature"},
				{"jkl012", "amend! Some feature"},
				{"mno345", "Some feature"},
			},
			index:        0,
			expectedHash: "jkl012",
		},
		{
			name: "amend commit with several prefixes applies to the innermost subject",
			commits: []commitDesc{
				{"abc123", "amend! amend! Some feature"},
				{"def456", "amend! Some feature"},
				{"ghi789", "Some feature"},
			},
			index:        0,
			expectedHash: "def456",
		},
		{
			name: "base commit is not in the list",
			commits: []commitDesc{
				{"abc123", "amend! Some feature"},
				{"def456", "Unrelated commit"},
			},
			index:        0,
			expectedHash: "",
		},
		{
			name: "base commit is above the amend commit",
			commits: []commitDesc{
				{"abc123", "Some feature"},
				{"def456", "amend! Some feature"},
			},
			index:        1,
			expectedHash: "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			commits := lo.Map(s.commits, func(desc commitDesc, _ int) *models.Commit {
				return models.NewCommit(hashPool, models.NewCommitOpts{Hash: desc.Hash, Name: desc.Name})
			})

			result, ok := findCommitWithPreviousMessage(commits, commits[s.index])

			if s.expectedHash == "" {
				assert.False(t, ok)
				assert.Nil(t, result)
			} else {
				assert.True(t, ok)
				assert.Equal(t, s.expectedHash, result.Hash())
			}
		})
	}
}

func TestFixupHelper_messageAfterAmending(t *testing.T) {
	scenarios := []struct {
		name            string
		message         string
		expectedMessage string
	}{
		{
			name:            "subject only",
			message:         "Some feature",
			expectedMessage: "Some feature",
		},
		{
			name:            "subject and body",
			message:         "Some feature\n\nSome description",
			expectedMessage: "Some feature\n\nSome description",
		},
		{
			name:            "amend commit",
			message:         "amend! Some feature\n\nA better subject\n\nSome description",
			expectedMessage: "A better subject\n\nSome description",
		},
		{
			name:            "amend commit without a replacement message",
			message:         "amend! Some feature",
			expectedMessage: "amend! Some feature",
		},
		{
			name:            "fixup commit",
			message:         "fixup! amend! Some feature\n\nSome description",
			expectedMessage: "fixup! amend! Some feature\n\nSome description",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expectedMessage, messageAfterAmending(s.message))
		})
	}
}

func TestFixupHelper_removeFixupCommits(t *testing.T) {
	hashPool := &utils.StringPool{}

	type commitDesc struct {
		Hash string
		Name string
	}

	scenarios := []struct {
		name           string
		commits        []commitDesc
		expectedResult []commitDesc
	}{
		{
			name:           "empty list",
			commits:        []commitDesc{},
			expectedResult: []commitDesc{},
		},
		{
			name: "single commit",
			commits: []commitDesc{
				{"abc123", "Some feature"},
			},
			expectedResult: []commitDesc{
				{"abc123", "Some feature"},
			},
		},
		{
			name: "two unrelated commits",
			commits: []commitDesc{
				{"abc123", "First feature"},
				{"def456", "Second feature"},
			},
			expectedResult: []commitDesc{
				{"abc123", "First feature"},
				{"def456", "Second feature"},
			},
		},
		{
			name: "fixup commit for last commit",
			commits: []commitDesc{
				{"abc123", "fixup! Some feature"},
				{"def456", "Some feature"},
			},
			expectedResult: []commitDesc{
				{"def456", "Some feature"},
			},
		},
		{
			name: "amend and squash commits for last commit",
			commits: []commitDesc{
				{"abc123", "squash! Some feature"},
				{"def456", "amend! Some feature"},
				{"ghi789", "Some feature"},
			},
			expectedResult: []commitDesc{
				{"ghi789", "Some feature"},
			},
		},
		{
			name: "fixup commit for different commit",
			commits: []commitDesc{
				{"abc123", "fixup! Other feature"},
				{"def456", "Some feature"},
			},
			expectedResult: []commitDesc{
				{"abc123", "fixup! Other feature"},
				{"def456", "Some feature"},
			},
		},
		{
			name: "last commit is a fixup itself",
			commits: []commitDesc{
				{"abc123", "fixup! Some feature"},
				{"def456", "fixup! Some feature"},
			},
			expectedResult: []commitDesc{
				{"abc123", "fixup! Some feature"},
				{"def456", "fixup! Some feature"},
			},
		},
		{
			name: "nested fixup commit",
			commits: []commitDesc{
				{"abc123", "fixup! fixup! Some feature"},
				{"def456", "amend! squash! fixup! Some feature"},
				{"ghi789", "Some feature"},
			},
			expectedResult: []commitDesc{
				{"ghi789", "Some feature"},
			},
		},
		{
			name: "fixup commits mixed with unrelated commits",
			commits: []commitDesc{
				{Hash: "abc123", Name: "fixup! Base commit"},
				{Hash: "def456", Name: "Unrelated commit"},
				{Hash: "ghi789", Name: "fixup! Base commit"},
				{Hash: "jkl012", Name: "Base commit"},
			},
			expectedResult: []commitDesc{
				{Hash: "def456", Name: "Unrelated commit"},
				{Hash: "jkl012", Name: "Base commit"},
			},
		},
		{
			name: "only fixup commits for last commit removed, others preserved",
			commits: []commitDesc{
				{Hash: "abc123", Name: "fixup! First feature"},
				{Hash: "def456", Name: "fixup! Second feature"},
				{Hash: "ghi789", Name: "Second feature"},
				{Hash: "jkl012", Name: "First feature"},
			},
			expectedResult: []commitDesc{
				{Hash: "def456", Name: "fixup! Second feature"},
				{Hash: "ghi789", Name: "Second feature"},
				{Hash: "jkl012", Name: "First feature"},
			},
		},
	}

	makeCommitFromDesc := func(desc commitDesc, _ int) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{Hash: desc.Hash, Name: desc.Name})
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			commits := lo.Map(s.commits, makeCommitFromDesc)
			result := removeFixupCommits(commits)
			expectedCommits := lo.Map(s.expectedResult, makeCommitFromDesc)
			assert.Equal(t, expectedCommits, result)
		})
	}
}
