package main

import (
  "fmt"
  "strings"
  "time"

  "github.com/fatih/color"
  "github.com/jesseduffield/gocui"
)

var cyclableViews = []string{"files", "branches", "commits", "stash"}

func refreshSidePanels(g *gocui.Gui, v *gocui.View) error {
  refreshBranches(g)
  refreshFiles(g)
  refreshCommits(g)
  return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
  var focusedViewName string
  if v == nil || v.Name() == cyclableViews[len(cyclableViews)-1] {
    focusedViewName = cyclableViews[0]
  } else {
    for i := range cyclableViews {
      if v.Name() == cyclableViews[i] {
        focusedViewName = cyclableViews[i+1]
        break
      }
      if i == len(cyclableViews)-1 {
        devLog(v.Name() + " is not in the list of views")
        return nil
      }
    }
  }
  focusedView, err := g.View(focusedViewName)
  if err != nil {
    panic(err)
    return err
  }
  return switchFocus(g, v, focusedView)
}

func newLineFocused(g *gocui.Gui, v *gocui.View) error {
  mainView, _ := g.View("main")
  mainView.SetOrigin(0, 0)

  switch v.Name() {
  case "files":
    return handleFileSelect(g, v)
  case "branches":
    return handleBranchSelect(g, v)
  case "confirmation":
    return nil
  case "main":
    return nil
  case "commits":
    return handleCommitSelect(g, v)
  case "stash":
    return handleStashEntrySelect(g, v)
  default:
    panic("No view matching newLineFocused switch statement")
  }
}

func returnFocus(g *gocui.Gui, v *gocui.View) error {
  previousView, err := g.View(state.PreviousView)
  if err != nil {
    panic(err)
  }
  return switchFocus(g, v, previousView)
}

func switchFocus(g *gocui.Gui, oldView, newView *gocui.View) error {
  // we assume we'll never want to return focus to a confirmation panel i.e.
  // we should never stack confirmation panels
  if oldView != nil && oldView.Name() != "confirmation" {
    oldView.Highlight = false
    devLog("setting previous view to:", oldView.Name())
    state.PreviousView = oldView.Name()
  }
  newView.Highlight = true
  devLog(newView.Name())
  if _, err := g.SetCurrentView(newView.Name()); err != nil {
    return err
  }
  g.Cursor = newView.Editable
  return newLineFocused(g, newView)
}

func getItemPosition(v *gocui.View) int {
  _, cy := v.Cursor()
  _, oy := v.Origin()
  return oy + cy
}

func trimmedContent(v *gocui.View) string {
  return strings.TrimSpace(v.Buffer())
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
  if v == nil {
    return nil
  }

  ox, oy := v.Origin()
  cx, cy := v.Cursor()
  if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
    if err := v.SetOrigin(ox, oy-1); err != nil {
      return err
    }
  }

  newLineFocused(g, v)
  return nil
}

func resetOrigin(v *gocui.View) error {
  if err := v.SetCursor(0, 0); err != nil {
    return err
  }
  return v.SetOrigin(0, 0)
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
  if v != nil {
    cx, cy := v.Cursor()
    ox, oy := v.Origin()
    if cy+oy >= len(v.BufferLines())-2 {
      return nil
    }
    if err := v.SetCursor(cx, cy+1); err != nil {
      if err := v.SetOrigin(ox, oy+1); err != nil {
        return err
      }
    }
  }

  newLineFocused(g, v)
  return nil
}

// if the cursor down past the last item, move it up one
func correctCursor(v *gocui.View) error {
  cx, cy := v.Cursor()
  _, oy := v.Origin()
  lineCount := len(v.BufferLines()) - 2
  if cy >= lineCount-oy {
    return v.SetCursor(cx, lineCount-oy)
  }
  return nil
}

func renderString(g *gocui.Gui, viewName, s string) error {
  g.Update(func(*gocui.Gui) error {
    timeStart := time.Now()
    colorLog(color.FgRed, viewName)
    v, err := g.View(viewName)
    if err != nil {
      panic(err)
    }
    v.Clear()
    fmt.Fprint(v, s)
    v.Wrap = true
    devLog("render time: ", time.Now().Sub(timeStart))
    return nil
  })
  return nil
}

func splitLines(multilineString string) []string {
  if multilineString == "" || multilineString == "\n" {
    return make([]string, 0)
  }
  lines := strings.Split(multilineString, "\n")
  if lines[len(lines)-1] == "" {
    return lines[:len(lines)-1]
  }
  return lines
}
