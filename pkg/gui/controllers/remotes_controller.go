package controllers

import (
	"fmt"
	"regexp"
	"slices"
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
		ListControllerTrait: NewListControllerTrait(
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
			Key:               opts.GetKey(opts.Config.Branches.AddForkRemote),
			Handler:           self.withItem(self.addFork),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.AddForkRemote,
			Tooltip:           self.c.Tr.AddForkRemoteTooltip,
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

func (self *RemotesController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			remote := self.context().GetSelected()
			if remote == nil {
				task = types.NewRenderStringTask("No remotes")
			} else {
				task = types.NewRenderStringTask(fmt.Sprintf("%s\nUrls:\n%s", style.FgGreen.Sprint(remote.Name), strings.Join(remote.Urls, "\n")))
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
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

	self.c.PostRefreshUpdate(remoteBranchesContext)

	self.c.Context().Push(remoteBranchesContext, types.OnFocusOpts{})
	return nil
}

// add remote and refresh the remotes list
func (self *RemotesController) addRemote(remoteName string, remoteUrl string) error {
	self.c.LogAction(self.c.Tr.Actions.AddRemote)
	err := self.c.Git().Remote.AddRemote(remoteName, remoteUrl)
	if err != nil {
		return err
	}

	// Do a sync refresh of the remotes so that we can select
	// the new one. Loading remotes is not expensive, so we can
	// afford it.
	self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.REMOTES},
		Mode:  types.SYNC,
	})
	return nil
}

func (self *RemotesController) selectRemoteAndCheckout(remoteName string, remoteUrl string, branchToCheckout string) error {
	// Select the remote
	for idx, remote := range self.c.Model().Remotes {
		if remote.Name == remoteName {
			self.c.Contexts().Remotes.SetSelection(idx)
			break
		}
	}

	// Fetch the remote
	return self.fetchAndCheckout(self.c.Contexts().Remotes.GetSelected(), branchToCheckout)
}

// add remote, then refresh and select it, and then fetch it and checkout the given branch if applicable
func (self *RemotesController) addRemoteHelper(remoteName string, remoteUrl string, branchToCheckout string) error {
	self.addRemote(remoteName, remoteUrl)

	return self.selectRemoteAndCheckout(remoteName, remoteUrl, branchToCheckout)
}

// the same as addRemoteHelper but we don't error if the remote with the already exists (the remote has to contain the given URL though)
func (self *RemotesController) addForkHelper(remoteName string, remoteUrl string, branchToCheckout string) error {
	for idx, remote := range self.c.Model().Remotes {
		if remote.Name == remoteName {
			hasUrl := slices.Contains(remote.Urls, remoteUrl)
			if !hasUrl {
				return fmt.Errorf("a remote named '%s' already exists with a different URL", remoteName)
			}
			self.c.Contexts().Remotes.SetSelection(idx)
			return self.fetchAndCheckout(remote, branchToCheckout)
		}
	}
	return self.addRemoteHelper(remoteName, remoteUrl, branchToCheckout)
}

func (self *RemotesController) add() error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.NewRemoteName,
		HandleConfirm: func(remoteName string) error {
			self.c.Prompt(types.PromptOpts{
				Title: self.c.Tr.NewRemoteUrl,
				HandleConfirm: func(remoteUrl string) error {
					return self.addRemoteHelper(remoteName, remoteUrl, "")
				},
			})

			return nil
		},
	})

	return nil
}

var (
	// 1. SCP-like SSH: git@host:owner[/subgroups]/repo(.git)
	sshScpRegex = regexp.MustCompile(`^(git@[^:]+:)([^/]+(?:/[^/]+)*)/([^/]+?)(\.git)?$`)

	// 2. SSH URL style: ssh://user@host[:port]/owner[/subgroups]/repo(.git)
	sshUrlRegex = regexp.MustCompile(`^(ssh://[^/]+/)([^/]+(?:/[^/]+)*)/([^/]+?)(\.git)?$`)

	// 3. HTTPS: https://host/owner[/subgroups]/repo(.git)
	httpRegex = regexp.MustCompile(`^(https?://[^/]+/)([^/]+(?:/[^/]+)*)/([^/]+?)(\.git)?$`)
)

