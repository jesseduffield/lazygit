// though this panel is called the merge panel, it's really going to use the main panel. This may change in the future

package gui

import (
	"fmt"
	"io/ioutil"
	"math"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
			gui.State.Panels.Merging.Lock()
			defer gui.State.Panels.Merging.Unlock()

			if !state.Active() {
				return nil
			}

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

func (gui *Gui) renderConflictsWithLock(hasFocus bool) error {
	return gui.withMergeConflictLock(func() error {
		return gui.renderConflicts(hasFocus)
	})
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
	if err := gui.refreshSidePanels(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}}); err != nil {
		return err
	}

	return gui.escapeMerge()
}

func (gui *Gui) onLastConflictResolved() error {
	// as part of refreshing files, we handle the situation where a file has had
	// its merge conflicts resolved.
	return gui.refreshSidePanels(types.RefreshOptions{mode: types.ASYNC, scope: []types.RefreshableView{types.FILES}})
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

func (gui *Gui) setMergeStateWithLock(path string) (bool, error) {
	gui.State.Panels.Merging.Lock()
	defer gui.State.Panels.Merging.Unlock()

	return gui.setMergeState(path)
}

func (gui *Gui) resetMergeStateWithLock() {
	gui.State.Panels.Merging.Lock()
	defer gui.State.Panels.Merging.Unlock()

	gui.resetMergeState()
}

func (gui *Gui) escapeMerge() error {
	gui.resetMergeState()

	// doing this in separate UI thread so that we're not still holding the lock by the time refresh the file
	gui.OnUIThread(func() error {
		return gui.pushContext(gui.State.Contexts.Files)
	})
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

func (gui *Gui) setConflictsAndRender(path string, hasFocus bool) (bool, error) {
	hasConflicts, err := gui.setMergeState(path)
	if err != nil {
		return false, err
	}

	// if we don't have conflicts we'll fall through and show the diff
	if hasConflicts {
		return true, gui.renderConflicts(hasFocus)
	}

	return false, nil
}

func (gui *Gui) setConflictsAndRenderWithLock(path string, hasFocus bool) (bool, error) {
	gui.State.Panels.Merging.Lock()
	defer gui.State.Panels.Merging.Unlock()

	return gui.setConflictsAndRender(path, hasFocus)
}

func (gui *Gui) refreshMergeState() error {
	gui.State.Panels.Merging.Lock()
	defer gui.State.Panels.Merging.Unlock()

	if gui.currentContext().GetKey() != MAIN_MERGING_CONTEXT_KEY {
		return nil
	}

	hasConflicts, err := gui.setConflictsAndRender(gui.State.Panels.Merging.GetPath(), true)
	if err != nil {
		return gui.surfaceError(err)
	}

	if !hasConflicts {
		return gui.escapeMerge()
	}

	return nil
}
