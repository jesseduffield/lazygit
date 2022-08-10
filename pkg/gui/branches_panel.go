package gui

import "github.com/jesseduffield/lazygit/pkg/gui/types"

func (gui *Gui) branchesRenderToMain() error {
	var task types.UpdateTask
	branch := gui.State.Contexts.Branches.GetSelected()
	if branch == nil {
		task = types.NewRenderStringTask(gui.c.Tr.NoBranchesThisRepo)
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(branch.FullRefName())

		task = types.NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: gui.c.Tr.LogTitle,
			Task:  task,
		},
	})
}
