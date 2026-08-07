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
// be the same as RepoPath(). It is empty for a bare repo, which has no
// worktree at all.
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

// Whether the repo has no worktree, so that there is nothing for lazygit to
// show. Note that this isn't quite git's core.bare: a repo that calls itself
// non-bare but doesn't have a worktree either counts as bare for us.
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
	gitDirOutput, err := callGitRevParseWithDir(cmd, dir, "--show-toplevel", "--absolute-git-dir", "--git-common-dir", "--show-superproject-working-tree")
	if err != nil {
		// --show-toplevel is the only one of these that needs a work tree, and
		// git makes it fatal when there isn't one. So this may just mean we're in
		// a repo that has no work tree.
		return getBareRepoPathsForDir(dir, cmd, err)
	}

	gitDirResults := strings.Split(utils.NormalizeLinefeeds(gitDirOutput), "\n")
	worktreePath := gitDirResults[0]
	worktreeGitDirPath := gitDirResults[1]
	repoGitDirPath := gitDirResults[2]

	// A worktree that has the repo's common git dir to itself is the repo's main
	// worktree, so it is the repoPath. That holds for a submodule as well: its
	// git dir lives under the superproject's .git/modules, but it is still the
	// submodule's own common dir.
	isMainWorktree := worktreeGitDirPath == repoGitDirPath

	// If we're in a submodule, --show-superproject-working-tree will return a
	// value, meaning gitDirResults will be length 4. That only tells us anything
	// new for a linked worktree of a submodule, which isMainWorktree misses.
	isSubmodule := len(gitDirResults) == 4

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
		isBareRepo:         false,
	}, nil
}

// getBareRepoPathsForDir is the fallback for when we couldn't ask git for the
// work tree. Everything but --show-toplevel works fine without one, so if the
// remaining queries succeed we are in a bare repo, and we return what we know
// about it with an empty worktreePath. If they fail too we simply aren't in a
// repo, and the caller's original error says so better than ours would.
func getBareRepoPathsForDir(
	dir string,
	cmd oscommands.ICmdObjBuilder,
	errWithWorktree error,
) (*RepoPaths, error) {
	output, err := callGitRevParseWithDir(cmd, dir, "--absolute-git-dir", "--git-common-dir")
	if err != nil {
		return nil, errWithWorktree
	}

	results := strings.Split(utils.NormalizeLinefeeds(output), "\n")
	repoGitDirPath := results[1]
	// A bare repo has no worktree, and so no repo path in the sense the caller
	// with a worktree means. It doesn't matter much what we say here, because
	// nobody reads it: whoever is handed a bare repo either offers to open a
	// recent one instead (app.setupRepo) or is turned away by NewGitCommand. The
	// directory holding the git dir is the nearest thing there is to a repo
	// path.
	repoPath := filepath.Dir(repoGitDirPath)

	return &RepoPaths{
		worktreePath:       "",
		worktreeGitDirPath: results[0],
		repoPath:           repoPath,
		repoGitDirPath:     repoGitDirPath,
		repoName:           filepath.Base(repoPath),
		isBareRepo:         true,
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
