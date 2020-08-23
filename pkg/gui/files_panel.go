package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getSelectedDirOrFile() (*commands.File, *commands.Dir) {
	currentView := gui.g.CurrentView()
	viewNameToCheck := ""
	if currentView != nil {
		viewNameToCheck = currentView.Name()
	}
	if viewNameToCheck != "extensiveFiles" && viewNameToCheck != "files" {
		viewNameToCheck = gui.State.PreviousView
	}

	if viewNameToCheck == "extensiveFiles" {
		selected := gui.State.Panels.ExtensiveFiles.Selected
		file, dir := gui.State.ExtensiveFiles.MatchPath(selected)

		return file, dir
	}

	selectedLine := gui.State.Panels.Files.SelectedLineIdx
	if selectedLine == -1 {
		return nil, nil
	}

	return gui.State.Files[selectedLine], nil
}

func (gui *Gui) selectFile(alreadySelected bool) error {
	gui.getFilesView().FocusPoint(0, gui.State.Panels.Files.SelectedLineIdx)

	v := gui.g.CurrentView()
	if gui.isExtensiveView(v) {
		return gui.handleExtensiveFileSelect(v, alreadySelected)
	}

	file, _ := gui.getSelectedDirOrFile()
	if file == nil {
		return gui.refreshMainViews(refreshMainOpts{
			main: &viewUpdateOpts{
				title: "",
				task:  gui.createRenderStringTask(gui.Tr.SLocalize("NoChangedFiles")),
			},
		})
	}

	if !alreadySelected {
		// TODO: pull into update task interface
		if err := gui.resetOrigin(gui.getMainView()); err != nil {
			return err
		}
		if err := gui.resetOrigin(gui.getSecondaryView()); err != nil {
			return err
		}
	}

	if file.HasInlineMergeConflicts {
		return gui.refreshMergePanel()
	}

	cmdStr := gui.GitCommand.WorktreeFileDiffCmdStr(file, false, !file.HasUnstagedChanges && file.HasStagedChanges)
	cmd := gui.OSCommand.ExecutableFromString(cmdStr)

	refreshOpts := refreshMainOpts{main: &viewUpdateOpts{
		title: gui.Tr.SLocalize("UnstagedChanges"),
		task:  gui.createRunPtyTask(cmd),
	}}

	if file.HasStagedChanges && file.HasUnstagedChanges {
		cmdStr := gui.GitCommand.WorktreeFileDiffCmdStr(file, false, true)
		cmd := gui.OSCommand.ExecutableFromString(cmdStr)

		refreshOpts.secondary = &viewUpdateOpts{
			title: gui.Tr.SLocalize("StagedChanges"),
			task:  gui.createRunPtyTask(cmd),
		}
	} else if !file.HasUnstagedChanges {
		refreshOpts.main.title = gui.Tr.SLocalize("StagedChanges")
	}

	return gui.refreshMainViews(refreshOpts)
}

func (gui *Gui) refreshFiles() error {
	gui.State.RefreshingFilesMutex.Lock()
	gui.State.IsRefreshingFiles = true
	defer func() {
		gui.State.IsRefreshingFiles = false
		gui.State.RefreshingFilesMutex.Unlock()
	}()

	isExtensiveFiles := gui.isExtensiveView(gui.g.CurrentView())

	selectedFile, selectedDir := gui.getSelectedDirOrFile()

	view := gui.getFilesView()
	if isExtensiveFiles {
		view = gui.GetExtendedFilesView()
	}

	if view == nil {
		// if the filesView hasn't been instantiated yet we just return
		return nil
	}
	if err := gui.refreshStateFiles(); err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		newSelectedFile, newSelectedDir := gui.getSelectedDirOrFile()

		if isExtensiveFiles {
			list := gui.State.ExtensiveFiles.Render(newSelectedFile, newSelectedDir)
			view.Clear()
			fmt.Fprint(view, list)
		} else if err := gui.Contexts.Files.Context.HandleRender(); err != nil {
			return err
		}

		if g.CurrentView() == view || (g.CurrentView() == gui.getMainView() && g.CurrentView().Context == MAIN_MERGING_CONTEXT_KEY) {
			alreadySelected := selectedFile != nil && newSelectedFile != nil && newSelectedFile.Name == selectedFile.Name
			return gui.selectFile(alreadySelected)
		}

		return gui.selectFile(
			(newSelectedFile != nil && selectedFile == newSelectedFile) ||
				(newSelectedDir != nil && selectedDir == newSelectedDir))
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
	result := make([]*commands.File, 0, len(files))
	for _, file := range files {
		if file.Tracked {
			result = append(result, file)
		}
	}
	return result
}

