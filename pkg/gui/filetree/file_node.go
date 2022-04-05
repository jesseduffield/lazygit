package filetree

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type FileNode struct {
	Children         []*FileNode
	File             *models.File
	Path             string // e.g. '/path/to/mydir'
	CompressionLevel int    // equal to the number of forward slashes you'll see in the path when it's rendered in tree mode
}

var (
	_ INode          = &FileNode{}
	_ types.ListItem = &FileNode{}
)

func (s *FileNode) ID() string {
	return s.GetPath()
}

func (s *FileNode) Description() string {
	return s.GetPath()
}

// methods satisfying INode interface

// interfaces values whose concrete value is nil are not themselves nil
// hence the existence of this method
func (s *FileNode) IsNil() bool {
	return s == nil
}

func (s *FileNode) IsLeaf() bool {
	return s.File != nil
}

func (s *FileNode) GetPath() string {
	return s.Path
}

func (s *FileNode) GetPreviousPath() string {
	if s.File != nil {
		return s.File.GetPreviousPath()
	}
	return ""
}

func (s *FileNode) GetChildren() []INode {
	return slices.Map(s.Children, func(child *FileNode) INode {
		return child
	})
}

func (s *FileNode) SetChildren(children []INode) {
	castChildren := slices.Map(children, func(child INode) *FileNode {
		return child.(*FileNode)
	})

	s.Children = castChildren
}

func (s *FileNode) GetCompressionLevel() int {
	return s.CompressionLevel
}

func (s *FileNode) SetCompressionLevel(level int) {
	s.CompressionLevel = level
}

// methods utilising generic functions for INodes

func (s *FileNode) Sort() {
	sortNode(s)
}

func (s *FileNode) ForEachFile(cb func(*models.File) error) error {
	return forEachLeaf(s, func(n INode) error {
		castNode := n.(*FileNode)
		return cb(castNode.File)
	})
}

func (s *FileNode) Any(test func(node *FileNode) bool) bool {
	return any(s, func(n INode) bool {
		castNode := n.(*FileNode)
		return test(castNode)
	})
}

func (n *FileNode) Flatten(collapsedPaths *CollapsedPaths) []*FileNode {
	results := flatten(n, collapsedPaths)
	return slices.Map(results, func(result INode) *FileNode {
		return result.(*FileNode)
	})
}

func (node *FileNode) GetNodeAtIndex(index int, collapsedPaths *CollapsedPaths) *FileNode {
	if node == nil {
		return nil
	}

	result := getNodeAtIndex(node, index, collapsedPaths)
	if result == nil {
		// not sure how this can be nil: we probably are missing a mutex somewhere
		return nil
	}

	return result.(*FileNode)
}

func (node *FileNode) GetIndexForPath(path string, collapsedPaths *CollapsedPaths) (int, bool) {
	return getIndexForPath(node, path, collapsedPaths)
}

func (node *FileNode) Size(collapsedPaths *CollapsedPaths) int {
	if node == nil {
		return 0
	}

	return size(node, collapsedPaths)
}

func (s *FileNode) Compress() {
	// with these functions I try to only have type conversion code on the actual struct,
	// but comparing interface values to nil is fraught with danger so I'm duplicating
	// that code here.
	if s == nil {
		return
	}

	compressAux(s)
}

func (node *FileNode) GetFilePathsMatching(test func(*models.File) bool) []string {
	return getPathsMatching(node, func(n INode) bool {
		castNode := n.(*FileNode)
		if castNode.File == nil {
			return false
		}
		return test(castNode.File)
	})
}

func (s *FileNode) GetLeaves() []*FileNode {
	leaves := getLeaves(s)

	return slices.Map(leaves, func(leaf INode) *FileNode {
		return leaf.(*FileNode)
	})
}

// extra methods

func (s *FileNode) GetHasUnstagedChanges() bool {
	return s.AnyFile(func(file *models.File) bool { return file.HasUnstagedChanges })
}

func (s *FileNode) GetHasStagedChanges() bool {
	return s.AnyFile(func(file *models.File) bool { return file.HasStagedChanges })
}

func (s *FileNode) GetHasInlineMergeConflicts() bool {
	return s.AnyFile(func(file *models.File) bool { return file.HasInlineMergeConflicts })
}

func (s *FileNode) GetIsTracked() bool {
	return s.AnyFile(func(file *models.File) bool { return file.Tracked })
}

func (s *FileNode) AnyFile(test func(file *models.File) bool) bool {
	return s.Any(func(node *FileNode) bool {
		return node.IsLeaf() && test(node.File)
	})
}

func (s *FileNode) NameAtDepth(depth int) string {
	splitName := split(s.Path)
	name := join(splitName[depth:])

	if s.File != nil && s.File.IsRename() {
		splitPrevName := split(s.File.PreviousName)

		prevName := s.File.PreviousName
		// if the file has just been renamed inside the same directory, we can shave off
		// the prefix for the previous path too. Otherwise we'll keep it unchanged
		sameParentDir := len(splitName) == len(splitPrevName) && join(splitName[0:depth]) == join(splitPrevName[0:depth])
		if sameParentDir {
			prevName = join(splitPrevName[depth:])
		}

		return prevName + " → " + name
	}

	return name
}
