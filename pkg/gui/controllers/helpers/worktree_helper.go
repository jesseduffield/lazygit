package helpers

import "github.com/jesseduffield/lazygit/pkg/gui/types"

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

func (self *WorktreeHelper) NewWorktree() error {
	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.NewWorktreePath,
		HandleConfirm: func(response string) error {
			self.c.LogAction(self.c.Tr.Actions.CreateWorktree)
			if err := self.c.Git().Worktree.New(sanitizedBranchName(response)); err != nil {
				return err
			}

			//if self.c.CurrentContext() != self.contexts.Worktrees {
			//	if err := self.c.PushContext(self.contexts.Worktrees); err != nil {
			//		return err
			//	}
			//}

			// self.contexts.LocalCommits.SetSelectedLineIdx(0)
			// self.contexts.Branches.SetSelectedLineIdx(0)

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
		},
	})
}

//func (self *WorktreeHelper) GetCurrentWorktreeName() string {
//	for _, worktree := range self.c.Model().Worktrees {
//		if worktree.Current() {
//			if worktree.Main() {
//				return ""
//			}
//			return worktree.Name()
//		}
//	}
//
//	return ""
//}
