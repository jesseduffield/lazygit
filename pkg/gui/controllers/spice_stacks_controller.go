package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type SpiceStacksController struct {
	baseController
	*ListControllerTrait[*models.SpiceStackItem]
	c *ControllerCommon
}

var _ types.IController = &SpiceStacksController{}

func NewSpiceStacksController(
	c *ControllerCommon,
) *SpiceStacksController {
	return &SpiceStacksController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().SpiceStacks,
			c.Contexts().SpiceStacks.GetSelected,
			c.Contexts().SpiceStacks.GetSelectedItems,
		),
		c: c,
	}
}

func (self *SpiceStacksController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.handleEnter,
			Description: self.c.Tr.ViewCommits,
			DescriptionFunc: func() string {
				item := self.context().GetSelected()
				if item == nil && self.c.Git().Spice != nil && self.c.Git().Spice.IsAvailable() && !self.c.Git().Spice.IsInitialized() {
					return self.c.Tr.SpiceInitialize
				}
				return self.c.Tr.ViewCommits
			},
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.SquashDown),
			Handler:           self.withItem(self.commitSquash),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.Squash,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MarkCommitAsFixup),
			Handler:           self.withItem(self.commitMarkFixup),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.Fixup,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.RenameCommit),
			Handler:           self.withItem(self.commitReword),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.Reword,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.commitDrop),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.Drop,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Edit),
			Handler:           self.withItem(self.commitEdit),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.Edit,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.AmendToCommit),
			Handler:           self.withItem(self.commitAmend),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.Amend,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:           self.withItem(self.commitReset),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.ViewResetOptions,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CherryPickCopy),
			Handler:           self.withItem(self.commitCopy),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.CherryPickCopy,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CopyCommitAttributeToClipboard),
			Handler:           self.withItem(self.copyCommitAttribute),
			GetDisabledReason: self.require(self.singleItemSelected(), self.commitSelected()),
			Description:       self.c.Tr.CopyCommitAttributeToClipboard,
			OpensMenu:         true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.press),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MoveDownCommit),
			Handler:           self.withItem(self.moveDown),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.MoveDownCommit, // fallback for cheatsheet
			DescriptionFunc: func() string {
				item := self.context().GetSelected()
				if item != nil && !item.IsCommit {
					return self.c.Tr.SpiceMoveBranchDown
				}
				return self.c.Tr.MoveDownCommit
			},
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MoveUpCommit),
			Handler:           self.withItem(self.moveUp),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.MoveUpCommit, // fallback for cheatsheet
			DescriptionFunc: func() string {
				item := self.context().GetSelected()
				if item != nil && !item.IsCommit {
					return self.c.Tr.SpiceMoveBranchUp
				}
				return self.c.Tr.MoveUpCommit
			},
		},
		{
			Key:             opts.GetKey("S"),
			Handler:         self.openStackOperationsMenu,
			Description:     self.c.Tr.SpiceStackOperations,
			OpensMenu:       true,
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("G"),
			Handler:         self.openNavigationMenu,
			Description:     self.c.Tr.SpiceStackNavigation,
			OpensMenu:       true,
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey("V"),
			Handler:     self.openLogFormatMenu,
			Description: self.c.Tr.ToggleSpiceLogFormat,
			Tooltip:     self.c.Tr.ToggleSpiceLogFormatTooltip,
			OpensMenu:   true,
		},
		{
			Key:               opts.GetKey("n"),
			Handler:           self.withItem(self.newBranchOnSelected),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.SpiceCreateBranch,
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey("N"),
			Handler:         self.newCommit,
			Description:     self.c.Tr.SpiceCreateCommit,
			DisplayOnScreen: true,
		},
	}

	return bindings
}

func (self *SpiceStacksController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {}
}

