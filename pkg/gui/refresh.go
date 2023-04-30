package gui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func getScopeNames(scopes []types.RefreshableView) []string {
	scopeNameMap := map[types.RefreshableView]string{
		types.COMMITS:         "commits",
		types.BRANCHES:        "branches",
		types.FILES:           "files",
		types.SUBMODULES:      "submodules",
		types.STASH:           "stash",
		types.REFLOG:          "reflog",
		types.TAGS:            "tags",
		types.REMOTES:         "remotes",
		types.STATUS:          "status",
		types.BISECT_INFO:     "bisect",
		types.STAGING:         "staging",
		types.MERGE_CONFLICTS: "mergeConflicts",
	}

	return slices.Map(scopes, func(scope types.RefreshableView) string {
		return scopeNameMap[scope]
	})
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
		var scopeSet *set.Set[types.RefreshableView]
		if len(options.Scope) == 0 {
			// not refreshing staging/patch-building unless explicitly requested because we only need
			// to refresh those while focused.
			scopeSet = set.NewFromSlice([]types.RefreshableView{
				types.COMMITS,
				types.BRANCHES,
				types.FILES,
				types.STASH,
				types.REFLOG,
				types.TAGS,
				types.REMOTES,
				types.STATUS,
				types.BISECT_INFO,
				types.STAGING,
			})
		} else {
			scopeSet = set.NewFromSlice(options.Scope)
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

		if scopeSet.Includes(types.COMMITS) || scopeSet.Includes(types.BRANCHES) || scopeSet.Includes(types.REFLOG) || scopeSet.Includes(types.BISECT_INFO) {
			refresh(gui.refreshCommits)
		} else if scopeSet.Includes(types.REBASE_COMMITS) {
			// the above block handles rebase commits so we only need to call this one
			// if we've asked specifically for rebase commits and not those other things
			refresh(func() { _ = gui.refreshRebaseCommits() })
		}

		// reason we're not doing this if the COMMITS type is included is that if the COMMITS type _is_ included we will refresh the commit files context anyway
		if scopeSet.Includes(types.COMMIT_FILES) && !scopeSet.Includes(types.COMMITS) {
			refresh(func() { _ = gui.refreshCommitFilesContext() })
		}

		if scopeSet.Includes(types.FILES) || scopeSet.Includes(types.SUBMODULES) {
			refresh(func() { _ = gui.refreshFilesAndSubmodules() })
		}

		if scopeSet.Includes(types.STASH) {
			refresh(func() { _ = gui.refreshStashEntries() })
		}

		if scopeSet.Includes(types.TAGS) {
			refresh(func() { _ = gui.refreshTags() })
		}

		if scopeSet.Includes(types.REMOTES) {
			refresh(func() { _ = gui.refreshRemotes() })
		}

		if scopeSet.Includes(types.STAGING) {
			refresh(func() { _ = gui.refreshStagingPanel(types.OnFocusOpts{}) })
		}

		if scopeSet.Includes(types.PATCH_BUILDING) {
			refresh(func() { _ = gui.refreshPatchBuildingPanel(types.OnFocusOpts{}) })
		}

		if scopeSet.Includes(types.MERGE_CONFLICTS) || scopeSet.Includes(types.FILES) {
			refresh(func() { _ = gui.refreshMergeState() })
		}

		wg.Wait()

		gui.refreshStatus()

		if options.Then != nil {
			options.Then()
		}
	}

	if options.Mode == types.BLOCK_UI {
		gui.c.OnUIThread(func() error {
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
		if ok && ctx.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
			// This makes sense when we've e.g. just amended a commit, meaning we get a new commit SHA at the same position.
			// However if we've just added a brand new commit, it pushes the list down by one and so we would end up
			// showing the contents of a different commit than the one we initially entered.
			// Ideally we would know when to refresh the commit files context and when not to,
			// or perhaps we could just pop that context off the stack whenever cycling windows.
			// For now the awkwardness remains.
			commit := gui.getSelectedLocalCommit()
			if commit != nil {
				gui.State.Contexts.CommitFiles.SetRef(commit)
				gui.State.Contexts.CommitFiles.SetTitleRef(commit.RefName())
				_ = gui.refreshCommitFilesContext()
			}
		}
		wg.Done()
	})

	wg.Wait()
}

