package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/sirupsen/logrus"
)

type CommitFileChangeManager struct {
	files          []*models.CommitFile
	tree           *CommitFileChangeNode
	showTree       bool
	log            *logrus.Entry
	collapsedPaths CollapsedPaths
	// parent is the identifier of the parent object e.g. a commit SHA if this commit file is for a commit, or a stash entry ref like 'stash@{1}'
	parent string
}

func (m *CommitFileChangeManager) GetParent() string {
	return m.parent
}

func NewCommitFileChangeManager(files []*models.CommitFile, log *logrus.Entry, showTree bool) *CommitFileChangeManager {
	return &CommitFileChangeManager{
		files:          files,
		log:            log,
		showTree:       showTree,
		collapsedPaths: CollapsedPaths{},
	}
}

func (m *CommitFileChangeManager) ExpandToPath(path string) {
	m.collapsedPaths.ExpandToPath(path)
}

func (m *CommitFileChangeManager) ToggleShowTree() {
	m.showTree = !m.showTree
	m.SetTree()
}

func (m *CommitFileChangeManager) GetItemAtIndex(index int) *CommitFileChangeNode {
	// need to traverse the three depth first until we get to the index.
	return m.tree.GetNodeAtIndex(index+1, m.collapsedPaths) // ignoring root
}

func (m *CommitFileChangeManager) GetIndexForPath(path string) (int, bool) {
	index, found := m.tree.GetIndexForPath(path, m.collapsedPaths)
	return index - 1, found
}

func (m *CommitFileChangeManager) GetAllItems() []*CommitFileChangeNode {
	if m.tree == nil {
		return nil
	}

	return m.tree.Flatten(m.collapsedPaths)[1:] // ignoring root
}

func (m *CommitFileChangeManager) GetItemsLength() int {
	return m.tree.Size(m.collapsedPaths) - 1 // ignoring root
}

func (m *CommitFileChangeManager) GetAllFiles() []*models.CommitFile {
	return m.files
}

func (m *CommitFileChangeManager) SetFiles(files []*models.CommitFile, parent string) {
	m.files = files
	m.parent = parent

	m.SetTree()
}

func (m *CommitFileChangeManager) SetTree() {
	if m.showTree {
		m.tree = BuildTreeFromCommitFiles(m.files)
	} else {
		m.tree = BuildFlatTreeFromCommitFiles(m.files)
	}
}

func (m *CommitFileChangeManager) IsCollapsed(path string) bool {
	return m.collapsedPaths.IsCollapsed(path)
}

func (m *CommitFileChangeManager) ToggleCollapsed(path string) {
	m.collapsedPaths.ToggleCollapsed(path)
}

func (m *CommitFileChangeManager) Render(diffName string, patchManager *patch.PatchManager) []string {
	return renderAux(m.tree, m.collapsedPaths, "", -1, func(n INode, depth int) string {
		castN := n.(*CommitFileChangeNode)
		return presentation.GetCommitFileLine(castN.NameAtDepth(depth), diffName, castN.File, patchManager, m.parent)
	})
}
