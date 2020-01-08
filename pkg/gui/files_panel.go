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

func (gui *Gui) selectFile(alreadySelected bool) error {
	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return gui.renderString(gui.g, "main", gui.Tr.SLocalize("NoChangedFiles"))
	}

	if err := gui.focusPoint(0, gui.State.Panels.Files.SelectedLine, len(gui.State.Files), gui.getFilesView()); err != nil {
		return err
	}

	if file.HasInlineMergeConflicts {
		gui.getMainView().Title = gui.Tr.SLocalize("MergeConflictsTitle")
		gui.State.SplitMainPanel = false
		return gui.refreshMergePanel()
	}

	content := gui.GitCommand.Diff(file, false, false)
	contentCached := gui.GitCommand.Diff(file, false, true)
	leftContent := content
	if file.HasStagedChanges && file.HasUnstagedChanges {
		gui.State.SplitMainPanel = true
		gui.getMainView().Title = gui.Tr.SLocalize("UnstagedChanges")
		gui.getSecondaryView().Title = gui.Tr.SLocalize("StagedChanges")
	} else {
		gui.State.SplitMainPanel = false
		if file.HasUnstagedChanges {
			leftContent = content
			gui.getMainView().Title = gui.Tr.SLocalize("UnstagedChanges")
		} else {
			leftContent = contentCached
			gui.getMainView().Title = gui.Tr.SLocalize("StagedChanges")
		}
	}

	if alreadySelected {
		gui.g.Update(func(*gocui.Gui) error {
			if err := gui.setViewContent(gui.g, gui.getSecondaryView(), contentCached); err != nil {
				return err
			}
			return gui.setViewContent(gui.g, gui.getMainView(), leftContent)
		})
		return nil
	}
	if err := gui.renderString(gui.g, "secondary", contentCached); err != nil {
		return err
	}
	return gui.renderString(gui.g, "main", leftContent)
}

func (gui *Gui) refreshFiles() error {
	gui.State.RefreshingFilesMutex.Lock()
	gui.State.IsRefreshingFiles = true
	defer func() {
		gui.State.IsRefreshingFiles = false
		gui.State.RefreshingFilesMutex.Unlock()
	}()

	selectedFile, _ := gui.getSelectedFile(gui.g)

	filesView := gui.getFilesView()
	if filesView == nil {
		// if the filesView hasn't been instantiated yet we just return
		return nil
	}
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

		if g.CurrentView() == filesView || (g.CurrentView() == gui.getMainView() && g.CurrentView().Context == "merging") {
			newSelectedFile, _ := gui.getSelectedFile(gui.g)
			alreadySelected := newSelectedFile.Name == selectedFile.Name
			return gui.selectFile(alreadySelected)
		}
		return nil
	})

	return nil
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
	return gui.enterFile(false, -1)
}

func (gui *Gui) enterFile(forceSecondaryFocused bool, selectedLineIdx int) error {
	file, err := gui.getSelectedFile(gui.g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}
	if file.HasInlineMergeConflicts {
		return gui.handleSwitchToMerge(gui.g, gui.getFilesView())
	}
	if file.HasMergeConflicts {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("FileStagingRequirements"))
	}
	if err := gui.changeMainViewsContext("staging"); err != nil {
		return err
	}
	if err := gui.switchFocus(gui.g, gui.getFilesView(), gui.getMainView()); err != nil {
		return err
	}
	return gui.refreshStagingPanel(forceSecondaryFocused, selectedLineIdx)
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

	return gui.selectFile(true)
}

func (gui *Gui) allFilesStaged() bool {
	for _, file := range gui.State.Files {
		if file.HasUnstagedChanges {
			return false
		}
	}
	return true
}

