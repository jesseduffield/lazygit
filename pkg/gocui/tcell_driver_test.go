package gocui

import (
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func TestFirstMouseMovementAfterPressIsDragEvent(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()

	pressEvent := gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 2, tcell.ButtonPrimary, tcell.ModNone))
	unchangedHeldEvent := gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 2, tcell.ButtonPrimary, tcell.ModNone))
	dragEvent := gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 3, tcell.ButtonPrimary, tcell.ModNone))

	assert.Equal(t, eventMouse, pressEvent.Type)
	assert.Equal(t, MouseLeft, pressEvent.Key.KeyName())
	assert.Equal(t, ModNone, pressEvent.Key.Mod())
	assert.Equal(t, eventNone, unchangedHeldEvent.Type)
	assert.Equal(t, eventMouse, dragEvent.Type)
	assert.Equal(t, MouseLeft, dragEvent.Key.KeyName())
	assert.Equal(t, ModMotion, dragEvent.Key.Mod())
}

func TestMouseReleaseAfterDragIsMouseEvent(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()

	gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 2, tcell.ButtonPrimary, tcell.ModNone))
	gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 3, tcell.ButtonPrimary, tcell.ModNone))
	releaseEvent := gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 3, tcell.ButtonNone, tcell.ModNone))

	assert.Equal(t, eventMouse, releaseEvent.Type)
	assert.Equal(t, MouseRelease, releaseEvent.Key.KeyName())
}

func TestMouseReleaseDoesNotKeepPressModifiers(t *testing.T) {
	t.Cleanup(resetMouseState)
	resetMouseState()

	gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 2, tcell.ButtonPrimary, tcell.ModAlt))
	gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 3, tcell.ButtonPrimary, tcell.ModAlt))
	releaseEvent := gocuiEventFromTcellEvent(tcell.NewEventMouse(1, 3, tcell.ButtonNone, tcell.ModAlt))

	assert.Equal(t, eventMouse, releaseEvent.Type)
	assert.Equal(t, MouseRelease, releaseEvent.Key.KeyName())
	assert.Equal(t, ModNone, releaseEvent.Key.Mod())
}

func resetMouseState() {
	lastMouseKey = tcell.ButtonNone
	dragState = NOT_DRAGGING
	lastX = 0
	lastY = 0
}
