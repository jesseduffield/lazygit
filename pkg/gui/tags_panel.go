package gui

func (gui *Gui) tagsRenderToMain() error {
	var task updateTask
	tag := gui.State.Contexts.Tags.GetSelected()
	if tag == nil {
		task = NewRenderStringTask("No tags")
	} else {
		cmdObj := gui.git.Branch.GetGraphCmdObj(tag.FullRefName())
		task = NewRunCommandTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		pair: gui.normalMainContextPair(),
		main: &viewUpdateOpts{
			title: "Tag",
			task:  task,
		},
	})
}
