// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	standardErrors "errors"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/go-errors/errors"
	"github.com/mattn/go-runewidth"
)

// OutputMode represents an output mode, which determines how colors
// are used.
type OutputMode int

var (
	// ErrAlreadyBlacklisted is returned when the keybinding is already blacklisted.
	ErrAlreadyBlacklisted = standardErrors.New("keybind already blacklisted")

	// ErrBlacklisted is returned when the keybinding being parsed / used is blacklisted.
	ErrBlacklisted = standardErrors.New("keybind blacklisted")

	// ErrNotBlacklisted is returned when a keybinding being whitelisted is not blacklisted.
	ErrNotBlacklisted = standardErrors.New("keybind not blacklisted")

	// ErrNoSuchKeybind is returned when the keybinding being parsed does not exist.
	ErrNoSuchKeybind = standardErrors.New("no such keybind")

	// ErrUnknownView allows to assert if a View must be initialized.
	ErrUnknownView = standardErrors.New("unknown view")

	// ErrQuit is used to decide if the MainLoop finished successfully.
	ErrQuit = standardErrors.New("quit")
)

const (
	// OutputNormal provides 8-colors terminal mode.
	OutputNormal OutputMode = iota

	// Output256 provides 256-colors terminal mode.
	Output256

	// Output216 provides 216 ansi color terminal mode.
	Output216

	// OutputGrayscale provides greyscale terminal mode.
	OutputGrayscale

	// OutputTrue provides 24bit color terminal mode.
	// This mode is recommended even if your terminal doesn't support
	// such mode. The colors are represented exactly as you
	// write them (no clamping or truncating). `tcell` should take care
	// of what your terminal can do.
	OutputTrue
)

type tabClickHandler func(int) error

type tabClickBinding struct {
	viewName string
	handler  tabClickHandler
}

// TODO: would be good to define inbound and outbound click handlers e.g.
// clicking on a file is an inbound thing where we don't care what context you're
// in when it happens, whereas clicking on the main view from the files view is an
// outbound click with a specific handler. But this requires more thinking about
// where handlers should live.
type ViewMouseBinding struct {
	// the view that is clicked
	ViewName string

	// the view that has focus when the click occurs.
	FocusedView string

	Handler func(ViewMouseBindingOpts) error

	Modifier Modifier

	// must be a mouse key
	Key Key
}

type ViewMouseBindingOpts struct {
	X int // i.e. origin x + cursor x
	Y int // i.e. origin y + cursor y
}

type GuiMutexes struct {
	// tickingMutex ensures we don't have two loops ticking. The point of 'ticking'
	// is to refresh the gui rapidly so that loader characters can be animated.
	tickingMutex sync.Mutex

	ViewsMutex sync.Mutex
}

type PlayMode int

const (
	NORMAL PlayMode = iota
	RECORDING
	REPLAYING
	// for the new form of integration tests
	REPLAYING_NEW
)

type Recording struct {
	KeyEvents    []*TcellKeyEventWrapper
	ResizeEvents []*TcellResizeEventWrapper
}

type replayedEvents struct {
	Keys    chan *TcellKeyEventWrapper
	Resizes chan *TcellResizeEventWrapper
}

type RecordingConfig struct {
	Speed  float64
	Leeway int
}

// Gui represents the whole User Interface, including the views, layouts
// and keybindings.
type Gui struct {
	RecordingConfig
	Recording *Recording
	// ReplayedEvents is for passing pre-recorded input events, for the purposes of testing
	ReplayedEvents replayedEvents
	PlayMode       PlayMode
	StartTime      time.Time

	tabClickBindings  []*tabClickBinding
	viewMouseBindings []*ViewMouseBinding
	gEvents           chan GocuiEvent
	userEvents        chan userEvent
	views             []*View
	currentView       *View
	managers          []Manager
	keybindings       []*keybinding
	maxX, maxY        int
	outputMode        OutputMode
	stop              chan struct{}
	blacklist         []Key

	// BgColor and FgColor allow to configure the background and foreground
	// colors of the GUI.
	BgColor, FgColor, FrameColor Attribute

	// SelBgColor and SelFgColor allow to configure the background and
	// foreground colors of the frame of the current view.
	SelBgColor, SelFgColor, SelFrameColor Attribute

	// If Highlight is true, Sel{Bg,Fg}Colors will be used to draw the
	// frame of the current view.
	Highlight bool

	// If ShowListFooter is true then show list footer (i.e. the part that says we're at item 5 out of 10)
	ShowListFooter bool

	// If Cursor is true then the cursor is enabled.
	Cursor bool

	// If Mouse is true then mouse events will be enabled.
	Mouse bool

	// If InputEsc is true, when ESC sequence is in the buffer and it doesn't
	// match any known sequence, ESC means KeyEsc.
	InputEsc bool

	// SupportOverlaps is true when we allow for view edges to overlap with other
	// view edges
	SupportOverlaps bool

	Mutexes GuiMutexes

	OnSearchEscape func() error
	// these keys must either be of type Key of rune
	SearchEscapeKey    interface{}
	NextSearchMatchKey interface{}
	PrevSearchMatchKey interface{}

	screen         tcell.Screen
	suspendedMutex sync.Mutex
	suspended      bool
}

