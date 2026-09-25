package helpers

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type DiffLineHelper struct {
	c *HelperCommon
}

func NewDiffLineHelper(c *HelperCommon) *DiffLineHelper {
	return &DiffLineHelper{c: c}
}

// GetDiffLineInfo recovers the identity — file, kind, and old/new line number —
// of the diff row at the given (wrapped) view line of the given view. It is the
// seam every consumer of a diff row goes through, so that how we recover that
// identity can change without them noticing: today the only way is to parse the
// view's contents as a unified diff, which works for the renderings that keep a
// diff's structure (no renderer, `git diff --color`, a renderer that only
// colorizes) and fails for the ones that restructure it.
//
// ok is false when the row's identity can't be recovered, in which case the
// caller must not act on the line at all.
func (self *DiffLineHelper) GetDiffLineInfo(view *gocui.View, viewLineIdx int) (types.DiffLineInfo, bool) {
	// The cursor and clicks land on a view line, which counts wrapped segments;
	// the contents are indexed by unwrapped buffer line.
	bufferLineIdx, ok := view.BufferLineForViewLine(viewLineIdx)
	if !ok {
		return types.DiffLineInfo{}, false
	}

	// The lines as written, not as shown: git ends the path field of a diff header
	// with a tab when the path contains a space, and the view shows a tab as spaces.
	parsed, ok := parseDiffLineFromBuffer(view.LinesAsWritten(), bufferLineIdx)
	if !ok {
		return types.DiffLineInfo{}, false
	}

	return self.diffLineInfoFromParsed(parsed), true
}

// diffLineInfoFromParsed turns the parser's repo-relative result into the
// absolute-path identity consumers work with.
func (self *DiffLineHelper) diffLineInfoFromParsed(parsed parsedDiffLine) types.DiffLineInfo {
	return types.DiffLineInfo{
		Path:    filepath.Join(self.c.Git().RepoPaths.WorktreePath(), parsed.RelPath),
		Type:    parsed.Type,
		NewLine: parsed.NewLine,
		OldLine: parsed.OldLine,
	}
}
