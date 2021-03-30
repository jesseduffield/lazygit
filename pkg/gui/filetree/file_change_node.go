package filetree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

type FileChangeNode struct {
	Children         []*FileChangeNode
	File             *models.File
	Path             string // e.g. '/path/to/mydir'
	CompressionLevel int    // equal to the number of forward slashes you'll see in the path when it's rendered in tree mode
}

// methods satisfying ListItem interface

func (s *FileChangeNode) ID() string {
	return s.GetPath()
}

func (s *FileChangeNode) Description() string {
	return s.GetPath()
}

// methods satisfying INode interface

func (s *FileChangeNode) IsLeaf() bool {
	return s.File != nil
}

func (s *FileChangeNode) GetPath() string {
	return s.Path
}

func (s *FileChangeNode) GetChildren() []INode {
	result := make([]INode, len(s.Children))
	for i, child := range s.Children {
		result[i] = child
	}

	return result
}

func (s *FileChangeNode) SetChildren(children []INode) {
	castChildren := make([]*FileChangeNode, len(children))
	for i, child := range children {
		castChildren[i] = child.(*FileChangeNode)
	}

	s.Children = castChildren
}

func (s *FileChangeNode) GetCompressionLevel() int {
	return s.CompressionLevel
}

func (s *FileChangeNode) SetCompressionLevel(level int) {
	s.CompressionLevel = level
}

// methods utilising generic functions for INodes

func (s *FileChangeNode) Sort() {
	sortNode(s)
}

func (s *FileChangeNode) ForEachFile(cb func(*models.File) error) error {
	return forEachLeaf(s, func(n INode) error {
		castNode := n.(*FileChangeNode)
		return cb(castNode.File)
	})
}

func (s *FileChangeNode) Any(test func(node *FileChangeNode) bool) bool {
	return any(s, func(n INode) bool {
		castNode := n.(*FileChangeNode)
		return test(castNode)
	})
}

func (n *FileChangeNode) Flatten(collapsedPaths map[string]bool) []*FileChangeNode {
	results := flatten(n, collapsedPaths)
	nodes := make([]*FileChangeNode, len(results))
	for i, result := range results {
		nodes[i] = result.(*FileChangeNode)
	}

	return nodes
}

func (node *FileChangeNode) GetNodeAtIndex(index int, collapsedPaths map[string]bool) *FileChangeNode {
	return getNodeAtIndex(node, index, collapsedPaths).(*FileChangeNode)
}

func (node *FileChangeNode) GetIndexForPath(path string, collapsedPaths map[string]bool) (int, bool) {
	return getIndexForPath(node, path, collapsedPaths)
}

func (node *FileChangeNode) Size(collapsedPaths map[string]bool) int {
	return size(node, collapsedPaths)
}

func (s *FileChangeNode) Compress() {
	// with these functions I try to only have type conversion code on the actual struct,
	// but comparing interface values to nil is fraught with danger so I'm duplicating
	// that code here.
	if s == nil {
		return
	}

	compressAux(s)
}

// This ignores the root
func (node *FileChangeNode) GetPathsMatching(test func(*FileChangeNode) bool) []string {
	return getPathsMatching(node, func(n INode) bool {
		return test(n.(*FileChangeNode))
	})
}

func (s *FileChangeNode) GetLeaves() []*FileChangeNode {
	leaves := getLeaves(s)
	castLeaves := make([]*FileChangeNode, len(leaves))
	for i := range leaves {
		castLeaves[i] = leaves[i].(*FileChangeNode)
	}

	return castLeaves
}

// extra methods

func (s *FileChangeNode) GetHasUnstagedChanges() bool {
	return s.AnyFile(func(file *models.File) bool { return file.HasUnstagedChanges })
}

func (s *FileChangeNode) GetHasStagedChanges() bool {
	return s.AnyFile(func(file *models.File) bool { return file.HasStagedChanges })
}

func (s *FileChangeNode) GetHasInlineMergeConflicts() bool {
	return s.AnyFile(func(file *models.File) bool { return file.HasInlineMergeConflicts })
}

func (s *FileChangeNode) GetIsTracked() bool {
	return s.AnyFile(func(file *models.File) bool { return file.Tracked })
}

func (s *FileChangeNode) AnyFile(test func(file *models.File) bool) bool {
	return s.Any(func(node *FileChangeNode) bool {
		return node.IsLeaf() && test(node.File)
	})
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

		return fmt.Sprintf("%s%s%s", prevName, " → ", name)
	}

	return name
}
