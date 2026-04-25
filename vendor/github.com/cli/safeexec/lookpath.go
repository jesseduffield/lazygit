//go:build !windows && go1.19
// +build !windows,go1.19

package safeexec

import (
	"errors"
	"os/exec"
)

func LookPath(file string) (string, error) {
	path, err := exec.LookPath(file)
	if errors.Is(err, exec.ErrDot) {
		return path, nil
	}
	return path, err
}
