package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
)

// list panel functions

func (gui *Gui) getSelectedSubCommit() *models.Commit {
	selectedLine := gui.State.Panels.SubCommits.SelectedLineIdx
	commits := gui.State.SubCommits
	if selectedLine == -1 || len(commits) == 0 {
		return nil
	}

	return commits[selectedLine]
}

func (gui *Gui) handleSubCommitSelect() error {
	commit := gui.getSelectedSubCommit()
	var task updateTask
	if commit == nil {
		task = NewRenderStringTask("No commits")
	} else {
		cmdObj := gui.Git.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())

		task = NewRunPtyTask(cmdObj)
	}

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Commit",
			task:  task,
		},
	})
}

func (gui *Gui) handleCheckoutSubCommit() error {
	commit := gui.getSelectedSubCommit()
	if commit == nil {
		return nil
	}

	err := gui.Ask(AskOpts{
		Title:  gui.Tr.LcCheckoutCommit,
		Prompt: gui.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{span: gui.Tr.Spans.CheckoutCommit})
		},
	})
	if err != nil {
		return err
	}

	gui.State.Panels.SubCommits.SelectedLineIdx = 0

	return nil
}

func (gui *Gui) handleCreateSubCommitResetMenu() error {
	commit := gui.getSelectedSubCommit()

	return gui.createResetMenu(commit.Sha)
}

func (gui *Gui) handleViewSubCommitFiles() error {
	commit := gui.getSelectedSubCommit()
	if commit == nil {
		return nil
	}

	return gui.switchToCommitFilesContext(commit.Sha, false, gui.State.Contexts.SubCommits, "branches")
}

func (gui *Gui) switchToSubCommitsContext(refName string) error {
	// need to populate my sub commits
	builder := commands.NewCommitListBuilder(gui.Log, gui.Git, gui.OS, gui.Tr)

	commits, err := builder.GetCommits(
		commands.GetCommitsOptions{
			Limit:                gui.State.Panels.Commits.LimitCommits,
			FilterPath:           gui.State.Modes.Filtering.GetPath(),
			IncludeRebaseCommits: false,
			RefName:              refName,
		},
	)
	if err != nil {
		return err
	}

	gui.State.SubCommits = commits
	gui.State.Panels.SubCommits.refName = refName
	gui.State.Panels.SubCommits.SelectedLineIdx = 0
	gui.State.Contexts.SubCommits.SetParentContext(gui.currentSideListContext())

	return gui.pushContext(gui.State.Contexts.SubCommits)
}

func (gui *Gui) handleSwitchToSubCommits() error {
	currentContext := gui.currentSideListContext()
	if currentContext == nil {
		return nil
	}

	return gui.switchToSubCommitsContext(currentContext.GetSelectedItemId())
}
