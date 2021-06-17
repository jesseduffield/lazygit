package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
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
		task = NewRunPtyTask(
			gui.Git.Stash().ShowEntryCmdObj(stashEntry.Index),
		)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Stash",
			task:  task,
		},
	})
}

func (gui *Gui) refreshStashEntries() error {
	gui.State.StashEntries = gui.Git.Stash().LoadEntries(gui.State.Modes.Filtering.GetPath())

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

	return gui.Ask(AskOpts{
		Title:  gui.Tr.StashApply,
		Prompt: gui.Tr.SureApplyStashEntry,
		HandleConfirm: func() error {
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

	return gui.Ask(AskOpts{
		Title:  gui.Tr.StashPop,
		Prompt: gui.Tr.SurePopStashEntry,
		HandleConfirm: func() error {
			return pop()
		},
	})
}

func (gui *Gui) handleStashDrop() error {
	return gui.Ask(AskOpts{
		Title:  gui.Tr.StashDrop,
		Prompt: gui.Tr.SureDropStashEntry,
		HandleConfirm: func() error {
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

		return gui.CreateErrorPanel(errorMessage)
	}
	if err := gui.Git.WithSpan(gui.Tr.Spans.Stash).Stash().Do(stashEntry.Index, method); err != nil {
		return gui.SurfaceError(err)
	}
	return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{STASH, FILES}})
}

func (gui *Gui) handleStashSave(stashFunc func(message string) error) error {
	if len(gui.trackedFiles()) == 0 && len(gui.stagedFiles()) == 0 {
		return gui.CreateErrorPanel(gui.Tr.NoTrackedStagedFilesStash)
	}

	return gui.Prompt(PromptOpts{
		Title: gui.Tr.StashChanges,
		HandleConfirm: func(stashComment string) error {
			if err := stashFunc(stashComment); err != nil {
				return gui.SurfaceError(err)
			}
			return gui.RefreshSidePanels(RefreshOptions{Scope: []RefreshableView{STASH, FILES}})
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