func (gui *Gui) stageSelectedFile(g *gocui.Gui) error {
	file, dir := gui.getSelectedDirOrFile()
	if file != nil {
		return gui.GitCommand.StageFile(file.Name)
	}
	if dir != nil {
		return gui.GitCommand.StageFile(dir.AbsolutePath())
	}
	return nil
}

func (gui *Gui) handleEnterFile(g *gocui.Gui, v *gocui.View) error {
	return gui.enterFile(false, -1)
}

func (gui *Gui) enterFile(forceSecondaryFocused bool, selectedLineIdx int) error {
	file, _ := gui.getSelectedDirOrFile()
	if file == nil {
		return nil
	}

	if file.HasInlineMergeConflicts {
		return gui.handleSwitchToMerge()
	}
	if file.HasMergeConflicts {
		return gui.createErrorPanel(gui.Tr.SLocalize("FileStagingRequirements"))
	}
	gui.switchContext(gui.Contexts.Staging.Context)

	return gui.refreshStagingPanel(forceSecondaryFocused, selectedLineIdx) // TODO: check if this is broken, try moving into context code
}

func (gui *Gui) handleFilePress() error {
	file, dir := gui.getSelectedDirOrFile()
	if file == nil && dir == nil {
		return nil
	}

	if file != nil {
		if file.HasInlineMergeConflicts {
			return gui.handleSwitchToMerge()
		}

		if file.HasUnstagedChanges {
			if err := gui.GitCommand.StageFile(file.Name); err != nil {
				return gui.surfaceError(err)
			}
		} else {
			if err := gui.GitCommand.UnStageFile(file.Name, file.Tracked); err != nil {
				return gui.surfaceError(err)
			}
		}
	} else {
		if dir.HasInlineMergeConflicts {
			return gui.handleSwitchToMerge()
		}

		if dir.HasUnstagedChanges {
			if err := gui.GitCommand.StageFile(dir.AbsolutePath()); err != nil {
				return gui.surfaceError(err)
			}
		} else {
			if err := gui.GitCommand.UnStageFile(dir.AbsolutePath(), dir.Tracked); err != nil {
				return gui.surfaceError(err)
			}
		}
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []int{FILES}}); err != nil {
		return err
	}

	return gui.selectFile(true)
}

func (gui *Gui) allFilesStaged(files []*commands.File) bool {
	for _, file := range files {
		if file.HasUnstagedChanges {
			return false
		}
	}
	return true
}

func (gui *Gui) focusAndSelectFile() error {
	return gui.selectFile(false)
}

func (gui *Gui) handleStageAll(g *gocui.Gui, v *gocui.View) error {
	var err error
	if gui.allFilesStaged(gui.State.Files) {
		err = gui.GitCommand.UnstageAll()
	} else {
		err = gui.GitCommand.StageAll()
	}
	if err != nil {
		_ = gui.surfaceError(err)
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []int{FILES}}); err != nil {
		return err
	}

	return gui.selectFile(false)
}

func (gui *Gui) handleIgnoreFile(g *gocui.Gui, v *gocui.View) error {
	file, dir := gui.getSelectedDirOrFile()
	if file == nil && dir == nil {
		return nil
	}

	if file != nil {
		if file.Tracked {
			return gui.ask(askOpts{
				title:  gui.Tr.SLocalize("IgnoreTracked"),
				prompt: gui.Tr.SLocalize("IgnoreTrackedPrompt"),
				handleConfirm: func() error {
					if err := gui.GitCommand.Ignore(file.Name); err != nil {
						return err
					}
					if err := gui.GitCommand.RemoveTrackedFiles(file.Name); err != nil {
						return err
					}
					return gui.refreshSidePanels(refreshOptions{scope: []int{FILES}})
				},
			})
		}

		if err := gui.GitCommand.Ignore(file.Name); err != nil {
			return gui.surfaceError(err)
		}
	} else if dir != nil {
		if dir.Tracked {
			return gui.ask(askOpts{
				title:  gui.Tr.SLocalize("IgnoreTracked"),
				prompt: gui.Tr.SLocalize("IgnoreTrackedPrompt"),
				handleConfirm: func() error {
					if err := gui.GitCommand.Ignore(dir.AbsolutePath()); err != nil {
						return err
					}
					if err := gui.GitCommand.RemoveTrackedFiles(dir.AbsolutePath()); err != nil {
						return err
					}
					return gui.refreshSidePanels(refreshOptions{scope: []int{FILES}})
				},
			})
		}

		if err := gui.GitCommand.Ignore(dir.AbsolutePath()); err != nil {
			return gui.surfaceError(err)
		}
	}

	return gui.refreshSidePanels(refreshOptions{scope: []int{FILES}})
}

