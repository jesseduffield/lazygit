package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileChangeNode struct {
	Children         []*FileChangeNode
	File             *File
	Path             string // e.g. '/path/to/mydir'
	Collapsed        bool
	CompressionLevel int // equal to the number of forward slashes you'll see in the path when it's rendered
}

func (s *FileChangeNode) GetHasUnstagedChanges() bool {
	if s.IsLeaf() {
		return s.File.HasUnstagedChanges
	}

	for _, child := range s.Children {
		if child.GetHasUnstagedChanges() {
			return true
		}
	}

	return false
}

func (s *FileChangeNode) GetHasStagedChanges() bool {
	if s.IsLeaf() {
		return s.File.HasStagedChanges
	}

	for _, child := range s.Children {
		if child.GetHasStagedChanges() {
			return true
		}
	}

	return false
}

func (s *FileChangeNode) GetNodeAtIndex(index int, collapsedPaths map[string]bool) *FileChangeNode {
	node, _ := s.getNodeAtIndexAux(index, collapsedPaths)

	return node
}

func (s *FileChangeNode) getNodeAtIndexAux(index int, collapsedPaths map[string]bool) (*FileChangeNode, int) {
	offset := 1

	if index == 0 {
		return s, offset
	}

	if !collapsedPaths[s.GetPath()] {
		for _, child := range s.Children {
			node, offsetChange := child.getNodeAtIndexAux(index-offset, collapsedPaths)
			offset += offsetChange
			if node != nil {
				return node, offset
			}
		}
	}

	return nil, offset
}

func (s *FileChangeNode) GetIndexForPath(path string, collapsedPaths map[string]bool) (int, bool) {
	return s.getIndexForPathAux(path, collapsedPaths)
}

func (s *FileChangeNode) getIndexForPathAux(path string, collapsedPaths map[string]bool) (int, bool) {
	offset := 0

	if s.Path == path {
		return offset, true
	}

	if !collapsedPaths[s.GetPath()] {
		for _, child := range s.Children {
			offsetChange, found := child.getIndexForPathAux(path, collapsedPaths)
			offset += offsetChange + 1
			if found {
				return offset, true
			}
		}
	}

	return offset, false
}

func (s *FileChangeNode) IsLeaf() bool {
	return s.File != nil
}

func (s *FileChangeNode) Size(collapsedPaths map[string]bool) int {
	output := 1

	if !collapsedPaths[s.GetPath()] {
		for _, child := range s.Children {
			output += child.Size(collapsedPaths)
		}
	}

	return output
}

func (s *FileChangeNode) Flatten(collapsedPaths map[string]bool) []*FileChangeNode {
	arr := []*FileChangeNode{s}

	if !collapsedPaths[s.GetPath()] {
		for _, child := range s.Children {
			arr = append(arr, child.Flatten(collapsedPaths)...)
		}
	}

	return arr
}

func (s *FileChangeNode) Sort() {
	s.sortChildren()

	for _, child := range s.Children {
		child.Sort()
	}
}

func (s *FileChangeNode) sortChildren() {
	if s.IsLeaf() {
		return
	}

	sortedChildren := make([]*FileChangeNode, len(s.Children))
	copy(sortedChildren, s.Children)

	sort.Slice(sortedChildren, func(i, j int) bool {
		if !sortedChildren[i].IsLeaf() && sortedChildren[j].IsLeaf() {
			return true
		}
		if sortedChildren[i].IsLeaf() && !sortedChildren[j].IsLeaf() {
			return false
		}

		return sortedChildren[i].Path < sortedChildren[j].Path
	})

	// TODO: think about making this in-place
	s.Children = sortedChildren
}

// returns true if any descendant file is tracked
func (s *FileChangeNode) GetIsTracked() bool {
	if s.File != nil {
		return s.File.GetIsTracked()
	}

	for _, child := range s.Children {
		if child.GetIsTracked() {
			return true
		}
	}

	return false
}

func (s *FileChangeNode) GetPath() string {
	return s.Path
}

func (s *FileChangeNode) Compress() {
	if s == nil {
		return
	}

	s.compressAux()
}

func (s *FileChangeNode) compressAux() *FileChangeNode {
	if s.IsLeaf() {
		return s
	}

	for i := range s.Children {
		for s.Children[i].HasExactlyOneChild() {
			prevCompressionLevel := s.Children[i].CompressionLevel
			grandchild := s.Children[i].Children[0]
			s.Children[i] = grandchild
			s.Children[i].CompressionLevel = prevCompressionLevel + 1
		}
	}

	for i, child := range s.Children {
		s.Children[i] = child.compressAux()
	}

	return s
}

func (s *FileChangeNode) HasExactlyOneChild() bool {
	return len(s.Children) == 1
}

// This ignores the root
func (s *FileChangeNode) GetPathsMatching(test func(*FileChangeNode) bool) []string {
	paths := []string{}

	if test(s) {
		paths = append(paths, s.GetPath())
	}

	for _, child := range s.Children {
		paths = append(paths, child.GetPathsMatching(test)...)
	}

	return paths
}

func (s *FileChangeNode) ID() string {
	return s.GetPath()
}

func (s *FileChangeNode) Description() string {
	return s.GetPath()
}

func (s *FileChangeNode) ForEachFile(cb func(*File) error) error {
	if s.File != nil {
		if err := cb(s.File); err != nil {
			return err
		}
	}

	for _, child := range s.Children {
		if err := child.ForEachFile(cb); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileChangeNode) GetLeaves() []*FileChangeNode {
	if s.IsLeaf() {
		return []*FileChangeNode{s}
	}

	output := []*FileChangeNode{}
	for _, child := range s.Children {
		output = append(output, child.GetLeaves()...)
	}

	return output
}

func (s *FileChangeNode) NameAtDepth(depth int) string {
	splitName := strings.Split(s.Path, string(os.PathSeparator))
	name := filepath.Join(splitName[depth:]...)

	if s.File != nil && s.File.IsRename() {
		splitPrevName := strings.Split(s.File.PreviousName, string(os.PathSeparator))

		prevName := s.File.PreviousName
		// if the file has just been renamed inside the same directory, we can shave off
		// the prefix for the previous path too. Otherwise we'll keep it unchanged
		sameParentDir := len(splitName) == len(splitPrevName) && filepath.Join(splitName[0:depth]...) == filepath.Join(splitPrevName[0:depth]...)
		if sameParentDir {
			prevName = filepath.Join(splitPrevName[depth:]...)
		}

		return fmt.Sprintf("%s%s%s", prevName, " -> ", name)
	}

	return name
}
