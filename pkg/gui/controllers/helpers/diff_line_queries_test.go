package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeBlockStart(t *testing.T) {
	// A diff with three change blocks separated by context:
	//   0 file header   1 hunk header   2 context
	//   3 +   4 +                         (block A)
	//   5 context
	//   6 -                               (block B)
	//   7 context
	//   8 +                               (block C)
	isChange := []bool{false, false, false, true, true, false, true, false, true}

	scenarios := []struct {
		name     string
		from     int
		forward  bool
		expected int
		found    bool
	}{
		{"forward from a header lands on the first block", 0, true, 3, true},
		{"forward from separating context lands on the next block", 5, true, 6, true},
		{"forward from the start of a block skips to the next", 3, true, 6, true},
		{"forward from inside a block skips the rest of it", 4, true, 6, true},
		{"forward from the last block finds nothing", 8, true, 0, false},
		{"backward from a later block lands on the previous one's start", 8, false, 6, true},
		{"backward from a block start lands on the previous block's start", 6, false, 3, true},
		{"backward from inside the first block finds nothing", 4, false, 0, false},
		{"backward from the first block's start finds nothing", 3, false, 0, false},
		{"backward from context lands on the preceding block's start", 7, false, 6, true},
		{"an anchor past the end finds nothing", 9, true, 0, false},
		{"a negative anchor finds nothing", -1, true, 0, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			got, found := changeBlockStart(isChange, s.from, s.forward)
			assert.Equal(t, s.found, found)
			if s.found {
				assert.Equal(t, s.expected, got)
			}
		})
	}
}

func TestFileStart(t *testing.T) {
	// A parseable two-file diff: every row carries its file's path, headers included,
	// as the buffer parser reports it.
	parseable := []string{"a", "a", "a", "a", "b", "b", "b", "b"}

	// The same diff as a renderer that doesn't say which file its headers belong to
	// emits it: only content lines carry the path, so navigation can land no higher
	// than each file's first content line.
	contentOnly := []string{"", "", "a", "a", "", "", "b", "b"}

	// Three such files, to exercise navigating from one file's untagged header to the
	// next: the row just above b's header is a's content, so the anchor's file has to
	// be found by scanning down (b) rather than up (a) — otherwise next-file would
	// jump back into b and a second press couldn't advance.
	contentOnlyThree := []string{"", "", "a", "a", "", "", "b", "b", "", "", "c", "c"}

	// A renderer that does tag its header rows: the file header and the hunk-header box
	// carry the file's path, but the blank separator rows around them carry nothing.
	// Navigation must land on the header's first row, not the blank line above it.
	//    0 blank   1-2 file hdr   3 blank   4-6 hunk hdr box   7 content
	//    8 blank   9-10 file hdr  11 blank  12-13 hunk hdr box  14 content
	headerTagged := []string{"", "a", "a", "", "a", "a", "a", "a", "", "b", "b", "", "b", "b", "b"}

	scenarios := []struct {
		name     string
		paths    []string
		from     int
		forward  bool
		expected int
		found    bool
	}{
		{"forward lands on the next file's header", parseable, 1, true, 4, true},
		{"forward from the last file finds nothing", parseable, 5, true, 0, false},
		{"backward lands on the previous file's header", parseable, 5, false, 0, true},
		{"backward from the first file finds nothing", parseable, 1, false, 0, false},
		{"forward lands on the next file's first content line", contentOnly, 2, true, 6, true},
		{"backward lands on the previous file's first content line", contentOnly, 7, false, 2, true},
		{"forward from an untagged header advances past it", contentOnly, 0, true, 6, true},
		{"a second forward press advances again", contentOnlyThree, 4, true, 10, true},
		{"forward lands on a tagged file header", headerTagged, 7, true, 9, true},
		{"backward lands on a tagged file header", headerTagged, 14, false, 1, true},
		{"a diff with no located rows finds nothing", []string{"", ""}, 0, true, 0, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			got, found := fileStart(s.paths, s.from, s.forward)
			assert.Equal(t, s.found, found)
			if s.found {
				assert.Equal(t, s.expected, got)
			}
		})
	}
}
