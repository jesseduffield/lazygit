package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// list panel functions

func (gui *Gui) getSelectedReflogCommit() *models.Commit {
	selectedLine := gui.State.Panels.ReflogCommits.SelectedLineIdx
	reflogComits := gui.State.Model.FilteredReflogCommits
	if selectedLine == -1 || len(reflogComits) == 0 {
		return nil
	}

	return reflogComits[selectedLine]
}

func (gui *Gui) reflogCommitsRenderToMain() error {
	commit := gui.getSelectedReflogCommit()
	var task updateTask
	if commit == nil {
		task = NewRenderStringTask("No reflog history")
	} else {
		cmdObj := gui.git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())

		task = NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Reflog Entry",
			task:  task,
		},
	})
}

func (gui *Gui) CheckoutReflogCommit() error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	err := gui.c.Ask(types.AskOpts{
		Title:  gui.c.Tr.LcCheckoutCommit,
		Prompt: gui.c.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.CheckoutReflogCommit)
			return gui.helpers.Refs.CheckoutRef(commit.Sha, types.CheckoutRefOptions{})
		},
	})
	if err != nil {
		return err
	}

	gui.State.Panels.ReflogCommits.SelectedLineIdx = 0

	return nil
}

func (gui *Gui) handleCreateReflogResetMenu() error {
	commit := gui.getSelectedReflogCommit()

	return gui.helpers.Refs.CreateGitResetMenu(commit.Sha)
}

func (gui *Gui) handleViewReflogCommitFiles() error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	return gui.SwitchToCommitFilesContext(controllers.SwitchToCommitFilesContextOpts{
		RefName:    commit.Sha,
		CanRebase:  false,
		Context:    gui.State.Contexts.ReflogCommits,
		WindowName: "commits",
	})
}

func (gui *Gui) handleCopyReflogCommit() error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	return gui.helpers.CherryPick.Copy(commit, gui.State.Model.FilteredReflogCommits, gui.State.Contexts.ReflogCommits)
}

func (gui *Gui) handleCopyReflogCommitRange() error {
	// just doing this to ensure something is selected
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	return gui.helpers.CherryPick.CopyRange(gui.State.Contexts.ReflogCommits.GetPanelState().GetSelectedLineIdx(), gui.State.Model.FilteredReflogCommits, gui.State.Contexts.ReflogCommits)
}
