package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
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

	dragAutoscroller    *helpers.DragAutoscroller
	draggingWithMouse   bool
	lineFlashGeneration uint64
}

const editedLineFlashDuration = 200 * time.Millisecond

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
			Description:     self.c.Tr.ToggleDiffPane,
			Tooltip:         self.c.Tr.ToggleDiffPaneTooltip,
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
			Keys:    opts.GetKeys(opts.Config.Universal.Select),
			Handler: self.primaryAction,
			// The description is of the working tree's diff, which is where the key does
			// the thing users know it for; over a commit's diff it says so for itself.
			Description:       self.c.Tr.Stage,
			DescriptionFunc:   self.diffActionDescription(self.c.Tr.Stage, self.c.Tr.ToggleSelectionForPatch),
			GetDisabledReason: self.diffSelectionDisabledReason,
			Tooltip:           self.c.Tr.StageSelectionTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Remove),
			Handler:           self.discardSelection,
			Description:       self.c.Tr.DiscardSelection,
			DescriptionFunc:   self.diffActionDescription(self.c.Tr.DiscardSelection, self.c.Tr.RemoveSelectionFromPatch),
			GetDisabledReason: self.discardSelectionDisabledReason,
			Tooltip:           self.c.Tr.DiscardSelectionTooltip,
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
			Keys:              opts.GetKeys(config.Keybinding{"f"}),
			Handler:           self.openJumpToFileMenu,
			Description:       "Jump to file",
			DescriptionFunc:   self.diffSelectionDescriptionText("Jump to file"),
			GetDisabledReason: self.diffSelectionDisabledReason,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Commits.OpenPullRequestInBrowser),
			Handler:           self.openPullRequestForSelectedLine,
			Description:       "Open pull request for selected line",
			DescriptionFunc:   self.diffSelectionDescriptionText("Open pull request for selected line"),
			Tooltip:           "Open a browser at the selected line in the diff of the current branch's pull request, so that you can comment on it. Only works for local branches that have a pull request on GitHub.",
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
		{
			Keys:            opts.GetKeys(opts.Config.Files.CommitChanges),
			Handler:         self.workingTreeAction(self.c.Helpers().WorkingTree.HandleCommitPress),
			Description:     self.c.Tr.Commit,
			DescriptionFunc: self.workingTreeActionDescription(self.c.Tr.Commit),
			Tooltip:         self.c.Tr.CommitTooltip,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Files.CommitChangesWithoutHook),
			Handler:         self.workingTreeAction(self.c.Helpers().WorkingTree.HandleWIPCommitPress),
			Description:     self.c.Tr.CommitChangesWithoutHook,
			DescriptionFunc: self.workingTreeActionDescription(self.c.Tr.CommitChangesWithoutHook),
		},
		{
			Keys:            opts.GetKeys(opts.Config.Files.CommitChangesWithEditor),
			Handler:         self.workingTreeAction(self.c.Helpers().WorkingTree.HandleCommitEditorPress),
			Description:     self.c.Tr.CommitChangesWithEditor,
			DescriptionFunc: self.workingTreeActionDescription(self.c.Tr.CommitChangesWithEditor),
		},
		{
			Keys:            opts.GetKeys(opts.Config.Files.FindBaseCommitForFixup),
			Handler:         self.workingTreeAction(self.c.Helpers().FixupHelper.HandleFindBaseCommitForFixupPress),
			Description:     self.c.Tr.FindBaseCommitForFixup,
			DescriptionFunc: self.workingTreeActionDescription(self.c.Tr.FindBaseCommitForFixup),
			Tooltip:         self.c.Tr.FindBaseCommitForFixupTooltip,
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
		{
			ViewName:                    self.context.GetViewName(),
			Key:                         gocui.MouseLeft,
			Modifier:                    gocui.ModAlt,
			Handler:                     self.editClickedLine,
			HandleWhenPopupPanelFocused: true,
		},
		{
			ViewName:                    self.context.GetViewName(),
			Key:                         gocui.MouseLeft,
			Modifier:                    gocui.ModShift,
			Handler:                     self.editClickedLine,
			HandleWhenPopupPanelFocused: true,
		},
	}
}

func (self *MainViewController) Context() types.Context {
	return self.context
}

