package context

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainContext struct {
	*SimpleContext
	*SearchTrait

	diffSelect types.DiffSelectState
	// dragAnchorViewLine is the view line a mouse-down landed on, remembered so that a
	// drag that follows can anchor its range there. The click may have selected a whole
	// hunk, whose range anchor is the block's far end, so the clicked line can't be
	// read back from the view.
	dragAnchorViewLine int
}

var (
	_ types.ISearchableContext = (*MainContext)(nil)
	_ types.DiffPaneContext    = (*MainContext)(nil)
)

// DiffSelectState returns the focused main view's selection mode state, for the
// controllers to read and mutate directly.
func (self *MainContext) DiffSelectState() *types.DiffSelectState {
	return &self.diffSelect
}

// ResetDiffSelectMode returns the pane's selection to the default mode — a single
// line, no range — for whenever it is established from scratch rather than moved. The
// view's range anchor is cleared too, so the next render highlights the cursor line
// only.
func (self *MainContext) ResetDiffSelectMode() {
	self.diffSelect.Mode = types.DiffSelectModeLine
	self.diffSelect.RangeIsSticky = false
	self.diffSelect.UserEnabledHunkMode = false
	self.GetView().CancelRangeSelect()
}

// SetDragAnchorViewLine records the view line a mouse-down landed on, so that a drag
// that follows can anchor its range there (see dragAnchorViewLine).
func (self *MainContext) SetDragAnchorViewLine(viewLine int) {
	self.dragAnchorViewLine = viewLine
}

// DragAnchorViewLine returns the view line the last mouse-down landed on.
func (self *MainContext) DragAnchorViewLine() int {
	return self.dragAnchorViewLine
}

func NewMainContext(
	view *gocui.View,
	windowName string,
	key types.ContextKey,
	c *ContextCommon,
) *MainContext {
	ctx := &MainContext{
		SimpleContext: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:       types.MAIN_CONTEXT,
				View:       view,
				WindowName: windowName,
				Key:        key,
				Focusable:  true,
			})),
		SearchTrait: NewSearchTrait(c),
	}

	return ctx
}

func (self *MainContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return nil
}

// OnSearchSelect collapses the selection to the match the search has just moved the
// cursor to. A range or a selected hunk means "these lines here", which a jump to a
// match somewhere else in the diff has nothing to do with: extending the range to the
// match, or holding on to a hunk the cursor has left, would both leave the user
// looking at a selection they didn't make.
func (self *MainContext) OnSearchSelect(int) {
	self.ResetDiffSelectMode()
}
