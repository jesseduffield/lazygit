package gui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func getScopeNames(scopes []types.RefreshableView) []string {
	scopeNameMap := map[types.RefreshableView]string{
		types.COMMITS:     "commits",
		types.BRANCHES:    "branches",
		types.FILES:       "files",
		types.SUBMODULES:  "submodules",
		types.STASH:       "stash",
		types.REFLOG:      "reflog",
		types.TAGS:        "tags",
		types.REMOTES:     "remotes",
		types.STATUS:      "status",
		types.BISECT_INFO: "bisect",
	}

	scopeNames := make([]string, len(scopes))
	for i, scope := range scopes {
		scopeNames[i] = scopeNameMap[scope]
	}

	return scopeNames
}

func getModeName(mode types.RefreshMode) string {
	switch mode {
	case types.SYNC:
		return "sync"
	case types.ASYNC:
		return "async"
	case types.BLOCK_UI:
		return "block-ui"
	default:
		return "unknown mode"
	}
}

func arrToMap(arr []types.RefreshableView) map[types.RefreshableView]bool {
	output := map[types.RefreshableView]bool{}
	for _, el := range arr {
		output[el] = true
	}
	return output
}

func (gui *Gui) Refresh(options types.RefreshOptions) error {
	if options.Scope == nil {
		gui.c.Log.Infof(
			"refreshing all scopes in %s mode",
			getModeName(options.Mode),
		)
	} else {
		gui.c.Log.Infof(
			"refreshing the following scopes in %s mode: %s",
			getModeName(options.Mode),
			strings.Join(getScopeNames(options.Scope), ","),
		)
	}

	wg := sync.WaitGroup{}

	f := func() {
		var scopeMap map[types.RefreshableView]bool
		if len(options.Scope) == 0 {
			scopeMap = arrToMap([]types.RefreshableView{
				types.COMMITS,
				types.BRANCHES,
				types.FILES,
				types.STASH,
				types.REFLOG,
				types.TAGS,
				types.REMOTES,
				types.STATUS,
				types.BISECT_INFO,
			})
		} else {
			scopeMap = arrToMap(options.Scope)
		}

		refresh := func(f func()) {
			wg.Add(1)
			func() {
				if options.Mode == types.ASYNC {
					go utils.Safe(f)
				} else {
					f()
				}
				wg.Done()
			}()
		}

		if scopeMap[types.COMMITS] || scopeMap[types.BRANCHES] || scopeMap[types.REFLOG] || scopeMap[types.BISECT_INFO] {
			refresh(gui.refreshCommits)
		} else if scopeMap[types.REBASE_COMMITS] {
			// the above block handles rebase commits so we only need to call this one
			// if we've asked specifically for rebase commits and not those other things
			refresh(func() { _ = gui.refreshRebaseCommits() })
		}

		if scopeMap[types.FILES] || scopeMap[types.SUBMODULES] {
			refresh(func() { _ = gui.refreshFilesAndSubmodules() })
		}

		if scopeMap[types.STASH] {
			refresh(func() { _ = gui.refreshStashEntries() })
		}

		if scopeMap[types.TAGS] {
			refresh(func() { _ = gui.refreshTags() })
		}

		if scopeMap[types.REMOTES] {
			refresh(func() { _ = gui.refreshRemotes() })
		}

		wg.Wait()

		gui.refreshStatus()

		if options.Then != nil {
			options.Then()
		}
	}

	if options.Mode == types.BLOCK_UI {
		gui.OnUIThread(func() error {
			f()
			return nil
		})
	} else {
		f()
	}

	return nil
}

// during startup, the bottleneck is fetching the reflog entries. We need these
// on startup to sort the branches by recency. So we have two phases: INITIAL, and COMPLETE.
// In the initial phase we don't get any reflog commits, but we asynchronously get them
// and refresh the branches after that
func (gui *Gui) refreshReflogCommitsConsideringStartup() {
	switch gui.State.StartupStage {
	case INITIAL:
		go utils.Safe(func() {
			_ = gui.refreshReflogCommits()
			gui.refreshBranches()
			gui.State.StartupStage = COMPLETE
		})

	case COMPLETE:
		_ = gui.refreshReflogCommits()
	}
}

