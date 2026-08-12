package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
)

// Both models.File and models.CommitFile satisfy this. Names returns the file's
// path, plus the path it was renamed from if it is a rename.
type fileWithNames[T any] interface {
	*T
	Names() []string
}

// pathsForDiff returns the paths to limit a diff command to for showing the
// changes of the given node. For a directory this is the directory itself,
// unless a filter is active, in which case we list the files that are visible
// under it.
func pathsForDiff[T any, PT fileWithNames[T]](node *filetree.Node[T], isFiltering bool) []string {
	if file := node.GetFile(); file != nil {
		return PT(file).Names()
	}

	if isFiltering {
		var paths []string
		_ = node.ForEachFile(func(file *T) error {
			// For a rename we need to pass both paths so that git detects it as
			// a rename rather than an unrelated delete and add.
			paths = append(paths, PT(file).Names()...)
			return nil
		})
		return paths
	}

	return []string{node.GetPath()}
}
