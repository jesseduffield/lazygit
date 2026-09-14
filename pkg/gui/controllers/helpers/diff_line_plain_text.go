package helpers

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Reading a rendering back to the diff it came from. What a diff view shows is a diff
// renderer's picture of a diff, and a picture is not what you want on your clipboard,
// or in a patch — so the lines of interest are located by identity in the diff itself,
// which the panel that rendered it hands out (types.FocusedMainViewDiffSource).

// PlainDiffOfSelection returns the text of the rows selected in view, in two parts.
//
// aboveDiff is what the selection covers above the diff, as it reads on screen; see
// textAboveDiff. fromDiff is the diff behind the rest: per file the selection touches,
// the run of diff lines from the first of its selected lines to the last, with the
// files in the order the selection meets them.
//
// A run, rather than the matched lines alone, so that what comes out reads as a diff:
// the lines between two selected ones come along even when the rendering didn't show
// them (difftastic leaves out whitespace-only changes) or showed them in another order
// (a side-by-side rendering groups the deletions of a hunk before its additions).
//
// Headers are selected lines like any other. A hunk header names the first line of its
// hunk, in the rendering as in the diff, so it is looked for the way a line of the file
// is; a file header names no line at all and a rendering may spread it over as many
// rows as it likes, so a selection touching one of them takes the whole header.
//
// plainDiff fetches the diff of the given repo-relative files, and is asked only for
// the files the selection touches, so that copying three lines of a commit's diff
// doesn't fetch the whole of it. fromDiff is "" when no selected row could be found in
// the diff, e.g. because the selection covers nothing but a renderer's decoration.
func (self *DiffLineHelper) PlainDiffOfSelection(
	view *gocui.View, first int, last int, plainDiff func(paths []string) string,
) (aboveDiff string, fromDiff string) {
	worktreePath := self.c.Git().RepoPaths.WorktreePath()
	aboveDiff = self.textAboveDiff(view, first, last)

	// The files in the order they are shown, and per file what the selection covers of
	// its diff.
	paths := []string{}
	selection := map[string]*selectedDiffLines{}
	for _, info := range self.DiffLinesInViewRange(view, first, last) {
		// A row that is neither a line of the file nor a header — the "\ No newline at
		// end of file" marker — names nothing to look for. It comes along anyway when it
		// falls within a run.
		if info.Type == types.DiffLineOther {
			continue
		}
		selected, ok := selection[info.Path]
		if !ok {
			paths = append(paths, info.Path)
			selected = &selectedDiffLines{lines: map[patchLine]bool{}}
			selection[info.Path] = selected
		}
		if info.Type == types.DiffLineFileHeader {
			selected.header = true
		} else {
			selected.lines[patchLineOf(info)] = true
		}
	}

	relPaths := repoRelativePaths(worktreePath, paths)
	if len(relPaths) == 0 {
		return aboveDiff, ""
	}

	diffLines := strings.Split(strings.TrimSuffix(plainDiff(relPaths), "\n"), "\n")
	runs := map[string][2]int{}
	for i, parsed := range parseAllDiffLinesFromBuffer(diffLines) {
		if !parsed.ok {
			continue
		}
		info := diffLineInfoIn(worktreePath, parsed.parsed)
		selected := selection[info.Path]
		if selected == nil || !selected.covers(info) {
			continue
		}
		if run, ok := runs[info.Path]; ok {
			runs[info.Path] = [2]int{run[0], i}
		} else {
			runs[info.Path] = [2]int{i, i}
		}
	}

	text := strings.Builder{}
	for _, path := range paths {
		run, ok := runs[path]
		if !ok {
			continue
		}
		for _, line := range diffLines[run[0] : run[1]+1] {
			text.WriteString(line)
			text.WriteString("\n")
		}
	}
	return aboveDiff, text.String()
}

// textAboveDiff returns the text of the selected rows that sit above the diff, as they
// read on screen. What a diff view shows before its first file — a commit's message and
// git's summary of it — is part of no file's diff, so there is nothing to look those
// rows up in and the rendering is all we have of them.
//
// It is "" for a view that shows no file at all: with no diff on screen there is
// nothing for the text to be above, and a rendering we couldn't read as a diff is
// exactly what we don't want on the clipboard.
func (self *DiffLineHelper) textAboveDiff(view *gocui.View, first int, last int) string {
	firstRow, ok := view.BufferLineForViewLine(first)
	if !ok {
		return ""
	}
	lastRow, ok := view.BufferLineForViewLine(last)
	if !ok {
		return ""
	}

	startOfDiff := slices.IndexFunc(self.filePaths(view), func(path string) bool { return path != "" })
	if startOfDiff == -1 {
		return ""
	}

	lastRow = min(lastRow, startOfDiff-1)
	if lastRow < firstRow {
		return ""
	}
	return strings.Join(view.BufferLines()[firstRow:lastRow+1], "\n") + "\n"
}

// selectedDiffLines is what a selection covers of one file's diff.
type selectedDiffLines struct {
	// The lines of the file to look for, hunk headers among them, by the identity that
	// names them in any rendering of the diff.
	lines map[patchLine]bool
	// Whether the file's header is covered, in whole or in part.
	header bool
}

// covers reports whether the given line of a file's diff is one the selection holds.
//
// A file header is answered for by kind rather than looked for: the two ways a row's
// identity is recovered disagree about what to call it — a parse of the diff says the
// file's first line, a renderer's record says the file has no line — and a rendering
// may show the header as one row or as ten. So a selection that touches it holds every
// line of it.
func (self *selectedDiffLines) covers(info types.DiffLineInfo) bool {
	if info.Type == types.DiffLineFileHeader {
		return self.header
	}
	return self.lines[patchLineOf(info)]
}

// repoRelativePaths turns the absolute paths a diff line's identity carries into the
// repo-relative ones git speaks, dropping any that lies outside the worktree — a diff
// renderer states the path however it likes, and one we can't place is one we can't
// ask git about.
func repoRelativePaths(worktreePath string, paths []string) []string {
	relPaths := make([]string, 0, len(paths))
	for _, path := range paths {
		relPath, err := filepath.Rel(worktreePath, path)
		if err != nil || strings.HasPrefix(relPath, "..") {
			continue
		}
		relPaths = append(relPaths, filepath.ToSlash(relPath))
	}
	return relPaths
}
