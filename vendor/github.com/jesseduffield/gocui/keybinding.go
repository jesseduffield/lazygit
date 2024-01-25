// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Key represents special keys or keys combinations.
type Key tcell.Key

// Modifier allows to define special keys combinations. They can be used
// in combination with Keys or Runes when a new keybinding is defined.
type Modifier tcell.ModMask

// Keybidings are used to link a given key-press event with a handler.
type keybinding struct {
	viewName string
	key      Key
	ch       rune
	mod      Modifier
	handler  func(*Gui, *View) error
}

// Parse takes the input string and extracts the keybinding.
// Returns a Key / rune, a Modifier and an error.
func Parse(input string) (interface{}, Modifier, error) {
	if len(input) == 1 {
		_, r, err := getKey(rune(input[0]))
		if err != nil {
			return nil, ModNone, err
		}
		return r, ModNone, nil
	}

	var modifier Modifier
	cleaned := make([]string, 0)

	tokens := strings.Split(input, "+")
	for _, t := range tokens {
		normalized := strings.Title(strings.ToLower(t))
		if t == "Alt" {
			modifier = ModAlt
			continue
		}
		cleaned = append(cleaned, normalized)
	}

	key, exist := translate[strings.Join(cleaned, "")]
	if !exist {
		return nil, ModNone, ErrNoSuchKeybind
	}

	return key, modifier, nil
}

// ParseAll takes an array of strings and returns a map of all keybindings.
func ParseAll(input []string) (map[interface{}]Modifier, error) {
	ret := make(map[interface{}]Modifier)
	for _, i := range input {
		k, m, err := Parse(i)
		if err != nil {
			return ret, err
		}
		ret[k] = m
	}
	return ret, nil
}

// MustParse takes the input string and returns a Key / rune and a Modifier.
// It will panic if any error occured.
func MustParse(input string) (interface{}, Modifier) {
	k, m, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return k, m
}

// MustParseAll takes an array of strings and returns a map of all keybindings.
// It will panic if any error occured.
func MustParseAll(input []string) map[interface{}]Modifier {
	result, err := ParseAll(input)
	if err != nil {
		panic(err)
	}
	return result
}

// newKeybinding returns a new Keybinding object.
func newKeybinding(viewname string, key Key, ch rune, mod Modifier, handler func(*Gui, *View) error) (kb *keybinding) {
	kb = &keybinding{
		viewName: viewname,
		key:      key,
		ch:       ch,
		mod:      mod,
		handler:  handler,
	}
	return kb
}

func eventMatchesKey(ev *GocuiEvent, key interface{}) bool {
	// assuming ModNone for now
	if Modifier(ev.Mod) != ModNone {
		return false
	}

	k, ch, err := getKey(key)
	if err != nil {
		return false
	}

	return k == Key(ev.Key) && ch == ev.Ch
}

// matchKeypress returns if the keybinding matches the keypress.
func (kb *keybinding) matchKeypress(key Key, ch rune, mod Modifier) bool {
	return kb.key == key && kb.ch == ch && kb.mod == mod
}

