package controllers

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type BranchesController struct {
	baseController
	*ListControllerTrait[*models.Branch]
	c *ControllerCommon
}

var _ types.IController = &BranchesController{}

func NewBranchesController(
	c *ControllerCommon,
) *BranchesController {
	return &BranchesController{
		baseController: baseController{},
		c:              c,
		ListControllerTrait: NewListControllerTrait[*models.Branch](
			c,
			c.Contexts().Branches,
			c.Contexts().Branches.GetSelected,
			c.Contexts().Branches.GetSelectedItems,
		),
	}
}

func (self *BranchesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.Select),
			Handler: self.withItem(self.press),
			GetDisabledReason: self.require(
				self.singleItemSelected(),
				self.notPulling,
			),
			Description:     self.c.Tr.Checkout,
			Tooltip:         self.c.Tr.CheckoutTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.New),
			Handler:           self.withItem(self.newBranch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewBranch,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.CreatePullRequest),
			Handler:           self.withItem(self.handleCreatePullRequest),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CreatePullRequest,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.ViewPullRequestOptions),
			Handler:           self.withItem(self.handleCreatePullRequestMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CreatePullRequestOptions,
			OpensMenu:         true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.CopyPullRequestURL),
			Handler:           self.copyPullRequestURL,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.CopyPullRequestURL,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.CheckoutBranchByName),
			Handler:     self.checkoutByName,
			Description: self.c.Tr.CheckoutByName,
			Tooltip:     self.c.Tr.CheckoutByNameTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.ForceCheckoutBranch),
			Handler:           self.forceCheckout,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ForceCheckout,
			Tooltip:           self.c.Tr.ForceCheckoutTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.delete),
			GetDisabledReason: self.require(self.itemRangeSelected(self.branchesAreReal)),
			Description:       self.c.Tr.Delete,
			Tooltip:           self.c.Tr.BranchDeleteTooltip,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranch),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.rebase)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.RebaseBranch,
			Tooltip:           self.c.Tr.RebaseBranchTooltip,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeIntoCurrentBranch),
			Handler:           opts.Guards.OutsideFilterMode(self.merge),
			GetDisabledReason: self.require(self.singleItemSelected(self.notMergingIntoYourself)),
			Description:       self.c.Tr.Merge,
			Tooltip:           self.c.Tr.MergeBranchTooltip,
			DisplayOnScreen:   true,
			OpensMenu:         true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.FastForward),
			Handler:           self.withItem(self.fastForward),
			GetDisabledReason: self.require(self.singleItemSelected(self.branchIsReal)),
			Description:       self.c.Tr.FastForward,
			Tooltip:           self.c.Tr.FastForwardTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.CreateTag),
			Handler:           self.withItem(self.createTag),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewTag,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.SortOrder),
			Handler:     self.createSortMenu,
			Description: self.c.Tr.SortOrder,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:           self.withItem(self.createResetMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ViewResetOptions,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RenameBranch),
			Handler:           self.withItem(self.rename),
			GetDisabledReason: self.require(self.singleItemSelected(self.branchIsReal)),
			Description:       self.c.Tr.RenameBranch,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.SetUpstream),
			Handler:           self.withItem(self.viewUpstreamOptions),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ViewBranchUpstreamOptions,
			Tooltip:           self.c.Tr.ViewBranchUpstreamOptionsTooltip,
			ShortDescription:  self.c.Tr.Upstream,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Key: opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler: self.withItem(func(selectedBranch *models.Branch) error {
				return self.c.Helpers().Diff.OpenDiffToolForRef(selectedBranch)
			}),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
	}
}

func (self *BranchesController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			branch := self.context().GetSelected()
			if branch == nil {
				task = types.NewRenderStringTask(self.c.Tr.NoBranchesThisRepo)
			} else {
				cmdObj := self.c.Git().Branch.GetGraphCmdObj(branch.FullRefName())

				task = types.NewRunPtyTask(cmdObj.GetCmd())
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: self.c.Tr.LogTitle,
					Task:  task,
				},
			})
		})
	}
}

