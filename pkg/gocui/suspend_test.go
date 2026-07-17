package gocui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// A flush while suspended must return without touching the screen: tcell
// releases the screen's cell buffer when disengaging, and drawing to a
// disengaged screen spins forever inside tcell while holding the screen lock,
// blocking the resume triggered by fg (#5309). The flush runs in a goroutine
// so that a regression fails the test instead of hanging the suite.
func TestFlushIsNoOpWhileSuspended(t *testing.T) {
	tests := []struct {
		name  string
		flush func(g *Gui) error
	}{
		{"flush", func(g *Gui) error { return g.flush() }},
		{"flushContentOnly", func(g *Gui) error { return g.flushContentOnly(g.views) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Deliberately not newTestGui: its cleanup closes the screen,
			// which would deadlock on the screen lock if a regression makes
			// the flush below spin.
			g, err := NewGui(NewGuiOpts{
				OutputMode: OutputNormal,
				Headless:   true,
				Width:      80,
				Height:     24,
			})
			assert.NoError(t, err)

			assert.NoError(t, g.Suspend())

			flushReturned := make(chan error, 1)
			go func() { flushReturned <- tc.flush(g) }()

			select {
			case err := <-flushReturned:
				assert.NoError(t, err)
			case <-time.After(time.Second):
				t.Fatal("flush touched the suspended screen and got stuck")
			}

			assert.NoError(t, g.Resume())
			g.Close()
		})
	}
}

func TestResumeSchedulesRedraw(t *testing.T) {
	g := newTestGui(t)

	assert.NoError(t, g.Suspend())
	assert.NoError(t, g.Resume())

	ev := GocuiEvent{Type: eventNone}
	select {
	case ev = <-g.gEvents:
	case <-time.After(100 * time.Millisecond):
	}

	assert.Equal(t, eventResize, ev.Type,
		"resuming must schedule a redraw; without one the screen stays blank until the next event arrives")
}
