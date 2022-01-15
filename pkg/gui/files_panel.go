package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

func (gui *Gui) filesRenderToMain() error {
	node := gui.getSelectedFileNode()

	if node == nil {
		return gui.refreshMainViews(refreshMainOpts{
			main: &viewUpdateOpts{
				title: "",
				task:  NewRenderStringTask(gui.Tr.NoChangedFiles),
			},
		})
	}

	if node.File != nil && node.File.HasInlineMergeConflicts {
		return gui.renderConflictsFromFilesPanel()
	}

	cmdObj := gui.Git.WorkingTree.WorktreeFileDiffCmdObj(node, false, !node.GetHasUnstagedChanges() && node.GetHasStagedChanges(), gui.State.IgnoreWhitespaceInDiffView)

	refreshOpts := refreshMainOpts{main: &viewUpdateOpts{
		title: gui.Tr.UnstagedChanges,
		task:  NewRunPtyTask(cmdObj.GetCmd()),
	}}

	if node.GetHasUnstagedChanges() {
		if node.GetHasStagedChanges() {
			cmdObj := gui.Git.WorkingTree.WorktreeFileDiffCmdObj(node, false, true, gui.State.IgnoreWhitespaceInDiffView)

			refreshOpts.secondary = &viewUpdateOpts{
				title: gui.Tr.StagedChanges,
				task:  NewRunPtyTask(cmdObj.GetCmd()),
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

	if err := gui.refreshStateSubmoduleConfigs(); err != nil {
		return err
	}
	if err := gui.refreshStateFiles(); err != nil {
		return err
	}

	gui.OnUIThread(func() error {
		if err := gui.postRefreshUpdate(gui.State.Contexts.Submodules); err != nil {
			gui.Log.Error(err)
		}

		if ContextKey(gui.Views.Files.Context) == FILES_CONTEXT_KEY {
			// doing this a little custom (as opposed to using gui.postRefreshUpdate) because we handle selecting the file explicitly below
			if err := gui.State.Contexts.Files.HandleRender(); err != nil {
				return err
			}
		}

		if gui.currentContext().GetKey() == FILES_CONTEXT_KEY || (gui.g.CurrentView() == gui.Views.Main && ContextKey(gui.g.CurrentView().Context) == MAIN_MERGING_CONTEXT_KEY) {
			newSelectedPath := gui.getSelectedPath()
			alreadySelected := selectedPath != "" && newSelectedPath == selectedPath
			if !alreadySelected {
				gui.takeOverMergeConflictScrolling()
			}

			gui.Views.Files.FocusPoint(0, gui.State.Panels.Files.SelectedLineIdx)
			return gui.filesRenderToMain()
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

	return gui.Git.WorkingTree.StageFile(file.Name)
}

func (gui *Gui) handleEnterFile() error {
	return gui.enterFile(OnFocusOpts{ClickedViewName: "", ClickedViewLineIdx: -1})
}

func (gui *Gui) enterFile(opts OnFocusOpts) error {
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
		return gui.switchToMerge()
	}
	if file.HasMergeConflicts {
		return gui.PopupHandler.ErrorMsg(gui.Tr.FileStagingRequirements)
	}

	return gui.pushContext(gui.State.Contexts.Staging, opts)
}

func (gui *Gui) handleFilePress() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.IsLeaf() {
		file := node.File

		if file.HasInlineMergeConflicts {
			return gui.switchToMerge()
		}

		if file.HasUnstagedChanges {
			gui.logAction(gui.Tr.Actions.StageFile)
			if err := gui.Git.WorkingTree.StageFile(file.Name); err != nil {
				return gui.PopupHandler.Error(err)
			}
		} else {
			gui.logAction(gui.Tr.Actions.UnstageFile)
			if err := gui.Git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return gui.PopupHandler.Error(err)
			}
		}
	} else {
		// if any files within have inline merge conflicts we can't stage or unstage,
		// or it'll end up with those >>>>>> lines actually staged
		if node.GetHasInlineMergeConflicts() {
			return gui.PopupHandler.ErrorMsg(gui.Tr.ErrStageDirWithInlineMergeConflicts)
		}

		if node.GetHasUnstagedChanges() {
			gui.logAction(gui.Tr.Actions.StageFile)
			if err := gui.Git.WorkingTree.StageFile(node.Path); err != nil {
				return gui.PopupHandler.Error(err)
			}
		} else {
			// pretty sure it doesn't matter that we're always passing true here
			gui.logAction(gui.Tr.Actions.UnstageFile)
			if err := gui.Git.WorkingTree.UnStageFile([]string{node.Path}, true); err != nil {
				return gui.PopupHandler.Error(err)
			}
		}
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}

	return gui.State.Contexts.Files.HandleFocus()
}

func (gui *Gui) allFilesStaged() bool {
	for _, file := range gui.State.FileManager.GetAllFiles() {
		if file.HasUnstagedChanges {
			return false
		}
	}
	return true
}

func (gui *Gui) onFocusFile() error {
	gui.takeOverMergeConflictScrolling()
	return nil
}

func (gui *Gui) handleStageAll() error {
	var err error
	if gui.allFilesStaged() {
		gui.logAction(gui.Tr.Actions.UnstageAllFiles)
		err = gui.Git.WorkingTree.UnstageAll()
	} else {
		gui.logAction(gui.Tr.Actions.StageAllFiles)
		err = gui.Git.WorkingTree.StageAll()
	}
	if err != nil {
		_ = gui.PopupHandler.Error(err)
	}

	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}

	return gui.State.Contexts.Files.HandleFocus()
}

func (gui *Gui) handleIgnoreFile() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.GetPath() == ".gitignore" {
		return gui.PopupHandler.ErrorMsg("Cannot ignore .gitignore")
	}

	unstageFiles := func() error {
		return node.ForEachFile(func(file *models.File) error {
			if file.HasStagedChanges {
				if err := gui.Git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if node.GetIsTracked() {
		return gui.PopupHandler.Ask(askOpts{
			title:  gui.Tr.IgnoreTracked,
			prompt: gui.Tr.IgnoreTrackedPrompt,
			handleConfirm: func() error {
				gui.logAction(gui.Tr.Actions.IgnoreFile)
				// not 100% sure if this is necessary but I'll assume it is
				if err := unstageFiles(); err != nil {
					return err
				}

				if err := gui.Git.WorkingTree.RemoveTrackedFiles(node.GetPath()); err != nil {
					return err
				}

				if err := gui.Git.WorkingTree.Ignore(node.GetPath()); err != nil {
					return err
				}
				return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}})
			},
		})
	}

	gui.logAction(gui.Tr.Actions.IgnoreFile)

	if err := unstageFiles(); err != nil {
		return err
	}

	if err := gui.Git.WorkingTree.Ignore(node.GetPath()); err != nil {
		return gui.PopupHandler.Error(err)
	}

	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}})
}

