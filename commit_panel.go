package main

import (
  "fmt"

  "github.com/jroimartin/gocui"
)

func handleCommitPress(g *gocui.Gui, currentView *gocui.View) error {
  if len(stagedFiles(state.GitFiles)) == 0 {
    return createSimpleConfirmationPanel(g, currentView, "Nothing to Commit", "There are no staged files to commit (esc)")
  }
  maxX, maxY := g.Size()
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
    panic(err)
  }
  refreshFiles(g)
  refreshCommits(g)
  return closeCommitPrompt(g, v)
}

func closeCommitPrompt(g *gocui.Gui, v *gocui.View) error {
  filesView, _ := g.View("files")
  // not passing in the view as oldView to switchFocus because we don't want a
  // reference pointing to a deleted view
  switchFocus(g, nil, filesView)
  if err := g.DeleteView("commit"); err != nil {
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
