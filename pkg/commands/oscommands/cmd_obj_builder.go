package oscommands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mgutz/str"
)

type ICmdObjBuilder interface {
	// NewFromArgs takes a slice of strings like []string{"git", "commit"} and returns a new command object.
	New(args []string) ICmdObj
	// NewShell takes a string like `git commit` and returns an executable shell command for it e.g. `sh -c 'git commit'`
	NewShell(commandStr string) ICmdObj
	// Like NewShell, but uses the user's shell rather than "bash", and passes -i to it
	NewInteractiveShell(commandStr string) ICmdObj
	// Quote wraps a string in quotes with any necessary escaping applied. The reason for bundling this up with the other methods in this interface is that we basically always need to make use of this when creating new command objects.
	Quote(str string) string
}

type CmdObjBuilder struct {
	runner   ICmdObjRunner
	platform *Platform
}

// poor man's version of explicitly saying that struct X implements interface Y
var _ ICmdObjBuilder = &CmdObjBuilder{}

func (self *CmdObjBuilder) New(args []string) ICmdObj {
	cmdObj := self.NewWithEnviron(args, os.Environ())
	return cmdObj
}

// A command with explicit environment from env
func (self *CmdObjBuilder) NewWithEnviron(args []string, env []string) ICmdObj {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = env

	return &CmdObj{
		cmd:    cmd,
		runner: self.runner,
	}
}

func (self *CmdObjBuilder) NewShell(commandStr string) ICmdObj {
	quotedCommand := self.quotedCommandString(commandStr)
	cmdArgs := str.ToArgv(fmt.Sprintf("%s %s %s", self.platform.Shell, self.platform.ShellArg, quotedCommand))

	return self.New(cmdArgs)
}

func (self *CmdObjBuilder) NewInteractiveShell(commandStr string) ICmdObj {
	quotedCommand := self.quotedCommandString(commandStr + self.platform.InteractiveShellExit)
	cmdArgs := str.ToArgv(fmt.Sprintf("%s %s %s %s", self.platform.InteractiveShell, self.platform.InteractiveShellArg, self.platform.ShellArg, quotedCommand))

	return self.New(cmdArgs)
}

func (self *CmdObjBuilder) quotedCommandString(commandStr string) string {
	// Windows does not seem to like quotes around the command
	if self.platform.OS == "windows" {
		return strings.NewReplacer(
			"^", "^^",
			"&", "^&",
			"|", "^|",
			"<", "^<",
			">", "^>",
			"%", "^%",
		).Replace(commandStr)
	}

	return self.Quote(commandStr)
}

func (self *CmdObjBuilder) CloneWithNewRunner(decorate func(ICmdObjRunner) ICmdObjRunner) *CmdObjBuilder {
	decoratedRunner := decorate(self.runner)

	return &CmdObjBuilder{
		runner:   decoratedRunner,
		platform: self.platform,
	}
}

const CHARS_REQUIRING_QUOTES = "\"\\$` "

// If you update this method, be sure to update CHARS_REQUIRING_QUOTES
func (self *CmdObjBuilder) Quote(message string) string {
	var quote string
	if self.platform.OS == "windows" {
		quote = `\"`
		message = strings.NewReplacer(
			`"`, `"'"'"`,
			`\"`, `\\"`,
		).Replace(message)
	} else {
		quote = `"`
		message = strings.NewReplacer(
			`\`, `\\`,
			`"`, `\"`,
			`$`, `\$`,
			"`", "\\`",
		).Replace(message)
	}
	return quote + message + quote
}
