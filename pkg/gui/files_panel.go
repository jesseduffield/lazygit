package gui

import (
	"fmt"
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
		gui.Log.Error(fmt.Sprintf("failed to get filesview in refreshfiles: %s", err))
		return err
	}

	err = gui.refreshStateFiles()
	if err != nil {
		gui.Log.Error(fmt.Sprintf("failed to refresh state files in refreshfiles: %s", err))
		return err
	}

	filesView.Clear()

	for _, file := range gui.State.Files {
		gui.renderFileName(file, filesView)
	}

	err = gui.correctCursor(filesView)
	if err != nil {
		gui.Log.Error(fmt.Sprintf("failed correctCursor in refreshfiles: %s", err))
		return err
	}

	if filesView == gui.g.CurrentView() {
		err = gui.handleFileSelect()
		if err != nil {
			gui.Log.Error(fmt.Sprintf("failed to get filesview in refreshfiles: %s", err))
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

	err := gui.updateHasMergeConflictStatus()
	if err != nil {
		gui.Log.Error(fmt.Sprintf("failed to updateHasMergedConflictStatus"+
			" in refreshStateFiles: %s", err))
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
		gui.Log.Errorf("Failed to get selected file at stageSelectedFile: %s\n", err)
		return err
	}

	err = gui.GitCommand.StageFile(file.Name)
	if err != nil {
		gui.Log.Errorf("Failed to stageFile at stageSelectedFile: %s\n", err)
	}

	return nil
}

// handleFilePress is called when the user selects a file.
// g and v are passed by the gocui library to the function but are not used.
// In case something goes wrong, this function returns an error.
func (gui *Gui) handleFilePress(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile()
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
		gui.Log.Errorf("Failed to getSelectedFile at handleFilePress")
		return err
	}

	if file.HasMergeConflicts {

		err = gui.handleSwitchToMerge(gui.g, v)
		if err != nil {
			gui.Log.Errorf("Failed to handleSwitchToMerge at handleFilePress: %s\n", err)
			return err
		}

		return nil
	}

	if file.HasUnstagedChanges {
		err = gui.GitCommand.StageFile(file.Name)
	} else {
		err = gui.GitCommand.UnStageFile(file.Name, file.Tracked)
	}
	if err != nil {
		gui.Log.Errorf("Failed to stage/unstage file at handleFilePress: %s\n", err)
		return err
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Error("Failed to refresh files at handleFilePress: ", err)
		return err
	}

	err = gui.handleFileSelect()
	if err != nil {
		gui.Log.Error("Failed to handleFileSelect at handleFilePress: ", err)
	}

	return nil
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
		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Error("Failed to create error panel in handleStageAll: ", err)
			return err
		}
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Error("Failed to refreshFiles in handleStageAll: ", err)
		return err
	}

	err = gui.handleFileSelect()
	if err != nil {
		gui.Log.Error("Failed to handleFileSelect in handleStageAll: ", err)
		return err
	}

	return nil
}

// handleAddPatch is called when a user wants to add a patch.
// g and v are passed by the gocui library
func (gui *Gui) handleAddPatch(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile()
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
		gui.Log.Errorf("Failed to getSelectedFile at handleAddPatch: %s\n", err)
		return err
	}

	if !file.HasUnstagedChanges {
		err = gui.createErrorPanel(gui.Tr.SLocalize("FileHasNoUnstagedChanges"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleAddPatch: %s\n", err)
			return err
		}
		return nil
	}

	if !file.Tracked {
		err = gui.createErrorPanel(gui.Tr.SLocalize("CannotGitAdd"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleAddPatch: %s\n", err)
			return err
		}

		return nil
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
		gui.Log.Errorf("Failed to get files view at getSelectedFile: %s\n", err)
		return commands.File{}, err
	}

	lineNumber := gui.getItemPosition(filesView)

	return gui.State.Files[lineNumber], nil
}

// handleFileRemoved gets called when a file is removed.
// g and v are passed by the gocui library.
// returns an error if something goes wrong
func (gui *Gui) handleFileRemove(g *gocui.Gui, v *gocui.View) error {

	var deleteVerb string

	file, err := gui.getSelectedFile()
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}

		gui.Log.Errorf("Failed to getSelectedFile at handleFileRemove: %s\n", err)
		return err
	}

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

	err = gui.createConfirmationPanel(v, strings.Title(deleteVerb)+" file", message,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.RemoveFile(file)
			if err != nil {
				gui.Log.Errorf("Failed to remove file: %s\n", err)
				return err
			}

			err = gui.refreshFiles()
			if err != nil {
				gui.Log.Errorf("Failed to refresh files: %s\n", err)
				return err
			}

			return nil
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to createConfirmationPanel at handleFileRemove: %s\n", err)
		return err
	}

	return nil
}

