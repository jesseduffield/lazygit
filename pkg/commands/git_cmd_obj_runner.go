package commands

import (
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
	return self.innerRunner.RunWithOutput(cmdObj)
}

func (self *gitCmdObjRunner) RunWithOutputs(cmdObj oscommands.ICmdObj) (string, string, error) {
	return self.innerRunner.RunWithOutputs(cmdObj)
}

func (self *gitCmdObjRunner) RunAndProcessLines(cmdObj oscommands.ICmdObj, onLine func(line string) (bool, error)) error {
	return self.innerRunner.RunAndProcessLines(cmdObj, onLine)
}
