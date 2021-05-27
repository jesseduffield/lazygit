package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/sirupsen/logrus"
)

type CommitFileManager struct {
	files          []*models.CommitFile
	tree           *CommitFileNode
	showTree       bool
	log            *logrus.Entry
	collapsedPaths CollapsedPaths
	// parent is the identifier of the parent object e.g. a commit SHA if this commit file is for a commit, or a stash entry ref like 'stash@{1}'
	parent string
}

func (m *CommitFileManager) GetParent() string {
	return m.parent
}

func NewCommitFileManager(files []*models.CommitFile, log *logrus.Entry, showTree bool) *CommitFileManager {
	return &CommitFileManager{
		files:          files,
		log:            log,
		showTree:       showTree,
		collapsedPaths: CollapsedPaths{},
	}
}

func (m *CommitFileManager) ExpandToPath(path string) {
	m.collapsedPaths.ExpandToPath(path)
}

func (m *CommitFileManager) ToggleShowTree() {
	m.showTree = !m.showTree
	m.SetTree()
}

func (m *CommitFileManager) GetItemAtIndex(index int) *CommitFileNode {
	// need to traverse the three depth first until we get to the index.
	return m.tree.GetNodeAtIndex(index+1, m.collapsedPaths) // ignoring root
}

func (m *CommitFileManager) GetIndexForPath(path string) (int, bool) {
	index, found := m.tree.GetIndexForPath(path, m.collapsedPaths)
	return index - 1, found
}

func (m *CommitFileManager) GetAllItems() []*CommitFileNode {
	if m.tree == nil {
		return nil
	}

	return m.tree.Flatten(m.collapsedPaths)[1:] // ignoring root
}

func (m *CommitFileManager) GetItemsLength() int {
	return m.tree.Size(m.collapsedPaths) - 1 // ignoring root
}

func (m *CommitFileManager) GetAllFiles() []*models.CommitFile {
	return m.files
}

func (m *CommitFileManager) SetFiles(files []*models.CommitFile, parent string) {
	m.files = files
	m.parent = parent

	m.SetTree()
}

func (m *CommitFileManager) SetTree() {
	if m.showTree {
		m.tree = BuildTreeFromCommitFiles(m.files)
	} else {
		m.tree = BuildFlatTreeFromCommitFiles(m.files)
	}
}

func (m *CommitFileManager) IsCollapsed(path string) bool {
	return m.collapsedPaths.IsCollapsed(path)
}

func (m *CommitFileManager) ToggleCollapsed(path string) {
	m.collapsedPaths.ToggleCollapsed(path)
}

func (m *CommitFileManager) Render(diffName string, patchManager *patch.PatchManager) []string {
	// can't rely on renderAux to check for nil because an interface won't be nil if its concrete value is nil
	if m.tree == nil {
		return []string{}
	}

	return renderAux(m.tree, m.collapsedPaths, "", -1, func(n INode, depth int) string {
		castN := n.(*CommitFileNode)

		// This is a little convoluted because we're dealing with either a leaf or a non-leaf.
		// But this code actually applies to both. If it's a leaf, the status will just
		// be whatever status it is, but if it's a non-leaf it will determine its status
		// based on the leaves of that subtree
		var status patch.PatchStatus
		if castN.EveryFile(func(file *models.CommitFile) bool {
			return patchManager.GetFileStatus(file.Name, m.parent) == patch.WHOLE
		}) {
			status = patch.WHOLE
		} else if castN.EveryFile(func(file *models.CommitFile) bool {
			return patchManager.GetFileStatus(file.Name, m.parent) == patch.UNSELECTED
		}) {
			status = patch.UNSELECTED
		} else {
			status = patch.PART
		}

		return presentation.GetCommitFileLine(castN.NameAtDepth(depth), diffName, castN.File, status)
	})
}
