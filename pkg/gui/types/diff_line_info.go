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
	// (for a deletion it is the new-file position the deletion sits at) and for
	// hunk headers (the first line of the hunk they head). Not meaningful for
	// file headers: a pager's `f` record carries no line number (so the OSC
	// backend reports 0), and the buffer parser reports 1.
	NewLine int
	// OldLine is the line's position in the old file. Set only for deletions.
	OldLine int
}

// IsChange reports whether the row is an added or deleted line, as opposed to a
// context line or a header. It mirrors patch.PatchLine.IsChange, and is how both
// the context-change scroll preservation (anchor on a surviving change line) and
// the focused main view's hunk navigation (a "hunk" is a block of consecutive
// change lines) classify a rendered row.
func (self DiffLineInfo) IsChange() bool {
	return self.Type == DiffLineAdded || self.Type == DiffLineDeleted
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

// SamePatchLine reports whether two identities point at the same source line of
// the same file. It is how the escape restore matches the line the patch explorer
// had selected against the rows of the focused main view as it re-renders: the
// comparison is in source-line-number space (PatchSelectLine), which survives the
// diff being regenerated, rather than a fragile view-line index.
//
// A backend that can't determine the side (delta's lazygit-edit hyperlinks report
// DiffLineOther) yields a non-deletion identity, so a deletion captured from a
// full-fidelity backend won't match such a row — the restore then just doesn't
// find its line, which is the acceptable degradation for that pager config.
//
// A header row shares its line number with a content row — a hunk header carries
// the hunk's first line — so headers only match headers of the same kind:
// otherwise a restore aiming at a hunk's first content line would land one row up
// on the header above it (or vice versa).
func (self DiffLineInfo) SamePatchLine(other DiffLineInfo) bool {
	if self.Path != other.Path {
		return false
	}
	if self.isHeader() != other.isHeader() {
		return false
	}
	if self.isHeader() && self.Type != other.Type {
		return false
	}
	selfLine, selfIsDeletion := self.PatchSelectLine()
	otherLine, otherIsDeletion := other.PatchSelectLine()
	return selfLine == otherLine && selfIsDeletion == otherIsDeletion
}

func (self DiffLineInfo) isHeader() bool {
	return self.Type == DiffLineFileHeader || self.Type == DiffLineHunkHeader
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
