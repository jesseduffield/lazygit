package gui

import (
	"errors"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

func (gui *Gui) refreshCommits(g *gocui.Gui) error {

	g.Update(func(*gocui.Gui) error {

		gui.State.Commits = gui.GitCommand.GetCommits()

		v, err := g.View("commits")
		if err != nil {
			panic(err)
		}

		v.Clear()

		red := color.New(color.FgRed)
		yellow := color.New(color.FgYellow)
		white := color.New(color.FgWhite)
		shaColor := white

		for _, commit := range gui.State.Commits {

			if commit.Pushed {
				shaColor = red
			} else {
				shaColor = yellow
			}

			shaColor.Fprint(v, commit.Sha+" ")
			white.Fprintln(v, commit.Name)

		}

		err = gui.refreshStatus(g)
		if err != nil {
			return err
		}

		if g.CurrentView().Name() == "commits" {
			gui.handleCommitSelect(g, v)
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

		err = gui.GitCommand.ResetToCommit(commit.Sha)
		if err != nil {
			return gui.createErrorPanel(g, err.Error())
		}

		err = gui.refreshCommits(g)
		if err != nil {
			panic(err)
		}

		err = gui.refreshFiles(g)
		if err != nil {
			panic(err)
		}

		err = gui.resetOrigin(commitView)
		if err != nil {
			return err
		}

		return gui.handleCommitSelect(g, nil)
	}, nil)
}

func (gui *Gui) renderCommitsOptions(g *gocui.Gui) error {

	return gui.renderOptionsMap(g, map[string]string{
		"s":       gui.Tr.SLocalize("squashDown"),
		"r":       gui.Tr.SLocalize("rename"),
		"g":       gui.Tr.SLocalize("resetToThisCommit"),
		"f":       gui.Tr.SLocalize("fixupCommit"),
		"← → ↑ ↓": gui.Tr.SLocalize("navigate"),
	})
}

func (gui *Gui) handleCommitSelect(g *gocui.Gui, v *gocui.View) error {

	if err := gui.renderCommitsOptions(g); err != nil {
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

func (gui *Gui) handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {

	if gui.getItemPosition(v) != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlySquashTopmostCommit"))
	}

	if len(gui.State.Commits) == 1 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}

	commit, err := gui.getSelectedCommit(g)
	if err != nil {
		return err
	}

	err = gui.GitCommand.SquashPreviousTwoCommits(commit.Name)
	if err != nil {
		return gui.createErrorPanel(g, err.Error())
	}

	err = gui.refreshCommits(g)
	if err != nil {
		panic(err)
	}

	gui.refreshStatus(g)

	return gui.handleCommitSelect(g, v)
}

// TODO: move to files panel
func (gui *Gui) anyUnStagedChanges(files []commands.File) bool {

	for _, file := range files {

		if file.Tracked && file.HasUnstagedChanges {
			return true
		}

	}

	return false
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

	err = gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("Fixup"), message,
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.SquashFixupCommit(branch.Name, commit.Sha)
			if err != nil {
				return gui.createErrorPanel(g, err.Error())
			}

			err = gui.refreshCommits(g)
			if err != nil {
				panic(err)
			}

			return gui.refreshStatus(g)

		}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleRenameCommit(g *gocui.Gui, v *gocui.View) error {

	if gui.getItemPosition(v) != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}

	err := gui.createPromptPanel(g, v, gui.Tr.SLocalize("RenameCommit"),
		func(g *gocui.Gui, v *gocui.View) error {

			err := gui.GitCommand.RenameCommit(v.Buffer())
			if err != nil {
				return gui.createErrorPanel(g, err.Error())
			}

			err = gui.refreshCommits(g)
			if err != nil {
				panic(err)
			}

			return gui.handleCommitSelect(g, v)
		})
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleRenameCommitEditor(g *gocui.Gui, v *gocui.View) error {

	if gui.getItemPosition(v) != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}

	gui.SubProcess = gui.GitCommand.PrepareCommitAmendSubProcess()

	g.Update(
		func(g *gocui.Gui) error {
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
