package controllers

import (
	"errors"
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
		// === NAVIGATION ===
		{
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ViewCommits,
		},
		// === COMMIT COMMANDS (only enabled when commit selected) ===
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

		// === BRANCH COMMANDS ===
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.press),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			DisplayOnScreen:   true,
		},

		// === MOVE UP/DOWN (context-sensitive) ===
		{
			Key:               opts.GetKey(opts.Config.Commits.MoveDownCommit),
			Handler:           self.withItem(self.moveDown),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.MoveDownCommit, // fallback for cheatsheet
			DescriptionFunc: func() string {
				item := self.context().GetSelected()
				if item != nil && !item.IsCommit {
					return "Move branch down in stack"
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
					return "Move branch up in stack"
				}
				return self.c.Tr.MoveUpCommit
			},
		},

		// === MENUS ===
		{
			Key:             opts.GetKey("S"),
			Handler:         self.openStackOperationsMenu,
			Description:     "Stack operations",
			OpensMenu:       true,
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey("G"),
			Handler:         self.openNavigationMenu,
			Description:     "Stack navigation",
			OpensMenu:       true,
			DisplayOnScreen: true,
		},

		// === VIEW OPTIONS ===
		{
			Key:         opts.GetKey("V"),
			Handler:     self.toggleLogFormat,
			Description: self.c.Tr.ToggleSpiceLogFormat,
			Tooltip:     self.c.Tr.ToggleSpiceLogFormatTooltip,
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
				task = types.NewRenderStringTask("No stacks")
				title = self.c.Tr.SpiceStacksTitle
			} else if item.IsCommit {
				// Show commit diff/patch
				cmdObj := self.c.Git().Commit.ShowCmdObj(item.CommitSha, nil)
				task = types.NewRunPtyTask(cmdObj.GetCmd())
				title = "Patch"
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

// === NAVIGATION HANDLERS ===

func (self *SpiceStacksController) enter(item *models.SpiceStackItem) error {
	if item.IsCommit {
		return self.viewCommitFiles(item)
	}
	return self.viewBranchCommits(item)
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
		return errors.New("Branch not found")
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
		return errors.New("Commit not found in commits list")
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

// === COMMIT COMMAND HANDLERS ===

// findCommitByHash searches for a commit in the model's commits list by SHA
// Uses prefix matching since SpiceStackItem stores short (7-char) hashes
func (self *SpiceStacksController) findCommitByHash(sha string) (*models.Commit, int) {
	for idx, commit := range self.c.Model().Commits {
		if strings.HasPrefix(commit.Hash(), sha) {
			return commit, idx
		}
	}
	return nil, -1
}

func (self *SpiceStacksController) commitSquash(item *models.SpiceStackItem) error {
	commit, commitIdx := self.findCommitByHash(item.CommitSha)
	if commit == nil {
		return errors.New("Commit not found in commits list")
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
	commit, commitIdx := self.findCommitByHash(item.CommitSha)
	if commit == nil {
		return errors.New("Commit not found in commits list")
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
	commit, commitIdx := self.findCommitByHash(item.CommitSha)
	if commit == nil {
		return errors.New("Commit not found in commits list")
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

	_, commitIdx := self.findCommitByHash(item.CommitSha)
	if commitIdx == -1 {
		return errors.New("Commit not found")
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
	_, commitIdx := self.findCommitByHash(item.CommitSha)
	if commitIdx == -1 {
		return errors.New("Commit not found in commits list")
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
	_, commitIdx := self.findCommitByHash(item.CommitSha)
	if commitIdx == -1 {
		return errors.New("Commit not found in commits list")
	}

	self.c.LogAction(self.c.Tr.Actions.EditCommit)
	err := self.c.Git().Rebase.InteractiveRebase(self.c.Model().Commits, commitIdx, commitIdx, todo.Edit)
	return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
}

func (self *SpiceStacksController) commitAmend(item *models.SpiceStackItem) error {
	_, commitIdx := self.findCommitByHash(item.CommitSha)
	if commitIdx == -1 {
		return errors.New("Commit not found in commits list")
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
	commit, _ := self.findCommitByHash(item.CommitSha)
	if commit == nil {
		return errors.New("Commit not found in commits list")
	}

	return self.c.Helpers().Refs.CreateCheckoutMenu(commit)
}

func (self *SpiceStacksController) commitReset(item *models.SpiceStackItem) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(item.CommitSha, item.CommitSha)
}

func (self *SpiceStacksController) commitCopy(item *models.SpiceStackItem) error {
	commit, _ := self.findCommitByHash(item.CommitSha)
	if commit == nil {
		return errors.New("Commit not found in commits list")
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

// === BRANCH COMMAND HANDLERS ===

func (self *SpiceStacksController) checkout(item *models.SpiceStackItem) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	if err := self.c.Git().Branch.Checkout(item.Name, git_commands.CheckoutOptions{Force: false}); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) restack(item *models.SpiceStackItem) error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceRestackingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Restack(item.Name)
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) restackAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceRestackingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Restack("")
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) submit(item *models.SpiceStackItem) error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceSubmittingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit(item.Name)
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) submitAll() error {
	return self.c.WithWaitingStatus(self.c.Tr.SpiceSubmittingStatus, func(task gocui.Task) error {
		err := self.c.Git().Spice.Submit("")
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
		return nil
	})
}

func (self *SpiceStacksController) navigateUp() error {
	if err := self.c.Git().Spice.NavigateUp(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) navigateDown() error {
	if err := self.c.Git().Spice.NavigateDown(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) navigateTop() error {
	if err := self.c.Git().Spice.NavigateTop(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) navigateBottom() error {
	if err := self.c.Git().Spice.NavigateBottom(); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) newBranch() error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.SpiceBranchNamePrompt,
		HandleConfirm: func(branchName string) error {
			if err := self.c.Git().Spice.CreateBranch(branchName); err != nil {
				return err
			}
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.FILES, types.SPICE_STACKS}})
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
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES, types.SPICE_STACKS}})
			return nil
		},
	})
	return nil
}

