package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ExpandHomeDir expands a leading "~" or "~/" in path to the user's
// home directory. Paths that do not begin with "~" are returned
// unchanged, including the empty string.
//
// We deliberately only support "~" for the current user (no
// "~user" syntax) because that is what the shells lazygit users
// expect to see lazygit honour, and the cross-platform story for
// other-user expansion is messy.
//
// If path begins with "~" but the home directory cannot be
// determined, the unexpanded path is returned along with the error.
func ExpandHomeDir(path string) (string, error) {
	if path == "" || path[0] != '~' {
		return path, nil
	}
	// Only "~" alone or "~/..." — not "~user/...".
	if path != "~" && !strings.HasPrefix(path, "~"+string(filepath.Separator)) && !strings.HasPrefix(path, "~/") {
		return path, errors.New("only the current user's home directory can be expanded with ~")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path, err
	}
	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, path[2:]), nil
}
