package gui

import (

	// "io"
	// "io/ioutil"

	// "strings"

	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mgutz/str"
)

// list panel functions

func (gui *Gui) getSelectedFileNode() *filetree.FileNode {
	selectedLine := gui.State.Panels.Files.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.FileManager.GetItemAtIndex(selectedLine)
}

func (gui *Gui) getSelectedFile() *models.File {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}
	return node.File
}

func (gui *Gui) getSelectedPath() string {
	node := gui.getSelectedFileNode()
	if node == nil {
		return ""
	}

	return node.GetPath()
}

func (gui *Gui) selectFile(alreadySelected bool) error {
	gui.getFilesView().FocusPoint(0, gui.State.Panels.Files.SelectedLineIdx)

	node := gui.getSelectedFileNode()

	if node == nil {
		return gui.refreshMainViews(refreshMainOpts{
			main: &viewUpdateOpts{
				title: "",
				task:  gui.createRenderStringTask(gui.Tr.NoChangedFiles),
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
		gui.takeOverMergeConflictScrolling()
	}

	if node.File != nil && node.File.HasInlineMergeConflicts {
		return gui.refreshMergePanelWithLock()
	}

	cmdStr := gui.GitCommand.WorktreeFileDiffCmdStr(node, false, !node.GetHasUnstagedChanges() && node.GetHasStagedChanges())
	cmd := gui.OSCommand.ExecutableFromString(cmdStr)

	refreshOpts := refreshMainOpts{main: &viewUpdateOpts{
		title: gui.Tr.UnstagedChanges,
		task:  gui.createRunPtyTask(cmd),
	}}

	if node.GetHasUnstagedChanges() {
		if node.GetHasStagedChanges() {
			cmdStr := gui.GitCommand.WorktreeFileDiffCmdStr(node, false, true)
			cmd := gui.OSCommand.ExecutableFromString(cmdStr)

			refreshOpts.secondary = &viewUpdateOpts{
				title: gui.Tr.StagedChanges,
				task:  gui.createRunPtyTask(cmd),
			}
		}
	} else {
		refreshOpts.main.title = gui.Tr.StagedChanges
	}

	return gui.refreshMainViews(refreshOpts)
}

func (gui *Gui) refreshFilesAndSubmodules() error {
	gui.Mutexes.RefreshingFilesMutex.Lock()
	gui.State.IsRefreshingFiles = true
	defer func() {
		gui.State.IsRefreshingFiles = false
		gui.Mutexes.RefreshingFilesMutex.Unlock()
	}()

	selectedPath := gui.getSelectedPath()

	filesView := gui.getFilesView()
	if filesView == nil {
		// if the filesView hasn't been instantiated yet we just return
		return nil
	}
	if err := gui.refreshStateSubmoduleConfigs(); err != nil {
		return err
	}
	if err := gui.refreshStateFiles(); err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		if err := gui.postRefreshUpdate(gui.Contexts.Submodules.Context); err != nil {
			gui.Log.Error(err)
		}

		if gui.getFilesView().Context == FILES_CONTEXT_KEY {
			// doing this a little custom (as opposed to using gui.postRefreshUpdate) because we handle selecting the file explicitly below
			if err := gui.Contexts.Files.Context.HandleRender(); err != nil {
				return err
			}
		}

		if gui.currentContext().GetKey() == FILES_CONTEXT_KEY || (g.CurrentView() == gui.getMainView() && g.CurrentView().Context == MAIN_MERGING_CONTEXT_KEY) {
			newSelectedPath := gui.getSelectedPath()
			alreadySelected := selectedPath != "" && newSelectedPath == selectedPath
			if err := gui.selectFile(alreadySelected); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

// specific functions

func (gui *Gui) stagedFiles() []*models.File {
	files := gui.State.FileManager.GetAllFiles()
	result := make([]*models.File, 0)
	for _, file := range files {
		if file.HasStagedChanges {
			result = append(result, file)
		}
	}
	return result
}

func (gui *Gui) trackedFiles() []*models.File {
	files := gui.State.FileManager.GetAllFiles()
	result := make([]*models.File, 0, len(files))
	for _, file := range files {
		if file.Tracked {
			result = append(result, file)
		}
	}
	return result
}

func (gui *Gui) stageSelectedFile() error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	return gui.GitCommand.StageFile(file.Name)
}

func (gui *Gui) handleEnterFile(g *gocui.Gui, v *gocui.View) error {
	return gui.enterFile(false, -1)
}

func (gui *Gui) enterFile(forceSecondaryFocused bool, selectedLineIdx int) error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return gui.handleToggleDirCollapsed()
	}

	file := node.File

	submoduleConfigs := gui.State.Submodules
	if file.IsSubmodule(submoduleConfigs) {
		submoduleConfig := file.SubmoduleConfig(submoduleConfigs)
		return gui.enterSubmodule(submoduleConfig)
	}

	if file.HasInlineMergeConflicts {
		return gui.handleSwitchToMerge()
	}
	if file.HasMergeConflicts {
		return gui.createErrorPanel(gui.Tr.FileStagingRequirements)
	}
	_ = gui.pushContext(gui.Contexts.Staging.Context)

	return gui.handleRefreshStagingPanel(forceSecondaryFocused, selectedLineIdx) // TODO: check if this is broken, try moving into context code
}

func (gui *Gui) handleFilePress() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.IsLeaf() {
		file := node.File

		if file.HasInlineMergeConflicts {
			return gui.handleSwitchToMerge()
		}

		if file.HasUnstagedChanges {
			if err := gui.GitCommand.StageFile(file.Name); err != nil {
				return gui.surfaceError(err)
			}
		} else {
			if err := gui.GitCommand.UnStageFile(file.Names(), file.Tracked); err != nil {
				return gui.surfaceError(err)
			}
		}
	} else {
		// if any files within have inline merge conflicts we can't stage or unstage,
		// or it'll end up with those >>>>>> lines actually staged
		if node.GetHasInlineMergeConflicts() {
			return gui.createErrorPanel(gui.Tr.ErrStageDirWithInlineMergeConflicts)
		}

		if node.GetHasUnstagedChanges() {
			if err := gui.GitCommand.StageFile(node.Path); err != nil {
				return gui.surfaceError(err)
			}
		} else {
			// pretty sure it doesn't matter that we're always passing true here
			if err := gui.GitCommand.UnStageFile([]string{node.Path}, true); err != nil {
				return gui.surfaceError(err)
			}
		}
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}

	return gui.selectFile(true)
}

func (gui *Gui) allFilesStaged() bool {
	for _, file := range gui.State.FileManager.GetAllFiles() {
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
	if gui.allFilesStaged() {
		err = gui.GitCommand.UnstageAll()
	} else {
		err = gui.GitCommand.StageAll()
	}
	if err != nil {
		_ = gui.surfaceError(err)
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}

	return gui.selectFile(false)
}

func (gui *Gui) handleIgnoreFile() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.GetPath() == ".gitignore" {
		return gui.createErrorPanel("Cannot ignore .gitignore")
	}

	unstageFiles := func() error {
		return node.ForEachFile(func(file *models.File) error {
			if file.HasStagedChanges {
				if err := gui.GitCommand.UnStageFile(file.Names(), file.Tracked); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if node.GetIsTracked() {
		return gui.ask(askOpts{
			title:  gui.Tr.IgnoreTracked,
			prompt: gui.Tr.IgnoreTrackedPrompt,
			handleConfirm: func() error {
				// not 100% sure if this is necessary but I'll assume it is
				if err := unstageFiles(); err != nil {
					return err
				}

				if err := gui.GitCommand.RemoveTrackedFiles(node.GetPath()); err != nil {
					return err
				}

				if err := gui.GitCommand.Ignore(node.GetPath()); err != nil {
					return err
				}
				return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}})
			},
		})
	}

	if err := unstageFiles(); err != nil {
		return err
	}

	if err := gui.GitCommand.Ignore(node.GetPath()); err != nil {
		return gui.surfaceError(err)
	}

	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}})
}