// whenever we change commits, we should update branches because the upstream/downstream
// counts can change. Whenever we change branches we should probably also change commits
// e.g. in the case of switching branches.
func (gui *Gui) refreshCommits() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go utils.Safe(func() {
		gui.refreshReflogCommitsConsideringStartup()

		gui.refreshBranches()
		wg.Done()
	})

	go utils.Safe(func() {
		_ = gui.refreshCommitsWithLimit()
		ctx, ok := gui.State.Contexts.CommitFiles.GetParentContext()
		if ok && ctx.GetKey() == context.BRANCH_COMMITS_CONTEXT_KEY {
			// This makes sense when we've e.g. just amended a commit, meaning we get a new commit SHA at the same position.
			// However if we've just added a brand new commit, it pushes the list down by one and so we would end up
			// showing the contents of a different commit than the one we initially entered.
			// Ideally we would know when to refresh the commit files context and when not to,
			// or perhaps we could just pop that context off the stack whenever cycling windows.
			// For now the awkwardness remains.
			commit := gui.getSelectedLocalCommit()
			if commit != nil {
				gui.State.Panels.CommitFiles.refName = commit.RefName()
				_ = gui.refreshCommitFilesView()
			}
		}
		wg.Done()
	})

	wg.Wait()
}

func (gui *Gui) refreshCommitsWithLimit() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	commits, err := gui.git.Loaders.Commits.GetCommits(
		loaders.GetCommitsOptions{
			Limit:                gui.State.Panels.Commits.LimitCommits,
			FilterPath:           gui.State.Modes.Filtering.GetPath(),
			IncludeRebaseCommits: true,
			RefName:              gui.refForLog(),
			All:                  gui.ShowWholeGitGraph,
		},
	)
	if err != nil {
		return err
	}
	gui.State.Commits = commits

	return gui.c.PostRefreshUpdate(gui.State.Contexts.BranchCommits)
}

func (gui *Gui) refreshRebaseCommits() error {
	gui.Mutexes.BranchCommitsMutex.Lock()
	defer gui.Mutexes.BranchCommitsMutex.Unlock()

	updatedCommits, err := gui.git.Loaders.Commits.MergeRebasingCommits(gui.State.Commits)
	if err != nil {
		return err
	}
	gui.State.Commits = updatedCommits

	return gui.c.PostRefreshUpdate(gui.State.Contexts.BranchCommits)
}

func (self *Gui) refreshTags() error {
	tags, err := self.git.Loaders.Tags.GetTags()
	if err != nil {
		return self.c.Error(err)
	}

	self.State.Tags = tags

	return self.postRefreshUpdate(self.State.Contexts.Tags)
}