func (self *SpiceStacksController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			var title string
			item := self.context().GetSelected()

			if item == nil {
				// Check if git-spice needs initialization
				if self.c.Git().Spice != nil && self.c.Git().Spice.IsAvailable() && !self.c.Git().Spice.IsInitialized() {
					task = types.NewRenderStringTask(self.c.Tr.SpiceNotInitialized)
				} else {
					task = types.NewRenderStringTask(self.c.Tr.SpiceNoStacks)
				}
				title = self.c.Tr.SpiceStacksTitle
			} else if item.IsCommit {
				// Show commit diff/patch
				cmdObj := self.c.Git().Commit.ShowCmdObj(item.CommitSha, nil)
				task = types.NewRunPtyTask(cmdObj.GetCmd())
				title = self.c.Tr.SpicePatch
			} else {
				// Show branch commit log
				cmdObj := self.c.Git().Branch.GetGraphCmdObj(item.FullRefName())
				task = types.NewRunPtyTask(cmdObj.GetCmd())
				title = self.c.Tr.LogTitle
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title:    title,
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Task:     task,
				},
			})
		})
	}
}

func (self *SpiceStacksController) enter(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return self.viewCommitFiles(item)
	}
	return self.viewBranchCommits(item)
}

func (self *SpiceStacksController) handleEnter() error {
	item := self.context().GetSelected()
	if item != nil {
		return self.enter(item)
	}

	// No item selected - check if we should offer initialization
	if self.c.Git().Spice != nil && self.c.Git().Spice.IsAvailable() && !self.c.Git().Spice.IsInitialized() {
		return self.initializeSpice()
	}

	return nil
}

func (self *SpiceStacksController) initializeSpice() error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.SpiceSelectTrunkBranch,
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetBranchNameSuggestionsFunc(),
		HandleConfirm: func(trunkBranch string) error {
			return self.runInitialization(trunkBranch)
		},
	})
	return nil
}

func (self *SpiceStacksController) runInitialization(trunk string) error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceInitializingStatus, func(gocui.Task) error {
		if err := self.c.Git().Spice.Init(trunk); err != nil {
			return err
		}

		// Clear the cache so IsInitialized() re-checks
		self.c.Git().Spice.ClearInitializedCache()

		// Refresh all relevant views
		self.c.Refresh(types.RefreshOptions{
			Mode:  types.ASYNC,
			Scope: []types.RefreshableView{types.SPICE_STACKS, types.BRANCHES},
		})

		return nil
	})
}

func (self *SpiceStacksController) press(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return self.commitCheckout(item)
	}
	return self.checkout(item)
}

func (self *SpiceStacksController) viewBranchCommits(item *models.SpiceStackItem) error {
	branch := self.findBranchByName(item.Name)
	if branch == nil {
		return errors.New(self.c.Tr.SpiceBranchNotFound)
	}

	return self.c.Helpers().SubCommits.ViewSubCommits(helpers.ViewSubCommitsOpts{
		Ref:             branch,
		TitleRef:        branch.RefName(),
		Context:         self.context(),
		ShowBranchHeads: false,
	})
}

func (self *SpiceStacksController) findBranchByName(name string) *models.Branch {
	for _, branch := range self.c.Model().Branches {
		if branch.Name == name {
			return branch
		}
	}
	return nil
}

func (self *SpiceStacksController) viewCommitFiles(item *models.SpiceStackItem) error {
	commit, _ := self.findCommitByHash(item.CommitSha)
	if commit == nil {
		return errors.New(self.c.Tr.SpiceCommitNotFound)
	}

	commitFilesContext := self.c.Contexts().CommitFiles
	commitFilesContext.ReInit(commit, nil)
	commitFilesContext.SetSelection(0)
	commitFilesContext.SetCanRebase(false)
	commitFilesContext.SetParentContext(self.context())
	commitFilesContext.SetWindowName(self.context().GetWindowName())
	commitFilesContext.ClearSearchString()

	self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	})

	self.c.Context().Push(commitFilesContext, types.OnFocusOpts{})
	return nil
}

