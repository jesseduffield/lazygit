package gocui

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func TestMouseReleaseDoesNotBreakDoubleClickDetection(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()
	g := newTestGui(t)
	view, _ := g.SetView("list", 0, 0, 20, 10, 0)
	doubleClicks := []bool{}
	g.SetViewClickBinding(&ViewMouseBinding{
		ViewName: "list",
		Key:      MouseLeft,
		Handler: func(opts ViewMouseBindingOpts) error {
			doubleClicks = append(doubleClicks, opts.IsDoubleClick)
			return nil
		},
	})

	for _, event := range []GocuiEvent{
		gocuiEventFromTcellEvent(tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonPrimary, tcell.ModNone)),
		gocuiEventFromTcellEvent(tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonNone, tcell.ModNone)),
		gocuiEventFromTcellEvent(tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonPrimary, tcell.ModNone)),
	} {
		assert.NoError(t, g.onKey(&event))
	}

	assert.Equal(t, []bool{false, true}, doubleClicks)
}

func TestASwallowedClickIsNoHalfOfADoubleClick(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()
	g := newTestGui(t)
	view, _ := g.SetView("list", 0, 0, 20, 10, 0)
	doubleClicks := []bool{}
	g.SetViewClickBinding(&ViewMouseBinding{
		ViewName: "list",
		Key:      MouseLeft,
		Handler: func(opts ViewMouseBindingOpts) error {
			doubleClicks = append(doubleClicks, opts.IsDoubleClick)
			return nil
		},
	})

	press := gocuiEventFromTcellEvent(
		tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonPrimary, tcell.ModNone))
	release := gocuiEventFromTcellEvent(
		tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonNone, tcell.ModNone))

	// A click the client rejects, as it does for one landing behind a popup panel.
	g.ShouldHandleMouseEvent = func(*View, KeyName) bool { return false }
	assert.NoError(t, g.onKey(&press))
	assert.NoError(t, g.onKey(&release))
	assert.Empty(t, doubleClicks)

	// The same spot clicked again once clicks are accepted. It is a click of its
	// own, not the second half of the one that nothing acted on.
	g.ShouldHandleMouseEvent = func(*View, KeyName) bool { return true }
	assert.NoError(t, g.onKey(&press))

	assert.Equal(t, []bool{false}, doubleClicks)
}

func TestDoubleClicksAreTimedByWhenTheClicksWereRead(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()
	g := newTestGui(t)
	view, _ := g.SetView("list", 0, 0, 20, 10, 0)
	doubleClicks := []bool{}
	g.SetViewClickBinding(&ViewMouseBinding{
		ViewName: "list",
		Key:      MouseLeft,
		Handler: func(opts ViewMouseBindingOpts) error {
			doubleClicks = append(doubleClicks, opts.IsDoubleClick)
			return nil
		},
	})

	click := func(readAt time.Time) {
		for _, buttons := range []tcell.ButtonMask{tcell.ButtonPrimary, tcell.ButtonNone} {
			tev := tcell.NewEventMouse(view.x0+1, view.y0+1, buttons, tcell.ModNone)
			tev.SetEventTime(readAt)
			event := gocuiEventFromTcellEvent(tev)
			assert.NoError(t, g.onKey(&event))
		}
	}

	// The clicks are handled right after each other, long after they were read,
	// as they are when they wait in the queue while the main loop is busy.
	readAt := time.Now().Add(-time.Minute)
	click(readAt)
	click(readAt.Add(DOUBLE_CLICK_THRESHOLD + time.Millisecond))
	click(readAt.Add(DOUBLE_CLICK_THRESHOLD + 100*time.Millisecond))

	assert.Equal(t,
		/* EXPECTED:
		[]bool{false, false, true},
		ACTUAL: */
		[]bool{false, true, true},
		doubleClicks)
}
