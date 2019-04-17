// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	if v, err := g.SetView("v1", 10, 2, 30, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "View #1")
	}
	if v, err := g.SetView("v2", 20, 4, 40, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "View #2")
	}
	if v, err := g.SetView("v3", 30, 6, 50, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "View #3")
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop("v1")
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '2', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop("v2")
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '3', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop("v3")
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
