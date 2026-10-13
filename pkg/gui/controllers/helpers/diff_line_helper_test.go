package helpers

import (
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestDiffLineInfoConverterIn(t *testing.T) {
	worktreePath := t.TempDir()
	absolutePath := filepath.Join(t.TempDir(), "d.txt")
	toDiffLineInfo := diffLineInfoConverterIn(worktreePath)

	lines := []parsedDiffLine{
		{Path: "a.txt", Type: types.DiffLineAdded, NewLine: 1},
		{Path: "a.txt", Type: types.DiffLineDeleted, NewLine: 2, OldLine: 1},
		{Path: "dir/b.txt", Type: types.DiffLineContext, NewLine: 3},
		{Path: absolutePath, Type: types.DiffLineAdded, NewLine: 4},
		{Path: "a.txt", Type: types.DiffLineAdded, NewLine: 5},
	}
	expected := []types.DiffLineInfo{
		{Path: filepath.Join(worktreePath, "a.txt"), Type: types.DiffLineAdded, NewLine: 1},
		{Path: filepath.Join(worktreePath, "a.txt"), Type: types.DiffLineDeleted, NewLine: 2, OldLine: 1},
		{Path: filepath.Join(worktreePath, "dir", "b.txt"), Type: types.DiffLineContext, NewLine: 3},
		{Path: absolutePath, Type: types.DiffLineAdded, NewLine: 4},
		{Path: filepath.Join(worktreePath, "a.txt"), Type: types.DiffLineAdded, NewLine: 5},
	}

	for i, line := range lines {
		assert.Equal(t, expected[i], toDiffLineInfo(line), "line %d", i)
	}
}