// handleIgnoreFile handles ignoring the file.
// g and v are passed by the gocui.
// returns an error if something goes wrong.
func (gui *Gui) handleIgnoreFile(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile()
	if err != nil {
		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleIgnoreFile: %s\n", err)
			return err
		}

		return nil
	}

	if file.Tracked {

		err = gui.createErrorPanel(gui.Tr.SLocalize("CantIgnoreTrackFiles"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleIgnoreFile: %s\n", err)
			return err
		}

		return nil
	}
	if err := gui.GitCommand.Ignore(file.Name); err != nil {

		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleIgnoreFile: %s\n", err)
			return err
		}

		return nil
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Errorf("Failed to refresh files at handleIgnoreFile: %s\n", err)
		return err
	}

	return nil
}

// handleFileSelect is called when a file is selected.
// It checks if there are any changed files and if there is one
// and it is selected, it gets rendered into the main view
func (gui *Gui) handleFileSelect() error {

	var content string

	file, err := gui.getSelectedFile()
	if err != nil {

		if err != gui.Errors.ErrNoFiles {
			gui.Log.Errorf("Failed to get selected file in handlefileselect: %s\n", err)
			return err
		}

		err = gui.renderString(gui.g, "main", gui.Tr.SLocalize("NoChangedFiles"))
		if err != nil {
			gui.Log.Errorf("Failed to render string in handlefileselect: %s\n", err)
			return err
		}

		err = gui.renderGlobalOptions()
		if err != nil {
			gui.Log.Errorf("Failed to renderfilesoptions in handlefileselect: %s\n", err)
			return err
		}

		return nil
	}

	err = gui.renderGlobalOptions()
	if err != nil {
		gui.Log.Errorf("Failed to renderfilesoptions in handlefileselect: %s\n", err)
		return err
	}

	if file.HasMergeConflicts {

		err = gui.refreshMergePanel(gui.g)
		if err != nil {
			gui.Log.Errorf("Failed to refreshmergepanel in handlefileselect: %s\n", err)
			return err
		}

		return nil
	}

	content = gui.GitCommand.Diff(file)

	err = gui.renderString(gui.g, "main", content)
	if err != nil {
		gui.Log.Errorf("Failed to render string in handlefileselect: %s\n", err)
		return err
	}

	return nil
}

// handleCommitPress is called when a user commits changes.
// g and v are passed by the gocui library.
// returns and error if something goes wrong.
func (gui *Gui) handleCommitPress(g *gocui.Gui, filesView *gocui.View) error {

	if len(gui.stagedFiles()) == 0 && !gui.State.HasMergeConflicts {

		err := gui.createErrorPanel(gui.Tr.SLocalize("NoStagedFilesToCommit"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleCommitPress: %s\n", err)
			return err
		}

		return nil
	}

	commitMessageView, err := gui.g.View("commitMessage")
	if err != nil {
		gui.Log.Errorf("Failed to get commitMessage view at handleCommitPress: %s\n", err)
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		_, err = g.SetViewOnTop("commitMessage")
		if err != nil {

		}

		err = gui.switchFocus(g, filesView, commitMessageView)
		if err != nil {

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
		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at editFile: %s\n", err)
			return err
		}

		return nil
	}

	if sub != nil {
		gui.SubProcess = sub
		return gui.Errors.ErrSubProcess
	}

	return nil
}

// handleFileEdit is called when a user wants to edit a file
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile()
	if err != nil {
		gui.Log.Errorf("Failed to get selected file at handleFileEdit: %s\n", err)
		return err
	}

	err = gui.editFile(file.Name)
	if err != nil {
		gui.Log.Errorf("Failed to editfile at handleFileEdit: %s\n", err)
		return err
	}

	return nil
}

