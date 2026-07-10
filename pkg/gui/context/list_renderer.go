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
	columnPositions []int
}

func (self *ListRenderer) GetList() types.IList {
	return self.list
}

func (self *ListRenderer) getNonModelItemList() []*NonModelItem {
	if self.getNonModelItems == nil {
		return nil
	}
	return self.getNonModelItems()
}

func (self *ListRenderer) ModelIndexToViewIndex(modelIndex int) int {
	return modelIndexToViewIndex(self.list.Len(), self.getNonModelItemList(), modelIndex)
}

func (self *ListRenderer) ViewIndexToModelIndex(viewIndex int) int {
	return viewIndexToModelIndex(self.list.Len(), self.getNonModelItemList(), viewIndex)
}

// modelToViewIndexConverter returns a model-to-view index conversion that
// reuses a single snapshot of the non-model items. Callers that convert many
// indices in a row (e.g. search, which converts every commit) should use this
// rather than calling ModelIndexToViewIndex per index, which would rebuild the
// non-model items each time.
func (self *ListRenderer) modelToViewIndexConverter() func(modelIndex int) int {
	listLength := self.list.Len()
	nonModelItems := self.getNonModelItemList()
	return func(modelIndex int) int {
		return modelIndexToViewIndex(listLength, nonModelItems, modelIndex)
	}
}

// The view shows the model items with the non-model items (e.g. section
// headers) inserted at their model indices. The two conversions below are
// computed directly from the current list length and non-model items, so they
// don't depend on the list having been rendered, and they can never be stale
// with respect to a model that changed since the last render (which used to
// cause both wrong results and index-out-of-range panics).
//
// The non-model items are assumed to be ordered by their Index, which is how
// all producers build them; the i-th one therefore ends up at view index
// Index+i.
func modelIndexToViewIndex(listLength int, nonModelItems []*NonModelItem, modelIndex int) int {
	modelIndex = lo.Clamp(modelIndex, 0, listLength)
	// Each non-model item inserted at or before this model item pushes it down
	// by one row in the view.
	viewIndex := modelIndex
	for _, item := range nonModelItems {
		if item.Index <= modelIndex {
			viewIndex++
		}
	}
	return viewIndex
}

func viewIndexToModelIndex(listLength int, nonModelItems []*NonModelItem, viewIndex int) int {
	viewIndex = lo.Clamp(viewIndex, 0, listLength+len(nonModelItems))
	// Subtract the non-model items that appear before this view index.
	modelIndex := viewIndex
	for i, item := range nonModelItems {
		if item.Index+i < viewIndex {
			modelIndex--
		}
	}
	return modelIndex
}

func (self *ListRenderer) ColumnPositions() []int {
	return self.columnPositions
}

// startIdx and endIdx are view indices, not model indices. If you want to
// render the whole list, pass -1 for both.
func (self *ListRenderer) renderLines(startIdx int, endIdx int) string {
	var columnAlignments []utils.Alignment
	if self.getColumnAlignments != nil {
		columnAlignments = self.getColumnAlignments()
	}
	nonModelItems := self.getNonModelItemList()
	startModelIdx := 0
	if startIdx == -1 {
		startIdx = 0
	} else {
		startModelIdx = viewIndexToModelIndex(self.list.Len(), nonModelItems, startIdx)
	}
	endModelIdx := self.list.Len()
	if endIdx == -1 {
		endIdx = endModelIdx + len(nonModelItems)
	} else {
		endModelIdx = viewIndexToModelIndex(self.list.Len(), nonModelItems, endIdx)
	}
	lines, columnPositions := utils.RenderDisplayStrings(
		self.getDisplayStrings(startModelIdx, endModelIdx),
		columnAlignments)
	self.columnPositions = columnPositions
	lines = self.insertNonModelItems(nonModelItems, endIdx, startIdx, lines, columnPositions)
	return strings.Join(lines, "\n")
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
			padding := ""
			if columnPositions != nil {
				padding = strings.Repeat(" ", columnPositions[item.Column])
			}
			lines = slices.Insert(lines, item.Index+offset-startIdx, padding+item.Content)
		}
		offset++
	}
	return lines
}
