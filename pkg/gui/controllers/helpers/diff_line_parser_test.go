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

func TestParseDiffLineMetadata(t *testing.T) {
	scenarios := []struct {
		name     string
		payload  string
		expected parsedDiffLine
		expectOk bool
	}{
		// These payloads are exactly what the patched delta emits (verified
		// against the real binary; see diff-line-metadata-notes.md §9).
		{"context", "1;c;1;;foo.txt", parsedDiffLine{RelPath: "foo.txt", Type: types.DiffLineContext, NewLine: 1}, true},
		{"added", "1;a;3;;foo.txt", parsedDiffLine{RelPath: "foo.txt", Type: types.DiffLineAdded, NewLine: 3}, true},
		// A deletion carries both numbers; two consecutive deletions share the
		// new-file line and differ only in the old-file line.
		{"first deletion", "1;d;2;2;foo.txt", parsedDiffLine{RelPath: "foo.txt", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 2}, true},
		{"second deletion", "1;d;2;3;foo.txt", parsedDiffLine{RelPath: "foo.txt", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 3}, true},
		// A whole-file deletion has new-file position 0 and the old path.
		{"deleted file", "1;d;0;1;gone.txt", parsedDiffLine{RelPath: "gone.txt", Type: types.DiffLineDeleted, NewLine: 0, OldLine: 1}, true},
		// The path is the last field, so a ';' within it is preserved.
		{"path with semicolon", "1;c;5;;weird;name.txt", parsedDiffLine{RelPath: "weird;name.txt", Type: types.DiffLineContext, NewLine: 5}, true},
		// A pager may emit an absolute path; the parser keeps it verbatim (the
		// caller decides whether to join the worktree path).
		{"absolute path", "1;a;7;;/abs/foo.txt", parsedDiffLine{RelPath: "/abs/foo.txt", Type: types.DiffLineAdded, NewLine: 7}, true},

		{"unknown version", "2;c;1;;foo.txt", parsedDiffLine{}, false},
		{"unknown type", "1;x;1;;foo.txt", parsedDiffLine{}, false},
		{"too few fields", "1;c;1", parsedDiffLine{}, false},
		{"non-numeric new-line", "1;c;x;;foo.txt", parsedDiffLine{}, false},
		{"non-numeric old-line", "1;d;2;y;foo.txt", parsedDiffLine{}, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			result, ok := parseDiffLineMetadata(s.payload)
			assert.Equal(t, s.expectOk, ok)
			if s.expectOk {
				assert.Equal(t, s.expected, result)
			}
		})
	}
}