func (gui *Gui) handleWIPCommitPress(g *gocui.Gui, filesView *gocui.View) error {
	skipHookPreifx := gui.Config.GetUserConfig().Git.SkipHookPrefix
	if skipHookPreifx == "" {
		return gui.createErrorPanel(gui.Tr.SkipHookPrefixNotConfigured)
	}

	_ = gui.renderStringSync("commitMessage", skipHookPreifx)
	if err := gui.getCommitMessageView().SetCursor(len(skipHookPreifx), 0); err != nil {
		return err
	}

	return gui.handleCommitPress()
}

func (gui *Gui) commitPrefixConfigForRepo() *config.CommitPrefixConfig {
	cfg, ok := gui.Config.GetUserConfig().Git.CommitPrefixes[utils.GetCurrentRepoName()]
	if !ok {
		return nil
	}

	return &cfg
}

func (gui *Gui) prepareFilesForCommit() error {
	noStagedFiles := len(gui.stagedFiles()) == 0
	if noStagedFiles && gui.Config.GetUserConfig().Gui.SkipNoStagedFilesWarning {
		err := gui.GitCommand.StageAll()
		if err != nil {
			return err
		}

		return gui.refreshFilesAndSubmodules()
	}

	return nil
}

func (gui *Gui) handleCommitPress() error {
	if err := gui.prepareFilesForCommit(); err != nil {
		return gui.surfaceError(err)
	}

	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(gui.handleCommitPress)
	}

	commitMessageView := gui.getCommitMessageView()
	commitPrefixConfig := gui.commitPrefixConfigForRepo()
	if commitPrefixConfig != nil {
		prefixPattern := commitPrefixConfig.Pattern
		prefixReplace := commitPrefixConfig.Replace
		rgx, err := regexp.Compile(prefixPattern)
		if err != nil {
			return gui.createErrorPanel(fmt.Sprintf("%s: %s", gui.Tr.LcCommitPrefixPatternError, err.Error()))
		}
		prefix := rgx.ReplaceAllString(gui.getCheckedOutBranch().Name, prefixReplace)
		gui.renderString("commitMessage", prefix)
		if err := commitMessageView.SetCursor(len(prefix), 0); err != nil {
			return err
		}
	}

	gui.g.Update(func(g *gocui.Gui) error {
		if err := gui.pushContext(gui.Contexts.CommitMessage.Context); err != nil {
			return err
		}

		gui.RenderCommitLength()
		return nil
	})
	return nil
}

