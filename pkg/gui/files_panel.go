package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
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

func (gui *Gui) refreshFilesAndSubmodules() error {
	gui.Mutexes.RefreshingFilesMutex.Lock()
	gui.State.IsRefreshingFiles = true
	defer func() {
		gui.State.IsRefreshingFiles = false
		gui.Mutexes.RefreshingFilesMutex.Unlock()
	}()

	prevSelectedPath := gui.getSelectedPath()

	if err := gui.refreshStateSubmoduleConfigs(); err != nil {
		return err
	}

	if err := gui.refreshMergeState(); err != nil {
		return err
	}

	if err := gui.refreshStateFiles(); err != nil {
		return err
	}

	gui.OnUIThread(func() error {
		if err := gui.c.PostRefreshUpdate(gui.State.Contexts.Submodules); err != nil {
			gui.c.Log.Error(err)
		}

		if types.ContextKey(gui.Views.Files.Context) == FILES_CONTEXT_KEY {
			// doing this a little custom (as opposed to using gui.c.PostRefreshUpdate) because we handle selecting the file explicitly below
			if err := gui.State.Contexts.Files.HandleRender(); err != nil {
				return err
			}
		}

		if gui.currentContext().GetKey() == FILES_CONTEXT_KEY {
			currentSelectedPath := gui.getSelectedPath()
			alreadySelected := prevSelectedPath != "" && currentSelectedPath == prevSelectedPath
			if !alreadySelected {
				gui.takeOverMergeConflictScrolling()
			}

			gui.Views.Files.FocusPoint(0, gui.State.Panels.Files.SelectedLineIdx)
			return gui.filesRenderToMain()
		}

		return nil
	})

	return nil
}

func (gui *Gui) refreshStateFiles() error {
	state := gui.State

	// keep track of where the cursor is currently and the current file names
	// when we refresh, go looking for a matching name
	// move the cursor to there.

	selectedNode := gui.getSelectedFileNode()

	prevNodes := gui.State.FileTreeViewModel.GetAllItems()
	prevSelectedLineIdx := gui.State.Panels.Files.SelectedLineIdx

	// If git thinks any of our files have inline merge conflicts, but they actually don't,
	// we stage them.
	// Note that if files with merge conflicts have both arisen and have been resolved
	// between refreshes, we won't stage them here. This is super unlikely though,
	// and this approach spares us from having to call `git status` twice in a row.
	// Although this also means that at startup we won't be staging anything until
	// we call git status again.
	pathsToStage := []string{}
	prevConflictFileCount := 0
	for _, file := range state.FileTreeViewModel.GetAllFiles() {
		if file.HasMergeConflicts {
			prevConflictFileCount++
		}
		if file.HasInlineMergeConflicts {
			hasConflicts, err := mergeconflicts.FileHasConflictMarkers(file.Name)
			if err != nil {
				gui.Log.Error(err)
			} else if !hasConflicts {
				pathsToStage = append(pathsToStage, file.Name)
			}
		}
	}

	if len(pathsToStage) > 0 {
		gui.c.LogAction(gui.Tr.Actions.StageResolvedFiles)
		if err := gui.git.WorkingTree.StageFiles(pathsToStage); err != nil {
			return gui.c.Error(err)
		}
	}

	files := gui.git.Loaders.Files.
		GetStatusFiles(loaders.GetStatusFileOptions{})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	if gui.git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE && conflictFileCount == 0 && prevConflictFileCount > 0 {
		gui.OnUIThread(func() error { return gui.promptToContinueRebase() })
	}

	// for when you stage the old file of a rename and the new file is in a collapsed dir
	state.FileTreeViewModel.RWMutex.Lock()
	for _, file := range files {
		if selectedNode != nil && selectedNode.Path != "" && file.PreviousName == selectedNode.Path {
			state.FileTreeViewModel.ExpandToPath(file.Name)
		}
	}

	// only taking over the filter if it hasn't already been set by the user.
	// Though this does make it impossible for the user to actually say they want to display all if
	// conflicts are currently being shown. Hmm. Worth it I reckon. If we need to add some
	// extra state here to see if the user's set the filter themselves we can do that, but
	// I'd prefer to maintain as little state as possible.
	if conflictFileCount > 0 {
		if state.FileTreeViewModel.GetFilter() == filetree.DisplayAll {
			state.FileTreeViewModel.SetFilter(filetree.DisplayConflicted)
		}
	} else if state.FileTreeViewModel.GetFilter() == filetree.DisplayConflicted {
		state.FileTreeViewModel.SetFilter(filetree.DisplayAll)
	}

	state.FileTreeViewModel.SetFiles(files)
	state.FileTreeViewModel.RWMutex.Unlock()

	if err := gui.fileWatcher.addFilesToFileWatcher(files); err != nil {
		return err
	}

	if selectedNode != nil {
		newIdx := gui.findNewSelectedIdx(prevNodes[prevSelectedLineIdx:], state.FileTreeViewModel.GetAllItems())
		if newIdx != -1 && newIdx != prevSelectedLineIdx {
			newNode := state.FileTreeViewModel.GetItemAtIndex(newIdx)
			// when not in tree mode, we show merge conflict files at the top, so you
			// can work through them one by one without having to sift through a large
			// set of files. If you have just fixed the merge conflicts of a file, we
			// actually don't want to jump to that file's new position, because that
			// file will now be ages away amidst the other files without merge
			// conflicts: the user in this case would rather work on the next file
			// with merge conflicts, which will have moved up to fill the gap left by
			// the last file, meaning the cursor doesn't need to move at all.
			leaveCursor := !state.FileTreeViewModel.InTreeMode() && newNode != nil &&
				selectedNode.File != nil && selectedNode.File.HasMergeConflicts &&
				newNode.File != nil && !newNode.File.HasMergeConflicts

			if !leaveCursor {
				state.Panels.Files.SelectedLineIdx = newIdx
			}
		}
	}

	gui.refreshSelectedLine(state.Panels.Files, state.FileTreeViewModel.GetItemsLength())
	return nil
}

// promptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (gui *Gui) promptToContinueRebase() error {
	gui.takeOverMergeConflictScrolling()

	return gui.PopupHandler.Ask(popup.AskOpts{
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