// findCommitByHash searches for a commit in the model's commits list by SHA.
// Uses prefix matching since SpiceStackItem stores short (7-char) hashes.
func (self *SpiceStacksController) findCommitByHash(sha string) (*models.Commit, int) {
	for idx, commit := range self.c.Model().Commits {
		if strings.HasPrefix(commit.Hash(), sha) {
			return commit, idx
		}
	}
	return nil, -1
}

// findCommitOrError is a helper that returns an error if the commit is not found.
func (self *SpiceStacksController) findCommitOrError(sha string) (*models.Commit, int, error) {
	commit, idx := self.findCommitByHash(sha)
	if commit == nil {
		return nil, -1, errors.New(self.c.Tr.SpiceCommitNotFound)
	}
	return commit, idx, nil
}

// refreshSpiceStacksOnly refreshes only the stacks view.
func (self *SpiceStacksController) refreshSpiceStacksOnly() {
	self.c.Refresh(types.RefreshOptions{
		Mode:  types.ASYNC,
		Scope: []types.RefreshableView{types.SPICE_STACKS},
	})
}

// refreshAllViews refreshes branches, commits, files, and stacks views.
func (self *SpiceStacksController) refreshAllViews() {
	self.c.Refresh(types.RefreshOptions{
		Mode:  types.ASYNC,
		Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS},
	})
}

func (self *SpiceStacksController) commitSquash(item *models.SpiceStackItem) error {
	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Squash,
		Prompt: self.c.Tr.SureSquashThisCommit,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.SquashingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.SquashCommitDown)
				err := self.c.Git().Rebase.InteractiveRebase(self.c.Model().Commits, commitIdx, commitIdx, todo.Squash)
				return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
	return nil
}

func (self *SpiceStacksController) commitMarkFixup(item *models.SpiceStackItem) error {
	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Fixup,
		Prompt: self.c.Tr.SureFixupThisCommit,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.FixingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.FixupCommit)
				err := self.c.Git().Rebase.InteractiveRebase(self.c.Model().Commits, commitIdx, commitIdx, todo.Fixup)
				return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
	return nil
}

func (self *SpiceStacksController) commitReword(item *models.SpiceStackItem) error {
	commit, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	commitMessage, err := self.c.Git().Commit.GetCommitMessage(commit.Hash())
	if err != nil {
		return err
	}

	self.c.Helpers().Commits.OpenCommitMessagePanel(
		&helpers.OpenCommitMessagePanelOpts{
			CommitIndex:      commitIdx,
			InitialMessage:   commitMessage,
			SummaryTitle:     self.c.Tr.Actions.RewordCommit,
			DescriptionTitle: self.c.Tr.CommitDescriptionTitle,
			PreserveMessage:  false,
			OnConfirm:        self.handleReword,
		},
	)
	return nil
}

func (self *SpiceStacksController) handleReword(summary string, description string) error {
	item := self.context().GetSelected()
	if item == nil {
		return nil
	}

	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	// Check if this is the head commit
	if commitIdx == 0 {
		return self.c.Helpers().GPG.WithGpgHandling(
			self.c.Git().Commit.RewordLastCommit(summary, description),
			git_commands.CommitGpgSign,
			self.c.Tr.RewordingStatus, nil, nil)
	}

	return self.c.WithWaitingStatus(self.c.Tr.RewordingStatus, func(gocui.Task) error {
		err := self.c.Git().Rebase.RewordCommit(self.c.Model().Commits, commitIdx, summary, description)
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
	})
}

func (self *SpiceStacksController) commitDrop(item *models.SpiceStackItem) error {
	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DropCommitTitle,
		Prompt: self.c.Tr.DropCommitPrompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DroppingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.DropCommit)
				err := self.c.Git().Rebase.InteractiveRebase(self.c.Model().Commits, commitIdx, commitIdx, todo.Drop)
				return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
	return nil
}

