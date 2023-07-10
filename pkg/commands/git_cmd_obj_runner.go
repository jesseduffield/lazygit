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

func (self *gitCmdObjRunner) Run(cmdObj oscommands.ICmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *gitCmdObjRunner) RunWithOutput(cmdObj oscommands.ICmdObj) (string, error) {
	var output string
	var err error
	for i := 0; i < RetryCount; i++ {
		newCmdObj := cmdObj.Clone()
		output, err = self.innerRunner.RunWithOutput(newCmdObj)

		if err == nil || !strings.Contains(output, ".git/index.lock") {
			return output, err
		}

		// if we have an error based on the index lock, we should wait a bit and then retry
		self.log.Warn("index.lock prevented command from running. Retrying command after a small wait")
		time.Sleep(WaitTime)
	}

	return output, err
}

func (self *gitCmdObjRunner) RunWithOutputs(cmdObj oscommands.ICmdObj) (string, string, error) {
	var stdout, stderr string
	var err error
	for i := 0; i < RetryCount; i++ {
		newCmdObj := cmdObj.Clone()
		stdout, stderr, err = self.innerRunner.RunWithOutputs(newCmdObj)

		if err == nil || !strings.Contains(stdout+stderr, ".git/index.lock") {
			return stdout, stderr, err
		}

		// if we have an error based on the index lock, we should wait a bit and then retry
		self.log.Warn("index.lock prevented command from running. Retrying command after a small wait")
		time.Sleep(WaitTime)
	}

	return stdout, stderr, err
}

// Retry logic not implemented here, but these commands typically don't need to obtain a lock.
func (self *gitCmdObjRunner) RunAndProcessLines(cmdObj oscommands.ICmdObj, onLine func(line string) (bool, error)) error {
	return self.innerRunner.RunAndProcessLines(cmdObj, onLine)
}
