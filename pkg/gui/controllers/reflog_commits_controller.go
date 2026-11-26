package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReflogCommitsController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c *ControllerCommon

	// what this panel offers on the diff it shows in the focused main view
	diffActions *CommitDiffActions
}

var _ types.IController = &ReflogCommitsController{}

func NewReflogCommitsController(
	c *ControllerCommon,
) *ReflogCommitsController {
	controller := &ReflogCommitsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().ReflogCommits,
			c.Contexts().ReflogCommits.GetSelected,
			c.Contexts().ReflogCommits.GetSelectedItems,
		),
		c: c,
	}
	controller.diffActions = NewCommitDiffActions(c, c.Contexts().ReflogCommits, controller.diffTarget)
	return controller
}

// diffTarget is the reflog entry the panel has selected, whose diff its main view
// shows. A reflog entry is never a commit of the checked-out branch as far as we are
// concerned, so nothing here may be rewritten.
func (self *ReflogCommitsController) diffTarget() *commitDiffTarget {
	commit := self.context().GetSelected()
	if commit == nil {
		return nil
	}
	return &commitDiffTarget{from: commit.ParentRefName(), to: commit.RefName()}
}

func (self *ReflogCommitsController) Context() types.Context {
	return self.context()
}

func (self *ReflogCommitsController) context() *context.ReflogCommitsContext {
	return self.c.Contexts().ReflogCommits
}

func (self *ReflogCommitsController) GetFocusedMainViewDiffSource() types.FocusedMainViewDiffSource {
	return self.diffActions
}

func (self *ReflogCommitsController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			commit := self.context().GetSelected()
			var task types.UpdateTask
			if commit == nil {
				task = types.NewRenderStringTask("No reflog history")
			} else {
				mode := self.c.Helpers().DiffLine.MainViewDiffMode()
				cmdObj := self.c.Git().Commit.ShowCmdObj(commit.Hash(), self.c.Helpers().Diff.FilterPathsForCommit(commit), mode)

				task = types.NewMainViewDiffTask(cmdObj.GetCmd(), mode)
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Reflog Entry",
					Task:  task,
				},
				Secondary: secondaryPatchPanelUpdateOpts(self.c),
			})
		})
	}
}
