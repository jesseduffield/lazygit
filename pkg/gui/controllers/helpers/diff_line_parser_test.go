package helpers

import (
	"slices"
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

func TestParseDiffLineFromBufferRename(t *testing.T) {
	// A rename with no content change has no hunks and no +++/--- lines, so the
	// path has to come from the "diff --git" line; a rename with a content
	// change has them, and they carry the new path.
	pureRename := strings.Split(`diff --git a/old.go b/new.go
similarity index 100%
rename from old.go
rename to new.go`, "\n")

	result, ok := parseDiffLineFromBuffer(pureRename, 2)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{RelPath: "new.go", Type: types.DiffLineFileHeader, NewLine: 1}, result)

	renameWithModification := strings.Split(`diff --git a/old.go b/new.go
similarity index 62%
rename from old.go
rename to new.go
index 1111111..2222222 100644
--- a/old.go
+++ b/new.go
@@ -1,2 +1,2 @@
 apple
-grape
+kiwi`, "\n")

	result, ok = parseDiffLineFromBuffer(renameWithModification, 10)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{RelPath: "new.go", Type: types.DiffLineAdded, NewLine: 2}, result)
}

func TestParseDiffLineFromBufferDeletedFile(t *testing.T) {
	// The new path is /dev/null, so the identity comes from the old path.
	deletedFile := strings.Split(`diff --git a/gone.go b/gone.go
deleted file mode 100644
index 1111111..0000000
--- a/gone.go
+++ /dev/null
@@ -1,2 +0,0 @@
-apple
-grape`, "\n")

	result, ok := parseDiffLineFromBuffer(deletedFile, 7)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{RelPath: "gone.go", Type: types.DiffLineDeleted, NewLine: 0, OldLine: 2}, result)
}

func TestParseDiffLineFromBufferSubmodule(t *testing.T) {
	// A submodule has no "diff --git" header and no hunks: git opens its section with
	// the commits it moved between and lists them below. So the section ends the one
	// above it, and every row of it belongs to the submodule as a whole.
	withSubmodule := strings.Split(`diff --git a/file.txt b/file.txt
index 1111111..2222222 100644
--- a/file.txt
+++ b/file.txt
@@ -1,2 +1,3 @@
 hello
 world
+world
Submodule modules/xyz a32f27c..2d9f921:
  > bump the thing`, "\n")

	for _, targetIdx := range []int{8, 9} {
		result, ok := parseDiffLineFromBuffer(withSubmodule, targetIdx)
		assert.True(t, ok)
		assert.Equal(t,
			parsedDiffLine{RelPath: "modules/xyz", Type: types.DiffLineFileHeader, NewLine: 1},
			result, "line %d", targetIdx)
	}

	result, ok := parseDiffLineFromBuffer(withSubmodule, 7)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{RelPath: "file.txt", Type: types.DiffLineAdded, NewLine: 3}, result)
}

func TestParseDiffLineFromBufferSubmoduleInARendering(t *testing.T) {
	// A diff renderer passes the lines git writes for a submodule through as they
	// are, while printing nothing below them that a section could end at. The
	// section has to end where what git writes for the submodule ends all the same:
	// the rows below belong to the other files of the diff, and can only be placed
	// by the records the renderer states for them.
	rendered := strings.Split(`Submodule modules/xyz a32f27c..2d9f921:
  > bump the thing

products/a.txt

 one
 two`, "\n")

	all := parseAllDiffLinesFromBuffer(rendered)
	assert.Equal(t, "modules/xyz", all[0].parsed.RelPath)
	assert.Equal(t, "modules/xyz", all[1].parsed.RelPath)
	for i := 2; i < len(rendered); i++ {
		assert.False(t, all[i].ok, "line %d: %q", i, rendered[i])
	}
}

