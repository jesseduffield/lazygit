package git_commands

import (
	"fmt"
	ioFs "io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/samber/lo"
	"github.com/spf13/afero"
)

type RepoPaths struct {
	currentPath        string
	worktreePath       string
	worktreeGitDirPath string
	repoPath           string
	repoGitDirPath     string
	repoName           string
}

// Current working directory of the program. Currently, this will always
// be the same as WorktreePath(), but in future we may support running
// lazygit from inside a subdirectory of the worktree.
func (self *RepoPaths) CurrentPath() string {
	return self.currentPath
}

// Path to the current worktree. If we're in the main worktree, this will
// be the same as RepoPath()
func (self *RepoPaths) WorktreePath() string {
	return self.worktreePath
}

// Path of the worktree's git dir.
// If we're in the main worktree, this will be the .git dir under the RepoPath().
// If we're in a linked worktree, it will be the directory pointed at by the worktree's .git file
func (self *RepoPaths) WorktreeGitDirPath() string {
	return self.worktreeGitDirPath
}

// Path of the repo. If we're in a the main worktree, this will be the same as WorktreePath()
// If we're in a bare repo, it will be the parent folder of the bare repo
func (self *RepoPaths) RepoPath() string {
	return self.repoPath
}

// path of the git-dir for the repo.
// If this is a bare repo, it will be the location of the bare repo
// If this is a non-bare repo, it will be the location of the .git dir in
// the main worktree.
func (self *RepoPaths) RepoGitDirPath() string {
	return self.repoGitDirPath
}

// Name of the repo. Basename of the folder containing the repo.
func (self *RepoPaths) RepoName() string {
	return self.repoName
}

// Returns the repo paths for a typical repo
func MockRepoPaths(currentPath string) *RepoPaths {
	return &RepoPaths{
		currentPath:        currentPath,
		worktreePath:       currentPath,
		worktreeGitDirPath: path.Join(currentPath, ".git"),
		repoPath:           currentPath,
		repoGitDirPath:     path.Join(currentPath, ".git"),
		repoName:           "lazygit",
	}
}

func GetRepoPaths(
	fs afero.Fs,
	currentPath string,
) (*RepoPaths, error) {
	return getRepoPathsAux(afero.NewOsFs(), resolveSymlink, currentPath)
}

func getRepoPathsAux(
	fs afero.Fs,
	resolveSymlinkFn func(string) (string, error),
	currentPath string,
) (*RepoPaths, error) {
	worktreePath := currentPath
	repoGitDirPath, repoPath, err := getCurrentRepoGitDirPath(fs, resolveSymlinkFn, currentPath)
	if err != nil {
		return nil, errors.Errorf("failed to get repo git dir path: %v", err)
	}

	var worktreeGitDirPath string
	if env.GetWorkTreeEnv() != "" {
		// This env is set when you pass --work-tree to lazygit. In that case,
		// we're not dealing with a linked work-tree, we're dealing with a 'specified'
		// worktree (for lack of a better term). In this case, the worktree has no
		// .git file and it just contains a bunch of files: it has no idea it's
		// pointed to by a bare repo. As such it does not have its own git dir within
		// the bare repo's git dir. Instead, we just use the bare repo's git dir.
		worktreeGitDirPath = repoGitDirPath
	} else {
		var err error
		worktreeGitDirPath, err = getWorktreeGitDirPath(fs, currentPath)
		if err != nil {
			return nil, errors.Errorf("failed to get worktree git dir path: %v", err)
		}
	}

	repoName := path.Base(repoPath)

	return &RepoPaths{
		currentPath:        currentPath,
		worktreePath:       worktreePath,
		worktreeGitDirPath: worktreeGitDirPath,
		repoPath:           repoPath,
		repoGitDirPath:     repoGitDirPath,
		repoName:           repoName,
	}, nil
}

// Returns the path of the git-dir for the worktree. For linked worktrees, the worktree has
// a .git file that points to the git-dir (which itself lives in the git-dir
// of the repo)
func getWorktreeGitDirPath(fs afero.Fs, worktreePath string) (string, error) {
	// if .git is a file, we're in a linked worktree, otherwise we're in
	// the main worktree
	dotGitPath := path.Join(worktreePath, ".git")
	gitFileInfo, err := fs.Stat(dotGitPath)
	if err != nil {
		return "", err
	}

	if gitFileInfo.IsDir() {
		return dotGitPath, nil
	}

	return linkedWorktreeGitDirPath(fs, worktreePath)
}

