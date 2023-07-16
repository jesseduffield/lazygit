package helpers

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IWorktreeHelper interface {
	GetMainWorktreeName() string
	GetCurrentWorktreeName() string
}

type WorktreeHelper struct {
	c           *HelperCommon
	reposHelper *ReposHelper
}

func NewWorktreeHelper(c *HelperCommon, reposHelper *ReposHelper) *WorktreeHelper {
	return &WorktreeHelper{
		c:           c,
		reposHelper: reposHelper,
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
		HandleConfirm: func(path string) error {
			return self.c.Prompt(types.PromptOpts{
				Title: self.c.Tr.NewWorktreePath,
				HandleConfirm: func(committish string) error {
					return self.c.WithWaitingStatus(self.c.Tr.AddingWorktree, func(gocui.Task) error {
						self.c.LogAction(self.c.Tr.Actions.AddWorktree)
						if err := self.c.Git().Worktree.New(sanitizedBranchName(path), committish); err != nil {
							return err
						}
						return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
					})
				},
			})
		},
	})
}

func (self *WorktreeHelper) Switch(worktree *models.Worktree, contextKey types.ContextKey) error {
	if self.c.Git().Worktree.IsCurrentWorktree(worktree) {
		return self.c.ErrorMsg(self.c.Tr.AlreadyInWorktree)
	}

	self.c.LogAction(self.c.Tr.SwitchToWorktree)

	// if we were in a submodule, we want to forget about that stack of repos
	// so that hitting escape in the new repo does nothing
	self.c.State().GetRepoPathStack().Clear()

	return self.reposHelper.DispatchSwitchTo(worktree.Path, true, self.c.Tr.ErrWorktreeMovedOrRemoved, contextKey)
}
