package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type GpgHelper struct {
	c *HelperCommon
}

func NewGpgHelper(c *HelperCommon) *GpgHelper {
	return &GpgHelper{
		c: c,
	}
}

func (self *GpgHelper) WithGpgHandling(
	cmdObj *oscommands.CmdObj,
	configKey git_commands.GpgConfigKey,
	waitingStatus string,
	onSuccess func() error,
	refreshScope []types.RefreshableView,
) error {
	refreshOptions := types.RefreshOptions{Mode: types.ASYNC, Scope: refreshScope}
	return self.withGpgHandling(
		cmdObj, configKey, waitingStatus, onSuccess, refreshOptions, refreshOptions)
}

// WithGpgHandlingAndSelectHeadCommit is like WithGpgHandling, but on success it
// selects the new HEAD commit rather than restoring the previous selection. For
// committing, where the commit we just created is the one we want selected.
func (self *GpgHelper) WithGpgHandlingAndSelectHeadCommit(
	cmdObj *oscommands.CmdObj,
	configKey git_commands.GpgConfigKey,
	waitingStatus string,
	onSuccess func() error,
) error {
	failureRefreshOptions := types.RefreshOptions{Mode: types.ASYNC}
	successRefreshOptions := types.RefreshOptions{Mode: types.ASYNC, CommitSelection: types.SelectHeadCommit}
	return self.withGpgHandling(
		cmdObj, configKey, waitingStatus, onSuccess, failureRefreshOptions, successRefreshOptions)
}

// Currently there is a bug where if we switch to a subprocess from within
// WithWaitingStatus we get stuck there and can't return to lazygit. We could
// fix this bug, or just stop running subprocesses from within there, given that
// we don't need to see a loading status if we're in a subprocess.
func (self *GpgHelper) withGpgHandling(
	cmdObj *oscommands.CmdObj,
	configKey git_commands.GpgConfigKey,
	waitingStatus string,
	onSuccess func() error,
	failureRefreshOptions types.RefreshOptions,
	successRefreshOptions types.RefreshOptions,
) error {
	needsGpg := self.c.Git().Config.IsGpgSignEnabled(configKey)
	if needsGpg && self.c.Git().Config.CanUseGpgLoopback() {
		return self.runWithLoopbackPinentry(
			cmdObj, waitingStatus, onSuccess, failureRefreshOptions, successRefreshOptions)
	}

	useSubprocess := needsGpg && self.c.Git().Config.NeedsGpgSubprocess(configKey)
	if useSubprocess {
		success, err := self.c.RunSubprocess(cmdObj)
		if success && onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}
		if success {
			self.c.Refresh(successRefreshOptions)
		} else {
			self.c.Refresh(failureRefreshOptions)
		}

		return err
	}

	return self.runAndStream(
		cmdObj, waitingStatus, onSuccess, failureRefreshOptions, successRefreshOptions)
}

// runWithLoopbackPinentry runs cmdObj with gpg's `--pinentry-mode=loopback`
// enabled (via AddGpgLoopbackEnvVars), so that a passphrase prompt appears as
// plain text on gpg's own stdio. We detect that prompt using the same
// mechanism as e.g. SSH credential prompts, and answer it from our own popup,
// so signing never has to hand off the terminal.
func (self *GpgHelper) runWithLoopbackPinentry(
	cmdObj *oscommands.CmdObj,
	waitingStatus string,
	onSuccess func() error,
	failureRefreshOptions types.RefreshOptions,
	successRefreshOptions types.RefreshOptions,
) error {
	self.c.Git().Config.AddGpgLoopbackEnvVars(cmdObj)

	return self.c.WithWaitingStatus(waitingStatus, func(task gocui.Task) error {
		if err := cmdObj.PromptOnCredentialRequest(task).Run(); err != nil {
			self.c.RefreshFromWorker(failureRefreshOptions)
			return fmt.Errorf(
				self.c.Tr.GitCommandFailed, self.c.UserConfig().Keybinding.Universal.ExtrasMenu,
			)
		}

		if onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}

		self.c.RefreshFromWorker(successRefreshOptions)
		return nil
	})
}

func (self *GpgHelper) runAndStream(
	cmdObj *oscommands.CmdObj,
	waitingStatus string,
	onSuccess func() error,
	failureRefreshOptions types.RefreshOptions,
	successRefreshOptions types.RefreshOptions,
) error {
	return self.c.WithWaitingStatus(waitingStatus, func(gocui.Task) error {
		if err := cmdObj.StreamOutput().Run(); err != nil {
			self.c.RefreshFromWorker(failureRefreshOptions)
			return fmt.Errorf(
				self.c.Tr.GitCommandFailed, self.c.UserConfig().Keybinding.Universal.ExtrasMenu,
			)
		}

		if onSuccess != nil {
			if err := onSuccess(); err != nil {
				return err
			}
		}

		self.c.RefreshFromWorker(successRefreshOptions)
		return nil
	})
}
