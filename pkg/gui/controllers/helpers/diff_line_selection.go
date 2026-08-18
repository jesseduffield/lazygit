package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// Putting a selection in the focused main view: where it starts out, and how it is
// widened to a whole change block. Both are answered from what the view is showing,
// as recovered by the queries next door.

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

	// The panel beneath renders a diff, but that diff may hold nothing to act on: a
	// binary file, or an empty commit. Rendering it worked that out, so the pane is
	// already showing no selection and there is nowhere to put one.
	if !self.ViewHasChangeLines(view) {
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

// changeToSelectOnScreen returns the change line keyboard focus establishes the
// selection on. In hunk mode that is the first block that begins on screen, so that
// the block being offered up is one the user can see the extent of, falling back to a
// block that reaches into the view from above — a change longer than the screen, where
// there is nothing else to offer. Line by line it is simply the first change line on
// screen. ok is false when the viewport shows no change at all.
func (self *DiffLineHelper) changeToSelectOnScreen(view *gocui.View) (int, bool) {
	if self.c.UserConfig().Gui.UseHunkModeInStagingView {
		return self.FirstChangeBlockInView(view)
	}
	return self.FirstChangeLineInView(view)
}

// hunkModeApplies reports whether an established selection should start out as the
// whole change block around the given change line. That's what the config asks for,
// except over a file shown as one solid block of changes, where it would select the
// whole file — see IsSingleHunkForWholeFile.
func (self *DiffLineHelper) hunkModeApplies(view *gocui.View, changeViewLine int) bool {
	return self.c.UserConfig().Gui.UseHunkModeInStagingView &&
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