func (gui *Gui) handleWIPCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	skipHookPreifx := gui.Config.GetUserConfig().GetString("git.skipHookPrefix")
	if skipHookPreifx == "" {
		return gui.createErrorPanel(gui.Tr.SLocalize("SkipHookPrefixNotConfigured"))
	}

	gui.renderStringSync("commitMessage", skipHookPreifx)
	if err := gui.getCommitMessageView().SetCursor(len(skipHookPreifx), 0); err != nil {
		return err
	}

	return gui.handleCommitPress()
}

func (gui *Gui) handleCommitPress() error {
	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(func() error {
			return gui.handleCommitPress()
		})
	}

	commitMessageView := gui.getCommitMessageView()
	prefixPattern := gui.Config.GetUserConfig().GetString("git.commitPrefixes." + utils.GetCurrentRepoName() + ".pattern")
	prefixReplace := gui.Config.GetUserConfig().GetString("git.commitPrefixes." + utils.GetCurrentRepoName() + ".replace")
	if len(prefixPattern) > 0 && len(prefixReplace) > 0 {
		rgx, err := regexp.Compile(prefixPattern)
		if err != nil {
			return gui.createErrorPanel(fmt.Sprintf("%s: %s", gui.Tr.SLocalize("commitPrefixPatternError"), err.Error()))
		}
		prefix := rgx.ReplaceAllString(gui.getCheckedOutBranch().Name, prefixReplace)
		gui.renderString("commitMessage", prefix)
		if err := commitMessageView.SetCursor(len(prefix), 0); err != nil {
			return err
		}
	}

	gui.g.Update(func(g *gocui.Gui) error {
		if err := gui.switchContext(gui.Contexts.CommitMessage.Context); err != nil {
			return err
		}

		gui.RenderCommitLength()
		return nil
	})
	return nil
}

func (gui *Gui) promptToStageAllAndRetry(retry func() error) error {
	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("NoFilesStagedTitle"),
		prompt: gui.Tr.SLocalize("NoFilesStagedPrompt"),
		handleConfirm: func() error {
			if err := gui.GitCommand.StageAll(); err != nil {
				return gui.surfaceError(err)
			}
			if err := gui.refreshFiles(); err != nil {
				return gui.surfaceError(err)
			}

			return retry()
		},
	})
}

func (gui *Gui) handleAmendCommitPress() error {
	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(func() error {
			return gui.handleAmendCommitPress()
		})
	}

	if len(gui.State.Commits) == 0 {
		return gui.createErrorPanel(gui.Tr.SLocalize("NoCommitToAmend"))
	}

	return gui.ask(askOpts{
		title:  strings.Title(gui.Tr.SLocalize("AmendLastCommit")),
		prompt: gui.Tr.SLocalize("SureToAmend"),
		handleConfirm: func() error {
			ok, err := gui.runSyncOrAsyncCommand(gui.GitCommand.AmendHead())
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}

			return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
		},
	})
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress() error {
	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(func() error {
			return gui.handleCommitEditorPress()
		})
	}

	gui.PrepareSubProcess("git", "commit")
	return nil
}

// PrepareSubProcess - prepare a subprocess for execution and tell the gui to switch to it
func (gui *Gui) PrepareSubProcess(commands ...string) {
	gui.SubProcess = gui.GitCommand.PrepareCommitSubProcess()
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})
}

func (gui *Gui) editFile(filename string) error {
	_, err := gui.runSyncOrAsyncCommand(gui.OSCommand.EditFile(filename))
	return err
}

func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {
	file, dir := gui.getSelectedDirOrFile()
	if file != nil {
		return gui.editFile(file.Name)
	}
	if dir != nil {
		return gui.editFile(dir.AbsolutePath())
	}
	return nil
}

