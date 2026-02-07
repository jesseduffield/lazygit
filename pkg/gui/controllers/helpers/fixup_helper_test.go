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

func TestFixupHelper_FindFixupBaseCommit(t *testing.T) {
	hashPool := &utils.StringPool{}

	type commitDesc struct {
		Hash string
		Name string
	}

	scenarios := []struct {
		subject    string
		commits    []commitDesc
		targetHash string
	}{
		{
			subject: "fixup! Simple feature",
			commits: []commitDesc{
				{Hash: "abc123", Name: "Simple feature"},
			},
			targetHash: "abc123",
		},
		{
			subject: "fixup! abc123",
			commits: []commitDesc{
				{Hash: "abc123", Name: "Something else"},
			},
			targetHash: "abc123",
		},
		{
			subject: "fixup! Partial match",
			commits: []commitDesc{
				{Hash: "def456", Name: "Partial match for this commit"},
			},
			targetHash: "def456",
		},
		{
			subject: "fixup! Multiple matches",
			commits: []commitDesc{
				{Hash: "111111", Name: "Multiple matches"},
				{Hash: "222222", Name: "Multiple matches"},
			},
			targetHash: "222222",
		},
		{
			subject: "fixup! Multiline",
			commits: []commitDesc{
				{Hash: "ghi789", Name: "Multiline\n\nDetailed description here"},
			},
			targetHash: "ghi789",
		},
		{
			subject: "fixup! No match",
			commits: []commitDesc{
				{Hash: "jkl012", Name: "Unrelated work"},
			},
			targetHash: "",
		},
		{
			subject: "fixup! 7777",
			commits: []commitDesc{
				{Hash: "77778888", Name: "Match by partial hash"},
			},
			targetHash: "77778888",
		},
		{
			subject: "fixup! Feature A",
			commits: []commitDesc{
				{Hash: "abc123", Name: "Feature A"},
				{Hash: "def456", Name: "Unrelated"},
				{Hash: "ghi789", Name: "Feature A"},
			},
			targetHash: "ghi789",
		},
	}

	makeCommitFromDesc := func(desc commitDesc, _ int) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{Hash: desc.Hash, Name: desc.Name})
	}

	for _, s := range scenarios {
		t.Run(s.subject, func(t *testing.T) {
			trimmedSubject, isFixupCommit := IsFixupCommit(s.subject)
			assert.Equal(t, true, isFixupCommit)

			commits := lo.Map(s.commits, makeCommitFromDesc)
			found := FindFixupBaseCommit(trimmedSubject, commits)
			if found == nil {
				assert.Equal(t, s.targetHash, "")
			} else {
				assert.Equal(t, s.targetHash, found.Hash())
			}
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
