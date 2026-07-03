package helpers

import (
	"strings"
	"sync"
	"time"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
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

	// Tracks repos for which the user has dismissed the "select base GitHub remote"
	// prompt, to avoid re-prompting on every subsequent refresh within the same session.
	// Keyed by repo path so that switching to a different repo while lazygit is running
	// still triggers the prompt there.
	githubBaseRemotePromptDismissed map[string]bool

	// Last observed refs+HEAD fingerprint, used by the background poller to
	// decide whether a real refresh is needed. Written at the end of every
	// refresh that re-read refs/commits, read by the poller.
	refsSnapshotMutex deadlock.Mutex
	refsSnapshot      string
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
				types.PULL_REQUESTS,
			})
		} else {
			scopeSet = set.NewFromSlice(options.Scope)
		}

		// Expand co-refreshing scopes up front so downstream conditions can be
		// simple single-scope checks. The relationships are:
		//   - whenever the reflog or bisect info changes, commits and branches
		//     can change too (e.g. switching branches updates the reflog and
		//     can move HEAD), so refresh commits + branches alongside
		//   - submodules are refreshed as part of the files refresh
		//   - merge conflicts are part of what the files refresh produces
		if scopeSet.Includes(types.REFLOG) || scopeSet.Includes(types.BISECT_INFO) {
			scopeSet.Add(types.COMMITS, types.BRANCHES)
		}
		if scopeSet.Includes(types.SUBMODULES) {
			scopeSet.Add(types.FILES)
		}
		if scopeSet.Includes(types.FILES) {
			scopeSet.Add(types.MERGE_CONFLICTS)
		}

		// Capture the refs snapshot now, before we start reading git's state
		// below, rather than after. This is important to guard against the race
		// of git's state changing externally while (or right after) we are
		// refreshing; the risk is one potential extra refresh, but capturing the
		// snapshot at the end would risk missing one, which is worse.
		self.updateRefsSnapshotIfRelevant(scopeSet)

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

		branchesAndRemotesWg := sync.WaitGroup{}
		includeWorktreesWithBranches := false
		if scopeSet.Includes(types.COMMITS) || scopeSet.Includes(types.BRANCHES) {
			// whenever we change commits, we should update branches because the upstream/downstream
			// counts can change. Whenever we change branches we should also change commits
			// e.g. in the case of switching branches.
			refresh("commits and commit files", func() {
				self.refreshCommitsAndCommitFiles(options.CommitSelection)
			})

			includeWorktreesWithBranches = scopeSet.Includes(types.WORKTREES)
			if self.c.UserConfig().Git.LocalBranchSortOrder == "recency" {
				branchesAndRemotesWg.Add(1)
				refresh("reflog and branches", func() {
					self.refreshReflogAndBranches(includeWorktreesWithBranches, options.KeepBranchSelectionIndex)
					branchesAndRemotesWg.Done()
				})
			} else {
				branchesAndRemotesWg.Add(1)
				refresh("branches", func() {
					// Not a recency sort, so branches doesn't depend on the reflog
					// being fresh; it runs concurrently with the reflog refresh
					// below and reads whatever's in the model, as it always has.
					self.refreshBranches(includeWorktreesWithBranches, options.KeepBranchSelectionIndex, true, self.c.Model().ReflogCommits)
					branchesAndRemotesWg.Done()
				})
				refresh("reflog", func() { _, _ = self.refreshReflogCommits() })
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
		if scopeSet.Includes(types.FILES) {
			fileWg.Add(1)
			refresh("files", func() {
				_ = self.refreshFilesAndSubmodules(options.Background)
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
			branchesAndRemotesWg.Add(1)
			refresh("remotes", func() {
				_ = self.refreshRemotes()
				branchesAndRemotesWg.Done()
			})
		}

		if scopeSet.Includes(types.PULL_REQUESTS) {
			refresh("pull requests", func() {
				branchesAndRemotesWg.Wait()
				self.refreshGithubPullRequests()
			})
		}

		if scopeSet.Includes(types.WORKTREES) && !includeWorktreesWithBranches {
			refresh("worktrees", func() { self.refreshWorktrees() })
		}

		if scopeSet.Includes(types.STAGING) {
			refresh("staging", func() {
				fileWg.Wait()
				// Bounce onto the UI thread so this runs after the files
				// scope's model-update bounce — RefreshStagingPanel reads
				// Model.Files (via Files.GetSelected) and would otherwise
				// see the pre-refresh model.
				self.c.OnUIThread(func() error {
					self.stagingHelper.RefreshStagingPanel(types.OnFocusOpts{})
					return nil
				})
			})
		}

		if scopeSet.Includes(types.PATCH_BUILDING) {
			refresh("patch building", func() { self.patchBuildingHelper.RefreshPatchBuildingPanel(types.OnFocusOpts{}) })
		}

		if scopeSet.Includes(types.MERGE_CONFLICTS) {
			refresh("merge conflicts", func() { _ = self.mergeConflictsHelper.RefreshMergeState() })
		}

		self.refreshStatus()

		wg.Wait()

		if options.Then != nil {
			// Queue Then via OnUIThread so it runs *after* the refresh-scope
			// functions' model-update bounces (which are already queued by
			// now), not synchronously here — at this point the workers have
			// returned but their bounces haven't been processed yet, so
			// invoking Then synchronously would run it on a model that's
			// still pre-refresh.
			self.c.OnUIThread(options.Then)
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

// SetRefsSnapshot stores the given snapshot as the last observed refs state.
// Called externally by the background poller at startup to seed the snapshot,
// and internally by Refresh at the end of a refs-touching refresh.
func (self *RefreshHelper) SetRefsSnapshot(snapshot string) {
	self.refsSnapshotMutex.Lock()
	defer self.refsSnapshotMutex.Unlock()
	self.refsSnapshot = snapshot
}

// RefsSnapshotChangedSince reports whether the given snapshot differs from
// the last observed one. Pure read; does not update internal state.
func (self *RefreshHelper) RefsSnapshotChangedSince(snapshot string) bool {
	self.refsSnapshotMutex.Lock()
	defer self.refsSnapshotMutex.Unlock()

	// An empty stored snapshot means no refresh has captured one yet, so we
	// have no baseline to compare against and report "unchanged" rather than
	// firing a spurious refresh. This can only be the unset zero value: a
	// snapshot we actually computed is never empty, because its HEAD component
	// is always non-empty (a branch ref when attached, a hash when detached —
	// even a repo with no commits yields "ref: refs/heads/main").
	if self.refsSnapshot == "" {
		return false
	}

	return snapshot != self.refsSnapshot
}

// updateRefsSnapshotIfRelevant captures a fresh refs snapshot from disk at the
// start of a refresh that re-reads refs/commits (see the call site for why we
// capture before reading the model rather than after). This keeps the
// background poller's stored snapshot in sync with what's been observed by the
// UI, so in-app commands and focus-in refreshes don't cause the next poll to
// spuriously re-trigger.
//
// We check just COMMITS and BRANCHES because the scope-expansion step at the
// top of Refresh has already added these whenever REFLOG or BISECT_INFO are
// in scope, and whenever a nil scope was passed.
func (self *RefreshHelper) updateRefsSnapshotIfRelevant(scopeSet *set.Set[types.RefreshableView]) {
	if !scopeSet.Includes(types.COMMITS) && !scopeSet.Includes(types.BRANCHES) {
		return
	}

	snapshot, err := self.c.Git().Status.RefsSnapshot()
	if err != nil {
		self.c.Log.Warnf("RefsSnapshot failed during refresh: %v", err)
		return
	}
	self.SetRefsSnapshot(snapshot)
}

func getScopeNames(scopes []types.RefreshableView) []string {
	scopeNameMap := map[types.RefreshableView]string{
		types.COMMITS:         "commits",
		types.REBASE_COMMITS:  "rebaseCommits",
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
		types.PATCH_BUILDING:  "patchBuilding",
		types.MERGE_CONFLICTS: "mergeConflicts",
		types.COMMIT_FILES:    "commitFiles",
		types.PULL_REQUESTS:   "pullRequests",
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
// refreshReflogCommitsConsideringStartup returns the reflog commits that the
// caller should hand to refreshBranches for recency sorting. In the COMPLETE
// (normal) case that's the freshly-loaded reflog; in the INITIAL case the
// reflog is loaded asynchronously (and drives its own branches refresh once
// ready), so we return the current model value for the immediate,
// non-recency-sorted branches refresh the caller does in the meantime.
func (self *RefreshHelper) refreshReflogCommitsConsideringStartup() []*models.Commit {
	switch self.c.State().GetRepoState().GetStartupStage() {
	case types.INITIAL:
		self.c.OnWorker(func(_ gocui.Task) error {
			reflogCommits, _ := self.refreshReflogCommits()
			self.refreshBranches(false, true, true, reflogCommits)
			self.c.State().GetRepoState().SetStartupStage(types.COMPLETE)
			return nil
		})

		return self.c.Model().ReflogCommits

	case types.COMPLETE:
		reflogCommits, _ := self.refreshReflogCommits()
		return reflogCommits
	}

	return self.c.Model().ReflogCommits
}

func (self *RefreshHelper) refreshReflogAndBranches(refreshWorktrees bool, keepBranchSelectionIndex bool) {
	loadBehindCounts := self.c.State().GetRepoState().GetStartupStage() == types.COMPLETE

	reflogCommits := self.refreshReflogCommitsConsideringStartup()

	self.refreshBranches(refreshWorktrees, keepBranchSelectionIndex, loadBehindCounts, reflogCommits)
}

func (self *RefreshHelper) refreshCommitsAndCommitFiles(commitSelection types.CommitSelectionBehavior) {
	_ = self.refreshCommitsWithLimit(commitSelection)
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

func (self *RefreshHelper) refreshCommitsWithLimit(commitSelection types.CommitSelectionBehavior) error {
	self.c.Mutexes().LocalCommitsMutex.Lock()
	defer self.c.Mutexes().LocalCommitsMutex.Unlock()

	var selectionRange *localCommitSelectionRange
	if commitSelection == types.KeepCommitSelectionByHash {
		selectedIdx, rangeStartIdx, rangeSelectMode := self.c.Contexts().LocalCommits.GetSelectionRangeAndMode()
		selectionRange = captureLocalCommitSelectionRange(self.c.Model().Commits, selectedIdx, rangeStartIdx, rangeSelectMode)
	}

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

	scrollSelectionIntoView := false
	switch commitSelection {
	case types.SelectHeadCommit:
		if headCommitIdx := models.HeadCommitIdx(commits); headCommitIdx >= 0 {
			self.c.Contexts().LocalCommits.SetSelection(headCommitIdx)
			scrollSelectionIntoView = true
		}
	case types.KeepCommitSelectionByHash:
		if selectionRange != nil {
			selectedIdx, rangeStartIdx, didMove, found := findLocalCommitSelectionRange(commits, selectionRange)
			if found {
				self.c.Contexts().LocalCommits.SetSelectionRangeAndMode(selectedIdx, rangeStartIdx, selectionRange.mode)
				scrollSelectionIntoView = didMove
			}
		}
	case types.KeepCommitSelectionIndex:
		// The caller set the selection index deliberately; leave it untouched.
	}

	self.refreshView(self.c.Contexts().LocalCommits)
	if scrollSelectionIntoView {
		self.c.OnUIThread(func() error {
			self.c.Contexts().LocalCommits.FocusLine(true)
			return nil
		})
	}
	return nil
}

type localCommitSelectionRange struct {
	selectedHash     string
	selectedIsTODO   bool
	rangeStartHash   string
	rangeStartIsTODO bool
	selectedIdx      int
	rangeStartIdx    int
	mode             traits.RangeSelectMode
}

func captureLocalCommitSelectionRange(
	commits []*models.Commit,
	selectedIdx int,
	rangeStartIdx int,
	mode traits.RangeSelectMode,
) *localCommitSelectionRange {
	if !hasRestorableCommitHash(commits, selectedIdx) || !hasRestorableCommitHash(commits, rangeStartIdx) {
		return nil
	}

	return &localCommitSelectionRange{
		selectedHash:     commits[selectedIdx].Hash(),
		selectedIsTODO:   commits[selectedIdx].IsTODO(),
		rangeStartHash:   commits[rangeStartIdx].Hash(),
		rangeStartIsTODO: commits[rangeStartIdx].IsTODO(),
		selectedIdx:      selectedIdx,
		rangeStartIdx:    rangeStartIdx,
		mode:             mode,
	}
}

func findLocalCommitSelectionRange(
	commits []*models.Commit,
	selectionRange *localCommitSelectionRange,
) (int, int, bool, bool) {
	selectedIdx, foundSelected := findCommitByHashPreferringTODOStatus(
		commits, selectionRange.selectedHash, selectionRange.selectedIsTODO)
	rangeStartIdx, foundRangeStart := findCommitByHashPreferringTODOStatus(
		commits, selectionRange.rangeStartHash, selectionRange.rangeStartIsTODO)
	if !foundSelected || !foundRangeStart {
		return 0, 0, false, false
	}

	didMove := selectedIdx != selectionRange.selectedIdx || rangeStartIdx != selectionRange.rangeStartIdx
	return selectedIdx, rangeStartIdx, didMove, true
}

// findCommitByHashPreferringTODOStatus finds the commit with the given hash.
// When both a TODO and a non-TODO commit share that hash - which happens while
// reverting or cherry-picking, where the rebase TODO entry has the same hash as
// the real commit - it returns the one whose TODO status matches isTODO. When
// only one commit has the hash, it is returned regardless of its TODO status,
// so that a selected commit which turned into a TODO entry across the refresh is
// still found (e.g. when starting an interactive rebase that stops to edit it).
func findCommitByHashPreferringTODOStatus(commits []*models.Commit, hash string, isTODO bool) (int, bool) {
	fallbackIdx := -1
	for idx, commit := range commits {
		if commit.Hash() != hash {
			continue
		}
		if commit.IsTODO() == isTODO {
			return idx, true
		}
		if fallbackIdx == -1 {
			fallbackIdx = idx
		}
	}

	return fallbackIdx, fallbackIdx != -1
}

func hasRestorableCommitHash(commits []*models.Commit, idx int) bool {
	return idx >= 0 && idx < len(commits) && commits[idx].Hash() != ""
}

func (self *RefreshHelper) refreshSubCommitsWithLimit() error {
	if self.c.Contexts().SubCommits.GetRef() == nil {
		return nil
	}

	self.c.Mutexes().SubCommitsMutex.Lock()
	defer self.c.Mutexes().SubCommitsMutex.Unlock()

	generation := self.c.State().GetRepoGeneration()

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
	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().SubCommits = commits
		self.RefreshAuthors(commits)
		return nil
	})

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
	generation := self.c.State().GetRepoGeneration()

	files, err := self.c.Git().Loaders.CommitFileLoader.GetFilesInDiff(from, to, reverse)
	if err != nil {
		return err
	}
	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().CommitFiles = files
		self.c.Contexts().CommitFiles.CommitFileTreeViewModel.SetTree()
		return nil
	})
	self.refreshView(self.c.Contexts().CommitFiles)
	return nil
}

func (self *RefreshHelper) refreshRebaseCommits() error {
	self.c.Mutexes().LocalCommitsMutex.Lock()
	defer self.c.Mutexes().LocalCommitsMutex.Unlock()

	generation := self.c.State().GetRepoGeneration()

	updatedCommits, err := self.c.Git().Loaders.CommitLoader.MergeRebasingCommits(self.c.Model().HashPool, self.c.Model().Commits)
	if err != nil {
		return err
	}
	workingTreeState := self.c.Git().Status.WorkingTreeState()

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().Commits = updatedCommits
		self.c.Model().WorkingTreeStateAtLastCommitRefresh = workingTreeState
		return nil
	})

	self.refreshView(self.c.Contexts().LocalCommits)
	return nil
}

func (self *RefreshHelper) refreshTags() error {
	generation := self.c.State().GetRepoGeneration()

	tags, err := self.c.Git().Loaders.TagLoader.GetTags()
	if err != nil {
		return err
	}

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().Tags = tags
		return nil
	})

	self.refreshView(self.c.Contexts().Tags)
	return nil
}

