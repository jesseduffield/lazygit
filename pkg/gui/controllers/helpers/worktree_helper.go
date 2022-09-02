package helpers

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