// NewGui returns a new Gui object with a given output mode.
func NewGui(mode OutputMode, supportOverlaps bool, playMode PlayMode, headless bool, runeReplacements map[rune]string) (*Gui, error) {
	g := &Gui{}

	var err error
	if headless {
		err = g.tcellInitSimulation()
	} else {
		err = g.tcellInit(runeReplacements)
	}
	if err != nil {
		return nil, err
	}

	if headless || runtime.GOOS == "windows" {
		g.maxX, g.maxY = g.screen.Size()
	} else {
		// TODO: find out if we actually need this bespoke logic for linux
		g.maxX, g.maxY, err = g.getTermWindowSize()
		if err != nil {
			return nil, err
		}
	}

	g.outputMode = mode

	g.stop = make(chan struct{})

	g.gEvents = make(chan GocuiEvent, 20)
	g.userEvents = make(chan userEvent, 20)

	if playMode == RECORDING {
		g.Recording = &Recording{
			KeyEvents:    []*TcellKeyEventWrapper{},
			ResizeEvents: []*TcellResizeEventWrapper{},
		}
	} else if playMode == REPLAYING || playMode == REPLAYING_NEW {
		g.ReplayedEvents = replayedEvents{
			Keys:    make(chan *TcellKeyEventWrapper),
			Resizes: make(chan *TcellResizeEventWrapper),
		}
	}

	g.BgColor, g.FgColor, g.FrameColor = ColorDefault, ColorDefault, ColorDefault
	g.SelBgColor, g.SelFgColor, g.SelFrameColor = ColorDefault, ColorDefault, ColorDefault

	// SupportOverlaps is true when we allow for view edges to overlap with other
	// view edges
	g.SupportOverlaps = supportOverlaps

	// default keys for when searching strings in a view
	g.SearchEscapeKey = KeyEsc
	g.NextSearchMatchKey = 'n'
	g.PrevSearchMatchKey = 'N'

	g.PlayMode = playMode

	return g, nil
}

// Close finalizes the library. It should be called after a successful
// initialization and when gocui is not needed anymore.
func (g *Gui) Close() {
	go func() {
		g.stop <- struct{}{}
	}()
	Screen.Fini()
}

// Size returns the terminal's size.
func (g *Gui) Size() (x, y int) {
	return g.maxX, g.maxY
}

// SetRune writes a rune at the given point, relative to the top-left
// corner of the terminal. It checks if the position is valid and applies
// the given colors.
func (g *Gui) SetRune(x, y int, ch rune, fgColor, bgColor Attribute) error {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		// swallowing error because it's not that big of a deal
		return nil
	}
	tcellSetCell(x, y, ch, fgColor, bgColor, g.outputMode)
	return nil
}

// Rune returns the rune contained in the cell at the given position.
// It checks if the position is valid.
func (g *Gui) Rune(x, y int) (rune, error) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return ' ', errors.New("invalid point")
	}
	c, _, _, _ := Screen.GetContent(x, y)
	return c, nil
}

// SetView creates a new view with its top-left corner at (x0, y0)
// and the bottom-right one at (x1, y1). If a view with the same name
// already exists, its dimensions are updated; otherwise, the error
// ErrUnknownView is returned, which allows to assert if the View must
// be initialized. It checks if the position is valid.
func (g *Gui) SetView(name string, x0, y0, x1, y1 int, overlaps byte) (*View, error) {
	if name == "" {
		return nil, errors.New("invalid name")
	}

	if v, err := g.View(name); err == nil {
		if v.x0 != x0 || v.x1 != x1 || v.y0 != y0 || v.y1 != y1 {
			v.clearViewLines()
		}

		v.x0 = x0
		v.y0 = y0
		v.x1 = x1
		v.y1 = y1
		return v, nil
	}

	g.Mutexes.ViewsMutex.Lock()

	v := newView(name, x0, y0, x1, y1, g.outputMode)
	v.BgColor, v.FgColor = g.BgColor, g.FgColor
	v.SelBgColor, v.SelFgColor = g.SelBgColor, g.SelFgColor
	v.Overlaps = overlaps
	g.views = append(g.views, v)

	g.Mutexes.ViewsMutex.Unlock()

	return v, errors.Wrap(ErrUnknownView, 0)
}

// SetViewBeneath sets a view stacked beneath another view
func (g *Gui) SetViewBeneath(name string, aboveViewName string, height int) (*View, error) {
	aboveView, err := g.View(aboveViewName)
	if err != nil {
		return nil, err
	}

	viewTop := aboveView.y1 + 1
	return g.SetView(name, aboveView.x0, viewTop, aboveView.x1, viewTop+height-1, 0)
}

// SetViewOnTop sets the given view on top of the existing ones.
func (g *Gui) SetViewOnTop(name string) (*View, error) {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	for i, v := range g.views {
		if v.name == name {
			s := append(g.views[:i], g.views[i+1:]...)
			g.views = append(s, v)
			return v, nil
		}
	}
	return nil, errors.Wrap(ErrUnknownView, 0)
}

// SetViewOnBottom sets the given view on bottom of the existing ones.
func (g *Gui) SetViewOnBottom(name string) (*View, error) {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	for i, v := range g.views {
		if v.name == name {
			s := append(g.views[:i], g.views[i+1:]...)
			g.views = append([]*View{v}, s...)
			return v, nil
		}
	}
	return nil, errors.Wrap(ErrUnknownView, 0)
}

