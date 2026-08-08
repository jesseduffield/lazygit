package git_commands

import (
	ioFs "io/fs"
	"os"
	"path/filepath"
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
		return nil, err
	}

	gitDirResults := strings.Split(utils.NormalizeLinefeeds(gitDirOutput), "\n")
	worktreePath := gitDirResults[0]
	worktreeGitDirPath := gitDirResults[1]
	repoGitDirPath := gitDirResults[2]
	isBareRepo := gitDirResults[3] == "true"

	// A worktree that has the repo's common git dir to itself is the repo's main
	// worktree, so it is the repoPath. That holds for a submodule as well: its
	// git dir lives under the superproject's .git/modules, but it is still the
	// submodule's own common dir.
	isMainWorktree := worktreeGitDirPath == repoGitDirPath

	// If we're in a submodule, --show-superproject-working-tree will return a
	// value, meaning gitDirResults will be length 5. That only tells us anything
	// new for a linked worktree of a submodule, which isMainWorktree misses.
	isSubmodule := len(gitDirResults) == 5

	// Otherwise we're in a linked worktree, and the repoPath is the repo's main
	// worktree. git won't tell us where that is: `git worktree list` reports it
	// as the common git dir with a trailing "/.git" removed, which is this same
	// derivation. So take the directory holding the common git dir. That is the
	// main worktree of an ordinary repo, and of a bare one it is the directory
	// its worktrees live in. It is not the main worktree of a repo that moved
	// that elsewhere with core.worktree; there we end up naming the git dir's
	// directory, which means that the repo name we display in the status panel
	// isn't correct, and we start looking for .lazygit.yml in the wrong place.
	// Both of those are not severe enough to justify the extra git call to get
	// the real main worktree, so we accept this for this rather niche use case.
	var repoPath string
	if isMainWorktree || isSubmodule {
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
