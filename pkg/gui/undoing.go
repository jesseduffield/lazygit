package gui

import (
	"regexp"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
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
	USER_ACTION = iota
	UNDO
	REDO
)

type reflogAction struct {
	regexStr string
	action   func(match []string, commitSha string, waitingStatus string, envVars []string, isUndo bool) error
	kind     int
}

func (gui *Gui) reflogActions() []reflogAction {
	return []reflogAction{
		{
			regexStr: `^checkout: moving from ([\S]+) to ([\S]+)`,
			kind:     USER_ACTION,
			action: func(match []string, commitSha string, waitingStatus string, envVars []string, isUndo bool) error {
				branchName := match[2]
				if isUndo {
					branchName = match[1]
				}
				return gui.handleCheckoutRef(branchName, handleCheckoutRefOptions{
					WaitingStatus: waitingStatus,
					EnvVars:       envVars,
				},
				)
			},
		},
		{
			regexStr: `^commit|^rebase -i \(start\)|^reset: moving to|^pull`,
			kind:     USER_ACTION,
			action: func(match []string, commitSha string, waitingStatus string, envVars []string, isUndo bool) error {
				return gui.handleHardResetWithAutoStash(commitSha, handleHardResetWithAutoStashOptions{EnvVars: envVars, WaitingStatus: waitingStatus})
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
}

func (gui *Gui) reflogUndo(g *gocui.Gui, v *gocui.View) error {
	return gui.iterateUserActions(func(match []string, reflogCommits []*commands.Commit, reflogIdx int, action reflogAction, counter int) (bool, error) {
		if counter == -1 {
			prevCommitSha := ""
			if len(reflogCommits)-1 >= reflogIdx+1 {
				prevCommitSha = reflogCommits[reflogIdx+1].Sha
			}
			return true, action.action(match, prevCommitSha, gui.Tr.SLocalize("UndoingStatus"), []string{"GIT_REFLOG_ACTION=[lazygit undo]"}, true)
		} else {
			return false, nil
		}
	})
}

func (gui *Gui) reflogRedo(g *gocui.Gui, v *gocui.View) error {
	return gui.iterateUserActions(func(match []string, reflogCommits []*commands.Commit, reflogIdx int, action reflogAction, counter int) (bool, error) {
		if counter == 0 {
			return true, action.action(match, reflogCommits[reflogIdx].Sha, gui.Tr.SLocalize("RedoingStatus"), []string{"GIT_REFLOG_ACTION=[lazygit redo]"}, false)
		} else if counter < 0 {
			return true, nil
		} else {
			return false, nil
		}
	})
}

func (gui *Gui) iterateUserActions(onUserAction func(match []string, reflogCommits []*commands.Commit, reflogIdx int, action reflogAction, counter int) (bool, error)) error {
	reflogCommits := gui.State.ReflogCommits

	counter := 0
	for i, reflogCommit := range reflogCommits {
		for _, action := range gui.reflogActions() {
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
				shouldReturn, err := onUserAction(match, reflogCommits, i, action, counter)
				if err != nil {
					return err
				}
				if shouldReturn {
					return nil
				}
			}
		}
	}
	return nil
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
