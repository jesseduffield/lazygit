package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type MainViewController struct {
	baseController
	c *ControllerCommon

	context      *context.MainContext
	otherContext *context.MainContext
}

var _ types.IController = &MainViewController{}

func NewMainViewController(
	c *ControllerCommon,
	context *context.MainContext,
	otherContext *context.MainContext,
) *MainViewController {
	return &MainViewController{
		baseController: baseController{},
		c:              c,
		context:        context,
		otherContext:   otherContext,
	}
}

func (self *MainViewController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Keys:            opts.GetKeys(opts.Config.Universal.TogglePanel),
			Handler:         self.togglePanel,
			Description:     self.c.Tr.ToggleStagingView,
			Tooltip:         self.c.Tr.ToggleStagingViewTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:    opts.GetKeys(opts.Config.Main.ToggleSelectHunk),
			Handler: self.toggleSelectHunk,
			DescriptionFunc: self.diffSelectionDescription(func() string {
				if self.diffSelectState().Mode == types.DiffSelectModeHunk {
					return self.c.Tr.SelectLineByLine
				}
				return self.c.Tr.SelectHunk
			}),
			Description:       self.c.Tr.ToggleSelectHunk,
			GetDisabledReason: self.diffSelectionDisabledReason,
			Tooltip:           self.c.Tr.ToggleSelectHunkTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.ToggleRangeSelect),
			Handler:           self.toggleRangeSelect,
			Description:       self.c.Tr.ToggleRangeSelect,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.ToggleRangeSelect),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Return),
			Handler:         self.escape,
			Description:     self.c.Tr.ExitFocusedMainView,
			DescriptionFunc: self.escapeDescription,
			DisplayOnScreen: true,
		},
		{
			// overriding this because we want to read all of the task's output before we start searching
			Keys:        opts.GetKeys(opts.Config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.StartSearch,
			Tag:         "navigation",
		},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.PrevItem), Handler: self.handlePrevLine},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.NextItem), Handler: self.handleNextLine},
		{
			Tag:               "navigation",
			Keys:              opts.GetKeys(opts.Config.Universal.RangeSelectUp),
			Handler:           self.extendRangeUp,
			Description:       self.c.Tr.RangeSelectUp,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.RangeSelectUp),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{
			Tag:               "navigation",
			Keys:              opts.GetKeys(opts.Config.Universal.RangeSelectDown),
			Handler:           self.extendRangeDown,
			Description:       self.c.Tr.RangeSelectDown,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.RangeSelectDown),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.PrevPage), Handler: self.handlePrevPage, Description: self.c.Tr.PrevPage},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.NextPage), Handler: self.handleNextPage, Description: self.c.Tr.NextPage},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.GotoTop), Handler: self.handleGotoTop, Description: self.c.Tr.GotoTop},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.GotoBottom), Handler: self.handleGotoBottom, Description: self.c.Tr.GotoBottom},
	}
}

func (self *MainViewController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInAlreadyFocusedView,
			FocusedView: self.context.GetViewName(),
		},
		{
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler:     self.onClickInOtherViewOfMainViewPair,
			FocusedView: self.otherContext.GetViewName(),
		},
	}
}

func (self *MainViewController) Context() types.Context {
	return self.context
}

func (self *MainViewController) togglePanel() error {
	if !self.otherContext.GetView().Visible {
		return nil
	}

	// Whether the pair holds a diff is decided by the side panel beneath, which
	// NextInStack only finds while our context is still the focused main view, so
	// read it before pushing the other pane.
	isDiff := self.isDiffView()
	self.c.Context().Push(self.otherContext, types.OnFocusOpts{})
	if isDiff {
		establishDiffSelection(self.c, self.otherContext, -1)
	}
	return nil
}

// escape dismisses the selection a step at a time before leaving the view: a range
// collapses to its cursor line, and hunk mode the user turned on goes back to
// line-by-line. Hunk mode that is merely the configured default is not something to
// escape from, so there escape leaves.
func (self *MainViewController) escape() error {
	if self.selectingRange() || self.selectingHunkEnabledByUser() {
		self.context.ResetDiffSelectMode()
		return nil
	}

	self.c.Context().Pop()
	return nil
}

