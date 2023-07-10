package helpers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
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

type RefreshHelper struct {
	c                    *HelperCommon
	refsHelper           *RefsHelper
	mergeAndRebaseHelper *MergeAndRebaseHelper
	patchBuildingHelper  *PatchBuildingHelper
	stagingHelper        *StagingHelper
	mergeConflictsHelper *MergeConflictsHelper
	fileWatcher          types.IFileWatcher
}

func NewRefreshHelper(
	c *HelperCommon,
	refsHelper *RefsHelper,
	mergeAndRebaseHelper *MergeAndRebaseHelper,
	patchBuildingHelper *PatchBuildingHelper,
	stagingHelper *StagingHelper,
	mergeConflictsHelper *MergeConflictsHelper,
	fileWatcher types.IFileWatcher,
) *RefreshHelper {
	return &RefreshHelper{
		c:                    c,
		refsHelper:           refsHelper,
		mergeAndRebaseHelper: mergeAndRebaseHelper,
		patchBuildingHelper:  patchBuildingHelper,
		stagingHelper:        stagingHelper,
		mergeConflictsHelper: mergeConflictsHelper,
		fileWatcher:          fileWatcher,
	}
}

func (self *RefreshHelper) Refresh(options types.RefreshOptions) error {
	if options.Scope == nil {
		self.c.Log.Infof(
			"refreshing all scopes in %s mode",
			getModeName(options.Mode),
		)
	} else {
		self.c.Log.Infof(
			"refreshing the following scopes in %s mode: %s",
			getModeName(options.Mode),
			strings.Join(getScopeNames(options.Scope), ","),
		)
	}

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
			if options.Mode == types.ASYNC {
				self.c.OnWorker(func(t gocui.Task) {
					f()
				})
			} else {
				f()
			}
		}

		if scopeSet.Includes(types.COMMITS) || scopeSet.Includes(types.BRANCHES) || scopeSet.Includes(types.REFLOG) || scopeSet.Includes(types.BISECT_INFO) {
			refresh(self.refreshCommits)
		} else if scopeSet.Includes(types.REBASE_COMMITS) {
			// the above block handles rebase commits so we only need to call this one
			// if we've asked specifically for rebase commits and not those other things
			refresh(func() { _ = self.refreshRebaseCommits() })
		}

		if scopeSet.Includes(types.SUB_COMMITS) {
			refresh(func() { _ = self.refreshSubCommitsWithLimit() })
		}

		// reason we're not doing this if the COMMITS type is included is that if the COMMITS type _is_ included we will refresh the commit files context anyway
		if scopeSet.Includes(types.COMMIT_FILES) && !scopeSet.Includes(types.COMMITS) {
			refresh(func() { _ = self.refreshCommitFilesContext() })
		}

		if scopeSet.Includes(types.FILES) || scopeSet.Includes(types.SUBMODULES) {
			refresh(func() { _ = self.refreshFilesAndSubmodules() })
		}

		if scopeSet.Includes(types.STASH) {
			refresh(func() { _ = self.refreshStashEntries() })
		}

		if scopeSet.Includes(types.TAGS) {
			refresh(func() { _ = self.refreshTags() })
		}

		if scopeSet.Includes(types.REMOTES) {
			refresh(func() { _ = self.refreshRemotes() })
		}

		if scopeSet.Includes(types.STAGING) {
			refresh(func() { _ = self.stagingHelper.RefreshStagingPanel(types.OnFocusOpts{}) })
		}

		if scopeSet.Includes(types.PATCH_BUILDING) {
			refresh(func() { _ = self.patchBuildingHelper.RefreshPatchBuildingPanel(types.OnFocusOpts{}) })
		}

		if scopeSet.Includes(types.MERGE_CONFLICTS) || scopeSet.Includes(types.FILES) {
			refresh(func() { _ = self.mergeConflictsHelper.RefreshMergeState() })
		}

		self.refreshStatus()

		if options.Then != nil {
			options.Then()
		}
	}

	if options.Mode == types.BLOCK_UI {
		self.c.OnUIThread(func() error {
			f()
			return nil
		})
	} else {
		f()
	}

	return nil
}

