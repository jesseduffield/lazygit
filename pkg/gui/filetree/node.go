package filetree

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Represents a file or directory in a file tree.
type Node[T any] struct {
	// File will be nil if the node is a directory.
	File *T

	// If the node is a directory, Children contains the contents of the directory,
	// otherwise it's nil.
	Children []*Node[T]

	// path of the file/directory
	Path string

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

func (self *Node[T]) GetPath() string {
	return self.Path
}

func (self *Node[T]) Sort() {
	self.SortChildren()

	for _, child := range self.Children {
		child.Sort()
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

func (self *Node[T]) SortChildren() {
	if self.IsFile() {
		return
	}

	children := slices.Clone(self.Children)

	slices.SortFunc(children, func(a, b *Node[T]) bool {
		if !a.IsFile() && b.IsFile() {
			return true
		}
		if a.IsFile() && !b.IsFile() {
			return false
		}

		return a.GetPath() < b.GetPath()
	})

	// TODO: think about making this in-place
	self.Children = children
}

func (self *Node[T]) Some(test func(*Node[T]) bool) bool {
	if test(self) {
		return true
	}

	for _, child := range self.Children {
		if child.Some(test) {
			return true
		}
	}

	return false
}

func (self *Node[T]) SomeFile(test func(*T) bool) bool {
	if self.IsFile() {
		if test(self.File) {
			return true
		}
	} else {
		for _, child := range self.Children {
			if child.SomeFile(test) {
				return true
			}
		}
	}

	return false
}

func (self *Node[T]) Every(test func(*Node[T]) bool) bool {
	if !test(self) {
		return false
	}

	for _, child := range self.Children {
		if !child.Every(test) {
			return false
		}
	}

	return true
}

func (self *Node[T]) EveryFile(test func(*T) bool) bool {
	if self.IsFile() {
		if !test(self.File) {
			return false
		}
	} else {
		for _, child := range self.Children {
			if !child.EveryFile(test) {
				return false
			}
		}
	}

	return true
}

func (self *Node[T]) Flatten(collapsedPaths *CollapsedPaths) []*Node[T] {
	result := []*Node[T]{self}

	if len(self.Children) > 0 && !collapsedPaths.IsCollapsed(self.GetPath()) {
		result = append(result, slices.FlatMap(self.Children, func(child *Node[T]) []*Node[T] {
			return child.Flatten(collapsedPaths)
		})...)
	}

	return result
}

func (self *Node[T]) GetNodeAtIndex(index int, collapsedPaths *CollapsedPaths) *Node[T] {
	if self == nil {
		return nil
	}

	node, _ := self.getNodeAtIndexAux(index, collapsedPaths)

	return node
}

func (self *Node[T]) getNodeAtIndexAux(index int, collapsedPaths *CollapsedPaths) (*Node[T], int) {
	offset := 1

	if index == 0 {
		return self, offset
	}

	if !collapsedPaths.IsCollapsed(self.GetPath()) {
		for _, child := range self.Children {
			foundNode, offsetChange := child.getNodeAtIndexAux(index-offset, collapsedPaths)
			offset += offsetChange
			if foundNode != nil {
				return foundNode, offset
			}
		}
	}

	return nil, offset
}

func (self *Node[T]) GetIndexForPath(path string, collapsedPaths *CollapsedPaths) (int, bool) {
	offset := 0

	if self.GetPath() == path {
		return offset, true
	}

	if !collapsedPaths.IsCollapsed(self.GetPath()) {
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

	if !collapsedPaths.IsCollapsed(self.GetPath()) {
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

func (self *Node[T]) GetPathsMatching(test func(*Node[T]) bool) []string {
	paths := []string{}

	if test(self) {
		paths = append(paths, self.GetPath())
	}

	for _, child := range self.Children {
		paths = append(paths, child.GetPathsMatching(test)...)
	}

	return paths
}

func (self *Node[T]) GetFilePathsMatching(test func(*T) bool) []string {
	matchingFileNodes := slices.Filter(self.GetLeaves(), func(node *Node[T]) bool {
		return test(node.File)
	})

	return slices.Map(matchingFileNodes, func(node *Node[T]) string {
		return node.GetPath()
	})
}

func (self *Node[T]) GetLeaves() []*Node[T] {
	if self.IsFile() {
		return []*Node[T]{self}
	}

	return slices.FlatMap(self.Children, func(child *Node[T]) []*Node[T] {
		return child.GetLeaves()
	})
}

func (self *Node[T]) ID() string {
	return self.GetPath()
}

func (self *Node[T]) Description() string {
	return self.GetPath()
}
