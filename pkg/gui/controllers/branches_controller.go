package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type BranchesController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &BranchesController{}

func NewBranchesController(
	common *controllerCommon,
) *BranchesController {
	return &BranchesController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *BranchesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelected(self.press),
			Description: self.c.Tr.LcCheckout,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.checkSelected(self.newBranch),
			Description: self.c.Tr.LcNewBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.CreatePullRequest),
			Handler:     self.checkSelected(self.handleCreatePullRequest),
			Description: self.c.Tr.LcCreatePullRequest,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.ViewPullRequestOptions),
			Handler:     self.checkSelected(self.handleCreatePullRequestMenu),
			Description: self.c.Tr.LcCreatePullRequestOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.CopyPullRequestURL),
			Handler:     self.copyPullRequestURL,
			Description: self.c.Tr.LcCopyPullRequestURL,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.CheckoutBranchByName),
			Handler:     self.checkoutByName,
			Description: self.c.Tr.LcCheckoutByName,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.ForceCheckoutBranch),
			Handler:     self.forceCheckout,
			Description: self.c.Tr.LcForceCheckout,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelectedAndReal(self.delete),
			Description: self.c.Tr.LcDeleteBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.RebaseBranch),
			Handler:     opts.Guards.OutsideFilterMode(self.rebase),
			Description: self.c.Tr.LcRebaseBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.MergeIntoCurrentBranch),
			Handler:     opts.Guards.OutsideFilterMode(self.merge),
			Description: self.c.Tr.LcMergeIntoCurrentBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.FastForward),
			Handler:     self.checkSelectedAndReal(self.fastForward),
			Description: self.c.Tr.FastForward,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.checkSelected(self.createResetMenu),
			Description: self.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.RenameBranch),
			Handler:     self.checkSelectedAndReal(self.rename),
			Description: self.c.Tr.LcRenameBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.SetUpstream),
			Handler:     self.checkSelected(self.setUpstream),
			Description: self.c.Tr.LcSetUnsetUpstream,
			OpensMenu:   true,
		},
	}
}

func (self *BranchesController) setUpstream(selectedBranch *models.Branch) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.SetUnsetUpstream,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.LcUnsetUpstream},
				OnPress: func() error {
					if err := self.git.Branch.UnsetUpstream(selectedBranch.Name); err != nil {
						return self.c.Error(err)
					}
					if err := self.c.Refresh(types.RefreshOptions{
						Mode: types.SYNC,
						Scope: []types.RefreshableView{
							types.BRANCHES,
							types.COMMITS,
						},
					}); err != nil {
						return self.c.Error(err)
					}
					return nil
				},
				Key: 'u',
			},
			{
				LabelColumns: []string{self.c.Tr.LcSetUpstream},
				OnPress: func() error {
					return self.helpers.Upstream.PromptForUpstreamWithoutInitialContent(selectedBranch, func(upstream string) error {
						upstreamRemote, upstreamBranch, err := self.helpers.Upstream.ParseUpstream(upstream)
						if err != nil {
							return self.c.Error(err)
						}

						if err := self.git.Branch.SetUpstream(upstreamRemote, upstreamBranch, selectedBranch.Name); err != nil {
							return self.c.Error(err)
						}
						if err := self.c.Refresh(types.RefreshOptions{
							Mode: types.SYNC,
							Scope: []types.RefreshableView{
								types.BRANCHES,
								types.COMMITS,
							},
						}); err != nil {
							return self.c.Error(err)
						}
						return nil
					})
				},
				Key: 's',
			},
		},
	})
}

func (self *BranchesController) Context() types.Context {
	return self.context()
}

func (self *BranchesController) context() *context.BranchesContext {
	return self.contexts.Branches
}

func (self *BranchesController) press(selectedBranch *models.Branch) error {
	if selectedBranch == self.helpers.Refs.GetCheckedOutRef() {
		return self.c.ErrorMsg(self.c.Tr.AlreadyCheckedOutBranch)
	}

	self.c.LogAction(self.c.Tr.Actions.CheckoutBranch)
	return self.helpers.Refs.CheckoutRef(selectedBranch.Name, types.CheckoutRefOptions{})
}

func (self *BranchesController) handleCreatePullRequest(selectedBranch *models.Branch) error {
	return self.createPullRequest(selectedBranch.Name, "")
}

func (self *BranchesController) handleCreatePullRequestMenu(selectedBranch *models.Branch) error {
	checkedOutBranch := self.helpers.Refs.GetCheckedOutRef()

	return self.createPullRequestMenu(selectedBranch, checkedOutBranch)
}

