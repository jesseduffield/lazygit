package gui

func (gui *Gui) branchesRenderToMain() error {
	var task updateTask
	branch := gui.State.Contexts.Branches.GetSelected()
	if branch == nil {
		task = NewRenderStringTask(gui.c.Tr.NoBranchesThisRepo)
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(branch.FullRefName())

		task = NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: gui.c.Tr.LogTitle,
			task:  task,
		},
	})
}
