package gocui

import "time"

// wheelGestureTimeout is the longest pause between two wheel events of the same
// scroll gesture.
const wheelGestureTimeout = 300 * time.Millisecond

// wheelAxisSwitchCount is the number of events in a row on the other axis that make a
// scroll gesture switch to that axis.
const wheelAxisSwitchCount = 3

type wheelAxis int

const (
	wheelAxisVertical wheelAxis = iota
	wheelAxisHorizontal
)

func wheelAxisOfKey(key KeyName) (wheelAxis, bool) {
	switch key {
	case MouseWheelUp, MouseWheelDown:
		return wheelAxisVertical, true
	case MouseWheelLeft, MouseWheelRight:
		return wheelAxisHorizontal, true
	default:
		return 0, false
	}
}

// wheelAxisLock keeps a scroll gesture on the axis that it started on.
//
// When the user swipes vertically on a trackpad, the fingers usually drift sideways a
// little, and the terminal reports the drift as horizontal wheel events mixed in with
// the vertical ones. The same happens the other way round. Terminals don't tell us
// when the fingers touch the trackpad or leave it, so a gesture is a run of wheel
// events with less than wheelGestureTimeout between neighbours. The events that the
// lock drops count too, so a gesture lasts as long as events keep coming.
//
// The first event of a gesture decides its axis, and the lock drops the events on the
// other axis. If wheelAxisSwitchCount events in a row arrive on the other axis, the
// gesture switches to that axis. This repairs a gesture whose first event was drift,
// and lets the user change direction without pausing.
type wheelAxisLock struct {
	axis           wheelAxis
	lastEventTime  time.Time
	otherAxisCount int
}

// accept reports whether a wheel event on the given axis should reach the bindings.
// when is the time at which the event was read.
func (l *wheelAxisLock) accept(axis wheelAxis, when time.Time) bool {
	startsGesture := l.lastEventTime.IsZero() || when.Sub(l.lastEventTime) >= wheelGestureTimeout
	l.lastEventTime = when

	if !startsGesture && axis != l.axis {
		l.otherAxisCount++
		if l.otherAxisCount < wheelAxisSwitchCount {
			return false
		}
	}

	l.axis = axis
	l.otherAxisCount = 0
	return true
}
