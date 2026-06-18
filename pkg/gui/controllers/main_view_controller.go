package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
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
	// A selection is shown whenever the main view holds a diff (see
	// sidePanelShowsDiff); we surface the bindings that act on it
	// (enter to dive into staging, e to edit the selected line, G to open the
	// line in the branch's pull request).
	selectionShown := self.context.GetView().Highlight

	var enterDescription string
	var editDescription string
	var editTooltip string
	var openPullRequestDescription string
	var openPullRequestTooltip string
	if selectionShown {
		enterDescription = self.c.Tr.EnterStaging
		editDescription = self.c.Tr.EditFile
		editTooltip = self.c.Tr.EditFileTooltip
		// TODO: i18n-ize these
		openPullRequestDescription = "Open pull request for selected line"
		openPullRequestTooltip = "Open a browser at the selected line in the diff of the current branch's pull request, so that you can comment on it. Only works for local branches that have a pull request on GitHub."
	}

	return []*types.Binding{
		{
			Keys:            opts.GetKeys(opts.Config.Universal.TogglePanel),
			Handler:         self.togglePanel,
			Description:     self.c.Tr.ToggleStagingView,
			Tooltip:         self.c.Tr.ToggleStagingViewTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Return),
			Handler:         self.escape,
			Description:     self.c.Tr.ExitFocusedMainView,
			DisplayOnScreen: true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.Select),
			Handler:         self.stageSelectedLine,
			Description:     self.c.Tr.Stage,
			Tooltip:         self.c.Tr.StageSelectionTooltip,
			DisplayOnScreen: selectionShown,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.GoInto),
			Handler:         self.enter,
			Description:     enterDescription,
			DisplayOnScreen: selectionShown,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.Edit),
			Handler:     self.editLine,
			Description: editDescription,
			Tooltip:     editTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Commits.OpenPullRequestInBrowser),
			Handler:     self.openPullRequestForSelectedLine,
			Description: openPullRequestDescription,
			Tooltip:     openPullRequestTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Main.ToggleSelectHunk),
			Handler:     self.toggleSelectHunk,
			Description: self.c.Tr.ToggleSelectHunk,
			DescriptionFunc: func() string {
				if self.sel().Mode == context.DiffSelectModeHunk {
					return self.c.Tr.SelectLineByLine
				}
				return self.c.Tr.SelectHunk
			},
			Tooltip:         self.c.Tr.ToggleSelectHunkTooltip,
			DisplayOnScreen: selectionShown,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.ToggleRangeSelect),
			Handler:     self.toggleRangeSelect,
			Description: self.c.Tr.ToggleRangeSelect,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Main.PrevHunk),
			Handler:     self.prevChangeBlock,
			Description: self.c.Tr.PrevHunk,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Main.NextHunk),
			Handler:     self.nextChangeBlock,
			Description: self.c.Tr.NextHunk,
		},
		{
			Keys:        opts.GetKeys(config.Keybinding{"N"}),
			Handler:     self.prevFile,
			Description: self.c.Tr.PrevFile,
		},
		{
			Keys:        opts.GetKeys(config.Keybinding{"n"}),
			Handler:     self.nextFile,
			Description: self.c.Tr.NextFile,
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
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.RangeSelectUp), Handler: self.extendRangeUp, Description: self.c.Tr.RangeSelectUp},
		{Tag: "navigation", Keys: opts.GetKeys(opts.Config.Universal.RangeSelectDown), Handler: self.extendRangeDown, Description: self.c.Tr.RangeSelectDown},
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
			// Alt- or shift-click anywhere on a diff line opens it in the editor,
			// without focusing the view or creating a selection. Two modifiers
			// because no single one survives every terminal: Ghostty forwards Alt
			// (and keeps shift for text selection), iTerm2 forwards only shift, and
			// VS Code forwards both. Whichever the terminal delivers triggers the
			// edit; the one it keeps for itself never arrives. No FocusedView, so it
			// fires whatever has focus, and HandleWhenPopupPanelFocused so it stays
			// live when a popup (e.g. the commit-message panel) is in front.
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

// Transient focus shifts (popups, search) leave HighlightInactive=true on our
// view (set by ContextMgr.Activate when a different view becomes current). Our
// context's highlightOnFocus is false, so SimpleContext.HandleFocus never
// resets it. Reset it here on the way back in, so that if we still hold a
// selection it's drawn as active. The flag is a no-op when Highlight is false.
func (self *MainViewController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		self.context.GetView().HighlightInactive = false
	}
}

