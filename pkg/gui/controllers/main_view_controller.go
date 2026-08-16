package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type MainViewController struct {
	baseController
	c *ControllerCommon

	context      *context.MainContext
	otherContext *context.MainContext

	dragAutoscroller  *helpers.DragAutoscroller
	draggingWithMouse bool
}

var _ types.IController = &MainViewController{}

func NewMainViewController(
	c *ControllerCommon,
	context *context.MainContext,
	otherContext *context.MainContext,
) *MainViewController {
	controller := &MainViewController{
		baseController: baseController{},
		c:              c,
		context:        context,
		otherContext:   otherContext,
	}
	controller.dragAutoscroller = helpers.NewDragAutoscroller(
		c.HelperCommon,
		context,
		controller.canDragAutoscroll,
		controller.handleDragAutoscroll,
	)
	return controller
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
			Keys:              opts.GetKeys(opts.Config.Universal.Edit),
			Handler:           self.editLine,
			Description:       self.c.Tr.EditFile,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.EditFile),
			GetDisabledReason: self.diffSelectionDisabledReason,
			Tooltip:           self.c.Tr.EditFileTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Select),
			Handler:           self.primaryAction,
			Description:       self.c.Tr.Stage,
			DescriptionFunc:   self.workingTreeActionDescription(self.c.Tr.Stage),
			GetDisabledReason: self.diffSelectionDisabledReason,
			Tooltip:           self.c.Tr.StageSelectionTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:           self.copySelection,
			Description:       self.c.Tr.CopySelectedTextToClipboard,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.CopySelectedTextToClipboard),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Main.PrevHunk),
			Handler:           self.prevChangeBlock,
			Description:       self.c.Tr.PrevHunk,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.PrevHunk),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Main.NextHunk),
			Handler:           self.nextChangeBlock,
			Description:       self.c.Tr.NextHunk,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.NextHunk),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Main.PrevFile),
			Handler:           self.prevFile,
			Description:       self.c.Tr.PrevFileInDiff,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.PrevFileInDiff),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Main.NextFile),
			Handler:           self.nextFile,
			Description:       self.c.Tr.NextFileInDiff,
			DescriptionFunc:   self.diffSelectionDescriptionText(self.c.Tr.NextFileInDiff),
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
		{
			// Dragging after a click extends a range selection from the clicked line.
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Modifier:    gocui.ModMotion,
			Handler:     self.onDragInFocusedView,
			FocusedView: self.context.GetViewName(),
		},
		{
			ViewName: self.context.GetViewName(),
			Key:      gocui.MouseRelease,
			Handler:  self.onDragRelease,
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

// diffSource returns the panel beneath the focused main view, as the thing that can
// hand out the diff it rendered there. nil when this pane isn't on the stack, or the
// panel beneath shows no diff.
func (self *MainViewController) diffSource() types.FocusedMainViewDiffSource {
	if !self.c.Context().IsInStack(self.context) {
		return nil
	}
	sidePanel := self.c.Context().NextInStack(self.context)
	if sidePanel == nil {
		return nil
	}
	return sidePanel.GetFocusedMainViewDiffSource()
}

// focusedMainViewActions returns what the panel beneath the focused main view does to
// a selection in its diff, or nil where it does nothing to it — a panel whose diff can
// be read and copied but not acted on.
func (self *MainViewController) focusedMainViewActions() types.FocusedMainViewActions {
	actions, _ := self.diffSource().(types.FocusedMainViewActions)
	return actions
}

// primaryAction acts on the selected diff lines, leaving what that means to the panel
// beneath — which also re-renders the diff, since it is the one that changed it.
func (self *MainViewController) primaryAction() error {
	actions := self.focusedMainViewActions()
	if actions == nil {
		return nil
	}
	first, last := self.context.GetView().SelectedLineRange()
	return actions.PrimaryAction(self.context, first, last)
}

// revealSelectionAfterAction moves the focused main view's selection to the change
// that takes the place of the one just acted on, once the changed diff has re-rendered.
// Call it from the panel's action handler with the pane it acted in and the first line
// of the selection, before triggering the re-render.
//
// The line acted on is gone from the diff, so what is remembered is its place among the
// diff's changes: the next change moves up into it, which is where you want to be to
// carry on. A range collapses to a single line at its start, and hunk mode selects the
// whole block it lands in, so that pressing the key again acts on the next hunk.
func revealSelectionAfterAction(c *ControllerCommon, pane types.DiffPaneContext, firstLineIdx int) {
	view := pane.GetView()
	ordinal, ok := c.Helpers().DiffLine.ChangeLineOrdinal(view, firstLineIdx)
	if !ok {
		return
	}

	sel := pane.DiffSelectState()
	if sel.Mode == types.DiffSelectModeRange {
		sel.Mode = types.DiffSelectModeLine
		sel.RangeIsSticky = false
	}
	selectHunk := sel.Mode == types.DiffSelectModeHunk

	c.Helpers().DiffLine.RevealChangeLineAtOrdinal(view, ordinal, func(viewLine int) {
		if selectHunk {
			selectDiffHunk(c, pane, viewLine, true)
			return
		}
		view.CancelRangeSelect()
		showSelectionAtLine(view, viewLine, true)
	})
}

// workingTreeActionDescription gives a command's description only where the command
// applies — over the working tree's diff — so that it is listed there and nowhere else.
func (self *MainViewController) workingTreeActionDescription(description string) func() string {
	return func() string {
		if self.diffMainViewType() != types.DiffMainViewTypeStaging {
			return ""
		}
		return description
	}
}

// copySelection copies the selected diff lines to the clipboard — not as the diff
// renderer drew them, but as they read in the diff itself, which is both what you meant
// to copy and the only form a renderer can't have mangled. A selection that is all
// additions or all deletions loses its +/- column, so that it can be pasted straight
// into code.
//
// The rows above the diff (a commit's message, git's summary of it) belong to no file,
// so they are copied as they stand on screen.
func (self *MainViewController) copySelection() error {
	text := self.textOfSelection()
	if text == "" {
		self.c.ErrorToast(self.c.Tr.SelectionNotFoundInDiffToast)
		return nil
	}

	self.c.LogAction(self.c.Tr.Actions.CopySelectedTextToClipboard)
	if err := self.c.OS().CopyToClipboard(text); err != nil {
		return err
	}
	self.c.Toast(self.c.Tr.SelectedDiffLinesCopiedToast)
	return nil
}

// textOfSelection is what copying the selection puts on the clipboard. It is "" when
// none of the selected rows could be placed in the diff, which is what a rendering's
// own decoration comes to.
func (self *MainViewController) textOfSelection() string {
	source := self.diffSource()
	if source == nil {
		return ""
	}
	view := self.context.GetView()
	first, last := view.SelectedLineRange()
	aboveDiff, fromDiff := self.c.Helpers().DiffLine.PlainDiffOfSelection(view, first, last,
		func(paths []string) string { return source.PlainDiff(self.context, paths) })
	if aboveDiff == "" {
		// Only text that is all diff has a +/- column to lose: a line of a commit message
		// may begin with a '-' without being a deletion of anything.
		fromDiff = dropDiffPrefix(fromDiff)
	}
	return aboveDiff + fromDiff
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

// onDragInFocusedView extends a range selection as the mouse is dragged after a
// click, anchored at the line the click landed on rather than wherever the click left
// the selection — a click can select a whole hunk, whose far end would otherwise
// become the anchor. Dragging turns hunk mode off: you get a plain range from the
// clicked line to the line under the cursor, which gocui has already moved here.
func (self *MainViewController) onDragInFocusedView(opts gocui.ViewMouseBindingOpts) error {
	view := self.context.GetView()
	if !self.isDiffView() || !view.Highlight {
		return nil
	}
	sel := self.diffSelectState()
	sel.Mode = types.DiffSelectModeRange
	sel.RangeIsSticky = false
	sel.UserEnabledHunkMode = false
	view.SetRangeSelectStart(self.context.DragAnchorViewLine())

	// A drag that reaches the edge of the view keeps going: mouse capture means the
	// pointer can be dragged past the edge, and there is more diff down there than
	// fits on screen. opts.Y is where the pointer is in the content, which the
	// autoscroller wants relative to the viewport.
	self.draggingWithMouse = true
	originY, _ := self.context.GetViewTrait().ViewPortYBounds()
	self.dragAutoscroller.Update(opts.Y - originY)
	return nil
}

func (self *MainViewController) onDragRelease(gocui.ViewMouseBindingOpts) error {
	self.draggingWithMouse = false
	self.dragAutoscroller.Cancel()

	// The drag moved the selection without going through showSelectionAtLine: gocui
	// moves the cursor for it. Let the search catch up with where it ended.
	self.context.GetView().SetNearestSearchPosition()
	return nil
}

// GetOnFocusLost stops an autoscroll that is still running when the view loses focus
// mid-drag, e.g. because a popup appeared, and gives up the mouse capture with it —
// otherwise the pointer would keep driving a view that no longer has focus.
func (self *MainViewController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		self.dragAutoscroller.Cancel()
		if self.draggingWithMouse {
			self.draggingWithMouse = false
			self.c.GocuiGui().CancelMouseCapture()
		}
	}
}

// canDragAutoscroll reports whether the autoscroller should run: only while a drag is
// actually extending a range in a diff. Scrolling down also has to keep the lazily
// loaded content ahead of the scroll, or it would stop at the loaded edge.
func (self *MainViewController) canDragAutoscroll(direction int) bool {
	if !self.draggingWithMouse || !self.isDiffView() {
		return false
	}
	view := self.context.GetView()
	if !view.Highlight || self.diffSelectState().Mode != types.DiffSelectModeRange {
		return false
	}
	if direction > 0 {
		self.c.ReadLinesToFillView(view)
	}
	return true
}

// handleDragAutoscroll extends the selection to the line the pointer ends up over
// after the autoscroller has scrolled, leaving the range anchored where the drag
// started. It reports whether the autoscroll should carry on.
//
// The pointer is usually outside the view by now — that is what mouse capture is for —
// so the line it is over is clamped to the visible ones, leaving the selection's far
// end at the edge the scroll is moving towards.
func (self *MainViewController) handleDragAutoscroll(viewLine int) bool {
	if !self.canDragAutoscroll(0) {
		return false
	}
	view := self.context.GetView()
	originY, viewportHeight := self.context.GetViewTrait().ViewPortYBounds()
	target := lo.Clamp(viewLine, 0, max(0, view.ViewLinesHeight()-1))
	view.SetCursorY(lo.Clamp(target-originY, 0, max(0, viewportHeight-1)))
	return true
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
	// Remember where the click landed so that a drag that follows anchors its range
	// there, even when this click selects a whole hunk.
	self.context.SetDragAnchorViewLine(viewLine)
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
		// Remember where the click landed so that a drag that follows anchors its range
		// there, even when this click selects a whole hunk.
		mainContext.SetDragAnchorViewLine(clickedViewLine)
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

	// A search carries on from where the selection now is, so that stepping to the
	// next match goes to the one after it rather than the one after the match the
	// user last stepped to.
	view.SetNearestSearchPosition()
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
	c *ControllerCommon, pane types.DiffPaneContext, changeViewLine int, scrollIntoView bool,
) {
	view := pane.GetView()
	start, end, ok := c.Helpers().DiffLine.ChangeBlockBounds(view, changeViewLine)
	if !ok {
		pane.DiffSelectState().Mode = types.DiffSelectModeLine
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
// alignTop says what a jump does with a target it has to scroll to: bring it to the
// top of the view, or leave the scrolling to place it as it sees fit.
func (self *MainViewController) navigate(find findDiffRowFn, forward bool, alignTop bool) {
	v := self.context.GetView()
	anchor := v.OriginY()
	if v.Highlight {
		anchor = v.SelectedLineIdx()
	}

	if target, ok := find(v, anchor, forward); ok {
		self.placeNavigationTarget(target, alignTop)
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
				self.placeNavigationTarget(target, alignTop)
			}
			return nil
		})
	})
}

