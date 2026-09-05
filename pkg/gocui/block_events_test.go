package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventWithheldWhileBlocking(t *testing.T) {
	scenarios := []struct {
		name     string
		event    GocuiEvent
		withheld bool
	}{
		{"key", GocuiEvent{Type: eventKey, Key: NewKeyRune('x')}, true},
		{"mouse click", GocuiEvent{Type: eventMouse, Key: NewKeyName(MouseLeft)}, true},
		{"mouse scroll", GocuiEvent{Type: eventMouse, Key: NewKeyName(MouseWheelDown)}, false},
		{"mouse move", GocuiEvent{Type: eventMouseMove}, true},
		{"resize", GocuiEvent{Type: eventResize}, false},
		{"focus", GocuiEvent{Type: eventFocus}, false},
		{"paste", GocuiEvent{Type: eventPaste}, false},
		{"error", GocuiEvent{Type: eventError}, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.withheld, eventWithheldWhileBlocking(&s.event))
		})
	}
}

// setupKeyRecorder wires a keybinding on a focused view that records each time
// it fires, and returns the key event that triggers it plus the record slice.
func setupKeyRecorder(t *testing.T, g *Gui) (GocuiEvent, *[]int) {
	t.Helper()

	_, _ = g.SetView("main", 0, 0, 80, 22, 0)
	_, err := g.SetCurrentView("main")
	assert.NoError(t, err)

	fired := []int{}
	callCount := 0
	key := NewKeyRune('x')
	g.SetKeybinding("main", key, func(*Gui, *View) error {
		callCount++
		fired = append(fired, callCount)
		return nil
	})

	return GocuiEvent{Type: eventKey, Key: key}, &fired
}

// runQueuedWork runs what the gui has queued for the following passes of the
// event loop, which is where EndBlockingEvents leaves the replay of the keys it
// buffered.
func runQueuedWork(t *testing.T, g *Gui) {
	t.Helper()

	for {
		ev, ok := g.userEvents.dequeue()
		if !ok {
			return
		}
		assert.NoError(t, ev.f(g))
		ev.task.Done()
	}
}

func TestBlockingEvents_KeysBufferedAndReplayed(t *testing.T) {
	g := newTestGui(t)
	keyEvent, fired := setupKeyRecorder(t, g)

	// Not blocking: the key dispatches immediately.
	assert.NoError(t, g.handleEvent(&keyEvent))
	assert.Len(t, *fired, 1)

	// While blocking: the key is buffered, not dispatched.
	g.BeginBlockingEvents()
	assert.NoError(t, g.handleEvent(&keyEvent))
	assert.NoError(t, g.handleEvent(&keyEvent))
	assert.Len(t, *fired, 1, "buffered keys must not dispatch while blocking")

	// Unblocking queues the replay rather than dispatching from here.
	g.EndBlockingEvents()
	assert.Len(t, *fired, 1, "the replay must wait for the event loop")

	runQueuedWork(t, g)
	assert.Len(t, *fired, 3, "both buffered keys should replay on unblock")
	assert.Empty(t, g.bufferedKeyEvents)
}

func TestBlockingEvents_KeysArrivingBeforeTheReplayGoBehindIt(t *testing.T) {
	g := newTestGui(t)
	keyEvent, fired := setupKeyRecorder(t, g)
	other := GocuiEvent{Type: eventKey, Key: NewKeyRune('y')}
	g.SetKeybinding("main", other.Key, func(*Gui, *View) error {
		*fired = append(*fired, 0)
		return nil
	})

	g.BeginBlockingEvents()
	assert.NoError(t, g.handleEvent(&keyEvent))
	g.EndBlockingEvents()

	// A key pressed while the replay is still queued joins the end of the buffer:
	// dispatching it now would put it ahead of the keys buffered before it.
	assert.NoError(t, g.handleEvent(&other))
	assert.Empty(t, *fired)

	runQueuedWork(t, g)
	assert.Equal(t, []int{1, 0}, *fired, "the keys should arrive in the order they were pressed")
}

func TestBlockingEvents_NestsWithCounter(t *testing.T) {
	g := newTestGui(t)
	keyEvent, fired := setupKeyRecorder(t, g)

	g.BeginBlockingEvents()
	g.BeginBlockingEvents()
	assert.NoError(t, g.handleEvent(&keyEvent))

	// The inner block ending still leaves us blocked: no replay yet.
	g.EndBlockingEvents()
	runQueuedWork(t, g)
	assert.Empty(t, *fired)

	// Only the outermost block ending replays.
	g.EndBlockingEvents()
	runQueuedWork(t, g)
	assert.Len(t, *fired, 1)
}

func TestBlockingEvents_MouseClicksDroppedNotBuffered(t *testing.T) {
	g := newTestGui(t)

	g.BeginBlockingEvents()
	click := GocuiEvent{Type: eventMouse, Key: NewKeyName(MouseLeft)}
	assert.NoError(t, g.handleEvent(&click))
	assert.Empty(t, g.bufferedKeyEvents, "mouse clicks must be dropped, not buffered")
	g.EndBlockingEvents()
	runQueuedWork(t, g)
}
