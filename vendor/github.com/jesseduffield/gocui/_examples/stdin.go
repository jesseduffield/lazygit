// Copyright 2015 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	if v, err := g.SetView("help", maxX-23, 0, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "↑ ↓: Seek input")
		fmt.Fprintln(v, "a: Enable autoscroll")
		fmt.Fprintln(v, "^C: Exit")
	}

	if v, err := g.SetView("stdin", 0, 0, 80, 35); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("stdin"); err != nil {
			return err
		}
		dumper := hex.Dumper(v)
		if _, err := io.Copy(dumper, os.Stdin); err != nil {
			return err
		}
		v.Wrap = true
	}

	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", 'a', gocui.ModNone, autoscroll); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func autoscroll(g *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = true
	return nil
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}
