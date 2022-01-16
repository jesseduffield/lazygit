package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SyncController struct {
	// I've said publicly that I'm against single-letter variable names but in this
	// case I would actually prefer a _zero_ letter variable name in the form of
	// struct embedding, but Go does not allow hiding public fields in an embedded struct
	// to the client
	c   *ControllerCommon
	git *commands.GitCommand

	getCheckedOutBranch func() *models.Branch
	suggestionsHelper   ISuggestionsHelper
	getSuggestedRemote  func() string
	checkMergeOrRebase  func(error) error
}

var _ types.IController = &SyncController{}

func NewSyncController(
	c *ControllerCommon,
	git *commands.GitCommand,
	getCheckedOutBranch func() *models.Branch,
	suggestionsHelper ISuggestionsHelper,
	getSuggestedRemote func() string,
	checkMergeOrRebase func(error) error,
) *SyncController {
	return &SyncController{
		c:   c,
		git: git,

		getCheckedOutBranch: getCheckedOutBranch,
		suggestionsHelper:   suggestionsHelper,
		getSuggestedRemote:  getSuggestedRemote,
		checkMergeOrRebase:  checkMergeOrRebase,
	}
}

func (self *SyncController) Keybindings(getKey func(key string) interface{}, config config.KeybindingConfig, guards types.KeybindingGuards) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         getKey(config.Universal.PushFiles),
			Handler:     guards.NoPopupPanel(self.HandlePush),
			Description: self.c.Tr.LcPush,
		},
		{
			Key:         getKey(config.Universal.PullFiles),
			Handler:     guards.NoPopupPanel(self.HandlePull),
			Description: self.c.Tr.LcPull,
		},
	}

	return bindings
}

func (self *SyncController) Context() types.Context {
	return nil
}

func (self *SyncController) HandlePush() error {
	return self.branchCheckedOut(self.push)()
}

func (self *SyncController) HandlePull() error {
	return self.branchCheckedOut(self.pull)()
}

func (self *SyncController) branchCheckedOut(f func(*models.Branch) error) func() error {
	return func() error {
		currentBranch := self.getCheckedOutBranch()
		if currentBranch == nil {
			// need to wait for branches to refresh
			return nil
		}

		return f(currentBranch)
	}
}

func (self *SyncController) push(currentBranch *models.Branch) error {
	// if we have pullables we'll ask if the user wants to force push
	if currentBranch.IsTrackingRemote() {
		opts := pushOpts{
			force:          false,
			upstreamRemote: currentBranch.UpstreamRemote,
			upstreamBranch: currentBranch.UpstreamBranch,
		}
		if currentBranch.HasCommitsToPull() {
			opts.force = true
			return self.requestToForcePush(opts)
		} else {
			return self.pushAux(opts)
		}
	} else {
		if self.git.Config.GetPushToCurrent() {
			return self.pushAux(pushOpts{setUpstream: true})
		} else {
			return self.promptForUpstream(currentBranch, func(upstream string) error {
				var upstreamBranch, upstreamRemote string
				split := strings.Split(upstream, " ")
				if len(split) == 2 {
					upstreamRemote = split[0]
					upstreamBranch = split[1]
				} else {
					upstreamRemote = upstream
					upstreamBranch = ""
				}

				return self.pushAux(pushOpts{
					force:          false,
					upstreamRemote: upstreamRemote,
					upstreamBranch: upstreamBranch,
					setUpstream:    true,
				})
			})
		}
	}
}

