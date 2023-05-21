package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/diffing"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type DiffHelper struct {
	c *HelperCommon
}

func NewDiffHelper(c *HelperCommon) *DiffHelper {
	return &DiffHelper{
		c: c,
	}
}

func (self *DiffHelper) DiffArgs() []string {
	output := []string{self.c.Modes().Diffing.Ref}

	right := self.currentDiffTerminal()
	if right != "" {
		output = append(output, right)
	}

	if self.c.Modes().Diffing.Reverse {
		output = append(output, "-R")
	}

	if self.c.State().GetIgnoreWhitespaceInDiffView() {
		output = append(output, "--ignore-all-space")
	}

	output = append(output, "--")

	file := self.currentlySelectedFilename()
	if file != "" {
		output = append(output, file)
	} else if self.c.Modes().Filtering.Active() {
		output = append(output, self.c.Modes().Filtering.GetPath())
	}

	return output
}

func (self *DiffHelper) ExitDiffMode() error {
	self.c.Modes().Diffing = diffing.New()
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *DiffHelper) RenderDiff() error {
	cmdObj := self.c.Git().Diff.DiffCmdObj(self.DiffArgs())
	task := types.NewRunPtyTask(cmdObj.GetCmd())

	return self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title:    "Diff",
			SubTitle: self.IgnoringWhitespaceSubTitle(),
			Task:     task,
		},
	})
}

// CurrentDiffTerminals returns the current diff terminals of the currently selected item.
// in the case of a branch it returns both the branch and it's upstream name,
// which becomes an option when you bring up the diff menu, but when you're just
// flicking through branches it will be using the local branch name.
func (self *DiffHelper) CurrentDiffTerminals() []string {
	c := self.c.CurrentSideContext()

	if c.GetKey() == "" {
		return nil
	}

	switch v := c.(type) {
	case types.DiffableContext:
		return v.GetDiffTerminals()
	}

	return nil
}

func (self *DiffHelper) currentDiffTerminal() string {
	names := self.CurrentDiffTerminals()
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func (self *DiffHelper) currentlySelectedFilename() string {
	currentContext := self.c.CurrentContext()

	switch currentContext := currentContext.(type) {
	case types.IListContext:
		if lo.Contains([]types.ContextKey{context.FILES_CONTEXT_KEY, context.COMMIT_FILES_CONTEXT_KEY}, currentContext.GetKey()) {
			return currentContext.GetSelectedItemId()
		}
	}

	return ""
}

func (self *DiffHelper) WithDiffModeCheck(f func() error) error {
	if self.c.Modes().Diffing.Active() {
		return self.RenderDiff()
	}

	return f()
}

func (self *DiffHelper) IgnoringWhitespaceSubTitle() string {
	if self.c.State().GetIgnoreWhitespaceInDiffView() {
		return self.c.Tr.IgnoreWhitespaceDiffViewSubTitle
	}

	return ""
}
