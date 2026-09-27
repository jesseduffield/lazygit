package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
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
			Keys:              opts.GetKeys(opts.Config.Universal.Push),
			Handler:           opts.Guards.NoPopupPanel(self.HandlePush),
			GetDisabledReason: self.getDisabledReasonForPushOrPull,
			Description:       self.c.Tr.Push,
			Tooltip:           self.c.Tr.PushTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Pull),
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
	branchesBelow := self.unpushedBranchesBelow(currentBranch)
	if len(branchesBelow) == 0 {
		return self.pushCurrentBranch(currentBranch)
	}

	branchName := map[string]string{"branchName": currentBranch.Name}
	return self.c.Menu(types.CreateMenuOptions{
		Title:  self.c.Tr.Push,
		Prompt: self.branchesBelowPrompt(currentBranch, branchesBelow),
		Items: []*types.MenuItem{
			{
				Label: utils.ResolvePlaceholderString(self.c.Tr.PushBranchAndBranchesBelow, branchName),
				OnPress: func() error {
					return self.pushWithBranchesBelow(currentBranch, branchesBelow)
				},
			},
			{
				Label: utils.ResolvePlaceholderString(self.c.Tr.PushOnlyCurrentBranch, branchName),
				OnPress: func() error {
					return self.pushCurrentBranch(currentBranch)
				},
			},
		},
	})
}

func (self *SyncController) pushCurrentBranch(currentBranch *models.Branch) error {
	return self.resolvePushOfCurrentBranch(currentBranch, func(opts pushOpts) error {
		// if we are behind our upstream branch we'll ask if the user wants to force push
		if currentBranch.IsBehindForPush() {
			return self.requestToForcePush(currentBranch, opts)
		}

		return self.pushAux(currentBranch, opts)
	})
}

// The branches stacked below the current one that have commits to push
func (self *SyncController) unpushedBranchesBelow(currentBranch *models.Branch) []*models.Branch {
	branchesBelow := helpers.BranchesBelowInStack(
		self.c.Model().Commits, self.c.Model().Branches, currentBranch, self.c.UserConfig().Git.MainBranches)

	return lo.Filter(branchesBelow, func(branch *models.Branch, _ int) bool {
		// Pushing a branch whose remote branch was deleted would recreate it.
		// A branch that is only behind its remote branch has nothing to push,
		// and force-pushing it would move the remote branch back to an older
		// commit.
		return branch.PushRemote != "" && !branch.UpstreamGone && branch.IsAheadForPush()
	})
}

func (self *SyncController) branchesBelowPrompt(currentBranch *models.Branch, branchesBelow []*models.Branch) string {
	intro := utils.ResolvePlaceholderString(
		self.c.Tr.BranchesBelowHaveCommitsToPush,
		map[string]string{"branchName": currentBranch.Name},
	)
	lines := lo.Map(branchesBelow, func(branch *models.Branch, _ int) string {
		divergence := "↑" + branch.AheadForPush
		if branch.IsBehindForPush() {
			divergence = "↓" + branch.BehindForPush + divergence
		}
		return fmt.Sprintf("%s %s", branch.Name, style.FgYellow.Sprint(divergence))
	})

	return intro + "\n\n  " + strings.Join(lines, "\n  ")
}

// Pushes the current branch and the branches stacked below it, after asking
// for confirmation if any of them has to be force-pushed
func (self *SyncController) pushWithBranchesBelow(currentBranch *models.Branch, branchesBelow []*models.Branch) error {
	if currentBranch.RemoteBranchStoredLocally() && currentBranch.PushRemote != "" {
		// We know where the current branch goes and whether it needs to be
		// forced, so it is pushed like the branches below it, in the same
		// command as those that go to the same remote
		branches := append([]*models.Branch{currentBranch}, branchesBelow...)
		return self.confirmForcePushIfNeeded(branches, func(forceWithLease bool) error {
			return self.pushBranchesAux(currentBranch, branchesBelow, forceWithLease)
		})
	}

	// The current branch has no upstream yet, or its remote branch isn't
	// stored locally. Push it the way it is pushed on its own, and the branches
	// below it after it.
	return self.resolvePushOfCurrentBranch(currentBranch, func(opts pushOpts) error {
		return self.confirmForcePushIfNeeded(branchesBelow, func(forceWithLease bool) error {
			opts.branchesBelow = branchesBelow
			opts.forceWithLeaseBelow = forceWithLease
			return self.pushAux(currentBranch, opts)
		})
	})
}

