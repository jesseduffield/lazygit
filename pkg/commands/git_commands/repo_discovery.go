package git_commands

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// DiscoverRepos returns the absolute paths of the git repos below root, down to
// maxDepth directory levels, sorted by path. It does not look inside a repo it
// found, it skips hidden directories, and it does not follow symlinks, so a
// symlink loop can't trap it. It skips directories that it can't read.
func DiscoverRepos(root string, maxDepth int) []string {
	var repos []string

	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			if _, err := os.Lstat(filepath.Join(path, ".git")); err == nil {
				repos = append(repos, path)
				continue
			}
			if depth < maxDepth {
				walk(path, depth+1)
			}
		}
	}
	walk(root, 1)

	slices.Sort(repos)
	return repos
}
