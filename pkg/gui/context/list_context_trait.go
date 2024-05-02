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
}

func (self *ListContextTrait) IsListContext() {}

func (self *ListContextTrait) FocusLine() {
	// Doing this at the end of the layout function because we need the view to be
	// resized before we focus the line, otherwise if we're in accordion mode
	// the view could be squashed and won't how to adjust the cursor/origin
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

		// If FocusPoint() caused the view to scroll (because the selected line
		// was out of view before), we need to rerender the view port again.
		// This can happen when pressing , or . to scroll by pages, or < or > to
		// jump to the top or bottom.
		newOrigin, _ := self.GetViewTrait().ViewPortYBounds()
		if self.refreshViewportOnChange && oldOrigin != newOrigin {
			self.refreshViewport()
		}
		return nil
	})

	self.setFooter()

	if self.refreshViewportOnChange {
		self.refreshViewport()
	}
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

func (self *ListContextTrait) HandleFocus(opts types.OnFocusOpts) error {
	self.FocusLine()

	self.GetViewTrait().SetHighlight(self.list.Len() > 0)

	return self.Context.HandleFocus(opts)
}

func (self *ListContextTrait) HandleFocusLost(opts types.OnFocusLostOpts) error {
	self.GetViewTrait().SetOriginX(0)

	if self.refreshViewportOnChange {
		self.refreshViewport()
	}

	return self.Context.HandleFocusLost(opts)
}

// OnFocus assumes that the content of the context has already been rendered to the view. OnRender is the function which actually renders the content to the view
func (self *ListContextTrait) HandleRender() error {
	self.list.ClampSelection()
	content := self.renderLines(-1, -1)
	self.GetViewTrait().SetContent(content)
	self.c.Render()
	self.setFooter()

	return nil
}

func (self *ListContextTrait) OnSearchSelect(selectedLineIdx int) error {
	self.GetList().SetSelection(selectedLineIdx)
	return self.HandleFocus(types.OnFocusOpts{})
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

// By default, list contexts supporta range select
func (self *ListContextTrait) RangeSelectEnabled() bool {
	return true
}