// GetOnFocus brings on the marks over the lines that are in the custom patch, which
// are an affordance of the focused view, so they arrive with the focus.
func (self *MainViewController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		self.c.Helpers().DiffLine.RefreshInclusionGutter()
	}
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
		self.c.Helpers().DiffLine.EstablishSelection(self.otherContext, -1)
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

// discardSelection takes the selected diff lines back out of what they are part of,
// which — like the primary action — is the panel's business, and so is the re-render
// that follows.
func (self *MainViewController) discardSelection() error {
	actions := self.focusedMainViewActions()
	if actions == nil {
		return nil
	}
	first, last := self.context.GetView().SelectedLineRange()
	return actions.DiscardSelection(self.context, first, last)
}

// workingTreeAction wraps a command that acts on the working tree — committing, finding
// the commit to fix up — so that it only runs while the focused main view is showing the
// working tree's diff. Over a commit's diff the key does nothing, so that browsing
// through history can't commit by accident. The check is per press, since what the main
// view shows changes as the user moves around while the keybindings are registered once.
func (self *MainViewController) workingTreeAction(action func() error) func() error {
	return func() error {
		if self.diffMainViewType() != types.DiffMainViewTypeStaging {
			return nil
		}
		return action()
	}
}

// workingTreeActionDescription gives a command's description only where the command
// applies — over the working tree's diff — so that it is listed there and nowhere else.
func (self *MainViewController) workingTreeActionDescription(description string) func() string {
	return self.diffActionDescription(description, "")
}

// diffActionDescription describes a command in the words of the diff it is over: acting
// on the working tree's diff stages, acting on a commit's builds a custom patch. Over
// content that is no diff at all the command doesn't apply, and describes itself as
// nothing, which keeps it out of the keybindings menu there.
func (self *MainViewController) diffActionDescription(staging string, patchBuilding string) func() string {
	return func() string {
		switch self.diffMainViewType() {
		case types.DiffMainViewTypeStaging:
			return staging
		case types.DiffMainViewTypePatchBuilding:
			return patchBuilding
		default:
			return ""
		}
	}
}

// copySelection copies the selected diff lines to the clipboard — not as the diff
// renderer drew them, but as they read in the diff itself, which is both what you meant
// to copy and the only form a renderer can't have mangled. A selection that is all
// additions or all deletions loses its +/- column, so that it can be pasted straight
// into code.
func (self *MainViewController) copySelection() error {
	source := self.diffSource()
	if source == nil {
		return nil
	}
	view := self.context.GetView()
	first, last := view.SelectedLineRange()
	text := self.c.Helpers().DiffLine.PlainDiffOfSelection(view, first, last,
		func(paths []string) string { return source.PlainDiff(self.context, paths) })
	if text == "" {
		return nil
	}

	self.c.LogAction(self.c.Tr.Actions.CopySelectedTextToClipboard)
	return self.c.OS().CopyToClipboard(dropDiffPrefix(text))
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

// discardSelectionDisabledReason disables discarding while there is nothing to discard,
// and where the panel beneath won't have it: taking lines out of a commit means
// rewriting it, which isn't always something we may do.
func (self *MainViewController) discardSelectionDisabledReason() *types.DisabledReason {
	if reason := self.diffSelectionDisabledReason(); reason != nil {
		return reason
	}
	if actions := self.focusedMainViewActions(); actions != nil {
		return actions.DiscardSelectionDisabledReason(self.context)
	}
	return nil
}

func (self *MainViewController) onClickInAlreadyFocusedView(opts gocui.ViewMouseBindingOpts) error {
	self.selectClickedDiffLine(opts.Y)
	return nil
}

func (self *MainViewController) editClickedLine(opts gocui.ViewMouseBindingOpts) error {
	var flashGeneration uint64
	err := self.editDiffLine(opts.Y, func() {
		self.lineFlashGeneration++
		flashGeneration = self.lineFlashGeneration
		self.context.GetView().SetLineFlash(opts.Y)
		self.c.GocuiGui().ForceFlushViewsContentOnly(self.c.GocuiGui().Views())
	})
	if flashGeneration != 0 {
		time.AfterFunc(editedLineFlashDuration, func() {
			self.c.OnUIThreadContentOnlyBackground(func() error {
				if self.lineFlashGeneration == flashGeneration {
					self.context.GetView().ClearLineFlash()
				}
				return nil
			})
		})
	}
	return err
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
		// Where the focus has gone is already known here, so asking again is what keeps
		// the patch marks over a move to the pane beside this one and takes them away
		// when the focus leaves the pair.
		self.c.Helpers().DiffLine.RefreshInclusionGutter()
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
			self.c.Helpers().DiffLine.ShowSelectionAtLine(view, viewLine, false)
			return
		}
		if self.c.Helpers().DiffLine.IsChangeLine(view, viewLine) {
			self.selectHunkAround(viewLine, false)
			return
		}
	}
	self.context.ResetDiffSelectMode()
	self.c.Helpers().DiffLine.ShowSelectionAtLine(view, viewLine, false)
}