func (gui *Gui) refreshCommitsWithLimit() error {
	gui.Mutexes.LocalCommitsMutex.Lock()
	defer gui.Mutexes.LocalCommitsMutex.Unlock()

	commits, err := gui.git.Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                gui.State.Contexts.LocalCommits.GetLimitCommits(),
			FilterPath:           gui.State.Modes.Filtering.GetPath(),
			IncludeRebaseCommits: true,
			RefName:              gui.refForLog(),
			All:                  gui.State.Contexts.LocalCommits.GetShowWholeGitGraph(),
		},
	)
	if err != nil {
		return err
	}
	gui.State.Model.Commits = commits
	gui.State.Model.WorkingTreeStateAtLastCommitRefresh = gui.git.Status.WorkingTreeState()

	return gui.c.PostRefreshUpdate(gui.State.Contexts.LocalCommits)
}

func (gui *Gui) refreshCommitFilesContext() error {
	ref := gui.State.Contexts.CommitFiles.GetRef()
	to := ref.RefName()
	from, reverse := gui.State.Modes.Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

	files, err := gui.git.Loaders.CommitFileLoader.GetFilesInDiff(from, to, reverse)
	if err != nil {
		return gui.c.Error(err)
	}
	gui.State.Model.CommitFiles = files
	gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.SetTree()

	return gui.c.PostRefreshUpdate(gui.State.Contexts.CommitFiles)
}

func (gui *Gui) refreshRebaseCommits() error {
	gui.Mutexes.LocalCommitsMutex.Lock()
	defer gui.Mutexes.LocalCommitsMutex.Unlock()

	updatedCommits, err := gui.git.Loaders.CommitLoader.MergeRebasingCommits(gui.State.Model.Commits)
	if err != nil {
		return err
	}
	gui.State.Model.Commits = updatedCommits
	gui.State.Model.WorkingTreeStateAtLastCommitRefresh = gui.git.Status.WorkingTreeState()

	return gui.c.PostRefreshUpdate(gui.State.Contexts.LocalCommits)
}

func (self *Gui) refreshTags() error {
	tags, err := self.git.Loaders.TagLoader.GetTags()
	if err != nil {
		return self.c.Error(err)
	}

	self.State.Model.Tags = tags

	return self.postRefreshUpdate(self.State.Contexts.Tags)
}

func (gui *Gui) refreshStateSubmoduleConfigs() error {
	configs, err := gui.git.Submodule.GetConfigs()
	if err != nil {
		return err
	}

	gui.State.Model.Submodules = configs

	return nil
}

// gui.refreshStatus is called at the end of this because that's when we can
// be sure there is a State.Model.Branches array to pick the current branch from
func (gui *Gui) refreshBranches() {
	reflogCommits := gui.State.Model.FilteredReflogCommits
	if gui.State.Modes.Filtering.Active() {
		// in filter mode we filter our reflog commits to just those containing the path
		// however we need all the reflog entries to populate the recencies of our branches
		// which allows us to order them correctly. So if we're filtering we'll just
		// manually load all the reflog commits here
		var err error
		reflogCommits, _, err = gui.git.Loaders.ReflogCommitLoader.GetReflogCommits(nil, "")
		if err != nil {
			gui.c.Log.Error(err)
		}
	}

	branches, err := gui.git.Loaders.BranchLoader.Load(reflogCommits)
	if err != nil {
		_ = gui.c.Error(err)
	}

	gui.State.Model.Branches = branches

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

	if err := gui.refreshStateSubmoduleConfigs(); err != nil {
		return err
	}

	if err := gui.refreshStateFiles(); err != nil {
		return err
	}

	gui.c.OnUIThread(func() error {
		if err := gui.c.PostRefreshUpdate(gui.State.Contexts.Submodules); err != nil {
			gui.c.Log.Error(err)
		}

		if err := gui.c.PostRefreshUpdate(gui.State.Contexts.Files); err != nil {
			gui.c.Log.Error(err)
		}

		return nil
	})

	return nil
}

