package gui

import "github.com/jesseduffield/lazygit/pkg/gui/types"

func (gui *Gui) tagsRenderToMain() error {
	var task types.UpdateTask
	tag := gui.State.Contexts.Tags.GetSelected()
	if tag == nil {
		task = types.NewRenderStringTask("No tags")
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(tag.FullRefName())
		task = types.NewRunCommandTask(cmdObj.GetCmd())
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: "Tag",
			Task:  task,
		},
	})
}
