package gui

import (
	// "io"
	// "io/ioutil"

	// "strings"

	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

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

func (gui *Gui) stageSelectedFile(g *gocui.Gui) error {

	file, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}

	err = gui.GitCommand.StageFile(file.Name)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleFilePress(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
		return err
	}

	if file.HasMergeConflicts {
		return gui.handleSwitchToMerge(g, v)
	}

	if file.HasUnstagedChanges {
		gui.GitCommand.StageFile(file.Name)
	} else {
		gui.GitCommand.UnStageFile(file.Name, file.Tracked)
	}

	err = gui.refreshFiles(g)
	if err != nil {
		return err
	}

	err = gui.handleFileSelect(g, v)
	if err != nil {
		return err
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

	err = gui.refreshFiles(g)
	if err != nil {
		return err
	}

	err = gui.handleFileSelect(g, v)
	if err != nil {
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

	var deleteVerb string

	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err == gui.Errors.ErrNoFiles {
			return nil
		}
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

	err = gui.createConfirmationPanel(g, v, strings.Title(deleteVerb)+" file", message,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.RemoveFile(file)
			if err != nil {
				return err
			}

			err = gui.refreshFiles(g)
			if err != nil {
				return err
			}

			return nil
		}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleIgnoreFile(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile(g)
	if err != nil {
		return gui.createErrorPanel(g, err.Error())
	}

	if file.Tracked {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantIgnoreTrackFiles"))
	}

	err = gui.GitCommand.Ignore(file.Name)
	if err != nil {
		return gui.createErrorPanel(g, err.Error())
	}

	err = gui.refreshFiles(g)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) renderfilesOptions(g *gocui.Gui, file *commands.File) error {
	optionsMap := map[string]string{
		"← → ↑ ↓":   gui.Tr.SLocalize("navigate"),
		"S":         gui.Tr.SLocalize("stashFiles"),
		"c":         gui.Tr.SLocalize("CommitChanges"),
		"o":         gui.Tr.SLocalize("open"),
		"i":         gui.Tr.SLocalize("ignore"),
		"d":         gui.Tr.SLocalize("delete"),
		"space":     gui.Tr.SLocalize("toggleStaged"),
		"R":         gui.Tr.SLocalize("refresh"),
		"t":         gui.Tr.SLocalize("addPatch"),
		"e":         gui.Tr.SLocalize("edit"),
		"a":         gui.Tr.SLocalize("toggleStagedAll"),
		"PgUp/PgDn": gui.Tr.SLocalize("scroll"),
	}

	if gui.State.HasMergeConflicts {
		optionsMap["a"] = gui.Tr.SLocalize("abortMerge")
		optionsMap["m"] = gui.Tr.SLocalize("resolveMergeConflicts")
	}

	if file == nil {
		return gui.renderOptionsMap(g, optionsMap)
	}

	if file.Tracked {
		optionsMap["d"] = gui.Tr.SLocalize("checkout")
	}

	err := gui.renderOptionsMap(g, optionsMap)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleFileSelect(g *gocui.Gui, v *gocui.View) error {

	var content string

	file, err := gui.getSelectedFile(g)
	if err != nil {

		if err != gui.Errors.ErrNoFiles {
			return err
		}

		err = gui.renderString(g, "main", gui.Tr.SLocalize("NoChangedFiles"))
		if err != nil {
			return err
		}

		err = gui.renderfilesOptions(g, nil)
		if err != nil {
			return err
		}

		return nil
	}

	err = gui.renderfilesOptions(g, &file)
	if err != nil {
		return nil
	}

	if file.HasMergeConflicts {
		return gui.refreshMergePanel(g)
	}

	content = gui.GitCommand.Diff(file)

	err = gui.renderString(g, "main", content)
	if err != nil {
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

		_, err := g.SetViewOnTop("commitMessage")
		if err != nil {
			return err
		}

		err = gui.switchFocus(g, filesView, commitMessageView)
		if err != nil {
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

		err = gui.createErrorPanel(gui.g, err.Error())
		if err != nil {
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

func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}

	err = gui.editFile(file.Name)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {

	file, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}

	err = gui.openFile(file.Name)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleRefreshFiles(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshFiles(g)
}

func (gui *Gui) refreshStateFiles() error {

	files := gui.GitCommand.GetStatusFiles()

	gui.State.Files = gui.GitCommand.MergeStatusFiles(gui.State.Files, files)

	gui.updateHasMergeConflictStatus()

	return nil
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
	// TODO check solution for this
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

		err = gui.renderString(g, "main", gui.Tr.SLocalize("NoFilesDisplay"))
		if err != nil {
			return "", err
		}

		return "", nil
	}

	if item.Type != "file" {

		err = gui.renderString(g, "main", gui.Tr.SLocalize("NotAFile"))
		if err != nil {
			return "", err
		}

		return "", nil
	}

	cat, err := gui.GitCommand.CatFile(item.Name)
	if err != nil {

		gui.Log.Error(err)

		err = gui.renderString(g, "main", err.Error())
		if err != nil {
			return "", err
		}

		return "", nil
	}

	return cat, nil
}

func (gui *Gui) refreshFiles(g *gocui.Gui) error {

	filesView, err := g.View("files")
	if err != nil {
		return err
	}

	err = gui.refreshStateFiles()
	if err != nil {
		return err
	}

	filesView.Clear()

	for _, file := range gui.State.Files {
		gui.renderFile(file, filesView)
	}

	err = gui.correctCursor(filesView)
	if err != nil {
		return err
	}

	if filesView == g.CurrentView() {
		err = gui.handleFileSelect(g, filesView)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) pullFiles(g *gocui.Gui, v *gocui.View) error {

	err := gui.createMessagePanel(g, v, "", gui.Tr.SLocalize("PullWait"))
	if err != nil {
		return err
	}

	go func() {

		err := gui.GitCommand.Pull()
		if err != nil {

			err = gui.createErrorPanel(g, err.Error())
			if err != nil {
				gui.Log.Error(err)
				return
			}

		} else {

			err = gui.closeConfirmationPrompt(g)
			if err != nil {
				gui.Log.Error(err)
				return
			}

			err = gui.refreshCommits(g)
			if err != nil {
				gui.Log.Error(err)
				return
			}

			err = gui.refreshStatus(g)
			if err != nil {
				gui.Log.Error(err)
				return
			}

		}

		err = gui.refreshFiles(g)
		if err != nil {
			gui.Log.Error(err)
			return
		}

	}()

	return nil
}

func (gui *Gui) pushWithForceFlag(currentView *gocui.View, force bool) error {

	err := gui.createMessagePanel(gui.g, currentView, "", gui.Tr.SLocalize("PushWait"))
	if err != nil {
		return err
	}

	go func() {

		branchName := gui.State.Branches[0].Name

		err := gui.GitCommand.Push(branchName, force)
		if err != nil {

			err = gui.createErrorPanel(gui.g, err.Error())
			if err != nil {
				gui.Log.Error(err)
				return
			}

		} else {

			err = gui.closeConfirmationPrompt(gui.g)
			if err != nil {
				gui.Log.Error(err)
				return
			}

			err = gui.refreshCommits(gui.g)
			if err != nil {
				gui.Log.Error(err)
				return
			}

			err = gui.refreshStatus(gui.g)
			if err != nil {
				gui.Log.Error(err)
				return
			}

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

	err := gui.createConfirmationPanel(g, nil, gui.Tr.SLocalize("ForcePush"), gui.Tr.SLocalize("ForcePushPrompt"),
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.pushWithForceFlag(v, true)
			if err != nil {
				return err
			}

			return nil
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

	err = gui.switchFocus(g, v, mergeView)
	if err != nil {
		return err
	}

	err = gui.refreshMergePanel(g)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleAbortMerge(g *gocui.Gui, v *gocui.View) error {

	err := gui.GitCommand.AbortMerge()
	if err != nil {

		err = gui.createErrorPanel(g, err.Error())
		if err != nil {
			return err
		}

		return nil
	}

	err = gui.createMessagePanel(g, v, "", gui.Tr.SLocalize("MergeAborted"))
	if err != nil {
		return err
	}

	err = gui.refreshStatus(g)
	if err != nil {
		return err
	}

	err = gui.refreshFiles(g)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleResetHard(g *gocui.Gui, v *gocui.View) error {

	err := gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("ClearFilePanel"), gui.Tr.SLocalize("SureResetHardHead"),
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.ResetHard()
			if err != nil {

				err = gui.createErrorPanel(g, err.Error())
				if err != nil {
					return err
				}

			}

			err = gui.refreshFiles(g)
			if err != nil {
				return err
			}

			return nil
		}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (gui *Gui) openFile(filename string) error {

	err := gui.OSCommand.OpenFile(filename)
	if err != nil {

		err = gui.createErrorPanel(gui.g, err.Error())
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