func (self *MainViewController) togglePanel() error {
	if !self.otherContext.GetView().Visible {
		return nil
	}

	// Capture diff-view-ness while our context is still the focused main view (so
	// NextInStack finds the side panel beneath it), before pushing the other pane.
	isDiff := self.isDiffView()
	self.c.Context().Push(self.otherContext, types.OnFocusOpts{})
	if isDiff {
		showInitialDiffSelection(self.c, self.otherContext)
	}
	return nil
}

// showInitialDiffSelection turns on the focused main view's selection when entering
// a diff view without pointing at a specific line: on the first change line already
// visible (so the view doesn't jump), falling back to the current top line when none
// is visible (scrolled into trailing context, or not loaded that far yet). The
// select mode is reset to its default (a single line).
func showInitialDiffSelection(c *ControllerCommon, mainContext *context.MainContext) {
	resetDiffSelectMode(mainContext)
	view := mainContext.GetView()
	target, ok := c.Helpers().Staging.FirstChangeLineInView(view)
	if !ok {
		target = view.OriginY()
	}
	showSelectionAtLine(view, target, true)
}

// resetDiffSelectMode returns the focused main view to its default select mode — a
// single line, no range — used whenever the selection is (re-)established from
// scratch (on focus, on a click). The view's range anchor is cleared too so the
// next render highlights only the cursor line.
func resetDiffSelectMode(mainContext *context.MainContext) {
	sel := mainContext.DiffSelectState()
	sel.Mode = context.DiffSelectModeLine
	sel.RangeIsSticky = false
	sel.UserEnabledHunkMode = false
	mainContext.GetView().CancelRangeSelect()
}

func (self *MainViewController) escape() error {
	self.c.Context().Pop()
	return nil
}

// isDiffView reports whether the focused main view currently shows a diff (so we
// show a selection in it). See sidePanelShowsDiff.
func (self *MainViewController) isDiffView() bool {
	return sidePanelShowsDiff(self.c.Context().NextInStack(self.context))
}

// stageSelectedLine stages (or unstages) the selected diff line(s) — a single line,
// a range, or a hunk — delegating to the side panel beneath the focused main view
// since what "stage" means is the panel's business (the working tree stages; later,
// commits add to a custom patch). Panels whose diff isn't stageable register no
// handler, so this is a no-op there.
func (self *MainViewController) stageSelectedLine() error {
	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext == nil {
		return nil
	}
	handler := sidePanelContext.GetOnStageFocusedMainView()
	if handler == nil {
		return nil
	}
	first, last := self.context.GetView().SelectedLineRange()
	return handler(self.context.GetViewName(), first, last)
}

func (self *MainViewController) enter() error {
	if !self.context.GetView().Highlight {
		return nil
	}
	return self.enterForLine(self.context.GetView().SelectedLineIdx())
}

// enterForLine dives into staging/patch-building for the given line, by
// delegating to the side panel beneath the focused main view (the same handler
// used when clicking).
func (self *MainViewController) enterForLine(lineIdx int) error {
	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext != nil && sidePanelContext.GetOnClickFocusedMainView() != nil {
		return sidePanelContext.GetOnClickFocusedMainView()(self.context.GetViewName(), lineIdx)
	}
	return nil
}

// showSelectionAtLine turns on the focused main view's selection and moves it to
// the given view line, clamped to the content. scrollIntoView scrolls the line into
// view if it's off-screen (used when navigating to it); a click leaves it false, the
// clicked line being visible already.
func showSelectionAtLine(view *gocui.View, lineIdx int, scrollIntoView bool) {
	view.Highlight = true
	view.HighlightInactive = false
	lineIdx = lo.Clamp(lineIdx, 0, view.ViewLinesHeight()-1)
	view.FocusPoint(0, lineIdx, scrollIntoView)
}

// sel returns our pane's diff selection mode state.
func (self *MainViewController) sel() *context.DiffSelectState {
	return self.context.DiffSelectState()
}

