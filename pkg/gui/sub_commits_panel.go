package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers"
	"github.com/jesseduffield/lazygit/pkg/gui/popup"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

func (gui *Gui) subCommitsRenderToMain() error {
	commit := gui.getSelectedSubCommit()
	var task updateTask
	if commit == nil {
		task = NewRenderStringTask("No commits")
	} else {
		cmdObj := gui.git.Commit.ShowCmdObj(commit.Sha, gui.State.Modes.Filtering.GetPath())

		task = NewRunPtyTask(cmdObj.GetCmd())
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

	err := gui.c.Ask(popup.AskOpts{
		Title:  gui.c.Tr.LcCheckoutCommit,
		Prompt: gui.c.Tr.SureCheckoutThisCommit,
		HandleConfirm: func() error {
			gui.c.LogAction(gui.c.Tr.Actions.CheckoutCommit)
			return gui.refHelper.CheckoutRef(commit.Sha, types.CheckoutRefOptions{})
		},
	})
	if err != nil {
		return err
	}

	gui.State.Contexts.SubCommits.GetPanelState().SetSelectedLineIdx(0)

	return nil
}

func (gui *Gui) handleCreateSubCommitResetMenu() error {
	commit := gui.getSelectedSubCommit()

	return gui.refHelper.CreateGitResetMenu(commit.Sha)
}

func (gui *Gui) handleViewSubCommitFiles() error {
	commit := gui.getSelectedSubCommit()
	if commit == nil {
		return nil
	}

	return gui.SwitchToCommitFilesContext(controllers.SwitchToCommitFilesContextOpts{
		RefName:    commit.Sha,
		CanRebase:  false,
		Context:    gui.State.Contexts.SubCommits,
		WindowName: "branches",
	})
}

func (gui *Gui) switchToSubCommitsContext(refName string) error {
	// need to populate my sub commits
	commits, err := gui.git.Loaders.Commits.GetCommits(
		loaders.GetCommitsOptions{
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
	gui.State.Contexts.SubCommits.GetPanelState().SetSelectedLineIdx(0)
	gui.State.Contexts.SubCommits.SetParentContext(gui.currentSideListContext())

	return gui.c.PushContext(gui.State.Contexts.SubCommits)
}

func (gui *Gui) handleSwitchToSubCommits() error {
	currentContext := gui.currentSideListContext()
	if currentContext == nil {
		return nil
	}

	return gui.switchToSubCommitsContext(currentContext.GetSelectedItemId())
}
