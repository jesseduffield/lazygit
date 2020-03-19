package gui

import (
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

// list panel functions

func (gui *Gui) getSelectedReflogCommit() *commands.Commit {
	selectedLine := gui.State.Panels.ReflogCommits.SelectedLine
	if selectedLine == -1 || len(gui.State.ReflogCommits) == 0 {
		return nil
	}

	return gui.State.ReflogCommits[selectedLine]
}

func (gui *Gui) handleReflogCommitSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	gui.getMainView().Title = "Reflog Entry"

	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return gui.newStringTask("main", "No reflog history")
	}
	v.FocusPoint(0, gui.State.Panels.ReflogCommits.SelectedLine)

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.ShowCmdStr(commit.Sha),
	)
	if err := gui.newPtyTask("main", cmd); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) refreshReflogCommits() error {
	commits, err := gui.GitCommand.GetReflogCommits()
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	gui.State.ReflogCommits = commits

	if gui.getCommitsView().Context == "reflog-commits" {
		return gui.renderReflogCommitsWithSelection()
	}

	return nil
}

func (gui *Gui) renderReflogCommitsWithSelection() error {
	commitsView := gui.getCommitsView()

	gui.refreshSelectedLine(&gui.State.Panels.ReflogCommits.SelectedLine, len(gui.State.ReflogCommits))
	displayStrings := presentation.GetCommitListDisplayStrings(gui.State.ReflogCommits, gui.State.ScreenMode != SCREEN_NORMAL)
	gui.renderDisplayStrings(commitsView, displayStrings)
	if gui.g.CurrentView() == commitsView && commitsView.Context == "reflog-commits" {
		if err := gui.handleReflogCommitSelect(gui.g, commitsView); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) handleCheckoutReflogCommit(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	err := gui.createConfirmationPanel(g, gui.getCommitsView(), true, gui.Tr.SLocalize("checkoutCommit"), gui.Tr.SLocalize("SureCheckoutThisCommit"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.handleCheckoutRef(commit.Sha)
	}, nil)
	if err != nil {
		return err
	}

	gui.State.Panels.ReflogCommits.SelectedLine = 0

	return nil
}

func (gui *Gui) handleCreateReflogResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedReflogCommit()

	return gui.createResetMenu(commit.Sha)
}

type reflogAction struct {
	regexStr string
	action   func(match []string, commitSha string, prevCommitSha string) (bool, error)
}

func (gui *Gui) reflogUndo(g *gocui.Gui, v *gocui.View) error {
	reflogActions := []reflogAction{
		{
			regexStr: `^checkout: moving from ([\S]+)`,
			action: func(match []string, commitSha string, prevCommitSha string) (bool, error) {
				if len(match) <= 1 {
					return false, nil
				}
				return true, gui.handleCheckoutRef(match[1])
			},
		},
		{
			regexStr: `^commit|^rebase -i \(start\)`,
			action: func(match []string, commitSha string, prevCommitSha string) (bool, error) {
				return true, gui.resetToRef(prevCommitSha, "hard")
			},
		},
	}

	for i, reflogCommit := range gui.State.ReflogCommits {
		for _, action := range reflogActions {
			re := regexp.MustCompile(action.regexStr)
			match := re.FindStringSubmatch(reflogCommit.Name)
			gui.Log.Warn(action.regexStr)
			gui.Log.Warn(spew.Sdump(match))
			if len(match) == 0 {
				continue
			}
			prevCommitSha := ""
			if len(gui.State.ReflogCommits)-1 >= i+1 {
				prevCommitSha = gui.State.ReflogCommits[i+1].Sha
			}
			gui.Log.Warn(prevCommitSha)

			done, err := action.action(match, reflogCommit.Sha, prevCommitSha)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
		}
	}

	return nil
}
