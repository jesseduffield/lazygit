package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemotesController struct {
	baseController
	*controllerCommon
	context *context.RemotesContext

	setRemoteBranches func([]*models.RemoteBranch)
}

var _ types.IController = &RemotesController{}

func NewRemotesController(
	common *controllerCommon,
	setRemoteBranches func([]*models.RemoteBranch),
) *RemotesController {
	return &RemotesController{
		baseController:    baseController{},
		controllerCommon:  common,
		context:           common.contexts.Remotes,
		setRemoteBranches: setRemoteBranches,
	}
}

func (self *RemotesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.GoInto),
			Handler: self.checkSelected(self.enter),
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.FetchRemote),
			Handler:     self.checkSelected(self.fetch),
			Description: self.c.Tr.LcFetchRemote,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.add,
			Description: self.c.Tr.LcAddNewRemote,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.remove),
			Description: self.c.Tr.LcRemoveRemote,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelected(self.edit),
			Description: self.c.Tr.LcEditRemote,
		},
	}

	return bindings
}

func (self *RemotesController) GetOnClick() func() error {
	return self.checkSelected(self.enter)
}

func (self *RemotesController) enter(remote *models.Remote) error {
	// naive implementation: get the branches from the remote and render them to the list, change the context
	self.setRemoteBranches(remote.Branches)

	newSelectedLine := 0
	if len(remote.Branches) == 0 {
		newSelectedLine = -1
	}
	self.contexts.RemoteBranches.SetSelectedLineIdx(newSelectedLine)
	self.contexts.RemoteBranches.SetTitleRef(remote.Name)

	if err := self.c.PostRefreshUpdate(self.contexts.RemoteBranches); err != nil {
		return err
	}

	return self.c.PushContext(self.contexts.RemoteBranches)
}

func (self *RemotesController) add() error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.LcNewRemoteName,
		HandleConfirm: func(remoteName string) error {
			return self.c.Prompt(types.PromptOpts{
				Title: self.c.Tr.LcNewRemoteUrl,
				HandleConfirm: func(remoteUrl string) error {
					self.c.LogAction(self.c.Tr.Actions.AddRemote)
					if err := self.git.Remote.AddRemote(remoteName, remoteUrl); err != nil {
						return err
					}
					return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.REMOTES}})
				},
			})
		},
	})
}

func (self *RemotesController) remove(remote *models.Remote) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.LcRemoveRemote,
		Prompt: self.c.Tr.LcRemoveRemotePrompt + " '" + remote.Name + "'?",
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.RemoveRemote)
			if err := self.git.Remote.RemoveRemote(remote.Name); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
		},
	})
}

func (self *RemotesController) edit(remote *models.Remote) error {
	editNameMessage := utils.ResolvePlaceholderString(
		self.c.Tr.LcEditRemoteName,
		map[string]string{
			"remoteName": remote.Name,
		},
	)

	return self.c.Prompt(types.PromptOpts{
		Title:          editNameMessage,
		InitialContent: remote.Name,
		HandleConfirm: func(updatedRemoteName string) error {
			if updatedRemoteName != remote.Name {
				self.c.LogAction(self.c.Tr.Actions.UpdateRemote)
				if err := self.git.Remote.RenameRemote(remote.Name, updatedRemoteName); err != nil {
					return self.c.Error(err)
				}
			}

			editUrlMessage := utils.ResolvePlaceholderString(
				self.c.Tr.LcEditRemoteUrl,
				map[string]string{
					"remoteName": updatedRemoteName,
				},
			)

			urls := remote.Urls
			url := ""
			if len(urls) > 0 {
				url = urls[0]
			}

			return self.c.Prompt(types.PromptOpts{
				Title:          editUrlMessage,
				InitialContent: url,
				HandleConfirm: func(updatedRemoteUrl string) error {
					self.c.LogAction(self.c.Tr.Actions.UpdateRemote)
					if err := self.git.Remote.UpdateRemoteUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
						return self.c.Error(err)
					}
					return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
				},
			})
		},
	})
}

func (self *RemotesController) fetch(remote *models.Remote) error {
	return self.c.WithWaitingStatus(self.c.Tr.FetchingRemoteStatus, func() error {
		err := self.git.Sync.FetchRemote(remote.Name)
		if err != nil {
			_ = self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
	})
}

func (self *RemotesController) checkSelected(callback func(*models.Remote) error) func() error {
	return func() error {
		file := self.context.GetSelected()
		if file == nil {
			return nil
		}

		return callback(file)
	}
}

func (self *RemotesController) Context() types.Context {
	return self.context
}