func (self *MainViewController) escapeDescription() string {
	if self.selectingRange() {
		return self.c.Tr.DismissRangeSelect
	}
	if self.selectingHunkEnabledByUser() {
		return self.c.Tr.SelectLineByLine
	}
	return self.c.Tr.ExitFocusedMainView
}

// selectingHunkEnabledByUser reports whether we are in hunk mode because the user
// asked for it, as opposed to it being the configured default.
func (self *MainViewController) selectingHunkEnabledByUser() bool {
	return self.diffSelectState().Mode == types.DiffSelectModeHunk && self.diffSelectState().UserEnabledHunkMode
}

// isDiffView reports whether the focused main view currently shows a diff, and so
// shows a selection. See types.DiffMainViewContext.
func (self *MainViewController) isDiffView() bool {
	return self.diffMainViewType() != types.DiffMainViewTypeNone
}

// diffMainViewType reports what the diff in the focused main view belongs to, taken
// from the side panel beneath it, or DiffMainViewTypeNone when this pane isn't on the
// stack or has no diff panel beneath it. The IsInStack guard is essential:
// NextInStack panics for a context that isn't in the stack, and GetKeybindings (which
// leads here) also runs for off-stack panes — at startup and while generating the
// cheatsheets, where the stack is empty.
func (self *MainViewController) diffMainViewType() types.DiffMainViewType {
	if !self.c.Context().IsInStack(self.context) {
		return types.DiffMainViewTypeNone
	}
	if diffContext, ok := self.c.Context().NextInStack(self.context).(types.DiffMainViewContext); ok {
		return diffContext.GetDiffMainViewType()
	}
	return types.DiffMainViewTypeNone
}

// diffSelectState returns this pane's diff selection mode state.
func (self *MainViewController) diffSelectState() *types.DiffSelectState {
	return self.context.DiffSelectState()
}

// diffSelectionDescription qualifies the description of a command that acts on the
// selection, so that it is listed only where it applies: the main view also shows
// content with nothing to select in it — a branch's commit log, the status dashboard —
// and a command with no description is left out of the keybindings menu.
//
// The static Description stays as it is: the cheatsheets are generated from that, and
// they document what a key does rather than when it applies.
func (self *MainViewController) diffSelectionDescription(describe func() string) func() string {
	return func() string {
		if !self.isDiffView() {
			return ""
		}
		return describe()
	}
}

func (self *MainViewController) diffSelectionDescriptionText(description string) func() string {
	return self.diffSelectionDescription(func() string { return description })
}

// diffSelectionDisabledReason disables the commands that act on the selection while
// there is none to act on: a diff view whose diff holds nothing selectable (a binary
// file, an empty commit) or which is showing a placeholder message.
func (self *MainViewController) diffSelectionDisabledReason() *types.DisabledReason {
	if !self.context.GetView().Highlight {
		return &types.DisabledReason{Text: self.c.Tr.NothingToSelectInDiff}
	}
	return nil
}

func (self *MainViewController) onClickInAlreadyFocusedView(opts gocui.ViewMouseBindingOpts) error {
	self.selectClickedDiffLine(opts.Y)
	return nil
}

func (self *MainViewController) onClickInOtherViewOfMainViewPair(opts gocui.ViewMouseBindingOpts) error {
	// Carry the select mode over from the pane we're leaving, so that clicking into
	// the other pane keeps hunk mode even the first time we enter it — its own mode
	// would otherwise still be the default single line until it had been focused at
	// least once. selectClickedDiffLine then keeps or collapses that mode depending on
	// where the click landed.
	*self.context.DiffSelectState() = *self.otherContext.DiffSelectState()
	self.c.Context().Push(self.context, types.OnFocusOpts{})
	self.selectClickedDiffLine(opts.Y)
	return nil
}

// selectClickedDiffLine sets the focused main view's selection from a click at the
// given view line. In hunk mode, clicking inside the selected block collapses it to
// that line; clicking a change line outside it keeps hunk mode and selects that block.
// A click on context, or any click outside hunk mode, selects just that line too.
func (self *MainViewController) selectClickedDiffLine(viewLine int) {
	if !self.isDiffView() {
		return
	}
	view := self.context.GetView()
	if self.diffSelectState().Mode == types.DiffSelectModeHunk {
		if start, end, ok := self.c.Helpers().DiffLine.SelectedHunkBounds(view); ok &&
			viewLine >= start && viewLine <= end {
			self.context.ResetDiffSelectMode()
			showSelectionAtLine(view, viewLine, false)
			return
		}
		if self.c.Helpers().DiffLine.IsChangeLine(view, viewLine) {
			self.selectHunkAround(viewLine, false)
			return
		}
	}
	self.context.ResetDiffSelectMode()
	showSelectionAtLine(view, viewLine, false)
}