func (self *MainViewController) selectHunkAround(changeViewLine int, scrollIntoView bool) {
	self.c.Helpers().DiffLine.SelectChangeBlock(self.context, changeViewLine, scrollIntoView)
}

// navigate moves the focused main view to the row find locates from the current
// anchor — the selected line when a selection is showing, otherwise the top visible
// line. With a selection we move it there and scroll it into view, like the staging
// view, re-selecting the whole block in hunk mode; with none we stay in scroll mode,
// bringing the target to the top without selecting anything.
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

func (self *MainViewController) nextChangeBlock() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentChangeBlock, true)
	return nil
}

func (self *MainViewController) prevChangeBlock() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentChangeBlock, false)
	return nil
}

func (self *MainViewController) nextFile() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentFile, true)
	return nil
}

func (self *MainViewController) prevFile() error {
	self.navigate(self.c.Helpers().DiffLine.AdjacentFile, false)
	return nil
}

func (self *MainViewController) placeNavigationTarget(target int) {
	v := self.context.GetView()
	if !v.Highlight {
		v.SetOrigin(0, target)
		return
	}
	// Jumping to another block or file moves the cursor without shift held, so a
	// range that only grows while it is collapses rather than stretching all the way
	// to the target. A sticky range does stretch — that is what makes it sticky.
	self.collapseNonStickyRange()
	if self.diffSelectState().Mode == types.DiffSelectModeHunk {
		self.selectHunkAround(target, true)
		return
	}
	// Line mode leaves a single-line selection at the target; an active range extends
	// to it, the anchor being untouched.
	self.c.Helpers().DiffLine.ShowSelectionAtLine(v, target, true)
}

// openJumpToFileMenu pops up a menu listing the files in the focused main view's diff, in
// the order they appear, as repo-relative paths; picking one jumps straight to it. It's a
// complement to n / N for a diff that spans many files. A no-op when the main view holds
// no diff (nothing to list).
//
// The diff loads lazily, so we read it to the end first — otherwise a file past the loaded
// portion of a long diff would be missing from the menu (and have no view line to jump to).
// Same as handleGotoBottom.
func (self *MainViewController) openJumpToFileMenu() error {
	manager := self.c.GetViewBufferManagerForView(self.context.GetView())
	if manager == nil {
		return nil
	}
	manager.ReadToEnd(func() {
		self.c.OnUIThread(func() error {
			return self.showJumpToFileMenu()
		})
	})
	return nil
}

func (self *MainViewController) showJumpToFileMenu() error {
	files := self.c.Helpers().DiffLine.FilesInDiff(self.context.GetView())
	if len(files) == 0 {
		return nil
	}

	worktreePath := self.c.Git().RepoPaths.WorktreePath()
	menuItems := lo.Map(files, func(file helpers.DiffFile, index int) *types.MenuItem {
		label := file.Path
		if rel, err := filepath.Rel(worktreePath, file.Path); err == nil {
			label = rel
		}
		firstViewLine := file.FirstViewLine
		return &types.MenuItem{
			Label:   label,
			OnPress: func() error { self.jumpToFile(firstViewLine); return nil },
		}
	})

	// TODO: i18n-ize this title
	return self.c.Menu(types.CreateMenuOptions{Title: "Jump to file", Items: menuItems, HideCancel: true, FilterAsYouType: true})
}