func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {
	file, _ := gui.getSelectedDirOrFile()
	if file == nil {
		return nil
	}
	return gui.openFile(file.Name)
}

func (gui *Gui) handleRefreshFiles(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshSidePanels(refreshOptions{scope: []int{FILES}})
}

func (gui *Gui) refreshStateFiles() error {
	// keep track of where the cursor is currently and the current file names
	// when we refresh, go looking for a matching name
	// move the cursor to there.
	selectedFile, _ := gui.getSelectedDirOrFile()

	// get files to stage
	files := gui.GitCommand.GetStatusFiles(commands.GetStatusFileOptions{})
	dir := commands.FilesToTree(gui.Log, files)

	gui.State.ExtensiveFiles = dir
	gui.State.Files = gui.GitCommand.MergeStatusFiles(gui.State.Files, files, selectedFile)

	if err := gui.fileWatcher.addFilesToFileWatcher(files); err != nil {
		return err
	}

	// let's try to find our file again and move the cursor to that
	for idx, f := range gui.State.Files {
		if selectedFile != nil && f.Matches(selectedFile) {
			gui.State.Panels.Files.SelectedLineIdx = idx
			break
		}
	}

	gui.refreshSelected(&gui.State.Panels.ExtensiveFiles.Selected, dir, 0)
	gui.refreshSelectedLine(gui.State.Panels.Files, len(gui.State.Files))
	return nil
}

func (gui *Gui) handlePullFiles(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	currentBranch := gui.currentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return nil
	}

	// if we have no upstream branch we need to set that first
	if currentBranch.Pullables == "?" {
		// see if we have this branch in our config with an upstream
		conf, err := gui.GitCommand.Repo.Config()
		if err != nil {
			return gui.surfaceError(err)
		}
		for branchName, branch := range conf.Branches {
			if branchName == currentBranch.Name {
				return gui.pullFiles(PullFilesOptions{RemoteName: branch.Remote, BranchName: branch.Name})
			}
		}

		return gui.prompt(gui.Tr.SLocalize("EnterUpstream"), "origin/"+currentBranch.Name, func(upstream string) error {
			if err := gui.GitCommand.SetUpstreamBranch(upstream); err != nil {
				errorMessage := err.Error()
				if strings.Contains(errorMessage, "does not exist") {
					errorMessage = fmt.Sprintf("upstream branch %s not found.\nIf you expect it to exist, you should fetch (with 'f').\nOtherwise, you should push (with 'shift+P')", upstream)
				}
				return gui.createErrorPanel(errorMessage)
			}
			return gui.pullFiles(PullFilesOptions{})
		})
	}

	return gui.pullFiles(PullFilesOptions{})
}

type PullFilesOptions struct {
	RemoteName string
	BranchName string
}

func (gui *Gui) pullFiles(opts PullFilesOptions) error {
	if err := gui.createLoaderPanel(gui.g.CurrentView(), gui.Tr.SLocalize("PullWait")); err != nil {
		return err
	}

	mode := gui.Config.GetUserConfig().GetString("git.pull.mode")

	go gui.pullWithMode(mode, opts)

	return nil
}

func (gui *Gui) pullWithMode(mode string, opts PullFilesOptions) error {
	err := gui.GitCommand.Fetch(
		commands.FetchOptions{
			PromptUserForCredential: gui.promptUserForCredential,
			RemoteName:              opts.RemoteName,
			BranchName:              opts.BranchName,
		},
	)
	gui.handleCredentialsPopup(err)
	if err != nil {
		return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
	}

	switch mode {
	case "rebase":
		err := gui.GitCommand.RebaseBranch("FETCH_HEAD")
		return gui.handleGenericMergeCommandResult(err)
	case "merge":
		err := gui.GitCommand.Merge("FETCH_HEAD", commands.MergeOpts{})
		return gui.handleGenericMergeCommandResult(err)
	case "ff-only":
		err := gui.GitCommand.Merge("FETCH_HEAD", commands.MergeOpts{FastForwardOnly: true})
		return gui.handleGenericMergeCommandResult(err)
	default:
		return gui.createErrorPanel(fmt.Sprintf("git pull mode '%s' unrecognised", mode))
	}
}