func getScopeNames(scopes []types.RefreshableView) []string {
	scopeNameMap := map[types.RefreshableView]string{
		types.COMMITS:         "commits",
		types.BRANCHES:        "branches",
		types.FILES:           "files",
		types.SUBMODULES:      "submodules",
		types.SUB_COMMITS:     "subCommits",
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

// during startup, the bottleneck is fetching the reflog entries. We need these
// on startup to sort the branches by recency. So we have two phases: INITIAL, and COMPLETE.
// In the initial phase we don't get any reflog commits, but we asynchronously get them
// and refresh the branches after that
func (self *RefreshHelper) refreshReflogCommitsConsideringStartup() {
	switch self.c.State().GetRepoState().GetStartupStage() {
	case types.INITIAL:
		self.c.OnWorker(func(_ gocui.Task) {
			_ = self.refreshReflogCommits()
			self.refreshBranches()
			self.c.State().GetRepoState().SetStartupStage(types.COMPLETE)
		})

	case types.COMPLETE:
		_ = self.refreshReflogCommits()
	}
}

// whenever we change commits, we should update branches because the upstream/downstream
// counts can change. Whenever we change branches we should probably also change commits
// e.g. in the case of switching branches.
func (self *RefreshHelper) refreshCommits() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go utils.Safe(func() {
		self.refreshReflogCommitsConsideringStartup()

		self.refreshBranches()
		wg.Done()
	})

	go utils.Safe(func() {
		_ = self.refreshCommitsWithLimit()
		ctx, ok := self.c.Contexts().CommitFiles.GetParentContext()
		if ok && ctx.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
			// This makes sense when we've e.g. just amended a commit, meaning we get a new commit SHA at the same position.
			// However if we've just added a brand new commit, it pushes the list down by one and so we would end up
			// showing the contents of a different commit than the one we initially entered.
			// Ideally we would know when to refresh the commit files context and when not to,
			// or perhaps we could just pop that context off the stack whenever cycling windows.
			// For now the awkwardness remains.
			commit := self.c.Contexts().LocalCommits.GetSelected()
			if commit != nil {
				self.c.Contexts().CommitFiles.SetRef(commit)
				self.c.Contexts().CommitFiles.SetTitleRef(commit.RefName())
				_ = self.refreshCommitFilesContext()
			}
		}
		wg.Done()
	})

	wg.Wait()
}

func (self *RefreshHelper) refreshCommitsWithLimit() error {
	self.c.Mutexes().LocalCommitsMutex.Lock()
	defer self.c.Mutexes().LocalCommitsMutex.Unlock()

	commits, err := self.c.Git().Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                self.c.Contexts().LocalCommits.GetLimitCommits(),
			FilterPath:           self.c.Modes().Filtering.GetPath(),
			IncludeRebaseCommits: true,
			RefName:              self.refForLog(),
			All:                  self.c.Contexts().LocalCommits.GetShowWholeGitGraph(),
		},
	)
	if err != nil {
		return err
	}
	self.c.Model().Commits = commits
	self.c.Model().WorkingTreeStateAtLastCommitRefresh = self.c.Git().Status.WorkingTreeState()

	return self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
}

func (self *RefreshHelper) refreshSubCommitsWithLimit() error {
	self.c.Mutexes().SubCommitsMutex.Lock()
	defer self.c.Mutexes().SubCommitsMutex.Unlock()

	commits, err := self.c.Git().Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                self.c.Contexts().SubCommits.GetLimitCommits(),
			FilterPath:           self.c.Modes().Filtering.GetPath(),
			IncludeRebaseCommits: false,
			RefName:              self.c.Contexts().SubCommits.GetRef().FullRefName(),
		},
	)
	if err != nil {
		return err
	}
	self.c.Model().SubCommits = commits

	return self.c.PostRefreshUpdate(self.c.Contexts().SubCommits)
}

func (self *RefreshHelper) refreshCommitFilesContext() error {
	ref := self.c.Contexts().CommitFiles.GetRef()
	to := ref.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

	files, err := self.c.Git().Loaders.CommitFileLoader.GetFilesInDiff(from, to, reverse)
	if err != nil {
		return self.c.Error(err)
	}
	self.c.Model().CommitFiles = files
	self.c.Contexts().CommitFiles.CommitFileTreeViewModel.SetTree()

	return self.c.PostRefreshUpdate(self.c.Contexts().CommitFiles)
}

func (self *RefreshHelper) refreshRebaseCommits() error {
	self.c.Mutexes().LocalCommitsMutex.Lock()
	defer self.c.Mutexes().LocalCommitsMutex.Unlock()

	updatedCommits, err := self.c.Git().Loaders.CommitLoader.MergeRebasingCommits(self.c.Model().Commits)
	if err != nil {
		return err
	}
	self.c.Model().Commits = updatedCommits
	self.c.Model().WorkingTreeStateAtLastCommitRefresh = self.c.Git().Status.WorkingTreeState()

	return self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
}

