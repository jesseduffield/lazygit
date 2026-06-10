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
	// A parseable two-file diff: every row carries its file's path (the headers
	// included), as the buffer parser reports.
	parseable := []string{"a", "a", "a", "a", "b", "b", "b", "b"}

	// The same diff as a restructuring pager emits it: only content lines carry the
	// path; the file/hunk header rows above each file are untagged (empty).
	tagged := []string{"", "", "a", "a", "", "", "b", "b"}

	// Three such files, to exercise navigating from one file's untagged header to the
	// next: the row just above b's header is a's content, so the anchor file must be
	// found by scanning down (b), not up (a) — otherwise next-file would jump back
	// into b and a second `n` couldn't advance.
	taggedThree := []string{"", "", "a", "a", "", "", "b", "b", "", "", "c", "c"}

	scenarios := []struct {
		name     string
		paths    []string
		from     int
		forward  bool
		expected int
		found    bool
	}{
		{"parseable: next file lands on its header", parseable, 1, true, 4, true},
		{"parseable: next from the last file finds nothing", parseable, 6, true, 0, false},
		{"parseable: previous file lands on its header", parseable, 5, false, 0, true},
		{"parseable: previous from the first file finds nothing", parseable, 1, false, 0, false},

		// With only content tagged, both directions still land on the file's top
		// (the untagged header rows), so navigation feels the same.
		{"tagged: next file lands on its header, not its first content", tagged, 2, true, 4, true},
		{"tagged: next from an untagged header still advances", tagged, 0, true, 4, true},
		{"tagged: previous file lands on its header", tagged, 7, false, 0, true},
		{"tagged: previous from the first file finds nothing", tagged, 2, false, 0, false},

		// From b's untagged header (row 4), the anchor file is b (below), so next goes
		// to c and previous goes to a — neither sticks on b.
		{"tagged: next from a middle file's header advances past it", taggedThree, 4, true, 8, true},
		{"tagged: previous from a middle file's header lands on the prior file", taggedThree, 4, false, 0, true},
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
