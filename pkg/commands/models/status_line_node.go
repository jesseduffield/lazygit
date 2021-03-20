package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type StatusLineNode struct {
	Children  []*StatusLineNode
	File      *File
	Name      string // e.g. 'mydir'
	Path      string // e.g. '/path/to/mydir'
	Collapsed bool
}

func (s *StatusLineNode) GetHasUnstagedChanges() bool {
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

func (s *StatusLineNode) GetHasStagedChanges() bool {
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

func (s *StatusLineNode) GetNodeAtIndex(index int, collapsedPaths map[string]bool) *StatusLineNode {
	node, _ := s.getNodeAtIndexAux(index, collapsedPaths)

	return node
}

func (s *StatusLineNode) getNodeAtIndexAux(index int, collapsedPaths map[string]bool) (*StatusLineNode, int) {
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

func (s *StatusLineNode) GetIndexForPath(path string, collapsedPaths map[string]bool) (int, bool) {
	return s.getIndexForPathAux(path, collapsedPaths)
}

func (s *StatusLineNode) getIndexForPathAux(path string, collapsedPaths map[string]bool) (int, bool) {
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

func (s *StatusLineNode) IsLeaf() bool {
	return len(s.Children) == 0
}

func (s *StatusLineNode) Size(collapsedPaths map[string]bool) int {
	output := 1

	if !collapsedPaths[s.GetPath()] {
		for _, child := range s.Children {
			output += child.Size(collapsedPaths)
		}
	}

	return output
}

func (s *StatusLineNode) Flatten(collapsedPaths map[string]bool) []*StatusLineNode {
	arr := []*StatusLineNode{s}

	if !collapsedPaths[s.GetPath()] {
		for _, child := range s.Children {
			arr = append(arr, child.Flatten(collapsedPaths)...)
		}
	}

	return arr
}

func (s *StatusLineNode) Sort() {
	s.sortChildren()

	for _, child := range s.Children {
		child.Sort()
	}
}

func (s *StatusLineNode) sortChildren() {
	if s.IsLeaf() {
		return
	}

	sortedChildren := make([]*StatusLineNode, len(s.Children))
	copy(sortedChildren, s.Children)

	sort.Slice(sortedChildren, func(i, j int) bool {
		if !sortedChildren[i].IsLeaf() && sortedChildren[j].IsLeaf() {
			return true
		}
		if sortedChildren[i].IsLeaf() && !sortedChildren[j].IsLeaf() {
			return false
		}

		return sortedChildren[i].Name < sortedChildren[j].Name
	})

	// TODO: think about making this in-place
	s.Children = sortedChildren
}

// returns true if any descendant file is tracked
func (s *StatusLineNode) GetIsTracked() bool {
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

func (s *StatusLineNode) GetPath() string {
	return s.Path
}

func (s *StatusLineNode) Compress() {
	if s == nil {
		return
	}

	s.compressAux()
}

func (s *StatusLineNode) compressAux() *StatusLineNode {
	if s.IsLeaf() {
		return s
	}

	for i := range s.Children {
		for s.Children[i].HasExactlyOneChild() {
			grandchild := s.Children[i].Children[0]
			s.Children[i] = grandchild
		}
	}

	for i, child := range s.Children {
		s.Children[i] = child.compressAux()
	}

	return s
}

func (s *StatusLineNode) HasExactlyOneChild() bool {
	return len(s.Children) == 1
}

// This ignores the root
func (s *StatusLineNode) GetPathsMatching(test func(*StatusLineNode) bool) []string {
	paths := []string{}

	if test(s) {
		paths = append(paths, s.GetPath())
	}

	for _, child := range s.Children {
		paths = append(paths, child.GetPathsMatching(test)...)
	}

	return paths
}

func (s *StatusLineNode) ID() string {
	return s.GetPath()
}

func (s *StatusLineNode) Description() string {
	return s.GetPath()
}

func (s *StatusLineNode) ForEachFile(cb func(*File) error) error {
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

func (s *StatusLineNode) NameAtDepth(depth int) string {
	splitName := strings.Split(s.Name, string(os.PathSeparator))
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
