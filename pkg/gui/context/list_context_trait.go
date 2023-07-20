package context

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ListContextTrait struct {
	types.Context

	c                 *ContextCommon
	list              types.IList
	getDisplayStrings func(startIdx int, length int) [][]string
	// Alignment for each column. If nil, the default is left alignment
	getColumnAlignments func() []utils.Alignment
	// Some contexts, like the commit context, will highlight the path from the selected commit
	// to its parents, because it's ambiguous otherwise. For these, we need to refresh the viewport
	// so that we show the highlighted path.
	// TODO: now that we allow scrolling, we should be smarter about what gets refreshed:
	// we should find out exactly which lines are now part of the path and refresh those.
	// We should also keep track of the previous path and refresh those lines too.
	refreshViewportOnChange bool
}

func (self *ListContextTrait) IsListContext() {}

func (self *ListContextTrait) GetList() types.IList {
	return self.list
}

func (self *ListContextTrait) FocusLine() {
	// Doing this at the end of the layout function because we need the view to be
	// resized before we focus the line, otherwise if we're in accordion mode
	// the view could be squashed and won't how to adjust the cursor/origin
	self.c.AfterLayout(func() error {
		self.GetViewTrait().FocusPoint(self.list.GetSelectedLineIdx())
		return nil
	})

	self.setFooter()

	if self.refreshViewportOnChange {
		self.refreshViewport()
	}
}

func (self *ListContextTrait) refreshViewport() {
	startIdx, length := self.GetViewTrait().ViewPortYBounds()
	displayStrings := self.getDisplayStrings(startIdx, length)
	content := utils.RenderDisplayStrings(displayStrings, nil)
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
	self.list.RefreshSelectedIdx()
	var columnAlignments []utils.Alignment
	if self.getColumnAlignments != nil {
		columnAlignments = self.getColumnAlignments()
	}
	content := utils.RenderDisplayStrings(
		self.getDisplayStrings(0, self.list.Len()),
		columnAlignments,
	)
	self.GetViewTrait().SetContent(content)
	self.c.Render()
	self.setFooter()

	return nil
}

func (self *ListContextTrait) OnSearchSelect(selectedLineIdx int) error {
	self.GetList().SetSelectedLineIdx(selectedLineIdx)
	return self.HandleFocus(types.OnFocusOpts{})
}
