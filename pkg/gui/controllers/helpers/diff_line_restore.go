package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/samber/lo"
)

// Keeping a diff view where it is when the same diff is rendered again differently.
// The line the user is on is remembered by identity (diff_line_helper.go), because a
// new rendering puts it on a different line of the view — and may not have it at all,
// which is what the fallbacks below are for.

// diffLineAnchor is a line for a restore to land on: the identity to find it by in
// the new rendering, and the screen row it was on, so that it can be put back there.
type diffLineAnchor struct {
	identity types.DiffLineInfo
	row      int
}

// PreserveDiffPositionOnRerender remembers where a diff view is and puts it back
// there as it next re-renders, instead of leaving the user at the top of a new
// rendering of the diff they were already reading. Call it on the view about to be
// re-rendered, right before triggering the re-render — on both panes of the main
// window where both are being rendered again, since either of them may hold the diff
// being read; a pane that isn't showing is left alone.
//
// The line to keep is the end of the selection that is on screen, and the middle
// visible line when there is no selection or the whole of it has been scrolled out of
// sight — what the user is looking at, rather than the view's top edge or a selection
// they have long since left behind. It may not survive the re-render: a context line
// goes when the context size shrinks, and a whole hunk or file goes when whitespace
// stops counting. So the lines around it come along as fallbacks and the view lands on
// the nearest one that is still there, put back on the screen row it was on. With none
// of them left — and with a renderer that says nothing about its rows there is nothing
// to look for in the first place — the view keeps the scroll offset it had, which is
// still nearer to what was being read than the top of the diff.
//
// An off-screen selection is still put back on the diff line it was on, wherever the
// new rendering has that; it is only the view that stays where it is.
//
// A range or hunk selection has a second end, which is remembered the same way, so
// that it still covers the same lines of the diff afterwards.
func (self *DiffLineHelper) PreserveDiffPositionOnRerender(view *gocui.View) {
	// A view that isn't the one its window is currently showing — the merge-conflicts
	// view takes the main window over — isn't the one about to be re-rendered, so a
	// restore installed on it would sit there and claim a later render instead.
	if !view.Visible {
		return
	}

	// The re-render is produced by a different command from the one behind what is on
	// screen — another context size, another renderer — so without being told otherwise
	// it would be taken for content the user has never seen and shown from the top.
	// Whether or not a line of the old rendering can be found in the new one, the offset
	// into it is nearer to where they were reading than the top is.
	if manager := self.c.GetViewBufferManagerForView(view); manager != nil {
		manager.SetKeepScrollPositionForNextTask()
	}

	showSelection := view.Highlight
	anchorViewLine := view.MiddleVisibleLineIdx()
	farEnd, hasFarEnd := types.DiffLineInfo{}, false
	// A cursor that has been scrolled away from is put back by its own lines rather
	// than by the anchor's, so that it comes out on the same line of the diff without
	// the view having to go there.
	var cursorCandidates []diffLineAnchor
	if showSelection {
		farEnd, hasFarEnd = self.selectionFarEndIdentity(view)
		if end, ok := visibleSelectionEnd(view); ok {
			anchorViewLine = end
		}
		if anchorViewLine != view.SelectedLineIdx() {
			cursorCandidates = self.nearbyDiffLines(view, view.SelectedLineIdx())
		}
	}

	self.restoreDiffLinePositionOnRerender(view, self.nearbyDiffLines(view, anchorViewLine),
		func(anchor diffLineAnchor, viewLine int) {
			// Put the line back on the screen row it was on, clamped into the view for
			// the fallback lines, which can come from off screen.
			row := lo.Clamp(anchor.row, 0, max(0, view.InnerHeight()-1))
			view.SetOrigin(0, max(0, viewLine-row))
			if showSelection {
				// Put the far end back before the cursor, so that the selection covers
				// the same lines again; a selection whose far end didn't survive the
				// re-render is left as the single line we landed on. The origin is
				// already where it should be, so moving the cursor mustn't scroll.
				view.CancelRangeSelect()
				cursorViewLine := self.selectionLine(view, cursorCandidates, viewLine)
				if hasFarEnd {
					if farEndViewLine, ok := self.findDiffLine(view, farEnd); ok {
						cursorViewLine, farEndViewLine = coverWholeLines(view, cursorViewLine, farEndViewLine)
						view.SetRangeSelectStart(farEndViewLine)
					}
				}
				view.FocusPoint(0, cursorViewLine, false)
			}
		})
}

