package controllers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemoteBranchesController struct {
	baseController
	*ListControllerTrait[*models.RemoteBranch]
	c *ControllerCommon
}

var _ types.IController = &RemoteBranchesController{}

func NewRemoteBranchesController(
	c *ControllerCommon,
) *RemoteBranchesController {
	return &RemoteBranchesController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.RemoteBranch](
			c,
			c.Contexts().RemoteBranches,
			c.Contexts().RemoteBranches.GetSelected,
			c.Contexts().RemoteBranches.GetSelectedItems,
		),
		c: c,
	}
}

func (self *RemoteBranchesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.checkoutBranch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.RemoteBranchCheckoutTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.New),
			Handler:           self.withItem(self.newLocalBranch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewBranch,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.MergeIntoCurrentBranch),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.merge)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Merge,
			Tooltip:           self.c.Tr.MergeBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.RebaseBranch),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.rebase)),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.RebaseBranch,
			Tooltip:           self.c.Tr.RebaseBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.delete),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Delete,
			Tooltip:           self.c.Tr.DeleteRemoteBranchTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.SetUpstream),
			Handler:           self.withItem(self.setAsUpstream),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.SetAsUpstream,
			Tooltip:           self.c.Tr.SetAsUpstreamTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.SortOrder),
			Handler:     self.createSortMenu,
			Description: self.c.Tr.SortOrder,
			OpensMenu:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:           self.withItem(self.createResetMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ViewResetOptions,
			Tooltip:           self.c.Tr.ResetTooltip,
			OpensMenu:         true,
		},
	}
}

func (self *RemoteBranchesController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			var task types.UpdateTask
			remoteBranch := self.context().GetSelected()
			if remoteBranch == nil {
				task = types.NewRenderStringTask("No branches for this remote")
			} else {
				cmdObj := self.c.Git().Branch.GetGraphCmdObj(remoteBranch.FullRefName())
				task = types.NewRunCommandTask(cmdObj.GetCmd())
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Remote Branch",
					Task:  task,
				},
			})
		})
	}
}

func (self *RemoteBranchesController) context() *context.RemoteBranchesContext {
	return self.c.Contexts().RemoteBranches
}

func (self *RemoteBranchesController) delete(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().BranchesHelper.ConfirmDeleteRemote(selectedBranch.RemoteName, selectedBranch.Name)
}

func (self *RemoteBranchesController) merge(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.MergeRefIntoCheckedOutBranch(selectedBranch.FullName())
}

func (self *RemoteBranchesController) rebase(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.RebaseOntoRef(selectedBranch.FullName())
}

func (self *RemoteBranchesController) createSortMenu() error {
	return self.c.Helpers().Refs.CreateSortOrderMenu([]string{"alphabetical", "date"}, func(sortOrder string) error {
		if self.c.GetAppState().RemoteBranchSortOrder != sortOrder {
			self.c.GetAppState().RemoteBranchSortOrder = sortOrder
			self.c.SaveAppStateAndLogError()
			self.c.Contexts().RemoteBranches.SetSelection(0)
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.REMOTES}})
		}
		return nil
	})
}

func (self *RemoteBranchesController) createResetMenu(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(selectedBranch.FullName())
}

func (self *RemoteBranchesController) setAsUpstream(selectedBranch *models.RemoteBranch) error {
	checkedOutBranch := self.c.Helpers().Refs.GetCheckedOutRef()

	message := utils.ResolvePlaceholderString(
		self.c.Tr.SetUpstreamMessage,
		map[string]string{
			"checkedOut": checkedOutBranch.Name,
			"selected":   selectedBranch.FullName(),
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.SetUpstreamTitle,
		Prompt: message,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.SetBranchUpstream)
			if err := self.c.Git().Branch.SetUpstream(selectedBranch.RemoteName, selectedBranch.Name, checkedOutBranch.Name); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
		},
	})
}

func (self *RemoteBranchesController) newLocalBranch(selectedBranch *models.RemoteBranch) error {
	// will set to the remote's branch name without the remote name
	nameSuggestion := strings.SplitAfterN(selectedBranch.RefName(), "/", 2)[1]

	return self.c.Helpers().Refs.NewBranch(selectedBranch.RefName(), selectedBranch.RefName(), nameSuggestion)
}

func (self *RemoteBranchesController) checkoutBranch(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().Refs.CheckoutRemoteBranch(selectedBranch.FullName(), selectedBranch.Name)
}
