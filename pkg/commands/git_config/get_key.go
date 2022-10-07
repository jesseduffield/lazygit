package git_config

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
)

// including license from https://github.com/tcnksm/go-gitconfig because this file is an adaptation of that repo's code
// Copyright (c) 2014 tcnksm

// MIT License

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

func runGitConfigCmd(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = io.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", fmt.Errorf("the key is not found for %s", cmd.Args)
			}
		}
		return "", err
	}

	return strings.TrimRight(stdout.String(), "\000"), nil
}

func getGitConfigCmd(key string) *exec.Cmd {
	gitArgs := []string{"config", "--get", "--null", key}
	return secureexec.Command("git", gitArgs...)
}

func getGitConfigGeneralCmd(args string) *exec.Cmd {
	gitArgs := append([]string{"config"}, strings.Split(args, " ")...)
	return secureexec.Command("git", gitArgs...)
}
