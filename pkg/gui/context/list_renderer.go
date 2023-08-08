package context

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ListRenderer struct {
	list              types.IList
	getDisplayStrings func(startIdx int, endIdx int) [][]string
	// Alignment for each column. If nil, the default is left alignment
	getColumnAlignments func() []utils.Alignment
}

func (self *ListRenderer) GetList() types.IList {
	return self.list
}

func (self *ListRenderer) renderLines(startIdx int, endIdx int) string {
	var columnAlignments []utils.Alignment
	if self.getColumnAlignments != nil {
		columnAlignments = self.getColumnAlignments()
	}
	lines := utils.RenderDisplayStrings(
		self.getDisplayStrings(startIdx, utils.Min(endIdx, self.list.Len())),
		columnAlignments)
	return strings.Join(lines, "\n")
}