func (self *RefreshHelper) refreshStateSubmoduleConfigs() ([]*models.SubmoduleConfig, error) {
	return self.c.Git().Submodule.GetConfigs(nil)
}

// self.refreshStatus is called at the end of this because that's when we can
// be sure there is a State.Model.Branches array to pick the current branch from
func (self *RefreshHelper) refreshBranches(refreshWorktrees bool, keepBranchSelectionIndex bool, loadBehindCounts bool, reflogCommits []*models.Commit) {
	self.c.Mutexes().RefreshingBranchesMutex.Lock()
	defer self.c.Mutexes().RefreshingBranchesMutex.Unlock()

	branches, err := self.c.Git().Loaders.BranchLoader.Load(
		reflogCommits,
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
	self.rebuildPullRequestsMap()

	if refreshWorktrees {
		// TODO: this synchronous worker write goes away when refreshBranches is
		// itself migrated to bouncing; for now it matches the rest of this
		// not-yet-bounced function.
		self.c.Model().Worktrees = self.loadWorktrees()
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
	self.c.OnUIThread(func() error {
		self.c.Mutexes().LocalCommitsMutex.Lock()
		self.c.Contexts().LocalCommits.HandleRender()
		self.c.Mutexes().LocalCommitsMutex.Unlock()
		return nil
	})

	self.refreshStatus()
}

func (self *RefreshHelper) refreshFilesAndSubmodules(background bool) error {
	configs, err := self.refreshStateSubmoduleConfigs()
	if err != nil {
		return err
	}

	if err := self.refreshStateFiles(background, configs); err != nil {
		return err
	}

	self.refreshView(self.c.Contexts().Submodules)
	self.refreshView(self.c.Contexts().Files)

	return nil
}

// onUIThreadUnlessRepoChanged bounces a refresh's model/view update onto the UI
// thread, but drops it if the repo was switched while the refresh was in flight.
// Refresh workers do their git work off the UI thread and enqueue their model
// writes here; a repo switch (which replaces the whole model and context tree)
// bumps the generation, so a write captured under the old generation must not
// clobber the new repo's state. Callers capture the generation with
// State().GetRepoGeneration() before doing their git work and pass it in.
func (self *RefreshHelper) onUIThreadUnlessRepoChanged(generation int, f func() error) {
	self.c.OnUIThread(func() error {
		if self.c.State().GetRepoGeneration() != generation {
			return nil
		}
		return f()
	})
}

func (self *RefreshHelper) refreshStateFiles(background bool, submoduleConfigs []*models.SubmoduleConfig) error {
	fileTreeViewModel := self.c.Contexts().Files.FileTreeViewModel
	generation := self.c.State().GetRepoGeneration()

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
			Background:         background,
		})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	repoState := self.c.State().GetRepoState()
	workingTreeState := self.c.Git().Status.WorkingTreeState()
	if workingTreeState.None() {
		// No operation is in progress (any more), so forget that we started one.
		// This also covers an operation that was finished or aborted externally.
		repoState.SetMergeOrRebaseStartedInLazygit(false)
	}

	if workingTreeState.Any() && conflictFileCount == 0 {
		if prevConflictFileCount > 0 && repoState.GetMergeOrRebaseStartedInLazygit() {
			// The conflicts of an operation we started have just been resolved
			// (e.g. in the user's editor). Offer to continue it. We only do this
			// for operations we started ourselves; prompting for one that was
			// started outside lazygit (e.g. by a coding agent) would be confusing.
			self.onUIThreadUnlessRepoChanged(generation, func() error {
				return self.mergeAndRebaseHelper.PromptToContinueRebase()
			})
		}
	} else {
		// Either there's no operation in progress any more, or new conflicts have
		// appeared. Either way, a "continue?" prompt we're showing is now stale
		// (e.g. the operation was continued or aborted outside lazygit), so
		// dismiss it rather than leave the user with a prompt that would fail.
		self.c.OnUIThread(func() error {
			self.mergeAndRebaseHelper.DismissContinueRebasePromptIfShowing()
			return nil
		})
	}

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		// only taking over the filter if it hasn't already been set by the user.
		if conflictFileCount > 0 && prevConflictFileCount == 0 {
			if fileTreeViewModel.GetStatusFilter() == filetree.DisplayAll {
				fileTreeViewModel.SetStatusFilter(filetree.DisplayConflicted)
				self.c.Contexts().Files.GetView().Subtitle = self.c.Tr.FilterLabelConflictingFiles
			}
		} else if conflictFileCount == 0 && fileTreeViewModel.GetStatusFilter() == filetree.DisplayConflicted {
			fileTreeViewModel.SetStatusFilter(filetree.DisplayAll)
			self.c.Contexts().Files.GetView().Subtitle = ""
		}

		self.c.Model().Submodules = submoduleConfigs
		self.c.Model().Files = files
		fileTreeViewModel.SetTree()
		return nil
	})

	return nil
}

