package gocui

import (
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func TestAModifiedClickNoBindingWantsLeavesTheViewAlone(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()
	g := newTestGui(t)
	view, _ := g.SetView("list", 0, 0, 20, 10, 0)

	clicks := 0
	g.SetViewClickBinding(&ViewMouseBinding{
		ViewName: "list",
		Key:      MouseLeft,
		Handler: func(ViewMouseBindingOpts) error {
			clicks++
			return nil
		},
	})

	click := func(y int, modifier tcell.ModMask) {
		for _, button := range []tcell.ButtonMask{tcell.ButtonPrimary, tcell.ButtonNone} {
			event := gocuiEventFromTcellEvent(
				tcell.NewEventMouse(view.x0+1, y, button, modifier))
			assert.NoError(t, g.onKey(&event))
		}
	}

	// A plain click is the binding's, and moves the cursor it acts on.
	click(view.y0+4, tcell.ModNone)
	assert.Equal(t, 1, clicks)
	assert.Equal(t, 3, view.CursorY())

	// An alt-click is nobody's here, bindings matching modifiers exactly. It has to
	// leave the cursor where the selection is, or a list view would draw its
	// selection on a line its owner never selected.
	click(view.y0+8, tcell.ModAlt)
	assert.Equal(t, 1, clicks)
	assert.Equal(t, 3, view.CursorY())
	assert.Nil(t, g.mouseCapture, "and no drag begins for it either")
}
