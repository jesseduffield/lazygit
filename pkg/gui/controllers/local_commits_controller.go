package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type (
	PullFilesFn func() error
)

type LocalCommitsController struct {
	baseController
	*controllerCommon

	pullFiles PullFilesFn
}

var _ types.IController = &LocalCommitsController{}

func NewLocalCommitsController(
	common *controllerCommon,
	pullFiles PullFilesFn,
) *LocalCommitsController {
	return &LocalCommitsController{
		baseController:   baseController{},
		controllerCommon: common,
		pullFiles:        pullFiles,
	}
}

func (self *LocalCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	outsideFilterModeBindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Commits.SquashDown),
			Handler:     self.checkSelected(self.squashDown),
			Description: self.c.Tr.LcSquashDown,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.MarkCommitAsFixup),
			Handler:     self.checkSelected(self.fixup),
			Description: self.c.Tr.LcFixupCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.RenameCommit),
			Handler:     self.checkSelected(self.reword),
			Description: self.c.Tr.LcRewordCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.RenameCommitWithEditor),
			Handler:     self.checkSelected(self.rewordEditor),
			Description: self.c.Tr.LcRenameCommitEditor,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.drop),
			Description: self.c.Tr.LcDeleteCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelected(self.edit),
			Description: self.c.Tr.LcEditCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.PickCommit),
			Handler:     self.checkSelected(self.pick),
			Description: self.c.Tr.LcPickCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.CreateFixupCommit),
			Handler:     self.checkSelected(self.createFixupCommit),
			Description: self.c.Tr.LcCreateFixupCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.SquashAboveCommits),
			Handler:     self.checkSelected(self.squashAllAboveFixupCommits),
			Description: self.c.Tr.LcSquashAboveCommits,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.MoveDownCommit),
			Handler:     self.checkSelected(self.moveDown),
			Description: self.c.Tr.LcMoveDownCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.MoveUpCommit),
			Handler:     self.checkSelected(self.moveUp),
			Description: self.c.Tr.LcMoveUpCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.PasteCommits),
			Handler:     opts.Guards.OutsideFilterMode(self.paste),
			Description: self.c.Tr.LcPasteCommits,
		},
		// overriding these navigation keybindings because we might need to load
		// more commits on demand
		{
			Key:         opts.GetKey(opts.Config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GotoBottom),
			Handler:     self.gotoBottom,
			Description: self.c.Tr.LcGotoBottom,
			Tag:         "navigation",
		},
	}

	for _, binding := range outsideFilterModeBindings {
		binding.Handler = opts.Guards.OutsideFilterMode(binding.Handler)
	}

	bindings := append(outsideFilterModeBindings, []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Commits.AmendToCommit),
			Handler:     self.checkSelected(self.amendTo),
			Description: self.c.Tr.LcAmendToCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ResetCommitAuthor),
			Handler:     self.checkSelected(self.amendAttribute),
			Description: self.c.Tr.LcResetCommitAuthor,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.RevertCommit),
			Handler:     self.checkSelected(self.revert),
			Description: self.c.Tr.LcRevertCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.TagCommit),
			Handler:     self.checkSelected(self.createTag),
			Description: self.c.Tr.LcTagCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.OpenLogMenu),
			Handler:     self.handleOpenLogMenu,
			Description: self.c.Tr.LcOpenLogMenu,
			OpensMenu:   true,
		},
	}...)

	return bindings
}

func (self *LocalCommitsController) squashDown(commit *models.Commit) error {
	if len(self.model.Commits) <= 1 {
		return self.c.ErrorMsg(self.c.Tr.YouNoCommitsToSquash)
	}

	applied, err := self.handleMidRebaseCommand("squash", commit)
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
			return self.c.WithWaitingStatus(self.c.Tr.SquashingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.SquashCommitDown)
				return self.interactiveRebase("squash")
			})
		},
	})
}

func (self *LocalCommitsController) fixup(commit *models.Commit) error {
	if len(self.model.Commits) <= 1 {
		return self.c.ErrorMsg(self.c.Tr.YouNoCommitsToSquash)
	}

	applied, err := self.handleMidRebaseCommand("fixup", commit)
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
			return self.c.WithWaitingStatus(self.c.Tr.FixingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.FixupCommit)
				return self.interactiveRebase("fixup")
			})
		},
	})
}

