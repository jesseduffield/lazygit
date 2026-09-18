package helpers

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/stretchr/testify/assert"
)

// link is the line the writer is expected to produce for a diffstat entry: the space
// it is indented by, the name linked, and the rest of the line as it came.
func link(name string, rest string) string {
	return " " + style.PrintHyperlink(name, DiffStatLinkScheme+name) + rest
}

// record is the OSC 1717 record a diff renderer speaking the given version of the
// protocol states a line of the given kind with, as it precedes that line in its
// output.
func record(version string, kind string) string {
	return fmt.Sprintf("%s%s;%s;;;pkg/gui.go\x1b\\", diffLineRecordOpener, version, kind)
}

// handshake is the record a renderer announces the protocol with: the version it
// speaks, and nothing about any line. Renderers end their records with either
// terminator, so both turn up.
var (
	handshake    = diffLineRecordOpener + "1\x1b\\"
	handshakeBel = diffLineRecordOpener + "1\x07"
)

func TestDiffStatLinkWriter(t *testing.T) {
	scenarios := []struct {
		name      string
		linkFiles bool
		lines     []string
		expected  []string
	}{
		{
			name:      "links the entries of the diffstat, and nothing after it",
			linkFiles: true,
			lines: []string{
				"commit 1234567",
				"",
				"    A commit message",
				"",
				" pkg/gui.go    | 12 ++++++------",
				" dir/other.go  |  3 ++-",
				" 2 files changed, 8 insertions(+), 7 deletions(-)",
				"",
				"diff --git a/pkg/gui.go b/pkg/gui.go",
				" a context line that reads like an entry | 3 ++-",
			},
			expected: []string{
				"commit 1234567",
				"",
				"    A commit message",
				"",
				link("pkg/gui.go", "    | 12 ++++++------"),
				link("dir/other.go", "  |  3 ++-"),
				" 2 files changed, 8 insertions(+), 7 deletions(-)",
				"",
				"diff --git a/pkg/gui.go b/pkg/gui.go",
				" a context line that reads like an entry | 3 ++-",
			},
		},
		{
			name:      "links a binary file, a file that changes nothing, and a name with spaces",
			linkFiles: true,
			lines: []string{
				" logo.png     | Bin 0 -> 1234 bytes",
				" script.sh    |   0",
				" my file.txt  |   2 +-",
			},
			expected: []string{
				link("logo.png", "     | Bin 0 -> 1234 bytes"),
				link("script.sh", "    |   0"),
				link("my file.txt", "  |   2 +-"),
			},
		},
		{
			name:      "stops looking once the diff itself has begun",
			linkFiles: true,
			lines: []string{
				"diff --git a/pkg/gui.go b/pkg/gui.go",
				" a context line that reads like an entry | 3 ++-",
			},
			expected: []string{
				"diff --git a/pkg/gui.go b/pkg/gui.go",
				" a context line that reads like an entry | 3 ++-",
			},
		},
		{
			name:      "stops looking at the first line a renderer states, whatever kind",
			linkFiles: true,
			lines: []string{
				" pkg/gui.go | 1 +",
				record("1", "c") + " a context line of the diff",
				" this/looks/like/a/diff/stat | 2 +",
			},
			expected: []string{
				link("pkg/gui.go", " | 1 +"),
				record("1", "c") + " a context line of the diff",
				" this/looks/like/a/diff/stat | 2 +",
			},
		},
		{
			name:      "stops for a record of a protocol version it doesn't read",
			linkFiles: true,
			lines: []string{
				record("7", "f") + "── pkg/gui.go ──",
				" this/looks/like/a/diff/stat | 2 +",
			},
			expected: []string{
				record("7", "f") + "── pkg/gui.go ──",
				" this/looks/like/a/diff/stat | 2 +",
			},
		},
		{
			name:      "keeps looking past the handshake, which states no line",
			linkFiles: true,
			lines: []string{
				handshake,
				" pkg/gui.go | 1 +",
			},
			expected: []string{
				handshake,
				link("pkg/gui.go", " | 1 +"),
			},
		},
		{
			// A diff with nothing above its diffstat. The handshake is written with no
			// newline after it, so it runs into the entry naming the first file.
			name:      "links an entry the handshake runs into",
			linkFiles: true,
			lines: []string{
				handshake + " pkg/gui.go   | 1 +",
				" dir/other.go | 2 +-",
			},
			expected: []string{
				handshake + link("pkg/gui.go", "   | 1 +"),
				link("dir/other.go", " | 2 +-"),
			},
		},
		{
			name:      "links an entry a handshake ended with a BEL runs into",
			linkFiles: true,
			lines: []string{
				handshakeBel + " pkg/gui.go | 1 +",
			},
			expected: []string{
				handshakeBel + link("pkg/gui.go", " | 1 +"),
			},
		},
		{
			name:      "leaves a render whose file names lead nowhere alone",
			linkFiles: false,
			lines: []string{
				" pkg/gui.go    | 12 ++++++------",
			},
			expected: []string{
				" pkg/gui.go    | 12 ++++++------",
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			writer := NewDiffStatLinkWriter(buffer)
			writer.BeginRender(scenario.linkFiles)

			for _, line := range scenario.lines {
				written, err := writer.Write([]byte(line + "\n"))
				assert.NoError(t, err)
				// The writer answers for what it was given, not for what it passed on.
				assert.Equal(t, len(line)+1, written)
			}

			assert.Equal(t, strings.Join(scenario.expected, "\n")+"\n", buffer.String())
		})
	}
}

