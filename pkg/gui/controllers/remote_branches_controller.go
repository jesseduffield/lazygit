package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemoteBranchesController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &RemoteBranchesController{}

func NewRemoteBranchesController(
	common *ControllerCommon,
) *RemoteBranchesController {
	return &RemoteBranchesController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *RemoteBranchesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key: opts.GetKey(opts.Config.Universal.Select),
			// gonna use the exact same handler as the 'n' keybinding because everybody wants this to happen when they checkout a remote branch
			Handler:     self.checkSelected(self.newLocalBranch),
			Description: self.c.Tr.Checkout,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.checkSelected(self.newLocalBranch),
			Description: self.c.Tr.NewBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.MergeIntoCurrentBranch),
			Handler:     opts.Guards.OutsideFilterMode(self.checkSelected(self.merge)),
			Description: self.c.Tr.MergeIntoCurrentBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.RebaseBranch),
			Handler:     opts.Guards.OutsideFilterMode(self.checkSelected(self.rebase)),
			Description: self.c.Tr.RebaseBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.delete),
			Description: self.c.Tr.DeleteBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.SetUpstream),
			Handler:     self.checkSelected(self.setAsUpstream),
			Description: self.c.Tr.SetAsUpstream,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.checkSelected(self.createResetMenu),
			Description: self.c.Tr.ViewResetOptions,
			OpensMenu:   true,
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

func (self *RemoteBranchesController) Context() types.Context {
	return self.context()
}

func (self *RemoteBranchesController) context() *context.RemoteBranchesContext {
	return self.c.Contexts().RemoteBranches
}

func (self *RemoteBranchesController) checkSelected(callback func(*models.RemoteBranch) error) func() error {
	return func() error {
		selectedItem := self.context().GetSelected()
		if selectedItem == nil {
			return nil
		}

		return callback(selectedItem)
	}
}

func (self *RemoteBranchesController) delete(selectedBranch *models.RemoteBranch) error {
	message := fmt.Sprintf("%s '%s'?", self.c.Tr.DeleteRemoteBranchMessage, selectedBranch.FullName())

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DeleteRemoteBranch,
		Prompt: message,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func(task gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.DeleteRemoteBranch)
				err := self.c.Git().Remote.DeleteRemoteBranch(task, selectedBranch.RemoteName, selectedBranch.Name)
				if err != nil {
					_ = self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			})
		},
	})
}

func (self *RemoteBranchesController) merge(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.MergeRefIntoCheckedOutBranch(selectedBranch.FullName())
}

func (self *RemoteBranchesController) rebase(selectedBranch *models.RemoteBranch) error {
	return self.c.Helpers().MergeAndRebase.RebaseOntoRef(selectedBranch.FullName())
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
