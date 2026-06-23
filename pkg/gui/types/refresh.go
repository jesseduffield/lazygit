package types

// models/views that we can refresh
type RefreshableView int

const (
	COMMITS RefreshableView = iota
	REBASE_COMMITS
	SUB_COMMITS
	BRANCHES
	FILES
	STASH
	REFLOG
	TAGS
	REMOTES
	WORKTREES
	STATUS
	SUBMODULES
	STAGING
	PATCH_BUILDING
	MERGE_CONFLICTS
	COMMIT_FILES
	// not actually views. Will refactor this later
	BISECT_INFO
	PULL_REQUESTS
)

type RefreshMode int

const (
	SYNC     RefreshMode = iota // wait until everything is done before returning
	ASYNC                       // return immediately, allowing each independent thing to update itself
	BLOCK_UI                    // wrap code in an update call to ensure UI updates all at once and keybindings aren't executed till complete
)

// CommitSelectionBehavior controls which local commit is selected after the
// commits list is reloaded by a refresh.
type CommitSelectionBehavior int

const (
	// Keep the same commit selected by hash (and the same range, when
	// range-selecting), restoring it at its new position if it moved. This is
	// the right default whenever the list reloads underneath a selection the
	// user hasn't deliberately changed.
	KeepCommitSelectionByHash CommitSelectionBehavior = iota

	// Leave the selection index untouched, because the caller set it itself
	// before refreshing. Used when jumping to the top of the list after a
	// checkout, and when following a commit that was just moved up or down.
	KeepCommitSelectionIndex

	// Select the HEAD commit. Used by operations that create a new commit at
	// HEAD (committing, merging, pulling with a merge); the by-hash behavior
	// can't restore a commit that didn't exist before the refresh.
	SelectHeadCommit
)

type RefreshOptions struct {
	Then  func()
	Scope []RefreshableView // e.g. []RefreshableView{COMMITS, BRANCHES}. Leave empty to refresh everything
	Mode  RefreshMode       // one of SYNC (default), ASYNC, and BLOCK_UI

	// Normally a refresh of the branches tries to keep the same branch selected
	// (by name); this is usually important in case the order of branches
	// changes. Passing true for KeepBranchSelectionIndex suppresses this and
	// keeps the selection index the same. Useful after checking out a detached
	// head, and selecting index 0.
	KeepBranchSelectionIndex bool

	// Controls which local commit is selected after the refresh. Defaults to
	// KeepCommitSelectionByHash.
	CommitSelection CommitSelectionBehavior

	// When true, this refresh was initiated by a background routine rather than
	// by a user action. We use it to keep background `git status` calls from
	// taking optional git locks, so they don't contend for index.lock with git
	// commands the user runs in a terminal. The cost is that such a status won't
	// persist git's refreshed stat-cache, which is the right trade-off for
	// unattended work; foreground refreshes leave this false so they do persist.
	Background bool
}
