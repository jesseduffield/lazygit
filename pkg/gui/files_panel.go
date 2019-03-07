package gui

import (

	// "io"
	// "io/ioutil"

	// "strings"

	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedFile(g *gocui.Gui) (*commands.File, error) {
	selectedLine := gui.State.Panels.Files.SelectedLine
	if selectedLine == -1 {
		return &commands.File{}, gui.Errors.ErrNoFiles
	}

	return gui.State.Files[selectedLine], nil
}

func (gui *Gui) handleFilesFocus(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	cx, cy := v.Cursor()
	_, oy := v.Origin()

	prevSelectedLine := gui.State.Panels.Files.SelectedLine
	newSelectedLine := cy - oy

	if newSelectedLine > len(gui.State.Files)-1 || len(utils.Decolorise(gui.State.Files[newSelectedLine].DisplayString)) < cx {
		return gui.handleFileSelect(gui.g, v, false)
	}

	gui.State.Panels.Files.SelectedLine = newSelectedLine

	if prevSelectedLine == newSelectedLine && gui.currentViewName() == v.Name() {
		return gui.handleFilePress(gui.g, v)
	} else {
		return gui.handleFileSelect(gui.g, v, true)
	}
}

func (gui *Gui) handleFileSelect(g *gocui.Gui, v *gocui.View, alreadySelected bool) error {
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoChangedFiles"))
	}

	if err := gui.focusPoint(0, gui.State.Panels.Files.SelectedLine, v); err != nil {
		return err
	}

	if file.HasInlineMergeConflicts {
		return gui.refreshMergePanel()
	}

	content := gui.GitCommand.Diff(file, false)
	if alreadySelected {
		g.Update(func(*gocui.Gui) error {
			return gui.setViewContent(gui.g, gui.getMainView(), content)
		})
		return nil
	}
	return gui.renderString(g, "main", content)
}

func (gui *Gui) refreshFiles() error {
	selectedFile, _ := gui.getSelectedFile(gui.g)

	filesView := gui.getFilesView()
	if err := gui.refreshStateFiles(); err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {

		filesView.Clear()
		isFocused := gui.g.CurrentView().Name() == "files"
		list, err := utils.RenderList(gui.State.Files, isFocused)
		if err != nil {
			return err
		}
		fmt.Fprint(filesView, list)

		if filesView == g.CurrentView() {
			newSelectedFile, _ := gui.getSelectedFile(gui.g)
			alreadySelected := newSelectedFile.Name == selectedFile.Name
			return gui.handleFileSelect(g, filesView, alreadySelected)
		}
		return nil
	})

	return nil
}

func (gui *Gui) handleFilesNextLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Files
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Files), false)

	return gui.handleFileSelect(gui.g, v, false)
}

func (gui *Gui) handleFilesPrevLine(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	panelState := gui.State.Panels.Files
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Files), true)

	return gui.handleFileSelect(gui.g, v, false)
}

// specific functions

func (gui *Gui) stagedFiles() []*commands.File {
	files := gui.State.Files
	result := make([]*commands.File, 0)
	for _, file := range files {
		if file.HasStagedChanges {
			result = append(result, file)
		}
	}
	return result
}

func (gui *Gui) trackedFiles() []*commands.File {
	files := gui.State.Files
	result := make([]*commands.File, 0)
	for _, file := range files {
		if file.Tracked {
			result = append(result, file)
		}
	}
	return result
}

func (gui *Gui) stageSelectedFile(g *gocui.Gui) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}
	return gui.GitCommand.StageFile(file.Name)
}

func (gui *Gui) handleEnterFile(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}
	if file.HasInlineMergeConflicts {
		return gui.handleSwitchToMerge(g, v)
	}
	if !file.HasUnstagedChanges || file.HasMergeConflicts {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("FileStagingRequirements"))
	}
	if err := gui.changeContext("main", "staging"); err != nil {
		return err
	}
	if err := gui.switchFocus(g, v, gui.getMainView()); err != nil {
		return err
	}
	return gui.refreshStagingPanel()
}

