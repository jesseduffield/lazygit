package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type SyncController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &SyncController{}

func NewSyncController(
	common *ControllerCommon,
) *SyncController {
	return &SyncController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *SyncController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Push),
			Handler:           opts.Guards.NoPopupPanel(self.HandlePush),
			GetDisabledReason: self.getDisabledReasonForPushOrPull,
			Description:       self.c.Tr.Push,
			Tooltip:           self.c.Tr.PushTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Pull),
			Handler:           opts.Guards.NoPopupPanel(self.HandlePull),
			GetDisabledReason: self.getDisabledReasonForPushOrPull,
			Description:       self.c.Tr.Pull,
			Tooltip:           self.c.Tr.PullTooltip,
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

func (self *SyncController) getDisabledReasonForPushOrPull() *types.DisabledReason {
	currentBranch := self.c.Helpers().Refs.GetCheckedOutRef()
	if currentBranch != nil {
		op := self.c.State().GetItemOperation(currentBranch)
		if op != types.ItemOperationNone {
			return &types.DisabledReason{Text: self.c.Tr.CantPullOrPushSameBranchTwice}
		}
	}

	return nil
}

func (self *SyncController) branchCheckedOut(f func(*models.Branch) error) func() error {
	return func() error {
		currentBranch := self.c.Helpers().Refs.GetCheckedOutRef()
		if currentBranch == nil {
			// need to wait for branches to refresh
			return nil
		}

		return f(currentBranch)
	}
}

func (self *SyncController) push(currentBranch *models.Branch) error {
	// if we are behind our upstream branch we'll ask if the user wants to force push
	if currentBranch.IsTrackingRemote() {
		opts := pushOpts{}
		if currentBranch.IsBehindForPush() {
			return self.requestToForcePush(currentBranch, opts)
		} else {
			return self.pushAux(currentBranch, opts)
		}
	} else {
		if self.c.Git().Config.GetPushToCurrent() {
			return self.pushAux(currentBranch, pushOpts{setUpstream: true})
		} else {
			return self.c.Helpers().Upstream.PromptForUpstreamWithInitialContent(currentBranch, func(upstream string) error {
				upstreamRemote, upstreamBranch, err := self.c.Helpers().Upstream.ParseUpstream(upstream)
				if err != nil {
					return err
				}

				return self.pushAux(currentBranch, pushOpts{
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
		return self.c.Helpers().Upstream.PromptForUpstreamWithInitialContent(currentBranch, func(upstream string) error {
			if err := self.setCurrentBranchUpstream(upstream); err != nil {
				return err
			}

			return self.PullAux(currentBranch, PullFilesOptions{Action: action})
		})
	}

	return self.PullAux(currentBranch, PullFilesOptions{Action: action})
}

func (self *SyncController) setCurrentBranchUpstream(upstream string) error {
	upstreamRemote, upstreamBranch, err := self.c.Helpers().Upstream.ParseUpstream(upstream)
	if err != nil {
		return err
	}

	if err := self.c.Git().Branch.SetCurrentBranchUpstream(upstreamRemote, upstreamBranch); err != nil {
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

func (self *SyncController) PullAux(currentBranch *models.Branch, opts PullFilesOptions) error {
	return self.c.WithInlineStatus(currentBranch, types.ItemOperationPulling, context.LOCAL_BRANCHES_CONTEXT_KEY, func(task gocui.Task) error {
		return self.pullWithLock(task, opts)
	})
}

func (self *SyncController) pullWithLock(task gocui.Task, opts PullFilesOptions) error {
	self.c.LogAction(opts.Action)

	err := self.c.Git().Sync.Pull(
		task,
		git_commands.PullOptions{
			RemoteName:      opts.UpstreamRemote,
			BranchName:      opts.UpstreamBranch,
			FastForwardOnly: opts.FastForwardOnly,
		},
	)

	return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
}

type pushOpts struct {
	force          bool
	upstreamRemote string
	upstreamBranch string
	setUpstream    bool
}

func (self *SyncController) pushAux(currentBranch *models.Branch, opts pushOpts) error {
	return self.c.WithInlineStatus(currentBranch, types.ItemOperationPushing, context.LOCAL_BRANCHES_CONTEXT_KEY, func(task gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.Push)
		err := self.c.Git().Sync.Push(
			task,
			git_commands.PushOpts{
				Force:          opts.force,
				UpstreamRemote: opts.upstreamRemote,
				UpstreamBranch: opts.upstreamBranch,
				SetUpstream:    opts.setUpstream,
			})
		if err != nil {
			if strings.Contains(err.Error(), "Updates were rejected") {
				return errors.New(self.c.Tr.UpdatesRejected)
			}
			return err
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *SyncController) requestToForcePush(currentBranch *models.Branch, opts pushOpts) error {
	forcePushDisabled := self.c.UserConfig.Git.DisableForcePushing
	if forcePushDisabled {
		return errors.New(self.c.Tr.ForcePushDisabled)
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.ForcePush,
		Prompt: self.forcePushPrompt(),
		HandleConfirm: func() error {
			opts.force = true
			return self.pushAux(currentBranch, opts)
		},
	})
}

func (self *SyncController) forcePushPrompt() string {
	return utils.ResolvePlaceholderString(
		self.c.Tr.ForcePushPrompt,
		map[string]string{
			"cancelKey":  self.c.UserConfig.Keybinding.Universal.Return,
			"confirmKey": self.c.UserConfig.Keybinding.Universal.Confirm,
		},
	)
}
