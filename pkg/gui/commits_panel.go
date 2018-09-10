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
			gui.Log.Errorf("Failed to get commits view: %s\n", err)
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

		err = gui.refreshStatus()
		if err != nil {
			gui.Log.Errorf("Failed to refresh status in refreshCommits: %s\n", err)
			return err
		}

		if gui.g.CurrentView().Name() == "commits" {

			err = gui.handleCommitSelect()
			if err != nil {
				gui.Log.Errorf("Failed to handleCommitSelect in refreshCommits: %s\n", err)
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

	err := gui.createConfirmationPanel(gui.g, v, gui.Tr.SLocalize("ResetToCommit"), gui.Tr.SLocalize("SureResetThisCommit"),
		func(g *gocui.Gui, v *gocui.View) error {

			commit, err := gui.getSelectedCommit()
			if err != nil {
				gui.Log.Errorf("Failed to get selected commit at handleResetToCommit: %s\n", err)
				return err
			}

			err = gui.GitCommand.ResetToCommit(commit.Sha)
			if err != nil {
				err = gui.createErrorPanel(gui.g, err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to create error panel at handleResetToCommit: %s\n", err)
					return err
				}
			}

			err = gui.refreshCommits()
			if err != nil {
				gui.Log.Errorf("Failed to refresh commits at handleResetToCommit: %s\n", err)
				return err
			}

			err = gui.refreshFiles()
			if err != nil {
				gui.Log.Errorf("Failed to refresh files at handleResetToCommit: %s\n", err)
				return err
			}

			err = gui.resetOrigin(v)
			if err != nil {
				gui.Log.Errorf("Failed to reset origin at handleResetToCommit %s\n", err)
				return err
			}

			err = gui.handleCommitSelect()
			if err != nil {
				gui.Log.Errorf("Failed to handle commit select at handleResetToCommit: %s\n", err)
				return err
			}

			return nil
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to create confirmation panel at handleResetToCommit: %s\n", err)
		return err
	}

	return nil
}

// handleCommitSelect gets called when a commit needs to be select.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitSelect() error {

	err := gui.renderGlobalOptions()
	if err != nil {
		gui.Log.Errorf("Failed to render global options at handleCommitSelect%s\n", err)
		return err
	}

	commit, err := gui.getSelectedCommit()
	if err != nil {

		if err.Error() != gui.Tr.SLocalize("NoCommitsThisBranch") {
			gui.Log.Errorf("Failed to select commit at handleResetToCommit: %s\n", err)
			return err
		}

		err = gui.renderString(gui.g, "main", gui.Tr.SLocalize("NoCommitsThisBranch"))
		if err != nil {
			gui.Log.Errorf("Failed to render string at handleResetToCommit: %s\n", err)
			return err
		}

		return nil
	}

	commitText := gui.GitCommand.Show(commit.Sha)

	err = gui.renderString(gui.g, "main", commitText)
	if err != nil {
		gui.Log.Errorf("Failed to render string at handleResetToCommit: %s\n", err)
		return err
	}

	return nil
}

// handleCommitSquashDown gets called when the user wants to squash down
// commits.
// g and v gets passed by gocui but g is not used.
// If anything goes wrong, it returns an error.
func (gui *Gui) handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {

	if gui.getItemPosition(v) != 0 {

		err := gui.createErrorPanel(gui.g, gui.Tr.SLocalize("OnlySquashTopmostCommit"))
		if err != nil {
			gui.Log.Errorf("Failed to create errorpanel at handleCommitSquashDown: %s\n", err)
			return err
		}

		return nil
	}

	if len(gui.State.Commits) == 1 {

		err := gui.createErrorPanel(gui.g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleCommitSquashDown: %s\n", err)
			return err
		}

		return nil
	}

	commit, err := gui.getSelectedCommit()
	if err != nil {
		gui.Log.Errorf("Failed to get selected commit at handleCommitSquashDown: %s\n", err)
		return err
	}

	err = gui.GitCommand.SquashPreviousTwoCommits(commit.Name)
	if err != nil {

		err = gui.createErrorPanel(gui.g, err.Error())
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleCommitSquashDown: %s\n", err)
			return err
		}

		return nil
	}

	err = gui.refreshCommits()
	if err != nil {
		gui.Log.Errorf("Failed to refresh commits at handleCommitSquashDown: %s\n", err)
		return err
	}

	err = gui.refreshStatus()
	if err != nil {
		gui.Log.Errorf("Failed to refresh status at handleCommitSquashDown: %s\n", err)
		return err
	}

	err = gui.handleCommitSelect()
	if err != nil {
		gui.Log.Errorf("Failed to handleCommitSelect at handleCommitSquashDown: %s\n", err)
		return err
	}

	return nil
}

