package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemoteBranchesController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &RemoteBranchesController{}

func NewRemoteBranchesController(
	common *controllerCommon,
) *RemoteBranchesController {
	return &RemoteBranchesController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *RemoteBranchesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key: opts.GetKey(opts.Config.Universal.Select),
			// gonna use the exact same handler as the 'n' keybinding because everybody wants this to happen when they checkout a remote branch
			Handler:     self.checkSelected(self.newLocalBranch),
			Description: self.c.Tr.LcCheckout,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.checkSelected(self.newLocalBranch),
			Description: self.c.Tr.LcNewBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.MergeIntoCurrentBranch),
			Handler:     opts.Guards.OutsideFilterMode(self.checkSelected(self.merge)),
			Description: self.c.Tr.LcMergeIntoCurrentBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.RebaseBranch),
			Handler:     opts.Guards.OutsideFilterMode(self.checkSelected(self.rebase)),
			Description: self.c.Tr.LcRebaseBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.delete),
			Description: self.c.Tr.LcDeleteBranch,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.SetUpstream),
			Handler:     self.checkSelected(self.setAsUpstream),
			Description: self.c.Tr.LcSetAsUpstream,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     self.escape,
			Description: self.c.Tr.ReturnToRemotesList,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.checkSelected(self.createResetMenu),
			Description: self.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
	}
}

func (self *RemoteBranchesController) Context() types.Context {
	return self.context()
}

func (self *RemoteBranchesController) context() *context.RemoteBranchesContext {
	return self.contexts.RemoteBranches
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

func (self *RemoteBranchesController) escape() error {
	return self.c.PushContext(self.contexts.Remotes)
}

func (self *RemoteBranchesController) delete(selectedBranch *models.RemoteBranch) error {
	message := fmt.Sprintf("%s '%s'?", self.c.Tr.DeleteRemoteBranchMessage, selectedBranch.FullName())

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DeleteRemoteBranch,
		Prompt: message,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.DeletingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.DeleteRemoteBranch)
				err := self.git.Remote.DeleteRemoteBranch(selectedBranch.RemoteName, selectedBranch.Name)
				if err != nil {
					_ = self.c.Error(err)
				}

				return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			})
		},
	})
}

func (self *RemoteBranchesController) merge(selectedBranch *models.RemoteBranch) error {
	return self.helpers.MergeAndRebase.MergeRefIntoCheckedOutBranch(selectedBranch.FullName())
}

func (self *RemoteBranchesController) rebase(selectedBranch *models.RemoteBranch) error {
	return self.helpers.MergeAndRebase.RebaseOntoRef(selectedBranch.FullName())
}

func (self *RemoteBranchesController) createResetMenu(selectedBranch *models.RemoteBranch) error {
	return self.helpers.Refs.CreateGitResetMenu(selectedBranch.FullName())
}

func (self *RemoteBranchesController) setAsUpstream(selectedBranch *models.RemoteBranch) error {
	checkedOutBranch := self.helpers.Refs.GetCheckedOutRef()

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
			if err := self.git.Branch.SetUpstream(selectedBranch.RemoteName, selectedBranch.Name, checkedOutBranch.Name); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
		},
	})
}

func (self *RemoteBranchesController) newLocalBranch(selectedBranch *models.RemoteBranch) error {
	// will set to the remote's branch name without the remote name
	nameSuggestion := strings.SplitAfterN(selectedBranch.RefName(), "/", 2)[1]

	return self.helpers.Refs.NewBranch(selectedBranch.RefName(), selectedBranch.RefName(), nameSuggestion)
}
