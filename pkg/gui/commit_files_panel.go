package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) getSelectedCommitFileNode() *filetree.CommitFileNode {
	selectedLine := gui.State.Panels.CommitFiles.SelectedLineIdx
	if selectedLine == -1 || selectedLine > gui.State.CommitFileTreeViewModel.GetItemsLength()-1 {
		return nil
	}

	return gui.State.CommitFileTreeViewModel.GetItemAtIndex(selectedLine)
}

func (gui *Gui) getSelectedCommitFile() *models.CommitFile {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}
	return node.File
}

func (gui *Gui) getSelectedCommitFilePath() string {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return ""
	}
	return node.GetPath()
}

func (gui *Gui) onCommitFileFocus() error {
	gui.escapeLineByLinePanel()
	return nil
}

func (gui *Gui) commitFilesRenderToMain() error {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	to := gui.State.CommitFileTreeViewModel.GetParent()
	from, reverse := gui.getFromAndReverseArgsForDiff(to)

	cmdObj := gui.Git.WorkingTree.ShowFileDiffCmdObj(from, to, reverse, node.GetPath(), false)
	task := NewRunPtyTask(cmdObj.GetCmd())

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Patch",
			task:  task,
		},
		secondary: gui.secondaryPatchPanelUpdateOpts(),
	})
}

func (gui *Gui) handleCheckoutCommitFile() error {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	gui.logAction(gui.Tr.Actions.CheckoutFile)
	if err := gui.Git.WorkingTree.CheckoutFile(gui.State.CommitFileTreeViewModel.GetParent(), node.GetPath()); err != nil {
		return gui.PopupHandler.Error(err)
	}

	return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) handleDiscardOldFileChange() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	fileName := gui.getSelectedCommitFileName()

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.DiscardFileChangesTitle,
		Prompt: gui.Tr.DiscardFileChangesPrompt,
		HandleConfirm: func() error {
			return gui.PopupHandler.WithWaitingStatus(gui.Tr.RebasingStatus, func() error {
				gui.logAction(gui.Tr.Actions.DiscardOldFileChange)
				if err := gui.Git.Rebase.DiscardOldFileChanges(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, fileName); err != nil {
					if err := gui.handleGenericMergeCommandResult(err); err != nil {
						return err
					}
				}

				return gui.refreshSidePanels(types.RefreshOptions{Mode: types.BLOCK_UI})
			})
		},
	})
}

func (gui *Gui) refreshCommitFilesView() error {
	currentSideContext := gui.currentSideContext()
	if currentSideContext.GetKey() == COMMIT_FILES_CONTEXT_KEY || currentSideContext.GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
		if err := gui.handleRefreshPatchBuildingPanel(-1); err != nil {
			return err
		}
	}

	to := gui.State.Panels.CommitFiles.refName
	from, reverse := gui.getFromAndReverseArgsForDiff(to)

	files, err := gui.Git.Loaders.CommitFiles.GetFilesInDiff(from, to, reverse)
	if err != nil {
		return gui.PopupHandler.Error(err)
	}
	gui.State.CommitFileTreeViewModel.SetParent(to)
	gui.State.CommitFileTreeViewModel.SetFiles(files)

	return gui.postRefreshUpdate(gui.State.Contexts.CommitFiles)
}

func (gui *Gui) handleOpenOldCommitFile() error {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	return gui.openFile(node.GetPath())
}

func (gui *Gui) handleEditCommitFile() error {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return gui.PopupHandler.ErrorMsg(gui.Tr.ErrCannotEditDirectory)
	}

	return gui.editFile(node.GetPath())
}

