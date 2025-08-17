package helpers

import (
	"strings"
	"sync"
	"time"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type RefreshHelper struct {
	c                    *HelperCommon
	refsHelper           *RefsHelper
	mergeAndRebaseHelper *MergeAndRebaseHelper
	patchBuildingHelper  *PatchBuildingHelper
	stagingHelper        *StagingHelper
	mergeConflictsHelper *MergeConflictsHelper
	worktreeHelper       *WorktreeHelper
	searchHelper         *SearchHelper
}

func NewRefreshHelper(
	c *HelperCommon,
	refsHelper *RefsHelper,
	mergeAndRebaseHelper *MergeAndRebaseHelper,
	patchBuildingHelper *PatchBuildingHelper,
	stagingHelper *StagingHelper,
	mergeConflictsHelper *MergeConflictsHelper,
	worktreeHelper *WorktreeHelper,
	searchHelper *SearchHelper,
) *RefreshHelper {
	return &RefreshHelper{
		c:                    c,
		refsHelper:           refsHelper,
		mergeAndRebaseHelper: mergeAndRebaseHelper,
		patchBuildingHelper:  patchBuildingHelper,
		stagingHelper:        stagingHelper,
		mergeConflictsHelper: mergeConflictsHelper,
		worktreeHelper:       worktreeHelper,
		searchHelper:         searchHelper,
	}
}

func (self *RefreshHelper) Refresh(options types.RefreshOptions) {
	if options.Mode == types.ASYNC && options.Then != nil {
		panic("RefreshOptions.Then doesn't work with mode ASYNC")
	}

	t := time.Now()
	defer func() {
		self.c.Log.Infof("Refresh took %s", time.Since(t))
	}()

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
				types.WORKTREES,
				types.STATUS,
				types.BISECT_INFO,
				types.STAGING,
			})
		} else {
			scopeSet = set.NewFromSlice(options.Scope)
		}

		wg := sync.WaitGroup{}
		refresh := func(name string, f func()) {
			// if we're in a demo we don't want any async refreshes because
			// everything happens fast and it's better to have everything update
			// in the one frame
			if !self.c.InDemo() && options.Mode == types.ASYNC {
				self.c.OnWorker(func(t gocui.Task) error {
					f()
					return nil
				})
			} else {
				wg.Add(1)
				go utils.Safe(func() {
					t := time.Now()
					defer wg.Done()
					f()
					self.c.Log.Infof("refreshed %s in %s", name, time.Since(t))
				})
			}
		}

		includeWorktreesWithBranches := false
		if scopeSet.Includes(types.COMMITS) || scopeSet.Includes(types.BRANCHES) || scopeSet.Includes(types.REFLOG) || scopeSet.Includes(types.BISECT_INFO) {
			// whenever we change commits, we should update branches because the upstream/downstream
			// counts can change. Whenever we change branches we should also change commits
			// e.g. in the case of switching branches.
			refresh("commits and commit files", self.refreshCommitsAndCommitFiles)

			includeWorktreesWithBranches = scopeSet.Includes(types.WORKTREES)
			if self.c.UserConfig().Git.LocalBranchSortOrder == "recency" {
				refresh("reflog and branches", func() { self.refreshReflogAndBranches(includeWorktreesWithBranches, options.KeepBranchSelectionIndex) })
			} else {
				refresh("branches", func() { self.refreshBranches(includeWorktreesWithBranches, options.KeepBranchSelectionIndex, true) })
				refresh("reflog", func() { _ = self.refreshReflogCommits() })
			}
		} else if scopeSet.Includes(types.REBASE_COMMITS) {
			// the above block handles rebase commits so we only need to call this one
			// if we've asked specifically for rebase commits and not those other things
			refresh("rebase commits", func() { _ = self.refreshRebaseCommits() })
		}

		if scopeSet.Includes(types.SUB_COMMITS) {
			refresh("sub commits", func() { _ = self.refreshSubCommitsWithLimit() })
		}

		// reason we're not doing this if the COMMITS type is included is that if the COMMITS type _is_ included we will refresh the commit files context anyway
		if scopeSet.Includes(types.COMMIT_FILES) && !scopeSet.Includes(types.COMMITS) {
			refresh("commit files", func() { _ = self.refreshCommitFilesContext() })
		}

		fileWg := sync.WaitGroup{}
		if scopeSet.Includes(types.FILES) || scopeSet.Includes(types.SUBMODULES) {
			fileWg.Add(1)
			refresh("files", func() {
				_ = self.refreshFilesAndSubmodules()
				fileWg.Done()
			})
		}

		if scopeSet.Includes(types.STASH) {
			refresh("stash", func() { self.refreshStashEntries() })
		}

		if scopeSet.Includes(types.TAGS) {
			refresh("tags", func() { _ = self.refreshTags() })
		}

		if scopeSet.Includes(types.REMOTES) {
			refresh("remotes", func() { _ = self.refreshRemotes() })
		}

		if scopeSet.Includes(types.WORKTREES) && !includeWorktreesWithBranches {
			refresh("worktrees", func() { self.refreshWorktrees() })
		}

		if scopeSet.Includes(types.STAGING) {
			refresh("staging", func() {
				fileWg.Wait()
				self.stagingHelper.RefreshStagingPanel(types.OnFocusOpts{})
			})
		}

		if scopeSet.Includes(types.PATCH_BUILDING) {
			refresh("patch building", func() { self.patchBuildingHelper.RefreshPatchBuildingPanel(types.OnFocusOpts{}) })
		}

		if scopeSet.Includes(types.MERGE_CONFLICTS) || scopeSet.Includes(types.FILES) {
			refresh("merge conflicts", func() { _ = self.mergeConflictsHelper.RefreshMergeState() })
		}

		self.refreshStatus()

		wg.Wait()

		if options.Then != nil {
			options.Then()
		}
	}

	if options.Mode == types.BLOCK_UI {
		self.c.OnUIThread(func() error {
			f()
			return nil
		})
		return
	}

	f()
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
		types.WORKTREES:       "worktrees",
		types.STATUS:          "status",
		types.BISECT_INFO:     "bisect",
		types.STAGING:         "staging",
		types.MERGE_CONFLICTS: "mergeConflicts",
	}

	return lo.Map(scopes, func(scope types.RefreshableView, _ int) string {
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
		self.c.OnWorker(func(_ gocui.Task) error {
			_ = self.refreshReflogCommits()
			self.refreshBranches(false, true, true)
			self.c.State().GetRepoState().SetStartupStage(types.COMPLETE)
			return nil
		})

	case types.COMPLETE:
		_ = self.refreshReflogCommits()
	}
}

