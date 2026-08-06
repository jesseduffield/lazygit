package gocui

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMouseCaptureRoutesMotionAndReleaseOutsideView(t *testing.T) {
	g := newTestGui(t)
	view, err := g.SetView("captured", 10, 5, 30, 15, 0)
	if err != nil && !errors.Is(err, ErrUnknownView) {
		assert.NoError(t, err)
		return
	}

	received := []ViewMouseBindingOpts{}
	for _, binding := range []*ViewMouseBinding{
		{
			ViewName: "captured",
			Key:      MouseLeft,
			Modifier: ModMotion,
			Handler: func(opts ViewMouseBindingOpts) error {
				received = append(received, opts)
				return nil
			},
		},
		{
			ViewName: "captured",
			Key:      MouseRelease,
			Handler: func(opts ViewMouseBindingOpts) error {
				assert.Nil(t, g.mouseCapture)
				received = append(received, opts)
				return nil
			},
		},
	} {
		assert.NoError(t, g.SetViewClickBinding(binding))
	}

	g.captureMouse(view)
	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: 0,
		MouseY: 0,
		Key:    NewKey(MouseLeft, "", ModMotion),
	}))
	assert.Equal(t, ViewMouseBindingOpts{X: -11, Y: -6, Key: MouseLeft}, received[0])
	assert.Equal(t, 0, view.CursorX())
	assert.Equal(t, 0, view.CursorY())

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: 79,
		MouseY: 23,
		Key:    NewKeyName(MouseRelease),
	}))
	assert.Equal(t, ViewMouseBindingOpts{X: 68, Y: 17, Key: MouseRelease}, received[1])
	assert.Equal(t, 0, view.CursorX())
	assert.Equal(t, 0, view.CursorY())
	assert.Nil(t, g.mouseCapture)
}

func TestPrimaryMouseDragStaysWithPressedView(t *testing.T) {
	g := newTestGui(t)
	left, _ := g.SetView("left", 0, 0, 20, 10, 0)
	_, _ = g.SetView("right", 21, 0, 41, 10, 0)

	receivedBy := ""
	for _, viewName := range []string{"left", "right"} {
		assert.NoError(t, g.SetViewClickBinding(&ViewMouseBinding{
			ViewName: viewName,
			Key:      MouseLeft,
			Modifier: ModMotion,
			Handler: func(ViewMouseBindingOpts) error {
				receivedBy = viewName
				return nil
			},
		}))
	}

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: left.x0 + 1,
		MouseY: left.y0 + 1,
		Key:    NewKeyName(MouseLeft),
	}))
	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: 22,
		MouseY: 1,
		Key:    NewKey(MouseLeft, "", ModMotion),
	}))

	assert.Equal(t, "left", receivedBy)
}

func TestPrimaryMouseDragDoesNotActivateTabs(t *testing.T) {
	g := newTestGui(t)
	view, _ := g.SetView("tabs", 0, 0, 40, 10, 0)
	view.Tabs = []string{"first", "second"}

	clickedTabs := []int{}
	assert.NoError(t, g.SetTabClickBinding("tabs", func(tabIndex int) error {
		clickedTabs = append(clickedTabs, tabIndex)
		return nil
	}))

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: view.x0 + 1,
		MouseY: view.y0 + 1,
		Key:    NewKeyName(MouseLeft),
	}))
	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: view.x0 + 3,
		MouseY: view.y0,
		Key:    NewKey(MouseLeft, "", ModMotion),
	}))

	assert.Empty(t, clickedTabs)

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: view.x0 + 3,
		MouseY: view.y0,
		Key:    NewKeyName(MouseRelease),
	}))
	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: view.x0 + 3,
		MouseY: view.y0,
		Key:    NewKeyName(MouseLeft),
	}))
	assert.Equal(t, []int{0}, clickedTabs)
}

func TestRejectedMouseReleaseClearsCapture(t *testing.T) {
	g := newTestGui(t)
	view, _ := g.SetView("captured", 0, 0, 20, 10, 0)
	g.captureMouse(view)
	g.ShouldHandleMouseEvent = func(*View, KeyName) bool { return false }

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: view.x0 + 1,
		MouseY: view.y0 + 1,
		Key:    NewKeyName(MouseRelease),
	}))

	assert.Nil(t, g.mouseCapture)
}

func TestDeleteViewClearsMouseState(t *testing.T) {
	g := newTestGui(t)
	view, _ := g.SetView("temporary", 0, 0, 20, 10, 0)
	g.captureMouse(view)
	g.lastHoverView = view

	assert.NoError(t, g.DeleteView("temporary"))

	assert.Nil(t, g.mouseCapture)
	assert.True(t, g.mouseGestureCanceled)
	assert.Nil(t, g.lastHoverView)
}

func TestCancelMouseCaptureSuppressesRemainingGesture(t *testing.T) {
	g := newTestGui(t)
	left, _ := g.SetView("left", 0, 0, 20, 10, 0)
	_, _ = g.SetView("right", 21, 0, 41, 10, 0)
	receivedBy := ""
	for _, viewName := range []string{"left", "right"} {
		assert.NoError(t, g.SetViewClickBinding(&ViewMouseBinding{
			ViewName: viewName,
			Key:      MouseLeft,
			Modifier: ModMotion,
			Handler: func(ViewMouseBindingOpts) error {
				receivedBy = viewName
				return nil
			},
		}))
	}

	g.captureMouse(left)
	g.CancelMouseCapture()
	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: 22,
		MouseY: 1,
		Key:    NewKey(MouseLeft, "", ModMotion),
	}))
	assert.Empty(t, receivedBy)

	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: 22,
		MouseY: 1,
		Key:    NewKeyName(MouseRelease),
	}))
	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: 22,
		MouseY: 1,
		Key:    NewKeyName(MouseLeft),
	}))
	assert.NoError(t, g.onKey(&GocuiEvent{
		Type:   eventMouse,
		MouseX: 23,
		MouseY: 1,
		Key:    NewKey(MouseLeft, "", ModMotion),
	}))

	assert.Equal(t, "right", receivedBy)
}