func (gui *Gui) refreshMergeState() error {
	gui.State.Contexts.MergeConflicts.GetMutex().Lock()
	defer gui.State.Contexts.MergeConflicts.GetMutex().Unlock()

	if gui.currentContext().GetKey() != context.MERGE_CONFLICTS_CONTEXT_KEY {
		return nil
	}

	hasConflicts, err := gui.helpers.MergeConflicts.SetConflictsAndRender(gui.State.Contexts.MergeConflicts.GetState().GetPath(), true)
	if err != nil {
		return gui.c.Error(err)
	}

	if !hasConflicts {
		return gui.helpers.MergeConflicts.EscapeMerge()
	}

	return nil
}

func (gui *Gui) refreshStateFiles() error {
	state := gui.State

	fileTreeViewModel := state.Contexts.Files.FileTreeViewModel

	// If git thinks any of our files have inline merge conflicts, but they actually don't,
	// we stage them.
	// Note that if files with merge conflicts have both arisen and have been resolved
	// between refreshes, we won't stage them here. This is super unlikely though,
	// and this approach spares us from having to call `git status` twice in a row.
	// Although this also means that at startup we won't be staging anything until
	// we call git status again.
	pathsToStage := []string{}
	prevConflictFileCount := 0
	for _, file := range gui.State.Model.Files {
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

	files := gui.git.Loaders.FileLoader.
		GetStatusFiles(git_commands.GetStatusFileOptions{})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	if gui.git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE && conflictFileCount == 0 && prevConflictFileCount > 0 {
		gui.c.OnUIThread(func() error { return gui.helpers.MergeAndRebase.PromptToContinueRebase() })
	}

	fileTreeViewModel.RWMutex.Lock()

	// only taking over the filter if it hasn't already been set by the user.
	// Though this does make it impossible for the user to actually say they want to display all if
	// conflicts are currently being shown. Hmm. Worth it I reckon. If we need to add some
	// extra state here to see if the user's set the filter themselves we can do that, but
	// I'd prefer to maintain as little state as possible.
	if conflictFileCount > 0 {
		if fileTreeViewModel.GetFilter() == filetree.DisplayAll {
			fileTreeViewModel.SetFilter(filetree.DisplayConflicted)
		}
	} else if fileTreeViewModel.GetFilter() == filetree.DisplayConflicted {
		fileTreeViewModel.SetFilter(filetree.DisplayAll)
	}

	state.Model.Files = files
	fileTreeViewModel.SetTree()
	fileTreeViewModel.RWMutex.Unlock()

	if err := gui.fileWatcher.addFilesToFileWatcher(files); err != nil {
		return err
	}

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
	if len(state.Model.ReflogCommits) > 0 {
		lastReflogCommit = state.Model.ReflogCommits[0]
	}

	refresh := func(stateCommits *[]*models.Commit, filterPath string) error {
		commits, onlyObtainedNewReflogCommits, err := gui.git.Loaders.ReflogCommitLoader.
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

	if err := refresh(&state.Model.ReflogCommits, ""); err != nil {
		return err
	}

	if gui.State.Modes.Filtering.Active() {
		if err := refresh(&state.Model.FilteredReflogCommits, state.Modes.Filtering.GetPath()); err != nil {
			return err
		}
	} else {
		state.Model.FilteredReflogCommits = state.Model.ReflogCommits
	}

	return gui.c.PostRefreshUpdate(gui.State.Contexts.ReflogCommits)
}

func (gui *Gui) refreshRemotes() error {
	prevSelectedRemote := gui.State.Contexts.Remotes.GetSelected()

	remotes, err := gui.git.Loaders.RemoteLoader.GetRemotes()
	if err != nil {
		return gui.c.Error(err)
	}

	gui.State.Model.Remotes = remotes

	// we need to ensure our selected remote branches aren't now outdated
	if prevSelectedRemote != nil && gui.State.Model.RemoteBranches != nil {
		// find remote now
		for _, remote := range remotes {
			if remote.Name == prevSelectedRemote.Name {
				gui.State.Model.RemoteBranches = remote.Branches
				break
			}
		}
	}

	if err := gui.c.PostRefreshUpdate(gui.State.Contexts.Remotes); err != nil {
		return err
	}

	if err := gui.c.PostRefreshUpdate(gui.State.Contexts.RemoteBranches); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) refreshStashEntries() error {
	gui.State.Model.StashEntries = gui.git.Loaders.StashLoader.
		GetStashEntries(gui.State.Modes.Filtering.GetPath())

	return gui.postRefreshUpdate(gui.State.Contexts.Stash)
}

// never call this on its own, it should only be called from within refreshCommits()
func (gui *Gui) refreshStatus() {
	gui.Mutexes.RefreshingStatusMutex.Lock()
	defer gui.Mutexes.RefreshingStatusMutex.Unlock()

	currentBranch := gui.helpers.Refs.GetCheckedOutRef()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return
	}
	status := ""

	if currentBranch.IsRealBranch() {
		status += presentation.ColoredBranchStatus(currentBranch, gui.Tr) + " "
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

func (gui *Gui) refreshStagingPanel(focusOpts types.OnFocusOpts) error {
	secondaryFocused := gui.secondaryStagingFocused()
	mainFocused := gui.mainStagingFocused()

	// this method could be called when the staging panel is not being used,
	// in which case we don't want to do anything.
	if !mainFocused && !secondaryFocused {
		return nil
	}

	mainSelectedLineIdx := -1
	secondarySelectedLineIdx := -1
	if focusOpts.ClickedViewLineIdx > 0 {
		if secondaryFocused {
			secondarySelectedLineIdx = focusOpts.ClickedViewLineIdx
		} else {
			mainSelectedLineIdx = focusOpts.ClickedViewLineIdx
		}
	}

	mainContext := gui.State.Contexts.Staging
	secondaryContext := gui.State.Contexts.StagingSecondary

	file := gui.getSelectedFile()
	if file == nil || (!file.HasUnstagedChanges && !file.HasStagedChanges) {
		return gui.handleStagingEscape()
	}

	mainDiff := gui.git.WorkingTree.WorktreeFileDiff(file, true, false, false)
	secondaryDiff := gui.git.WorkingTree.WorktreeFileDiff(file, true, true, false)

	// grabbing locks here and releasing before we finish the function
	// because pushing say the secondary context could mean entering this function
	// again, and we don't want to have a deadlock
	mainContext.GetMutex().Lock()
	secondaryContext.GetMutex().Lock()

	mainContext.SetState(
		patch_exploring.NewState(mainDiff, mainSelectedLineIdx, mainContext.GetState(), gui.Log),
	)

	secondaryContext.SetState(
		patch_exploring.NewState(secondaryDiff, secondarySelectedLineIdx, secondaryContext.GetState(), gui.Log),
	)

	mainState := mainContext.GetState()
	secondaryState := secondaryContext.GetState()

	mainContent := mainContext.GetContentToRender(!secondaryFocused)
	secondaryContent := secondaryContext.GetContentToRender(secondaryFocused)

	mainContext.GetMutex().Unlock()
	secondaryContext.GetMutex().Unlock()

	if mainState == nil && secondaryState == nil {
		return gui.handleStagingEscape()
	}

	if mainState == nil && !secondaryFocused {
		return gui.c.PushContext(secondaryContext, focusOpts)
	}

	if secondaryState == nil && secondaryFocused {
		return gui.c.PushContext(mainContext, focusOpts)
	}

	if secondaryFocused {
		gui.State.Contexts.StagingSecondary.FocusSelection()
	} else {
		gui.State.Contexts.Staging.FocusSelection()
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Staging,
		Main: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(mainContent),
			Title: gui.Tr.UnstagedChanges,
		},
		Secondary: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(secondaryContent),
			Title: gui.Tr.StagedChanges,
		},
	})
}

