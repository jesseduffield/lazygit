// though this panel is called the merge panel, it's really going to use the main panel. This may change in the future

package gui

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) findConflicts(content string) ([]commands.Conflict, error) {
	conflicts := make([]commands.Conflict, 0)
	var newConflict commands.Conflict
	for i, line := range utils.SplitLines(content) {
		trimmedLine := strings.TrimPrefix(line, "++")
		gui.Log.Info(trimmedLine)
		if trimmedLine == "<<<<<<< HEAD" || trimmedLine == "<<<<<<< MERGE_HEAD" || trimmedLine == "<<<<<<< Updated upstream" {
			newConflict = commands.Conflict{Start: i}
		} else if trimmedLine == "=======" {
			newConflict.Middle = i
		} else if strings.HasPrefix(trimmedLine, ">>>>>>> ") {
			newConflict.End = i
			conflicts = append(conflicts, newConflict)
		}
	}
	return conflicts, nil
}

func (gui *Gui) shiftConflict(conflicts []commands.Conflict) (commands.Conflict, []commands.Conflict) {
	return conflicts[0], conflicts[1:]
}

func (gui *Gui) shouldHighlightLine(index int, conflict commands.Conflict, top bool) bool {
	return (index >= conflict.Start && index <= conflict.Middle && top) || (index >= conflict.Middle && index <= conflict.End && !top)
}

func (gui *Gui) coloredConflictFile(content string, conflicts []commands.Conflict, conflictIndex int, conflictTop, hasFocus bool) (string, error) {
	if len(conflicts) == 0 {
		return content, nil
	}
	conflict, remainingConflicts := gui.shiftConflict(conflicts)
	var outputBuffer bytes.Buffer
	for i, line := range utils.SplitLines(content) {
		colourAttr := color.FgWhite
		if i == conflict.Start || i == conflict.Middle || i == conflict.End {
			colourAttr = color.FgRed
		}
		colour := color.New(colourAttr)
		if hasFocus && conflictIndex < len(conflicts) && conflicts[conflictIndex] == conflict && gui.shouldHighlightLine(i, conflict, conflictTop) {
			colour.Add(color.Bold)
		}
		if i == conflict.End && len(remainingConflicts) > 0 {
			conflict, remainingConflicts = gui.shiftConflict(remainingConflicts)
		}
		outputBuffer.WriteString(utils.ColoredStringDirect(line, colour) + "\n")
	}
	return outputBuffer.String(), nil
}

func (gui *Gui) handleSelectTop(g *gocui.Gui, v *gocui.View) error {
	gui.State.Panels.Merging.ConflictTop = true
	return gui.refreshMergePanel()
}

func (gui *Gui) handleSelectBottom(g *gocui.Gui, v *gocui.View) error {
	gui.State.Panels.Merging.ConflictTop = false
	return gui.refreshMergePanel()
}

func (gui *Gui) handleSelectNextConflict(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Merging.ConflictIndex >= len(gui.State.Panels.Merging.Conflicts)-1 {
		return nil
	}
	gui.State.Panels.Merging.ConflictIndex++
	return gui.refreshMergePanel()
}

func (gui *Gui) handleSelectPrevConflict(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Merging.ConflictIndex <= 0 {
		return nil
	}
	gui.State.Panels.Merging.ConflictIndex--
	return gui.refreshMergePanel()
}

func (gui *Gui) isIndexToDelete(i int, conflict commands.Conflict, pick string) bool {
	return i == conflict.Middle ||
		i == conflict.Start ||
		i == conflict.End ||
		pick != "both" &&
			(pick == "bottom" && i > conflict.Start && i < conflict.Middle) ||
		(pick == "top" && i > conflict.Middle && i < conflict.End)
}

func (gui *Gui) resolveConflict(g *gocui.Gui, conflict commands.Conflict, pick string) error {
	gitFile, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}
	file, err := os.Open(gitFile.Name)
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
		if !gui.isIndexToDelete(i, conflict, pick) {
			output += line
		}
	}
	gui.Log.Info(output)
	return ioutil.WriteFile(gitFile.Name, []byte(output), 0644)
}

func (gui *Gui) pushFileSnapshot(g *gocui.Gui) error {
	gitFile, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}
	content, err := gui.GitCommand.CatFile(gitFile.Name)
	if err != nil {
		return err
	}
	gui.State.Panels.Merging.EditHistory.Push(content)
	return nil
}

func (gui *Gui) handlePopFileSnapshot(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Merging.EditHistory.Len() == 0 {
		return nil
	}
	prevContent := gui.State.Panels.Merging.EditHistory.Pop().(string)
	gitFile, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}
	ioutil.WriteFile(gitFile.Name, []byte(prevContent), 0644)
	return gui.refreshMergePanel()
}

