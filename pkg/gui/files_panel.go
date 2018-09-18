package gui

import (
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// refreshFiles refreshes the files view.
// returns an error if something goes wrong.
func (gui *Gui) refreshFiles() error {
	filesView, err := gui.g.View("files")
	if err != nil {
		return err
	}

	if err = gui.refreshStateFiles(); err != nil {
		return err
	}

	filesView.Clear()

	for _, file := range gui.State.Files {
		gui.renderFileName(file, filesView)
	}

	if err = gui.correctCursor(filesView); err != nil {
		return err
	}

	if filesView == gui.g.CurrentView() {
		if err = gui.handleFileSelect(); err != nil {
			return err
		}
	}

	return nil
}

// refreshStateFiles refreshes the state files.
// returns an error if something goes wrong.
func (gui *Gui) refreshStateFiles() error {
	files := gui.GitCommand.GetStatusFiles()
	gui.State.Files = gui.GitCommand.MergeStatusFiles(gui.State.Files, files)

	if err := gui.updateHasMergeConflictStatus(); err != nil {
		return err
	}

	return nil
}

// stagedFiles returns the staged files.
// returns a file array.
func (gui *Gui) stagedFiles() []commands.File {
	files := gui.State.Files
	result := make([]commands.File, 0)

	for _, file := range files {
		if file.HasStagedChanges {
			result = append(result, file)
		}
	}

	return result
}

// trackedFiles returns the tracked files.
// returns a file array.
func (gui *Gui) trackedFiles() []commands.File {
	files := gui.State.Files
	result := make([]commands.File, 0)

	for _, file := range files {
		if file.Tracked {
			result = append(result, file)
		}
	}

	return result
}

// stageSelectedFile stages the selected file.
// returns an error if something went wrong.
func (gui *Gui) stageSelectedFile() error {
	file, err := gui.getSelectedFile()
	if err != nil {
		return err
	}

	if err := gui.GitCommand.StageFile(file.Name); err != nil {
		return err
	}

	return nil
}

// handleFilePress is called when the user selects a file.
// g and v are passed by the gocui library to the function but are not used.
// In case something goes wrong, this function returns an error.
func (gui *Gui) handleFilePress(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}

	if file.HasMergeConflicts {
		return gui.handleSwitchToMerge(gui.g, v)
	}

	if file.HasUnstagedChanges {
		err = gui.GitCommand.StageFile(file.Name)
	} else {
		err = gui.GitCommand.UnStageFile(file.Name, file.Tracked)
	}
	if err != nil {
		return err
	}

	if err = gui.refreshFiles(); err != nil {
		return err
	}

	return gui.handleFileSelect()
}

// allFilesStage returns whether or not all the files are staged.
func (gui *Gui) allFilesStaged() bool {
	for _, file := range gui.State.Files {
		if file.HasUnstagedChanges {
			return false
		}
	}

	return true
}

// handleStageAll is called when the user pressed the stage all key
// in the gui.
// g and v are passed by the gocui library bubt are not used.
// In case something goes wrong, it returns an error.
func (gui *Gui) handleStageAll(g *gocui.Gui, v *gocui.View) error {
	var err error
	if gui.allFilesStaged() {
		err = gui.GitCommand.UnstageAll()
	} else {
		err = gui.GitCommand.StageAll()
	}
	if err != nil {
		return gui.createErrorPanel(err.Error())
	}

	if err = gui.refreshFiles(); err != nil {
		return err
	}

	return gui.handleFileSelect()
}

// handleAddPatch is called when a user wants to add a patch.
// g and v are passed by the gocui library
func (gui *Gui) handleAddPatch(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}

	if !file.HasUnstagedChanges {
		return gui.createErrorPanel(gui.Tr.SLocalize("FileHasNoUnstagedChanges"))
	}

	if !file.Tracked {
		return gui.createErrorPanel(gui.Tr.SLocalize("CannotGitAdd"))
	}

	gui.SubProcess = gui.GitCommand.AddPatch(file.Name)
	return gui.Errors.ErrSubProcess
}

// getSelectedFile returns the selected files
// returns the file and an error if something goes wrong.
func (gui *Gui) getSelectedFile() (commands.File, error) {
	if len(gui.State.Files) == 0 {
		return commands.File{}, gui.Errors.ErrNoFiles
	}

	filesView, err := gui.g.View("files")
	if err != nil {
		return commands.File{}, err
	}

	lineNumber := gui.getItemPosition(filesView)
	return gui.State.Files[lineNumber], nil
}

