// Copyright 2023 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm
// +build js,wasm

package tcell

import (
	"errors"
	"strings"
	"sync"
	"syscall/js"
	"unicode/utf8"
)

func NewTerminfoScreen() (Screen, error) {
	t := &wScreen{}
	t.fallback = make(map[rune]string)

	return t, nil
}

type wScreen struct {
	w, h  int
	style Style
	cells CellBuffer

	running      bool
	clear        bool
	flagsPresent bool
	pasteEnabled bool
	mouseFlags   MouseFlags

	cursorStyle CursorStyle

	quit     chan struct{}
	evch     chan Event
	fallback map[rune]string

	sync.Mutex
}

func (t *wScreen) Init() error {
	t.w, t.h = 80, 24 // default for html as of now
	t.evch = make(chan Event, 10)
	t.quit = make(chan struct{})

	t.Lock()
	t.running = true
	t.style = StyleDefault
	t.cells.Resize(t.w, t.h)
	t.Unlock()

	js.Global().Set("onKeyEvent", js.FuncOf(t.onKeyEvent))

	return nil
}

func (t *wScreen) Fini() {
	close(t.quit)
}

func (t *wScreen) SetStyle(style Style) {
	t.Lock()
	t.style = style
	t.Unlock()
}

func (t *wScreen) Clear() {
	t.Fill(' ', t.style)
}

func (t *wScreen) Fill(r rune, style Style) {
	t.Lock()
	t.cells.Fill(r, style)
	t.Unlock()
}

func (t *wScreen) SetContent(x, y int, mainc rune, combc []rune, style Style) {
	t.Lock()
	t.cells.SetContent(x, y, mainc, combc, style)
	t.Unlock()
}

func (t *wScreen) GetContent(x, y int) (rune, []rune, Style, int) {
	t.Lock()
	mainc, combc, style, width := t.cells.GetContent(x, y)
	t.Unlock()
	return mainc, combc, style, width
}

func (t *wScreen) SetCell(x, y int, style Style, ch ...rune) {
	if len(ch) > 0 {
		t.SetContent(x, y, ch[0], ch[1:], style)
	} else {
		t.SetContent(x, y, ' ', nil, style)
	}
}

// paletteColor gives a more natural palette color actually matching
// typical XTerm.  We might in the future want to permit styling these
// via CSS.

var palette = map[Color]int32{
	ColorBlack: 0x000000,
	ColorMaroon: 0xcd0000,
	ColorGreen: 0x00cd00,
	ColorOlive: 0xcdcd00,
	ColorNavy: 0x0000ee,
	ColorPurple: 0xcd00cd,
	ColorTeal: 0x00cdcd,
	ColorSilver: 0xe5e5e5,
	ColorGray: 0x7f7f7f,
	ColorRed: 0xff0000,
	ColorLime: 0x00ff00,
	ColorYellow: 0xffff00,
	ColorBlue: 0x5c5cff,
	ColorFuchsia: 0xff00ff,
	ColorAqua: 0x00ffff,
	ColorWhite: 0xffffff,
}

func paletteColor(c Color) int32 {
	if (c.IsRGB()) {
		return int32(c & 0xffffff);
	}
	if (c >= ColorBlack && c <= ColorWhite) {
		return palette[c]
	}
	return c.Hex()
}

func (t *wScreen) drawCell(x, y int) int {
	mainc, combc, style, width := t.cells.GetContent(x, y)

	if !t.cells.Dirty(x, y) {
		return width
	}

	if style == StyleDefault {
		style = t.style
	}

	fg, bg := paletteColor(style.fg), paletteColor(style.bg)
	if (fg == -1) {
		fg = 0xe5e5e5;
	}
	if (bg == -1) {
		bg = 0x000000;
	}

	var combcarr []interface{} = make([]interface{}, len(combc))
	for i, c := range combc {
		combcarr[i] = c
	}

	t.cells.SetDirty(x, y, false)
	js.Global().Call("drawCell", x, y, mainc, combcarr, fg, bg, int(style.attrs))

	return width
}

