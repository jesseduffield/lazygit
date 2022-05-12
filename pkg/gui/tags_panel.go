package gui

func (self *Gui) tagsRenderToMain() error {
	var task updateTask
	tag := self.State.Contexts.Tags.GetSelected()
	if tag == nil {
		task = NewRenderStringTask("No tags")
	} else {
		cmdObj := self.git.Branch.GetGraphCmdObj(tag.FullRefName())
		task = NewRunCommandTask(cmdObj.GetCmd())
	}

	return self.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Tag",
			task:  task,
		},
	})
}
