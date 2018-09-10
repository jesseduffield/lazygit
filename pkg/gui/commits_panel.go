package gui

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) refreshCommits() error {

	gui.g.Update(func(*gocui.Gui) error {

		red := color.New(color.FgRed)
		yellow := color.New(color.FgYellow)
		white := color.New(color.FgWhite)
		shaColor := white

		gui.State.Commits = gui.GitCommand.GetCommits()

		v, err := gui.g.View("commits")
		if err != nil {
			gui.Log.Error(fmt.Sprintf("Failed to get commits view: %s\n", err))
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
			gui.Log.Error(fmt.Sprintf("Failed to refresh status in refreshCommits: %s\n", err))
			return err
		}

		if gui.g.CurrentView().Name() == "commits" {
			err = gui.handleCommitSelect(gui.g, v)
			if err != nil {
				gui.Log.Error(fmt.Sprintf("Failed to handleCommitSelect in refreshCommits: %s\n", err))
				return err
			}
		}

		return nil
	})

	return nil
}

func (gui *Gui) handleResetToCommit(g *gocui.Gui, commitView *gocui.View) error {
	return gui.createConfirmationPanel(g, commitView, gui.Tr.SLocalize("ResetToCommit"), gui.Tr.SLocalize("SureResetThisCommit"), func(g *gocui.Gui, v *gocui.View) error {
		commit, err := gui.getSelectedCommit(g)
		if err != nil {
			panic(err)
		}
		if err := gui.GitCommand.ResetToCommit(commit.Sha); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(); err != nil {
			panic(err)
		}
		if err := gui.refreshFiles(); err != nil {
			panic(err)
		}
		gui.resetOrigin(commitView)
		return gui.handleCommitSelect(g, nil)
	}, nil)
}

func (gui *Gui) handleCommitSelect(g *gocui.Gui, v *gocui.View) error {
	err := gui.renderGlobalOptions()
	if err != nil {
		return err
	}
	commit, err := gui.getSelectedCommit(g)
	if err != nil {
		if err.Error() != gui.Tr.SLocalize("NoCommitsThisBranch") {
			return err
		}
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoCommitsThisBranch"))
	}
	commitText := gui.GitCommand.Show(commit.Sha)
	return gui.renderString(g, "main", commitText)
}

// handleCommitSquashDown gets called when the user wants to squash down
// commits.
// g and v gets passed by gocui but g is not used.
// If anything goes wrong, it returns an error
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

	commit, err := gui.getSelectedCommit(gui.g)
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

	err = gui.handleCommitSelect(gui.g, v)
	if err != nil {
		gui.Log.Errorf("Failed to handleCommitSelect at handleCommitSquashDown: %s\n", err)
		return err
	}

	return nil
}

func (gui *Gui) handleCommitFixup(g *gocui.Gui, v *gocui.View) error {
	if len(gui.State.Commits) == 1 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}
	if gui.anyUnStagedChanges(gui.State.Files) {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantFixupWhileUnstagedChanges"))
	}
	branch := gui.State.Branches[0]
	commit, err := gui.getSelectedCommit(g)
	if err != nil {
		return err
	}
	message := gui.Tr.SLocalize("SureFixupThisCommit")
	gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("Fixup"), message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.SquashFixupCommit(branch.Name, commit.Sha); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(); err != nil {
			panic(err)
		}
		return gui.refreshStatus()
	}, nil)
	return nil
}

func (gui *Gui) handleRenameCommit(g *gocui.Gui, v *gocui.View) error {
	if gui.getItemPosition(gui.getCommitsView(g)) != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}
	return gui.createPromptPanel(g, v, gui.Tr.SLocalize("renameCommit"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.RenameCommit(v.Buffer()); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(); err != nil {
			panic(err)
		}
		return gui.handleCommitSelect(g, v)
	})
	return nil
}

func (gui *Gui) handleRenameCommitEditor(g *gocui.Gui, v *gocui.View) error {
	if gui.getItemPosition(gui.getCommitsView(g)) != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}

	gui.SubProcess = gui.GitCommand.PrepareCommitAmendSubProcess()
	g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})

	return nil
}

func (gui *Gui) getSelectedCommit(g *gocui.Gui) (commands.Commit, error) {
	v, err := g.View("commits")
	if err != nil {
		panic(err)
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
