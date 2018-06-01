package main

import (
  "fmt"
  "time"

  "github.com/fatih/color"
  "github.com/jroimartin/gocui"
)

func refreshStatus(g *gocui.Gui) error {
  v, err := g.View("status")
  if err != nil {
    return err
  }
  v.Clear()
  up, down := gitUpstreamDifferenceCount()
  fmt.Fprint(v, "↑"+up+"↓"+down)
  branches := state.Branches
  if len(branches) == 0 {
    return nil
  }
  branch := branches[0]
  // utilising the fact these all have padding to only grab the name
  // from the display string with the existing coloring applied
  fmt.Fprint(v, " "+branch.DisplayString[4:])

  colorLog(color.FgCyan, time.Now().Sub(StartTime))
  return nil
}
