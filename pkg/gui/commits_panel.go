package gui

import (
	"errors"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

// refreshCommits refreshes the commits view.
// If something goes wrong, it returns an error
func (gui *Gui) refreshCommits() error {
	gui.g.Update(func(*gocui.Gui) error {
		red := color.New(color.FgRed)
		yellow := color.New(color.FgYellow)
		white := color.New(color.FgWhite)
		shaColor := white

		gui.State.Commits = gui.GitCommand.GetCommits()

		v, err := gui.g.View("commits")
		if err != nil {
			return err
		}

		v.Clear()

		for _, commit := range gui.State.Commits {
			if commit.Pushed {
				shaColor = red
			} else {
				shaColor = yellow
			}

			shaColor.Fprint(v, commit.Sha+" ")
			white.Fprintln(v, commit.Name)
		}

		if err := gui.refreshStatus(); err != nil {
			return err
		}

		if gui.g.CurrentView().Name() == "commits" {
			if err := gui.handleCommitSelect(); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

// handleResetToCommit is called when the user wants to reset to a commit.
// g and v are passed by the gocui, but only v is used.
// If anything goes wrong it returns an error.
func (gui *Gui) handleResetToCommit(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel(v, gui.Tr.SLocalize("ResetToCommit"), gui.Tr.SLocalize("SureResetThisCommit"),
		func(g *gocui.Gui, v *gocui.View) error {

			commit, err := gui.getSelectedCommit()
			if err != nil {
				return err
			}

			if err := gui.GitCommand.ResetToCommit(commit.Sha); err != nil {
				return gui.createErrorPanel(err.Error())
			}

			if err := gui.refreshCommits(); err != nil {
				return err
			}

			if err := gui.refreshFiles(); err != nil {
				return err
			}

			if err := gui.resetOrigin(v); err != nil {
				return err
			}

			return gui.handleCommitSelect()
		}, nil)
}

// handleCommitSelect gets called when a commit needs to be select.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitSelect() error {
	if err := gui.renderGlobalOptions(); err != nil {
		return err
	}

	commit, err := gui.getSelectedCommit()
	if err != nil {
		if err.Error() != gui.Tr.SLocalize("NoCommitsThisBranch") {
			return err
		}
		return gui.renderString("main", gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	commitText := gui.GitCommand.Show(commit.Sha)
	return gui.renderString("main", commitText)
}

// handleCommitSquashDown gets called when the user wants to squash down
// commits.
// g and v gets passed by gocui but g is not used.
// If anything goes wrong, it returns an error.
func (gui *Gui) handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {
	if gui.getItemPosition(v) != 0 {
		return gui.createErrorPanel(gui.Tr.SLocalize("OnlySquashTopmostCommit"))
	}

	if len(gui.State.Commits) == 1 {
		return gui.createErrorPanel(gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}

	commit, err := gui.getSelectedCommit()
	if err != nil {
		return err
	}

	if err := gui.GitCommand.SquashPreviousTwoCommits(commit.Name); err != nil {
		return gui.createErrorPanel(err.Error())
	}

	if err := gui.refreshCommits(); err != nil {
		return err
	}

	if err := gui.refreshStatus(); err != nil {
		return err
	}

	return gui.handleCommitSelect()
}

// handleCommitFixup is called when a user wants to fix a commit.
// g and v are passed to by the gocui library but only v is used.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitFixup(g *gocui.Gui, v *gocui.View) error {
	if len(gui.State.Commits) == 1 {
		return gui.createErrorPanel(gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}

	if gui.anyUnStagedChanges(gui.State.Files) {
		return gui.createErrorPanel(gui.Tr.SLocalize("CantFixupWhileUnstagedChanges"))
	}

	commit, err := gui.getSelectedCommit()
	if err != nil {
		return err
	}

	branch := gui.State.Branches[0]
	message := gui.Tr.SLocalize("SureFixupThisCommit")

	return gui.createConfirmationPanel(v, gui.Tr.SLocalize("Fixup"), message,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := gui.GitCommand.SquashFixupCommit(branch.Name, commit.Sha); err != nil {
				return gui.createErrorPanel(err.Error())
			}

			if err := gui.refreshCommits(); err != nil {
				return err
			}

			return gui.refreshStatus()
		}, nil)
}

// handleRenameCommit is called when a user wants to rename a commit.
// g and v are passed by the gocui library but only v is used.
// If anything goes wrong it returns an error.
func (gui *Gui) handleRenameCommit(g *gocui.Gui, v *gocui.View) error {
	if gui.getItemPosition(v) != 0 {
		return gui.createErrorPanel(gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}

	return gui.createPromptPanel(v, gui.Tr.SLocalize("renameCommit"),
		func(g *gocui.Gui, v *gocui.View) error {
			if err := gui.GitCommand.RenameCommit(v.Buffer()); err != nil {
				return gui.createErrorPanel(err.Error())
			}

			if err := gui.refreshCommits(); err != nil {
				return err
			}

			return gui.handleCommitSelect()
		})
}

// handleRenameCommitEditor is called when the user wants to edit the
// commit naming in an editor.
// g and v are passed by the gocui library, but only v is used.
// If anything goes wrong, it returns an error.
func (gui *Gui) handleRenameCommitEditor(g *gocui.Gui, v *gocui.View) error {
	if gui.getItemPosition(v) != 0 {
		return gui.createErrorPanel(gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}

	gui.SubProcess = gui.GitCommand.PrepareCommitAmendSubProcess()

	gui.g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})

	return nil
}

// getSelectedCommit gets the selected commit.
// returns the commit that is currently selected and an error if something
// went wrong.
func (gui *Gui) getSelectedCommit() (commands.Commit, error) {
	v, err := gui.g.View("commits")
	if err != nil {
		return commands.Commit{}, err
	}

	if len(gui.State.Commits) == 0 {
		return commands.Commit{}, errors.New(gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	lineNumber := gui.getItemPosition(v)
	if lineNumber > len(gui.State.Commits)-1 {
		gui.Log.Info(gui.Tr.SLocalize("PotentialErrInGetselectedCommit"), gui.State.Commits, lineNumber)
		return gui.State.Commits[len(gui.State.Commits)-1], nil
	}

	return gui.State.Commits[lineNumber], nil
}
