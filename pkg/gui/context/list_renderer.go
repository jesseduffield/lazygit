package context

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type NonModelItem struct {
	// Where in the model this should be inserted
	Index int
	// Content to render
	Content string
	// The column from which to render the item
	Column int
}

type ListRenderer struct {
	list types.IList
	// Function to get the display strings for each model item in the given
	// range. startIdx and endIdx are model indices. For each model item, return
	// an array of strings, one for each column; the list renderer will take
	// care of aligning the columns appropriately.
	getDisplayStrings func(startIdx int, endIdx int) [][]string
	// Alignment for each column. If nil, the default is left alignment
	getColumnAlignments func() []utils.Alignment
	// Function to insert non-model items (e.g. section headers). If nil, no
	// such items are inserted
	getNonModelItems func() []*NonModelItem

	// The remaining fields are private and shouldn't be initialized by clients
	numNonModelItems        int
	viewIndicesByModelIndex []int
	modelIndicesByViewIndex []int
}

func (self *ListRenderer) GetList() types.IList {
	return self.list
}

func (self *ListRenderer) ModelIndexToViewIndex(modelIndex int) int {
	modelIndex = lo.Clamp(modelIndex, 0, self.list.Len())
	if self.viewIndicesByModelIndex != nil {
		return self.viewIndicesByModelIndex[modelIndex]
	}

	return modelIndex
}

func (self *ListRenderer) ViewIndexToModelIndex(viewIndex int) int {
	viewIndex = utils.Clamp(viewIndex, 0, self.list.Len()+self.numNonModelItems)
	if self.modelIndicesByViewIndex != nil {
		return self.modelIndicesByViewIndex[viewIndex]
	}

	return viewIndex
}

// startIdx and endIdx are view indices, not model indices. If you want to
// render the whole list, pass -1 for both.
func (self *ListRenderer) renderLines(startIdx int, endIdx int) string {
	var columnAlignments []utils.Alignment
	if self.getColumnAlignments != nil {
		columnAlignments = self.getColumnAlignments()
	}
	nonModelItems := []*NonModelItem{}
	self.numNonModelItems = 0
	if self.getNonModelItems != nil {
		nonModelItems = self.getNonModelItems()
		self.prepareConversionArrays(nonModelItems)
	}
	startModelIdx := 0
	if startIdx == -1 {
		startIdx = 0
	} else {
		startModelIdx = self.ViewIndexToModelIndex(startIdx)
	}
	endModelIdx := self.list.Len()
	if endIdx == -1 {
		endIdx = endModelIdx + len(nonModelItems)
	} else {
		endModelIdx = self.ViewIndexToModelIndex(endIdx)
	}
	lines, columnPositions := utils.RenderDisplayStrings(
		self.getDisplayStrings(startModelIdx, endModelIdx),
		columnAlignments)
	lines = self.insertNonModelItems(nonModelItems, endIdx, startIdx, lines, columnPositions)
	return strings.Join(lines, "\n")
}

func (self *ListRenderer) prepareConversionArrays(nonModelItems []*NonModelItem) {
	self.numNonModelItems = len(nonModelItems)
	self.viewIndicesByModelIndex = lo.Range(self.list.Len() + 1)
	self.modelIndicesByViewIndex = lo.Range(self.list.Len() + 1)
	offset := 0
	for _, item := range nonModelItems {
		for i := item.Index; i <= self.list.Len(); i++ {
			self.viewIndicesByModelIndex[i]++
		}
		self.modelIndicesByViewIndex = slices.Insert(
			self.modelIndicesByViewIndex, item.Index+offset, self.modelIndicesByViewIndex[item.Index+offset])
		offset++
	}
}

func (self *ListRenderer) insertNonModelItems(
	nonModelItems []*NonModelItem, endIdx int, startIdx int, lines []string, columnPositions []int,
) []string {
	offset := 0
	for _, item := range nonModelItems {
		if item.Index+offset >= endIdx {
			break
		}
		if item.Index+offset >= startIdx {
			padding := strings.Repeat(" ", columnPositions[item.Column])
			lines = slices.Insert(lines, item.Index+offset-startIdx, padding+item.Content)
		}
		offset++
	}
	return lines
}