func (t *wScreen) ShowCursor(x, y int) {
	t.Lock()
	js.Global().Call("showCursor", x, y)
	t.Unlock()
}

func (t *wScreen) SetCursorStyle(cs CursorStyle) {
	t.Lock()
	js.Global().Call("setCursorStyle", curStyleClasses[cs])
	t.Unlock()
}

func (t *wScreen) HideCursor() {
	t.ShowCursor(-1, -1)
}

func (t *wScreen) Show() {
	t.Lock()
	t.resize()
	t.draw()
	t.Unlock()
}

func (t *wScreen) clearScreen() {
	js.Global().Call("clearScreen", t.style.fg.Hex(), t.style.bg.Hex())
	t.clear = false
}

func (t *wScreen) draw() {
	if t.clear {
		t.clearScreen()
	}

	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			width := t.drawCell(x, y)
			x += width - 1
		}
	}

	js.Global().Call("show")
}

func (t *wScreen) EnableMouse(flags ...MouseFlags) {
	var f MouseFlags
	flagsPresent := false
	for _, flag := range flags {
		f |= flag
		flagsPresent = true
	}
	if !flagsPresent {
		f = MouseMotionEvents | MouseDragEvents | MouseButtonEvents
	}

	t.Lock()
	t.mouseFlags = f
	t.enableMouse(f)
	t.Unlock()
}

func (t *wScreen) enableMouse(f MouseFlags) {
	if f&MouseButtonEvents != 0 {
		js.Global().Set("onMouseClick", js.FuncOf(t.onMouseEvent))
	} else {
		js.Global().Set("onMouseClick", js.FuncOf(t.unset))
	}

	if f&MouseDragEvents != 0 || f&MouseMotionEvents != 0 {
		js.Global().Set("onMouseMove", js.FuncOf(t.onMouseEvent))
	} else {
		js.Global().Set("onMouseMove", js.FuncOf(t.unset))
	}
}

func (t *wScreen) DisableMouse() {
	t.Lock()
	t.mouseFlags = 0
	t.enableMouse(0)
	t.Unlock()
}

func (t *wScreen) EnablePaste() {
	t.Lock()
	t.pasteEnabled = true
	t.enablePasting(true)
	t.Unlock()
}

func (t *wScreen) DisablePaste() {
	t.Lock()
	t.pasteEnabled = false
	t.enablePasting(false)
	t.Unlock()
}

func (t *wScreen) enablePasting(on bool) {
	if on {
		js.Global().Set("onPaste", js.FuncOf(t.onPaste))
	} else {
		js.Global().Set("onPaste", js.FuncOf(t.unset))
	}
}

func (t *wScreen) Size() (int, int) {
	t.Lock()
	w, h := t.w, t.h
	t.Unlock()
	return w, h
}

// resize does nothing, as asking the web window to resize
// without a specified width or height will cause no change.
func (t *wScreen) resize() {}

func (t *wScreen) Colors() int {
	return 16777216 // 256 ^ 3
}

func (t *wScreen) ChannelEvents(ch chan<- Event, quit <-chan struct{}) {
	defer close(ch)
	for {
		select {
		case <-quit:
			return
		case <-t.quit:
			return
		case ev := <-t.evch:
			select {
			case <-quit:
				return
			case <-t.quit:
				return
			case ch <- ev:
			}
		}
	}
}

func (t *wScreen) PollEvent() Event {
	select {
	case <-t.quit:
		return nil
	case ev := <-t.evch:
		return ev
	}
}

func (t *wScreen) HasPendingEvent() bool {
	return len(t.evch) > 0
}

func (t *wScreen) PostEventWait(ev Event) {
	t.evch <- ev
}