func (self *BranchesController) viewUpstreamOptions(selectedBranch *models.Branch) error {
	upstream := lo.Ternary(selectedBranch.RemoteBranchStoredLocally(),
		selectedBranch.ShortUpstreamRefName(),
		self.c.Tr.UpstreamGenericName)

	viewDivergenceItem := &types.MenuItem{
		LabelColumns: []string{self.c.Tr.ViewDivergenceFromUpstream},
		OnPress: func() error {
			branch := self.context().GetSelected()
			if branch == nil {
				return nil
			}

			return self.c.Helpers().SubCommits.ViewSubCommits(helpers.ViewSubCommitsOpts{
				Ref:                     branch,
				TitleRef:                fmt.Sprintf("%s <-> %s", branch.RefName(), upstream),
				RefToShowDivergenceFrom: branch.FullUpstreamRefName(),
				Context:                 self.context(),
				ShowBranchHeads:         false,
			})
		},
	}

	var disabledReason *types.DisabledReason
	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(selectedBranch, self.c.Model().MainBranches)
	if err != nil {
		return err
	}
	if baseBranch == "" {
		baseBranch = self.c.Tr.CouldNotDetermineBaseBranch
		disabledReason = &types.DisabledReason{Text: self.c.Tr.CouldNotDetermineBaseBranch}
	}
	shortBaseBranchName := helpers.ShortBranchName(baseBranch)
	label := utils.ResolvePlaceholderString(
		self.c.Tr.ViewDivergenceFromBaseBranch,
		map[string]string{"baseBranch": shortBaseBranchName},
	)
	viewDivergenceFromBaseBranchItem := &types.MenuItem{
		LabelColumns: []string{label},
		Key:          'b',
		OnPress: func() error {
			branch := self.context().GetSelected()
			if branch == nil {
				return nil
			}

			return self.c.Helpers().SubCommits.ViewSubCommits(helpers.ViewSubCommitsOpts{
				Ref:                     branch,
				TitleRef:                fmt.Sprintf("%s <-> %s", branch.RefName(), shortBaseBranchName),
				RefToShowDivergenceFrom: baseBranch,
				Context:                 self.context(),
				ShowBranchHeads:         false,
			})
		},
		DisabledReason: disabledReason,
	}

	unsetUpstreamItem := &types.MenuItem{
		LabelColumns: []string{self.c.Tr.UnsetUpstream},
		OnPress: func() error {
			if err := self.c.Git().Branch.UnsetUpstream(selectedBranch.Name); err != nil {
				return err
			}
			if err := self.c.Refresh(types.RefreshOptions{
				Mode: types.SYNC,
				Scope: []types.RefreshableView{
					types.BRANCHES,
					types.COMMITS,
				},
			}); err != nil {
				return err
			}
			return nil
		},
		Key: 'u',
	}

	setUpstreamItem := &types.MenuItem{
		LabelColumns: []string{self.c.Tr.SetUpstream},
		OnPress: func() error {
			return self.c.Helpers().Upstream.PromptForUpstreamWithoutInitialContent(selectedBranch, func(upstream string) error {
				upstreamRemote, upstreamBranch, err := self.c.Helpers().Upstream.ParseUpstream(upstream)
				if err != nil {
					return err
				}

				if err := self.c.Git().Branch.SetUpstream(upstreamRemote, upstreamBranch, selectedBranch.Name); err != nil {
					return err
				}
				if err := self.c.Refresh(types.RefreshOptions{
					Mode: types.SYNC,
					Scope: []types.RefreshableView{
						types.BRANCHES,
						types.COMMITS,
					},
				}); err != nil {
					return err
				}
				return nil
			})
		},
		Key: 's',
	}

	upstreamResetOptions := utils.ResolvePlaceholderString(
		self.c.Tr.ViewUpstreamResetOptions,
		map[string]string{"upstream": upstream},
	)
	upstreamResetTooltip := utils.ResolvePlaceholderString(
		self.c.Tr.ViewUpstreamResetOptionsTooltip,
		map[string]string{"upstream": upstream},
	)

	upstreamRebaseOptions := utils.ResolvePlaceholderString(
		self.c.Tr.ViewUpstreamRebaseOptions,
		map[string]string{"upstream": upstream},
	)
	upstreamRebaseTooltip := utils.ResolvePlaceholderString(
		self.c.Tr.ViewUpstreamRebaseOptionsTooltip,
		map[string]string{"upstream": upstream},
	)

	upstreamResetItem := &types.MenuItem{
		LabelColumns: []string{upstreamResetOptions},
		OpensMenu:    true,
		OnPress: func() error {
			err := self.c.Helpers().Refs.CreateGitResetMenu(upstream)
			if err != nil {
				return err
			}
			return nil
		},
		Tooltip: upstreamResetTooltip,
		Key:     'g',
	}

	upstreamRebaseItem := &types.MenuItem{
		LabelColumns: []string{upstreamRebaseOptions},
		OpensMenu:    true,
		OnPress: func() error {
			if err := self.c.Helpers().MergeAndRebase.RebaseOntoRef(upstream); err != nil {
				return err
			}
			return nil
		},
		Tooltip: upstreamRebaseTooltip,
		Key:     'r',
	}

	if !selectedBranch.IsTrackingRemote() {
		unsetUpstreamItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.UpstreamNotSetError}
	}

	if !selectedBranch.RemoteBranchStoredLocally() {
		viewDivergenceItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.UpstreamNotSetError}
		upstreamResetItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.UpstreamNotSetError}
		upstreamRebaseItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.UpstreamNotSetError}
	}

	options := []*types.MenuItem{
		viewDivergenceItem,
		viewDivergenceFromBaseBranchItem,
		unsetUpstreamItem,
		setUpstreamItem,
		upstreamResetItem,
		upstreamRebaseItem,
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.BranchUpstreamOptionsTitle,
		Items: options,
	})
}

