// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"github.com/gdamore/tcell/v3"
)

// KeyName represents special keys or keys combinations.
type KeyName tcell.Key

// Modifier allows to define special keys combinations. They can be used
// in combination with Keys or Runes when a new keybinding is defined.
type Modifier tcell.ModMask

// Keybindings are used to link a given key-press event with a handler.
type keybinding struct {
	viewName string
	key      Key
	handler  func(*Gui, *View) error
}

// newKeybinding returns a new Keybinding object.
func newKeybinding(viewname string, key Key, handler func(*Gui, *View) error) (kb *keybinding) {
	kb = &keybinding{
		viewName: viewname,
		key:      key,
		handler:  handler,
	}
	return kb
}

// matchKeypress returns if the keybinding matches the keypress.
func (kb *keybinding) matchKeypress(key Key) bool {
	return kb.key.Equals(key)
}

// Special keys.
const (
	KeyF1             KeyName = KeyName(tcell.KeyF1)
	KeyF2                     = KeyName(tcell.KeyF2)
	KeyF3                     = KeyName(tcell.KeyF3)
	KeyF4                     = KeyName(tcell.KeyF4)
	KeyF5                     = KeyName(tcell.KeyF5)
	KeyF6                     = KeyName(tcell.KeyF6)
	KeyF7                     = KeyName(tcell.KeyF7)
	KeyF8                     = KeyName(tcell.KeyF8)
	KeyF9                     = KeyName(tcell.KeyF9)
	KeyF10                    = KeyName(tcell.KeyF10)
	KeyF11                    = KeyName(tcell.KeyF11)
	KeyF12                    = KeyName(tcell.KeyF12)
	KeyInsert                 = KeyName(tcell.KeyInsert)
	KeyDelete                 = KeyName(tcell.KeyDelete)
	KeyHome                   = KeyName(tcell.KeyHome)
	KeyEnd                    = KeyName(tcell.KeyEnd)
	KeyPgdn                   = KeyName(tcell.KeyPgDn)
	KeyPgup                   = KeyName(tcell.KeyPgUp)
	KeyArrowUp                = KeyName(tcell.KeyUp)
	KeyShiftArrowUp           = KeyName(tcell.KeyF62)
	KeyArrowDown              = KeyName(tcell.KeyDown)
	KeyShiftArrowDown         = KeyName(tcell.KeyF63)
	KeyArrowLeft              = KeyName(tcell.KeyLeft)
	KeyArrowRight             = KeyName(tcell.KeyRight)
)

// Keys combinations.
const (
	KeyCtrlTilde = KeyName(tcell.KeyF64) // arbitrary assignment
	KeyBackspace = KeyName(tcell.KeyBackspace)
	KeyTab       = KeyName(tcell.KeyTab)
	KeyBacktab   = KeyName(tcell.KeyBacktab)
	KeyEnter     = KeyName(tcell.KeyEnter)
	KeyEsc       = KeyName(tcell.KeyEscape)

	// The following assignments were used in termbox implementation.
	// In tcell, these are not keys per se. But in gocui we have them
	// mapped to the keys so we have to use placeholder keys.

	KeyAltEnter     = KeyName(tcell.KeyF64) // arbitrary assignments
	MouseLeft       = KeyName(tcell.KeyF63)
	MouseRight      = KeyName(tcell.KeyF62)
	MouseMiddle     = KeyName(tcell.KeyF61)
	MouseRelease    = KeyName(tcell.KeyF60)
	MouseWheelUp    = KeyName(tcell.KeyF59)
	MouseWheelDown  = KeyName(tcell.KeyF58)
	MouseWheelLeft  = KeyName(tcell.KeyF57)
	MouseWheelRight = KeyName(tcell.KeyF56)
)

// Modifiers.
const (
	ModNone   Modifier = Modifier(0)
	ModShift           = Modifier(tcell.ModShift)
	ModCtrl            = Modifier(tcell.ModCtrl)
	ModAlt             = Modifier(tcell.ModAlt)
	ModMeta            = Modifier(tcell.ModMeta)
	ModMotion          = Modifier(16) // just picking an arbitrary number here that doesn't clash with tcell's modifiers
)