func (self *SpiceStacksController) commitEdit(item *models.SpiceStackItem) error {
	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.EditCommit)
	err = self.c.Git().Rebase.InteractiveRebase(self.c.Model().Commits, commitIdx, commitIdx, todo.Edit)
	return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
}

func (self *SpiceStacksController) commitAmend(item *models.SpiceStackItem) error {
	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	// If it's the head commit, use the amend helper
	if commitIdx == 0 {
		return self.c.Helpers().AmendHelper.AmendHead()
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AmendCommitTitle,
		Prompt: self.c.Tr.AmendCommitPrompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.AmendCommit)
				err := self.c.Git().Rebase.AmendTo(self.c.Model().Commits, commitIdx)
				return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
	return nil
}

func (self *SpiceStacksController) commitCheckout(item *models.SpiceStackItem) error {
	commit, _, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	return self.c.Helpers().Refs.CreateCheckoutMenu(commit)
}

func (self *SpiceStacksController) commitReset(item *models.SpiceStackItem) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(item.CommitSha, item.CommitSha)
}

func (self *SpiceStacksController) commitCopy(item *models.SpiceStackItem) error {
	commit, _, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	// Directly manipulate cherry-pick data since CopyRange uses context selection
	// which doesn't apply to SpiceStacks context
	cherryPicking := self.c.Modes().CherryPicking
	if cherryPicking.ContextKey != string(self.c.Contexts().SpiceStacks.GetKey()) {
		cherryPicking.ContextKey = string(self.c.Contexts().SpiceStacks.GetKey())
		cherryPicking.CherryPickedCommits = nil
	}

	// Toggle: if already copied, remove it; otherwise add it
	if cherryPicking.SelectedHashSet().Includes(commit.Hash()) {
		cherryPicking.Remove(commit, self.c.Model().Commits)
	} else {
		cherryPicking.Add(commit, self.c.Model().Commits)
	}

	return nil
}

func (self *SpiceStacksController) copyCommitAttribute(item *models.SpiceStackItem) error {
	commit, _, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.CopyCommitAttributeToClipboard,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.CommitHash,
				OnPress: func() error {
					return self.copyCommitHashToClipboard(commit)
				},
			},
			{
				Label: self.c.Tr.CommitSubject,
				Key:   's',
				OnPress: func() error {
					return self.copyCommitSubjectToClipboard(commit)
				},
			},
			{
				Label: self.c.Tr.CommitMessage,
				Key:   'm',
				OnPress: func() error {
					return self.copyCommitMessageToClipboard(commit)
				},
			},
			{
				Label: self.c.Tr.CommitURL,
				Key:   'u',
				OnPress: func() error {
					return self.copyCommitURLToClipboard(commit)
				},
			},
			{
				Label: self.c.Tr.CommitDiff,
				Key:   'd',
				OnPress: func() error {
					return self.copyCommitDiffToClipboard(commit)
				},
			},
			{
				Label: self.c.Tr.CommitAuthor,
				Key:   'a',
				OnPress: func() error {
					return self.copyAuthorToClipboard(commit)
				},
			},
		},
	})
}

func (self *SpiceStacksController) copyCommitHashToClipboard(commit *models.Commit) error {
	self.c.LogAction(self.c.Tr.Actions.CopyCommitHashToClipboard)
	if err := self.c.OS().CopyToClipboard(commit.Hash()); err != nil {
		return err
	}
	self.c.Toast(fmt.Sprintf("'%s' %s", commit.Hash(), self.c.Tr.CopiedToClipboard))
	return nil
}