func (self *SyncController) pull(currentBranch *models.Branch) error {
	action := self.c.Tr.Actions.Pull

	// if we have no upstream branch we need to set that first
	if !currentBranch.IsTrackingRemote() {
		return self.promptForUpstream(currentBranch, func(upstream string) error {
			var upstreamBranch, upstreamRemote string
			split := strings.Split(upstream, " ")
			if len(split) != 2 {
				return self.c.ErrorMsg(self.c.Tr.InvalidUpstream)
			}

			upstreamRemote = split[0]
			upstreamBranch = split[1]

			if err := self.git.Branch.SetCurrentBranchUpstream(upstreamRemote, upstreamBranch); err != nil {
				errorMessage := err.Error()
				if strings.Contains(errorMessage, "does not exist") {
					errorMessage = fmt.Sprintf("upstream branch %s not found.\nIf you expect it to exist, you should fetch (with 'f').\nOtherwise, you should push (with 'shift+P')", upstream)
				}
				return self.c.ErrorMsg(errorMessage)
			}
			return self.PullAux(PullFilesOptions{UpstreamRemote: upstreamRemote, UpstreamBranch: upstreamBranch, Action: action})
		})
	}

	return self.PullAux(PullFilesOptions{UpstreamRemote: currentBranch.UpstreamRemote, UpstreamBranch: currentBranch.UpstreamBranch, Action: action})
}

func (self *SyncController) promptForUpstream(currentBranch *models.Branch, onConfirm func(string) error) error {
	suggestedRemote := self.getSuggestedRemote()

	return self.c.Prompt(popup.PromptOpts{
		Title:               self.c.Tr.EnterUpstream,
		InitialContent:      suggestedRemote + " " + currentBranch.Name,
		FindSuggestionsFunc: self.suggestionsHelper.GetRemoteBranchesSuggestionsFunc(" "),
		HandleConfirm:       onConfirm,
	})
}

type PullFilesOptions struct {
	UpstreamRemote  string
	UpstreamBranch  string
	FastForwardOnly bool
	Action          string
}

func (self *SyncController) PullAux(opts PullFilesOptions) error {
	return self.c.WithLoaderPanel(self.c.Tr.PullWait, func() error {
		return self.pullWithLock(opts)
	})
}

func (self *SyncController) pullWithLock(opts PullFilesOptions) error {
	self.c.LogAction(opts.Action)

	err := self.git.Sync.Pull(
		git_commands.PullOptions{
			RemoteName:      opts.UpstreamRemote,
			BranchName:      opts.UpstreamBranch,
			FastForwardOnly: opts.FastForwardOnly,
		},
	)

	return self.checkMergeOrRebase(err)
}

type pushOpts struct {
	force          bool
	upstreamRemote string
	upstreamBranch string
	setUpstream    bool
}

func (self *SyncController) pushAux(opts pushOpts) error {
	return self.c.WithLoaderPanel(self.c.Tr.PushWait, func() error {
		self.c.LogAction(self.c.Tr.Actions.Push)
		err := self.git.Sync.Push(git_commands.PushOpts{
			Force:          opts.force,
			UpstreamRemote: opts.upstreamRemote,
			UpstreamBranch: opts.upstreamBranch,
			SetUpstream:    opts.setUpstream,
		})

		if err != nil {
			if !opts.force && strings.Contains(err.Error(), "Updates were rejected") {
				forcePushDisabled := self.c.UserConfig.Git.DisableForcePushing
				if forcePushDisabled {
					_ = self.c.ErrorMsg(self.c.Tr.UpdatesRejectedAndForcePushDisabled)
					return nil
				}
				_ = self.c.Ask(popup.AskOpts{
					Title:  self.c.Tr.ForcePush,
					Prompt: self.c.Tr.ForcePushPrompt,
					HandleConfirm: func() error {
						newOpts := opts
						newOpts.force = true

						return self.pushAux(newOpts)
					},
				})
				return nil
			}
			_ = self.c.Error(err)
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *SyncController) requestToForcePush(opts pushOpts) error {
	forcePushDisabled := self.c.UserConfig.Git.DisableForcePushing
	if forcePushDisabled {
		return self.c.ErrorMsg(self.c.Tr.ForcePushDisabled)
	}

	return self.c.Ask(popup.AskOpts{
		Title:  self.c.Tr.ForcePush,
		Prompt: self.c.Tr.ForcePushPrompt,
		HandleConfirm: func() error {
			return self.pushAux(opts)
		},
	})
}