func (self *BranchesController) Context() types.Context {
	return self.context()
}

func (self *BranchesController) context() *context.BranchesContext {
	return self.c.Contexts().Branches
}

func (self *BranchesController) press(selectedBranch *models.Branch) error {
	if selectedBranch == self.c.Helpers().Refs.GetCheckedOutRef() {
		return errors.New(self.c.Tr.AlreadyCheckedOutBranch)
	}

	worktreeForRef, ok := self.worktreeForBranch(selectedBranch)
	if ok && !worktreeForRef.IsCurrent {
		return self.promptToCheckoutWorktree(worktreeForRef)
	}

	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	return self.c.Helpers().Refs.CheckoutRef(selectedBranch.Name, types.CheckoutRefOptions{})
}

func (self *BranchesController) notPulling() *types.DisabledReason {
	currentBranch := self.c.Helpers().Refs.GetCheckedOutRef()
	if currentBranch != nil {
		op := self.c.State().GetItemOperation(currentBranch)
		if op == types.ItemOperationFastForwarding || op == types.ItemOperationPulling {
			return &types.DisabledReason{Text: self.c.Tr.CantCheckoutBranchWhilePulling}
		}
	}

	return nil
}

func (self *BranchesController) worktreeForBranch(branch *models.Branch) (*models.Worktree, bool) {
	return git_commands.WorktreeForBranch(branch, self.c.Model().Worktrees)
}

func (self *BranchesController) promptToCheckoutWorktree(worktree *models.Worktree) error {
	prompt := utils.ResolvePlaceholderString(self.c.Tr.AlreadyCheckedOutByWorktree, map[string]string{
		"worktreeName": worktree.Name,
	})

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.SwitchToWorktree,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.Helpers().Worktree.Switch(worktree, context.LOCAL_BRANCHES_CONTEXT_KEY)
		},
	})

	return nil
}

