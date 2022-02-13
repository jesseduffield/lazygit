package gui

func (gui *Gui) stashRenderToMain() error {
	var task updateTask
	stashEntry := gui.State.Contexts.Stash.GetSelected()
	if stashEntry == nil {
		task = NewRenderStringTask(gui.c.Tr.NoStashEntries)
	} else {
		task = NewRunPtyTask(gui.git.Stash.ShowStashEntryCmdObj(stashEntry.Index).GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Stash",
			task:  task,
		},
	})
}
