package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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

	// The commit and find-fixup-base commands act on the working tree, so — like the
	// old staging panel, which was always the working-tree diff — we only surface them
	// while the focused main view plays that role (see workingTreeAction).
	stagingMode := self.diffMainViewType() == types.DiffMainViewTypeStaging

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

	var commitDescription string
	var commitTooltip string
	var commitWithoutHookDescription string
	var commitWithEditorDescription string
	var findBaseCommitDescription string
	var findBaseCommitTooltip string
	if stagingMode {
		commitDescription = self.c.Tr.Commit
		commitTooltip = self.c.Tr.CommitTooltip
		commitWithoutHookDescription = self.c.Tr.CommitChangesWithoutHook
		commitWithEditorDescription = self.c.Tr.CommitChangesWithEditor
		findBaseCommitDescription = self.c.Tr.FindBaseCommitForFixup
		findBaseCommitTooltip = self.c.Tr.FindBaseCommitForFixupTooltip
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
			Handler:         self.primaryAction,
			Description:     self.c.Tr.Stage,
			Tooltip:         self.c.Tr.StageSelectionTooltip,
			DisplayOnScreen: selectionShown,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Remove),
			Handler:           self.discardSelection,
			GetDisabledReason: self.discardSelectionDisabledReason,
			Description:       self.c.Tr.DiscardSelection,
			Tooltip:           self.c.Tr.DiscardSelectionTooltip,
			DisplayOnScreen:   selectionShown,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Universal.CopyToClipboard),
			Handler:         self.copySelection,
			Description:     self.c.Tr.CopySelectedTextToClipboard,
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
			Keys:        opts.GetKeys(opts.Config.Files.CommitChanges),
			Handler:     self.workingTreeAction(self.c.Helpers().WorkingTree.HandleCommitPress),
			Description: commitDescription,
			Tooltip:     commitTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.CommitChangesWithoutHook),
			Handler:     self.workingTreeAction(self.c.Helpers().WorkingTree.HandleWIPCommitPress),
			Description: commitWithoutHookDescription,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.CommitChangesWithEditor),
			Handler:     self.workingTreeAction(self.c.Helpers().WorkingTree.HandleCommitEditorPress),
			Description: commitWithEditorDescription,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.FindBaseCommitForFixup),
			Handler:     self.workingTreeAction(self.c.Helpers().FixupHelper.HandleFindBaseCommitForFixupPress),
			Description: findBaseCommitDescription,
			Tooltip:     findBaseCommitTooltip,
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
			// Dragging after a click extends a range selection from the clicked line.
			ViewName:    self.context.GetViewName(),
			Key:         gocui.MouseLeft,
			Modifier:    gocui.ModMotion,
			Handler:     self.onDragInFocusedView,
			FocusedView: self.context.GetViewName(),
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
		// The inclusion gutter is shown only while the focused main view holds focus, so
		// (re-)establish it now (a no-op unless the panel beneath builds a custom patch).
		self.c.Helpers().Staging.RefreshInclusionGutter()
	}
}

