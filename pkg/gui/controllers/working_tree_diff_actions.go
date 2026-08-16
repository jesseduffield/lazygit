package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// WorkingTreeDiffActions implements what the files panel offers on the diff it renders
// into the focused main view: the diff itself, for the commands that need to read lines
// out of it rather than off the screen.
type WorkingTreeDiffActions struct {
	c *ControllerCommon
}

var _ types.FocusedMainViewDiffSource = &WorkingTreeDiffActions{}

func NewWorkingTreeDiffActions(c *ControllerCommon) *WorkingTreeDiffActions {
	return &WorkingTreeDiffActions{c: c}
}

func (self *WorkingTreeDiffActions) context() *context.WorkingTreeContext {
	return self.c.Contexts().Files
}

// PlainDiff hands out the working tree's diff for the given files, taken from the
// side of the index that the asking pane shows.
func (self *WorkingTreeDiffActions) PlainDiff(pane types.DiffPaneContext, paths []string) string {
	node := self.context().GetSelected()
	if node == nil {
		return ""
	}
	// An error means there is no diff to be had, which for our purposes is the same as
	// an empty one.
	diff, _ := self.c.Git().WorkingTree.
		WorktreeFileDiffCmdObj(node, true, self.showsStagedSide(pane), paths).
		RunWithOutput()
	return diff
}

// showsStagedSide reports whether the given main pane is the one showing the staged
// side of a file's diff, which is always the lower one.
func (self *WorkingTreeDiffActions) showsStagedSide(pane types.DiffPaneContext) bool {
	return pane.GetKey() == self.c.Contexts().NormalSecondary.GetKey()
}