func (gui *Gui) pushWithForceFlag(v *gocui.View, force bool, upstream string, args string) error {
	if err := gui.createLoaderPanel(v, gui.Tr.SLocalize("PushWait")); err != nil {
		return err
	}
	go func() {
		branchName := gui.getCheckedOutBranch().Name
		err := gui.GitCommand.Push(branchName, force, upstream, args, gui.promptUserForCredential)
		if err != nil && !force && strings.Contains(err.Error(), "Updates were rejected") {
			gui.ask(askOpts{
				title:  gui.Tr.SLocalize("ForcePush"),
				prompt: gui.Tr.SLocalize("ForcePushPrompt"),
				handleConfirm: func() error {
					return gui.pushWithForceFlag(v, true, upstream, args)
				},
			})

			return
		}
		gui.handleCredentialsPopup(err)
		_ = gui.refreshSidePanels(refreshOptions{mode: ASYNC})
	}()
	return nil
}

func (gui *Gui) pushFiles(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	// if we have pullables we'll ask if the user wants to force push
	currentBranch := gui.currentBranch()

	if currentBranch.Pullables == "?" {
		// see if we have this branch in our config with an upstream
		conf, err := gui.GitCommand.Repo.Config()
		if err != nil {
			return gui.surfaceError(err)
		}
		for branchName, branch := range conf.Branches {
			if branchName == currentBranch.Name {
				return gui.pushWithForceFlag(v, false, "", fmt.Sprintf("%s %s", branch.Remote, branchName))
			}
		}

		if gui.GitCommand.PushToCurrent {
			return gui.pushWithForceFlag(v, false, "", "--set-upstream")
		}

		return gui.prompt(gui.Tr.SLocalize("EnterUpstream"), "origin "+currentBranch.Name, func(response string) error {
			return gui.pushWithForceFlag(v, false, response, "")
		})
	} else if currentBranch.Pullables == "0" {
		return gui.pushWithForceFlag(v, false, "", "")
	}

	return gui.ask(askOpts{
		title:  gui.Tr.SLocalize("ForcePush"),
		prompt: gui.Tr.SLocalize("ForcePushPrompt"),
		handleConfirm: func() error {
			return gui.pushWithForceFlag(v, true, "", "")
		},
	})
}

func (gui *Gui) handleSwitchToMerge() error {
	file, dir := gui.getSelectedDirOrFile()
	if file == nil && dir == nil {
		return nil
	}

	if (file != nil && !file.HasInlineMergeConflicts) || (dir != nil && !dir.HasInlineMergeConflicts) {
		return gui.createErrorPanel(gui.Tr.SLocalize("FileNoMergeCons"))
	}

	return gui.switchContext(gui.Contexts.Merging.Context)
}

