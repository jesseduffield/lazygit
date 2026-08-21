package context

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainContext struct {
	*SimpleContext
	*SearchTrait

	diffSelect DiffSelectState
	// selectableContentRenderKey names the render whose content HasSelectableContent
	// was worked out from. What there is to select is a property of the content, so an
	// answer about the content of another render says nothing about this one.
	selectableContentRenderKey string
}

var _ types.ISearchableContext = (*MainContext)(nil)

// DiffSelectMode is how the focused main view's diff selection extends from the
// cursor: a single line, a range from a fixed anchor, or the change block (hunk)
// around the cursor.
type DiffSelectMode int

const (
	DiffSelectModeLine DiffSelectMode = iota
	DiffSelectModeRange
	DiffSelectModeHunk
)

// DiffSelectState holds the *mode* of the focused main view's diff selection. The
// selected line and the range anchor themselves live in the gocui view (its cursor
// and range-select start), so only the mode lives here. It's on the context rather
// than on the controller because the controller that drives the selection, the
// controller that establishes it on focus, and the pane-toggle that seeds it on the
// other pane all reach the pane through its context.
type DiffSelectState struct {
	Mode DiffSelectMode
	// When a range is sticky, moving the cursor without holding shift extends the
	// range; otherwise it collapses the range back to a single line.
	RangeIsSticky bool
	// Whether hunk mode was turned on by the user rather than being the configured
	// default. This decides whether escape leaves hunk mode or leaves the view.
	UserEnabledHunkMode bool
}

// DiffSelectState returns the focused main view's selection mode state, for the
// controllers to read and mutate directly.
func (self *MainContext) DiffSelectState() *DiffSelectState {
	return &self.diffSelect
}

// ResetDiffSelectMode returns the pane's selection to the default mode — a single
// line, no range — for whenever it is established from scratch rather than moved. The
// view's range anchor is cleared too, so the next render highlights the cursor line
// only.
func (self *MainContext) ResetDiffSelectMode() {
	self.diffSelect.Mode = DiffSelectModeLine
	self.diffSelect.RangeIsSticky = false
	self.diffSelect.UserEnabledHunkMode = false
	self.GetView().CancelRangeSelect()
}

// SelectableContentRenderKey returns the render HasSelectableContent describes (see
// selectableContentRenderKey).
func (self *MainContext) SelectableContentRenderKey() string {
	return self.selectableContentRenderKey
}

// SetSelectableContentRenderKey records which render HasSelectableContent describes.
func (self *MainContext) SetSelectableContentRenderKey(key string) {
	self.selectableContentRenderKey = key
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

// When selecting a search result, collapse a range selection (whether sticky or not)
// or a hunk selection to just the matching line.
func (self *MainContext) OnSearchSelect(int) {
	self.ResetDiffSelectMode()
}
