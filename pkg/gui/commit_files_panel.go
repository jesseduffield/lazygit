package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) getSelectedCommitFile() *models.CommitFile {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}
	return node.File
}

func (gui *Gui) getSelectedCommitFilePath() string {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return ""
	}
	return node.GetPath()
}

// TODO: do we need this?
func (gui *Gui) onCommitFileFocus() error {
	gui.escapeLineByLinePanel()
	return nil
}

func (gui *Gui) commitFilesRenderToMain() error {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}

	to := gui.State.Contexts.CommitFiles.GetRefName()
	from, reverse := gui.getFromAndReverseArgsForDiff(to)

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

func (gui *Gui) handleCheckoutCommitFile() error {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}

	gui.c.LogAction(gui.c.Tr.Actions.CheckoutFile)
	if err := gui.git.WorkingTree.CheckoutFile(gui.State.Contexts.CommitFiles.GetRefName(), node.GetPath()); err != nil {
		return gui.c.Error(err)
	}

	return gui.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (gui *Gui) handleDiscardOldFileChange() error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	fileName := gui.getSelectedCommitFileName()

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.DiscardFileChangesTitle,
		Prompt: gui.c.Tr.DiscardFileChangesPrompt,
		HandleConfirm: func() error {
			return gui.c.WithWaitingStatus(gui.c.Tr.RebasingStatus, func() error {
				gui.c.LogAction(gui.c.Tr.Actions.DiscardOldFileChange)
				if err := gui.git.Rebase.DiscardOldFileChanges(gui.State.Model.Commits, gui.State.Contexts.LocalCommits.GetSelectedLineIdx(), fileName); err != nil {
					if err := gui.helpers.MergeAndRebase.CheckMergeOrRebase(err); err != nil {
						return err
					}
				}

				return gui.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI})
			})
		},
	})
}

func (gui *Gui) refreshCommitFilesView() error {
	currentSideContext := gui.currentSideContext()
	if currentSideContext.GetKey() == context.COMMIT_FILES_CONTEXT_KEY || currentSideContext.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
		if err := gui.handleRefreshPatchBuildingPanel(-1); err != nil {
			return err
		}
	}

	to := gui.State.Contexts.CommitFiles.GetRefName()
	from, reverse := gui.getFromAndReverseArgsForDiff(to)

	files, err := gui.git.Loaders.CommitFiles.GetFilesInDiff(from, to, reverse)
	if err != nil {
		return gui.c.Error(err)
	}
	gui.State.Model.CommitFiles = files
	gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.SetTree()

	return gui.c.PostRefreshUpdate(gui.State.Contexts.CommitFiles)
}

func (gui *Gui) handleOpenOldCommitFile() error {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}

	return gui.helpers.Files.OpenFile(node.GetPath())
}

func (gui *Gui) handleEditCommitFile() error {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return gui.c.ErrorMsg(gui.c.Tr.ErrCannotEditDirectory)
	}

	return gui.helpers.Files.EditFile(node.GetPath())
}

func (gui *Gui) handleToggleFileForPatch() error {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}

	toggleTheFile := func() error {
		if !gui.git.Patch.PatchManager.Active() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		// if there is any file that hasn't been fully added we'll fully add everything,
		// otherwise we'll remove everything
		adding := node.AnyFile(func(file *models.CommitFile) bool {
			return gui.git.Patch.PatchManager.GetFileStatus(file.Name, gui.State.Contexts.CommitFiles.GetRefName()) != patch.WHOLE
		})

		err := node.ForEachFile(func(file *models.CommitFile) error {
			if adding {
				return gui.git.Patch.PatchManager.AddFileWhole(file.Name)
			} else {
				return gui.git.Patch.PatchManager.RemoveFile(file.Name)
			}
		})

		if err != nil {
			return gui.c.Error(err)
		}

		if gui.git.Patch.PatchManager.IsEmpty() {
			gui.git.Patch.PatchManager.Reset()
		}

		return gui.c.PostRefreshUpdate(gui.State.Contexts.CommitFiles)
	}

	if gui.git.Patch.PatchManager.Active() && gui.git.Patch.PatchManager.To != gui.State.Contexts.CommitFiles.GetRefName() {
		return gui.c.Ask(types.AskOpts{
			Title:  gui.c.Tr.DiscardPatch,
			Prompt: gui.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				gui.git.Patch.PatchManager.Reset()
				return toggleTheFile()
			},
		})
	}

	return toggleTheFile()
}

func (gui *Gui) startPatchManager() error {
	commitFilesContext := gui.State.Contexts.CommitFiles

	canRebase := commitFilesContext.GetCanRebase()
	to := commitFilesContext.GetRefName()

	from, reverse := gui.getFromAndReverseArgsForDiff(to)

	gui.git.Patch.PatchManager.Start(from, to, reverse, canRebase)
	return nil
}

func (gui *Gui) handleEnterCommitFile() error {
	return gui.enterCommitFile(types.OnFocusOpts{ClickedViewName: "", ClickedViewLineIdx: -1})
}

func (gui *Gui) enterCommitFile(opts types.OnFocusOpts) error {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return gui.handleToggleCommitFileDirCollapsed()
	}

	enterTheFile := func() error {
		if !gui.git.Patch.PatchManager.Active() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		return gui.c.PushContext(gui.State.Contexts.PatchBuilding, opts)
	}

	if gui.git.Patch.PatchManager.Active() && gui.git.Patch.PatchManager.To != gui.State.Contexts.CommitFiles.GetRefName() {
		return gui.c.Ask(types.AskOpts{
			Title:  gui.c.Tr.DiscardPatch,
			Prompt: gui.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				gui.git.Patch.PatchManager.Reset()
				return enterTheFile()
			},
		})
	}

	return enterTheFile()
}

func (gui *Gui) handleToggleCommitFileDirCollapsed() error {
	node := gui.State.Contexts.CommitFiles.GetSelectedFileNode()
	if node == nil {
		return nil
	}

	gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.ToggleCollapsed(node.GetPath())

	if err := gui.c.PostRefreshUpdate(gui.State.Contexts.CommitFiles); err != nil {
		gui.c.Log.Error(err)
	}

	return nil
}

func (gui *Gui) SwitchToCommitFilesContext(opts controllers.SwitchToCommitFilesContextOpts) error {
	// sometimes the commitFiles view is already shown in another window, so we need to ensure that window
	// no longer considers the commitFiles view as its main view.
	gui.resetWindowContext(gui.State.Contexts.CommitFiles)

	gui.State.Contexts.CommitFiles.SetSelectedLineIdx(0)
	gui.State.Contexts.CommitFiles.SetRefName(opts.RefName)
	gui.State.Contexts.CommitFiles.SetCanRebase(opts.CanRebase)
	gui.State.Contexts.CommitFiles.SetParentContext(opts.Context)
	gui.State.Contexts.CommitFiles.SetWindowName(opts.Context.GetWindowName())

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	return gui.c.PushContext(gui.State.Contexts.CommitFiles)
}

// NOTE: this is very similar to handleToggleFileTreeView, could be DRY'd with generics
func (gui *Gui) handleToggleCommitFileTreeView() error {
	gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.ToggleShowTree()

	return gui.c.PostRefreshUpdate(gui.State.Contexts.CommitFiles)
}
