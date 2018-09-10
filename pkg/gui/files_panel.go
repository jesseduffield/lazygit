package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// refreshFiles refreshes the files view
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
		gui.renderFile(file, filesView)
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

// refreshStateFiles refreshes the state files
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

func (gui *Gui) stageSelectedFile() error {
	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		return err
	}
	return gui.GitCommand.StageFile(file.Name)
}

// handleFilePress is called when the user selects a file.
// g and v are passed by the gocui library to the function but are not used.
// In case something goes wrong, this function returns an error
func (gui *Gui) handleFilePress(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
		return err
	}

	if file.HasMergeConflicts {
		return gui.handleSwitchToMerge(gui.g, v)
	}

	if file.HasUnstagedChanges {
		gui.GitCommand.StageFile(file.Name)
	} else {
		gui.GitCommand.UnStageFile(file.Name, file.Tracked)
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
// In case something goes wrong, it returns an error
func (gui *Gui) handleStageAll(g *gocui.Gui, v *gocui.View) error {

	var err error

	if gui.allFilesStaged() {
		err = gui.GitCommand.UnstageAll()
	} else {
		err = gui.GitCommand.StageAll()
	}

	if err != nil {
		err = gui.createErrorPanel(gui.g, err.Error())
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

func (gui *Gui) getSelectedFile(g *gocui.Gui) (commands.File, error) {
	if len(gui.State.Files) == 0 {
		return commands.File{}, gui.Errors.ErrNoFiles
	}
	filesView, err := g.View("files")
	if err != nil {
		panic(err)
	}
	lineNumber := gui.getItemPosition(filesView)
	return gui.State.Files[lineNumber], nil
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

// handleFileSelect is called when a file is selected.
// It checks if there are any changed files and if there is one
// and it is selected, it gets rendered into the main view
func (gui *Gui) handleFileSelect() error {

	var content string

	file, err := gui.getSelectedFile(gui.g)
	if err != nil {

		if err != gui.Errors.ErrNoFiles {
			gui.Log.Error("Failed to get selected file in handlefileselect: ", err)
			return err
		}

		err = gui.renderString(gui.g, "main", gui.Tr.SLocalize("NoChangedFiles"))
		if err != nil {
			gui.Log.Error("Failed to render string in handlefileselect: ", err)
			return err
		}

		err = gui.renderGlobalOptions()
		if err != nil {
			gui.Log.Error("Failed to renderfilesoptions in handlefileselect: ", err)
			return err
		}

		return nil
	}

	err = gui.renderGlobalOptions()
	if err != nil {
		gui.Log.Error("Failed to renderfilesoptions in handlefileselect: ", err)
		return err
	}

	if file.HasMergeConflicts {

		err = gui.refreshMergePanel(gui.g)
		if err != nil {
			gui.Log.Error("Failed to refreshmergepanel in handlefileselect: ", err)
			return err
		}

		return nil
	}

	content = gui.GitCommand.Diff(file)

	err = gui.renderString(gui.g, "main", content)
	if err != nil {
		gui.Log.Error("Failed to render string in handlefileselect : ", err)
		return err
	}

	return nil
}

func (gui *Gui) handleCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && !gui.State.HasMergeConflicts {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("NoStagedFilesToCommit"))
	}
	commitMessageView := gui.getCommitMessageView(g)
	g.Update(func(g *gocui.Gui) error {
		g.SetViewOnTop("commitMessage")
		gui.switchFocus(g, filesView, commitMessageView)
		gui.RenderCommitLength()
		return nil
	})
	return nil
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && !gui.State.HasMergeConflicts {
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
	sub, err := gui.OSCommand.EditFile(filename)
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}
	if sub != nil {
		gui.SubProcess = sub
		return gui.Errors.ErrSubProcess
	}
	return nil
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

func (gui *Gui) updateHasMergeConflictStatus() error {
	merging, err := gui.GitCommand.IsInMergeState()
	if err != nil {
		return err
	}
	gui.State.HasMergeConflicts = merging
	return nil
}

func (gui *Gui) renderFile(file commands.File, filesView *gocui.View) {
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

	err := gui.createMessagePanel(g, v, "", gui.Tr.SLocalize("PullWait"))
	if err != nil {
		gui.Log.Error(err)
		return err
	}

	go func() {

		err := gui.GitCommand.Pull()
		if err != nil {
			_ = gui.createErrorPanel(gui.g, err.Error())
		} else {
			_ = gui.closeConfirmationPrompt(gui.g)
			_ = gui.refreshCommits()
			_ = gui.refreshStatus()
		}

		err = gui.refreshFiles()
		if err != nil {
			gui.Log.Error(err)
			return
		}

	}()

	return nil
}

func (gui *Gui) pushWithForceFlag(currentView *gocui.View, force bool) error {
	if err := gui.createMessagePanel(gui.g, currentView, "", gui.Tr.SLocalize("PushWait")); err != nil {
		return err
	}
	go func() {
		branchName := gui.State.Branches[0].Name
		if err := gui.GitCommand.Push(branchName, force); err != nil {
			_ = gui.createErrorPanel(gui.g, err.Error())
		} else {
			_ = gui.closeConfirmationPrompt(gui.g)
			_ = gui.refreshCommits()
			_ = gui.refreshStatus()
		}
	}()
	return nil
}

func (gui *Gui) pushFiles(g *gocui.Gui, v *gocui.View) error {
	// if we have pullables we'll ask if the user wants to force push
	_, pullables := gui.GitCommand.UpstreamDifferenceCount()
	if pullables == "?" || pullables == "0" {
		return gui.pushWithForceFlag(v, false)
	}
	err := gui.createConfirmationPanel(g, nil, gui.Tr.SLocalize("ForcePush"), gui.Tr.SLocalize("ForcePushPrompt"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.pushWithForceFlag(v, true)
	}, nil)
	return err
}

func (gui *Gui) handleSwitchToMerge(g *gocui.Gui, v *gocui.View) error {
	mergeView, err := g.View("main")
	if err != nil {
		return err
	}
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}
	if !file.HasMergeConflicts {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("FileNoMergeCons"))
	}
	gui.switchFocus(g, v, mergeView)
	return gui.refreshMergePanel(g)
}

func (gui *Gui) handleAbortMerge(g *gocui.Gui, v *gocui.View) error {
	if err := gui.GitCommand.AbortMerge(); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	gui.createMessagePanel(g, v, "", gui.Tr.SLocalize("MergeAborted"))
	gui.refreshStatus()
	return gui.refreshFiles()
}

func (gui *Gui) handleResetHard(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("ClearFilePanel"), gui.Tr.SLocalize("SureResetHardHead"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.ResetHard(); err != nil {
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

func (gui *Gui) anyUnStagedChanges(files []commands.File) bool {
	for _, file := range files {
		if file.Tracked && file.HasUnstagedChanges {
			return true
		}
	}
	return false
}
