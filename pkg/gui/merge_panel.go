// though this panel is called the merge panel, it's really going to use the main panel. This may change in the future

package gui

import (
	"fmt"
	"io/ioutil"
	"math"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) handleSelectTop() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectTopOption()
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handleSelectBottom() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()
		gui.State.Panels.Merging.SelectBottomOption()
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
	gui.OnRunCommand(oscommands.NewCmdLogEntry("Undoing last conflict resolution", "Undo merge conflict resolution", false))
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
		}
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handlePickBothHunks() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverMergeConflictScrolling()

		ok, err := gui.resolveConflict(mergeconflicts.BOTH)
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
		return false, gui.SurfaceError(err)
	}

	var logStr string
	switch selection {
	case mergeconflicts.TOP:
		logStr = "Picking top hunk"
	case mergeconflicts.BOTTOM:
		logStr = "Picking bottom hunk"
	case mergeconflicts.BOTH:
		logStr = "Picking both hunks"
	}
	gui.OnRunCommand(oscommands.NewCmdLogEntry(logStr, "Resolve merge conflict", false))
	return true, ioutil.WriteFile(gitFile.Name, []byte(output), 0644)
}

func (gui *Gui) refreshMergePanelWithLock() error {
	return gui.withMergeConflictLock(gui.refreshMergePanel)
}

func (gui *Gui) refreshMergePanel() error {
	panelState := gui.State.Panels.Merging
	cat, err := gui.catSelectedFile()
	if err != nil {
		return gui.refreshMainViews(refreshMainOpts{
			main: &viewUpdateOpts{
				title: "",
				task:  NewRenderStringTask(err.Error()),
			},
		})
	}

	panelState.SetConflictsFromCat(cat)

	if panelState.NoConflicts() {
		return gui.handleCompleteMerge()
	}

	hasFocus := gui.currentViewName() == "main"
	content := mergeconflicts.ColoredConflictFile(cat, panelState.State, hasFocus)

	if err := gui.scrollToConflict(); err != nil {
		return err
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title:  gui.Tr.MergeConflictsTitle,
			task:   NewRenderStringWithoutScrollTask(content),
			noWrap: true,
		},
	})
}

func (gui *Gui) catSelectedFile() (string, error) {
	item := gui.getSelectedFile()
	if item == nil {
		return "", errors.New(gui.Tr.NoFilesDisplay)
	}

	if item.Type != "file" {
		return "", errors.New(gui.Tr.NotAFile)
	}

	cat, err := gui.GetOS().CatFile(item.Name)
	if err != nil {
		gui.Log.Error(err)
		return "", err
	}
	return cat, nil
}

func (gui *Gui) scrollToConflict() error {
	if gui.State.Panels.Merging.UserScrolling {
		return nil
	}

	panelState := gui.State.Panels.Merging
	if panelState.NoConflicts() {
		return nil
	}

	gui.centerYPos(gui.Views.Main, panelState.GetConflictMiddle())

	return nil
}

func (gui *Gui) centerYPos(view *gocui.View, y int) {
	ox, _ := view.Origin()
	_, height := view.Size()
	newOriginY := int(math.Max(0, float64(y-(height/2))))
	gui.g.Update(func(g *gocui.Gui) error {
		return view.SetOrigin(ox, newOriginY)
	})
}

func (gui *Gui) getMergingOptions() map[string]string {
	keybindingConfig := gui.Config.GetUserConfig().Keybinding

	return map[string]string{
		fmt.Sprintf("%s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevItem), gui.getKeyDisplay(keybindingConfig.Universal.NextItem)):   gui.Tr.LcSelectHunk,
		fmt.Sprintf("%s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevBlock), gui.getKeyDisplay(keybindingConfig.Universal.NextBlock)): gui.Tr.LcNavigateConflicts,
		gui.getKeyDisplay(keybindingConfig.Universal.Select):   gui.Tr.LcPickHunk,
		gui.getKeyDisplay(keybindingConfig.Main.PickBothHunks): gui.Tr.LcPickBothHunks,
		gui.getKeyDisplay(keybindingConfig.Universal.Undo):     gui.Tr.LcUndo,
	}
}

func (gui *Gui) handleEscapeMerge() error {
	gui.takeOverMergeConflictScrolling()

	gui.State.Panels.Merging.Reset()
	if err := gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{FILES}}); err != nil {
		return err
	}
	// it's possible this method won't be called from the merging view so we need to
	// ensure we only 'return' focus if we already have it
	if gui.g.CurrentView() == gui.Views.Main {
		return gui.pushContext(gui.State.Contexts.Files)
	}
	return nil
}

func (gui *Gui) handleCompleteMerge() error {
	if err := gui.stageSelectedFile(); err != nil {
		return err
	}

	if err := gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{FILES}}); err != nil {
		return err
	}

	// if we got conflicts after unstashing, we don't want to call any git
	// commands to continue rebasing/merging here
	if gui.Git.Status().InNormalWorkingTreeState() {
		return gui.handleEscapeMerge()
	}

	// if there are no more files with merge conflicts, we should ask whether the user wants to continue
	if !gui.anyFilesWithMergeConflicts() {
		return gui.promptToContinueRebase()
	}

	return gui.handleEscapeMerge()
}

// promptToContinueRebase asks the user if they want to continue the rebase/merge that's in progress
func (gui *Gui) promptToContinueRebase() error {
	gui.takeOverMergeConflictScrolling()

	return gui.Ask(AskOpts{
		Title:               "continue",
		Prompt:              gui.Tr.ConflictsResolved,
		HandlersManageFocus: true,
		HandleConfirm: func() error {
			if err := gui.pushContext(gui.State.Contexts.Files); err != nil {
				return err
			}

			return gui.genericMergeCommand("continue")
		},
		HandleClose: func() error {
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
	gui.State.Panels.Merging.UserScrolling = false
}
