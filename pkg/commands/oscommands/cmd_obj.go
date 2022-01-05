package oscommands

import (
	"os/exec"
)

// A command object is a general way to represent a command to be run on the
// command line.
type ICmdObj interface {
	GetCmd() *exec.Cmd
	// outputs string representation of command. Note that if the command was built
	// using NewFromArgs, the output won't be quite the same as what you would type
	// into a terminal e.g. 'sh -c git commit' as opposed to 'sh -c "git commit"'
	ToString() string

	AddEnvVars(...string) ICmdObj
	GetEnvVars() []string

	// runs the command and returns an error if any
	Run() error
	// runs the command and returns the output as a string, and an error if any
	RunWithOutput() (string, error)
	// runs the command and runs a callback function on each line of the output. If the callback returns true for the boolean value, we kill the process and return.
	RunAndProcessLines(onLine func(line string) (bool, error)) error

	// Be calling DontLog(), we're saying that once we call Run(), we don't want to
	// log the command in the UI (it'll still be logged in the log file). The general rule
	// is that if a command doesn't change the git state (e.g. read commands like `git diff`)
	// then we don't want to log it. If we are changing something (e.g. `git add .`) then
	// we do. The only exception is if we're running a command in the background periodically
	// like `git fetch`, which technically does mutate stuff but isn't something we need
	// to notify the user about.
	DontLog() ICmdObj

	// This returns false if DontLog() was called
	ShouldLog() bool
}

type CmdObj struct {
	cmdStr string
	cmd    *exec.Cmd

	runner ICmdObjRunner

	// if set to true, we don't want to log the command to the user.
	dontLog bool
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

func (self *CmdObj) DontLog() ICmdObj {
	self.dontLog = true
	return self
}

func (self *CmdObj) ShouldLog() bool {
	return !self.dontLog
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
