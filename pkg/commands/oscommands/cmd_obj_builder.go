package oscommands

import (
	"fmt"
	"os"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/mgutz/str"
)

type ICmdObjBuilder interface {
	// New returns a new command object based on the string provided
	New(cmdStr string) ICmdObj
	// NewShell takes a string like `git commit` and returns an executable shell command for it e.g. `sh -c 'git commit'`
	NewShell(commandStr string) ICmdObj
	// NewFromArgs takes a slice of strings like []string{"git", "commit"} and returns a new command object. This can be useful when you don't want to worry about whitespace and quoting and stuff.
	NewFromArgs(args []string) ICmdObj
	// Quote wraps a string in quotes with any necessary escaping applied. The reason for bundling this up with the other methods in this interface is that we basically always need to make use of this when creating new command objects.
	Quote(str string) string
}

type CmdObjBuilder struct {
	runner   ICmdObjRunner
	platform *Platform
}

// poor man's version of explicitly saying that struct X implements interface Y
var _ ICmdObjBuilder = &CmdObjBuilder{}

func (self *CmdObjBuilder) New(cmdStr string) ICmdObj {
	args := str.ToArgv(cmdStr)
	cmd := secureexec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()

	return &CmdObj{
		cmdStr: cmdStr,
		cmd:    cmd,
		runner: self.runner,
	}
}

func (self *CmdObjBuilder) NewFromArgs(args []string) ICmdObj {
	cmd := secureexec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()

	return &CmdObj{
		cmdStr: strings.Join(args, " "),
		cmd:    cmd,
		runner: self.runner,
	}
}

func (self *CmdObjBuilder) NewShell(commandStr string) ICmdObj {
	var quotedCommand string
	// Windows does not seem to like quotes around the command
	if self.platform.OS == "windows" {
		quotedCommand = strings.NewReplacer(
			"^", "^^",
			"&", "^&",
			"|", "^|",
			"<", "^<",
			">", "^>",
			"%", "^%",
		).Replace(commandStr)
	} else {
		quotedCommand = self.Quote(commandStr)
	}

	shellCommand := fmt.Sprintf("%s %s %s", self.platform.Shell, self.platform.ShellArg, quotedCommand)
	return self.New(shellCommand)
}

func (self *CmdObjBuilder) CloneWithNewRunner(decorate func(ICmdObjRunner) ICmdObjRunner) *CmdObjBuilder {
	decoratedRunner := decorate(self.runner)

	return &CmdObjBuilder{
		runner:   decoratedRunner,
		platform: self.platform,
	}
}

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