func (self *SpiceStacksController) moveBranchUp(item *models.SpiceStackItem) error {
	if err := self.c.Git().Spice.MoveBranchUp(item.Name); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) moveBranchDown(item *models.SpiceStackItem) error {
	if err := self.c.Git().Spice.MoveBranchDown(item.Name); err != nil {
		return err
	}
	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.SPICE_STACKS}})
	return nil
}

func (self *SpiceStacksController) commitMoveDown(item *models.SpiceStackItem) error {
	_, commitIdx := self.findCommitByHash(item.CommitSha)
	if commitIdx == -1 {
		return errors.New("Commit not found")
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
	_, commitIdx := self.findCommitByHash(item.CommitSha)
	if commitIdx == -1 {
		return errors.New("Commit not found")
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

// === MENU HANDLERS ===

func (self *SpiceStacksController) openStackOperationsMenu() error {
	item := self.context().GetSelected()

	menuItems := []*types.MenuItem{
		{
			Label:   "Restack branch",
			Key:     'r',
			OnPress: func() error { return self.restack(item) },
			DisabledReason: self.branchSelectedReason(item),
		},
		{
			Label:   "Restack all",
			Key:     'R',
			OnPress: func() error { return self.restackAll() },
		},
		{
			Label:   "Submit branch",
			Key:     's',
			OnPress: func() error { return self.submit(item) },
			DisabledReason: self.branchSelectedReason(item),
		},
		{
			Label:   "Submit all",
			Key:     'S',
			OnPress: func() error { return self.submitAll() },
		},
		{
			Label:   "Create branch",
			Key:     'c',
			OnPress: func() error { return self.newBranch() },
		},
		{
			Label:   "Delete branch",
			Key:     'd',
			OnPress: func() error { return self.delete(item) },
			DisabledReason: self.branchSelectedReason(item),
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: "Stack Operations",
		Items: menuItems,
	})
}

func (self *SpiceStacksController) openNavigationMenu() error {
	menuItems := []*types.MenuItem{
		{
			Label:   "Up",
			Key:     'u',
			Tooltip: "Navigate up one branch in the stack (gs up)",
			OnPress: func() error { return self.navigateUp() },
		},
		{
			Label:   "Down",
			Key:     'd',
			Tooltip: "Navigate down one branch in the stack (gs down)",
			OnPress: func() error { return self.navigateDown() },
		},
		{
			Label:   "Top",
			Key:     't',
			Tooltip: "Navigate to the top of the stack (gs top)",
			OnPress: func() error { return self.navigateTop() },
		},
		{
			Label:   "Bottom",
			Key:     'b',
			Tooltip: "Navigate to the bottom of the stack (gs bottom)",
			OnPress: func() error { return self.navigateBottom() },
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: "Stack Navigation",
		Items: menuItems,
	})
}

func (self *SpiceStacksController) branchSelectedReason(item *models.SpiceStackItem) *types.DisabledReason {
	if item == nil || item.IsCommit {
		return &types.DisabledReason{Text: self.c.Tr.SpiceBranchOnly}
	}
	return nil
}

func (self *SpiceStacksController) toggleLogFormat() error {
	currentFormat := self.c.UserConfig().Git.Spice.LogFormat

	if currentFormat == "long" {
		self.c.UserConfig().Git.Spice.LogFormat = "short"
	} else {
		self.c.UserConfig().Git.Spice.LogFormat = "long"
	}

	// Refresh the spice stacks view
	self.c.Refresh(types.RefreshOptions{
		Mode:  types.ASYNC,
		Scope: []types.RefreshableView{types.SPICE_STACKS},
	})
	return nil
}

// === HELPER METHODS ===

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
			return &types.DisabledReason{Text: "No item selected"}
		}
		return nil
	}
}

func (self *SpiceStacksController) branchSelected() func() *types.DisabledReason {
	return func() *types.DisabledReason {
		item := self.context().GetSelected()
		if item == nil || item.IsCommit {
			return &types.DisabledReason{Text: self.c.Tr.SpiceBranchOnly}
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
