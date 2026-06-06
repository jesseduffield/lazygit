package helpers

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

// A two-file commit diff as it appears (decolorized) in the main view. file1 has
// two consecutive deletions (grape, pear) that share a new-file line number;
// file2 has two consecutive additions.
const twoFileDiff = `diff --git a/file1.go b/file1.go
index 1111111..2222222 100644
--- a/file1.go
+++ b/file1.go
@@ -1,4 +1,2 @@
 apple
-grape
-pear
 lemon
diff --git a/dir/file2.go b/dir/file2.go
index 3333333..4444444 100644
--- a/dir/file2.go
+++ b/dir/file2.go
@@ -10,2 +9,4 @@ func foo() {
 ctx
+added1
+added2
 ctx2`

func TestParseDiffLineFromBuffer(t *testing.T) {
	bufferLines := strings.Split(twoFileDiff, "\n")

	scenarios := []struct {
		name      string
		targetIdx int
		expected  parsedDiffLine
		expectOk  bool
	}{
		{"file header", 0, parsedDiffLine{RelPath: "file1.go", Type: types.DiffLineFileHeader, NewLine: 1}, true},
		{"hunk header", 4, parsedDiffLine{RelPath: "file1.go", Type: types.DiffLineHunkHeader, NewLine: 1}, true},
		{"context line", 5, parsedDiffLine{RelPath: "file1.go", Type: types.DiffLineContext, NewLine: 1}, true},
		// The two deletions share new-file line 2 but have distinct old-file lines.
		{"first deletion", 6, parsedDiffLine{RelPath: "file1.go", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 2}, true},
		{"second deletion", 7, parsedDiffLine{RelPath: "file1.go", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 3}, true},
		// The second file: its path comes from the second "diff --git" section,
		// and its additions get distinct new-file line numbers.
		{"first addition", 15, parsedDiffLine{RelPath: "dir/file2.go", Type: types.DiffLineAdded, NewLine: 10}, true},
		{"second addition", 16, parsedDiffLine{RelPath: "dir/file2.go", Type: types.DiffLineAdded, NewLine: 11}, true},
		{"out of range", 999, parsedDiffLine{}, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result, ok := parseDiffLineFromBuffer(bufferLines, s.targetIdx)
			assert.Equal(t, s.expectOk, ok)
			if s.expectOk {
				assert.Equal(t, s.expected, result)
			}
		})
	}
}

func TestParseDiffLineFromBufferNotADiff(t *testing.T) {
	// A rendering with no "diff --git" line (e.g. delta's default mode) can't be
	// parsed, so the caller falls back to another backend.
	bufferLines := []string{"some", "lines", "that", "are not a diff"}
	_, ok := parseDiffLineFromBuffer(bufferLines, 2)
	assert.False(t, ok)
}

func TestParseDiffLineFromBufferGutterMangled(t *testing.T) {
	// delta with line-number gutters keeps the diff/hunk headers but pushes the
	// +/- markers off the start of each body line, so every line reads as
	// context. The body no longer matches the hunk header, so we refuse to parse
	// (and the caller falls back) rather than return a confident mis-parse.
	mangled := strings.Split(`diff --git a/file1.txt b/file1.txt
index 1111111..2222222 100644
--- a/file1.txt
+++ b/file1.txt
@@ -1,5 +1,3 @@
  1 ⋮  1 │ apple
  2 ⋮    │-grape
  3 ⋮    │-pear
  4 ⋮  2 │ lemon
  5 ⋮  3 │ mango`, "\n")

	_, ok := parseDiffLineFromBuffer(mangled, 6)
	assert.False(t, ok)
}
