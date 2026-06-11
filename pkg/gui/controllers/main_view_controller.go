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
	// When a selection is shown, we surface the bindings that act on it
	// (enter to dive into staging, e to edit the selected line, G to open the
	// line in the branch's pull request, escape to hide the selection).
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
			Handler:         self.toggleSelection,
			Description:     self.c.Tr.ToggleSelectionInFocusedMainView,
			DisplayOnScreen: !selectionShown,
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
	if self.otherContext.GetView().Visible {
		self.c.Context().Push(self.otherContext, types.OnFocusOpts{})
	}

	return nil
}

func (self *MainViewController) escape() error {
	v := self.context.GetView()
	if v.Highlight {
		v.Highlight = false
		return nil
	}
	self.c.Context().Pop()
	return nil
}

func (self *MainViewController) toggleSelection() error {
	v := self.context.GetView()
	if v.Highlight {
		v.Highlight = false
		return nil
	}
	// Start the selection in the middle of the visible area.
	showSelectionAtLine(v, v.OriginY()+v.InnerHeight()/2, false)
	return nil
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

// navigate jumps the focused main view by file or change block (hunk), using find to
// locate the target row from the current anchor. The anchor is the selected line if a
// selection is showing, otherwise the top visible line. With a selection showing we
// move it to the target and scroll it into view, like the staging view; with none we
// stay in scroll mode, bringing the target to the top without selecting anything.
func (self *MainViewController) navigate(find func(*gocui.View, int, bool) (int, bool), forward bool) error {
	v := self.context.GetView()
	showSelection := v.Highlight
	anchor := v.OriginY()
	if showSelection {
		anchor = v.SelectedLineIdx()
	}

	target, ok := find(v, anchor, forward)
	if !ok {
		return nil
	}

	if showSelection {
		showSelectionAtLine(v, target, true)
	} else {
		v.SetOrigin(0, target)
	}
	return nil
}

func (self *MainViewController) nextChangeBlock() error {
	return self.navigate(self.c.Helpers().Staging.AdjacentChangeBlock, true)
}

func (self *MainViewController) prevChangeBlock() error {
	return self.navigate(self.c.Helpers().Staging.AdjacentChangeBlock, false)
}

func (self *MainViewController) nextFile() error {
	return self.navigate(self.c.Helpers().Staging.AdjacentFile, true)
}

func (self *MainViewController) prevFile() error {
	return self.navigate(self.c.Helpers().Staging.AdjacentFile, false)
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
	// A click points at a line, so it sets the selection there; a double-click
	// additionally dives into staging/patch-building for that line.
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
