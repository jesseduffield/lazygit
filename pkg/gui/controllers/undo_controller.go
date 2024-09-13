package controllers

import (
	"errors"
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Quick summary of how this all works:
// when you want to undo or redo, we start from the top of the reflog and work
// down until we've reached the last user-initiated reflog entry that hasn't already been undone
// we then do the reverse of what that reflog describes.
// When we do this, we create a new reflog entry, and tag it as either an undo or redo
// Then, next time we want to undo, we'll use those entries to know which user-initiated
// actions we can skip. E.g. if I do three things, A, B, and C, and hit undo twice,
// the reflog will read UUCBA, and when I read the first two undos, I know to skip the following
// two user actions, meaning we end up undoing reflog entry C. Redoing works in a similar way.

type UndoController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &UndoController{}

func NewUndoController(
	c *ControllerCommon,
) *UndoController {
	return &UndoController{
		baseController: baseController{},
		c:              c,
	}
}

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

func (self *UndoController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Undo),
			Handler:     self.reflogUndo,
			Description: self.c.Tr.UndoReflog,
			Tooltip:     self.c.Tr.UndoTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Redo),
			Handler:     self.reflogRedo,
			Description: self.c.Tr.RedoReflog,
			Tooltip:     self.c.Tr.RedoTooltip,
		},
	}

	return bindings
}

func (self *UndoController) Context() types.Context {
	return nil
}

func (self *UndoController) reflogUndo() error {
	undoEnvVars := []string{"GIT_REFLOG_ACTION=[lazygit undo]"}
	undoingStatus := self.c.Tr.UndoingStatus

	if self.c.Git().Status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
		return errors.New(self.c.Tr.CantUndoWhileRebasing)
	}

	return self.parseReflogForActions(func(counter int, action reflogAction) (bool, error) {
		if counter != 0 {
			return false, nil
		}

		switch action.kind {
		case COMMIT, REBASE:
			self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.Actions.Undo,
				Prompt: fmt.Sprintf(self.c.Tr.HardResetAutostashPrompt, action.from),
				HandleConfirm: func() error {
					self.c.LogAction(self.c.Tr.Actions.Undo)
					return self.hardResetWithAutoStash(action.from, hardResetOptions{
						EnvVars:       undoEnvVars,
						WaitingStatus: undoingStatus,
					})
				},
			})
			return true, nil

		case CHECKOUT:
			self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.Actions.Undo,
				Prompt: fmt.Sprintf(self.c.Tr.CheckoutPrompt, action.from),
				HandleConfirm: func() error {
					self.c.LogAction(self.c.Tr.Actions.Undo)
					return self.c.Helpers().Refs.CheckoutRef(action.from, types.CheckoutRefOptions{
						EnvVars:       undoEnvVars,
						WaitingStatus: undoingStatus,
					})
				},
			})
			return true, nil

		case CURRENT_REBASE:
			// do nothing
		}

		self.c.Log.Error("didn't match on the user action when trying to undo")
		return true, nil
	})
}

func (self *UndoController) reflogRedo() error {
	redoEnvVars := []string{"GIT_REFLOG_ACTION=[lazygit redo]"}
	redoingStatus := self.c.Tr.RedoingStatus

	if self.c.Git().Status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
		return errors.New(self.c.Tr.CantRedoWhileRebasing)
	}

	return self.parseReflogForActions(func(counter int, action reflogAction) (bool, error) {
		// if we're redoing and the counter is zero, we just return
		if counter == 0 {
			return true, nil
		} else if counter > 1 {
			return false, nil
		}

		switch action.kind {
		case COMMIT, REBASE:
			self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.Actions.Redo,
				Prompt: fmt.Sprintf(self.c.Tr.HardResetAutostashPrompt, action.to),
				HandleConfirm: func() error {
					self.c.LogAction(self.c.Tr.Actions.Redo)
					return self.hardResetWithAutoStash(action.to, hardResetOptions{
						EnvVars:       redoEnvVars,
						WaitingStatus: redoingStatus,
					})
				},
			})
			return true, nil

		case CHECKOUT:
			self.c.Confirm(types.ConfirmOpts{
				Title:  self.c.Tr.Actions.Redo,
				Prompt: fmt.Sprintf(self.c.Tr.CheckoutPrompt, action.to),
				HandleConfirm: func() error {
					self.c.LogAction(self.c.Tr.Actions.Redo)
					return self.c.Helpers().Refs.CheckoutRef(action.to, types.CheckoutRefOptions{
						EnvVars:       redoEnvVars,
						WaitingStatus: redoingStatus,
					})
				},
			})
			return true, nil

		case CURRENT_REBASE:
			// do nothing
		}

		self.c.Log.Error("didn't match on the user action when trying to redo")
		return true, nil
	})
}

