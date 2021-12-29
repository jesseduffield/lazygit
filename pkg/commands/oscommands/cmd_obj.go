package oscommands

import (
	"os/exec"
)

// A command object is a general way to represent a command to be run on the
// command line.
type ICmdObj interface {
	GetCmd() *exec.Cmd
	ToString() string
	AddEnvVars(...string) ICmdObj
	GetEnvVars() []string

	Run() error
	RunWithOutput() (string, error)
	RunLineOutputCmd(onLine func(line string) (bool, error)) error

	// logs command
	Log()
}

type CmdObj struct {
	cmdStr string
	cmd    *exec.Cmd

	runner     ICmdObjRunner
	logCommand func(ICmdObj)
}

func (self *CmdObj) GetCmd() *exec.Cmd {
	return self.cmd
}

func (self *CmdObj) ToString() string {
	return self.cmdStr
}

func (self *CmdObj) AddEnvVars(vars ...string) ICmdObj {
	self.cmd.Env = append(self.cmd.Env, vars...)

	return self
}

func (self *CmdObj) GetEnvVars() []string {
	return self.cmd.Env
}

func (self *CmdObj) Log() {
	self.logCommand(self)
}

func (self *CmdObj) Run() error {
	return self.runner.Run(self)
}

func (self *CmdObj) RunWithOutput() (string, error) {
	return self.runner.RunWithOutput(self)
}

func (self *CmdObj) RunLineOutputCmd(onLine func(line string) (bool, error)) error {
	return self.runner.RunLineOutputCmd(self, onLine)
}
