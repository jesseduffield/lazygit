package controllers

import (
	"fmt"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// after selecting the 200th commit, we'll load in all the rest
const COMMIT_THRESHOLD = 200

type (
	PullFilesFn func() error
)

type LocalCommitsController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c *ControllerCommon

	pullFiles PullFilesFn
}

var _ types.IController = &LocalCommitsController{}

func NewLocalCommitsController(
	c *ControllerCommon,
	pullFiles PullFilesFn,
) *LocalCommitsController {
	return &LocalCommitsController{
		baseController: baseController{},
		c:              c,
		pullFiles:      pullFiles,
		ListControllerTrait: NewListControllerTrait[*models.Commit](
			c,
			c.Contexts().LocalCommits,
			c.Contexts().LocalCommits.GetSelected,
		),
	}
}

func (self *LocalCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	editCommitKey := opts.Config.Universal.Edit

	outsideFilterModeBindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Commits.SquashDown),
			Handler: self.withItem(self.squashDown),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.getDisabledReasonForSquashDown),
			),
			Description: self.c.Tr.SquashDown,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.MarkCommitAsFixup),
			Handler: self.withItem(self.fixup),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.getDisabledReasonForFixup),
			),
			Description: self.c.Tr.FixupCommit,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.RenameCommit),
			Handler: self.withItem(self.reword),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.rebaseCommandEnabled(todo.Reword)),
			),
			Description: self.c.Tr.RewordCommit,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.RenameCommitWithEditor),
			Handler: self.withItem(self.rewordEditor),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.rebaseCommandEnabled(todo.Reword)),
			),
			Description: self.c.Tr.RenameCommitEditor,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Remove),
			Handler: self.withItem(self.drop),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.rebaseCommandEnabled(todo.Drop)),
			),
			Description: self.c.Tr.DeleteCommit,
		},
		{
			Key:     opts.GetKey(editCommitKey),
			Handler: self.withItem(self.edit),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.rebaseCommandEnabled(todo.Edit)),
			),
			Description: self.c.Tr.EditCommit,
		},
		{
			// The user-facing description here is 'Start interactive rebase' but internally
			// we're calling it 'quick-start interactive rebase' to differentiate it from
			// when you manually select the base commit.
			Key:               opts.GetKey(opts.Config.Commits.StartInteractiveRebase),
			Handler:           self.withItem(self.quickStartInteractiveRebase),
			GetDisabledReason: self.require(self.notMidRebase, self.canFindCommitForQuickStart),
			Description:       self.c.Tr.QuickStartInteractiveRebase,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.QuickStartInteractiveRebaseTooltip, map[string]string{
				"editKey": keybindings.Label(editCommitKey),
			}),
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.PickCommit),
			Handler: self.withItem(self.pick),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.rebaseCommandEnabled(todo.Pick)),
			),
			Description: self.c.Tr.PickCommit,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CreateFixupCommit),
			Handler:           self.withItem(self.createFixupCommit),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CreateFixupCommitDescription,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.SquashAboveCommits),
			Handler: self.withItem(self.squashAllAboveFixupCommits),
			GetDisabledReason: self.require(
				self.notMidRebase,
				self.singleItemSelected(),
			),
			Description: self.c.Tr.SquashAboveCommits,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MoveDownCommit),
			Handler:           self.withItem(self.moveDown),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.MoveDownCommit,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MoveUpCommit),
			Handler:           self.withItem(self.moveUp),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.MoveUpCommit,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.PasteCommits),
			Handler:           self.paste,
			GetDisabledReason: self.require(self.canPaste),
			Description:       self.c.Tr.PasteCommits,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MarkCommitAsBaseForRebase),
			Handler:           self.withItem(self.markAsBaseCommit),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.MarkAsBaseCommit,
			Tooltip:           self.c.Tr.MarkAsBaseCommitTooltip,
		},
		// overriding these navigation keybindings because we might need to load
		// more commits on demand
		{
			Key:         opts.GetKey(opts.Config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.StartSearch,
			Tag:         "navigation",
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GotoBottom),
			Handler:     self.gotoBottom,
			Description: self.c.Tr.GotoBottom,
			Tag:         "navigation",
		},
	}

	for _, binding := range outsideFilterModeBindings {
		binding.Handler = opts.Guards.OutsideFilterMode(binding.Handler)
	}

	bindings := append(outsideFilterModeBindings, []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Commits.AmendToCommit),
			Handler:           self.withItem(self.amendTo),
			GetDisabledReason: self.require(self.singleItemSelected(self.canAmend)),
			Description:       self.c.Tr.AmendToCommit,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.ResetCommitAuthor),
			Handler:           self.withItem(self.amendAttribute),
			GetDisabledReason: self.require(self.singleItemSelected(self.canAmend)),
			Description:       self.c.Tr.SetResetCommitAuthor,
			OpensMenu:         true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.RevertCommit),
			Handler:           self.withItem(self.revert),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.RevertCommit,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CreateTag),
			Handler:           self.withItem(self.createTag),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.TagCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.OpenLogMenu),
			Handler:     self.handleOpenLogMenu,
			Description: self.c.Tr.OpenLogMenu,
			OpensMenu:   true,
		},
	}...)

	return bindings
}

