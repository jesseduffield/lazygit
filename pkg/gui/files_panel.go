package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedFileNode() *filetree.FileNode {
	selectedLine := gui.State.Panels.Files.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.FileTreeViewModel.GetItemAtIndex(selectedLine)
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

	cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, !node.GetHasUnstagedChanges() && node.GetHasStagedChanges(), gui.IgnoreWhitespaceInDiffView)

	refreshOpts := refreshMainOpts{main: &viewUpdateOpts{
		title: gui.c.Tr.UnstagedChanges,
		task:  NewRunPtyTask(cmdObj.GetCmd()),
	}}

	if node.GetHasUnstagedChanges() {
		if node.GetHasStagedChanges() {
			cmdObj := gui.git.WorkingTree.WorktreeFileDiffCmdObj(node, false, true, gui.IgnoreWhitespaceInDiffView)

			refreshOpts.secondary = &viewUpdateOpts{
				title: gui.c.Tr.StagedChanges,
				task:  NewRunPtyTask(cmdObj.GetCmd()),
			}
		}
	} else {
		refreshOpts.main.title = gui.c.Tr.StagedChanges
	}

	return gui.refreshMainViews(refreshOpts)
}

// promptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (gui *Gui) promptToContinueRebase() error {
	gui.takeOverMergeConflictScrolling()

	return gui.PopupHandler.Ask(types.AskOpts{
		Title:  "continue",
		Prompt: gui.Tr.ConflictsResolved,
		HandleConfirm: func() error {
			return gui.genericMergeCommand(REBASE_OPTION_CONTINUE)
		},
	})
}

// Let's try to find our file again and move the cursor to that.
// If we can't find our file, it was probably just removed by the user. In that
// case, we go looking for where the next file has been moved to. Given that the
// user could have removed a whole directory, we continue iterating through the old
// nodes until we find one that exists in the new set of nodes, then move the cursor
// to that.
// prevNodes starts from our previously selected node because we don't need to consider anything above that
func (gui *Gui) findNewSelectedIdx(prevNodes []*filetree.FileNode, currNodes []*filetree.FileNode) int {
	getPaths := func(node *filetree.FileNode) []string {
		if node == nil {
			return nil
		}
		if node.File != nil && node.File.IsRename() {
			return node.File.Names()
		} else {
			return []string{node.Path}
		}
	}

	for _, prevNode := range prevNodes {
		selectedPaths := getPaths(prevNode)

		for idx, node := range currNodes {
			paths := getPaths(node)

			// If you started off with a rename selected, and now it's broken in two, we want you to jump to the new file, not the old file.
			// This is because the new should be in the same position as the rename was meaning less cursor jumping
			foundOldFileInRename := prevNode.File != nil && prevNode.File.IsRename() && node.Path == prevNode.File.PreviousName
			foundNode := utils.StringArraysOverlap(paths, selectedPaths) && !foundOldFileInRename
			if foundNode {
				return idx
			}
		}
	}

	return -1
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
