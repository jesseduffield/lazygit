package gocui

import (
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func TestMouseReleaseDoesNotBreakDoubleClickDetection(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()
	g := newTestGui(t)
	view, _ := g.SetView("list", 0, 0, 20, 10, 0)
	doubleClicks := []bool{}
	assert.NoError(t, g.SetViewClickBinding(&ViewMouseBinding{
		ViewName: "list",
		Key:      MouseLeft,
		Handler: func(opts ViewMouseBindingOpts) error {
			doubleClicks = append(doubleClicks, opts.IsDoubleClick)
			return nil
		},
	}))

	for _, event := range []GocuiEvent{
		gocuiEventFromTcellEvent(tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonPrimary, tcell.ModNone)),
		gocuiEventFromTcellEvent(tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonNone, tcell.ModNone)),
		gocuiEventFromTcellEvent(tcell.NewEventMouse(view.x0+1, view.y0+1, tcell.ButtonPrimary, tcell.ModNone)),
	} {
		assert.NoError(t, g.onKey(&event))
	}

	assert.Equal(t, []bool{false, true}, doubleClicks)
}