// Calls push right away if none of the branches has diverged from its remote
// branch, and after the user confirmed force-pushing otherwise
func (self *SyncController) confirmForcePushIfNeeded(branches []*models.Branch, push func(forceWithLease bool) error) error {
	diverged := lo.Filter(branches, func(branch *models.Branch, _ int) bool {
		return branch.IsBehindForPush()
	})
	if len(diverged) == 0 {
		return push(false)
	}

	if self.c.UserConfig().Git.DisableForcePushing {
		return errors.New(self.c.Tr.ForcePushBranchesDisabled)
	}

	self.c.Confirm(types.ConfirmOpts{
		Title: self.c.Tr.ForcePush,
		Prompt: utils.ResolvePlaceholderString(
			self.c.Tr.ForcePushBranchesPrompt,
			map[string]string{
				"branches":   "  " + strings.Join(lo.Map(diverged, func(branch *models.Branch, _ int) string { return branch.Name }), "\n  "),
				"cancelKey":  self.c.UserConfig().Keybinding.Universal.Return.String(),
				"confirmKey": self.c.UserConfig().Keybinding.Universal.Confirm.String(),
			},
		),
		HandleConfirm: func() error {
			return push(true)
		},
	})

	return nil
}

// Works out where the current branch is pushed to: to its upstream, to a
// branch of the same name if push.default is "current", or to an upstream the
// user enters in a prompt. Calls onResolved with the options for that push.
func (self *SyncController) resolvePushOfCurrentBranch(currentBranch *models.Branch, onResolved func(pushOpts) error) error {
	if currentBranch.IsTrackingRemote() {
		return onResolved(pushOpts{remoteBranchStoredLocally: currentBranch.RemoteBranchStoredLocally()})
	}

	if self.c.Git().Config.GetPushToCurrent() {
		return onResolved(pushOpts{setUpstream: true})
	}

	return self.c.Helpers().Upstream.PromptForUpstreamWithInitialContent(currentBranch, func(upstream string) error {
		upstreamRemote, upstreamBranch, err := self.c.Helpers().Upstream.ParseUpstream(upstream)
		if err != nil {
			return err
		}

		return onResolved(pushOpts{
			setUpstream:    true,
			upstreamRemote: upstreamRemote,
			upstreamBranch: upstreamBranch,
		})
	})
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

	return self.c.Helpers().MergeAndRebase.CheckMergeOrRebaseAndSelectHeadCommit(err)
}

type pushOpts struct {
	force          bool
	forceWithLease bool
	upstreamRemote string
	upstreamBranch string
	setUpstream    bool

	// If this is false, we can't tell ahead of time whether a force-push will
	// be necessary, so we start with a normal push and offer to force-push if
	// the server rejected. If this is true, we don't offer to force-push if the
	// server rejected, but rather ask the user to fetch.
	remoteBranchStoredLocally bool

	// Branches stacked below the current one, pushed after it with one command
	// per remote, with --force-with-lease if forceWithLeaseBelow is set. The
	// force options above apply to the current branch's own push only.
	branchesBelow       []*models.Branch
	forceWithLeaseBelow bool
}

