package types

// DiffLineType classifies a rendered diff row. It mirrors the per-line "type"
// of the metadata model described in diff-line-metadata-notes.md, so that the
// host-side buffer parser (mechanism #1) and the future pager-emitted OSC
// metadata (#2) produce the same shape.
type DiffLineType int

const (
	DiffLineFileHeader DiffLineType = iota
	DiffLineHunkHeader
	DiffLineContext
	DiffLineAdded
	DiffLineDeleted
	// DiffLineOther is anything that isn't one of the above (e.g. the
	// "\ No newline at end of file" marker). It is also what a backend that
	// can't determine the side reports — delta's lazygit-edit hyperlinks carry
	// no side — so consumers treat it like a non-deletion content line.
	DiffLineOther
)

// DiffLineInfo is the patch-space identity of a rendered diff row, as recovered
// by StagingHelper.GetDiffLineInfo. It is the single shape the focused main view
// and patch explorer consumers act on, regardless of which backend produced it.
type DiffLineInfo struct {
	// Path is the absolute path of the file the line belongs to.
	Path string
	Type DiffLineType
	// NewLine is the line's position in the new file. Set for all content lines
	// (for a deletion it is the new-file position the deletion sits at).
	NewLine int
	// OldLine is the line's position in the old file. Set only for deletions.
	OldLine int
}

// PatchSelectLine returns the source line to land on when diving into the patch
// explorer for this row, in source-line-number space so it survives the patch
// being regenerated. For a deletion it is the old-file line number — two
// consecutive deletions share a new-file line number, so only the old-file
// number tells them apart — and for everything else the new-file line number.
func (self DiffLineInfo) PatchSelectLine() (lineNumber int, isDeletion bool) {
	if self.Type == DiffLineDeleted {
		return self.OldLine, true
	}
	return self.NewLine, false
}

// PullRequestAnchor returns the side ("L"/"R") and line number to anchor a
// GitHub PR deep-link at: the left/old side for a deletion, the right/new side
// otherwise.
func (self DiffLineInfo) PullRequestAnchor() (side string, lineNumber int) {
	if self.Type == DiffLineDeleted {
		return "L", self.OldLine
	}
	return "R", self.NewLine
}
