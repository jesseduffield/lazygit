package filetree

import (
	"path"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// Represents a file or directory in a file tree.
type Node[T any] struct {
	// File will be nil if the node is a directory.
	File *T

	// If the node is a directory, Children contains the contents of the directory,
	// otherwise it's nil.
	Children []*Node[T]

	// path of the file/directory
	// private; use either GetPath() or GetInternalPath() to access
	path string

	// rather than render a tree as:
	// a/
	//   b/
	//     file.blah
	//
	// we instead render it as:
	// a/b/
	//	 file.blah
	// This saves vertical space. The CompressionLevel of a node is equal to the
	// number of times a 'compression' like the above has happened, where two
	// nodes are squished into one.
	CompressionLevel int
}

var _ types.ListItem = &Node[models.File]{}

func (self *Node[T]) IsFile() bool {
	return self.File != nil
}

func (self *Node[T]) GetFile() *T {
	return self.File
}

// This returns the logical path from the user's point of view. It is the
// relative path from the root of the repository.
// Use this for display, or when you want to perform some action on the path
// (e.g. a git command).
func (self *Node[T]) GetPath() string {
	return strings.TrimPrefix(self.path, "./")
}

// This returns the internal path from the tree's point of view. It's the same
// as GetPath(), but prefixed with "./" for the root item.
// Use this when interacting with the tree itself, e.g. when calling
// ToggleCollapsed.
func (self *Node[T]) GetInternalPath() string {
	return self.path
}

func (self *Node[T]) Sort(cmp func(a, b *Node[T]) int) {
	self.SortChildren(cmp)

	for _, child := range self.Children {
		child.Sort(cmp)
	}
}

// NodeSortComparator returns a comparator function for sorting tree nodes
// based on the given sort order and case sensitivity.
// sortOrder must be one of: "mixed", "filesFirst", "foldersFirst".
func NodeSortComparator[T any](sortOrder string, caseSensitive bool) func(a, b *Node[T]) int {
	strCmp := strings.Compare
	if !caseSensitive {
		strCmp = func(a, b string) int {
			return strings.Compare(strings.ToLower(a), strings.ToLower(b))
		}
	}

	// dirVsFileOrder is the return value when a is a directory and b is a file.
	// -1 means directories come first, 1 means files come first.
	dirVsFileOrder := 0
	switch sortOrder {
	case "foldersFirst":
		dirVsFileOrder = -1
	case "filesFirst":
		dirVsFileOrder = 1
	}

	if dirVsFileOrder != 0 {
		return func(a, b *Node[T]) int {
			aIsDir := !a.IsFile()
			bIsDir := !b.IsFile()
			if aIsDir != bIsDir {
				if aIsDir {
					return dirVsFileOrder
				}
				return -dirVsFileOrder
			}
			return strCmp(a.path, b.path)
		}
	}

	// "mixed": sort by path only
	return func(a, b *Node[T]) int {
		return strCmp(a.path, b.path)
	}
}

func (self *Node[T]) ForEachFile(cb func(*T) error) error {
	if self.IsFile() {
		if err := cb(self.File); err != nil {
			return err
		}
	}

	for _, child := range self.Children {
		if err := child.ForEachFile(cb); err != nil {
			return err
		}
	}

	return nil
}

func (self *Node[T]) SortChildren(cmp func(a, b *Node[T]) int) {
	if self.IsFile() {
		return
	}

	children := slices.Clone(self.Children)

	slices.SortFunc(children, cmp)

	// TODO: think about making this in-place
	self.Children = children
}

func (self *Node[T]) Some(predicate func(*Node[T]) bool) bool {
	if predicate(self) {
		return true
	}

	for _, child := range self.Children {
		if child.Some(predicate) {
			return true
		}
	}

	return false
}

func (self *Node[T]) SomeFile(predicate func(*T) bool) bool {
	if self.IsFile() {
		if predicate(self.File) {
			return true
		}
	} else {
		for _, child := range self.Children {
			if child.SomeFile(predicate) {
				return true
			}
		}
	}

	return false
}

func (self *Node[T]) Every(predicate func(*Node[T]) bool) bool {
	if !predicate(self) {
		return false
	}

	for _, child := range self.Children {
		if !child.Every(predicate) {
			return false
		}
	}

	return true
}

func (self *Node[T]) EveryFile(predicate func(*T) bool) bool {
	if self.IsFile() {
		if !predicate(self.File) {
			return false
		}
	} else {
		for _, child := range self.Children {
			if !child.EveryFile(predicate) {
				return false
			}
		}
	}

	return true
}

func (self *Node[T]) FindFirstFileBy(predicate func(*T) bool) *T {
	if self.IsFile() {
		if predicate(self.File) {
			return self.File
		}
	} else {
		for _, child := range self.Children {
			if file := child.FindFirstFileBy(predicate); file != nil {
				return file
			}
		}
	}

	return nil
}

func (self *Node[T]) Flatten(collapsedPaths *CollapsedPaths) []*Node[T] {
	result := []*Node[T]{self}

	if len(self.Children) > 0 && !collapsedPaths.IsCollapsed(self.path) {
		result = append(result, lo.FlatMap(self.Children, func(child *Node[T], _ int) []*Node[T] {
			return child.Flatten(collapsedPaths)
		})...)
	}

	return result
}

func (self *Node[T]) GetNodeAtIndex(index int, collapsedPaths *CollapsedPaths) *Node[T] {
	if self == nil {
		return nil
	}

	node, _, _ := self.getNodeAtIndexAux(index, collapsedPaths, -1)

	return node
}

// GetVisualDepthAtIndex returns the visual depth (indentation level) of the
// node at the given flat index. Visual depth differs from tree depth because
// compressed nodes (e.g. "a/b/") count as a single visual level.
// Returns -1 if the index is out of range.
func (self *Node[T]) GetVisualDepthAtIndex(index int, collapsedPaths *CollapsedPaths) int {
	if self == nil {
		return -1
	}

	_, _, depth := self.getNodeAtIndexAux(index, collapsedPaths, -1)

	return depth
}

func (self *Node[T]) getNodeAtIndexAux(index int, collapsedPaths *CollapsedPaths, visualDepth int) (*Node[T], int, int) {
	offset := 1

	if index == 0 {
		return self, offset, visualDepth
	}

	if !collapsedPaths.IsCollapsed(self.path) {
		for _, child := range self.Children {
			foundNode, offsetChange, depth := child.getNodeAtIndexAux(index-offset, collapsedPaths, visualDepth+1)
			offset += offsetChange
			if foundNode != nil {
				return foundNode, offset, depth
			}
		}
	}

	return nil, offset, -1
}

func (self *Node[T]) GetIndexForPath(path string, collapsedPaths *CollapsedPaths) (int, bool) {
	offset := 0

	if self.path == path {
		return offset, true
	}

	if !collapsedPaths.IsCollapsed(self.path) {
		for _, child := range self.Children {
			offsetChange, found := child.GetIndexForPath(path, collapsedPaths)
			offset += offsetChange + 1
			if found {
				return offset, true
			}
		}
	}

	return offset, false
}

func (self *Node[T]) Size(collapsedPaths *CollapsedPaths) int {
	if self == nil {
		return 0
	}

	output := 1

	if !collapsedPaths.IsCollapsed(self.path) {
		for _, child := range self.Children {
			output += child.Size(collapsedPaths)
		}
	}

	return output
}

func (self *Node[T]) Compress() {
	if self == nil {
		return
	}

	self.compressAux()
}

func (self *Node[T]) compressAux() *Node[T] {
	if self.IsFile() {
		return self
	}

	children := self.Children
	for i := range children {
		grandchildren := children[i].Children
		for len(grandchildren) == 1 && !grandchildren[0].IsFile() {
			grandchildren[0].CompressionLevel = children[i].CompressionLevel + 1
			children[i] = grandchildren[0]
			grandchildren = children[i].Children
		}
	}

	for i := range children {
		children[i] = children[i].compressAux()
	}

	self.Children = children

	return self
}

func (self *Node[T]) GetPathsMatching(predicate func(*Node[T]) bool) []string {
	paths := []string{}

	if predicate(self) {
		paths = append(paths, self.GetPath())
	}

	for _, child := range self.Children {
		paths = append(paths, child.GetPathsMatching(predicate)...)
	}

	return paths
}

func (self *Node[T]) GetFilePathsMatching(predicate func(*T) bool) []string {
	return lo.FilterMap(self.GetLeaves(), func(node *Node[T], _ int) (string, bool) {
		return node.GetPath(), predicate(node.File)
	})
}

func (self *Node[T]) GetLeaves() []*Node[T] {
	if self.IsFile() {
		return []*Node[T]{self}
	}

	return lo.FlatMap(self.Children, func(child *Node[T], _ int) []*Node[T] {
		return child.GetLeaves()
	})
}

func (self *Node[T]) ID() string {
	return self.GetPath()
}

func (self *Node[T]) Description() string {
	return self.GetPath()
}

func (self *Node[T]) Name() string {
	return path.Base(self.path)
}
