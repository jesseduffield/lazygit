//go:build !windows && !go1.19
// +build !windows,!go1.19

package safeexec

import "os/exec"

func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
