package patch

import "github.com/samber/lo"

type PatchLineKind int

const (
	PATCH_HEADER PatchLineKind = iota
	HUNK_HEADER
	ADDITION
	DELETION
	CONTEXT
	NEWLINE_MESSAGE
)

type PatchLine struct {
	Kind    PatchLineKind
	Content string // something like '+ hello' (note the first character is not removed)
}

func (self *PatchLine) isChange() bool {
	return self.Kind == ADDITION || self.Kind == DELETION
}

// Returns the number of lines in the given slice that have one of the given kinds
func nLinesWithKind(lines []*PatchLine, kinds []PatchLineKind) int {
	return lo.CountBy(lines, func(line *PatchLine) bool {
		return lo.Contains(kinds, line.Kind)
	})
}
