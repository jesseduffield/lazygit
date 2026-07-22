package helpers

import (
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands"
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

	// branchLoadSeq hands out a monotonically increasing sequence number to
	// each branch load (via Add, on the worker); appliedBranchLoadSeq is the
	// highest sequence whose result has been written to the model (touched only
	// on the UI thread, inside the bounce). Together they let a branch load's
	// bounce drop its write if a later-started load has already applied, so
	// concurrent branch loads don't clobber each other out of order.
	branchLoadSeq        atomic.Int64
	appliedBranchLoadSeq int64
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
	self.performRefresh(options, false, false)
}

// RefreshBlockingInput is Refresh for handlers whose next keypress may depend
// on the state the refresh produces. See IGuiCommon.RefreshBlockingInput.
func (self *RefreshHelper) RefreshBlockingInput(options types.RefreshOptions) {
	self.performRefresh(options, false, true)
}

// RefreshFromWorker is Refresh for callers already running on a worker
// goroutine (e.g. inside a WithWaitingStatus handler) rather than the UI
// thread. See IGuiCommon.RefreshFromWorker.
func (self *RefreshHelper) RefreshFromWorker(options types.RefreshOptions) {
	self.performRefresh(options, true, false)
}

type refreshEnv struct {
	// Whether everything this refresh dispatches uses the background task
	// variants, which don't count towards lazygit being busy — so the refresh
	// doesn't block switching repos. Set for refreshes initiated by a
	// background routine, and for foreground ones that opted in via
	// RefreshOptions.DontBlockRepoSwitch.
	background bool

	// Whether the refresh was initiated by an unattended background routine
	// (RefreshOptions.Background) rather than by user activity. The files
	// refresh uses this to decide whether git may take optional locks and
	// persist its refreshed stat cache.
	backgroundRoutine bool

	// the repo generation captured when the refresh started
	generation int

	// the git command instance captured when the refresh started. The refresh
	// workers run their git commands through this rather than reading the live
	// instance: a repo switch mid-refresh replaces the live instance (and the
	// process cwd), while this one keeps addressing the repo the refresh was
	// started for (its commands are pinned to that repo's directory).
	git *commands.GitCommand

	// When non-nil, each scope's UI-thread bounce is collected here instead of
	// being dispatched as it's produced, so they can all be applied in a single
	// frame once the whole refresh is done (see RefreshOptions.BatchUIUpdates).
	// Held by pointer so the copies of env that flow through the scope functions
	// all share the one batch.
	batch *refreshBounceBatch
}

// refreshBounceBatch collects the UI-thread bounces of a batched refresh so they
// can be applied together in one frame rather than one scope at a time. The
// scopes run on separate worker goroutines and add concurrently, hence the
// mutex. Once the refresh starts flushing it closes the batch, so that any
// bounces enqueued afterwards — the nested ones a flushed bounce produces in
// turn, e.g. scrolling the selection into view — are dispatched immediately as
// ordinary follow-ups instead of being collected into a batch that nothing
// will drain.
type refreshBounceBatch struct {
	mutex  deadlock.Mutex
	funcs  []func()
	closed bool
}

// add collects f and returns true. Once the batch is closed it collects nothing
// and returns false, telling the caller to dispatch f immediately instead.
func (self *refreshBounceBatch) add(f func()) bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.closed {
		return false
	}
	self.funcs = append(self.funcs, f)
	return true
}

// close marks the batch flushed and returns everything collected so far.
func (self *refreshBounceBatch) close() []func() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.closed = true
	return self.funcs
}

