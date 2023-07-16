package git_commands

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

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

func (self *WorktreeCommands) New(worktreePath string, committish string) error {
	cmdArgs := NewGitCmd("worktree").Arg("add", worktreePath, committish).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *WorktreeCommands) Delete(worktreePath string, force bool) error {
	cmdArgs := NewGitCmd("worktree").Arg("remove").ArgIf(force, "-f").Arg(worktreePath).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *WorktreeCommands) IsCurrentWorktree(w *models.Worktree) bool {
	return IsCurrentWorktree(w)
}

func IsCurrentWorktree(w *models.Worktree) bool {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return EqualPath(pwd, w.Path)
}

func (self *WorktreeCommands) IsWorktreePathMissing(w *models.Worktree) bool {
	if _, err := os.Stat(w.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return true
		}
		log.Fatalln(fmt.Errorf("failed to check if worktree path `%s` exists\n%w", w.Path, err).Error())
	}
	return false
}

// checks if two paths are equal
// TODO: support relative paths
func EqualPath(a string, b string) bool {
	return a == b
}
