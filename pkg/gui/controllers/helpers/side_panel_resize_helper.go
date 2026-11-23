package helpers

import (
	"maps"
	"time"

	"github.com/jesseduffield/lazycore/pkg/boxlayout"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

// SidePanelResizeHelper lets the user resize the side panels by dragging the
// boundary between them and the main view with the mouse. The boundary consists
// of the edges of the side panels that face the main view and the edges of the
// main view and the command log that face the side panels, so it is two columns
// wide (or two rows high if the side panels are above the main view).
type SidePanelResizeHelper struct {
	c                  *HelperCommon
	windowHelper       *WindowHelper
	confirmationHelper *ConfirmationHelper

	// The boundary as of the last layout; nil if only the side panels or only
	// the main view are visible
	boundary *sidePanelBoundary

	// The last known position of the mouse pointer
	pointerKnown       bool
	pointerX, pointerY int
	hovering           bool

	dragging   bool
	dragLayout config.SidePanelLayout
	// The size of the side panels that the drag has set, and the screen size
	// it refers to; 0 if the drag hasn't changed the size yet
	dragSize, dragScreenSize int
	// The sizes for dragLayout before the drag started
	sizesBeforeDrag []config.SidePanelSize
	// Where in the boundary the drag started: 0 for the column (or row) of the
	// side panels, 1 for the one of the main view
	grabOffset int
	// When the user last pressed the mouse button on the boundary without
	// dragging it; zero if the press doesn't count towards a double click
	lastClickTime time.Time
}

type sidePanelBoundary struct {
	arrangement SidePanelArrangement
	// The last column (or row) of the side panels. The boundary consists of
	// this one and the next one.
	pos int
	// The first and the last row (or column) of the boundary
	start, end int
	sideViews  []*gocui.View
	mainViews  []*gocui.View
}

func NewSidePanelResizeHelper(
	c *HelperCommon,
	windowHelper *WindowHelper,
	confirmationHelper *ConfirmationHelper,
) *SidePanelResizeHelper {
	return &SidePanelResizeHelper{
		c:                  c,
		windowHelper:       windowHelper,
		confirmationHelper: confirmationHelper,
	}
}

// coordinates returns the given screen position as a coordinate across the
// boundary and one along it.
func (self *sidePanelBoundary) coordinates(x, y int) (across int, along int) {
	if self.arrangement.SidePanelsAbove {
		return y, x
	}
	return x, y
}

func (self *sidePanelBoundary) contains(x, y int) bool {
	across, along := self.coordinates(x, y)
	return (across == self.pos || across == self.pos+1) && along >= self.start && along <= self.end
}

// IsActive returns true while the mouse hovers over the boundary or drags it.
func (self *SidePanelResizeHelper) IsActive() bool {
	return self.boundary != nil && (self.hovering || self.dragging)
}

// Layout finds the boundary after the views have been laid out, and highlights
// it while the mouse hovers over it or drags it.
func (self *SidePanelResizeHelper) Layout(
	windowDimensions map[string]boxlayout.Dimensions,
	arrangement SidePanelArrangement,
	hasArrangement bool,
) {
	self.boundary = nil
	if hasArrangement {
		self.boundary = self.findBoundary(windowDimensions, arrangement)
	}
	self.updateHovering()
	self.highlightBoundary()
}

func (self *SidePanelResizeHelper) OnMouseMove(x, y int) {
	self.pointerKnown = true
	self.pointerX, self.pointerY = x, y
	self.updateHovering()
}

// OnFocusLost stops highlighting the boundary. The terminal doesn't report it
// when the pointer leaves its window, so if the pointer leaves the window right
// from the boundary, the boundary stays highlighted until the window loses focus
// or the pointer comes back.
func (self *SidePanelResizeHelper) OnFocusLost() {
	self.pointerKnown = false
	self.updateHovering()
}

func (self *SidePanelResizeHelper) updateHovering() {
	self.hovering = self.pointerKnown &&
		self.boundary != nil &&
		!self.confirmationHelper.IsPopupPanelFocused() &&
		self.boundary.contains(self.pointerX, self.pointerY)
}

// InterceptMouse handles the mouse events of gestures that start on the
// boundary; see gocui.Gui.InterceptMouse.
func (self *SidePanelResizeHelper) InterceptMouse(ev gocui.MouseEvent) bool {
	switch {
	case ev.Key == gocui.MouseLeft && ev.Mod == gocui.ModNone:
		return self.onPress(ev)
	case ev.Key == gocui.MouseLeft && ev.Mod == gocui.ModMotion:
		return self.onDrag(ev)
	case ev.Key == gocui.MouseRelease:
		return self.onRelease(ev)
	}
	return false
}

func (self *SidePanelResizeHelper) onPress(ev gocui.MouseEvent) bool {
	boundary := self.boundary
	if boundary == nil ||
		self.confirmationHelper.IsPopupPanelFocused() ||
		!boundary.contains(ev.X, ev.Y) {
		return false
	}

	self.pointerKnown = true
	self.pointerX, self.pointerY = ev.X, ev.Y
	layout := boundary.arrangement.Layout

	if !self.lastClickTime.IsZero() && ev.When.Sub(self.lastClickTime) < gocui.DOUBLE_CLICK_THRESHOLD {
		self.lastClickTime = time.Time{}
		self.resetSizes(layout)
		return true
	}

	across, _ := boundary.coordinates(ev.X, ev.Y)
	self.dragging = true
	self.dragSize, self.dragScreenSize = 0, 0
	self.dragLayout = layout
	self.sizesBeforeDrag = self.c.State().GetSidePanelSizes()[layout]
	self.grabOffset = across - boundary.pos
	self.lastClickTime = ev.When
	return true
}

func (self *SidePanelResizeHelper) onDrag(ev gocui.MouseEvent) bool {
	if !self.dragging {
		return false
	}

	self.pointerX, self.pointerY = ev.X, ev.Y
	boundary := self.boundary
	if boundary == nil || boundary.arrangement.Layout != self.dragLayout {
		// The layout changed under the drag, e.g. because a refresh split the
		// main view. Ignore the rest of the gesture.
		return true
	}

	// Several events can be handled before the next layout, so the boundary
	// may not show the size of the previous drag event yet.
	currentSize := boundary.pos + 1
	if self.dragSize != 0 {
		currentSize = self.dragSize
	}
	across, _ := boundary.coordinates(ev.X, ev.Y)
	screenSize := boundary.arrangement.ScreenSize
	size := clampSideSectionSize(across-self.grabOffset+1, screenSize)
	if size == currentSize {
		return true
	}

	self.dragSize, self.dragScreenSize = size, screenSize
	self.lastClickTime = time.Time{}
	self.setSizes(self.dragLayout, recordSidePanelSize(self.sizesBeforeDrag, screenSize, size))
	return true
}

func (self *SidePanelResizeHelper) onRelease(ev gocui.MouseEvent) bool {
	if !self.dragging {
		return false
	}

	self.dragging = false
	self.pointerX, self.pointerY = ev.X, ev.Y
	self.updateHovering()

	if self.dragSize != 0 {
		layout, screenSize, size := self.dragLayout, self.dragScreenSize, self.dragSize
		self.c.UpdateAppState(func(state *config.AppState) {
			if state.SidePanelSizes == nil {
				state.SidePanelSizes = map[config.SidePanelLayout][]config.SidePanelSize{}
			}
			sizes := ValidSidePanelSizes(state.SidePanelSizes)[layout]
			state.SidePanelSizes[layout] = recordSidePanelSize(sizes, screenSize, size)
		})
	}
	return true
}

// resetSizes forgets the sizes that the user set for the given layout, so that
// the default applies again.
func (self *SidePanelResizeHelper) resetSizes(layout config.SidePanelLayout) {
	self.setSizes(layout, nil)
	self.c.UpdateAppState(func(state *config.AppState) {
		delete(state.SidePanelSizes, layout)
	})
}

func (self *SidePanelResizeHelper) setSizes(layout config.SidePanelLayout, sizes []config.SidePanelSize) {
	sizesByLayout := maps.Clone(self.c.State().GetSidePanelSizes())
	if sizesByLayout == nil {
		sizesByLayout = map[config.SidePanelLayout][]config.SidePanelSize{}
	}
	if len(sizes) == 0 {
		delete(sizesByLayout, layout)
	} else {
		sizesByLayout[layout] = sizes
	}
	self.c.State().SetSidePanelSizes(sizesByLayout)
}

type windowView struct {
	dimensions boxlayout.Dimensions
	view       *gocui.View
}

// visibleWindowViews returns the views shown in the given windows, if the
// windows take up any space.
func (self *SidePanelResizeHelper) visibleWindowViews(
	windowDimensions map[string]boxlayout.Dimensions,
	windows []string,
) []windowView {
	result := []windowView{}
	for _, window := range windows {
		dimensions, ok := windowDimensions[window]
		if !ok || dimensions.X1 < dimensions.X0 || dimensions.Y1 < dimensions.Y0 {
			continue
		}
		view, err := self.c.GocuiGui().View(self.windowHelper.GetViewNameForWindow(window))
		if err != nil || !view.Visible {
			continue
		}
		result = append(result, windowView{dimensions: dimensions, view: view})
	}
	return result
}

func (self *SidePanelResizeHelper) findBoundary(
	windowDimensions map[string]boxlayout.Dimensions,
	arrangement SidePanelArrangement,
) *sidePanelBoundary {
	sideWindowViews := self.visibleWindowViews(windowDimensions, self.windowHelper.SideWindows())
	mainWindowViews := self.visibleWindowViews(windowDimensions, []string{"main", "secondary", "extras"})
	if len(sideWindowViews) == 0 || len(mainWindowViews) == 0 {
		return nil
	}

	// The coordinates of a window across the boundary, and along it
	extent := func(d boxlayout.Dimensions) (across0, across1, along0, along1 int) {
		if arrangement.SidePanelsAbove {
			return d.Y0, d.Y1, d.X0, d.X1
		}
		return d.X0, d.X1, d.Y0, d.Y1
	}

	boundary := &sidePanelBoundary{arrangement: arrangement, pos: -1, start: -1, end: -1}
	for _, w := range sideWindowViews {
		_, across1, _, _ := extent(w.dimensions)
		boundary.pos = max(boundary.pos, across1)
	}
	addAlong := func(along0, along1 int) {
		if boundary.start == -1 || along0 < boundary.start {
			boundary.start = along0
		}
		boundary.end = max(boundary.end, along1)
	}
	for _, w := range sideWindowViews {
		_, across1, along0, along1 := extent(w.dimensions)
		if across1 == boundary.pos {
			boundary.sideViews = append(boundary.sideViews, w.view)
			addAlong(along0, along1)
		}
	}
	for _, w := range mainWindowViews {
		across0, _, along0, along1 := extent(w.dimensions)
		if across0 == boundary.pos+1 {
			boundary.mainViews = append(boundary.mainViews, w.view)
			addAlong(along0, along1)
		}
	}
	if len(boundary.mainViews) == 0 {
		return nil
	}
	return boundary
}

func (self *SidePanelResizeHelper) highlightBoundary() {
	for _, view := range self.c.GocuiGui().Views() {
		view.EdgeHighlight = nil
	}
	if !self.IsActive() {
		return
	}

	for _, view := range self.boundary.sideViews {
		view.EdgeHighlight = self.edgeHighlight(true)
	}
	for _, view := range self.boundary.mainViews {
		view.EdgeHighlight = self.edgeHighlight(false)
	}
}

// edgeHighlight returns the highlight for the edge of a side view or a main view
// that is part of the boundary. It draws the edge with heavy lines, and uses
// corners that join them to the light lines of the rest of the frame. There are
// no such corners for double or bold lines, so with these borders it only
// changes the color of the edge.
func (self *SidePanelResizeHelper) edgeHighlight(isSideView bool) *gocui.EdgeHighlight {
	highlight := &gocui.EdgeHighlight{Color: theme.ActiveBorderColor}
	above := self.boundary.arrangement.SidePanelsAbove
	switch {
	case !above && isSideView:
		highlight.Edge = gocui.EdgeRight
	case !above:
		highlight.Edge = gocui.EdgeLeft
	case isSideView:
		highlight.Edge = gocui.EdgeBottom
	default:
		highlight.Edge = gocui.EdgeTop
	}

	switch self.c.UserConfig().Gui.Border {
	case "single", "rounded":
		switch highlight.Edge {
		case gocui.EdgeRight:
			highlight.Rune, highlight.StartCorner, highlight.EndCorner = '┃', '┒', '┚'
		case gocui.EdgeLeft:
			highlight.Rune, highlight.StartCorner, highlight.EndCorner = '┃', '┎', '┖'
		case gocui.EdgeBottom:
			highlight.Rune, highlight.StartCorner, highlight.EndCorner = '━', '┕', '┙'
		case gocui.EdgeTop:
			highlight.Rune, highlight.StartCorner, highlight.EndCorner = '━', '┍', '┑'
		}
	case "hidden":
		// The frame is invisible, so draw just the edge
		highlight.Rune = '┃'
		if above {
			highlight.Rune = '━'
		}
		highlight.StartCorner, highlight.EndCorner = highlight.Rune, highlight.Rune
	}
	return highlight
}