func (self *RefreshHelper) refreshTags() error {
	tags, err := self.c.Git().Loaders.TagLoader.GetTags()
	if err != nil {
		return self.c.Error(err)
	}

	self.c.Model().Tags = tags

	return self.c.PostRefreshUpdate(self.c.Contexts().Tags)
}

func (self *RefreshHelper) refreshStateSubmoduleConfigs() error {
	configs, err := self.c.Git().Submodule.GetConfigs()
	if err != nil {
		return err
	}

	self.c.Model().Submodules = configs

	return nil
}

// self.refreshStatus is called at the end of this because that's when we can
// be sure there is a State.Model.Branches array to pick the current branch from
func (self *RefreshHelper) refreshBranches() {
	self.c.Mutexes().RefreshingBranchesMutex.Lock()
	defer self.c.Mutexes().RefreshingBranchesMutex.Unlock()

	reflogCommits := self.c.Model().FilteredReflogCommits
	if self.c.Modes().Filtering.Active() {
		// in filter mode we filter our reflog commits to just those containing the path
		// however we need all the reflog entries to populate the recencies of our branches
		// which allows us to order them correctly. So if we're filtering we'll just
		// manually load all the reflog commits here
		var err error
		reflogCommits, _, err = self.c.Git().Loaders.ReflogCommitLoader.GetReflogCommits(nil, "")
		if err != nil {
			self.c.Log.Error(err)
		}
	}

	branches, err := self.c.Git().Loaders.BranchLoader.Load(reflogCommits)
	if err != nil {
		_ = self.c.Error(err)
	}

	self.c.Model().Branches = branches

	if err := self.c.PostRefreshUpdate(self.c.Contexts().Branches); err != nil {
		self.c.Log.Error(err)
	}

	self.refreshStatus()
}

func (self *RefreshHelper) refreshFilesAndSubmodules() error {
	self.c.Mutexes().RefreshingFilesMutex.Lock()
	self.c.State().SetIsRefreshingFiles(true)
	defer func() {
		self.c.State().SetIsRefreshingFiles(false)
		self.c.Mutexes().RefreshingFilesMutex.Unlock()
	}()

	if err := self.refreshStateSubmoduleConfigs(); err != nil {
		return err
	}

	if err := self.refreshStateFiles(); err != nil {
		return err
	}

	self.c.OnUIThread(func() error {
		if err := self.c.PostRefreshUpdate(self.c.Contexts().Submodules); err != nil {
			self.c.Log.Error(err)
		}

		if err := self.c.PostRefreshUpdate(self.c.Contexts().Files); err != nil {
			self.c.Log.Error(err)
		}

		return nil
	})

	return nil
}

