// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

// Editor interface must be satisfied by gocui editors.
type Editor interface {
	Edit(v *View, key Key, ch rune, mod Modifier) bool
}

// The EditorFunc type is an adapter to allow the use of ordinary functions as
// Editors. If f is a function with the appropriate signature, EditorFunc(f)
// is an Editor object that calls f.
type EditorFunc func(v *View, key Key, ch rune, mod Modifier) bool

// Edit calls f(v, key, ch, mod)
func (f EditorFunc) Edit(v *View, key Key, ch rune, mod Modifier) bool {
	return f(v, key, ch, mod)
}

// DefaultEditor is the default editor.
var DefaultEditor Editor = EditorFunc(SimpleEditor)

// SimpleEditor is used as the default gocui editor.
func SimpleEditor(v *View, key Key, ch rune, mod Modifier) bool {
	switch {
	case (key == KeyBackspace || key == KeyBackspace2) && (mod&ModAlt) != 0,
		key == KeyCtrlW:
		v.TextArea.BackSpaceWord()
	case key == KeyBackspace || key == KeyBackspace2 || key == KeyCtrlH:
		v.TextArea.BackSpaceChar()
	case key == KeyCtrlD || key == KeyDelete:
		v.TextArea.DeleteChar()
	case key == KeyArrowDown:
		v.TextArea.MoveCursorDown()
	case key == KeyArrowUp:
		v.TextArea.MoveCursorUp()
	case (key == KeyArrowLeft || ch == 'b') && (mod&ModAlt) != 0:
		v.TextArea.MoveLeftWord()
	case key == KeyArrowLeft || key == KeyCtrlB:
		v.TextArea.MoveCursorLeft()
	case (key == KeyArrowRight || ch == 'f') && (mod&ModAlt) != 0:
		v.TextArea.MoveRightWord()
	case key == KeyArrowRight || key == KeyCtrlF:
		v.TextArea.MoveCursorRight()
	case key == KeyEnter:
		v.TextArea.TypeCharacter("\n")
	case key == KeySpace:
		v.TextArea.TypeCharacter(" ")
	case key == KeyInsert:
		v.TextArea.ToggleOverwrite()
	case key == KeyCtrlU:
		v.TextArea.DeleteToStartOfLine()
	case key == KeyCtrlK:
		v.TextArea.DeleteToEndOfLine()
	case key == KeyCtrlA || key == KeyHome:
		v.TextArea.GoToStartOfLine()
	case key == KeyCtrlE || key == KeyEnd:
		v.TextArea.GoToEndOfLine()
	case key == KeyCtrlW:
		v.TextArea.BackSpaceWord()
	case key == KeyCtrlY:
		v.TextArea.Yank()
	case ch != 0:
		v.TextArea.TypeCharacter(string(ch))
	default:
		return false
	}

	v.RenderTextArea()

	return true
}
