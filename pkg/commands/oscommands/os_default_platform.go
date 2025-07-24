//go:build !windows
// +build !windows

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

func Kill(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		// You can't kill a person with no body
		return nil
	}

	if cmd.ProcessState != nil {
		// The process has already exited
		return nil
	}

	// Terminate the process gracefully, and hope that it handles the signal properly. We use it
	// only for git, and we know that git is well-behaved in this regard. If we were to use this for
	// other commands, we may need to implement a more robust solution, e.g. waiting for a while and
	// then killing the process forcefully if it didn't terminate.
	return cmd.Process.Signal(syscall.SIGTERM)
}
