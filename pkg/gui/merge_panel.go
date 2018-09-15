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

// findConflicts searches for conflicts in the given content.
// returns the conflicts and if something went wrong an error.
func (gui *Gui) findConflicts(content string) ([]commands.Conflict, error) {
	var (
		newConflict commands.Conflict
		conflicts   = make([]commands.Conflict, 0)
	)

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

// shiftConflicts shifts the conflicts one position to the right
func (gui *Gui) shiftConflict(conflicts []commands.Conflict) (commands.Conflict, []commands.Conflict) {
	return conflicts[0], conflicts[1:]
}

// shouldHighlightLine returns whether or not to highlight the line.
// index: the index.
// conflict: the conflict.
// top: the conflictsTop boolean
// returns a boolean.
func (gui *Gui) shouldHighlightLine(index int, conflict commands.Conflict, top bool) bool {
	return (index >= conflict.Start && index <= conflict.Middle && top) || (index >= conflict.Middle && index <= conflict.End && !top)
}

// coloredConflictFile creates a color representation of the conflicts file
func (gui *Gui) coloredConflictFile(content string, conflicts []commands.Conflict, conflictIndex int, conflictTop, hasFocus bool) (string, error) {

	var (
		outputBuffer                 bytes.Buffer
		conflict, remainingConflicts = gui.shiftConflict(conflicts)
	)

	if len(conflicts) == 0 {
		return content, nil
	}

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

// handleSelectTop is called when the user wants to go to the top
// of the conflicts file.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleSelectTop(g *gocui.Gui, v *gocui.View) error {

	gui.State.ConflictTop = true

	err := gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at handleSelectTop: %s\n", err)
		return err
	}

	return nil
}

// handleSelectBottom is called when the user wants to go to the bottom
// of the conflicts file.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleSelectBottom(g *gocui.Gui, v *gocui.View) error {

	gui.State.ConflictTop = false

	err := gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at handleSelectBottom: %s\n", err)
		return err
	}

	return nil
}

// handleSelectNextConflict is called when the user wants to go to the next
// conflict.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleSelectNextConflict(g *gocui.Gui, v *gocui.View) error {

	if gui.State.ConflictIndex >= len(gui.State.Conflicts)-1 {
		return nil
	}

	gui.State.ConflictIndex++

	err := gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at handleSelectNextConflict: %s\n", err)
		return err
	}

	return nil
}

// handleSelectPrevConflict is called when the user wants to go to the previous
// conflict.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleSelectPrevConflict(g *gocui.Gui, v *gocui.View) error {

	if gui.State.ConflictIndex <= 0 {
		return nil
	}

	gui.State.ConflictIndex--

	err := gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at handleSelectPrevConflict: %s\n", err)
		return err
	}

	return nil
}

// isIndexToDelete returns if the index should be deleted
// i: TODO
// conflict: TODO
// pick: TODO
// returns a boolean indicating whether or not to do it.
func (gui *Gui) isIndexToDelete(i int, conflict commands.Conflict, pick string) bool {
	return i == conflict.Middle ||
		i == conflict.Start ||
		i == conflict.End ||
		(pick == "bottom" && i > conflict.Start && i < conflict.Middle) ||
		(pick == "top" && i > conflict.Middle && i < conflict.End)
}

