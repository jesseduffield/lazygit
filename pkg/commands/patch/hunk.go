package patch

import "fmt"

// Example hunk:
// @@ -16,2 +14,3 @@ func (f *CommitFile) Description() string {
// 	return f.Name
// -}
// +
// +// test

type Hunk struct {
	// the line number of the first line in the old file ('16' in the above example)
	oldStart int
	// the line number of the first line in the new file ('14' in the above example)
	newStart int
	// the context at the end of the header line (' func (f *CommitFile) Description() string {' in the above example)
	headerContext string
	// the body of the hunk, excluding the header line
	bodyLines []*PatchLine
}

// Returns the number of lines in the hunk in the original file ('2' in the above example)
func (self *Hunk) oldLength() int {
	return nLinesWithKind(self.bodyLines, []PatchLineKind{CONTEXT, DELETION})
}

// Returns the number of lines in the hunk in the new file ('3' in the above example)
func (self *Hunk) newLength() int {
	return nLinesWithKind(self.bodyLines, []PatchLineKind{CONTEXT, ADDITION})
}

// Returns true if the hunk contains any changes (i.e. if it's not just a context hunk).
// We'll end up with a context hunk if we're transforming a patch and one of the hunks
// has no selected lines.
func (self *Hunk) containsChanges() bool {
	return nLinesWithKind(self.bodyLines, []PatchLineKind{ADDITION, DELETION}) > 0
}

// Returns the number of lines in the hunk, including the header line
func (self *Hunk) lineCount() int {
	return len(self.bodyLines) + 1
}

// Returns all lines in the hunk, including the header line
func (self *Hunk) allLines() []*PatchLine {
	lines := []*PatchLine{{Content: self.formatHeaderLine(), Kind: HUNK_HEADER}}
	lines = append(lines, self.bodyLines...)
	return lines
}

// Returns the header line, including the unified diff header and the context
func (self *Hunk) formatHeaderLine() string {
	return fmt.Sprintf("%s%s", self.formatHeaderStart(), self.headerContext)
}

// Returns the first part of the header line i.e. the unified diff part (excluding any context)
func (self *Hunk) formatHeaderStart() string {
	newLengthDisplay := ""
	newLength := self.newLength()
	// if the new length is 1, it's omitted
	if newLength != 1 {
		newLengthDisplay = fmt.Sprintf(",%d", newLength)
	}

	return fmt.Sprintf("@@ -%d,%d +%d%s @@", self.oldStart, self.oldLength(), self.newStart, newLengthDisplay)
}
