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
	task := NewRunPtyTask(cmdObj.GetCmd())

	pair := gui.normalMainContextPair()
	if node.File != nil {
		pair = gui.patchBuildingMainContextPair()
	}

	return gui.refreshMainViews(refreshMainOpts{
		pair: pair,
		main: &viewUpdateOpts{
			title: gui.Tr.Patch,
			task:  task,
		},
		secondary: gui.secondaryPatchPanelUpdateOpts(),
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