// the reflogs panel is the only panel where we cache data, in that we only
// load entries that have been created since we last ran the call. This means
// we need to be more careful with how we use this, and to ensure we're emptying
// the reflogs array when changing contexts.
// This method also manages two things: ReflogCommits and FilteredReflogCommits.
// FilteredReflogCommits are rendered in the reflogs panel, and ReflogCommits
// are used by the branches panel to obtain recency values for sorting.
// refreshReflogCommits returns the (non-filtered) ReflogCommits it loaded, so
// that a subsequent branches refresh can use them for recency sorting without
// having to read them back out of the model.
func (self *RefreshHelper) refreshReflogCommits() ([]*models.Commit, error) {
	generation := self.c.State().GetRepoGeneration()
	// pulling state into its own variable in case it gets swapped out for another state
	// and we get an out of bounds exception
	model := self.c.Model()

	// load does the git work on the worker and returns the new value for a
	// reflog slice, reading the existing slice for the incremental fetch. The
	// caller writes the result in the bounce.
	load := func(existing []*models.Commit, filterPath string, filterAuthor string) ([]*models.Commit, error) {
		var lastReflogCommit *models.Commit
		if filterPath == "" && filterAuthor == "" && len(existing) > 0 {
			lastReflogCommit = existing[0]
		}

		commits, onlyObtainedNewReflogCommits, err := self.c.Git().Loaders.ReflogCommitLoader.
			GetReflogCommits(model.HashPool, lastReflogCommit, filterPath, filterAuthor)
		if err != nil {
			return nil, err
		}

		if onlyObtainedNewReflogCommits {
			return append(commits, existing...), nil
		}
		return commits, nil
	}

	reflogCommits, err := load(model.ReflogCommits, "", "")
	if err != nil {
		return nil, err
	}

	filteredReflogCommits := reflogCommits
	if self.c.Modes().Filtering.Active() {
		filteredReflogCommits, err = load(model.FilteredReflogCommits, self.c.Modes().Filtering.GetPath(), self.c.Modes().Filtering.GetAuthor())
		if err != nil {
			return nil, err
		}
	}

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		model.ReflogCommits = reflogCommits
		model.FilteredReflogCommits = filteredReflogCommits
		return nil
	})

	self.refreshView(self.c.Contexts().ReflogCommits)
	return reflogCommits, nil
}

