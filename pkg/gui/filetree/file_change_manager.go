package filetree

import (
	"fmt"
	"os"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/sirupsen/logrus"
)

const EXPANDED_ARROW = "▼"
const COLLAPSED_ARROW = "►"

const INNER_ITEM = "├─ "
const LAST_ITEM = "└─ "
const NESTED = "│  "
const NOTHING = "   "

type FileChangeManager struct {
	files          []*models.File
	tree           *FileChangeNode
	showTree       bool
	log            *logrus.Entry
	collapsedPaths map[string]bool
}

func NewFileChangeManager(files []*models.File, log *logrus.Entry, showTree bool) *FileChangeManager {
	return &FileChangeManager{
		files:          files,
		log:            log,
		showTree:       showTree,
		collapsedPaths: map[string]bool{},
	}
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

func (m *FileChangeManager) Render(diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	return m.renderAux(m.tree, "", -1, diffName, submoduleConfigs)
}

func (m *FileChangeManager) IsCollapsed(s *FileChangeNode) bool {
	return m.collapsedPaths[s.GetPath()]
}

func (m *FileChangeManager) ToggleCollapsed(s *FileChangeNode) {
	m.collapsedPaths[s.GetPath()] = !m.collapsedPaths[s.GetPath()]
}

func (m *FileChangeManager) renderAux(s *FileChangeNode, prefix string, depth int, diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
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

func (m *FileChangeManager) ExpandToPath(path string) {
	// need every directory along the way
	split := strings.Split(path, string(os.PathSeparator))
	for i := range split {
		dir := strings.Join(split[0:i+1], string(os.PathSeparator))
		m.collapsedPaths[dir] = false
	}
}
