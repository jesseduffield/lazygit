// Copyright 2020 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// We probably don't want this being a global variable for YOLO for now
var Screen tcell.Screen

// oldStyle is a representation of how a cell would be styled when we were using termbox
type oldStyle struct {
	fg         Attribute
	bg         Attribute
	outputMode OutputMode
}

var runeReplacements = map[rune]string{
	'┌': "+",
	'┐': "+",
	'└': "+",
	'┘': "+",
	'╭': "+",
	'╮': "+",
	'╰': "+",
	'╯': "+",
	'─': "-",
	'═': "-",
	'║': "|",
	'╔': "+",
	'╗': "+",
	'╚': "+",
	'╝': "+",

	// using a hyphen here actually looks weird.
	// We see these characters when in portrait mode
	'╶': " ",
	'╴': " ",

	'┴': "+",
	'┬': "+",
	'╷': "|",
	'├': "+",
	'│': "|",
	'▼': "v",
	'►': ">",
	'▲': "^",
	'◄': "<",
}

// tcellInit initializes tcell screen for use.
func (g *Gui) tcellInit(runeReplacements map[rune]string) error {
	runewidth.DefaultCondition.EastAsianWidth = false
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if s, e := tcell.NewScreen(); e != nil {
		return e
	} else if e = s.Init(); e != nil {
		return e
	} else {
		registerRuneFallbacks(s, runeReplacements)

		g.screen = s
		Screen = s
		return nil
	}
}

func registerRuneFallbacks(s tcell.Screen, additional map[rune]string) {
	for before, after := range runeReplacements {
		s.RegisterRuneFallback(before, after)
	}

	for before, after := range additional {
		s.RegisterRuneFallback(before, after)
	}
}

// tcellInitSimulation initializes tcell screen for use.
func (g *Gui) tcellInitSimulation(width int, height int) error {
	s := tcell.NewSimulationScreen("")
	if e := s.Init(); e != nil {
		return e
	} else {
		g.screen = s
		Screen = s
		// setting to a larger value than the typical terminal size
		// so that during a test we're more likely to see an item to select in a view.
		s.SetSize(width, height)
		s.Sync()
		return nil
	}
}

// tcellSetCell sets the character cell at a given location to the given
// content (rune) and attributes using provided OutputMode
func tcellSetCell(x, y int, ch rune, fg, bg Attribute, outputMode OutputMode) {
	st := getTcellStyle(oldStyle{fg: fg, bg: bg, outputMode: outputMode})
	Screen.SetContent(x, y, ch, nil, st)
}

// getTcellStyle creates tcell.Style from Attributes
func getTcellStyle(input oldStyle) tcell.Style {
	st := tcell.StyleDefault

	// extract colors and attributes
	if input.fg != ColorDefault {
		st = st.Foreground(getTcellColor(input.fg, input.outputMode))
		st = setTcellFontEffectStyle(st, input.fg)
	}
	if input.bg != ColorDefault {
		st = st.Background(getTcellColor(input.bg, input.outputMode))
		st = setTcellFontEffectStyle(st, input.bg)
	}

	return st
}

// setTcellFontEffectStyle add additional attributes to tcell.Style
func setTcellFontEffectStyle(st tcell.Style, attr Attribute) tcell.Style {
	if attr&AttrBold != 0 {
		st = st.Bold(true)
	}
	if attr&AttrUnderline != 0 {
		st = st.Underline(true)
	}
	if attr&AttrReverse != 0 {
		st = st.Reverse(true)
	}
	if attr&AttrBlink != 0 {
		st = st.Blink(true)
	}
	if attr&AttrDim != 0 {
		st = st.Dim(true)
	}
	if attr&AttrItalic != 0 {
		st = st.Italic(true)
	}
	if attr&AttrStrikeThrough != 0 {
		st = st.StrikeThrough(true)
	}
	return st
}

// gocuiEventType represents the type of event.
type gocuiEventType uint8

// GocuiEvent represents events like a keys, mouse actions, or window resize.
//
//	The 'Mod', 'Key' and 'Ch' fields are valid if 'Type' is 'eventKey'.
//	The 'MouseX' and 'MouseY' fields are valid if 'Type' is 'eventMouse'.
//	The 'Width' and 'Height' fields are valid if 'Type' is 'eventResize'.
//	The 'Err' field is valid if 'Type' is 'eventError'.
type GocuiEvent struct {
	Type   gocuiEventType
	Mod    Modifier
	Key    Key
	Ch     rune
	Width  int
	Height int
	Err    error
	MouseX int
	MouseY int
	N      int
}

// Event types.
const (
	eventNone gocuiEventType = iota
	eventKey
	eventResize
	eventMouse
	eventInterrupt
	eventError
	eventRaw
)

const (
	NOT_DRAGGING int = iota
	MAYBE_DRAGGING
	DRAGGING
)

var (
	lastMouseKey tcell.ButtonMask = tcell.ButtonNone
	lastMouseMod tcell.ModMask    = tcell.ModNone
	dragState    int              = NOT_DRAGGING
	lastX        int              = 0
	lastY        int              = 0
)

// this wrapper struct has public keys so we can easily serialize/deserialize to JSON
type TcellKeyEventWrapper struct {
	Timestamp int64
	Mod       tcell.ModMask
	Key       tcell.Key
	Ch        rune
}

