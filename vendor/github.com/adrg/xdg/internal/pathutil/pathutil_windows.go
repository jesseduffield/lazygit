package pathutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// UserHomeDir returns the home directory of the current user.
func UserHomeDir() string {
	return KnownFolder(windows.FOLDERID_Profile, []string{"USERPROFILE"}, nil)
}

// Exists returns true if the specified path exists.
func Exists(path string) bool {
	fi, err := os.Lstat(path)
	if fi != nil && fi.Mode()&os.ModeSymlink != 0 {
		_, err = filepath.EvalSymlinks(path)
	}

	return err == nil || errors.Is(err, fs.ErrExist)
}

// ExpandHome substitutes `%USERPROFILE%` at the start of the specified `path`.
func ExpandHome(path string) string {
	home := UserHomeDir()
	if path == "" || home == "" {
		return path
	}
	if strings.HasPrefix(path, `%USERPROFILE%`) {
		return filepath.Join(home, path[13:])
	}

	return path
}

// KnownFolder returns the location of the folder with the specified ID.
// If that fails, the folder location is determined by reading the provided
// environment variables (the first non-empty read value is returned).
// If that fails as well, the first non-empty fallback is returned.
// If all of the above fails, the function returns an empty string.
func KnownFolder(id *windows.KNOWNFOLDERID, envVars []string, fallbacks []string) string {
	if id != nil {
		flags := []uint32{windows.KF_FLAG_DEFAULT, windows.KF_FLAG_DEFAULT_PATH}
		for _, flag := range flags {
			p, _ := windows.KnownFolderPath(id, flag|windows.KF_FLAG_DONT_VERIFY)
			if p != "" {
				return p
			}
		}
	}

	for _, envVar := range envVars {
		p := os.Getenv(envVar)
		if p != "" {
			return p
		}
	}

	for _, fallback := range fallbacks {
		if fallback != "" {
			return fallback
		}
	}

	return ""
}
