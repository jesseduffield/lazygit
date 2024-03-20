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
}

var _ types.IController = &ReflogCommitsController{}

func NewReflogCommitsController(
	c *ControllerCommon,
) *ReflogCommitsController {
	return &ReflogCommitsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.Commit](
			c,
			c.Contexts().ReflogCommits,
			c.Contexts().ReflogCommits.GetSelected,
			c.Contexts().ReflogCommits.GetSelectedItems,
		),
		c: c,
	}
}

func (self *ReflogCommitsController) Context() types.Context {
	return self.context()
}

func (self *ReflogCommitsController) context() *context.ReflogCommitsContext {
	return self.c.Contexts().ReflogCommits
}

func (self *ReflogCommitsController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			commit := self.context().GetSelected()
			var task types.UpdateTask
			if commit == nil {
				task = types.NewRenderStringTask("No reflog history")
			} else {
				cmdObj := self.c.Git().Commit.ShowCmdObj(commit.Hash, self.c.Modes().Filtering.GetPath())

				task = types.NewRunPtyTask(cmdObj.GetCmd())
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Reflog Entry",
					Task:  task,
				},
			})
		})
	}
}
