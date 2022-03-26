package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
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

	to := gui.State.Contexts.CommitFiles.GetRefName()
	from, reverse := gui.State.Modes.Diffing.GetFromAndReverseArgsForDiff(to)

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
	gui.State.Contexts.CommitFiles.SetRefName(opts.RefName)
	gui.State.Contexts.CommitFiles.SetTitleRef(opts.RefDescription)
	gui.State.Contexts.CommitFiles.SetCanRebase(opts.CanRebase)
	gui.State.Contexts.CommitFiles.SetParentContext(opts.Context)
	gui.State.Contexts.CommitFiles.SetWindowName(opts.Context.GetWindowName())

	if err := gui.refreshCommitFilesContext(); err != nil {
		return err
	}

	return gui.c.PushContext(gui.State.Contexts.CommitFiles)
}

func (gui *Gui) refreshCommitFilesContext() error {
	currentSideContext := gui.currentSideContext()
	if currentSideContext.GetKey() == context.COMMIT_FILES_CONTEXT_KEY ||
		currentSideContext.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
		if err := gui.handleRefreshPatchBuildingPanel(-1); err != nil {
			return err
		}
	}

	to := gui.State.Contexts.CommitFiles.GetRefName()
	from, reverse := gui.State.Modes.Diffing.GetFromAndReverseArgsForDiff(to)

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
