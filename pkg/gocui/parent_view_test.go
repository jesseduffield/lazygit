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

func TestEmbeddedViewsAreFocusedTogether(t *testing.T) {
	g := newTestGui(t)
	parent, child := setupParentAndChildView(t, g)
	sibling, _ := g.SetView("sibling", 0, 12, 20, 14, 0)
	sibling.ParentView = parent
	unrelated, _ := g.SetView("unrelated", 30, 0, 50, 10, 0)

	assert.True(t, g.hasFocus(child))
	assert.True(t, g.hasFocus(parent))
	assert.True(t, g.hasFocus(sibling))
	assert.False(t, g.hasFocus(unrelated))

	_, err := g.SetCurrentView(unrelated.Name())
	assert.NoError(t, err)

	assert.True(t, g.hasFocus(unrelated))
	assert.False(t, g.hasFocus(parent))
	assert.False(t, g.hasFocus(child))
}

func TestPrintableKeysGoToTheFieldBeingTypedIn(t *testing.T) {
	for _, test := range []struct {
		name              string
		keybindOnEdit     bool
		declineKeybinding bool
		expectedPresses   int
		expectedEdits     int
	}{
		{name: "the field gets the key", expectedEdits: 1},
		{name: "the parent view gets the key", keybindOnEdit: true, expectedPresses: 1},
		{
			name:              "the field gets the key the parent view declined",
			keybindOnEdit:     true,
			declineKeybinding: true,
			expectedPresses:   1,
			expectedEdits:     1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			g := newTestGui(t)
			parent, child := setupParentAndChildView(t, g)
			child.Editable = true
			child.KeybindOnEdit = test.keybindOnEdit

			edits := 0
			child.Editor = EditorFunc(func(*View, Key) bool {
				edits++
				return true
			})
			presses := 0
			g.SetKeybinding(parent.Name(), NewKeyRune('j'), func(*Gui, *View) error {
				presses++
				if test.declineKeybinding {
					return ErrKeybindingNotHandled
				}
				return nil
			})

			assert.NoError(t, g.onKey(&GocuiEvent{Type: eventKey, Key: NewKeyRune('j')}))

			assert.Equal(t, test.expectedPresses, presses)
			assert.Equal(t, test.expectedEdits, edits)
		})
	}
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
