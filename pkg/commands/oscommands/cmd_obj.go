package oscommands

import (
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/mgutz/str"
)

type CredentialsStrategy int

const (
	IGNORE CredentialsStrategy = iota
	// for PROMPT and FAIL we expect that sometimes a credential may be requested.
	// In the case of PROMPT, we forward the request on to the user via the UI.
	// In the case of FAIL we just enter a newline so that the program fails.
	// We should only use PROMPT/FAIL when we actually expect a credential to be
	// requested because it's slightly less efficient to check for a password request
	PROMPT
	FAIL
)

type CmdObj struct {
	CmdStr string
	Cmd    *exec.Cmd
	CredentialsStrategy
}

func (self *CmdObj) ToCmd() *exec.Cmd {
	if self.Cmd == nil {
		self.Cmd = executableFromString(self.CmdStr)
	}

	return self.Cmd
}

func (self *CmdObj) ToString() string {
	if self.CmdStr != "" {
		return self.CmdStr
	}

	return strings.Join(self.Cmd.Args, " ")
}

func (self *CmdObj) AddEnvVars(vars ...string) {
	self.ToCmd()

	self.Cmd.Env = append(self.Cmd.Env, vars...)
}

func executableFromString(cmdStr string) *exec.Cmd {
	splitCmd := str.ToArgv(cmdStr)
	cmd := secureexec.Command(splitCmd[0], splitCmd[1:]...)
	return cmd
}