func (self *BranchesController) handleCreatePullRequest(selectedBranch *models.Branch) error {
	if !selectedBranch.IsTrackingRemote() {
		return errors.New(self.c.Tr.PullRequestNoUpstream)
	}
	return self.createPullRequest(selectedBranch.UpstreamBranch, "")
}

func (self *BranchesController) handleCreatePullRequestMenu(selectedBranch *models.Branch) error {
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef()

	return self.createPullRequestMenu(selectedBranch, checkedOutBranch)
}

func (self *BranchesController) copyPullRequestURL() error {
	branch := self.context().GetSelected()

	branchExistsOnRemote := self.c.Git().Remote.CheckRemoteBranchExists(branch.Name)

	if !branchExistsOnRemote {
		return errors.New(self.c.Tr.NoBranchOnRemote)
	}

	url, err := self.c.Helpers().Host.GetPullRequestURL(branch.Name, "")
	if err != nil {
		return err
	}
	self.c.LogAction(self.c.Tr.Actions.CopyPullRequestURL)
	if err := self.c.OS().CopyToClipboard(url); err != nil {
		return err
	}

	self.c.Toast(self.c.Tr.PullRequestURLCopiedToClipboard)

	return nil
}

func (self *BranchesController) forceCheckout() error {
	branch := self.context().GetSelected()
	message := self.c.Tr.SureForceCheckout
	title := self.c.Tr.ForceCheckoutBranch

	self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ForceCheckoutBranch)
			if err := self.c.Git().Branch.Checkout(branch.Name, git_commands.CheckoutOptions{Force: true}); err != nil {
				return err
			}
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})

	return nil
}

func (self *BranchesController) checkoutByName() error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.BranchName + ":",
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRefsSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			self.c.LogAction("Checkout branch")
			_, branchName, found := self.c.Helpers().Refs.ParseRemoteBranchName(response)
			if found {
				return self.c.Helpers().Refs.CheckoutRemoteBranch(response, branchName)
			}
			return self.c.Helpers().Refs.CheckoutRef(response, types.CheckoutRefOptions{
				OnRefNotFound: func(ref string) error {
					self.c.Confirm(types.ConfirmOpts{
						Title:  self.c.Tr.BranchNotFoundTitle,
						Prompt: fmt.Sprintf("%s %s%s", self.c.Tr.BranchNotFoundPrompt, ref, "?"),
						HandleConfirm: func() error {
							return self.createNewBranchWithName(ref)
						},
					})

					return nil
				},
			})
		},
	},
	)

	return nil
}

func (self *BranchesController) createNewBranchWithName(newBranchName string) error {
	branch := self.context().GetSelected()
	if branch == nil {
		return nil
	}

	if err := self.c.Git().Branch.New(newBranchName, branch.FullRefName()); err != nil {
		return err
	}

	self.context().SetSelection(0)
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, KeepBranchSelectionIndex: true})
}

func (self *BranchesController) localDelete(branches []*models.Branch) error {
	return self.c.Helpers().BranchesHelper.ConfirmLocalDelete(branches)
}

func (self *BranchesController) remoteDelete(branches []*models.Branch) error {
	remoteBranches := lo.Map(branches, func(branch *models.Branch, _ int) *models.RemoteBranch {
		return &models.RemoteBranch{Name: branch.UpstreamBranch, RemoteName: branch.UpstreamRemote}
	})
	return self.c.Helpers().BranchesHelper.ConfirmDeleteRemote(remoteBranches)
}

func (self *BranchesController) localAndRemoteDelete(branches []*models.Branch) error {
	return self.c.Helpers().BranchesHelper.ConfirmLocalAndRemoteDelete(branches)
}

