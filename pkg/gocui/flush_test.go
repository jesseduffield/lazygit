package gocui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestGui(t *testing.T) *Gui {
	t.Helper()
	g, err := NewGui(NewGuiOpts{
		OutputMode: OutputNormal,
		Headless:   true,
		Width:      80,
		Height:     24,
	})
	assert.NoError(t, err)
	t.Cleanup(func() { g.Close() })
	return g
}

// setupViews creates a few views and does an initial full flush so all views
// start in a clean (non-tainted) state.
func setupViews(t *testing.T, g *Gui) (*View, *View) {
	t.Helper()

	status, _ := g.SetView("status", 0, 22, 40, 24, 0)
	status.Frame = false
	main, _ := g.SetView("main", 0, 0, 80, 22, 0)

	// Initial content
	status.SetContent("Ready")
	main.SetContent("hello world")

	// Full flush to draw everything and clear tainted flags
	assert.NoError(t, g.flush())

	return status, main
}

// pushContentOnly pushes a content-only event directly to the channel
// (synchronous, deterministic — unlike Update which spawns a goroutine).
func pushContentOnly(g *Gui, f func(*Gui) error) {
	g.userEvents <- userEvent{f: f, task: g.NewTask(), contentOnly: true}
}

// pushRegular pushes a regular event directly to the channel.
func pushRegular(g *Gui, f func(*Gui) error) {
	g.userEvents <- userEvent{f: f, task: g.NewTask(), contentOnly: false}
}

func TestFlushContentOnly_SkipsUntaintedViews(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	// After initial flush, both views should be untainted
	assert.False(t, status.IsTainted(), "status view should not be tainted after flush")
	assert.False(t, main.IsTainted(), "main view should not be tainted after flush")

	// Modify only the status view
	status.SetContent("Fetching /")

	assert.True(t, status.IsTainted(), "status view should be tainted after SetContent")
	assert.False(t, main.IsTainted(), "main view should not be tainted (was not modified)")

	// flushContentOnly should succeed and clear status tainted flag
	assert.NoError(t, g.flushContentOnly(g.views))

	assert.False(t, status.IsTainted(), "status view should not be tainted after flushContentOnly")
	assert.False(t, main.IsTainted(), "main view should not be tainted after flushContentOnly")
}

func TestFlushContentOnly_WritesCorrectContent(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	status.SetContent("Fetching |")
	assert.NoError(t, g.flushContentOnly(g.views))

	assert.Equal(t, "Fetching |", status.Buffer())
}

func TestProcessEvent_ContentOnlyEvent_SkipsTaintedCheck(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	// Send a content-only event that modifies only the status view
	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("Fetching /")
		return nil
	})

	assert.NoError(t, g.processEvent())

	// status was modified and drawn → tainted cleared
	assert.False(t, status.IsTainted(), "status should not be tainted after processEvent with contentOnly")
	// main was NOT modified → should still be untainted
	assert.False(t, main.IsTainted(), "main should not be tainted after processEvent with contentOnly")
}

func TestProcessEvent_RegularEvent_UsesFullFlush(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	// Regular event (not content-only) should trigger full flush
	pushRegular(g, func(gui *Gui) error {
		status.SetContent("Fetching \\")
		return nil
	})

	assert.NoError(t, g.processEvent())

	assert.False(t, status.IsTainted(), "status should not be tainted after full flush")
}

func TestProcessEvent_MixedBatch_UsesFullFlush(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	// Queue a content-only event followed by a regular event.
	// processEvent picks up the first; processRemainingEvents picks up
	// the second. Since the second is not contentOnly, full flush runs.
	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("Fetching -")
		return nil
	})
	pushRegular(g, func(gui *Gui) error {
		main.SetContent("updated main")
		return nil
	})

	assert.NoError(t, g.processEvent())

	// Both views were modified and should have been drawn by full flush
	assert.False(t, status.IsTainted(), "status should not be tainted after full flush")
	assert.False(t, main.IsTainted(), "main should not be tainted after full flush")
}

func TestProcessEvent_RegularThenContentOnly_UsesFullFlush(t *testing.T) {
	g := newTestGui(t)
	status, main := setupViews(t, g)

	// Even if a regular event comes first and the remaining are contentOnly,
	// the batch must use full flush.
	pushRegular(g, func(gui *Gui) error {
		main.SetContent("new main content")
		return nil
	})
	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("Fetching |")
		return nil
	})

	assert.NoError(t, g.processEvent())

	assert.False(t, status.IsTainted(), "status should not be tainted after full flush")
	assert.False(t, main.IsTainted(), "main should not be tainted after full flush")
}

func TestProcessRemainingEvents_AllContentOnly_ReturnsTrue(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("a")
		return nil
	})
	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("b")
		return nil
	})

	contentOnly, err := g.processRemainingEvents()
	assert.NoError(t, err)
	assert.True(t, contentOnly, "should return true when all events are contentOnly")
}

