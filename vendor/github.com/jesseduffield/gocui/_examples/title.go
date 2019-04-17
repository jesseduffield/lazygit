// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
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

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Overlap (front)
	if v, err := g.SetView("v1", 10, 2, 30, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v2", 20, 4, 40, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}

	// Overlap (back)
	if v, err := g.SetView("v3", 60, 4, 80, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v4", 50, 2, 70, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}

	// Overlap (frame)
	if v, err := g.SetView("v15", 90, 2, 110, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v16", 100, 5, 120, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v17", 140, 5, 160, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v18", 130, 2, 150, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}

	// Long title
	if v, err := g.SetView("v5", 10, 12, 30, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Long long long long title"
	}

	// No title
	if v, err := g.SetView("v6", 35, 12, 55, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = ""
	}
	if _, err := g.SetView("v7", 60, 12, 80, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	// Small view
	if v, err := g.SetView("v8", 85, 12, 88, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}

	// Screen borders
	if v, err := g.SetView("v9", -10, 20, 10, 24); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v10", maxX-10, 20, maxX+10, 24); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}

	// Out of screen
	if v, err := g.SetView("v11", -21, 28, -1, 32); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v12", maxX, 28, maxX+20, 32); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v13", 10, -7, 30, -1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}
	if v, err := g.SetView("v14", 10, maxY, 30, maxY+6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Regular title"
	}

	return nil
}