// navigate jumps the focused main view by file or change block (hunk), using find to
// locate the target row from the current anchor. The anchor is the selected line if a
// selection is showing, otherwise the top visible line. With a selection showing we
// move it to the target and scroll it into view, like the staging view — re-selecting
// the whole block in hunk mode; with none we stay in scroll mode, bringing the target
// to the top without selecting anything.
func (self *MainViewController) navigate(find func(*gocui.View, int, bool) (int, bool), forward bool) {
	v := self.context.GetView()
	if !v.Highlight {
		if target, ok := find(v, v.OriginY(), forward); ok {
			v.SetOrigin(0, target)
		}
		return
	}

	target, ok := find(v, v.SelectedLineIdx(), forward)
	if !ok {
		return
	}
	if self.sel().Mode == context.DiffSelectModeHunk {
		self.selectHunkAround(target)
	} else {
		// Line mode leaves a single-line selection at the target; an active range
		// extends to it, since the anchor is untouched.
		showSelectionAtLine(v, target, true)
	}
}

func (self *MainViewController) nextChangeBlock() error {
	self.navigate(self.c.Helpers().Staging.AdjacentChangeBlock, true)
	return nil
}

func (self *MainViewController) prevChangeBlock() error {
	self.navigate(self.c.Helpers().Staging.AdjacentChangeBlock, false)
	return nil
}

func (self *MainViewController) nextFile() error {
	self.navigate(self.c.Helpers().Staging.AdjacentFile, true)
	return nil
}

func (self *MainViewController) prevFile() error {
	self.navigate(self.c.Helpers().Staging.AdjacentFile, false)
	return nil
}

// selectHunkAround re-selects the whole change block around the given change line,
// for hunk mode: the cursor goes to the block's first line and the range anchor to
// its last, so the native range highlight spans the block. With no change block
// there (the line is trailing context) it falls back to a single-line selection.
func (self *MainViewController) selectHunkAround(changeViewLine int) {
	v := self.context.GetView()
	start, end, ok := self.c.Helpers().Staging.ChangeBlockBounds(v, changeViewLine)
	if !ok {
		self.sel().Mode = context.DiffSelectModeLine
		v.CancelRangeSelect()
		showSelectionAtLine(v, changeViewLine, true)
		return
	}
	v.SetRangeSelectStart(end)
	showSelectionAtLine(v, start, true)
}

