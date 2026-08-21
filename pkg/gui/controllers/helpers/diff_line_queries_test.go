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