func (self *BranchesController) delete(branches []*models.Branch) error {
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef()
	isBranchCheckedOut := lo.SomeBy(branches, func(branch *models.Branch) bool {
		return checkedOutBranch.Name == branch.Name
	})
	hasUpstream := lo.EveryBy(branches, func(branch *models.Branch) bool {
		return branch.IsTrackingRemote() && !branch.UpstreamGone
	})

	localDeleteItem := &types.MenuItem{
		Label: lo.Ternary(len(branches) > 1, self.c.Tr.DeleteLocalBranches, self.c.Tr.DeleteLocalBranch),
		Key:   'c',
		OnPress: func() error {
			return self.localDelete(branches)
		},
	}
	if isBranchCheckedOut {
		localDeleteItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.CantDeleteCheckOutBranch}
	}

	remoteDeleteItem := &types.MenuItem{
		Label: lo.Ternary(len(branches) > 1, self.c.Tr.DeleteRemoteBranches, self.c.Tr.DeleteRemoteBranch),
		Key:   'r',
		OnPress: func() error {
			return self.remoteDelete(branches)
		},
	}
	if !hasUpstream {
		remoteDeleteItem.DisabledReason = &types.DisabledReason{
			Text: lo.Ternary(len(branches) > 1, self.c.Tr.UpstreamsNotSetError, self.c.Tr.UpstreamNotSetError),
		}
	}

	deleteBothItem := &types.MenuItem{
		Label: lo.Ternary(len(branches) > 1, self.c.Tr.DeleteLocalAndRemoteBranches, self.c.Tr.DeleteLocalAndRemoteBranch),
		Key:   'b',
		OnPress: func() error {
			return self.localAndRemoteDelete(branches)
		},
	}
	if isBranchCheckedOut {
		deleteBothItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.CantDeleteCheckOutBranch}
	} else if !hasUpstream {
		deleteBothItem.DisabledReason = &types.DisabledReason{
			Text: lo.Ternary(len(branches) > 1, self.c.Tr.UpstreamsNotSetError, self.c.Tr.UpstreamNotSetError),
		}
	}

	var menuTitle string
	if len(branches) == 1 {
		menuTitle = utils.ResolvePlaceholderString(
			self.c.Tr.DeleteBranchTitle,
			map[string]string{
				"selectedBranchName": branches[0].Name,
			},
		)
	} else {
		menuTitle = self.c.Tr.DeleteBranchesTitle
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: menuTitle,
		Items: []*types.MenuItem{localDeleteItem, remoteDeleteItem, deleteBothItem},
	})
}

func (self *BranchesController) merge() error {
	selectedBranchName := self.context().GetSelected().Name
	return self.c.Helpers().MergeAndRebase.MergeRefIntoCheckedOutBranch(selectedBranchName)
}

func (self *BranchesController) rebase(branch *models.Branch) error {
	return self.c.Helpers().MergeAndRebase.RebaseOntoRef(branch.Name)
}

func (self *BranchesController) fastForward(branch *models.Branch) error {
	if !branch.IsTrackingRemote() {
		return errors.New(self.c.Tr.FwdNoUpstream)
	}
	if !branch.RemoteBranchStoredLocally() {
		return errors.New(self.c.Tr.FwdNoLocalUpstream)
	}
	if branch.IsAheadForPull() {
		return errors.New(self.c.Tr.FwdCommitsToPush)
	}

	action := self.c.Tr.Actions.FastForwardBranch

	return self.c.WithInlineStatus(branch, types.ItemOperationFastForwarding, context.LOCAL_BRANCHES_CONTEXT_KEY, func(task gocui.Task) error {
		worktree, ok := self.worktreeForBranch(branch)
		if ok {
			self.c.LogAction(action)

			worktreeGitDir := ""
			worktreePath := ""
			// if it is the current worktree path, no need to specify the path
			if !worktree.IsCurrent {
				worktreeGitDir = worktree.GitDir
				worktreePath = worktree.Path
			}

			err := self.c.Git().Sync.Pull(
				task,
				git_commands.PullOptions{
					RemoteName:      branch.UpstreamRemote,
					BranchName:      branch.UpstreamBranch,
					FastForwardOnly: true,
					WorktreeGitDir:  worktreeGitDir,
					WorktreePath:    worktreePath,
				},
			)
			_ = self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
			return err
		} else {
			self.c.LogAction(action)

			err := self.c.Git().Sync.FastForward(
				task, branch.Name, branch.UpstreamRemote, branch.UpstreamBranch,
			)
			_ = self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
			return err
		}
	})
}

