// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"unicode"
)

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
	case key == KeyBackspace || key == KeyBackspace2:
		v.TextArea.BackSpaceChar()
	case key == KeyCtrlD || key == KeyDelete:
		v.TextArea.DeleteChar()
	case key == KeyArrowDown:
		v.TextArea.MoveCursorDown()
	case key == KeyArrowUp:
		v.TextArea.MoveCursorUp()
	case key == KeyArrowLeft && (mod&ModAlt) != 0:
		v.TextArea.MoveLeftWord()
	case key == KeyArrowLeft:
		v.TextArea.MoveCursorLeft()
	case key == KeyArrowRight && (mod&ModAlt) != 0:
		v.TextArea.MoveRightWord()
	case key == KeyArrowRight:
		v.TextArea.MoveCursorRight()
	case key == KeyEnter:
		v.TextArea.TypeRune('\n')
	case key == KeySpace:
		v.TextArea.TypeRune(' ')
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

		// TODO: see if we need all three of these conditions: maybe the final one is sufficient
	case ch != 0 && mod == 0 && unicode.IsPrint(ch):
		v.TextArea.TypeRune(ch)
	default:
		return false
	}

	v.RenderTextArea()

	return true
}