func (self *RefreshHelper) refreshReflogAndBranches(refreshWorktrees bool, keepBranchSelectionIndex bool) {
	loadBehindCounts := self.c.State().GetRepoState().GetStartupStage() == types.COMPLETE

	self.refreshReflogCommitsConsideringStartup()

	self.refreshBranches(refreshWorktrees, keepBranchSelectionIndex, loadBehindCounts)
}

func (self *RefreshHelper) refreshCommitsAndCommitFiles() {
	_ = self.refreshCommitsWithLimit()
	ctx := self.c.Contexts().CommitFiles.GetParentContext()
	if ctx != nil && ctx.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY {
		// This makes sense when we've e.g. just amended a commit, meaning we get a new commit hash at the same position.
		// However if we've just added a brand new commit, it pushes the list down by one and so we would end up
		// showing the contents of a different commit than the one we initially entered.
		// Ideally we would know when to refresh the commit files context and when not to,
		// or perhaps we could just pop that context off the stack whenever cycling windows.
		// For now the awkwardness remains.
		commit := self.c.Contexts().LocalCommits.GetSelected()
		if commit != nil && commit.RefName() != "" {
			refRange := self.c.Contexts().LocalCommits.GetSelectedRefRangeForDiffFiles()
			self.c.Contexts().CommitFiles.ReInit(commit, refRange)
			_ = self.refreshCommitFilesContext()
		}
	}
}

func (self *RefreshHelper) determineCheckedOutRef() models.Ref {
	if rebasedBranch := self.c.Git().Status.BranchBeingRebased(); rebasedBranch != "" {
		// During a rebase we're on a detached head, so cannot determine the
		// branch name in the usual way. We need to read it from the
		// ".git/rebase-merge/head-name" file instead.
		return &models.Branch{Name: strings.TrimPrefix(rebasedBranch, "refs/heads/")}
	}

	if bisectInfo := self.c.Git().Bisect.GetInfo(); bisectInfo.Bisecting() && bisectInfo.GetStartHash() != "" {
		// Likewise, when we're bisecting we're on a detached head as well. In
		// this case we read the branch name from the ".git/BISECT_START" file.
		return &models.Branch{Name: bisectInfo.GetStartHash()}
	}

	// In all other cases, get the branch name by asking git what branch is
	// checked out. Note that if we're on a detached head (for reasons other
	// than rebasing or bisecting, i.e. it was explicitly checked out), then
	// this will return an empty string.
	if branchName, err := self.c.Git().Branch.CurrentBranchName(); err == nil && branchName != "" {
		return &models.Branch{Name: branchName}
	}

	// Should never get here unless the working copy is corrupt
	return nil
}

