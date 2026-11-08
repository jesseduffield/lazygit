package git_commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// IsTrackedByEnclosingRepo reports whether git tracks the file at path in the
// repository whose working tree contains it. If neither the file's directory
// nor any directory above it has a .git, there is no such repository, and we
// don't run git at all.
//
// We can't ask git about path directly, because the index may spell it
// differently. On a case-insensitive file system, a tracked .LAZYGIT.YML is the
// file .lazygit.yml, and on macOS, the index may have a directory name in a
// different Unicode normalization form than path. So we ask git for all tracked
// files that have path's base name in any letter case, and check whether one of
// them is the same file as path. This only works if the base name has the same
// spelling in path and in the index, so it must consist of ASCII characters.
func IsTrackedByEnclosingRepo(cmd oscommands.ICmdObjBuilder, path string) (bool, error) {
	repoDir, found := findDirWithGitEntry(filepath.Dir(path))
	if !found {
		return false, nil
	}

	baseName := filepath.Base(path)
	if strings.ContainsFunc(baseName, func(r rune) bool { return r > unicode.MaxASCII }) {
		return false, fmt.Errorf("can't tell whether git tracks %s, because its name isn't ASCII", path)
	}

	pathInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	cmdObj := cmd.New(
		NewGitCmd("ls-files").
			Arg("-z", "--", ":(icase,glob)**/"+escapeGlob(baseName)).
			Dir(repoDir).
			ToArgv(),
	).DontLog()
	// With GIT_LITERAL_PATHSPECS set, git would look for a file whose name
	// starts with ":(icase,glob)", and find nothing.
	stdout, _, err := ForOtherRepo(cmdObj).RemoveEnvVar("GIT_LITERAL_PATHSPECS").RunWithOutputs()
	if err != nil {
		return false, err
	}

	for _, trackedPath := range strings.Split(stdout, "\x00") {
		if trackedPath == "" {
			continue
		}
		trackedInfo, err := os.Lstat(filepath.Join(repoDir, trackedPath))
		if err == nil && os.SameFile(trackedInfo, pathInfo) {
			return true, nil
		}
	}

	return false, nil
}

// Returns the first directory from dir upwards that has a .git; this is the
// root of the working tree that dir is in.
func findDirWithGitEntry(dir string) (string, bool) {
	for {
		// If we can't tell whether there is a .git, we let git decide.
		if _, err := os.Lstat(filepath.Join(dir, ".git")); !errors.Is(err, os.ErrNotExist) {
			return dir, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func escapeGlob(name string) string {
	var result strings.Builder
	for _, r := range name {
		if strings.ContainsRune(`\*?[`, r) {
			result.WriteRune('\\')
		}
		result.WriteRune(r)
	}
	return result.String()
}