func (gui *Gui) promptToStageAllAndRetry(retry func() error) error {
	return gui.ask(askOpts{
		title:  gui.Tr.NoFilesStagedTitle,
		prompt: gui.Tr.NoFilesStagedPrompt,
		handleConfirm: func() error {
			if err := gui.GitCommand.StageAll(); err != nil {
				return gui.surfaceError(err)
			}
			if err := gui.refreshFilesAndSubmodules(); err != nil {
				return gui.surfaceError(err)
			}

			return retry()
		},
	})
}

func (gui *Gui) handleAmendCommitPress() error {
	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(gui.handleAmendCommitPress)
	}

	if len(gui.State.Commits) == 0 {
		return gui.createErrorPanel(gui.Tr.NoCommitToAmend)
	}

	return gui.ask(askOpts{
		title:  strings.Title(gui.Tr.AmendLastCommit),
		prompt: gui.Tr.SureToAmend,
		handleConfirm: func() error {
			return gui.WithWaitingStatus(gui.Tr.AmendingStatus, func() error {
				ok, err := gui.runSyncOrAsyncCommand(gui.GitCommand.AmendHead())
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}

				return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
			})
		},
	})
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress() error {
	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(gui.handleCommitEditorPress)
	}

	gui.PrepareSubProcess("git commit")
	return nil
}

// PrepareSubProcess - prepare a subprocess for execution and tell the gui to switch to it
func (gui *Gui) PrepareSubProcess(command string) {
	splitCmd := str.ToArgv(command)
	gui.SubProcess = gui.OSCommand.PrepareSubProcess(splitCmd[0], splitCmd[1:]...)
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})
}

func (gui *Gui) PrepareShellSubProcess(command string) {
	gui.SubProcess = gui.OSCommand.PrepareShellSubProcess(command)
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})
}

func (gui *Gui) editFile(filename string) error {
	_, err := gui.runSyncOrAsyncCommand(gui.GitCommand.EditFile(filename))
	return err
}

func (gui *Gui) handleFileEdit(g *gocui.Gui, v *gocui.View) error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return gui.createErrorPanel(gui.Tr.ErrCannotEditDirectory)
	}

	return gui.editFile(node.GetPath())
}

func (gui *Gui) handleFileOpen(g *gocui.Gui, v *gocui.View) error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	return gui.openFile(node.GetPath())
}

func (gui *Gui) handleRefreshFiles(g *gocui.Gui, v *gocui.View) error {
	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}})
}