// translations for strings to keys
var translate = map[string]Key{
	"F1":             KeyF1,
	"F2":             KeyF2,
	"F3":             KeyF3,
	"F4":             KeyF4,
	"F5":             KeyF5,
	"F6":             KeyF6,
	"F7":             KeyF7,
	"F8":             KeyF8,
	"F9":             KeyF9,
	"F10":            KeyF10,
	"F11":            KeyF11,
	"F12":            KeyF12,
	"Insert":         KeyInsert,
	"Delete":         KeyDelete,
	"Home":           KeyHome,
	"End":            KeyEnd,
	"Pgup":           KeyPgup,
	"Pgdn":           KeyPgdn,
	"ArrowUp":        KeyArrowUp,
	"ShiftArrowUp":   KeyShiftArrowUp,
	"ArrowDown":      KeyArrowDown,
	"ShiftArrowDown": KeyShiftArrowDown,
	"ArrowLeft":      KeyArrowLeft,
	"ArrowRight":     KeyArrowRight,
	"CtrlTilde":      KeyCtrlTilde,
	"Ctrl2":          KeyCtrl2,
	"CtrlSpace":      KeyCtrlSpace,
	"CtrlA":          KeyCtrlA,
	"CtrlB":          KeyCtrlB,
	"CtrlC":          KeyCtrlC,
	"CtrlD":          KeyCtrlD,
	"CtrlE":          KeyCtrlE,
	"CtrlF":          KeyCtrlF,
	"CtrlG":          KeyCtrlG,
	"Backspace":      KeyBackspace,
	"CtrlH":          KeyCtrlH,
	"Tab":            KeyTab,
	"BackTab":        KeyBacktab,
	"CtrlI":          KeyCtrlI,
	"CtrlJ":          KeyCtrlJ,
	"CtrlK":          KeyCtrlK,
	"CtrlL":          KeyCtrlL,
	"Enter":          KeyEnter,
	"CtrlM":          KeyCtrlM,
	"CtrlN":          KeyCtrlN,
	"CtrlO":          KeyCtrlO,
	"CtrlP":          KeyCtrlP,
	"CtrlQ":          KeyCtrlQ,
	"CtrlR":          KeyCtrlR,
	"CtrlS":          KeyCtrlS,
	"CtrlT":          KeyCtrlT,
	"CtrlU":          KeyCtrlU,
	"CtrlV":          KeyCtrlV,
	"CtrlW":          KeyCtrlW,
	"CtrlX":          KeyCtrlX,
	"CtrlY":          KeyCtrlY,
	"CtrlZ":          KeyCtrlZ,
	"Esc":            KeyEsc,
	"CtrlLsqBracket": KeyCtrlLsqBracket,
	"Ctrl3":          KeyCtrl3,
	"Ctrl4":          KeyCtrl4,
	"CtrlBackslash":  KeyCtrlBackslash,
	"Ctrl5":          KeyCtrl5,
	"CtrlRsqBracket": KeyCtrlRsqBracket,
	"Ctrl6":          KeyCtrl6,
	"Ctrl7":          KeyCtrl7,
	"CtrlSlash":      KeyCtrlSlash,
	"CtrlUnderscore": KeyCtrlUnderscore,
	"Space":          KeySpace,
	"Backspace2":     KeyBackspace2,
	"Ctrl8":          KeyCtrl8,
	"Mouseleft":      MouseLeft,
	"Mousemiddle":    MouseMiddle,
	"Mouseright":     MouseRight,
	"Mouserelease":   MouseRelease,
	"MousewheelUp":   MouseWheelUp,
	"MousewheelDown": MouseWheelDown,
}

// Special keys.
const (
	KeyF1             Key = Key(tcell.KeyF1)
	KeyF2                 = Key(tcell.KeyF2)
	KeyF3                 = Key(tcell.KeyF3)
	KeyF4                 = Key(tcell.KeyF4)
	KeyF5                 = Key(tcell.KeyF5)
	KeyF6                 = Key(tcell.KeyF6)
	KeyF7                 = Key(tcell.KeyF7)
	KeyF8                 = Key(tcell.KeyF8)
	KeyF9                 = Key(tcell.KeyF9)
	KeyF10                = Key(tcell.KeyF10)
	KeyF11                = Key(tcell.KeyF11)
	KeyF12                = Key(tcell.KeyF12)
	KeyInsert             = Key(tcell.KeyInsert)
	KeyDelete             = Key(tcell.KeyDelete)
	KeyHome               = Key(tcell.KeyHome)
	KeyEnd                = Key(tcell.KeyEnd)
	KeyPgdn               = Key(tcell.KeyPgDn)
	KeyPgup               = Key(tcell.KeyPgUp)
	KeyArrowUp            = Key(tcell.KeyUp)
	KeyShiftArrowUp       = Key(tcell.KeyF62)
	KeyArrowDown          = Key(tcell.KeyDown)
	KeyShiftArrowDown     = Key(tcell.KeyF63)
	KeyArrowLeft          = Key(tcell.KeyLeft)
	KeyArrowRight         = Key(tcell.KeyRight)
)