func (self *RefreshHelper) performRefresh(options types.RefreshOptions, calledFromWorker bool, blockInput bool) {
	startTime := time.Now()

	// A refresh from a worker blocks that worker until it's done; one from the
	// UI thread returns immediately and finishes in the background.
	syncOrAsync := "async"
	if calledFromWorker {
		syncOrAsync = "sync"
	}
	if options.Scope == nil {
		self.c.Log.Infof("refreshing all scopes (%s)", syncOrAsync)
	} else {
		self.c.Log.Infof(
			"refreshing the following scopes (%s): %s",
			syncOrAsync,
			strings.Join(getScopeNames(options.Scope), ","),
		)
	}

	// Debug-only guard: every refresh must be issued from the entry point that
	// matches its goroutine — Refresh on the UI thread, RefreshFromWorker on a
	// worker. goid stays out of production control flow (debug only).
	if self.c.GetConfig().GetDebug() && self.c.GocuiGui().IsUIThread() == calledFromWorker {
		panic("Refresh called from a worker, or RefreshFromWorker called from the UI thread")
	}

	if options.Then != nil && options.DontBlockRepoSwitch {
		// Then is not generation-guarded, so if a switch crossed the refresh it
		// would run against the newly switched-to repo. A refresh carrying a
		// Then must keep blocking switches.
		panic("a refresh with a Then callback must not set DontBlockRepoSwitch")
	}

	// A RefreshBlockingInput caller wants keyboard input withheld until the
	// refreshed state is in place (see IGuiCommon.RefreshBlockingInput). Begin
	// the block synchronously here in the calling handler, so that no keypress
	// can slip through before it; the finishing step ends it from a callback
	// queued behind the refresh's own updates (see waitAndFinalize). Demos
	// take the blocking inline path below and need none of this.
	blockInputUntilDone := blockInput && !self.c.InDemo()
	if blockInputUntilDone {
		self.c.GocuiGui().BeginBlockingEvents()
	}

	// Capture the refresh's baseline once, here at the start: the repo
	// generation that every scope's bounce is guarded against, and the git
	// command instance the scopes run their commands through. The two are
	// captured together on the UI thread so that they can't straddle a repo
	// switch (which runs on the UI thread): pairing the old repo's instance
	// with the new repo's generation would let a refresh compute data from
	// the old repo and write it into the new repo's model unguarded. With a
	// consistent pair, a switch-crossing refresh keeps running its commands
	// against the repo it started in, and the generation guard drops its
	// writes.
	env := refreshEnv{
		background:        options.Background || options.DontBlockRepoSwitch,
		backgroundRoutine: options.Background,
	}
	self.captureOnUIThread(calledFromWorker, env.background, func() {
		env.generation = self.c.State().GetRepoGeneration()
		env.git = self.c.Git()
	})
	if options.BatchUIUpdates {
		env.batch = &refreshBounceBatch{}
	}

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
	//   - pull requests are fetched for the tracking branches against the
	//     remotes, so refresh both alongside to fetch against fresh data
	if scopeSet.Includes(types.REFLOG) || scopeSet.Includes(types.BISECT_INFO) {
		scopeSet.Add(types.COMMITS, types.BRANCHES)
	}
	if scopeSet.Includes(types.SUBMODULES) {
		scopeSet.Add(types.FILES)
	}
	if scopeSet.Includes(types.FILES) {
		scopeSet.Add(types.MERGE_CONFLICTS)
	}
	if scopeSet.Includes(types.PULL_REQUESTS) {
		scopeSet.Add(types.BRANCHES, types.REMOTES)
	}

	// Capture the refs snapshot now, before we start reading git's state
	// below, rather than after. This is important to guard against the race
	// of git's state changing externally while (or right after) we are
	// refreshing; the risk is one potential extra refresh, but capturing the
	// snapshot at the end would risk missing one, which is worse.
	self.updateRefsSnapshotIfRelevant(scopeSet, env)

	wg := sync.WaitGroup{}
	refresh := func(name string, f func()) {
		wg.Add(1)
		// Each scope runs on its own goroutine, joined by the wg.Wait in
		// waitAndFinalize. They don't need to be registered as gocui tasks for
		// repo-switch safety: performRefresh always runs under a task that stays
		// busy until that wg.Wait returns — the calling worker's task when
		// called from a worker, or the waitAndFinalize worker task when called
		// from the UI thread (created before the triggering event's task ends,
		// so there's no gap) — and that task already covers the whole refresh.
		go utils.Safe(func() {
			t := time.Now()
			defer wg.Done()
			f()
			self.c.Log.Infof("refreshed %s in %s", name, time.Since(t))
		})
	}

	branchesAndRemotesWg := sync.WaitGroup{}
	// The pull-request fetch (below) needs the just-loaded branches and
	// remotes. Their model writes are bounced onto the UI thread, so the
	// fetch worker can't read them back from the model without racing (and
	// would see the pre-refresh values); instead the branches and remotes
	// loads stash what they loaded here, and the wait on
	// branchesAndRemotesWg gives the fetch the happens-before to read them.
	var loadedBranches []*models.Branch
	var loadedRemotes []*models.Remote
	includeWorktreesWithBranches := false
	if scopeSet.Includes(types.COMMITS) || scopeSet.Includes(types.BRANCHES) {
		// whenever we change commits, we should update branches because the upstream/downstream
		// counts can change. Whenever we change branches we should also change commits
		// e.g. in the case of switching branches.
		// Capture the commits, reflog and branches refresh inputs (model,
		// contexts, modes) on the UI thread, before the git work is dispatched
		// to a worker, so the workers compute from an immutable snapshot
		// instead of reading state the UI thread concurrently mutates.
		var capturedCommits capturedCommitState
		var capturedReflog capturedReflogState
		var capturedBranches capturedBranchState
		self.captureOnUIThread(calledFromWorker, env.background, func() {
			capturedCommits = self.captureCommitsState(options.CommitSelection)
			capturedReflog = self.captureReflogState()
			capturedBranches = self.captureBranchState()
		})
		refresh("commits and commit files", func() {
			self.refreshCommitsAndCommitFiles(capturedCommits, options.CommitSelection, env)
		})

		includeWorktreesWithBranches = scopeSet.Includes(types.WORKTREES)
		if self.c.UserConfig().Git.LocalBranchSortOrder == "recency" {
			branchesAndRemotesWg.Add(1)
			refresh("reflog and branches", func() {
				loadedBranches = self.refreshReflogAndBranches(capturedReflog, capturedBranches, includeWorktreesWithBranches, options.BranchSelection, options.SelectTopReflogCommit, env)
				branchesAndRemotesWg.Done()
			})
		} else {
			branchesAndRemotesWg.Add(1)
			refresh("branches", func() {
				// Not a recency sort, so branches doesn't depend on the reflog
				// being fresh; it runs concurrently with the reflog refresh
				// below and uses the reflog we captured up front, as it always has.
				loadedBranches = self.refreshBranches(capturedBranches, includeWorktreesWithBranches, options.BranchSelection, true, capturedReflog.reflogCommits, env)
				branchesAndRemotesWg.Done()
			})
			refresh("reflog", func() {
				_, _ = self.refreshReflogCommits(capturedReflog, env, options.SelectTopReflogCommit)
			})
		}
	} else if scopeSet.Includes(types.REBASE_COMMITS) {
		// the above block handles rebase commits so we only need to call this one
		// if we've asked specifically for rebase commits and not those other things
		var rebaseHashPool *utils.StringPool
		var rebaseCommits []*models.Commit
		self.captureOnUIThread(calledFromWorker, env.background, func() {
			rebaseHashPool, rebaseCommits = self.captureRebaseCommitState()
		})
		refresh("rebase commits", func() { _ = self.refreshRebaseCommits(rebaseHashPool, rebaseCommits, env) })
	}

	if scopeSet.Includes(types.SUB_COMMITS) {
		var capturedSubCommits capturedSubCommitState
		self.captureOnUIThread(calledFromWorker, env.background, func() {
			capturedSubCommits = self.captureSubCommitState()
		})
		refresh("sub commits", func() { _ = self.refreshSubCommitsWithLimit(capturedSubCommits, env) })
	}

	// reason we're not doing this if the COMMITS type is included is that if the COMMITS type _is_ included we will refresh the commit files context anyway
	if scopeSet.Includes(types.COMMIT_FILES) && !scopeSet.Includes(types.COMMITS) {
		var capturedCommitFiles capturedCommitFilesState
		self.captureOnUIThread(calledFromWorker, env.background, func() {
			capturedCommitFiles = self.captureCommitFilesState()
		})
		refresh("commit files", func() { _ = self.refreshCommitFilesContext(capturedCommitFiles, env) })
	}

	fileWg := sync.WaitGroup{}
	if scopeSet.Includes(types.FILES) {
		var capturedFiles capturedFilesState
		self.captureOnUIThread(calledFromWorker, env.background, func() {
			capturedFiles = self.captureFilesState()
		})
		fileWg.Add(1)
		refresh("files", func() {
			_ = self.refreshFilesAndSubmodules(capturedFiles, env)
			fileWg.Done()
		})
	}

	if scopeSet.Includes(types.STASH) {
		var stashFilterPath string
		self.captureOnUIThread(calledFromWorker, env.background, func() {
			stashFilterPath = self.c.Modes().Filtering.GetPath()
		})
		refresh("stash", func() { self.refreshStashEntries(stashFilterPath, env) })
	}

	if scopeSet.Includes(types.TAGS) {
		refresh("tags", func() { _ = self.refreshTags(env) })
	}

	if scopeSet.Includes(types.REMOTES) {
		// Capture the previously-selected remote on the UI thread; the worker
		// needs it to keep the remote-branches selection valid, and reading
		// the Remotes context off the UI thread races its render.
		var prevSelectedRemote *models.Remote
		self.captureOnUIThread(calledFromWorker, env.background, func() {
			prevSelectedRemote = self.c.Contexts().Remotes.GetSelected()
		})
		branchesAndRemotesWg.Add(1)
		refresh("remotes", func() {
			loadedRemotes, _ = self.refreshRemotes(prevSelectedRemote, env)
			branchesAndRemotesWg.Done()
		})
	}

	if scopeSet.Includes(types.PULL_REQUESTS) {
		// Fetching pull requests talks to the GitHub API over the network; on
		// a bad connection that request can stall for a long time. It runs no
		// git commands against the repo, and its model writes are guarded by
		// the repo generation (a repo switch mid-fetch simply drops the
		// result), so it is safe to run as a background task even when the
		// enclosing refresh is a foreground one — a foreground task would
		// block repo switching for as long as the request takes. The env copy
		// makes the downstream UI-thread bounces background as well.
		prEnv := env
		prEnv.background = true
		self.c.OnWorkerBackground(func(gocui.Task) error {
			branchesAndRemotesWg.Wait()

			t := time.Now()
			// Use the branches and remotes the loads above stashed, not
			// Model().Branches/Remotes: those writes are bounced onto the
			// UI thread and may not have landed on this worker yet. The
			// wait above orders us after both loads have stashed theirs.
			self.refreshGithubPullRequests(loadedBranches, loadedRemotes, prEnv)
			self.c.Log.Infof("refreshed pull requests in %s", time.Since(t))
			return nil
		})
	}

	if scopeSet.Includes(types.WORKTREES) && !includeWorktreesWithBranches {
		refresh("worktrees", func() { self.refreshWorktrees(env) })
	}

	if scopeSet.Includes(types.STAGING) {
		refresh("staging", func() {
			fileWg.Wait()
			// Bounce onto the UI thread so this runs after the files
			// scope's model-update bounce — RefreshStagingPanel reads
			// Model.Files (via Files.GetSelected) and would otherwise
			// see the pre-refresh model. Guard on the generation so a
			// repo switch mid-refresh drops it, like the model bounces.
			self.onUIThreadUnlessRepoChanged(env, func() {
				self.stagingHelper.RefreshStagingPanel(types.OnFocusOpts{})
			})
		})
	}

	if scopeSet.Includes(types.PATCH_BUILDING) {
		refresh("patch building", func() {
			// Bounce onto the UI thread, like the staging panel above:
			// RefreshPatchBuildingPanel reads the commit-files selection and
			// sets the patch view's origin, neither of which may run off the UI
			// thread. Guard on the generation so a repo switch mid-refresh drops
			// it, like the model bounces.
			self.onUIThreadUnlessRepoChanged(env, func() {
				self.patchBuildingHelper.RefreshPatchBuildingPanel(types.OnFocusOpts{})
			})
		})
	}

	if scopeSet.Includes(types.MERGE_CONFLICTS) {
		refresh("merge conflicts", func() {
			// Bounce onto the UI thread, like the staging and patch-building
			// panels above: RefreshMergeState reads the current context and
			// renders (or escapes) the merge-conflicts view, none of which may
			// run off the UI thread.
			self.onUIThreadUnlessRepoChanged(env, func() {
				_ = self.mergeConflictsHelper.RefreshMergeState()
			})
		})
	}

	self.refreshStatus(env)

	waitAndFinalize := func() {
		wg.Wait()

		if env.batch != nil {
			// Apply all the scopes' collected bounces in a single UI-thread task,
			// so they land in one frame: gocui drains every queued event before it
			// redraws, so one task means one repaint. Bounces enqueued from within
			// these (see refreshBounceBatch) run as ordinary follow-ups.
			bounces := env.batch.close()
			self.onUIThread(env.background, func() error {
				for _, bounce := range bounces {
					bounce()
				}
				return nil
			})
		}

		if options.Then != nil {
			// Queue Then via OnUIThread so it runs *after* the refresh-scope
			// functions' model-update bounces (which are already queued by
			// now), not synchronously here — at this point the workers have
			// returned but their bounces haven't been processed yet, so
			// invoking Then synchronously would run it on a model that's
			// still pre-refresh.
			self.onUIThread(env.background, options.Then)
		}

		if blockInputUntilDone {
			// Queued after the scopes' model bounces and Then, so by the time
			// this runs — and the keys buffered during the refresh replay —
			// the refreshed state is in place.
			self.c.OnUIThread(func() error {
				return self.c.GocuiGui().EndBlockingEvents()
			})
		}

		self.c.Log.Infof("Refresh took %s", time.Since(startTime))
	}

	// waitAndFinalize blocks until every scope is done. Run it inline when we're
	// already on a worker (or in a demo, for a deterministic single frame); when
	// we're on the UI thread, dispatch it to a worker so it doesn't block the UI.
	if calledFromWorker || self.c.InDemo() {
		waitAndFinalize()
	} else {
		self.onWorker(env.background, func(t gocui.Task) error {
			waitAndFinalize()
			return nil
		})
	}
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
func (self *RefreshHelper) updateRefsSnapshotIfRelevant(scopeSet *set.Set[types.RefreshableView], env refreshEnv) {
	if !scopeSet.Includes(types.COMMITS) && !scopeSet.Includes(types.BRANCHES) {
		return
	}

	snapshot, err := env.git.Status.RefsSnapshot()
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

// During startup, the bottleneck is fetching the reflog entries, which we need
// in order to sort the branches by recency. So we have two phases: INITIAL and
// COMPLETE. In the INITIAL phase we don't have any reflog commits yet, so we
// show the branches right away sorted by whatever we have (typically nothing,
// i.e. not by recency), then load the reflog on a worker and refresh the
// branches again, this time recency-sorted. From then on we're in the COMPLETE
// phase and load the reflog synchronously before refreshing the branches.
//
// The immediate refresh must run before we spawn the async one, not after: that
// order gives the immediate (non-recency) load a lower branch-load sequence
// than the async (recency) load, so the sequence guard in refreshBranches keeps
// the recency-sorted result even if the two loads' bounces land out of order.
// capturedReflogState holds the reflog refresh's model/mode inputs, gathered on
// the UI thread before the git work runs. The existing reflog slices feed the
// incremental fetch (we only load entries newer than the ones we already have).
type capturedReflogState struct {
	reflogCommits         []*models.Commit
	filteredReflogCommits []*models.Commit
	hashPool              *utils.StringPool
	filteringActive       bool
	filterPath            string
	filterAuthor          string
}

// captureReflogState reads the reflog refresh's inputs into an immutable
// snapshot. It must run on the UI thread.
func (self *RefreshHelper) captureReflogState() capturedReflogState {
	return capturedReflogState{
		reflogCommits:         self.c.Model().ReflogCommits,
		filteredReflogCommits: self.c.Model().FilteredReflogCommits,
		hashPool:              self.c.Model().HashPool,
		filteringActive:       self.c.Modes().Filtering.Active(),
		filterPath:            self.c.Modes().Filtering.GetPath(),
		filterAuthor:          self.c.Modes().Filtering.GetAuthor(),
	}
}

// capturedBranchState holds the branches refresh's model inputs, gathered on the
// UI thread before the git work runs. oldBranches is used only to carry over the
// previous BehindBaseBranch values (to reduce flicker) — an atomic each, so a
// pre-refresh snapshot serves both the immediate and recency loads identically.
type capturedBranchState struct {
	mainBranches *git_commands.MainBranches
	oldBranches  []*models.Branch
}

// captureBranchState reads the branches refresh's model inputs into an immutable
// snapshot. It must run on the UI thread.
func (self *RefreshHelper) captureBranchState() capturedBranchState {
	return capturedBranchState{
		mainBranches: self.c.Model().MainBranches,
		oldBranches:  self.c.Model().Branches,
	}
}

func (self *RefreshHelper) refreshReflogAndBranches(capturedReflog capturedReflogState, capturedBranches capturedBranchState, refreshWorktrees bool, branchSelection types.BranchSelectionBehavior, selectTopReflogCommit bool, env refreshEnv) []*models.Branch {
	switch self.c.State().GetRepoState().GetStartupStage() {
	case types.INITIAL:
		// Return the immediate (non-recency) load's branches; the recency-sorted
		// reload below runs on its own worker after we return. Both hold the same
		// set of branches, which is all the caller (the PR fetch) needs.
		branches := self.refreshBranches(capturedBranches, refreshWorktrees, branchSelection, false, capturedReflog.reflogCommits, env)

		self.onWorker(env.background, func(_ gocui.Task) error {
			reflogCommits, _ := self.refreshReflogCommits(capturedReflog, env, false)
			self.refreshBranches(capturedBranches, false, types.SelectCheckedOutBranch, true, reflogCommits, env)
			self.c.State().GetRepoState().SetStartupStage(types.COMPLETE)
			return nil
		})

		return branches

	case types.COMPLETE:
		reflogCommits, _ := self.refreshReflogCommits(capturedReflog, env, selectTopReflogCommit)
		return self.refreshBranches(capturedBranches, refreshWorktrees, branchSelection, true, reflogCommits, env)
	}

	return nil
}

// capturedCommitState holds everything the commits refresh reads from the
// model, contexts, and modes. It is gathered on the UI thread (see
// captureCommitsState) before the git work is dispatched to a worker, so the
// worker computes from an immutable snapshot rather than reading state the UI
// thread concurrently mutates.
type capturedCommitState struct {
	selectionRange       *localCommitSelectionRange
	limitCommits         bool
	showWholeGitGraph    bool
	filterRefs           []string
	filterPath           string
	filterAuthor         string
	mainBranches         *git_commands.MainBranches
	hashPool             *utils.StringPool
	parentIsLocalCommits bool
}

// captureCommitsState reads the commits refresh's model/context/mode inputs
// into an immutable snapshot. It must run on the UI thread.
func (self *RefreshHelper) captureCommitsState(commitSelection types.CommitSelectionBehavior) capturedCommitState {
	var selectionRange *localCommitSelectionRange
	if commitSelection == types.KeepCommitSelectionByHash {
		selectedIdx, rangeStartIdx, rangeSelectMode := self.c.Contexts().LocalCommits.GetSelectionRangeAndMode()
		selectionRange = captureLocalCommitSelectionRange(self.c.Model().Commits, selectedIdx, rangeStartIdx, rangeSelectMode)
	}

	parentCtx := self.c.Contexts().CommitFiles.GetParentContext()

	return capturedCommitState{
		selectionRange:       selectionRange,
		limitCommits:         self.c.Contexts().LocalCommits.GetLimitCommits(),
		showWholeGitGraph:    self.c.Contexts().LocalCommits.GetShowWholeGitGraph(),
		filterRefs:           self.c.Contexts().LocalCommits.GetFilterRefs(),
		filterPath:           self.c.Modes().Filtering.GetPath(),
		filterAuthor:         self.c.Modes().Filtering.GetAuthor(),
		mainBranches:         self.c.Model().MainBranches,
		hashPool:             self.c.Model().HashPool,
		parentIsLocalCommits: parentCtx != nil && parentCtx.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY,
	}
}

func (self *RefreshHelper) refreshCommitsAndCommitFiles(captured capturedCommitState, commitSelection types.CommitSelectionBehavior, env refreshEnv) {
	_ = self.refreshCommitsWithLimit(captured, commitSelection, env)
	if captured.parentIsLocalCommits {
		// This makes sense when we've e.g. just amended a commit, meaning we get a new commit hash at the same position.
		// However if we've just added a brand new commit, it pushes the list down by one and so we would end up
		// showing the contents of a different commit than the one we initially entered.
		// Ideally we would know when to refresh the commit files context and when not to,
		// or perhaps we could just pop that context off the stack whenever cycling windows.
		// For now the awkwardness remains.
		//
		// The commit selection is restored in refreshCommitsWithLimit's bounce,
		// so read it on the UI thread after that bounce; then load the commit
		// files back on a worker (refreshCommitFilesContext does git work).
		self.onUIThreadUnlessRepoChanged(env, func() {
			commit := self.c.Contexts().LocalCommits.GetSelected()
			if commit != nil && commit.RefName() != "" {
				refRange := self.c.Contexts().LocalCommits.GetSelectedRefRangeForDiffFiles()
				self.c.Contexts().CommitFiles.ReInit(commit, refRange)
				// Capture the diff endpoints here, on the UI thread and after
				// ReInit has set them, before dispatching the git work.
				capturedCommitFiles := self.captureCommitFilesState()
				self.onWorker(env.background, func(gocui.Task) error {
					_ = self.refreshCommitFilesContext(capturedCommitFiles, env)
					return nil
				})
			}
		})
	}
}

func (self *RefreshHelper) determineCheckedOutRef(env refreshEnv) models.Ref {
	if rebasedBranch := env.git.Status.BranchBeingRebased(); rebasedBranch != "" {
		// During a rebase we're on a detached head, so cannot determine the
		// branch name in the usual way. We need to read it from the
		// ".git/rebase-merge/head-name" file instead.
		return &models.Branch{Name: strings.TrimPrefix(rebasedBranch, "refs/heads/")}
	}

	if bisectInfo := env.git.Bisect.GetInfo(); bisectInfo.Bisecting() && bisectInfo.GetStartHash() != "" {
		// Likewise, when we're bisecting we're on a detached head as well. In
		// this case we read the branch name from the ".git/BISECT_START" file.
		return &models.Branch{Name: bisectInfo.GetStartHash()}
	}

	// In all other cases, get the branch name by asking git what branch is
	// checked out. Note that if we're on a detached head (for reasons other
	// than rebasing or bisecting, i.e. it was explicitly checked out), then
	// this will return an empty string.
	if branchName, err := env.git.Branch.CurrentBranchName(); err == nil && branchName != "" {
		return &models.Branch{Name: branchName}
	}

	// Should never get here unless the working copy is corrupt
	return nil
}

func (self *RefreshHelper) refreshCommitsWithLimit(captured capturedCommitState, commitSelection types.CommitSelectionBehavior, env refreshEnv) error {
	checkedOutRef := self.determineCheckedOutRef(env)
	refName, bisectInfo := self.refForLog(env)
	commits, err := env.git.Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                captured.limitCommits,
			FilterPath:           captured.filterPath,
			FilterAuthor:         captured.filterAuthor,
			IncludeRebaseCommits: true,
			RefName:              refName,
			RefForPushedStatus:   checkedOutRef,
			All:                  captured.showWholeGitGraph,
			FilterRefs:           captured.filterRefs,
			MainBranches:         captured.mainBranches,
			HashPool:             captured.hashPool,
		},
	)
	if err != nil {
		return err
	}
	workingTreeState := env.git.Status.WorkingTreeState()

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().BisectInfo = bisectInfo
		self.c.Model().Commits = commits
		self.RefreshAuthors(commits)
		self.c.Model().WorkingTreeStateAtLastCommitRefresh = workingTreeState
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
			if captured.selectionRange != nil {
				selectedIdx, rangeStartIdx, didMove, found := findLocalCommitSelectionRange(commits, captured.selectionRange)
				if found {
					self.c.Contexts().LocalCommits.SetSelectionRangeAndMode(selectedIdx, rangeStartIdx, captured.selectionRange.mode)
					scrollSelectionIntoView = didMove
				}
			}
		case types.KeepCommitSelectionIndex:
			// The caller set the selection index deliberately; leave it untouched.
		}

		if scrollSelectionIntoView {
			// Enqueued from within this bounce so it runs after refreshView's
			// render below (which was enqueued first), matching the previous
			// ordering where FocusLine ran after the view was re-rendered.
			self.onUIThreadUnlessRepoChanged(env, func() {
				self.c.Contexts().LocalCommits.FocusLine(true)
			})
		}
	})

	self.refreshView(self.c.Contexts().LocalCommits, env)
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

