package controllers

import (
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
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
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().LocalCommits,
			c.Contexts().LocalCommits.GetSelected,
			c.Contexts().LocalCommits.GetSelectedItems,
		),
	}
}

func (self *LocalCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	editCommitKey := opts.Config.Universal.Edit

	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Commits.SquashDown),
			Handler: opts.Guards.OutsideFilterMode(self.withItemsRange(self.squashDown)),
			GetDisabledReason: self.require(
				self.itemRangeSelected(
					self.midRebaseCommandEnabled,
					self.canSquashOrFixup,
				),
			),
			Description:     self.c.Tr.Squash,
			Tooltip:         self.c.Tr.SquashTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.MarkCommitAsFixup),
			Handler: opts.Guards.OutsideFilterMode(self.withItemsRange(self.fixup)),
			GetDisabledReason: self.require(
				self.itemRangeSelected(
					self.midRebaseCommandEnabled,
					self.canSquashOrFixup,
				),
			),
			Description:     self.c.Tr.Fixup,
			Tooltip:         self.c.Tr.FixupTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.RenameCommit),
			Handler: self.withItem(self.reword),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.rewordEnabled),
			),
			Description:     self.c.Tr.Reword,
			Tooltip:         self.c.Tr.CommitRewordTooltip,
			DisplayOnScreen: true,
			OpensMenu:       true,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.RenameCommitWithEditor),
			Handler: self.withItem(self.rewordEditor),
			GetDisabledReason: self.require(
				self.singleItemSelected(self.rewordEnabled),
			),
			Description: self.c.Tr.RewordCommitEditor,
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Remove),
			Handler: self.withItemsRange(self.drop),
			GetDisabledReason: self.require(
				self.itemRangeSelected(
					self.canDropCommits,
				),
			),
			Description:     self.c.Tr.DropCommit,
			Tooltip:         self.c.Tr.DropCommitTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:     opts.GetKey(editCommitKey),
			Handler: opts.Guards.OutsideFilterMode(self.withItemsRange(self.edit)),
			GetDisabledReason: self.require(
				self.itemRangeSelected(self.midRebaseCommandEnabled),
			),
			Description:      self.c.Tr.EditCommit,
			ShortDescription: self.c.Tr.Edit,
			Tooltip:          self.c.Tr.EditCommitTooltip,
			DisplayOnScreen:  true,
		},
		{
			// The user-facing description here is 'Start interactive rebase' but internally
			// we're calling it 'quick-start interactive rebase' to differentiate it from
			// when you manually select the base commit.
			Key:               opts.GetKey(opts.Config.Commits.StartInteractiveRebase),
			Handler:           opts.Guards.OutsideFilterMode(self.quickStartInteractiveRebase),
			GetDisabledReason: self.require(self.notMidRebase(self.c.Tr.AlreadyRebasing), self.canFindCommitForQuickStart),
			Description:       self.c.Tr.QuickStartInteractiveRebase,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.QuickStartInteractiveRebaseTooltip, map[string]string{
				"editKey": keybindings.Label(editCommitKey),
			}),
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.PickCommit),
			Handler: opts.Guards.OutsideFilterMode(self.withItems(self.pick)),
			GetDisabledReason: self.require(
				self.itemRangeSelected(self.pickEnabled),
			),
			Description: self.c.Tr.Pick,
			Tooltip:     self.c.Tr.PickCommitTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CreateFixupCommit),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.createFixupCommit)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CreateFixupCommit,
			Tooltip: utils.ResolvePlaceholderString(
				self.c.Tr.CreateFixupCommitTooltip,
				map[string]string{
					"squashAbove": keybindings.Label(opts.Config.Commits.SquashAboveCommits),
				},
			),
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.SquashAboveCommits),
			Handler: opts.Guards.OutsideFilterMode(self.squashFixupCommits),
			GetDisabledReason: self.require(
				self.notMidRebase(self.c.Tr.AlreadyRebasing),
			),
			Description: self.c.Tr.SquashAboveCommits,
			Tooltip:     self.c.Tr.SquashAboveCommitsTooltip,
			OpensMenu:   true,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.MoveDownCommit),
			Handler: opts.Guards.OutsideFilterMode(self.withItemsRange(self.moveDown)),
			GetDisabledReason: self.require(self.itemRangeSelected(
				self.midRebaseMoveCommandEnabled,
				self.canMoveDown,
			)),
			Description: self.c.Tr.MoveDownCommit,
		},
		{
			Key:     opts.GetKey(opts.Config.Commits.MoveUpCommit),
			Handler: opts.Guards.OutsideFilterMode(self.withItemsRange(self.moveUp)),
			GetDisabledReason: self.require(self.itemRangeSelected(
				self.midRebaseMoveCommandEnabled,
				self.canMoveUp,
			)),
			Description: self.c.Tr.MoveUpCommit,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.PasteCommits),
			Handler:           opts.Guards.OutsideFilterMode(self.paste),
			GetDisabledReason: self.require(self.canPaste),
			Description:       self.c.Tr.PasteCommits,
			DisplayStyle:      &style.FgCyan,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.MarkCommitAsBaseForRebase),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.markAsBaseCommit)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.MarkAsBaseCommit,
			Tooltip:           self.c.Tr.MarkAsBaseCommitTooltip,
		},
		// overriding this navigation keybinding because we might need to load
		// more commits on demand
		{
			Key:         opts.GetKey(opts.Config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.StartSearch,
			Tag:         "navigation",
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.AmendToCommit),
			Handler:           self.withItem(self.amendTo),
			GetDisabledReason: self.require(self.singleItemSelected(self.canAmend)),
			Description:       self.c.Tr.Amend,
			Tooltip:           self.c.Tr.AmendCommitTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.ResetCommitAuthor),
			Handler:           self.withItemsRange(self.amendAttribute),
			GetDisabledReason: self.require(self.itemRangeSelected(self.canAmendRange)),
			Description:       self.c.Tr.AmendCommitAttribute,
			Tooltip:           self.c.Tr.AmendCommitAttributeTooltip,
			OpensMenu:         true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.RevertCommit),
			Handler:           self.withItemsRange(self.revert),
			GetDisabledReason: self.require(self.itemRangeSelected()),
			Description:       self.c.Tr.Revert,
			Tooltip:           self.c.Tr.RevertCommitTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.CreateTag),
			Handler:           self.withItem(self.createTag),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.TagCommit,
			Tooltip:           self.c.Tr.TagCommitTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.OpenLogMenu),
			Handler:     self.handleOpenLogMenu,
			Description: self.c.Tr.OpenLogMenu,
			Tooltip:     self.c.Tr.OpenLogMenuTooltip,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *LocalCommitsController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			commit := self.context().GetSelected()
			if commit == nil {
				task = types.NewRenderStringTask(self.c.Tr.NoCommitsThisBranch)
			} else if commit.Action == todo.UpdateRef {
				task = types.NewRenderStringTask(
					utils.ResolvePlaceholderString(
						self.c.Tr.UpdateRefHere,
						map[string]string{
							"ref": strings.TrimPrefix(commit.Name, "refs/heads/"),
						}))
			} else if commit.Action == todo.Exec {
				task = types.NewRenderStringTask(
					self.c.Tr.ExecCommandHere + "\n\n" + commit.Name)
			} else {
				refRange := self.context().GetSelectedRefRangeForDiffFiles()
				task = self.c.Helpers().Diff.GetUpdateTaskForRenderingCommitsDiff(commit, refRange)
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
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

func (self *LocalCommitsController) squashDown(selectedCommits []*models.Commit, startIdx int, endIdx int) error {
	if self.isRebasing() {
		return self.updateTodos(todo.Squash, selectedCommits)
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Squash,
		Prompt: self.c.Tr.SureSquashThisCommit,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.SquashingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.SquashCommitDown)
				return self.interactiveRebase(todo.Squash, startIdx, endIdx)
			})
		},
	})

	return nil
}

