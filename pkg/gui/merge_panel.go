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
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) findConflicts(content string) ([]commands.Conflict, error) {
	conflicts := make([]commands.Conflict, 0)
	var newConflict commands.Conflict
	for i, line := range utils.SplitLines(content) {
		if line == "<<<<<<< HEAD" || line == "<<<<<<< MERGE_HEAD" || line == "<<<<<<< Updated upstream" {
			newConflict = commands.Conflict{Start: i}
		} else if line == "=======" {
			newConflict.Middle = i
		} else if strings.HasPrefix(line, ">>>>>>> ") {
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
	gui.State.ConflictTop = true
	return gui.refreshMergePanel(g)
}

func (gui *Gui) handleSelectBottom(g *gocui.Gui, v *gocui.View) error {
	gui.State.ConflictTop = false
	return gui.refreshMergePanel(g)
}

func (gui *Gui) handleSelectNextConflict(g *gocui.Gui, v *gocui.View) error {
	if gui.State.ConflictIndex >= len(gui.State.Conflicts)-1 {
		return nil
	}
	gui.State.ConflictIndex++
	return gui.refreshMergePanel(g)
}

func (gui *Gui) handleSelectPrevConflict(g *gocui.Gui, v *gocui.View) error {
	if gui.State.ConflictIndex <= 0 {
		return nil
	}
	gui.State.ConflictIndex--
	return gui.refreshMergePanel(g)
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
	gui.State.EditHistory.Push(content)
	return nil
}

func (gui *Gui) handlePopFileSnapshot(g *gocui.Gui, v *gocui.View) error {
	if gui.State.EditHistory.Len() == 0 {
		return nil
	}
	prevContent := gui.State.EditHistory.Pop().(string)
	gitFile, err := gui.getSelectedFile(g)
	if err != nil {
		return err
	}
	ioutil.WriteFile(gitFile.Name, []byte(prevContent), 0644)
	return gui.refreshMergePanel(g)
}

func (gui *Gui) handlePickHunk(g *gocui.Gui, v *gocui.View) error {
	conflict := gui.State.Conflicts[gui.State.ConflictIndex]
	gui.pushFileSnapshot(g)
	pick := "bottom"
	if gui.State.ConflictTop {
		pick = "top"
	}
	err := gui.resolveConflict(g, conflict, pick)
	if err != nil {
		panic(err)
	}
	gui.refreshMergePanel(g)
	return nil
}

func (gui *Gui) handlePickBothHunks(g *gocui.Gui, v *gocui.View) error {
	conflict := gui.State.Conflicts[gui.State.ConflictIndex]
	gui.pushFileSnapshot(g)
	err := gui.resolveConflict(g, conflict, "both")
	if err != nil {
		panic(err)
	}
	return gui.refreshMergePanel(g)
}

func (gui *Gui) refreshMergePanel(g *gocui.Gui) error {
	cat, err := gui.catSelectedFile(g)
	if err != nil {
		return err
	}
	if cat == "" {
		return nil
	}
	gui.State.Conflicts, err = gui.findConflicts(cat)
	if err != nil {
		return err
	}

	if len(gui.State.Conflicts) == 0 {
		return gui.handleCompleteMerge(g)
	} else if gui.State.ConflictIndex > len(gui.State.Conflicts)-1 {
		gui.State.ConflictIndex = len(gui.State.Conflicts) - 1
	}
	hasFocus := gui.currentViewName(g) == "main"
	if hasFocus {
		gui.renderMergeOptions(g)
	}
	content, err := gui.coloredConflictFile(cat, gui.State.Conflicts, gui.State.ConflictIndex, gui.State.ConflictTop, hasFocus)
	if err != nil {
		return err
	}
	if err := gui.scrollToConflict(g); err != nil {
		return err
	}
	return gui.renderString(g, "main", content)
}

func (gui *Gui) scrollToConflict(g *gocui.Gui) error {
	mainView, err := g.View("main")
	if err != nil {
		return err
	}
	if len(gui.State.Conflicts) == 0 {
		return nil
	}
	conflict := gui.State.Conflicts[gui.State.ConflictIndex]
	ox, _ := mainView.Origin()
	_, height := mainView.Size()
	conflictMiddle := (conflict.End + conflict.Start) / 2
	newOriginY := int(math.Max(0, float64(conflictMiddle-(height/2))))
	return mainView.SetOrigin(ox, newOriginY)
}

func (gui *Gui) switchToMerging(g *gocui.Gui) error {
	gui.State.ConflictIndex = 0
	gui.State.ConflictTop = true
	_, err := g.SetCurrentView("main")
	if err != nil {
		return err
	}
	return gui.refreshMergePanel(g)
}

func (gui *Gui) renderMergeOptions(g *gocui.Gui) error {
	return gui.renderOptionsMap(g, map[string]string{
		"↑ ↓":   gui.Tr.SLocalize("selectHunk"),
		"← →":   gui.Tr.SLocalize("navigateConflicts"),
		"space": gui.Tr.SLocalize("pickHunk"),
		"b":     gui.Tr.SLocalize("pickBothHunks"),
		"z":     gui.Tr.SLocalize("undo"),
	})
}

func (gui *Gui) handleEscapeMerge(g *gocui.Gui, v *gocui.View) error {
	filesView, err := g.View("files")
	if err != nil {
		return err
	}
	gui.refreshFiles(g)
	return gui.switchFocus(g, v, filesView)
}

func (gui *Gui) handleCompleteMerge(g *gocui.Gui) error {
	filesView, err := g.View("files")
	if err != nil {
		return err
	}
	gui.stageSelectedFile(g)
	gui.refreshFiles(g)
	return gui.switchFocus(g, nil, filesView)
}
