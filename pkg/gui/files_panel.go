package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
		return gui.c.RenderToMainViews(types.RefreshMainOpts{
			Pair: gui.c.MainViewPairs().Normal,
			Main: &types.ViewUpdateOpts{
				Title: gui.c.Tr.DiffTitle,
				Task:  types.NewRenderStringTask(gui.c.Tr.NoChangedFiles),
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

	pair := gui.c.MainViewPairs().Normal
	if node.File != nil {
		pair = gui.c.MainViewPairs().Staging
	}

	split := gui.c.UserConfig.Gui.SplitDiff == "always" || (node.GetHasUnstagedChanges() && node.GetHasStagedChanges())
	mainShowsStaged := !split && node.GetHasStagedChanges()

	cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, mainShowsStaged, gui.IgnoreWhitespaceInDiffView)
	title := gui.c.Tr.UnstagedChanges
	if mainShowsStaged {
		title = gui.c.Tr.StagedChanges
	}
	refreshOpts := types.RefreshMainOpts{
		Pair: pair,
		Main: &types.ViewUpdateOpts{
			Task:  types.NewRunPtyTask(cmdObj.GetCmd()),
			Title: title,
		},
	}

	if split {
		cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, true, gui.IgnoreWhitespaceInDiffView)

		title := gui.c.Tr.StagedChanges
		if mainShowsStaged {
			title = gui.c.Tr.UnstagedChanges
		}

		refreshOpts.Secondary = &types.ViewUpdateOpts{
			Title: title,
			Task:  types.NewRunPtyTask(cmdObj.GetCmd()),
		}
	}

	return gui.c.RenderToMainViews(refreshOpts)
}

func (gui *Gui) getSetTextareaTextFn(getView func() *gocui.View) func(string) {
	return func(text string) {
		// using a getView function so that we don't need to worry about when the view is created
		view := getView()
		view.ClearTextArea()
		view.TextArea.TypeString(text)
		_ = gui.resizePopupPanel(view, view.TextArea.GetContent())
		view.RenderTextArea()
	}
}