func (self *RefreshHelper) refreshStateFiles() error {
	fileTreeViewModel := self.c.Contexts().Files.FileTreeViewModel

	// If git thinks any of our files have inline merge conflicts, but they actually don't,
	// we stage them.
	// Note that if files with merge conflicts have both arisen and have been resolved
	// between refreshes, we won't stage them here. This is super unlikely though,
	// and this approach spares us from having to call `git status` twice in a row.
	// Although this also means that at startup we won't be staging anything until
	// we call git status again.
	pathsToStage := []string{}
	prevConflictFileCount := 0
	for _, file := range self.c.Model().Files {
		if file.HasMergeConflicts {
			prevConflictFileCount++
		}
		if file.HasInlineMergeConflicts {
			hasConflicts, err := mergeconflicts.FileHasConflictMarkers(file.Name)
			if err != nil {
				self.c.Log.Error(err)
			} else if !hasConflicts {
				pathsToStage = append(pathsToStage, file.Name)
			}
		}
	}

	if len(pathsToStage) > 0 {
		self.c.LogAction(self.c.Tr.Actions.StageResolvedFiles)
		if err := self.c.Git().WorkingTree.StageFiles(pathsToStage); err != nil {
			return self.c.Error(err)
		}
	}

	files := self.c.Git().Loaders.FileLoader.
		GetStatusFiles(git_commands.GetStatusFileOptions{})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	if self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE && conflictFileCount == 0 && prevConflictFileCount > 0 {
		self.c.OnUIThread(func() error { return self.mergeAndRebaseHelper.PromptToContinueRebase() })
	}

	fileTreeViewModel.RWMutex.Lock()

	// only taking over the filter if it hasn't already been set by the user.
	// Though this does make it impossible for the user to actually say they want to display all if
	// conflicts are currently being shown. Hmm. Worth it I reckon. If we need to add some
	// extra state here to see if the user's set the filter themselves we can do that, but
	// I'd prefer to maintain as little state as possible.
	if conflictFileCount > 0 {
		if fileTreeViewModel.GetFilter() == filetree.DisplayAll {
			fileTreeViewModel.SetStatusFilter(filetree.DisplayConflicted)
		}
	} else if fileTreeViewModel.GetFilter() == filetree.DisplayConflicted {
		fileTreeViewModel.SetStatusFilter(filetree.DisplayAll)
	}

	self.c.Model().Files = files
	fileTreeViewModel.SetTree()
	fileTreeViewModel.RWMutex.Unlock()

	if err := self.fileWatcher.AddFilesToFileWatcher(files); err != nil {
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
func (self *RefreshHelper) refreshReflogCommits() error {
	// pulling state into its own variable incase it gets swapped out for another state
	// and we get an out of bounds exception
	model := self.c.Model()
	var lastReflogCommit *models.Commit
	if len(model.ReflogCommits) > 0 {
		lastReflogCommit = model.ReflogCommits[0]
	}

	refresh := func(stateCommits *[]*models.Commit, filterPath string) error {
		commits, onlyObtainedNewReflogCommits, err := self.c.Git().Loaders.ReflogCommitLoader.
			GetReflogCommits(lastReflogCommit, filterPath)
		if err != nil {
			return self.c.Error(err)
		}

		if onlyObtainedNewReflogCommits {
			*stateCommits = append(commits, *stateCommits...)
		} else {
			*stateCommits = commits
		}
		return nil
	}

	if err := refresh(&model.ReflogCommits, ""); err != nil {
		return err
	}

	if self.c.Modes().Filtering.Active() {
		if err := refresh(&model.FilteredReflogCommits, self.c.Modes().Filtering.GetPath()); err != nil {
			return err
		}
	} else {
		model.FilteredReflogCommits = model.ReflogCommits
	}

	return self.c.PostRefreshUpdate(self.c.Contexts().ReflogCommits)
}

func (self *RefreshHelper) refreshRemotes() error {
	prevSelectedRemote := self.c.Contexts().Remotes.GetSelected()

	remotes, err := self.c.Git().Loaders.RemoteLoader.GetRemotes()
	if err != nil {
		return self.c.Error(err)
	}

	self.c.Model().Remotes = remotes

	// we need to ensure our selected remote branches aren't now outdated
	if prevSelectedRemote != nil && self.c.Model().RemoteBranches != nil {
		// find remote now
		for _, remote := range remotes {
			if remote.Name == prevSelectedRemote.Name {
				self.c.Model().RemoteBranches = remote.Branches
				break
			}
		}
	}

	if err := self.c.PostRefreshUpdate(self.c.Contexts().Remotes); err != nil {
		return err
	}

	if err := self.c.PostRefreshUpdate(self.c.Contexts().RemoteBranches); err != nil {
		return err
	}

	return nil
}

func (self *RefreshHelper) refreshStashEntries() error {
	self.c.Model().StashEntries = self.c.Git().Loaders.StashLoader.
		GetStashEntries(self.c.Modes().Filtering.GetPath())

	return self.c.PostRefreshUpdate(self.c.Contexts().Stash)
}

// never call this on its own, it should only be called from within refreshCommits()
func (self *RefreshHelper) refreshStatus() {
	self.c.Mutexes().RefreshingStatusMutex.Lock()
	defer self.c.Mutexes().RefreshingStatusMutex.Unlock()

	currentBranch := self.refsHelper.GetCheckedOutRef()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return
	}
	status := ""

	if currentBranch.IsRealBranch() {
		status += presentation.ColoredBranchStatus(currentBranch, self.c.Tr) + " "
	}

	workingTreeState := self.c.Git().Status.WorkingTreeState()
	if workingTreeState != enums.REBASE_MODE_NONE {
		status += style.FgYellow.Sprintf("(%s) ", presentation.FormatWorkingTreeStateLower(self.c.Tr, workingTreeState))
	}

	name := presentation.GetBranchTextStyle(currentBranch.Name).Sprint(currentBranch.Name)
	repoName := utils.GetCurrentRepoName()
	status += fmt.Sprintf("%s → %s ", repoName, name)

	self.c.SetViewContent(self.c.Views().Status, status)
}

func (self *RefreshHelper) refForLog() string {
	bisectInfo := self.c.Git().Bisect.GetInfo()
	self.c.Model().BisectInfo = bisectInfo

	if !bisectInfo.Started() {
		return "HEAD"
	}

	// need to see if our bisect's current commit is reachable from our 'new' ref.
	if bisectInfo.Bisecting() && !self.c.Git().Bisect.ReachableFromStart(bisectInfo) {
		return bisectInfo.GetNewSha()
	}

	return bisectInfo.GetStartSha()
}