func (self *LocalCommitsController) fixup(selectedCommits []*models.Commit, startIdx int, endIdx int) error {
	if self.isRebasing() {
		return self.updateTodos(todo.Fixup, selectedCommits)
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Fixup,
		Prompt: self.c.Tr.SureFixupThisCommit,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.FixingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.FixupCommit)
				return self.interactiveRebase(todo.Fixup, startIdx, endIdx)
			})
		},
	})

	return nil
}

func (self *LocalCommitsController) reword(commit *models.Commit) error {
	commitIdx := self.context().GetSelectedLineIdx()
	if self.c.Git().Config.NeedsGpgSubprocessForCommit() && !self.isHeadCommit(commitIdx) {
		return errors.New(self.c.Tr.DisabledForGPG)
	}
	commitMessage, err := self.c.Git().Commit.GetCommitMessage(commit.Hash())
	if err != nil {
		return err
	}
	if self.c.UserConfig().Git.Commit.AutoWrapCommitMessage {
		commitMessage = helpers.TryRemoveHardLineBreaks(commitMessage, self.c.UserConfig().Git.Commit.AutoWrapWidth)
	}
	self.c.Helpers().Commits.OpenCommitMessagePanel(
		&helpers.OpenCommitMessagePanelOpts{
			CommitIndex:      commitIdx,
			InitialMessage:   commitMessage,
			SummaryTitle:     self.c.Tr.Actions.RewordCommit,
			DescriptionTitle: self.c.Tr.CommitDescriptionTitle,
			PreserveMessage:  false,
			OnConfirm:        self.handleReword,
			OnSwitchToEditor: self.switchFromCommitMessagePanelToEditor,
		},
	)

	return nil
}

func (self *LocalCommitsController) switchFromCommitMessagePanelToEditor(filepath string) error {
	if self.isSelectedHeadCommit() {
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

	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
}

func (self *LocalCommitsController) handleReword(summary string, description string) error {
	if models.IsHeadCommit(self.c.Model().Commits, self.c.Contexts().LocalCommits.GetSelectedLineIdx()) {
		// we've selected the top commit so no rebase is required
		return self.c.Helpers().GPG.WithGpgHandling(self.c.Git().Commit.RewordLastCommit(summary, description),
			git_commands.CommitGpgSign,
			self.c.Tr.RewordingStatus, nil, nil)
	}

	return self.c.WithWaitingStatus(self.c.Tr.RewordingStatus, func(gocui.Task) error {
		err := self.c.Git().Rebase.RewordCommit(self.c.Model().Commits, self.c.Contexts().LocalCommits.GetSelectedLineIdx(), summary, description)
		if err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	})
}

func (self *LocalCommitsController) doRewordEditor() error {
	self.c.LogAction(self.c.Tr.Actions.RewordCommit)

	if self.isSelectedHeadCommit() {
		return self.c.RunSubprocessAndRefresh(self.c.Git().Commit.RewordLastCommitInEditorCmdObj())
	}

	subProcess, err := self.c.Git().Rebase.RewordCommitInEditor(
		self.c.Model().Commits, self.context().GetSelectedLineIdx(),
	)
	if err != nil {
		return err
	}
	if subProcess != nil {
		return self.c.RunSubprocessAndRefresh(subProcess)
	}

	return nil
}

func (self *LocalCommitsController) rewordEditor(commit *models.Commit) error {
	return self.c.ConfirmIf(!self.c.UserConfig().Gui.SkipRewordInEditorWarning,
		types.ConfirmOpts{
			Title:         self.c.Tr.RewordInEditorTitle,
			Prompt:        self.c.Tr.RewordInEditorPrompt,
			HandleConfirm: self.doRewordEditor,
		})
}

func (self *LocalCommitsController) drop(selectedCommits []*models.Commit, startIdx int, endIdx int) error {
	if self.isRebasing() {
		groupedTodos := lo.GroupBy(selectedCommits, func(c *models.Commit) bool {
			return c.Action == todo.UpdateRef
		})
		updateRefTodos := groupedTodos[true]
		nonUpdateRefTodos := groupedTodos[false]

		if len(updateRefTodos) > 0 {
			self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.DropCommitTitle,
				Prompt: self.c.Tr.DropUpdateRefPrompt,
				HandleConfirm: func() error {
					selectedIdx, rangeStartIdx, rangeSelectMode := self.context().GetSelectionRangeAndMode()

					if err := self.c.Git().Rebase.DeleteUpdateRefTodos(updateRefTodos); err != nil {
						return err
					}

					if selectedIdx > rangeStartIdx {
						selectedIdx = max(selectedIdx-len(updateRefTodos), rangeStartIdx)
					} else {
						rangeStartIdx = max(rangeStartIdx-len(updateRefTodos), selectedIdx)
					}

					self.context().SetSelectionRangeAndMode(selectedIdx, rangeStartIdx, rangeSelectMode)

					return self.updateTodos(todo.Drop, nonUpdateRefTodos)
				},
			})

			return nil
		}

		return self.updateTodos(todo.Drop, selectedCommits)
	}

	isMerge := selectedCommits[0].IsMerge()

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DropCommitTitle,
		Prompt: lo.Ternary(isMerge, self.c.Tr.DropMergeCommitPrompt, self.c.Tr.DropCommitPrompt),
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DroppingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.DropCommit)
				if isMerge {
					return self.dropMergeCommit(startIdx)
				}
				return self.interactiveRebase(todo.Drop, startIdx, endIdx)
			})
		},
	})

	return nil
}