func (gui *Gui) handleFilePress(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
		return err
	}

	if file.HasInlineMergeConflicts {
		return gui.handleSwitchToMerge(g, v)
	}

	if file.HasUnstagedChanges {
		gui.GitCommand.StageFile(file.Name)
	} else {
		gui.GitCommand.UnStageFile(file.Name, file.Tracked)
	}

	if err := gui.refreshFiles(); err != nil {
		return err
	}

	return gui.handleFileSelect(g, v, true)
}

func (gui *Gui) allFilesStaged() bool {
	for _, file := range gui.State.Files {
		if file.HasUnstagedChanges {
			return false
		}
	}
	return true
}

func (gui *Gui) handleStageAll(g *gocui.Gui, v *gocui.View) error {
	var err error
	if gui.allFilesStaged() {
		err = gui.GitCommand.UnstageAll()
	} else {
		err = gui.GitCommand.StageAll()
	}
	if err != nil {
		_ = gui.createErrorPanel(g, err.Error())
	}

	if err := gui.refreshFiles(); err != nil {
		return err
	}

	return gui.handleFileSelect(g, v, false)
}

func (gui *Gui) handleAddPatch(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
		return err
	}
	if !file.HasUnstagedChanges {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("FileHasNoUnstagedChanges"))
	}
	if !file.Tracked {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CannotGitAdd"))
	}

	gui.SubProcess = gui.GitCommand.AddPatch(file.Name)
	return gui.Errors.ErrSubProcess
}

func (gui *Gui) handleFileRemove(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
		return err
	}
	var deleteVerb string
	if file.Tracked {
		deleteVerb = gui.Tr.SLocalize("checkout")
	} else {
		deleteVerb = gui.Tr.SLocalize("delete")
	}
	message := gui.Tr.TemplateLocalize(
		"SureTo",
		Teml{
			"deleteVerb": deleteVerb,
			"fileName":   file.Name,
		},
	)
	return gui.createConfirmationPanel(g, v, strings.Title(deleteVerb)+" file", message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.RemoveFile(file); err != nil {
			return err
		}
		return gui.refreshFiles()
	}, nil)
}

func (gui *Gui) handleIgnoreFile(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	if file.Tracked {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantIgnoreTrackFiles"))
	}
	if err := gui.GitCommand.Ignore(file.Name); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	return gui.refreshFiles()
}

func (gui *Gui) handleCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && gui.State.WorkingTreeState == "normal" {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("NoStagedFilesToCommit"))
	}
	commitMessageView := gui.getCommitMessageView()
	g.Update(func(g *gocui.Gui) error {
		g.SetViewOnTop("commitMessage")
		gui.switchFocus(g, filesView, commitMessageView)
		gui.RenderCommitLength()
		return nil
	})
	return nil
}

func (gui *Gui) handleAmendCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && gui.State.WorkingTreeState == "normal" {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("NoStagedFilesToCommit"))
	}
	if len(gui.State.Commits) == 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("NoCommitToAmend"))
	}

	title := strings.Title(gui.Tr.SLocalize("AmendLastCommit"))
	question := gui.Tr.SLocalize("SureToAmend")

	return gui.createConfirmationPanel(g, filesView, title, question, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.runSyncOrAsyncCommand(gui.GitCommand.AmendHead()); err != nil {
			return err
		}

		return gui.refreshSidePanels(g)
	}, nil)
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && gui.State.WorkingTreeState == "normal" {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("NoStagedFilesToCommit"))
	}
	gui.PrepareSubProcess(g, "git", "commit")
	return nil
}

// PrepareSubProcess - prepare a subprocess for execution and tell the gui to switch to it
func (gui *Gui) PrepareSubProcess(g *gocui.Gui, commands ...string) {
	gui.SubProcess = gui.GitCommand.PrepareCommitSubProcess()
	g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})
}

func (gui *Gui) editFile(filename string) error {
	return gui.runSyncOrAsyncCommand(gui.OSCommand.EditFile(filename))
}

func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}

	return gui.editFile(file.Name)
}

func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}
	return gui.openFile(file.Name)
}

func (gui *Gui) handleRefreshFiles(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshFiles()
}