func (self *BranchesController) copyPullRequestURL() error {
	branch := self.context().GetSelected()

	branchExistsOnRemote := self.git.Remote.CheckRemoteBranchExists(branch.Name)

	if !branchExistsOnRemote {
		return self.c.Error(errors.New(self.c.Tr.NoBranchOnRemote))
	}

	url, err := self.helpers.Host.GetPullRequestURL(branch.Name, "")
	if err != nil {
		return self.c.Error(err)
	}
	self.c.LogAction(self.c.Tr.Actions.CopyPullRequestURL)
	if err := self.os.CopyToClipboard(url); err != nil {
		return self.c.Error(err)
	}

	self.c.Toast(self.c.Tr.PullRequestURLCopiedToClipboard)

	return nil
}

func (self *BranchesController) forceCheckout() error {
	branch := self.context().GetSelected()
	message := self.c.Tr.SureForceCheckout
	title := self.c.Tr.ForceCheckoutBranch

	return self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ForceCheckoutBranch)
			if err := self.git.Branch.Checkout(branch.Name, git_commands.CheckoutOptions{Force: true}); err != nil {
				_ = self.c.Error(err)
			}
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

func (self *BranchesController) checkoutByName() error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.BranchName + ":",
		FindSuggestionsFunc: self.helpers.Suggestions.GetRefsSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			self.c.LogAction("Checkout branch")
			return self.helpers.Refs.CheckoutRef(response, types.CheckoutRefOptions{
				OnRefNotFound: func(ref string) error {
					return self.c.Confirm(types.ConfirmOpts{
						Title:  self.c.Tr.BranchNotFoundTitle,
						Prompt: fmt.Sprintf("%s %s%s", self.c.Tr.BranchNotFoundPrompt, ref, "?"),
						HandleConfirm: func() error {
							return self.createNewBranchWithName(ref)
						},
					})
				},
			})
		},
	},
	)
}

func (self *BranchesController) createNewBranchWithName(newBranchName string) error {
	branch := self.context().GetSelected()
	if branch == nil {
		return nil
	}

	if err := self.git.Branch.New(newBranchName, branch.FullRefName()); err != nil {
		return self.c.Error(err)
	}

	self.context().SetSelectedLineIdx(0)
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *BranchesController) delete(branch *models.Branch) error {
	checkedOutBranch := self.helpers.Refs.GetCheckedOutRef()
	if checkedOutBranch.Name == branch.Name {
		return self.c.ErrorMsg(self.c.Tr.CantDeleteCheckOutBranch)
	}
	return self.deleteWithForce(branch, false)
}

func (self *BranchesController) deleteWithForce(selectedBranch *models.Branch, force bool) error {
	title := self.c.Tr.DeleteBranch
	var templateStr string
	if force {
		templateStr = self.c.Tr.ForceDeleteBranchMessage
	} else {
		templateStr = self.c.Tr.DeleteBranchMessage
	}
	message := utils.ResolvePlaceholderString(
		templateStr,
		map[string]string{
			"selectedBranchName": selectedBranch.Name,
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.DeleteBranch)
			if err := self.git.Branch.Delete(selectedBranch.Name, force); err != nil {
				errMessage := err.Error()
				if !force && strings.Contains(errMessage, "git branch -D ") {
					return self.deleteWithForce(selectedBranch, true)
				}
				return self.c.ErrorMsg(errMessage)
			}
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		},
	})
}

func (self *BranchesController) merge() error {
	selectedBranchName := self.context().GetSelected().Name
	return self.helpers.MergeAndRebase.MergeRefIntoCheckedOutBranch(selectedBranchName)
}

func (self *BranchesController) rebase() error {
	selectedBranchName := self.context().GetSelected().Name
	return self.helpers.MergeAndRebase.RebaseOntoRef(selectedBranchName)
}

func (self *BranchesController) fastForward(branch *models.Branch) error {
	if !branch.IsTrackingRemote() {
		return self.c.ErrorMsg(self.c.Tr.FwdNoUpstream)
	}
	if !branch.RemoteBranchStoredLocally() {
		return self.c.ErrorMsg(self.c.Tr.FwdNoLocalUpstream)
	}
	if branch.HasCommitsToPush() {
		return self.c.ErrorMsg(self.c.Tr.FwdCommitsToPush)
	}

	action := self.c.Tr.Actions.FastForwardBranch

	message := utils.ResolvePlaceholderString(
		self.c.Tr.Fetching,
		map[string]string{
			"from": fmt.Sprintf("%s/%s", branch.UpstreamRemote, branch.UpstreamBranch),
			"to":   branch.Name,
		},
	)

	return self.c.WithLoaderPanel(message, func() error {
		if branch == self.helpers.Refs.GetCheckedOutRef() {
			self.c.LogAction(action)

			err := self.git.Sync.Pull(
				git_commands.PullOptions{
					RemoteName:      branch.UpstreamRemote,
					BranchName:      branch.UpstreamBranch,
					FastForwardOnly: true,
				},
			)
			if err != nil {
				_ = self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		} else {
			self.c.LogAction(action)
			err := self.git.Sync.FastForward(branch.Name, branch.UpstreamRemote, branch.UpstreamBranch)
			if err != nil {
				_ = self.c.Error(err)
			}
			_ = self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.BRANCHES}})
		}

		return nil
	})
}

