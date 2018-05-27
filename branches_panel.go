// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (

  // "io"
  // "io/ioutil"

  // "strings"

  "github.com/fatih/color"
  "github.com/jroimartin/gocui"
)

func handleBranchPress(g *gocui.Gui, v *gocui.View) error {
  branch := getSelectedBranch(v)
  if output, err := gitCheckout(branch.Name, false); err != nil {
    createSimpleConfirmationPanel(g, v, "Error", output)
  }
  return refreshSidePanels(g, v)
}

func handleForceCheckout(g *gocui.Gui, v *gocui.View) error {
  branch := getSelectedBranch(v)
  return createConfirmationPanel(g, v, "Force Checkout Branch", "Are you sure you want force checkout? You will lose all local changes (y/n)", func(g *gocui.Gui, v *gocui.View) error {
    if output, err := gitCheckout(branch.Name, true); err != nil {
      createSimpleConfirmationPanel(g, v, "Error", output)
    }
    return refreshSidePanels(g, v)
  }, nil)
}

func getSelectedBranch(v *gocui.View) Branch {
  lineNumber := getItemPosition(v)
  return state.Branches[lineNumber]
}

func handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
  renderString(g, "options", "space: checkout, s: squash down")
  lineNumber := getItemPosition(v)
  branch := state.Branches[lineNumber]
  diff, _ := getBranchDiff(branch.Name, branch.BaseBranch)
  if err := renderString(g, "main", diff); err != nil {
    return err
  }
  return nil
}

func refreshBranches(g *gocui.Gui) error {
  v, err := g.View("branches")
  if err != nil {
    panic(err)
  }
  state.Branches = getGitBranches()
  yellow := color.New(color.FgYellow)
  red := color.New(color.FgRed)
  white := color.New(color.FgWhite)
  green := color.New(color.FgGreen)

  v.Clear()
  for _, branch := range state.Branches {
    if branch.Type == "feature" {
      green.Fprintln(v, branch.DisplayString)
      continue
    }
    if branch.Type == "bugfix" {
      yellow.Fprintln(v, branch.DisplayString)
      continue
    }
    if branch.Type == "hotfix" {
      red.Fprintln(v, branch.DisplayString)
      continue
    }
    white.Fprintln(v, branch.DisplayString)
  }
  resetOrigin(v)
  return nil
}