// resolveConflicts is called to resolve the commit.
// conflict: TODO
// pick: TODO
// returns an error if something goes wrong.
func (gui *Gui) resolveConflict(conflict commands.Conflict, pick string) error {

	gitFile, err := gui.getSelectedFile()
	if err != nil {
		gui.Log.Errorf("Failed to getSelectedFiles at resolveConflic: %s\n", err)
		return err
	}

	file, err := os.Open(gitFile.Name)
	if err != nil {
		gui.Log.Errorf("Failed to open file at resolveConflict: %s\n", err)
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

	err = ioutil.WriteFile(gitFile.Name, []byte(output), 0644)
	if err != nil {
		gui.Log.Errorf("Failed to writeFile at resolveConflict: %s\n", err)
		return err
	}

	return nil
}

// pushFileSnapshot TODO
// returns an error when something goes wrong.
func (gui *Gui) pushFileSnapshot() error {

	gitFile, err := gui.getSelectedFile()
	if err != nil {
		gui.Log.Errorf("Failed to getSelectedFile: %s\n", err)
		return err
	}

	content, err := gui.GitCommand.CatFile(gitFile.Name)
	if err != nil {
		gui.Log.Errorf("Failed to cat file at pushFileSnapshot: %s\n", err)
		return err
	}

	gui.State.EditHistory.Push(content)

	return nil
}

// handlePopFileSnapshot is called user presses the keybind.
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
func (gui *Gui) handlePopFileSnapshot(g *gocui.Gui, v *gocui.View) error {

	if gui.State.EditHistory.Len() == 0 {
		return nil
	}

	prevContent := gui.State.EditHistory.Pop().(string)

	gitFile, err := gui.getSelectedFile()
	if err != nil {
		gui.Log.Errorf("Failed to getSelectedFile at handlePopFileSnapshot: %s\n", err)
		return err
	}

	err = ioutil.WriteFile(gitFile.Name, []byte(prevContent), 0644)
	if err != nil {
		gui.Log.Errorf("Failed to writeFile at handlePopFileSnapshot: %s\n", err)
		return err
	}

	err = gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at handlePopFileSnapshot: %s\n", err)
		return err
	}

	return nil
}

// handlePickHunk is called when the user selects a hunk in the merge panel.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handlePickHunk(g *gocui.Gui, v *gocui.View) error {
	var (
		conflict = gui.State.Conflicts[gui.State.ConflictIndex]
		pick     = "bottom"
	)

	err := gui.pushFileSnapshot()
	if err != nil {
		gui.Log.Errorf("Failed to pushFileSnapshot at handlePickHunk: %s\n", err)
		return err
	}

	if gui.State.ConflictTop {
		pick = "top"
	}

	err = gui.resolveConflict(conflict, pick)
	if err != nil {
		gui.Log.Errorf("Failed to resolveConflict at handlePickHunk: %s\n", err)
		return err
	}

	err = gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at handlePickHunk: %s\n", err)
		return err
	}

	return nil
}

// handlePickBothHunks is called when the user wants to pick both hunks
// when resolving a merge conflict.
// g and v are passed by the gocui library.
func (gui *Gui) handlePickBothHunks(g *gocui.Gui, v *gocui.View) error {

	conflict := gui.State.Conflicts[gui.State.ConflictIndex]

	err := gui.pushFileSnapshot()
	if err != nil {
		gui.Log.Errorf("Failed to pushFileSnapshot at handlePickBothHunks: %s\n", err)
		return err
	}

	err = gui.resolveConflict(conflict, "both")
	if err != nil {
		gui.Log.Errorf("Failed to resolveConflict at handlePickBothHunks: %s\n", err)
		return err
	}

	err = gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel: %s\n", err)
		return err
	}

	return nil
}

