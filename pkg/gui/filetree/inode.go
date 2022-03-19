package filetree

import "github.com/jesseduffield/generics/slices"

type INode interface {
	IsNil() bool
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

	sortedChildren := slices.Clone(node.GetChildren())

	slices.SortFunc(sortedChildren, func(a, b INode) bool {
		if !a.IsLeaf() && b.IsLeaf() {
			return true
		}
		if a.IsLeaf() && !b.IsLeaf() {
			return false
		}

		return a.GetPath() < b.GetPath()
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

func flatten(node INode, collapsedPaths *CollapsedPaths) []INode {
	result := []INode{}
	result = append(result, node)

	if !collapsedPaths.IsCollapsed(node.GetPath()) {
		for _, child := range node.GetChildren() {
			result = append(result, flatten(child, collapsedPaths)...)
		}
	}

	return result
}

func getNodeAtIndex(node INode, index int, collapsedPaths *CollapsedPaths) INode {
	foundNode, _ := getNodeAtIndexAux(node, index, collapsedPaths)

	return foundNode
}

func getNodeAtIndexAux(node INode, index int, collapsedPaths *CollapsedPaths) (INode, int) {
	offset := 1

	if index == 0 {
		return node, offset
	}

	if !collapsedPaths.IsCollapsed(node.GetPath()) {
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

func getIndexForPath(node INode, path string, collapsedPaths *CollapsedPaths) (int, bool) {
	offset := 0

	if node.GetPath() == path {
		return offset, true
	}

	if !collapsedPaths.IsCollapsed(node.GetPath()) {
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

func size(node INode, collapsedPaths *CollapsedPaths) int {
	output := 1

	if !collapsedPaths.IsCollapsed(node.GetPath()) {
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
		for len(grandchildren) == 1 && !grandchildren[0].IsLeaf() {
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

	return slices.FlatMap(node.GetChildren(), func(child INode) []INode {
		return getLeaves(child)
	})
}
