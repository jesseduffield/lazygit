package main

import (
  "fmt"
  "time"

  "github.com/fatih/color"
  "github.com/jesseduffield/gocui"
)

func refreshStatus(g *gocui.Gui) error {
  v, err := g.View("status")
  if err != nil {
    panic(err)
  }
  // for some reason if this isn't wrapped in an update the clear seems to
  // be applied after the other things or something like that; the panel's
  // contents end up cleared
  g.Update(func(*gocui.Gui) error {
    v.Clear()
    pushables, pullables := gitUpstreamDifferenceCount()
    fmt.Fprint(v, "↑"+pushables+"↓"+pullables)
    branches := state.Branches
    if err := updateHasMergeConflictStatus(); err != nil {
      return err
    }
    if state.HasMergeConflicts {
      colour := color.New(color.FgYellow)
      fmt.Fprint(v, coloredString(" (merging)", colour))
    }
    if len(branches) == 0 {
      return nil
    }
    branch := branches[0]
    // utilising the fact these all have padding to only grab the name
    // from the display string with the existing coloring applied
    fmt.Fprint(v, " "+branch.DisplayString[4:])
    colorLog(color.FgCyan, time.Now().Sub(startTime))
    return nil
  })

  return nil
}