// findDiffRowFn locates a row of the rendered diff to navigate to, given the view,
// the anchor view line to start from, and the direction.
type findDiffRowFn func(view *gocui.View, anchorViewLine int, forward bool) (int, bool)

func (self *MainViewController) nextChangeBlock() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentChangeBlock, true, false)
	return nil
}

func (self *MainViewController) prevChangeBlock() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentChangeBlock, false, false)
	return nil
}

// nextFile and prevFile bring the file they go to to the top of the view, since what
// you are going there for is the file, and the more of it is on screen the better.
func (self *MainViewController) nextFile() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentFile, true, true)
	return nil
}

func (self *MainViewController) prevFile() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentFile, false, true)
	return nil
}

// placeNavigationTarget moves the selection to the row a jump found, bringing it on
// screen if it isn't already.
//
// alignTop asks for the target to become the view's top line, so that everything that
// begins there is on screen. It only applies to a target the view has to scroll to: a
// jump to something already on screen leaves the view alone, there being nothing to
// gain from moving what the user is looking at. In hunk mode what ends up selected is
// the first change block at or below the target, which a large context size can put
// further down than a screenful; the selection is then scrolled into view as any other
// jump's is, and the alignment gives way to that.
func (self *MainViewController) placeNavigationTarget(target int, alignTop bool) {
	v := self.context.GetView()
	if !v.Highlight {
		v.SetOrigin(0, target)
		return
	}
	if alignTop {
		self.scrollTargetToTop(target)
	}
	// Jumping to another block or file moves the cursor without shift held, so a
	// range that grows only while shift is held collapses rather than stretching all
	// the way to the target. A sticky range stretches instead; this is the point of
	// being sticky.
	self.collapseNonStickyRange()
	if self.diffSelectState().Mode == types.DiffSelectModeHunk {
		self.selectHunkAround(target, true)
		return
	}
	// Line mode leaves a single-line selection at the target; an active range extends
	// to it, the anchor being untouched.
	showSelectionAtLine(v, target, true)
}