func (self *RefreshHelper) refreshRemotes() error {
	generation := self.c.State().GetRepoGeneration()
	prevSelectedRemote := self.c.Contexts().Remotes.GetSelected()

	remotes, err := self.c.Git().Loaders.RemoteLoader.GetRemotes()
	if err != nil {
		return err
	}

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().Remotes = remotes

		hadPrs := len(self.c.Model().PullRequestsMap) != 0
		self.rebuildPullRequestsMap()
		if !hadPrs && len(self.c.Model().PullRequestsMap) != 0 {
			// if we didn't have PRs in the map before but now we do, we need to redraw the branches view
			self.refreshView(self.c.Contexts().Branches)
		}

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
		return nil
	})

	self.refreshView(self.c.Contexts().Remotes)
	self.refreshView(self.c.Contexts().RemoteBranches)
	return nil
}

func (self *RefreshHelper) loadWorktrees() []*models.Worktree {
	worktrees, err := self.c.Git().Loaders.Worktrees.GetWorktrees()
	if err != nil {
		self.c.Log.Error(err)
		return []*models.Worktree{}
	}
	return worktrees
}

func (self *RefreshHelper) refreshWorktrees() {
	generation := self.c.State().GetRepoGeneration()

	worktrees := self.loadWorktrees()

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().Worktrees = worktrees
		return nil
	})

	// need to refresh branches because the branches view shows worktrees against
	// branches
	self.refreshView(self.c.Contexts().Branches)
	self.refreshView(self.c.Contexts().Worktrees)
}

