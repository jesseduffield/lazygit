package helpers

import (
	"testing"

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