func (self *SpiceStacksController) copyCommitSubjectToClipboard(commit *models.Commit) error {
	message, err := self.c.Git().Commit.GetCommitMessage(commit.Hash())
	if err != nil {
		return err
	}
	subject, _ := self.c.Helpers().Commits.SplitCommitMessageAndDescription(message)
	self.c.LogAction(self.c.Tr.Actions.CopyCommitSubjectToClipboard)
	if err := self.c.OS().CopyToClipboard(subject); err != nil {
		return err
	}
	self.c.Toast(fmt.Sprintf("'%s' %s", subject, self.c.Tr.CopiedToClipboard))
	return nil
}

func (self *SpiceStacksController) copyCommitMessageToClipboard(commit *models.Commit) error {
	message, err := self.c.Git().Commit.GetCommitMessage(commit.Hash())
	if err != nil {
		return err
	}
	self.c.LogAction(self.c.Tr.Actions.CopyCommitMessageToClipboard)
	if err := self.c.OS().CopyToClipboard(message); err != nil {
		return err
	}
	self.c.Toast(self.c.Tr.CommitMessageCopiedToClipboard)
	return nil
}

func (self *SpiceStacksController) copyCommitURLToClipboard(commit *models.Commit) error {
	url, err := self.c.Helpers().Host.GetCommitURL(commit.Hash())
	if err != nil {
		return err
	}
	self.c.LogAction(self.c.Tr.Actions.CopyCommitURLToClipboard)
	if err := self.c.OS().CopyToClipboard(url); err != nil {
		return err
	}
	self.c.Toast(self.c.Tr.CommitURLCopiedToClipboard)
	return nil
}

func (self *SpiceStacksController) copyCommitDiffToClipboard(commit *models.Commit) error {
	diff, err := self.c.Git().Commit.GetCommitDiff(commit.Hash())
	if err != nil {
		return err
	}
	self.c.LogAction(self.c.Tr.Actions.CopyCommitDiffToClipboard)
	if err := self.c.OS().CopyToClipboard(diff); err != nil {
		return err
	}
	self.c.Toast(self.c.Tr.CommitDiffCopiedToClipboard)
	return nil
}

func (self *SpiceStacksController) copyAuthorToClipboard(commit *models.Commit) error {
	author, err := self.c.Git().Commit.GetCommitAuthor(commit.Hash())
	if err != nil {
		return err
	}
	authorStr := fmt.Sprintf("%s <%s>", author.Name, author.Email)
	self.c.LogAction(self.c.Tr.Actions.CopyCommitAuthorToClipboard)
	if err := self.c.OS().CopyToClipboard(authorStr); err != nil {
		return err
	}
	self.c.Toast(fmt.Sprintf("'%s' %s", authorStr, self.c.Tr.CopiedToClipboard))
	return nil
}

func (self *SpiceStacksController) checkout(item *models.SpiceStackItem) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	if err := self.c.Git().Branch.Checkout(item.Name, git_commands.CheckoutOptions{Force: false}); err != nil {
		return err
	}
	self.refreshAllViews()
	return nil
}

func (self *SpiceStacksController) restack(item *models.SpiceStackItem) error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceRestackingStatus, func(task gocui.Task) error {
		if err := self.c.Git().Spice.Restack(item.Name); err != nil {
			return err
		}
		self.refreshSpiceStacksOnly()
		return nil
	})
}

func (self *SpiceStacksController) restackAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceRestackingStatus, func(task gocui.Task) error {
		if err := self.c.Git().Spice.Restack(""); err != nil {
			return err
		}
		self.refreshSpiceStacksOnly()
		return nil
	})
}

func (self *SpiceStacksController) submit(item *models.SpiceStackItem, opts git_commands.SubmitOpts) error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceSubmittingStatus, func(task gocui.Task) error {
		if err := self.c.Git().Spice.Submit(item.Name, opts); err != nil {
			return err
		}
		self.refreshSpiceStacksOnly()
		return nil
	})
}

func (self *SpiceStacksController) submitAll(opts git_commands.SubmitOpts) error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceSubmittingStatus, func(task gocui.Task) error {
		if err := self.c.Git().Spice.Submit("", opts); err != nil {
			return err
		}
		self.refreshSpiceStacksOnly()
		return nil
	})
}

