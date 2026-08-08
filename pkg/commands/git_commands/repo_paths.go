package git_commands

import (
	ioFs "io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/env"
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
	gitLocationEnvVars []string
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

// Whether we found no worktree, so that there is nothing for lazygit to show.
// Note that this isn't quite git's core.bare: a repo that calls itself non-bare
// but whose worktree we couldn't find counts as bare for us too. Concretely,
// this is true when we're in
//
//   - a genuinely bare repo;
//   - the git dir of a linked worktree (.git/worktrees/x), whose worktree is
//     recorded but not somewhere we look;
//   - a repo that keeps its worktree somewhere only GIT_WORK_TREE knows, such
//     as a vcsh-style dotfiles repo that hasn't been given core.worktree.
//
// The .git dir of an ordinary repo is not one of them: GetRepoPathsForDir
// notices the worktree holding it and hands back that repo instead.
func (self *RepoPaths) IsBareRepo() bool {
	return self.isBareRepo
}

// The environment that tells git where this repo is, as "NAME=value" entries.
// It is empty for the vast majority of repos, which git finds for itself by
// looking for a .git in the directory a command runs in. It is only non-empty
// when that doesn't work — when the git dir lives somewhere else entirely,
// because of core.worktree or --work-tree — and then every command addressing
// the repo has to carry it.
func (self *RepoPaths) GitLocationEnvVars() []string {
	return self.gitLocationEnvVars
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
	repoPaths, err := repoPathsForDir(dir, cmd)
	if err != nil || !repoPaths.IsBareRepo() {
		return repoPaths, err
	}

	// We're in a git dir rather than in a working tree, which usually just means
	// somebody ran lazygit in the .git of an ordinary repo. git's convention is
	// that a git dir called .git belongs to the directory holding it, so look
	// there: if that is a working tree, it is the repo we were asked about, and
	// there's no reason to make the user go up a directory and try again.
	//
	// The git dirs that aren't called .git keep the paths we have. A linked
	// worktree's (.git/worktrees/x) and a submodule's (.git/modules/x) do have a
	// working tree, but only the directory holding a .git tells us where, so we
	// would be guessing. A bare repo's has none to find.
	if filepath.Base(repoPaths.WorktreeGitDirPath()) != ".git" {
		return repoPaths, nil
	}

	pathsFromWorkTree, err := repoPathsForDir(filepath.Dir(repoPaths.WorktreeGitDirPath()), cmd)
	if err != nil || pathsFromWorkTree.IsBareRepo() {
		return repoPaths, nil
	}
	return pathsFromWorkTree, nil
}

// repoPathsForDir asks git about the repo at dir, and reports a bare repo when
// there is no working tree there. Unlike GetRepoPathsForDir it never looks
// anywhere but dir, which is what keeps that one from going round in circles.
func repoPathsForDir(
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
		gitLocationEnvVars: gitLocationEnvVars(cmd, worktreePath, worktreeGitDirPath),
	}, nil
}

// gitLocationEnvVars works out whether git can find the repo by itself when a
// command runs in its worktree, and if it can't, returns the environment that
// tells git where it is. See RepoPaths.GitLocationEnvVars.
func gitLocationEnvVars(
	cmd oscommands.ICmdObjBuilder,
	worktreePath string,
	worktreeGitDirPath string,
) []string {
	// The ordinary repo, where the git dir sits in the worktree. Both paths are
	// git's own answers from the same invocation, so they are spelled alike and
	// comparing them is safe.
	if worktreeGitDirPath == filepath.Join(worktreePath, ".git") {
		return nil
	}

	// A linked worktree or a submodule instead has a .git file naming its git
	// dir, and git follows that just as happily. We could read the file, but the
	// path in it may well name the same directory differently than git did
	// above, so ask git to resolve it — from the worktree and nothing else.
	discoveredGitDirPath, err := callGitRevParseInOtherRepo(cmd, worktreePath, "--absolute-git-dir")
	if err == nil && discoveredGitDirPath == worktreeGitDirPath {
		return nil
	}

	return []string{
		env.GitDirEnvVar + "=" + worktreeGitDirPath,
		env.GitWorkTreeEnvVar + "=" + worktreePath,
	}
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

// Asks git about the repo at dir. This is how we find our own repo, so it has
// to be answered the way git itself would answer it there, GIT_DIR and
// GIT_WORK_TREE included.
func callGitRevParseWithDir(
	cmd oscommands.ICmdObjBuilder,
	dir string,
	gitRevArgs ...string,
) (string, error) {
	return runGitRevParse(newGitRevParseCmd(cmd, dir, gitRevArgs...))
}

// Asks git about a repo that isn't the one we have open; see forOtherRepo.
func callGitRevParseInOtherRepo(
	cmd oscommands.ICmdObjBuilder,
	dir string,
	gitRevArgs ...string,
) (string, error) {
	return runGitRevParse(forOtherRepo(newGitRevParseCmd(cmd, dir, gitRevArgs...)))
}

func newGitRevParseCmd(
	cmd oscommands.ICmdObjBuilder,
	dir string,
	gitRevArgs ...string,
) *oscommands.CmdObj {
	gitRevParse := NewGitCmd("rev-parse").Arg("--path-format=absolute").Arg(gitRevArgs...)
	if dir != "" {
		gitRevParse.Dir(dir)
	}

	return cmd.New(gitRevParse.ToArgv()).DontLog()
}

func runGitRevParse(gitCmd *oscommands.CmdObj) (string, error) {
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