// jumpToFile moves the focused main view to the given view line exactly the way
// next/previous-file navigation does (see navigate), so jumping from the menu and
// stepping with n / N land identically — scrolling the file to the top with no selection,
// or moving the selection (or hunk selection) to it when one is showing.
func (self *MainViewController) jumpToFile(firstViewLine int) {
	self.navigate(func(*gocui.View, int, bool) (int, bool) {
		return firstViewLine, true
	}, true)
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
	v.FocusPoint(0, after, true)
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
	v := self.context.GetView()
	v.FocusPoint(0, lo.Clamp(target, 0, v.ViewLinesHeight()-1), true)
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
	return self.editDiffLine(view.SelectedLineIdx(), nil)
}

func (self *MainViewController) editDiffLine(viewLine int, beforeEdit func()) error {
	info, ok := self.c.Helpers().DiffLine.GetDiffLineInfo(self.context.GetView(), viewLine)
	if !ok {
		return nil
	}
	if beforeEdit != nil {
		beforeEdit()
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

func (self *MainViewController) openPullRequestForSelectedLine() error {
	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext == nil {
		return nil
	}

	// The branch whose PR to open depends on where we navigated from: the
	// checked-out branch when looking at its own commits, but the branch we
	// drilled into when in the sub-commits or commit-files panels.
	branchName, ok := self.branchForPullRequest(sidePanelContext)
	if !ok {
		return nil
	}

	pr, ok := self.c.Model().PullRequestsMap[branchName]
	if !ok {
		return errors.New(self.c.Tr.NoPullRequestForBranch)
	}

	// The diff shown is the diff of a particular commit, so we deep-link into
	// that commit's view of the PR; its right-side line numbers match what we're
	// showing, so (unlike editLine) no line-number adjustment is needed.
	diffableContext, ok := sidePanelContext.(types.DiffableContext)
	if !ok {
		return nil
	}
	commitSha := diffableContext.RefForAdjustingLineNumberInDiff()
	if commitSha == "" {
		return nil
	}

	info, ok := self.c.Helpers().DiffLine.GetDiffLineInfo(self.context.GetView(), self.context.GetView().SelectedLineIdx())
	if !ok {
		return nil
	}

	relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), info.Path)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.OpenPullRequest)
	return self.c.OS().OpenLink(
		githubPullRequestLineURL(pr.Url, commitSha, filepath.ToSlash(relativePath), info.NewLine))
}

// branchForPullRequest returns the local branch whose pull request applies to
// the diff currently shown in the focused main view, given the side panel
// beneath it. It returns false for contexts that don't map to a local branch
// (e.g. the working-tree files panel, stashes, tags, or remote branches).
func (self *MainViewController) branchForPullRequest(sidePanelContext types.Context) (string, bool) {
	switch sidePanelContext.GetKey() {
	case context.LOCAL_COMMITS_CONTEXT_KEY:
		return self.c.Model().CheckedOutBranch, true
	case context.SUB_COMMITS_CONTEXT_KEY:
		ref := self.c.Contexts().SubCommits.GetRef()
		if ref == nil {
			return "", false
		}
		return ref.RefName(), true
	case context.COMMIT_FILES_CONTEXT_KEY:
		// The commit files panel doesn't itself know which branch it belongs to;
		// that's determined by the panel we entered it from.
		parent := self.c.Contexts().CommitFiles.GetParentContext()
		if parent == nil {
			return "", false
		}
		return self.branchForPullRequest(parent)
	default:
		return "", false
	}
}

// githubPullRequestLineURL builds a URL that opens the given line of a file in
// the diff of a specific commit within a GitHub pull request. The file is
// identified by the SHA-256 of its repo-relative path, and R<line> targets the
// right (new) side of the diff. See
// https://github.com/orgs/community/discussions/55764.
func githubPullRequestLineURL(prURL string, commitSha string, relativePath string, lineNumber int) string {
	pathHash := sha256.Sum256([]byte(relativePath))
	anchor := fmt.Sprintf("diff-%sR%d", hex.EncodeToString(pathHash[:]), lineNumber)
	return fmt.Sprintf("%s/changes/%s#%s", prURL, commitSha, anchor)
}

func (self *MainViewController) openSearch() error {
	if manager := self.c.GetViewBufferManagerForView(self.context.GetView()); manager != nil {
		manager.ReadToEnd(func() {
			self.c.OnUIThread(func() error {
				return self.c.Helpers().Search.OpenSearchPrompt(self.context)
			})
		})
	}

	return nil
}
