// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "fmt"

  "github.com/jroimartin/gocui"
)

func refreshLogs(g *gocui.Gui) error {
  // here is where you want to pickup from
  // state.Logs = getGitLogs(nil)
  s := getLog()
  g.Update(func(*gocui.Gui) error {
    v, err := g.View("logs")
    v.Clear()
    if err != nil {
      panic(err)
    }
    v.Clear()
    fmt.Fprint(v, s)
    return nil
  })
  return nil
}
