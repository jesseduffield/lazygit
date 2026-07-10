package commands

import (
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/sirupsen/logrus"
)

// here we're wrapping the default command runner in some git-specific stuff e.g. retry logic if we get an error due to the presence of .git/index.lock

const (
	WaitTime   = 50 * time.Millisecond
	RetryCount = 5
)

type gitCmdObjRunner struct {
	log         *logrus.Entry
	innerRunner oscommands.ICmdObjRunner
}

// isRetryableError returns true if a failed command hit a transient
// lock-related condition that may succeed on retry. The lock message can reach
// us either in the command's captured output or, for streamed commands whose
// output we don't capture, only in the returned error, so we check both.
//
// We match the bare "index.lock" fragment rather than a fuller path or message
// so we catch the lock wherever git puts it: the main .git dir, a linked
// worktree's git dir (.git/worktrees/<name>/index.lock), or a submodule's git
// dir.
func isRetryableError(output string, err error) bool {
	text := output
	if err != nil {
		text += "\n" + err.Error()
	}
	return strings.Contains(text, "index.lock") ||
		strings.Contains(text, "cannot lock ref")
}

func (self *gitCmdObjRunner) Run(cmdObj *oscommands.CmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *gitCmdObjRunner) RunWithOutput(cmdObj *oscommands.CmdObj) (string, error) {
	return self.retryOnLockError(func() (string, error) {
		return self.innerRunner.RunWithOutput(cmdObj.Clone())
	})
}

func (self *gitCmdObjRunner) RunWithOutputs(cmdObj *oscommands.CmdObj) (string, string, error) {
	var stdout, stderr string
	_, err := self.retryOnLockError(func() (string, error) {
		var runErr error
		stdout, stderr, runErr = self.innerRunner.RunWithOutputs(cmdObj.Clone())
		return stdout + stderr, runErr
	})
	return stdout, stderr, err
}

// retryOnLockError runs the given function, retrying if it fails with a
// transient lock error (see isRetryableError). The string returned by run is
// the command output we inspect to classify the failure. We clone the command
// for each attempt (inside run) because an *exec.Cmd can only be run once.
func (self *gitCmdObjRunner) retryOnLockError(run func() (string, error)) (string, error) {
	var output string
	var err error
	for range RetryCount {
		output, err = run()

		if err == nil || !isRetryableError(output, err) {
			return output, err
		}

		// if we have an error based on a lock, we should wait a bit and then retry
		self.log.Warn("lock error prevented command from running. Retrying command after a small wait")
		time.Sleep(WaitTime)
	}

	return output, err
}

// Retry logic not implemented here, but these commands typically don't need to obtain a lock.
func (self *gitCmdObjRunner) RunAndProcessLines(cmdObj *oscommands.CmdObj, onLine func(line string) (bool, error)) error {
	return self.innerRunner.RunAndProcessLines(cmdObj, onLine)
}