func (self *LocalCommitsController) reword(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand("reword", commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	message, err := self.git.Commit.GetCommitMessage(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	// TODO: use the commit message panel here
	return self.c.Prompt(types.PromptOpts{
		Title:          self.c.Tr.LcRewordCommit,
		InitialContent: message,
		HandleConfirm: func(response string) error {
			self.c.LogAction(self.c.Tr.Actions.RewordCommit)
			if err := self.git.Rebase.RewordCommit(self.model.Commits, self.context().GetSelectedLineIdx(), response); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

func (self *LocalCommitsController) rewordEditor(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand("reword", commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.RewordInEditorTitle,
		Prompt: self.c.Tr.RewordInEditorPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.RewordCommit)
			subProcess, err := self.git.Rebase.RewordCommitInEditor(
				self.model.Commits, self.context().GetSelectedLineIdx(),
			)
			if err != nil {
				return self.c.Error(err)
			}
			if subProcess != nil {
				return self.c.RunSubprocessAndRefresh(subProcess)
			}

			return nil
		},
	})
}

func (self *LocalCommitsController) drop(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand("drop", commit)
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
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.DropCommit)
				return self.interactiveRebase("drop")
			})
		},
	})
}

func (self *LocalCommitsController) edit(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand("edit", commit)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.EditCommit)
		return self.interactiveRebase("edit")
	})
}

func (self *LocalCommitsController) pick(commit *models.Commit) error {
	applied, err := self.handleMidRebaseCommand("pick", commit)
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

func (self *LocalCommitsController) interactiveRebase(action string) error {
	err := self.git.Rebase.InteractiveRebase(self.model.Commits, self.context().GetSelectedLineIdx(), action)
	return self.helpers.MergeAndRebase.CheckMergeOrRebase(err)
}

// handleMidRebaseCommand sees if the selected commit is in fact a rebasing
// commit meaning you are trying to edit the todo file rather than actually
// begin a rebase. It then updates the todo file with that action
func (self *LocalCommitsController) handleMidRebaseCommand(action string, commit *models.Commit) (bool, error) {
	if commit.Status != "rebasing" {
		return false, nil
	}

	// for now we do not support setting 'reword' because it requires an editor
	// and that means we either unconditionally wait around for the subprocess to ask for
	// our input or we set a lazygit client as the EDITOR env variable and have it
	// request us to edit the commit message when prompted.
	if action == "reword" {
		return true, self.c.ErrorMsg(self.c.Tr.LcRewordNotSupported)
	}

	self.c.LogAction("Update rebase TODO")
	self.c.LogCommand(
		fmt.Sprintf("Updating rebase action of commit %s to '%s'", commit.ShortSha(), action),
		false,
	)

	if err := self.git.Rebase.EditRebaseTodo(
		self.context().GetSelectedLineIdx(), action,
	); err != nil {
		return false, self.c.Error(err)
	}

	return true, self.c.Refresh(types.RefreshOptions{
		Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
	})
}

func (self *LocalCommitsController) moveDown(commit *models.Commit) error {
	index := self.context().GetSelectedLineIdx()
	commits := self.model.Commits
	if commit.Status == "rebasing" {
		if commits[index+1].Status != "rebasing" {
			return nil
		}

		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)
		self.c.LogCommand(fmt.Sprintf("Moving commit %s down", commit.ShortSha()), false)

		if err := self.git.Rebase.MoveTodoDown(index); err != nil {
			return self.c.Error(err)
		}
		self.context().MoveSelectedLine(1)
		return self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
	}

	return self.c.WithWaitingStatus(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)
		err := self.git.Rebase.MoveCommitDown(self.model.Commits, index)
		if err == nil {
			self.context().MoveSelectedLine(1)
		}
		return self.helpers.MergeAndRebase.CheckMergeOrRebase(err)
	})
}

func (self *LocalCommitsController) moveUp(commit *models.Commit) error {
	index := self.context().GetSelectedLineIdx()
	if index == 0 {
		return nil
	}

	if commit.Status == "rebasing" {
		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		self.c.LogCommand(
			fmt.Sprintf("Moving commit %s up", commit.ShortSha()),
			false,
		)

		if err := self.git.Rebase.MoveTodoDown(index - 1); err != nil {
			return self.c.Error(err)
		}
		self.context().MoveSelectedLine(-1)
		return self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
	}

	return self.c.WithWaitingStatus(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		err := self.git.Rebase.MoveCommitDown(self.model.Commits, index-1)
		if err == nil {
			self.context().MoveSelectedLine(-1)
		}
		return self.helpers.MergeAndRebase.CheckMergeOrRebase(err)
	})
}

