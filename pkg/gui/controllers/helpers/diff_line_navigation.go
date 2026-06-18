package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
)

// AdjacentChangeBlock returns the view line to move to for next/previous change-block
// navigation in view's displayed diff, starting from anchorViewLine. A "change block"
// is lazygit's notion of a hunk — a run of consecutive added/deleted lines separated
// by context, of which there may be several within one git @@ hunk — matching what
// the staging view's hunk navigation jumps between. forward=true targets the start of
// the next block; forward=false the start of the previous one (from mid-block this
// skips to the previous block, mirroring State.SelectPreviousHunk). ok is false when
// there's no further block, so the caller leaves the view where it is.
func (self *StagingHelper) AdjacentChangeBlock(view *gocui.View, anchorViewLine int, forward bool) (int, bool) {
	contents := view.DiffLineContents()
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, false
	}

	resolved := self.resolveDiffLines(contents)
	isChange := make([]bool, len(resolved))
	for i, r := range resolved {
		isChange[i] = r.ok && r.info.IsChange()
	}

	target, ok := changeBlockStart(isChange, anchor, forward)
	if !ok {
		return 0, false
	}
	return view.ViewLineForBufferLine(target)
}

// AdjacentFile returns the view line to move to for next/previous file navigation in
// view's (possibly multi-file) displayed diff, starting from anchorViewLine: the
// first row belonging to the next/previous file, found where the per-row metadata's
// file changes. ok is false at the first/last file.
func (self *StagingHelper) AdjacentFile(view *gocui.View, anchorViewLine int, forward bool) (int, bool) {
	contents := view.DiffLineContents()
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, false
	}

	resolved := self.resolveDiffLines(contents)
	paths := make([]string, len(resolved))
	for i, r := range resolved {
		if r.ok {
			paths[i] = r.info.Path
		}
	}

	target, ok := fileStart(paths, anchor, forward)
	if !ok {
		return 0, false
	}
	return view.ViewLineForBufferLine(target)
}

// FirstChangeLineInView returns the view line of the first change line at or below
// the top of the viewport, for placing the initial selection when focusing the main
// view: we select the first change the user can already see rather than jumping to
// the top of the diff (which would be jarring when the view is scrolled down). If
// the top visible line is itself mid-change-block, that line is returned, so the
// selection stays put. ok is false when no change line is loaded at or below the top
// (e.g. scrolled into trailing context, or the diff isn't loaded that far yet),
// leaving the caller to fall back.
func (self *StagingHelper) FirstChangeLineInView(view *gocui.View) (int, bool) {
	top, ok := view.BufferLineForViewLine(view.OriginY())
	if !ok {
		return 0, false
	}

	resolved := self.resolveDiffLines(view.DiffLineContents())
	for i := top; i < len(resolved); i++ {
		if resolved[i].ok && resolved[i].info.IsChange() {
			return view.ViewLineForBufferLine(i)
		}
	}
	return 0, false
}

// changeBlockStart finds, in a diff whose lines are flagged by isChange, the first
// line of the change block adjacent to `from` in the given direction. It is the pure
// index arithmetic behind AdjacentChangeBlock, mirroring the staging view's
// State.SelectNextHunk / SelectPreviousHunk line by line.
func changeBlockStart(isChange []bool, from int, forward bool) (int, bool) {
	if forward {
		i := from
		for i < len(isChange) && isChange[i] { // leave the current block
			i++
		}
		for i < len(isChange) && !isChange[i] { // skip the separating context
			i++
		}
		if i < len(isChange) {
			return i, true
		}
		return 0, false
	}

	i := from
	for i >= 0 && isChange[i] { // leave the current block
		i--
	}
	for i >= 0 && !isChange[i] { // skip context, landing on the previous block's last line
		i--
	}
	if i < 0 {
		return 0, false
	}
	for i > 0 && isChange[i-1] { // walk back to that block's first line
		i--
	}
	return i, true
}

// fileStart finds, in a diff whose lines carry the file path they belong to (empty
// for a row no backend could place, e.g. a restructuring pager's file headers), the
// top row of the file adjacent to `from` in the given direction. It is the pure
// index arithmetic behind AdjacentFile. A file is identified by its path, so we find
// where the path changes and then back up over the neighbouring file's unplaced
// header rows, landing on its first row — the `diff --git`/`@@` header when the
// buffer is parseable, or whatever the pager renders above the file's first tagged
// line otherwise.
func fileStart(paths []string, from int, forward bool) (int, bool) {
	anchorPath, ok := anchorFilePath(paths, from)
	if !ok {
		return 0, false
	}

	if forward {
		for i := from; i < len(paths); i++ {
			if paths[i] != "" && paths[i] != anchorPath {
				return backUpOverHeader(paths, i), true
			}
		}
		return 0, false
	}

	// Walk back past the current file (its rows and any unplaced rows) to the
	// previous file's last located row, then back over that whole file to its top.
	i := from
	for i >= 0 && (paths[i] == "" || paths[i] == anchorPath) {
		i--
	}
	if i < 0 {
		return 0, false
	}
	prevPath := paths[i]
	for i > 0 && (paths[i-1] == "" || paths[i-1] == prevPath) {
		i--
	}
	return i, true
}

// backUpOverHeader moves from a file's first located row up over the unplaced header
// rows directly above it, to the file's top. It stops at the previous file's last
// located row, so it never crosses into it.
func backUpOverHeader(paths []string, firstLocated int) int {
	i := firstLocated
	for i > 0 && paths[i-1] == "" {
		i--
	}
	return i
}

// anchorFilePath returns the path of the file the anchor sits in: the first row at or
// below it that carries a path — the file whose content is at or below the top of the
// view — falling back to the nearest above when there's nothing below. Scanning down
// first matters because the anchor is often an untagged file-header row whose nearest
// tagged row is the *previous* file's content just above it; taking that would make
// next-file navigation jump back into the file just left (so a second `n` wouldn't
// advance). ok is false when no row carries a path.
func anchorFilePath(paths []string, from int) (string, bool) {
	for i := from; i < len(paths); i++ {
		if paths[i] != "" {
			return paths[i], true
		}
	}
	for i := from - 1; i >= 0; i-- {
		if paths[i] != "" {
			return paths[i], true
		}
	}
	return "", false
}
