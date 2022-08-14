package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) commitFilesRenderToMain() error {
	node := gui.State.Contexts.CommitFiles.GetSelected()
	if node == nil {
		return nil
	}

	ref := gui.State.Contexts.CommitFiles.GetRef()
	to := ref.RefName()
	from, reverse := gui.State.Modes.Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

	cmdObj := gui.git.WorkingTree.ShowFileDiffCmdObj(from, to, reverse, node.GetPath(), false)
	task := types.NewRunPtyTask(cmdObj.GetCmd())

	pair := gui.c.MainViewPairs().Normal
	if node.File != nil {
		pair = gui.c.MainViewPairs().PatchBuilding
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: pair,
		Main: &types.ViewUpdateOpts{
			Title: gui.Tr.Patch,
			Task:  task,
		},
		Secondary: gui.secondaryPatchPanelUpdateOpts(),
	})
}

func (gui *Gui) SwitchToCommitFilesContext(opts controllers.SwitchToCommitFilesContextOpts) error {
	gui.State.Contexts.CommitFiles.SetSelectedLineIdx(0)
	gui.State.Contexts.CommitFiles.SetRef(opts.Ref)
	gui.State.Contexts.CommitFiles.SetTitleRef(opts.Ref.Description())
	gui.State.Contexts.CommitFiles.SetCanRebase(opts.CanRebase)
	gui.State.Contexts.CommitFiles.SetParentContext(opts.Context)
	gui.State.Contexts.CommitFiles.SetWindowName(opts.Context.GetWindowName())

	if err := gui.c.Refresh(types.RefreshOptions{
		Scope: []types.RefreshableView{types.COMMIT_FILES},
	}); err != nil {
		return err
	}

	return gui.c.PushContext(gui.State.Contexts.CommitFiles)
}