func (gui *Gui) handleWIPCommitPress() error {
	skipHookPrefix := gui.UserConfig.Git.SkipHookPrefix
	if skipHookPrefix == "" {
		return gui.PopupHandler.ErrorMsg(gui.Tr.SkipHookPrefixNotConfigured)
	}

	textArea := gui.Views.CommitMessage.TextArea
	textArea.Clear()
	textArea.TypeString(skipHookPrefix)
	gui.Views.CommitMessage.RenderTextArea()

	return gui.handleCommitPress()
}

func (gui *Gui) commitPrefixConfigForRepo() *config.CommitPrefixConfig {
	cfg, ok := gui.UserConfig.Git.CommitPrefixes[utils.GetCurrentRepoName()]
	if !ok {
		return nil
	}

	return &cfg
}

func (gui *Gui) prepareFilesForCommit() error {
	noStagedFiles := len(gui.stagedFiles()) == 0
	if noStagedFiles && gui.UserConfig.Gui.SkipNoStagedFilesWarning {
		gui.logAction(gui.Tr.Actions.StageAllFiles)
		err := gui.Git.WorkingTree.StageAll()
		if err != nil {
			return err
		}

		return gui.refreshFilesAndSubmodules()
	}

	return nil
}

func (gui *Gui) handleCommitPress() error {
	if err := gui.prepareFilesForCommit(); err != nil {
		return gui.PopupHandler.Error(err)
	}

	if gui.State.FileManager.GetItemsLength() == 0 {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoFilesStagedTitle)
	}

	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(gui.handleCommitPress)
	}

	if len(gui.State.failedCommitMessage) > 0 {
		gui.Views.CommitMessage.ClearTextArea()
		gui.Views.CommitMessage.TextArea.TypeString(gui.State.failedCommitMessage)
		gui.Views.CommitMessage.RenderTextArea()
	} else {
		commitPrefixConfig := gui.commitPrefixConfigForRepo()
		if commitPrefixConfig != nil {
			prefixPattern := commitPrefixConfig.Pattern
			prefixReplace := commitPrefixConfig.Replace
			rgx, err := regexp.Compile(prefixPattern)
			if err != nil {
				return gui.PopupHandler.ErrorMsg(fmt.Sprintf("%s: %s", gui.Tr.LcCommitPrefixPatternError, err.Error()))
			}
			prefix := rgx.ReplaceAllString(gui.getCheckedOutBranch().Name, prefixReplace)
			gui.Views.CommitMessage.ClearTextArea()
			gui.Views.CommitMessage.TextArea.TypeString(prefix)
			gui.Views.CommitMessage.RenderTextArea()
		}
	}

	if err := gui.pushContext(gui.State.Contexts.CommitMessage); err != nil {
		return err
	}

	gui.RenderCommitLength()
	return nil
}