// capturedSubCommitState holds the sub-commits refresh's model/context/mode
// inputs, gathered on the UI thread (see captureSubCommitState) before the git
// work is dispatched to a worker.
type capturedSubCommitState struct {
	ref                     models.Ref
	limitCommits            bool
	refToShowDivergenceFrom string
	filterPath              string
	filterAuthor            string
	mainBranches            *git_commands.MainBranches
	hashPool                *utils.StringPool
}

// captureSubCommitState reads the sub-commits refresh's inputs into an immutable
// snapshot. It must run on the UI thread.
func (self *RefreshHelper) captureSubCommitState() capturedSubCommitState {
	return capturedSubCommitState{
		ref:                     self.c.Contexts().SubCommits.GetRef(),
		limitCommits:            self.c.Contexts().SubCommits.GetLimitCommits(),
		refToShowDivergenceFrom: self.c.Contexts().SubCommits.GetRefToShowDivergenceFrom(),
		filterPath:              self.c.Modes().Filtering.GetPath(),
		filterAuthor:            self.c.Modes().Filtering.GetAuthor(),
		mainBranches:            self.c.Model().MainBranches,
		hashPool:                self.c.Model().HashPool,
	}
}

func (self *RefreshHelper) refreshSubCommitsWithLimit(captured capturedSubCommitState, env refreshEnv) error {
	if captured.ref == nil {
		return nil
	}

	commits, err := env.git.Loaders.CommitLoader.GetCommits(
		git_commands.GetCommitsOptions{
			Limit:                   captured.limitCommits,
			FilterPath:              captured.filterPath,
			FilterAuthor:            captured.filterAuthor,
			IncludeRebaseCommits:    false,
			RefName:                 captured.ref.FullRefName(),
			RefToShowDivergenceFrom: captured.refToShowDivergenceFrom,
			RefForPushedStatus:      captured.ref,
			MainBranches:            captured.mainBranches,
			HashPool:                captured.hashPool,
		},
	)
	if err != nil {
		return err
	}
	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().SubCommits = commits
		self.RefreshAuthors(commits)
	})

	self.refreshView(self.c.Contexts().SubCommits, env)
	return nil
}

