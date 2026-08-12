package controllers

import (
	"path"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/samber/lo"
)

// Both models.File and models.CommitFile satisfy this. Names returns the file's
// path, plus the path it was renamed from if it is a rename.
type fileWithNames[T any] interface {
	*T
	GetPath() string
	GetPreviousPath() string
	Names() []string
}

// diffPathsForNode returns the paths to limit a diff command to for showing the
// changes of the given node. files are all the files that the diff contains,
// while root is the root of the tree the node belongs to, which holds only the
// files matching the text filter when there is one.
func diffPathsForNode[T any, PT fileWithNames[T]](node *filetree.Node[T], root *filetree.Node[T], files []*T, isFiltering bool) []string {
	if file := node.GetFile(); file != nil {
		return PT(file).Names()
	}

	dir := node.GetPath()

	if isFiltering {
		// Passing the directory would bring back the files that the filter hides,
		// so we spell out the ones it leaves.
		var paths []string
		for _, file := range filesInDir[T, PT](filesInTree(root), dir) {
			paths = append(paths, PT(file).Names()...)
		}
		return paths
	}

	// The directory covers everything below it, but git only pairs up the two
	// ends of a rename if both are in the pathspec, and one end can well be
	// outside the directory. Without that end we would get an addition or a
	// deletion where the diff has a rename.
	var outsidePaths []string
	for _, f := range filesInDir[T, PT](files, dir) {
		file := PT(f)
		if p := file.GetPath(); !isInDir(p, dir) {
			outsidePaths = append(outsidePaths, p)
		}
		if p := file.GetPreviousPath(); p != "" && !isInDir(p, dir) {
			outsidePaths = append(outsidePaths, p)
		}
	}

	return dropContainedPaths(append([]string{dir}, collapseToDirs[T, PT](outsidePaths, files, dir)...))
}

// dropContainedPaths removes the paths that another one of them contains, since
// a pathspec that matches a directory matches everything below it anyway.
func dropContainedPaths(paths []string) []string {
	return lo.Filter(paths, func(p string, _ int) bool {
		return !lo.SomeBy(paths, func(other string) bool {
			return other != p && isInDir(p, other)
		})
	})
}

// collapseToDirs replaces each of the given paths with the highest directory
// that can stand in for it, so that moving a whole directory elsewhere costs a
// single pathspec rather than one per file. There is a limit to how long a
// command line may get, and a commit can move a great many files at once.
func collapseToDirs[T any, PT fileWithNames[T]](paths []string, files []*T, dir string) []string {
	if len(paths) == 0 {
		return nil
	}

	// A directory can stand in for the paths under it as long as everything it
	// contains ends up in the diff anyway, which is to say as long as all of it
	// is in the directory we are diffing too.
	canStandIn := make(map[string]bool)
	standsIn := func(candidate string) bool {
		if result, ok := canStandIn[candidate]; ok {
			return result
		}

		result := lo.EveryBy(files, func(file *T) bool {
			return !fileIsInDir[T, PT](file, candidate) || fileIsInDir[T, PT](file, dir)
		})
		canStandIn[candidate] = result
		return result
	}

	return lo.Uniq(lo.Map(paths, func(p string, _ int) string {
		// A directory that can't stand in for the path rules out its parents
		// too, since they contain everything it contains. We stop short of the
		// repository root: it would leave the command with nothing to say about
		// the directory whose diff we are showing.
		for candidate := path.Dir(p); candidate != "." && standsIn(candidate); candidate = path.Dir(candidate) {
			p = candidate
		}
		return p
	}))
}

func filesInTree[T any](root *filetree.Node[T]) []*T {
	files := []*T{}
	_ = root.ForEachFile(func(file *T) error {
		files = append(files, file)
		return nil
	})
	return files
}

// filesInDir returns the files that the given directory contains, either at
// their current or at their previous path.
func filesInDir[T any, PT fileWithNames[T]](files []*T, dir string) []*T {
	return lo.Filter(files, func(file *T, _ int) bool {
		return fileIsInDir[T, PT](file, dir)
	})
}

func fileIsInDir[T any, PT fileWithNames[T]](f *T, dir string) bool {
	file := PT(f)
	previousPath := file.GetPreviousPath()
	return isInDir(file.GetPath(), dir) || (previousPath != "" && isInDir(previousPath, dir))
}

func isInDir(path string, dir string) bool {
	// "." is the root item, which contains every file
	return dir == "." || strings.HasPrefix(path, dir+"/")
}
