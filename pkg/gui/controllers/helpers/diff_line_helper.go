package helpers

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type DiffLineHelper struct {
	c *HelperCommon

	// What the probe said about the diff renderer that rendererSignature names, or nil
	// before it has been asked about any (see diffRendererEmitsMetadata).
	rendererEmitsMetadata *bool
	rendererSignature     string
}

func NewDiffLineHelper(c *HelperCommon) *DiffLineHelper {
	return &DiffLineHelper{c: c}
}

// GetDiffLineInfo recovers the identity — file, kind, and old/new line number —
// of the diff row at the given (wrapped) view line of the given view. It is the
// seam every consumer of a diff row goes through, so that how we recover that
// identity can change without them noticing.
//
// There are two ways. A diff renderer that speaks the OSC 1717 protocol states
// the identity of each line it renders, which is the only way to recover it from
// a rendering that doesn't look like a diff any more — columns, or +/- markers
// replaced by colour. Otherwise we parse the view's contents as a unified diff,
// which works for the renderings that keep a diff's structure (no renderer, `git
// diff --color`, a renderer that only colorizes) and fails for the rest.
//
// ok is false when the row's identity can't be recovered, in which case the
// caller must not act on the line at all.
func (self *DiffLineHelper) GetDiffLineInfo(view *gocui.View, viewLineIdx int) (types.DiffLineInfo, bool) {
	identities, ok := self.diffLineIdentitiesAt(view, viewLineIdx)
	if !ok {
		return types.DiffLineInfo{}, false
	}
	return identities[0], true
}

// diffLineIdentitiesAt recovers every diff line the row at the given (wrapped) view
// line shows, left to right. It is GetDiffLineInfo's form for a reader that can't
// settle for the line the row leads with: an end of a selection covers its whole
// row, so where a rendering puts a modification's two halves side by side it covers
// both of them. ok is false when the row's identity can't be recovered at all.
func (self *DiffLineHelper) diffLineIdentitiesAt(
	view *gocui.View, viewLineIdx int,
) ([]types.DiffLineInfo, bool) {
	// The cursor and clicks land on a view line, which counts wrapped segments;
	// the contents are indexed by unwrapped buffer line.
	bufferLineIdx, ok := view.BufferLineForViewLine(viewLineIdx)
	if !ok {
		return nil, false
	}

	contents := view.DiffLineContents()
	if bufferLineIdx >= len(contents) {
		return nil, false
	}

	if identities := self.diffLineIdentitiesFromRecords(contents[bufferLineIdx].Metadata); len(identities) > 0 {
		return self.inRepoTerms(view, identities), true
	}

	parsed, ok := parseDiffLineFromBuffer(diffLineTexts(contents), bufferLineIdx)
	if !ok {
		return nil, false
	}

	return self.inRepoTerms(view, []types.DiffLineInfo{self.diffLineInfo(parsed)}), true
}

// diffLineInfoFromRecords recovers a row's identity from the records the diff
// renderer stated for it, and is what takes precedence over the buffer parse. ok is
// false when the row carries no record we understand, leaving the caller to parse.
//
// A row can carry more than one record, when the rendering puts two diff lines on it
// (a side-by-side row shows a deletion and the addition replacing it); the leftmost
// is the one a reader would call the row's own, so it is the row's identity.
func (self *DiffLineHelper) diffLineInfoFromRecords(metadata []string) (types.DiffLineInfo, bool) {
	identities := self.diffLineIdentitiesFromRecords(metadata)
	if len(identities) == 0 {
		return types.DiffLineInfo{}, false
	}
	return identities[0], true
}

// diffLineIdentitiesFromRecords recovers the identity of every diff line the row's
// records state, left to right. Which of them a reader is after depends on the
// reader: the one the row leads with is the row's own identity (see
// diffLineInfoFromRecords), while a reader looking for a particular line has to
// consider them all, since which of a modification's two halves leads a row is up to
// the rendering.
func (self *DiffLineHelper) diffLineIdentitiesFromRecords(metadata []string) []types.DiffLineInfo {
	identities := make([]types.DiffLineInfo, 0, len(metadata))
	for _, record := range metadata {
		if parsed, ok := parseDiffLineMetadata(record); ok {
			identities = append(identities, self.diffLineInfo(parsed))
		}
	}
	return identities
}

// resolvedDiffLine is one rendered row's recovered identity, plus whether it could
// be recovered at all — the element of the table resolveDiffLines produces.
type resolvedDiffLine struct {
	info types.DiffLineInfo
	ok   bool
}

// resolveDiffLines recovers the identity of every row of a rendered diff in one
// pass, indexed 1:1 with contents. It is the batch form of GetDiffLineInfo, for the
// whole-buffer scans (which change lines are where, which file each row belongs
// to). Resolving row by row would re-run the buffer parser's whole-section parse
// once per row — O(n²) on a large single-file diff — so the buffer parser runs once
// for the whole buffer and the per-row metadata takes precedence on top.
func (self *DiffLineHelper) resolveDiffLines(contents []gocui.DiffLineContent) []resolvedDiffLine {
	bufferParsed := parseAllDiffLinesFromBuffer(diffLineTexts(contents))
	resolved := make([]resolvedDiffLine, len(contents))
	for i, content := range contents {
		if info, ok := self.diffLineInfoFromRecords(content.Metadata); ok {
			resolved[i] = resolvedDiffLine{info, true}
		} else if bufferParsed[i].ok {
			resolved[i] = resolvedDiffLine{self.diffLineInfo(bufferParsed[i].parsed), true}
		}
	}
	return resolved
}

// resolveDiffLineIdentities recovers every diff line each row of a rendered diff
// shows, in one pass, indexed 1:1 with contents. It is resolveDiffLines' form for the
// readers that can't settle for the line a row leads with: looking for a remembered
// line in a new rendering has to consider both halves of a modification, since a
// side-by-side row leads with the deletion whose addition was what got remembered
// under a unified one.
func (self *DiffLineHelper) resolveDiffLineIdentities(contents []gocui.DiffLineContent) [][]types.DiffLineInfo {
	bufferParsed := parseAllDiffLinesFromBuffer(diffLineTexts(contents))
	identities := make([][]types.DiffLineInfo, len(contents))
	for i, content := range contents {
		if fromRecords := self.diffLineIdentitiesFromRecords(content.Metadata); len(fromRecords) > 0 {
			identities[i] = fromRecords
		} else if bufferParsed[i].ok {
			identities[i] = []types.DiffLineInfo{self.diffLineInfo(bufferParsed[i].parsed)}
		}
	}
	return identities
}

// diffLineInfo turns a parser's result into the absolute-path identity consumers
// work with. The path arrives repo-relative from the diff header, but a renderer
// states it however it likes, absolute paths included.
func (self *DiffLineHelper) diffLineInfo(parsed parsedDiffLine) types.DiffLineInfo {
	return diffLineInfoIn(self.c.Git().RepoPaths.WorktreePath(), parsed)
}

// diffLineInfoIn is diffLineInfo against a given worktree, for the callers that can't
// ask which repo we are in where they run: a repo switch replaces it, so only the UI
// thread may read it.
func diffLineInfoIn(worktreePath string, parsed parsedDiffLine) types.DiffLineInfo {
	path := parsed.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(worktreePath, path)
	}

	return types.DiffLineInfo{
		Path:    path,
		Type:    parsed.Type,
		NewLine: parsed.NewLine,
		OldLine: parsed.OldLine,
	}
}
