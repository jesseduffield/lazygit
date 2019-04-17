// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

type Label struct {
	name string
	w, h int
	body string
}

func NewLabel(name string, body string) *Label {
	lines := strings.Split(body, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	return &Label{name: name, w: w, h: h, body: body}
}

func (w *Label) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, 0, 0, w.w, w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}

func flowLayout(g *gocui.Gui) error {
	views := g.Views()
	x := 0
	for _, v := range views {
		w, h := v.Size()
		_, err := g.SetView(v.Name(), x, 0, x+w+1, h+1)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		x += w + 2
	}
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	l1 := NewLabel("l1", "This")
	l2 := NewLabel("l2", "is")
	l3 := NewLabel("l3", "a")
	l4 := NewLabel("l4", "flow\nlayout")
	l5 := NewLabel("l5", "!")
	fl := gocui.ManagerFunc(flowLayout)
	g.SetManager(l1, l2, l3, l4, l5, fl)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