func TestProcessRemainingEvents_MixedEvents_ReturnsFalse(t *testing.T) {
	g := newTestGui(t)
	status, _ := setupViews(t, g)

	pushContentOnly(g, func(gui *Gui) error {
		status.SetContent("a")
		return nil
	})
	pushRegular(g, func(gui *Gui) error {
		status.SetContent("b")
		return nil
	})

	contentOnly, err := g.processRemainingEvents()
	assert.NoError(t, err)
	assert.False(t, contentOnly, "should return false when any event is not contentOnly")
}

func TestProcessRemainingEvents_EmptyQueue_ReturnsTrue(t *testing.T) {
	g := newTestGui(t)

	contentOnly, err := g.processRemainingEvents()
	assert.NoError(t, err)
	assert.True(t, contentOnly, "should return true when no events are queued")
}

// Ensure an overlapping view that is not tainted does not get overdrawn
func TestFlushContentOnly_DoesNotOverdrawHigherZViews(t *testing.T) {
	g := newTestGui(t)

	// Base view
	list, _ := g.SetView("list", 0, 0, 79, 23, 0)
	list.Frame = false
	list.SetContent(strings.Repeat("LIST LINE FILLER FILLER FILLER FILLER FILLER FILLER FILLER FILLER FILLER\n", 22))

	// Overlapping 'popup'
	popup, _ := g.SetView("popup", 20, 8, 60, 16, 0)
	popup.Frame = false
	popupLine := strings.Repeat("P", 60)
	popup.SetContent(strings.Repeat(popupLine+"\n", 16))

	// Full flush — popup ends up on top.
	assert.NoError(t, g.flush())

	cellAt := func(x, y int) string {
		s, _, _ := g.screen.Get(x, y)
		return s
	}

	// Taint only the list view
	list.SetContent(strings.Repeat(strings.Repeat("X", 80)+"\n", 22))
	assert.True(t, list.IsTainted(), "list should be tainted after SetContent")
	assert.False(t, popup.IsTainted(), "popup should not be tainted")

	// flushContentOnly is what spinner ticks ultimately invoke.
	assert.NoError(t, g.flushContentOnly(g.views))

	assert.Equal(t, "P", cellAt(21, 9),
		"popup region must still show popup content after flushContentOnly; "+
			"if this fails the popup-overdraw bug is present")

	// Additional checks to be sure
	assert.Equal(t, "P", cellAt(40, 11), "interior popup cell should still show popup content")
	assert.Equal(t, "P", cellAt(58, 14), "near-edge popup cell should still show popup content")

	// Ensure tainted view was updated
	assert.Equal(t, "X", cellAt(5, 5), "list cell outside popup should show new list content")
	assert.Equal(t, "X", cellAt(70, 20), "list cell outside popup should show new list content")
}

// Ensure transitive overlap: with views in z-order [a, b, c] where b overlaps a
// and c overlaps b but c does NOT overlap a, tainting a must redraw all three —
// otherwise b's redraw paints over c.
func TestFlushContentOnly_RedrawsTransitivelyOverlappingViews(t *testing.T) {
	g := newTestGui(t)

	// Geometry: b straddles a and c; a and c are disjoint.
	//   a: (0,0)-(40,10)    b: (30,5)-(60,15)    c: (50,12)-(75,20)
	a, _ := g.SetView("a", 0, 0, 40, 10, 0)
	a.Frame = false
	a.SetContent(strings.Repeat(strings.Repeat("A", 60)+"\n", 20))

	b, _ := g.SetView("b", 30, 5, 60, 15, 0)
	b.Frame = false
	b.SetContent(strings.Repeat(strings.Repeat("B", 60)+"\n", 20))

	c, _ := g.SetView("c", 50, 12, 75, 20, 0)
	c.Frame = false
	c.SetContent(strings.Repeat(strings.Repeat("C", 60)+"\n", 20))

	assert.NoError(t, g.flush())

	cellAt := func(x, y int) string {
		s, _, _ := g.screen.Get(x, y)
		return s
	}

	// Taint only a.
	a.SetContent(strings.Repeat(strings.Repeat("X", 60)+"\n", 20))
	assert.True(t, a.IsTainted())
	assert.False(t, b.IsTainted())
	assert.False(t, c.IsTainted())

	assert.NoError(t, g.flushContentOnly(g.views))

	// a redrawn (direct).
	assert.Equal(t, "X", cellAt(5, 5), "a should be redrawn (tainted)")
	// b redrawn (overlaps a).
	assert.Equal(t, "B", cellAt(45, 7), "b should be redrawn (overlaps a)")
	// c redrawn transitively (overlaps b, which overlaps a). Without the
	// transitive case, b's redraw would paint over c at this cell.
	assert.Equal(t, "C", cellAt(55, 14),
		"c should be redrawn transitively; if 'B' here, b's redraw painted over c")
}