func (self *RefreshHelper) refreshCommitsWithLimit() error {
	self.c.Mutexes().LocalCommitsMutex.Lock()
	defer self.c.Mutexes().LocalCommitsMutex.Unlock()

	checkedOutRef := self.determineCheckedOutRef()
	commits, err := self.c.Git().Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                self.c.Contexts().LocalCommits.GetLimitCommits(),
			FilterPath:           self.c.Modes().Filtering.GetPath(),
			FilterAuthor:         self.c.Modes().Filtering.GetAuthor(),
			IncludeRebaseCommits: true,
			RefName:              self.refForLog(),
			RefForPushedStatus:   checkedOutRef,
			All:                  self.c.Contexts().LocalCommits.GetShowWholeGitGraph(),
			MainBranches:         self.c.Model().MainBranches,
			HashPool:             self.c.Model().HashPool,
		},
	)
	if err != nil {
		return err
	}
	self.c.Model().Commits = commits
	self.RefreshAuthors(commits)
	self.c.Model().WorkingTreeStateAtLastCommitRefresh = self.c.Git().Status.WorkingTreeState()
	if checkedOutRef != nil {
		self.c.Model().CheckedOutBranch = checkedOutRef.RefName()
	} else {
		self.c.Model().CheckedOutBranch = ""
	}

	self.refreshView(self.c.Contexts().LocalCommits)
	return nil
}

func (self *RefreshHelper) refreshSubCommitsWithLimit() error {
	if self.c.Contexts().SubCommits.GetRef() == nil {
		return nil
	}

	self.c.Mutexes().SubCommitsMutex.Lock()
	defer self.c.Mutexes().SubCommitsMutex.Unlock()

	commits, err := self.c.Git().Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                   self.c.Contexts().SubCommits.GetLimitCommits(),
			FilterPath:              self.c.Modes().Filtering.GetPath(),
			FilterAuthor:            self.c.Modes().Filtering.GetAuthor(),
			IncludeRebaseCommits:    false,
			RefName:                 self.c.Contexts().SubCommits.GetRef().FullRefName(),
			RefToShowDivergenceFrom: self.c.Contexts().SubCommits.GetRefToShowDivergenceFrom(),
			RefForPushedStatus:      self.c.Contexts().SubCommits.GetRef(),
			MainBranches:            self.c.Model().MainBranches,
			HashPool:                self.c.Model().HashPool,
		},
	)
	if err != nil {
		return err
	}
	self.c.Model().SubCommits = commits
	self.RefreshAuthors(commits)

	self.refreshView(self.c.Contexts().SubCommits)
	return nil
}

func (self *RefreshHelper) RefreshAuthors(commits []*models.Commit) {
	self.c.Mutexes().AuthorsMutex.Lock()
	defer self.c.Mutexes().AuthorsMutex.Unlock()

	authors := self.c.Model().Authors
	for _, commit := range commits {
		if _, ok := authors[commit.AuthorEmail]; !ok {
			authors[commit.AuthorEmail] = &models.Author{
				Email: commit.AuthorEmail,
				Name:  commit.AuthorName,
			}
		}
	}
}

func (self *RefreshHelper) refreshCommitFilesContext() error {
	from, to := self.c.Contexts().CommitFiles.GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)

	files, err := self.c.Git().Loaders.CommitFileLoader.GetFilesInDiff(from, to, reverse)
	if err != nil {
		return err
	}
	self.c.Model().CommitFiles = files
	self.c.Contexts().CommitFiles.CommitFileTreeViewModel.SetTree()

	self.refreshView(self.c.Contexts().CommitFiles)
	return nil
}

func (self *RefreshHelper) refreshRebaseCommits() error {
	self.c.Mutexes().LocalCommitsMutex.Lock()
	defer self.c.Mutexes().LocalCommitsMutex.Unlock()

	updatedCommits, err := self.c.Git().Loaders.CommitLoader.MergeRebasingCommits(self.c.Model().HashPool, self.c.Model().Commits)
	if err != nil {
		return err
	}
	self.c.Model().Commits = updatedCommits
	self.c.Model().WorkingTreeStateAtLastCommitRefresh = self.c.Git().Status.WorkingTreeState()

	self.refreshView(self.c.Contexts().LocalCommits)
	return nil
}