func (self *RefreshHelper) refreshStashEntries() {
	generation := self.c.State().GetRepoGeneration()

	stashEntries := self.c.Git().Loaders.StashLoader.
		GetStashEntries(self.c.Modes().Filtering.GetPath())

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().StashEntries = stashEntries
		return nil
	})

	self.refreshView(self.c.Contexts().Stash)
}

// never call this on its own, it should only be called from within refreshCommits()
func (self *RefreshHelper) refreshStatus() {
	self.c.Mutexes().RefreshingStatusMutex.Lock()
	defer self.c.Mutexes().RefreshingStatusMutex.Unlock()

	generation := self.c.State().GetRepoGeneration()

	currentBranch := self.refsHelper.GetCheckedOutRef()
	if currentBranch == nil {
		// need to wait for branches to refresh
		return
	}

	workingTreeState := self.c.Git().Status.WorkingTreeState()
	linkedWorktreeName := self.worktreeHelper.GetLinkedWorktreeName()

	repoName := self.c.Git().RepoPaths.RepoName()

	status := presentation.FormatStatus(repoName, currentBranch, types.ItemOperationNone, linkedWorktreeName, workingTreeState, self.c.Tr, self.c.UserConfig())

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.SetViewContent(self.c.Views().Status, status)
		return nil
	})
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
	// refreshView is called from the worker goroutine that drives async
	// refreshes, so bounce to the UI thread before mutating view content.
	self.c.OnUIThread(func() error {
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
		return nil
	})
}

