package oscommands

import (
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/mgutz/str"
)

type CredentialsStrategy int

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

func (self *CmdObj) AddEnvVars(vars ...string) {
	self.cmd.Env = append(self.cmd.Env, vars...)
}

func NewCmdObjFromStr(cmdStr string) *CmdObj {
	args := str.ToArgv(cmdStr)
	cmd := secureexec.Command(args[0], args[1:]...)

	return &CmdObj{
		cmdStr: cmdStr,
		cmd:    cmd,
	}
}

func NewCmdObjFromArgs(args []string) *CmdObj {
	cmd := secureexec.Command(args[0], args[1:]...)

	return &CmdObj{
		cmdStr: strings.Join(args, " "),
		cmd:    cmd,
	}
}

func NewCmdObj(cmd *exec.Cmd) *CmdObj {
	return &CmdObj{
		cmdStr: strings.Join(cmd.Args, " "),
		cmd:    cmd,
	}
}
