package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
)

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

func (gui *Gui) filesRenderToMain() error {
	node := gui.getSelectedFileNode()

	if node == nil {
		return gui.refreshMainViews(refreshMainOpts{
			pair: gui.normalMainContextPair(),
			main: &viewUpdateOpts{
				title: gui.c.Tr.DiffTitle,
				task:  NewRenderStringTask(gui.c.Tr.NoChangedFiles),
			},
		})
	}

	if node.File != nil && node.File.HasInlineMergeConflicts {
		hasConflicts, err := gui.helpers.MergeConflicts.SetMergeState(node.GetPath())
		if err != nil {
			return err
		}

		if hasConflicts {
			return gui.refreshMergePanel(false)
		}
	}

	gui.helpers.MergeConflicts.ResetMergeState()

	pair := gui.normalMainContextPair()
	if node.File != nil {
		pair = gui.stagingMainContextPair()
	}

	split := gui.c.UserConfig.Gui.SplitDiff == "always" || (node.GetHasUnstagedChanges() && node.GetHasStagedChanges())
	mainShowsStaged := !split && node.GetHasStagedChanges()

	cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, mainShowsStaged, gui.IgnoreWhitespaceInDiffView)
	title := gui.c.Tr.UnstagedChanges
	if mainShowsStaged {
		title = gui.c.Tr.StagedChanges
	}
	refreshOpts := refreshMainOpts{
		pair: pair,
		main: &viewUpdateOpts{
			task:  NewRunPtyTask(cmdObj.GetCmd()),
			title: title,
		},
	}

	if split {
		cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, true, gui.IgnoreWhitespaceInDiffView)

		title := gui.c.Tr.StagedChanges
		if mainShowsStaged {
			title = gui.c.Tr.UnstagedChanges
		}

		refreshOpts.secondary = &viewUpdateOpts{
			title: title,
			task:  NewRunPtyTask(cmdObj.GetCmd()),
		}
	}

	return gui.refreshMainViews(refreshOpts)
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
