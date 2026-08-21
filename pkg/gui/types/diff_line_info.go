package types

import "github.com/jesseduffield/lazygit/pkg/commands/patch"

// DiffLineType classifies a row of a rendered diff.
type DiffLineType int

const (
	DiffLineFileHeader DiffLineType = iota
	DiffLineHunkHeader
	DiffLineContext
	DiffLineAdded
	DiffLineDeleted
	// DiffLineOther is anything that isn't one of the above, e.g. the
	// "\ No newline at end of file" marker.
	DiffLineOther
)

// DiffLineInfo is the identity of a row of a rendered diff in terms of the patch
// it was rendered from: which file the row belongs to, what kind of row it is,
// and where the line sits in the old and new versions of that file. It is what
// lets us act on the line the user is pointing at in a diff view — stage it,
// open it in an editor, keep the cursor on it across a re-render — none of which
// the rendered text alone can tell us.
type DiffLineInfo struct {
	// Path is the absolute path of the file the line belongs to.
	Path string
	Type DiffLineType
	// NewLine is the line's position in the new version of the file. Set for all
	// content lines (for a deletion it is the position the deletion sits at) and
	// for hunk headers (the first line of the hunk they head).
	NewLine int
	// OldLine is the line's position in the old version of the file. Set only
	// for deletions, which are the only rows that need it: two consecutive
	// deletions share a new-file position and differ only here.
	OldLine int
}

// IsChange reports whether the row is an added or deleted line, as opposed to a
// context line or a header. It mirrors patch.PatchLine.IsChange: those are the
// rows a patch is built from, and the rows navigation moves between.
func (self DiffLineInfo) IsChange() bool {
	return self.Type == DiffLineAdded || self.Type == DiffLineDeleted
}

// PatchLineIdentity says which change line of the file the row is, in the terms a patch
// of that file is built and read in: an addition by where it sits in the new version of
// the file, a deletion by where it sat in the old one. Only meaningful for a change row.
func (self DiffLineInfo) PatchLineIdentity() patch.LineIdentity {
	if self.Type == DiffLineDeleted {
		return patch.LineIdentity{LineNumber: self.OldLine, IsDeletion: true}
	}
	return patch.LineIdentity{LineNumber: self.NewLine}
}

// IsContent reports whether the row is a line of the file itself — a change or a
// context line — as opposed to a header or a marker. Those are the rows that have a
// position in the file, and so can be looked for in another rendering of the same
// diff, or in the diff itself.
func (self DiffLineInfo) IsContent() bool {
	return self.IsChange() || self.Type == DiffLineContext
}
