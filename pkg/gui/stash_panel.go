package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// list panel functions

func (gui *Gui) getSelectedStashEntry() *models.StashEntry {
	selectedLine := gui.State.Panels.Stash.SelectedLineIdx
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Model.StashEntries[selectedLine]
}

func (gui *Gui) stashRenderToMain() error {
	var task updateTask
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		task = NewRenderStringTask(gui.c.Tr.NoStashEntries)
	} else {
		task = NewRunPtyTask(gui.git.Stash.ShowStashEntryCmdObj(stashEntry.Index).GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Stash",
			task:  task,
		},
	})
}

// specific functions

func (gui *Gui) handleStashApply() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	skipStashWarning := gui.c.UserConfig.Gui.SkipStashWarning

	apply := func() error {
		gui.c.LogAction(gui.c.Tr.Actions.Stash)
		err := gui.git.Stash.Apply(stashEntry.Index)
		_ = gui.postStashRefresh()
		if err != nil {
			return gui.c.Error(err)
		}
		return nil
	}

	if skipStashWarning {
		return apply()
	}

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.StashApply,
		Prompt: gui.c.Tr.SureApplyStashEntry,
		HandleConfirm: func() error {
			return apply()
		},
	})
}

func (gui *Gui) handleStashPop() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	skipStashWarning := gui.c.UserConfig.Gui.SkipStashWarning

	pop := func() error {
		gui.c.LogAction(gui.c.Tr.Actions.Stash)
		err := gui.git.Stash.Pop(stashEntry.Index)
		_ = gui.postStashRefresh()
		if err != nil {
			return gui.c.Error(err)
		}
		return nil
	}

	if skipStashWarning {
		return pop()
	}

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.StashPop,
		Prompt: gui.c.Tr.SurePopStashEntry,
		HandleConfirm: func() error {
			return pop()
		},
	})
}

func (gui *Gui) handleStashDrop() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	return gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.StashDrop,
		Prompt: gui.c.Tr.SureDropStashEntry,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.Stash)
			err := gui.git.Stash.Drop(stashEntry.Index)
			_ = gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH}})
			if err != nil {
				return gui.c.Error(err)
			}
			return nil
		},
	})
}

func (gui *Gui) postStashRefresh() error {
	return gui.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
}

func (gui *Gui) handleViewStashFiles() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	return gui.SwitchToCommitFilesContext(controllers.SwitchToCommitFilesContextOpts{
		RefName:    stashEntry.RefName(),
		CanRebase:  false,
		Context:    gui.State.Contexts.Stash,
		WindowName: "stash",
	})
}

func (gui *Gui) handleNewBranchOffStashEntry() error {
	stashEntry := gui.getSelectedStashEntry()
	if stashEntry == nil {
		return nil
	}

	return gui.helpers.Refs.NewBranch(stashEntry.RefName(), stashEntry.Description(), "")
}
