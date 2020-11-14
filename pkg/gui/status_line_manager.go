package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/sirupsen/logrus"
)

const EXPANDED_ARROW = "▼"
const COLLAPSED_ARROW = "►"

type StatusLineManager struct {
	Files    []*models.File
	Tree     *models.StatusLineNode
	TreeMode bool
	Log      *logrus.Entry
}

func (m *StatusLineManager) GetItemAtIndex(index int) *models.StatusLineNode {
	if m.TreeMode {
		// need to traverse the three depth first until we get to the index.
		return m.Tree.GetNodeAtIndex(index + 1) // ignoring root
	}

	m.Log.Warn(index)
	if index > len(m.Files)-1 {
		return nil
	}

	return &models.StatusLineNode{File: m.Files[index]}
}

func (m *StatusLineManager) GetAllItems() []*models.StatusLineNode {
	return m.Tree.Flatten()[1:] // ignoring root
}

func (m *StatusLineManager) GetItemsLength() int {
	return m.Tree.Size() - 1 // ignoring root
}

func (m *StatusLineManager) GetAllFiles() []*models.File {
	return m.Files
}

func (m *StatusLineManager) SetFiles(files []*models.File) {
	m.Files = files
	m.Tree = GetTreeFromStatusFiles(files)
}

func (m *StatusLineManager) Render(diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	return m.renderAux(m.Tree, -1, diffName, submoduleConfigs)
}

func (m *StatusLineManager) renderAux(s *models.StatusLineNode, depth int, diffName string, submoduleConfigs []*models.SubmoduleConfig) []string {
	if s == nil {
		return []string{}
	}

	getLine := func() string {
		return strings.Repeat("  ", depth) + presentation.GetStatusNodeLine(s.HasUnstagedChanges(), s.GetShortStatus(), s.Name, diffName, submoduleConfigs, s.File)
	}

	if s.IsLeaf() {
		if depth == -1 {
			return []string{}
		}
		return []string{getLine()}
	}

	if s.Collapsed {
		return []string{fmt.Sprintf("%s%s %s", strings.Repeat("  ", depth), s.Name, COLLAPSED_ARROW)}
	}

	arr := []string{}
	if depth > -1 {
		arr = append(arr, fmt.Sprintf("%s%s %s", strings.Repeat("  ", depth), s.Name, EXPANDED_ARROW))
	}

	for _, child := range s.Children {
		arr = append(arr, m.renderAux(child, depth+1, diffName, submoduleConfigs)...)
	}

	return arr
}