func (self *LocalCommitsController) amendTo(commit *models.Commit) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AmendCommitTitle,
		Prompt: self.c.Tr.AmendCommitPrompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.AmendCommit)
				err := self.git.Rebase.AmendTo(commit.Sha)
				return self.helpers.MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
}

func (self *LocalCommitsController) amendAttribute(commit *models.Commit) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: "Amend commit attribute",
		Items: []*types.MenuItem{
			{
				Label:   "reset author",
				OnPress: self.resetAuthor,
				Key:     'a',
				Tooltip: "Reset the commit's author to the currently configured user. This will also renew the author timestamp",
			},
			{
				Label:   "set author",
				OnPress: self.setAuthor,
				Key:     'A',
				Tooltip: "Set the author based on a prompt",
			},
		},
	})
}

func (self *LocalCommitsController) resetAuthor() error {
	return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.ResetCommitAuthor)
		if err := self.git.Rebase.ResetCommitAuthor(self.model.Commits, self.context().GetSelectedLineIdx()); err != nil {
			return self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *LocalCommitsController) setAuthor() error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.SetAuthorPromptTitle,
		FindSuggestionsFunc: self.helpers.Suggestions.GetAuthorsSuggestionsFunc(),
		HandleConfirm: func(value string) error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.SetCommitAuthor)
				if err := self.git.Rebase.SetCommitAuthor(self.model.Commits, self.context().GetSelectedLineIdx(), value); err != nil {
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
				if err := self.git.Commit.Revert(commit.Sha); err != nil {
					return self.c.Error(err)
				}
				return self.afterRevertCommit()
			},
		})
	}
}

func (self *LocalCommitsController) createRevertMergeCommitMenu(commit *models.Commit) error {
	menuItems := make([]*types.MenuItem, len(commit.Parents))
	for i, parentSha := range commit.Parents {
		i := i
		message, err := self.git.Commit.GetCommitMessageFirstLine(parentSha)
		if err != nil {
			return self.c.Error(err)
		}

		menuItems[i] = &types.MenuItem{
			Label: fmt.Sprintf("%s: %s", utils.SafeTruncate(parentSha, 8), message),
			OnPress: func() error {
				parentNumber := i + 1
				self.c.LogAction(self.c.Tr.Actions.RevertCommit)
				if err := self.git.Commit.RevertMerge(commit.Sha, parentNumber); err != nil {
					return self.c.Error(err)
				}
				return self.afterRevertCommit()
			},
		}
	}

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.SelectParentCommitForMerge, Items: menuItems})
}

func (self *LocalCommitsController) afterRevertCommit() error {
	self.context().MoveSelectedLine(1)
	return self.c.Refresh(types.RefreshOptions{
		Mode: types.BLOCK_UI, Scope: []types.RefreshableView{types.COMMITS, types.BRANCHES},
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
			self.c.LogAction(self.c.Tr.Actions.CreateFixupCommit)
			if err := self.git.Commit.CreateFixupCommit(commit.Sha); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
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
			return self.c.WithWaitingStatus(self.c.Tr.SquashingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.SquashAllAboveFixupCommits)
				err := self.git.Rebase.SquashAllAboveFixupCommits(commit.Sha)
				return self.helpers.MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
}

func (self *LocalCommitsController) createTag(commit *models.Commit) error {
	return self.helpers.Tags.CreateTagMenu(commit.Sha, func() {})
}

func (self *LocalCommitsController) openSearch() error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if self.context().GetLimitCommits() {
		self.context().SetLimitCommits(false)
		if err := self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
			return err
		}
	}

	self.c.OpenSearch()

	return nil
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

					return self.c.WithWaitingStatus(self.c.Tr.LcLoadingCommits, func() error {
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
							return self.c.WithWaitingStatus(self.c.Tr.LcLoadingCommits, func() error {
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

func (self *LocalCommitsController) checkSelected(callback func(*models.Commit) error) func() error {
	return func() error {
		commit := self.context().GetSelected()
		if commit == nil {
			return nil
		}

		return callback(commit)
	}
}

func (self *LocalCommitsController) Context() types.Context {
	return self.context()
}

func (self *LocalCommitsController) context() *context.LocalCommitsContext {
	return self.contexts.LocalCommits
}

func (self *LocalCommitsController) paste() error {
	return self.helpers.CherryPick.Paste()
}
