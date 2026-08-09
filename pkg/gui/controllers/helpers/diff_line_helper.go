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
// identity can change without them noticing.
//
// There are two ways, and which one is used is settled for the rendering as a
// whole (see renderingStatesDiffLines). A diff renderer that speaks the OSC 1717
// protocol states the identity of each line it renders. That is the only way to
// recover it from a rendering that doesn't look like a diff any more (columns, or
// +/- markers replaced by colour), and a row such a renderer says nothing about
// has no identity. Otherwise we parse the view's contents as a unified diff; this
// works for the renderings that keep a diff's structure (no renderer, `git diff
// --color`, a renderer that only colorizes) and fails for the rest.
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

	contents := view.DiffLineContents()
	if bufferLineIdx >= len(contents) {
		return types.DiffLineInfo{}, false
	}

	if renderingStatesDiffLines(contents) {
		// A row can carry more than one record, when the rendering puts two diff
		// lines on it; the first one is the row's identity, and the leftmost record
		// is the one a reader would call the row's own.
		if metadata := contents[bufferLineIdx].Metadata; len(metadata) > 0 {
			if parsed, ok := parseDiffLineMetadata(metadata[0]); ok {
				return self.diffLineInfo(parsed), true
			}
		}
		return types.DiffLineInfo{}, false
	}

	parsed, ok := parseDiffLineFromBuffer(diffLineTexts(contents), bufferLineIdx)
	if !ok {
		return types.DiffLineInfo{}, false
	}

	return self.diffLineInfo(parsed), true
}

// diffLineInfo turns a parser's result into the absolute-path identity consumers
// work with. The path arrives repo-relative from the diff header, but a renderer
// states it however it likes, absolute paths included.
func (self *DiffLineHelper) diffLineInfo(parsed parsedDiffLine) types.DiffLineInfo {
	path := parsed.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(self.c.Git().RepoPaths.WorktreePath(), path)
	}

	return types.DiffLineInfo{
		Path:    path,
		Type:    parsed.Type,
		NewLine: parsed.NewLine,
		OldLine: parsed.OldLine,
	}
}