func (self *RefreshHelper) RefreshAuthors(commits []*models.Commit) {
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

// capturedCommitFilesState holds the commit-files refresh's context/mode inputs
// (the diff endpoints), gathered on the UI thread before the git work runs.
type capturedCommitFilesState struct {
	from    string
	to      string
	reverse bool
}

// captureCommitFilesState reads the commit-files refresh's diff endpoints into
// an immutable snapshot. It must run on the UI thread.
func (self *RefreshHelper) captureCommitFilesState() capturedCommitFilesState {
	from, to := self.c.Contexts().CommitFiles.GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	return capturedCommitFilesState{from: from, to: to, reverse: reverse}
}

func (self *RefreshHelper) refreshCommitFilesContext(captured capturedCommitFilesState, env refreshEnv) error {
	files, err := env.git.Loaders.CommitFileLoader.GetFilesInDiff(captured.from, captured.to, captured.reverse)
	if err != nil {
		return err
	}
	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().CommitFiles = files
		self.c.Contexts().CommitFiles.CommitFileTreeViewModel.SetTree()
	})
	self.refreshView(self.c.Contexts().CommitFiles, env)
	return nil
}

// captureRebaseCommitState reads the rebase-commits refresh's model inputs into
// an immutable snapshot. It must run on the UI thread.
func (self *RefreshHelper) captureRebaseCommitState() (hashPool *utils.StringPool, commits []*models.Commit) {
	return self.c.Model().HashPool, self.c.Model().Commits
}

