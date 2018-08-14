package gui

import (

	// "io"
	// "io/ioutil"

	// "strings"

	"errors"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/i18n"
)

var tr *i18n.Localizer

var (
	errNoFiles    = errors.New(tr.SLocalize("NoChangedFiles", "No changed files"))
	errNoUsername = errors.New(tr.SLocalize("NoUsernameSetErr", `No username set. Please do: git config --global user.name "Your Name"`))
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
	return gui.GitCommand.StageFile(file.Name)
}

func (gui *Gui) handleFilePress(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err == errNoFiles {
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

	if err := gui.refreshFiles(g); err != nil {
		return err
	}

	return gui.handleFileSelect(g, v)
}

func (gui *Gui) handleAddPatch(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err == errNoFiles {
			return nil
		}
		return err
	}
	if !file.HasUnstagedChanges {
		return gui.createErrorPanel(g, tr.SLocalize("FileHasNoUnstagedChanges", "File has no unstaged changes to add"))
	}
	if !file.Tracked {
		return gui.createErrorPanel(g, tr.SLocalize("CannotGitAdd", "Cannot git add --patch untracked files"))
	}
	sub, err := gui.GitCommand.AddPatch(file.Name)
	if err != nil {
		return err
	}
	gui.SubProcess = sub
	return ErrSubProcess
}

func (gui *Gui) getSelectedFile(g *gocui.Gui) (commands.File, error) {
	if len(gui.State.Files) == 0 {
		return commands.File{}, errNoFiles
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
		if err == errNoFiles {
			return nil
		}
		return err
	}
	var deleteVerb string
	if file.Tracked {
		deleteVerb = "checkout"
	} else {
		deleteVerb = "delete"
	}
	return gui.createConfirmationPanel(g, v, strings.Title(deleteVerb)+" file", "Are you sure you want to "+deleteVerb+" "+file.Name+" (you will lose your changes)?", func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.RemoveFile(file); err != nil {
			panic(err)
		}
		return gui.refreshFiles(g)
	}, nil)
}

func (gui *Gui) handleIgnoreFile(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	if file.Tracked {
		return gui.createErrorPanel(g, tr.SLocalize("CantIgnoreTrackFiles", "Cannot ignore tracked files"))
	}
	gui.GitCommand.Ignore(file.Name)
	return gui.refreshFiles(g)
}

func (gui *Gui) renderfilesOptions(g *gocui.Gui, file *commands.File) error {
	optionsMap := map[string]string{
		"← → ↑ ↓":   gui.Tr.SLocalize("navigate", "navigate"),
		"S":         gui.Tr.SLocalize("stashFiles", "stash files"),
		"c":         gui.Tr.SLocalize("CommitChanges", "commit changes"),
		"o":         gui.Tr.SLocalize("open", "open"),
		"i":         gui.Tr.SLocalize("ignore", "ignore"),
		"d":         gui.Tr.SLocalize("delete", "delete"),
		"space":     gui.Tr.SLocalize("toggleStaged", "toggle staged"),
		"R":         gui.Tr.SLocalize("refresh", "refresh"),
		"t":         gui.Tr.SLocalize("addPatch", "add patch"),
		"e":         gui.Tr.SLocalize("edit", "edit"),
		"PgUp/PgDn": gui.Tr.SLocalize("scroll", "scroll"),
	}
	if gui.State.HasMergeConflicts {
		optionsMap["a"] = gui.Tr.SLocalize("abortMerge", "abort merge")
		optionsMap["m"] = gui.Tr.SLocalize("resolveMergeConflicts", "resolve merge conflicts")
	}
	if file == nil {
		return gui.renderOptionsMap(g, optionsMap)
	}
	if file.Tracked {
		optionsMap["d"] = gui.Tr.SLocalize("checkout", "checkout")
	}
	return gui.renderOptionsMap(g, optionsMap)
}

func (gui *Gui) handleFileSelect(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != errNoFiles {
			return err
		}
		gui.renderString(g, "main", "No changed files")
		return gui.renderfilesOptions(g, nil)
	}
	gui.renderfilesOptions(g, &file)
	var content string
	if file.HasMergeConflicts {
		return gui.refreshMergePanel(g)
	}

	content = gui.GitCommand.Diff(file)
	return gui.renderString(g, "main", content)
}

func (gui *Gui) handleCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && !gui.State.HasMergeConflicts {
		return gui.createErrorPanel(g, "There are no staged files to commit")
	}
	commitMessageView := gui.getCommitMessageView(g)
	g.Update(func(g *gocui.Gui) error {
		g.SetViewOnTop("commitMessage")
		gui.switchFocus(g, filesView, commitMessageView)
		return nil
	})
	return nil
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress(g *gocui.Gui, filesView *gocui.View) error {
	if len(gui.stagedFiles()) == 0 && !gui.State.HasMergeConflicts {
		return gui.createErrorPanel(g, "There are no staged files to commit")
	}
	gui.PrepareSubProcess(g, "git", "commit")
	return nil
}

