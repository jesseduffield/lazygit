package gui

import (
	"regexp"

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
	displayStrings := presentation.GetReflogCommitListDisplayStrings(gui.State.ReflogCommits, gui.State.ScreenMode != SCREEN_NORMAL)
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
		return gui.handleCheckoutRef(commit.Sha, handleCheckoutRefOptions{})
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

const (
	USER_ACTION = iota
	UNDO
	REDO
)

type reflogAction struct {
	regexStr string
	action   func(match []string, commitSha string) error
	kind     int
}

func (gui *Gui) reflogUndo(g *gocui.Gui, v *gocui.View) error {
	reflogCommits := gui.State.ReflogCommits

	envVars := []string{"GIT_REFLOG_ACTION=[lazygit undo]"}

	// first step, work out what it is, second step, do an action if necessary

	reflogActions := []reflogAction{
		{
			regexStr: `^checkout: moving from ([\S]+)`,
			kind:     USER_ACTION,
			action: func(match []string, commitSha string) error {
				return gui.handleCheckoutRef(match[1], handleCheckoutRefOptions{
					WaitingStatus: gui.Tr.SLocalize("UndoingStatus"),
					EnvVars:       envVars,
				},
				)
			},
		},
		{
			regexStr: `^commit|^rebase -i \(start\)|^reset: moving to|^pull`,
			kind:     USER_ACTION,
			action: func(match []string, commitSha string) error {
				return gui.handleHardResetWithAutoStash(commitSha, handleHardResetWithAutoStashOptions{EnvVars: envVars})
			},
		},
		{
			regexStr: `^\[lazygit undo\]`,
			kind:     UNDO,
		},
		{
			regexStr: `^\[lazygit redo\]`,
			kind:     REDO,
		},
	}

	counter := 0
	for i, reflogCommit := range reflogCommits {
		for _, action := range reflogActions {
			re := regexp.MustCompile(action.regexStr)
			match := re.FindStringSubmatch(reflogCommit.Name)
			if len(match) == 0 {
				continue
			}

			switch action.kind {
			case UNDO:
				counter++
			case REDO:
				counter--
			case USER_ACTION:
				counter--
				if counter == -1 {
					prevCommitSha := ""
					if len(reflogCommits)-1 >= i+1 {
						prevCommitSha = reflogCommits[i+1].Sha
					}
					return action.action(match, prevCommitSha)
				}
			}
		}
	}

	return nil
}

func (gui *Gui) reflogRedo(g *gocui.Gui, v *gocui.View) error {
	reflogCommits := gui.State.ReflogCommits

	envVars := []string{"GIT_REFLOG_ACTION=[lazygit redo]"}

	reflogActions := []reflogAction{
		{
			regexStr: `^checkout: moving from [\S]+ to ([\S]+)`,
			action: func(match []string, commitSha string) error {
				return gui.handleCheckoutRef(match[1], handleCheckoutRefOptions{
					WaitingStatus: gui.Tr.SLocalize("RedoingStatus"),
					EnvVars:       envVars,
				},
				)
			},
		},
		{
			regexStr: `^commit|^rebase -i \(start\)|^reset: moving to|^pull`,
			action: func(match []string, commitSha string) error {
				return gui.handleHardResetWithAutoStash(commitSha, handleHardResetWithAutoStashOptions{EnvVars: envVars})
			},
		},
		{
			regexStr: `^\[lazygit undo\]`,
			kind:     UNDO,
		},
		{
			regexStr: `^\[lazygit redo\]`,
			kind:     REDO,
		},
	}

	counter := 0
	for _, reflogCommit := range reflogCommits {
		for _, action := range reflogActions {
			re := regexp.MustCompile(action.regexStr)
			match := re.FindStringSubmatch(reflogCommit.Name)
			if len(match) == 0 {
				continue
			}

			switch action.kind {
			case UNDO:
				counter++
			case REDO:
				counter--
			case USER_ACTION:
				counter--
				if counter == 0 {
					return action.action(match, reflogCommit.Sha)
				} else if counter < 0 {
					return nil
				}
			}
		}
	}

	return nil
}

type handleHardResetWithAutoStashOptions struct {
	EnvVars []string
}

// only to be used in the undo flow for now
func (gui *Gui) handleHardResetWithAutoStash(commitSha string, options handleHardResetWithAutoStashOptions) error {
	// if we have any modified tracked files we need to ask the user if they want us to stash for them
	dirtyWorkingTree := false
	for _, file := range gui.State.Files {
		if file.Tracked {
			dirtyWorkingTree = true
			break
		}
	}

	reset := func() error {
		if err := gui.resetToRef(commitSha, "hard", commands.RunCommandOptions{EnvVars: options.EnvVars}); err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		return nil
	}

	if dirtyWorkingTree {
		// offer to autostash changes
		return gui.createConfirmationPanel(gui.g, gui.getBranchesView(), true, gui.Tr.SLocalize("AutoStashTitle"), gui.Tr.SLocalize("AutoStashPrompt"), func(g *gocui.Gui, v *gocui.View) error {
			return gui.WithWaitingStatus(gui.Tr.SLocalize("UndoingStatus"), func() error {
				if err := gui.GitCommand.StashSave(gui.Tr.SLocalize("StashPrefix") + commitSha); err != nil {
					return gui.createErrorPanel(g, err.Error())
				}
				if err := reset(); err != nil {
					return err
				}

				if err := gui.GitCommand.StashDo(0, "pop"); err != nil {
					if err := gui.refreshSidePanels(g); err != nil {
						return err
					}
					return gui.createErrorPanel(g, err.Error())
				}
				return gui.refreshSidePanels(g)
			})
		}, nil)
	}

	return gui.WithWaitingStatus(gui.Tr.SLocalize("UndoingStatus"), func() error {
		if err := reset(); err != nil {
			return err
		}
		return gui.refreshSidePanels(gui.g)
	})
}
