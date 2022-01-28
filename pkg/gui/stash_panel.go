package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

// list panel functions

func (gui *Gui) getSelectedStashEntry() *models.StashEntry {
	selectedLine := gui.State.Panels.Stash.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.StashEntries[selectedLine]
}

func (gui *Gui) stashRenderToMain() error {
	var task updateTask
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		task = NewRenderStringTask(gui.Tr.NoStashEntries)
	} else {
		task = NewRunPtyTask(gui.Git.Stash.ShowStashEntryCmdObj(stashEntry.Index).GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Stash",
			task:  task,
		},
	})
}

func (gui *Gui) refreshStashEntries() error {
	gui.State.StashEntries = gui.Git.Loaders.Stash.
		GetStashEntries(gui.State.Modes.Filtering.GetPath())

	return gui.postRefreshUpdate(gui.State.Contexts.Stash)
}

// specific functions

func (gui *Gui) handleStashApply() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	skipStashWarning := gui.UserConfig.Gui.SkipStashWarning

	apply := func() error {
		gui.logAction(gui.Tr.Actions.Stash)
		err := gui.Git.Stash.Apply(stashEntry.Index)
		_ = gui.postStashRefresh()
		if err != nil {
			return gui.surfaceError(err)
		}
		return nil
	}

	if skipStashWarning {
		return apply()
	}

	return gui.ask(askOpts{
		title:  gui.Tr.StashApply,
		prompt: gui.Tr.SureApplyStashEntry,
		handleConfirm: func() error {
			return apply()
		},
	})
}

func (gui *Gui) handleStashPop() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	skipStashWarning := gui.UserConfig.Gui.SkipStashWarning

	pop := func() error {
		gui.logAction(gui.Tr.Actions.Stash)
		err := gui.Git.Stash.Pop(stashEntry.Index)
		_ = gui.postStashRefresh()
		if err != nil {
			return gui.surfaceError(err)
		}
		return nil
	}

	if skipStashWarning {
		return pop()
	}

	return gui.ask(askOpts{
		title:  gui.Tr.StashPop,
		prompt: gui.Tr.SurePopStashEntry,
		handleConfirm: func() error {
			return pop()
		},
	})
}

func (gui *Gui) handleStashDrop() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	return gui.ask(askOpts{
		title:  gui.Tr.StashDrop,
		prompt: gui.Tr.SureDropStashEntry,
		handleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.Stash)
			err := gui.Git.Stash.Drop(stashEntry.Index)
			_ = gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{STASH}})
			if err != nil {
				return gui.surfaceError(err)
			}
			return nil
		},
	})
}

func (gui *Gui) postStashRefresh() error {
	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{STASH, FILES}})
}

func (gui *Gui) handleStashSave(stashFunc func(message string) error) error {
	if len(gui.trackedFiles()) == 0 && len(gui.stagedFiles()) == 0 {
		return gui.createErrorPanel(gui.Tr.NoTrackedStagedFilesStash)
	}

	return gui.prompt(promptOpts{
		title: gui.Tr.StashChanges,
		handleConfirm: func(stashComment string) error {
			err := stashFunc(stashComment)
			_ = gui.postStashRefresh()
			if err != nil {
				return gui.surfaceError(err)
			}
			return nil
		},
	})
}

func (gui *Gui) handleViewStashFiles() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	return gui.switchToCommitFilesContext(stashEntry.RefName(), false, gui.State.Contexts.Stash, "stash")
}
