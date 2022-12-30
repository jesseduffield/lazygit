package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReflogCommitsController struct {
	baseController
	*controllerCommon
	context *context.ReflogCommitsContext
}

var _ types.IController = &ReflogCommitsController{}

func NewReflogCommitsController(
	common *controllerCommon,
	context *context.ReflogCommitsContext,
) *ReflogCommitsController {
	return &ReflogCommitsController{
		baseController:   baseController{},
		controllerCommon: common,
		context:          context,
	}
}

func (self *ReflogCommitsController) Context() types.Context {
	return self.context
}

func (self *ReflogCommitsController) GetOnRenderToMain() func() error {
	return func() error {
		return self.helpers.Diff.WithDiffModeCheck(func() error {
			commit := self.context.GetSelected()
			var task types.UpdateTask
			if commit == nil {
				task = types.NewRenderStringTask("No reflog history")
			} else {
				cmdObj := self.c.Git().Commit.ShowCmdObj(commit.Sha, self.c.Modes().Filtering.GetPath(), self.c.State().GetIgnoreWhitespaceInDiffView())

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
