package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) onSubCommitFocus() error {
	context := gui.State.Contexts.SubCommits
	if context.GetSelectedLineIdx() > COMMIT_THRESHOLD && context.GetLimitCommits() {
		context.SetLimitCommits(false)
		go utils.Safe(func() {
			if err := gui.refreshSubCommitsWithLimit(); err != nil {
				_ = gui.c.Error(err)
			}
		})
	}

	return nil
}

func (gui *Gui) subCommitsRenderToMain() error {
	commit := gui.State.Contexts.SubCommits.GetSelected()
	var task types.UpdateTask
	if commit == nil {
		task = types.NewRenderStringTask("No commits")
	} else {
		cmdObj := gui.git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath(),
			gui.IgnoreWhitespaceInDiffView)

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
