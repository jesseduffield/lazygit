package helpers

import (
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
)

// Reading a rendering back to the diff it came from. What a diff view shows is a diff
// renderer's picture of a diff, and a picture is not what you want on your clipboard,
// or in a patch — so the lines of interest are located by identity in the diff itself,
// which the panel that rendered it hands out (types.FocusedMainViewDiffSource).

// PlainDiffOfSelection returns the text of the diff behind the rows selected in view:
// per file the selection touches, the run of diff lines from the first of its selected
// lines to the last, with the files in the order the selection meets them.
//
// A run, rather than the matched lines alone, so that what comes out reads as a diff:
// the lines between two selected ones come along even when the rendering didn't show
// them (difftastic leaves out whitespace-only changes) or showed them in another order
// (a side-by-side rendering groups the deletions of a hunk before its additions).
//
// plainDiff fetches the diff of the given repo-relative files, and is asked only for
// the files the selection touches, so that copying three lines of a commit's diff
// doesn't fetch the whole of it. It returns "" when no selected row could be found in
// the diff, e.g. because the selection covers nothing but a renderer's decoration.
func (self *DiffLineHelper) PlainDiffOfSelection(
	view *gocui.View, first int, last int, plainDiff func(paths []string) string,
) string {
	worktreePath := self.c.Git().RepoPaths.WorktreePath()

	// The files in the order they are shown, and per file the lines to look for. Only
	// content lines: a header names no line of the file, so it can't be looked for, and
	// the headers within a run come along with it anyway.
	paths := []string{}
	selected := map[string]map[patchLine]bool{}
	for _, info := range self.DiffLinesInViewRange(view, first, last) {
		if !info.IsContent() {
			continue
		}
		if _, ok := selected[info.Path]; !ok {
			paths = append(paths, info.Path)
			selected[info.Path] = map[patchLine]bool{}
		}
		selected[info.Path][patchLineOf(info)] = true
	}

	relPaths := repoRelativePaths(worktreePath, paths)
	if len(relPaths) == 0 {
		return ""
	}

	diffLines := strings.Split(strings.TrimSuffix(plainDiff(relPaths), "\n"), "\n")
	runs := map[string][2]int{}
	for i, parsed := range parseAllDiffLinesFromBuffer(diffLines) {
		if !parsed.ok {
			continue
		}
		info := diffLineInfoIn(worktreePath, parsed.parsed)
		if !selected[info.Path][patchLineOf(info)] {
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
	return text.String()
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
