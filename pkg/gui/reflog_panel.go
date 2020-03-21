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

type reflogAction struct {
	regexStr string
	action   func(match []string, commitSha string, onDone func()) (bool, error)
}

func (gui *Gui) reflogKey(reflogCommit *commands.Commit) string {
	return reflogCommit.Date + reflogCommit.Name
}

func (gui *Gui) idxOfUndoReflogKey(key string) int {
	for i, reflogCommit := range gui.State.ReflogCommits {
		if gui.reflogKey(reflogCommit) == key {
			return i
		}
	}
	return -1
}

func (gui *Gui) setUndoReflogKey(key string) {
	gui.State.Undo.ReflogKey = key
	// adding one because this is called before we actually refresh the reflog on our end
	// so the index will soon change.
	gui.State.Undo.ReflogIdx = gui.idxOfUndoReflogKey(key) + 1
}

func (gui *Gui) reflogUndo(g *gocui.Gui, v *gocui.View) error {
	reflogCommits := gui.State.ReflogCommits

	// if the index of the previous reflog entry has changed, we need to start from the beginning, because it means there's been user input.
	startIndex := gui.State.Undo.ReflogIdx
	if gui.idxOfUndoReflogKey(gui.State.Undo.ReflogKey) != gui.State.Undo.ReflogIdx {
		gui.State.Undo.UndoCount = 0
		startIndex = 0
	}

	reflogActions := []reflogAction{
		{
			regexStr: `^checkout: moving from ([\S]+)`,
			action: func(match []string, commitSha string, onDone func()) (bool, error) {
				if len(match) <= 1 {
					return false, nil
				}
				return true, gui.handleCheckoutRef(match[1], handleCheckoutRefOptions{
					OnDone:        onDone,
					WaitingStatus: gui.Tr.SLocalize("UndoingStatus"),
					EnvVars:       []string{"GIT_REFLOG_ACTION=[lazygit]"},
				},
				)
			},
		},
		{
			regexStr: `^commit|^rebase -i \(start\)|^reset: moving to|^pull`,
			action: func(match []string, commitSha string, onDone func()) (bool, error) {
				return true, gui.handleHardResetWithAutoStash(commitSha, onDone)
			},
		},
	}

	for offsetIdx, reflogCommit := range reflogCommits[startIndex:] {
		i := offsetIdx + startIndex
		for _, action := range reflogActions {
			re := regexp.MustCompile(action.regexStr)
			match := re.FindStringSubmatch(reflogCommit.Name)
			if len(match) == 0 {
				continue
			}
			prevCommitSha := ""
			if len(reflogCommits)-1 >= i+1 {
				prevCommitSha = reflogCommits[i+1].Sha
			}

			nextKey := gui.reflogKey(reflogCommits[i+1])
			onDone := func() {
				gui.setUndoReflogKey(nextKey)
				gui.State.Undo.UndoCount++
			}

			isMatchingAction, err := action.action(match, prevCommitSha, onDone)
			if !isMatchingAction {
				continue
			}

			return err
		}
	}

	return nil
}

func (gui *Gui) reflogRedo(g *gocui.Gui, v *gocui.View) error {
	reflogCommits := gui.State.ReflogCommits

	// if the index of the previous reflog entry has changed there is nothing to redo because there's been a user action
	startIndex := gui.State.Undo.ReflogIdx
	if gui.idxOfUndoReflogKey(gui.State.Undo.ReflogKey) != gui.State.Undo.ReflogIdx || startIndex == 0 || gui.State.Undo.UndoCount == 0 {
		return nil
	}

	reflogActions := []reflogAction{
		{
			regexStr: `^checkout: moving from [\S]+ to ([\S]+)`,
			action: func(match []string, commitSha string, onDone func()) (bool, error) {
				if len(match) <= 1 {
					return false, nil
				}
				return true, gui.handleCheckoutRef(match[1], handleCheckoutRefOptions{
					OnDone:        onDone,
					WaitingStatus: gui.Tr.SLocalize("RedoingStatus"),
					EnvVars:       []string{"GIT_REFLOG_ACTION=[lazygit]"},
				},
				)
			},
		},
		{
			regexStr: `^commit|^rebase -i \(start\)|^reset: moving to|^pull`,
			action: func(match []string, commitSha string, onDone func()) (bool, error) {
				return true, gui.handleHardResetWithAutoStash(commitSha, onDone)
			},
		},
	}

	for i := startIndex - 1; i > 0; i++ {
		reflogCommit := reflogCommits[i]

		for _, action := range reflogActions {
			re := regexp.MustCompile(action.regexStr)
			match := re.FindStringSubmatch(reflogCommit.Name)
			if len(match) == 0 {
				continue
			}

			prevKey := gui.reflogKey(reflogCommits[i-1])
			onDone := func() {
				gui.setUndoReflogKey(prevKey)
				gui.State.Undo.UndoCount--
			}

			isMatchingAction, err := action.action(match, reflogCommit.Sha, onDone)
			if !isMatchingAction {
				continue
			}

			return err
		}
	}

	return nil
}

// only to be used in the undo flow for now
func (gui *Gui) handleHardResetWithAutoStash(commitSha string, onDone func()) error {
	// if we have any modified tracked files we need to ask the user if they want us to stash for them
	dirtyWorkingTree := false
	for _, file := range gui.State.Files {
		if file.Tracked {
			dirtyWorkingTree = true
			break
		}
	}

	reset := func() error {
		if err := gui.resetToRef(commitSha, "hard", commands.RunCommandOptions{EnvVars: []string{"GIT_REFLOG_ACTION=[lazygit]"}}); err != nil {
			return gui.createErrorPanel(gui.g, err.Error())
		}
		onDone()
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
