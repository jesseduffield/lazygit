// though this panel is called the merge panel, it's really going to use the main panel. This may change in the future

package gui

import (
	"fmt"
	"io/ioutil"
	"math"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
)

func (gui *Gui) handleSelectPrevConflictHunk() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectPrevConflictHunk()
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handleSelectNextConflictHunk() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectNextConflictHunk()
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handleSelectNextConflict() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectNextConflict()
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handleSelectPrevConflict() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectPrevConflict()
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) pushFileSnapshot() error {
	content, err := gui.catSelectedFile()
	if err != nil {
		return err
	}
	gui.State.Panels.Merging.PushFileSnapshot(content)
	return nil
}

func (gui *Gui) handlePopFileSnapshot() error {
	prevContent, ok := gui.State.Panels.Merging.PopFileSnapshot()
	if !ok {
		return nil
	}

	gitFile := gui.getSelectedFile()
	if gitFile == nil {
		return nil
	}
	gui.logAction("Restoring file to previous state")
	gui.logCommand("Undoing last conflict resolution", false)
	if err := ioutil.WriteFile(gitFile.Name, []byte(prevContent), 0644); err != nil {
		return err
	}

	return gui.refreshMergePanel()
}

func (gui *Gui) handlePickHunk() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()

		ok, err := gui.resolveConflict(gui.State.Panels.Merging.Selection())
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		if gui.State.Panels.Merging.IsFinalConflict() {
			if err := gui.handleCompleteMerge(); err != nil {
				return err
			}
			return nil
		}
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handlePickAllHunks() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()

		ok, err := gui.resolveConflict(mergeconflicts.ALL)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		return gui.refreshMergePanel()
	})
}

func (gui *Gui) resolveConflict(selection mergeconflicts.Selection) (bool, error) {
	gitFile := gui.getSelectedFile()
	if gitFile == nil {
		return false, nil
	}

	ok, output, err := gui.State.Panels.Merging.ContentAfterConflictResolve(gitFile.Name, selection)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	if err := gui.pushFileSnapshot(); err != nil {
		return false, gui.PopupHandler.Error(err)
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
	return true, ioutil.WriteFile(gitFile.Name, []byte(output), 0644)
}

func (gui *Gui) refreshMergePanelWithLock() error {
	return gui.withMergeConflictLock(gui.refreshMergePanel)
}

// not re-using state here because we can run into issues with mutexes when
// doing that.
func (gui *Gui) renderConflictsFromFilesPanel() error {
	state := mergeconflicts.NewState()
	_, err := gui.renderConflicts(state, false)

	return err
}

func (gui *Gui) renderConflicts(state *mergeconflicts.State, hasFocus bool) (bool, error) {
	cat, err := gui.catSelectedFile()
	if err != nil {
		return false, gui.refreshMainViews(refreshMainOpts{
			main: &viewUpdateOpts{
				title: "",
				task:  NewRenderStringTask(err.Error()),
			},
		})
	}

	state.SetConflictsFromCat(cat)

	if state.NoConflicts() {
		return false, gui.handleCompleteMerge()
	}

	content := mergeconflicts.ColoredConflictFile(cat, state, hasFocus)

	if !gui.State.Panels.Merging.UserVerticalScrolling {
		// TODO: find a way to not have to do this OnUIThread thing. Why doesn't it work
		// without it given that we're calling the 'no scroll' variant below?
		gui.OnUIThread(func() error {
			gui.centerYPos(gui.Views.Main, state.GetConflictMiddle())
			return nil
		})
	}

	return true, gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title:  gui.Tr.MergeConflictsTitle,
			task:   NewRenderStringWithoutScrollTask(content),
			noWrap: true,
		},
	})
}

func (gui *Gui) refreshMergePanel() error {
	conflictsFound, err := gui.renderConflicts(gui.State.Panels.Merging.State, true)
	if err != nil {
		return err
	}

	if !conflictsFound {
		return gui.handleCompleteMerge()
	}

	return nil
}

func (gui *Gui) catSelectedFile() (string, error) {
	item := gui.getSelectedFile()
	if item == nil {
		return "", errors.New(gui.Tr.NoFilesDisplay)
	}

	if item.Type != "file" {
		return "", errors.New(gui.Tr.NotAFile)
	}

	cat, err := gui.Git.File.Cat(item.Name)
	if err != nil {
		gui.Log.Error(err)
		return "", err
	}
	return cat, nil
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

func (gui *Gui) handleCompleteMerge() error {
	if err := gui.stageSelectedFile(); err != nil {
		return err
	}
	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}

	// if there are no more files with merge conflicts, we should ask whether the user wants to continue
	if gui.Git.Status.WorkingTreeState() != enums.REBASE_MODE_NONE && !gui.anyFilesWithMergeConflicts() {
		return gui.promptToContinueRebase()
	}

	return gui.escapeMerge()
}

func (gui *Gui) escapeMerge() error {
	gui.takeOverMergeConflictScrolling()

	gui.State.Panels.Merging.Reset()

	// it's possible this method won't be called from the merging view so we need to
	// ensure we only 'return' focus if we already have it
	if gui.g.CurrentView() == gui.Views.Main {
		return gui.pushContext(gui.State.Contexts.Files)
	}
	return nil
}

// promptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (gui *Gui) promptToContinueRebase() error {
	gui.takeOverMergeConflictScrolling()

	return gui.PopupHandler.Ask(askOpts{
		title:               "continue",
		prompt:              gui.Tr.ConflictsResolved,
		handlersManageFocus: true,
		handleConfirm: func() error {
			if err := gui.pushContext(gui.State.Contexts.Files); err != nil {
				return err
			}

			return gui.genericMergeCommand(REBASE_OPTION_CONTINUE)
		},
		handleClose: func() error {
			return gui.pushContext(gui.State.Contexts.Files)
		},
	})
}

func (gui *Gui) canScrollMergePanel() bool {
	currentView := gui.g.CurrentView()
	if currentView != gui.Views.Main && currentView != gui.Views.Files {
		return false
	}

	file := gui.getSelectedFile()
	if file == nil {
		return false
	}

	return file.HasInlineMergeConflicts
}

func (gui *Gui) withMergeConflictLock(f func() error) error {
	gui.State.Panels.Merging.Lock()
	defer gui.State.Panels.Merging.Unlock()

	return f()
}

func (gui *Gui) takeOverMergeConflictScrolling() {
	gui.State.Panels.Merging.UserVerticalScrolling = false
}
