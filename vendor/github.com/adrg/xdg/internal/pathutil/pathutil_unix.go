//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || nacl || linux || netbsd || openbsd || solaris

package pathutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// UserHomeDir returns the home directory of the current user.
func UserHomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	return "/"
}

// Exists returns true if the specified path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, fs.ErrExist)
}

// ExpandHome substitutes `~` and `$HOME` at the start of the specified `path`.
func ExpandHome(path string) string {
	home := UserHomeDir()
	if path == "" || home == "" {
		return path
	}
	if path[0] == '~' {
		return filepath.Join(home, path[1:])
	}
	if strings.HasPrefix(path, "$HOME") {
		return filepath.Join(home, path[5:])
	}

	return path
}