func linkedWorktreeGitDirPath(fs afero.Fs, worktreePath string) (string, error) {
	dotGitPath := path.Join(worktreePath, ".git")
	gitFileContents, err := afero.ReadFile(fs, dotGitPath)
	if err != nil {
		return "", err
	}

	// The file will have `gitdir: /path/to/.git/worktrees/<worktree-name>`
	gitDirLine := lo.Filter(strings.Split(string(gitFileContents), "\n"), func(line string, _ int) bool {
		return strings.HasPrefix(line, "gitdir: ")
	})

	if len(gitDirLine) == 0 {
		return "", errors.New(fmt.Sprintf("%s is a file which suggests we are in a submodule or a worktree but the file's contents do not contain a gitdir pointing to the actual .git directory", dotGitPath))
	}

	gitDir := strings.TrimPrefix(gitDirLine[0], "gitdir: ")

	// For windows support
	gitDir = filepath.ToSlash(gitDir)

	return gitDir, nil
}

func getCurrentRepoGitDirPath(
	fs afero.Fs,
	resolveSymlinkFn func(string) (string, error),
	currentPath string,
) (string, string, error) {
	var unresolvedGitPath string
	if env.GetGitDirEnv() != "" {
		unresolvedGitPath = env.GetGitDirEnv()
	} else {
		unresolvedGitPath = path.Join(currentPath, ".git")
	}

	gitPath, err := resolveSymlinkFn(unresolvedGitPath)
	if err != nil {
		return "", "", err
	}

	// check if .git is a file or a directory
	gitFileInfo, err := fs.Stat(gitPath)
	if err != nil {
		return "", "", err
	}

	if gitFileInfo.IsDir() {
		// must be in the main worktree
		return gitPath, path.Dir(gitPath), nil
	}

	// either in a submodule, or worktree
	worktreeGitPath, err := linkedWorktreeGitDirPath(fs, currentPath)
	if err != nil {
		return "", "", errors.Errorf("could not find git dir for %s: %v", currentPath, err)
	}

	_, err = fs.Stat(worktreeGitPath)
	if err != nil {
		if os.IsNotExist(err) {
			// hardcoding error to get around windows-specific error message
			return "", "", errors.Errorf("could not find git dir for %s. %s does not exist", currentPath, worktreeGitPath)
		}
		return "", "", errors.Errorf("could not find git dir for %s: %v", currentPath, err)
	}

	// confirm whether the next directory up is the worktrees directory
	parent := path.Dir(worktreeGitPath)
	if path.Base(parent) == "worktrees" {
		gitDirPath := path.Dir(parent)
		return gitDirPath, path.Dir(gitDirPath), nil
	}

	// Unlike worktrees, submodules can be nested arbitrarily deep, so we check
	// if the `modules` directory is anywhere up the chain.
	if strings.Contains(worktreeGitPath, "/modules/") {
		// For submodules, we just return the path directly
		return worktreeGitPath, currentPath, nil
	}

	// If this error causes issues, we could relax the constraint and just always
	// return the path
	return "", "", errors.Errorf("could not find git dir for %s: path is not under `worktrees` or `modules` directories", currentPath)
}

// takes a path containing a symlink and returns the true path
func resolveSymlink(path string) (string, error) {
	l, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if l.Mode()&os.ModeSymlink == 0 {
		return path, nil
	}

	return filepath.EvalSymlinks(path)
}

// Returns the paths of linked worktrees
func linkedWortkreePaths(fs afero.Fs, repoGitDirPath string) []string {
	result := []string{}
	// For each directory in this path we're going to cat the `gitdir` file and append its contents to our result
	// That file points us to the `.git` file in the worktree.
	worktreeGitDirsPath := path.Join(repoGitDirPath, "worktrees")

	// ensure the directory exists
	_, err := fs.Stat(worktreeGitDirsPath)
	if err != nil {
		return result
	}

	_ = afero.Walk(fs, worktreeGitDirsPath, func(currPath string, info ioFs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		gitDirPath := path.Join(currPath, "gitdir")
		gitDirBytes, err := afero.ReadFile(fs, gitDirPath)
		if err != nil {
			// ignoring error
			return nil
		}
		trimmedGitDir := strings.TrimSpace(string(gitDirBytes))
		// removing the .git part
		worktreeDir := path.Dir(trimmedGitDir)
		result = append(result, worktreeDir)
		return nil
	})

	return result
}
