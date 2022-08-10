package gui

import "github.com/jesseduffield/lazygit/pkg/gui/types"

// list panel functions

func (gui *Gui) subCommitsRenderToMain() error {
	commit := gui.State.Contexts.SubCommits.GetSelected()
	var task types.UpdateTask
	if commit == nil {
		task = types.NewRenderStringTask("No commits")
	} else {
		cmdObj := gui.git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())

		task = types.NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: "Commit",
			Task:  task,
		},
	})
}
