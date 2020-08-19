package gui

import (
	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

func (gui *Gui) getSelectedCommitFile() *commands.CommitFile {
	selectedLine := gui.State.Panels.CommitFiles.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.CommitFiles[selectedLine]
}

func (gui *Gui) handleCommitFileSelect() error {
	if gui.popupPanelFocused() {
		return nil
	}

	if gui.currentViewName() == "commitFiles" {
		gui.handleEscapeLineByLinePanel()
	}

	commitFile := gui.getSelectedCommitFile()
	if commitFile == nil {
		// TODO: consider making it so that we can also render strings to our own view through some common interface, or just render this to the main view for consistency
		gui.renderString("commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
		return nil
	}

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.ShowCommitFileCmdStr(commitFile.Sha, commitFile.Name, false),
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

func (gui *Gui) handleSwitchToCommitsPanel(g *gocui.Gui, v *gocui.View) error {
	return gui.switchContext(gui.Contexts.BranchCommits.Context)
}

func (gui *Gui) handleCheckoutCommitFile(g *gocui.Gui, v *gocui.View) error {
	file := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine]

	if err := gui.GitCommand.CheckoutFile(file.Sha, file.Name); err != nil {
		return gui.surfaceError(err)
	}

	return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
}

func (gui *Gui) handleDiscardOldFileChange(g *gocui.Gui, v *gocui.View) error {
	if ok, err := gui.validateNormalWorkingTreeState(); !ok {
		return err
	}

	fileName := gui.State.CommitFiles[gui.State.Panels.CommitFiles.SelectedLine].Name

	return gui.ask(askOpts{
		returnToView:       v,
		returnFocusOnClose: true,
		title:              gui.Tr.SLocalize("DiscardFileChangesTitle"),
		prompt:             gui.Tr.SLocalize("DiscardFileChangesPrompt"),
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("RebasingStatus"), func() error {
				if err := gui.GitCommand.DiscardOldFileChanges(gui.State.Commits, gui.State.Panels.Commits.SelectedLine, fileName); err != nil {
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

	commit := gui.getSelectedCommit()
	if commit == nil {
		return nil
	}

	files, err := gui.GitCommand.GetCommitFiles(commit.Sha, gui.GitCommand.PatchManager)
	if err != nil {
		return gui.surfaceError(err)
	}
	gui.State.CommitFiles = files

	return gui.postRefreshUpdate(gui.Contexts.BranchCommits.Files.Context)
}

func (gui *Gui) renderCommitFiles() error {
	gui.refreshSelectedLine(&gui.State.Panels.CommitFiles.SelectedLine, len(gui.State.CommitFiles))

	commitsFileView := gui.getCommitFilesView()
	displayStrings := presentation.GetCommitFileListDisplayStrings(gui.State.CommitFiles, gui.State.Diff.Ref)
	gui.renderDisplayStrings(commitsFileView, displayStrings)

	return nil
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
		if !gui.GitCommand.PatchManager.CommitSelected() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		gui.GitCommand.PatchManager.ToggleFileWhole(commitFile.Name)

		return gui.refreshCommitFilesView()
	}

	if gui.GitCommand.PatchManager.CommitSelected() && gui.GitCommand.PatchManager.CommitSha != commitFile.Sha {
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
	for _, commitFile := range gui.State.CommitFiles {
		commitText, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name, true)
		if err != nil {
			return err
		}
		diffMap[commitFile.Name] = commitText
	}

	commit := gui.getSelectedCommit()
	if commit == nil {
		return errors.New("No commit selected")
	}

	gui.GitCommand.PatchManager.Start(commit.Sha, diffMap)
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
		if !gui.GitCommand.PatchManager.CommitSelected() {
			if err := gui.startPatchManager(); err != nil {
				return err
			}
		}

		if err := gui.switchContext(gui.Contexts.PatchBuilding.Context); err != nil {
			return err
		}
		return gui.refreshPatchBuildingPanel(selectedLineIdx)
	}

	if gui.GitCommand.PatchManager.CommitSelected() && gui.GitCommand.PatchManager.CommitSha != commitFile.Sha {
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
				return gui.switchContext(gui.Contexts.BranchCommits.Files.Context)
			},
		})
	}

	return enterTheFile(selectedLineIdx)
}
