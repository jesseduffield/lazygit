// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "errors"

  "github.com/fatih/color"
  "github.com/jroimartin/gocui"
)

var (
  // ErrNoCommits : When no commits are found for the branch
  ErrNoCommits = errors.New("No commits for this branch")
)

func refreshCommits(g *gocui.Gui) error {
  state.Commits = getCommits()
  g.Update(func(*gocui.Gui) error {
    v, err := g.View("commits")
    if err != nil {
      panic(err)
    }
    v.Clear()
    red := color.New(color.FgRed)
    yellow := color.New(color.FgYellow)
    white := color.New(color.FgWhite)
    shaColor := white
    for _, commit := range state.Commits {
      if commit.Pushed {
        shaColor = red
      } else {
        shaColor = yellow
      }
      shaColor.Fprint(v, commit.Sha+" ")
      white.Fprintln(v, commit.Name)
    }
    return nil
  })
  return nil
}

func handleCommitSelect(g *gocui.Gui, v *gocui.View) error {
  renderString(g, "options", "s: squash down, r: rename")
  commit, err := getSelectedCommit(v)
  if err != nil {
    if err != ErrNoCommits {
      return err
    }
    return renderString(g, "main", "No commits for this branch")
  }
  commitText := gitShow(commit.Sha)
  return renderString(g, "main", commitText)
}

func handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {
  if getItemPosition(v) != 0 {
    return createSimpleConfirmationPanel(g, v, "Error", "Can only squash topmost commit")
  }
  commit, err := getSelectedCommit(v)
  if err != nil {
    return err
  }
  if output, err := gitSquashPreviousTwoCommits(commit.Name); err != nil {
    return createSimpleConfirmationPanel(g, v, "Error", output)
  }
  if err := refreshCommits(g); err != nil {
    panic(err)
  }
  return handleCommitSelect(g, v)
}

func handleRenameCommit(g *gocui.Gui, v *gocui.View) error {
  if getItemPosition(v) != 0 {
    return createSimpleConfirmationPanel(g, v, "Error", "Can only rename topmost commit")
  }
  createPromptPanel(g, v, "Rename Commit", func(g *gocui.Gui, v *gocui.View) error {
    if output, err := gitRenameCommit(v.Buffer()); err != nil {
      return createSimpleConfirmationPanel(g, v, "Error", output)
    }
    if err := refreshCommits(g); err != nil {
      panic(err)
    }
    return handleCommitSelect(g, v)
  })
  return nil
}

func getSelectedCommit(v *gocui.View) (Commit, error) {
  if len(state.Commits) == 0 {
    return Commit{}, ErrNoCommits
  }
  lineNumber := getItemPosition(v)
  return state.Commits[lineNumber], nil
}