func (self *LocalCommitsController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			var task types.UpdateTask
			commit := self.context().GetSelected()
			if commit == nil {
				task = types.NewRenderStringTask(self.c.Tr.NoCommitsThisBranch)
			} else if commit.Action == todo.UpdateRef {
				task = types.NewRenderStringTask(
					utils.ResolvePlaceholderString(
						self.c.Tr.UpdateRefHere,
						map[string]string{
							"ref": commit.Name,
						}))
			} else {
				cmdObj := self.c.Git().Commit.ShowCmdObj(commit.Sha, self.c.Modes().Filtering.GetPath())
				task = types.NewRunPtyTask(cmdObj.GetCmd())
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title:    "Patch",
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Task:     task,
				},
				Secondary: secondaryPatchPanelUpdateOpts(self.c),
			})
		})
	}
}

func secondaryPatchPanelUpdateOpts(c *ControllerCommon) *types.ViewUpdateOpts {
	if c.Git().Patch.PatchBuilder.Active() {
		patch := c.Git().Patch.PatchBuilder.RenderAggregatedPatch(false)

		return &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(patch),
			Title: c.Tr.CustomPatch,
		}
	}

	return nil
}

func (self *LocalCommitsController) squashDown(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand(todo.Squash, commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Squash,
		Prompt: self.c.Tr.SureSquashThisCommit,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.SquashingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.SquashCommitDown)
				return self.interactiveRebase(todo.Squash)
			})
		},
	})
}

func (self *LocalCommitsController) getDisabledReasonForSquashDown(commit *models.Commit) *types.DisabledReason {
	if self.context().GetSelectedLineIdx() >= len(self.c.Model().Commits)-1 {
		return &types.DisabledReason{Text: self.c.Tr.CannotSquashOrFixupFirstCommit}
	}

	return self.rebaseCommandEnabled(todo.Squash)(commit)
}

func (self *LocalCommitsController) fixup(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand(todo.Fixup, commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Fixup,
		Prompt: self.c.Tr.SureFixupThisCommit,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.FixingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.FixupCommit)
				return self.interactiveRebase(todo.Fixup)
			})
		},
	})
}

func (self *LocalCommitsController) getDisabledReasonForFixup(commit *models.Commit) *types.DisabledReason {
	if self.context().GetSelectedLineIdx() >= len(self.c.Model().Commits)-1 {
		return &types.DisabledReason{Text: self.c.Tr.CannotSquashOrFixupFirstCommit}
	}

	return self.rebaseCommandEnabled(todo.Squash)(commit)
}