func (self *RefreshHelper) refreshTags() error {
	tags, err := self.c.Git().Loaders.TagLoader.GetTags()
	if err != nil {
		return err
	}

	self.c.Model().Tags = tags

	self.refreshView(self.c.Contexts().Tags)
	return nil
}

func (self *RefreshHelper) refreshStateSubmoduleConfigs() error {
	configs, err := self.c.Git().Submodule.GetConfigs(nil)
	if err != nil {
		return err
	}

	self.c.Model().Submodules = configs

	return nil
}

// self.refreshStatus is called at the end of this because that's when we can
// be sure there is a State.Model.Branches array to pick the current branch from
func (self *RefreshHelper) refreshBranches(refreshWorktrees bool, keepBranchSelectionIndex bool, loadBehindCounts bool) {
	self.c.Mutexes().RefreshingBranchesMutex.Lock()
	defer self.c.Mutexes().RefreshingBranchesMutex.Unlock()

	branches, err := self.c.Git().Loaders.BranchLoader.Load(
		self.c.Model().ReflogCommits,
		self.c.Model().MainBranches,
		self.c.Model().Branches,
		loadBehindCounts,
		func(f func() error) {
			self.c.OnWorker(func(_ gocui.Task) error {
				return f()
			})
		},
		func() {
			self.c.OnUIThread(func() error {
				self.c.Contexts().Branches.HandleRender()
				self.refreshStatus()
				return nil
			})
		})
	if err != nil {
		self.c.Log.Error(err)
	}

	prevSelectedBranch := self.c.Contexts().Branches.GetSelected()

	self.c.Model().Branches = branches

	if refreshWorktrees {
		self.loadWorktrees()
		self.refreshView(self.c.Contexts().Worktrees)
	}

	if !keepBranchSelectionIndex && prevSelectedBranch != nil {
		self.searchHelper.ReApplyFilter(self.c.Contexts().Branches)

		_, idx, found := lo.FindIndexOf(self.c.Contexts().Branches.GetItems(),
			func(b *models.Branch) bool { return b.Name == prevSelectedBranch.Name })
		if found {
			self.c.Contexts().Branches.SetSelectedLineIdx(idx)
		}
	}

	self.refreshView(self.c.Contexts().Branches)

	// Need to re-render the commits view because the visualization of local
	// branch heads might have changed
	self.c.Mutexes().LocalCommitsMutex.Lock()
	self.c.Contexts().LocalCommits.HandleRender()
	self.c.Mutexes().LocalCommitsMutex.Unlock()

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
		self.refreshView(self.c.Contexts().Submodules)
		self.refreshView(self.c.Contexts().Files)
		return nil
	})

	return nil
}

