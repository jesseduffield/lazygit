package gocui

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setUpViewWithClickBindings creates a view that records the mouse keys of the
// bindings it receives.
func setUpViewWithClickBindings(t *testing.T, g *Gui) *[]KeyName {
	t.Helper()
	_, err := g.SetView("view", 10, 5, 30, 15, 0)
	if err != nil && !errors.Is(err, ErrUnknownView) {
		assert.NoError(t, err)
	}

	received := []KeyName{}
	for _, binding := range []*ViewMouseBinding{
		{ViewName: "view", Key: MouseLeft},
		{ViewName: "view", Key: MouseLeft, Modifier: ModMotion},
		{ViewName: "view", Key: MouseRelease},
		{ViewName: "view", Key: MouseWheelDown},
	} {
		binding.Handler = func(opts ViewMouseBindingOpts) error {
			received = append(received, opts.Key)
			return nil
		}
		g.SetViewClickBinding(binding)
	}
	return &received
}

// mouseEvent returns a mouse event in row 7, the third row of the view
func mouseEvent(x int, key KeyName, mod Modifier) *GocuiEvent {
	return &GocuiEvent{Type: eventMouse, MouseX: x, MouseY: 7, Key: NewKey(key, "", mod)}
}

func TestInterceptMouseOwnsTheGestureOfAPressItHandles(t *testing.T) {
	g := newTestGui(t)
	viewReceived := setUpViewWithClickBindings(t, g)
	intercepted := []MouseEvent{}
	g.InterceptMouse = func(ev MouseEvent) bool {
		intercepted = append(intercepted, ev)
		return ev.X == 10 && ev.Key == MouseLeft && ev.Mod == ModNone
	}

	assert.NoError(t, g.onKey(mouseEvent(10, MouseLeft, ModNone)))
	assert.NoError(t, g.onKey(mouseEvent(15, MouseLeft, ModMotion)))
	assert.NoError(t, g.onKey(mouseEvent(16, MouseRelease, ModNone)))

	assert.Equal(t, []MouseEvent{
		{X: 10, Y: 7, Key: MouseLeft, Mod: ModNone},
		{X: 15, Y: 7, Key: MouseLeft, Mod: ModMotion},
		{X: 16, Y: 7, Key: MouseRelease, Mod: ModNone},
	}, intercepted)
	assert.Empty(t, *viewReceived)

	// The next gesture reaches the view if InterceptMouse doesn't handle it
	assert.NoError(t, g.onKey(mouseEvent(15, MouseLeft, ModNone)))
	assert.NoError(t, g.onKey(mouseEvent(15, MouseRelease, ModNone)))
	assert.Equal(t, []KeyName{MouseLeft, MouseRelease}, *viewReceived)
}

func TestInterceptMouseDoesNotOwnTheGestureOfAWheelEventItHandles(t *testing.T) {
	g := newTestGui(t)
	viewReceived := setUpViewWithClickBindings(t, g)
	g.InterceptMouse = func(ev MouseEvent) bool {
		return ev.Key == MouseWheelDown
	}

	assert.NoError(t, g.onKey(mouseEvent(15, MouseWheelDown, ModNone)))
	assert.NoError(t, g.onKey(mouseEvent(15, MouseLeft, ModNone)))

	assert.Equal(t, []KeyName{MouseLeft}, *viewReceived)
}

func TestInterceptMouseIsNotAskedWhileAViewHasCapturedTheMouse(t *testing.T) {
	g := newTestGui(t)
	viewReceived := setUpViewWithClickBindings(t, g)
	asked := false
	g.InterceptMouse = func(ev MouseEvent) bool {
		asked = true
		return false
	}

	assert.NoError(t, g.onKey(mouseEvent(15, MouseLeft, ModNone)))
	asked = false
	g.InterceptMouse = func(ev MouseEvent) bool {
		asked = true
		return true
	}
	assert.NoError(t, g.onKey(mouseEvent(16, MouseLeft, ModMotion)))
	assert.NoError(t, g.onKey(mouseEvent(16, MouseRelease, ModNone)))

	assert.False(t, asked)
	assert.Equal(t, []KeyName{MouseLeft, MouseLeft, MouseRelease}, *viewReceived)
}

func TestOnMouseMoveSeesMovesOutsideViews(t *testing.T) {
	g := newTestGui(t)
	setUpViewWithClickBindings(t, g)
	type position struct{ x, y int }
	moves := []position{}
	g.OnMouseMove = func(x, y int) {
		moves = append(moves, position{x, y})
	}

	assert.NoError(t, g.onKey(&GocuiEvent{Type: eventMouseMove, MouseX: 15, MouseY: 7}))
	assert.NoError(t, g.onKey(&GocuiEvent{Type: eventMouseMove, MouseX: 50, MouseY: 20}))

	assert.Equal(t, []position{{15, 7}, {50, 20}}, moves)
}
