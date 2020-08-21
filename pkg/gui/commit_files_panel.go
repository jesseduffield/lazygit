package gui

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) getSelectedCommitFile() *commands.CommitFile {
	selectedLine := gui.State.Panels.CommitFiles.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.CommitFiles[selectedLine]
}

func (gui *Gui) handleCommitFileSelect() error {
	gui.handleEscapeLineByLinePanel()

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		// TODO: consider making it so that we can also render strings to our own view through some common interface, or just render this to the main view for consistency
		gui.renderString("commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
		return nil
	}

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.ShowCommitFileCmdStr(commitFile.Parent, commitFile.Name, false),
	)
	task := gui.createRunPtyTask(cmd)

	return gui.refreshMain(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Patch",
			task:  task,
		},
		secondary: gui.secondaryPatchPanelUpdateOpts(),
	})
}

func (gui *Gui) handleCheckoutCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLineIdx]

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
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("DiscardFileChangesTitle"),
		prompt:             gui.Tr.SLocalize("DiscardFileChangesPrompt"),
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

	files, err := gui.GitCommand.GetFilesInRef(gui.State.Panels.CommitFiles.refName, gui.State.Panels.CommitFiles.isStash, gui.GitCommand.PatchManager)
	if err != nil {
		return gui.surfaceError(err)
	}
	gui.State.CommitFiles = files

	gui.Log.Warn(spew.Sdump(files))

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
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		gui.renderString("commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
		return nil
	}

	toggleTheFile := func() error {
		if !gui.GitCommand.PatchManager.Active() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		gui.GitCommand.PatchManager.ToggleFileWhole(commitFile.Name)

		return gui.refreshCommitFilesView()
	}

	if gui.GitCommand.PatchManager.Active() && gui.GitCommand.PatchManager.Parent != commitFile.Parent {
		return gui.ask(askOpts{
			returnToView:       v,
			returnFocusOnClose: true,
			title:              gui.Tr.SLocalize("DiscardPatch"),
			prompt:             gui.Tr.SLocalize("DiscardPatchConfirm"),
			handleConfirm: func() error {
				gui.GitCommand.PatchManager.Reset()
				return toggleTheFile()
			},
		})
	}

	return toggleTheFile()
}

func (gui *Gui) startPatchManager() error {
	diffMap := map[string]string{}
	// TODO: only load these files as we need to
	for _, commitFile := range gui.State.CommitFiles {
		commitText, err := gui.GitCommand.ShowCommitFile(commitFile.Parent, commitFile.Name, true)
		if err != nil {
			return err
		}
		diffMap[commitFile.Name] = commitText
	}

	gui.GitCommand.PatchManager.Start(gui.State.Panels.CommitFiles.refName, diffMap)
	return nil
}

func (gui *Gui) handleEnterCommitFile(g *gocui.Gui, v *gocui.View) error {
	return gui.enterCommitFile(-1)
}

func (gui *Gui) enterCommitFile(selectedLineIdx int) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		gui.renderString("commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
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

	if gui.GitCommand.PatchManager.Active() && gui.GitCommand.PatchManager.Parent != commitFile.Parent {
		return gui.ask(askOpts{
			returnToView:       gui.getCommitFilesView(),
			returnFocusOnClose: false,
			title:              gui.Tr.SLocalize("DiscardPatch"),
			prompt:             gui.Tr.SLocalize("DiscardPatchConfirm"),
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

func (gui *Gui) switchToCommitFilesContext(refName string, isStash bool, context Context, windowName string) error {
	// sometimes the commitFiles view is already shown in another window, so we need to ensure that window
	// no longer considers the commitFiles view as its main view.
	gui.resetWindowForView("commitFiles")

	gui.State.Panels.CommitFiles.SelectedLineIdx = 0
	gui.State.Panels.CommitFiles.refName = refName
	gui.State.Panels.CommitFiles.isStash = isStash
	gui.Contexts.CommitFiles.Context.SetParentContext(context)
	gui.Contexts.CommitFiles.Context.SetWindowName(windowName)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	return gui.switchContext(gui.Contexts.CommitFiles.Context)
}