func (gui *Gui) refreshStateFiles() error {
	// keep track of where the cursor is currently and the current file names
	// when we refresh, go looking for a matching name
	// move the cursor to there.

	selectedNode := gui.getSelectedFileNode()

	prevNodes := gui.State.FileManager.GetAllItems()
	prevSelectedLineIdx := gui.State.Panels.Files.SelectedLineIdx

	files := gui.GitCommand.GetStatusFiles(commands.GetStatusFileOptions{})

	// for when you stage the old file of a rename and the new file is in a collapsed dir
	for _, file := range files {
		if selectedNode != nil && selectedNode.Path != "" && file.PreviousName == selectedNode.Path {
			gui.State.FileManager.ExpandToPath(file.Name)
		}
	}

	gui.State.FileManager.SetFiles(files)

	if err := gui.fileWatcher.addFilesToFileWatcher(files); err != nil {
		return err
	}

	if selectedNode != nil {
		newIdx := gui.findNewSelectedIdx(prevNodes[prevSelectedLineIdx:], gui.State.FileManager.GetAllItems())
		if newIdx != -1 && newIdx != prevSelectedLineIdx {
			newNode := gui.State.FileManager.GetItemAtIndex(newIdx)
			// when not in tree mode, we show merge conflict files at the top, so you
			// can work through them one by one without having to sift through a large
			// set of files. If you have just fixed the merge conflicts of a file, we
			// actually don't want to jump to that file's new position, because that
			// file will now be ages away amidst the other files without merge
			// conflicts: the user in this case would rather work on the next file
			// with merge conflicts, which will have moved up to fill the gap left by
			// the last file, meaning the cursor doesn't need to move at all.
			leaveCursor := !gui.State.FileManager.InTreeMode() && newNode != nil &&
				selectedNode.File != nil && selectedNode.File.HasMergeConflicts &&
				newNode.File != nil && !newNode.File.HasMergeConflicts

			if !leaveCursor {
				gui.State.Panels.Files.SelectedLineIdx = newIdx
			}
		}
	}

	gui.refreshSelectedLine(gui.State.Panels.Files, gui.State.FileManager.GetItemsLength())
	return nil
}

// Let's try to find our file again and move the cursor to that.
// If we can't find our file, it was probably just removed by the user. In that
// case, we go looking for where the next file has been moved to. Given that the
// user could have removed a whole directory, we continue iterating through the old
// nodes until we find one that exists in the new set of nodes, then move the cursor
// to that.
// prevNodes starts from our previously selected node because we don't need to consider anything above that
func (gui *Gui) findNewSelectedIdx(prevNodes []*filetree.FileNode, currNodes []*filetree.FileNode) int {
	getPaths := func(node *filetree.FileNode) []string {
		if node == nil {
			return nil
		}
		if node.File != nil && node.File.IsRename() {
			return node.File.Names()
		} else {
			return []string{node.Path}
		}
	}

	for _, prevNode := range prevNodes {
		selectedPaths := getPaths(prevNode)

		for idx, node := range currNodes {
			paths := getPaths(node)

			// If you started off with a rename selected, and now it's broken in two, we want you to jump to the new file, not the old file.
			// This is because the new should be in the same position as the rename was meaning less cursor jumping
			foundOldFileInRename := prevNode.File != nil && prevNode.File.IsRename() && node.Path == prevNode.File.PreviousName
			foundNode := utils.StringArraysOverlap(paths, selectedPaths) && !foundOldFileInRename
			if foundNode {
				return idx
			}
		}
	}

	return -1
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

		return gui.prompt(promptOpts{
			title:          gui.Tr.EnterUpstream,
			initialContent: "origin/" + currentBranch.Name,
			handleConfirm: func(upstream string) error {
				if err := gui.GitCommand.SetUpstreamBranch(upstream); err != nil {
					errorMessage := err.Error()
					if strings.Contains(errorMessage, "does not exist") {
						errorMessage = fmt.Sprintf("upstream branch %s not found.\nIf you expect it to exist, you should fetch (with 'f').\nOtherwise, you should push (with 'shift+P')", upstream)
					}
					return gui.createErrorPanel(errorMessage)
				}
				return gui.pullFiles(PullFilesOptions{})
			},
		})
	}

	return gui.pullFiles(PullFilesOptions{})
}

type PullFilesOptions struct {
	RemoteName string
	BranchName string
}

func (gui *Gui) pullFiles(opts PullFilesOptions) error {
	if err := gui.createLoaderPanel(gui.Tr.PullWait); err != nil {
		return err
	}

	mode := gui.Config.GetUserConfig().Git.Pull.Mode

	go utils.Safe(func() { _ = gui.pullWithMode(mode, opts) })

	return nil
}