// GetOnFocusLost re-evaluates the inclusion gutter as the focused main view loses focus.
// It's a focused-main-view affordance of the whole pair, so it must persist when tabbing
// between the two panes but hide when focus leaves the pair entirely. By the time this runs
// the new context is already current, so RefreshInclusionGutter decides correctly (and so
// the gutter doesn't flicker off-then-on across a pane switch).
func (self *MainViewController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		self.c.Helpers().Staging.RefreshInclusionGutter()
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

// showInitialDiffSelection turns on the focused main view's selection when entering a
// diff view by keyboard, without pointing at a specific line. See
// establishFocusedDiffSelection, which it defers to (clicks go through there too,
// passing the clicked line).
func showInitialDiffSelection(c *ControllerCommon, mainContext *context.MainContext) {
	establishFocusedDiffSelection(c, mainContext, -1)
}

// establishFocusedDiffSelection turns on the focused main view's selection after the
// view is focused. Under a pager that doesn't speak the metadata protocol (probed, so
// known up front; see StagingHelper.DiffMainViewShouldRenderRaw) the diff on screen was
// rendered pretty for browsing and isn't resolvable, so it re-renders it raw and places
// the selection once that lands; otherwise it places the selection on the diff directly.
// clickedViewLine is the view line a click pointed at, or -1 for keyboard focus (start
// at the first change block).
func establishFocusedDiffSelection(c *ControllerCommon, mainContext *context.MainContext, clickedViewLine int) {
	staging := c.Helpers().Staging
	resetDiffSelectMode(mainContext)

	if staging.DiffMainViewShouldRenderRaw() {
		sidePanel := c.Context().NextInStack(mainContext)
		staging.RenderFocusedMainViewRaw(mainContext.GetView(), sidePanel, func() {
			placeOrHideInitialDiffSelection(c, mainContext, clickedViewLine, true)
		})
		return
	}

	placeOrHideInitialDiffSelection(c, mainContext, clickedViewLine, clickedViewLine < 0)
}

// placeOrHideInitialDiffSelection puts the focused main view's selection on the clicked
// line (clickedViewLine >= 0) or, for keyboard focus, on the first change line at or
// below the top of the viewport — so the view barely moves — falling back to the top
// line when none is visible. With hunk mode configured as the default the selection
// widens to the whole change block, like entering the staging view does: keyboard focus
// snaps to the block at the first visible change, and a click on a change line selects
// that line's block, ready to stage. A click on context still selects just that line —
// the click points at it precisely, so it can be edited (e) rather than staged.
// When the diff has nothing to act on (a placeholder, a binary file, an all-context
// diff) it shows no selection rather than highlighting a stray line.
func placeOrHideInitialDiffSelection(c *ControllerCommon, mainContext *context.MainContext, clickedViewLine int, scrollIntoView bool) {
	view := mainContext.GetView()
	if !c.Helpers().Staging.ViewHasChangeLines(view) {
		view.Highlight = false
		return
	}
	if clickedViewLine >= 0 {
		// Remember where the click landed so a drag that follows anchors its range there,
		// even when the click selects a whole hunk (whose anchor is the block's far end).
		mainContext.SetDragAnchorViewLine(clickedViewLine)
		if c.UserConfig().Gui.UseHunkModeInStagingView && c.Helpers().Staging.IsChangeLine(view, clickedViewLine) {
			mainContext.DiffSelectState().Mode = context.DiffSelectModeHunk
			selectDiffHunk(c, mainContext, clickedViewLine)
			return
		}
		showSelectionAtLine(view, clickedViewLine, scrollIntoView)
		return
	}
	target, ok := c.Helpers().Staging.FirstChangeLineInView(view)
	if !ok {
		showSelectionAtLine(view, view.OriginY(), true)
		return
	}
	if c.UserConfig().Gui.UseHunkModeInStagingView {
		mainContext.DiffSelectState().Mode = context.DiffSelectModeHunk
		selectDiffHunk(c, mainContext, target)
		return
	}
	showSelectionAtLine(view, target, true)
}

// updateFocusedMainViewSelectionVisibility shows or hides the focused-main-view selection
// to match what a side panel is rendering into the main view, called from the panel's
// render-to-main so the selection tracks content changes (a refresh after the last change
// is discarded, or changes vanishing / appearing outside lazygit). A selection is shown
// only on the main pane that currently holds focus, and only when it's rendering a diff
// (something to act on) — never over "No changed files" or a merge-conflict message.
// normalHasDiff/secondaryHasDiff say whether each pane is being given a diff; the caller
// knows this from which content it's about to render (the rendered content can't be read
// here, since the render it triggers is asynchronous). Initial keyboard/click focus is
// handled separately by showInitialDiffSelection, since focusing reuses the already-
// rendered content rather than re-rendering.
func updateFocusedMainViewSelectionVisibility(c *ControllerCommon, normalHasDiff bool, secondaryHasDiff bool) {
	focusedKey := c.Context().CurrentStatic().GetKey()
	normal := c.Contexts().Normal
	secondary := c.Contexts().NormalSecondary
	normal.GetView().Highlight = normalHasDiff && focusedKey == normal.GetKey()
	secondary.GetView().Highlight = secondaryHasDiff && focusedKey == secondary.GetKey()
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

// focusedMainViewActions returns the actions the side panel beneath the focused main
// view offers on its diff (diving in, staging, patch toggling), or nil when there is no
// panel beneath or its diff offers none.
func (self *MainViewController) focusedMainViewActions() types.FocusedMainViewActions {
	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext == nil {
		return nil
	}
	return sidePanelContext.GetFocusedMainViewActions()
}

// diffMainViewType reports what the side panel beneath the focused main view makes its
// primary action mean — staging (the working tree), patch-building, or nothing — or
// DiffMainViewTypeNone when this pane isn't focused or has no diff panel beneath it. The
// IsInStack guard is essential: NextInStack panics for a context that isn't in the stack,
// and GetKeybindings (which calls this) runs for off-stack panes too — at startup and
// during cheatsheet generation, where the stack is empty.
func (self *MainViewController) diffMainViewType() types.DiffMainViewType {
	if !self.c.Context().IsInStack(self.context) {
		return types.DiffMainViewTypeNone
	}
	if diffContext, ok := self.c.Context().NextInStack(self.context).(types.DiffMainViewContext); ok {
		return diffContext.GetDiffMainViewType()
	}
	return types.DiffMainViewTypeNone
}

// workingTreeAction wraps a working-tree command (commit, find-fixup-base) so it runs
// only while the focused main view shows the working-tree diff — the role the staging
// panel used to play, and the only place these commands mean anything. Over any other
// diff (a commit, a stash) the key is a no-op, so committing can't be triggered by
// accident while browsing history. Re-checks on each press rather than capturing the
// state at keybinding-registration time.
func (self *MainViewController) workingTreeAction(action func() error) func() error {
	return func() error {
		if self.diffMainViewType() != types.DiffMainViewTypeStaging {
			return nil
		}
		return action()
	}
}

// primaryAction acts on the selected diff line(s) — a single line, a range, or a hunk —
// delegating to the side panel beneath the focused main view, since what the action means
// is the panel's business: the working tree stages, while commits toggle the selection
// into a custom patch. The handler does its own re-render and re-establishes the selection
// afterwards (see revealSelectionAfterPrimaryAction), so the dispatcher just hands over
// the selected range. A no-op when the panel beneath offers no actions.
func (self *MainViewController) primaryAction() error {
	actions := self.focusedMainViewActions()
	if actions == nil {
		return nil
	}
	v := self.context.GetView()
	first, last := v.SelectedLineRange()
	return actions.PrimaryAction(self.context.GetViewName(), first, last)
}

// discardSelection discards the selected diff line(s), delegating to the panel beneath
// (working-tree discard for files, line-removal from the commit for the commit panels).
// As with the primary action, the handler does its own re-render. A no-op when the panel
// beneath offers no actions; discardSelectionDisabledReason gates the cases where the
// panel has discard but it isn't applicable right now.
func (self *MainViewController) discardSelection() error {
	actions := self.focusedMainViewActions()
	if actions == nil {
		return nil
	}
	v := self.context.GetView()
	first, last := v.SelectedLineRange()
	return actions.DiscardSelection(self.context.GetViewName(), first, last)
}

func (self *MainViewController) discardSelectionDisabledReason() *types.DisabledReason {
	actions := self.focusedMainViewActions()
	if actions == nil {
		return nil
	}
	return actions.DiscardSelectionDisabledReason(self.context.GetViewName())
}

// copySelection copies the selected diff line(s) to the clipboard. Unlike the primary
// action and discard, this is the same for every diff panel (it just reads the text the
// focused main view shows), so it lives here rather than on FocusedMainViewActions — and
// so it works over panels with no actions, like the reflog. A no-op when there's no
// selection (non-diff content).
//
// With no pager the main view shows the raw diff, so the +/-/space column is stripped
// from a homogeneous selection (dropDiffPrefix) to ease pasting into code, mirroring the
// staging view's copy. With a pager configured we can't assume that column is present —
// some pagers keep it and only colorize (e.g. ydiff), others drop or restructure it (e.g.
// delta) — and we can't tell which, so we conservatively don't strip and copy verbatim.
// (The cost is missing the convenience for column-preserving pagers.)
func (self *MainViewController) copySelection() error {
	v := self.context.GetView()
	if !v.Highlight {
		return nil
	}

	// The selection is in (wrapped) view-line space; map it to the unwrapped buffer lines
	// it covers and copy those, so a wrapped line is copied whole and once — not indexed
	// into the buffer by a view-line number (which copies the wrong line, or panics when
	// the view line is past the buffer's length).
	firstView, lastView := v.SelectedLineRange()
	firstBuffer, ok := v.BufferLineForViewLine(firstView)
	if !ok {
		return nil
	}
	lastBuffer, ok := v.BufferLineForViewLine(lastView)
	if !ok {
		return nil
	}

	contents := v.DiffLineContents()
	lastBuffer = min(lastBuffer, len(contents)-1)
	if firstBuffer > lastBuffer {
		return nil
	}
	lines := make([]string, 0, lastBuffer-firstBuffer+1)
	for i := firstBuffer; i <= lastBuffer; i++ {
		lines = append(lines, contents[i].Text)
	}

	// Trailing newline included so the last line is terminated too (dropDiffPrefix keeps
	// it; the pager-verbatim path needs it added here).
	selected := strings.Join(lines, "\n") + "\n"
	if !self.usingExternalDiff() {
		selected = dropDiffPrefix(selected)
	}

	self.c.LogAction(self.c.Tr.Actions.CopySelectedTextToClipboard)
	return self.c.OS().CopyToClipboard(selected)
}

// usingExternalDiff reports whether the focused main view's diff is produced by an
// external diff command (a pager), in which case the rendered lines may not carry the raw
// +/-/space column (see copySelection).
func (self *MainViewController) usingExternalDiff() bool {
	pagerConfig := self.c.State().GetPagerConfig()
	return pagerConfig.GetExternalDiffCommand() != "" || pagerConfig.GetUseExternalDiffGitConfig()
}

// revealSelectionAfterPrimaryAction re-establishes the focused-main-view selection after a
// primary action (staging or a patch toggle) re-renders the diff. The selection's
// change-line ordinal is read from the source pane (still showing the pre-action diff
// until the queued re-render) and re-applied once the target pane re-renders — so the
// selection lands on the change nearest the one acted on rather than at a stale position.
// sourceViewName and targetViewName are usually the same pane, but staging can move the
// acted-on side to the other pane (passing that pane as the target). The target inherits
// the (collapsed) select mode; a range collapses back to a single line, hunk mode stays
// on to land on the next hunk. advanceBy is forwarded to RevealSelectionAfterStaging: 0 for
// staging/removal (the acted-on lines are consumed), the toggled change-line count for a
// custom-patch toggle (which leaves the diff intact, so the selection must be advanced past
// the just-toggled lines to land on the next stageable hunk).
func revealSelectionAfterPrimaryAction(c *ControllerCommon, sourceViewName string, targetViewName string, firstLineIdx int, advanceBy int) {
	sourceContext := mainContextForViewName(c, sourceViewName)
	targetContext := mainContextForViewName(c, targetViewName)

	sel := sourceContext.DiffSelectState()
	if sel.Mode == context.DiffSelectModeRange {
		sel.Mode = context.DiffSelectModeLine
		sel.RangeIsSticky = false
	}
	*targetContext.DiffSelectState() = *sel
	mode := sel.Mode

	sourceView := sourceContext.GetView()
	targetView := targetContext.GetView()
	c.Helpers().Staging.RevealSelectionAfterStaging(sourceView, targetView, firstLineIdx, advanceBy, func(viewLine int) {
		if mode == context.DiffSelectModeHunk {
			selectDiffHunk(c, targetContext, viewLine)
		} else {
			targetView.CancelRangeSelect()
			showSelectionAtLine(targetView, viewLine, true)
		}
	})
}

// preserveFocusedMainViewSelectionAcrossContentChange keeps the focused main view's
// selection sensible when its diff is about to re-render with *different* content — e.g. a
// commit's diff shrinking after the custom patch built from it is moved into the index, or
// the whole diff changing after an undo. Those mutations don't run through the
// focused-main-view action handlers (which install their own reveal; see
// revealSelectionAfterPrimaryAction), so without this the stale selection — often a large,
// now-meaningless range — would be left painted over the new content.
//
// It is the command-agnostic counterpart of revealSelectionAfterPrimaryAction: rather than
// teaching every mutating command to capture the selection, the focused main view preserves
// it itself, by its change-line ordinal, as the diff re-renders. Call it from a diff side
// panel's render-to-main with the task about to render into the main pane, before
// triggering the render so the restore rides it.
//
// It stands down — leaving the selection exactly as is — unless all of these hold, so it
// neither fights a more precise restore nor disturbs the selection needlessly:
//   - the focused main view is the current context and shows a selection;
//   - no precise restore is already pending (escape / post-stage reveal / context-size),
//     which places the selection more precisely than an ordinal can;
//   - the diff command is changing. A plain refresh re-renders the *same* commit's diff,
//     identical content under an identical command, and there the selection — range and
//     all — must be left alone; only a command change (e.g. a rebase rewriting the commit's
//     hash) means the content, and so the old selection's position, actually moved.
func preserveFocusedMainViewSelectionAcrossContentChange(c *ControllerCommon, mainTask types.UpdateTask) {
	mainContext := c.Contexts().Normal
	if c.Context().Current() != mainContext {
		return
	}
	view := mainContext.GetView()
	if !view.Highlight {
		return
	}
	manager := c.GetViewBufferManagerForView(view)
	if manager == nil || manager.GetRestoreForNextTask() != nil {
		return
	}
	newKey, ok := diffTaskCommandKey(mainTask)
	if !ok || newKey == manager.GetTaskKey() {
		return
	}

	// Anchor on the selection's first line, matching every revealSelectionAfterPrimaryAction
	// caller, so a range collapses to its top and the ordinal lands on the nearest surviving
	// change there.
	first, _ := view.SelectedLineRange()
	revealSelectionAfterPrimaryAction(c, mainContext.GetViewName(), mainContext.GetViewName(), first, 0)
}

// diffTaskCommandKey returns the buffer-manager task key a diff render task will register
// under — the joined command line, as newCmdTask/newPtyTask derive it — or false for a
// non-command task (a placeholder string), which has no command to compare.
func diffTaskCommandKey(task types.UpdateTask) (string, bool) {
	switch t := task.(type) {
	case *types.RunCommandTask:
		return strings.Join(t.Cmd.Args, " "), true
	case *types.RunPtyTask:
		return strings.Join(t.Cmd.Args, " "), true
	}
	return "", false
}

func (self *MainViewController) enter() error {
	if !self.context.GetView().Highlight {
		return nil
	}
	return self.enterForLine(self.context.GetView().SelectedLineIdx())
}

// enterForLine dives into staging/patch-building for the given line, by
// delegating to the side panel beneath the focused main view (the same action
// taken when clicking).
func (self *MainViewController) enterForLine(lineIdx int) error {
	if actions := self.focusedMainViewActions(); actions != nil {
		return actions.OnClick(self.context.GetViewName(), lineIdx)
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
	selectDiffHunk(self.c, self.context, changeViewLine)
}

// selectDiffHunk is the body of selectHunkAround, as a free function so the focus
// entry point (showInitialDiffSelection) can establish a hunk selection too.
func selectDiffHunk(c *ControllerCommon, mainContext *context.MainContext, changeViewLine int) {
	view := mainContext.GetView()
	start, end, ok := c.Helpers().Staging.ChangeBlockBounds(view, changeViewLine)
	if !ok {
		mainContext.DiffSelectState().Mode = context.DiffSelectModeLine
		view.CancelRangeSelect()
		showSelectionAtLine(view, changeViewLine, true)
		return
	}
	view.SetRangeSelectStart(end)
	showSelectionAtLine(view, start, true)
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

// mainContextForViewName maps a focused main view's view name (as passed to the
// side-panel handlers) to its main context — the secondary pane for the secondary
// view name, the primary pane otherwise.
func mainContextForViewName(c *ControllerCommon, viewName string) *context.MainContext {
	if viewName == c.Contexts().NormalSecondary.GetViewName() {
		return c.Contexts().NormalSecondary
	}
	return c.Contexts().Normal
}

// focusedMainViewContextForViewName is mainContextForViewName as a types.Context, for
// callers that only need the interface.
func focusedMainViewContextForViewName(c *ControllerCommon, viewName string) types.Context {
	return mainContextForViewName(c, viewName)
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
	return self.selectClickedDiffLine(opts)
}

func (self *MainViewController) editClickedLine(opts gocui.ViewMouseBindingOpts) error {
	return self.editDiffLine(opts.Y)
}

func (self *MainViewController) onClickInOtherViewOfMainViewPair(opts gocui.ViewMouseBindingOpts) error {
	// Carry the select mode over from the pane we're leaving, so clicking into the
	// other pane keeps hunk mode even the first time we enter it — without this its
	// mode would still be the default single line until it had been focused at least
	// once (e.g. by tabbing to it). selectClickedDiffLine then preserves or collapses
	// that mode depending on where we clicked.
	*self.context.DiffSelectState() = *self.otherContext.DiffSelectState()
	self.c.Context().Push(self.context, types.OnFocusOpts{})
	return self.selectClickedDiffLine(opts)
}

// selectClickedDiffLine sets the focused main view's selection from a click at view
// line opts.Y. In hunk mode a click on a change line keeps hunk mode and selects that
// whole block, so clicking from hunk to hunk stays ready to stage; a click on context
// drops to a single line, as does any click when we weren't in hunk mode — the click
// points at the line precisely (e.g. to edit it). A double-click additionally dives
// into staging/patch-building for that line.
func (self *MainViewController) selectClickedDiffLine(opts gocui.ViewMouseBindingOpts) error {
	if !self.isDiffView() {
		return nil
	}
	view := self.context.GetView()
	// Remember where the click landed so a drag that follows anchors its range there,
	// even when this click selects a whole hunk (whose anchor is the block's far end).
	self.context.SetDragAnchorViewLine(opts.Y)
	if self.sel().Mode == context.DiffSelectModeHunk && self.c.Helpers().Staging.IsChangeLine(view, opts.Y) {
		self.selectHunkAround(opts.Y)
	} else {
		resetDiffSelectMode(self.context)
		showSelectionAtLine(view, opts.Y, false)
	}
	if opts.IsDoubleClick {
		return self.enterForLine(opts.Y)
	}
	return nil
}

// onDragInFocusedView extends a range selection as the mouse is dragged after a click,
// anchored at the line the click landed on (DragAnchorViewLine) — not wherever the click
// left the selection, since a click can select a whole hunk whose far end would otherwise
// become the anchor. Dragging turns hunk mode off: you get a plain range from the clicked
// line to the line under the cursor, which gocui has already moved here. A no-op when the
// view holds no selectable diff.
func (self *MainViewController) onDragInFocusedView(gocui.ViewMouseBindingOpts) error {
	view := self.context.GetView()
	if !self.isDiffView() || !view.Highlight {
		return nil
	}
	sel := self.sel()
	sel.Mode = context.DiffSelectModeRange
	sel.RangeIsSticky = false
	sel.UserEnabledHunkMode = false
	view.SetRangeSelectStart(self.context.DragAnchorViewLine())
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
