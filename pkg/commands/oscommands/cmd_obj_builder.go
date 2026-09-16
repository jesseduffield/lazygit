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
	New(args []string) *CmdObj
	// NewShell takes a string like `git commit` and returns an executable shell command for it e.g. `sh -c 'git commit'`
	// shellFunctionsFile is an optional file path that will be sourced before executing the command. Callers should pass
	// the value of UserConfig.OS.ShellFunctionsFile.
	NewShell(commandStr string, shellFunctionsFile string) *CmdObj
	// Quote wraps a string in quotes with any necessary escaping applied. The reason for bundling this up with the other methods in this interface is that we basically always need to make use of this when creating new command objects.
	Quote(str string) string
}

type CmdObjBuilder struct {
	runner   ICmdObjRunner
	platform *Platform
}

// poor man's version of explicitly saying that struct X implements interface Y
var _ ICmdObjBuilder = &CmdObjBuilder{}

func (self *CmdObjBuilder) New(args []string) *CmdObj {
	cmdObj := self.NewWithEnviron(args, os.Environ())
	return cmdObj
}

// A command with explicit environment from env
func (self *CmdObjBuilder) NewWithEnviron(args []string, env []string) *CmdObj {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = env

	return &CmdObj{
		cmd:    cmd,
		runner: self.runner,
	}
}

func (self *CmdObjBuilder) NewShell(commandStr string, shellFunctionsFile string) *CmdObj {
	if len(shellFunctionsFile) > 0 {
		commandStr = fmt.Sprintf("%ssource %s\n%s", self.platform.PrefixForShellFunctionsFile, shellFunctionsFile, commandStr)
	}

	if self.platform.OS == "windows" {
		return self.newWindowsShell(commandStr)
	}

	quotedCommand := self.Quote(commandStr)
	cmdArgs := str.ToArgv(fmt.Sprintf("%s %s %s", self.platform.Shell, self.platform.ShellArg, quotedCommand))

	return self.New(cmdArgs)
}

// newWindowsShell wraps the command in `cmd.exe /s /c "<command>"`. The /s
// flag tells cmd to strip exactly the outermost pair of quotes and pass the
// rest through unchanged, which preserves any quoting the command itself
// contains (e.g. `"C:\Program Files\my-editor.exe" file.txt`). Without /s,
// cmd's default rules drop the wrong quotes once the command line contains
// more than two of them.
//
// We bypass Go's standard arg quoting via SysProcAttr.CmdLine: it follows the
// CommandLineToArgvW convention (`\"` for inner quotes), but cmd.exe doesn't.
func (self *CmdObjBuilder) newWindowsShell(commandStr string) *CmdObj {
	args := []string{self.platform.Shell, "/s", self.platform.ShellArg, commandStr}
	cmdObj := self.New(args)

	cmdLine := fmt.Sprintf(`%s /s %s "%s"`, self.platform.Shell, self.platform.ShellArg, commandStr)
	setRawCmdLine(cmdObj.GetCmd(), cmdLine)

	return cmdObj
}

func (self *CmdObjBuilder) CloneWithNewRunner(decorate func(ICmdObjRunner) ICmdObjRunner) *CmdObjBuilder {
	decoratedRunner := decorate(self.runner)

	return &CmdObjBuilder{
		runner:   decoratedRunner,
		platform: self.platform,
	}
}

func (self *CmdObjBuilder) Quote(message string) string {
	if self.platform.OS == "windows" {
		return quoteForWindows(message)
	}
	message = strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		`$`, `\$`,
		"`", "\\`",
	).Replace(message)
	return `"` + message + `"`
}

// quoteForWindows encodes a value using the standard Windows command-line
// convention (the algorithm behind syscall.EscapeArg, reimplemented here so
// it's available on all platforms). The result is always wrapped in double
// quotes so cmd.exe and CommandLineToArgvW treat it as a single argument
// regardless of what shell metacharacters it contains.
func quoteForWindows(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	slashes := 0
	for i := range len(s) {
		c := s[i]
		switch c {
		case '\\':
			slashes++
			b.WriteByte(c)
		case '"':
			for ; slashes > 0; slashes-- {
				b.WriteByte('\\')
			}
			b.WriteByte('\\')
			b.WriteByte(c)
		default:
			slashes = 0
			b.WriteByte(c)
		}
	}
	for ; slashes > 0; slashes-- {
		b.WriteByte('\\')
	}
	b.WriteByte('"')
	return b.String()
}
