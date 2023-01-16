# _Tcell_ Tutorial

_Tcell_ provides a low-level, portable API for building terminal-based programs.
A [terminal emulator](https://en.wikipedia.org/wiki/Terminal_emulator)
(or a real terminal such as a DEC VT-220) is used to interact with such a program.

_Tcell_'s interface is fairly low-level.
While it provides a reasonably portable way of dealing with all the usual terminal
features, it may be easier to utilize a higher level framework.
A number of such frameworks are listed on the _Tcell_ main [README](README.md).

This tutorial provides the details of _Tcell_, and is appropriate for developers
wishing to create their own application frameworks or needing more direct access
to the terminal capabilities.

## Resize events

Applications receive an event of type `EventResize` when they are first initialized and each time the terminal is resized.
The new size is available as `Size`.

```go
switch ev := ev.(type) {
case *tcell.EventResize:
	w, h := ev.Size()
	logMessage(fmt.Sprintf("Resized to %dx%d", w, h))
}
```

## Key events

When a key is pressed, applications receive an event of type `EventKey`.
This event describes the modifier keys pressed (if any) and the pressed key or rune.

When a rune key is pressed, an event with its `Key` set to `KeyRune` is dispatched.

When a non-rune key is pressed, it is available as the `Key` of the event.

```go
switch ev := ev.(type) {
case *tcell.EventKey:
    mod, key, ch := ev.Mod(), ev.Key(), ev.Rune()
    logMessage(fmt.Sprintf("EventKey Modifiers: %d Key: %d Rune: %d", mod, key, ch))
}
```

### Key event restrictions

Terminal-based programs have less visibility into keyboard activity than graphical applications.

When a key is pressed and held, additional key press events are sent by the terminal emulator.
The rate of these repeated events depends on the emulator's configuration.
Key release events are not available.

It is not possible to distinguish runes typed while holding shift and runes typed using caps lock.
Capital letters are reported without the Shift modifier.

## Mouse events

Applications receive an event of type `EventMouse` when the mouse moves, or a mouse button is pressed or released.
Mouse events are only delivered if
`EnableMouse` has been called.

The mouse buttons being pressed (if any) are available as `Buttons`, and the position of the mouse is available as `Position`.

```go
switch ev := ev.(type) {
case *tcell.EventMouse:
	mod := ev.Modifiers()
	btns := ev.Buttons()
	x, y := ev.Position()
	logMessage(fmt.Sprintf("EventMouse Modifiers: %d Buttons: %d Position: %d,%d", mod, btns, x, y))
}
```

### Mouse buttons

Identifier | Alias           | Description
-----------|-----------------|-----------
Button1    | ButtonPrimary   | Left button
Button2    | ButtonSecondary | Right button
Button3    | ButtonMiddle    | Middle button
Button4    |                 | Side button (thumb/next)
Button5    |                 | Side button (thumb/prev)
WheelUp    |                 | Scroll wheel up
WheelDown  |                 | Scroll wheel down
WheelLeft  |                 | Horizontal wheel left
WheelRight |                 | Horizontal wheel right

## Usage

To create a _Tcell_ application, first initialize a screen to hold it.

```go
s, err := tcell.NewScreen()
if err != nil {
	log.Fatalf("%+v", err)
}
if err := s.Init(); err != nil {
	log.Fatalf("%+v", err)
}

// Set default text style
defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
s.SetStyle(defStyle)

// Clear screen
s.Clear()
```

Text may be drawn on the screen using `SetContent`.

```go
s.SetContent(0, 0, 'H', nil, defStyle)
s.SetContent(1, 0, 'i', nil, defStyle)
s.SetContent(2, 0, '!', nil, defStyle)
```

To draw text more easily, define a render function.

```go
func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
```

Lastly, define an event loop to handle user input and update application state.

```go
quit := func() {
    s.Fini()
    os.Exit(0)
}
for {
    // Update screen
    s.Show()

    // Poll event
    ev := s.PollEvent()

    // Process event
    switch ev := ev.(type) {
    case *tcell.EventResize:
        s.Sync()
    case *tcell.EventKey:
        if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
            quit()
        }
    }
}
```

## Demo application

The following demonstrates how to initialize a screen, draw text/graphics and handle user input.

```go
package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	// Draw initial boxes
	drawBox(s, 1, 1, 42, 7, boxStyle, "Click and drag to draw a box")
	drawBox(s, 5, 9, 32, 14, boxStyle, "Press C to reset")

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.
	// xmax, ymax := s.Size()

	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

	// Event loop
	ox, oy := -1, -1
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
		case *tcell.EventMouse:
			x, y := ev.Position()

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:
				if ox < 0 {
					ox, oy = x, y // record location when click started
				}

			case tcell.ButtonNone:
				if ox >= 0 {
					label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
					drawBox(s, ox, oy, x, y, boxStyle, label)
					ox, oy = -1, -1
				}
			}
		}
	}
}
```