func (self *LocalCommitsController) dropMergeCommit(commitIdx int) error {
	err := self.c.Git().Rebase.DropMergeCommit(self.c.Model().Commits, commitIdx)
	return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
}

func (self *LocalCommitsController) edit(selectedCommits []*models.Commit, startIdx int, endIdx int) error {
	if self.isRebasing() {
		return self.updateTodos(todo.Edit, selectedCommits)
	}

	commits := self.c.Model().Commits
	if !commits[endIdx].IsMerge() {
		selectionRangeAndMode := self.getSelectionRangeAndMode()
		err := self.c.Git().Rebase.InteractiveRebase(commits, startIdx, endIdx, todo.Edit)
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err,
			types.RefreshOptions{
				Mode: types.BLOCK_UI, Then: func() {
					self.restoreSelectionRangeAndMode(selectionRangeAndMode)
				},
			})
	}

	return self.startInteractiveRebaseWithEdit(selectedCommits)
}

func (self *LocalCommitsController) quickStartInteractiveRebase() error {
	commitToEdit, err := self.findCommitForQuickStartInteractiveRebase()
	if err != nil {
		return err
	}

	return self.startInteractiveRebaseWithEdit([]*models.Commit{commitToEdit})
}

func (self *LocalCommitsController) startInteractiveRebaseWithEdit(
	commitsToEdit []*models.Commit,
) error {
	return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.EditCommit)
		selectionRangeAndMode := self.getSelectionRangeAndMode()
		err := self.c.Git().Rebase.EditRebase(commitsToEdit[len(commitsToEdit)-1].Hash())
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err,
			types.RefreshOptions{Mode: types.BLOCK_UI, Then: func() {
				todos := make([]*models.Commit, 0, len(commitsToEdit)-1)
				for _, c := range commitsToEdit[:len(commitsToEdit)-1] {
					// Merge commits can't be set to "edit", so just skip them
					if !c.IsMerge() {
						todos = append(todos, models.NewCommit(self.c.Model().HashPool, models.NewCommitOpts{Hash: c.Hash(), Action: todo.Pick}))
					}
				}
				if len(todos) > 0 {
					err := self.updateTodos(todo.Edit, todos)
					if err != nil {
						self.c.Log.Errorf("error when updating todos: %v", err)
					}
				}

				self.restoreSelectionRangeAndMode(selectionRangeAndMode)
			}})
	})
}

type SelectionRangeAndMode struct {
	selectedHash   string
	rangeStartHash string
	mode           traits.RangeSelectMode
}

func (self *LocalCommitsController) getSelectionRangeAndMode() SelectionRangeAndMode {
	selectedIdx, rangeStartIdx, rangeSelectMode := self.context().GetSelectionRangeAndMode()
	commits := self.c.Model().Commits
	selectedHash := commits[selectedIdx].Hash()
	rangeStartHash := commits[rangeStartIdx].Hash()
	return SelectionRangeAndMode{selectedHash, rangeStartHash, rangeSelectMode}
}

func (self *LocalCommitsController) restoreSelectionRangeAndMode(selectionRangeAndMode SelectionRangeAndMode) {
	// We need to select the same commit range again because after starting a rebase,
	// new lines can be added for update-ref commands in the TODO file, due to
	// stacked branches. So the selected commits may be in different positions in the list.
	_, newSelectedIdx, ok1 := lo.FindIndexOf(self.c.Model().Commits, func(c *models.Commit) bool {
		return c.Hash() == selectionRangeAndMode.selectedHash
	})
	_, newRangeStartIdx, ok2 := lo.FindIndexOf(self.c.Model().Commits, func(c *models.Commit) bool {
		return c.Hash() == selectionRangeAndMode.rangeStartHash
	})
	if ok1 && ok2 {
		self.context().SetSelectionRangeAndMode(newSelectedIdx, newRangeStartIdx, selectionRangeAndMode.mode)
		self.context().HandleFocus(types.OnFocusOpts{})
	}
}