func (self *BranchesController) createTag(branch *models.Branch) error {
	return self.c.Helpers().Tags.OpenCreateTagPrompt(branch.FullRefName(), func() {})
}

func (self *BranchesController) createSortMenu() error {
	return self.c.Helpers().Refs.CreateSortOrderMenu([]string{"recency", "alphabetical", "date"}, func(sortOrder string) error {
		if self.c.GetAppState().LocalBranchSortOrder != sortOrder {
			self.c.GetAppState().LocalBranchSortOrder = sortOrder
			self.c.SaveAppStateAndLogError()
			self.c.Contexts().Branches.SetSelection(0)
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		}
		return nil
	},
		self.c.GetAppState().LocalBranchSortOrder)
}

func (self *BranchesController) createResetMenu(selectedBranch *models.Branch) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(selectedBranch.Name)
}

func (self *BranchesController) rename(branch *models.Branch) error {
	promptForNewName := func() error {
		self.c.Prompt(types.PromptOpts{
			Title:          self.c.Tr.NewBranchNamePrompt + " " + branch.Name + ":",
			InitialContent: branch.Name,
			HandleConfirm: func(newBranchName string) error {
				self.c.LogAction(self.c.Tr.Actions.RenameBranch)
				if err := self.c.Git().Branch.Rename(branch.Name, helpers.SanitizedBranchName(newBranchName)); err != nil {
					return err
				}

				// need to find where the branch is now so that we can re-select it. That means we need to refetch the branches synchronously and then find our branch
				_ = self.c.Refresh(types.RefreshOptions{
					Mode:  types.SYNC,
					Scope: []types.RefreshableView{types.BRANCHES, types.WORKTREES},
				})

				// now that we've got our stuff again we need to find that branch and reselect it.
				for i, newBranch := range self.c.Model().Branches {
					if newBranch.Name == newBranchName {
						self.context().SetSelection(i)
						self.context().HandleRender()
					}
				}

				return nil
			},
		})

		return nil
	}

	// I could do an explicit check here for whether the branch is tracking a remote branch
	// but if we've selected it we'll already know that via Pullables and Pullables.
	// Bit of a hack but I'm lazy.
	if !branch.IsTrackingRemote() {
		return promptForNewName()
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:         self.c.Tr.RenameBranch,
		Prompt:        self.c.Tr.RenameBranchWarning,
		HandleConfirm: promptForNewName,
	})

	return nil
}

func (self *BranchesController) newBranch(selectedBranch *models.Branch) error {
	return self.c.Helpers().Refs.NewBranch(selectedBranch.FullRefName(), selectedBranch.RefName(), "")
}

