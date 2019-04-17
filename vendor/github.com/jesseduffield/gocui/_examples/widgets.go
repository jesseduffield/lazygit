// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/go-errors/errors"
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

const delta = 0.2

type HelpWidget struct {
	name string
	x, y int
	w, h int
	body string
}

func NewHelpWidget(name string, x, y int, body string) *HelpWidget {
	lines := strings.Split(body, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	return &HelpWidget{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (w *HelpWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}

type StatusbarWidget struct {
	name string
	x, y int
	w    int
	val  float64
}

func NewStatusbarWidget(name string, x, y, w int) *StatusbarWidget {
	return &StatusbarWidget{name: name, x: x, y: y, w: w}
}

func (w *StatusbarWidget) SetVal(val float64) error {
	if val < 0 || val > 1 {
		return errors.New("invalid value")
	}
	w.val = val
	return nil
}

func (w *StatusbarWidget) Val() float64 {
	return w.val
}

func (w *StatusbarWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()

	rep := int(w.val * float64(w.w-1))
	fmt.Fprint(v, strings.Repeat("▒", rep))
	return nil
}

type ButtonWidget struct {
	name    string
	x, y    int
	w       int
	label   string
	handler func(g *gocui.Gui, v *gocui.View) error
}

func NewButtonWidget(name string, x, y int, label string, handler func(g *gocui.Gui, v *gocui.View) error) *ButtonWidget {
	return &ButtonWidget{name: name, x: x, y: y, w: len(label) + 1, label: label, handler: handler}
}

func (w *ButtonWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
			return err
		}
		fmt.Fprint(v, w.label)
	}
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorRed

	help := NewHelpWidget("help", 1, 1, helpText)
	status := NewStatusbarWidget("status", 1, 7, 50)
	butdown := NewButtonWidget("butdown", 52, 7, "DOWN", statusDown(status))
	butup := NewButtonWidget("butup", 58, 7, "UP", statusUp(status))
	g.SetManager(help, status, butdown, butup)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggleButton); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func toggleButton(g *gocui.Gui, v *gocui.View) error {
	nextview := "butdown"
	if v != nil && v.Name() == "butdown" {
		nextview = "butup"
	}
	_, err := g.SetCurrentView(nextview)
	return err
}

func statusUp(status *StatusbarWidget) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return statusSet(status, delta)
	}
}

func statusDown(status *StatusbarWidget) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return statusSet(status, -delta)
	}
}

func statusSet(sw *StatusbarWidget, inc float64) error {
	val := sw.Val() + inc
	if val < 0 || val > 1 {
		return nil
	}
	return sw.SetVal(val)
}

const helpText = `KEYBINDINGS
Tab: Move between buttons
Enter: Push button
^C: Exit`
