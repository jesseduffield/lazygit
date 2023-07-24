package models

// A git worktree
type Worktree struct {
	// if false, this is a linked worktree
	IsMain bool
	// path to the directory of the worktree i.e. the directory that contains all the user's files
	Path string
	// path of the git directory for this worktree. The equivalent of the .git directory
	// in the main worktree. For linked worktrees this would be <repo_path>/.git/worktrees/<name>
	GitDir string
	// If the worktree has a branch checked out, this field will be set to the branch name.
	// A branch is considered 'checked out' if:
	// * the worktree is directly on the branch
	// * the worktree is mid-rebase on the branch
	// * the worktree is mid-bisect on the branch
	Branch string
	// based on the path, but uniquified. Not the same name that git uses in the worktrees/ folder (no good reason for this,
	// I just prefer my naming convention better)
	NameField string
}

func (w *Worktree) RefName() string {
	return w.Name()
}

func (w *Worktree) ID() string {
	return w.Path
}

func (w *Worktree) Description() string {
	return w.RefName()
}

func (w *Worktree) Name() string {
	return w.NameField
}

func (w *Worktree) Main() bool {
	return w.IsMain
}
