package controllers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type (
	CheckoutRefFn                func(refName string, opts types.CheckoutRefOptions) error
	CreateGitResetMenuFn         func(refName string) error
	SwitchToCommitFilesContextFn func(SwitchToCommitFilesContextOpts) error
	GetHostingServiceMgrFn       func() *hosting_service.HostingServiceMgr
	PullFilesFn                  func() error
	CheckMergeOrRebase           func(error) error
)

type LocalCommitsController struct {
	c                *types.ControllerCommon
	getContext       func() types.IListContext
	os               *oscommands.OSCommand
	git              *commands.GitCommand
	tagsHelper       *TagsHelper
	refsHelper       IRefsHelper
	cherryPickHelper *CherryPickHelper
	rebaseHelper     *RebaseHelper

	getSelectedLocalCommit     func() *models.Commit
	getCommits                 func() []*models.Commit
	getSelectedLocalCommitIdx  func() int
	CheckMergeOrRebase         CheckMergeOrRebase
	pullFiles                  PullFilesFn
	getHostingServiceMgr       GetHostingServiceMgrFn
	switchToCommitFilesContext SwitchToCommitFilesContextFn
	getLimitCommits            func() bool
	setLimitCommits            func(bool)
	getShowWholeGitGraph       func() bool
	setShowWholeGitGraph       func(bool)
}

var _ types.IController = &LocalCommitsController{}

func NewLocalCommitsController(
	c *types.ControllerCommon,
	getContext func() types.IListContext,
	os *oscommands.OSCommand,
	git *commands.GitCommand,
	tagsHelper *TagsHelper,
	refsHelper IRefsHelper,
	cherryPickHelper *CherryPickHelper,
	rebaseHelper *RebaseHelper,
	getSelectedLocalCommit func() *models.Commit,
	getCommits func() []*models.Commit,
	getSelectedLocalCommitIdx func() int,
	CheckMergeOrRebase CheckMergeOrRebase,
	pullFiles PullFilesFn,
	getHostingServiceMgr GetHostingServiceMgrFn,
	switchToCommitFilesContext SwitchToCommitFilesContextFn,
	getLimitCommits func() bool,
	setLimitCommits func(bool),
	getShowWholeGitGraph func() bool,
	setShowWholeGitGraph func(bool),
) *LocalCommitsController {
	return &LocalCommitsController{
		c:                          c,
		getContext:                 getContext,
		os:                         os,
		git:                        git,
		tagsHelper:                 tagsHelper,
		refsHelper:                 refsHelper,
		cherryPickHelper:           cherryPickHelper,
		rebaseHelper:               rebaseHelper,
		getSelectedLocalCommit:     getSelectedLocalCommit,
		getCommits:                 getCommits,
		getSelectedLocalCommitIdx:  getSelectedLocalCommitIdx,
		CheckMergeOrRebase:         CheckMergeOrRebase,
		pullFiles:                  pullFiles,
		getHostingServiceMgr:       getHostingServiceMgr,
		switchToCommitFilesContext: switchToCommitFilesContext,
		getLimitCommits:            getLimitCommits,
		setLimitCommits:            setLimitCommits,
		getShowWholeGitGraph:       getShowWholeGitGraph,
		setShowWholeGitGraph:       setShowWholeGitGraph,
	}
}