func (t *wScreen) PostEvent(ev Event) error {
	select {
	case t.evch <- ev:
		return nil
	default:
		return ErrEventQFull
	}
}

func (t *wScreen) clip(x, y int) (int, int) {
	w, h := t.cells.Size()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > w-1 {
		x = w - 1
	}
	if y > h-1 {
		y = h - 1
	}
	return x, y
}

func (t *wScreen) onMouseEvent(this js.Value, args []js.Value) interface{} {
	mod := ModNone
	button := ButtonNone

	switch args[2].Int() {
	case 0:
		if t.mouseFlags&MouseMotionEvents == 0 {
			// don't want this event! is a mouse motion event, but user has asked not.
			return nil
		}
		button = ButtonNone
	case 1:
		button = Button1
	case 2:
		button = Button3 // Note we prefer to treat right as button 2
	case 3:
		button = Button2 // And the middle button as button 3
	}

	if args[3].Bool() { // mod shift
		mod |= ModShift
	}

	if args[4].Bool() { // mod alt
		mod |= ModAlt
	}

	if args[5].Bool() { // mod ctrl
		mod |= ModCtrl
	}

	t.PostEventWait(NewEventMouse(args[0].Int(), args[1].Int(), button, mod))
	return nil
}

func (t *wScreen) onKeyEvent(this js.Value, args []js.Value) interface{} {
	key := args[0].String()

	// don't accept any modifier keys as their own
	if key == "Control" || key == "Alt" || key == "Meta" || key == "Shift" {
		return nil
	}

	mod := ModNone
	if args[1].Bool() { // mod shift
		mod |= ModShift
	}

	if args[2].Bool() { // mod alt
		mod |= ModAlt
	}

	if args[3].Bool() { // mod ctrl
		mod |= ModCtrl
	}

	if args[4].Bool() { // mod meta
		mod |= ModMeta
	}

	// check for special case of Ctrl + key
	if mod == ModCtrl {
		if k, ok := WebKeyNames["Ctrl-"+strings.ToLower(key)]; ok {
			t.PostEventWait(NewEventKey(k, 0, mod))
			return nil
		}
	}

	// next try function keys
	if k, ok := WebKeyNames[key]; ok {
		t.PostEventWait(NewEventKey(k, 0, mod))
		return nil
	}

	// finally try normal, printable chars
	r, _ := utf8.DecodeRuneInString(key)
	t.PostEventWait(NewEventKey(KeyRune, r, mod))
	return nil
}

func (t *wScreen) onPaste(this js.Value, args []js.Value) interface{} {
	t.PostEventWait(NewEventPaste(args[0].Bool()))
	return nil
}

// unset is a dummy function for js when we want nothing to
// happen when javascript calls a function (for example, when
// mouse input is disabled, when onMouseEvent() is called from
// js, it redirects here and does nothing).
func (t *wScreen) unset(this js.Value, args []js.Value) interface{} {
	return nil
}

func (t *wScreen) Sync() {
	t.Lock()
	t.resize()
	t.clear = true
	t.cells.Invalidate()
	t.draw()
	t.Unlock()
}

func (t *wScreen) CharacterSet() string {
	return "UTF-8"
}

func (t *wScreen) RegisterRuneFallback(orig rune, fallback string) {
	t.Lock()
	t.fallback[orig] = fallback
	t.Unlock()
}

func (t *wScreen) UnregisterRuneFallback(orig rune) {
	t.Lock()
	delete(t.fallback, orig)
	t.Unlock()
}

func (t *wScreen) CanDisplay(r rune, checkFallbacks bool) bool {
	if utf8.ValidRune(r) {
		return true
	}
	if !checkFallbacks {
		return false
	}
	if _, ok := t.fallback[r]; ok {
		return true
	}
	return false
}

func (t *wScreen) HasMouse() bool {
	return true
}

func (t *wScreen) HasKey(k Key) bool {
	return true
}

