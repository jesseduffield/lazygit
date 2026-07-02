package context

import (
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type MainContext struct {
	*SimpleContext
	*SearchTrait

	diffSelect DiffSelectState
	// dragAnchorViewLine is the view line a mouse-down landed on, remembered so a drag
	// that follows can anchor its range there. The click can show a whole hunk (a range
	// with a different anchor), so the clicked line can't be read back from the view.
	dragAnchorViewLine int
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
// selected line and the range anchor themselves live in the gocui view (its
// cursor and rangeSelectStartY, i.e. native range select); only the mode lives
// here. It's on the context, not the controller, so that the main view controller
// (which drives the selection), the focus controller (which resets it on focus),
// and togglePanel (which sets it on the other pane) can all reach it via the
// context they already hold.
type DiffSelectState struct {
	Mode DiffSelectMode
	// When a range is sticky, moving the cursor without holding shift extends the
	// range; otherwise it collapses the range back to a single line.
	RangeIsSticky bool
	// Whether the user turned on hunk mode explicitly, as opposed to it being the
	// configured default; this decides whether escape leaves hunk mode.
	UserEnabledHunkMode bool
}

// DiffSelectState returns the focused main view's selection mode state, for the
// controllers to read and mutate directly.
func (self *MainContext) DiffSelectState() *DiffSelectState {
	return &self.diffSelect
}

// SetDragAnchorViewLine records the view line a mouse-down landed on, so a drag that
// follows can anchor its range there (see dragAnchorViewLine).
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
				Kind:             types.MAIN_CONTEXT,
				View:             view,
				WindowName:       windowName,
				Key:              key,
				Focusable:        true,
				HighlightOnFocus: false,
			})),
		SearchTrait: NewSearchTrait(c),
	}

	return ctx
}

func (self *MainContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return nil
}

func (self *MainContext) OnSearchSelect(int) {
}
