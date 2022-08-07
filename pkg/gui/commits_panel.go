package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// after selecting the 200th commit, we'll load in all the rest
const COMMIT_THRESHOLD = 200

// list panel functions

func (gui *Gui) getSelectedLocalCommit() *models.Commit {
	return gui.State.Contexts.LocalCommits.GetSelected()
}

func (gui *Gui) onCommitFocus() error {
	context := gui.State.Contexts.LocalCommits
	if context.GetSelectedLineIdx() > COMMIT_THRESHOLD && context.GetLimitCommits() {
		context.SetLimitCommits(false)
		go utils.Safe(func() {
			if err := gui.refreshCommitsWithLimit(); err != nil {
				_ = gui.c.Error(err)
			}
		})
	}

	return nil
}

func (gui *Gui) branchCommitsRenderToMain() error {
	var task types.UpdateTask
	commit := gui.State.Contexts.LocalCommits.GetSelected()
	if commit == nil {
		task = types.NewRenderStringTask(gui.c.Tr.NoCommitsThisBranch)
	} else {
		cmdObj := gui.git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())
		task = types.NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: "Patch",
			Task:  task,
		},
		Secondary: gui.secondaryPatchPanelUpdateOpts(),
	})
}

func (gui *Gui) secondaryPatchPanelUpdateOpts() *types.ViewUpdateOpts {
	if gui.git.Patch.PatchManager.Active() {
		patch := gui.git.Patch.PatchManager.RenderAggregatedPatchColored(false)

		return &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(patch),
			Title: gui.Tr.CustomPatch,
		}
	}

	return nil
}

func (gui *Gui) refForLog() string {
	bisectInfo := gui.git.Bisect.GetInfo()
	gui.State.Model.BisectInfo = bisectInfo

	if !bisectInfo.Started() {
		return "HEAD"
	}

	// need to see if our bisect's current commit is reachable from our 'new' ref.
	if bisectInfo.Bisecting() && !gui.git.Bisect.ReachableFromStart(bisectInfo) {
		return bisectInfo.GetNewSha()
	}

	return bisectInfo.GetStartSha()
}