// scrollTargetToTop scrolls the given row of the diff to the top of the view, leaving
// the view where it is when that row is on screen already. The last screenful of the
// diff is as far as it goes, so that the view doesn't scroll past the end of what it is
// showing.
func (self *MainViewController) scrollTargetToTop(target int) {
	view := self.context.GetView()
	originY, height := self.context.GetViewTrait().ViewPortYBounds()
	if target >= originY && target < originY+height {
		return
	}
	view.SetOriginY(min(target, max(0, view.ViewLinesHeight()-height)))
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
	if sel.Mode == types.DiffSelectModeHunk {
		sel.Mode = types.DiffSelectModeLine
		self.context.GetView().CancelRangeSelect()
		return
	}
	self.collapseNonStickyRange()
}

// collapseNonStickyRange drops a range that only grows while shift is held back to a
// single line at the cursor.
func (self *MainViewController) collapseNonStickyRange() {
	sel := self.diffSelectState()
	if sel.Mode == types.DiffSelectModeRange && !sel.RangeIsSticky {
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
		self.navigate(self.c.Helpers().DiffLine.AdjacentChangeBlock, delta > 0, false)
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

func (self *MainViewController) editLine() error {
	view := self.context.GetView()
	if !view.Highlight {
		return nil
	}

	info, ok := self.c.Helpers().DiffLine.GetDiffLineInfo(view, view.SelectedLineIdx())
	if !ok {
		return nil
	}

	// A file-header row points at the file as a whole rather than at a line in it, so
	// it opens the file without jumping anywhere — as pressing edit on a file in a side
	// panel does.
	if info.Type == types.DiffLineFileHeader {
		return self.c.Helpers().Files.EditFiles([]string{info.Path})
	}

	// The diff may be of an older commit, whose line numbers aren't the file's current
	// ones, so they have to be carried forward before we can point an editor at them.
	lineNumber := self.c.Helpers().Diff.AdjustLineNumber(info.Path, info.NewLine, self.context.GetViewName())
	return self.c.Helpers().Files.EditFileAtLine(info.Path, lineNumber)
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
