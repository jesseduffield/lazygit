package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c *ControllerCommon
}

var _ types.IController = &SubCommitsController{}

func NewSubCommitsController(
	c *ControllerCommon,
) *SubCommitsController {
	return &SubCommitsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.Commit](
			c,
			c.Contexts().SubCommits,
			c.Contexts().SubCommits.GetSelected,
			c.Contexts().SubCommits.GetSelectedItems,
		),
		c: c,
	}
}

func (self *SubCommitsController) Context() types.Context {
	return self.context()
}

func (self *SubCommitsController) context() *context.SubCommitsContext {
	return self.c.Contexts().SubCommits
}

func (self *SubCommitsController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			commit := self.context().GetSelected()
			var task types.UpdateTask
			if commit == nil {
				task = types.NewRenderStringTask("No commits")
			} else {
				cmdObj := self.c.Git().Commit.ShowCmdObj(commit.Hash, self.c.Modes().Filtering.GetPath())

				task = types.NewRunPtyTask(cmdObj.GetCmd())
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title:    "Commit",
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Task:     task,
				},
			})
		})
	}
}

func (self *SubCommitsController) GetOnFocus() func(types.OnFocusOpts) error {
	return func(types.OnFocusOpts) error {
		context := self.context()
		if context.GetSelectedLineIdx() > COMMIT_THRESHOLD && context.GetLimitCommits() {
			context.SetLimitCommits(false)
			self.c.OnWorker(func(_ gocui.Task) error {
				return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.SUB_COMMITS}})
			})
		}

		return nil
	}
}