func (self *RefreshHelper) refreshGithubPullRequests() {
	self.c.Mutexes().RefreshingPullRequestsMutex.Lock()
	defer self.c.Mutexes().RefreshingPullRequestsMutex.Unlock()

	generation := self.c.State().GetRepoGeneration()

	clearPullRequests := func() {
		self.onUIThreadUnlessRepoChanged(generation, func() error {
			self.c.Model().PullRequests = nil
			self.c.Model().PullRequestsMap = nil
			return nil
		})
	}

	githubRemotes := getAuthenticatedGithubRemotes(self.getGithubRemotes(), self.c.Git().GitHub.GetAuthToken)
	if len(githubRemotes) == 0 {
		clearPullRequests()
		return
	}

	baseInfo := getGithubBaseRemote(githubRemotes, self.c.Git().GitHub.ConfiguredBaseRemoteName())
	if baseInfo == nil {
		clearPullRequests()

		if !self.githubBaseRemotePromptDismissed[self.c.Git().RepoPaths.RepoPath()] {
			self.promptForBaseGithubRepo(githubRemotes)
		}
		return
	}

	self.setGithubPullRequests(baseInfo)
}

type githubRemoteInfo struct {
	remote      *models.Remote
	serviceInfo hosting_service.ServiceInfo
	authToken   string
}

