package filetree

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/sirupsen/logrus"
)

const EXPANDED_ARROW = "▼"
const COLLAPSED_ARROW = "►"

type FileChangeManager struct {
	Files          []*models.File
	Tree           *models.FileChangeNode
	ShowTree       bool
	Log            *logrus.Entry
	CollapsedPaths map[string]bool
}

func NewFileChangeManager(files []*models.File, log *logrus.Entry, showTree bool) *FileChangeManager {
	return &FileChangeManager{
		Files:          files,
		Log:            log,
		ShowTree:       showTree,
		CollapsedPaths: map[string]bool{},
	}
}

func (m *FileChangeManager) GetItemAtIndex(index int) *models.FileChangeNode {
	// need to traverse the three depth first until we get to the index.
	return m.Tree.GetNodeAtIndex(index+1, m.CollapsedPaths) // ignoring root
}

func (m *FileChangeManager) GetIndexForPath(path string) (int, bool) {
	index, found := m.Tree.GetIndexForPath(path, m.CollapsedPaths)
	return index - 1, found
}

func (m *FileChangeManager) GetAllItems() []*models.FileChangeNode {
	if m.Tree == nil {
		return nil
	}

	return m.Tree.Flatten(m.CollapsedPaths)[1:] // ignoring root
}

func (m *FileChangeManager) GetItemsLength() int {
	return m.Tree.Size(m.CollapsedPaths) - 1 // ignoring root
}

func (m *FileChangeManager) GetAllFiles() []*models.File {
	return m.Files
}

func (m *FileChangeManager) SetFiles(files []*models.File) {
	m.Files = files

	m.SetTree()
}

func (m *FileChangeManager) SetTree() {
	if m.ShowTree {
		m.Tree = BuildTreeFromFiles(m.Files)
	} else {
		m.Tree = BuildFlatTreeFromFiles(m.Files)
	}
}

func (m *FileChangeManager) Render(diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	return m.renderAux(m.Tree, "", -1, diffName, submoduleConfigs)
}

const INNER_ITEM = "├─ "
const LAST_ITEM = "└─ "
const NESTED = "│  "
const NOTHING = "   "

func (m *FileChangeManager) IsCollapsed(s *models.FileChangeNode) bool {
	return m.CollapsedPaths[s.GetPath()]
}

func (m *FileChangeManager) ToggleCollapsed(s *models.FileChangeNode) {
	m.CollapsedPaths[s.GetPath()] = !m.CollapsedPaths[s.GetPath()]
}

func (m *FileChangeManager) renderAux(s *models.FileChangeNode, prefix string, depth int, diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	isRoot := depth == -1
	if s == nil {
		return []string{}
	}

	getLine := func() string {
		return prefix + presentation.GetFileLine(s.GetHasUnstagedChanges(), s.GetHasStagedChanges(), s.NameAtDepth(depth), diffName, submoduleConfigs, s.File)
	}

	if s.IsLeaf() {
		if isRoot {
			return []string{}
		}
		return []string{getLine()}
	}

	if m.IsCollapsed(s) {
		return []string{fmt.Sprintf("%s %s", getLine(), COLLAPSED_ARROW)}
	}

	arr := []string{}
	if !isRoot {
		arr = append(arr, fmt.Sprintf("%s %s", getLine(), EXPANDED_ARROW))
	}

	newPrefix := prefix
	if strings.HasSuffix(prefix, LAST_ITEM) {
		newPrefix = strings.TrimSuffix(prefix, LAST_ITEM) + NOTHING
	} else if strings.HasSuffix(prefix, INNER_ITEM) {
		newPrefix = strings.TrimSuffix(prefix, INNER_ITEM) + NESTED
	}

	for i, child := range s.Children {
		isLast := i == len(s.Children)-1

		var childPrefix string
		if isRoot {
			childPrefix = newPrefix
		} else if isLast {
			childPrefix = newPrefix + LAST_ITEM
		} else {
			childPrefix = newPrefix + INNER_ITEM
		}

		arr = append(arr, m.renderAux(child, childPrefix, depth+1+s.CompressionLevel, diffName, submoduleConfigs)...)
	}

	return arr
}