func (gui *Gui) pullWithMode(mode string, opts PullFilesOptions) error {
	gui.Mutexes.FetchMutex.Lock()
	defer gui.Mutexes.FetchMutex.Unlock()

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
	if err := gui.createLoaderPanel(gui.Tr.PushWait); err != nil {
		return err
	}
	go utils.Safe(func() {
		branchName := gui.getCheckedOutBranch().Name
		err := gui.GitCommand.Push(branchName, force, upstream, args, gui.promptUserForCredential)
		if err != nil && !force && strings.Contains(err.Error(), "Updates were rejected") {
			forcePushDisabled := gui.Config.GetUserConfig().Git.DisableForcePushing
			if forcePushDisabled {
				_ = gui.createErrorPanel(gui.Tr.UpdatesRejectedAndForcePushDisabled)
				return
			}
			_ = gui.ask(askOpts{
				title:  gui.Tr.ForcePush,
				prompt: gui.Tr.ForcePushPrompt,
				handleConfirm: func() error {
					return gui.pushWithForceFlag(v, true, upstream, args)
				},
			})
			return
		}
		gui.handleCredentialsPopup(err)
		_ = gui.refreshSidePanels(refreshOptions{mode: ASYNC})
	})
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
		} else {
			return gui.prompt(promptOpts{
				title:          gui.Tr.EnterUpstream,
				initialContent: "origin " + currentBranch.Name,
				handleConfirm: func(response string) error {
					return gui.pushWithForceFlag(v, false, response, "")
				},
			})
		}
	} else if currentBranch.Pullables == "0" {
		return gui.pushWithForceFlag(v, false, "", "")
	}

	forcePushDisabled := gui.Config.GetUserConfig().Git.DisableForcePushing
	if forcePushDisabled {
		return gui.createErrorPanel(gui.Tr.ForcePushDisabled)
	}

	return gui.ask(askOpts{
		title:  gui.Tr.ForcePush,
		prompt: gui.Tr.ForcePushPrompt,
		handleConfirm: func() error {
			return gui.pushWithForceFlag(v, true, "", "")
		},
	})
}

func (gui *Gui) handleSwitchToMerge() error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	if !file.HasInlineMergeConflicts {
		return gui.createErrorPanel(gui.Tr.FileNoMergeCons)
	}

	return gui.pushContext(gui.Contexts.Merging.Context)
}

func (gui *Gui) openFile(filename string) error {
	if err := gui.OSCommand.OpenFile(filename); err != nil {
		return gui.surfaceError(err)
	}
	return nil
}

func (gui *Gui) anyFilesWithMergeConflicts() bool {
	for _, file := range gui.State.FileManager.GetAllFiles() {
		if file.HasMergeConflicts {
			return true
		}
	}
	return false
}

func (gui *Gui) handleCustomCommand(g *gocui.Gui, v *gocui.View) error {
	return gui.prompt(promptOpts{
		title: gui.Tr.CustomCommand,
		handleConfirm: func(command string) error {
			gui.SubProcess = gui.OSCommand.PrepareShellSubProcess(command)
			return gui.Errors.ErrSubProcess
		},
	})
}

func (gui *Gui) handleCreateStashMenu(g *gocui.Gui, v *gocui.View) error {
	menuItems := []*menuItem{
		{
			displayString: gui.Tr.LcStashAllChanges,
			onPress: func() error {
				return gui.handleStashSave(gui.GitCommand.StashSave)
			},
		},
		{
			displayString: gui.Tr.LcStashStagedChanges,
			onPress: func() error {
				return gui.handleStashSave(gui.GitCommand.StashSaveStagedChanges)
			},
		},
	}

	return gui.createMenu(gui.Tr.LcStashOptions, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) handleStashChanges(g *gocui.Gui, v *gocui.View) error {
	return gui.handleStashSave(gui.GitCommand.StashSave)
}

func (gui *Gui) handleCreateResetToUpstreamMenu(g *gocui.Gui, v *gocui.View) error {
	return gui.createResetMenu("@{upstream}")
}

func (gui *Gui) handleToggleDirCollapsed() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	gui.State.FileManager.ToggleCollapsed(node.GetPath())

	if err := gui.postRefreshUpdate(gui.Contexts.Files.Context); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) handleToggleFileTreeView() error {
	// get path of currently selected file
	path := gui.getSelectedPath()

	gui.State.FileManager.ToggleShowTree()

	// find that same node in the new format and move the cursor to it
	if path != "" {
		gui.State.FileManager.ExpandToPath(path)
		index, found := gui.State.FileManager.GetIndexForPath(path)
		if found {
			gui.filesListContext().GetPanelState().SetSelectedLineIdx(index)
		}
	}

	if gui.getFilesView().Context == FILES_CONTEXT_KEY {
		if err := gui.Contexts.Files.Context.HandleRender(); err != nil {
			return err
		}
	}

	return nil
}
