package git_commands

import (
	"errors"
	"fmt"
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
	cmdArgs := NewGitCmd("checkout").Arg("--detach").ToArgv()

	return self.cmd.New(cmdArgs).SetWd(worktreePath).Run()
}

func (self *WorktreeCommands) IsCurrentWorktree(path string) bool {
	return IsCurrentWorktree(path)
}

func IsCurrentWorktree(path string) bool {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return EqualPath(pwd, path)
}

func (self *WorktreeCommands) IsWorktreePathMissing(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return true
		}
		log.Fatalln(fmt.Errorf("failed to check if worktree path `%s` exists\n%w", path, err).Error())
	}
	return false
}

// checks if two paths are equal
// TODO: support relative paths
func EqualPath(a string, b string) bool {
	return a == b
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

	return !IsCurrentWorktree(worktree.Path)
}

func GetCurrentRepoPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// check if .git is a file or a directory
	gitPath := filepath.Join(pwd, ".git")
	gitFileInfo, err := os.Stat(gitPath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if gitFileInfo.IsDir() {
		// must be in the main worktree
		return currentPath()
	}

	// must be a worktree or bare repo
	worktreeGitPath, ok := WorktreeGitPath(pwd)
	if !ok {
		// fallback
		return currentPath()
	}

	// now we just jump up three directories to get the repo name
	return filepath.Dir(filepath.Dir(filepath.Dir(worktreeGitPath)))
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
		if errors.Is(err, fs.ErrNotExist) {
			return result
		}
		log.Fatalln(err.Error())
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
		log.Fatalln(err.Error())
	}

	return result
}
