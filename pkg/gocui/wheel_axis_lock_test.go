package gocui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWheelAxisLock(t *testing.T) {
	const (
		V = wheelAxisVertical
		H = wheelAxisHorizontal
	)

	type event struct {
		axis wheelAxis
		// milliseconds since the first event
		ms       int
		accepted bool
	}

	scenarios := []struct {
		name   string
		events []event
	}{
		{
			name:   "drift during a vertical gesture is dropped",
			events: []event{{V, 0, true}, {H, 10, false}, {V, 20, true}, {H, 30, false}, {V, 40, true}},
		},
		{
			name:   "drift during a horizontal gesture is dropped",
			events: []event{{H, 0, true}, {V, 10, false}, {H, 20, true}, {V, 30, false}, {H, 40, true}},
		},
		{
			name:   "a pause starts a new gesture",
			events: []event{{V, 0, true}, {H, 300, true}, {V, 310, false}},
		},
		{
			name:   "dropped events keep the gesture going",
			events: []event{{V, 0, true}, {H, 250, false}, {H, 500, false}},
		},
		{
			name:   "a run of events on the other axis switches the gesture to that axis",
			events: []event{{V, 0, true}, {H, 10, false}, {H, 20, false}, {H, 30, true}, {V, 40, false}, {H, 50, true}},
		},
		{
			name:   "an event on the gesture's axis interrupts the run",
			events: []event{{V, 0, true}, {H, 10, false}, {H, 20, false}, {V, 30, true}, {H, 40, false}, {H, 50, false}, {H, 60, true}},
		},
	}

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			lock := wheelAxisLock{}
			for i, e := range s.events {
				when := start.Add(time.Duration(e.ms) * time.Millisecond)
				assert.Equal(t, e.accepted, lock.accept(e.axis, when), "event %d", i)
			}
		})
	}
}
