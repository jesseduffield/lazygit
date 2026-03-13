package git_commands

import (
	ioFs "io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spf13/afero"
)

type RepoPaths struct {
	worktreePath       string
	worktreeGitDirPath string
	repoPath           string
	repoGitDirPath     string
	repoName           string
	isBareRepo         bool
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

func (self *RepoPaths) IsBareRepo() bool {
	return self.isBareRepo
}

// Returns the repo paths for a typical repo
func MockRepoPaths(currentPath string) *RepoPaths {
	return &RepoPaths{
		worktreePath:       currentPath,
		worktreeGitDirPath: filepath.Join(currentPath, ".git"),
		repoPath:           currentPath,
		repoGitDirPath:     filepath.Join(currentPath, ".git"),
		repoName:           "lazygit",
		isBareRepo:         false,
	}
}

func GetRepoPaths(
	cmd oscommands.ICmdObjBuilder,
	version *GitVersion,
) (*RepoPaths, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return GetRepoPathsForDir(cwd, cmd)
}

func GetRepoPathsForDir(
	dir string,
	cmd oscommands.ICmdObjBuilder,
) (*RepoPaths, error) {
	gitDirOutput, err := callGitRevParseWithDir(cmd, dir, "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--is-bare-repository", "--show-superproject-working-tree")
	if err != nil {
		if strings.Contains(err.Error(), "must be run in a work tree") {
			repoPaths, bareErr := handleBareRepoSetup(dir, cmd)
			if bareErr == nil {
				return repoPaths, nil
			}
		}

		return nil, err
	}

	gitDirResults := strings.Split(utils.NormalizeLinefeeds(gitDirOutput), "\n")
	worktreePath := gitDirResults[0]
	worktreeGitDirPath := gitDirResults[1]
	repoGitDirPath := gitDirResults[2]
	isBareRepo := gitDirResults[3] == "true"

	// If we're in a submodule, --show-superproject-working-tree will return
	// a value, meaning gitDirResults will be length 5. In that case
	// return the worktree path as the repoPath. Otherwise we're in a
	// normal repo or a worktree so return the parent of the git common
	// dir (repoGitDirPath)
	isSubmodule := len(gitDirResults) == 5

	var repoPath string
	if isSubmodule {
		repoPath = worktreePath
	} else {
		repoPath = filepath.Dir(repoGitDirPath)
	}
	repoName := filepath.Base(repoPath)

	return &RepoPaths{
		worktreePath:       worktreePath,
		worktreeGitDirPath: worktreeGitDirPath,
		repoPath:           repoPath,
		repoGitDirPath:     repoGitDirPath,
		repoName:           repoName,
		isBareRepo:         isBareRepo,
	}, nil
}

func callGitRevParseWithDir(
	cmd oscommands.ICmdObjBuilder,
	dir string,
	gitRevArgs ...string,
) (string, error) {
	gitRevParse := NewGitCmd("rev-parse").Arg("--path-format=absolute").Arg(gitRevArgs...)
	if dir != "" {
		gitRevParse.Dir(dir)
	}

	gitCmd := cmd.New(gitRevParse.ToArgv()).DontLog()
	res, err := gitCmd.RunWithOutput()
	if err != nil {
		return "", errors.Errorf("'%s' failed: %v", gitCmd.ToString(), err)
	}
	return strings.TrimSpace(res), nil
}

// Returns the paths of linked worktrees
func linkedWortkreePaths(fs afero.Fs, repoGitDirPath string) []string {
	result := []string{}
	// For each directory in this path we're going to cat the `gitdir` file and append its contents to our result
	// That file points us to the `.git` file in the worktree.
	worktreeGitDirsPath := filepath.Join(repoGitDirPath, "worktrees")

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

		gitDirPath := filepath.Join(currPath, "gitdir")
		gitDirBytes, err := afero.ReadFile(fs, gitDirPath)
		if err != nil {
			// ignoring error
			return nil
		}
		trimmedGitDir := strings.TrimSpace(string(gitDirBytes))
		// removing the .git part
		worktreeDir := filepath.Dir(trimmedGitDir)
		result = append(result, worktreeDir)
		return nil
	})

	return result
}

func handleBareRepoSetup(dir string, cmd oscommands.ICmdObjBuilder) (*RepoPaths, error) {
	if strings.Contains(dir, "..") || strings.Contains(dir, "~") {
		return nil, errors.New("invalid directory path: potential path traversal detected")
	}

	commonPatterns := []string{".bare", "bare.git", ".git"}

	for _, pattern := range commonPatterns {
		bareDir := filepath.Join(dir, pattern)

		if isBareRepo(bareDir, cmd) {
			return handleBareRepoWithWorktrees(dir, bareDir, cmd)
		}
	}

	bareDir, err := findBareRepoInDir(dir, cmd)
	if err == nil && bareDir != "" {
		return handleBareRepoWithWorktrees(dir, bareDir, cmd)
	}

	return nil, errors.New("no bare repo setup detected")
}

