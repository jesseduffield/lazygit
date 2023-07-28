package git_commands

import (
	"path/filepath"

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
