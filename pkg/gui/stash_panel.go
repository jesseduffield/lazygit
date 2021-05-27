package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedStashEntry() *models.StashEntry {
	selectedLine := gui.State.Panels.Stash.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.StashEntries[selectedLine]
}

func (gui *Gui) handleStashEntrySelect() error {
	var task updateTask
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		task = NewRenderStringTask(gui.Tr.NoStashEntries)
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.ShowStashEntryCmdStr(stashEntry.Index),
		)
		task = NewRunPtyTask(cmd)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Stash",
			task:  task,
		},
	})
}

func (gui *Gui) refreshStashEntries() error {
	gui.State.StashEntries = gui.GitCommand.GetStashEntries(gui.State.Modes.Filtering.GetPath())

	return gui.State.Contexts.Stash.HandleRender()
}

// specific functions

func (gui *Gui) handleStashApply() error {
	skipStashWarning := gui.Config.GetUserConfig().Gui.SkipStashWarning

	apply := func() error {
		return gui.stashDo("apply")
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
	skipStashWarning := gui.Config.GetUserConfig().Gui.SkipStashWarning

	pop := func() error {
		return gui.stashDo("pop")
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
	return gui.ask(askOpts{
		title:  gui.Tr.StashDrop,
		prompt: gui.Tr.SureDropStashEntry,
		handleConfirm: func() error {
			return gui.stashDo("drop")
		},
	})
}

func (gui *Gui) stashDo(method string) error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		errorMessage := utils.ResolvePlaceholderString(
			gui.Tr.NoStashTo,
			map[string]string{
				"method": method,
			},
		)

		return gui.createErrorPanel(errorMessage)
	}
	if err := gui.GitCommand.WithSpan(gui.Tr.Spans.Stash).StashDo(stashEntry.Index, method); err != nil {
		return gui.surfaceError(err)
	}
	return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{STASH, FILES}})
}

func (gui *Gui) handleStashSave(stashFunc func(message string) error) error {
	if len(gui.trackedFiles()) == 0 && len(gui.stagedFiles()) == 0 {
		return gui.createErrorPanel(gui.Tr.NoTrackedStagedFilesStash)
	}

	return gui.prompt(promptOpts{
		title: gui.Tr.StashChanges,
		handleConfirm: func(stashComment string) error {
			if err := stashFunc(stashComment); err != nil {
				return gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{scope: []RefreshableView{STASH, FILES}})
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