func (gui *Gui) promptToStageAllAndRetry(retry func() error) error {
	return gui.PopupHandler.Ask(askOpts{
		title:  gui.Tr.NoFilesStagedTitle,
		prompt: gui.Tr.NoFilesStagedPrompt,
		handleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.StageAllFiles)
			if err := gui.Git.WorkingTree.StageAll(); err != nil {
				return gui.PopupHandler.Error(err)
			}
			if err := gui.refreshFilesAndSubmodules(); err != nil {
				return gui.PopupHandler.Error(err)
			}

			return retry()
		},
	})
}

func (gui *Gui) handleAmendCommitPress() error {
	if gui.State.FileManager.GetItemsLength() == 0 {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoFilesStagedTitle)
	}

	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(gui.handleAmendCommitPress)
	}

	if len(gui.State.Commits) == 0 {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoCommitToAmend)
	}

	return gui.PopupHandler.Ask(askOpts{
		title:  strings.Title(gui.Tr.AmendLastCommit),
		prompt: gui.Tr.SureToAmend,
		handleConfirm: func() error {
			cmdObj := gui.Git.Commit.AmendHeadCmdObj()
			gui.logAction(gui.Tr.Actions.AmendCommit)
			return gui.withGpgHandling(cmdObj, gui.Tr.AmendingStatus, nil)
		},
	})
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress() error {
	if gui.State.FileManager.GetItemsLength() == 0 {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoFilesStagedTitle)
	}

	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(gui.handleCommitEditorPress)
	}

	gui.logAction(gui.Tr.Actions.Commit)
	return gui.runSubprocessWithSuspenseAndRefresh(
		gui.Git.Commit.CommitEditorCmdObj(),
	)
}

func (gui *Gui) handleStatusFilterPressed() error {
	return gui.PopupHandler.Menu(createMenuOptions{
		title: gui.Tr.FilteringMenuTitle,
		items: []*menuItem{
			{
				displayString: gui.Tr.FilterStagedFiles,
				onPress: func() error {
					return gui.setStatusFiltering(filetree.DisplayStaged)
				},
			},
			{
				displayString: gui.Tr.FilterUnstagedFiles,
				onPress: func() error {
					return gui.setStatusFiltering(filetree.DisplayUnstaged)
				},
			},
			{
				displayString: gui.Tr.ResetCommitFilterState,
				onPress: func() error {
					return gui.setStatusFiltering(filetree.DisplayAll)
				},
			},
		},
	})
}

func (gui *Gui) setStatusFiltering(filter filetree.FileManagerDisplayFilter) error {
	state := gui.State
	state.FileManager.SetDisplayFilter(filter)
	return gui.handleRefreshFiles()
}

func (gui *Gui) editFile(filename string) error {
	return gui.editFileAtLine(filename, 1)
}

