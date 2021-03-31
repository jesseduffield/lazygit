package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/sirupsen/logrus"
)

type FileChangeManager struct {
	files          []*models.File
	tree           *FileChangeNode
	showTree       bool
	log            *logrus.Entry
	collapsedPaths CollapsedPaths
}

func NewFileChangeManager(files []*models.File, log *logrus.Entry, showTree bool) *FileChangeManager {
	return &FileChangeManager{
		files:          files,
		log:            log,
		showTree:       showTree,
		collapsedPaths: CollapsedPaths{},
	}
}

func (m *FileChangeManager) ExpandToPath(path string) {
	m.collapsedPaths.ExpandToPath(path)
}

func (m *FileChangeManager) ToggleShowTree() {
	m.showTree = !m.showTree
	m.SetTree()
}

func (m *FileChangeManager) GetItemAtIndex(index int) *FileChangeNode {
	// need to traverse the three depth first until we get to the index.
	return m.tree.GetNodeAtIndex(index+1, m.collapsedPaths) // ignoring root
}

func (m *FileChangeManager) GetIndexForPath(path string) (int, bool) {
	index, found := m.tree.GetIndexForPath(path, m.collapsedPaths)
	return index - 1, found
}

func (m *FileChangeManager) GetAllItems() []*FileChangeNode {
	if m.tree == nil {
		return nil
	}

	return m.tree.Flatten(m.collapsedPaths)[1:] // ignoring root
}

func (m *FileChangeManager) GetItemsLength() int {
	return m.tree.Size(m.collapsedPaths) - 1 // ignoring root
}

func (m *FileChangeManager) GetAllFiles() []*models.File {
	return m.files
}

func (m *FileChangeManager) SetFiles(files []*models.File) {
	m.files = files

	m.SetTree()
}

func (m *FileChangeManager) SetTree() {
	if m.showTree {
		m.tree = BuildTreeFromFiles(m.files)
	} else {
		m.tree = BuildFlatTreeFromFiles(m.files)
	}
}

func (m *FileChangeManager) IsCollapsed(path string) bool {
	return m.collapsedPaths.IsCollapsed(path)
}

func (m *FileChangeManager) ToggleCollapsed(path string) {
	m.collapsedPaths.ToggleCollapsed(path)
}

func (m *FileChangeManager) Render(diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	return renderAux(m.tree, m.collapsedPaths, "", -1, func(n INode, depth int) string {
		castN := n.(*FileChangeNode)
		return presentation.GetFileLine(castN.GetHasUnstagedChanges(), castN.GetHasStagedChanges(), castN.NameAtDepth(depth), diffName, submoduleConfigs, castN.File)
	})
}
