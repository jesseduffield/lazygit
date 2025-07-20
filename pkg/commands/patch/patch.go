package patch

import (
	"github.com/samber/lo"
)

type Patch struct {
	// header of the patch (split on newlines) e.g.
	// diff --git a/filename b/filename
	// index dcd3485..1ba5540 100644
	// --- a/filename
	// +++ b/filename
	header []string
	// hunks of the patch
	hunks []*Hunk
}

// Returns a new patch with the specified transformation applied (e.g.
// only selecting a subset of changes).
// Leaves the original patch unchanged.
func (self *Patch) Transform(opts TransformOpts) *Patch {
	return transform(self, opts)
}

// Returns the patch as a plain string
func (self *Patch) FormatPlain() string {
	return formatPlain(self)
}

// Returns a range of lines from the patch as a plain string (range is inclusive)
func (self *Patch) FormatRangePlain(startIdx int, endIdx int) string {
	return formatRangePlain(self, startIdx, endIdx)
}

// Returns the patch as a string with ANSI color codes for displaying in a view
func (self *Patch) FormatView(opts FormatViewOpts) string {
	return formatView(self, opts)
}

// Returns the lines of the patch
func (self *Patch) Lines() []*PatchLine {
	lines := []*PatchLine{}
	for _, line := range self.header {
		lines = append(lines, &PatchLine{Content: line, Kind: PATCH_HEADER})
	}

	for _, hunk := range self.hunks {
		lines = append(lines, hunk.allLines()...)
	}

	return lines
}

// Returns the patch line index of the first line in the given hunk
func (self *Patch) HunkStartIdx(hunkIndex int) int {
	hunkIndex = lo.Clamp(hunkIndex, 0, len(self.hunks)-1)

	result := len(self.header)
	for i := range hunkIndex {
		result += self.hunks[i].lineCount()
	}
	return result
}

// Returns the patch line index of the last line in the given hunk
func (self *Patch) HunkEndIdx(hunkIndex int) int {
	hunkIndex = lo.Clamp(hunkIndex, 0, len(self.hunks)-1)

	return self.HunkStartIdx(hunkIndex) + self.hunks[hunkIndex].lineCount() - 1
}

func (self *Patch) ContainsChanges() bool {
	return lo.SomeBy(self.hunks, func(hunk *Hunk) bool {
		return hunk.containsChanges()
	})
}

// Takes a line index in the patch and returns the line number in the new file.
// If the line is a header line, returns 1.
// If the line is a hunk header line, returns the first file line number in that hunk.
// If the line is out of range below, returns the last file line number in the last hunk.
func (self *Patch) LineNumberOfLine(idx int) int {
	if idx < len(self.header) || len(self.hunks) == 0 {
		return 1
	}

	hunkIdx := self.HunkContainingLine(idx)
	// cursor out of range, return last file line number
	if hunkIdx == -1 {
		lastHunk := self.hunks[len(self.hunks)-1]
		return lastHunk.newStart + lastHunk.newLength() - 1
	}

	hunk := self.hunks[hunkIdx]
	hunkStartIdx := self.HunkStartIdx(hunkIdx)
	idxInHunk := idx - hunkStartIdx

	if idxInHunk == 0 {
		return hunk.newStart
	}

	lines := hunk.bodyLines[:idxInHunk-1]
	offset := nLinesWithKind(lines, []PatchLineKind{ADDITION, CONTEXT})
	return hunk.newStart + offset
}

// Returns hunk index containing the line at the given patch line index
func (self *Patch) HunkContainingLine(idx int) int {
	for hunkIdx, hunk := range self.hunks {
		hunkStartIdx := self.HunkStartIdx(hunkIdx)
		if idx >= hunkStartIdx && idx < hunkStartIdx+hunk.lineCount() {
			return hunkIdx
		}
	}
	return -1
}

// Returns the patch line index of the next change (i.e. addition or deletion)
// that matches the same "included" state, given the includedLines. If you don't
// care about included states, pass nil for includedLines and false for included.
func (self *Patch) GetNextChangeIdxOfSameIncludedState(idx int, includedLines []int, included bool) (int, bool) {
	idx = lo.Clamp(idx, 0, self.LineCount()-1)

	lines := self.Lines()

	isMatch := func(i int, line *PatchLine) bool {
		sameIncludedState := lo.Contains(includedLines, i) == included
		return line.IsChange() && sameIncludedState
	}

	for i, line := range lines[idx:] {
		if isMatch(i+idx, line) {
			return i + idx, true
		}
	}

	// there are no changes from the cursor onwards so we'll instead
	// return the index of the last change
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if isMatch(i, line) {
			return i, true
		}
	}

	return 0, false
}

// Returns the patch line index of the next change (i.e. addition or deletion).
func (self *Patch) GetNextChangeIdx(idx int) int {
	result, _ := self.GetNextChangeIdxOfSameIncludedState(idx, nil, false)
	return result
}

// Returns the length of the patch in lines
func (self *Patch) LineCount() int {
	count := len(self.header)
	for _, hunk := range self.hunks {
		count += hunk.lineCount()
	}
	return count
}

// Returns the number of hunks of the patch
func (self *Patch) HunkCount() int {
	return len(self.hunks)
}

// Adjust the given line number (one-based) according to the current patch. The
// patch is supposed to be a diff of an old file state against the working
// directory; the line number is a line number in that old file, and the
// function returns the corresponding line number in the working directory file.
func (self *Patch) AdjustLineNumber(lineNumber int) int {
	adjustedLineNumber := lineNumber
	for _, hunk := range self.hunks {
		if hunk.oldStart >= lineNumber {
			break
		}

		if hunk.oldStart+hunk.oldLength() > lineNumber {
			return hunk.newStart
		}

		adjustedLineNumber += hunk.newLength() - hunk.oldLength()
	}

	return adjustedLineNumber
}

func (self *Patch) IsSingleHunkForWholeFile() bool {
	if len(self.hunks) != 1 {
		return false
	}

	// We consider a patch to be a single hunk for the whole file if it has only additions or
	// deletions but not both, and no context lines. This not quite correct, because it will also
	// return true for a block of added or deleted lines if the diff context size is 0, but in this
	// case you wouldn't be able to stage things anyway, so it doesn't matter.
	bodyLines := self.hunks[0].bodyLines
	return nLinesWithKind(bodyLines, []PatchLineKind{DELETION, CONTEXT}) == 0 ||
		nLinesWithKind(bodyLines, []PatchLineKind{ADDITION, CONTEXT}) == 0
}