// Keys combinations.
const (
	KeyCtrlTilde      = Key(tcell.KeyF64) // arbitrary assignment
	KeyCtrlSpace      = Key(tcell.KeyCtrlSpace)
	KeyCtrlA          = Key(tcell.KeyCtrlA)
	KeyCtrlB          = Key(tcell.KeyCtrlB)
	KeyCtrlC          = Key(tcell.KeyCtrlC)
	KeyCtrlD          = Key(tcell.KeyCtrlD)
	KeyCtrlE          = Key(tcell.KeyCtrlE)
	KeyCtrlF          = Key(tcell.KeyCtrlF)
	KeyCtrlG          = Key(tcell.KeyCtrlG)
	KeyBackspace      = Key(tcell.KeyBackspace)
	KeyCtrlH          = Key(tcell.KeyCtrlH)
	KeyTab            = Key(tcell.KeyTab)
	KeyBacktab        = Key(tcell.KeyBacktab)
	KeyCtrlI          = Key(tcell.KeyCtrlI)
	KeyCtrlJ          = Key(tcell.KeyCtrlJ)
	KeyCtrlK          = Key(tcell.KeyCtrlK)
	KeyCtrlL          = Key(tcell.KeyCtrlL)
	KeyEnter          = Key(tcell.KeyEnter)
	KeyCtrlM          = Key(tcell.KeyCtrlM)
	KeyCtrlN          = Key(tcell.KeyCtrlN)
	KeyCtrlO          = Key(tcell.KeyCtrlO)
	KeyCtrlP          = Key(tcell.KeyCtrlP)
	KeyCtrlQ          = Key(tcell.KeyCtrlQ)
	KeyCtrlR          = Key(tcell.KeyCtrlR)
	KeyCtrlS          = Key(tcell.KeyCtrlS)
	KeyCtrlT          = Key(tcell.KeyCtrlT)
	KeyCtrlU          = Key(tcell.KeyCtrlU)
	KeyCtrlV          = Key(tcell.KeyCtrlV)
	KeyCtrlW          = Key(tcell.KeyCtrlW)
	KeyCtrlX          = Key(tcell.KeyCtrlX)
	KeyCtrlY          = Key(tcell.KeyCtrlY)
	KeyCtrlZ          = Key(tcell.KeyCtrlZ)
	KeyEsc            = Key(tcell.KeyEscape)
	KeyCtrlUnderscore = Key(tcell.KeyCtrlUnderscore)
	KeySpace          = Key(32)
	KeyBackspace2     = Key(tcell.KeyBackspace2)
	KeyCtrl8          = Key(tcell.KeyBackspace2) // same key as in termbox-go

	// The following assignments were used in termbox implementation.
	// In tcell, these are not keys per se. But in gocui we have them
	// mapped to the keys so we have to use placeholder keys.

	KeyAltEnter       = Key(tcell.KeyF64) // arbitrary assignments
	MouseLeft         = Key(tcell.KeyF63)
	MouseRight        = Key(tcell.KeyF62)
	MouseMiddle       = Key(tcell.KeyF61)
	MouseRelease      = Key(tcell.KeyF60)
	MouseWheelUp      = Key(tcell.KeyF59)
	MouseWheelDown    = Key(tcell.KeyF58)
	MouseWheelLeft    = Key(tcell.KeyF57)
	MouseWheelRight   = Key(tcell.KeyF56)
	KeyCtrl2          = Key(tcell.KeyNUL) // termbox defines theses
	KeyCtrl3          = Key(tcell.KeyEscape)
	KeyCtrl4          = Key(tcell.KeyCtrlBackslash)
	KeyCtrl5          = Key(tcell.KeyCtrlRightSq)
	KeyCtrl6          = Key(tcell.KeyCtrlCarat)
	KeyCtrl7          = Key(tcell.KeyCtrlUnderscore)
	KeyCtrlSlash      = Key(tcell.KeyCtrlUnderscore)
	KeyCtrlRsqBracket = Key(tcell.KeyCtrlRightSq)
	KeyCtrlBackslash  = Key(tcell.KeyCtrlBackslash)
	KeyCtrlLsqBracket = Key(tcell.KeyCtrlLeftSq)
)

// Modifiers.
const (
	ModNone   Modifier = Modifier(0)
	ModAlt             = Modifier(tcell.ModAlt)
	ModMotion          = Modifier(2) // just picking an arbitrary number here that doesn't clash with tcell.ModAlt
	// ModCtrl doesn't work with keyboard keys. Use CtrlKey in Key and ModNone. This is was for mouse clicks only (tcell.v1)
	// ModCtrl = Modifier(tcell.ModCtrl)
)