func (self *LocalCommitsController) findCommitForQuickStartInteractiveRebase() (*models.Commit, error) {
	commit, index, ok := lo.FindIndexOf(self.c.Model().Commits, func(c *models.Commit) bool {
		return c.IsMerge() || c.Status == models.StatusMerged
	})

	if !ok || index == 0 {
		errorMsg := utils.ResolvePlaceholderString(self.c.Tr.CannotQuickStartInteractiveRebase, map[string]string{
			"editKey": keybindings.Label(self.c.UserConfig().Keybinding.Universal.Edit),
		})

		return nil, errors.New(errorMsg)
	}

	return commit, nil
}

func (self *LocalCommitsController) pick(selectedCommits []*models.Commit) error {
	if self.isRebasing() {
		return self.updateTodos(todo.Pick, selectedCommits)
	}

	panic("should be disabled when not rebasing")
}

func (self *LocalCommitsController) interactiveRebase(action todo.TodoCommand, startIdx int, endIdx int) error {
	// When performing an action that will remove the selected commits, we need to select the
	// next commit down (which will end up at the start index after the action is performed)
	if action == todo.Drop || action == todo.Fixup || action == todo.Squash {
		self.context().SetSelection(startIdx)
	}

	err := self.c.Git().Rebase.InteractiveRebase(self.c.Model().Commits, startIdx, endIdx, action)

	return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
}

// updateTodos sees if the selected commit is in fact a rebasing
// commit meaning you are trying to edit the todo file rather than actually
// begin a rebase. It then updates the todo file with that action
func (self *LocalCommitsController) updateTodos(action todo.TodoCommand, selectedCommits []*models.Commit) error {
	if err := self.c.Git().Rebase.EditRebaseTodo(selectedCommits, action); err != nil {
		return err
	}

	self.c.Refresh(types.RefreshOptions{
		Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
	})

	return nil
}

func (self *LocalCommitsController) rewordEnabled(commit *models.Commit) *types.DisabledReason {
	// for now we do not support setting 'reword' on TODO commits because it requires an editor
	// and that means we either unconditionally wait around for the subprocess to ask for
	// our input or we set a lazygit client as the EDITOR env variable and have it
	// request us to edit the commit message when prompted.
	if commit.IsTODO() {
		return &types.DisabledReason{Text: self.c.Tr.RewordNotSupported}
	}

	// If we are in a rebase, the only action that is allowed for
	// non-todo commits is rewording the current head commit
	if self.isRebasing() && !self.isSelectedHeadCommit() {
		return &types.DisabledReason{Text: self.c.Tr.AlreadyRebasing}
	}

	return nil
}

func (self *LocalCommitsController) isRebasing() bool {
	return self.c.Model().WorkingTreeStateAtLastCommitRefresh.Any()
}

func (self *LocalCommitsController) isCherryPickingOrReverting() bool {
	return self.c.Model().WorkingTreeStateAtLastCommitRefresh.CherryPicking ||
		self.c.Model().WorkingTreeStateAtLastCommitRefresh.Reverting
}

func (self *LocalCommitsController) moveDown(selectedCommits []*models.Commit, startIdx int, endIdx int) error {
	if self.isRebasing() {
		if err := self.c.Git().Rebase.MoveTodosDown(selectedCommits); err != nil {
			return err
		}
		self.context().MoveSelection(1)

		self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
		return nil
	}

	return self.c.WithWaitingStatusSync(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)
		err := self.c.Git().Rebase.MoveCommitsDown(self.c.Model().Commits, startIdx, endIdx)
		if err == nil {
			self.context().MoveSelection(1)
		}
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err, types.RefreshOptions{Mode: types.SYNC})
	})
}

func (self *LocalCommitsController) moveUp(selectedCommits []*models.Commit, startIdx int, endIdx int) error {
	if self.isRebasing() {
		if err := self.c.Git().Rebase.MoveTodosUp(selectedCommits); err != nil {
			return err
		}
		self.context().MoveSelection(-1)

		self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
		return nil
	}

	return self.c.WithWaitingStatusSync(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		err := self.c.Git().Rebase.MoveCommitsUp(self.c.Model().Commits, startIdx, endIdx)
		if err == nil {
			self.context().MoveSelection(-1)
		}
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err, types.RefreshOptions{Mode: types.SYNC})
	})
}

func (self *LocalCommitsController) amendTo(commit *models.Commit) error {
	var handleCommit func() error

	if self.isSelectedHeadCommit() {
		handleCommit = func() error {
			return self.c.Helpers().WorkingTree.WithEnsureCommittableFiles(func() error {
				if err := self.c.Helpers().AmendHelper.AmendHead(); err != nil {
					return err
				}
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				return nil
			})
		}
	} else {
		handleCommit = func() error {
			return self.c.Helpers().WorkingTree.WithEnsureCommittableFiles(func() error {
				return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
					self.c.LogAction(self.c.Tr.Actions.AmendCommit)
					err := self.c.Git().Rebase.AmendTo(self.c.Model().Commits, self.context().GetView().SelectedLineIdx())
					return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
				})
			})
		}
	}

	return self.c.ConfirmIf(!self.c.UserConfig().Gui.SkipAmendWarning,
		types.ConfirmOpts{
			Title:         self.c.Tr.AmendCommitTitle,
			Prompt:        self.c.Tr.AmendCommitPrompt,
			HandleConfirm: handleCommit,
		})
}

func (self *LocalCommitsController) canAmendRange(commits []*models.Commit, start, end int) *types.DisabledReason {
	if (start != end || !self.isHeadCommit(start)) && self.isRebasing() {
		return &types.DisabledReason{Text: self.c.Tr.AlreadyRebasing}
	}

	return nil
}

func (self *LocalCommitsController) canAmend(_ *models.Commit) *types.DisabledReason {
	idx := self.context().GetSelectedLineIdx()
	return self.canAmendRange(self.c.Model().Commits, idx, idx)
}