func (self *BranchesController) createPullRequestMenu(selectedBranch *models.Branch, checkedOutBranch *models.Branch) error {
	menuItems := make([]*types.MenuItem, 0, 4)

	fromToLabelColumns := func(from string, to string) []string {
		return []string{fmt.Sprintf("%s → %s", from, to)}
	}

	menuItemsForBranch := func(branch *models.Branch) []*types.MenuItem {
		return []*types.MenuItem{
			{
				LabelColumns: fromToLabelColumns(branch.Name, self.c.Tr.DefaultBranch),
				OnPress: func() error {
					return self.handleCreatePullRequest(branch)
				},
			},
			{
				LabelColumns: fromToLabelColumns(branch.Name, self.c.Tr.SelectBranch),
				OnPress: func() error {
					if !branch.IsTrackingRemote() {
						return errors.New(self.c.Tr.PullRequestNoUpstream)
					}

					if len(self.c.Model().Remotes) == 1 {
						toRemote := self.c.Model().Remotes[0].Name
						self.c.Log.Debugf("PR will target the only existing remote '%s'", toRemote)
						return self.promptForTargetBranchNameAndCreatePullRequest(branch, toRemote)
					}

					self.c.Prompt(types.PromptOpts{
						Title:               self.c.Tr.SelectTargetRemote,
						FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
						HandleConfirm: func(toRemote string) error {
							self.c.Log.Debugf("PR will target remote '%s'", toRemote)

							return self.promptForTargetBranchNameAndCreatePullRequest(branch, toRemote)
						},
					})

					return nil
				},
			},
		}
	}

	if selectedBranch != checkedOutBranch {
		menuItems = append(menuItems,
			&types.MenuItem{
				LabelColumns: fromToLabelColumns(checkedOutBranch.Name, selectedBranch.Name),
				OnPress: func() error {
					if !checkedOutBranch.IsTrackingRemote() || !selectedBranch.IsTrackingRemote() {
						return errors.New(self.c.Tr.PullRequestNoUpstream)
					}
					return self.createPullRequest(checkedOutBranch.UpstreamBranch, selectedBranch.UpstreamBranch)
				},
			},
		)
		menuItems = append(menuItems, menuItemsForBranch(checkedOutBranch)...)
	}

	menuItems = append(menuItems, menuItemsForBranch(selectedBranch)...)

	return self.c.Menu(types.CreateMenuOptions{Title: fmt.Sprint(self.c.Tr.CreatePullRequestOptions), Items: menuItems})
}

func (self *BranchesController) promptForTargetBranchNameAndCreatePullRequest(fromBranch *models.Branch, toRemote string) error {
	remoteDoesNotExist := lo.NoneBy(self.c.Model().Remotes, func(remote *models.Remote) bool {
		return remote.Name == toRemote
	})
	if remoteDoesNotExist {
		return fmt.Errorf(self.c.Tr.NoValidRemoteName, toRemote)
	}

	self.c.Prompt(types.PromptOpts{
		Title:               fmt.Sprintf("%s → %s/", fromBranch.UpstreamBranch, toRemote),
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteBranchesForRemoteSuggestionsFunc(toRemote),
		HandleConfirm: func(toBranch string) error {
			self.c.Log.Debugf("PR will target branch '%s' on remote '%s'", toBranch, toRemote)
			return self.createPullRequest(fromBranch.UpstreamBranch, toBranch)
		},
	})

	return nil
}

func (self *BranchesController) createPullRequest(from string, to string) error {
	url, err := self.c.Helpers().Host.GetPullRequestURL(from, to)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.OpenPullRequest)

	if err := self.c.OS().OpenLink(url); err != nil {
		return err
	}

	return nil
}

func (self *BranchesController) branchIsReal(branch *models.Branch) *types.DisabledReason {
	if !branch.IsRealBranch() {
		return &types.DisabledReason{Text: self.c.Tr.SelectedItemIsNotABranch}
	}

	return nil
}

func (self *BranchesController) branchesAreReal(selectedBranches []*models.Branch, startIdx int, endIdx int) *types.DisabledReason {
	if !lo.EveryBy(selectedBranches, func(branch *models.Branch) bool {
		return branch.IsRealBranch()
	}) {
		return &types.DisabledReason{Text: self.c.Tr.SelectedItemIsNotABranch}
	}

	return nil
}

func (self *BranchesController) notMergingIntoYourself(branch *models.Branch) *types.DisabledReason {
	selectedBranchName := branch.Name
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef().Name

	if checkedOutBranch == selectedBranchName {
		return &types.DisabledReason{Text: self.c.Tr.CantMergeBranchIntoItself}
	}

	return nil
}