func NewTcellKeyEventWrapper(event *tcell.EventKey, timestamp int64) *TcellKeyEventWrapper {
	return &TcellKeyEventWrapper{
		Timestamp: timestamp,
		Mod:       event.Modifiers(),
		Key:       event.Key(),
		Ch:        event.Rune(),
	}
}

func (wrapper TcellKeyEventWrapper) toTcellEvent() tcell.Event {
	return tcell.NewEventKey(wrapper.Key, wrapper.Ch, wrapper.Mod)
}

type TcellResizeEventWrapper struct {
	Timestamp int64
	Width     int
	Height    int
}

func NewTcellResizeEventWrapper(event *tcell.EventResize, timestamp int64) *TcellResizeEventWrapper {
	w, h := event.Size()

	return &TcellResizeEventWrapper{
		Timestamp: timestamp,
		Width:     w,
		Height:    h,
	}
}

func (wrapper TcellResizeEventWrapper) toTcellEvent() tcell.Event {
	return tcell.NewEventResize(wrapper.Width, wrapper.Height)
}

// pollEvent get tcell.Event and transform it into gocuiEvent
func (g *Gui) pollEvent() GocuiEvent {
	var tev tcell.Event
	if g.playRecording {
		select {
		case ev := <-g.ReplayedEvents.Keys:
			tev = (ev).toTcellEvent()
		case ev := <-g.ReplayedEvents.Resizes:
			tev = (ev).toTcellEvent()
		}
	} else {
		tev = Screen.PollEvent()
	}

	switch tev := tev.(type) {
	case *tcell.EventInterrupt:
		return GocuiEvent{Type: eventInterrupt}
	case *tcell.EventResize:
		w, h := tev.Size()
		return GocuiEvent{Type: eventResize, Width: w, Height: h}
	case *tcell.EventKey:
		k := tev.Key()
		ch := rune(0)
		if k == tcell.KeyRune {
			k = 0 // if rune remove key (so it can match rune instead of key)
			ch = tev.Rune()
			if ch == ' ' {
				// special handling for spacebar
				k = 32 // tcell keys ends at 31 or starts at 256
				ch = rune(0)
			}
		}
		mod := tev.Modifiers()
		// remove control modifier and setup special handling of ctrl+spacebar, etc.
		if mod == tcell.ModCtrl && k == 32 {
			mod = 0
			ch = rune(0)
			k = tcell.KeyCtrlSpace
		} else if mod == tcell.ModCtrl || mod == tcell.ModShift {
			// remove Ctrl or Shift if specified
			// - shift - will be translated to the final code of rune
			// - ctrl  - is translated in the key
			mod = 0
		} else if mod == tcell.ModAlt && k == tcell.KeyEnter {
			// for the sake of convenience I'm having a KeyAltEnter key. I will likely
			// regret this laziness in the future. We're arbitrarily mapping that to tcell's
			// KeyF64.
			mod = 0
			k = tcell.KeyF64
		}

		return GocuiEvent{
			Type: eventKey,
			Key:  Key(k),
			Ch:   ch,
			Mod:  Modifier(mod),
		}
	case *tcell.EventMouse:
		x, y := tev.Position()
		button := tev.Buttons()
		mouseKey := MouseRelease
		mouseMod := ModNone
		// process mouse wheel
		if button&tcell.WheelUp != 0 {
			mouseKey = MouseWheelUp
		}
		if button&tcell.WheelDown != 0 {
			mouseKey = MouseWheelDown
		}
		if button&tcell.WheelLeft != 0 {
			mouseKey = MouseWheelLeft
		}
		if button&tcell.WheelRight != 0 {
			mouseKey = MouseWheelRight
		}

		wheeling := mouseKey == MouseWheelUp || mouseKey == MouseWheelDown || mouseKey == MouseWheelLeft || mouseKey == MouseWheelRight

		// process button events (not wheel events)
		button &= tcell.ButtonMask(0xff)
		if button != tcell.ButtonNone && lastMouseKey == tcell.ButtonNone {
			lastMouseKey = button
			lastMouseMod = tev.Modifiers()
			switch button {
			case tcell.ButtonPrimary:
				mouseKey = MouseLeft
				dragState = MAYBE_DRAGGING
				lastX = x
				lastY = y
			case tcell.ButtonSecondary:
				mouseKey = MouseRight
			case tcell.ButtonMiddle:
				mouseKey = MouseMiddle
			}
		}

		switch tev.Buttons() {
		case tcell.ButtonNone:
			if lastMouseKey != tcell.ButtonNone {
				switch lastMouseKey {
				case tcell.ButtonPrimary:
					dragState = NOT_DRAGGING
				case tcell.ButtonSecondary:
				case tcell.ButtonMiddle:
				}
				mouseMod = Modifier(lastMouseMod)
				lastMouseMod = tcell.ModNone
				lastMouseKey = tcell.ButtonNone
			}
		}

		if !wheeling {
			switch dragState {
			case NOT_DRAGGING:
				return GocuiEvent{Type: eventNone}
			// if we haven't released the left mouse button and we've moved the cursor then we're dragging
			case MAYBE_DRAGGING:
				if x != lastX || y != lastY {
					dragState = DRAGGING
				}
			case DRAGGING:
				mouseMod = ModMotion
				mouseKey = MouseLeft
			}
		}

		return GocuiEvent{
			Type:   eventMouse,
			MouseX: x,
			MouseY: y,
			Key:    mouseKey,
			Ch:     0,
			Mod:    mouseMod,
		}
	default:
		return GocuiEvent{Type: eventNone}
	}
}
