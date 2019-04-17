// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// WARNING: tricky code just for testing purposes, do not use as reference.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var vbuf, buf string

func quit(g *gocui.Gui, v *gocui.View) error {
	vbuf = v.ViewBuffer()
	buf = v.Buffer()
	return gocui.ErrQuit
}

func overwrite(g *gocui.Gui, v *gocui.View) error {
	v.Overwrite = !v.Overwrite
	return nil
}

func layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("main", 0, 0, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("main", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlI, gocui.ModNone, overwrite); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	g.Close()

	fmt.Printf("VBUF:\n%s\n", vbuf)
	fmt.Printf("BUF:\n%s\n", buf)
}