func (g *Gui) SetViewOnTopOf(toMove string, other string) error {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	if toMove == other {
		return nil
	}

	// need to find the two current positions and then move toMove before other in the list.
	toMoveIndex := -1
	otherIndex := -1

	for i, v := range g.views {
		if v.name == toMove {
			toMoveIndex = i
		}

		if v.name == other {
			otherIndex = i
		}
	}

	if toMoveIndex == -1 || otherIndex == -1 {
		return errors.Wrap(ErrUnknownView, 0)
	}

	// already on top
	if toMoveIndex > otherIndex {
		return nil
	}

	// need to actually do it the other way around. Last is highest
	viewToMove := g.views[toMoveIndex]

	g.views = append(g.views[:toMoveIndex], g.views[toMoveIndex+1:]...)
	g.views = append(g.views[:otherIndex], append([]*View{viewToMove}, g.views[otherIndex:]...)...)
	return nil
}

// replaces the content in toView with the content in fromView
func (g *Gui) CopyContent(fromView *View, toView *View) {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	toView.clear()

	toView.lines = fromView.lines
	toView.viewLines = fromView.viewLines
	toView.ox = fromView.ox
	toView.oy = fromView.oy
	toView.cx = fromView.cx
	toView.cy = fromView.cy
}

// Views returns all the views in the GUI.
func (g *Gui) Views() []*View {
	return g.views
}

// View returns a pointer to the view with the given name, or error
// ErrUnknownView if a view with that name does not exist.
func (g *Gui) View(name string) (*View, error) {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	for _, v := range g.views {
		if v.name == name {
			return v, nil
		}
	}
	return nil, errors.Wrap(ErrUnknownView, 0)
}

// VisibleViewByPosition returns a pointer to a view matching the given position, or
// error ErrUnknownView if a view in that position does not exist.
func (g *Gui) VisibleViewByPosition(x, y int) (*View, error) {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	// traverse views in reverse order checking top views first
	for i := len(g.views); i > 0; i-- {
		v := g.views[i-1]

		if !v.Visible {
			continue
		}

		frameOffset := 0
		if v.Frame {
			frameOffset = 1
		}
		if x > v.x0-frameOffset && x < v.x1+frameOffset && y > v.y0-frameOffset && y < v.y1+frameOffset {
			return v, nil
		}
	}
	return nil, errors.Wrap(ErrUnknownView, 0)
}

// ViewPosition returns the coordinates of the view with the given name, or
// error ErrUnknownView if a view with that name does not exist.
func (g *Gui) ViewPosition(name string) (x0, y0, x1, y1 int, err error) {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	for _, v := range g.views {
		if v.name == name {
			return v.x0, v.y0, v.x1, v.y1, nil
		}
	}
	return 0, 0, 0, 0, errors.Wrap(ErrUnknownView, 0)
}

// DeleteView deletes a view by name.
func (g *Gui) DeleteView(name string) error {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	for i, v := range g.views {
		if v.name == name {
			g.views = append(g.views[:i], g.views[i+1:]...)
			return nil
		}
	}
	return errors.Wrap(ErrUnknownView, 0)
}

// SetCurrentView gives the focus to a given view.
func (g *Gui) SetCurrentView(name string) (*View, error) {
	g.Mutexes.ViewsMutex.Lock()
	defer g.Mutexes.ViewsMutex.Unlock()

	for _, v := range g.views {
		if v.name == name {
			g.currentView = v
			return v, nil
		}
	}
	return nil, errors.Wrap(ErrUnknownView, 0)
}

// CurrentView returns the currently focused view, or nil if no view
// owns the focus.
func (g *Gui) CurrentView() *View {
	return g.currentView
}

// SetKeybinding creates a new keybinding. If viewname equals to ""
// (empty string) then the keybinding will apply to all views. key must
// be a rune or a Key.
//
// When mouse keys are used (MouseLeft, MouseRight, ...), modifier might not work correctly.
// It behaves differently on different platforms. Somewhere it doesn't register Alt key press,
// on others it might report Ctrl as Alt. It's not consistent and therefore it's not recommended
// to use with mouse keys.
func (g *Gui) SetKeybinding(viewname string, key interface{}, mod Modifier, handler func(*Gui, *View) error) error {
	var kb *keybinding

	k, ch, err := getKey(key)
	if err != nil {
		return err
	}

	if g.isBlacklisted(k) {
		return ErrBlacklisted
	}

	kb = newKeybinding(viewname, k, ch, mod, handler)
	g.keybindings = append(g.keybindings, kb)
	return nil
}

// DeleteKeybinding deletes a keybinding.
func (g *Gui) DeleteKeybinding(viewname string, key interface{}, mod Modifier) error {
	k, ch, err := getKey(key)
	if err != nil {
		return err
	}

	for i, kb := range g.keybindings {
		if kb.viewName == viewname && kb.ch == ch && kb.key == k && kb.mod == mod {
			g.keybindings = append(g.keybindings[:i], g.keybindings[i+1:]...)
			return nil
		}
	}
	return errors.New("keybinding not found")
}

// DeleteKeybindings deletes all keybindings of view.
func (g *Gui) DeleteAllKeybindings() {
	g.keybindings = []*keybinding{}
	g.tabClickBindings = []*tabClickBinding{}
	g.viewMouseBindings = []*ViewMouseBinding{}
}

// DeleteKeybindings deletes all keybindings of view.
func (g *Gui) DeleteViewKeybindings(viewname string) {
	var s []*keybinding
	for _, kb := range g.keybindings {
		if kb.viewName != viewname {
			s = append(s, kb)
		}
	}
	g.keybindings = s
}

// SetTabClickBinding sets a binding for a tab click event
func (g *Gui) SetTabClickBinding(viewName string, handler tabClickHandler) error {
	g.tabClickBindings = append(g.tabClickBindings, &tabClickBinding{
		viewName: viewName,
		handler:  handler,
	})

	return nil
}

