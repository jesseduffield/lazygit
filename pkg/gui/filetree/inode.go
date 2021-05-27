package filetree

import (
	"fmt"
	"sort"
	"strings"
)

type INode interface {
	IsLeaf() bool
	GetPath() string
	GetChildren() []INode
	SetChildren([]INode)
	GetCompressionLevel() int
	SetCompressionLevel(int)
}

func sortNode(node INode) {
	sortChildren(node)

	for _, child := range node.GetChildren() {
		sortNode(child)
	}
}

func sortChildren(node INode) {
	if node.IsLeaf() {
		return
	}

	children := node.GetChildren()
	sortedChildren := make([]INode, len(children))
	copy(sortedChildren, children)

	sort.Slice(sortedChildren, func(i, j int) bool {
		if !sortedChildren[i].IsLeaf() && sortedChildren[j].IsLeaf() {
			return true
		}
		if sortedChildren[i].IsLeaf() && !sortedChildren[j].IsLeaf() {
			return false
		}

		return sortedChildren[i].GetPath() < sortedChildren[j].GetPath()
	})

	// TODO: think about making this in-place
	node.SetChildren(sortedChildren)
}

func forEachLeaf(node INode, cb func(INode) error) error {
	if node.IsLeaf() {
		if err := cb(node); err != nil {
			return err
		}
	}

	for _, child := range node.GetChildren() {
		if err := forEachLeaf(child, cb); err != nil {
			return err
		}
	}

	return nil
}

func any(node INode, test func(INode) bool) bool {
	if test(node) {
		return true
	}

	for _, child := range node.GetChildren() {
		if any(child, test) {
			return true
		}
	}

	return false
}

func every(node INode, test func(INode) bool) bool {
	if !test(node) {
		return false
	}

	for _, child := range node.GetChildren() {
		if !every(child, test) {
			return false
		}
	}

	return true
}

func flatten(node INode, collapsedPaths map[string]bool) []INode {
	result := []INode{}
	result = append(result, node)

	if !collapsedPaths[node.GetPath()] {
		for _, child := range node.GetChildren() {
			result = append(result, flatten(child, collapsedPaths)...)
		}
	}

	return result
}

func getNodeAtIndex(node INode, index int, collapsedPaths map[string]bool) INode {
	foundNode, _ := getNodeAtIndexAux(node, index, collapsedPaths)

	return foundNode
}

func getNodeAtIndexAux(node INode, index int, collapsedPaths map[string]bool) (INode, int) {
	offset := 1

	if index == 0 {
		return node, offset
	}

	if !collapsedPaths[node.GetPath()] {
		for _, child := range node.GetChildren() {
			foundNode, offsetChange := getNodeAtIndexAux(child, index-offset, collapsedPaths)
			offset += offsetChange
			if foundNode != nil {
				return foundNode, offset
			}
		}
	}

	return nil, offset
}

func getIndexForPath(node INode, path string, collapsedPaths map[string]bool) (int, bool) {
	offset := 0

	if node.GetPath() == path {
		return offset, true
	}

	if !collapsedPaths[node.GetPath()] {
		for _, child := range node.GetChildren() {
			offsetChange, found := getIndexForPath(child, path, collapsedPaths)
			offset += offsetChange + 1
			if found {
				return offset, true
			}
		}
	}

	return offset, false
}

func size(node INode, collapsedPaths map[string]bool) int {
	output := 1

	if !collapsedPaths[node.GetPath()] {
		for _, child := range node.GetChildren() {
			output += size(child, collapsedPaths)
		}
	}

	return output
}

func compressAux(node INode) INode {
	if node.IsLeaf() {
		return node
	}

	children := node.GetChildren()
	for i := range children {
		grandchildren := children[i].GetChildren()
		for len(grandchildren) == 1 {
			grandchildren[0].SetCompressionLevel(children[i].GetCompressionLevel() + 1)
			children[i] = grandchildren[0]
			grandchildren = children[i].GetChildren()
		}
	}

	for i := range children {
		children[i] = compressAux(children[i])
	}

	node.SetChildren(children)

	return node
}

func getPathsMatching(node INode, test func(INode) bool) []string {
	paths := []string{}

	if test(node) {
		paths = append(paths, node.GetPath())
	}

	for _, child := range node.GetChildren() {
		paths = append(paths, getPathsMatching(child, test)...)
	}

	return paths
}

func getLeaves(node INode) []INode {
	if node.IsLeaf() {
		return []INode{node}
	}

	output := []INode{}
	for _, child := range node.GetChildren() {
		output = append(output, getLeaves(child)...)
	}

	return output
}

func renderAux(s INode, collapsedPaths CollapsedPaths, prefix string, depth int, renderLine func(INode, int) string) []string {
	isRoot := depth == -1

	renderLineWithPrefix := func() string {
		return prefix + renderLine(s, depth)
	}

	if s.IsLeaf() {
		if isRoot {
			return []string{}
		}
		return []string{renderLineWithPrefix()}
	}

	if collapsedPaths.IsCollapsed(s.GetPath()) {
		return []string{fmt.Sprintf("%s %s", renderLineWithPrefix(), COLLAPSED_ARROW)}
	}

	arr := []string{}
	if !isRoot {
		arr = append(arr, fmt.Sprintf("%s %s", renderLineWithPrefix(), EXPANDED_ARROW))
	}

	newPrefix := prefix
	if strings.HasSuffix(prefix, LAST_ITEM) {
		newPrefix = strings.TrimSuffix(prefix, LAST_ITEM) + NOTHING
	} else if strings.HasSuffix(prefix, INNER_ITEM) {
		newPrefix = strings.TrimSuffix(prefix, INNER_ITEM) + NESTED
	}

	for i, child := range s.GetChildren() {
		isLast := i == len(s.GetChildren())-1

		var childPrefix string
		if isRoot {
			childPrefix = newPrefix
		} else if isLast {
			childPrefix = newPrefix + LAST_ITEM
		} else {
			childPrefix = newPrefix + INNER_ITEM
		}

		arr = append(arr, renderAux(child, collapsedPaths, childPrefix, depth+1+s.GetCompressionLevel(), renderLine)...)
	}

	return arr
}