func (self *LocalCommitsController) Keybindings(
	getKey func(key string) interface{},
	config config.KeybindingConfig,
	guards types.KeybindingGuards,
) []*types.Binding {
	outsideFilterModeBindings := []*types.Binding{
		{
			Key:         getKey(config.Commits.SquashDown),
			Handler:     self.squashDown,
			Description: self.c.Tr.LcSquashDown,
		},
		{
			Key:         getKey(config.Commits.MarkCommitAsFixup),
			Handler:     self.fixup,
			Description: self.c.Tr.LcFixupCommit,
		},
		{
			Key:         getKey(config.Commits.RenameCommit),
			Handler:     self.checkSelected(self.reword),
			Description: self.c.Tr.LcRewordCommit,
		},
		{
			Key:         getKey(config.Commits.RenameCommitWithEditor),
			Handler:     self.rewordEditor,
			Description: self.c.Tr.LcRenameCommitEditor,
		},
		{
			Key:         getKey(config.Universal.Remove),
			Handler:     self.drop,
			Description: self.c.Tr.LcDeleteCommit,
		},
		{
			Key:         getKey(config.Universal.Edit),
			Handler:     self.edit,
			Description: self.c.Tr.LcEditCommit,
		},
		{
			Key:         getKey(config.Commits.PickCommit),
			Handler:     self.pick,
			Description: self.c.Tr.LcPickCommit,
		},
		{
			Key:         getKey(config.Commits.CreateFixupCommit),
			Handler:     self.checkSelected(self.handleCreateFixupCommit),
			Description: self.c.Tr.LcCreateFixupCommit,
		},
		{
			Key:         getKey(config.Commits.SquashAboveCommits),
			Handler:     self.checkSelected(self.handleSquashAllAboveFixupCommits),
			Description: self.c.Tr.LcSquashAboveCommits,
		},
		{
			Key:         getKey(config.Commits.MoveDownCommit),
			Handler:     self.handleCommitMoveDown,
			Description: self.c.Tr.LcMoveDownCommit,
		},
		{
			Key:         getKey(config.Commits.MoveUpCommit),
			Handler:     self.handleCommitMoveUp,
			Description: self.c.Tr.LcMoveUpCommit,
		},
		{
			Key:         getKey(config.Commits.AmendToCommit),
			Handler:     self.handleCommitAmendTo,
			Description: self.c.Tr.LcAmendToCommit,
		},
		{
			Key:         getKey(config.Commits.RevertCommit),
			Handler:     self.checkSelected(self.handleCommitRevert),
			Description: self.c.Tr.LcRevertCommit,
		},
		{
			Key:         getKey(config.Universal.New),
			Modifier:    gocui.ModNone,
			Handler:     self.checkSelected(self.newBranch),
			Description: self.c.Tr.LcCreateNewBranchFromCommit,
		},
		{
			Key:         getKey(config.Commits.CherryPickCopy),
			Handler:     self.checkSelected(self.copy),
			Description: self.c.Tr.LcCherryPickCopy,
		},
		{
			Key:         getKey(config.Commits.CherryPickCopyRange),
			Handler:     self.checkSelected(self.copyRange),
			Description: self.c.Tr.LcCherryPickCopyRange,
		},
		{
			Key:         getKey(config.Commits.PasteCommits),
			Handler:     guards.OutsideFilterMode(self.paste),
			Description: self.c.Tr.LcPasteCommits,
		},
		// overriding these navigation keybindings because we might need to load
		// more commits on demand
		{
			Key:         getKey(config.Universal.StartSearch),
			Handler:     self.openSearch,
			Description: self.c.Tr.LcStartSearch,
			Tag:         "navigation",
		},
		{
			Key:         getKey(config.Universal.GotoBottom),
			Handler:     self.gotoBottom,
			Description: self.c.Tr.LcGotoBottom,
			Tag:         "navigation",
		},
		{
			Key:     gocui.MouseLeft,
			Handler: func() error { return self.getContext().HandleClick(self.checkSelected(self.enter)) },
		},
	}

	for _, binding := range outsideFilterModeBindings {
		binding.Handler = guards.OutsideFilterMode(binding.Handler)
	}

	bindings := append(outsideFilterModeBindings, []*types.Binding{
		{
			Key:         getKey(config.Commits.OpenLogMenu),
			Handler:     self.handleOpenLogMenu,
			Description: self.c.Tr.LcOpenLogMenu,
			OpensMenu:   true,
		},
		{
			Key:         getKey(config.Commits.ViewResetOptions),
			Handler:     self.checkSelected(self.handleCreateCommitResetMenu),
			Description: self.c.Tr.LcResetToThisCommit,
		},
		{
			Key:         getKey(config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.LcViewCommitFiles,
		},
		{
			Key:         getKey(config.Commits.CheckoutCommit),
			Handler:     self.checkSelected(self.handleCheckoutCommit),
			Description: self.c.Tr.LcCheckoutCommit,
		},
		{
			Key:         getKey(config.Commits.TagCommit),
			Handler:     self.checkSelected(self.handleTagCommit),
			Description: self.c.Tr.LcTagCommit,
		},
		{
			Key:         getKey(config.Commits.CopyCommitMessageToClipboard),
			Handler:     self.checkSelected(self.handleCopySelectedCommitMessageToClipboard),
			Description: self.c.Tr.LcCopyCommitMessageToClipboard,
		},
		{
			Key:         getKey(config.Commits.OpenInBrowser),
			Handler:     self.checkSelected(self.handleOpenCommitInBrowser),
			Description: self.c.Tr.LcOpenCommitInBrowser,
		},
	}...)

	return append(bindings, self.getContext().Keybindings(getKey, config, guards)...)
}

func (self *LocalCommitsController) squashDown() error {
	if len(self.getCommits()) <= 1 {
		return self.c.ErrorMsg(self.c.Tr.YouNoCommitsToSquash)
	}

	applied, err := self.handleMidRebaseCommand("squash")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.Ask(types.AskOpts{
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

func (self *LocalCommitsController) fixup() error {
	if len(self.getCommits()) <= 1 {
		return self.c.ErrorMsg(self.c.Tr.YouNoCommitsToSquash)
	}

	applied, err := self.handleMidRebaseCommand("fixup")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.Ask(types.AskOpts{
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
	applied, err := self.handleMidRebaseCommand("reword")
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
			if err := self.git.Rebase.RewordCommit(self.getCommits(), self.getSelectedLocalCommitIdx(), response); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

func (self *LocalCommitsController) rewordEditor() error {
	applied, err := self.handleMidRebaseCommand("reword")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	self.c.LogAction(self.c.Tr.Actions.RewordCommit)
	subProcess, err := self.git.Rebase.RewordCommitInEditor(
		self.getCommits(), self.getSelectedLocalCommitIdx(),
	)
	if err != nil {
		return self.c.Error(err)
	}
	if subProcess != nil {
		return self.c.RunSubprocessAndRefresh(subProcess)
	}

	return nil
}

func (self *LocalCommitsController) drop() error {
	applied, err := self.handleMidRebaseCommand("drop")
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	return self.c.Ask(types.AskOpts{
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

func (self *LocalCommitsController) edit() error {
	applied, err := self.handleMidRebaseCommand("edit")
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

func (self *LocalCommitsController) pick() error {
	applied, err := self.handleMidRebaseCommand("pick")
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
	err := self.git.Rebase.InteractiveRebase(self.getCommits(), self.getSelectedLocalCommitIdx(), action)
	return self.CheckMergeOrRebase(err)
}

// handleMidRebaseCommand sees if the selected commit is in fact a rebasing
// commit meaning you are trying to edit the todo file rather than actually
// begin a rebase. It then updates the todo file with that action
func (self *LocalCommitsController) handleMidRebaseCommand(action string) (bool, error) {
	selectedCommit := self.getSelectedLocalCommit()
	if selectedCommit.Status != "rebasing" {
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
		fmt.Sprintf("Updating rebase action of commit %s to '%s'", selectedCommit.ShortSha(), action),
		false,
	)

	if err := self.git.Rebase.EditRebaseTodo(
		self.getSelectedLocalCommitIdx(), action,
	); err != nil {
		return false, self.c.Error(err)
	}

	return true, self.c.Refresh(types.RefreshOptions{
		Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
	})
}

func (self *LocalCommitsController) handleCommitMoveDown() error {
	index := self.getContext().GetPanelState().GetSelectedLineIdx()
	commits := self.getCommits()
	selectedCommit := self.getCommits()[index]
	if selectedCommit.Status == "rebasing" {
		if commits[index+1].Status != "rebasing" {
			return nil
		}

		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)
		self.c.LogCommand(fmt.Sprintf("Moving commit %s down", selectedCommit.ShortSha()), false)

		if err := self.git.Rebase.MoveTodoDown(index); err != nil {
			return self.c.Error(err)
		}
		// TODO: use MoveSelectedLine
		_ = self.getContext().HandleNextLine()
		return self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
	}

	return self.c.WithWaitingStatus(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitDown)
		err := self.git.Rebase.MoveCommitDown(self.getCommits(), index)
		if err == nil {
			// TODO: use MoveSelectedLine
			_ = self.getContext().HandleNextLine()
		}
		return self.CheckMergeOrRebase(err)
	})
}

func (self *LocalCommitsController) handleCommitMoveUp() error {
	index := self.getContext().GetPanelState().GetSelectedLineIdx()
	if index == 0 {
		return nil
	}

	selectedCommit := self.getCommits()[index]
	if selectedCommit.Status == "rebasing" {
		// logging directly here because MoveTodoDown doesn't have enough information
		// to provide a useful log
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		self.c.LogCommand(
			fmt.Sprintf("Moving commit %s up", selectedCommit.ShortSha()),
			false,
		)

		if err := self.git.Rebase.MoveTodoDown(index - 1); err != nil {
			return self.c.Error(err)
		}
		_ = self.getContext().HandlePrevLine()
		return self.c.Refresh(types.RefreshOptions{
			Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
		})
	}

	return self.c.WithWaitingStatus(self.c.Tr.MovingStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.MoveCommitUp)
		err := self.git.Rebase.MoveCommitDown(self.getCommits(), index-1)
		if err == nil {
			_ = self.getContext().HandlePrevLine()
		}
		return self.CheckMergeOrRebase(err)
	})
}

func (self *LocalCommitsController) handleCommitAmendTo() error {
	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.AmendCommitTitle,
		Prompt: self.c.Tr.AmendCommitPrompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.AmendingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.AmendCommit)
				err := self.git.Rebase.AmendTo(self.getSelectedLocalCommit().Sha)
				return self.CheckMergeOrRebase(err)
			})
		},
	})
}

