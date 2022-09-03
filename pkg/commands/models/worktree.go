package models

import (
	"fmt"
	"github.com/go-errors/errors"
	"io/fs"
	"log"
	"os"
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

func (w *Worktree) Current() bool {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return pwd == w.Path
}

func (w *Worktree) Missing() bool {
	if _, err := os.Stat(w.Path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return true
		}
		log.Fatalln(fmt.Errorf("failed to check if worktree path `%s` exists\n%w", w.Path, err).Error())
	}
	return false
}
