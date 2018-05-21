// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "fmt"
  "strings"
  // "io"
  // "io/ioutil"
  "log"
  // "strings"
  "os"

  "github.com/fatih/color"
  "github.com/jroimartin/gocui"
)

type stateType struct {
  GitFiles []GitFile
  Branches []Branch
}

var state = stateType{GitFiles: make([]GitFile, 0)}

var cyclableViews = []string{"files", "branches"}

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
        panic(v.Name() + " is not in the list of views")
      }
    }
  }
  focusedView, err := g.View(focusedViewName)
  if err != nil {
    panic(err)
    return err
  }
  if v != nil {
    v.Highlight = false
  }
  focusedView.Highlight = true
  devLog(focusedViewName)
  _, err = g.SetCurrentView(focusedViewName)
  itemSelected(g, focusedView)
  showViewOptions(g, focusedViewName)
  return err
}

func showViewOptions(g *gocui.Gui, viewName string) error {
  optionsMap := map[string]string{
    "files":    "space: toggle staged, c: commit changes",
    "branches": "space: checkout",
  }
  g.Update(func(*gocui.Gui) error {
    v, err := g.View("options")
    if err != nil {
      panic(err)
    }
    v.Clear()
    fmt.Fprint(v, optionsMap[viewName])
    return nil
  })
  return nil
}

func getItemPosition(v *gocui.View) int {
  _, cy := v.Cursor()
  _, oy := v.Origin()
  return oy + cy
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

  itemSelected(g, v)
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

  itemSelected(g, v)
  return nil
}

func itemSelected(g *gocui.Gui, v *gocui.View) error {
  mainView, _ := g.View("main")
  mainView.SetOrigin(0, 0)

  switch v.Name() {
  case "files":
    return handleFileSelect(g, v)
  case "branches":
    return handleBranchSelect(g, v)
  default:
    panic("No view matching itemSelected switch statement")
  }
}

func scrollUp(g *gocui.Gui, v *gocui.View) error {
  mainView, _ := g.View("main")
  ox, oy := mainView.Origin()
  if oy >= 1 {
    return mainView.SetOrigin(ox, oy-1)
  }
  return nil
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
  mainView, _ := g.View("main")
  ox, oy := mainView.Origin()
  if oy < len(mainView.BufferLines()) {
    return mainView.SetOrigin(ox, oy+1)
  }
  return nil
}

func devLog(s string) {
  f, _ := os.OpenFile("/Users/jesseduffieldduffield/go/src/github.com/jesseduffield/gitgot/development.log", os.O_APPEND|os.O_WRONLY, 0644)
  defer f.Close()

  f.WriteString(s + "\n")
}

func handleBranchPress(g *gocui.Gui, v *gocui.View) error {
  branch := getSelectedBranch(v)
  if err := gitCheckout(branch.Name, false); err != nil {
    return err
  }
  refreshBranches(v)
  refreshFiles(g)
  return nil
}

func handleFilePress(g *gocui.Gui, v *gocui.View) error {
  file := getSelectedFile(v)

  if file.HasUnstagedChanges {
    stageFile(file.Name)
  } else {
    unStageFile(file.Name)
  }

  if err := refreshFiles(g); err != nil {
    return err
  }
  if err := handleFileSelect(g, v); err != nil {
    return err
  }

  return nil
}

func getSelectedFile(v *gocui.View) GitFile {
  lineNumber := getItemPosition(v)
  return state.GitFiles[lineNumber]
}

func getSelectedBranch(v *gocui.View) Branch {
  lineNumber := getItemPosition(v)
  return state.Branches[lineNumber]
}

func handleBranchSelect(g *gocui.Gui, v *gocui.View) error {
  lineNumber := getItemPosition(v)
  branch := state.Branches[lineNumber]
  diff, _ := getBranchDiff(branch.Name, branch.BaseBranch)
  if err := renderString(g, diff); err != nil {
    return err
  }
  return nil
}