// establishDiffSelection turns on the focused main view's selection once the view has
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
func establishDiffSelection(c *ControllerCommon, mainContext *context.MainContext, clickedViewLine int) {
	mainContext.ResetDiffSelectMode()
	view := mainContext.GetView()

	// The panel beneath renders a diff, but that diff may hold nothing to act on: a
	// binary file, or an empty commit. Rendering it worked that out, so the pane is
	// already showing no selection and there is nowhere to put one.
	if !c.Helpers().DiffLine.ViewHasChangeLines(view) {
		return
	}

	if clickedViewLine >= 0 {
		if hunkModeApplies(c, view, clickedViewLine) &&
			c.Helpers().DiffLine.IsChangeLine(view, clickedViewLine) {
			mainContext.DiffSelectState().Mode = types.DiffSelectModeHunk
			selectDiffHunk(c, mainContext, clickedViewLine, false)
			return
		}
		showSelectionAtLine(view, clickedViewLine, false)
		return
	}

	target, ok := changeToSelectOnScreen(c, view)
	if !ok {
		showSelectionAtLine(view, view.MiddleVisibleLineIdx(), false)
		return
	}
	if hunkModeApplies(c, view, target) {
		mainContext.DiffSelectState().Mode = types.DiffSelectModeHunk
		selectDiffHunk(c, mainContext, target, false)
		return
	}
	showSelectionAtLine(view, target, false)
}

// changeToSelectOnScreen returns the change line keyboard focus establishes the
// selection on. In hunk mode that is the first block that begins on screen, so that
// the block being offered up is one the user can see the extent of, falling back to a
// block that reaches into the view from above — a change longer than the screen, where
// there is nothing else to offer. Line by line it is simply the first change line on
// screen. ok is false when the viewport shows no change at all.
func changeToSelectOnScreen(c *ControllerCommon, view *gocui.View) (int, bool) {
	if c.UserConfig().Gui.UseHunkModeInStagingView {
		return c.Helpers().DiffLine.FirstChangeBlockInView(view)
	}
	return c.Helpers().DiffLine.FirstChangeLineInView(view)
}

// hunkModeApplies reports whether an established selection should start out as the
// whole change block around the given change line. That's what the config asks for,
// except over a file shown as one solid block of changes, where it would select the
// whole file — see DiffLineHelper.IsSingleHunkForWholeFile.
func hunkModeApplies(c *ControllerCommon, view *gocui.View, changeViewLine int) bool {
	return c.UserConfig().Gui.UseHunkModeInStagingView &&
		!c.Helpers().DiffLine.IsSingleHunkForWholeFile(view, changeViewLine)
}

// showSelectionAtLine moves the focused main view's selection to the given view line,
// clamped to the content. scrollIntoView scrolls the line into view when it's
// off-screen, for navigating to it; a click leaves it false, the clicked line being on
// screen already.
func showSelectionAtLine(view *gocui.View, lineIdx int, scrollIntoView bool) {
	view.FocusPoint(0, lo.Clamp(lineIdx, 0, max(0, view.ViewLinesHeight()-1)), scrollIntoView)
}

func (self *MainViewController) selectHunkAround(changeViewLine int, scrollIntoView bool) {
	selectDiffHunk(self.c, self.context, changeViewLine, scrollIntoView)
}

// selectDiffHunk selects the whole change block around the given change line, for
// hunk mode: the cursor goes to the block's first line and the range anchor to its
// last, so the native range highlight spans the block. With no block to be found —
// a diff with no changes in it — it falls back to a single-line selection.
//
// scrollIntoView brings the block's first line on screen, for the commands that mean
// to go there; a click leaves it false, so that the view doesn't move under the mouse
// when the block the click landed in starts above the viewport.
func selectDiffHunk(
	c *ControllerCommon, mainContext *context.MainContext, changeViewLine int, scrollIntoView bool,
) {
	view := mainContext.GetView()
	start, end, ok := c.Helpers().DiffLine.ChangeBlockBounds(view, changeViewLine)
	if !ok {
		mainContext.DiffSelectState().Mode = types.DiffSelectModeLine
		view.CancelRangeSelect()
		showSelectionAtLine(view, changeViewLine, scrollIntoView)
		return
	}
	view.SetRangeSelectStart(end)
	showSelectionAtLine(view, start, scrollIntoView)
}

