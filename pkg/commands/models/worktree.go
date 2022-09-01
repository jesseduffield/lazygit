package models

// Worktree : A git worktree
type Worktree struct {
	Name    string
	Main    bool
	Current bool
	Path    string
}

func (w *Worktree) RefName() string {
	return w.Name
}

func (w *Worktree) ID() string {
	return w.RefName()
}

func (w *Worktree) Description() string {
	return w.RefName()
}