// PrepareSubProcess - prepare a subprocess for execution and tell the gui to switch to it
func (gui *Gui) PrepareSubProcess(g *gocui.Gui, commands ...string) error {
	sub, err := gui.GitCommand.PrepareCommitSubProcess()
	if err != nil {
		return err
	}
	gui.SubProcess = sub
	g.Update(func(g *gocui.Gui) error {
		return ErrSubProcess
	})
	return nil
}

func (gui *Gui) genericFileOpen(g *gocui.Gui, v *gocui.View, open func(string) (*exec.Cmd, error)) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != errNoFiles {
			return err
		}
		return nil
	}
	sub, err := open(file.Name)
	if err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	if sub != nil {
		gui.SubProcess = sub
		return ErrSubProcess
	}
	return nil
}

func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {
	return gui.genericFileOpen(g, v, gui.OSCommand.EditFile)
}

func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {
	return gui.genericFileOpen(g, v, gui.OSCommand.OpenFile)
}

func (gui *Gui) handleSublimeFileOpen(g *gocui.Gui, v *gocui.View) error {
	return gui.genericFileOpen(g, v, gui.OSCommand.SublimeOpenFile)
}

func (gui *Gui) handleVsCodeFileOpen(g *gocui.Gui, v *gocui.View) error {
	return gui.genericFileOpen(g, v, gui.OSCommand.VsCodeOpenFile)
}

func (gui *Gui) handleRefreshFiles(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshFiles(g)
}

func (gui *Gui) refreshStateFiles() {
	// get files to stage
	files := gui.GitCommand.GetStatusFiles()
	gui.State.Files = gui.GitCommand.MergeStatusFiles(gui.State.Files, files)
	gui.updateHasMergeConflictStatus()
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
		if err != errNoFiles {
			return "", err
		}
		return "", gui.renderString(g, "main", "No file to display")
	}
	cat, err := gui.GitCommand.CatFile(item.Name)
	if err != nil {
		panic(err)
	}
	return cat, nil
}

func (gui *Gui) refreshFiles(g *gocui.Gui) error {
	filesView, err := g.View("files")
	if err != nil {
		return err
	}
	gui.refreshStateFiles()
	filesView.Clear()
	for _, file := range gui.State.Files {
		gui.renderFile(file, filesView)
	}
	gui.correctCursor(filesView)
	if filesView == g.CurrentView() {
		gui.handleFileSelect(g, filesView)
	}
	return nil
}

func (gui *Gui) pullFiles(g *gocui.Gui, v *gocui.View) error {
	gui.createMessagePanel(g, v, "", "Pulling...")
	go func() {
		if err := gui.GitCommand.Pull(); err != nil {
			gui.createErrorPanel(g, err.Error())
		} else {
			gui.closeConfirmationPrompt(g)
			gui.refreshCommits(g)
			gui.refreshStatus(g)
		}
		gui.refreshFiles(g)
	}()
	return nil
}

func (gui *Gui) pushFiles(g *gocui.Gui, v *gocui.View) error {
	gui.createMessagePanel(g, v, "", "Pushing...")
	go func() {
		branchName := gui.State.Branches[0].Name
		if err := gui.GitCommand.Push(branchName); err != nil {
			gui.createErrorPanel(g, err.Error())
		} else {
			gui.closeConfirmationPrompt(g)
			gui.refreshCommits(g)
			gui.refreshStatus(g)
		}
	}()
	return nil
}

func (gui *Gui) handleSwitchToMerge(g *gocui.Gui, v *gocui.View) error {
	mergeView, err := g.View("main")
	if err != nil {
		return err
	}
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != errNoFiles {
			return err
		}
		return nil
	}
	if !file.HasMergeConflicts {
		return gui.createErrorPanel(g, "This file has no merge conflicts")
	}
	gui.switchFocus(g, v, mergeView)
	return gui.refreshMergePanel(g)
}

func (gui *Gui) handleAbortMerge(g *gocui.Gui, v *gocui.View) error {
	if err := gui.GitCommand.AbortMerge(); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	gui.createMessagePanel(g, v, "", "Merge aborted")
	gui.refreshStatus(g)
	return gui.refreshFiles(g)
}

func (gui *Gui) handleResetHard(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(g, v, "Clear file panel", "Are you sure you want `reset --hard HEAD`? You may lose changes", func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.ResetHard(); err != nil {
			gui.createErrorPanel(g, err.Error())
		}
		return gui.refreshFiles(g)
	}, nil)
}