// handleFileRemoved gets called when a file is removed.
// g and v are passed by the gocui library.
// returns an error if something goes wrong
func (gui *Gui) handleFileRemove(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
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

	return gui.createConfirmationPanel(v, strings.Title(deleteVerb)+" file", message,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := gui.GitCommand.RemoveFile(file); err != nil {
				return err
			}

			return gui.refreshFiles()
		}, nil)
}

// handleIgnoreFile handles ignoring the file.
// g and v are passed by the gocui.
// returns an error if something goes wrong.
func (gui *Gui) handleIgnoreFile(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		return gui.createErrorPanel(err.Error())
	}

	if file.Tracked {
		return gui.createErrorPanel(gui.Tr.SLocalize("CantIgnoreTrackFiles"))
	}

	if err := gui.GitCommand.Ignore(file.Name); err != nil {
		return gui.createErrorPanel(err.Error())
	}

	return gui.refreshFiles()
}

// handleFileSelect is called when a file is selected.
// It checks if there are any changed files and if there is one
// and it is selected, it gets rendered into the main view
func (gui *Gui) handleFileSelect() error {
	if err := gui.renderGlobalOptions(); err != nil {
		return err
	}

	file, err := gui.getSelectedFile()
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return gui.renderString("main", gui.Tr.SLocalize("NoChangedFiles"))
	}

	if file.HasMergeConflicts {
		return gui.refreshMergePanel()
	}

	content := gui.GitCommand.Diff(file)
	return gui.renderString("main", content)
}

// handleCommitPress is called when a user commits changes.
// g and v are passed by the gocui library.
// returns and error if something goes wrong.
func (gui *Gui) handleCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && !gui.State.HasMergeConflicts {
		return gui.createErrorPanel(gui.Tr.SLocalize("NoStagedFilesToCommit"))
	}

	commitMessageView, err := gui.g.View("commitMessage")
	if err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		if _, err = gui.g.SetViewOnTop("commitMessage"); err != nil {
			return err
		}

		if err = gui.switchFocus(filesView, commitMessageView); err != nil {
			return err
		}

		gui.RenderCommitLength()
		return nil
	})

	return nil
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && !gui.State.HasMergeConflicts {
		return gui.createErrorPanel(gui.Tr.SLocalize("NoStagedFilesToCommit"))
	}

	gui.PrepareSubProcess(gui.g, "git", "commit")
	return nil
}

// PrepareSubProcess - prepare a subprocess for execution and tell the gui to switch to it
func (gui *Gui) PrepareSubProcess(g *gocui.Gui, commands ...string) {
	gui.SubProcess = gui.GitCommand.PrepareCommitSubProcess()
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})
}

// editFile edits a file.
// This is achieved with a subprocess
func (gui *Gui) editFile(filename string) error {
	sub, err := gui.OSCommand.EditFile(filename)
	if err != nil {
		return gui.createErrorPanel(err.Error())
	}

	if sub == nil {
		return nil
	}

	gui.SubProcess = sub
	return gui.Errors.ErrSubProcess
}

// handleFileEdit is called when a user wants to edit a file
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		return err
	}

	return gui.editFile(file.Name)
}

// handleFileOpen is called when a user opens a file
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		return err
	}

	return gui.openFile(file.Name)
}

// handleRefreshFiles is a macro for the keybindings
func (gui *Gui) handleRefreshFiles(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshFiles()
}

// updateHasMergeConflictStatus updates the status.
// returns an error if something goes wrong.
func (gui *Gui) updateHasMergeConflictStatus() error {
	merging, err := gui.GitCommand.IsInMergeState()
	if err != nil {
		return err
	}

	gui.State.HasMergeConflicts = merging
	return nil
}

// renderFileName renders a filename.
// file is the file to render.
// filesview is where the view is located.
func (gui *Gui) renderFileName(file commands.File, filesView *gocui.View) {
	// potentially inefficient to be instantiating these color
	// objects with each render
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)

	if !file.Tracked && !file.HasStagedChanges {
		red.Fprintln(filesView, file.DisplayString)
		return
	}

	green.Fprint(filesView, file.DisplayString[0:1])
	red.Fprint(filesView, file.DisplayString[1:3])

	if file.HasUnstagedChanges {
		red.Fprintln(filesView, file.Name)
	} else {
		green.Fprintln(filesView, file.Name)
	}
}