func (gui *Gui) editFileAtLine(filename string, lineNumber int) error {
	cmdStr, err := gui.Git.File.GetEditCmdStr(filename, lineNumber)
	if err != nil {
		return gui.PopupHandler.Error(err)
	}

	gui.logAction(gui.Tr.Actions.EditFile)
	return gui.runSubprocessWithSuspenseAndRefresh(
		gui.OSCommand.Cmd.NewShell(cmdStr),
	)
}

func (gui *Gui) handleFileEdit() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return gui.PopupHandler.ErrorMsg(gui.Tr.ErrCannotEditDirectory)
	}

	return gui.editFile(node.GetPath())
}

func (gui *Gui) handleFileOpen() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	return gui.openFile(node.GetPath())
}

func (gui *Gui) handleRefreshFiles() error {
	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}})
}

func (gui *Gui) refreshStateFiles() error {
	state := gui.State

	// keep track of where the cursor is currently and the current file names
	// when we refresh, go looking for a matching name
	// move the cursor to there.

	selectedNode := gui.getSelectedFileNode()

	prevNodes := gui.State.FileManager.GetAllItems()
	prevSelectedLineIdx := gui.State.Panels.Files.SelectedLineIdx

	files := gui.Git.Loaders.Files.
		GetStatusFiles(loaders.GetStatusFileOptions{})

	// for when you stage the old file of a rename and the new file is in a collapsed dir
	state.FileManager.RWMutex.Lock()
	for _, file := range files {
		if selectedNode != nil && selectedNode.Path != "" && file.PreviousName == selectedNode.Path {
			state.FileManager.ExpandToPath(file.Name)
		}
	}

	state.FileManager.SetFiles(files)
	state.FileManager.RWMutex.Unlock()

	if err := gui.fileWatcher.addFilesToFileWatcher(files); err != nil {
		return err
	}

	if selectedNode != nil {
		newIdx := gui.findNewSelectedIdx(prevNodes[prevSelectedLineIdx:], state.FileManager.GetAllItems())
		if newIdx != -1 && newIdx != prevSelectedLineIdx {
			newNode := state.FileManager.GetItemAtIndex(newIdx)
			// when not in tree mode, we show merge conflict files at the top, so you
			// can work through them one by one without having to sift through a large
			// set of files. If you have just fixed the merge conflicts of a file, we
			// actually don't want to jump to that file's new position, because that
			// file will now be ages away amidst the other files without merge
			// conflicts: the user in this case would rather work on the next file
			// with merge conflicts, which will have moved up to fill the gap left by
			// the last file, meaning the cursor doesn't need to move at all.
			leaveCursor := !state.FileManager.InTreeMode() && newNode != nil &&
				selectedNode.File != nil && selectedNode.File.HasMergeConflicts &&
				newNode.File != nil && !newNode.File.HasMergeConflicts

			if !leaveCursor {
				state.Panels.Files.SelectedLineIdx = newIdx
			}
		}
	}

	gui.refreshSelectedLine(state.Panels.Files, state.FileManager.GetItemsLength())
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

func (gui *Gui) handlePullFiles() error {
	if gui.popupPanelFocused() {
		return nil
	}

	action := gui.Tr.Actions.Pull

	currentBranch := gui.currentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return nil
	}

	// if we have no upstream branch we need to set that first
	if !currentBranch.IsTrackingRemote() {
		suggestedRemote := getSuggestedRemote(gui.State.Remotes)

		return gui.PopupHandler.Prompt(promptOpts{
			title:               gui.Tr.EnterUpstream,
			initialContent:      suggestedRemote + " " + currentBranch.Name,
			findSuggestionsFunc: gui.getRemoteBranchesSuggestionsFunc(" "),
			handleConfirm: func(upstream string) error {
				var upstreamBranch, upstreamRemote string
				split := strings.Split(upstream, " ")
				if len(split) != 2 {
					return gui.PopupHandler.ErrorMsg(gui.Tr.InvalidUpstream)
				}

				upstreamRemote = split[0]
				upstreamBranch = split[1]

				if err := gui.Git.Branch.SetCurrentBranchUpstream(upstreamRemote, upstreamBranch); err != nil {
					errorMessage := err.Error()
					if strings.Contains(errorMessage, "does not exist") {
						errorMessage = fmt.Sprintf("upstream branch %s not found.\nIf you expect it to exist, you should fetch (with 'f').\nOtherwise, you should push (with 'shift+P')", upstream)
					}
					return gui.PopupHandler.ErrorMsg(errorMessage)
				}
				return gui.pullFiles(PullFilesOptions{UpstreamRemote: upstreamRemote, UpstreamBranch: upstreamBranch, action: action})
			},
		})
	}

	return gui.pullFiles(PullFilesOptions{UpstreamRemote: currentBranch.UpstreamRemote, UpstreamBranch: currentBranch.UpstreamBranch, action: action})
}