func (gui *Gui) refreshStateFiles() error {
	// get files to stage
	files := gui.GitCommand.GetStatusFiles()
	gui.State.Files = gui.GitCommand.MergeStatusFiles(gui.State.Files, files)
	gui.refreshSelectedLine(&gui.State.Panels.Files.SelectedLine, len(gui.State.Files))
	return gui.updateWorkTreeState()
}

func (gui *Gui) catSelectedFile(g *gocui.Gui) (string, error) {
	item, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return "", err
		}
		return "", gui.renderString(g, "main", gui.Tr.SLocalize("NoFilesDisplay"))
	}
	if item.Type != "file" {
		return "", gui.renderString(g, "main", gui.Tr.SLocalize("NotAFile"))
	}
	cat, err := gui.GitCommand.CatFile(item.Name)
	if err != nil {
		gui.Log.Error(err)
		return "", gui.renderString(g, "main", err.Error())
	}
	return cat, nil
}

func (gui *Gui) pullFiles(g *gocui.Gui, v *gocui.View) error {
	if err := gui.createLoaderPanel(gui.g, v, gui.Tr.SLocalize("PullWait")); err != nil {
		return err
	}

	go func() {
		unamePassOpend := false
		err := gui.GitCommand.Pull(func(passOrUname string) string {
			unamePassOpend = true
			return gui.waitForPassUname(g, v, passOrUname)
		})
		gui.HandleCredentialsPopup(g, unamePassOpend, err)
	}()
	return nil
}

func (gui *Gui) pushWithForceFlag(g *gocui.Gui, v *gocui.View, force bool) error {
	if err := gui.createLoaderPanel(gui.g, v, gui.Tr.SLocalize("PushWait")); err != nil {
		return err
	}
	go func() {
		unamePassOpend := false
		branchName := gui.State.Branches[0].Name
		err := gui.GitCommand.Push(branchName, force, func(passOrUname string) string {
			unamePassOpend = true
			return gui.waitForPassUname(g, v, passOrUname)
		})
		gui.HandleCredentialsPopup(g, unamePassOpend, err)
	}()
	return nil
}

func (gui *Gui) pushFiles(g *gocui.Gui, v *gocui.View) error {
	// if we have pullables we'll ask if the user wants to force push
	_, pullables := gui.GitCommand.GetCurrentBranchUpstreamDifferenceCount()
	if pullables == "?" || pullables == "0" {
		return gui.pushWithForceFlag(g, v, false)
	}
	err := gui.createConfirmationPanel(g, nil, gui.Tr.SLocalize("ForcePush"), gui.Tr.SLocalize("ForcePushPrompt"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.pushWithForceFlag(g, v, true)
	}, nil)
	return err
}

func (gui *Gui) handleSwitchToMerge(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}
	if !file.HasInlineMergeConflicts {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("FileNoMergeCons"))
	}
	if err := gui.changeContext("main", "merging"); err != nil {
		return err
	}
	if err := gui.switchFocus(g, v, gui.getMainView()); err != nil {
		return err
	}
	return gui.refreshMergePanel()
}

func (gui *Gui) handleAbortMerge(g *gocui.Gui, v *gocui.View) error {
	if err := gui.GitCommand.AbortMerge(); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	gui.createMessagePanel(g, v, "", gui.Tr.SLocalize("MergeAborted"))
	gui.refreshStatus(g)
	return gui.refreshFiles()
}

func (gui *Gui) handleResetAndClean(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("ClearFilePanel"), gui.Tr.SLocalize("SureResetHardHead"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.ResetAndClean(); err != nil {
			gui.createErrorPanel(g, err.Error())
		}
		return gui.refreshFiles()
	}, nil)
}

func (gui *Gui) openFile(filename string) error {
	if err := gui.OSCommand.OpenFile(filename); err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}
	return nil
}

func (gui *Gui) anyFilesWithMergeConflicts() bool {
	for _, file := range gui.State.Files {
		if file.HasMergeConflicts {
			return true
		}
	}
	return false
}

func (gui *Gui) handleSoftReset(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("SoftReset"), gui.Tr.SLocalize("ConfirmSoftReset"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.SoftReset("HEAD^"); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}

		if err := gui.refreshCommits(gui.g); err != nil {
			return err
		}
		return gui.refreshFiles()
	}, nil)
}
