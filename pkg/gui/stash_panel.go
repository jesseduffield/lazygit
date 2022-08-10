package gui

import "github.com/jesseduffield/lazygit/pkg/gui/types"

func (gui *Gui) stashRenderToMain() error {
	var task types.UpdateTask
	stashEntry := gui.State.Contexts.Stash.GetSelected()
	if stashEntry == nil {
		task = types.NewRenderStringTask(gui.c.Tr.NoStashEntries)
	} else {
		task = types.NewRunPtyTask(gui.git.Stash.ShowStashEntryCmdObj(stashEntry.Index).GetCmd())
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: "Stash",
			Task:  task,
		},
	})
}