func (gui *Gui) handleStagingEscape() error {
	return gui.c.PushContext(gui.State.Contexts.Files)
}

func (gui *Gui) secondaryStagingFocused() bool {
	return gui.currentStaticContext().GetKey() == gui.State.Contexts.StagingSecondary.GetKey()
}

func (gui *Gui) mainStagingFocused() bool {
	return gui.currentStaticContext().GetKey() == gui.State.Contexts.Staging.GetKey()
}

func (gui *Gui) refreshPatchBuildingPanel(opts types.OnFocusOpts) error {
	selectedLineIdx := -1
	if opts.ClickedWindowName == "main" {
		selectedLineIdx = opts.ClickedViewLineIdx
	}

	if !gui.git.Patch.PatchBuilder.Active() {
		return gui.helpers.PatchBuilding.Escape()
	}

	// get diff from commit file that's currently selected
	path := gui.State.Contexts.CommitFiles.GetSelectedPath()
	if path == "" {
		return nil
	}

	ref := gui.State.Contexts.CommitFiles.CommitFileTreeViewModel.GetRef()
	to := ref.RefName()
	from, reverse := gui.State.Modes.Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())
	diff, err := gui.git.WorkingTree.ShowFileDiff(from, to, reverse, path, true,
		gui.IgnoreWhitespaceInDiffView)
	if err != nil {
		return err
	}

	secondaryDiff := gui.git.Patch.PatchBuilder.RenderPatchForFile(path, false, false)
	if err != nil {
		return err
	}

	context := gui.State.Contexts.CustomPatchBuilder

	oldState := context.GetState()

	state := patch_exploring.NewState(diff, selectedLineIdx, oldState, gui.Log)
	context.SetState(state)
	if state == nil {
		return gui.helpers.PatchBuilding.Escape()
	}

	gui.State.Contexts.CustomPatchBuilder.FocusSelection()

	mainContent := context.GetContentToRender(true)

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().PatchBuilding,
		Main: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(mainContent),
			Title: gui.Tr.Patch,
		},
		Secondary: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(secondaryDiff),
			Title: gui.Tr.CustomPatch,
		},
	})
}

func (gui *Gui) refreshMergePanel(isFocused bool) error {
	content := gui.State.Contexts.MergeConflicts.GetContentToRender(isFocused)

	var task types.UpdateTask
	if gui.State.Contexts.MergeConflicts.IsUserScrolling() {
		task = types.NewRenderStringWithoutScrollTask(content)
	} else {
		originY := gui.State.Contexts.MergeConflicts.GetOriginY()
		task = types.NewRenderStringWithScrollTask(content, 0, originY)
	}

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().MergeConflicts,
		Main: &types.ViewUpdateOpts{
			Task: task,
		},
	})
}

func (gui *Gui) refreshSubCommitsWithLimit() error {
	gui.Mutexes.SubCommitsMutex.Lock()
	defer gui.Mutexes.SubCommitsMutex.Unlock()

	context := gui.State.Contexts.SubCommits

	commits, err := gui.git.Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                context.GetLimitCommits(),
			FilterPath:           gui.State.Modes.Filtering.GetPath(),
			IncludeRebaseCommits: false,
			RefName:              context.GetRef().FullRefName(),
		},
	)
	if err != nil {
		return err
	}
	gui.State.Model.SubCommits = commits

	return gui.c.PostRefreshUpdate(gui.State.Contexts.SubCommits)
}