func (self *RefreshHelper) refreshRebaseCommits(hashPool *utils.StringPool, commits []*models.Commit, env refreshEnv) error {
	updatedCommits, err := env.git.Loaders.CommitLoader.MergeRebasingCommits(hashPool, commits)
	if err != nil {
		return err
	}
	workingTreeState := env.git.Status.WorkingTreeState()

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().Commits = updatedCommits
		self.c.Model().WorkingTreeStateAtLastCommitRefresh = workingTreeState
	})

	self.refreshView(self.c.Contexts().LocalCommits, env)
	return nil
}

func (self *RefreshHelper) refreshTags(env refreshEnv) error {
	tags, err := env.git.Loaders.TagLoader.GetTags()
	if err != nil {
		return err
	}

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().Tags = tags
	})

	self.refreshView(self.c.Contexts().Tags, env)
	return nil
}

func (self *RefreshHelper) refreshStateSubmoduleConfigs(env refreshEnv) ([]*models.SubmoduleConfig, error) {
	return env.git.Submodule.GetConfigs(nil)
}

// self.refreshStatus is called at the end of this because that's when we can
// be sure there is a State.Model.Branches array to pick the current branch from
func (self *RefreshHelper) refreshBranches(captured capturedBranchState, refreshWorktrees bool, branchSelection types.BranchSelectionBehavior, loadBehindCounts bool, reflogCommits []*models.Commit, env refreshEnv) []*models.Branch {
	loadSeq := self.branchLoadSeq.Add(1)

	branches, err := env.git.Loaders.BranchLoader.Load(
		reflogCommits,
		captured.mainBranches,
		captured.oldBranches,
		loadBehindCounts,
		func(f func() error) {
			self.onWorker(env.background, func(_ gocui.Task) error {
				err := f()
				if err != nil && self.c.State().GetRepoGeneration() != env.generation {
					// An error returned from a worker is shown in a popup. Don't
					// do that if the repo was switched while this worker was in
					// flight: its results are dropped anyway, and the error
					// concerns a repo the user has already left — e.g. failing to
					// compute the behind-counts for a worktree that was deleted
					// after switching away from it.
					self.c.Log.Warnf("dropping error from a stale refresh worker after a repo switch: %v", err)
					return nil
				}
				return err
			})
		},
		func() {
			self.onUIThreadUnlessRepoChanged(env, func() {
				self.c.Contexts().Branches.HandleRender()
				self.refreshStatus(env)
			})
		})
	if err != nil {
		self.c.Log.Error(err)
	}

	var worktrees []*models.Worktree
	if refreshWorktrees {
		worktrees = self.loadWorktrees(env)
	}

	self.onUIThreadUnlessRepoChanged(env, func() {
		// Drop this write if a branch load that started later has already applied
		// its result. At the INITIAL startup stage an immediate load (not
		// recency-sorted) and an async recency-sorted load run concurrently; this
		// makes the later-started (recency-sorted) one win regardless of which
		// finishes first, so its result isn't clobbered by the stale immediate one.
		if loadSeq < self.appliedBranchLoadSeq {
			return
		}
		self.appliedBranchLoadSeq = loadSeq

		// Read the currently-selected branch before overwriting the list, so we
		// can restore it by name below. Reading it here in the bounce keeps it on
		// the UI thread.
		prevSelectedBranch := self.c.Contexts().Branches.GetSelected()

		self.c.Model().Branches = branches
		// Rebuilding here (rather than on the worker) means the map is built from
		// the branches we just wrote, on the UI thread.
		self.rebuildPullRequestsMap()

		if refreshWorktrees {
			self.c.Model().Worktrees = worktrees
			self.refreshView(self.c.Contexts().Worktrees, env)
		}

		// Setting the selection here, in the same bounce that writes the list,
		// keeps it on the UI thread and keeps the list and selection updating in
		// the same frame.
		switch branchSelection {
		case types.KeepBranchSelectionByName:
			if prevSelectedBranch != nil {
				self.searchHelper.ReApplyFilter(self.c.Contexts().Branches)

				_, idx, found := lo.FindIndexOf(self.c.Contexts().Branches.GetItems(),
					func(b *models.Branch) bool { return b.Name == prevSelectedBranch.Name })
				if found {
					self.c.Contexts().Branches.SetSelectedLineIdx(idx)
				}
			}
		case types.SelectCheckedOutBranch:
			// The checked-out branch is always at the top of the list. Setting
			// the selection doesn't scroll the view, so also reset the origin.
			self.c.Contexts().Branches.SetSelectedLineIdx(0)
			self.c.Contexts().Branches.GetView().SetOriginY(0)
		}

		// Need to re-render the commits view because the visualization of local
		// branch heads might have changed
		self.c.Contexts().LocalCommits.HandleRender()
	})

	self.refreshView(self.c.Contexts().Branches, env)

	self.refreshStatus(env)

	// Return the freshly-loaded branches so the caller can hand them to the PR
	// fetch without reading them back from the (bounce-written) model.
	return branches
}

