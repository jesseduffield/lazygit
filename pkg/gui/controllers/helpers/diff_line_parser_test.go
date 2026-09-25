package helpers

import (
	"slices"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gocui"
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
		{"file header", 0, parsedDiffLine{Path: "file1.go", Type: types.DiffLineFileHeader, NewLine: 1}, true},
		{"hunk header", 4, parsedDiffLine{Path: "file1.go", Type: types.DiffLineHunkHeader, NewLine: 1}, true},
		{"context line", 5, parsedDiffLine{Path: "file1.go", Type: types.DiffLineContext, NewLine: 1}, true},
		// The two deletions share new-file line 2 but have distinct old-file lines.
		{"first deletion", 6, parsedDiffLine{Path: "file1.go", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 2}, true},
		{"second deletion", 7, parsedDiffLine{Path: "file1.go", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 3}, true},
		// The second file: its path comes from the second "diff --git" section,
		// and its additions get distinct new-file line numbers.
		{"first addition", 15, parsedDiffLine{Path: "dir/file2.go", Type: types.DiffLineAdded, NewLine: 10}, true},
		{"second addition", 16, parsedDiffLine{Path: "dir/file2.go", Type: types.DiffLineAdded, NewLine: 11}, true},
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
	assert.Equal(t, parsedDiffLine{Path: "new.go", Type: types.DiffLineFileHeader, NewLine: 1}, result)

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
	assert.Equal(t, parsedDiffLine{Path: "new.go", Type: types.DiffLineAdded, NewLine: 2}, result)
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
	assert.Equal(t, parsedDiffLine{Path: "gone.go", Type: types.DiffLineDeleted, NewLine: 0, OldLine: 2}, result)
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
			parsedDiffLine{Path: "modules/xyz", Type: types.DiffLineFileHeader, NewLine: 1},
			result, "line %d", targetIdx)
	}

	result, ok := parseDiffLineFromBuffer(withSubmodule, 7)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{Path: "file.txt", Type: types.DiffLineAdded, NewLine: 3}, result)
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
	assert.Equal(t, "modules/xyz", all[0].parsed.Path)
	assert.Equal(t, "modules/xyz", all[1].parsed.Path)
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
	assert.Equal(t, parsedDiffLine{Path: "dir/file2.go", Type: types.DiffLineAdded, NewLine: 10}, result)

	// Only the section the buffer breaks off in is read that way. One that another
	// section follows is all there, so a hunk short of what its header declares means
	// the rendering restructured the diff, and none of it is resolved.
	shortFirstSection := append(slices.Clone(lines[:8]), lines[9:]...)
	_, ok = parseDiffLineFromBuffer(shortFirstSection, 5)
	assert.False(t, ok)
}

func TestPathFromDiffHeaderField(t *testing.T) {
	scenarios := []struct {
		name     string
		field    string
		expected string
	}{
		{"new side", "b/file.go", "file.go"},
		{"old side", "a/file.go", "file.go"},
		{"a missing file", "/dev/null", "/dev/null"},
		// git terminates the field with a tab when the path has a space in it.
		{"path with a space", "b/with space.go\t", "with space.go"},
		// With core.quotePath enabled (the default) every non-ASCII byte is
		// escaped, and the field is quoted as a whole, prefix included.
		{"non-ASCII path", `"b/caf\303\251.go"`, "café.go"},
		{"non-ASCII path with a space", "\"b/caf\\303\\251 x.go\"\t", "café x.go"},
		{"path with a double quote", `"b/we\"ird.go"`, `we"ird.go`},
		{"path with a backslash", `"b/back\\slash.go"`, `back\slash.go`},
		{"path with a tab", `"b/tab\there.go"`, "tab\there.go"},
		{"undecodable", `"b/unterminated`, ""},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, pathFromDiffHeaderField(s.field))
		})
	}
}

func TestParseDiffLineFromBufferQuotedPath(t *testing.T) {
	// A rename of a file whose name needs quoting, with a content change: the
	// path is quoted on the "diff --git" line and on both of the +++/--- lines.
	renamed := []string{
		`diff --git "a/caf\303\251 old.go" "b/caf\303\251 new.go"`,
		"similarity index 62%",
		`rename from "caf\303\251 old.go"`,
		`rename to "caf\303\251 new.go"`,
		"index 1111111..2222222 100644",
		"--- \"a/caf\\303\\251 old.go\"\t",
		"+++ \"b/caf\\303\\251 new.go\"\t",
		"@@ -1,2 +1,2 @@",
		" apple",
		"-grape",
		"+kiwi",
	}

	result, ok := parseDiffLineFromBuffer(renamed, 10)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{Path: "café new.go", Type: types.DiffLineAdded, NewLine: 2}, result)

	// The same rename without a content change has no +++/--- lines, so the path
	// comes from the "diff --git" line, where both paths are quoted.
	result, ok = parseDiffLineFromBuffer(renamed[:4], 2)
	assert.True(t, ok)
	assert.Equal(t, parsedDiffLine{Path: "café new.go", Type: types.DiffLineFileHeader, NewLine: 1}, result)
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