// coverWholeLines moves the two ends of a restored selection out to the edges of the
// diff lines they are on, so that the selection covers those lines whole. Both ends
// arrive on the first view line of their diff line, which is where looking one up by
// identity lands, and the view draws a line it wraps as several — of which a
// selection of that line means all.
func coverWholeLines(view *gocui.View, cursorViewLine int, farEndViewLine int) (int, int) {
	if cursorViewLine <= farEndViewLine {
		return cursorViewLine, lastViewLineOfSameDiffLine(view, farEndViewLine)
	}
	return lastViewLineOfSameDiffLine(view, cursorViewLine), farEndViewLine
}

// lastViewLineOfSameDiffLine returns the last view line showing the same line of the
// diff as the given one, which is that line itself unless the view wrapped it.
func lastViewLineOfSameDiffLine(view *gocui.View, viewLine int) int {
	bufferLine, ok := view.BufferLineForViewLine(viewLine)
	if !ok {
		return viewLine
	}
	if last, ok := view.LastViewLineForBufferLine(bufferLine); ok {
		return last
	}
	return viewLine
}

// visibleSelectionEnd returns the end of the selection to keep in place across a
// re-render: the selected line when it is on screen, and the range's other end when
// that is and the selected line isn't — a range can be long enough for the user to be
// looking at one end of it with the other far away. ok is false when the whole
// selection is off screen, and there is nothing of it to keep in place.
func visibleSelectionEnd(view *gocui.View) (int, bool) {
	if view.IsLineVisible(view.SelectedLineIdx()) {
		return view.SelectedLineIdx(), true
	}
	if farEnd, _, ok := selectionFarEndViewLine(view); ok && view.IsLineVisible(farEnd) {
		return farEnd, true
	}
	return 0, false
}

// selectionLine returns the line to put the cursor on once a re-render is on screen:
// the line the position anchor landed on, which is the selected one whenever it was
// on screen, and otherwise the nearest surviving line to where the selection was —
// found among its own candidates, since the anchor's are a search of the diff from
// somewhere else entirely.
func (self *DiffLineHelper) selectionLine(
	view *gocui.View, candidates []diffLineAnchor, anchorViewLine int,
) int {
	if len(candidates) == 0 {
		return anchorViewLine
	}
	_, bufferLine := self.nearestSurvivingCandidate(view.DiffLineContents(), candidates)
	if bufferLine == -1 {
		return anchorViewLine
	}
	if viewLine, ok := view.ViewLineForBufferLine(bufferLine); ok {
		return viewLine
	}
	return anchorViewLine
}

// selectionFarEndIdentity returns the identity of the end of a range or hunk
// selection the cursor isn't on, so that a re-render can put it back. ok is false for
// a selection that is only a cursor, where restoring that is the whole job, and for
// an end that resolves to no diff line.
//
// An end covers the whole of its row, so where the row shows more than one diff line
// — a rendering that puts a modification's two halves side by side, or a word diff
// that puts both on the one line it changed — the end takes the outermost of them:
// the last for the range's lower end and the first for its upper one. Otherwise a
// rendering that splits them apart again would get back only the half the row led
// with, and half a change selected where a whole one was.
func (self *DiffLineHelper) selectionFarEndIdentity(view *gocui.View) (types.DiffLineInfo, bool) {
	farEnd, isLowerEnd, ok := selectionFarEndViewLine(view)
	if !ok {
		return types.DiffLineInfo{}, false
	}
	identities, ok := self.diffLineIdentitiesAt(view, farEnd)
	if !ok {
		return types.DiffLineInfo{}, false
	}
	if isLowerEnd {
		return identities[len(identities)-1], true
	}
	return identities[0], true
}

