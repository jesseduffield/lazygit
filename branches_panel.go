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
  if err := gitCheckout(branch.Name, false); err != nil {
    panic(err)
  }
  refreshBranches(v)
  refreshFiles(g)
  refreshLogs(g)
  return nil
}

func getSelectedBranch(v *gocui.View) Branch {
  lineNumber := getItemPosition(v)
  return state.Branches[lineNumber]
}

func handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
  renderString(g, "options", "space: checkout")
  lineNumber := getItemPosition(v)
  branch := state.Branches[lineNumber]
  diff, _ := getBranchDiff(branch.Name, branch.BaseBranch)
  if err := renderString(g, "main", diff); err != nil {
    return err
  }
  return nil
}

func refreshBranches(v *gocui.View) error {
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
