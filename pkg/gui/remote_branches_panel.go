package gui

func (gui *Gui) remoteBranchesRenderToMain() error {
	var task updateTask
	remoteBranch := gui.State.Contexts.RemoteBranches.GetSelected()
	if remoteBranch == nil {
		task = NewRenderStringTask("No branches for this remote")
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(remoteBranch.FullRefName())
		task = NewRunCommandTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Remote Branch",
			task:  task,
		},
	})
}
