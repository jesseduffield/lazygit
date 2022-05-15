package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
)

// list panel functions

func (gui *Gui) getSelectedFileNode() *filetree.FileNode {
	return gui.State.Contexts.Files.GetSelected()
}

func (gui *Gui) getSelectedFile() *models.File {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}
	return node.File
}

func (gui *Gui) getSelectedPath() string {
	node := gui.getSelectedFileNode()
	if node == nil {
		return ""
	}

	return node.GetPath()
}

func (gui *Gui) filesRenderToMain() error {
	node := gui.getSelectedFileNode()

	if node == nil {
		return gui.refreshMainViews(refreshMainOpts{
			main: &viewUpdateOpts{
				title: "",
				task:  NewRenderStringTask(gui.c.Tr.NoChangedFiles),
			},
		})
	}

	if node.File != nil && node.File.HasInlineMergeConflicts {
		ok, err := gui.setConflictsAndRenderWithLock(node.GetPath(), false)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	gui.resetMergeStateWithLock()

	mainContext := gui.State.Contexts.Normal
	if node.File != nil {
		mainContext = gui.State.Contexts.Staging
	}

	split := gui.c.UserConfig.Gui.SplitDiff == "always" || (node.GetHasUnstagedChanges() && node.GetHasStagedChanges())
	mainShowsStaged := !split && node.GetHasStagedChanges()

	cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, mainShowsStaged, gui.IgnoreWhitespaceInDiffView)
	refreshOpts := refreshMainOpts{main: &viewUpdateOpts{
		title:   gui.c.Tr.UnstagedChanges,
		task:    NewRunPtyTask(cmdObj.GetCmd()),
		context: mainContext,
	}}
	if mainShowsStaged {
		refreshOpts.main.title = gui.c.Tr.StagedChanges
	}

	if split {
		cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, true, gui.IgnoreWhitespaceInDiffView)

		refreshOpts.secondary = &viewUpdateOpts{
			title:   gui.c.Tr.StagedChanges,
			task:    NewRunPtyTask(cmdObj.GetCmd()),
			context: mainContext,
		}
	}

	return gui.refreshMainViews(refreshOpts)
}

func (gui *Gui) onFocusFile() error {
	gui.takeOverMergeConflictScrolling()
	return nil
}

func (gui *Gui) getSetTextareaTextFn(getView func() *gocui.View) func(string) {
	return func(text string) {
		// using a getView function so that we don't need to worry about when the view is created
		view := getView()
		view.ClearTextArea()
		view.TextArea.TypeString(text)
		view.RenderTextArea()
	}
}
