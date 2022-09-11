package models

import (
	"path/filepath"
)

// Worktree : A git worktree
type Worktree struct {
	Id     int
	Path   string
	Branch string
}

func (w *Worktree) RefName() string {
	return w.Name()
}

func (w *Worktree) ID() string {
	return w.RefName()
}

func (w *Worktree) Description() string {
	return w.RefName()
}

func (w *Worktree) Name() string {
	return filepath.Base(w.Path)
}

func (w *Worktree) Main() bool {
	return w.Id == 0
}
