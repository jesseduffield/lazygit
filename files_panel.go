package main

import (

  // "io"
  // "io/ioutil"

  // "strings"

  "errors"
  "strings"

  "github.com/fatih/color"
  "github.com/jroimartin/gocui"
)

var (
  // ErrNoFiles : when there are no modified files in the repo
  ErrNoFiles = errors.New("No changed files")
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
  file, err := getSelectedFile(v)
  if err != nil {
    return err
  }

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

func getSelectedFile(v *gocui.View) (GitFile, error) {
  if len(state.GitFiles) == 0 {
    return GitFile{}, ErrNoFiles
  }
  lineNumber := getItemPosition(v)
  return state.GitFiles[lineNumber], nil
}

func handleFileRemove(g *gocui.Gui, v *gocui.View) error {
  file, err := getSelectedFile(v)
  if err != nil {
    return err
  }
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

func handleIgnoreFile(g *gocui.Gui, v *gocui.View) error {
  file, err := getSelectedFile(v)
  if err != nil {
    return err
  }
  if file.Tracked {
    return createErrorPanel(g, "Cannot ignore tracked files")
  }
  gitIgnore(file.Name)
  return refreshFiles(g)
}

func handleFileSelect(g *gocui.Gui, v *gocui.View) error {
  baseString := "tab: switch to branches, space: toggle staged, c: commit changes, o: open, s: open in sublime, i: ignore"
  item, err := getSelectedFile(v)
  if err != nil {
    if err != ErrNoFiles {
      return err
    }
    renderString(g, "main", "No changed files")
    colorLog(color.FgRed, "error")
    return renderString(g, "options", baseString)
  }
  var optionsString string
  if item.Tracked {
    optionsString = baseString + ", r: checkout"
  } else {
    optionsString = baseString + ", r: delete"
  }
  renderString(g, "options", optionsString)
  diff := getDiff(item)
  return renderString(g, "main", diff)
}

func genericFileOpen(g *gocui.Gui, v *gocui.View, open func(string) (string, error)) error {
  file, err := getSelectedFile(v)
  if err != nil {
    return err
  }
  _, err = open(file.Name)
  return err
}

func handleFileOpen(g *gocui.Gui, v *gocui.View) error {
  return genericFileOpen(g, v, openFile)
}
func handleSublimeFileOpen(g *gocui.Gui, v *gocui.View) error {
  return genericFileOpen(g, v, sublimeOpenFile)
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
      refreshCommits(g)
      refreshFiles(g)
      refreshStatus(g)
      devLog("pulled.")
    }
  }()
  return nil
}

func pushFiles(g *gocui.Gui, v *gocui.View) error {
  devLog("pushing...")
  createSimpleConfirmationPanel(g, v, "", "Pushing...")
  go func() {
    if output, err := gitPush(); err != nil {
      createSimpleConfirmationPanel(g, v, "Error", output)
    } else {
      closeConfirmationPrompt(g)
      refreshCommits(g)
      refreshStatus(g)
      devLog("pushed.")
    }
  }()
  return nil
}
