package gui

func (gui *Gui) branchesRenderToMain() error {
	var task updateTask
	branch := gui.State.Contexts.Branches.GetSelected()
	if branch == nil {
		task = NewRenderStringTask(gui.c.Tr.NoBranchesThisRepo)
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(branch.Name)

		task = NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Log",
			task:  task,
		},
	})
}
