//go:build !windows

package oscommands

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func GetPlatform() *Platform {
	shell := getUserShell()

	prefixForShellFunctionsFile := ""
	if strings.HasSuffix(shell, "bash") {
		prefixForShellFunctionsFile = "shopt -s expand_aliases\n"
	}

	return &Platform{
		OS:                          runtime.GOOS,
		Shell:                       shell,
		ShellArg:                    "-c",
		PrefixForShellFunctionsFile: prefixForShellFunctionsFile,
		OpenCommand:                 "open {{filename}}",
		OpenLinkCommand:             "open {{link}}",
	}
}

func getUserShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}

	return "bash"
}

func (c *OSCommand) UpdateWindowTitle() error {
	return nil
}

// setRawCmdLine is the non-Windows no-op counterpart of the Windows shim
// (see the comment there). NewShell's shell-building logic is portable, so
// this call is reached on every host; only the Windows build does anything.
func setRawCmdLine(cmd *exec.Cmd, cmdLine string) {}

func TerminateProcessGracefully(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}

	return cmd.Process.Signal(syscall.SIGTERM)
}