func (gui *Gui) focusAndSelectFile(g *gocui.Gui, v *gocui.View) error {
	if _, err := gui.g.SetCurrentView("files"); err != nil {
		return err
	}

	return gui.selectFile(false)
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

	return gui.selectFile(false)
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

func (gui *Gui) handleWIPCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	skipHookPreifx := gui.Config.GetUserConfig().GetString("git.skipHookPrefix")
	if skipHookPreifx == "" {
		return gui.createErrorPanel(gui.g, gui.Tr.SLocalize("SkipHookPrefixNotConfigured"))
	}

	if err := gui.renderString(g, "commitMessage", skipHookPreifx); err != nil {
		return err
	}
	if err := gui.getCommitMessageView().SetCursor(len(skipHookPreifx), 0); err != nil {
		return err
	}

	return gui.handleCommitPress(g, filesView)
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

	return gui.createConfirmationPanel(g, filesView, true, title, question, func(g *gocui.Gui, v *gocui.View) error {
		ok, err := gui.runSyncOrAsyncCommand(gui.GitCommand.AmendHead())
		if err != nil {
			return err
		}
		if !ok {
			return nil
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
	_, err := gui.runSyncOrAsyncCommand(gui.OSCommand.EditFile(filename))
	return err
}

func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	return gui.editFile(file.Name)
}

func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
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

	if err := gui.fileWatcher.addFilesToFileWatcher(files); err != nil {
		return err
	}

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

func (gui *Gui) handlePullFiles(g *gocui.Gui, v *gocui.View) error {
	// if we have no upstream branch we need to set that first
	_, pullables := gui.GitCommand.GetCurrentBranchUpstreamDifferenceCount()
	currentBranchName, err := gui.GitCommand.CurrentBranchName()
	if err != nil {
		return err
	}
	if pullables == "?" {
		return gui.createPromptPanel(g, v, gui.Tr.SLocalize("EnterUpstream"), "origin/"+currentBranchName, func(g *gocui.Gui, v *gocui.View) error {
			upstream := gui.trimmedContent(v)
			if err := gui.GitCommand.SetUpstreamBranch(upstream); err != nil {
				errorMessage := err.Error()
				if strings.Contains(errorMessage, "does not exist") {
					errorMessage = fmt.Sprintf("upstream branch %s not found.\nIf you expect it to exist, you should fetch (with 'f').\nOtherwise, you should push (with 'shift+P')", upstream)
				}
				return gui.createErrorPanel(gui.g, errorMessage)
			}
			return gui.pullFiles(v)
		})
	}

	return gui.pullFiles(v)
}

func (gui *Gui) pullFiles(v *gocui.View) error {
	if err := gui.createLoaderPanel(gui.g, v, gui.Tr.SLocalize("PullWait")); err != nil {
		return err
	}

	go func() {
		unamePassOpend := false
		err := gui.GitCommand.Pull(func(passOrUname string) string {
			unamePassOpend = true
			return gui.waitForPassUname(gui.g, v, passOrUname)
		})
		gui.HandleCredentialsPopup(gui.g, unamePassOpend, err)
	}()

	return nil
}

func (gui *Gui) pushWithForceFlag(g *gocui.Gui, v *gocui.View, force bool, upstream string) error {
	if err := gui.createLoaderPanel(gui.g, v, gui.Tr.SLocalize("PushWait")); err != nil {
		return err
	}
	go func() {
		unamePassOpend := false
		branchName := gui.getCheckedOutBranch().Name
		err := gui.GitCommand.Push(branchName, force, upstream, func(passOrUname string) string {
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
	currentBranchName, err := gui.GitCommand.CurrentBranchName()
	if err != nil {
		return err
	}

	if pullables == "?" {
		return gui.createPromptPanel(g, v, gui.Tr.SLocalize("EnterUpstream"), "origin "+currentBranchName, func(g *gocui.Gui, v *gocui.View) error {
			return gui.pushWithForceFlag(g, v, false, gui.trimmedContent(v))
		})
	} else if pullables == "0" {
		return gui.pushWithForceFlag(g, v, false, "")
	}
	return gui.createConfirmationPanel(g, nil, true, gui.Tr.SLocalize("ForcePush"), gui.Tr.SLocalize("ForcePushPrompt"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.pushWithForceFlag(g, v, true, "")
	}, nil)
}

func (gui *Gui) handleSwitchToMerge(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		return nil
	}
	if !file.HasInlineMergeConflicts {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("FileNoMergeCons"))
	}
	if err := gui.changeMainViewsContext("merging"); err != nil {
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

type discardOption struct {
	handler     func(fileName *commands.File) error
	description string
}

// GetDisplayStrings is a function.
func (r *discardOption) GetDisplayStrings(isFocused bool) []string {
	return []string{r.description}
}

func (gui *Gui) handleCreateDiscardMenu(g *gocui.Gui, v *gocui.View) error {
	file, err := gui.getSelectedFile(g)
	if err != nil {
		if err != gui.Errors.ErrNoFiles {
			return err
		}
		return nil
	}

	options := []*discardOption{
		{
			description: gui.Tr.SLocalize("discardAllChanges"),
			handler: func(file *commands.File) error {
				return gui.GitCommand.DiscardAllFileChanges(file)
			},
		},
		{
			description: gui.Tr.SLocalize("cancel"),
			handler: func(file *commands.File) error {
				return nil
			},
		},
	}

	if file.HasStagedChanges && file.HasUnstagedChanges {
		discardUnstagedChanges := &discardOption{
			description: gui.Tr.SLocalize("discardUnstagedChanges"),
			handler: func(file *commands.File) error {
				return gui.GitCommand.DiscardUnstagedFileChanges(file)
			},
		}

		options = append(options[:1], append([]*discardOption{discardUnstagedChanges}, options[1:]...)...)
	}

	handleMenuPress := func(index int) error {
		file, err := gui.getSelectedFile(g)
		if err != nil {
			return err
		}

		if err := options[index].handler(file); err != nil {
			return err
		}

		return gui.refreshFiles()
	}

	return gui.createMenu(file.Name, options, len(options), handleMenuPress)
}

func (gui *Gui) handleCustomCommand(g *gocui.Gui, v *gocui.View) error {
	return gui.createPromptPanel(g, v, gui.Tr.SLocalize("CustomCommand"), "", func(g *gocui.Gui, v *gocui.View) error {
		command := gui.trimmedContent(v)
		gui.SubProcess = gui.OSCommand.RunCustomCommand(command)
		return gui.Errors.ErrSubProcess
	})
}

type stashOption struct {
	description string
	handler     func() error
}

// GetDisplayStrings is a function.
func (o *stashOption) GetDisplayStrings(isFocused bool) []string {
	return []string{o.description}
}

func (gui *Gui) handleCreateStashMenu(g *gocui.Gui, v *gocui.View) error {
	options := []*stashOption{
		{
			description: gui.Tr.SLocalize("stashAllChanges"),
			handler: func() error {
				return gui.handleStashSave(gui.GitCommand.StashSave)
			},
		},
		{
			description: gui.Tr.SLocalize("stashStagedChanges"),
			handler: func() error {
				return gui.handleStashSave(gui.GitCommand.StashSaveStagedChanges)
			},
		},
		{
			description: gui.Tr.SLocalize("cancel"),
			handler: func() error {
				return nil
			},
		},
	}

	handleMenuPress := func(index int) error {
		return options[index].handler()
	}

	return gui.createMenu(gui.Tr.SLocalize("stashOptions"), options, len(options), handleMenuPress)
}

func (gui *Gui) handleStashChanges(g *gocui.Gui, v *gocui.View) error {
	return gui.handleStashSave(gui.GitCommand.StashSave)
}
