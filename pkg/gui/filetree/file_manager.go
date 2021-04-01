package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/sirupsen/logrus"
)

type FileManager struct {
	files          []*models.File
	tree           *FileNode
	showTree       bool
	log            *logrus.Entry
	collapsedPaths CollapsedPaths
}

func NewFileManager(files []*models.File, log *logrus.Entry, showTree bool) *FileManager {
	return &FileManager{
		files:          files,
		log:            log,
		showTree:       showTree,
		collapsedPaths: CollapsedPaths{},
	}
}

func (m *FileManager) InTreeMode() bool {
	return m.showTree
}

func (m *FileManager) ExpandToPath(path string) {
	m.collapsedPaths.ExpandToPath(path)
}

func (m *FileManager) ToggleShowTree() {
	m.showTree = !m.showTree
	m.SetTree()
}

func (m *FileManager) GetItemAtIndex(index int) *FileNode {
	// need to traverse the three depth first until we get to the index.
	return m.tree.GetNodeAtIndex(index+1, m.collapsedPaths) // ignoring root
}

func (m *FileManager) GetIndexForPath(path string) (int, bool) {
	index, found := m.tree.GetIndexForPath(path, m.collapsedPaths)
	return index - 1, found
}

func (m *FileManager) GetAllItems() []*FileNode {
	if m.tree == nil {
		return nil
	}

	return m.tree.Flatten(m.collapsedPaths)[1:] // ignoring root
}

func (m *FileManager) GetItemsLength() int {
	return m.tree.Size(m.collapsedPaths) - 1 // ignoring root
}

func (m *FileManager) GetAllFiles() []*models.File {
	return m.files
}

func (m *FileManager) SetFiles(files []*models.File) {
	m.files = files

	m.SetTree()
}

func (m *FileManager) SetTree() {
	if m.showTree {
		m.tree = BuildTreeFromFiles(m.files)
	} else {
		m.tree = BuildFlatTreeFromFiles(m.files)
	}
}

func (m *FileManager) IsCollapsed(path string) bool {
	return m.collapsedPaths.IsCollapsed(path)
}

func (m *FileManager) ToggleCollapsed(path string) {
	m.collapsedPaths.ToggleCollapsed(path)
}

func (m *FileManager) Render(diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	// can't rely on renderAux to check for nil because an interface won't be nil if its concrete value is nil
	if m.tree == nil {
		return []string{}
	}

	return renderAux(m.tree, m.collapsedPaths, "", -1, func(n INode, depth int) string {
		castN := n.(*FileNode)
		return presentation.GetFileLine(castN.GetHasUnstagedChanges(), castN.GetHasStagedChanges(), castN.NameAtDepth(depth), diffName, submoduleConfigs, castN.File)
	})
}
