// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

// Editor interface must be satisfied by gocui editors.
type Editor interface {
	Edit(v *View, key Key) bool
}

// The EditorFunc type is an adapter to allow the use of ordinary functions as
// Editors. If f is a function with the appropriate signature, EditorFunc(f)
// is an Editor object that calls f.
type EditorFunc func(v *View, key Key) bool

// Edit calls f(v, key, mod)
func (f EditorFunc) Edit(v *View, key Key) bool {
	return f(v, key)
}

// DefaultEditor is the default editor.
var DefaultEditor Editor = EditorFunc(SimpleEditor)

var (
	moveWordLeftKeybinding      = NewKey(KeyArrowLeft, "", ModCtrl)
	moveWordRightKeybinding     = NewKey(KeyArrowRight, "", ModCtrl)
	backspaceWordKeybinding     = NewKey(KeyBackspace, "", ModCtrl)
	forwardDeleteWordKeybinding = NewKey(KeyDelete, "", ModCtrl)
)

// SimpleEditor is used as the default gocui editor.
func SimpleEditor(v *View, key Key) bool {
	switch {
	case key.Equals(backspaceWordKeybinding),
		key.Equals(NewKeyStrMod("w", ModCtrl)):
		v.TextArea.BackSpaceWord()
	case key.Equals(forwardDeleteWordKeybinding),
		key.Equals(NewKeyStrMod("d", ModAlt)):
		v.TextArea.ForwardDeleteWord()
	case key.Equals(NewKeyName(KeyBackspace)),
		key.Equals(NewKeyStrMod("h", ModCtrl)):
		v.TextArea.BackSpaceChar()
	case key.Equals(NewKeyStrMod("d", ModCtrl)),
		key.Equals(NewKeyName(KeyDelete)):
		v.TextArea.DeleteChar()
	case key.Equals(NewKeyName(KeyArrowDown)):
		v.TextArea.MoveCursorDown()
	case key.Equals(NewKeyName(KeyArrowUp)):
		v.TextArea.MoveCursorUp()
	case key.Equals(NewKeyStrMod("b", ModAlt)),
		key.Equals(moveWordLeftKeybinding):
		v.TextArea.MoveLeftWord()
	case key.Equals(NewKeyName(KeyArrowLeft)),
		key.Equals(NewKeyStrMod("b", ModCtrl)):
		v.TextArea.MoveCursorLeft()
	case key.Equals(NewKeyStrMod("f", ModAlt)),
		key.Equals(moveWordRightKeybinding):
		v.TextArea.MoveRightWord()
	case key.Equals(NewKeyName(KeyArrowRight)),
		key.Equals(NewKeyStrMod("f", ModCtrl)):
		v.TextArea.MoveCursorRight()
	case key.Equals(NewKeyName(KeyEnter)):
		v.TextArea.TypeCharacter("\n")
	case key.Equals(NewKeyName(KeyInsert)):
		v.TextArea.ToggleOverwrite()
	case key.Equals(NewKeyStrMod("u", ModCtrl)):
		v.TextArea.DeleteToStartOfLine()
	case key.Equals(NewKeyStrMod("k", ModCtrl)):
		v.TextArea.DeleteToEndOfLine()
	case key.Equals(NewKeyStrMod("a", ModCtrl)),
		key.Equals(NewKeyName(KeyHome)):
		v.TextArea.GoToStartOfLine()
	case key.Equals(NewKeyStrMod("e", ModCtrl)),
		key.Equals(NewKeyName(KeyEnd)):
		v.TextArea.GoToEndOfLine()
	case key.Equals(NewKeyStrMod("y", ModCtrl)):
		v.TextArea.Yank()
	case key.Str() != "" && key.Mod() == 0:
		v.TextArea.TypeCharacter(key.Str())
	default:
		return false
	}

	v.RenderTextArea()

	return true
}
