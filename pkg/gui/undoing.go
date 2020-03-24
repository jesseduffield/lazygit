package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Quick summary of how this all works:
// when you want to undo or redo, we start from the top of the reflog and work
// down until we've reached the last user-initiated reflog entry that hasn't already been undone
// we then do the reverse of what that reflog describes.
// When we do this, we create a new reflog entry, and tag it as either an undo or redo
// Then, next time we want to undo, we'll use those entries to know which user-initiated
// actions we can skip. E.g. if I do do three things, A, B, and C, and hit undo twice,
// the reflog will read UUCBA, and when I read the first two undos, I know to skip the following
// two user actions, meaning we end up undoing reflog entry C. Redoing works in a similar way.

const (
	CHECKOUT = iota
	COMMIT
	REBASE
	CURRENT_REBASE
)

type reflogAction struct {
	kind int // one of CHECKOUT, REBASE, and COMMIT
	from string
	to   string
}

// Here we're going through the reflog and maintaining a counter that represents how many
// undos/redos/user actions we've seen. when we hit a user action we call the callback specifying
// what the counter is up to and the nature of the action.
// We can't take you from a non-interactive rebase state into an interactive rebase state, so if we hit
// a 'finish' or an 'abort' entry, we ignore everything else until we find the corresponding 'start' entry.
// If we find ourselves already in an interactive rebase and we've hit the start entry,
// we can't really do an undo because there's no way to redo back into the rebase.
// instead we just ask the user if they want to abort the rebase instead.
func (gui *Gui) parseReflogForActions(onUserAction func(counter int, action reflogAction) (bool, error)) error {
	counter := 0
	reflogCommits := gui.State.ReflogCommits
	rebaseFinishCommitSha := ""
	var action *reflogAction
	for reflogCommitIdx, reflogCommit := range reflogCommits {
		action = nil

		prevCommitSha := ""
		if len(reflogCommits)-1 >= reflogCommitIdx+1 {
			prevCommitSha = reflogCommits[reflogCommitIdx+1].Sha
		}

		if rebaseFinishCommitSha == "" {
			if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^\[lazygit undo\]`); ok {
				counter++
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^\[lazygit redo\]`); ok {
				counter--
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^rebase -i \(abort\)|^rebase -i \(finish\)`); ok {
				rebaseFinishCommitSha = reflogCommit.Sha
			} else if ok, match := utils.FindStringSubmatch(reflogCommit.Name, `^checkout: moving from ([\S]+) to ([\S]+)`); ok {
				action = &reflogAction{kind: CHECKOUT, from: match[1], to: match[2]}
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^commit|^reset: moving to|^pull`); ok {
				action = &reflogAction{kind: COMMIT, from: prevCommitSha, to: reflogCommit.Sha}
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^rebase -i \(start\)`); ok {
				// if we're here then we must be currently inside an interactive rebase
				action = &reflogAction{kind: CURRENT_REBASE}
			}
		} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^rebase -i \(start\)`); ok {
			action = &reflogAction{kind: REBASE, from: prevCommitSha, to: rebaseFinishCommitSha}
		}

		if action != nil {
			ok, err := onUserAction(counter, *action)
			if ok {
				return err
			}
			counter--
			if action.kind == REBASE {
				rebaseFinishCommitSha = ""
			}
		}
	}
	return nil
}

func (gui *Gui) reflogUndo(g *gocui.Gui, v *gocui.View) error {
	undoEnvVars := []string{"GIT_REFLOG_ACTION=[lazygit undo]"}
	undoingStatus := gui.Tr.SLocalize("UndoingStatus")

	return gui.parseReflogForActions(func(counter int, action reflogAction) (bool, error) {
		if counter != 0 {
			return false, nil
		}

		switch action.kind {
		case COMMIT, REBASE:
			return true, gui.handleHardResetWithAutoStash(action.from, handleHardResetWithAutoStashOptions{
				EnvVars:       undoEnvVars,
				WaitingStatus: undoingStatus,
			})
		case CURRENT_REBASE:
			return true, gui.createConfirmationPanel(g, v, true, gui.Tr.SLocalize("AbortRebase"), gui.Tr.SLocalize("UndoOutOfRebaseWarning"), func(g *gocui.Gui, v *gocui.View) error {
				return gui.genericMergeCommand("abort")
			}, nil)
		case CHECKOUT:
			return true, gui.handleCheckoutRef(action.from, handleCheckoutRefOptions{
				EnvVars:       undoEnvVars,
				WaitingStatus: undoingStatus,
			})
		}

		gui.Log.Error("didn't match on the user action when trying to undo")
		return true, nil
	})
}

func (gui *Gui) reflogRedo(g *gocui.Gui, v *gocui.View) error {
	redoEnvVars := []string{"GIT_REFLOG_ACTION=[lazygit redo]"}
	redoingStatus := gui.Tr.SLocalize("RedoingStatus")

	return gui.parseReflogForActions(func(counter int, action reflogAction) (bool, error) {
		// if we're redoing and the counter is zero, we just return
		if counter == 0 {
			return true, nil
		} else if counter > 1 {
			return false, nil
		}

		switch action.kind {
		case COMMIT, REBASE:
			return true, gui.handleHardResetWithAutoStash(action.to, handleHardResetWithAutoStashOptions{
				EnvVars:       redoEnvVars,
				WaitingStatus: redoingStatus,
			})
		case CURRENT_REBASE:
			// no idea if this is even possible but you certainly can't redo into the end of a rebase if you're still in the rebase
			return true, nil
		case CHECKOUT:
			return true, gui.handleCheckoutRef(action.to, handleCheckoutRefOptions{
				EnvVars:       redoEnvVars,
				WaitingStatus: redoingStatus,
			})
		}

		gui.Log.Error("didn't match on the user action when trying to redo")
		return true, nil
	})
}

type handleHardResetWithAutoStashOptions struct {
	WaitingStatus string
	EnvVars       []string
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
			return gui.WithWaitingStatus(options.WaitingStatus, func() error {
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

	return gui.WithWaitingStatus(options.WaitingStatus, func() error {
		if err := reset(); err != nil {
			return err
		}
		return gui.refreshSidePanels(gui.g)
	})
}
