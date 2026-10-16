package patch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileVersions(t *testing.T) {
	scenarios := []struct {
		testName string
		diff     string
		expected FileVersions
	}{
		{
			testName: "modified file",
			diff: `diff --git a/file1 b/file1
index 9f439a7..20141c9 100644
--- a/file1
+++ b/file1
@@ -1,3 +1,4 @@
 one
+index 1234567..89abcde 100755
 two
 three
`,
			expected: FileVersions{OldBlob: "9f439a7", NewBlob: "20141c9", OldMode: "100644", NewMode: "100644"},
		},
		{
			testName: "added file",
			diff: `diff --git a/file1 b/file1
new file mode 100644
index 0000000..3e75765
--- /dev/null
+++ b/file1
@@ -0,0 +1 @@
+new
`,
			expected: FileVersions{OldBlob: "0000000", NewBlob: "3e75765", NewMode: "100644"},
		},
		{
			testName: "deleted file",
			diff: `diff --git a/file1 b/file1
deleted file mode 100644
index 3e75765..0000000
--- a/file1
+++ /dev/null
@@ -1 +0,0 @@
-new
`,
			expected: FileVersions{OldBlob: "3e75765", NewBlob: "0000000", OldMode: "100644"},
		},
		{
			testName: "changed mode and content",
			diff: `diff --git a/file1 b/file1
old mode 100644
new mode 100755
index 9f439a7..20141c9
--- a/file1
+++ b/file1
@@ -1 +1 @@
-one
+ONE
`,
			expected: FileVersions{OldBlob: "9f439a7", NewBlob: "20141c9", OldMode: "100644", NewMode: "100755"},
		},
		{
			testName: "renamed file with changes",
			diff: `diff --git a/file1 b/file2
similarity index 80%
rename from file1
rename to file2
index 9f439a7..20141c9 100644
--- a/file1
+++ b/file2
@@ -1 +1 @@
-one
+ONE
`,
			expected: FileVersions{OldBlob: "9f439a7", NewBlob: "20141c9", OldMode: "100644", NewMode: "100644"},
		},
		{
			testName: "changed mode only",
			diff: `diff --git a/file1 b/file1
old mode 100644
new mode 100755
`,
			expected: FileVersions{OldMode: "100644", NewMode: "100755"},
		},
		{
			testName: "no header",
			diff:     "",
			expected: FileVersions{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.Equal(t, s.expected, ParseFileVersions(strings.Split(s.diff, "\n")))
			assert.Equal(t, s.expected, Parse(s.diff).FileVersions())
		})
	}
}

func TestFileVersionsSame(t *testing.T) {
	versions := FileVersions{OldBlob: "9f439a7", NewBlob: "20141c9", OldMode: "100644", NewMode: "100644"}

	scenarios := []struct {
		testName string
		other    FileVersions
		expected bool
	}{
		{
			testName: "equal",
			other:    versions,
			expected: true,
		},
		{
			testName: "ids abbreviated to more digits",
			other:    FileVersions{OldBlob: "9f439a739f", NewBlob: "20141c9eb2", OldMode: "100644", NewMode: "100644"},
			expected: true,
		},
		{
			testName: "different new blob",
			other:    FileVersions{OldBlob: "9f439a7", NewBlob: "3e75765", OldMode: "100644", NewMode: "100644"},
			expected: false,
		},
		{
			testName: "different old blob",
			other:    FileVersions{OldBlob: "3e75765", NewBlob: "20141c9", OldMode: "100644", NewMode: "100644"},
			expected: false,
		},
		{
			testName: "different mode",
			other:    FileVersions{OldBlob: "9f439a7", NewBlob: "20141c9", OldMode: "100644", NewMode: "100755"},
			expected: false,
		},
		{
			testName: "no blobs",
			other:    FileVersions{OldMode: "100644", NewMode: "100644"},
			expected: false,
		},
		{
			testName: "nothing at all",
			other:    FileVersions{},
			expected: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.Equal(t, s.expected, versions.Same(s.other))
			assert.Equal(t, s.expected, s.other.Same(versions))
		})
	}

	assert.True(t, FileVersions{}.Same(FileVersions{}))
}