func (self *LocalCommitsController) handleCommitRevert(commit *models.Commit) error {
	if commit.IsMerge() {
		return self.createRevertMergeCommitMenu(commit)
	} else {
		return self.c.Ask(types.AskOpts{
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
			DisplayString: fmt.Sprintf("%s: %s", utils.SafeTruncate(parentSha, 8), message),
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
	_ = self.getContext().HandleNextLine()
	return self.c.Refresh(types.RefreshOptions{
		Mode: types.BLOCK_UI, Scope: []types.RefreshableView{types.COMMITS, types.BRANCHES},
	})
}

func (self *LocalCommitsController) enter(commit *models.Commit) error {
	return self.switchToCommitFilesContext(SwitchToCommitFilesContextOpts{
		RefName:    commit.Sha,
		CanRebase:  true,
		Context:    self.getContext(),
		WindowName: "commits",
	})
}

func (self *LocalCommitsController) handleCreateFixupCommit(commit *models.Commit) error {
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.SureCreateFixupCommit,
		map[string]string{
			"commit": commit.Sha,
		},
	)

	return self.c.Ask(types.AskOpts{
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

func (self *LocalCommitsController) handleSquashAllAboveFixupCommits(commit *models.Commit) error {
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.SureSquashAboveCommits,
		map[string]string{
			"commit": commit.Sha,
		},
	)

	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.SquashAboveCommits,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.SquashingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.SquashAllAboveFixupCommits)
				err := self.git.Rebase.SquashAllAboveFixupCommits(commit.Sha)
				return self.CheckMergeOrRebase(err)
			})
		},
	})
}