func (g *Gui) SetViewClickBinding(binding *ViewMouseBinding) error {
	g.viewMouseBindings = append(g.viewMouseBindings, binding)

	return nil
}

// BlackListKeybinding adds a keybinding to the blacklist
func (g *Gui) BlacklistKeybinding(k Key) error {
	for _, j := range g.blacklist {
		if j == k {
			return ErrAlreadyBlacklisted
		}
	}
	g.blacklist = append(g.blacklist, k)
	return nil
}

// WhiteListKeybinding removes a keybinding from the blacklist
func (g *Gui) WhitelistKeybinding(k Key) error {
	for i, j := range g.blacklist {
		if j == k {
			g.blacklist = append(g.blacklist[:i], g.blacklist[i+1:]...)
			return nil
		}
	}
	return ErrNotBlacklisted
}

// getKey takes an empty interface with a key and returns the corresponding
// typed Key or rune.
func getKey(key interface{}) (Key, rune, error) {
	switch t := key.(type) {
	case nil: // Ignore keybinding if `nil`
		return 0, 0, nil
	case Key:
		return t, 0, nil
	case rune:
		return 0, t, nil
	default:
		return 0, 0, errors.New("unknown type")
	}
}

// userEvent represents an event triggered by the user.
type userEvent struct {
	f func(*Gui) error
}

// Update executes the passed function. This method can be called safely from a
// goroutine in order to update the GUI. It is important to note that the
// passed function won't be executed immediately, instead it will be added to
// the user events queue. Given that Update spawns a goroutine, the order in
// which the user events will be handled is not guaranteed.
func (g *Gui) Update(f func(*Gui) error) {
	go g.UpdateAsync(f)
}

// UpdateAsync is a version of Update that does not spawn a go routine, it can
// be a bit more efficient in cases where Update is called many times like when
// tailing a file.  In general you should use Update()
func (g *Gui) UpdateAsync(f func(*Gui) error) {
	g.userEvents <- userEvent{f: f}
}

// A Manager is in charge of GUI's layout and can be used to build widgets.
type Manager interface {
	// Layout is called every time the GUI is redrawn, it must contain the
	// base views and its initializations.
	Layout(*Gui) error
}

// The ManagerFunc type is an adapter to allow the use of ordinary functions as
// Managers. If f is a function with the appropriate signature, ManagerFunc(f)
// is an Manager object that calls f.
type ManagerFunc func(*Gui) error

// Layout calls f(g)
func (f ManagerFunc) Layout(g *Gui) error {
	return f(g)
}

// SetManager sets the given GUI managers. It deletes all views and
// keybindings.
func (g *Gui) SetManager(managers ...Manager) {
	g.managers = managers
	g.currentView = nil
	g.views = nil
	g.keybindings = nil
	g.tabClickBindings = nil

	go func() { g.gEvents <- GocuiEvent{Type: eventResize} }()
}

// SetManagerFunc sets the given manager function. It deletes all views and
// keybindings.
func (g *Gui) SetManagerFunc(manager func(*Gui) error) {
	g.SetManager(ManagerFunc(manager))
}

// MainLoop runs the main loop until an error is returned. A successful
// finish should return ErrQuit.
func (g *Gui) MainLoop() error {
	g.StartTime = time.Now()
	if g.PlayMode == REPLAYING {
		go g.replayRecording()
	}

	go func() {
		for {
			select {
			case <-g.stop:
				return
			default:
				g.gEvents <- g.pollEvent()
			}
		}
	}()

	if g.Mouse {
		Screen.EnableMouse()
	}

	for {
		select {
		case ev := <-g.gEvents:
			if err := g.handleEvent(&ev); err != nil {
				return err
			}
		case ev := <-g.userEvents:
			if err := ev.f(g); err != nil {
				return err
			}
		}
		if err := g.consumeevents(); err != nil {
			return err
		}
		if err := g.flush(); err != nil {
			return err
		}
	}
}

