package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// list panel functions

func (gui *Gui) reflogCommitsRenderToMain() error {
	commit := gui.State.Contexts.ReflogCommits.GetSelected()
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
	commit := gui.State.Contexts.ReflogCommits.GetSelected()
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

	return nil
}

func (gui *Gui) handleCreateReflogResetMenu() error {
	commit := gui.State.Contexts.ReflogCommits.GetSelected()

	return gui.helpers.Refs.CreateGitResetMenu(commit.Sha)
}

func (gui *Gui) handleViewReflogCommitFiles() error {
	commit := gui.State.Contexts.ReflogCommits.GetSelected()
	if commit == nil {
		return nil
	}

	return gui.SwitchToCommitFilesContext(controllers.SwitchToCommitFilesContextOpts{
		RefName:   commit.Sha,
		CanRebase: false,
		Context:   gui.State.Contexts.ReflogCommits,
	})
}

func (gui *Gui) handleCopyReflogCommit() error {
	commit := gui.State.Contexts.ReflogCommits.GetSelected()
	if commit == nil {
		return nil
	}

	return gui.helpers.CherryPick.Copy(commit, gui.State.Model.FilteredReflogCommits, gui.State.Contexts.ReflogCommits)
}

func (gui *Gui) handleCopyReflogCommitRange() error {
	// just doing this to ensure something is selected
	commit := gui.State.Contexts.ReflogCommits.GetSelected()
	if commit == nil {
		return nil
	}

	return gui.helpers.CherryPick.CopyRange(gui.State.Contexts.ReflogCommits.GetSelectedLineIdx(), gui.State.Model.FilteredReflogCommits, gui.State.Contexts.ReflogCommits)
}
