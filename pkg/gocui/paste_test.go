package gocui

import (
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/stretchr/testify/assert"
)

func TestPasteIntoViewWithPasteHandler(t *testing.T) {
	g := newTestGui(t)
	v, _ := g.SetView("editor", 0, 0, 40, 10, 0)
	v.Editable = true
	_, err := g.SetCurrentView("editor")
	assert.NoError(t, err)

	pasted := []string{}
	v.OnPaste = func(text string) error {
		pasted = append(pasted, text)
		return nil
	}
	escPressed := false
	g.SetKeybinding("editor", NewKeyName(KeyEsc), func(*Gui, *View) error {
		escPressed = true
		return nil
	})

	for _, tev := range []tcell.Event{
		tcell.NewEventPaste(true),
		tcell.NewEventKey(tcell.KeyRune, "a", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyEnter, "", tcell.ModNone),
		// Ghostty pastes newlines as ctrl-j
		tcell.NewEventKey(tcell.KeyRune, "j", tcell.ModCtrl),
		tcell.NewEventKey(tcell.KeyTab, "", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyEsc, "", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, "a", tcell.ModCtrl),
		tcell.NewEventKey(tcell.KeyRune, "b", tcell.ModNone),
		tcell.NewEventPaste(false),
	} {
		ev := gocuiEventFromTcellEvent(tev)
		assert.NoError(t, g.handleEvent(&ev))
	}

	assert.Equal(t, []string{"a\n\n\tb"}, pasted)
	assert.False(t, escPressed)
	assert.Equal(t, "", v.TextArea.GetContent())
	assert.False(t, g.IsPasting)
}

func TestPasteIntoViewWithoutPasteHandler(t *testing.T) {
	g := newTestGui(t)
	v, _ := g.SetView("editor", 0, 0, 40, 10, 0)
	v.Editable = true
	v.Editor = DefaultEditor
	_, err := g.SetCurrentView("editor")
	assert.NoError(t, err)

	for _, tev := range []tcell.Event{
		tcell.NewEventPaste(true),
		tcell.NewEventKey(tcell.KeyRune, "a", tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, "b", tcell.ModNone),
		tcell.NewEventPaste(false),
	} {
		ev := gocuiEventFromTcellEvent(tev)
		assert.NoError(t, g.handleEvent(&ev))
	}

	assert.Equal(t, "ab", v.TextArea.GetContent())
}
