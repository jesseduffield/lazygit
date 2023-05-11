package types

import "github.com/jesseduffield/gocui"

type Key interface{} // FIXME: find out how to get `gocui.Key | rune`

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Handler     func() error
	Key         Key
	Modifier    gocui.Modifier
	Description string
	Alternative string
	Tag         string // e.g. 'navigation'. Used for grouping things in the cheatsheet
	OpensMenu   bool

	// If true, the keybinding will appear at the bottom of the screen. If
	// the given view has no bindings with Display: true, the default keybindings
	// will be displayed instead.
	// TODO: implement this
	Display bool

	// to be displayed if the keybinding is highlighted from within a menu
	Tooltip string
}

// A guard is a decorator which checks something before executing a handler
// and potentially early-exits if some precondition hasn't been met.
type Guard func(func() error) func() error

type KeybindingGuards struct {
	OutsideFilterMode Guard
	NoPopupPanel      Guard
}