func (self *LocalCommitsController) reword(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand(todo.Reword, commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	commitMessage, err := self.c.Git().Commit.GetCommitMessage(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	return self.c.Helpers().Commits.OpenCommitMessagePanel(
		&helpers.OpenCommitMessagePanelOpts{
			CommitIndex:      self.context().GetSelectedLineIdx(),
			InitialMessage:   commitMessage,
			SummaryTitle:     self.c.Tr.Actions.RewordCommit,
			DescriptionTitle: self.c.Tr.CommitDescriptionTitle,
			PreserveMessage:  false,
			OnConfirm:        self.handleReword,
			OnSwitchToEditor: self.switchFromCommitMessagePanelToEditor,
		},
	)
}

func (self *LocalCommitsController) switchFromCommitMessagePanelToEditor(filepath string) error {
	if self.isHeadCommit() {
		return self.c.RunSubprocessAndRefresh(
			self.c.Git().Commit.RewordLastCommitInEditorWithMessageFileCmdObj(filepath))
	}

	err := self.c.Git().Rebase.BeginInteractiveRebaseForCommit(self.c.Model().Commits, self.context().GetSelectedLineIdx(), false)
	if err != nil {
		return err
	}

	// now the selected commit should be our head so we'll amend it with the new message
	err = self.c.RunSubprocessAndRefresh(
		self.c.Git().Commit.RewordLastCommitInEditorWithMessageFileCmdObj(filepath))
	if err != nil {
		return err
	}

	err = self.c.Git().Rebase.ContinueRebase()
	if err != nil {
		return err
	}

	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *LocalCommitsController) handleReword(summary string, description string) error {
	err := self.c.Git().Rebase.RewordCommit(self.c.Model().Commits, self.c.Contexts().LocalCommits.GetSelectedLineIdx(), summary, description)
	if err != nil {
		return self.c.Error(err)
	}
	self.c.Helpers().Commits.OnCommitSuccess()
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *LocalCommitsController) doRewordEditor() error {
	self.c.LogAction(self.c.Tr.Actions.RewordCommit)

	if self.isHeadCommit() {
		return self.c.RunSubprocessAndRefresh(self.c.Git().Commit.RewordLastCommitInEditorCmdObj())
	}

	subProcess, err := self.c.Git().Rebase.RewordCommitInEditor(
		self.c.Model().Commits, self.context().GetSelectedLineIdx(),
	)
	if err != nil {
		return self.c.Error(err)
	}
	if subProcess != nil {
		return self.c.RunSubprocessAndRefresh(subProcess)
	}

	return nil
}

func (self *LocalCommitsController) rewordEditor(commit *models.Commit) error {
	midRebase, err := self.handleMidRebaseCommand(todo.Reword, commit)
	if err != nil {
		return err
	}
	if midRebase {
		return nil
	}

	if self.c.UserConfig.Gui.SkipRewordInEditorWarning {
		return self.doRewordEditor()
	} else {
		return self.c.Confirm(types.ConfirmOpts{
			Title:         self.c.Tr.RewordInEditorTitle,
			Prompt:        self.c.Tr.RewordInEditorPrompt,
			HandleConfirm: self.doRewordEditor,
		})
	}
}

func (self *LocalCommitsController) drop(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand(todo.Drop, commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DeleteCommitTitle,
		Prompt: self.c.Tr.DeleteCommitPrompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.DropCommit)
				return self.interactiveRebase(todo.Drop)
			})
		},
	})
}

func (self *LocalCommitsController) edit(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand(todo.Edit, commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.startInteractiveRebaseWithEdit(commit, commit)
}

func (self *LocalCommitsController) quickStartInteractiveRebase(selectedCommit *models.Commit) error {
	commitToEdit, err := self.findCommitForQuickStartInteractiveRebase()
	if err != nil {
		return self.c.Error(err)
	}

	return self.startInteractiveRebaseWithEdit(commitToEdit, selectedCommit)
}

func (self *LocalCommitsController) startInteractiveRebaseWithEdit(
	commitToEdit *models.Commit,
	selectedCommit *models.Commit,
) error {
	return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.EditCommit)
		err := self.c.Git().Rebase.EditRebase(commitToEdit.Sha)
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err,
			types.RefreshOptions{Mode: types.BLOCK_UI, Then: func() {
				// We need to select the same commit again because after starting a rebase,
				// new lines can be added for update-ref commands in the TODO file, due to
				// stacked branches. So the commit may be in a different position in the list.
				_, index, ok := lo.FindIndexOf(self.c.Model().Commits, func(c *models.Commit) bool {
					return c.Sha == selectedCommit.Sha
				})
				if ok {
					self.context().SetSelection(index)
				}
			}})
	})
}

func (self *LocalCommitsController) findCommitForQuickStartInteractiveRebase() (*models.Commit, error) {
	commit, index, ok := lo.FindIndexOf(self.c.Model().Commits, func(c *models.Commit) bool {
		return c.IsMerge() || c.Status == models.StatusMerged
	})

	if !ok || index == 0 {
		errorMsg := utils.ResolvePlaceholderString(self.c.Tr.CannotQuickStartInteractiveRebase, map[string]string{
			"editKey": keybindings.Label(self.c.UserConfig.Keybinding.Universal.Edit),
		})

		return nil, errors.New(errorMsg)
	}

	return commit, nil
}