func (self *LocalCommitsController) handleTagCommit(commit *models.Commit) error {
	return self.tagsHelper.CreateTagMenu(commit.Sha, func() {})
}

func (self *LocalCommitsController) handleCheckoutCommit(commit *models.Commit) error {
	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.LcCheckoutCommit,
		Prompt: self.c.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.CheckoutCommit)
			return self.refsHelper.CheckoutRef(commit.Sha, types.CheckoutRefOptions{})
		},
	})
}

func (self *LocalCommitsController) handleCreateCommitResetMenu(commit *models.Commit) error {
	return self.refsHelper.CreateGitResetMenu(commit.Sha)
}

func (self *LocalCommitsController) openSearch() error {
	// we usually lazyload these commits but now that we're searching we need to load them now
	if self.getLimitCommits() {
		self.setLimitCommits(false)
		if err := self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
			return err
		}
	}

	self.c.OpenSearch()

	return nil
}

func (self *LocalCommitsController) gotoBottom() error {
	// we usually lazyload these commits but now that we're jumping to the bottom we need to load them now
	if self.getLimitCommits() {
		self.setLimitCommits(false)
		if err := self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}}); err != nil {
			return err
		}
	}

	_ = self.getContext().HandleGotoBottom()

	return nil
}

func (self *LocalCommitsController) handleCopySelectedCommitMessageToClipboard(commit *models.Commit) error {
	message, err := self.git.Commit.GetCommitMessage(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.CopyCommitMessageToClipboard)
	if err := self.os.CopyToClipboard(message); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.CommitMessageCopiedToClipboard)

	return nil
}

