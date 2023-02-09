package gui

import "github.com/jesseduffield/lazygit/pkg/gui/types"

func (gui *Gui) reflogCommitsRenderToMain() error {
	commit := gui.State.Contexts.ReflogCommits.GetSelected()
	var task types.UpdateTask
	if commit == nil {
		task = types.NewRenderStringTask("No reflog history")
	} else {
		cmdObj := gui.git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath(),
			gui.IgnoreWhitespaceInDiffView)

		task = types.NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: "Reflog Entry",
			Task:  task,
		},
	})
}