func (self *RefreshHelper) getGithubRemotes() []githubRemoteInfo {
	return lo.FilterMap(self.c.Model().Remotes, func(remote *models.Remote, _ int) (githubRemoteInfo, bool) {
		if len(remote.Urls) == 0 {
			return githubRemoteInfo{}, false
		}
		serviceInfo, err := self.c.Git().HostingService.GetServiceInfo(remote.Urls[0])
		if err != nil || serviceInfo.Provider != "github" {
			return githubRemoteInfo{}, false
		}
		return githubRemoteInfo{remote: remote, serviceInfo: serviceInfo}, true
	})
}

// getAuthenticatedGithubRemotes drops remotes for which no auth token is
// available and attaches the resolved token to the rest. Token lookups are
// cached by host so that multiple remotes pointing at the same instance
// (e.g. origin + a fork on github.com) only trigger one lookup.
func getAuthenticatedGithubRemotes(githubRemotes []githubRemoteInfo, getAuthToken func(host string) string) []githubRemoteInfo {
	tokensByHost := map[string]string{}
	return lo.FilterMap(githubRemotes, func(info githubRemoteInfo, _ int) (githubRemoteInfo, bool) {
		host := info.serviceInfo.WebDomain
		token, cached := tokensByHost[host]
		if !cached {
			token = getAuthToken(host)
			tokensByHost[host] = token
		}
		if token == "" {
			return githubRemoteInfo{}, false
		}
		info.authToken = token
		return info, true
	})
}