func (self *SpiceStacksController) navigateUp() error {
	if err := self.c.Git().Spice.NavigateUp(); err != nil {
		return err
	}
	self.refreshAllViews()
	return nil
}

func (self *SpiceStacksController) navigateDown() error {
	if err := self.c.Git().Spice.NavigateDown(); err != nil {
		return err
	}
	self.refreshAllViews()
	return nil
}

func (self *SpiceStacksController) navigateTop() error {
	if err := self.c.Git().Spice.NavigateTop(); err != nil {
		return err
	}
	self.refreshAllViews()
	return nil
}

func (self *SpiceStacksController) navigateBottom() error {
	if err := self.c.Git().Spice.NavigateBottom(); err != nil {
		return err
	}
	self.refreshAllViews()
	return nil
}

func (self *SpiceStacksController) newBranch() error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.SpiceBranchNamePrompt,
		HandleConfirm: func(branchName string) error {
			if err := self.c.Git().Spice.CreateBranch(branchName, ""); err != nil {
				return err
			}
			self.refreshAllViews()
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) delete(item *models.SpiceStackItem) error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.SpiceDeleteConfirmTitle,
		Prompt: self.c.Tr.SpiceDeleteConfirmPrompt,
		HandleConfirm: func() error {
			if err := self.c.Git().Spice.DeleteBranch(item.Name); err != nil {
				return err
			}
			self.c.Refresh(types.RefreshOptions{
				Mode:  types.ASYNC,
				Scope: []types.RefreshableView{types.BRANCHES, types.SPICE_STACKS},
			})
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) moveBranchUp(item *models.SpiceStackItem) error {
	if err := self.c.Git().Spice.MoveBranchUp(item.Name); err != nil {
		return err
	}
	self.refreshSpiceStacksOnly()
	return nil
}

func (self *SpiceStacksController) moveBranchDown(item *models.SpiceStackItem) error {
	if err := self.c.Git().Spice.MoveBranchDown(item.Name); err != nil {
		return err
	}
	self.refreshSpiceStacksOnly()
	return nil
}

// getBranchNameFromItem returns the branch name for an item.
// For branch items, it returns the branch name directly.
// For commit items, it finds the parent branch by searching upward.
func (self *SpiceStacksController) getBranchNameFromItem(item *models.SpiceStackItem) string {
	if !item.IsCommit {
		return item.Name
	}
	// Find parent branch for commit items
	items := self.context().GetItems()
	idx := self.context().GetSelectedLineIdx()
	for i := idx - 1; i >= 0; i-- {
		if !items[i].IsCommit {
			return items[i].Name
		}
	}
	return ""
}

func (self *SpiceStacksController) newBranchOnSelected(item *models.SpiceStackItem) error {
	targetBranch := self.getBranchNameFromItem(item)
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.SpiceBranchNamePrompt,
		HandleConfirm: func(branchName string) error {
			if err := self.c.Git().Spice.CreateBranch(branchName, targetBranch); err != nil {
				return err
			}
			self.refreshAllViews()
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) newCommit() error {
	if err := self.c.Git().Spice.CreateCommit(); err != nil {
		return err
	}
	self.refreshAllViews()
	return nil
}

func (self *SpiceStacksController) commitMoveDown(item *models.SpiceStackItem) error {
	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	return self.c.WithWaitingStatusSync(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)
		err := self.c.Git().Rebase.MoveCommitsDown(self.c.Model().Commits, commitIdx, commitIdx)
		if err == nil {
			self.context().MoveSelection(1)
		}
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err, types.RefreshOptions{Mode: types.SYNC})
	})
}

