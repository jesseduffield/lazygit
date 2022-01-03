package gui

import (
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

type ReflogActionKind int

const (
	CHECKOUT ReflogActionKind = iota
	COMMIT
	REBASE
	CURRENT_REBASE
)

type reflogAction struct {
	kind ReflogActionKind
	from string
	to   string
}

// Here we're going through the reflog and maintaining a counter that represents how many
// undos/redos/user actions we've seen. when we hit a user action we call the callback specifying
// what the counter is up to and the nature of the action.
// If we find ourselves mid-rebase, we just return because undo/redo mid rebase
// requires knowledge of previous TODO file states, which you can't just get from the reflog.
// Though we might support this later, hence the use of the CURRENT_REBASE action kind.
func (gui *Gui) parseReflogForActions(onUserAction func(counter int, action reflogAction) (bool, error)) error {
	counter := 0
	reflogCommits := gui.State.FilteredReflogCommits
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
				action = &reflogAction{kind: CURRENT_REBASE, from: prevCommitSha}
			}
		} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^rebase -i \(start\)`); ok {
			action = &reflogAction{kind: REBASE, from: prevCommitSha, to: rebaseFinishCommitSha}
			rebaseFinishCommitSha = ""
		}

		if action != nil {
			if action.kind != CURRENT_REBASE && action.from == action.to {
				// if we're going from one place to the same place we'll ignore the action.
				continue
			}
			ok, err := onUserAction(counter, *action)
			if ok {
				return err
			}
			counter--
		}
	}
	return nil
}

func (gui *Gui) reflogUndo() error {
	undoEnvVars := []string{"GIT_REFLOG_ACTION=[lazygit undo]"}
	undoingStatus := gui.Tr.UndoingStatus

	if gui.GitCommand.WorkingTreeState() == commands.REBASE_MODE_REBASING {
		return gui.createErrorPanel(gui.Tr.LcCantUndoWhileRebasing)
	}

	span := gui.Tr.Spans.Undo

	return gui.parseReflogForActions(func(counter int, action reflogAction) (bool, error) {
		if counter != 0 {
			return false, nil
		}

		switch action.kind {
		case COMMIT, REBASE:
			return true, gui.handleHardResetWithAutoStash(action.from, handleHardResetWithAutoStashOptions{
				EnvVars:       undoEnvVars,
				WaitingStatus: undoingStatus,
				span:          span,
			})
		case CHECKOUT:
			return true, gui.handleCheckoutRef(action.from, handleCheckoutRefOptions{
				EnvVars:       undoEnvVars,
				WaitingStatus: undoingStatus,
				span:          span,
			})
		}

		gui.Log.Error("didn't match on the user action when trying to undo")
		return true, nil
	})
}

func (gui *Gui) reflogRedo() error {
	redoEnvVars := []string{"GIT_REFLOG_ACTION=[lazygit redo]"}
	redoingStatus := gui.Tr.RedoingStatus

	if gui.GitCommand.WorkingTreeState() == commands.REBASE_MODE_REBASING {
		return gui.createErrorPanel(gui.Tr.LcCantRedoWhileRebasing)
	}

	span := gui.Tr.Spans.Redo

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
				span:          span,
			})
		case CHECKOUT:
			return true, gui.handleCheckoutRef(action.to, handleCheckoutRefOptions{
				EnvVars:       redoEnvVars,
				WaitingStatus: redoingStatus,
				span:          span,
			})
		}

		gui.Log.Error("didn't match on the user action when trying to redo")
		return true, nil
	})
}

type handleHardResetWithAutoStashOptions struct {
	WaitingStatus string
	EnvVars       []string
	span          string
}

// only to be used in the undo flow for now
func (gui *Gui) handleHardResetWithAutoStash(commitSha string, options handleHardResetWithAutoStashOptions) error {
	gitCommand := gui.GitCommand.WithSpan(options.span)

	reset := func() error {
		if err := gui.resetToRef(commitSha, "hard", options.span, options.EnvVars); err != nil {
			return gui.surfaceError(err)
		}
		return nil
	}

	// if we have any modified tracked files we need to ask the user if they want us to stash for them
	dirtyWorkingTree := len(gui.trackedFiles()) > 0 || len(gui.stagedFiles()) > 0
	if dirtyWorkingTree {
		// offer to autostash changes
		return gui.ask(askOpts{
			title:  gui.Tr.AutoStashTitle,
			prompt: gui.Tr.AutoStashPrompt,
			handleConfirm: func() error {
				return gui.WithWaitingStatus(options.WaitingStatus, func() error {
					if err := gitCommand.StashSave(gui.Tr.StashPrefix + commitSha); err != nil {
						return gui.surfaceError(err)
					}
					if err := reset(); err != nil {
						return err
					}

					err := gitCommand.StashDo(0, "pop")
					if err := gui.refreshSidePanels(refreshOptions{}); err != nil {
						return err
					}
					if err != nil {
						return gui.surfaceError(err)
					}
					return nil
				})
			},
		})
	}

	return gui.WithWaitingStatus(options.WaitingStatus, func() error {
		return reset()
	})
}