func (self *LocalCommitsController) amendAttribute(commits []*models.Commit, start, end int) error {
	opts := self.c.KeybindingsOpts()
	return self.c.Menu(types.CreateMenuOptions{
		Title: "Amend commit attribute",
		Items: []*types.MenuItem{
			{
				Label:   self.c.Tr.ResetAuthor,
				OnPress: func() error { return self.resetAuthor(start, end) },
				Key:     opts.GetKey(opts.Config.AmendAttribute.ResetAuthor),
				Tooltip: self.c.Tr.ResetAuthorTooltip,
			},
			{
				Label:   self.c.Tr.SetAuthor,
				OnPress: func() error { return self.setAuthor(start, end) },
				Key:     opts.GetKey(opts.Config.AmendAttribute.SetAuthor),
				Tooltip: self.c.Tr.SetAuthorTooltip,
			},
			{
				Label:   self.c.Tr.AddCoAuthor,
				OnPress: func() error { return self.addCoAuthor(start, end) },
				Key:     opts.GetKey(opts.Config.AmendAttribute.AddCoAuthor),
				Tooltip: self.c.Tr.AddCoAuthorTooltip,
			},
		},
	})
}

func (self *LocalCommitsController) resetAuthor(start, end int) error {
	return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.ResetCommitAuthor)
		if err := self.c.Git().Rebase.ResetCommitAuthor(self.c.Model().Commits, start, end); err != nil {
			return err
		}

		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		return nil
	})
}

func (self *LocalCommitsController) setAuthor(start, end int) error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.SetAuthorPromptTitle,
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc(),
		HandleConfirm: func(value string) error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.SetCommitAuthor)
				if err := self.c.Git().Rebase.SetCommitAuthor(self.c.Model().Commits, start, end, value); err != nil {
					return err
				}

				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				return nil
			})
		},
	})

	return nil
}

func (self *LocalCommitsController) addCoAuthor(start, end int) error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.AddCoAuthorPromptTitle,
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetAuthorsSuggestionsFunc(),
		HandleConfirm: func(value string) error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.AddCommitCoAuthor)
				if err := self.c.Git().Rebase.AddCommitCoAuthor(self.c.Model().Commits, start, end, value); err != nil {
					return err
				}
				self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
				return nil
			})
		},
	})

	return nil
}

func (self *LocalCommitsController) revert(commits []*models.Commit, start, end int) error {
	var promptText string
	if len(commits) == 1 {
		promptText = utils.ResolvePlaceholderString(
			self.c.Tr.ConfirmRevertCommit,
			map[string]string{
				"selectedCommit": commits[0].ShortHash(),
			})
	} else {
		promptText = self.c.Tr.ConfirmRevertCommitRange
	}
	hashes := lo.Map(commits, func(c *models.Commit, _ int) string { return c.Hash() })
	isMerge := lo.SomeBy(commits, func(c *models.Commit) bool { return c.IsMerge() })

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Actions.RevertCommit,
		Prompt: promptText,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.RevertCommit)
			return self.c.WithWaitingStatusSync(self.c.Tr.RevertingStatus, func() error {
				mustStash := helpers.IsWorkingTreeDirty(self.c.Model().Files)

				if mustStash {
					if err := self.c.Git().Stash.Push(self.c.Tr.AutoStashForReverting); err != nil {
						return err
					}
				}

				result := self.c.Git().Commit.Revert(hashes, isMerge)
				if err := self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(result, types.RefreshOptions{Mode: types.SYNC}); err != nil {
					return err
				}
				self.context().MoveSelection(len(commits))
				self.context().FocusLine()

				if mustStash {
					if err := self.c.Git().Stash.Pop(0); err != nil {
						return err
					}
					self.c.Refresh(types.RefreshOptions{
						Scope: []types.RefreshableView{types.STASH, types.FILES},
					})
				}

				return nil
			})
		},
	})

	return nil
}

func (self *LocalCommitsController) createFixupCommit(commit *models.Commit) error {
	var disabledReasonWhenFilesAreNeeded *types.DisabledReason
	if len(self.c.Model().Files) == 0 {
		disabledReasonWhenFilesAreNeeded = &types.DisabledReason{
			Text:             self.c.Tr.NoFilesStagedTitle,
			ShowErrorInPanel: true,
		}
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CreateFixupCommit,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.FixupMenu_Fixup,
				Key:   'f',
				OnPress: func() error {
					return self.c.Helpers().WorkingTree.WithEnsureCommittableFiles(func() error {
						self.c.LogAction(self.c.Tr.Actions.CreateFixupCommit)
						return self.c.WithWaitingStatusSync(self.c.Tr.CreatingFixupCommitStatus, func() error {
							if err := self.c.Git().Commit.CreateFixupCommit(commit.Hash()); err != nil {
								return err
							}

							if err := self.moveFixupCommitToOwnerStackedBranch(commit); err != nil {
								return err
							}

							self.context().MoveSelectedLine(1)
							self.c.Refresh(types.RefreshOptions{Mode: types.SYNC})
							return nil
						})
					})
				},
				DisabledReason: disabledReasonWhenFilesAreNeeded,
				Tooltip:        self.c.Tr.FixupMenu_FixupTooltip,
			},
			{
				Label: self.c.Tr.FixupMenu_AmendWithChanges,
				Key:   'a',
				OnPress: func() error {
					return self.c.Helpers().WorkingTree.WithEnsureCommittableFiles(func() error {
						return self.createAmendCommit(commit, true)
					})
				},
				DisabledReason: disabledReasonWhenFilesAreNeeded,
				Tooltip:        self.c.Tr.FixupMenu_AmendWithChangesTooltip,
			},
			{
				Label:   self.c.Tr.FixupMenu_AmendWithoutChanges,
				Key:     'r',
				OnPress: func() error { return self.createAmendCommit(commit, false) },
				Tooltip: self.c.Tr.FixupMenu_AmendWithoutChangesTooltip,
			},
		},
	})
}

