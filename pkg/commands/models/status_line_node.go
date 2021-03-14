package models

import (
	"fmt"
	"sort"
)

type StatusLineNode struct {
	Children  []*StatusLineNode
	File      *File
	Name      string // e.g. 'mydir'
	Path      string // e.g. '/path/to/mydir'
	Collapsed bool
}

func (s *StatusLineNode) GetShortStatus() string {
	// need to see if any child has unstaged changes.
	if s.IsLeaf() {
		return s.File.ShortStatus
	}

	firstChar := " "
	secondChar := " "
	if s.GetHasStagedChanges() {
		firstChar = "M"
	}
	if s.GetHasUnstagedChanges() {
		secondChar = "M"
	}

	return firstChar + secondChar
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

func (s *StatusLineNode) GetNodeAtIndex(index int) *StatusLineNode {
	node, _ := s.getNodeAtIndexAux(index)

	return node
}

func (s *StatusLineNode) getNodeAtIndexAux(index int) (*StatusLineNode, int) {
	offset := 1

	if index == 0 {
		return s, offset
	}

	for _, child := range s.Children {
		node, offsetChange := child.getNodeAtIndexAux(index - offset)
		offset += offsetChange
		if node != nil {
			return node, offset
		}
	}

	return nil, offset
}

func (s *StatusLineNode) IsLeaf() bool {
	return len(s.Children) == 0
}

func (s *StatusLineNode) Size() int {
	output := 1

	for _, child := range s.Children {
		output += child.Size()
	}

	return output
}

func (s *StatusLineNode) Flatten() []*StatusLineNode {
	arr := []*StatusLineNode{s}

	for _, child := range s.Children {
		arr = append(arr, child.Flatten()...)
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

func (s *StatusLineNode) GetIsTracked() bool {
	if s.File != nil {
		return s.File.GetIsTracked()
	}

	// pretty sure I'm allowed to do this
	return true
}

func (s *StatusLineNode) GetPath() string {
	return s.Path
}

func (s *StatusLineNode) HasExactlyOneChild() bool {
	return len(s.Children) == 1
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

	for i, child := range s.Children {
		if child.HasExactlyOneChild() {
			grandchild := child.Children[0]
			grandchild.Name = fmt.Sprintf("%s/%s", child.Name, grandchild.Name)
			s.Children[i] = grandchild
		}
	}

	for i, child := range s.Children {
		s.Children[i] = child.compressAux()
	}

	return s
}