// handleFileOpen is called when a user opens a file
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile()
	if err != nil {
		gui.Log.Errorf("Failed to getSelectedFile at handleFileOpen: %s\n", err)
		return err
	}

	err = gui.openFile(file.Name)
	if err != nil {
		gui.Log.Errorf("Failed to openFile at handleFileOpen: %s\n", err)
		return err
	}

	return nil
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
		gui.Log.Errorf("Failed to get merge state at updateHasMergeConflictStatus: %s\n", err)
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
			gui.Log.Errorf("FAiled to getSelectedFile: %s\n", err)
			return "", err
		}

		err = gui.renderString(gui.g, "main", gui.Tr.SLocalize("NoFilesDisplay"))
		if err != nil {
			gui.Log.Errorf("Failed to renderString at catSelectedFile: %s\n", err)
			return "", err
		}

		return "", nil
	}

	if item.Type != "file" {

		err = gui.renderString(gui.g, "main", gui.Tr.SLocalize("NotAFile"))
		if err != nil {
			gui.Log.Errorf("Failed to renderString at catSelectedFile: %s\n", err)
			return "", err
		}

		return "", nil
	}

	cat, err := gui.GitCommand.CatFile(item.Name)
	if err != nil {
		gui.Log.Error(err)
		return "", gui.renderString(gui.g, "main", err.Error())
	}

	return cat, nil
}

// pullFiles pulls the files.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) pullFiles(g *gocui.Gui, v *gocui.View) error {

	err := gui.createMessagePanel(v, "", gui.Tr.SLocalize("PullWait"))
	if err != nil {
		gui.Log.Errorf("Failed to createMessagePanel at pullFiles: %s\n", err)
		return err
	}

	go func() {

		err := gui.GitCommand.Pull()
		if err != nil {
			err = gui.createErrorPanel(err.Error())
			if err != nil {
				gui.Log.Errorf("Failed to create error panel at pullFiles: %s\n", err)
				return
			}
		} else {

			err = gui.closeConfirmationPrompt()
			if err != nil {
				gui.Log.Errorf("Failed to create confirmation panel at pullFiles: %s\n", err)
				return
			}

			err = gui.refreshCommits()
			if err != nil {
				gui.Log.Errorf("Failed to refresh commits at pullFiles: %s\n", err)
				return
			}

			err = gui.refreshStatus()
			if err != nil {
				gui.Log.Errorf("Failed to refresh status at pullFiles: %s\n", err)
				return
			}

		}

		err = gui.refreshFiles()
		if err != nil {
			gui.Log.Errorf("Failed to refresh files at pullfiles: %s\n", err)
			return
		}

	}()

	return nil
}

// pushWithForceFlag does what it says, it pushes with or without force.
// currentView is used to return focus to the view.
// force indicates whether or not to push with foce.
func (gui *Gui) pushWithForceFlag(currentView *gocui.View, force bool) error {

	err := gui.createMessagePanel(currentView, "", gui.Tr.SLocalize("PushWait"))
	if err != nil {
		gui.Log.Errorf("Failed to createMessagePanel at pushWithForceFlag: %s\n", err)
		return err
	}

	go func() {
		branchName := gui.State.Branches[0].Name
		err := gui.GitCommand.Push(branchName, force)
		if err != nil {
			err = gui.createErrorPanel(err.Error())
			if err != nil {
				gui.Log.Errorf("Failed to create error panel at pushWithForceFlag: %s\n")
			}
		} else {
			err = gui.closeConfirmationPrompt()
			if err != nil {
				gui.Log.Errorf("Failed to closeConfirmationPrompt at pushWithForceFlag: %s\n", err)
			}

			err = gui.refreshCommits()
			if err != nil {
				gui.Log.Errorf("Failed to refreshCommits at pushWithForceFlag: %s\n", err)
			}

			err = gui.refreshStatus()
			if err != nil {
				gui.Log.Errorf("Failed to refreshStatus at pushWithForceFlag: %s\n", err)
			}
		}
	}()

	return nil
}

