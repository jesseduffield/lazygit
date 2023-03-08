package patch

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	hunkIndex = utils.Clamp(hunkIndex, 0, len(self.hunks)-1)

	result := len(self.header)
	for i := 0; i < hunkIndex; i++ {
		result += self.hunks[i].lineCount()
	}
	return result
}

// Returns the patch line index of the last line in the given hunk
func (self *Patch) HunkEndIdx(hunkIndex int) int {
	hunkIndex = utils.Clamp(hunkIndex, 0, len(self.hunks)-1)

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
		return hunk.oldStart
	}

	lines := hunk.bodyLines[:idxInHunk-1]
	offset := nLinesWithKind(lines, []PatchLineKind{ADDITION, CONTEXT})
	return hunk.oldStart + offset
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

// Returns the patch line index of the next change (i.e. addition or deletion).
func (self *Patch) GetNextChangeIdx(idx int) int {
	idx = utils.Clamp(idx, 0, self.LineCount()-1)

	lines := self.Lines()

	for i, line := range lines[idx:] {
		if line.isChange() {
			return i + idx
		}
	}

	// there are no changes from the cursor onwards so we'll instead
	// return the index of the last change
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if line.isChange() {
			return i
		}
	}

	// should not be possible
	return 0
}

// Returns the length of the patch in lines
func (self *Patch) LineCount() int {
	count := len(self.header)
	for _, hunk := range self.hunks {
		count += hunk.lineCount()
	}
	return count
}