// refreshMergePanel refreshes the mergePanel.
// returns an error if something goes wrong.
func (gui *Gui) refreshMergePanel() error {

	cat, err := gui.catSelectedFile()
	if err != nil {
		gui.Log.Errorf("Failed to catSelectedFile at refreshMergePanel: %s\n", err)
		return err
	}

	if cat == "" {
		return nil
	}

	gui.State.Conflicts, err = gui.findConflicts(cat)
	if err != nil {
		gui.Log.Errorf("Failed to findConflicts at refreshMergePanel: %s\n", err)
		return err
	}

	if len(gui.State.Conflicts) == 0 {

		err = gui.handleCompleteMerge()
		if err != nil {
			gui.Log.Errorf("Failed to handleCompleteMerge at refreshMergePanel: %s\n", err)
			return err
		}

		return nil
	} else if gui.State.ConflictIndex > len(gui.State.Conflicts)-1 {
		gui.State.ConflictIndex = len(gui.State.Conflicts) - 1
	}

	hasFocus := gui.currentViewName(gui.g) == "main"
	if hasFocus {
		err = gui.renderMergeOptions()
		if err != nil {
			gui.Log.Errorf("Failed to renderMergeOptions at refreshMergePanel: %s\n", err)
			return err
		}
	}

	content, err := gui.coloredConflictFile(cat, gui.State.Conflicts, gui.State.ConflictIndex, gui.State.ConflictTop, hasFocus)
	if err != nil {
		gui.Log.Errorf("Failed to get coloredConflictFile at refreshMergePanel: %s\n", err)
		return err
	}

	err = gui.scrollToConflict()
	if err != nil {
		gui.Log.Errorf("Failed ")
		return err
	}

	err = gui.renderString(gui.g, "main", content)
	if err != nil {
		gui.Log.Errorf("Failed to renderString at refreshMergePanel: %s\n", err)
		return err
	}

	return nil
}

// scrollToConflict scrolls to the conflict.
// returns an error if something goes wrong.
func (gui *Gui) scrollToConflict() error {

	mainView, err := gui.g.View("main")
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

	err = mainView.SetOrigin(ox, newOriginY)
	if err != nil {
		gui.Log.Errorf("Failed to SetOrigin at scrollToConflict: %s\n", err)
		return err
	}

	return nil
}

// switchToMerging switches the main view to merging.
// returns an error when something goes wrong.
func (gui *Gui) switchToMerging() error {
	gui.State.ConflictIndex = 0
	gui.State.ConflictTop = true
	_, err := gui.g.SetCurrentView("main")
	if err != nil {
		return err
	}

	err = gui.refreshMergePanel()
	if err != nil {
		gui.Log.Errorf("Failed to refreshMergePanel at switchToMerging: %s\n", err)
		return err
	}

	return nil
}

// renderMergeOptions renders the options.
// Acts as a macro for optionsmap.
// returns an error if something went wrong.
func (gui *Gui) renderMergeOptions() error {
	return gui.renderOptionsMap(gui.g, map[string]string{
		"↑ ↓":   gui.Tr.SLocalize("selectHunk"),
		"← →":   gui.Tr.SLocalize("navigateConflicts"),
		"space": gui.Tr.SLocalize("pickHunk"),
		"b":     gui.Tr.SLocalize("pickBothHunks"),
		"z":     gui.Tr.SLocalize("undo"),
	})
}

// handleEscapeMerge is called when a user presses escape while merging.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleEscapeMerge(g *gocui.Gui, v *gocui.View) error {

	filesView, err := gui.g.View("files")
	if err != nil {
		gui.Log.Errorf("Failed to get files view at handleEscapeMerge: %s\n", err)
		return err
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Errorf("Failed to refresh files at handleEscapeMerge: %s\n", err)
		return err
	}

	err = gui.switchFocus(gui.g, v, filesView)
	if err != nil {
		gui.Log.Errorf("Failed to switchFocus at handleEscapeMerge: %s\n", err)
		return err
	}

	return nil
}

// handleCompleteMerge is called when the user completes the merge.
// returns an error if something goes wrong.
func (gui *Gui) handleCompleteMerge() error {

	v, err := gui.g.View("files")
	if err != nil {
		gui.Log.Errorf("Failed to get files view at handleCompleteMerge: %s\n", err)
		return err
	}

	err = gui.stageSelectedFile()
	if err != nil {
		gui.Log.Errorf("Failed to stageSelectedFile at handleCompleteMerge: %s\n", err)
		return err
	}

	err = gui.refreshFiles()
	if err != nil {
		gui.Log.Errorf("Failed to refreshFiles at handleCompleteMerge:%s\n", err)
		return err
	}

	err = gui.switchFocus(gui.g, nil, v)
	if err != nil {
		gui.Log.Errorf("Failed to switchFocus at handleCompleteMerge: %s\n", err)
		return err
	}

	return nil
}
