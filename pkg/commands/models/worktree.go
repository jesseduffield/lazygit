package models

// A git worktree
type Worktree struct {
	// if false, this is a linked worktree
	IsMain bool
	Path   string
	Branch string
	// based on the path, but uniquified
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