func (self *LocalCommitsController) moveFixupCommitToOwnerStackedBranch(targetCommit *models.Commit) error {
	if self.c.Git().Version.IsOlderThan(2, 38, 0) {
		// Git 2.38.0 introduced the `rebase.updateRefs` config option. Don't
		// move the commit down with older versions, as it would break the stack.
		return nil
	}

	if self.c.Git().Status.WorkingTreeState().Any() {
		// Can't move commits while rebasing
		return nil
	}

	if targetCommit.Status == models.StatusMerged {
		// Target commit is already on main. It's a bit questionable that we
		// allow creating a fixup commit for it in the first place, but we
		// always did, so why restrict that now; however, it doesn't make sense
		// to move the created fixup commit down in that case.
		return nil
	}

	if !self.c.Git().Config.GetRebaseUpdateRefs() {
		// If the user has disabled rebase.updateRefs, we don't move the fixup
		// because this would break the stack of branches (presumably they like
		// to manage it themselves manually, or something).
		return nil
	}

	headOfOwnerBranchIdx := -1
	for i := self.context().GetSelectedLineIdx(); i > 0; i-- {
		if lo.SomeBy(self.c.Model().Branches, func(b *models.Branch) bool {
			return b.CommitHash == self.c.Model().Commits[i].Hash()
		}) {
			headOfOwnerBranchIdx = i
			break
		}
	}

	if headOfOwnerBranchIdx == -1 {
		return nil
	}

	return self.c.Git().Rebase.MoveFixupCommitDown(self.c.Model().Commits, headOfOwnerBranchIdx)
}

func (self *LocalCommitsController) createAmendCommit(commit *models.Commit, includeFileChanges bool) error {
	commitMessage, err := self.c.Git().Commit.GetCommitMessage(commit.Hash())
	if err != nil {
		return err
	}
	if self.c.UserConfig().Git.Commit.AutoWrapCommitMessage {
		commitMessage = helpers.TryRemoveHardLineBreaks(commitMessage, self.c.UserConfig().Git.Commit.AutoWrapWidth)
	}
	originalSubject, _, _ := strings.Cut(commitMessage, "\n")
	self.c.Helpers().Commits.OpenCommitMessagePanel(
		&helpers.OpenCommitMessagePanelOpts{
			CommitIndex:      self.context().GetSelectedLineIdx(),
			InitialMessage:   commitMessage,
			SummaryTitle:     self.c.Tr.CreateAmendCommit,
			DescriptionTitle: self.c.Tr.CommitDescriptionTitle,
			PreserveMessage:  false,
			OnConfirm: func(summary string, description string) error {
				self.c.LogAction(self.c.Tr.Actions.CreateFixupCommit)
				return self.c.WithWaitingStatusSync(self.c.Tr.CreatingFixupCommitStatus, func() error {
					if err := self.c.Git().Commit.CreateAmendCommit(originalSubject, summary, description, includeFileChanges); err != nil {
						return err
					}

					if err := self.moveFixupCommitToOwnerStackedBranch(commit); err != nil {
						return err
					}

					self.context().MoveSelectedLine(1)
					self.c.Refresh(types.RefreshOptions{Mode: types.SYNC})
					return nil
				})
			},
			OnSwitchToEditor: nil,
		},
	)

	return nil
}

func (self *LocalCommitsController) squashFixupCommits() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SquashAboveCommits,
		Items: []*types.MenuItem{
			{
				Label:          self.c.Tr.SquashCommitsInCurrentBranch,
				OnPress:        self.squashAllFixupsInCurrentBranch,
				DisabledReason: self.canFindCommitForSquashFixupsInCurrentBranch(),
				Key:            'b',
				Tooltip:        self.c.Tr.SquashCommitsInCurrentBranchTooltip,
			},
			{
				Label:          self.c.Tr.SquashCommitsAboveSelectedCommit,
				OnPress:        self.withItem(self.squashAllFixupsAboveSelectedCommit),
				DisabledReason: self.singleItemSelected()(),
				Key:            'a',
				Tooltip:        self.c.Tr.SquashCommitsAboveSelectedTooltip,
			},
		},
	})
}

func (self *LocalCommitsController) squashAllFixupsAboveSelectedCommit(commit *models.Commit) error {
	return self.squashFixupsImpl(commit, self.context().GetSelectedLineIdx())
}

func (self *LocalCommitsController) squashAllFixupsInCurrentBranch() error {
	commit, rebaseStartIdx, err := self.findCommitForSquashFixupsInCurrentBranch()
	if err != nil {
		return err
	}

	return self.squashFixupsImpl(commit, rebaseStartIdx)
}

func (self *LocalCommitsController) squashFixupsImpl(commit *models.Commit, rebaseStartIdx int) error {
	selectionOffset := countSquashableCommitsAbove(self.c.Model().Commits, self.context().GetSelectedLineIdx(), rebaseStartIdx)
	return self.c.WithWaitingStatusSync(self.c.Tr.SquashingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.SquashAllAboveFixupCommits)
		err := self.c.Git().Rebase.SquashAllAboveFixupCommits(commit)
		self.context().MoveSelectedLine(-selectionOffset)
		return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseWithRefreshOptions(
			err, types.RefreshOptions{Mode: types.SYNC})
	})
}

func (self *LocalCommitsController) findCommitForSquashFixupsInCurrentBranch() (*models.Commit, int, error) {
	commits := self.c.Model().Commits
	_, index, ok := lo.FindIndexOf(commits, func(c *models.Commit) bool {
		return c.IsMerge() || c.Status == models.StatusMerged
	})

	if !ok || index == 0 {
		return nil, -1, errors.New(self.c.Tr.CannotSquashCommitsInCurrentBranch)
	}

	return commits[index-1], index - 1, nil
}