// navigate moves the focused main view to the row find locates from the current
// anchor — the selected line when a selection is showing, otherwise the top visible
// line. With a selection we move it there and scroll it into view, re-selecting the
// whole block in hunk mode; with none we stay in scroll mode, bringing the target
// to the top without selecting anything.
func (self *MainViewController) navigate(find findDiffRowFn, forward bool) {
	v := self.context.GetView()
	anchor := v.OriginY()
	if v.Highlight {
		anchor = v.SelectedLineIdx()
	}

	if target, ok := find(v, anchor, forward); ok {
		self.placeNavigationTarget(target)
		return
	}
	if !forward {
		// Everything above the anchor has loaded, so a backward target that wasn't
		// found doesn't exist.
		return
	}

	// The diff loads lazily, so a target below the loaded portion isn't there to be
	// found yet. Read the rest of it in and look again before concluding there is none.
	manager := self.c.GetViewBufferManagerForView(v)
	if manager == nil {
		return
	}
	manager.ReadToEnd(func() {
		self.c.OnUIThread(func() error {
			if target, ok := find(v, anchor, forward); ok {
				self.placeNavigationTarget(target)
			}
			return nil
		})
	})
}

// findDiffRowFn locates a row of the rendered diff to navigate to, given the view,
// the anchor view line to start from, and the direction.
type findDiffRowFn func(view *gocui.View, anchorViewLine int, forward bool) (int, bool)

func (self *MainViewController) placeNavigationTarget(target int) {
	v := self.context.GetView()
	if !v.Highlight {
		v.SetOrigin(0, target)
		return
	}
	if self.diffSelectState().Mode == types.DiffSelectModeHunk {
		self.selectHunkAround(target, true)
		return
	}
	// Line mode leaves a single-line selection at the target; an active range extends
	// to it, the anchor being untouched.
	showSelectionAtLine(v, target, true)
}

// moveCursor moves the selection cursor by delta view lines (negative = up), with the
// configured scroll-off margin, reading more content in first when moving down. The
// range anchor is left untouched, so this extends or contracts a range and just moves
// the selected line otherwise.
func (self *MainViewController) moveCursor(delta int) {
	v := self.context.GetView()
	if delta > 0 {
		self.c.ReadLinesToFillView(v)
	}
	before := v.SelectedLineIdx()
	after := lo.Clamp(before+delta, 0, v.ViewLinesHeight()-1)
	if delta == -1 {
		checkScrollUp(self.context.GetViewTrait(), self.c.UserConfig(), before, after)
	} else if delta == 1 {
		checkScrollDown(self.context.GetViewTrait(), self.c.UserConfig(), before, after)
	}
	showSelectionAtLine(v, after, true)
}

// collapseForLineMove drops hunk mode, and a non-sticky range, back to a single-line
// selection — what a plain (non-shift, non-hunk-step) move does before moving. A
// sticky range is kept, so the move extends it.
func (self *MainViewController) collapseForLineMove() {
	sel := self.diffSelectState()
	if sel.Mode == types.DiffSelectModeHunk ||
		(sel.Mode == types.DiffSelectModeRange && !sel.RangeIsSticky) {
		sel.Mode = types.DiffSelectModeLine
		self.context.GetView().CancelRangeSelect()
	}
}

// adjustSelection moves the selection by delta view lines, for the plain up/down and
// page keys. In hunk mode a single-line step jumps to the adjacent block, while a
// larger page step drops out of hunk mode first. A non-sticky range collapses back to
// a single line on a plain move. With no selection — non-diff content — it scrolls.
func (self *MainViewController) adjustSelection(delta int) {
	if !self.context.GetView().Highlight {
		self.handleLineChange(delta)
		return
	}
	if self.diffSelectState().Mode == types.DiffSelectModeHunk && (delta == 1 || delta == -1) {
		self.navigate(self.c.Helpers().DiffLine.AdjacentChangeBlock, delta > 0)
		return
	}
	self.collapseForLineMove()
	self.moveCursor(delta)
}

