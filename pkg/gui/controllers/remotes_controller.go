package controllers

import (
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemotesController struct {
	c       *ControllerCommon
	context types.IListContext
	git     *commands.GitCommand

	getSelectedRemote func() *models.Remote
	setRemoteBranches func([]*models.RemoteBranch)
	allContexts       context.ContextTree
	fetchMutex        *sync.Mutex
}

var _ types.IController = &RemotesController{}

func NewRemotesController(
	c *ControllerCommon,
	context types.IListContext,
	git *commands.GitCommand,
	allContexts context.ContextTree,
	getSelectedRemote func() *models.Remote,
	setRemoteBranches func([]*models.RemoteBranch),
	fetchMutex *sync.Mutex,
) *RemotesController {
	return &RemotesController{
		c:                 c,
		git:               git,
		allContexts:       allContexts,
		context:           context,
		getSelectedRemote: getSelectedRemote,
		setRemoteBranches: setRemoteBranches,
		fetchMutex:        fetchMutex,
	}
}

func (self *RemotesController) Keybindings(getKey func(key string) interface{}, config config.KeybindingConfig, guards types.KeybindingGuards) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     getKey(config.Universal.GoInto),
			Handler: self.checkSelected(self.enter),
		},
		{
			Key:     gocui.MouseLeft,
			Handler: func() error { return self.context.HandleClick(self.checkSelected(self.enter)) },
		},
		{
			Key:         getKey(config.Branches.FetchRemote),
			Handler:     self.checkSelected(self.fetch),
			Description: self.c.Tr.LcFetchRemote,
		},
		{
			Key:         getKey(config.Universal.New),
			Handler:     self.add,
			Description: self.c.Tr.LcAddNewRemote,
		},
		{
			Key:         getKey(config.Universal.Remove),
			Handler:     self.checkSelected(self.remove),
			Description: self.c.Tr.LcRemoveRemote,
		},
		{
			Key:         getKey(config.Universal.Edit),
			Handler:     self.checkSelected(self.edit),
			Description: self.c.Tr.LcEditRemote,
		},
	}

	return append(bindings, self.context.Keybindings(getKey, config, guards)...)
}

func (self *RemotesController) enter(remote *models.Remote) error {
	// naive implementation: get the branches from the remote and render them to the list, change the context
	self.setRemoteBranches(remote.Branches)

	newSelectedLine := 0
	if len(remote.Branches) == 0 {
		newSelectedLine = -1
	}
	self.allContexts.RemoteBranches.GetPanelState().SetSelectedLineIdx(newSelectedLine)

	return self.c.PushContext(self.allContexts.RemoteBranches)
}

func (self *RemotesController) add() error {
	return self.c.Prompt(popup.PromptOpts{
		Title: self.c.Tr.LcNewRemoteName,
		HandleConfirm: func(remoteName string) error {
			return self.c.Prompt(popup.PromptOpts{
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
	return self.c.Ask(popup.AskOpts{
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

	return self.c.Prompt(popup.PromptOpts{
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

			return self.c.Prompt(popup.PromptOpts{
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
		self.fetchMutex.Lock()
		defer self.fetchMutex.Unlock()

		err := self.git.Sync.FetchRemote(remote.Name)
		if err != nil {
			_ = self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
	})
}

func (self *RemotesController) checkSelected(callback func(*models.Remote) error) func() error {
	return func() error {
		file := self.getSelectedRemote()
		if file == nil {
			return nil
		}

		return callback(file)
	}
}

func (self *RemotesController) Context() types.Context {
	return self.context
}