func TestParseDiffLineMetadata(t *testing.T) {
	scenarios := []struct {
		name     string
		payload  string
		expected parsedDiffLine
		expectOk bool
	}{
		{"context", "1;c;1;;foo.txt", parsedDiffLine{Path: "foo.txt", Type: types.DiffLineContext, NewLine: 1}, true},
		{"added", "1;a;3;;foo.txt", parsedDiffLine{Path: "foo.txt", Type: types.DiffLineAdded, NewLine: 3}, true},
		// A deletion carries both numbers; two consecutive deletions share the
		// new-file line and differ only in the old-file one.
		{"first deletion", "1;d;2;2;foo.txt", parsedDiffLine{Path: "foo.txt", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 2}, true},
		{"second deletion", "1;d;2;3;foo.txt", parsedDiffLine{Path: "foo.txt", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 3}, true},
		// A whole-file deletion has new-file position 0 and the old path.
		{"deleted file", "1;d;0;1;gone.txt", parsedDiffLine{Path: "gone.txt", Type: types.DiffLineDeleted, NewLine: 0, OldLine: 1}, true},
		// The path is the last field, so a ';' within it survives.
		{"path with semicolon", "1;c;5;;weird;name.txt", parsedDiffLine{Path: "weird;name.txt", Type: types.DiffLineContext, NewLine: 5}, true},
		// A renderer may state the path absolutely; the parser keeps it verbatim
		// and leaves resolving it to the caller.
		{"absolute path", "1;a;7;;/abs/foo.txt", parsedDiffLine{Path: "/abs/foo.txt", Type: types.DiffLineAdded, NewLine: 7}, true},
		// A file header has no line number; a hunk header carries the new-file
		// line of the hunk's first line (0 for a whole-file deletion, mirroring
		// `@@ -1,N +0,0 @@`).
		{"file header", "1;f;;;foo.txt", parsedDiffLine{Path: "foo.txt", Type: types.DiffLineFileHeader}, true},
		{"hunk header", "1;h;10;;foo.txt", parsedDiffLine{Path: "foo.txt", Type: types.DiffLineHunkHeader, NewLine: 10}, true},
		{"hunk header of a deleted file", "1;h;0;;gone.txt", parsedDiffLine{Path: "gone.txt", Type: types.DiffLineHunkHeader, NewLine: 0}, true},
		// A file header's line number is always empty, but a renderer that fills
		// it in anyway is taken at its word rather than rejected.
		{"file header with a line number", "1;f;10;;foo.txt", parsedDiffLine{Path: "foo.txt", Type: types.DiffLineFileHeader, NewLine: 10}, true},

		{"unknown version", "2;c;1;;foo.txt", parsedDiffLine{}, false},
		{"unknown type", "1;x;1;;foo.txt", parsedDiffLine{}, false},
		{"too few fields", "1;c;1", parsedDiffLine{}, false},
		{"non-numeric new-line", "1;c;x;;foo.txt", parsedDiffLine{}, false},
		{"non-numeric old-line", "1;d;2;y;foo.txt", parsedDiffLine{}, false},
		// Only a file header may omit the new-file line; on any other kind the
		// record is malformed, and rejecting it falls the row back to the diff
		// text rather than acting on a line number we don't have.
		{"empty new-line on a content line", "1;c;;;foo.txt", parsedDiffLine{}, false},
		{"empty new-line on a hunk header", "1;h;;;foo.txt", parsedDiffLine{}, false},
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

func TestRenderingStatesDiffLines(t *testing.T) {
	row := func(text string, records ...string) gocui.DiffLineContent {
		return gocui.DiffLineContent{Text: text, Metadata: records}
	}

	scenarios := []struct {
		name     string
		contents []gocui.DiffLineContent
		expected bool
	}{
		{
			name:     "a rendering without records is read as a diff",
			contents: []gocui.DiffLineContent{row("diff --git a/foo.txt b/foo.txt"), row("+one")},
			expected: false,
		},
		{
			name:     "a record on any row makes the records the source",
			contents: []gocui.DiffLineContent{row("foo.txt"), row("one", "1;c;1;;foo.txt"), row("")},
			expected: true,
		},
		{
			// A renderer announces the protocol with a record that names no line; a
			// rendering with nothing but that one says nothing about its rows.
			name:     "the version-only handshake record doesn't count",
			contents: []gocui.DiffLineContent{row("foo.txt", "1"), row("one")},
			expected: false,
		},
		{
			name:     "records of a version we don't understand don't count",
			contents: []gocui.DiffLineContent{row("one", "2;c;1;;foo.txt")},
			expected: false,
		},
		{
			name:     "an empty rendering states nothing",
			contents: nil,
			expected: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, renderingStatesDiffLines(s.contents))
		})
	}
}