func (self *LocalCommitsController) handleOpenLogMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.LogMenuTitle,
		Items: []*types.MenuItem{
			{
				DisplayString: self.c.Tr.ToggleShowGitGraphAll,
				OnPress: func() error {
					self.setShowWholeGitGraph(!self.getShowWholeGitGraph())

					if self.getShowWholeGitGraph() {
						self.setLimitCommits(false)
					}

					return self.c.WithWaitingStatus(self.c.Tr.LcLoadingCommits, func() error {
						return self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}})
					})
				},
			},
			{
				DisplayString: self.c.Tr.ShowGitGraph,
				OpensMenu:     true,
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
								DisplayString: "always",
								OnPress:       onPress("always"),
							},
							{
								DisplayString: "never",
								OnPress:       onPress("never"),
							},
							{
								DisplayString: "when maximised",
								OnPress:       onPress("when-maximised"),
							},
						},
					})
				},
			},
			{
				DisplayString: self.c.Tr.SortCommits,
				OpensMenu:     true,
				OnPress: func() error {
					onPress := func(value string) func() error {
						return func() error {
							self.c.UserConfig.Git.Log.Order = value
							return self.c.WithWaitingStatus(self.c.Tr.LcLoadingCommits, func() error {
								return self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.COMMITS}})
							})
						}
					}

					return self.c.Menu(types.CreateMenuOptions{
						Title: self.c.Tr.LogMenuTitle,
						Items: []*types.MenuItem{
							{
								DisplayString: "topological (topo-order)",
								OnPress:       onPress("topo-order"),
							},
							{
								DisplayString: "date-order",
								OnPress:       onPress("date-order"),
							},
							{
								DisplayString: "author-date-order",
								OnPress:       onPress("author-date-order"),
							},
						},
					})
				},
			},
		},
	})
}

func (self *LocalCommitsController) handleOpenCommitInBrowser(commit *models.Commit) error {
	hostingServiceMgr := self.getHostingServiceMgr()

	url, err := hostingServiceMgr.GetCommitURL(commit.Sha)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.OpenCommitInBrowser)
	if err := self.os.OpenLink(url); err != nil {
		return self.c.Error(err)
	}

	return nil
}

func (self *LocalCommitsController) checkSelected(callback func(*models.Commit) error) func() error {
	return func() error {
		commit := self.getSelectedLocalCommit()
		if commit == nil {
			return nil
		}

		return callback(commit)
	}
}

func (self *LocalCommitsController) Context() types.Context {
	return self.getContext()
}

func (self *LocalCommitsController) newBranch(commit *models.Commit) error {
	return self.refsHelper.NewBranch(commit.RefName(), commit.Description(), "")
}

func (self *LocalCommitsController) copy(commit *models.Commit) error {
	return self.cherryPickHelper.Copy(commit, self.getCommits(), self.getContext())
}

func (self *LocalCommitsController) copyRange(*models.Commit) error {
	return self.cherryPickHelper.CopyRange(self.getContext().GetPanelState().GetSelectedLineIdx(), self.getCommits(), self.getContext())
}

func (self *LocalCommitsController) paste() error {
	return self.cherryPickHelper.Paste()
}