type PullFilesOptions struct {
	UpstreamRemote  string
	UpstreamBranch  string
	FastForwardOnly bool
	action          string
}

func (gui *Gui) pullFiles(opts PullFilesOptions) error {
	return gui.PopupHandler.WithLoaderPanel(gui.Tr.PullWait, func() error {
		return gui.pullWithLock(opts)
	})
}

func (gui *Gui) pullWithLock(opts PullFilesOptions) error {
	gui.Mutexes.FetchMutex.Lock()
	defer gui.Mutexes.FetchMutex.Unlock()

	gui.logAction(opts.action)

	err := gui.Git.Sync.Pull(
		git_commands.PullOptions{
			RemoteName:      opts.UpstreamRemote,
			BranchName:      opts.UpstreamBranch,
			FastForwardOnly: opts.FastForwardOnly,
		},
	)
	if err == nil {
		_ = gui.closeConfirmationPrompt(false)
	}
	return gui.handleGenericMergeCommandResult(err)
}

type pushOpts struct {
	force          bool
	upstreamRemote string
	upstreamBranch string
	setUpstream    bool
}

func (gui *Gui) push(opts pushOpts) error {
	return gui.PopupHandler.WithLoaderPanel(gui.Tr.PushWait, func() error {
		gui.logAction(gui.Tr.Actions.Push)
		err := gui.Git.Sync.Push(git_commands.PushOpts{
			Force:          opts.force,
			UpstreamRemote: opts.upstreamRemote,
			UpstreamBranch: opts.upstreamBranch,
			SetUpstream:    opts.setUpstream,
		})

		if err != nil && !opts.force && strings.Contains(err.Error(), "Updates were rejected") {
			forcePushDisabled := gui.UserConfig.Git.DisableForcePushing
			if forcePushDisabled {
				_ = gui.PopupHandler.ErrorMsg(gui.Tr.UpdatesRejectedAndForcePushDisabled)
				return nil
			}
			_ = gui.PopupHandler.Ask(askOpts{
				title:  gui.Tr.ForcePush,
				prompt: gui.Tr.ForcePushPrompt,
				handleConfirm: func() error {
					newOpts := opts
					newOpts.force = true

					return gui.push(newOpts)
				},
			})
			return nil
		}
		gui.handleCredentialsPopup(err)
		return gui.refreshSidePanels(refreshOptions{mode: ASYNC})
	})
}

func (gui *Gui) pushFiles() error {
	if gui.popupPanelFocused() {
		return nil
	}

	// if we have pullables we'll ask if the user wants to force push
	currentBranch := gui.currentBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return nil
	}

	if currentBranch.IsTrackingRemote() {
		opts := pushOpts{
			force:          false,
			upstreamRemote: currentBranch.UpstreamRemote,
			upstreamBranch: currentBranch.UpstreamBranch,
		}
		if currentBranch.HasCommitsToPull() {
			opts.force = true
			return gui.requestToForcePush(opts)
		} else {
			return gui.push(opts)
		}
	} else {
		suggestedRemote := getSuggestedRemote(gui.State.Remotes)

		if gui.Git.Config.GetPushToCurrent() {
			return gui.push(pushOpts{setUpstream: true})
		} else {
			return gui.PopupHandler.Prompt(promptOpts{
				title:               gui.Tr.EnterUpstream,
				initialContent:      suggestedRemote + " " + currentBranch.Name,
				findSuggestionsFunc: gui.getRemoteBranchesSuggestionsFunc(" "),
				handleConfirm: func(upstream string) error {
					var upstreamBranch, upstreamRemote string
					split := strings.Split(upstream, " ")
					if len(split) == 2 {
						upstreamRemote = split[0]
						upstreamBranch = split[1]
					} else {
						upstreamRemote = upstream
						upstreamBranch = ""
					}

					return gui.push(pushOpts{
						force:          false,
						upstreamRemote: upstreamRemote,
						upstreamBranch: upstreamBranch,
						setUpstream:    true,
					})
				},
			})
		}
	}
}