// selectionFarEndViewLine returns the view line of the end of a range or hunk
// selection the cursor isn't on, and whether that is the lower of the two ends. ok
// is false when there is no range at all, only a cursor.
//
// A range whose two ends are on the same view line still has one, and is not the
// same thing as a cursor sitting there: it covers everything that row shows, which
// may be two lines of the diff at once.
func selectionFarEndViewLine(view *gocui.View) (int, bool, bool) {
	if !view.HasRangeSelect() {
		return 0, false, false
	}
	first, last := view.SelectedLineRange()
	if view.SelectedLineIdx() == first {
		return last, true, true
	}
	return first, false, true
}

// findDiffLine returns the view line showing the given diff line in what view is
// displaying now, for placing a remembered line once the re-render is on screen.
func (self *DiffLineHelper) findDiffLine(view *gocui.View, identity types.DiffLineInfo) (int, bool) {
	bufferLine, ok := self.patchLineRows(view.DiffLineContents())[patchLineOf(identity)]
	if !ok {
		return 0, false
	}
	return view.ViewLineForBufferLine(bufferLine)
}

// restoreDiffLinePositionOnRerender arranges for view's next re-render to land on the
// first of the given candidate lines the new rendering still has, calling place with
// that candidate and the view line it ended up on. The candidates are in priority
// order (see nearbyDiffLines); if the rendering has none of them, place isn't called
// and the view re-renders as it otherwise would.
//
// The nearest candidate is looked for as the content loads, so that the re-render can
// be revealed at the right position as soon as that line and a screenful below it
// have arrived. Only the nearest one, because the candidates aren't in load order: a
// farther one can load first, and landing on it while a nearer one is still on its
// way would be settling for worse. The rest are considered together once the whole
// rendering is there.
func (self *DiffLineHelper) restoreDiffLinePositionOnRerender(
	view *gocui.View, candidates []diffLineAnchor, place func(anchor diffLineAnchor, viewLine int),
) {
	if len(candidates) == 0 {
		return
	}

	// The search of the loading content runs on the task's own goroutine, where the
	// repo we are in may not be read — a repo switch replaces it — so take it here, on
	// the UI thread, for the search to work from.
	worktreePath := self.c.Git().RepoPaths.WorktreePath()

	// Which candidate the search settled on, for place to put back where it was.
	found := diffLineAnchor{}

	self.installDiffLineRestore(view,
		func(rows []gocui.DiffLineContent, offset int) (int, bool) {
			for i, row := range rows {
				if rowShowsDiffLine(row, worktreePath, candidates[0].identity) {
					found = candidates[0]
					return offset + i, true
				}
			}
			return 0, false
		},
		func(contents []gocui.DiffLineContent) (int, bool) {
			anchor, bufferLine := self.nearestSurvivingCandidate(contents, candidates)
			if bufferLine == -1 {
				return 0, false
			}
			found = anchor
			return bufferLine, true
		},
		func(viewLine int) { place(found, viewLine) },
		nil,
	)
}

// ChangeLineOrdinal returns how many change lines of view's rendered diff come before
// the one at the given view line — that line's place in the sequence of changes. ok is
// false when the view line belongs to no row of the content.
//
// It is how a place in a diff is remembered across acting on it: an action consumes
// the lines it acted on, so the identity of the line the user was on is gone, but the
// place it left behind is the same one that identity used to have.
func (self *DiffLineHelper) ChangeLineOrdinal(view *gocui.View, viewLine int) (int, bool) {
	bufferLine, ok := view.BufferLineForViewLine(viewLine)
	if !ok {
		return 0, false
	}

	ordinal := 0
	for i, row := range self.resolveDiffLines(view.DiffLineContents()) {
		if i >= bufferLine {
			break
		}
		if row.ok && row.info.IsChange() {
			ordinal++
		}
	}
	return ordinal, true
}

