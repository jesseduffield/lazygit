package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSamePatchLine(t *testing.T) {
	scenarios := []struct {
		name     string
		a, b     DiffLineInfo
		expected bool
	}{
		{
			"content lines at the same new-file line match",
			DiffLineInfo{Path: "foo", Type: DiffLineContext, NewLine: 10},
			DiffLineInfo{Path: "foo", Type: DiffLineAdded, NewLine: 10},
			true,
		},
		{
			"different files don't match",
			DiffLineInfo{Path: "foo", Type: DiffLineContext, NewLine: 10},
			DiffLineInfo{Path: "bar", Type: DiffLineContext, NewLine: 10},
			false,
		},
		{
			"a deletion doesn't match a non-deletion at the same position",
			DiffLineInfo{Path: "foo", Type: DiffLineDeleted, NewLine: 10, OldLine: 10},
			DiffLineInfo{Path: "foo", Type: DiffLineContext, NewLine: 10},
			false,
		},
		// A hunk header carries the new-file line of the hunk's first line, so
		// it shares its number with that content line; the header/content guard
		// is what keeps a restore aiming at one from landing on the other.
		{
			"a hunk header doesn't match the hunk's first content line",
			DiffLineInfo{Path: "foo", Type: DiffLineHunkHeader, NewLine: 10},
			DiffLineInfo{Path: "foo", Type: DiffLineContext, NewLine: 10},
			false,
		},
		{
			"a content line doesn't match a hunk header at its line",
			DiffLineInfo{Path: "foo", Type: DiffLineContext, NewLine: 10},
			DiffLineInfo{Path: "foo", Type: DiffLineHunkHeader, NewLine: 10},
			false,
		},
		{
			"hunk headers of the same hunk match",
			DiffLineInfo{Path: "foo", Type: DiffLineHunkHeader, NewLine: 10},
			DiffLineInfo{Path: "foo", Type: DiffLineHunkHeader, NewLine: 10},
			true,
		},
		{
			"hunk headers of different hunks don't match",
			DiffLineInfo{Path: "foo", Type: DiffLineHunkHeader, NewLine: 10},
			DiffLineInfo{Path: "foo", Type: DiffLineHunkHeader, NewLine: 25},
			false,
		},
		{
			"file headers of the same file match",
			DiffLineInfo{Path: "foo", Type: DiffLineFileHeader},
			DiffLineInfo{Path: "foo", Type: DiffLineFileHeader},
			true,
		},
		// A whole-file deletion's hunk header carries new-line 0, the same
		// number a file header reports; only the type tells them apart.
		{
			"a file header doesn't match a deleted file's hunk header",
			DiffLineInfo{Path: "foo", Type: DiffLineFileHeader},
			DiffLineInfo{Path: "foo", Type: DiffLineHunkHeader, NewLine: 0},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, s.a.SamePatchLine(s.b))
		})
	}
}