// moveCursor moves the selection cursor by delta view lines (negative = up), with
// the configured scroll-off margin, reading more content in first when moving down.
// The range anchor is left untouched, so in a range this extends/contracts it; in
// line mode it just moves the selected line.
func (self *MainViewController) moveCursor(delta int) {
	v := self.context.GetView()
	if delta > 0 {
		if manager := self.c.GetViewBufferManagerForView(v); manager != nil {
			manager.ReadLines(delta)
		}
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

// scroll moves the view by delta lines without a selection, for non-diff main
// content where there's nothing to select.
func (self *MainViewController) scroll(delta int) {
	v := self.context.GetView()
	if delta > 0 {
		if manager := self.c.GetViewBufferManagerForView(v); manager != nil {
			manager.ReadLines(delta)
		}
		v.ScrollDown(delta)
	} else {
		v.ScrollUp(-delta)
	}
}

// collapseForLineMove drops hunk mode, and a non-sticky range, back to a single-line
// selection — what a plain (non-shift, non-hunk-step) move does before moving. A
// sticky range is kept so the move extends it.
func (self *MainViewController) collapseForLineMove() {
	sel := self.sel()
	if sel.Mode == context.DiffSelectModeHunk ||
		(sel.Mode == context.DiffSelectModeRange && !sel.RangeIsSticky) {
		sel.Mode = context.DiffSelectModeLine
		self.context.GetView().CancelRangeSelect()
	}
}

// adjustSelection moves the selection by delta view lines for the plain up/down and
// page keys. In hunk mode a single-line step (delta ±1) jumps to the adjacent hunk; a
// larger page step drops out of hunk mode first, like the staging view. A non-sticky
// range collapses back to a single line on a plain move. With no selection (non-diff
// content) it scrolls.
func (self *MainViewController) adjustSelection(delta int) {
	v := self.context.GetView()
	if !v.Highlight {
		self.scroll(delta)
		return
	}
	if self.sel().Mode == context.DiffSelectModeHunk && (delta == 1 || delta == -1) {
		self.navigate(self.c.Helpers().Staging.AdjacentChangeBlock, delta > 0)
		return
	}
	self.collapseForLineMove()
	self.moveCursor(delta)
}

// selectAbsoluteLine moves the selection to a specific view line (the top or bottom
// of the diff), dropping hunk mode and a non-sticky range like a plain move does.
func (self *MainViewController) selectAbsoluteLine(target int) {
	self.collapseForLineMove()
	v := self.context.GetView()
	v.FocusPoint(0, lo.Clamp(target, 0, v.ViewLinesHeight()-1), true)
}

// selectingRange reports whether a range selection is currently active: we're in
// range mode and either it's sticky or the anchor and cursor differ (a non-sticky
// range that has actually been extended).
func (self *MainViewController) selectingRange() bool {
	if self.sel().Mode != context.DiffSelectModeRange {
		return false
	}
	start, end := self.context.GetView().SelectedLineRange()
	return self.sel().RangeIsSticky || start != end
}

// toggleSelectHunk switches between selecting the change block (hunk) around the
// cursor and a single line, mirroring the staging view's `a`.
func (self *MainViewController) toggleSelectHunk() error {
	v := self.context.GetView()
	if !v.Highlight {
		return nil
	}
	sel := self.sel()
	if sel.Mode == context.DiffSelectModeHunk {
		sel.Mode = context.DiffSelectModeLine
		v.CancelRangeSelect()
	} else {
		sel.Mode = context.DiffSelectModeHunk
		sel.UserEnabledHunkMode = true
		self.selectHunkAround(v.SelectedLineIdx())
	}
	return nil
}

// toggleRangeSelect starts or cancels a sticky range selection (extended by plain
// up/down), mirroring the staging view's `v`.
func (self *MainViewController) toggleRangeSelect() error {
	v := self.context.GetView()
	if !v.Highlight {
		return nil
	}
	sel := self.sel()
	if self.selectingRange() {
		sel.Mode = context.DiffSelectModeLine
		sel.RangeIsSticky = false
		v.CancelRangeSelect()
	} else {
		sel.Mode = context.DiffSelectModeRange
		sel.RangeIsSticky = true
		v.SetRangeSelectStart(v.SelectedLineIdx())
	}
	return nil
}

// extendRange grows a (non-sticky) range selection by one line in response to
// shift+up/down, starting one at the cursor if there isn't one yet. Mirrors the
// staging view's range-select keys.
func (self *MainViewController) extendRange(forward bool) error {
	v := self.context.GetView()
	if !v.Highlight {
		return nil
	}
	sel := self.sel()
	if !self.selectingRange() {
		sel.Mode = context.DiffSelectModeRange
		v.SetRangeSelectStart(v.SelectedLineIdx())
	}
	sel.RangeIsSticky = false
	delta := 1
	if !forward {
		delta = -1
	}
	self.moveCursor(delta)
	return nil
}

func (self *MainViewController) extendRangeUp() error {
	return self.extendRange(false)
}

func (self *MainViewController) extendRangeDown() error {
	return self.extendRange(true)
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
		self.scroll(-v.ViewLinesHeight())
		return nil
	}
	self.selectAbsoluteLine(0)
	return nil
}

func (self *MainViewController) handleGotoBottom() error {
	manager := self.c.GetViewBufferManagerForView(self.context.GetView())
	if manager == nil {
		return nil
	}
	manager.ReadToEnd(func() {
		self.c.OnUIThread(func() error {
			v := self.context.GetView()
			if !v.Highlight {
				self.scroll(v.ViewLinesHeight())
				return nil
			}
			self.selectAbsoluteLine(v.ViewLinesHeight() - 1)
			return nil
		})
	})
	return nil
}

// sidePanelShowsDiff reports whether the given side panel's focused main view
// shows a diff, which is when we show a selection in it (so the user can stage a
// line, edit it, jump by hunk/file, or open it in a PR). Panels whose main view
// shows non-diff content (a branch's commit log, the status dashboard, …) show no
// selection because there's nothing to act on. See types.DiffMainViewContext.
func sidePanelShowsDiff(sidePanel types.Context) bool {
	_, ok := sidePanel.(types.DiffMainViewContext)
	return ok
}

// focusedMainViewContextForViewName maps a focused main view's view name (as
// passed to GetOnClickFocusedMainView) to its context.
func focusedMainViewContextForViewName(c *ControllerCommon, viewName string) types.Context {
	if viewName == c.Contexts().NormalSecondary.GetViewName() {
		return c.Contexts().NormalSecondary
	}
	return c.Contexts().Normal
}

// focusedMainViewSnapshot records the focused main view to return to when diving
// into a patch explorer from it, so escaping can come back with the main view
// focused. sidePanel is the panel to land on first (which re-renders the
// content); for commits/stash it's the originating panel, skipping the commit
// files panel we pass through. Where to scroll to and select on return isn't
// captured: escape lands on the line the explorer ended up on (see
// EscapeFromPatchExplorer). Call this before any mutation that might re-render
// the main view.
func focusedMainViewSnapshot(c *ControllerCommon, mainViewName string, sidePanel types.Context) *types.FocusedMainViewSnapshot {
	mainView := focusedMainViewContextForViewName(c, mainViewName)
	sidePanelSelectedLineIdx := -1
	if listContext, ok := sidePanel.(types.IListContext); ok {
		sidePanelSelectedLineIdx = listContext.GetList().GetSelectedLineIdx()
	}
	return &types.FocusedMainViewSnapshot{
		SidePanel:                sidePanel,
		SidePanelSelectedLineIdx: sidePanelSelectedLineIdx,
		MainView:                 mainView,
	}
}

func (self *MainViewController) editLine() error {
	if !self.context.GetView().Highlight {
		return nil
	}
	return self.editDiffLine(self.context.GetView().SelectedLineIdx())
}

// editDiffLine opens the file the given diff line belongs to in the editor, at
// that line. The file and line are resolved the same way entering staging does.
func (self *MainViewController) editDiffLine(viewLineIdx int) error {
	info, ok := self.c.Helpers().Staging.GetDiffLineInfo(self.context.GetViewName(), viewLineIdx)
	if !ok {
		return nil
	}
	lineNumber := self.c.Helpers().Diff.AdjustLineNumber(info.Path, info.NewLine, self.context.GetViewName())
	return self.c.Helpers().Files.EditFileAtLine(info.Path, lineNumber)
}

func (self *MainViewController) openPullRequestForSelectedLine() error {
	if !self.context.GetView().Highlight {
		return nil
	}

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

	// Figure out the clicked file and line the same way entering staging does.
	info, ok := self.c.Helpers().Staging.GetDiffLineInfo(
		self.context.GetViewName(), self.context.GetView().SelectedLineIdx())
	if !ok {
		return nil
	}

	relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), info.Path)
	if err != nil {
		return err
	}

	// A deletion isn't on the right (new) side of the diff, so anchor it on the
	// left (old) side; everything else on the right.
	side, lineNumber := info.PullRequestAnchor()

	self.c.LogAction(self.c.Tr.Actions.OpenPullRequest)
	return self.c.OS().OpenLink(
		githubPullRequestLineURL(pr.Url, commitSha, filepath.ToSlash(relativePath), side, lineNumber))
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
// identified by the SHA-256 of its repo-relative path, and side ("R"/"L")
// selects the right (new) or left (old) side of the diff. See
// https://github.com/orgs/community/discussions/55764.
func githubPullRequestLineURL(prURL string, commitSha string, relativePath string, side string, lineNumber int) string {
	pathHash := sha256.Sum256([]byte(relativePath))
	anchor := fmt.Sprintf("diff-%s%s%d", hex.EncodeToString(pathHash[:]), side, lineNumber)
	return fmt.Sprintf("%s/changes/%s#%s", prURL, commitSha, anchor)
}

func (self *MainViewController) onClickInAlreadyFocusedView(opts gocui.ViewMouseBindingOpts) error {
	if !self.isDiffView() {
		return nil
	}
	// A click points at a line, so it sets a single-line selection there; a
	// double-click additionally dives into staging/patch-building for that line.
	resetDiffSelectMode(self.context)
	showSelectionAtLine(self.context.GetView(), opts.Y, false)
	if opts.IsDoubleClick {
		return self.enterForLine(opts.Y)
	}
	return nil
}

func (self *MainViewController) editClickedLine(opts gocui.ViewMouseBindingOpts) error {
	return self.editDiffLine(opts.Y)
}

func (self *MainViewController) onClickInOtherViewOfMainViewPair(opts gocui.ViewMouseBindingOpts) error {
	self.c.Context().Push(self.context, types.OnFocusOpts{})
	if !self.isDiffView() {
		return nil
	}
	resetDiffSelectMode(self.context)
	showSelectionAtLine(self.context.GetView(), opts.Y, false)
	if opts.IsDoubleClick {
		return self.enterForLine(opts.Y)
	}
	return nil
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