// Here we're going through the reflog and maintaining a counter that represents how many
// undos/redos/user actions we've seen. when we hit a user action we call the callback specifying
// what the counter is up to and the nature of the action.
// If we find ourselves mid-rebase, we just return because undo/redo mid rebase
// requires knowledge of previous TODO file states, which you can't just get from the reflog.
// Though we might support this later, hence the use of the CURRENT_REBASE action kind.
func (self *UndoController) parseReflogForActions(onUserAction func(counter int, action reflogAction) (bool, error)) error {
	counter := 0
	reflogCommits := self.c.Model().FilteredReflogCommits
	rebaseFinishCommitHash := ""
	var action *reflogAction
	for reflogCommitIdx, reflogCommit := range reflogCommits {
		action = nil

		prevCommitHash := ""
		if len(reflogCommits)-1 >= reflogCommitIdx+1 {
			prevCommitHash = reflogCommits[reflogCommitIdx+1].Hash
		}

		if rebaseFinishCommitHash == "" {
			if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^\[lazygit undo\]`); ok {
				counter++
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^\[lazygit redo\]`); ok {
				counter--
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^rebase (-i )?\(abort\)|^rebase (-i )?\(finish\)`); ok {
				rebaseFinishCommitHash = reflogCommit.Hash
			} else if ok, match := utils.FindStringSubmatch(reflogCommit.Name, `^checkout: moving from ([\S]+) to ([\S]+)`); ok {
				action = &reflogAction{kind: CHECKOUT, from: match[1], to: match[2]}
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^commit|^reset: moving to|^pull`); ok {
				action = &reflogAction{kind: COMMIT, from: prevCommitHash, to: reflogCommit.Hash}
			} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^rebase (-i )?\(start\)`); ok {
				// if we're here then we must be currently inside an interactive rebase
				action = &reflogAction{kind: CURRENT_REBASE, from: prevCommitHash}
			}
		} else if ok, _ := utils.FindStringSubmatch(reflogCommit.Name, `^rebase (-i )?\(start\)`); ok {
			action = &reflogAction{kind: REBASE, from: prevCommitHash, to: rebaseFinishCommitHash}
			rebaseFinishCommitHash = ""
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

type hardResetOptions struct {
	WaitingStatus string
	EnvVars       []string
}

// only to be used in the undo flow for now (does an autostash)
func (self *UndoController) hardResetWithAutoStash(commitHash string, options hardResetOptions) error {
	reset := func() error {
		return self.c.Helpers().Refs.ResetToRef(commitHash, "hard", options.EnvVars)
	}

	// if we have any modified tracked files we need to ask the user if they want us to stash for them
	dirtyWorkingTree := self.c.Helpers().WorkingTree.IsWorkingTreeDirty()
	if dirtyWorkingTree {
		// offer to autostash changes
		self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.AutoStashTitle,
			Prompt: self.c.Tr.AutoStashPrompt,
			HandleConfirm: func() error {
				return self.c.WithWaitingStatus(options.WaitingStatus, func(gocui.Task) error {
					if err := self.c.Git().Stash.Push(self.c.Tr.StashPrefix + commitHash); err != nil {
						return err
					}
					if err := reset(); err != nil {
						return err
					}

					err := self.c.Git().Stash.Pop(0)
					if err != nil {
						return err
					}
					return self.c.Refresh(types.RefreshOptions{})
				})
			},
		})
		return nil
	}

	return self.c.WithWaitingStatus(options.WaitingStatus, func(gocui.Task) error {
		return reset()
	})
}