func (self *LocalCommitsController) pick(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand(todo.Pick, commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	// at this point we aren't actually rebasing so we will interpret this as an
	// attempt to pull. We might revoke this later after enabling configurable keybindings
	return self.pullFiles()
}

func (self *LocalCommitsController) interactiveRebase(action todo.TodoCommand) error {
	err := self.c.Git().Rebase.InteractiveRebase(self.c.Model().Commits, self.context().GetSelectedLineIdx(), action)
	return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
}

// handleMidRebaseCommand sees if the selected commit is in fact a rebasing
// commit meaning you are trying to edit the todo file rather than actually
// begin a rebase. It then updates the todo file with that action
func (self *LocalCommitsController) handleMidRebaseCommand(action todo.TodoCommand, commit *models.Commit) (bool, error) {
	if !commit.IsTODO() {
		return false, nil
	}

	self.c.LogAction("Update rebase TODO")

	msg := utils.ResolvePlaceholderString(
		self.c.Tr.Log.HandleMidRebaseCommand,
		map[string]string{
			"shortSha": commit.ShortSha(),
			"action":   action.String(),
		},
	)
	self.c.LogCommand(msg, false)

	if err := self.c.Git().Rebase.EditRebaseTodo(commit, action); err != nil {
		return false, self.c.Error(err)
	}

	return true, self.c.Refresh(types.RefreshOptions{
		Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
	})
}

func (self *LocalCommitsController) rebaseCommandEnabled(action todo.TodoCommand) func(*models.Commit) *types.DisabledReason {
	return func(commit *models.Commit) *types.DisabledReason {
		if commit.Action == models.ActionConflict {
			return &types.DisabledReason{Text: self.c.Tr.ChangingThisActionIsNotAllowed}
		}

		if !commit.IsTODO() {
			if self.c.Model().WorkingTreeStateAtLastCommitRefresh != enums.REBASE_MODE_NONE {
				// If we are in a rebase, the only action that is allowed for
				// non-todo commits is rewording the current head commit
				if !(action == todo.Reword && self.isHeadCommit()) {
					return &types.DisabledReason{Text: self.c.Tr.AlreadyRebasing}
				}
			}

			return nil
		}

		// for now we do not support setting 'reword' because it requires an editor
		// and that means we either unconditionally wait around for the subprocess to ask for
		// our input or we set a lazygit client as the EDITOR env variable and have it
		// request us to edit the commit message when prompted.
		if action == todo.Reword {
			return &types.DisabledReason{Text: self.c.Tr.RewordNotSupported}
		}

		if allowed := isChangeOfRebaseTodoAllowed(action); !allowed {
			return &types.DisabledReason{Text: self.c.Tr.ChangingThisActionIsNotAllowed}
		}

		return nil
	}
}

func (self *LocalCommitsController) moveDown(commit *models.Commit) error {
	index := self.context().GetSelectedLineIdx()
	commits := self.c.Model().Commits

	// can't move past the initial commit
	if index >= len(commits)-1 {
		return nil
	}

	if commit.IsTODO() {
		if !commits[index+1].IsTODO() || commits[index+1].Action == models.ActionConflict {
			return nil
		}

		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)

		msg := utils.ResolvePlaceholderString(
			self.c.Tr.Log.MovingCommitDown,
			map[string]string{
				"shortSha": commit.ShortSha(),
			},
		)
		self.c.LogCommand(msg, false)

		if err := self.c.Git().Rebase.MoveTodoDown(commit); err != nil {
			return self.c.Error(err)
		}
		self.context().MoveSelectedLine(1)
		return self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
	}

	if self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return self.c.ErrorMsg(self.c.Tr.AlreadyRebasing)
	}

	return self.c.WithWaitingStatusSync(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)
		err := self.c.Git().Rebase.MoveCommitDown(self.c.Model().Commits, index)
		if err == nil {
			self.context().MoveSelectedLine(1)
		}
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err, types.RefreshOptions{Mode: types.SYNC})
	})
}