func (self *SpiceStacksController) commitMoveUp(item *models.SpiceStackItem) error {
	_, commitIdx, err := self.findCommitOrError(item.CommitSha)
	if err != nil {
		return err
	}

	return self.c.WithWaitingStatusSync(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		err := self.c.Git().Rebase.MoveCommitsUp(self.c.Model().Commits, commitIdx, commitIdx)
		if err == nil {
			self.context().MoveSelection(-1)
		}
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err, types.RefreshOptions{Mode: types.SYNC})
	})
}

// Unified move handlers that dispatch based on selection type
func (self *SpiceStacksController) moveDown(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return self.commitMoveDown(item)
	}
	return self.moveBranchDown(item)
}

func (self *SpiceStacksController) moveUp(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return self.commitMoveUp(item)
	}
	return self.moveBranchUp(item)
}

func (self *SpiceStacksController) openStackOperationsMenu() error {
	item := self.context().GetSelected()

	menuItems := []*types.MenuItem{
		{
			Label:          self.c.Tr.SpiceRestackBranch,
			Key:            'r',
			OnPress:        func() error { return self.restack(item) },
			DisabledReason: self.branchSelectedReason(item),
		},
		{
			Label:   self.c.Tr.SpiceRestackAll,
			Key:     'R',
			OnPress: func() error { return self.restackAll() },
		},
		{
			Label:          self.c.Tr.SpiceSubmitBranch,
			Key:            's',
			OnPress:        func() error { return self.submit(item, git_commands.SubmitOpts{}) },
			DisabledReason: self.branchSelectedReason(item),
		},
		{
			Label:          self.c.Tr.SpiceSubmitBranchOptions,
			Key:            'o',
			OnPress:        func() error { return self.openSubmitBranchMenu(item) },
			OpensMenu:      true,
			DisabledReason: self.branchSelectedReason(item),
		},
		{
			Label:   self.c.Tr.SpiceSubmitAll,
			Key:     'S',
			OnPress: func() error { return self.submitAll(git_commands.SubmitOpts{}) },
		},
		{
			Label:     self.c.Tr.SpiceSubmitAllOptions,
			Key:       'O',
			OnPress:   func() error { return self.openSubmitAllMenu() },
			OpensMenu: true,
		},
		{
			Label:   self.c.Tr.SpiceNewBranch,
			Key:     'c',
			OnPress: func() error { return self.newBranch() },
		},
		{
			Label:          self.c.Tr.SpiceDeleteBranch,
			Key:            'd',
			OnPress:        func() error { return self.delete(item) },
			DisabledReason: self.branchSelectedReason(item),
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SpiceStackOperationsMenuTitle,
		Items: menuItems,
	})
}

func (self *SpiceStacksController) openNavigationMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label:   self.c.Tr.SpiceNavigateUp,
			Key:     'u',
			Tooltip: self.c.Tr.SpiceNavigateUpTooltip,
			OnPress: func() error { return self.navigateUp() },
		},
		{
			Label:   self.c.Tr.SpiceNavigateDown,
			Key:     'd',
			Tooltip: self.c.Tr.SpiceNavigateDownTooltip,
			OnPress: func() error { return self.navigateDown() },
		},
		{
			Label:   self.c.Tr.SpiceNavigateTop,
			Key:     't',
			Tooltip: self.c.Tr.SpiceNavigateTopTooltip,
			OnPress: func() error { return self.navigateTop() },
		},
		{
			Label:   self.c.Tr.SpiceNavigateBottom,
			Key:     'b',
			Tooltip: self.c.Tr.SpiceNavigateBottomTooltip,
			OnPress: func() error { return self.navigateBottom() },
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SpiceStackNavigationMenuTitle,
		Items: menuItems,
	})
}

