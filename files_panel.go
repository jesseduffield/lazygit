// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (

  // "io"
  // "io/ioutil"

  // "strings"

  "strings"

  "github.com/fatih/color"
  "github.com/jroimartin/gocui"
)

func stagedFiles(files []GitFile) []GitFile {
  result := make([]GitFile, 0)
  for _, file := range files {
    if file.HasStagedChanges {
      result = append(result, file)
    }
  }
  return result
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
  if len(state.GitFiles) == 0 {
    // find a way to not have to do this
    return GitFile{
      Name:               "noFile",
      DisplayString:      "none",
      HasStagedChanges:   false,
      HasUnstagedChanges: false,
      Tracked:            false,
      Deleted:            false,
    }
  }
  return state.GitFiles[lineNumber]
}

func handleFileRemove(g *gocui.Gui, v *gocui.View) error {
  file := getSelectedFile(v)
  var deleteVerb string
  if file.Tracked {
    deleteVerb = "checkout"
  } else {
    deleteVerb = "delete"
  }
  return createConfirmationPanel(g, v, strings.Title(deleteVerb)+" file", "Are you sure you want to "+deleteVerb+" "+file.Name+" (you will lose your changes)? (y/n)", func(g *gocui.Gui, v *gocui.View) error {
    if err := removeFile(file); err != nil {
      panic(err)
    }
    return refreshFiles(g)
  }, nil)
}

func handleFileSelect(g *gocui.Gui, v *gocui.View) error {
  item := getSelectedFile(v)
  var optionsString string
  baseString := "space: toggle staged, c: commit changes, option+o: open"
  if item.Tracked {
    optionsString = baseString + ", option+d: checkout"
  } else {
    optionsString = baseString + ", option+d: delete"
  }
  renderString(g, "options", optionsString)
  diff := getDiff(item)
  return renderString(g, "main", diff)
}

func handleFileOpen(g *gocui.Gui, v *gocui.View) error {
  file := getSelectedFile(v)
  _, err := openFile(file.Name)
  return err
}

func handleSublimeFileOpen(g *gocui.Gui, v *gocui.View) error {
  file := getSelectedFile(v)
  _, err := sublimeOpenFile(file.Name)
  return err
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
  correctCursor(filesView)
  return nil
}

func pullFiles(g *gocui.Gui, v *gocui.View) error {
  devLog("pulling...")
  createSimpleConfirmationPanel(g, v, "", "Pulling...")
  go func() {
    if output, err := gitPull(); err != nil {
      createSimpleConfirmationPanel(g, v, "Error", output)
    } else {
      closeConfirmationPrompt(g)
    }
  }()
  devLog("pulled.")
  return refreshFiles(g)
}

func pushFiles(g *gocui.Gui, v *gocui.View) error {
  devLog("pushing...")
  createSimpleConfirmationPanel(g, v, "", "Pushing...")
  go func() {
    if output, err := gitPush(); err != nil {
      createSimpleConfirmationPanel(g, v, "Error", output)
    } else {
      closeConfirmationPrompt(g)
    }
  }()
  devLog("pushed.")
  return nil
}
