package oscommands

import (
	"os/exec"
)

// A command object is a general way to represent a command to be run on the
// command line. If you want to log the command you'll use .ToString() and
// if you want to run it you'll use .GetCmd()
type ICmdObj interface {
	GetCmd() *exec.Cmd
	ToString() string
	AddEnvVars(...string) ICmdObj
	GetEnvVars() []string
}

type CmdObj struct {
	cmdStr string
	cmd    *exec.Cmd
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
