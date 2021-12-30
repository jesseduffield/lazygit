package commands

import (
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/sirupsen/logrus"
)

// here we're wrapping the default command runner in some git-specific stuff e.g. retry logic if we get an error due to the presence of .git/index.lock

type gitCmdObjRunner struct {
	log         *logrus.Entry
	innerRunner oscommands.ICmdObjRunner
}

func (self *gitCmdObjRunner) Run(cmdObj oscommands.ICmdObj) error {
	_, err := self.RunWithOutput(cmdObj)
	return err
}

func (self *gitCmdObjRunner) RunWithOutput(cmdObj oscommands.ICmdObj) (string, error) {
	// TODO: have this retry logic in other places we run the command
	waitTime := 50 * time.Millisecond
	retryCount := 5
	attempt := 0

	for {
		output, err := self.innerRunner.RunWithOutput(cmdObj)
		if err != nil {
			// if we have an error based on the index lock, we should wait a bit and then retry
			if strings.Contains(output, ".git/index.lock") {
				self.log.Error(output)
				self.log.Info("index.lock prevented command from running. Retrying command after a small wait")
				attempt++
				time.Sleep(waitTime)
				if attempt < retryCount {
					continue
				}
			}
		}
		return output, err
	}
}

func (self *gitCmdObjRunner) RunAndProcessLines(cmdObj oscommands.ICmdObj, onLine func(line string) (bool, error)) error {
	return self.innerRunner.RunAndProcessLines(cmdObj, onLine)
}
