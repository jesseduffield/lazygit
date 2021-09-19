//go:build !windows
// +build !windows

package secureexec

import (
	"os/exec"
)

func Command(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
