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

// BranchSelectionBehavior controls which local branch is selected after the
// branches list is reloaded by a refresh.
type BranchSelectionBehavior int

const (
	// Keep the same branch selected by name, restoring it at its new position if
	// the order changed. This is the right default whenever the list reloads
	// underneath a selection the user hasn't deliberately changed.
	KeepBranchSelectionByName BranchSelectionBehavior = iota

	// Select the checked-out branch (the one at the top of the list). Used after
	// operations that check something out - checkout, creating a branch, moving
	// commits to a new branch - so the newly checked-out ref ends up selected.
	SelectCheckedOutBranch
)

type RefreshOptions struct {
	Then  func() error
	Scope []RefreshableView // e.g. []RefreshableView{COMMITS, BRANCHES}. Leave empty to refresh everything

	// If true, hold off on updating the UI until all scopes have finished
	// refreshing and then apply them together in a single frame, rather than
	// letting each scope update the UI as soon as it's done.
	BatchUIUpdates bool

	// Controls which local branch is selected after the refresh. Defaults to
	// KeepBranchSelectionByName.
	BranchSelection BranchSelectionBehavior

	// Controls which local commit is selected after the refresh. Defaults to
	// KeepCommitSelectionByHash.
	CommitSelection CommitSelectionBehavior

	// When true, select the top (most recent) reflog entry after the refresh.
	// Used alongside SelectCheckedOutBranch by operations that check something
	// out, since the checkout adds a new reflog entry at the top. Defaults to
	// keeping the reflog selection where it is.
	SelectTopReflogCommit bool

	// When true, this refresh was initiated by a background routine rather than
	// by a user action. Every git command suppresses optional locks by default
	// so it can't contend for index.lock (see git_commands.OptionalLocksEnvVar);
	// a foreground files refresh (this false) is the one command that opts back
	// in, so it persists git's refreshed stat-cache and keeps later status calls
	// fast. Background refreshes leave the suppression in place: not persisting
	// the stat-cache is the right trade-off for unattended work.
	Background bool
}
