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

	// runs the command and returns an error if any
	Run() error
	// runs the command and returns the output as a string, and an error if any
	RunWithOutput() (string, error)
	// runs the command and runs a callback function on each line of the output. If the callback returns true for the boolean value, we kill the process and return.
	RunAndProcessLines(onLine func(line string) (bool, error)) error

	// logs command
	Log() ICmdObj
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

func (self *CmdObj) Log() ICmdObj {
	self.logCommand(self)

	return self
}

func (self *CmdObj) Run() error {
	return self.runner.Run(self)
}

func (self *CmdObj) RunWithOutput() (string, error) {
	return self.runner.RunWithOutput(self)
}

func (self *CmdObj) RunAndProcessLines(onLine func(line string) (bool, error)) error {
	return self.runner.RunAndProcessLines(self, onLine)
}
