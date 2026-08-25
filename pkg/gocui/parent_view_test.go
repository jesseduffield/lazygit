package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// A view and its parent view, with the child holding the focus.
func setupParentAndChildView(t *testing.T, g *Gui) (*View, *View) {
	t.Helper()

	parent, _ := g.SetView("parent", 0, 0, 20, 10, 0)
	child, _ := g.SetView("child", 0, 10, 20, 12, 0)
	child.ParentView = parent
	_, err := g.SetCurrentView(child.Name())
	assert.NoError(t, err)

	return parent, child
}

func TestKeybindingOfParentViewIsUsedWhenChildHasNone(t *testing.T) {
	g := newTestGui(t)
	parent, child := setupParentAndChildView(t, g)

	pressed := []string{}
	g.SetKeybinding(parent.Name(), NewKeyName(KeyArrowDown), func(*Gui, *View) error {
		pressed = append(pressed, "parent")
		return nil
	})
	g.SetKeybinding(child.Name(), NewKeyName(KeyEnter), func(*Gui, *View) error {
		pressed = append(pressed, "child")
		return nil
	})

	assert.NoError(t, g.onKey(&GocuiEvent{Type: eventKey, Key: NewKeyName(KeyArrowDown)}))
	assert.NoError(t, g.onKey(&GocuiEvent{Type: eventKey, Key: NewKeyName(KeyEnter)}))

	assert.Equal(t, []string{"parent", "child"}, pressed)
}

func TestFirstMatchingKeybindingOfParentViewWins(t *testing.T) {
	g := newTestGui(t)
	parent, _ := setupParentAndChildView(t, g)

	pressed := []string{}
	for _, name := range []string{"first", "second"} {
		g.SetKeybinding(parent.Name(), NewKeyName(KeyArrowDown), func(*Gui, *View) error {
			pressed = append(pressed, name)
			return nil
		})
	}

	assert.NoError(t, g.onKey(&GocuiEvent{Type: eventKey, Key: NewKeyName(KeyArrowDown)}))

	assert.Equal(t, []string{"first"}, pressed)
}

func TestUnhandledKeybindingOfParentViewFallsThroughToEditor(t *testing.T) {
	g := newTestGui(t)
	parent, child := setupParentAndChildView(t, g)

	edited := []Key{}
	child.Editable = true
	child.Editor = EditorFunc(func(_ *View, key Key) bool {
		edited = append(edited, key)
		return true
	})
	g.SetKeybinding(parent.Name(), NewKeyName(KeyArrowDown), func(*Gui, *View) error {
		return ErrKeybindingNotHandled
	})

	assert.NoError(t, g.onKey(&GocuiEvent{Type: eventKey, Key: NewKeyName(KeyArrowDown)}))

	assert.Equal(t, []Key{NewKeyName(KeyArrowDown)}, edited)
}