func (gui *Gui) refreshStateSubmoduleConfigs() error {
	configs, err := gui.git.Submodule.GetConfigs()
	if err != nil {
		return err
	}

	gui.State.Submodules = configs

	return nil
}

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a state.Branches array to pick the current branch from
func (gui *Gui) refreshBranches() {
	reflogCommits := gui.State.FilteredReflogCommits
	if gui.State.Modes.Filtering.Active() {
		// in filter mode we filter our reflog commits to just those containing the path
		// however we need all the reflog entries to populate the recencies of our branches
		// which allows us to order them correctly. So if we're filtering we'll just
		// manually load all the reflog commits here
		var err error
		reflogCommits, _, err = gui.git.Loaders.ReflogCommits.GetReflogCommits(nil, "")
		if err != nil {
			gui.c.Log.Error(err)
		}
	}

	branches, err := gui.git.Loaders.Branches.Load(reflogCommits)
	if err != nil {
		_ = gui.c.Error(err)
	}

	gui.State.Branches = branches

	if err := gui.c.PostRefreshUpdate(gui.State.Contexts.Branches); err != nil {
		gui.c.Log.Error(err)
	}

	gui.refreshStatus()
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
		if err := gui.c.PostRefreshUpdate(gui.State.Contexts.Submodules); err != nil {
			gui.c.Log.Error(err)
		}

		if types.ContextKey(gui.Views.Files.Context) == context.FILES_CONTEXT_KEY {
			// doing this a little custom (as opposed to using gui.c.PostRefreshUpdate) because we handle selecting the file explicitly below
			if err := gui.State.Contexts.Files.HandleRender(); err != nil {
				return err
			}
		}

		if gui.currentContext().GetKey() == context.FILES_CONTEXT_KEY {
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

func (gui *Gui) refreshMergeState() error {
	gui.State.Panels.Merging.Lock()
	defer gui.State.Panels.Merging.Unlock()

	if gui.currentContext().GetKey() != context.MAIN_MERGING_CONTEXT_KEY {
		return nil
	}

	hasConflicts, err := gui.setConflictsAndRender(gui.State.Panels.Merging.GetPath(), true)
	if err != nil {
		return gui.c.Error(err)
	}

	if !hasConflicts {
		return gui.escapeMerge()
	}

	return nil
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
		gui.c.LogAction(gui.Tr.Actions.StageResolvedFiles)
		if err := gui.git.WorkingTree.StageFiles(pathsToStage); err != nil {
			return gui.c.Error(err)
		}
	}

	files := gui.git.Loaders.Files.
		GetStatusFiles(loaders.GetStatusFileOptions{})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	if gui.git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE && conflictFileCount == 0 && prevConflictFileCount > 0 {
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

// the reflogs panel is the only panel where we cache data, in that we only
// load entries that have been created since we last ran the call. This means
// we need to be more careful with how we use this, and to ensure we're emptying
// the reflogs array when changing contexts.
// This method also manages two things: ReflogCommits and FilteredReflogCommits.
// FilteredReflogCommits are rendered in the reflogs panel, and ReflogCommits
// are used by the branches panel to obtain recency values for sorting.
func (gui *Gui) refreshReflogCommits() error {
	// pulling state into its own variable incase it gets swapped out for another state
	// and we get an out of bounds exception
	state := gui.State
	var lastReflogCommit *models.Commit
	if len(state.ReflogCommits) > 0 {
		lastReflogCommit = state.ReflogCommits[0]
	}

	refresh := func(stateCommits *[]*models.Commit, filterPath string) error {
		commits, onlyObtainedNewReflogCommits, err := gui.git.Loaders.ReflogCommits.
			GetReflogCommits(lastReflogCommit, filterPath)
		if err != nil {
			return gui.c.Error(err)
		}

		if onlyObtainedNewReflogCommits {
			*stateCommits = append(commits, *stateCommits...)
		} else {
			*stateCommits = commits
		}
		return nil
	}

	if err := refresh(&state.ReflogCommits, ""); err != nil {
		return err
	}

	if gui.State.Modes.Filtering.Active() {
		if err := refresh(&state.FilteredReflogCommits, state.Modes.Filtering.GetPath()); err != nil {
			return err
		}
	} else {
		state.FilteredReflogCommits = state.ReflogCommits
	}

	return gui.c.PostRefreshUpdate(gui.State.Contexts.ReflogCommits)
}

func (gui *Gui) refreshRemotes() error {
	prevSelectedRemote := gui.getSelectedRemote()

	remotes, err := gui.git.Loaders.Remotes.GetRemotes()
	if err != nil {
		return gui.c.Error(err)
	}

	gui.State.Remotes = remotes

	// we need to ensure our selected remote branches aren't now outdated
	if prevSelectedRemote != nil && gui.State.RemoteBranches != nil {
		// find remote now
		for _, remote := range remotes {
			if remote.Name == prevSelectedRemote.Name {
				gui.State.RemoteBranches = remote.Branches
			}
		}
	}

	return gui.c.PostRefreshUpdate(gui.mustContextForContextKey(types.ContextKey(gui.Views.Branches.Context)))
}

func (gui *Gui) refreshStashEntries() error {
	gui.State.StashEntries = gui.git.Loaders.Stash.
		GetStashEntries(gui.State.Modes.Filtering.GetPath())

	return gui.postRefreshUpdate(gui.State.Contexts.Stash)
}

// never call this on its own, it should only be called from within refreshCommits()
func (gui *Gui) refreshStatus() {
	gui.Mutexes.RefreshingStatusMutex.Lock()
	defer gui.Mutexes.RefreshingStatusMutex.Unlock()

	currentBranch := gui.getCheckedOutBranch()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return
	}
	status := ""

	if currentBranch.IsRealBranch() {
		status += presentation.ColoredBranchStatus(currentBranch) + " "
	}

	workingTreeState := gui.git.Status.WorkingTreeState()
	if workingTreeState != enums.REBASE_MODE_NONE {
		status += style.FgYellow.Sprintf("(%s) ", formatWorkingTreeState(workingTreeState))
	}

	name := presentation.GetBranchTextStyle(currentBranch.Name).Sprint(currentBranch.Name)
	repoName := utils.GetCurrentRepoName()
	status += fmt.Sprintf("%s → %s ", repoName, name)

	gui.setViewContent(gui.Views.Status, status)
}
