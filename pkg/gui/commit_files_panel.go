package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
)

// TODO: do we need this?
func (gui *Gui) onCommitFileFocus() error {
	gui.escapeLineByLinePanel()
	return nil
}

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

	mainContext := gui.State.Contexts.Normal
	if node.File != nil {
		mainContext = gui.State.Contexts.PatchBuilding
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title:   "Patch",
			task:    task,
			context: mainContext,
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

	if err := gui.refreshCommitFilesContext(); err != nil {
		return err
	}

	return gui.c.PushContext(gui.State.Contexts.CommitFiles)
}

func (gui *Gui) refreshCommitFilesContext() error {
	ref := gui.State.Contexts.CommitFiles.GetRef()
	to := ref.RefName()
	from, reverse := gui.State.Modes.Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

	files, err := gui.git.Loaders.CommitFiles.GetFilesInDiff(from, to, reverse)
	if err != nil {
		return gui.c.Error(err)
	}
	gui.State.Model.CommitFiles = files
	gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.SetTree()

	return gui.c.PostRefreshUpdate(gui.State.Contexts.CommitFiles)
}

func (gui *Gui) getSelectedCommitFileName() string {
	node := gui.State.Contexts.CommitFiles.GetSelected()
	if node == nil {
		return ""
	}

	return node.Path
}
