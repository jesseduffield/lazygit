package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RemotesController struct {
	baseController
	c *ControllerCommon

	setRemoteBranches func([]*models.RemoteBranch)
}

var _ types.IController = &RemotesController{}

func NewRemotesController(
	common *ControllerCommon,
	setRemoteBranches func([]*models.RemoteBranch),
) *RemotesController {
	return &RemotesController{
		baseController:    baseController{},
		c:                 common,
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
			Description: self.c.Tr.FetchRemote,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.add,
			Description: self.c.Tr.AddNewRemote,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.remove),
			Description: self.c.Tr.RemoveRemote,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelected(self.edit),
			Description: self.c.Tr.EditRemote,
		},
	}

	return bindings
}

func (self *RemotesController) Context() types.Context {
	return self.context()
}

func (self *RemotesController) context() *context.RemotesContext {
	return self.c.Contexts().Remotes
}

func (self *RemotesController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			var task types.UpdateTask
			remote := self.context().GetSelected()
			if remote == nil {
				task = types.NewRenderStringTask("No remotes")
			} else {
				task = types.NewRenderStringTask(fmt.Sprintf("%s\nUrls:\n%s", style.FgGreen.Sprint(remote.Name), strings.Join(remote.Urls, "\n")))
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Remote",
					Task:  task,
				},
			})
		})
	}
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
	remoteBranchesContext := self.c.Contexts().RemoteBranches
	remoteBranchesContext.SetSelectedLineIdx(newSelectedLine)
	remoteBranchesContext.SetTitleRef(remote.Name)
	remoteBranchesContext.SetParentContext(self.Context())

	if err := self.c.PostRefreshUpdate(remoteBranchesContext); err != nil {
		return err
	}

	return self.c.PushContext(remoteBranchesContext)
}

func (self *RemotesController) add() error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.NewRemoteName,
		HandleConfirm: func(remoteName string) error {
			return self.c.Prompt(types.PromptOpts{
				Title: self.c.Tr.NewRemoteUrl,
				HandleConfirm: func(remoteUrl string) error {
					self.c.LogAction(self.c.Tr.Actions.AddRemote)
					if err := self.c.Git().Remote.AddRemote(remoteName, remoteUrl); err != nil {
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
		Title:  self.c.Tr.RemoveRemote,
		Prompt: self.c.Tr.RemoveRemotePrompt + " '" + remote.Name + "'?",
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.RemoveRemote)
			if err := self.c.Git().Remote.RemoveRemote(remote.Name); err != nil {
				return self.c.Error(err)
			}

			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
		},
	})
}

func (self *RemotesController) edit(remote *models.Remote) error {
	editNameMessage := utils.ResolvePlaceholderString(
		self.c.Tr.EditRemoteName,
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
				if err := self.c.Git().Remote.RenameRemote(remote.Name, updatedRemoteName); err != nil {
					return self.c.Error(err)
				}
			}

			editUrlMessage := utils.ResolvePlaceholderString(
				self.c.Tr.EditRemoteUrl,
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
					if err := self.c.Git().Remote.UpdateRemoteUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
						return self.c.Error(err)
					}
					return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
				},
			})
		},
	})
}

func (self *RemotesController) fetch(remote *models.Remote) error {
	return self.c.WithWaitingStatus(self.c.Tr.FetchingRemoteStatus, func(task gocui.Task) error {
		err := self.c.Git().Sync.FetchRemote(task, remote.Name)
		if err != nil {
			_ = self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
	})
}

func (self *RemotesController) checkSelected(callback func(*models.Remote) error) func() error {
	return func() error {
		file := self.context().GetSelected()
		if file == nil {
			return nil
		}

		return callback(file)
	}
}
