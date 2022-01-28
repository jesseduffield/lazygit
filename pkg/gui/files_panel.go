package gui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedFileNode() *filetree.FileNode {
	selectedLine := gui.State.Panels.Files.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.FileTreeViewModel.GetItemAtIndex(selectedLine)
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
		ok, err := gui.setConflictsAndRenderWithLock(node.GetPath(), false)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	gui.resetMergeStateWithLock()

	cmdObj := gui.Git.WorkingTree.WorktreeFileDiffCmdObj(node, false, !node.GetHasUnstagedChanges() && node.GetHasStagedChanges(), gui.IgnoreWhitespaceInDiffView)

	refreshOpts := refreshMainOpts{main: &viewUpdateOpts{
		title: gui.Tr.UnstagedChanges,
		task:  NewRunPtyTask(cmdObj.GetCmd()),
	}}

	if node.GetHasUnstagedChanges() {
		if node.GetHasStagedChanges() {
			cmdObj := gui.Git.WorkingTree.WorktreeFileDiffCmdObj(node, false, true, gui.IgnoreWhitespaceInDiffView)

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

	prevSelectedPath := gui.getSelectedPath()

	if err := gui.refreshStateSubmoduleConfigs(); err != nil {
		return err
	}

	if err := gui.refreshMergeState(); err != nil {
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

		if gui.currentContext().GetKey() == FILES_CONTEXT_KEY {
			currentSelectedPath := gui.getSelectedPath()
			alreadySelected := prevSelectedPath != "" && currentSelectedPath == prevSelectedPath
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
	files := gui.State.FileTreeViewModel.GetAllFiles()
	result := make([]*models.File, 0)
	for _, file := range files {
		if file.HasStagedChanges {
			result = append(result, file)
		}
	}
	return result
}

func (gui *Gui) trackedFiles() []*models.File {
	files := gui.State.FileTreeViewModel.GetAllFiles()
	result := make([]*models.File, 0, len(files))
	for _, file := range files {
		if file.Tracked {
			result = append(result, file)
		}
	}
	return result
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

	if err := gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}}); err != nil {
		return err
	}

	return gui.State.Contexts.Files.HandleFocus()
}

func (gui *Gui) allFilesStaged() bool {
	for _, file := range gui.State.FileTreeViewModel.GetAllFiles() {
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

	if err := gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}}); err != nil {
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
		return gui.PopupHandler.Ask(popup.AskOpts{
			Title:  gui.Tr.IgnoreTracked,
			Prompt: gui.Tr.IgnoreTrackedPrompt,
			HandleConfirm: func() error {
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
				return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
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

	return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
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

	if gui.State.FileTreeViewModel.GetItemsLength() == 0 {
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
	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.NoFilesStagedTitle,
		Prompt: gui.Tr.NoFilesStagedPrompt,
		HandleConfirm: func() error {
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
	if gui.State.FileTreeViewModel.GetItemsLength() == 0 {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoFilesStagedTitle)
	}

	if len(gui.stagedFiles()) == 0 {
		return gui.promptToStageAllAndRetry(gui.handleAmendCommitPress)
	}

	if len(gui.State.Commits) == 0 {
		return gui.PopupHandler.ErrorMsg(gui.Tr.NoCommitToAmend)
	}

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  strings.Title(gui.Tr.AmendLastCommit),
		Prompt: gui.Tr.SureToAmend,
		HandleConfirm: func() error {
			cmdObj := gui.Git.Commit.AmendHeadCmdObj()
			gui.logAction(gui.Tr.Actions.AmendCommit)
			return gui.withGpgHandling(cmdObj, gui.Tr.AmendingStatus, nil)
		},
	})
}

// handleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (gui *Gui) handleCommitEditorPress() error {
	if gui.State.FileTreeViewModel.GetItemsLength() == 0 {
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
	return gui.PopupHandler.Menu(popup.CreateMenuOptions{
		Title: gui.Tr.FilteringMenuTitle,
		Items: []*popup.MenuItem{
			{
				DisplayString: gui.Tr.FilterStagedFiles,
				OnPress: func() error {
					return gui.setStatusFiltering(filetree.DisplayStaged)
				},
			},
			{
				DisplayString: gui.Tr.FilterUnstagedFiles,
				OnPress: func() error {
					return gui.setStatusFiltering(filetree.DisplayUnstaged)
				},
			},
			{
				DisplayString: gui.Tr.ResetCommitFilterState,
				OnPress: func() error {
					return gui.setStatusFiltering(filetree.DisplayAll)
				},
			},
		},
	})
}

func (gui *Gui) setStatusFiltering(filter filetree.FileTreeDisplayFilter) error {
	state := gui.State
	state.FileTreeViewModel.SetFilter(filter)
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
	return gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
}

func (gui *Gui) refreshStateFiles() error {
	state := gui.State

	// keep track of where the cursor is currently and the current file names
	// when we refresh, go looking for a matching name
	// move the cursor to there.

	selectedNode := gui.getSelectedFileNode()

	prevNodes := gui.State.FileTreeViewModel.GetAllItems()
	prevSelectedLineIdx := gui.State.Panels.Files.SelectedLineIdx

	// If git thinks any of our files have inline merge conflicts, but they actually don't,
	// we stage them.
	// Note that if files with merge conflicts have both arisen and have been resolved
	// between refreshes, we won't stage them here. This is super unlikely though,
	// and this approach spares us from having to call `git status` twice in a row.
	// Although this also means that at startup we won't be staging anything until
	// we call git status again.
	pathsToStage := []string{}
	prevConflictFileCount := 0
	for _, file := range state.FileTreeViewModel.GetAllFiles() {
		if file.HasMergeConflicts {
			prevConflictFileCount++
		}
		if file.HasInlineMergeConflicts {
			hasConflicts, err := mergeconflicts.FileHasConflictMarkers(file.Name)
			if err != nil {
				gui.Log.Error(err)
			} else if !hasConflicts {
				pathsToStage = append(pathsToStage, file.Name)
			}
		}
	}

	if len(pathsToStage) > 0 {
		gui.logAction(gui.Tr.Actions.StageResolvedFiles)
		if err := gui.Git.WorkingTree.StageFiles(pathsToStage); err != nil {
			return gui.surfaceError(err)
		}
	}

	files := gui.Git.Loaders.Files.
		GetStatusFiles(loaders.GetStatusFileOptions{})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	if gui.Git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE && conflictFileCount == 0 && prevConflictFileCount > 0 {
		gui.OnUIThread(func() error { return gui.promptToContinueRebase() })
	}

	// for when you stage the old file of a rename and the new file is in a collapsed dir
	state.FileTreeViewModel.RWMutex.Lock()
	for _, file := range files {
		if selectedNode != nil && selectedNode.Path != "" && file.PreviousName == selectedNode.Path {
			state.FileTreeViewModel.ExpandToPath(file.Name)
		}
	}

	// only taking over the filter if it hasn't already been set by the user.
	// Though this does make it impossible for the user to actually say they want to display all if
	// conflicts are currently being shown. Hmm. Worth it I reckon. If we need to add some
	// extra state here to see if the user's set the filter themselves we can do that, but
	// I'd prefer to maintain as little state as possible.
	if conflictFileCount > 0 {
		if state.FileTreeViewModel.GetFilter() == filetree.DisplayAll {
			state.FileTreeViewModel.SetFilter(filetree.DisplayConflicted)
		}
	} else if state.FileTreeViewModel.GetFilter() == filetree.DisplayConflicted {
		state.FileTreeViewModel.SetFilter(filetree.DisplayAll)
	}

	state.FileTreeViewModel.SetFiles(files)
	state.FileTreeViewModel.RWMutex.Unlock()

	if err := gui.fileWatcher.addFilesToFileWatcher(files); err != nil {
		return err
	}

	if selectedNode != nil {
		newIdx := gui.findNewSelectedIdx(prevNodes[prevSelectedLineIdx:], state.FileTreeViewModel.GetAllItems())
		if newIdx != -1 && newIdx != prevSelectedLineIdx {
			newNode := state.FileTreeViewModel.GetItemAtIndex(newIdx)
			// when not in tree mode, we show merge conflict files at the top, so you
			// can work through them one by one without having to sift through a large
			// set of files. If you have just fixed the merge conflicts of a file, we
			// actually don't want to jump to that file's new position, because that
			// file will now be ages away amidst the other files without merge
			// conflicts: the user in this case would rather work on the next file
			// with merge conflicts, which will have moved up to fill the gap left by
			// the last file, meaning the cursor doesn't need to move at all.
			leaveCursor := !state.FileTreeViewModel.InTreeMode() && newNode != nil &&
				selectedNode.File != nil && selectedNode.File.HasMergeConflicts &&
				newNode.File != nil && !newNode.File.HasMergeConflicts

			if !leaveCursor {
				state.Panels.Files.SelectedLineIdx = newIdx
			}
		}
	}

	gui.refreshSelectedLine(state.Panels.Files, state.FileTreeViewModel.GetItemsLength())
	return nil
}

// promptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (gui *Gui) promptToContinueRebase() error {
	gui.takeOverMergeConflictScrolling()

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  "continue",
		Prompt: gui.Tr.ConflictsResolved,
		HandleConfirm: func() error {
			return gui.genericMergeCommand(REBASE_OPTION_CONTINUE)
		},
	})
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

		return gui.PopupHandler.Prompt(popup.PromptOpts{
			Title:               gui.Tr.EnterUpstream,
			InitialContent:      suggestedRemote + " " + currentBranch.Name,
			FindSuggestionsFunc: gui.getRemoteBranchesSuggestionsFunc(" "),
			HandleConfirm: func(upstream string) error {
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

		if err != nil {
			if !opts.force && strings.Contains(err.Error(), "Updates were rejected") {
				forcePushDisabled := gui.UserConfig.Git.DisableForcePushing
				if forcePushDisabled {
					_ = gui.PopupHandler.ErrorMsg(gui.Tr.UpdatesRejectedAndForcePushDisabled)
					return nil
				}
				_ = gui.PopupHandler.Ask(popup.AskOpts{
					Title:  gui.Tr.ForcePush,
					Prompt: gui.Tr.ForcePushPrompt,
					HandleConfirm: func() error {
						newOpts := opts
						newOpts.force = true

						return gui.push(newOpts)
					},
				})
				return nil
			}
			_ = gui.PopupHandler.Error(err)
		}
		return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC})
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
			return gui.PopupHandler.Prompt(popup.PromptOpts{
				Title:               gui.Tr.EnterUpstream,
				InitialContent:      suggestedRemote + " " + currentBranch.Name,
				FindSuggestionsFunc: gui.getRemoteBranchesSuggestionsFunc(" "),
				HandleConfirm: func(upstream string) error {
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

	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.ForcePush,
		Prompt: gui.Tr.ForcePushPrompt,
		HandleConfirm: func() error {
			return gui.push(opts)
		},
	})
}

func (gui *Gui) switchToMerge() error {
	file := gui.getSelectedFile()
	if file == nil {
		return nil
	}

	gui.takeOverMergeConflictScrolling()

	if gui.State.Panels.Merging.GetPath() != file.Name {
		hasConflicts, err := gui.setMergeStateWithLock(file.Name)
		if err != nil {
			return err
		}
		if !hasConflicts {
			return nil
		}
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

func (gui *Gui) handleCustomCommand() error {
	return gui.PopupHandler.Prompt(popup.PromptOpts{
		Title:               gui.Tr.CustomCommand,
		FindSuggestionsFunc: gui.getCustomCommandsHistorySuggestionsFunc(),
		HandleConfirm: func(command string) error {
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
	return gui.PopupHandler.Menu(popup.CreateMenuOptions{
		Title: gui.Tr.LcStashOptions,
		Items: []*popup.MenuItem{
			{
				DisplayString: gui.Tr.LcStashAllChanges,
				OnPress: func() error {
					gui.logAction(gui.Tr.Actions.StashAllChanges)
					return gui.handleStashSave(gui.Git.Stash.Save)
				},
			},
			{
				DisplayString: gui.Tr.LcStashStagedChanges,
				OnPress: func() error {
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

	gui.State.FileTreeViewModel.ToggleCollapsed(node.GetPath())

	if err := gui.postRefreshUpdate(gui.State.Contexts.Files); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) handleToggleFileTreeView() error {
	// get path of currently selected file
	path := gui.getSelectedPath()

	gui.State.FileTreeViewModel.ToggleShowTree()

	// find that same node in the new format and move the cursor to it
	if path != "" {
		gui.State.FileTreeViewModel.ExpandToPath(path)
		index, found := gui.State.FileTreeViewModel.GetIndexForPath(path)
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
	return gui.PopupHandler.Ask(popup.AskOpts{
		Title:  gui.Tr.MergeToolTitle,
		Prompt: gui.Tr.MergeToolPrompt,
		HandleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.OpenMergeTool)
			return gui.runSubprocessWithSuspenseAndRefresh(
				gui.Git.WorkingTree.OpenMergeToolCmdObj(),
			)
		},
	})
}

func (gui *Gui) resetSubmodule(submodule *models.SubmoduleConfig) error {
	return gui.PopupHandler.WithWaitingStatus(gui.Tr.LcResettingSubmoduleStatus, func() error {
		gui.logAction(gui.Tr.Actions.ResetSubmodule)

		file := gui.fileForSubmodule(submodule)
		if file != nil {
			if err := gui.Git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return gui.PopupHandler.Error(err)
			}
		}

		if err := gui.Git.Submodule.Stash(submodule); err != nil {
			return gui.PopupHandler.Error(err)
		}
		if err := gui.Git.Submodule.Reset(submodule); err != nil {
			return gui.PopupHandler.Error(err)
		}

		return gui.refreshSidePanels(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.SUBMODULES}})
	})
}

func (gui *Gui) fileForSubmodule(submodule *models.SubmoduleConfig) *models.File {
	for _, file := range gui.State.FileManager.GetAllFiles() {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}
