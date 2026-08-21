package types

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
	// default, which is what decides whether escape leaves hunk mode or leaves the
	// view.
	UserEnabledHunkMode bool
}