func handleFileSelect(g *gocui.Gui, v *gocui.View) error {
  item := getSelectedFile(v)
  diff := getDiff(item)
  if err := renderString(g, diff); err != nil {
    return err
  }

  // maxX, maxY := g.Size()
  // if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
  //   if err != gocui.ErrUnknownView {
  //     return errkjhgkhj
  //   }
  //   fmt.Fprintln(v, l)
  //   if _, err := g.SetCurrentView("msg"); err != nil {
  //     return err
  //   }
  // }
  return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
  if err := g.DeleteView("msg"); err != nil {
    return err
  }
  if _, err := g.SetCurrentView("files"); err != nil {
    return err
  }
  return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
  return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
  for _, view := range cyclableViews {
    if err := g.SetKeybinding(view, gocui.KeyTab, gocui.ModNone, nextView); err != nil {
      return err
    }
    if err := g.SetKeybinding(view, 'q', gocui.ModNone, quit); err != nil {
      return err
    }
    if err := g.SetKeybinding(view, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
      return err
    }
    if err := g.SetKeybinding(view, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
      return err
    }
    if err := g.SetKeybinding(view, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
      return err
    }
    if err := g.SetKeybinding(view, gocui.KeyPgup, gocui.ModNone, scrollUp); err != nil {
      return err
    }
    if err := g.SetKeybinding(view, gocui.KeyPgdn, gocui.ModNone, scrollDown); err != nil {
      return err
    }
  }
  if err := g.SetKeybinding("files", gocui.KeySpace, gocui.ModNone, handleFilePress); err != nil {
    return err
  }
  if err := g.SetKeybinding("branches", gocui.KeySpace, gocui.ModNone, handleBranchPress); err != nil {
    return err
  }
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

func refreshFiles(g *gocui.Gui) error {
  filesView, err := g.View("files")
  if err != nil {
    return err
  }

  // get files to stage
  gitFiles := getGitStatusFiles()
  state.GitFiles = mergeGitStatusFiles(state.GitFiles, gitFiles)

  filesView.Clear()
  red := color.New(color.FgRed)
  green := color.New(color.FgGreen)
  for _, gitFile := range state.GitFiles {
    if !gitFile.Tracked {
      red.Fprintln(filesView, gitFile.DisplayString)
      continue
    }
    green.Fprint(filesView, gitFile.DisplayString[0:1])
    red.Fprint(filesView, gitFile.DisplayString[1:3])
    if gitFile.HasUnstagedChanges {
      red.Fprintln(filesView, gitFile.Name)
    } else {
      green.Fprintln(filesView, gitFile.Name)
    }
  }
  return nil
}

func layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  leftSideWidth := maxX / 3
  filesBranchesBoundary := maxY - 10

  optionsTop := maxY - 3
  // hiding options if there's not enough space
  if maxY < 30 {
    optionsTop = maxY
  }

  sideView, err := g.SetView("files", 0, 0, leftSideWidth, filesBranchesBoundary-1)
  if err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    sideView.Highlight = true
    sideView.Title = "Files"
    refreshFiles(g)
  }

  if v, err := g.SetView("main", leftSideWidth+2, 0, maxX-1, optionsTop-1); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Diff"
    v.Wrap = true
    if _, err := g.SetCurrentView("files"); err != nil {
      return err
    }
    handleFileSelect(g, sideView)
  }

  if v, err := g.SetView("branches", 0, filesBranchesBoundary, leftSideWidth, optionsTop-1); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Branches"

    // these are only called once
    refreshBranches(v)
    nextView(g, nil)
  }

  if v, err := g.SetView("options", 0, optionsTop, maxX-1, optionsTop+2); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Options"
  }

  return nil
}

func renderString(g *gocui.Gui, s string) error {
  g.Update(func(*gocui.Gui) error {
    v, err := g.View("main")
    if err != nil {
      panic(err)
    }
    v.Clear()
    fmt.Fprint(v, s)
    v.Wrap = true
    return nil
  })
  return nil
}

func run() {
  g, err := gocui.NewGui(gocui.OutputNormal)
  if err != nil {
    log.Panicln(err)
  }
  defer g.Close()

  // g.Cursor = true

  g.SetManagerFunc(layout)

  if err := keybindings(g); err != nil {
    log.Panicln(err)
  }

  if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
    log.Panicln(err)
  }
}

// const mcRide = "
//                                                     `.-::-`
//                                                  -/o+oossys+:.
//                                                `+o++++++osssyys/`
//                    ://-:+.`     .::-.       . `++oyyo/-/+oooosyhy.
//                     `-+sy::-:::/+o+yss+-...   /s++ss/:/:+osyoosydh`
//                        `-:/+o/:/+:/+-s/:s/o+`/++s++/:--/+shds+++yd:
//                            `y+/+soy:+/-o++y+yhyyyo/---/oyhddo/::od-
//                           .+o-``-+syysy//o:-oo+oyyyo+oyhyddds/oshy
//                        `:o++o+/-....-:/+oooyyh+:ooshhhhhhdddssyyy`
//                      .:o+/++ooosso//:::+yo.::hs+++:yhhhhdddhoyhh:
//                  `-/+so///+osyso-.:://++-` `:hhhdsohddhhhdddssh+
//                -+oso++ssoyys:.`              ydddddddddddhho+yd+
//             `:sysssssssydh:`    `-:::-..-...`ydddddddddyso++shds
//           `/syyysssyyhhdd+``..://+ooo/++ssssoyddddddhho/:::oyhdhs-`
//          -syyyysssyhhddhyo++++/::+/+/-:::///+sddddhs//+o+/ososyhhs+/.`
//        `+hhyyyyyyyhddhs+///://///+ooo/::+o++osyhyyys+--+//o//oosyys++++:..``
//       .sddhyhyyyhddyso++/::://////+syo/:osssssyhsssoooosoo//+ossssyssooooo+++:.
//       .hdhhhhhhhhhysssssysssssssyyyhddso+soyhhhsssooosyyssso+syysoososoo/++osyo/
//        -syyyyyyyyyyyyyyyyyyo/::----:shdsyo+yysyyyssssosyysos+/+++/+ooo++:/+/ooss/
//          `........----..``           odhyyyhhsysoss++oysso++s/++++syys++/:::/:+sy-
//                                      `ydyssyysyoyyo+sysyys++s+++++ooo+osss+/+++syy
//                                       /dysyssoyyoo+oyyshss//:---:/++++oshhysooosyh`
//                                       .dhhhyysyyys++yyyyss+--:::/:///oshddhhyo+osy`
//                                        yddhhyyssy+//ssyyso/-:://+ooosyhddhsoo+/+so
//                                        +ddhhyysss+osyyysss:::/oyyhhyhddddds+///oy/
//                                        /dddhhyyyssysssssss+++ooyhdddddddhdyo///yyo
//                                        /dddhyyyyyysssoo+/:-/oshhdddddddssdds+//sys
//                                        +ddhhyyhhy/oo+/:::::+syhddddddds -hdyo++ohh`
//                                        sddhhysyysoys/:::::osyhdddddddy`  sdhsosohh:
//                                       `dddddhhhhhhhyo:-/ossoshddddhhd-   .ddyssohh/"
