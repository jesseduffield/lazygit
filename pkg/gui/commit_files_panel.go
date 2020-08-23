package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) getSelectedCommitFile() *commands.CommitFile {
	selectedLine := gui.State.Panels.CommitFiles.SelectedLineIdx
	if selectedLine == -1 || selectedLine > len(gui.State.CommitFiles)-1 {
		return nil
	}

	return gui.State.CommitFiles[selectedLine]
}

func (gui *Gui) handleCommitFileSelect() error {
	gui.handleEscapeLineByLinePanel()

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		return nil
	}

	to := commitFile.Parent
	from, reverse := gui.getFromAndReverseArgsForDiff(to)

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.ShowFileDiffCmdStr(from, to, reverse, commitFile.Name, false),
	)
	task := gui.createRunPtyTask(cmd)

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Patch",
			task:  task,
		},
		secondary: gui.secondaryPatchPanelUpdateOpts(),
	})
}

func (gui *Gui) handleCheckoutCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.getSelectedCommitFile()
	if file == nil {
		return nil
	}

	if err := gui.GitCommand.CheckoutFile(file.Parent, file.Name); err != nil {
		return gui.surfaceError(err)
	}

	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleDiscardOldFileChange(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	fileName := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLineIdx].Name

	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("DiscardFileChangesTitle"),
		prompt: gui.Tr.SLocalize("DiscardFileChangesPrompt"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
				if err := gui.GitCommand.DiscardOldFileChanges(gui.State.Commits, gui.State.Panels.Commits.SelectedLineIdx, fileName); err != nil {
					if err := gui.handleGenericMergeCommandResult(err); err != nil {
						return err
					}
				}

				return gui.refreshSidePanels(refreshOptions{mode: BLOCK_UI})
			})
		},
	})
}

func (gui *Gui) refreshCommitFilesView() error {
	if err := gui.refreshPatchBuildingPanel(-1); err != nil {
		return err
	}

	to := gui.State.Panels.CommitFiles.refName
	from, reverse := gui.getFromAndReverseArgsForDiff(to)

	files, err := gui.GitCommand.GetFilesInDiff(from, to, reverse, gui.GitCommand.PatchManager)
	if err != nil {
		return gui.surfaceError(err)
	}
	gui.State.CommitFiles = files

	return gui.postRefreshUpdate(gui.Contexts.CommitFiles.Context)
}

func (gui *Gui) handleOpenOldCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.getSelectedCommitFile()
	if file == nil {
		return nil
	}

	return gui.openFile(file.Name)
}

func (gui *Gui) handleEditCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.getSelectedCommitFile()
	if file == nil {
		return nil
	}

	return gui.editFile(file.Name)
}

func (gui *Gui) handleToggleFileForPatch(g *gocui.Gui, v *gocui.View) error {
	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		return nil
	}

	toggleTheFile := func() error {
		if !gui.GitCommand.PatchManager.Active() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		if err := gui.GitCommand.PatchManager.ToggleFileWhole(commitFile.Name); err != nil {
			return err
		}

		if gui.GitCommand.PatchManager.IsEmpty() {
			gui.GitCommand.PatchManager.Reset()
		}

		return gui.refreshCommitFilesView()
	}

	if gui.GitCommand.PatchManager.Active() && gui.GitCommand.PatchManager.To != commitFile.Parent {
		return gui.ask(askOpts{
			title:  gui.Tr.SLocalize("DiscardPatch"),
			prompt: gui.Tr.SLocalize("DiscardPatchConfirm"),
			handleConfirm: func() error {
				gui.GitCommand.PatchManager.Reset()
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

	gui.GitCommand.PatchManager.Start(from, to, reverse, canRebase)
	return nil
}

func (gui *Gui) handleEnterCommitFile(g *gocui.Gui, v *gocui.View) error {
	return gui.enterCommitFile(-1)
}

func (gui *Gui) enterCommitFile(selectedLineIdx int) error {
	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		return nil
	}

	enterTheFile := func(selectedLineIdx int) error {
		if !gui.GitCommand.PatchManager.Active() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		if err := gui.switchContext(gui.Contexts.PatchBuilding.Context); err != nil {
			return err
		}
		return gui.refreshPatchBuildingPanel(selectedLineIdx)
	}

	if gui.GitCommand.PatchManager.Active() && gui.GitCommand.PatchManager.To != commitFile.Parent {
		return gui.ask(askOpts{
			title:               gui.Tr.SLocalize("DiscardPatch"),
			prompt:              gui.Tr.SLocalize("DiscardPatchConfirm"),
			handlersManageFocus: true,
			handleConfirm: func() error {
				gui.GitCommand.PatchManager.Reset()
				return enterTheFile(selectedLineIdx)
			},
			handleClose: func() error {
				return gui.switchContext(gui.Contexts.CommitFiles.Context)
			},
		})
	}

	return enterTheFile(selectedLineIdx)
}

func (gui *Gui) switchToCommitFilesContext(refName string, canRebase bool, context Context, windowName string) error {
	// sometimes the commitFiles view is already shown in another window, so we need to ensure that window
	// no longer considers the commitFiles view as its main view.
	gui.resetWindowForView("commitFiles")

	gui.State.Panels.CommitFiles.SelectedLineIdx = 0
	gui.State.Panels.CommitFiles.refName = refName
	gui.State.Panels.CommitFiles.canRebase = canRebase
	gui.Contexts.CommitFiles.Context.SetParentContext(context)
	gui.Contexts.CommitFiles.Context.SetWindowName(windowName)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	return gui.switchContext(gui.Contexts.CommitFiles.Context)
}
