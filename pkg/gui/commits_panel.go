package gui

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// list panel functions

func (gui *Gui) getSelectedCommit(g *gocui.Gui) *commands.Commit {
	selectedLine := gui.State.Panels.Commits.SelectedLine
	if selectedLine == -1 {
		return nil
	}

	return gui.State.Commits[selectedLine]
}

func (gui *Gui) handleCommitSelect(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return gui.renderString(g, "main", gui.Tr.SLocalize("NoCommitsThisBranch"))
	}

	if err := gui.focusPoint(0, gui.State.Panels.Commits.SelectedLine, v); err != nil {
		return err
	}
	commitText, err := gui.GitCommand.Show(commit.Sha)
	if err != nil {
		return err
	}
	return gui.renderString(g, "main", commitText)
}

func (gui *Gui) refreshCommits(g *gocui.Gui) error {
	g.Update(func(*gocui.Gui) error {
		commits, err := gui.GitCommand.GetCommits()
		if err != nil {
			return err
		}
		gui.State.Commits = commits

		gui.refreshSelectedLine(&gui.State.Panels.Commits.SelectedLine, len(gui.State.Commits))

		list, err := utils.RenderList(gui.State.Commits)
		if err != nil {
			return err
		}

		v := gui.getCommitsView()
		v.Clear()
		fmt.Fprint(v, list)

		gui.refreshStatus(g)
		if v == g.CurrentView() {
			gui.handleCommitSelect(g, v)
		}
		return nil
	})
	return nil
}

func (gui *Gui) handleCommitsNextLine(g *gocui.Gui, v *gocui.View) error {
	panelState := gui.State.Panels.Commits
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Commits), false)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleCommitSelect(gui.g, v)
}

func (gui *Gui) handleCommitsPrevLine(g *gocui.Gui, v *gocui.View) error {
	panelState := gui.State.Panels.Commits
	gui.changeSelectedLine(&panelState.SelectedLine, len(gui.State.Commits), true)

	if err := gui.resetOrigin(gui.getMainView()); err != nil {
		return err
	}
	return gui.handleCommitSelect(gui.g, v)
}

// specific functions

func (gui *Gui) handleResetToCommit(g *gocui.Gui, commitView *gocui.View) error {
	return gui.createConfirmationPanel(g, commitView, gui.Tr.SLocalize("ResetToCommit"), gui.Tr.SLocalize("SureResetThisCommit"), func(g *gocui.Gui, v *gocui.View) error {
		commit := gui.getSelectedCommit(g)
		if commit == nil {
			panic(errors.New(gui.Tr.SLocalize("NoCommitsThisBranch")))
		}
		if err := gui.GitCommand.ResetToCommit(commit.Sha); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(g); err != nil {
			panic(err)
		}
		if err := gui.refreshFiles(); err != nil {
			panic(err)
		}
		gui.resetOrigin(commitView)
		gui.State.Panels.Commits.SelectedLine = 0
		return gui.handleCommitSelect(g, commitView)
	}, nil)
}

func (gui *Gui) handleCommitSquashDown(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Commits.SelectedLine != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlySquashTopmostCommit"))
	}
	if len(gui.State.Commits) <= 1 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return errors.New(gui.Tr.SLocalize("NoCommitsThisBranch"))
	}
	if err := gui.GitCommand.SquashPreviousTwoCommits(commit.Name); err != nil {
		return gui.createErrorPanel(g, err.Error())
	}
	if err := gui.refreshCommits(g); err != nil {
		panic(err)
	}
	gui.refreshStatus(g)
	return gui.handleCommitSelect(g, v)
}

// TODO: move to files panel
func (gui *Gui) anyUnStagedChanges(files []*commands.File) bool {
	for _, file := range files {
		if file.Tracked && file.HasUnstagedChanges {
			return true
		}
	}
	return false
}

func (gui *Gui) handleCommitFixup(g *gocui.Gui, v *gocui.View) error {
	if len(gui.State.Commits) <= 1 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("YouNoCommitsToSquash"))
	}
	if gui.anyUnStagedChanges(gui.State.Files) {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("CantFixupWhileUnstagedChanges"))
	}
	branch := gui.State.Branches[0]
	commit := gui.getSelectedCommit(g)
	if commit == nil {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("NoCommitsThisBranch"))
	}
	message := gui.Tr.SLocalize("SureFixupThisCommit")
	gui.createConfirmationPanel(g, v, gui.Tr.SLocalize("Fixup"), message, func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.SquashFixupCommit(branch.Name, commit.Sha); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(g); err != nil {
			panic(err)
		}
		return gui.refreshStatus(g)
	}, nil)
	return nil
}

func (gui *Gui) handleRenameCommit(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Commits.SelectedLine != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}
	return gui.createPromptPanel(g, v, gui.Tr.SLocalize("renameCommit"), func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.GitCommand.RenameCommit(v.Buffer()); err != nil {
			return gui.createErrorPanel(g, err.Error())
		}
		if err := gui.refreshCommits(g); err != nil {
			panic(err)
		}
		return gui.handleCommitSelect(g, v)
	})
}

func (gui *Gui) handleRenameCommitEditor(g *gocui.Gui, v *gocui.View) error {
	if gui.State.Panels.Commits.SelectedLine != 0 {
		return gui.createErrorPanel(g, gui.Tr.SLocalize("OnlyRenameTopCommit"))
	}

	gui.SubProcess = gui.GitCommand.PrepareCommitAmendSubProcess()
	g.Update(func(g *gocui.Gui) error {
		return gui.Errors.ErrSubProcess
	})

	return nil
}