func findBareRepoInDir(dir string, cmd oscommands.ICmdObjBuilder) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		candidateDir := filepath.Join(dir, entry.Name())
		if isBareRepo(candidateDir, cmd) {
			return candidateDir, nil
		}
	}

	return "", errors.New("no bare repository found in directory")
}

func isBareRepo(gitDir string, cmd oscommands.ICmdObjBuilder) bool {
	gitCmd := cmd.New(NewGitCmd("rev-parse").Arg("--is-bare-repository").Dir(gitDir).ToArgv()).DontLog()
	output, err := gitCmd.RunWithOutput()
	return err == nil && strings.TrimSpace(output) == "true"
}

func handleBareRepoWithWorktrees(parentDir, bareDir string, cmd oscommands.ICmdObjBuilder) (*RepoPaths, error) {
	gitCmd := cmd.New(NewGitCmd("worktree").Arg("list").Arg("--porcelain").Dir(bareDir).ToArgv()).DontLog()
	output, err := gitCmd.RunWithOutput()
	if err != nil {
		return nil, errors.Errorf("failed to list worktrees: %v", err)
	}

	worktrees := parseWorktreeList(output)
	if len(worktrees) == 0 {
		return nil, errors.New("no worktrees found for bare repository")
	}

	selectedWorktree := selectBestWorktree(worktrees, parentDir, bareDir, cmd)

	if !strings.HasPrefix(selectedWorktree.Path, parentDir) {
		return nil, errors.Errorf("worktree path %s is not under parent directory %s", selectedWorktree.Path, parentDir)
	}

	return GetRepoPathsForDir(selectedWorktree.Path, cmd)
}

type WorktreeInfo struct {
	Path   string
	Head   string
	Branch string
}

func parseWorktreeList(output string) []WorktreeInfo {
	var worktrees []WorktreeInfo
	lines := strings.Split(strings.TrimSpace(output), "\n")

	var current WorktreeInfo
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = WorktreeInfo{}
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "HEAD ") {
			current.Head = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch ")
		}
	}

	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees
}

func GetDefaultBranch(gitDir string, cmd oscommands.ICmdObjBuilder) string {
	// Try to get the default branch from the remote HEAD
	gitCmd := cmd.New(NewGitCmd("symbolic-ref").
		Arg("refs/remotes/origin/HEAD", "--short").
		Dir(gitDir).ToArgv()).DontLog()
	output, err := gitCmd.RunWithOutput()
	if err == nil {
		branchName := strings.TrimSpace(output)
		if strings.HasPrefix(branchName, "origin/") {
			return strings.TrimPrefix(branchName, "origin/")
		}

		return branchName
	}

	// Try to get the default branch from git config
	gitCmd = cmd.New(NewGitCmd("config").
		Arg("init.defaultBranch").
		Dir(gitDir).ToArgv()).DontLog()
	output, err = gitCmd.RunWithOutput()
	if err == nil {
		branchName := strings.TrimSpace(output)
		if branchName != "" {
			return branchName
		}
	}

	return "main"
}

func selectBestWorktree(worktrees []WorktreeInfo, parentDir string, bareDir string, cmd oscommands.ICmdObjBuilder) WorktreeInfo {
	if len(worktrees) == 1 {
		return worktrees[0]
	}

	var candidates []WorktreeInfo

	for _, wt := range worktrees {
		if strings.HasPrefix(wt.Path, parentDir) {
			candidates = append(candidates, wt)
		}
	}

	if len(candidates) == 0 {
		return worktrees[0]
	}

	defaultBranch := GetDefaultBranch(bareDir, cmd)

	for _, wt := range candidates {
		branch := strings.TrimPrefix(wt.Branch, "refs/heads/")
		if branch == defaultBranch {
			return wt
		}
	}

	// Fall back to checking for "main" or "master" if default branch not found
	for _, wt := range candidates {
		branch := strings.TrimPrefix(wt.Branch, "refs/heads/")
		if branch == "main" || branch == "master" {
			return wt
		}
	}

	for _, wt := range candidates {
		dirName := filepath.Base(wt.Path)
		if dirName == defaultBranch {
			return wt
		}
	}

	// Fall back to checking for "main" or "master" directory names
	for _, wt := range candidates {
		dirName := filepath.Base(wt.Path)
		if dirName == "main" || dirName == "master" {
			return wt
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return filepath.Base(candidates[i].Path) < filepath.Base(candidates[j].Path)
	})

	return candidates[0]
}
