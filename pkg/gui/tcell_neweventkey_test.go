package gui_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// Regression (issue #5253): Shift+NumpadSubtract must not normalize to the same
// event as plain '-', or it wrongly triggers plain '-' keybindings.
func TestTcellNewEventKey_preservesShiftForHyphenRune(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyRune, '-', tcell.ModShift)
	if ev.Modifiers() != tcell.ModShift {
		t.Fatalf("want ModShift, got %v", ev.Modifiers())
	}
}