func (self *SyncController) pushAux(currentBranch *models.Branch, opts pushOpts) error {
	return self.withPushingStatus(currentBranch, opts.branchesBelow, func(task gocui.Task) error {
		refspecs := []string{}
		if opts.upstreamBranch != "" {
			refspecs = append(refspecs, fmt.Sprintf("refs/heads/%s:%s", currentBranch.Name, opts.upstreamBranch))
		}
		err := self.c.Git().Sync.Push(
			task,
			git_commands.PushOpts{
				Force:          opts.force,
				ForceWithLease: opts.forceWithLease,
				SetUpstream:    opts.setUpstream,
				Remote:         opts.upstreamRemote,
				Refspecs:       refspecs,
			})
		if err != nil {
			if !opts.force && !opts.forceWithLease && strings.Contains(err.Error(), "Updates were rejected") {
				if opts.remoteBranchStoredLocally {
					return errors.New(self.c.Tr.UpdatesRejected)
				}

				forcePushDisabled := self.c.UserConfig().Git.DisableForcePushing
				if forcePushDisabled {
					return errors.New(self.c.Tr.UpdatesRejectedAndForcePushDisabled)
				}
				self.c.Confirm(types.ConfirmOpts{
					Title:  self.c.Tr.ForcePush,
					Prompt: self.forcePushPrompt(),
					HandleConfirm: func() error {
						newOpts := opts
						newOpts.force = true

						return self.pushAux(currentBranch, newOpts)
					},
				})
				return nil
			}
			return err
		}

		err = self.pushBranches(task, opts.branchesBelow, opts.forceWithLeaseBelow)
		self.c.RefreshFromWorker(types.RefreshOptions{})
		return err
	})
}

// Pushes the current branch along with the branches stacked below it, all of
// them with explicit refspecs
func (self *SyncController) pushBranchesAux(currentBranch *models.Branch, branchesBelow []*models.Branch, forceWithLease bool) error {
	return self.withPushingStatus(currentBranch, branchesBelow, func(task gocui.Task) error {
		branches := append([]*models.Branch{currentBranch}, branchesBelow...)
		err := self.pushBranches(task, branches, forceWithLease)
		self.c.RefreshFromWorker(types.RefreshOptions{})
		return err
	})
}

// Runs f as a push of the current branch, showing it and the other branches
// as being pushed while it runs
func (self *SyncController) withPushingStatus(currentBranch *models.Branch, otherBranches []*models.Branch, f func(gocui.Task) error) error {
	return self.c.WithInlineStatus(currentBranch, types.ItemOperationPushing, context.LOCAL_BRANCHES_CONTEXT_KEY, func(task gocui.Task) error {
		for _, branch := range otherBranches {
			self.c.State().SetItemOperation(branch, types.ItemOperationPushing)
		}
		defer func() {
			for _, branch := range otherBranches {
				self.c.State().ClearItemOperation(branch)
			}
		}()

		self.c.LogAction(self.c.Tr.Actions.Push)
		return f(task)
	})
}

// Pushes the branches to their push destinations, one command per remote
func (self *SyncController) pushBranches(task gocui.Task, branches []*models.Branch, forceWithLease bool) error {
	remotes := lo.Uniq(lo.Map(branches, func(branch *models.Branch, _ int) string { return branch.PushRemote }))
	for _, remote := range remotes {
		refspecs := lo.FilterMap(branches, func(branch *models.Branch, _ int) (string, bool) {
			return fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch.Name, branch.PushBranch), branch.PushRemote == remote
		})
		err := self.c.Git().Sync.Push(task, git_commands.PushOpts{
			ForceWithLease: forceWithLease,
			Remote:         remote,
			Refspecs:       refspecs,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *SyncController) requestToForcePush(currentBranch *models.Branch, opts pushOpts) error {
	forcePushDisabled := self.c.UserConfig().Git.DisableForcePushing
	if forcePushDisabled {
		return errors.New(self.c.Tr.ForcePushDisabled)
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.ForcePush,
		Prompt: self.forcePushPrompt(),
		HandleConfirm: func() error {
			opts.forceWithLease = true
			return self.pushAux(currentBranch, opts)
		},
	})

	return nil
}

func (self *SyncController) forcePushPrompt() string {
	return utils.ResolvePlaceholderString(
		self.c.Tr.ForcePushPrompt,
		map[string]string{
			"cancelKey":  self.c.UserConfig().Keybinding.Universal.Return.String(),
			"confirmKey": self.c.UserConfig().Keybinding.Universal.Confirm.String(),
		},
	)
}