func (self *BranchesController) createResetMenu(selectedBranch *models.Branch) error {
	return self.helpers.Refs.CreateGitResetMenu(selectedBranch.Name)
}

func (self *BranchesController) rename(branch *models.Branch) error {
	promptForNewName := func() error {
		return self.c.Prompt(types.PromptOpts{
			Title:          self.c.Tr.NewBranchNamePrompt + " " + branch.Name + ":",
			InitialContent: branch.Name,
			HandleConfirm: func(newBranchName string) error {
				self.c.LogAction(self.c.Tr.Actions.RenameBranch)
				if err := self.git.Branch.Rename(branch.Name, newBranchName); err != nil {
					return self.c.Error(err)
				}

				// need to find where the branch is now so that we can re-select it. That means we need to refetch the branches synchronously and then find our branch
				_ = self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.BRANCHES}})

				// now that we've got our stuff again we need to find that branch and reselect it.
				for i, newBranch := range self.model.Branches {
					if newBranch.Name == newBranchName {
						self.context().SetSelectedLineIdx(i)
						if err := self.context().HandleRender(); err != nil {
							return err
						}
					}
				}

				return nil
			},
		})
	}

	// I could do an explicit check here for whether the branch is tracking a remote branch
	// but if we've selected it we'll already know that via Pullables and Pullables.
	// Bit of a hack but I'm lazy.
	if !branch.IsTrackingRemote() {
		return promptForNewName()
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:         self.c.Tr.LcRenameBranch,
		Prompt:        self.c.Tr.RenameBranchWarning,
		HandleConfirm: promptForNewName,
	})
}

func (self *BranchesController) newBranch(selectedBranch *models.Branch) error {
	return self.helpers.Refs.NewBranch(selectedBranch.FullRefName(), selectedBranch.RefName(), "")
}

func (self *BranchesController) createPullRequestMenu(selectedBranch *models.Branch, checkedOutBranch *models.Branch) error {
	menuItems := make([]*types.MenuItem, 0, 4)

	fromToLabelColumns := func(from string, to string) []string {
		return []string{fmt.Sprintf("%s → %s", from, to)}
	}

	menuItemsForBranch := func(branch *models.Branch) []*types.MenuItem {
		return []*types.MenuItem{
			{
				LabelColumns: fromToLabelColumns(branch.Name, self.c.Tr.LcDefaultBranch),
				OnPress: func() error {
					return self.createPullRequest(branch.Name, "")
				},
			},
			{
				LabelColumns: fromToLabelColumns(branch.Name, self.c.Tr.LcSelectBranch),
				OnPress: func() error {
					return self.c.Prompt(types.PromptOpts{
						Title:               branch.Name + " →",
						FindSuggestionsFunc: self.helpers.Suggestions.GetBranchNameSuggestionsFunc(),
						HandleConfirm: func(targetBranchName string) error {
							return self.createPullRequest(branch.Name, targetBranchName)
						},
					})
				},
			},
		}
	}

	if selectedBranch != checkedOutBranch {
		menuItems = append(menuItems,
			&types.MenuItem{
				LabelColumns: fromToLabelColumns(checkedOutBranch.Name, selectedBranch.Name),
				OnPress: func() error {
					return self.createPullRequest(checkedOutBranch.Name, selectedBranch.Name)
				},
			},
		)
		menuItems = append(menuItems, menuItemsForBranch(checkedOutBranch)...)
	}

	menuItems = append(menuItems, menuItemsForBranch(selectedBranch)...)

	return self.c.Menu(types.CreateMenuOptions{Title: fmt.Sprintf(self.c.Tr.CreatePullRequestOptions), Items: menuItems})
}

func (self *BranchesController) createPullRequest(from string, to string) error {
	url, err := self.helpers.Host.GetPullRequestURL(from, to)
	if err != nil {
		return self.c.Error(err)
	}

	self.c.LogAction(self.c.Tr.Actions.OpenPullRequest)

	if err := self.os.OpenLink(url); err != nil {
		return self.c.Error(err)
	}

	return nil
}

func (self *BranchesController) checkSelected(callback func(*models.Branch) error) func() error {
	return func() error {
		selectedItem := self.context().GetSelected()
		if selectedItem == nil {
			return nil
		}

		return callback(selectedItem)
	}
}

func (self *BranchesController) checkSelectedAndReal(callback func(*models.Branch) error) func() error {
	return func() error {
		selectedItem := self.context().GetSelected()
		if selectedItem == nil || !selectedItem.IsRealBranch() {
			return nil
		}

		return callback(selectedItem)
	}
}