func TestDiffStatLinkWriterStartsLookingAgainWithEachRender(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewDiffStatLinkWriter(buffer)

	for range 2 {
		buffer.Reset()
		writer.BeginRender(true)
		_, _ = writer.Write([]byte(" pkg/gui.go | 1 +\n"))
		_, _ = writer.Write([]byte(" 1 file changed, 1 insertion(+)\n"))

		assert.Equal(t, link("pkg/gui.go", " | 1 +")+
			"\n 1 file changed, 1 insertion(+)\n", buffer.String())
	}
}

func TestFileNamedByDiffStatEntry(t *testing.T) {
	paths := []string{
		"pkg/gui.go",
		"pkg/integration/tests/main_view/jump_to_a_file_of_the_diff.go",
		"vendor/github.com/gdamore/tcell/v3/AUTHORS",
		"pkg/gocui/AUTHORS",
		"renamed.txt",
		"a/very/deeply/nested/directory/structure/some_long_file_name.txt",
	}

	scenarios := []struct {
		name     string
		entry    string
		expected string
	}{
		{
			name:     "a path the diffstat had room for",
			entry:    "pkg/gui.go",
			expected: "pkg/gui.go",
		},
		{
			name:     "a path cut off at a directory boundary",
			entry:    ".../tests/main_view/jump_to_a_file_of_the_diff.go",
			expected: "pkg/integration/tests/main_view/jump_to_a_file_of_the_diff.go",
		},
		{
			name:     "a path cut off inside the file name",
			entry:    "..._long_file_name.txt",
			expected: "a/very/deeply/nested/directory/structure/some_long_file_name.txt",
		},
		{
			name:     "a rename, stated as the part the two paths share",
			entry:    "vendor/github.com/gdamore/tcell/{v2 => v3}/AUTHORS",
			expected: "vendor/github.com/gdamore/tcell/v3/AUTHORS",
		},
		{
			name:     "a rename whose shared part was cut off along with the brace",
			entry:    ".../github.com/jesseduffield => pkg}/gocui/AUTHORS",
			expected: "pkg/gocui/AUTHORS",
		},
		{
			name:     "a rename of paths that share nothing",
			entry:    "original.txt => renamed.txt",
			expected: "renamed.txt",
		},
		{
			name:     "a name of no file of the diff",
			entry:    "pkg/nowhere.go",
			expected: "",
		},
		{
			name:     "a name several files of the diff end with",
			entry:    "AUTHORS",
			expected: "",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			index, ok := fileNamedByDiffStatEntry(scenario.entry, paths)
			if scenario.expected == "" {
				assert.False(t, ok)
				return
			}
			assert.True(t, ok)
			assert.Equal(t, scenario.expected, paths[index])
		})
	}
}