func (self *RefreshHelper) refreshStateFiles() error {
	fileTreeViewModel := self.c.Contexts().Files.FileTreeViewModel

	prevConflictFileCount := 0
	if self.c.UserConfig().Git.AutoStageResolvedConflicts {
		// If git thinks any of our files have inline merge conflicts, but they actually don't,
		// we stage them.
		// Note that if files with merge conflicts have both arisen and have been resolved
		// between refreshes, we won't stage them here. This is super unlikely though,
		// and this approach spares us from having to call `git status` twice in a row.
		// Although this also means that at startup we won't be staging anything until
		// we call git status again.
		pathsToStage := []string{}
		for _, file := range self.c.Model().Files {
			if file.HasMergeConflicts {
				prevConflictFileCount++
			}
			if file.HasInlineMergeConflicts {
				hasConflicts, err := mergeconflicts.FileHasConflictMarkers(file.Path)
				if err != nil {
					self.c.Log.Error(err)
				} else if !hasConflicts {
					pathsToStage = append(pathsToStage, file.Path)
				}
			}
		}

		if len(pathsToStage) > 0 {
			self.c.LogAction(self.c.Tr.Actions.StageResolvedFiles)
			if err := self.c.Git().WorkingTree.StageFiles(pathsToStage, nil); err != nil {
				return err
			}
		}
	}

	files := self.c.Git().Loaders.FileLoader.
		GetStatusFiles(git_commands.GetStatusFileOptions{
			ForceShowUntracked: self.c.Contexts().Files.ForceShowUntracked(),
		})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	if self.c.Git().Status.WorkingTreeState().Any() && conflictFileCount == 0 && prevConflictFileCount > 0 {
		self.c.OnUIThread(func() error { return self.mergeAndRebaseHelper.PromptToContinueRebase() })
	}

	fileTreeViewModel.RWMutex.Lock()

	// only taking over the filter if it hasn't already been set by the user.
	if conflictFileCount > 0 && prevConflictFileCount == 0 {
		if fileTreeViewModel.GetFilter() == filetree.DisplayAll {
			fileTreeViewModel.SetStatusFilter(filetree.DisplayConflicted)
			self.c.Contexts().Files.GetView().Subtitle = self.c.Tr.FilterLabelConflictingFiles
		}
	} else if conflictFileCount == 0 && fileTreeViewModel.GetFilter() == filetree.DisplayConflicted {
		fileTreeViewModel.SetStatusFilter(filetree.DisplayAll)
		self.c.Contexts().Files.GetView().Subtitle = ""
	}

	self.c.Model().Files = files
	fileTreeViewModel.SetTree()
	fileTreeViewModel.RWMutex.Unlock()

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
	// pulling state into its own variable in case it gets swapped out for another state
	// and we get an out of bounds exception
	model := self.c.Model()

	refresh := func(stateCommits *[]*models.Commit, filterPath string, filterAuthor string) error {
		var lastReflogCommit *models.Commit
		if filterPath == "" && filterAuthor == "" && len(*stateCommits) > 0 {
			lastReflogCommit = (*stateCommits)[0]
		}

		commits, onlyObtainedNewReflogCommits, err := self.c.Git().Loaders.ReflogCommitLoader.
			GetReflogCommits(self.c.Model().HashPool, lastReflogCommit, filterPath, filterAuthor)
		if err != nil {
			return err
		}

		if onlyObtainedNewReflogCommits {
			*stateCommits = append(commits, *stateCommits...)
		} else {
			*stateCommits = commits
		}
		return nil
	}

	if err := refresh(&model.ReflogCommits, "", ""); err != nil {
		return err
	}

	if self.c.Modes().Filtering.Active() {
		if err := refresh(&model.FilteredReflogCommits, self.c.Modes().Filtering.GetPath(), self.c.Modes().Filtering.GetAuthor()); err != nil {
			return err
		}
	} else {
		model.FilteredReflogCommits = model.ReflogCommits
	}

	self.refreshView(self.c.Contexts().ReflogCommits)
	return nil
}

func (self *RefreshHelper) refreshRemotes() error {
	prevSelectedRemote := self.c.Contexts().Remotes.GetSelected()

	remotes, err := self.c.Git().Loaders.RemoteLoader.GetRemotes()
	if err != nil {
		return err
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

	self.refreshView(self.c.Contexts().Remotes)
	self.refreshView(self.c.Contexts().RemoteBranches)
	return nil
}

func (self *RefreshHelper) loadWorktrees() {
	worktrees, err := self.c.Git().Loaders.Worktrees.GetWorktrees()
	if err != nil {
		self.c.Log.Error(err)
		self.c.Model().Worktrees = []*models.Worktree{}
	}

	self.c.Model().Worktrees = worktrees
}

func (self *RefreshHelper) refreshWorktrees() {
	self.loadWorktrees()

	// need to refresh branches because the branches view shows worktrees against
	// branches
	self.refreshView(self.c.Contexts().Branches)
	self.refreshView(self.c.Contexts().Worktrees)
}

func (self *RefreshHelper) refreshStashEntries() {
	self.c.Model().StashEntries = self.c.Git().Loaders.StashLoader.
		GetStashEntries(self.c.Modes().Filtering.GetPath())

	self.refreshView(self.c.Contexts().Stash)
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

	workingTreeState := self.c.Git().Status.WorkingTreeState()
	linkedWorktreeName := self.worktreeHelper.GetLinkedWorktreeName()

	repoName := self.c.Git().RepoPaths.RepoName()

	status := presentation.FormatStatus(repoName, currentBranch, types.ItemOperationNone, linkedWorktreeName, workingTreeState, self.c.Tr, self.c.UserConfig())

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
		return bisectInfo.GetNewHash()
	}

	return bisectInfo.GetStartHash()
}

func (self *RefreshHelper) refreshView(context types.Context) {
	// Re-applying the filter must be done before re-rendering the view, so that
	// the filtered list model is up to date for rendering.
	self.searchHelper.ReApplyFilter(context)

	self.c.PostRefreshUpdate(context)

	self.c.AfterLayout(func() error {
		// Re-applying the search must be done after re-rendering the view though,
		// so that the "x of y" status is shown correctly.
		//
		// Also, it must be done after layout, because otherwise FocusPoint
		// hasn't been called yet (see ListContextTrait.FocusLine), which means
		// that the scroll position might be such that the entire visible
		// content is outside the viewport. And this would cause problems in
		// searchModelCommits.
		self.searchHelper.ReApplySearch(context)
		return nil
	})
}
