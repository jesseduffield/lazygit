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

	// weighted selects the rule of acceptWeighted, for comparing the two rules.
	// LAZYGIT_WHEEL_LOCK=weighted turns it on.
	weighted bool
	score    int
}

// accept reports whether a wheel event on the given axis should reach the bindings.
// when is the time at which the event was read.
func (l *wheelAxisLock) accept(axis wheelAxis, when time.Time) bool {
	if l.weighted {
		return l.acceptWeighted(axis, when)
	}

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

// wheelAxisSwitchScore is the score at which acceptWeighted switches a gesture to the
// other axis.
const wheelAxisSwitchScore = 3

// acceptWeighted works like accept, except for how a gesture switches axis. An event
// on the other axis adds its weight to a score, and an event on the gesture's axis
// subtracts its weight, down to zero. Once the score reaches wheelAxisSwitchScore, the
// gesture switches to the other axis. A vertical event weighs 2 and a horizontal one
// 1, so it takes two vertical events to switch a horizontal gesture, even with a
// horizontal one between them. A vertical gesture switches only if horizontal events
// outnumber vertical ones by more than 2 to 1, or if three come in a row.
func (l *wheelAxisLock) acceptWeighted(axis wheelAxis, when time.Time) bool {
	startsGesture := l.lastEventTime.IsZero() || when.Sub(l.lastEventTime) >= wheelGestureTimeout
	l.lastEventTime = when

	if startsGesture {
		l.axis = axis
		l.score = 0
		return true
	}

	if axis == l.axis {
		l.score = max(0, l.score-wheelAxisWeight(axis))
		return true
	}

	l.score += wheelAxisWeight(axis)
	if l.score < wheelAxisSwitchScore {
		return false
	}
	l.axis = axis
	l.score = 0
	return true
}

func wheelAxisWeight(axis wheelAxis) int {
	if axis == wheelAxisVertical {
		return 2
	}
	return 1
}