func getGithubBaseRemote(githubRemotes []githubRemoteInfo, configuredRemoteName string) *githubRemoteInfo {
	findRemoteByName := func(name string) *githubRemoteInfo {
		info, ok := lo.Find(githubRemotes, func(info githubRemoteInfo) bool {
			return info.remote.Name == name
		})
		if !ok {
			return nil
		}
		return &info
	}

	if configuredRemoteName != "" {
		return findRemoteByName(configuredRemoteName)
	}

	if len(githubRemotes) == 1 {
		return &githubRemotes[0]
	}

	// Not sure if "upstream" is really a common convention for the name of the remote that PRs are
	// made against, but if it exists it's pretty likely to be the one we want.
	if info := findRemoteByName("upstream"); info != nil {
		return info
	}

	return nil
}

func (self *RefreshHelper) promptForBaseGithubRepo(githubRemotes []githubRemoteInfo) {
	menuItems := lo.Map(githubRemotes, func(info githubRemoteInfo, _ int) *types.MenuItem {
		return &types.MenuItem{
			LabelColumns: []string{info.remote.Name, style.FgCyan.Sprint(info.serviceInfo.RepoName)},
			OnPress: func() error {
				return self.c.WithWaitingStatus(self.c.Tr.FetchingPullRequests, func(gocui.Task) error {
					if err := self.c.Git().GitHub.SetConfiguredBaseRemoteName(info.remote.Name); err != nil {
						self.c.Log.Error(err)
					}

					self.setGithubPullRequests(&info)
					return nil
				})
			},
		}
	})

	_ = self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.SelectRemoteRepository,
		Items: menuItems,
		OnCancel: func() error {
			if self.githubBaseRemotePromptDismissed == nil {
				self.githubBaseRemotePromptDismissed = make(map[string]bool)
			}
			self.githubBaseRemotePromptDismissed[self.c.Git().RepoPaths.RepoPath()] = true
			return nil
		},
	})
}

func (self *RefreshHelper) rebuildPullRequestsMap() {
	self.c.Model().PullRequestsMap = git_commands.GenerateGithubPullRequestMap(
		self.c.Model().PullRequests,
		self.c.Model().Branches,
		self.c.Model().Remotes,
	)
}

func (self *RefreshHelper) setGithubPullRequests(baseInfo *githubRemoteInfo) {
	generation := self.c.State().GetRepoGeneration()

	if len(self.c.Model().Branches) == 0 {
		return
	}

	branches := lo.Filter(self.c.Model().Branches, func(branch *models.Branch, _ int) bool {
		return branch.IsTrackingRemote()
	})
	branchNames := lo.Map(branches, func(branch *models.Branch, _ int) string {
		return branch.UpstreamBranch
	})

	prs, err := self.c.Git().GitHub.FetchRecentPRs(branchNames, &baseInfo.serviceInfo, baseInfo.authToken)
	if err != nil {
		self.c.Log.Error("error fetching pull requests from GitHub: " + err.Error())
		return
	}

	self.savePullRequestsToCache(prs)

	self.onUIThreadUnlessRepoChanged(generation, func() error {
		self.c.Model().PullRequests = prs
		// Rebuilding here rather than on the worker means the map is built from
		// the branches and remotes as they are on the UI thread, after their
		// own refreshes' bounces have applied.
		self.rebuildPullRequestsMap()
		self.c.PostRefreshUpdate(self.c.Contexts().Branches)
		return nil
	})
}

func (self *RefreshHelper) savePullRequestsToCache(prs []*models.GithubPullRequest) {
	repoPath := self.c.Git().RepoPaths.RepoPath()
	cached := lo.Map(prs, func(pr *models.GithubPullRequest, _ int) config.CachedPullRequest {
		return config.CachedPullRequest{
			HeadRefName:         pr.HeadRefName,
			Number:              pr.Number,
			Title:               pr.Title,
			State:               pr.State,
			Url:                 pr.Url,
			HeadRepositoryOwner: pr.HeadRepositoryOwner.Login,
		}
	})

	appState := self.c.GetAppState()
	if appState.GithubPullRequests == nil {
		appState.GithubPullRequests = make(map[string][]config.CachedPullRequest)
	}
	appState.GithubPullRequests[repoPath] = cached
	self.c.SaveAppStateAndLogError()
}