func (self *LocalCommitsController) moveUp(commit *models.Commit) error {
	index := self.context().GetSelectedLineIdx()
	if index == 0 {
		return nil
	}

	if commit.IsTODO() {
		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		msg := utils.ResolvePlaceholderString(
			self.c.Tr.Log.MovingCommitUp,
			map[string]string{
				"shortSha": commit.ShortSha(),
			},
		)
		self.c.LogCommand(msg, false)

		if err := self.c.Git().Rebase.MoveTodoUp(self.c.Model().Commits[index]); err != nil {
			return self.c.Error(err)
		}
		self.context().MoveSelectedLine(-1)
		return self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
	}

	if self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return self.c.ErrorMsg(self.c.Tr.AlreadyRebasing)
	}

	return self.c.WithWaitingStatusSync(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		err := self.c.Git().Rebase.MoveCommitUp(self.c.Model().Commits, index)
		if err == nil {
			self.context().MoveSelectedLine(-1)
		}
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err, types.RefreshOptions{Mode: types.SYNC})
	})
}

func (self *LocalCommitsController) amendTo(commit *models.Commit) error {
	if self.isHeadCommit() {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.AmendCommitTitle,
			Prompt: self.c.Tr.AmendCommitPrompt,
			HandleConfirm: func() error {
				return self.c.Helpers().WorkingTree.WithEnsureCommitableFiles(func() error {
					if err := self.c.Helpers().AmendHelper.AmendHead(); err != nil {
						return err
					}
					return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				})
			},
		})
	}

	if self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE {
		return self.c.ErrorMsg(self.c.Tr.AlreadyRebasing)
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AmendCommitTitle,
		Prompt: self.c.Tr.AmendCommitPrompt,
		HandleConfirm: func() error {
			return self.c.Helpers().WorkingTree.WithEnsureCommitableFiles(func() error {
				return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
					self.c.LogAction(self.c.Tr.Actions.AmendCommit)
					err := self.c.Git().Rebase.AmendTo(self.c.Model().Commits, self.context().GetView().SelectedLineIdx())
					return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
				})
			})
		},
	})
}

func (self *LocalCommitsController) canAmend(commit *models.Commit) *types.DisabledReason {
	if !self.isHeadCommit() && self.c.Model().WorkingTreeStateAtLastCommitRefresh != enums.REBASE_MODE_NONE {
		return &types.DisabledReason{Text: self.c.Tr.AlreadyRebasing}
	}

	return nil
}

func (self *LocalCommitsController) amendAttribute(commit *models.Commit) error {
	if self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE && !self.isHeadCommit() {
		return self.c.ErrorMsg(self.c.Tr.AlreadyRebasing)
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: "Amend commit attribute",
		Items: []*types.MenuItem{
			{
				Label:   self.c.Tr.ResetAuthor,
				OnPress: self.resetAuthor,
				Key:     'a',
				Tooltip: "Reset the commit's author to the currently configured user. This will also renew the author timestamp",
			},
			{
				Label:   self.c.Tr.SetAuthor,
				OnPress: self.setAuthor,
				Key:     'A',
				Tooltip: "Set the author based on a prompt",
			},
			{
				Label:   self.c.Tr.AddCoAuthor,
				OnPress: self.addCoAuthor,
				Key:     'c',
				Tooltip: self.c.Tr.AddCoAuthorTooltip,
			},
		},
	})
}

func (self *LocalCommitsController) resetAuthor() error {
	return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.ResetCommitAuthor)
		if err := self.c.Git().Rebase.ResetCommitAuthor(self.c.Model().Commits, self.context().GetSelectedLineIdx()); err != nil {
			return self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *LocalCommitsController) setAuthor() error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.SetAuthorPromptTitle,
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc(),
		HandleConfirm: func(value string) error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.SetCommitAuthor)
				if err := self.c.Git().Rebase.SetCommitAuthor(self.c.Model().Commits, self.context().GetSelectedLineIdx(), value); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			})
		},
	})
}

func (self *LocalCommitsController) addCoAuthor() error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.AddCoAuthorPromptTitle,
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc(),
		HandleConfirm: func(value string) error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.AddCommitCoAuthor)
				if err := self.c.Git().Rebase.AddCommitCoAuthor(self.c.Model().Commits, self.context().GetSelectedLineIdx(), value); err != nil {
					return self.c.Error(err)
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			})
		},
	})
}