// selectAbsoluteLine moves the selection to a specific view line — the top or bottom
// of the diff — dropping hunk mode and a non-sticky range like a plain move does.
func (self *MainViewController) selectAbsoluteLine(target int) {
	self.collapseForLineMove()
	showSelectionAtLine(self.context.GetView(), target, true)
}

// selectingRange reports whether a range selection is currently active: we're in
// range mode and either it's sticky or the anchor and cursor differ, i.e. a
// non-sticky range that has actually been extended.
func (self *MainViewController) selectingRange() bool {
	if self.diffSelectState().Mode != types.DiffSelectModeRange {
		return false
	}
	start, end := self.context.GetView().SelectedLineRange()
	return self.diffSelectState().RangeIsSticky || start != end
}

// toggleSelectHunk switches between selecting the change block around the cursor and
// a single line.
func (self *MainViewController) toggleSelectHunk() error {
	v := self.context.GetView()
	if !v.Highlight {
		return nil
	}
	sel := self.diffSelectState()
	if sel.Mode == types.DiffSelectModeHunk {
		sel.Mode = types.DiffSelectModeLine
		v.CancelRangeSelect()
	} else {
		sel.Mode = types.DiffSelectModeHunk
		sel.UserEnabledHunkMode = true
		self.selectHunkAround(v.SelectedLineIdx(), true)
	}
	return nil
}

// toggleRangeSelect starts or cancels a sticky range selection, which the plain
// up/down keys extend.
func (self *MainViewController) toggleRangeSelect() error {
	v := self.context.GetView()
	if !v.Highlight {
		return nil
	}
	sel := self.diffSelectState()
	if self.selectingRange() {
		sel.Mode = types.DiffSelectModeLine
		sel.RangeIsSticky = false
		v.CancelRangeSelect()
	} else {
		sel.Mode = types.DiffSelectModeRange
		sel.RangeIsSticky = true
		v.SetRangeSelectStart(v.SelectedLineIdx())
	}
	return nil
}

// extendRange grows a non-sticky range selection by one line in response to
// shift+up/down, starting one at the cursor if there isn't one yet.
func (self *MainViewController) extendRange(forward bool) error {
	v := self.context.GetView()
	if !v.Highlight {
		return nil
	}
	sel := self.diffSelectState()
	if !self.selectingRange() {
		sel.Mode = types.DiffSelectModeRange
		v.SetRangeSelectStart(v.SelectedLineIdx())
	}
	sel.RangeIsSticky = false
	if forward {
		self.moveCursor(1)
	} else {
		self.moveCursor(-1)
	}
	return nil
}

func (self *MainViewController) extendRangeUp() error {
	return self.extendRange(false)
}

func (self *MainViewController) extendRangeDown() error {
	return self.extendRange(true)
}

func (self *MainViewController) handleLineChange(delta int) {
	v := self.context.GetView()
	if delta < 0 {
		v.ScrollUp(-delta)
	} else {
		v.ScrollDown(delta)
		self.c.ReadLinesToFillView(v)
	}
}

func (self *MainViewController) handlePrevLine() error {
	self.adjustSelection(-1)
	return nil
}

func (self *MainViewController) handleNextLine() error {
	self.adjustSelection(1)
	return nil
}

func (self *MainViewController) handlePrevPage() error {
	self.adjustSelection(-self.context.GetViewTrait().PageDelta())
	return nil
}

func (self *MainViewController) handleNextPage() error {
	self.adjustSelection(self.context.GetViewTrait().PageDelta())
	return nil
}

func (self *MainViewController) handleGotoTop() error {
	v := self.context.GetView()
	if !v.Highlight {
		self.handleLineChange(-v.ViewLinesHeight())
		return nil
	}
	self.selectAbsoluteLine(0)
	return nil
}

func (self *MainViewController) handleGotoBottom() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				v := self.context.GetView()
				if !v.Highlight {
					self.handleLineChange(v.ViewLinesHeight())
					return nil
				}
				self.selectAbsoluteLine(v.ViewLinesHeight() - 1)
				return nil
			})
		})
	}

	return nil
}

func (self *MainViewController) openSearch() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				self.c.Helpers().Search.OpenSearchPrompt(self.context)
				return nil
			})
		})
	}

	return nil
}