func (self *RefreshHelper) refreshFilesAndSubmodules(captured capturedFilesState, env refreshEnv) error {
	configs, err := self.refreshStateSubmoduleConfigs(env)
	if err != nil {
		return err
	}

	if err := self.refreshStateFiles(captured, env, configs); err != nil {
		return err
	}

	self.refreshView(self.c.Contexts().Submodules, env)
	self.refreshView(self.c.Contexts().Files, env)

	return nil
}

// onUIThreadUnlessRepoChanged bounces a refresh's model/view update onto the UI
// thread, but drops it if the repo was switched while the refresh was in flight.
// Refresh workers do their git work off the UI thread and enqueue their model
// writes here; a repo switch (which replaces the whole model and context tree)
// bumps the generation, so a write captured under the old generation must not
// clobber the new repo's state. The generation is captured once at the start of
// the refresh and carried in env (see refreshEnv).
func (self *RefreshHelper) onUIThreadUnlessRepoChanged(env refreshEnv, f func()) {
	wrapper := func() {
		if self.c.State().GetRepoGeneration() != env.generation {
			return
		}
		f()
	}

	// A batched refresh collects its bounces and fires them together at the end
	// (see refreshBounceBatch); add reports false once the batch is flushing, so
	// bounces enqueued from within a flushed bounce dispatch immediately.
	if env.batch != nil && env.batch.add(wrapper) {
		return
	}

	self.onUIThread(env.background, func() error { wrapper(); return nil })
}

// onWorker and onUIThread pick the foreground or background variant of the
// corresponding dispatch method depending on whether we're servicing a
// background refresh. Background refreshes (auto-fetch and friends) must not
// count towards lazygit being busy, or they'd spuriously block a repo switch;
// see the *Background methods on gocui.Gui.
func (self *RefreshHelper) onWorker(background bool, f func(gocui.Task) error) {
	if background {
		self.c.OnWorkerBackground(f)
	} else {
		self.c.OnWorker(f)
	}
}

func (self *RefreshHelper) onUIThread(background bool, f func() error) {
	if background {
		self.c.OnUIThreadBackground(f)
	} else {
		self.c.OnUIThread(f)
	}
}

