package filetree

import "github.com/jesseduffield/lazygit/pkg/commands/models"

// CommitFileNode wraps a node and provides some commit-file-specific methods for it.
type CommitFileNode struct {
	*Node[models.CommitFile]
}

func NewCommitFileNode(node *Node[models.CommitFile]) *CommitFileNode {
	if node == nil {
		return nil
	}

	return &CommitFileNode{Node: node}
}

// returns the underlying node, without any commit-file-specific methods attached
func (self *CommitFileNode) Raw() *Node[models.CommitFile] {
	if self == nil {
		return nil
	}

	return self.Node
}
