// though this panel is called the merge panel, it's really going to use the main panel. This may change in the future

package gui

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/go-errors/errors"
	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
)

func (gui *Gui) handleSelectTop() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverScrolling()
		gui.State.Panels.Merging.ConflictTop = true
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handleSelectBottom() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverScrolling()
		gui.State.Panels.Merging.ConflictTop = false
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handleSelectNextConflict() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverScrolling()
		if gui.State.Panels.Merging.ConflictIndex >= len(gui.State.Panels.Merging.Conflicts)-1 {
			return nil
		}
		gui.State.Panels.Merging.ConflictIndex++
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) handleSelectPrevConflict() error {
	return gui.withMergeConflictLock(func() error {
		gui.takeOverScrolling()
		if gui.State.Panels.Merging.ConflictIndex <= 0 {
			return nil
		}
		gui.State.Panels.Merging.ConflictIndex--
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) pushFileSnapshot() error {
	gitFile := gui.getSelectedFile()
	if gitFile == nil {
		return nil
	}
	content, err := gui.GitCommand.CatFile(gitFile.Name)
	if err != nil {
		return err
	}
	gui.State.Panels.Merging.EditHistory.Push(content)
	return nil
}

func (gui *Gui) handlePopFileSnapshot() error {
	if gui.State.Panels.Merging.EditHistory.Len() == 0 {
		return nil
	}
	prevContent := gui.State.Panels.Merging.EditHistory.Pop().(string)
	gitFile := gui.getSelectedFile()
	if gitFile == nil {
		return nil
	}
	if err := ioutil.WriteFile(gitFile.Name, []byte(prevContent), 0644); err != nil {
		return err
	}

	return gui.refreshMergePanel()
}

func (gui *Gui) handlePickHunk() error {
	return gui.withMergeConflictLock(func() error {
		conflict := gui.getCurrentConflict()
		if conflict == nil {
			return nil
		}

		gui.takeOverScrolling()

		if err := gui.pushFileSnapshot(); err != nil {
			return err
		}

		selection := mergeconflicts.BOTTOM
		if gui.State.Panels.Merging.ConflictTop {
			selection = mergeconflicts.TOP
		}
		err := gui.resolveConflict(*conflict, selection)
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
	})
}

func (gui *Gui) handlePickBothHunks() error {
	return gui.withMergeConflictLock(func() error {
		conflict := gui.getCurrentConflict()
		if conflict == nil {
			return nil
		}

		gui.takeOverScrolling()

		if err := gui.pushFileSnapshot(); err != nil {
			return err
		}
		err := gui.resolveConflict(*conflict, mergeconflicts.BOTH)
		if err != nil {
			panic(err)
		}
		return gui.refreshMergePanel()
	})
}

func (gui *Gui) getCurrentConflict() *commands.Conflict {
	if len(gui.State.Panels.Merging.Conflicts) == 0 {
		return nil
	}

	return &gui.State.Panels.Merging.Conflicts[gui.State.Panels.Merging.ConflictIndex]
}

func (gui *Gui) resolveConflict(conflict commands.Conflict, selection mergeconflicts.Selection) error {
	gitFile := gui.getSelectedFile()
	if gitFile == nil {
		return nil
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
		if !mergeconflicts.IsIndexToDelete(i, conflict, selection) {
			output += line
		}
	}
	return ioutil.WriteFile(gitFile.Name, []byte(output), 0644)
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
				task:  gui.createRenderStringTask(err.Error()),
			},
		})
	}

	panelState.Conflicts = mergeconflicts.FindConflicts(cat)

	// handle potential fixes that the user made in their editor since we last refreshed
	if len(panelState.Conflicts) == 0 {
		return gui.handleCompleteMerge()
	} else if panelState.ConflictIndex > len(panelState.Conflicts)-1 {
		panelState.ConflictIndex = len(panelState.Conflicts) - 1
	}

	hasFocus := gui.currentViewName() == "main"
	content := mergeconflicts.ColoredConflictFile(cat, panelState.Conflicts, panelState.ConflictIndex, panelState.ConflictTop, hasFocus)

	if err := gui.scrollToConflict(); err != nil {
		return err
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title:  gui.Tr.MergeConflictsTitle,
			task:   gui.createRenderStringWithoutScrollTask(content),
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

	cat, err := gui.GitCommand.CatFile(item.Name)
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
	gui.takeOverScrolling()

	gui.State.Panels.Merging.EditHistory = stack.New()
	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}
	// it's possible this method won't be called from the merging view so we need to
	// ensure we only 'return' focus if we already have it
	if gui.g.CurrentView() == gui.getMainView() {
		return gui.pushContext(gui.Contexts.Files.Context)
	}
	return nil
}

func (gui *Gui) handleCompleteMerge() error {
	if err := gui.stageSelectedFile(); err != nil {
		return err
	}
	if err := gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{FILES}}); err != nil {
		return err
	}
	// if we got conflicts after unstashing, we don't want to call any git
	// commands to continue rebasing/merging here
	if gui.GitCommand.WorkingTreeState() == commands.REBASE_MODE_NORMAL {
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
	gui.takeOverScrolling()

	return gui.ask(askOpts{
		title:               "continue",
		prompt:              gui.Tr.ConflictsResolved,
		handlersManageFocus: true,
		handleConfirm: func() error {
			if err := gui.pushContext(gui.Contexts.Files.Context); err != nil {
				return err
			}

			return gui.genericMergeCommand("continue")
		},
		handleClose: func() error {
			return gui.pushContext(gui.Contexts.Files.Context)
		},
	})
}

func (gui *Gui) canScrollMergePanel() bool {
	currentViewName := gui.currentViewName()
	if currentViewName != "main" {
		return false
	}

	file := gui.getSelectedFile()
	if file == nil {
		return false
	}

	return file.HasInlineMergeConflicts
}

func (gui *Gui) withMergeConflictLock(f func() error) error {
	gui.State.Panels.Merging.ConflictsMutex.Lock()
	defer gui.State.Panels.Merging.ConflictsMutex.Unlock()

	return f()
}

func (gui *Gui) takeOverScrolling() {
	gui.State.Panels.Merging.UserScrolling = false
}
