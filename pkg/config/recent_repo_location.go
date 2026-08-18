package config

// RecentRepoLocation is one entry of AppState.RecentRepoLocations: the
// directory to change back to, plus the GIT_DIR/GIT_WORK_TREE environment a
// repo opened with --git-dir/--work-tree (or core.worktree) needs to be
// found again. GitLocationEnvVars is empty for every repo git can find from
// the work tree alone (#5942).
type RecentRepoLocation struct {
	Path               string
	GitLocationEnvVars []string
}