// Anticipate how many commits above the selectedIdx are going to get squashed
// by the SquashAllAboveFixupCommits call, so that we can adjust the selection
// afterwards. Let's hope we're matching git's behavior correctly here.
func countSquashableCommitsAbove(commits []*models.Commit, selectedIdx int, rebaseStartIdx int) int {
	result := 0

	// For each commit _above_ the selection, ...
	for i, commit := range commits[0:selectedIdx] {
		// ... see if it is a fixup commit, and get the base subject it applies to
		if baseSubject, isFixup := isFixupCommit(commit.Name); isFixup {
			// Then, for each commit after the fixup, up to and including the
			// rebase start commit, see if we find the base commit
			for _, baseCommit := range commits[i+1 : rebaseStartIdx+1] {
				if strings.HasPrefix(baseCommit.Name, baseSubject) {
					result++
				}
			}
		}
	}
	return result
}

// Check whether the given subject line is the subject of a fixup commit, and
// returns (trimmedSubject, true) if so (where trimmedSubject is the subject
// with all fixup prefixes removed), or (subject, false) if not.
func isFixupCommit(subject string) (string, bool) {
	prefixes := []string{"fixup! ", "squash! ", "amend! "}
	trimPrefix := func(s string) (string, bool) {
		for _, prefix := range prefixes {
			if strings.HasPrefix(s, prefix) {
				return strings.TrimPrefix(s, prefix), true
			}
		}
		return s, false
	}

	if subject, wasTrimmed := trimPrefix(subject); wasTrimmed {
		for {
			// handle repeated prefixes like "fixup! amend! fixup! Subject"
			if subject, wasTrimmed = trimPrefix(subject); !wasTrimmed {
				break
			}
		}
		return subject, true
	}

	return subject, false
}

func (self *LocalCommitsController) createTag(commit *models.Commit) error {
	return self.c.Helpers().Tags.OpenCreateTagPrompt(commit.Hash(), func() {})
}

func (self *LocalCommitsController) openSearch() error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if self.context().GetLimitCommits() {
		self.context().SetLimitCommits(false)
		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS}})
	}

	return self.c.Helpers().Search.OpenSearchPrompt(self.context())
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
						self.c.Refresh(
							types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}},
						)
						return nil
					})
				},
			},
			{
				Label:     self.c.Tr.ShowGitGraph,
				Tooltip:   self.c.Tr.ShowGitGraphTooltip,
				OpensMenu: true,
				OnPress: func() error {
					currentValue := self.c.UserConfig().Git.Log.ShowGraph
					onPress := func(value string) func() error {
						return func() error {
							self.c.UserConfig().Git.Log.ShowGraph = value
							self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
							self.c.PostRefreshUpdate(self.c.Contexts().SubCommits)
							return nil
						}
					}
					return self.c.Menu(types.CreateMenuOptions{
						Title: self.c.Tr.LogMenuTitle,
						Items: []*types.MenuItem{
							{
								Label:   "always",
								OnPress: onPress("always"),
								Widget:  types.MakeMenuRadioButton(currentValue == "always"),
							},
							{
								Label:   "never",
								OnPress: onPress("never"),
								Widget:  types.MakeMenuRadioButton(currentValue == "never"),
							},
							{
								Label:   "when maximised",
								OnPress: onPress("when-maximised"),
								Widget:  types.MakeMenuRadioButton(currentValue == "when-maximised"),
							},
						},
					})
				},
			},
			{
				Label:     self.c.Tr.SortCommits,
				Tooltip:   self.c.Tr.SortCommitsTooltip,
				OpensMenu: true,
				OnPress: func() error {
					currentValue := self.c.UserConfig().Git.Log.Order
					onPress := func(value string) func() error {
						return func() error {
							self.c.UserConfig().Git.Log.Order = value
							return self.c.WithWaitingStatus(self.c.Tr.LoadingCommits, func(gocui.Task) error {
								self.c.Refresh(
									types.RefreshOptions{
										Mode:  types.SYNC,
										Scope: []types.RefreshableView{types.COMMITS},
									},
								)
								return nil
							})
						}
					}

					return self.c.Menu(types.CreateMenuOptions{
						Title: self.c.Tr.LogMenuTitle,
						Items: []*types.MenuItem{
							{
								Label:   "topological (topo-order)",
								OnPress: onPress("topo-order"),
								Widget:  types.MakeMenuRadioButton(currentValue == "topo-order"),
							},
							{
								Label:   "date-order",
								OnPress: onPress("date-order"),
								Widget:  types.MakeMenuRadioButton(currentValue == "date-order"),
							},
							{
								Label:   "author-date-order",
								OnPress: onPress("author-date-order"),
								Widget:  types.MakeMenuRadioButton(currentValue == "author-date-order"),
							},
							{
								Label:   "default",
								OnPress: onPress("default"),
								Widget:  types.MakeMenuRadioButton(currentValue == "default"),
							},
						},
					})
				},
			},
		},
	})
}

func (self *LocalCommitsController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		context := self.context()
		if context.GetSelectedLineIdx() > COMMIT_THRESHOLD && context.GetLimitCommits() {
			context.SetLimitCommits(false)
			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS}})
		}
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
	if commit.Hash() == self.c.Modes().MarkedBaseCommit.GetHash() {
		// Reset when invoking it again on the marked commit
		self.c.Modes().MarkedBaseCommit.SetHash("")
	} else {
		self.c.Modes().MarkedBaseCommit.SetHash(commit.Hash())
	}
	self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
	return nil
}

func (self *LocalCommitsController) isHeadCommit(idx int) bool {
	return models.IsHeadCommit(self.c.Model().Commits, idx)
}

func (self *LocalCommitsController) isSelectedHeadCommit() bool {
	return self.isHeadCommit(self.context().GetSelectedLineIdx())
}

func (self *LocalCommitsController) notMidRebase(message string) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		if self.isRebasing() {
			return &types.DisabledReason{Text: message}
		}

		return nil
	}
}

func (self *LocalCommitsController) canFindCommitForQuickStart() *types.DisabledReason {
	if _, err := self.findCommitForQuickStartInteractiveRebase(); err != nil {
		return &types.DisabledReason{Text: err.Error(), ShowErrorInPanel: true}
	}

	return nil
}