func TestSubmodulePath(t *testing.T) {
	scenarios := []struct {
		name     string
		line     string
		expected string
	}{
		{"moved on", "Submodule modules/xyz a32f27c..2d9f921:", "modules/xyz"},
		{"moved back", "Submodule modules/xyz 2d9f921..a32f27c (rewind):", "modules/xyz"},
		{"moved sideways", "Submodule modules/xyz a32f27c...2d9f921:", "modules/xyz"},
		{"added", "Submodule modules/xyz 0000000...2d9f921 (new submodule)", "modules/xyz"},
		{"removed", "Submodule modules/xyz a32f27c...0000000 (submodule deleted)", "modules/xyz"},
		{"commits missing", "Submodule modules/xyz a32f27c...2d9f921 (commits not present)", "modules/xyz"},
		{"dirty", "Submodule modules/xyz contains modified content", "modules/xyz"},
		{"with something new in it", "Submodule modules/xyz contains untracked content", "modules/xyz"},
		// git writes the path unquoted, so one with a space in it, or one ending in
		// something that reads like a range of commits, is told apart by matching the
		// last range on the line.
		{"path with a space", "Submodule my modules/xyz a32f27c..2d9f921:", "my modules/xyz"},
		{"path reading like a range", "Submodule a32f27c..2d9f921 deadbee..fa1afe1:", "a32f27c..2d9f921"},
		// A line of a file that reads like one of these is indented by the column the
		// diff states the line's side in, so it cannot be mistaken for one.
		{"a line of a file", " Submodule modules/xyz contains modified content", ""},
		{"something else entirely", "Submodule support was added", ""},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, submodulePath(s.line))
		})
	}
}

func TestParseDiffLineFromBufferNotADiff(t *testing.T) {
	// A rendering with no "diff --git" line can't be parsed, so the caller falls
	// back rather than acting on the line.
	bufferLines := []string{"some", "lines", "that", "are not a diff"}
	_, ok := parseDiffLineFromBuffer(bufferLines, 2)
	assert.False(t, ok)
}

func TestParseDiffLineFromBufferGutterMangled(t *testing.T) {
	// A diff renderer that moves the line numbers into a gutter keeps the diff
	// and hunk headers but pushes the +/- markers off the start of each body
	// line, so every line reads as context. The body no longer matches the hunk
	// header, so we refuse to parse rather than return a confident mis-parse.
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

func TestParseDiffLineFromBufferReadInPart(t *testing.T) {
	lines := strings.Split(twoFileDiff, "\n")

	// A long diff is read a screenful at a time, so the buffer breaks off part way
	// through a hunk. The lines that did arrive are resolved all the same, since
	// holding out for the whole hunk would leave the diff on screen with nothing to
	// act on.
	cutShort := lines[:len(lines)-1]
	result, ok := parseDiffLineFromBuffer(cutShort, 15)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{RelPath: "dir/file2.go", Type: types.DiffLineAdded, NewLine: 10}, result)

	// Only the section the buffer breaks off in is read that way. One that another
	// section follows is all there, so a hunk short of what its header declares means
	// the rendering restructured the diff, and none of it is resolved.
	shortFirstSection := append(slices.Clone(lines[:8]), lines[9:]...)
	_, ok = parseDiffLineFromBuffer(shortFirstSection, 5)
	assert.False(t, ok)
}

func TestParseAllDiffLinesFromBuffer(t *testing.T) {
	// Some decoration above the diff, which belongs to no file section: a commit
	// message and a diffstat, as `git show` renders them.
	bufferLines := append(
		[]string{"commit 1234567", "", "    do a thing", "", " file1.go | 2 --", ""},
		strings.Split(twoFileDiff, "\n")...,
	)

	all := parseAllDiffLinesFromBuffer(bufferLines)

	// The batch parse resolves each file section once, and has to agree with
	// resolving the lines one at a time.
	assert.Len(t, all, len(bufferLines))
	for i := range bufferLines {
		parsed, ok := parseDiffLineFromBuffer(bufferLines, i)
		assert.Equal(t, bufferLineParse{parsed, ok}, all[i], "line %d: %q", i, bufferLines[i])
	}

	// The lines above the first file section are left unresolved.
	for i := range 6 {
		assert.False(t, all[i].ok)
	}
	assert.True(t, all[6].ok)
}
