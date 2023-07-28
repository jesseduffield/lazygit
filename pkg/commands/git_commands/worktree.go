package git_commands

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type WorktreeCommands struct {
	*GitCommon
}

func NewWorktreeCommands(gitCommon *GitCommon) *WorktreeCommands {
	return &WorktreeCommands{
		GitCommon: gitCommon,
	}
}

type NewWorktreeOpts struct {
	// required. The path of the new worktree.
	Path string
	// required. The base branch/ref.
	Base string

	// if true, ends up with a detached head
	Detach bool

	// optional. if empty, and if detach is false, we will checkout the base
	Branch string
}

func (self *WorktreeCommands) New(opts NewWorktreeOpts) error {
	if opts.Detach && opts.Branch != "" {
		panic("cannot specify branch when detaching")
	}

	cmdArgs := NewGitCmd("worktree").Arg("add").
		ArgIf(opts.Detach, "--detach").
		ArgIf(opts.Branch != "", "-b", opts.Branch).
		Arg(opts.Path, opts.Base)

	return self.cmd.New(cmdArgs.ToArgv()).Run()
}

func (self *WorktreeCommands) Delete(worktreePath string, force bool) error {
	cmdArgs := NewGitCmd("worktree").Arg("remove").ArgIf(force, "-f").Arg(worktreePath).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *WorktreeCommands) Detach(worktreePath string) error {
	cmdArgs := NewGitCmd("checkout").Arg("--detach").GitDir(filepath.Join(worktreePath, ".git")).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func WorktreeForBranch(branch *models.Branch, worktrees []*models.Worktree) (*models.Worktree, bool) {
	for _, worktree := range worktrees {
		if worktree.Branch == branch.Name {
			return worktree, true
		}
	}

	return nil, false
}

func CheckedOutByOtherWorktree(branch *models.Branch, worktrees []*models.Worktree) bool {
	worktree, ok := WorktreeForBranch(branch, worktrees)
	if !ok {
		return false
	}

	return !worktree.IsCurrent
}

// If in a non-bare repo, this returns the path of the main worktree
// TODO: see if this works with a bare repo.
func GetCurrentRepoPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// check if .git is a file or a directory
	gitPath := filepath.Join(pwd, ".git")
	gitFileInfo, err := os.Stat(gitPath)
	if err != nil {
		// fallback
		return currentPath()
	}

	if gitFileInfo.IsDir() {
		// must be in the main worktree
		return currentPath()
	}

	// either in a submodule, a worktree, or a bare repo
	worktreeGitPath, ok := LinkedWorktreeGitPath(pwd)
	if !ok {
		// fallback
		return currentPath()
	}

	// confirm whether the next directory up is the 'worktrees' directory
	parent := filepath.Dir(worktreeGitPath)
	if filepath.Base(parent) != "worktrees" {
		// fallback
		return currentPath()
	}

	// now we just jump up two more directories to get the repo name
	return filepath.Dir(filepath.Dir(parent))
}

func GetCurrentRepoName() string {
	return filepath.Base(GetCurrentRepoPath())
}

func currentPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return pwd
}

func linkedWortkreePaths() []string {
	// first we need to get the repo dir
	repoPath := GetCurrentRepoPath()
	result := []string{}
	worktreePath := filepath.Join(repoPath, ".git", "worktrees")
	// for each directory in this path we're going to cat the `gitdir` file and append its contents to our result

	// ensure the directory exists
	_, err := os.Stat(worktreePath)
	if err != nil {
		return result
	}

	err = filepath.Walk(worktreePath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			gitDirPath := filepath.Join(path, "gitdir")
			gitDirBytes, err := os.ReadFile(gitDirPath)
			if err != nil {
				// ignoring error
				return nil
			}
			trimmedGitDir := strings.TrimSpace(string(gitDirBytes))
			// removing the .git part
			worktreeDir := filepath.Dir(trimmedGitDir)
			result = append(result, worktreeDir)
		}
		return nil
	})
	if err != nil {
		return result
	}

	return result
}