func (t *wScreen) SetSize(w, h int) {
	if w == t.w && h == t.h {
		return
	}

	t.cells.Invalidate()
	t.cells.Resize(w, h)
	js.Global().Call("resize", w, h)
	t.w, t.h = w, h
	t.PostEvent(NewEventResize(w, h))
}

func (t *wScreen) Resize(int, int, int, int) {}

// Suspend simply pauses all input and output, and clears the screen.
// There isn't a "default terminal" to go back to.
func (t *wScreen) Suspend() error {
	t.Lock()
	if !t.running {
		t.Unlock()
		return nil
	}
	t.running = false
	t.clearScreen()
	t.enableMouse(0)
	t.enablePasting(false)
	js.Global().Set("onKeyEvent", js.FuncOf(t.unset)) // stop keypresses
	return nil
}

func (t *wScreen) Resume() error {
	t.Lock()

	if t.running {
		return errors.New("already engaged")
	}
	t.running = true

	t.enableMouse(t.mouseFlags)
	t.enablePasting(t.pasteEnabled)

	js.Global().Set("onKeyEvent", js.FuncOf(t.onKeyEvent))

	t.Unlock()
	return nil
}

func (t *wScreen) Beep() error {
	js.Global().Call("beep")
	return nil
}

// WebKeyNames maps string names reported from HTML
// (KeyboardEvent.key) to tcell accepted keys.
var WebKeyNames = map[string]Key{
	"Enter":      KeyEnter,
	"Backspace":  KeyBackspace,
	"Tab":        KeyTab,
	"Backtab":    KeyBacktab,
	"Escape":     KeyEsc,
	"Backspace2": KeyBackspace2,
	"Delete":     KeyDelete,
	"Insert":     KeyInsert,
	"ArrowUp":    KeyUp,
	"ArrowDown":  KeyDown,
	"ArrowLeft":  KeyLeft,
	"ArrowRight": KeyRight,
	"Home":       KeyHome,
	"End":        KeyEnd,
	"UpLeft":     KeyUpLeft,    // not supported by HTML
	"UpRight":    KeyUpRight,   // not supported by HTML
	"DownLeft":   KeyDownLeft,  // not supported by HTML
	"DownRight":  KeyDownRight, // not supported by HTML
	"Center":     KeyCenter,
	"PgDn":       KeyPgDn,
	"PgUp":       KeyPgUp,
	"Clear":      KeyClear,
	"Exit":       KeyExit,
	"Cancel":     KeyCancel,
	"Pause":      KeyPause,
	"Print":      KeyPrint,
	"F1":         KeyF1,
	"F2":         KeyF2,
	"F3":         KeyF3,
	"F4":         KeyF4,
	"F5":         KeyF5,
	"F6":         KeyF6,
	"F7":         KeyF7,
	"F8":         KeyF8,
	"F9":         KeyF9,
	"F10":        KeyF10,
	"F11":        KeyF11,
	"F12":        KeyF12,
	"F13":        KeyF13,
	"F14":        KeyF14,
	"F15":        KeyF15,
	"F16":        KeyF16,
	"F17":        KeyF17,
	"F18":        KeyF18,
	"F19":        KeyF19,
	"F20":        KeyF20,
	"F21":        KeyF21,
	"F22":        KeyF22,
	"F23":        KeyF23,
	"F24":        KeyF24,
	"F25":        KeyF25,
	"F26":        KeyF26,
	"F27":        KeyF27,
	"F28":        KeyF28,
	"F29":        KeyF29,
	"F30":        KeyF30,
	"F31":        KeyF31,
	"F32":        KeyF32,
	"F33":        KeyF33,
	"F34":        KeyF34,
	"F35":        KeyF35,
	"F36":        KeyF36,
	"F37":        KeyF37,
	"F38":        KeyF38,
	"F39":        KeyF39,
	"F40":        KeyF40,
	"F41":        KeyF41,
	"F42":        KeyF42,
	"F43":        KeyF43,
	"F44":        KeyF44,
	"F45":        KeyF45,
	"F46":        KeyF46,
	"F47":        KeyF47,
	"F48":        KeyF48,
	"F49":        KeyF49,
	"F50":        KeyF50,
	"F51":        KeyF51,
	"F52":        KeyF52,
	"F53":        KeyF53,
	"F54":        KeyF54,
	"F55":        KeyF55,
	"F56":        KeyF56,
	"F57":        KeyF57,
	"F58":        KeyF58,
	"F59":        KeyF59,
	"F60":        KeyF60,
	"F61":        KeyF61,
	"F62":        KeyF62,
	"F63":        KeyF63,
	"F64":        KeyF64,
	"Ctrl-a":     KeyCtrlA,          // not reported by HTML- need to do special check
	"Ctrl-b":     KeyCtrlB,          // not reported by HTML- need to do special check
	"Ctrl-c":     KeyCtrlC,          // not reported by HTML- need to do special check
	"Ctrl-d":     KeyCtrlD,          // not reported by HTML- need to do special check
	"Ctrl-e":     KeyCtrlE,          // not reported by HTML- need to do special check
	"Ctrl-f":     KeyCtrlF,          // not reported by HTML- need to do special check
	"Ctrl-g":     KeyCtrlG,          // not reported by HTML- need to do special check
	"Ctrl-j":     KeyCtrlJ,          // not reported by HTML- need to do special check
	"Ctrl-k":     KeyCtrlK,          // not reported by HTML- need to do special check
	"Ctrl-l":     KeyCtrlL,          // not reported by HTML- need to do special check
	"Ctrl-n":     KeyCtrlN,          // not reported by HTML- need to do special check
	"Ctrl-o":     KeyCtrlO,          // not reported by HTML- need to do special check
	"Ctrl-p":     KeyCtrlP,          // not reported by HTML- need to do special check
	"Ctrl-q":     KeyCtrlQ,          // not reported by HTML- need to do special check
	"Ctrl-r":     KeyCtrlR,          // not reported by HTML- need to do special check
	"Ctrl-s":     KeyCtrlS,          // not reported by HTML- need to do special check
	"Ctrl-t":     KeyCtrlT,          // not reported by HTML- need to do special check
	"Ctrl-u":     KeyCtrlU,          // not reported by HTML- need to do special check
	"Ctrl-v":     KeyCtrlV,          // not reported by HTML- need to do special check
	"Ctrl-w":     KeyCtrlW,          // not reported by HTML- need to do special check
	"Ctrl-x":     KeyCtrlX,          // not reported by HTML- need to do special check
	"Ctrl-y":     KeyCtrlY,          // not reported by HTML- need to do special check
	"Ctrl-z":     KeyCtrlZ,          // not reported by HTML- need to do special check
	"Ctrl- ":     KeyCtrlSpace,      // not reported by HTML- need to do special check
	"Ctrl-_":     KeyCtrlUnderscore, // not reported by HTML- need to do special check
	"Ctrl-]":     KeyCtrlRightSq,    // not reported by HTML- need to do special check
	"Ctrl-\\":    KeyCtrlBackslash,  // not reported by HTML- need to do special check
	"Ctrl-^":     KeyCtrlCarat,      // not reported by HTML- need to do special check
}

var curStyleClasses = map[CursorStyle]string{
	CursorStyleDefault:           "cursor-blinking-block",
	CursorStyleBlinkingBlock:     "cursor-blinking-block",
	CursorStyleSteadyBlock:       "cursor-steady-block",
	CursorStyleBlinkingUnderline: "cursor-blinking-underline",
	CursorStyleSteadyUnderline:   "cursor-steady-underline",
	CursorStyleBlinkingBar:       "cursor-blinking-bar",
	CursorStyleSteadyBar:         "cursor-steady-bar",
}