func (self *SpiceStacksController) openSubmitBranchMenu(item *models.SpiceStackItem) error {
	menuItems := []*types.MenuItem{
		{
			Label:   self.c.Tr.SpiceNoPublish,
			Key:     'n',
			Tooltip: self.c.Tr.SpiceNoPublishTooltip,
			OnPress: func() error { return self.submit(item, git_commands.SubmitOpts{NoPublish: true}) },
		},
		{
			Label:   self.c.Tr.SpiceUpdateOnly,
			Key:     'u',
			Tooltip: self.c.Tr.SpiceUpdateOnlyTooltip,
			OnPress: func() error { return self.submit(item, git_commands.SubmitOpts{UpdateOnly: true}) },
		},
		{
			Label:   self.c.Tr.SpiceSubmitDefault,
			Key:     's',
			OnPress: func() error { return self.submit(item, git_commands.SubmitOpts{}) },
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SpiceSubmitBranchOptionsMenuTitle,
		Items: menuItems,
	})
}

func (self *SpiceStacksController) openSubmitAllMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label:   self.c.Tr.SpiceNoPublish,
			Key:     'n',
			Tooltip: self.c.Tr.SpiceNoPublishTooltip,
			OnPress: func() error { return self.submitAll(git_commands.SubmitOpts{NoPublish: true}) },
		},
		{
			Label:   self.c.Tr.SpiceUpdateOnly,
			Key:     'u',
			Tooltip: self.c.Tr.SpiceUpdateOnlyTooltip,
			OnPress: func() error { return self.submitAll(git_commands.SubmitOpts{UpdateOnly: true}) },
		},
		{
			Label:   self.c.Tr.SpiceSubmitAllDefault,
			Key:     's',
			OnPress: func() error { return self.submitAll(git_commands.SubmitOpts{}) },
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SpiceSubmitAllOptionsMenuTitle,
		Items: menuItems,
	})
}

func (self *SpiceStacksController) branchSelectedReason(item *models.SpiceStackItem) *types.DisabledReason {
	if item == nil || item.IsCommit {
		return &types.DisabledReason{Text: self.c.Tr.SpiceBranchOnly}
	}
	return nil
}

func (self *SpiceStacksController) openLogFormatMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SpiceLogFormatMenuTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.SpiceLogFormatShort,
				Key:   's',
				OnPress: func() error {
					return self.setLogFormat("short")
				},
			},
			{
				Label: self.c.Tr.SpiceLogFormatLong,
				Key:   'l',
				OnPress: func() error {
					return self.setLogFormat("long")
				},
			},
			{
				Label: self.c.Tr.SpiceLogFormatDefault,
				Key:   'd',
				OnPress: func() error {
					return self.setLogFormat("") // empty = use config default
				},
			},
		},
	})
}

func (self *SpiceStacksController) setLogFormat(format string) error {
	self.c.GetAppState().Spice.LogFormat = format
	self.c.SaveAppStateAndLogError()
	self.refreshSpiceStacksOnly()
	return nil
}

func (self *SpiceStacksController) context() *context.SpiceStacksContext {
	return self.c.Contexts().SpiceStacks
}

func (self *SpiceStacksController) withItem(f func(item *models.SpiceStackItem) error) func() error {
	return func() error {
		item := self.context().GetSelected()
		if item == nil {
			return nil
		}
		return f(item)
	}
}

func (self *SpiceStacksController) singleItemSelected() func() *types.DisabledReason {
	return func() *types.DisabledReason {
		if self.context().GetSelected() == nil {
			return &types.DisabledReason{Text: self.c.Tr.SpiceNoItemSelected}
		}
		return nil
	}
}

func (self *SpiceStacksController) commitSelected() func() *types.DisabledReason {
	return func() *types.DisabledReason {
		item := self.context().GetSelected()
		if item == nil || !item.IsCommit {
			return &types.DisabledReason{Text: self.c.Tr.SpiceCommitOnly}
		}
		return nil
	}
}

func (self *SpiceStacksController) require(conditions ...func() *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		for _, condition := range conditions {
			if reason := condition(); reason != nil {
				return reason
			}
		}
		return nil
	}
}