// RevealChangeLineAtOrdinal arranges for view's next re-render to be shown with the
// change line at the given ordinal placed by place — the diff having changed under the
// user, this is where what they were doing carries on. When the new diff has fewer
// changes than that, because the ones acted on were its last, it lands on the last
// change left.
//
// done is called once the selection is where it belongs, or once it turns out that no
// render is coming to put it there, for a caller that must not let the user act again
// in between.
func (self *DiffLineHelper) RevealChangeLineAtOrdinal(
	view *gocui.View, ordinal int, place func(viewLine int), done func(),
) {
	// How many change lines the incremental search has passed, so that it can carry on
	// counting where it left off.
	seen := 0

	self.installDiffLineRestore(view,
		func(rows []gocui.DiffLineContent, offset int) (int, bool) {
			for i, row := range rows {
				if info, ok := self.diffLineInfoFromRecords(row.Metadata); ok && info.IsChange() {
					if seen == ordinal {
						return offset + i, true
					}
					seen++
				}
			}
			return 0, false
		},
		func(contents []gocui.DiffLineContent) (int, bool) {
			last, count := -1, 0
			for i, row := range self.resolveDiffLines(contents) {
				if !row.ok || !row.info.IsChange() {
					continue
				}
				if count == ordinal {
					return i, true
				}
				count++
				last = i
			}
			return last, last != -1
		},
		place,
		done,
	)
}

// installDiffLineRestore is what the restores are built on: it arranges for view's
// next re-render to be revealed with the row a search finds in it placed by place,
// instead of from the top.
//
// The search comes in two halves, because the content arrives a line at a time.
// findEarly is given the rows that have loaded since it last looked, so that the
// re-render can be revealed as soon as the row is there rather than waiting for the
// rest of a long diff; it can only go by what the renderer states about a row, a
// partly-loaded diff being unparseable. findComplete is given the whole rendering at
// the swap, for a target the incremental search couldn't settle on. Either returns the
// buffer line it found, and place is not called at all when neither does.
func (self *DiffLineHelper) installDiffLineRestore(
	view *gocui.View,
	findEarly func(rows []gocui.DiffLineContent, offset int) (int, bool),
	findComplete func(contents []gocui.DiffLineContent) (int, bool),
	place func(viewLine int),
	done func(),
) {
	// Get-or-create, because the pane may not have rendered anything yet: a file whose
	// diff has only just become split has a second pane whose first render is the one
	// this restore is for.
	manager := self.c.GetOrCreateViewBufferManagerForView(view)
	if manager == nil {
		if done != nil {
			done()
		}
		return
	}

	// The readiness check below runs on the task's own goroutine, which may not read
	// the view's dimensions, so take them here, on the UI thread.
	viewHeight := view.InnerHeight()

	// What the search of the loading content has found, and how far it has looked, so
	// that each line is looked at once.
	foundLine := -1
	scanned := 0

	manager.SetRestoreForNextTask(&tasks.RenderRestore{
		FirstPaintReady: func() bool {
			if foundLine == -1 {
				rows := view.OffscreenDiffLineContentsFrom(scanned)
				if bufferLine, ok := findEarly(rows, scanned); ok {
					foundLine = bufferLine
				}
				scanned += len(rows)
				if foundLine == -1 {
					return false
				}
			}
			// Wait for a screenful below the line as well, so that the re-render isn't
			// revealed with it stranded at the bottom of a half-filled view.
			return view.OffscreenLineCount() >= foundLine+viewHeight
		},
		Apply: func(swapIn func()) {
			bufferLine := foundLine
			if bufferLine == -1 {
				if line, ok := findComplete(view.OffscreenDiffLineContents()); ok {
					bufferLine = line
				}
			}

			swapIn()

			if bufferLine == -1 {
				return
			}
			if viewLine, ok := view.ViewLineForBufferLine(bufferLine); ok {
				place(viewLine)
			}
		},
		Done: done,
	})
}

