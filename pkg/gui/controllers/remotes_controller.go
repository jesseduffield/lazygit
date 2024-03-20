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
	*ListControllerTrait[*models.Remote]
	c *ControllerCommon

	setRemoteBranches func([]*models.RemoteBranch)
}

var _ types.IController = &RemotesController{}

func NewRemotesController(
	c *ControllerCommon,
	setRemoteBranches func([]*models.RemoteBranch),
) *RemotesController {
	return &RemotesController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.Remote](
			c,
			c.Contexts().Remotes,
			c.Contexts().Remotes.GetSelected,
			c.Contexts().Remotes.GetSelectedItems,
		),
		c:                 c,
		setRemoteBranches: setRemoteBranches,
	}
}

func (self *RemotesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ViewBranches,
			DisplayOnScreen:   true,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.New),
			Handler:         self.add,
			Description:     self.c.Tr.NewRemote,
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.remove),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Remove,
			Tooltip:           self.c.Tr.RemoveRemoteTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Edit),
			Handler:           self.withItem(self.edit),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Edit,
			Tooltip:           self.c.Tr.EditRemoteTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Branches.FetchRemote),
			Handler:           self.withItem(self.fetch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Fetch,
			Tooltip:           self.c.Tr.FetchRemoteTooltip,
			DisplayOnScreen:   true,
		},
	}

	return bindings
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
	return self.withItemGraceful(self.enter)
}

func (self *RemotesController) enter(remote *models.Remote) error {
	// naive implementation: get the branches from the remote and render them to the list, change the context
	self.setRemoteBranches(remote.Branches)

	newSelectedLine := 0
	if len(remote.Branches) == 0 {
		newSelectedLine = -1
	}
	remoteBranchesContext := self.c.Contexts().RemoteBranches
	remoteBranchesContext.SetSelection(newSelectedLine)
	remoteBranchesContext.SetTitleRef(remote.Name)
	remoteBranchesContext.SetParentContext(self.Context())
	remoteBranchesContext.GetView().TitlePrefix = self.Context().GetView().TitlePrefix

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

					// Do a sync refresh of the remotes so that we can select
					// the new one. Loading remotes is not expensive, so we can
					// afford it.
					if err := self.c.Refresh(types.RefreshOptions{
						Scope: []types.RefreshableView{types.REMOTES},
						Mode:  types.SYNC,
					}); err != nil {
						return err
					}

					// Select the new remote
					for idx, remote := range self.c.Model().Remotes {
						if remote.Name == remoteName {
							self.c.Contexts().Remotes.SetSelection(idx)
							break
						}
					}

					// Fetch the new remote
					return self.fetch(self.c.Contexts().Remotes.GetSelected())
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
	return self.c.WithInlineStatus(remote, types.ItemOperationFetching, context.REMOTES_CONTEXT_KEY, func(task gocui.Task) error {
		err := self.c.Git().Sync.FetchRemote(task, remote.Name)
		if err != nil {
			_ = self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{
			Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES},
			Mode:  types.ASYNC,
		})
	})
}