func (gui *Gui) handleToggleFileForPatch() error {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	toggleTheFile := func() error {
		if !gui.Git.Patch.PatchManager.Active() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		// if there is any file that hasn't been fully added we'll fully add everything,
		// otherwise we'll remove everything
		adding := node.AnyFile(func(file *models.CommitFile) bool {
			return gui.Git.Patch.PatchManager.GetFileStatus(file.Name, gui.State.CommitFileTreeViewModel.GetParent()) != patch.WHOLE
		})

		err := node.ForEachFile(func(file *models.CommitFile) error {
			if adding {
				return gui.Git.Patch.PatchManager.AddFileWhole(file.Name)
			} else {
				return gui.Git.Patch.PatchManager.RemoveFile(file.Name)
			}
		})

		if err != nil {
			return gui.PopupHandler.Error(err)
		}

		if gui.Git.Patch.PatchManager.IsEmpty() {
			gui.Git.Patch.PatchManager.Reset()
		}

		return gui.postRefreshUpdate(gui.State.Contexts.CommitFiles)
	}

	if gui.Git.Patch.PatchManager.Active() && gui.Git.Patch.PatchManager.To != gui.State.CommitFileTreeViewModel.GetParent() {
		return gui.PopupHandler.Ask(popup.AskOpts{
			Title:  gui.Tr.DiscardPatch,
			Prompt: gui.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				gui.Git.Patch.PatchManager.Reset()
				return toggleTheFile()
			},
		})
	}

	return toggleTheFile()
}

func (gui *Gui) startPatchManager() error {
	canRebase := gui.State.Panels.CommitFiles.canRebase

	to := gui.State.Panels.CommitFiles.refName
	from, reverse := gui.getFromAndReverseArgsForDiff(to)

	gui.Git.Patch.PatchManager.Start(from, to, reverse, canRebase)
	return nil
}

func (gui *Gui) handleEnterCommitFile() error {
	return gui.enterCommitFile(OnFocusOpts{ClickedViewName: "", ClickedViewLineIdx: -1})
}

func (gui *Gui) enterCommitFile(opts OnFocusOpts) error {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return gui.handleToggleCommitFileDirCollapsed()
	}

	enterTheFile := func() error {
		if !gui.Git.Patch.PatchManager.Active() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		return gui.pushContext(gui.State.Contexts.PatchBuilding, opts)
	}

	if gui.Git.Patch.PatchManager.Active() && gui.Git.Patch.PatchManager.To != gui.State.CommitFileTreeViewModel.GetParent() {
		return gui.PopupHandler.Ask(popup.AskOpts{
			Title:  gui.Tr.DiscardPatch,
			Prompt: gui.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				gui.Git.Patch.PatchManager.Reset()
				return enterTheFile()
			},
		})
	}

	return enterTheFile()
}

func (gui *Gui) handleToggleCommitFileDirCollapsed() error {
	node := gui.getSelectedCommitFileNode()
	if node == nil {
		return nil
	}

	gui.State.CommitFileTreeViewModel.ToggleCollapsed(node.GetPath())

	if err := gui.postRefreshUpdate(gui.State.Contexts.CommitFiles); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) switchToCommitFilesContext(refName string, canRebase bool, context Context, windowName string) error {
	// sometimes the commitFiles view is already shown in another window, so we need to ensure that window
	// no longer considers the commitFiles view as its main view.
	gui.resetWindowForView(gui.Views.CommitFiles)

	gui.State.Panels.CommitFiles.SelectedLineIdx = 0
	gui.State.Panels.CommitFiles.refName = refName
	gui.State.Panels.CommitFiles.canRebase = canRebase
	gui.State.Contexts.CommitFiles.SetParentContext(context)
	gui.State.Contexts.CommitFiles.SetWindowName(windowName)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	return gui.pushContext(gui.State.Contexts.CommitFiles)
}

// NOTE: this is very similar to handleToggleFileTreeView, could be DRY'd with generics
func (gui *Gui) handleToggleCommitFileTreeView() error {
	path := gui.getSelectedCommitFilePath()

	gui.State.CommitFileTreeViewModel.ToggleShowTree()

	// find that same node in the new format and move the cursor to it
	if path != "" {
		gui.State.CommitFileTreeViewModel.ExpandToPath(path)
		index, found := gui.State.CommitFileTreeViewModel.GetIndexForPath(path)
		if found {
			gui.State.Contexts.CommitFiles.GetPanelState().SetSelectedLineIdx(index)
		}
	}

	if err := gui.State.Contexts.CommitFiles.HandleRender(); err != nil {
		return err
	}
	if err := gui.State.Contexts.CommitFiles.HandleFocus(); err != nil {
		return err
	}

	return nil
}