func (gui *Gui) handlePickHunk(g *gocui.Gui, v *gocui.View) error {
	conflict := gui.State.Panels.Merging.Conflicts[gui.State.Panels.Merging.ConflictIndex]
	gui.pushFileSnapshot(g)
	pick := "bottom"
	if gui.State.Panels.Merging.ConflictTop {
		pick = "top"
	}
	err := gui.resolveConflict(g, conflict, pick)
	if err != nil {
		panic(err)
	}

	// if that was the last conflict, finish the merge for this file
	if len(gui.State.Panels.Merging.Conflicts) == 1 {
		if err := gui.handleCompleteMerge(); err != nil {
			return err
		}
	}
	return gui.refreshMergePanel()
}

func (gui *Gui) handlePickBothHunks(g *gocui.Gui, v *gocui.View) error {
	conflict := gui.State.Panels.Merging.Conflicts[gui.State.Panels.Merging.ConflictIndex]
	gui.pushFileSnapshot(g)
	err := gui.resolveConflict(g, conflict, "both")
	if err != nil {
		panic(err)
	}
	return gui.refreshMergePanel()
}

func (gui *Gui) refreshMergePanel() error {
	panelState := gui.State.Panels.Merging
	cat, err := gui.catSelectedFile(gui.g)
	if err != nil {
		return err
	}
	if cat == "" {
		return nil
	}
	panelState.Conflicts, err = gui.findConflicts(cat)
	if err != nil {
		return err
	}

	// handle potential fixes that the user made in their editor since we last refreshed
	if len(panelState.Conflicts) == 0 {
		return gui.handleCompleteMerge()
	} else if panelState.ConflictIndex > len(panelState.Conflicts)-1 {
		panelState.ConflictIndex = len(panelState.Conflicts) - 1
	}

	hasFocus := gui.currentViewName() == "main"
	content, err := gui.coloredConflictFile(cat, panelState.Conflicts, panelState.ConflictIndex, panelState.ConflictTop, hasFocus)
	if err != nil {
		return err
	}
	if err := gui.renderString(gui.g, "main", content); err != nil {
		return err
	}
	if err := gui.scrollToConflict(gui.g); err != nil {
		return err
	}

	mainView := gui.getMainView()
	mainView.Wrap = false

	return nil
}

func (gui *Gui) scrollToConflict(g *gocui.Gui) error {
	panelState := gui.State.Panels.Merging
	if len(panelState.Conflicts) == 0 {
		return nil
	}
	mergingView := gui.getMainView()
	conflict := panelState.Conflicts[panelState.ConflictIndex]
	ox, _ := mergingView.Origin()
	_, height := mergingView.Size()
	conflictMiddle := (conflict.End + conflict.Start) / 2
	newOriginY := int(math.Max(0, float64(conflictMiddle-(height/2))))
	gui.g.Update(func(g *gocui.Gui) error {
		return mergingView.SetOrigin(ox, newOriginY)
	})
	return nil
}

func (gui *Gui) renderMergeOptions() error {
	return gui.renderOptionsMap(map[string]string{
		"↑ ↓":   gui.Tr.SLocalize("selectHunk"),
		"← →":   gui.Tr.SLocalize("navigateConflicts"),
		"space": gui.Tr.SLocalize("pickHunk"),
		"b":     gui.Tr.SLocalize("pickBothHunks"),
		"z":     gui.Tr.SLocalize("undo"),
	})
}

func (gui *Gui) handleEscapeMerge(g *gocui.Gui, v *gocui.View) error {
	gui.State.Panels.Merging.EditHistory = stack.New()
	if err := gui.refreshFiles(); err != nil {
		return err
	}
	// it's possible this method won't be called from the merging view so we need to
	// ensure we only 'return' focus if we already have it
	if gui.g.CurrentView() == gui.getMainView() {
		return gui.switchFocus(g, v, gui.getFilesView())
	}
	return nil
}

func (gui *Gui) handleCompleteMerge() error {
	if err := gui.stageSelectedFile(gui.g); err != nil {
		return err
	}
	if err := gui.refreshFiles(); err != nil {
		return err
	}
	// if we got conflicts after unstashing, we don't want to call any git
	// commands to continue rebasing/merging here
	if gui.State.WorkingTreeState == "normal" {
		return gui.handleEscapeMerge(gui.g, gui.getMainView())
	}
	// if there are no more files with merge conflicts, we should ask whether the user wants to continue
	if !gui.anyFilesWithMergeConflicts() {
		return gui.promptToContinue()
	}
	return gui.handleEscapeMerge(gui.g, gui.getMainView())
}

// promptToContinue asks the user if they want to continue the rebase/merge that's in progress
func (gui *Gui) promptToContinue() error {
	return gui.createConfirmationPanel(gui.g, gui.getFilesView(), "continue", gui.Tr.SLocalize("ConflictsResolved"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.genericMergeCommand("continue")
	}, nil)
}
