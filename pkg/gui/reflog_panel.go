package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// list panel functions

func (gui *Gui) getSelectedReflogCommit() *commands.Commit {
	selectedLine := gui.State.Panels.ReflogCommits.SelectedLineIdx
	reflogComits := gui.State.FilteredReflogCommits
	if selectedLine == -1 || len(reflogComits) == 0 {
		return nil
	}

	return reflogComits[selectedLine]
}

func (gui *Gui) handleReflogCommitSelect() error {
	commit := gui.getSelectedReflogCommit()
	var task updateTask
	if commit == nil {
		task = gui.createRenderStringTask("No reflog history")
	} else {
		cmd := gui.OSCommand.ExecutableFromString(
			gui.GitCommand.ShowCmdStr(commit.Sha, gui.State.Modes.Filtering.Path),
		)

		task = gui.createRunPtyTask(cmd)
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
	var lastReflogCommit *commands.Commit
	if len(state.ReflogCommits) > 0 {
		lastReflogCommit = state.ReflogCommits[0]
	}

	refresh := func(stateCommits *[]*commands.Commit, filterPath string) error {
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
		if err := refresh(&state.FilteredReflogCommits, state.Modes.Filtering.Path); err != nil {
			return err
		}
	} else {
		state.FilteredReflogCommits = state.ReflogCommits
	}

	return gui.postRefreshUpdate(gui.Contexts.ReflogCommits.Context)
}

func (gui *Gui) handleCheckoutReflogCommit(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	err := gui.ask(askOpts{
		title:  gui.Tr.SLocalize("checkoutCommit"),
		prompt: gui.Tr.SLocalize("SureCheckoutThisCommit"),
		handleConfirm: func() error {
			return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{})
		},
	})
	if err != nil {
		return err
	}

	gui.State.Panels.ReflogCommits.SelectedLineIdx = 0

	return nil
}

func (gui *Gui) handleCreateReflogResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedReflogCommit()

	return gui.createResetMenu(commit.Sha)
}

func (gui *Gui) handleViewReflogCommitFiles() error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	return gui.switchToCommitFilesContext(commit.Sha, false, gui.Contexts.ReflogCommits.Context, "commits")
}