func (self *LocalCommitsController) canFindCommitForSquashFixupsInCurrentBranch() *types.DisabledReason {
	if _, _, err := self.findCommitForSquashFixupsInCurrentBranch(); err != nil {
		return &types.DisabledReason{Text: err.Error()}
	}

	return nil
}

func (self *LocalCommitsController) canSquashOrFixup(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	if endIdx >= len(self.c.Model().Commits)-1 {
		return &types.DisabledReason{Text: self.c.Tr.CannotSquashOrFixupFirstCommit}
	}

	if lo.SomeBy(selectedCommits, func(c *models.Commit) bool { return c.IsMerge() }) {
		return &types.DisabledReason{Text: self.c.Tr.CannotSquashOrFixupMergeCommit}
	}

	return nil
}

func (self *LocalCommitsController) canMoveDown(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	if endIdx >= len(self.c.Model().Commits)-1 {
		return &types.DisabledReason{Text: self.c.Tr.CannotMoveAnyFurther}
	}

	if self.isRebasing() {
		commits := self.c.Model().Commits

		if !commits[endIdx+1].IsTODO() || commits[endIdx+1].Status == models.StatusConflicted {
			return &types.DisabledReason{Text: self.c.Tr.CannotMoveAnyFurther}
		}
	}

	return nil
}

func (self *LocalCommitsController) canMoveUp(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	if startIdx == 0 {
		return &types.DisabledReason{Text: self.c.Tr.CannotMoveAnyFurther}
	}

	if self.isRebasing() {
		commits := self.c.Model().Commits

		if !commits[startIdx-1].IsTODO() || commits[startIdx-1].Status == models.StatusConflicted {
			return &types.DisabledReason{Text: self.c.Tr.CannotMoveAnyFurther}
		}
	}

	return nil
}

// Ensures that if we are mid-rebase, we're only selecting valid commits (non-conflict TODO commits)
func (self *LocalCommitsController) midRebaseCommandEnabled(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	if self.isCherryPickingOrReverting() {
		return &types.DisabledReason{Text: self.c.Tr.NotAllowedMidCherryPickOrRevert}
	}

	if !self.isRebasing() {
		return nil
	}

	for _, commit := range selectedCommits {
		if !commit.IsTODO() {
			return &types.DisabledReason{Text: self.c.Tr.MustSelectTodoCommits}
		}

		if !isChangeOfRebaseTodoAllowed(commit.Action) {
			return &types.DisabledReason{Text: self.c.Tr.ChangingThisActionIsNotAllowed}
		}
	}

	return nil
}

// Ensures that if we are mid-rebase, we're only selecting commits that can be moved
func (self *LocalCommitsController) midRebaseMoveCommandEnabled(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	if self.isCherryPickingOrReverting() {
		return &types.DisabledReason{Text: self.c.Tr.NotAllowedMidCherryPickOrRevert}
	}

	if !self.isRebasing() {
		if lo.SomeBy(selectedCommits, func(c *models.Commit) bool { return c.IsMerge() }) {
			return &types.DisabledReason{Text: self.c.Tr.CannotMoveMergeCommit}
		}

		return nil
	}

	for _, commit := range selectedCommits {
		if !commit.IsTODO() {
			return &types.DisabledReason{Text: self.c.Tr.MustSelectTodoCommits}
		}

		// All todo types that can be edited are allowed to be moved, plus
		// update-ref todos
		if !isChangeOfRebaseTodoAllowed(commit.Action) && commit.Action != todo.UpdateRef {
			return &types.DisabledReason{Text: self.c.Tr.ChangingThisActionIsNotAllowed}
		}
	}

	return nil
}

func (self *LocalCommitsController) canDropCommits(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	if self.isCherryPickingOrReverting() {
		return &types.DisabledReason{Text: self.c.Tr.NotAllowedMidCherryPickOrRevert}
	}

	if !self.isRebasing() {
		if len(selectedCommits) > 1 && lo.SomeBy(selectedCommits, func(c *models.Commit) bool { return c.IsMerge() }) {
			return &types.DisabledReason{Text: self.c.Tr.DroppingMergeRequiresSingleSelection}
		}

		return nil
	}

	nonUpdateRefTodos := lo.Filter(selectedCommits, func(c *models.Commit, _ int) bool {
		return c.Action != todo.UpdateRef
	})

	for _, commit := range nonUpdateRefTodos {
		if !commit.IsTODO() {
			return &types.DisabledReason{Text: self.c.Tr.MustSelectTodoCommits}
		}

		if !isChangeOfRebaseTodoAllowed(commit.Action) {
			return &types.DisabledReason{Text: self.c.Tr.ChangingThisActionIsNotAllowed}
		}
	}

	return nil
}

// These actions represent standard things you might want to do with a commit,
// as opposed to TODO actions like 'merge', 'update-ref', etc.
var standardActions = []todo.TodoCommand{
	todo.Pick,
	todo.Drop,
	todo.Edit,
	todo.Fixup,
	todo.Squash,
	todo.Reword,
}

func isChangeOfRebaseTodoAllowed(oldAction todo.TodoCommand) bool {
	// Only allow updating a standard action, meaning we disallow
	// updating a merge commit or update ref commit (until we decide what would be sensible
	// to do in those cases)
	return lo.Contains(standardActions, oldAction)
}

func (self *LocalCommitsController) pickEnabled(selectedCommits []*models.Commit, startIdx int, endIdx int) *types.DisabledReason {
	if self.isCherryPickingOrReverting() {
		return &types.DisabledReason{Text: self.c.Tr.NotAllowedMidCherryPickOrRevert}
	}

	if !self.isRebasing() {
		return &types.DisabledReason{Text: self.c.Tr.PickIsOnlyAllowedDuringRebase, AllowFurtherDispatching: true}
	}

	return self.midRebaseCommandEnabled(selectedCommits, startIdx, endIdx)
}