func (gui *Gui) openFile(filename string) error {
	if err := gui.OSCommand.OpenFile(filename); err != nil {
		return gui.surfaceError(err)
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

func (gui *Gui) handleCustomCommand(g *gocui.Gui, v *gocui.View) error {
	return gui.prompt(gui.Tr.SLocalize("CustomCommand"), "", func(command string) error {
		gui.SubProcess = gui.OSCommand.RunCustomCommand(command)
		return gui.Errors.ErrSubProcess
	})
}

func (gui *Gui) handleCreateStashMenu(g *gocui.Gui, v *gocui.View) error {
	menuItems := []*menuItem{
		{
			displayString: gui.Tr.SLocalize("stashAllChanges"),
			onPress: func() error {
				return gui.handleStashSave(gui.GitCommand.StashSave)
			},
		},
		{
			displayString: gui.Tr.SLocalize("stashStagedChanges"),
			onPress: func() error {
				return gui.handleStashSave(gui.GitCommand.StashSaveStagedChanges)
			},
		},
	}

	return gui.createMenu(gui.Tr.SLocalize("stashOptions"), menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) handleStashChanges(g *gocui.Gui, v *gocui.View) error {
	return gui.handleStashSave(gui.GitCommand.StashSave)
}

func (gui *Gui) handleCreateResetToUpstreamMenu(g *gocui.Gui, v *gocui.View) error {
	return gui.createResetMenu("@{upstream}")
}

func (gui *Gui) isExtensiveView(v *gocui.View) bool {
	return v != nil && v.Name() == "extensiveFiles"
}

func (gui *Gui) handleExtensiveFilesFocus(v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	cx, cy := v.Cursor()
	_, oy := v.Origin()

	// prevSelectedLine := gui.State.Panels.ExtensiveFiles.Selected
	newSelectedLine := cy - oy

	if newSelectedLine > len(gui.State.Files)-1 || len(utils.Decolorise(gui.State.Files[newSelectedLine].DisplayString)) < cx {
		return gui.selectFile(false)
	}

	gui.State.Panels.Files.SelectedLineIdx = newSelectedLine

	return nil
}

func (gui *Gui) handleCloseExtensiveView(g *gocui.Gui, filesView *gocui.View) error {
	viewNames := []string{
		"status",
		"branches",
		"commits",
		"stash",
		"files", // files needs to be last in this array to give the focus back on files
	}
	var v *gocui.View
	var err error
	for _, viewName := range viewNames {
		v, err = gui.g.SetViewOnTop(viewName)
		if err != nil {
			return err
		}
	}

	err = gui.switchFocus(gui.g.CurrentView(), v)
	if err != nil {
		return err
	}
	return gui.refreshFiles()
}

func (gui *Gui) handleOpenExtensiveView(g *gocui.Gui, filesView *gocui.View) error {
	v, err := gui.g.SetViewOnTop("extensiveFiles")
	if err != nil {
		return err
	}
	err = gui.switchFocus(gui.g.CurrentView(), v)
	if err != nil {
		return err
	}
	return gui.refreshFiles()
}

// handleFilesGoInsideFolder handles the arrow right
func (gui *Gui) handleFilesGoInsideFolder(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	dir := gui.State.ExtensiveFiles
	gui.refreshSelected(&gui.State.Panels.ExtensiveFiles.Selected, dir, 'r')

	return gui.handleExtensiveFileSelect(v, false)
}

// handleFilesGoToFolderParent handles the arrow left
func (gui *Gui) handleFilesGoToFolderParent(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	dir := gui.State.ExtensiveFiles
	gui.refreshSelected(&gui.State.Panels.ExtensiveFiles.Selected, dir, 'l')

	return gui.handleExtensiveFileSelect(v, false)
}

// handleFilesNextFileOrFolder handles the arrow down
func (gui *Gui) handleFilesNextFileOrFolder(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	dir := gui.State.ExtensiveFiles
	gui.refreshSelected(&gui.State.Panels.ExtensiveFiles.Selected, dir, 'd')

	return gui.handleExtensiveFileSelect(v, false)
}

// handleFilesPrevFileOrFolder handles the arrow up
func (gui *Gui) handleFilesPrevFileOrFolder(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	dir := gui.State.ExtensiveFiles
	gui.refreshSelected(&gui.State.Panels.ExtensiveFiles.Selected, dir, 'u')

	return gui.handleExtensiveFileSelect(v, false)
}

func (gui *Gui) handleExtensiveFileSelect(v *gocui.View, alreadySelected bool) error {
	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	file, dir := gui.State.ExtensiveFiles.MatchPath(gui.State.Panels.ExtensiveFiles.Selected)

	y := 0
	if file != nil {
		y = file.GetY()
	} else if dir != nil {
		y = dir.GetY()
	}

	gui.GetExtendedFilesView().FocusPoint(0, y)

	if file != nil {
		if file.HasInlineMergeConflicts {
			return gui.refreshMergePanel()
		}

		cmdStr := gui.GitCommand.WorktreeFileDiffCmdStr(file, false, !file.HasUnstagedChanges && file.HasStagedChanges)
		cmd := gui.OSCommand.ExecutableFromString(cmdStr)

		refreshOpts := refreshMainOpts{main: &viewUpdateOpts{
			title: gui.Tr.SLocalize("UnstagedChanges"),
			task:  gui.createRunPtyTask(cmd),
		}}

		if file.HasStagedChanges && file.HasUnstagedChanges {
			cmdStr := gui.GitCommand.WorktreeFileDiffCmdStr(file, false, true)
			cmd := gui.OSCommand.ExecutableFromString(cmdStr)

			refreshOpts.secondary = &viewUpdateOpts{
				title: gui.Tr.SLocalize("StagedChanges"),
				task:  gui.createRunPtyTask(cmd),
			}
		} else if !file.HasUnstagedChanges {
			refreshOpts.main.title = gui.Tr.SLocalize("StagedChanges")
		}

		return gui.refreshMainViews(refreshOpts)
	}

	return nil
}
