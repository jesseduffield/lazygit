package controllers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
)

// Both models.File and models.CommitFile satisfy this. Names returns the file's
// path, plus the path it was renamed from if it is a rename.
type fileWithNames[T any] interface {
	*T
	GetPath() string
	GetPreviousPath() string
	Names() []string
}

// pathsForDiff returns the paths to limit a diff command to for showing the
// changes of the given node. root is the root of the tree that the node belongs
// to, and isFiltering says whether that tree is reduced to the files matching a
// text filter.
func pathsForDiff[T any, PT fileWithNames[T]](node *filetree.Node[T], root *filetree.Node[T], isFiltering bool) []string {
	if file := node.GetFile(); file != nil {
		return PT(file).Names()
	}

	dir := node.GetPath()

	if isFiltering {
		// Passing the directory would bring back the files that the filter hides,
		// so we spell out the ones it leaves.
		var paths []string
		forEachFileInDir[T, PT](root, dir, func(file PT) {
			paths = append(paths, file.Names()...)
		})
		return paths
	}

	// The directory covers everything below it, but git only pairs up the two
	// ends of a rename if both are in the pathspec, and one end can well be
	// outside the directory. Without that end we would get an addition or a
	// deletion where the diff has a rename.
	paths := []string{dir}
	forEachFileInDir[T, PT](root, dir, func(file PT) {
		if path := file.GetPath(); !isInDir(path, dir) {
			paths = append(paths, path)
		}
		if previousPath := file.GetPreviousPath(); previousPath != "" && !isInDir(previousPath, dir) {
			paths = append(paths, previousPath)
		}
	})
	return paths
}

// forEachFileInDir calls cb for each file in the tree that the given directory
// contains, either at its current or at its previous path.
func forEachFileInDir[T any, PT fileWithNames[T]](root *filetree.Node[T], dir string, cb func(PT)) {
	_ = root.ForEachFile(func(f *T) error {
		file := PT(f)
		previousPath := file.GetPreviousPath()
		if isInDir(file.GetPath(), dir) || (previousPath != "" && isInDir(previousPath, dir)) {
			cb(file)
		}
		return nil
	})
}

func isInDir(path string, dir string) bool {
	// "." is the root item, which contains every file
	return dir == "." || strings.HasPrefix(path, dir+"/")
}
