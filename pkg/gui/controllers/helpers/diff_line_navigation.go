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

// DiffFile is a file shown in a (possibly multi-file) diff: its absolute path and the
// view line its section starts at — the row that next/previous-file navigation lands on.
type DiffFile struct {
	Path          string
	FirstViewLine int
}

// FilesInDiff lists the files shown in view's diff, in display order, each paired with
// the view line its section starts at. It is the jump-to-file menu's source: jumping to
// a file goes to its FirstViewLine — the file's first located row, the same row that
// AdjacentFile lands on, so the menu and n/N agree on where each file begins. A file
// whose start row isn't currently mapped to a view line (not loaded yet) is skipped.
func (self *StagingHelper) FilesInDiff(view *gocui.View) []DiffFile {
	resolved := self.resolveDiffLines(view.DiffLineContents())
	paths := make([]string, len(resolved))
	for i, r := range resolved {
		if r.ok {
			paths[i] = r.info.Path
		}
	}

	var files []DiffFile
	seen := map[string]bool{}
	for i, path := range paths {
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true
		if viewLine, ok := view.ViewLineForBufferLine(i); ok {
			files = append(files, DiffFile{Path: path, FirstViewLine: viewLine})
		}
	}
	return files
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

// ViewHasChangeLines reports whether view's displayed diff contains any change line
// (an addition or deletion), i.e. whether there's anything to select. It's false when
// the main view shows a non-diff placeholder ("No changed files", a merge message) or a
// diff that is somehow all context — the cases where the focused main view should show
// no selection.
func (self *StagingHelper) ViewHasChangeLines(view *gocui.View) bool {
	for _, resolved := range self.resolveDiffLines(view.DiffLineContents()) {
		if resolved.ok && resolved.info.IsChange() {
			return true
		}
	}
	return false
}

// IsChangeLine reports whether the given view line of view's displayed diff is a
// change line (an addition or deletion) rather than context, a header, or an
// unparseable row — i.e. whether a click there points at something stageable.
func (self *StagingHelper) IsChangeLine(view *gocui.View, viewLineIdx int) bool {
	info, ok := self.GetDiffLineInfoForView(view, viewLineIdx)
	return ok && info.IsChange()
}

// ChangeBlockBounds returns the view-line range [start, end] of the change block
// (lazygit's notion of a hunk; see AdjacentChangeBlock) to select when entering or
// moving in hunk mode in view's displayed diff. The block is the one containing
// anchorViewLine, or — when that line is context — the first block at or below it
// (matching how toggling hunk mode in the staging view snaps to the next change).
// ok is false when no change line lies at or below the anchor (e.g. scrolled into
// trailing context, or the diff isn't loaded that far yet).
func (self *StagingHelper) ChangeBlockBounds(view *gocui.View, anchorViewLine int) (int, int, bool) {
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, 0, false
	}

	resolved := self.resolveDiffLines(view.DiffLineContents())
	isChange := make([]bool, len(resolved))
	for i, r := range resolved {
		isChange[i] = r.ok && r.info.IsChange()
	}

	// Snap to the first change line at or after the anchor, then expand over the
	// whole contiguous run of change lines around it.
	start := anchor
	for start < len(isChange) && !isChange[start] {
		start++
	}
	if start >= len(isChange) {
		return 0, 0, false
	}
	end := start
	for start > 0 && isChange[start-1] {
		start--
	}
	for end < len(isChange)-1 && isChange[end+1] {
		end++
	}

	startView, ok1 := view.ViewLineForBufferLine(start)
	endView, ok2 := view.ViewLineForBufferLine(end)
	if !ok1 || !ok2 {
		return 0, 0, false
	}
	// ViewLineForBufferLine gives the first view line of the block's last buffer
	// line; extend over any further view lines it wrapped to, so the highlight
	// covers the whole block.
	for endView < view.ViewLinesHeight()-1 {
		next, ok := view.BufferLineForViewLine(endView + 1)
		if !ok || next != end {
			break
		}
		endView++
	}
	return startView, endView, true
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
// for a row no backend could place), the first located row of the file adjacent to
// `from` in the given direction — the row file navigation lands on. It is the pure
// index arithmetic behind AdjacentFile. A file is identified by its path, so we
// find where the path changes, skipping unlocated rows: those are blank separator
// rows between files, or the header rows of a pager that doesn't emit `f`/`h`
// records. So the landing row is the file's header for any conforming source (a
// parseable buffer, or a pager tagging its headers), and the first content line
// under a pager that leaves its headers untagged — an accepted degradation. (An
// earlier version instead backed up over the untagged rows above the first located
// one, to reach the file's top under such pagers; but that overshoots onto the
// blank line above the header whenever the headers themselves are tagged.)
func fileStart(paths []string, from int, forward bool) (int, bool) {
	anchorPath, ok := anchorFilePath(paths, from)
	if !ok {
		return 0, false
	}

	if forward {
		for i := from; i < len(paths); i++ {
			if paths[i] != "" && paths[i] != anchorPath {
				return i, true
			}
		}
		return 0, false
	}

	// Walk back past the current file (its rows and any unlocated rows) to the
	// previous file's last located row, then back over that whole file, landing
	// on its first located row.
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
	for paths[i] != prevPath {
		i++
	}
	return i, true
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