func (self *LocalCommitsController) revert(commit *models.Commit) error {
	if commit.IsMerge() {
		return self.createRevertMergeCommitMenu(commit)
	} else {
		return self.c.Confirm(types.ConfirmOpts{
			Title: self.c.Tr.Actions.RevertCommit,
			Prompt: utils.ResolvePlaceholderString(
				self.c.Tr.ConfirmRevertCommit,
				map[string]string{
					"selectedCommit": commit.ShortSha(),
				}),
			HandleConfirm: func() error {
				self.c.LogAction(self.c.Tr.Actions.RevertCommit)
				return self.c.WithWaitingStatusSync(self.c.Tr.RevertingStatus, func() error {
					if err := self.c.Git().Commit.Revert(commit.Sha); err != nil {
						return err
					}
					return self.afterRevertCommit()
				})
			},
		})
	}
}

func (self *LocalCommitsController) createRevertMergeCommitMenu(commit *models.Commit) error {
	menuItems := make([]*types.MenuItem, len(commit.Parents))
	for i, parentSha := range commit.Parents {
		i := i
		message, err := self.c.Git().Commit.GetCommitMessageFirstLine(parentSha)
		if err != nil {
			return self.c.Error(err)
		}

		menuItems[i] = &types.MenuItem{
			Label: fmt.Sprintf("%s: %s", utils.SafeTruncate(parentSha, 8), message),
			OnPress: func() error {
				parentNumber := i + 1
				self.c.LogAction(self.c.Tr.Actions.RevertCommit)
				return self.c.WithWaitingStatusSync(self.c.Tr.RevertingStatus, func() error {
					if err := self.c.Git().Commit.RevertMerge(commit.Sha, parentNumber); err != nil {
						return err
					}
					return self.afterRevertCommit()
				})
			},
		}
	}

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.SelectParentCommitForMerge, Items: menuItems})
}

func (self *LocalCommitsController) afterRevertCommit() error {
	self.context().MoveSelectedLine(1)
	return self.c.Refresh(types.RefreshOptions{
		Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS, types.BRANCHES},
	})
}

func (self *LocalCommitsController) createFixupCommit(commit *models.Commit) error {
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.SureCreateFixupCommit,
		map[string]string{
			"commit": commit.Sha,
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.CreateFixupCommit,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.Helpers().WorkingTree.WithEnsureCommitableFiles(func() error {
				self.c.LogAction(self.c.Tr.Actions.CreateFixupCommit)
				if err := self.c.Git().Commit.CreateFixupCommit(commit.Sha); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			})
		},
	})
}

func (self *LocalCommitsController) squashAllAboveFixupCommits(commit *models.Commit) error {
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.SureSquashAboveCommits,
		map[string]string{"commit": commit.Sha},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.SquashAboveCommits,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.SquashingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.SquashAllAboveFixupCommits)
				err := self.c.Git().Rebase.SquashAllAboveFixupCommits(commit)
				return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
}

// For getting disabled reason
func (self *LocalCommitsController) notMidRebase() *types.DisabledReason {
	if self.c.Model().WorkingTreeStateAtLastCommitRefresh != enums.REBASE_MODE_NONE {
		return &types.DisabledReason{Text: self.c.Tr.AlreadyRebasing}
	}

	return nil
}

// For getting disabled reason
func (self *LocalCommitsController) canFindCommitForQuickStart() *types.DisabledReason {
	if _, err := self.findCommitForQuickStartInteractiveRebase(); err != nil {
		return &types.DisabledReason{Text: err.Error(), ShowErrorInPanel: true}
	}

	return nil
}

func (self *LocalCommitsController) createTag(commit *models.Commit) error {
	return self.c.Helpers().Tags.OpenCreateTagPrompt(commit.Sha, func() {})
}

func (self *LocalCommitsController) openSearch() error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if self.context().GetLimitCommits() {
		self.context().SetLimitCommits(false)
		if err := self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
			return err
		}
	}

	return self.c.Helpers().Search.OpenSearchPrompt(self.context())
}