// captureOnUIThread runs fn on the UI thread and returns once it has run. fn
// reads the model/context/mode state a refresh scope needs into locals, so the
// worker that follows computes from an immutable snapshot instead of reading
// state the UI thread concurrently mutates. When the enclosing refresh function
// runs on the UI thread (calledFromWorker is false) fn runs inline; when it runs
// on a worker, fn is dispatched to the UI thread and we block for it.
//
// The inline case matters for correctness as much as the hop: OnUIThreadAndWait
// must not be called from the UI thread itself (it would park the thread
// waiting for a callback that only it can run), and capturing inline also
// guarantees the snapshot reflects the state at the moment Refresh was called,
// before the calling handler regains control and can mutate it.
func (self *RefreshHelper) captureOnUIThread(calledFromWorker bool, background bool, fn func()) {
	if !calledFromWorker {
		fn()
		return
	}

	wrapped := func() error {
		fn()
		return nil
	}
	if background {
		_ = self.c.GocuiGui().OnUIThreadAndWaitBackground(wrapped)
	} else {
		_ = self.c.GocuiGui().OnUIThreadAndWait(wrapped)
	}
}

// capturedFilesState holds the files refresh's context/model inputs, gathered
// on the UI thread before the git work runs: the previous files list (to detect
// resolved conflicts and drive the auto-stage), and whether untracked files are
// force-shown.
type capturedFilesState struct {
	prevFiles          []*models.File
	forceShowUntracked bool
}

// captureFilesState reads the files refresh's inputs into an immutable snapshot.
// It must run on the UI thread.
func (self *RefreshHelper) captureFilesState() capturedFilesState {
	return capturedFilesState{
		prevFiles:          self.c.Model().Files,
		forceShowUntracked: self.c.Contexts().Files.ForceShowUntracked(),
	}
}

func (self *RefreshHelper) refreshStateFiles(captured capturedFilesState, env refreshEnv, submoduleConfigs []*models.SubmoduleConfig) error {
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
		for _, file := range captured.prevFiles {
			if file.HasMergeConflicts {
				prevConflictFileCount++
			}
			if file.HasInlineMergeConflicts {
				// Join with the refresh's repo root rather than relying on the
				// process working directory, which may already point at another
				// repo if the user switched while this refresh was in flight.
				hasConflicts, err := mergeconflicts.FileHasConflictMarkers(
					filepath.Join(env.git.RepoPaths.WorktreePath(), file.Path))
				if err != nil {
					self.c.Log.Error(err)
				} else if !hasConflicts {
					pathsToStage = append(pathsToStage, file.Path)
				}
			}
		}

		if len(pathsToStage) > 0 {
			self.c.LogAction(self.c.Tr.Actions.StageResolvedFiles)
			if err := env.git.WorkingTree.StageFiles(pathsToStage, nil); err != nil {
				return err
			}
		}
	}

	files := env.git.Loaders.FileLoader.
		GetStatusFiles(git_commands.GetStatusFileOptions{
			ForceShowUntracked: captured.forceShowUntracked,
			Background:         env.backgroundRoutine,
		})

	conflictFileCount := 0
	for _, file := range files {
		if file.HasMergeConflicts {
			conflictFileCount++
		}
	}

	repoState := self.c.State().GetRepoState()
	workingTreeState := env.git.Status.WorkingTreeState()
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
			self.onUIThreadUnlessRepoChanged(env, func() {
				// The merge-conflicts scope of this refresh also notices that
				// the conflicts are gone and escapes from the merge conflicts
				// view to the files context (see RefreshMergeState), but it
				// runs concurrently with us, and its escape refuses to push
				// the files context over a popup. So if our prompt opens
				// first, the escape does nothing, and closing the prompt
				// would land the user in the dead merge conflicts view.
				// Escape it ourselves before opening the prompt, so that the
				// prompt always opens on top of the files context.
				if self.c.Context().IsCurrent(self.c.Contexts().MergeConflicts) {
					self.mergeConflictsHelper.ResetMergeState()
					self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
				}
				self.mergeAndRebaseHelper.PromptToContinueRebase()
			})
		}
	} else {
		// Either there's no operation in progress any more, or new conflicts have
		// appeared. Either way, a "continue?" prompt we're showing is now stale
		// (e.g. the operation was continued or aborted outside lazygit), so
		// dismiss it rather than leave the user with a prompt that would fail.
		// Guard on the generation like the sibling PromptToContinueRebase
		// bounce above: if the repo was switched while this refresh was in
		// flight, a prompt showing now belongs to the new repo, so leave it be.
		self.onUIThreadUnlessRepoChanged(env, func() {
			self.mergeAndRebaseHelper.DismissContinueRebasePromptIfShowing()
		})
	}

	self.onUIThreadUnlessRepoChanged(env, func() {
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
func (self *RefreshHelper) refreshReflogCommits(captured capturedReflogState, env refreshEnv, selectTopEntry bool) ([]*models.Commit, error) {
	// load does the git work on the worker and returns the new value for a
	// reflog slice, reading the existing slice (captured on the UI thread) for
	// the incremental fetch. The caller writes the result in the bounce.
	load := func(existing []*models.Commit, filterPath string, filterAuthor string) ([]*models.Commit, error) {
		var lastReflogCommit *models.Commit
		if filterPath == "" && filterAuthor == "" && len(existing) > 0 {
			lastReflogCommit = existing[0]
		}

		commits, onlyObtainedNewReflogCommits, err := env.git.Loaders.ReflogCommitLoader.
			GetReflogCommits(captured.hashPool, lastReflogCommit, filterPath, filterAuthor)
		if err != nil {
			return nil, err
		}

		if onlyObtainedNewReflogCommits {
			return append(commits, existing...), nil
		}
		return commits, nil
	}

	reflogCommits, err := load(captured.reflogCommits, "", "")
	if err != nil {
		return nil, err
	}

	filteredReflogCommits := reflogCommits
	if captured.filteringActive {
		filteredReflogCommits, err = load(captured.filteredReflogCommits, captured.filterPath, captured.filterAuthor)
		if err != nil {
			return nil, err
		}
	}

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().ReflogCommits = reflogCommits
		self.c.Model().FilteredReflogCommits = filteredReflogCommits
		// Setting the selection here, in the same bounce that writes the list,
		// keeps it on the UI thread and atomic with the list update. Setting the
		// selection doesn't scroll the view, so also reset the origin.
		if selectTopEntry {
			self.c.Contexts().ReflogCommits.SetSelectedLineIdx(0)
			self.c.Contexts().ReflogCommits.GetView().SetOriginY(0)
		}
	})

	self.refreshView(self.c.Contexts().ReflogCommits, env)
	return reflogCommits, nil
}

func (self *RefreshHelper) refreshRemotes(prevSelectedRemote *models.Remote, env refreshEnv) ([]*models.Remote, error) {
	remotes, err := env.git.Loaders.RemoteLoader.GetRemotes()
	if err != nil {
		return nil, err
	}

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().Remotes = remotes

		hadPrs := len(self.c.Model().PullRequestsMap) != 0
		self.rebuildPullRequestsMap()
		if !hadPrs && len(self.c.Model().PullRequestsMap) != 0 {
			// if we didn't have PRs in the map before but now we do, we need to redraw the branches view
			self.refreshView(self.c.Contexts().Branches, env)
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
	})

	self.refreshView(self.c.Contexts().Remotes, env)
	self.refreshView(self.c.Contexts().RemoteBranches, env)
	return remotes, nil
}

