package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// Putting a selection in the focused main view: where it starts out, where a jump
// leaves it, and how it is widened to a whole change block. All three are answered
// from what the view is showing, as recovered by the queries next door.

// EstablishSelection turns on the focused main view's selection once the view has
// been focused. clickedViewLine is the view line a click pointed at, or -1 for
// keyboard focus, which points at no particular line and so starts at the first
// change line on screen.
//
// Focusing never moves the view: you focus the diff you are reading in order to point
// at something in it, so the selection goes where you are looking rather than the
// view going where the selection would like to be. With no change line on screen at
// all — a long stretch of context — it lands on the middle visible line, the likeliest
// one to be the one being read.
//
// With hunk mode configured as the default the selection widens to the whole change
// block: keyboard focus lands on the first block on screen, and a click on a change
// line selects that line's block, ready to act on. A click on context still selects
// just that line — the click points at it precisely, so it stays editable.
func (self *DiffLineHelper) EstablishSelection(mainContext *context.MainContext, clickedViewLine int) {
	mainContext.ResetDiffSelectMode()
	view := mainContext.GetView()

	// The pane may hold nothing to act on: a diff of a binary file or an empty commit,
	// or content that is no diff of the panel's at all. Rendering it worked that out,
	// so the pane is already showing no selection and there is nowhere to put one.
	if !mainContext.HasSelectableContent() {
		return
	}

	if clickedViewLine >= 0 {
		// Remember where the click landed so that a drag that follows anchors its range
		// there, even when this click selects a whole hunk.
		mainContext.SetDragAnchorViewLine(clickedViewLine)
		if self.hunkModeApplies(view, clickedViewLine) && self.IsChangeLine(view, clickedViewLine) {
			mainContext.DiffSelectState().Mode = types.DiffSelectModeHunk
			self.SelectChangeBlock(mainContext, clickedViewLine, false)
			return
		}
		self.ShowSelectionAtLine(view, clickedViewLine, false)
		return
	}

	target, ok := self.changeToSelectOnScreen(view)
	if !ok {
		self.ShowSelectionAtLine(view, view.MiddleVisibleLineIdx(), false)
		return
	}
	if self.hunkModeApplies(view, target) {
		mainContext.DiffSelectState().Mode = types.DiffSelectModeHunk
		self.SelectChangeBlock(mainContext, target, false)
		return
	}
	self.ShowSelectionAtLine(view, target, false)
}

// PlaceNavigationTarget moves the pane's selection to the row a jump found, bringing
// it on screen if it isn't already. With no selection to move — a pane that isn't
// focused, or one showing a diff with nothing selectable in it — the row goes to the
// top of the view instead, that being all a jump can do there.
//
// alignTop asks for the target to become the view's top line, so that everything that
// begins there is on screen. It only applies to a target the view has to scroll to: a
// jump to something already on screen leaves the view alone, there being nothing to
// gain from moving what the user is looking at. In hunk mode what ends up selected is
// the first change block at or below the target, which a large context size can put
// further down than a screenful; the selection is then scrolled into view as any other
// jump's is, and the alignment gives way to that.
func (self *DiffLineHelper) PlaceNavigationTarget(
	pane types.DiffPaneContext, target int, alignTop bool,
) {
	view := pane.GetView()
	if !view.Highlight {
		view.SetOrigin(0, target)
		return
	}
	if alignTop {
		self.scrollTargetToTop(pane, target)
	}
	// Jumping to another block or file moves the cursor without shift held, so a
	// range that grows only while shift is held collapses rather than stretching all
	// the way to the target. A sticky range stretches instead; this is the point of
	// being sticky.
	self.CollapseNonStickyRange(pane)
	if pane.DiffSelectState().Mode == types.DiffSelectModeHunk {
		self.SelectChangeBlock(pane, target, true)
		return
	}
	// Line mode leaves a single-line selection at the target; an active range extends
	// to it, the anchor being untouched.
	self.ShowSelectionAtLine(view, target, true)
}

// scrollTargetToTop scrolls the given row of the diff to the top of the view, leaving
// the view where it is when that row is on screen already. The last screenful of the
// diff is as far as it goes, so that the view doesn't scroll past the end of what it is
// showing.
func (self *DiffLineHelper) scrollTargetToTop(pane types.DiffPaneContext, target int) {
	view := pane.GetView()
	originY, height := pane.GetViewTrait().ViewPortYBounds()
	if target >= originY && target < originY+height {
		return
	}
	view.SetOriginY(min(target, max(0, view.ViewLinesHeight()-height)))
}

// CollapseNonStickyRange drops a range that only grows while shift is held back to a
// single line at the cursor.
func (self *DiffLineHelper) CollapseNonStickyRange(pane types.DiffPaneContext) {
	sel := pane.DiffSelectState()
	if sel.Mode == types.DiffSelectModeRange && !sel.RangeIsSticky {
		sel.Mode = types.DiffSelectModeLine
		pane.GetView().CancelRangeSelect()
	}
}

// changeToSelectOnScreen returns the change line keyboard focus establishes the
// selection on. In hunk mode that is the first block that begins on screen, so that
// the block being offered up is one the user can see the extent of, falling back to a
// block that reaches into the view from above — a change longer than the screen, where
// there is nothing else to offer. Line by line it is simply the first change line on
// screen. ok is false when the viewport shows no change at all.
func (self *DiffLineHelper) changeToSelectOnScreen(view *gocui.View) (int, bool) {
	if self.c.UserConfig().Gui.UseHunkModeInDiffView {
		return self.FirstChangeBlockInView(view)
	}
	return self.FirstChangeLineInView(view)
}

