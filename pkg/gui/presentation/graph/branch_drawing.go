package graph

import (
	"slices"
	"strconv"
	"strings"
)

// The branch drawing symbols are characters in the Unicode Private Use Area
// (U+F5D0 to U+F60D) for drawing git graphs, introduced by kitty in
// https://github.com/kovidgoyal/kitty/pull/7681. Unlike the box drawing
// characters, they can show how the lines in a cell connect. For example, there
// is a symbol for a cell in which a line from above bends to the left and a
// line from the left bends down.

const branchDrawingHorizontal = "\uf5d0"

// The lines in a cell that isn't a commit
type lineCourses struct {
	up, down   verticalLine
	horizontal bool
}

var branchDrawingConnections = map[lineCourses]string{
	{}:                 " ",
	{horizontal: true}: branchDrawingHorizontal,
	// A horizontal line passing behind a vertical one only shows in the
	// neighbouring cells
	{up: straightLine, down: straightLine}:                   "\uf5d1", // │
	{up: straightLine, down: straightLine, horizontal: true}: "\uf5d1", // │
	{up: lineToLeft}:                                       "\uf5d9", // ╯
	{up: lineToRight}:                                      "\uf5d8", // ╰
	{down: lineToLeft}:                                     "\uf5d7", // ╮
	{down: lineToRight}:                                    "\uf5d6", // ╭
	{up: lineToLeft, horizontal: true}:                     "\uf5e3", // ╯ on ─
	{up: lineToRight, horizontal: true}:                    "\uf5e4", // ╰ on ─
	{down: lineToLeft, horizontal: true}:                   "\uf5e0", // ╮ on ─
	{down: lineToRight, horizontal: true}:                  "\uf5e1", // ╭ on ─
	{up: lineToLeft, down: lineToLeft}:                     "\uf5df", // ╯ and ╮
	{up: lineToRight, down: lineToRight}:                   "\uf5dc", // ╰ and ╭
	{up: lineToLeft, down: lineToLeft, horizontal: true}:   "\uf5e8", // ╯ and ╮ on ─
	{up: lineToRight, down: lineToRight, horizontal: true}: "\uf5e9", // ╰ and ╭ on ─
	{up: lineToLeft, down: lineToRight, horizontal: true}:  "\uf5ec", // ╯ and ╭ on ─
	{up: lineToRight, down: lineToLeft, horizontal: true}:  "\uf5ed", // ╰ and ╮ on ─
}

// The edges of a commit's cell that its lines touch
type commitEdges struct {
	up, down, left, right bool
}

// Commits are drawn as hollow circles, merge commits as filled ones. Both
// connect to the lines at the given edges.
var branchDrawingCommitSymbols = map[commitEdges]struct{ commit, merge string }{
	{}:                                    {"\uf5ef", "\uf5ee"},
	{right: true}:                         {"\uf5f1", "\uf5f0"},
	{left: true}:                          {"\uf5f3", "\uf5f2"},
	{left: true, right: true}:             {"\uf5f5", "\uf5f4"},
	{down: true}:                          {"\uf5f7", "\uf5f6"},
	{up: true}:                            {"\uf5f9", "\uf5f8"},
	{up: true, down: true}:                {"\uf5fb", "\uf5fa"},
	{down: true, right: true}:             {"\uf5fd", "\uf5fc"},
	{down: true, left: true}:              {"\uf5ff", "\uf5fe"},
	{up: true, right: true}:               {"\uf601", "\uf600"},
	{up: true, left: true}:                {"\uf603", "\uf602"},
	{up: true, down: true, right: true}:   {"\uf605", "\uf604"},
	{up: true, down: true, left: true}:    {"\uf607", "\uf606"},
	{down: true, left: true, right: true}: {"\uf609", "\uf608"},
	{up: true, left: true, right: true}:   {"\uf60b", "\uf60a"},
	{up: true, down: true, left: true, right: true}: {"\uf60d", "\uf60c"},
}

func (cell *Cell) branchDrawingChars() (string, string) {
	second := " "
	if cell.right {
		second = branchDrawingHorizontal
	}

	switch cell.cellType {
	case COMMIT, MERGE:
		symbols := branchDrawingCommitSymbols[commitEdges{
			up:    cell.up != noLine,
			down:  cell.down != noLine,
			left:  cell.left,
			right: cell.right,
		}]
		if cell.cellType == MERGE {
			return symbols.merge, second
		}
		return symbols.commit, second
	default:
		if cell.horizontalOnTop {
			return branchDrawingHorizontal, second
		}
		if first, ok := branchDrawingConnections[lineCourses{cell.up, cell.down, cell.horizontal}]; ok {
			return first, second
		}
		// There is no symbol for lines that the layout of the graph never
		// produces, such as a lone half of a line
		first, _ := getBoxDrawingChars(cell.up != noLine, cell.down != noLine, cell.left, cell.right)
		return first, second
	}
}

// The terminals that draw the branch drawing symbols themselves, with the
// first version that draws all of them
var terminalsWithBranchDrawingSymbols = map[string][]int{
	"kitty":   {0, 36, 2},
	"ghostty": {1, 0, 0},
}

// TerminalDrawsBranchDrawingSymbols says whether a terminal is known to draw
// the branch drawing symbols itself, given the name and version it reports
func TerminalDrawsBranchDrawingSymbols(name, version string) bool {
	minVersion, ok := terminalsWithBranchDrawingSymbols[strings.ToLower(name)]
	return ok && slices.Compare(versionNumbers(version), minVersion) >= 0
}

// The numbers at the start of the dot-separated parts of a version, e.g.
// [1 3 0] for "1.3.0-dev+abc"
func versionNumbers(version string) []int {
	var numbers []int
	for part := range strings.SplitSeq(version, ".") {
		digits := part[:len(part)-len(strings.TrimLeft(part, "0123456789"))]
		number, err := strconv.Atoi(digits)
		if err != nil {
			break
		}
		numbers = append(numbers, number)
	}
	return numbers
}
