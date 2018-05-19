// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "fmt"
  // "io"
  // "io/ioutil"
  "log"
  // "strings"
  "os"
  "github.com/jroimartin/gocui"
  "github.com/fatih/color"
)

type gitFile struct {
  Name string
  Staged bool
}

var gitFiles []gitFile

func nextView(g *gocui.Gui, v *gocui.View) error {
  if v == nil || v.Name() == "side" {
    _, err := g.SetCurrentView("main")
    return err
  }
  _, err := g.SetCurrentView("side")
  return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
  if v != nil {
    cx, cy := v.Cursor()
    if err := v.SetCursor(cx, cy+1); err != nil {
      ox, oy := v.Origin()
      if err := v.SetOrigin(ox, oy+1); err != nil {
        return err
      }
    }
  }

  // refresh main panel's text to match newly selected item
  return handleItemSelect(g, v)
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
  if v != nil {
    ox, oy := v.Origin()
    cx, cy := v.Cursor()
    if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
      if err := v.SetOrigin(ox, oy-1); err != nil {
        return err
      }
    }
  }

  // refresh main panel's text to match newly selected item
  return handleItemSelect(g, v)
}

func devLog(s string) {
  f, _ := os.OpenFile("development.log", os.O_APPEND|os.O_WRONLY, 0644)
  defer f.Close()

  f.WriteString(s + "\n")
}

func handleItemPress(g *gocui.Gui, v *gocui.View) error {
  item := getItem(v)

  if item.Staged {
    unStageFile(item.Name)
  } else {
    stageFile(item.Name)
  }

  if err := refreshList(v); err != nil {
    return err
  }
  return nil
}

func getItem(v *gocui.View) gitFile {
  _, lineNumber := v.Cursor()
  if lineNumber >= len(gitFiles) {
    return gitFiles[len(gitFiles) - 1]
  }
  return gitFiles[lineNumber]
}

func handleItemSelect(g *gocui.Gui, v *gocui.View) error {
  item := getItem(v)
  diff := getDiff(item.Name, item.Staged)
  devLog(diff)
  if err := renderString(g, diff); err != nil {
    return err
  }

  // maxX, maxY := g.Size()
  // if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
  //   if err != gocui.ErrUnknownView {
  //     return err
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
  if _, err := g.SetCurrentView("side"); err != nil {
    return err
  }
  return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
  return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
  if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
    return err
  }
  if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
    return err
  }
  if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
    return err
  }
  if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
    return err
  }
  if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
    return err
  }
  if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
    return err
  }
  if err := g.SetKeybinding("side", gocui.KeySpace, gocui.ModNone, handleItemPress); err != nil {
    return err
  }
  // if err := g.SetKeybinding("msg", gocui.KeySpace, gocui.ModNone, delMsg); err != nil {
  //   return err
  // }
  return nil
}

func refreshList(v *gocui.View) error {
  // get files to stage
  statusString, _ := runCommand("git status")
  filesToStage := getFilesToStage(statusString)
  filesToUnstage := getFilesToUnstage(statusString)
  // v.Highlight = true
  // v.SelBgColor = gocui.ColorWhite
  // v.SelFgColor = gocui.ColorBlack
  v.Clear()
  gitFiles = make([]gitFile, 0)
  red := color.New(color.FgRed)
  for _, file := range filesToStage {
    gitFiles = append(gitFiles, gitFile{file, false})
    red.Fprintln(v, file)
  }
  green := color.New(color.FgGreen)
  for _, file := range filesToUnstage {
    gitFiles = append(gitFiles, gitFile{file, true})
    green.Fprintln(v, file)
  }
  devLog(fmt.Sprint(gitFiles))
  return nil
}

func layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  sideView, err := g.SetView("side", -1, -1, 30, maxY)
  if err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    sideView.Title = "Files"
    devLog("test")
    refreshList(sideView)
  }

  if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Editable = true
    v.Wrap = true
    if _, err := g.SetCurrentView("side"); err != nil {
      return err
    }
    handleItemSelect(g, sideView)
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

  g.Cursor = true

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
