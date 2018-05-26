// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (

  // "io"
  // "io/ioutil"

  "math"
  // "strings"

  "github.com/jroimartin/gocui"
)

func wrappedConfirmationFunction(function func(*gocui.Gui, *gocui.View) error) func(*gocui.Gui, *gocui.View) error {
  return func(g *gocui.Gui, v *gocui.View) error {
    if function != nil {
      if err := function(g, v); err != nil {
        panic(err)
      }
    }
    if err := returnFocus(g, v); err != nil {
      panic(err)
    }
    g.DeleteKeybindings("confirmation")
    return g.DeleteView("confirmation")
  }
}

func getConfirmationPanelDimensions(g *gocui.Gui, prompt string) (int, int, int, int) {
  width, height := g.Size()
  panelWidth := 60
  panelHeight := int(math.Ceil(float64(len(prompt)) / float64(panelWidth)))
  return width/2 - panelWidth/2,
    height/2 - panelHeight/2 - panelHeight%2 - 1,
    width/2 + panelWidth/2,
    height/2 + panelHeight/2
}

func createConfirmationPanel(g *gocui.Gui, sourceView *gocui.View, title, prompt string, handleYes, handleNo func(*gocui.Gui, *gocui.View) error) error {
  x0, y0, x1, y1 := getConfirmationPanelDimensions(g, prompt)
  if v, err := g.SetView("confirmation", x0, y0, x1, y1); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = title
    renderString(g, "confirmation", prompt+" (y/n)")
    switchFocus(g, sourceView, v)
    if err := g.SetKeybinding("confirmation", 'n', gocui.ModNone, wrappedConfirmationFunction(handleNo)); err != nil {
      return err
    }
    if err := g.SetKeybinding("confirmation", 'y', gocui.ModNone, wrappedConfirmationFunction(handleYes)); err != nil {
      return err
    }
  }
  return nil
}
