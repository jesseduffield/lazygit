package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnlyBindingsThatOptedInFireBehindAFocusedPopup(t *testing.T) {
	g := newTestGui(t)
	_, _ = g.SetView("main", 0, 0, 20, 10, 0)

	fired := []string{}
	g.SetViewClickBinding(&ViewMouseBinding{
		ViewName:                    "main",
		Key:                         MouseLeft,
		Modifier:                    ModAlt,
		HandleWhenPopupPanelFocused: true,
		Handler: func(ViewMouseBindingOpts) error {
			fired = append(fired, "opted in")
			return nil
		},
	})
	g.SetViewClickBinding(&ViewMouseBinding{
		ViewName: "main",
		Key:      MouseLeft,
		Handler: func(ViewMouseBindingOpts) error {
			fired = append(fired, "ordinary")
			return nil
		},
	})

	// This is how a client reports that a popup panel has the focus and the click
	// landed on a view behind it.
	g.ShouldHandleMouseEvent = func(*View, KeyName) bool { return false }

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type: eventMouse, MouseX: 3, MouseY: 4,
		Key: NewKey(MouseLeft, "", ModAlt),
	}))
	assert.Equal(t, []string{"opted in"}, fired)

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type: eventMouse, MouseX: 3, MouseY: 4,
		Key: NewKey(MouseLeft, "", ModNone),
	}))
	assert.Equal(t, []string{"opted in"}, fired,
		"a binding that didn't opt in must still be swallowed")
}
