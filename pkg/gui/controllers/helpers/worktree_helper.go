package helpers

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IWorktreeHelper interface {
	GetMainWorktreeName() string
	GetCurrentWorktreeName() string
}

type WorktreeHelper struct {
	c *HelperCommon
}

func NewWorktreeHelper(c *HelperCommon) *WorktreeHelper {
	return &WorktreeHelper{
		c: c,
	}
}

func (self *WorktreeHelper) GetMainWorktreeName() string {
	for _, worktree := range self.c.Model().Worktrees {
		if worktree.Main() {
			return worktree.Name()
		}
	}

	return ""
}

func (self *WorktreeHelper) IsCurrentWorktree(w *models.Worktree) bool {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return pwd == w.Path
}

func (self *WorktreeHelper) IsWorktreePathMissing(w *models.Worktree) bool {
	if _, err := os.Stat(w.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return true
		}
		log.Fatalln(fmt.Errorf("failed to check if worktree path `%s` exists\n%w", w.Path, err).Error())
	}
	return false
}

func (self *WorktreeHelper) NewWorktree() error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.NewWorktreePath,
		HandleConfirm: func(response string) error {
			self.c.LogAction(self.c.Tr.Actions.CreateWorktree)
			if err := self.c.Git().Worktree.New(sanitizedBranchName(response)); err != nil {
				return err
			}
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}
