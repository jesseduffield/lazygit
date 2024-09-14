package context

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ListContextTrait struct {
	types.Context
	ListRenderer

	c *ContextCommon
	// Some contexts, like the commit context, will highlight the path from the selected commit
	// to its parents, because it's ambiguous otherwise. For these, we need to refresh the viewport
	// so that we show the highlighted path.
	// TODO: now that we allow scrolling, we should be smarter about what gets refreshed:
	// we should find out exactly which lines are now part of the path and refresh those.
	// We should also keep track of the previous path and refresh those lines too.
	refreshViewportOnChange bool
	// If this is true, we only render the visible lines of the list. Useful for lists that can
	// get very long, because it can save a lot of memory
	renderOnlyVisibleLines bool
}

func (self *ListContextTrait) IsListContext() {}

func (self *ListContextTrait) FocusLine() {
	// Doing this at the end of the layout function because we need the view to be
	// resized before we focus the line, otherwise if we're in accordion mode
	// the view could be squashed and won't how to adjust the cursor/origin.
	// Also, refreshing the viewport needs to happen after the view has been resized.
	self.c.AfterLayout(func() error {
		oldOrigin, _ := self.GetViewTrait().ViewPortYBounds()

		self.GetViewTrait().FocusPoint(
			self.ModelIndexToViewIndex(self.list.GetSelectedLineIdx()))

		selectRangeIndex, isSelectingRange := self.list.GetRangeStartIdx()
		if isSelectingRange {
			selectRangeIndex = self.ModelIndexToViewIndex(selectRangeIndex)
			self.GetViewTrait().SetRangeSelectStart(selectRangeIndex)
		} else {
			self.GetViewTrait().CancelRangeSelect()
		}

		if self.refreshViewportOnChange {
			self.refreshViewport()
		} else if self.renderOnlyVisibleLines {
			newOrigin, _ := self.GetViewTrait().ViewPortYBounds()
			if oldOrigin != newOrigin {
				self.HandleRender()
			}
		}
		return nil
	})

	self.setFooter()
}

func (self *ListContextTrait) refreshViewport() {
	startIdx, length := self.GetViewTrait().ViewPortYBounds()
	content := self.renderLines(startIdx, startIdx+length)
	self.GetViewTrait().SetViewPortContent(content)
}

func (self *ListContextTrait) setFooter() {
	self.GetViewTrait().SetFooter(formatListFooter(self.list.GetSelectedLineIdx(), self.list.Len()))
}

func formatListFooter(selectedLineIdx int, length int) string {
	return fmt.Sprintf("%d of %d", selectedLineIdx+1, length)
}

func (self *ListContextTrait) HandleFocus(opts types.OnFocusOpts) {
	self.FocusLine()

	self.GetViewTrait().SetHighlight(self.list.Len() > 0)

	self.Context.HandleFocus(opts)
}

func (self *ListContextTrait) HandleFocusLost(opts types.OnFocusLostOpts) {
	self.GetViewTrait().SetOriginX(0)

	if self.refreshViewportOnChange {
		self.refreshViewport()
	}

	self.Context.HandleFocusLost(opts)
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (self *ListContextTrait) HandleRender() {
	self.list.ClampSelection()
	if self.renderOnlyVisibleLines {
		// Rendering only the visible area can save a lot of cell memory for
		// those views that support it.
		totalLength := self.list.Len()
		if self.getNonModelItems != nil {
			totalLength += len(self.getNonModelItems())
		}
		self.GetViewTrait().SetContentLineCount(totalLength)
		startIdx, length := self.GetViewTrait().ViewPortYBounds()
		content := self.renderLines(startIdx, startIdx+length)
		self.GetViewTrait().SetViewPortContentAndClearEverythingElse(content)
	} else {
		content := self.renderLines(-1, -1)
		self.GetViewTrait().SetContent(content)
	}
	self.c.Render()
	self.setFooter()
}

func (self *ListContextTrait) OnSearchSelect(selectedLineIdx int) error {
	self.GetList().SetSelection(self.ViewIndexToModelIndex(selectedLineIdx))
	self.HandleFocus(types.OnFocusOpts{})
	return nil
}

func (self *ListContextTrait) IsItemVisible(item types.HasUrn) bool {
	startIdx, length := self.GetViewTrait().ViewPortYBounds()
	selectionStart := self.ViewIndexToModelIndex(startIdx)
	selectionEnd := self.ViewIndexToModelIndex(startIdx + length)
	for i := selectionStart; i < selectionEnd; i++ {
		iterItem := self.GetList().GetItem(i)
		if iterItem != nil && iterItem.URN() == item.URN() {
			return true
		}
	}
	return false
}

// By default, list contexts supports range select
func (self *ListContextTrait) RangeSelectEnabled() bool {
	return true
}

func (self *ListContextTrait) RenderOnlyVisibleLines() bool {
	return self.renderOnlyVisibleLines
}

func (self *ListContextTrait) TotalContentHeight() int {
	result := self.list.Len()
	if self.getNonModelItems != nil {
		result += len(self.getNonModelItems())
	}
	return result
}
