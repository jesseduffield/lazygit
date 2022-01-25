// though this panel is called the merge panel, it's really going to use the main panel. This may change in the future

package gui

import (
	"fmt"
	"io/ioutil"
	"math"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
)

func (gui *Gui) handleSelectPrevConflictHunk() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectPrevConflictHunk()
		return gui.renderConflictsWithFocus()
	})
}

func (gui *Gui) handleSelectNextConflictHunk() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectNextConflictHunk()
		return gui.renderConflictsWithFocus()
	})
}

func (gui *Gui) handleSelectNextConflict() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectNextConflict()
		return gui.renderConflictsWithFocus()
	})
}

func (gui *Gui) handleSelectPrevConflict() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectPrevConflict()
		return gui.renderConflictsWithFocus()
	})
}

func (gui *Gui) handleMergeConflictUndo() error {
	state := gui.State.Panels.Merging

	ok := state.Undo()
	if !ok {
		return nil
	}

	gui.logAction("Restoring file to previous state")
	gui.logCommand("Undoing last conflict resolution", false)
	if err := ioutil.WriteFile(state.GetPath(), []byte(state.GetContent()), 0644); err != nil {
		return err
	}

	return gui.renderConflictsWithFocus()
}

func (gui *Gui) handlePickHunk() error {
	return gui.withMergeConflictLock(func() error {
		ok, err := gui.resolveConflict(gui.State.Panels.Merging.Selection())
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		if gui.State.Panels.Merging.AllConflictsResolved() {
			return gui.onLastConflictResolved()
		}

		return gui.renderConflictsWithFocus()
	})
}

func (gui *Gui) handlePickAllHunks() error {
	return gui.withMergeConflictLock(func() error {
		ok, err := gui.resolveConflict(mergeconflicts.ALL)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		if gui.State.Panels.Merging.AllConflictsResolved() {
			return gui.onLastConflictResolved()
		}

		return gui.renderConflictsWithFocus()
	})
}

func (gui *Gui) resolveConflict(selection mergeconflicts.Selection) (bool, error) {
	gui.takeOverMergeConflictScrolling()

	state := gui.State.Panels.Merging

	ok, content, err := state.ContentAfterConflictResolve(selection)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	var logStr string
	switch selection {
	case mergeconflicts.TOP:
		logStr = "Picking top hunk"
	case mergeconflicts.MIDDLE:
		logStr = "Picking middle hunk"
	case mergeconflicts.BOTTOM:
		logStr = "Picking bottom hunk"
	case mergeconflicts.ALL:
		logStr = "Picking all hunks"
	}
	gui.logAction("Resolve merge conflict")
	gui.logCommand(logStr, false)
	state.PushContent(content)
	return true, ioutil.WriteFile(state.GetPath(), []byte(content), 0644)
}

// precondition: we actually have conflicts to render
func (gui *Gui) renderConflicts(hasFocus bool) error {
	state := gui.State.Panels.Merging.State
	content := mergeconflicts.ColoredConflictFile(state, hasFocus)

	if !gui.State.Panels.Merging.UserVerticalScrolling {
		// TODO: find a way to not have to do this OnUIThread thing. Why doesn't it work
		// without it given that we're calling the 'no scroll' variant below?
		gui.OnUIThread(func() error {
			gui.centerYPos(gui.Views.Main, state.GetConflictMiddle())
			return nil
		})
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title:  gui.Tr.MergeConflictsTitle,
			task:   NewRenderStringWithoutScrollTask(content),
			noWrap: true,
		},
	})
}

func (gui *Gui) renderConflictsWithFocus() error {
	return gui.renderConflicts(true)
}

func (gui *Gui) centerYPos(view *gocui.View, y int) {
	ox, _ := view.Origin()
	_, height := view.Size()
	newOriginY := int(math.Max(0, float64(y-(height/2))))
	_ = view.SetOrigin(ox, newOriginY)
}

func (gui *Gui) getMergingOptions() map[string]string {
	keybindingConfig := gui.UserConfig.Keybinding

	return map[string]string{
		fmt.Sprintf("%s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevItem), gui.getKeyDisplay(keybindingConfig.Universal.NextItem)):   gui.Tr.LcSelectHunk,
		fmt.Sprintf("%s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevBlock), gui.getKeyDisplay(keybindingConfig.Universal.NextBlock)): gui.Tr.LcNavigateConflicts,
		gui.getKeyDisplay(keybindingConfig.Universal.Select):   gui.Tr.LcPickHunk,
		gui.getKeyDisplay(keybindingConfig.Main.PickBothHunks): gui.Tr.LcPickAllHunks,
		gui.getKeyDisplay(keybindingConfig.Universal.Undo):     gui.Tr.LcUndo,
	}
}

func (gui *Gui) handleEscapeMerge() error {
	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}

	return gui.escapeMerge()
}

func (gui *Gui) onLastConflictResolved() error {
	// as part of refreshing files, we handle the situation where a file has had
	// its merge conflicts resolved.
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{FILES}})
}

func (gui *Gui) resetMergeState() {
	gui.takeOverMergeConflictScrolling()
	gui.State.Panels.Merging.Reset()
}

func (gui *Gui) setMergeState(path string) (bool, error) {
	content, err := gui.Git.File.Cat(path)
	if err != nil {
		return false, err
	}

	gui.State.Panels.Merging.SetContent(content, path)

	return !gui.State.Panels.Merging.NoConflicts(), nil
}

func (gui *Gui) escapeMerge() error {
	gui.resetMergeState()

	// it's possible this method won't be called from the merging view so we need to
	// ensure we only 'return' focus if we already have it

	if gui.currentContext().GetKey() == MAIN_MERGING_CONTEXT_KEY {
		return gui.pushContext(gui.State.Contexts.Files)
	}
	return nil
}

func (gui *Gui) renderingConflicts() bool {
	currentView := gui.g.CurrentView()
	if currentView != gui.Views.Main && currentView != gui.Views.Files {
		return false
	}

	return gui.State.Panels.Merging.Active()
}

func (gui *Gui) withMergeConflictLock(f func() error) error {
	gui.State.Panels.Merging.Lock()
	defer gui.State.Panels.Merging.Unlock()

	return f()
}

func (gui *Gui) takeOverMergeConflictScrolling() {
	gui.State.Panels.Merging.UserVerticalScrolling = false
}
