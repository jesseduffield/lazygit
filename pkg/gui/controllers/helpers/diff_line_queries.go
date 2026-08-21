package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/samber/lo"
)

// The questions a diff view can be asked about what it is showing — where the change
// lines are, which block or file a row belongs to — answered in the view-line terms a
// cursor and a click speak. They are all built on the identities recovered in
// diff_line_helper.go, which is where the answering stops and the recovering starts.

// changeLines resolves view's rendered diff to one flag per buffer line: whether
// that row is a change line (an addition or a deletion), as opposed to context, a
// header, or a row whose identity couldn't be recovered. Those are the rows a
// selection is anchored on and navigation moves between.
func (self *DiffLineHelper) changeLines(view *gocui.View) []bool {
	resolved := self.resolveDiffLines(view.DiffLineContents())
	isChange := make([]bool, len(resolved))
	for i, r := range resolved {
		isChange[i] = r.ok && r.info.IsChange()
	}
	return isChange
}

// FirstChangeLineInView returns the view line of the first change line on screen. It
// is where the selection goes when the main view is focused by keyboard: focusing a
// diff you are reading points at something in it without moving it, so the search
// stops at the bottom of the viewport rather than going after a change further down.
// ok is false when the viewport holds no change line — scrolled into a long stretch
// of context, or past the last change.
func (self *DiffLineHelper) FirstChangeLineInView(view *gocui.View) (int, bool) {
	top, bottom, ok := visibleBufferLines(view)
	if !ok {
		return 0, false
	}

	isChange := self.changeLines(view)
	for i := top; i <= min(bottom, len(isChange)-1); i++ {
		if isChange[i] {
			return view.ViewLineForBufferLine(i)
		}
	}
	return 0, false
}

// visibleBufferLines returns the first and last line of view's content that the
// viewport shows any part of, for the queries that only care about what the user can
// see. The last line is the one at the bottom edge, or the content's last when the
// content ends above it. ok is false for a view showing no content at all.
func visibleBufferLines(view *gocui.View) (int, int, bool) {
	top, ok := view.BufferLineForViewLine(view.OriginY())
	if !ok {
		return 0, 0, false
	}

	lastVisible := min(view.OriginY()+view.InnerHeight(), view.ViewLinesHeight()) - 1
	bottom, ok := view.BufferLineForViewLine(lastVisible)
	if !ok {
		return top, top, true
	}
	return top, bottom, true
}

// ViewHasChangeLines reports whether view's rendered diff holds any change line at
// all, i.e. whether there is anything to select. It is false over a non-diff
// placeholder, and over a diff with nothing in it — an empty commit, a binary file —
// which are the cases where the focused main view shows no selection.
func (self *DiffLineHelper) ViewHasChangeLines(view *gocui.View) bool {
	return lo.Contains(self.changeLines(view), true)
}

// IsChangeLine reports whether the given view line of view's rendered diff is a
// change line rather than context, a header, or an unresolvable row — i.e. whether
// pointing at it points at something a patch could be built from.
func (self *DiffLineHelper) IsChangeLine(view *gocui.View, viewLineIdx int) bool {
	info, ok := self.GetDiffLineInfo(view, viewLineIdx)
	return ok && info.IsChange()
}

// ChangeBlockBounds returns the inclusive view-line range of the change block to
// select in hunk mode around anchorViewLine. A change block is lazygit's notion of a
// hunk — a run of consecutive added or deleted lines bounded by context, of which a
// single git @@ hunk may hold several. When the anchor is context, the block used is
// the first at or below it, or — with nothing below, the cursor sitting past the last
// change — the nearest above, so that hunk mode always has a block to select. ok is
// false only when the diff holds no change line at all.
func (self *DiffLineHelper) ChangeBlockBounds(view *gocui.View, anchorViewLine int) (int, int, bool) {
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, 0, false
	}

	isChange := self.changeLines(view)
	start := anchor
	for start < len(isChange) && !isChange[start] {
		start++
	}
	if start >= len(isChange) {
		for start = min(anchor, len(isChange)-1); start >= 0 && !isChange[start]; start-- {
		}
		if start < 0 {
			return 0, 0, false
		}
	}
	end := start
	for start > 0 && isChange[start-1] {
		start--
	}
	for end < len(isChange)-1 && isChange[end+1] {
		end++
	}

	startView, startOk := view.ViewLineForBufferLine(start)
	// The block's last line goes to its last view line, so that a line the view
	// wrapped is highlighted to its end rather than only where it begins.
	endView, endOk := view.LastViewLineForBufferLine(end)
	if !startOk || !endOk {
		return 0, 0, false
	}
	return startView, endView, true
}

// AdjacentChangeBlock returns the view line to move to for next/previous change-block
// navigation in view's rendered diff, starting from anchorViewLine. A change block is
// lazygit's notion of a hunk (see ChangeBlockBounds). forward=true targets the start
// of the next block, forward=false the start of the previous one — from mid-block that
// means the previous block, rather than the one we are in. ok is false when there's no
// further block, so the caller leaves the view where it is.
func (self *DiffLineHelper) AdjacentChangeBlock(view *gocui.View, anchorViewLine int, forward bool) (int, bool) {
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, false
	}

	target, ok := changeBlockStart(self.changeLines(view), anchor, forward)
	if !ok {
		return 0, false
	}
	return view.ViewLineForBufferLine(target)
}

// changeBlockStart finds, in a diff whose lines are flagged by isChange, the first
// line of the change block adjacent to `from` in the given direction. It is the pure
// index arithmetic behind AdjacentChangeBlock.
func changeBlockStart(isChange []bool, from int, forward bool) (int, bool) {
	if from < 0 || from >= len(isChange) {
		return 0, false
	}

	if forward {
		i := from
		for i < len(isChange) && isChange[i] { // leave the current block
			i++
		}
		for i < len(isChange) && !isChange[i] { // skip the separating context
			i++
		}
		if i == len(isChange) {
			return 0, false
		}
		return i, true
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
