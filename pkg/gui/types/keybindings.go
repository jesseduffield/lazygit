package types

import "github.com/jesseduffield/gocui"

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Contexts    []string
	Handler     func() error
	Key         interface{} // FIXME: find out how to get `gocui.Key | rune`
	Modifier    gocui.Modifier
	Description string
	Alternative string
	Tag         string // e.g. 'navigation'. Used for grouping things in the cheatsheet
	OpensMenu   bool
}

// A guard is a decorator which checks something before executing a handler
// and potentially early-exits if some precondition hasn't been met.
type Guard func(func() error) func() error

type KeybindingGuards struct {
	OutsideFilterMode Guard
	NoPopupPanel      Guard
}