// handleCommitFixup is called when a user wants to fix a commit.
// g and v are passed to by the gocui library but only v is used.
// If anything goes wrong it returns an error.
func (gui *Gui) handleCommitFixup(g *gocui.Gui, v *gocui.View) error {

	if len(gui.State.Commits) == 1 {

		err := gui.createErrorPanel(gui.g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleCommitFixup: %s\n", err)
			return err
		}

		return nil
	}

	if gui.anyUnStagedChanges(gui.State.Files) {
		err := gui.createErrorPanel(gui.g, gui.Tr.SLocalize("CantFixupWhileUnstagedChanges"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleCommitFixup: %s\n", err)
			return err
		}

		return nil
	}

	branch := gui.State.Branches[0]

	commit, err := gui.getSelectedCommit()
	if err != nil {
		gui.Log.Errorf("Failed to get selected commit: %s\n", err)
		return err
	}

	message := gui.Tr.SLocalize("SureFixupThisCommit")

	err = gui.createConfirmationPanel(gui.g, v, gui.Tr.SLocalize("Fixup"), message,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.SquashFixupCommit(branch.Name, commit.Sha)
			if err != nil {

				err = gui.createErrorPanel(gui.g, err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to create error panel at handleCommitFixup: %s\n", err)
					return err
				}

				return nil
			}

			err = gui.refreshCommits()
			if err != nil {
				gui.Log.Errorf("Failed to refresh commits at handleCommitFixup: %s\n", err)
				return err
			}

			err = gui.refreshStatus()
			if err != nil {
				gui.Log.Errorf("Failed to refresh status at handleCommitFixup: %s\n", err)
				return err
			}

			return nil
		}, nil)
	if err != nil {
		gui.Log.Errorf("Failed to create confirmation panel at handleCommitFixup: %s\n", err)
		return err
	}

	return nil
}

// handleRenameCommit is called when a user wants to rename a commit.
// g and v are passed by the gocui library but only v is used.
// If anything goes wrong it returns an error.
func (gui *Gui) handleRenameCommit(g *gocui.Gui, v *gocui.View) error {

	if gui.getItemPosition(v) != 0 {

		err := gui.createErrorPanel(gui.g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleRenameCommit: %s\n", err)
			return err
		}

		return nil
	}

	err := gui.createPromptPanel(gui.g, v, gui.Tr.SLocalize("renameCommit"),
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.RenameCommit(v.Buffer())
			if err != nil {

				err = gui.createErrorPanel(gui.g, err.Error())
				if err != nil {
					gui.Log.Errorf("Failed to create error panel at handleRenameCommit: %s\n", err)
					return err
				}

				return nil
			}

			err = gui.refreshCommits()
			if err != nil {
				gui.Log.Errorf("Failed to refresh commits at handleRenameCommit: %s\n", err)
				return err
			}

			err = gui.handleCommitSelect()
			if err != nil {
				gui.Log.Errorf("Failed to handleCommitSelect at handleRenameCommit: %s\n", err)
				return err
			}

			return nil
		})
	if err != nil {
		gui.Log.Errorf("Failed to create prompt panel at handleRenameCommit: %s\n", err)
		return err
	}

	return nil
}

// handleRenameCommitEditor is called when the user wants to edit the
// commit naming in an editor.
// g and v are passed by the gocui library, but only v is used.
// If anything goes wrong, it returns an error.
func (gui *Gui) handleRenameCommitEditor(g *gocui.Gui, v *gocui.View) error {

	if gui.getItemPosition(v) != 0 {

		err := gui.createErrorPanel(gui.g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
		if err != nil {
			gui.Log.Errorf("Failed to create error panel at handleRenameCommitEditor: %s\n", err)
			return err
		}

		return nil
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
		gui.Log.Errorf("Failed to get the commits view at getSelectedCommit: %s\n", err)
		return commands.Commit{}, err
	}

	if len(gui.State.Commits) == 0 {
		gui.Log.Errorf(gui.Tr.SLocalize("NoCommitsThisBranch"))
		return commands.Commit{}, errors.New(gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	lineNumber := gui.getItemPosition(v)
	if lineNumber > len(gui.State.Commits)-1 {
		gui.Log.Info(gui.Tr.SLocalize("PotentialErrInGetselectedCommit"), gui.State.Commits, lineNumber)
		return gui.State.Commits[len(gui.State.Commits)-1], nil
	}

	return gui.State.Commits[lineNumber], nil
}