// consumeevents handles the remaining events in the events pool.
func (g *Gui) consumeevents() error {
	for {
		select {
		case ev := <-g.gEvents:
			if err := g.handleEvent(&ev); err != nil {
				return err
			}
		case ev := <-g.userEvents:
			if err := ev.f(g); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// handleEvent handles an event, based on its type (key-press, error,
// etc.)
func (g *Gui) handleEvent(ev *GocuiEvent) error {
	switch ev.Type {
	case eventKey, eventMouse:
		return g.onKey(ev)
	case eventError:
		return ev.Err
	case eventResize:
		g.onResize()
		return nil
	default:
		return nil
	}
}

func (g *Gui) onResize() {
	// not sure if we actually need this
	// g.screen.Sync()
}

func (g *Gui) clear(fg, bg Attribute) (int, int) {
	st := getTcellStyle(oldStyle{fg: fg, bg: bg, outputMode: g.outputMode})
	w, h := Screen.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			Screen.SetContent(col, row, ' ', nil, st)
		}
	}
	return w, h
}

// drawFrameEdges draws the horizontal and vertical edges of a view.
func (g *Gui) drawFrameEdges(v *View, fgColor, bgColor Attribute) error {
	runeH, runeV := '─', '│'
	if len(v.FrameRunes) >= 2 {
		runeH, runeV = v.FrameRunes[0], v.FrameRunes[1]
	}

	for x := v.x0 + 1; x < v.x1 && x < g.maxX; x++ {
		if x < 0 {
			continue
		}
		if v.y0 > -1 && v.y0 < g.maxY {
			if err := g.SetRune(x, v.y0, runeH, fgColor, bgColor); err != nil {
				return err
			}
		}
		if v.y1 > -1 && v.y1 < g.maxY {
			if err := g.SetRune(x, v.y1, runeH, fgColor, bgColor); err != nil {
				return err
			}
		}
	}

	showScrollbar, realScrollbarStart, realScrollbarEnd := calcRealScrollbarStartEnd(v)
	for y := v.y0 + 1; y < v.y1 && y < g.maxY; y++ {
		if y < 0 {
			continue
		}
		if v.x0 > -1 && v.x0 < g.maxX {
			if err := g.SetRune(v.x0, y, runeV, fgColor, bgColor); err != nil {
				return err
			}
		}
		if v.x1 > -1 && v.x1 < g.maxX {
			runeToPrint := calcScrollbarRune(showScrollbar, realScrollbarStart, realScrollbarEnd, v.y0+1, v.y1-1, y, runeV)

			if err := g.SetRune(v.x1, y, runeToPrint, fgColor, bgColor); err != nil {
				return err
			}
		}
	}
	return nil
}

func calcScrollbarRune(showScrollbar bool, scrollbarStart int, scrollbarEnd int, rangeStart int, rangeEnd int, position int, runeV rune) rune {
	if !showScrollbar {
		return runeV
	} else if position == rangeStart {
		return '▲'
	} else if position == rangeEnd {
		return '▼'
	} else if position > scrollbarStart && position < scrollbarEnd {
		return '█'
	} else if position > rangeStart && position < rangeEnd {
		// keeping this as a separate branch in case we later want to render something different here.
		return runeV
	} else {
		return runeV
	}
}

func calcRealScrollbarStartEnd(v *View) (bool, int, int) {
	height := v.InnerHeight() + 1
	fullHeight := v.ViewLinesHeight() - v.scrollMargin()

	if v.CanScrollPastBottom {
		fullHeight += height
	}

	if height < 2 || height >= fullHeight {
		return false, 0, 0
	}

	originY := v.OriginY()
	scrollbarStart, scrollbarHeight := calcScrollbar(fullHeight, height, originY, height-1)
	top := v.y0 + 1
	realScrollbarStart := top + scrollbarStart
	realScrollbarEnd := realScrollbarStart + scrollbarHeight

	return true, realScrollbarStart, realScrollbarEnd
}

func cornerRune(index byte) rune {
	return []rune{' ', '│', '│', '│', '─', '┘', '┐', '┤', '─', '└', '┌', '├', '├', '┴', '┬', '┼'}[index]
}

// cornerCustomRune returns rune from `v.FrameRunes` slice. If the length of slice is less than 11
// all the missing runes will be translated to the default `cornerRune()`
func cornerCustomRune(v *View, index byte) rune {
	// Translate `cornerRune()` index
	//  0    1    2    3    4    5    6    7    8    9    10   11   12   13   14   15
	// ' ', '│', '│', '│', '─', '┘', '┐', '┤', '─', '└', '┌', '├', '├', '┴', '┬', '┼'
	// into `FrameRunes` index
	//  0    1    2    3    4    5    6    7    8    9    10
	// '─', '│', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼'
	switch index {
	case 1, 2, 3:
		return v.FrameRunes[1]
	case 4, 8:
		return v.FrameRunes[0]
	case 5:
		return v.FrameRunes[5]
	case 6:
		return v.FrameRunes[3]
	case 7:
		if len(v.FrameRunes) < 8 {
			break
		}
		return v.FrameRunes[7]
	case 9:
		return v.FrameRunes[4]
	case 10:
		return v.FrameRunes[2]
	case 11, 12:
		if len(v.FrameRunes) < 7 {
			break
		}
		return v.FrameRunes[6]
	case 13:
		if len(v.FrameRunes) < 10 {
			break
		}
		return v.FrameRunes[9]
	case 14:
		if len(v.FrameRunes) < 9 {
			break
		}
		return v.FrameRunes[8]
	case 15:
		if len(v.FrameRunes) < 11 {
			break
		}
		return v.FrameRunes[10]
	default:
		return ' ' // cornerRune(0)
	}
	return cornerRune(index)
}

func corner(v *View, directions byte) rune {
	index := v.Overlaps | directions
	if len(v.FrameRunes) >= 6 {
		return cornerCustomRune(v, index)
	}
	return cornerRune(index)
}

// drawFrameCorners draws the corners of the view.
func (g *Gui) drawFrameCorners(v *View, fgColor, bgColor Attribute) error {
	if v.y0 == v.y1 {
		if !g.SupportOverlaps && v.x0 >= 0 && v.x1 >= 0 && v.y0 >= 0 && v.x0 < g.maxX && v.x1 < g.maxX && v.y0 < g.maxY {
			if err := g.SetRune(v.x0, v.y0, '╶', fgColor, bgColor); err != nil {
				return err
			}
			if err := g.SetRune(v.x1, v.y0, '╴', fgColor, bgColor); err != nil {
				return err
			}
		}
		return nil
	}

	runeTL, runeTR, runeBL, runeBR := '┌', '┐', '└', '┘'
	if len(v.FrameRunes) >= 6 {
		runeTL, runeTR, runeBL, runeBR = v.FrameRunes[2], v.FrameRunes[3], v.FrameRunes[4], v.FrameRunes[5]
	}
	if g.SupportOverlaps {
		runeTL = corner(v, BOTTOM|RIGHT)
		runeTR = corner(v, BOTTOM|LEFT)
		runeBL = corner(v, TOP|RIGHT)
		runeBR = corner(v, TOP|LEFT)
	}

	corners := []struct {
		x, y int
		ch   rune
	}{{v.x0, v.y0, runeTL}, {v.x1, v.y0, runeTR}, {v.x0, v.y1, runeBL}, {v.x1, v.y1, runeBR}}

	for _, c := range corners {
		if c.x >= 0 && c.y >= 0 && c.x < g.maxX && c.y < g.maxY {
			if err := g.SetRune(c.x, c.y, c.ch, fgColor, bgColor); err != nil {
				return err
			}
		}
	}
	return nil
}

// drawTitle draws the title of the view.
func (g *Gui) drawTitle(v *View, fgColor, bgColor Attribute) error {
	if v.y0 < 0 || v.y0 >= g.maxY {
		return nil
	}

	tabs := v.Tabs
	separator := " - "
	charIndex := 0
	currentTabStart := -1
	currentTabEnd := -1
	if len(tabs) == 0 {
		tabs = []string{v.Title}
	} else {
		for i, tab := range tabs {
			if i == v.TabIndex {
				currentTabStart = charIndex
				currentTabEnd = charIndex + len(tab)
				break
			}
			charIndex += len(tab)
			if i < len(tabs)-1 {
				charIndex += len(separator)
			}
		}
	}

	str := strings.Join(tabs, separator)

	x := v.x0 + 2
	for i, ch := range str {
		if x < 0 {
			continue
		} else if x > v.x1-2 || x >= g.maxX {
			break
		}
		currentFgColor := fgColor
		currentBgColor := bgColor
		// if you are the current view and you have multiple tabs, de-highlight the non-selected tabs
		if v == g.currentView && len(v.Tabs) > 0 {
			currentFgColor = v.FgColor
			currentBgColor = v.BgColor
		}

		if i >= currentTabStart && i <= currentTabEnd {
			currentFgColor = v.SelFgColor
			if v != g.currentView {
				currentFgColor -= AttrBold
			}
		}
		if err := g.SetRune(x, v.y0, ch, currentFgColor, currentBgColor); err != nil {
			return err
		}
		x += runewidth.RuneWidth(ch)
	}
	return nil
}

// drawSubtitle draws the subtitle of the view.
func (g *Gui) drawSubtitle(v *View, fgColor, bgColor Attribute) error {
	if v.y0 < 0 || v.y0 >= g.maxY {
		return nil
	}

	start := v.x1 - 5 - len(v.Subtitle)
	if start < v.x0 {
		return nil
	}
	x := start
	for _, ch := range v.Subtitle {
		if x >= v.x1 {
			break
		}
		if err := g.SetRune(x, v.y0, ch, fgColor, bgColor); err != nil {
			return err
		}
		x += runewidth.RuneWidth(ch)
	}
	return nil
}

// drawListFooter draws the footer of a list view, showing something like '1 of 10'
func (g *Gui) drawListFooter(v *View, fgColor, bgColor Attribute) error {
	if len(v.lines) == 0 {
		return nil
	}

	message := v.Footer

	if v.y1 < 0 || v.y1 >= g.maxY {
		return nil
	}

	start := v.x1 - 1 - len(message)
	if start < v.x0 {
		return nil
	}
	x := start
	for _, ch := range message {
		if x >= v.x1 {
			break
		}
		if err := g.SetRune(x, v.y1, ch, fgColor, bgColor); err != nil {
			return err
		}
		x += runewidth.RuneWidth(ch)
	}
	return nil
}

// flush updates the gui, re-drawing frames and buffers.
func (g *Gui) flush() error {
	// pretty sure we don't need this, but keeping it here in case we get weird visual artifacts
	// g.clear(g.FgColor, g.BgColor)

	maxX, maxY := Screen.Size()
	// if GUI's size has changed, we need to redraw all views
	if maxX != g.maxX || maxY != g.maxY {
		for _, v := range g.views {
			v.clearViewLines()
		}
	}
	g.maxX, g.maxY = maxX, maxY

	for _, m := range g.managers {
		if err := m.Layout(g); err != nil {
			return err
		}
	}
	for _, v := range g.views {
		if err := g.draw(v); err != nil {
			return err
		}
	}

	Screen.Show()
	return nil
}

// draw manages the cursor and calls the draw function of a view.
func (g *Gui) draw(v *View) error {
	if g.suspended {
		return nil
	}

	if !v.Visible || v.y1 < v.y0 || v.x1 < v.x0 {
		return nil
	}

	if g.Cursor {
		if curview := g.currentView; curview != nil {
			vMaxX, vMaxY := curview.Size()
			if curview.cx < 0 {
				curview.cx = 0
			} else if curview.cx >= vMaxX {
				curview.cx = vMaxX - 1
			}
			if curview.cy < 0 {
				curview.cy = 0
			} else if curview.cy >= vMaxY {
				curview.cy = vMaxY - 1
			}

			gMaxX, gMaxY := g.Size()
			cx, cy := curview.x0+curview.cx+1, curview.y0+curview.cy+1
			// This test probably doesn't need to be here.
			// tcell is hiding cursor by setting coordinates outside of screen.
			// Keeping it here for now, as I'm not 100% sure :)
			if cx >= 0 && cx < gMaxX && cy >= 0 && cy < gMaxY {
				Screen.ShowCursor(cx, cy)
			} else {
				Screen.HideCursor()
			}
		}
	} else {
		Screen.HideCursor()
	}

	if err := v.draw(); err != nil {
		return err
	}

	if v.Frame {
		var fgColor, bgColor, frameColor Attribute
		if g.Highlight && v == g.currentView {
			fgColor = g.SelFgColor
			bgColor = g.SelBgColor
			frameColor = g.SelFrameColor
		} else {
			bgColor = g.BgColor
			if v.TitleColor != ColorDefault {
				fgColor = v.TitleColor
			} else {
				fgColor = g.FgColor
			}
			if v.FrameColor != ColorDefault {
				frameColor = v.FrameColor
			} else {
				frameColor = g.FrameColor
			}
		}

		if err := g.drawFrameEdges(v, frameColor, bgColor); err != nil {
			return err
		}
		if err := g.drawFrameCorners(v, frameColor, bgColor); err != nil {
			return err
		}
		if v.Title != "" || len(v.Tabs) > 0 {
			if err := g.drawTitle(v, fgColor, bgColor); err != nil {
				return err
			}
		}
		if v.Subtitle != "" {
			if err := g.drawSubtitle(v, fgColor, bgColor); err != nil {
				return err
			}
		}
		if v.Footer != "" && g.ShowListFooter {
			if err := g.drawListFooter(v, fgColor, bgColor); err != nil {
				return err
			}
		}
	}

	return nil
}

// onKey manages key-press events. A keybinding handler is called when
// a key-press or mouse event satisfies a configured keybinding. Furthermore,
// currentView's internal buffer is modified if currentView.Editable is true.
func (g *Gui) onKey(ev *GocuiEvent) error {
	switch ev.Type {
	case eventKey:

		_, err := g.execKeybindings(g.currentView, ev)
		if err != nil {
			return err
		}

	case eventMouse:
		mx, my := ev.MouseX, ev.MouseY
		v, err := g.VisibleViewByPosition(mx, my)
		if err != nil {
			break
		}
		if v.Frame && my == v.y0 {
			if len(v.Tabs) > 0 {
				tabIndex := v.GetClickedTabIndex(mx - v.x0)

				if tabIndex >= 0 {
					for _, binding := range g.tabClickBindings {
						if binding.viewName == v.Name() {
							return binding.handler(tabIndex)
						}
					}
				}
			}
		}

		newCx := mx - v.x0 - 1
		newCy := my - v.y0 - 1
		// if view  is editable don't go further than the furthest character for that line
		if v.Editable && newCy >= 0 && newCy <= len(v.lines)-1 {
			lastCharForLine := len(v.lines[newCy])
			if lastCharForLine < newCx {
				newCx = lastCharForLine
			}
		}
		if err := v.SetCursor(newCx, newCy); err != nil {
			return err
		}

		if IsMouseKey(ev.Key) {
			opts := ViewMouseBindingOpts{X: newCx + v.ox, Y: newCy + v.oy}
			matched, err := g.execMouseKeybindings(v, ev, opts)
			if err != nil {
				return err
			}
			if matched {
				return nil
			}
		}

		if _, err := g.execKeybindings(v, ev); err != nil {
			return err
		}
	}

	return nil
}

func (g *Gui) execMouseKeybindings(view *View, ev *GocuiEvent, opts ViewMouseBindingOpts) (bool, error) {
	isMatch := func(binding *ViewMouseBinding) bool {
		return binding.ViewName == view.Name() &&
			ev.Key == binding.Key &&
			ev.Mod == binding.Modifier
	}

	// first pass looks for ones that match the focused view
	for _, binding := range g.viewMouseBindings {
		if isMatch(binding) && binding.FocusedView != "" && binding.FocusedView == g.currentView.Name() {
			return true, binding.Handler(opts)
		}
	}

	for _, binding := range g.viewMouseBindings {
		if isMatch(binding) && binding.FocusedView == "" {
			return true, binding.Handler(opts)
		}
	}

	return false, nil
}

func IsMouseKey(key interface{}) bool {
	switch key {
	case
		MouseLeft,
		MouseRight,
		MouseMiddle,
		MouseRelease,
		MouseWheelUp,
		MouseWheelDown,
		MouseWheelLeft,
		MouseWheelRight:
		return true
	default:
		return false
	}
}

// execKeybindings executes the keybinding handlers that match the passed view
// and event. The value of matched is true if there is a match and no errors.
func (g *Gui) execKeybindings(v *View, ev *GocuiEvent) (matched bool, err error) {
	var globalKb *keybinding
	var matchingParentViewKb *keybinding

	// if we're searching, and we've hit n/N/Esc, we ignore the default keybinding
	if v != nil && v.IsSearching() && Modifier(ev.Mod) == ModNone {
		if eventMatchesKey(ev, g.NextSearchMatchKey) {
			return true, v.gotoNextMatch()
		} else if eventMatchesKey(ev, g.PrevSearchMatchKey) {
			return true, v.gotoPreviousMatch()
		} else if eventMatchesKey(ev, g.SearchEscapeKey) {
			v.searcher.clearSearch()
			if g.OnSearchEscape != nil {
				if err := g.OnSearchEscape(); err != nil {
					return true, err
				}
			}
			return true, nil
		}
	}

	for _, kb := range g.keybindings {
		if kb.handler == nil {
			continue
		}
		if !kb.matchKeypress(Key(ev.Key), ev.Ch, Modifier(ev.Mod)) {
			continue
		}
		if g.matchView(v, kb) {
			return g.execKeybinding(v, kb)
		}
		if v != nil && g.matchView(v.ParentView, kb) {
			matchingParentViewKb = kb
		}
		if globalKb == nil && kb.viewName == "" && ((v != nil && !v.Editable) || (kb.ch == 0 && kb.key != KeyCtrlU && kb.key != KeyCtrlA && kb.key != KeyCtrlE)) {
			globalKb = kb
		}
	}
	if matchingParentViewKb != nil {
		return g.execKeybinding(v.ParentView, matchingParentViewKb)
	}

	if g.currentView != nil && g.currentView.Editable && g.currentView.Editor != nil {
		matched := g.currentView.Editor.Edit(g.currentView, Key(ev.Key), ev.Ch, Modifier(ev.Mod))
		if matched {
			return true, nil
		}
	}

	if globalKb != nil {
		return g.execKeybinding(v, globalKb)
	}
	return false, nil
}

// execKeybinding executes a given keybinding
func (g *Gui) execKeybinding(v *View, kb *keybinding) (bool, error) {
	if g.isBlacklisted(kb.key) {
		return true, nil
	}

	if err := kb.handler(g, v); err != nil {
		return false, err
	}
	return true, nil
}

func (g *Gui) StartTicking() {
	go func() {
		g.Mutexes.tickingMutex.Lock()
		defer g.Mutexes.tickingMutex.Unlock()
		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()
	outer:
		for {
			select {
			case <-ticker.C:
				// I'm okay with having a data race here: there's no harm in letting one of these updates through
				if g.suspended {
					continue outer
				}

				for _, view := range g.Views() {
					if view.HasLoader {
						g.userEvents <- userEvent{func(g *Gui) error { return nil }}
						continue outer
					}
				}
				return
			case <-g.stop:
				return
			}
		}
	}()
}

// isBlacklisted reports whether the key is blacklisted
func (g *Gui) isBlacklisted(k Key) bool {
	for _, j := range g.blacklist {
		if j == k {
			return true
		}
	}
	return false
}

// IsUnknownView reports whether the contents of an error is "unknown view".
func IsUnknownView(err error) bool {
	return err != nil && err.Error() == ErrUnknownView.Error()
}

// IsQuit reports whether the contents of an error is "quit".
func IsQuit(err error) bool {
	return err != nil && err.Error() == ErrQuit.Error()
}

func (g *Gui) replayRecording() {
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(2)

	// lots of duplication here due to lack of generics. Also we don't support mouse
	// events because it would be awkward to replicate but it would be trivial to add
	// support
	go func() {
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()

		// The playback could be paused at any time because integration tests run concurrently.
		// Therefore we can't just check for a given event whether we've passed its timestamp,
		// or else we'll have an explosion of keypresses after the test is resumed.
		// We need to check if we've waited long enough since the last event was replayed.
		for i, event := range g.Recording.KeyEvents {
			var prevEventTimestamp int64 = 0
			if i > 0 {
				prevEventTimestamp = g.Recording.KeyEvents[i-1].Timestamp
			}
			timeToWait := float64(event.Timestamp-prevEventTimestamp) / g.RecordingConfig.Speed
			if i == 0 {
				timeToWait += float64(g.RecordingConfig.Leeway)
			}
			var timeWaited float64 = 0
		middle:
			for {
				select {
				case <-ticker.C:
					timeWaited += 1
					if timeWaited >= timeToWait {
						g.ReplayedEvents.Keys <- event
						break middle
					}
				case <-g.stop:
					return
				}
			}
		}

		waitGroup.Done()
	}()

	go func() {
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()

		// duplicating until Go gets generics
		for i, event := range g.Recording.ResizeEvents {
			var prevEventTimestamp int64 = 0
			if i > 0 {
				prevEventTimestamp = g.Recording.ResizeEvents[i-1].Timestamp
			}
			timeToWait := float64(event.Timestamp-prevEventTimestamp) / g.RecordingConfig.Speed
			if i == 0 {
				timeToWait += float64(g.RecordingConfig.Leeway)
			}
			var timeWaited float64 = 0
		middle2:
			for {
				select {
				case <-ticker.C:
					timeWaited += 1
					if timeWaited >= timeToWait {
						g.ReplayedEvents.Resizes <- event
						break middle2
					}
				case <-g.stop:
					return
				}
			}
		}

		waitGroup.Done()
	}()

	waitGroup.Wait()

	// leaving some time for any handlers to execute before quitting
	time.Sleep(time.Second * 1)

	g.Update(func(*Gui) error {
		return ErrQuit
	})

	time.Sleep(time.Second * 1)

	log.Fatal("gocui should have already exited")
}

func (g *Gui) Suspend() error {
	g.suspendedMutex.Lock()
	defer g.suspendedMutex.Unlock()

	if g.suspended {
		return errors.New("Already suspended")
	}

	g.suspended = true

	return g.screen.Suspend()
}

func (g *Gui) Resume() error {
	g.suspendedMutex.Lock()
	defer g.suspendedMutex.Unlock()

	if !g.suspended {
		return errors.New("Cannot resume because we are not suspended")
	}

	g.suspended = false

	return g.screen.Resume()
}

// matchView returns if the keybinding matches the current view (and the view's context)
func (g *Gui) matchView(v *View, kb *keybinding) bool {
	// if the user is typing in a field, ignore char keys
	if v == nil {
		return false
	}
	if v.Editable == true && kb.ch != 0 {
		return false
	}
	if kb.viewName != v.name {
		return false
	}
	return true
}
