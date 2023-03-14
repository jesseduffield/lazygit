package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SyncController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &SyncController{}

func NewSyncController(
	common *controllerCommon,
) *SyncController {
	return &SyncController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *SyncController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Push),
			Handler:     opts.Guards.NoPopupPanel(self.HandlePush),
			Description: self.c.Tr.LcPush,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Pull),
			Handler:     opts.Guards.NoPopupPanel(self.HandlePull),
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
		currentBranch := self.helpers.Refs.GetCheckedOutRef()
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
		opts := pushOpts{}
		if currentBranch.HasCommitsToPull() {
			return self.requestToForcePush(opts)
		} else {
			return self.pushAux(opts)
		}
	} else {
		if self.git.Config.GetPushToCurrent() {
			return self.pushAux(pushOpts{setUpstream: true})
		} else {
			return self.helpers.Upstream.PromptForUpstreamWithInitialContent(currentBranch, func(upstream string) error {
				upstreamRemote, upstreamBranch, err := self.helpers.Upstream.ParseUpstream(upstream)
				if err != nil {
					return self.c.Error(err)
				}

				return self.pushAux(pushOpts{
					setUpstream:    true,
					upstreamRemote: upstreamRemote,
					upstreamBranch: upstreamBranch,
				})
			})
		}
	}
}

func (self *SyncController) pull(currentBranch *models.Branch) error {
	action := self.c.Tr.Actions.Pull

	// if we have no upstream branch we need to set that first
	if !currentBranch.IsTrackingRemote() {
		return self.helpers.Upstream.PromptForUpstreamWithInitialContent(currentBranch, func(upstream string) error {
			if err := self.setCurrentBranchUpstream(upstream); err != nil {
				return self.c.Error(err)
			}

			return self.PullAux(PullFilesOptions{Action: action})
		})
	}

	return self.PullAux(PullFilesOptions{Action: action})
}

func (self *SyncController) setCurrentBranchUpstream(upstream string) error {
	upstreamRemote, upstreamBranch, err := self.helpers.Upstream.ParseUpstream(upstream)
	if err != nil {
		return err
	}

	if err := self.git.Branch.SetCurrentBranchUpstream(upstreamRemote, upstreamBranch); err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf(
				"upstream branch %s/%s not found.\nIf you expect it to exist, you should fetch (with 'f').\nOtherwise, you should push (with 'shift+P')",
				upstreamRemote, upstreamBranch,
			)
		}
		return err
	}
	return nil
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

	return self.helpers.MergeAndRebase.CheckMergeOrRebase(err)
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
				_ = self.c.Confirm(types.ConfirmOpts{
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

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.ForcePush,
		Prompt: self.c.Tr.ForcePushPrompt,
		HandleConfirm: func() error {
			opts.force = true
			return self.pushAux(opts)
		},
	})
}
