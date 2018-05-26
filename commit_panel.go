// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "fmt"

  "github.com/jroimartin/gocui"
)

func handleCommitPress(g *gocui.Gui, currentView *gocui.View) error {
  devLog(stagedFiles(state.GitFiles))
  if len(stagedFiles(state.GitFiles)) == 0 {
    return nil
  }
  maxX, maxY := g.Size()
  // var v *gocui.View
  if v, err := g.SetView("commit", maxX/2-30, maxY/2-1, maxX/2+30, maxY/2+1); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Commit Message"
    v.Editable = true
    if _, err := g.SetCurrentView("commit"); err != nil {
      return err
    }
    switchFocus(g, currentView, v)
  }
  return nil
}

func handleCommitSubmit(g *gocui.Gui, v *gocui.View) error {
  if len(v.BufferLines()) == 0 {
    return closeCommitPrompt(g, v)
  }
  message := fmt.Sprint(v.BufferLines()[0])
  // for whatever reason, a successful commit returns an error, so we're not
  // going to check for an error here
  if err := gitCommit(message); err != nil {
    devLog(err)
    panic(err)
  }
  refreshFiles(g)
  refreshLogs(g)
  return closeCommitPrompt(g, v)
}

func closeCommitPrompt(g *gocui.Gui, v *gocui.View) error {
  filesView, _ := g.View("files")
  switchFocus(g, v, filesView)
  devLog("test prompt close")
  if _, err := g.SetViewOnBottom("commit"); err != nil {
    return err
  }
  if _, err := g.SetCurrentView(state.PreviousView); err != nil {
    return err
  }
  return nil
}

func handleCommitPromptFocus(g *gocui.Gui, v *gocui.View) error {
  return renderString(g, "options", "esc: close, enter: commit")
}