// pushFiles pushes the files...
// g and v are added by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) pushFiles(g *gocui.Gui, v *gocui.View) error {

	_, pullables := gui.GitCommand.UpstreamDifferenceCount()
	if pullables == "?" || pullables == "0" {

		err := gui.pushWithForceFlag(v, false)
		if err != nil {
			gui.Log.Errorf("Failed to push with force at pushFiles: %s\n", err)
			return err
		}

		return nil
	}

	err := gui.createConfirmationPanel(nil, gui.Tr.SLocalize("ForcePush"), gui.Tr.SLocalize("ForcePushPrompt"),
		func(g *gocui.Gui, v *gocui.View) error {
			err := gui.pushWithForceFlag(v, true)
			if err != nil {
				gui.Log.Errorf("Failed to pushWithForceFlag at pushFiles: %s\n", err)
				return err
			}

			return nil
		}, nil)

	return err
}

// handleSwitchToMerge is called when a user wants to start merging.
// v is used to refocus after merging is done.
// returns an error if something goes wrong.
func (gui *Gui) handleSwitchToMerge(g *gocui.Gui, v *gocui.View) error {

	mergeView, err := gui.g.View("main")
	if err != nil {
		return err
	}

	file, err := gui.getSelectedFile()
	if err != nil {

		if err != gui.Errors.ErrNoFiles {
			gui.Log.Errorf("Failed to get selected file at handleSwitchToMerge: %s\n", err)
			return err
		}

		return nil
	}

	if !file.HasMergeConflicts {
		err = gui.createErrorPanel(gui.Tr.SLocalize("FileNoMergeCons"))
		if err != nil {
			gui.Log.Errorf("Failed to createErrorPanel at handleSwitchToMerge: %s\n", err)
			return err
		}

		return nil
	}

	err = gui.switchFocus(gui.g, v, mergeView)
	if err != nil {
		gui.Log.Errorf("Failed to switchFocus at handleSwitchToMerge: %s\n", err)
		return err
	}

	err = gui.refreshMergePanel(gui.g)
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at handleSwitchToMerge: %s\n", err)
		return err
	}

	return nil
}

// handleAbortMerge is called when someone aborts the merge... duuhh.
// g and v are passed by the gocui library but g is not used.
// If anything goes wrong, it returns an error
func (gui *Gui) handleAbortMerge(g *gocui.Gui, v *gocui.View) error {

	err := gui.GitCommand.AbortMerge()
	if err != nil {

		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleAbortMerge: %s\n", err)
			return err
		}
	}

	err = gui.createMessagePanel(v, "", gui.Tr.SLocalize("MergeAborted"))
	if err != nil {
		gui.Log.Errorf("Failed to create message panel at handleAbortMerge: %s\n", err)
		return err
	}

	err = gui.refreshStatus()
	if err != nil {
		gui.Log.Errorf("Failed to refresh status at handleAbortMerge: %s\n", err)
		return err
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Errorf("Failed to refresh files at handleAbortMerge: %s\n", err)
		return err
	}

	return nil
}

// handleResetHard is called when the user wants to hard reset to a commit.
// g and v are passed by the gocui library.
// returns an error if something went wrong.
func (gui *Gui) handleResetHard(g *gocui.Gui, v *gocui.View) error {

	err := gui.createConfirmationPanel(v, gui.Tr.SLocalize("ClearFilePanel"), gui.Tr.SLocalize("SureResetHardHead"), func(g *gocui.Gui, v *gocui.View) error {
		err := gui.GitCommand.ResetHard()
		if err != nil {
			err = gui.createErrorPanel(err.Error())
			if err != nil {
				gui.Log.Errorf("Failed to create error panel at handleHardResest: %s\n", err)
				return err
			}
		}

		err = gui.refreshFiles()
		if err != nil {
			gui.Log.Errorf("Failed to refresh files at handleHardReset: %s\n", err)
			return err
		}

		return nil
	}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to createConfirmationPanel at handleHardReset: %s\n", err)
		return err
	}

	return nil
}

// openFile opens a file
func (gui *Gui) openFile(filename string) error {

	err := gui.OSCommand.OpenFile(filename)
	if err != nil {

		err = gui.createErrorPanel(err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at openFile: %s\n", err)
			return err
		}

		return nil
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