// replaceForkUsername rewrites a Git remote URL to use the given fork username,
// keeping the repo name and host intact. Supports SCP-like SSH, SSH URL style, and HTTPS.
func replaceForkUsername(remoteUrl, forkUsername string) (string, error) {
	if forkUsername == "" {
		return "", fmt.Errorf("fork username cannot be empty")
	}
	if remoteUrl == "" {
		return "", fmt.Errorf("remote URL cannot be empty")
	}

	switch {
	case sshScpRegex.MatchString(remoteUrl):
		return sshScpRegex.ReplaceAllString(remoteUrl, "${1}"+forkUsername+"/$3$4"), nil
	case sshUrlRegex.MatchString(remoteUrl):
		return sshUrlRegex.ReplaceAllString(remoteUrl, "${1}"+forkUsername+"/$3$4"), nil
	case httpRegex.MatchString(remoteUrl):
		return httpRegex.ReplaceAllString(remoteUrl, "${1}"+forkUsername+"/$3$4"), nil
	default:
		return "", fmt.Errorf("unsupported or invalid remote URL: %s", remoteUrl)
	}
}

func (self *RemotesController) addFork(baseRemote *models.Remote) error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.AddForkRemoteUsername,
		HandleConfirm: func(forkUsername string) error {
			branchToCheckout := ""

			parts := strings.SplitN(forkUsername, ":", 2)
			if len(parts) == 2 {
				forkUsername = parts[0]
				branchToCheckout = parts[1]
			}
			baseUrl := baseRemote.Urls[0]
			remoteUrl, err := replaceForkUsername(baseUrl, forkUsername)
			if err != nil {
				return err
			}

			return self.addForkHelper(forkUsername, remoteUrl, branchToCheckout)
		},
	})

	return nil
}

func (self *RemotesController) remove(remote *models.Remote) error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.RemoveRemote,
		Prompt: self.c.Tr.RemoveRemotePrompt + " '" + remote.Name + "'?",
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.RemoveRemote)
			if err := self.c.Git().Remote.RemoveRemote(remote.Name); err != nil {
				return err
			}

			self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
			return nil
		},
	})

	return nil
}

func (self *RemotesController) edit(remote *models.Remote) error {
	editNameMessage := utils.ResolvePlaceholderString(
		self.c.Tr.EditRemoteName,
		map[string]string{
			"remoteName": remote.Name,
		},
	)

	self.c.Prompt(types.PromptOpts{
		Title:          editNameMessage,
		InitialContent: remote.Name,
		HandleConfirm: func(updatedRemoteName string) error {
			if updatedRemoteName != remote.Name {
				self.c.LogAction(self.c.Tr.Actions.UpdateRemote)
				if err := self.c.Git().Remote.RenameRemote(remote.Name, updatedRemoteName); err != nil {
					return err
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

			self.c.Prompt(types.PromptOpts{
				Title:          editUrlMessage,
				InitialContent: url,
				HandleConfirm: func(updatedRemoteUrl string) error {
					self.c.LogAction(self.c.Tr.Actions.UpdateRemote)
					if err := self.c.Git().Remote.UpdateRemoteUrl(updatedRemoteName, updatedRemoteUrl); err != nil {
						return err
					}
					self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES}})
					return nil
				},
			})

			return nil
		},
	})

	return nil
}

func (self *RemotesController) fetch(remote *models.Remote) error {
	return self.fetchAndCheckout(remote, "")
}

func (self *RemotesController) fetchAndCheckout(remote *models.Remote, branchName string) error {
	return self.c.WithInlineStatus(remote, types.ItemOperationFetching, context.REMOTES_CONTEXT_KEY, func(task gocui.Task) error {
		err := self.c.Git().Sync.FetchRemote(task, remote.Name)
		if err != nil {
			return err
		}
		if branchName != "" {
			err = self.c.Git().Branch.New(branchName, remote.Name+"/"+branchName)
		}
		self.c.Refresh(types.RefreshOptions{
			Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES},
			Mode:  types.ASYNC,
		})
		return err
	})
}