// catSelectedFiles reads the file and "cats" the data.
// returns the string and an error if anything goes wrong.
func (gui *Gui) catSelectedFile() (string, error) {
	item, err := gui.getSelectedFile()
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return "", err
		}
		return "", gui.renderString("main", gui.Tr.SLocalize("NoFilesDisplay"))
	}

	if item.Type != "file" {
		return "", gui.renderString("main", gui.Tr.SLocalize("NotAFile"))
	}

	cat, err := gui.GitCommand.CatFile(item.Name)
	if err != nil {
		return "", gui.renderString("main", err.Error())
	}

	return cat, nil
}

// pullFiles pulls the files.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) pullFiles(g *gocui.Gui, v *gocui.View) error {
	if err := gui.createMessagePanel(v, "", gui.Tr.SLocalize("PullWait")); err != nil {
		return err
	}

	go func() {
		if err := gui.GitCommand.Pull(); err != nil {
			_ = gui.createErrorPanel(err.Error())
		} else {
			_ = gui.closeConfirmationPrompt()
			_ = gui.refreshCommits()
			_ = gui.refreshStatus()
		}
		_ = gui.refreshFiles()
	}()

	return nil
}

// pushWithForceFlag does what it says, it pushes with or without force.
// currentView is used to return focus to the view.
// force indicates whether or not to push with foce.
func (gui *Gui) pushWithForceFlag(currentView *gocui.View, force bool) error {
	if err := gui.createMessagePanel(currentView, "", gui.Tr.SLocalize("PushWait")); err != nil {
		return err
	}

	go func() {
		branchName := gui.State.Branches[0].Name
		if err := gui.GitCommand.Push(branchName, force); err != nil {
			_ = gui.createErrorPanel(err.Error())
		} else {
			_ = gui.closeConfirmationPrompt()
			_ = gui.refreshCommits()
			_ = gui.refreshStatus()
		}
		_ = gui.refreshFiles()
	}()

	return nil
}

// pushFiles pushes the files...
// g and v are added by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) pushFiles(g *gocui.Gui, v *gocui.View) error {
	_, pullables := gui.GitCommand.UpstreamDifferenceCount()
	if pullables == "?" || pullables == "0" {
		return gui.pushWithForceFlag(v, false)
	}

	return gui.createConfirmationPanel(nil, gui.Tr.SLocalize("ForcePush"), gui.Tr.SLocalize("ForcePushPrompt"),
		func(g *gocui.Gui, v *gocui.View) error {
			return gui.pushWithForceFlag(v, true)
		}, nil)
}

// handleSwitchToMerge is called when a user wants to start merging.
// v is used to refocus after merging is done.
// returns an error if something goes wrong.
func (gui *Gui) handleSwitchToMerge(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}

	mergeView, err := gui.g.View("main")
	if err != nil {
		return err
	}

	if !file.HasMergeConflicts {
		return gui.createErrorPanel(gui.Tr.SLocalize("FileNoMergeCons"))
	}

	if err = gui.switchFocus(v, mergeView); err != nil {
		return err
	}

	return gui.refreshMergePanel()
}

// handleAbortMerge is called when someone aborts the merge... duuhh.
// g and v are passed by the gocui library but g is not used.
// If anything goes wrong, it returns an error
func (gui *Gui) handleAbortMerge(g *gocui.Gui, v *gocui.View) error {
	if err := gui.GitCommand.AbortMerge(); err != nil {
		return gui.createErrorPanel(err.Error())
	}

	if err := gui.createMessagePanel(v, "", gui.Tr.SLocalize("MergeAborted")); err != nil {
		return err
	}

	if err := gui.refreshStatus(); err != nil {
		return err
	}

	return gui.refreshFiles()
}

// handleResetHard is called when the user wants to hard reset to a commit.
// g and v are passed by the gocui library.
// returns an error if something went wrong.
func (gui *Gui) handleResetHard(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(v, gui.Tr.SLocalize("ClearFilePanel"), gui.Tr.SLocalize("SureResetHardHead"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.ResetHard(); err != nil {
			if err := gui.createErrorPanel(err.Error()); err != nil {
				return err
			}
		}

		return gui.refreshFiles()
	}, nil)
}

// openFile opens a file
func (gui *Gui) openFile(filename string) error {
	if err := gui.OSCommand.OpenFile(filename); err != nil {
		return gui.createErrorPanel(err.Error())
	}
	return nil
}

// anyUnStageChanges returns whether or not there are any unstage changes.
func (gui *Gui) anyUnStagedChanges(files []commands.File) bool {
	for _, file := range files {
		if file.Tracked && file.HasUnstagedChanges {
			return true
		}
	}
	return false
}
