package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

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
	v.Highlight = true
	v.HighlightInactive = false
	lineIdx := v.OriginY() + v.InnerHeight()/2
	lineIdx = lo.Clamp(lineIdx, 0, v.ViewLinesHeight()-1)
	v.FocusPoint(0, lineIdx, false)
	return nil
}

func (self *MainViewController) enter() error {
	if !self.context.GetView().Highlight {
		return nil
	}
	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext != nil && sidePanelContext.GetOnClickFocusedMainView() != nil {
		return sidePanelContext.GetOnClickFocusedMainView()(
			self.context.GetViewName(), self.context.GetView().SelectedLineIdx())
	}
	return nil
}

func (self *MainViewController) editLine() error {
	if !self.context.GetView().Highlight {
		return nil
	}
	// Figure out the clicked file and line the same way entering staging does.
	path, lineNumber, ok := self.c.Helpers().Staging.GetFileAndLineForClickedDiffLine(
		self.context.GetViewName(), self.context.GetView().SelectedLineIdx())
	if !ok {
		return nil
	}
	lineNumber = self.c.Helpers().Diff.AdjustLineNumber(path, lineNumber, self.context.GetViewName())
	return self.c.Helpers().Files.EditFileAtLine(path, lineNumber)
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
	path, lineNumber, ok := self.c.Helpers().Staging.GetFileAndLineForClickedDiffLine(
		self.context.GetViewName(), self.context.GetView().SelectedLineIdx())
	if !ok {
		return nil
	}

	relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), path)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.OpenPullRequest)
	return self.c.OS().OpenLink(
		githubPullRequestLineURL(pr.Url, commitSha, filepath.ToSlash(relativePath), lineNumber))
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

func (self *MainViewController) onClickInAlreadyFocusedView(opts gocui.ViewMouseBindingOpts) error {
	if self.context.GetView().Highlight && !opts.IsDoubleClick {
		return nil
	}

	sidePanelContext := self.c.Context().NextInStack(self.context)
	if sidePanelContext != nil && sidePanelContext.GetOnClickFocusedMainView() != nil {
		return sidePanelContext.GetOnClickFocusedMainView()(self.context.GetViewName(), opts.Y)
	}
	return nil
}

func (self *MainViewController) onClickInOtherViewOfMainViewPair(opts gocui.ViewMouseBindingOpts) error {
	self.c.Context().Push(self.context, types.OnFocusOpts{
		ClickedWindowName:  self.context.GetWindowName(),
		ClickedViewLineIdx: opts.Y,
	})

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
