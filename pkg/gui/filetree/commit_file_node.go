package filetree

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitFileNode struct {
	Children         []*CommitFileNode
	File             *models.CommitFile
	Path             string // e.g. '/path/to/mydir'
	CompressionLevel int    // equal to the number of forward slashes you'll see in the path when it's rendered in tree mode
}

var (
	_ INode          = &CommitFileNode{}
	_ types.ListItem = &CommitFileNode{}
)

func (s *CommitFileNode) ID() string {
	return s.GetPath()
}

func (s *CommitFileNode) Description() string {
	return s.GetPath()
}

// methods satisfying INode interface

func (s *CommitFileNode) IsNil() bool {
	return s == nil
}

func (s *CommitFileNode) IsLeaf() bool {
	return s.File != nil
}

func (s *CommitFileNode) GetPath() string {
	return s.Path
}

func (s *CommitFileNode) GetChildren() []INode {
	return slices.Map(s.Children, func(child *CommitFileNode) INode {
		return child
	})
}

func (s *CommitFileNode) SetChildren(children []INode) {
	castChildren := slices.Map(children, func(child INode) *CommitFileNode {
		return child.(*CommitFileNode)
	})

	s.Children = castChildren
}

func (s *CommitFileNode) GetCompressionLevel() int {
	return s.CompressionLevel
}

func (s *CommitFileNode) SetCompressionLevel(level int) {
	s.CompressionLevel = level
}

// methods utilising generic functions for INodes

func (s *CommitFileNode) Sort() {
	sortNode(s)
}

func (s *CommitFileNode) ForEachFile(cb func(*models.CommitFile) error) error {
	return forEachLeaf(s, func(n INode) error {
		castNode := n.(*CommitFileNode)
		return cb(castNode.File)
	})
}

func (s *CommitFileNode) Any(test func(node *CommitFileNode) bool) bool {
	return any(s, func(n INode) bool {
		castNode := n.(*CommitFileNode)
		return test(castNode)
	})
}

func (s *CommitFileNode) Every(test func(node *CommitFileNode) bool) bool {
	return every(s, func(n INode) bool {
		castNode := n.(*CommitFileNode)
		return test(castNode)
	})
}

func (s *CommitFileNode) EveryFile(test func(file *models.CommitFile) bool) bool {
	return every(s, func(n INode) bool {
		castNode := n.(*CommitFileNode)

		return castNode.File == nil || test(castNode.File)
	})
}

func (n *CommitFileNode) Flatten(collapsedPaths *CollapsedPaths) []*CommitFileNode {
	results := flatten(n, collapsedPaths)

	return slices.Map(results, func(result INode) *CommitFileNode {
		return result.(*CommitFileNode)
	})
}

func (node *CommitFileNode) GetNodeAtIndex(index int, collapsedPaths *CollapsedPaths) *CommitFileNode {
	if node == nil {
		return nil
	}

	result := getNodeAtIndex(node, index, collapsedPaths)
	if result == nil {
		// not sure how this can be nil: we probably are missing a mutex somewhere
		return nil
	}

	return result.(*CommitFileNode)
}

func (node *CommitFileNode) GetIndexForPath(path string, collapsedPaths *CollapsedPaths) (int, bool) {
	return getIndexForPath(node, path, collapsedPaths)
}

func (node *CommitFileNode) Size(collapsedPaths *CollapsedPaths) int {
	if node == nil {
		return 0
	}

	return size(node, collapsedPaths)
}

func (s *CommitFileNode) Compress() {
	// with these functions I try to only have type conversion code on the actual struct,
	// but comparing interface values to nil is fraught with danger so I'm duplicating
	// that code here.
	if s == nil {
		return
	}

	compressAux(s)
}

func (s *CommitFileNode) GetLeaves() []*CommitFileNode {
	leaves := getLeaves(s)

	return slices.Map(leaves, func(leaf INode) *CommitFileNode {
		return leaf.(*CommitFileNode)
	})
}

// extra methods

func (s *CommitFileNode) AnyFile(test func(file *models.CommitFile) bool) bool {
	return s.Any(func(node *CommitFileNode) bool {
		return node.IsLeaf() && test(node.File)
	})
}

func (s *CommitFileNode) NameAtDepth(depth int) string {
	splitName := split(s.Path)
	name := join(splitName[depth:])

	return name
}
