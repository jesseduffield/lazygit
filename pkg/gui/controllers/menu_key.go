package controllers

import "github.com/jesseduffield/lazygit/pkg/gocui"

// menuKey is a shorthand for constructing a key value for a menu item from a single rune literal,
// avoiding the noise of `gocui.NewKeyRune('a')` at every call site. There is an intentionally
// identical helper in the helpers package so that callers in either package can use the unqualified
// form.
func menuKey(r rune) gocui.Key {
	return gocui.NewKeyRune(r)
}
