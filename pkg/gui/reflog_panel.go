package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

// list panel functions

func (gui *Gui) getSelectedReflogCommit() *models.Commit {
	selectedLine := gui.State.Panels.ReflogCommits.SelectedLineIdx
	reflogComits := gui.State.FilteredReflogCommits
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
		cmdObj := gui.GitCommand.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())

		task = NewRunPtyTask(cmdObj.GetCmd())
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Reflog Entry",
			task:  task,
		},
	})
}

// the reflogs panel is the only panel where we cache data, in that we only
// load entries that have been created since we last ran the call. This means
// we need to be more careful with how we use this, and to ensure we're emptying
// the reflogs array when changing contexts.
// This method also manages two things: ReflogCommits and FilteredReflogCommits.
// FilteredReflogCommits are rendered in the reflogs panel, and ReflogCommits
// are used by the branches panel to obtain recency values for sorting.
func (gui *Gui) refreshReflogCommits() error {
	// pulling state into its own variable incase it gets swapped out for another state
	// and we get an out of bounds exception
	state := gui.State
	var lastReflogCommit *models.Commit
	if len(state.ReflogCommits) > 0 {
		lastReflogCommit = state.ReflogCommits[0]
	}

	refresh := func(stateCommits *[]*models.Commit, filterPath string) error {
		commits, onlyObtainedNewReflogCommits, err := gui.GitCommand.GetReflogCommits(lastReflogCommit, filterPath)
		if err != nil {
			return gui.surfaceError(err)
		}

		if onlyObtainedNewReflogCommits {
			*stateCommits = append(commits, *stateCommits...)
		} else {
			*stateCommits = commits
		}
		return nil
	}

	if err := refresh(&state.ReflogCommits, ""); err != nil {
		return err
	}

	if gui.State.Modes.Filtering.Active() {
		if err := refresh(&state.FilteredReflogCommits, state.Modes.Filtering.GetPath()); err != nil {
			return err
		}
	} else {
		state.FilteredReflogCommits = state.ReflogCommits
	}

	return gui.postRefreshUpdate(gui.State.Contexts.ReflogCommits)
}

func (gui *Gui) handleCheckoutReflogCommit() error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	err := gui.ask(askOpts{
		title:  gui.Tr.LcCheckoutCommit,
		prompt: gui.Tr.SureCheckoutThisCommit,
		handleConfirm: func() error {
			return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{span: gui.Tr.Spans.CheckoutReflogCommit})
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

	return gui.createResetMenu(commit.Sha)
}

func (gui *Gui) handleViewReflogCommitFiles() error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	return gui.switchToCommitFilesContext(commit.Sha, false, gui.State.Contexts.ReflogCommits, "commits")
}
