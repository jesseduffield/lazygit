package filtering

import "testing"

func TestFilteringScreenModeChangedReset(t *testing.T) {
	mode := New("", "")
	mode.SetPath("foo")
	mode.SetScreenModeChanged(true)

	if !mode.ScreenModeChanged() {
		t.Fatalf("expected screenModeChanged to be true")
	}

	mode.Reset()

	if mode.ScreenModeChanged() {
		t.Fatalf("expected screenModeChanged to be false after Reset")
	}
}
