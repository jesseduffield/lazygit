package gocui

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

// wheelGestureLog collects the wheel events of a gesture, and writes one line for the
// whole gesture to the debug log once the wheel has been idle for
// wheelGestureTimeout.
//
// Each event is a letter: U, D, L or R if the lock let it through, and u, d, l or r
// if it dropped it. A | before a letter means that the gesture switched to that
// event's axis. A number in parentheses before a letter is the time in
// milliseconds since the previous event, shown if it was 40ms or more.
type wheelGestureLog struct {
	mutex      sync.Mutex
	generation int

	pauseBefore time.Duration // negative for the first gesture
	start, last time.Time
	events      strings.Builder
	eventCount  int
	dropped     int
	mods        []string
}

func (l *wheelGestureLog) record(
	log func(string), ev *GocuiEvent, previousEventTime time.Time, previousAxis wheelAxis, accepted bool, axis wheelAxis,
) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	startsGesture := previousEventTime.IsZero() || ev.When.Sub(previousEventTime) >= wheelGestureTimeout
	if startsGesture {
		l.flushLocked(log)
		l.pauseBefore = -1
		if !previousEventTime.IsZero() {
			l.pauseBefore = ev.When.Sub(previousEventTime)
		}
		l.start = ev.When
	} else if gap := ev.When.Sub(l.last); gap >= 40*time.Millisecond {
		fmt.Fprintf(&l.events, "(%d)", gap.Milliseconds())
	}
	l.last = ev.When

	if accepted && !startsGesture && axis != previousAxis {
		l.events.WriteString("|")
	}
	letter := map[KeyName]string{
		MouseWheelUp:    "U",
		MouseWheelDown:  "D",
		MouseWheelLeft:  "L",
		MouseWheelRight: "R",
	}[ev.Key.KeyName()]
	if !accepted {
		letter = strings.ToLower(letter)
		l.dropped++
	}
	l.events.WriteString(letter)
	l.eventCount++

	for _, m := range []struct {
		mod  Modifier
		name string
	}{{ModShift, "shift"}, {ModCtrl, "ctrl"}, {ModAlt, "alt"}, {ModMeta, "meta"}} {
		if ev.WheelMod&m.mod != 0 && !slices.Contains(l.mods, m.name) {
			l.mods = append(l.mods, m.name)
		}
	}

	l.generation++
	generation := l.generation
	go safe(func() {
		time.Sleep(wheelGestureTimeout)
		l.mutex.Lock()
		defer l.mutex.Unlock()
		if l.generation == generation {
			l.flushLocked(log)
		}
	})
}

func (l *wheelGestureLog) flushLocked(log func(string)) {
	if l.eventCount == 0 {
		return
	}

	pause := "start"
	if l.pauseBefore >= 0 {
		pause = fmt.Sprintf("%dms", l.pauseBefore.Milliseconds())
	}
	details := fmt.Sprintf("%dms, %d events", l.last.Sub(l.start).Milliseconds(), l.eventCount)
	if l.dropped > 0 {
		details += fmt.Sprintf(", %d dropped", l.dropped)
	}
	if len(l.mods) > 0 {
		details += ", mods " + strings.Join(l.mods, "+")
	}
	log(fmt.Sprintf("WHEEL after %s: %s  (%s)", pause, l.events.String(), details))

	l.events.Reset()
	l.eventCount = 0
	l.dropped = 0
	l.mods = nil
}