// hunkModeApplies reports whether an established selection should start out as the
// whole change block around the given change line. That's what the config asks for,
// except over a file shown as one solid block of changes, where it would select the
// whole file — see IsSingleHunkForWholeFile.
func (self *DiffLineHelper) hunkModeApplies(view *gocui.View, changeViewLine int) bool {
	return self.c.UserConfig().Gui.UseHunkModeInDiffView &&
		!self.IsSingleHunkForWholeFile(view, changeViewLine)
}

// ShowSelectionAtLine moves the focused main view's selection to the given view line,
// clamped to the content. scrollIntoView scrolls the line into view when it's
// off-screen, for navigating to it; a click leaves it false, the clicked line being on
// screen already.
func (self *DiffLineHelper) ShowSelectionAtLine(view *gocui.View, lineIdx int, scrollIntoView bool) {
	view.FocusPoint(0, lo.Clamp(lineIdx, 0, max(0, view.ViewLinesHeight()-1)), scrollIntoView)

	// A search carries on from where the selection now is, so that stepping to the
	// next match goes to the one after it rather than the one after the match the
	// user last stepped to.
	view.SetNearestSearchPosition()
}

// SelectChangeBlock selects the whole change block around the given change line, for
// hunk mode: the cursor goes to the block's first line and the range anchor to its
// last, so the native range highlight spans the block. With no block to be found —
// a diff with no changes in it — it falls back to a single-line selection.
//
// scrollIntoView brings the block's first line on screen, for the commands that mean
// to go there; a click leaves it false, so that the view doesn't move under the mouse
// when the block the click landed in starts above the viewport.
func (self *DiffLineHelper) SelectChangeBlock(
	pane types.DiffPaneContext, changeViewLine int, scrollIntoView bool,
) {
	view := pane.GetView()
	start, end, ok := self.ChangeBlockBounds(view, changeViewLine)
	if !ok {
		pane.DiffSelectState().Mode = types.DiffSelectModeLine
		view.CancelRangeSelect()
		self.ShowSelectionAtLine(view, changeViewLine, scrollIntoView)
		return
	}
	view.SetRangeSelectStart(end)
	self.ShowSelectionAtLine(view, start, scrollIntoView)
}

// SelectedHunkBounds returns the change block selected in hunk mode. The range
// anchor stays on the block's far end when a click moves the cursor before its
// handler runs, so it still identifies the selected block.
func (self *DiffLineHelper) SelectedHunkBounds(view *gocui.View) (int, int, bool) {
	anchor := view.RangeSelectStartY()
	if anchor < 0 {
		return 0, 0, false
	}
	return self.ChangeBlockBounds(view, anchor)
}

// RefreshInclusionGutter updates the marks drawn over the diff in the main pane, which
// say which of its lines are in the custom patch being built from it.
//
// They are shown while the focused main view holds the focus — either of its panes, so
// that moving between the diff and the patch previewed beside it doesn't make them come
// and go — and only over a diff a patch is being built from: a patch built from some
// other commit says nothing about the lines of this one.
//
// Call it whenever either of those can have changed: as a pane's content settles, when
// the focus arrives or leaves, and when the patch itself changes.
func (self *DiffLineHelper) RefreshInclusionGutter() {
	view := self.c.Contexts().Normal.GetView()

	included := self.patchInclusion()
	if included == nil {
		view.SetInclusionGutter(false, nil)
		return
	}

	resolved := self.resolveDiffLines(view.DiffLineContents())
	marks := make([]bool, len(resolved))
	showsChanges := false
	for i, row := range resolved {
		if !row.ok || !row.info.IsChange() {
			continue
		}
		showsChanges = true
		marks[i] = included(row.info)
	}

	// Nothing to mark and nowhere to mark it: the pane is showing a message rather than
	// a diff, or a diff with nothing in it.
	if !showsChanges {
		view.SetInclusionGutter(false, nil)
		return
	}
	view.SetInclusionGutter(true, marks)
}

// patchInclusion asks the panel whose diff the focused main view is showing which of
// that diff's lines are in the custom patch being built from it, and answers nil where
// there is no such patch — including when the focus is elsewhere, the marks being an
// affordance of the focused view.
func (self *DiffLineHelper) patchInclusion() func(types.DiffLineInfo) bool {
	if !self.mainViewIsFocused() {
		return nil
	}
	// The panel beneath is found from the pane that holds the focus, which is not always
	// the one the diff is in: moving to the pane beside it takes the other off the stack.
	sidePanel := self.c.Context().NextInStack(self.c.Context().CurrentStatic())
	if sidePanel == nil {
		return nil
	}
	actions, ok := sidePanel.GetFocusedMainViewDiffSource().(types.FocusedMainViewActions)
	if !ok {
		return nil
	}
	return actions.PatchInclusion()
}

// ShowsCustomPatch reports whether the given view is the one previewing the custom patch
// being built, which is the lower pane while a patch is being built from the diff in the
// upper one.
func (self *DiffLineHelper) ShowsCustomPatch(view *gocui.View) bool {
	return view == self.c.Contexts().NormalSecondary.GetView() && self.patchInclusion() != nil
}
