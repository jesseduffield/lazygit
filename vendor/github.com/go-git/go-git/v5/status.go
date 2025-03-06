package git

import (
	"bytes"
	"fmt"
	"path/filepath"

	mindex "github.com/go-git/go-git/v5/utils/merkletrie/index"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
)

// Status represents the current status of a Worktree.
// The key of the map is the path of the file.
type Status map[string]*FileStatus

// File returns the FileStatus for a given path, if the FileStatus doesn't
// exists a new FileStatus is added to the map using the path as key.
func (s Status) File(path string) *FileStatus {
	if _, ok := (s)[path]; !ok {
		s[path] = &FileStatus{Worktree: Untracked, Staging: Untracked}
	}

	return s[path]
}

// IsUntracked checks if file for given path is 'Untracked'
func (s Status) IsUntracked(path string) bool {
	stat, ok := (s)[filepath.ToSlash(path)]
	return ok && stat.Worktree == Untracked
}

// IsClean returns true if all the files are in Unmodified status.
func (s Status) IsClean() bool {
	for _, status := range s {
		if status.Worktree != Unmodified || status.Staging != Unmodified {
			return false
		}
	}

	return true
}

func (s Status) String() string {
	buf := bytes.NewBuffer(nil)
	for path, status := range s {
		if status.Staging == Unmodified && status.Worktree == Unmodified {
			continue
		}

		if status.Staging == Renamed {
			path = fmt.Sprintf("%s -> %s", path, status.Extra)
		}

		fmt.Fprintf(buf, "%c%c %s\n", status.Staging, status.Worktree, path)
	}

	return buf.String()
}

// FileStatus contains the status of a file in the worktree
type FileStatus struct {
	// Staging is the status of a file in the staging area
	Staging StatusCode
	// Worktree is the status of a file in the worktree
	Worktree StatusCode
	// Extra contains extra information, such as the previous name in a rename
	Extra string
}

// StatusCode status code of a file in the Worktree
type StatusCode byte

const (
	Unmodified         StatusCode = ' '
	Untracked          StatusCode = '?'
	Modified           StatusCode = 'M'
	Added              StatusCode = 'A'
	Deleted            StatusCode = 'D'
	Renamed            StatusCode = 'R'
	Copied             StatusCode = 'C'
	UpdatedButUnmerged StatusCode = 'U'
)

// StatusStrategy defines the different types of strategies when processing
// the worktree status.
type StatusStrategy int

const (
	// TODO: (V6) Review the default status strategy.
	// TODO: (V6) Review the type used to represent Status, to enable lazy
	// processing of statuses going direct to the backing filesystem.
	defaultStatusStrategy = Empty

	// Empty starts its status map from empty. Missing entries for a given
	// path means that the file is untracked. This causes a known issue (#119)
	// whereby unmodified files can be incorrectly reported as untracked.
	//
	// This can be used when returning the changed state within a modified Worktree.
	// For example, to check whether the current worktree is clean.
	Empty StatusStrategy = 0
	// Preload goes through all existing nodes from the index and add them to the
	// status map as unmodified. This is currently the most reliable strategy
	// although it comes at a performance cost in large repositories.
	//
	// This method is recommended when fetching the status of unmodified files.
	// For example, to confirm the status of a specific file that is either
	// untracked or unmodified.
	Preload StatusStrategy = 1
)

func (s StatusStrategy) new(w *Worktree) (Status, error) {
	switch s {
	case Preload:
		return preloadStatus(w)
	case Empty:
		return make(Status), nil
	}
	return nil, fmt.Errorf("%w: %+v", ErrUnsupportedStatusStrategy, s)
}

func preloadStatus(w *Worktree) (Status, error) {
	idx, err := w.r.Storer.Index()
	if err != nil {
		return nil, err
	}

	idxRoot := mindex.NewRootNode(idx)
	nodes := []noder.Noder{idxRoot}

	status := make(Status)
	for len(nodes) > 0 {
		var node noder.Noder
		node, nodes = nodes[0], nodes[1:]
		if node.IsDir() {
			children, err := node.Children()
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, children...)
			continue
		}
		fs := status.File(node.Name())
		fs.Worktree = Unmodified
		fs.Staging = Unmodified
	}

	return status, nil
}
