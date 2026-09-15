package helpers

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const (
	dragAutoscrollInitialDelay     = 300 * time.Millisecond
	dragAutoscrollSlowInterval     = 250 * time.Millisecond
	dragAutoscrollFastInterval     = 100 * time.Millisecond
	dragAutoscrollVeryFastInterval = 50 * time.Millisecond
)

// All state is UI-thread-owned. Timer goroutines only enqueue tick back onto
// the UI thread, where generation changes and scroll callbacks are serialized
// with mouse handlers and focus changes.
type DragAutoscroller struct {
	c       *HelperCommon
	context types.Context

	canScroll func(direction int) bool
	onScroll  func(viewIndex int) bool

	// Incremented whenever the scroll direction changes or the autoscroller
	// is canceled. A scheduled tick carries the generation it was created
	// for, so stale ticks can be told apart from the one that is current.
	generation uint64
	direction  int
	interval   time.Duration
	// Last known pointer position relative to the viewport; used by ticks to
	// compute which line ends up under the pointer after scrolling.
	pointerViewportY int
}

func NewDragAutoscroller(
	c *HelperCommon,
	context types.Context,
	canScroll func(direction int) bool,
	onScroll func(viewIndex int) bool,
) *DragAutoscroller {
	return &DragAutoscroller{
		c:         c,
		context:   context,
		canScroll: canScroll,
		onScroll:  onScroll,
	}
}

// Update is called with the pointer position of every drag event. Entering a
// scroll zone arms a timer (with an initial delay, so that merely passing
// through the zone doesn't scroll); once armed, scrolling continues on its
// own until the pointer leaves the zone, the drag ends, or a callback stops
// it.
func (self *DragAutoscroller) Update(pointerViewportY int) {
	_, viewportHeight := self.context.GetViewTrait().ViewPortYBounds()
	direction, interval := dragAutoscrollZone(viewportHeight, pointerViewportY)
	if direction != 0 && self.canScroll != nil && !self.canScroll(direction) {
		direction = 0
		interval = 0
	}

	self.pointerViewportY = pointerViewportY
	generation, schedule := self.updateState(direction, interval)
	if schedule {
		self.schedule(generation, dragAutoscrollInitialDelay)
	}
}

func (self *DragAutoscroller) Direction() int {
	return self.direction
}

func (self *DragAutoscroller) updateState(direction int, interval time.Duration) (uint64, bool) {
	if direction == self.direction {
		self.interval = interval
		return self.generation, false
	}

	self.generation++
	self.direction = direction
	self.interval = interval
	return self.generation, direction != 0
}

func (self *DragAutoscroller) Cancel() {
	self.generation++
	self.direction = 0
	self.interval = 0
}

func (self *DragAutoscroller) schedule(generation uint64, delay time.Duration) {
	time.AfterFunc(delay, func() {
		self.c.OnUIThreadBackground(func() error {
			self.tick(generation)
			return nil
		})
	})
}

func (self *DragAutoscroller) tick(generation uint64) {
	if generation != self.generation {
		return
	}
	if self.direction == 0 ||
		self.canScroll != nil && !self.canScroll(self.direction) {
		self.Cancel()
		return
	}

	view := self.context.GetViewTrait()
	oldOriginY, _ := view.ViewPortYBounds()
	if self.direction < 0 {
		view.ScrollUp(1)
	} else {
		view.ScrollDown(1)
	}
	newOriginY, _ := view.ViewPortYBounds()
	if newOriginY == oldOriginY {
		self.Cancel()
		return
	}

	if !self.onScroll(newOriginY + self.pointerViewportY) {
		self.Cancel()
		return
	}

	self.schedule(generation, self.interval)
}

// dragAutoscrollZone returns the scroll direction and tick interval for a
// pointer position: anything beyond the view scrolls very fast, the outermost
// viewport row scrolls fast, the row just inside it scrolls slowly, and anything
// further inside doesn't scroll at all.
func dragAutoscrollZone(viewportHeight int, pointerViewportY int) (int, time.Duration) {
	switch {
	case pointerViewportY < 0:
		return -1, dragAutoscrollVeryFastInterval
	case pointerViewportY == 0:
		return -1, dragAutoscrollFastInterval
	case pointerViewportY == 1:
		return -1, dragAutoscrollSlowInterval
	case pointerViewportY > viewportHeight-1:
		return 1, dragAutoscrollVeryFastInterval
	case pointerViewportY == viewportHeight-1:
		return 1, dragAutoscrollFastInterval
	case pointerViewportY == viewportHeight-2:
		return 1, dragAutoscrollSlowInterval
	default:
		return 0, 0
	}
}
