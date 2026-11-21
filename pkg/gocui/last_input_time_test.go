package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastInputTime(t *testing.T) {
	scenarios := []struct {
		name  string
		event GocuiEvent
		input bool
	}{
		{"key", GocuiEvent{Type: eventKey, Key: NewKeyRune('x')}, true},
		{"mouse click", GocuiEvent{Type: eventMouse, Key: NewKeyName(MouseLeft)}, true},
		{"mouse scroll", GocuiEvent{Type: eventMouse, Key: NewKeyName(MouseWheelDown)}, true},
		{"mouse move", GocuiEvent{Type: eventMouseMove}, false},
		{"paste", GocuiEvent{Type: eventPaste, Start: true}, true},
		{"resize", GocuiEvent{Type: eventResize}, false},
		{"focus", GocuiEvent{Type: eventFocus}, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			g := newTestGui(t)
			assert.True(t, g.LastInputTime().IsZero())

			_ = g.handleEvent(&s.event)
			assert.Equal(t, s.input, !g.LastInputTime().IsZero())
		})
	}
}
