package mergeconflicts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindConflicts(t *testing.T) {
	type scenario struct {
		name     string
		content  string
		expected []*mergeConflict
	}

	scenarios := []scenario{
		{
			name:     "empty",
			content:  "",
			expected: []*mergeConflict{},
		},
		{
			name: "various conflicts",
			content: `++<<<<<<< HEAD
foo
++=======
bar
++>>>>>>> branch

<<<<<<< HEAD: foo/bar/baz.go
foo
bar
=======
baz
>>>>>>> branch

++<<<<<<< MERGE_HEAD
foo
++=======
bar
++>>>>>>> branch

++<<<<<<< Updated upstream
foo
++=======
bar
++>>>>>>> branch

++<<<<<<< ours
foo
++=======
bar
++>>>>>>> branch

<<<<<<< Updated upstream: foo/bar/baz.go
foo
bar
=======
baz
>>>>>>> branch
`,
			expected: []*mergeConflict{
				{
					start:  0,
					middle: 2,
					end:    4,
				},
				{
					start:  6,
					middle: 9,
					end:    11,
				},
				{
					start:  13,
					middle: 15,
					end:    17,
				},
				{
					start:  19,
					middle: 21,
					end:    23,
				},
				{
					start:  25,
					middle: 27,
					end:    29,
				},
				{
					start:  31,
					middle: 34,
					end:    36,
				},
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			assert.EqualValues(t, s.expected, findConflicts(s.content))
		})
	}
}