func getSuggestedRemote(remotes []*models.Remote) string {
	if len(remotes) == 0 {
		return "origin"
	}

	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote.Name
		}
	}

	return remotes[0].Name
}

func (gui *Gui) requestToForcePush(opts pushOpts) error {
	forcePushDisabled := gui.UserConfig.Git.DisableForcePushing
	if forcePushDisabled {
		return gui.PopupHandler.ErrorMsg(gui.Tr.ForcePushDisabled)
	}

	return gui.PopupHandler.Ask(askOpts{
		title:  gui.Tr.ForcePush,
		prompt: gui.Tr.ForcePushPrompt,
		handleConfirm: func() error {
			return gui.push(opts)
		},
	})
}

func (gui *Gui) switchToMerge() error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	return gui.pushContext(gui.State.Contexts.Merging)
}

func (gui *Gui) openFile(filename string) error {
	gui.logAction(gui.Tr.Actions.OpenFile)
	if err := gui.OSCommand.OpenFile(filename); err != nil {
		return gui.PopupHandler.Error(err)
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

func (gui *Gui) handleCustomCommand() error {
	return gui.PopupHandler.Prompt(promptOpts{
		title:               gui.Tr.CustomCommand,
		findSuggestionsFunc: gui.getCustomCommandsHistorySuggestionsFunc(),
		handleConfirm: func(command string) error {
			gui.Config.GetAppState().CustomCommandsHistory = utils.Limit(
				utils.Uniq(
					append(gui.Config.GetAppState().CustomCommandsHistory, command),
				),
				1000,
			)

			err := gui.Config.SaveAppState()
			if err != nil {
				gui.Log.Error(err)
			}

			gui.logAction(gui.Tr.Actions.CustomCommand)
			return gui.runSubprocessWithSuspenseAndRefresh(
				gui.OSCommand.Cmd.NewShell(command),
			)
		},
	})
}

func (gui *Gui) handleCreateStashMenu() error {
	return gui.PopupHandler.Menu(createMenuOptions{
		title: gui.Tr.LcStashOptions,
		items: []*menuItem{
			{
				displayString: gui.Tr.LcStashAllChanges,
				onPress: func() error {
					gui.logAction(gui.Tr.Actions.StashAllChanges)
					return gui.handleStashSave(gui.Git.Stash.Save)
				},
			},
			{
				displayString: gui.Tr.LcStashStagedChanges,
				onPress: func() error {
					gui.logAction(gui.Tr.Actions.StashStagedChanges)
					return gui.handleStashSave(gui.Git.Stash.SaveStagedChanges)
				},
			},
		},
	})
}

func (gui *Gui) handleStashChanges() error {
	return gui.handleStashSave(gui.Git.Stash.Save)
}

func (gui *Gui) handleCreateResetToUpstreamMenu() error {
	return gui.createResetMenu("@{upstream}")
}

func (gui *Gui) handleToggleDirCollapsed() error {
	node := gui.getSelectedFileNode()
	if node == nil {
		return nil
	}

	gui.State.FileManager.ToggleCollapsed(node.GetPath())

	if err := gui.postRefreshUpdate(gui.State.Contexts.Files); err != nil {
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

	if ContextKey(gui.Views.Files.Context) == FILES_CONTEXT_KEY {
		if err := gui.State.Contexts.Files.HandleRender(); err != nil {
			return err
		}
		if err := gui.State.Contexts.Files.HandleFocus(); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) handleOpenMergeTool() error {
	return gui.PopupHandler.Ask(askOpts{
		title:  gui.Tr.MergeToolTitle,
		prompt: gui.Tr.MergeToolPrompt,
		handleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.OpenMergeTool)
			return gui.runSubprocessWithSuspenseAndRefresh(
				gui.Git.WorkingTree.OpenMergeToolCmdObj(),
			)
		},
	})
}
