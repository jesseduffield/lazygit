package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

const (
	// these are the possible ref types for refs that you can view files of.
	// for local commits, we're allowed to build a patch and do things involving rebasing
	// with that patch
	REF_TYPE_LOCAL_COMMIT = iota

	// for other kinds of commits like reflog commits, we can't do anything rebasey
	REF_TYPE_OTHER_COMMIT

	// for stash entries we can't do anything rebasey, and the command for
	// obtaining the files is slightly different
	REF_TYPE_STASH
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

	isStash := gui.State.Panels.CommitFiles.refType == REF_TYPE_STASH
	files, err := gui.GitCommand.GetFilesInRef(gui.State.Panels.CommitFiles.refName, isStash, gui.GitCommand.PatchManager)
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

	canRebase := gui.State.Panels.CommitFiles.refType == REF_TYPE_LOCAL_COMMIT
	gui.GitCommand.PatchManager.Start(gui.State.Panels.CommitFiles.refName, diffMap, canRebase)
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

func (gui *Gui) switchToCommitFilesContext(refName string, refType int, context Context, windowName string) error {
	// sometimes the commitFiles view is already shown in another window, so we need to ensure that window
	// no longer considers the commitFiles view as its main view.
	gui.resetWindowForView("commitFiles")

	gui.State.Panels.CommitFiles.SelectedLineIdx = 0
	gui.State.Panels.CommitFiles.refName = refName
	gui.State.Panels.CommitFiles.refType = refType
	gui.Contexts.CommitFiles.Context.SetParentContext(context)
	gui.Contexts.CommitFiles.Context.SetWindowName(windowName)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	return gui.switchContext(gui.Contexts.CommitFiles.Context)
}
