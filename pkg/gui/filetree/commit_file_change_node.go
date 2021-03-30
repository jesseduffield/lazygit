package filetree

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type CommitFileChangeNode struct {
	Children         []*CommitFileChangeNode
	File             *models.CommitFile
	Path             string // e.g. '/path/to/mydir'
	CompressionLevel int    // equal to the number of forward slashes you'll see in the path when it's rendered in tree mode
}

// methods satisfying ListItem interface

func (s *CommitFileChangeNode) ID() string {
	return s.GetPath()
}

func (s *CommitFileChangeNode) Description() string {
	return s.GetPath()
}

// methods satisfying INode interface

func (s *CommitFileChangeNode) IsLeaf() bool {
	return s.File != nil
}

func (s *CommitFileChangeNode) GetPath() string {
	return s.Path
}

func (s *CommitFileChangeNode) GetChildren() []INode {
	result := make([]INode, len(s.Children))
	for i, child := range s.Children {
		result[i] = child
	}

	return result
}

func (s *CommitFileChangeNode) SetChildren(children []INode) {
	castChildren := make([]*CommitFileChangeNode, len(children))
	for i, child := range children {
		castChildren[i] = child.(*CommitFileChangeNode)
	}

	s.Children = castChildren
}

func (s *CommitFileChangeNode) GetCompressionLevel() int {
	return s.CompressionLevel
}

func (s *CommitFileChangeNode) SetCompressionLevel(level int) {
	s.CompressionLevel = level
}

// methods utilising generic functions for INodes

func (s *CommitFileChangeNode) Sort() {
	sortNode(s)
}

func (s *CommitFileChangeNode) ForEachFile(cb func(*models.CommitFile) error) error {
	return forEachLeaf(s, func(n INode) error {
		castNode := n.(*CommitFileChangeNode)
		return cb(castNode.File)
	})
}

func (s *CommitFileChangeNode) Any(test func(node *CommitFileChangeNode) bool) bool {
	return any(s, func(n INode) bool {
		castNode := n.(*CommitFileChangeNode)
		return test(castNode)
	})
}

func (n *CommitFileChangeNode) Flatten(collapsedPaths map[string]bool) []*CommitFileChangeNode {
	results := flatten(n, collapsedPaths)
	nodes := make([]*CommitFileChangeNode, len(results))
	for i, result := range results {
		nodes[i] = result.(*CommitFileChangeNode)
	}

	return nodes
}

func (node *CommitFileChangeNode) GetNodeAtIndex(index int, collapsedPaths map[string]bool) *CommitFileChangeNode {
	return getNodeAtIndex(node, index, collapsedPaths).(*CommitFileChangeNode)
}

func (node *CommitFileChangeNode) GetIndexForPath(path string, collapsedPaths map[string]bool) (int, bool) {
	return getIndexForPath(node, path, collapsedPaths)
}

func (node *CommitFileChangeNode) Size(collapsedPaths map[string]bool) int {
	return size(node, collapsedPaths)
}

func (s *CommitFileChangeNode) Compress() {
	// with these functions I try to only have type conversion code on the actual struct,
	// but comparing interface values to nil is fraught with danger so I'm duplicating
	// that code here.
	if s == nil {
		return
	}

	compressAux(s)
}

// This ignores the root
func (node *CommitFileChangeNode) GetPathsMatching(test func(*CommitFileChangeNode) bool) []string {
	return getPathsMatching(node, func(n INode) bool {
		return test(n.(*CommitFileChangeNode))
	})
}

func (s *CommitFileChangeNode) GetLeaves() []*CommitFileChangeNode {
	leaves := getLeaves(s)
	castLeaves := make([]*CommitFileChangeNode, len(leaves))
	for i := range leaves {
		castLeaves[i] = leaves[i].(*CommitFileChangeNode)
	}

	return castLeaves
}

// extra methods

func (s *CommitFileChangeNode) AnyFile(test func(file *models.CommitFile) bool) bool {
	return s.Any(func(node *CommitFileChangeNode) bool {
		return node.IsLeaf() && test(node.File)
	})
}

func (s *CommitFileChangeNode) NameAtDepth(depth int) string {
	splitName := strings.Split(s.Path, string(os.PathSeparator))
	name := filepath.Join(splitName[depth:]...)

	return name
}
