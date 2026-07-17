package gocui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResumeSchedulesRedraw(t *testing.T) {
	g := newTestGui(t)

	assert.NoError(t, g.Suspend())
	assert.NoError(t, g.Resume())

	ev := GocuiEvent{Type: eventNone}
	select {
	case ev = <-g.gEvents:
	case <-time.After(100 * time.Millisecond):
	}

	/* EXPECTED:
	assert.Equal(t, eventResize, ev.Type,
		"resuming must schedule a redraw; without one the screen stays blank until the next event arrives")
	ACTUAL: */
	assert.Equal(t, eventNone, ev.Type)
}
