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

	if renderingStatesDiffLines(contents) {
		if identities := self.diffLineIdentitiesFromRecords(contents[bufferLineIdx].Metadata); len(identities) > 0 {
			return self.inRepoTerms(view, identities), true
		}
		return nil, false
	}

	parsed, ok := parseDiffLineFromBuffer(diffLineTexts(contents), bufferLineIdx)
	if !ok {
		return nil, false
	}

	return self.inRepoTerms(view, []types.DiffLineInfo{self.diffLineInfo(parsed)}), true
}

// diffLineIdentitiesFromRecords recovers the identity of every diff line the row's
// records state, left to right. A row can carry more than one record, when the
// rendering puts two diff lines on it (a side-by-side row shows a deletion and the
// addition replacing it). Which of them a reader is after depends on the reader: the
// one the row leads with is the row's own identity (see GetDiffLineInfo and
// resolveDiffLines), while a reader looking for a particular line has to consider
// them all, since which of a modification's two halves leads a row is up to the
// rendering.
func (self *DiffLineHelper) diffLineIdentitiesFromRecords(metadata []string) []types.DiffLineInfo {
	return self.diffLineInfos(parseDiffLineRecords(metadata))
}

// diffLineInfoFromRecords recovers a row's own identity from the records the diff
// renderer stated for it. That is the line the row leads with, of the ones
// diffLineIdentitiesFromRecords finds on it. ok is false when the row carries no
// record we understand.
func (self *DiffLineHelper) diffLineInfoFromRecords(metadata []string) (types.DiffLineInfo, bool) {
	identities := self.diffLineIdentitiesFromRecords(metadata)
	if len(identities) == 0 {
		return types.DiffLineInfo{}, false
	}
	return identities[0], true
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
// to). A row's identity is the line it leads with, of those resolveDiffLineIdentities
// finds on it.
func (self *DiffLineHelper) resolveDiffLines(contents []gocui.DiffLineContent) []resolvedDiffLine {
	resolved := make([]resolvedDiffLine, len(contents))
	for i, identities := range self.resolveDiffLineIdentities(contents) {
		if len(identities) > 0 {
			resolved[i] = resolvedDiffLine{identities[0], true}
		}
	}
	return resolved
}

// resolveDiffLineIdentities recovers every diff line each row of a rendered diff
// shows, in one pass, indexed 1:1 with contents. It reads the rendering the way
// GetDiffLineInfo does, by the renderer's records or by parsing it as a unified diff
// (see parseDiffLineIdentities), and is the form of the batch resolver for the
// readers that can't settle for the line a row leads with: looking for a remembered
// line in a new rendering has to consider both halves of a modification, since a
// side-by-side row leads with the deletion whose addition was what got remembered
// under a unified one.
func (self *DiffLineHelper) resolveDiffLineIdentities(contents []gocui.DiffLineContent) [][]types.DiffLineInfo {
	identities := make([][]types.DiffLineInfo, len(contents))
	for i, parsed := range parseDiffLineIdentities(contents) {
		if len(parsed) > 0 {
			identities[i] = self.diffLineInfos(parsed)
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

// diffLineInfos is diffLineInfo over every line of a row.
func (self *DiffLineHelper) diffLineInfos(parsed []parsedDiffLine) []types.DiffLineInfo {
	infos := make([]types.DiffLineInfo, len(parsed))
	for i, line := range parsed {
		infos[i] = self.diffLineInfo(line)
	}
	return infos
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
