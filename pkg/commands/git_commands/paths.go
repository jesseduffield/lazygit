package git_commands

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/samber/lo"
)

type RepoPaths interface {
	// Current working directory of the program. Currently, this will always
	// be the same as WorktreePath(), but in future we may support running
	// lazygit from inside a subdirectory of the worktree.
	CurrentPath() string
	// Path to the current worktree. If we're in the main worktree, this will
	// be the same as RepoPath()
	WorktreePath() string
	// Path of the worktree's git dir.
	// If we're in the main worktree, this will be the .git dir under the RepoPath().
	// If we're in a linked worktree, it will be the directory pointed at by the worktree's .git file
	WorktreeGitDirPath() string
	// Path of the repo. If we're in a the main worktree, this will be the same as WorktreePath()
	// If we're in a bare repo, it will be the parent folder of the bare repo
	RepoPath() string
	// path of the git-dir for the repo.
	// If this is a bare repo, it will be the location of the bare repo
	// If this is a non-bare repo, it will be the location of the .git dir in
	// the main worktree.
	RepoGitDirPath() string
	// Name of the repo. Basename of the folder containing the repo.
	RepoName() string
}

type RepoDirsImpl struct {
	currentPath        string
	worktreePath       string
	worktreeGitDirPath string
	repoPath           string
	repoGitDirPath     string
	repoName           string
}

var _ RepoPaths = &RepoDirsImpl{}

func (self *RepoDirsImpl) CurrentPath() string {
	return self.currentPath
}

func (self *RepoDirsImpl) WorktreePath() string {
	return self.worktreePath
}

func (self *RepoDirsImpl) WorktreeGitDirPath() string {
	return self.worktreeGitDirPath
}

func (self *RepoDirsImpl) RepoPath() string {
	return self.repoPath
}

func (self *RepoDirsImpl) RepoGitDirPath() string {
	return self.repoGitDirPath
}

func (self *RepoDirsImpl) RepoName() string {
	return self.repoName
}

func GetRepoPaths() (RepoPaths, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return &RepoDirsImpl{}, errors.Errorf("failed to get current path: %v", err)
	}

	// converting to forward slashes for the sake of windows (which uses backwards slashes). We want everything
	// to have forward slashes internally
	currentPath = filepath.ToSlash(currentPath)

	worktreePath := currentPath
	repoGitDirPath, repoPath, err := GetCurrentRepoGitDirPath(currentPath)
	if err != nil {
		return &RepoDirsImpl{}, errors.Errorf("failed to get repo git dir path: %v", err)
	}
	worktreeGitDirPath, err := worktreeGitDirPath(currentPath)
	if err != nil {
		return &RepoDirsImpl{}, errors.Errorf("failed to get worktree git dir path: %v", err)
	}
	repoName := path.Base(repoPath)

	return &RepoDirsImpl{
		currentPath:        currentPath,
		worktreePath:       worktreePath,
		worktreeGitDirPath: worktreeGitDirPath,
		repoPath:           repoPath,
		repoGitDirPath:     repoGitDirPath,
		repoName:           repoName,
	}, nil
}

// Returns the paths of linked worktrees
func linkedWortkreePaths(repoGitDirPath string) []string {
	result := []string{}
	// For each directory in this path we're going to cat the `gitdir` file and append its contents to our result
	// That file points us to the `.git` file in the worktree.
	worktreeGitDirsPath := path.Join(repoGitDirPath, "worktrees")

	// ensure the directory exists
	_, err := os.Stat(worktreeGitDirsPath)
	if err != nil {
		return result
	}

	_ = filepath.Walk(worktreeGitDirsPath, func(currPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		gitDirPath := path.Join(currPath, "gitdir")
		gitDirBytes, err := os.ReadFile(gitDirPath)
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

// Returns the path of the git-dir for the worktree. For linked worktrees, the worktree has
// a .git file that points to the git-dir (which itself lives in the git-dir
// of the repo)
func worktreeGitDirPath(worktreePath string) (string, error) {
	// if .git is a file, we're in a linked worktree, otherwise we're in
	// the main worktree
	dotGitPath := path.Join(worktreePath, ".git")
	gitFileInfo, err := os.Stat(dotGitPath)
	if err != nil {
		return "", err
	}

	if gitFileInfo.IsDir() {
		return dotGitPath, nil
	}

	return linkedWorktreeGitDirPath(worktreePath)
}

func linkedWorktreeGitDirPath(worktreePath string) (string, error) {
	dotGitPath := path.Join(worktreePath, ".git")
	gitFileContents, err := os.ReadFile(dotGitPath)
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
	return gitDir, nil
}

func GetCurrentRepoGitDirPath(currentPath string) (string, string, error) {
	var unresolvedGitPath string
	if env.GetGitDirEnv() != "" {
		unresolvedGitPath = env.GetGitDirEnv()
	} else {
		unresolvedGitPath = path.Join(currentPath, ".git")
	}

	gitPath, err := resolveSymlink(unresolvedGitPath)
	if err != nil {
		return "", "", err
	}

	// check if .git is a file or a directory
	gitFileInfo, err := os.Stat(gitPath)
	if err != nil {
		return "", "", err
	}

	if gitFileInfo.IsDir() {
		// must be in the main worktree
		return gitPath, path.Dir(gitPath), nil
	}

	// either in a submodule, or worktree
	worktreeGitPath, err := linkedWorktreeGitDirPath(currentPath)
	if err != nil {
		return "", "", errors.Errorf("could not find git dir for %s: %v", currentPath, err)
	}

	// confirm whether the next directory up is the worktrees/submodules directory
	parent := path.Dir(worktreeGitPath)
	if path.Base(parent) != "worktrees" && path.Base(parent) != "modules" {
		return "", "", errors.Errorf("could not find git dir for %s", currentPath)
	}

	// if it's a submodule, we treat it as its own repo
	if path.Base(parent) == "modules" {
		return worktreeGitPath, currentPath, nil
	}

	gitDirPath := path.Dir(parent)
	return gitDirPath, path.Dir(gitDirPath), nil
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