func (self *RefreshHelper) loadWorktrees(env refreshEnv) []*models.Worktree {
	worktrees, err := env.git.Loaders.Worktrees.GetWorktrees()
	if err != nil {
		self.c.Log.Error(err)
		return []*models.Worktree{}
	}
	return worktrees
}

func (self *RefreshHelper) refreshWorktrees(env refreshEnv) {
	worktrees := self.loadWorktrees(env)

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().Worktrees = worktrees
	})

	// need to refresh branches because the branches view shows worktrees against
	// branches
	self.refreshView(self.c.Contexts().Branches, env)
	self.refreshView(self.c.Contexts().Worktrees, env)
}

func (self *RefreshHelper) refreshStashEntries(filterPath string, env refreshEnv) {
	stashEntries := env.git.Loaders.StashLoader.
		GetStashEntries(filterPath)

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().StashEntries = stashEntries
	})

	self.refreshView(self.c.Contexts().Stash, env)
}

// never call this on its own, it should only be called from within refreshCommits()
func (self *RefreshHelper) refreshStatus(env refreshEnv) {
	workingTreeState := env.git.Status.WorkingTreeState()
	repoName := env.git.RepoPaths.RepoName()

	self.onUIThreadUnlessRepoChanged(env, func() {
		// Read the checked-out branch and the linked worktree name here on the UI
		// thread: both derive from models (Branches, Worktrees) that their
		// refreshes now write via bounces, so reading them on the worker would
		// see stale values from before those bounces applied.
		currentBranch := self.refsHelper.GetCheckedOutRef()
		if currentBranch == nil {
			// need to wait for branches to refresh
			return
		}
		linkedWorktreeName := self.worktreeHelper.GetLinkedWorktreeName()

		status := presentation.FormatStatus(repoName, currentBranch, types.ItemOperationNone, linkedWorktreeName, workingTreeState, self.c.Tr, self.c.UserConfig())
		self.c.SetViewContent(self.c.Views().Status, status)
	})
}

// refForLog returns the ref to log commits from, along with the bisect info it
// read to decide that. The caller writes the bisect info to the model (in its
// bounce) rather than refForLog doing it, so the model write stays on the UI
// thread.
func (self *RefreshHelper) refForLog(env refreshEnv) (string, *git_commands.BisectInfo) {
	bisectInfo := env.git.Bisect.GetInfo()

	if !bisectInfo.Started() {
		return "HEAD", bisectInfo
	}

	// need to see if our bisect's current commit is reachable from our 'new' ref.
	if bisectInfo.Bisecting() && !env.git.Bisect.ReachableFromStart(bisectInfo) {
		return bisectInfo.GetNewHash(), bisectInfo
	}

	return bisectInfo.GetStartHash(), bisectInfo
}

func (self *RefreshHelper) refreshView(context types.Context, env refreshEnv) {
	// refreshView is called from the worker goroutine that drives async
	// refreshes, so bounce to the UI thread before mutating view content. Guard
	// on the generation like the model-update bounces do: if the repo was
	// switched while the refresh was in flight, its model write was already
	// dropped, so there's nothing fresh to render — and the captured context
	// belongs to the old repo's now-replaced context tree anyway.
	self.onUIThreadUnlessRepoChanged(env, func() {
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
	})
}

func (self *RefreshHelper) refreshGithubPullRequests(branches []*models.Branch, remotes []*models.Remote, env refreshEnv) {
	clearPullRequests := func() {
		self.onUIThreadUnlessRepoChanged(env, func() {
			self.c.Model().PullRequests = nil
			self.c.Model().PullRequestsMap = nil
		})
	}

	githubRemotes := getAuthenticatedGithubRemotes(self.getGithubRemotes(remotes, env), env.git.GitHub.GetAuthToken)
	if len(githubRemotes) == 0 {
		clearPullRequests()
		return
	}

	baseInfo := getGithubBaseRemote(githubRemotes, env.git.GitHub.ConfiguredBaseRemoteName())
	if baseInfo == nil {
		clearPullRequests()

		if !self.githubBaseRemotePromptDismissed[env.git.RepoPaths.RepoPath()] {
			self.promptForBaseGithubRepo(githubRemotes)
		}
		return
	}

	self.setGithubPullRequests(baseInfo, branches, env)
}

type githubRemoteInfo struct {
	remote      *models.Remote
	serviceInfo hosting_service.ServiceInfo
	authToken   string
}

func (self *RefreshHelper) getGithubRemotes(remotes []*models.Remote, env refreshEnv) []githubRemoteInfo {
	return lo.FilterMap(remotes, func(remote *models.Remote, _ int) (githubRemoteInfo, bool) {
		if len(remote.Urls) == 0 {
			return githubRemoteInfo{}, false
		}
		serviceInfo, err := env.git.HostingService.GetServiceInfo(remote.Urls[0])
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

					self.c.RefreshFromWorker(types.RefreshOptions{Scope: []types.RefreshableView{types.PULL_REQUESTS}})
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

func (self *RefreshHelper) setGithubPullRequests(baseInfo *githubRemoteInfo, branches []*models.Branch, env refreshEnv) {
	if len(branches) == 0 {
		return
	}

	trackingBranches := lo.Filter(branches, func(branch *models.Branch, _ int) bool {
		return branch.IsTrackingRemote()
	})
	branchNames := lo.Map(trackingBranches, func(branch *models.Branch, _ int) string {
		return branch.UpstreamBranch
	})

	prs, err := env.git.GitHub.FetchRecentPRs(branchNames, &baseInfo.serviceInfo, baseInfo.authToken)
	if err != nil {
		self.c.Log.Error("error fetching pull requests from GitHub: " + err.Error())
		return
	}

	self.savePullRequestsToCache(prs, env)

	self.onUIThreadUnlessRepoChanged(env, func() {
		self.c.Model().PullRequests = prs
		// Rebuilding here rather than on the worker means the map is built from
		// the branches and remotes as they are on the UI thread, after their
		// own refreshes' bounces have applied.
		self.rebuildPullRequestsMap()
		self.c.PostRefreshUpdate(self.c.Contexts().Branches)
	})
}

func (self *RefreshHelper) savePullRequestsToCache(prs []*models.GithubPullRequest, env refreshEnv) {
	// Key the cache by the repo the refresh was started for, not the live one:
	// this runs on a worker, and if the user switched repos while the fetch was
	// in flight, the live instance would file the old repo's pull requests
	// under the new repo's path.
	repoPath := env.git.RepoPaths.RepoPath()
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
