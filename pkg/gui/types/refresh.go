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
	// not actually a view. Will refactor this later
	BISECT_INFO
)

type RefreshMode int

const (
	SYNC     RefreshMode = iota // wait until everything is done before returning
	ASYNC                       // return immediately, allowing each independent thing to update itself
	BLOCK_UI                    // wrap code in an update call to ensure UI updates all at once and keybindings aren't executed till complete
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
}
