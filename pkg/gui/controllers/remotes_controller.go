package controllers

import (
	"fmt"
	"net/url"
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

func (self *RemotesController) addRemoteHelper(remoteName string, remoteUrl string) error {
	self.c.LogAction(self.c.Tr.Actions.AddRemote)
	if err := self.c.Git().Remote.AddRemote(remoteName, remoteUrl); err != nil {
		return err
	}

	// Do a sync refresh of the remotes so that we can select
	// the new one. Loading remotes is not expensive, so we can
	// afford it.
	self.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.REMOTES},
		Mode:  types.SYNC,
	})

	// Select the new remote
	for idx, remote := range self.c.Model().Remotes {
		if remote.Name == remoteName {
			self.c.Contexts().Remotes.SetSelection(idx)
			break
		}
	}

	// Fetch the new remote
	return self.fetch(self.c.Contexts().Remotes.GetSelected())
}

func (self *RemotesController) add() error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.NewRemoteName,
		HandleConfirm: func(remoteName string) error {
			self.c.Prompt(types.PromptOpts{
				Title: self.c.Tr.NewRemoteUrl,
				HandleConfirm: func(remoteUrl string) error {
					return self.addRemoteHelper(remoteName, remoteUrl)
				},
			})

			return nil
		},
	})

	return nil
}

// replaceForkUsername replaces the "owner" part of a git remote URL with forkUsername,
// preserving the repo name (last path segment) and everything else (host, scheme, port, .git suffix).
// Supported forms:
//   - SSH scp-like:   git@host:owner[/subgroups]/repo(.git)
//   - HTTPS/HTTP:     https://host/owner[/subgroups]/repo(.git)
//
// Rules:
//   - If there are fewer than 2 path segments (i.e., no clear owner+repo), return an error.
//   - For multi-segment paths (e.g., group/subgroup/repo), the entire prefix is replaced by forkUsername.
func replaceForkUsername(remoteUrl, forkUsername string) (string, error) {
	if forkUsername == "" {
		return "", fmt.Errorf("Fork username cannot be empty")
	}
	if remoteUrl == "" {
		return "", fmt.Errorf("Remote url cannot be empty")
	}

	// SSH scp-like (most common): git@host:path
	if isScpLikeSSH(remoteUrl) {
		colon := strings.IndexByte(remoteUrl, ':')
		if colon == -1 {
			return "", fmt.Errorf("Invalid SSH remote URL (missing ':'): %s", remoteUrl)
		}
		path := remoteUrl[colon+1:] // e.g. owner/repo(.git) or group/sub/repo(.git)
		segments := splitNonEmpty(path, "/")
		if len(segments) < 2 {
			return "", fmt.Errorf("Remote URL must include owner and repo: %s", remoteUrl)
		}
		last := segments[len(segments)-1] // repo(.git)
		newPath := forkUsername + "/" + last
		return remoteUrl[:colon+1] + newPath, nil
	}

	// Try URL parsing for http(s) (and reject anything else).
	u, err := url.Parse(remoteUrl)
	if err != nil {
		return "", fmt.Errorf("Invalid remote URL: %w", err)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return "", fmt.Errorf("Unsupported remote URL scheme: %s", u.Scheme)
	}

	// u.Path like "/owner[/subgroups]/repo(.git)" or "" or "/"
	path := strings.Trim(u.Path, "/")
	segments := splitNonEmpty(path, "/")
	if len(segments) < 2 {
		return "", fmt.Errorf("Remote URL must include owner and repo: %s", remoteUrl)
	}

	last := segments[len(segments)-1] // repo(.git)
	u.Path = "/" + forkUsername + "/" + last

	// Preserve trailing slash only if it existed and wasn't empty
	// (remotes rarely care, but we'll avoid adding one)
	return u.String(), nil
}

func isScpLikeSSH(s string) bool {
	// Minimal heuristic: "<user>@<host>:<path>"
	at := strings.IndexByte(s, '@')
	colon := strings.IndexByte(s, ':')
	return at > 0 && colon > at
}

func splitNonEmpty(s, sep string) []string {
	raw := strings.Split(s, sep)
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (self *RemotesController) addFork(baseRemote *models.Remote) error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.AddForkRemoteUsername,
		HandleConfirm: func(forkUsername string) error {
			self.c.Prompt(types.PromptOpts{
				Title:          self.c.Tr.NewRemoteName,
				InitialContent: forkUsername,
				HandleConfirm: func(remoteName string) error {
					if forkUsername == "" {
						return fmt.Errorf("Fork username cannot be empty")
					}
					if len(baseRemote.Urls) == 0 {
						return fmt.Errorf("Base remote must have url")
					}
					url := baseRemote.Urls[0]
					if url == "" {
						return fmt.Errorf("Base remote url cannot be empty")
					}
					remoteUrl, err := replaceForkUsername(url, forkUsername)
					if err != nil {
						return fmt.Errorf("Failed to replace fork username in remote URL: `%w`, make sure it's a valid url", err)
					}

					return self.addRemoteHelper(remoteName, remoteUrl)
				},
			})

			return nil
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
	return self.c.WithInlineStatus(remote, types.ItemOperationFetching, context.REMOTES_CONTEXT_KEY, func(task gocui.Task) error {
		err := self.c.Git().Sync.FetchRemote(task, remote.Name)
		if err != nil {
			return err
		}

		self.c.Refresh(types.RefreshOptions{
			Scope: []types.RefreshableView{types.BRANCHES, types.REMOTES},
			Mode:  types.ASYNC,
		})
		return nil
	})
}