// nearbyDiffLines collects the lines of view's rendered diff as candidates for a
// restore to land on, ordered by proximity to the anchor line — the anchor itself
// first, then outward, preferring at-or-below on ties — each tagged with the screen
// row it is on. A restore lands on the first of them its re-render still has, so this
// order is what makes it land as near as possible to where the user was.
//
// The walk covers the whole diff rather than stopping at the change lines on either
// side of the anchor, which a context-size change always keeps: ignoring whitespace
// keeps nothing in particular, and can take a hunk or a whole file out of the diff,
// leaving the nearest surviving line in a neighbouring file.
func (self *DiffLineHelper) nearbyDiffLines(view *gocui.View, anchorViewLine int) []diffLineAnchor {
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return nil
	}
	resolved := self.resolveDiffLines(view.DiffLineContents())
	if anchor >= len(resolved) {
		return nil
	}
	rows := screenRows(view, len(resolved))

	candidates := make([]diffLineAnchor, 0, len(resolved))
	collect := func(bufferLine int) {
		if line := resolved[bufferLine]; line.ok {
			candidates = append(candidates, diffLineAnchor{identity: line.info, row: rows[bufferLine]})
		}
	}
	collect(anchor)
	for below, above := anchor+1, anchor-1; below < len(resolved) || above >= 0; below, above = below+1, above-1 {
		if below < len(resolved) {
			collect(below)
		}
		if above >= 0 {
			collect(above)
		}
	}
	return candidates
}

// screenRows maps each line of view's content to the screen row it is drawn on. The
// lines above the visible ones get -1 and those below them the view's height, so that
// putting one of them back where it was lands it at the top or bottom edge.
func screenRows(view *gocui.View, bufferLineCount int) []int {
	height := view.InnerHeight()
	originY := view.OriginY()

	rows := make([]int, bufferLineCount)
	for i := range rows {
		rows[i] = -1
	}
	lastVisible := -1
	for y := originY; y < min(originY+height, view.ViewLinesHeight()); y++ {
		bufferLine, ok := view.BufferLineForViewLine(y)
		if !ok || bufferLine >= bufferLineCount {
			continue
		}
		if rows[bufferLine] == -1 {
			rows[bufferLine] = y - originY
		}
		lastVisible = bufferLine
	}
	for i := lastVisible + 1; i < bufferLineCount; i++ {
		rows[i] = height
	}
	return rows
}

// nearestSurvivingCandidate returns the first of the candidates that the given
// rendering still shows, and the line of it that does. The rendering is indexed
// first, rather than searched once per candidate: the candidate list is as long as
// the diff, and so is the rendering.
func (self *DiffLineHelper) nearestSurvivingCandidate(
	contents []gocui.DiffLineContent, candidates []diffLineAnchor,
) (diffLineAnchor, int) {
	rows := self.patchLineRows(contents)
	for _, candidate := range candidates {
		if line, ok := rows[patchLineOf(candidate.identity)]; ok {
			return candidate, line
		}
	}
	return diffLineAnchor{}, -1
}

// patchLineRows indexes a rendering by the diff lines it shows: for each of them, the
// first of its rows that does. A row can show more than one, and each is then a way
// of finding that row again.
func (self *DiffLineHelper) patchLineRows(contents []gocui.DiffLineContent) map[patchLine]int {
	rows := map[patchLine]int{}
	for i, identities := range self.resolveDiffLineIdentities(contents) {
		for _, identity := range identities {
			if _, seen := rows[patchLineOf(identity)]; !seen {
				rows[patchLineOf(identity)] = i
			}
		}
	}
	return rows
}

