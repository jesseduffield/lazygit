package patch

import (
	"strings"

	"github.com/samber/lo"
)

type patchTransformer struct {
	patch *Patch
	opts  TransformOpts
}

type TransformOpts struct {
	// Create a patch that will applied in reverse with `git apply --reverse`.
	// This affects how unselected lines are treated when only parts of a hunk
	// are selected: usually, for unselected lines we change '-' lines to
	// context lines and remove '+' lines, but when Reverse is true we need to
	// turn '+' lines into context lines and remove '-' lines.
	Reverse bool

	// If set, we will replace the original header with one referring to this file name.
	// For staging/unstaging lines we don't want the original header because
	// it makes git confused e.g. when dealing with deleted/added files
	// but with building and applying patches the original header gives git
	// information it needs to cleanly apply patches
	FileNameOverride string

	// Custom patches tend to work better when treating new files as diffs
	// against an empty file. The only case where we need this to be false is
	// when moving a custom patch to an earlier commit; in that case the patch
	// command would fail with the error "file does not exist in index" if we
	// treat it as a diff against an empty file.
	TurnAddedFilesIntoDiffAgainstEmptyFile bool

	// The indices of lines that should be included in the patch.
	IncludedLineIndices []int
}

func transform(patch *Patch, opts TransformOpts) *Patch {
	transformer := &patchTransformer{
		patch: patch,
		opts:  opts,
	}

	return transformer.transform()
}

// helper function that takes a start and end index and returns a slice of all
// indexes inbetween (inclusive)
func ExpandRange(start int, end int) []int {
	expanded := []int{}
	for i := start; i <= end; i++ {
		expanded = append(expanded, i)
	}
	return expanded
}

func (self *patchTransformer) transform() *Patch {
	header := self.transformHeader()
	hunks := self.transformHunks()

	return &Patch{
		header: header,
		hunks:  hunks,
	}
}

func (self *patchTransformer) transformHeader() []string {
	if self.opts.FileNameOverride != "" {
		return []string{
			"--- a/" + self.opts.FileNameOverride,
			"+++ b/" + self.opts.FileNameOverride,
		}
	} else if self.opts.TurnAddedFilesIntoDiffAgainstEmptyFile {
		result := make([]string, 0, len(self.patch.header))
		for idx, line := range self.patch.header {
			if strings.HasPrefix(line, "new file mode") {
				continue
			}
			if line == "--- /dev/null" && strings.HasPrefix(self.patch.header[idx+1], "+++ b/") {
				line = "--- a/" + self.patch.header[idx+1][6:]
			}
			result = append(result, line)
		}
		return result
	} else {
		return self.patch.header
	}
}

func (self *patchTransformer) transformHunks() []*Hunk {
	newHunks := make([]*Hunk, 0, len(self.patch.hunks))

	startOffset := 0
	var formattedHunk *Hunk
	for i, hunk := range self.patch.hunks {
		startOffset, formattedHunk = self.transformHunk(
			hunk,
			startOffset,
			self.patch.HunkStartIdx(i),
		)
		if formattedHunk.containsChanges() {
			newHunks = append(newHunks, formattedHunk)
		}
	}

	return newHunks
}

func (self *patchTransformer) transformHunk(hunk *Hunk, startOffset int, firstLineIdx int) (int, *Hunk) {
	newLines := self.transformHunkLines(hunk, firstLineIdx)
	newNewStart, newStartOffset := self.transformHunkHeader(newLines, hunk.oldStart, startOffset)

	newHunk := &Hunk{
		bodyLines:     newLines,
		oldStart:      hunk.oldStart,
		newStart:      newNewStart,
		headerContext: hunk.headerContext,
	}

	return newStartOffset, newHunk
}

func (self *patchTransformer) transformHunkLines(hunk *Hunk, firstLineIdx int) []*PatchLine {
	skippedNewlineMessageIndex := -1
	newLines := []*PatchLine{}

	for i, line := range hunk.bodyLines {
		lineIdx := i + firstLineIdx + 1 // plus one for header line
		if line.Content == "" {
			break
		}
		isLineSelected := lo.Contains(self.opts.IncludedLineIndices, lineIdx)

		if isLineSelected || (line.Kind == NEWLINE_MESSAGE && skippedNewlineMessageIndex != lineIdx) || line.Kind == CONTEXT {
			newLines = append(newLines, line)
			continue
		}

		if (line.Kind == DELETION && !self.opts.Reverse) || (line.Kind == ADDITION && self.opts.Reverse) {
			content := " " + line.Content[1:]
			newLines = append(newLines, &PatchLine{
				Kind:    CONTEXT,
				Content: content,
			})
			continue
		}

		if line.Kind == ADDITION {
			// we don't want to include the 'newline at end of file' line if it involves an addition we're not including
			skippedNewlineMessageIndex = lineIdx + 1
		}
	}

	return newLines
}

func (self *patchTransformer) transformHunkHeader(newBodyLines []*PatchLine, oldStart int, startOffset int) (int, int) {
	oldLength := nLinesWithKind(newBodyLines, []PatchLineKind{CONTEXT, DELETION})
	newLength := nLinesWithKind(newBodyLines, []PatchLineKind{CONTEXT, ADDITION})

	var newStartOffset int
	// if the hunk went from zero to positive length, we need to increment the starting point by one
	// if the hunk went from positive to zero length, we need to decrement the starting point by one
	if oldLength == 0 {
		newStartOffset = 1
	} else if newLength == 0 {
		newStartOffset = -1
	} else {
		newStartOffset = 0
	}

	newStart := oldStart + startOffset + newStartOffset

	newStartOffset = startOffset + newLength - oldLength

	return newStart, newStartOffset
}
