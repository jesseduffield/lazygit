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

// isRetryableError returns true if the error output indicates a transient
// lock-related error that may succeed on retry
func isRetryableError(output string) bool {
	return strings.Contains(output, ".git/index.lock") ||
		strings.Contains(output, "cannot lock ref") ||
		strings.Contains(output, "cannot update ref")
}

func (self *gitCmdObjRunner) Run(cmdObj *oscommands.CmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *gitCmdObjRunner) RunWithOutput(cmdObj *oscommands.CmdObj) (string, error) {
	var output string
	var err error
	for range RetryCount {
		newCmdObj := cmdObj.Clone()
		output, err = self.innerRunner.RunWithOutput(newCmdObj)

		if err == nil || !isRetryableError(output) {
			return output, err
		}

		// if we have an error based on a lock, we should wait a bit and then retry
		self.log.Warn("lock error prevented command from running. Retrying command after a small wait")
		time.Sleep(WaitTime)
	}

	return output, err
}

func (self *gitCmdObjRunner) RunWithOutputs(cmdObj *oscommands.CmdObj) (string, string, error) {
	var stdout, stderr string
	var err error
	for range RetryCount {
		newCmdObj := cmdObj.Clone()
		stdout, stderr, err = self.innerRunner.RunWithOutputs(newCmdObj)

		if err == nil || !isRetryableError(stdout+stderr) {
			return stdout, stderr, err
		}

		// if we have an error based on a lock, we should wait a bit and then retry
		self.log.Warn("lock error prevented command from running. Retrying command after a small wait")
		time.Sleep(WaitTime)
	}

	return stdout, stderr, err
}

// Retry logic not implemented here, but these commands typically don't need to obtain a lock.
func (self *gitCmdObjRunner) RunAndProcessLines(cmdObj *oscommands.CmdObj, onLine func(line string) (bool, error)) error {
	return self.innerRunner.RunAndProcessLines(cmdObj, onLine)
}