// rowShowsDiffLine reports whether the given row of a rendering shows the given diff
// line — among any others it shows, since a side-by-side rendering puts a deletion
// beside the addition replacing it. It only knows what the renderer states about the
// row, since the alternative, parsing the rendering as a diff, needs whole hunks and
// this is asked of content that is still loading. It takes the repo's worktree path
// rather than reading it, being asked off the UI thread.
func rowShowsDiffLine(row gocui.DiffLineContent, worktreePath string, target types.DiffLineInfo) bool {
	return lo.SomeBy(row.Metadata, func(record string) bool {
		parsed, ok := parseDiffLineMetadata(record)
		return ok && patchLineOf(diffLineInfoIn(worktreePath, parsed)) == patchLineOf(target)
	})
}

// patchLine is what stays the same about a diff line when the same diff is rendered
// again differently: which file it belongs to, the line number that identifies it on
// the side it belongs to, and what kind of line it is.
type patchLine struct {
	path string
	// Every kind of content line collapses into DiffLineContext, since an addition
	// and the context line it turns into when whitespace stops counting are the same
	// line of the same file. The header rows keep their kind: a file's header and the
	// first line of the file it heads are not the same place.
	kind types.DiffLineType
	// The old file's line number for a deletion, since two consecutive deletions
	// share a new-file position and differ only here; the new file's otherwise.
	line       int
	isDeletion bool
}

func patchLineOf(info types.DiffLineInfo) patchLine {
	switch info.Type {
	case types.DiffLineFileHeader, types.DiffLineHunkHeader:
		return patchLine{path: info.Path, kind: info.Type, line: info.NewLine}
	case types.DiffLineDeleted:
		return patchLine{path: info.Path, kind: types.DiffLineContext, line: info.OldLine, isDeletion: true}
	default:
		return patchLine{path: info.Path, kind: types.DiffLineContext, line: info.NewLine}
	}
}

// RevealSelectionAfterAction moves a diff pane's selection to the change that takes the
// place of the one just acted on, once the changed diff has re-rendered. Call it with
// the pane acted in, the pane the work carries on in, and the first line of the
// selection, before triggering the re-render.
//
// The line acted on is gone from the diff, so what is remembered is its place among the
// diff's changes: the next change moves up into it, which is where you want to be to
// carry on. A range collapses to a single line at its start, and hunk mode selects the
// whole block it lands in, so that pressing the key again acts on the next hunk. The
// target pane inherits that select mode, this being the same piece of work continuing
// in another pane — and shows no selection until the restore places one, so that what
// it was left showing the last time it was used doesn't appear for a frame.
//
// advanceBy moves on by that many changes past the place remembered, for an action that
// leaves the diff as it was: lines taken into a custom patch are still in the commit's
// diff, so the place remembered is still the line acted on, and carrying on means going
// past the lines just dealt with rather than staying on them.
//
// done, which may be nil, is called once the selection is where it belongs, or once it
// turns out that no render is coming to put it there — for a caller that must not let
// the user act again in between.
func (self *DiffLineHelper) RevealSelectionAfterAction(
	source types.DiffPaneContext, target types.DiffPaneContext, firstLineIdx int, advanceBy int, done func(),
) {
	ordinal, ok := self.ChangeLineOrdinal(source.GetView(), firstLineIdx)
	if !ok {
		if done != nil {
			done()
		}
		return
	}

	sel := source.DiffSelectState()
	if sel.Mode == types.DiffSelectModeRange {
		sel.Mode = types.DiffSelectModeLine
		sel.RangeIsSticky = false
	}
	*target.DiffSelectState() = *sel
	selectHunk := sel.Mode == types.DiffSelectModeHunk

	targetView := target.GetView()
	if target != source {
		target.SetHasSelectableContent(false)
		self.c.Context().UpdateSelectionHighlights()
	}

	self.RevealChangeLineAtOrdinal(targetView, ordinal+advanceBy, func(viewLine int) {
		if selectHunk {
			self.SelectChangeBlock(target, viewLine, true)
			return
		}
		targetView.CancelRangeSelect()
		self.ShowSelectionAtLine(targetView, viewLine, true)
	}, done)
}