func (self *LocalCommitsController) gotoBottom() error {
	// we usually lazyload these commits but now that we're jumping to the bottom we need to load them now
	if self.context().GetLimitCommits() {
		self.context().SetLimitCommits(false)
		if err := self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
			return err
		}
	}

	self.context().SetSelectedLineIdx(self.context().Len() - 1)

	return nil
}

func (self *LocalCommitsController) handleOpenLogMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.LogMenuTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.ToggleShowGitGraphAll,
				OnPress: func() error {
					self.context().SetShowWholeGitGraph(!self.context().GetShowWholeGitGraph())

					if self.context().GetShowWholeGitGraph() {
						self.context().SetLimitCommits(false)
					}

					return self.c.WithWaitingStatus(self.c.Tr.LoadingCommits, func(gocui.Task) error {
						return self.c.Refresh(
							types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}},
						)
					})
				},
			},
			{
				Label:     self.c.Tr.ShowGitGraph,
				OpensMenu: true,
				OnPress: func() error {
					onPress := func(value string) func() error {
						return func() error {
							self.c.UserConfig.Git.Log.ShowGraph = value
							return nil
						}
					}
					return self.c.Menu(types.CreateMenuOptions{
						Title: self.c.Tr.LogMenuTitle,
						Items: []*types.MenuItem{
							{
								Label:   "always",
								OnPress: onPress("always"),
							},
							{
								Label:   "never",
								OnPress: onPress("never"),
							},
							{
								Label:   "when maximised",
								OnPress: onPress("when-maximised"),
							},
						},
					})
				},
			},
			{
				Label:     self.c.Tr.SortCommits,
				OpensMenu: true,
				OnPress: func() error {
					onPress := func(value string) func() error {
						return func() error {
							self.c.UserConfig.Git.Log.Order = value
							return self.c.WithWaitingStatus(self.c.Tr.LoadingCommits, func(gocui.Task) error {
								return self.c.Refresh(
									types.RefreshOptions{
										Mode:  types.SYNC,
										Scope: []types.RefreshableView{types.COMMITS},
									},
								)
							})
						}
					}

					return self.c.Menu(types.CreateMenuOptions{
						Title: self.c.Tr.LogMenuTitle,
						Items: []*types.MenuItem{
							{
								Label:   "topological (topo-order)",
								OnPress: onPress("topo-order"),
							},
							{
								Label:   "date-order",
								OnPress: onPress("date-order"),
							},
							{
								Label:   "author-date-order",
								OnPress: onPress("author-date-order"),
							},
						},
					})
				},
			},
		},
	})
}

func (self *LocalCommitsController) GetOnFocus() func(types.OnFocusOpts) error {
	return func(types.OnFocusOpts) error {
		context := self.context()
		if context.GetSelectedLineIdx() > COMMIT_THRESHOLD && context.GetLimitCommits() {
			context.SetLimitCommits(false)
			self.c.OnWorker(func(_ gocui.Task) {
				if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
					_ = self.c.Error(err)
				}
			})
		}

		return nil
	}
}

func (self *LocalCommitsController) context() *context.LocalCommitsContext {
	return self.c.Contexts().LocalCommits
}

func (self *LocalCommitsController) paste() error {
	return self.c.Helpers().CherryPick.Paste()
}

func (self *LocalCommitsController) canPaste() *types.DisabledReason {
	if !self.c.Helpers().CherryPick.CanPaste() {
		return &types.DisabledReason{Text: self.c.Tr.NoCopiedCommits}
	}

	return nil
}

func (self *LocalCommitsController) markAsBaseCommit(commit *models.Commit) error {
	if commit.Sha == self.c.Modes().MarkedBaseCommit.GetSha() {
		// Reset when invoking it again on the marked commit
		self.c.Modes().MarkedBaseCommit.SetSha("")
	} else {
		self.c.Modes().MarkedBaseCommit.SetSha(commit.Sha)
	}
	return self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
}

func (self *LocalCommitsController) isHeadCommit() bool {
	return models.IsHeadCommit(self.c.Model().Commits, self.context().GetSelectedLineIdx())
}

func isChangeOfRebaseTodoAllowed(action todo.TodoCommand) bool {
	allowedActions := []todo.TodoCommand{
		todo.Pick,
		todo.Drop,
		todo.Edit,
		todo.Fixup,
		todo.Squash,
		todo.Reword,
	}

	return lo.Contains(allowedActions, action)
}
