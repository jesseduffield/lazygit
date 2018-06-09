// though this panel is called the merge panel, it's really going to use the main panel. This may change in the future

package main

import (
  "bufio"
  "bytes"
  "io/ioutil"
  "os"
  "strings"

  "github.com/fatih/color"
  "github.com/jesseduffield/gocui"
)

func findConflicts(content string) ([]conflict, error) {
  conflicts := make([]conflict, 0)
  var newConflict conflict
  for i, line := range splitLines(content) {
    if line == "<<<<<<< HEAD" {
      newConflict = conflict{start: i}
    } else if line == "=======" {
      newConflict.middle = i
    } else if strings.HasPrefix(line, ">>>>>>> ") {
      newConflict.end = i
      conflicts = append(conflicts, newConflict)
    }
  }
  return conflicts, nil
}

func shiftConflict(conflicts []conflict) (conflict, []conflict) {
  return conflicts[0], conflicts[1:]
}

func shouldHighlightLine(index int, conflict conflict, top bool) bool {
  return (index >= conflict.start && index <= conflict.middle && top) || (index >= conflict.middle && index <= conflict.end && !top)
}

func coloredConflictFile(content string, conflicts []conflict, conflictIndex int, conflictTop, hasFocus bool) (string, error) {
  if len(conflicts) == 0 {
    return content, nil
  }
  conflict, remainingConflicts := shiftConflict(conflicts)
  var outputBuffer bytes.Buffer
  for i, line := range splitLines(content) {
    colourAttr := color.FgWhite
    if i == conflict.start || i == conflict.middle || i == conflict.end {
      colourAttr = color.FgRed
    }
    colour := color.New(colourAttr)
    if hasFocus && conflictIndex < len(conflicts) && conflicts[conflictIndex] == conflict && shouldHighlightLine(i, conflict, conflictTop) {
      colour.Add(color.Bold)
    }
    if i == conflict.end && len(remainingConflicts) > 0 {
      conflict, remainingConflicts = shiftConflict(remainingConflicts)
    }
    outputBuffer.WriteString(coloredString(line, colour) + "\n")
  }
  return outputBuffer.String(), nil
}

func handleSelectTop(g *gocui.Gui, v *gocui.View) error {
  state.ConflictTop = true
  return refreshMergePanel(g)
}

func handleSelectBottom(g *gocui.Gui, v *gocui.View) error {
  state.ConflictTop = false
  return refreshMergePanel(g)
}

func handleSelectNextConflict(g *gocui.Gui, v *gocui.View) error {
  if state.ConflictIndex >= len(state.Conflicts)-1 {
    return nil
  }
  state.ConflictIndex++
  return refreshMergePanel(g)
}

func handleSelectPrevConflict(g *gocui.Gui, v *gocui.View) error {
  if state.ConflictIndex <= 0 {
    return nil
  }
  state.ConflictIndex--
  return refreshMergePanel(g)
}

func isIndexToDelete(i int, conflict conflict, top bool) bool {
  return i == conflict.middle ||
    i == conflict.start ||
    i == conflict.end ||
    (!top && i > conflict.start && i < conflict.middle) ||
    (top && i > conflict.middle && i < conflict.end)
}

func resolveConflict(filename string, conflict conflict, top bool) error {
  file, err := os.Open(filename)
  if err != nil {
    return err
  }
  defer file.Close()

  reader := bufio.NewReader(file)
  output := ""
  for i := 0; true; i++ {
    line, err := reader.ReadString('\n')
    if err != nil {
      break
    }
    if !isIndexToDelete(i, conflict, top) {
      output += line
    }
  }
  devLog(output)
  return ioutil.WriteFile(filename, []byte(output), 0644)
}

func pushFileSnapshot(filename string) error {
  content, err := catFile(filename)
  if err != nil {
    return err
  }
  state.EditHistory.Push(content)
  return nil
}

func handlePopFileSnapshot(g *gocui.Gui, v *gocui.View) error {
  colorLog(color.FgCyan, "IM HERE")
  if state.EditHistory.Len() == 0 {
    return nil
  }
  prevContent := state.EditHistory.Pop().(string)
  gitFile, err := getSelectedFile(g)
  if err != nil {
    return err
  }
  ioutil.WriteFile(gitFile.Name, []byte(prevContent), 0644)
  return refreshMergePanel(g)
}

func handlePickConflict(g *gocui.Gui, v *gocui.View) error {
  conflict := state.Conflicts[state.ConflictIndex]
  gitFile, err := getSelectedFile(g)
  if err != nil {
    return err
  }
  pushFileSnapshot(gitFile.Name)
  err = resolveConflict(gitFile.Name, conflict, state.ConflictTop)
  if err != nil {
    panic(err)
  }
  return refreshMergePanel(g)
}

func currentViewName(g *gocui.Gui) string {
  currentView := g.CurrentView()
  return currentView.Name()
}

func refreshMergePanel(g *gocui.Gui) error {
  cat, err := catSelectedFile(g)
  if err != nil {
    return err
  }
  state.Conflicts, err = findConflicts(cat)
  if err != nil {
    return err
  }

  if len(state.Conflicts) == 0 {
    state.ConflictIndex = 0
  } else if state.ConflictIndex > len(state.Conflicts)-1 {
    state.ConflictIndex = len(state.Conflicts) - 1
  }
  hasFocus := currentViewName(g) == "main"
  if hasFocus {
    renderMergeOptions(g)
  }
  content, err := coloredConflictFile(cat, state.Conflicts, state.ConflictIndex, state.ConflictTop, hasFocus)
  if err != nil {
    return err
  }
  return renderString(g, "main", content)
}

func switchToMerging(g *gocui.Gui) error {
  state.ConflictIndex = 0
  state.ConflictTop = true
  _, err := g.SetCurrentView("main")
  if err != nil {
    return err
  }
  return refreshMergePanel(g)
}

func renderMergeOptions(g *gocui.Gui) error {
  return renderOptionsMap(g, map[string]string{
    "up/down":    "pick hunk",
    "left/right": "previous/next commit",
    "space":      "pick hunk",
    "z":          "undo",
  })
}

func handleEscapeMerge(g *gocui.Gui, v *gocui.View) error {
  filesView, err := g.View("files")
  if err != nil {
    return err
  }
  refreshFiles(g)
  return switchFocus(g, v, filesView)
}
