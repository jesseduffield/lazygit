package main

import (

  // "io"
  // "io/ioutil"

  "log"
  "time"
  // "strings"
  "github.com/jesseduffield/gocui"
)

type stateType struct {
  GitFiles     []GitFile
  Branches     []Branch
  Commits      []Commit
  StashEntries []StashEntry
  PreviousView string
}

var state = stateType{
  GitFiles:     make([]GitFile, 0),
  PreviousView: "files",
  Commits:      make([]Commit, 0),
  StashEntries: make([]StashEntry, 0),
}

func scrollUpMain(g *gocui.Gui, v *gocui.View) error {
  mainView, _ := g.View("main")
  ox, oy := mainView.Origin()
  if oy >= 1 {
    return mainView.SetOrigin(ox, oy-1)
  }
  return nil
}

func scrollDownMain(g *gocui.Gui, v *gocui.View) error {
  mainView, _ := g.View("main")
  ox, oy := mainView.Origin()
  if oy < len(mainView.BufferLines()) {
    return mainView.SetOrigin(ox, oy+1)
  }
  return nil
}

func keybindings(g *gocui.Gui) error {
  if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
    return err
  }
  if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
    return err
  }
  if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
    return err
  }
  if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
    return err
  }
  if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
    return err
  }
  if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, scrollUpMain); err != nil {
    return err
  }
  if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, scrollDownMain); err != nil {
    return err
  }
  if err := g.SetKeybinding("files", 'c', gocui.ModNone, handleCommitPress); err != nil {
    return err
  }
  if err := g.SetKeybinding("files", gocui.KeySpace, gocui.ModNone, handleFilePress); err != nil {
    return err
  }
  if err := g.SetKeybinding("files", 'd', gocui.ModNone, handleFileRemove); err != nil {
    return err
  }
  if err := g.SetKeybinding("files", 'o', gocui.ModNone, handleFileOpen); err != nil {
    return err
  }
  if err := g.SetKeybinding("files", 's', gocui.ModNone, handleSublimeFileOpen); err != nil {
    return err
  }
  if err := g.SetKeybinding("", 'P', gocui.ModNone, pushFiles); err != nil {
    return err
  }
  if err := g.SetKeybinding("", 'p', gocui.ModNone, pullFiles); err != nil {
    return err
  }
  if err := g.SetKeybinding("files", 'i', gocui.ModNone, handleIgnoreFile); err != nil {
    return err
  }
  if err := g.SetKeybinding("files", 'S', gocui.ModNone, handleStashSave); err != nil {
    return err
  }
  if err := g.SetKeybinding("branches", gocui.KeySpace, gocui.ModNone, handleBranchPress); err != nil {
    return err
  }
  if err := g.SetKeybinding("branches", 'F', gocui.ModNone, handleForceCheckout); err != nil {
    return err
  }
  if err := g.SetKeybinding("branches", 'n', gocui.ModNone, handleNewBranch); err != nil {
    return err
  }
  if err := g.SetKeybinding("commits", 's', gocui.ModNone, handleCommitSquashDown); err != nil {
    return err
  }
  if err := g.SetKeybinding("commits", 'r', gocui.ModNone, handleRenameCommit); err != nil {
    return err
  }
  if err := g.SetKeybinding("commits", 'g', gocui.ModNone, handleResetToCommit); err != nil {
    return err
  }
  if err := g.SetKeybinding("stash", gocui.KeySpace, gocui.ModNone, handleStashApply); err != nil {
    return err
  }
  // TODO: come up with a better keybinding (p/P used for pushing/pulling which
  // I'd like to be global. Perhaps all global keybindings should use a modifier
  // like command? But then there's gonna be hotkey conflicts with the terminal
  if err := g.SetKeybinding("stash", 'k', gocui.ModNone, handleStashPop); err != nil {
    return err
  }
  if err := g.SetKeybinding("stash", 'd', gocui.ModNone, handleStashDrop); err != nil {
    return err
  }
  return nil
}

func layout(g *gocui.Gui) error {
  width, height := g.Size()
  leftSideWidth := width / 3
  statusFilesBoundary := 2
  filesBranchesBoundary := height - 20
  commitsBranchesBoundary := height - 10
  commitsStashBoundary := height - 5

  optionsTop := height - 2
  // hiding options if there's not enough space
  if height < 30 {
    optionsTop = height - 1
  }

  filesView, err := g.SetView("files", 0, statusFilesBoundary+1, leftSideWidth, filesBranchesBoundary-1)
  if err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    filesView.Highlight = true
    filesView.Title = "Files"
  }

  if v, err := g.SetView("status", 0, statusFilesBoundary-2, leftSideWidth, statusFilesBoundary); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Status"
  }

  mainView, err := g.SetView("main", leftSideWidth+1, 0, width-1, optionsTop)
  if err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    mainView.Title = "Diff"
    mainView.Wrap = true
  }

  if v, err := g.SetView("branches", 0, filesBranchesBoundary, leftSideWidth, commitsBranchesBoundary-1); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Branches"

  }

  if v, err := g.SetView("commits", 0, commitsBranchesBoundary, leftSideWidth, commitsStashBoundary-1); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Commits"

  }

  if v, err := g.SetView("stash", 0, commitsStashBoundary, leftSideWidth, optionsTop); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.Title = "Stash"
  }

  if v, err := g.SetView("options", -1, optionsTop, width, optionsTop+2); err != nil {
    if err != gocui.ErrUnknownView {
      return err
    }
    v.BgColor = gocui.ColorBlue
    v.Frame = false
    v.Title = "Options"

    // these are only called once
    handleFileSelect(g, filesView)
    refreshFiles(g)
    refreshBranches(g)
    refreshCommits(g)
    refreshStashEntries(g)
    nextView(g, nil)
  }

  return nil
}

func fetch(g *gocui.Gui) {
  gitFetch()
  refreshStatus(g)
}

func run() {
  g, err := gocui.NewGui(gocui.OutputNormal)
  if err != nil {
    log.Panicln(err)
  }
  defer g.Close()

  // periodically fetching to check for upstream differences
  go func() {
    for range time.Tick(time.Second * 60) {
      fetch(g)
    }
  }()

  g.SetManagerFunc(layout)

  if err := keybindings(g); err != nil {
    log.Panicln(err)
  }

  if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
    log.Panicln(err)
  }
}

func quit(g *gocui.Gui, v *gocui.View) error {
  return gocui.ErrQuit
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
