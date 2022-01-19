package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (gui *Gui) handleOpenBisectMenu() error {
	if ok, err := gui.validateNotInFilterMode(); err != nil || !ok {
		return err
	}

	// no shame in getting this directly rather than using the cached value
	// given how cheap it is to obtain
	info := gui.Git.Bisect.GetInfo()
	commit := gui.getSelectedLocalCommit()
	if info.Started() {
		return gui.openMidBisectMenu(info, commit)
	} else {
		return gui.openStartBisectMenu(info, commit)
	}
}

func (gui *Gui) openMidBisectMenu(info *git_commands.BisectInfo, commit *models.Commit) error {
	// if there is not yet a 'current' bisect commit, or if we have
	// selected the current commit, we need to jump to the next 'current' commit
	// after we perform a bisect action. The reason we don't unconditionally jump
	// is that sometimes the user will want to go and mark a few commits as skipped
	// in a row and they wouldn't want to be jumped back to the current bisect
	// commit each time.
	// Originally we were allowing the user to, from the bisect menu, select whether
	// they were talking about the selected commit or the current bisect commit,
	// and that was a bit confusing (and required extra keypresses).
	selectCurrentAfter := info.GetCurrentSha() == "" || info.GetCurrentSha() == commit.Sha

	menuItems := []*menuItem{
		{
			displayString: fmt.Sprintf(gui.Tr.Bisect.Mark, commit.ShortSha(), info.NewTerm()),
			onPress: func() error {
				gui.logAction(gui.Tr.Actions.BisectMark)
				if err := gui.Git.Bisect.Mark(commit.Sha, info.NewTerm()); err != nil {
					return gui.surfaceError(err)
				}

				return gui.afterMark(selectCurrentAfter)
			},
		},
		{
			displayString: fmt.Sprintf(gui.Tr.Bisect.Mark, commit.ShortSha(), info.OldTerm()),
			onPress: func() error {
				gui.logAction(gui.Tr.Actions.BisectMark)
				if err := gui.Git.Bisect.Mark(commit.Sha, info.OldTerm()); err != nil {
					return gui.surfaceError(err)
				}

				return gui.afterMark(selectCurrentAfter)
			},
		},
		{
			displayString: fmt.Sprintf(gui.Tr.Bisect.Skip, commit.ShortSha()),
			onPress: func() error {
				gui.logAction(gui.Tr.Actions.BisectSkip)
				if err := gui.Git.Bisect.Skip(commit.Sha); err != nil {
					return gui.surfaceError(err)
				}

				return gui.afterMark(selectCurrentAfter)
			},
		},
		{
			displayString: gui.Tr.Bisect.ResetOption,
			onPress: func() error {
				return gui.resetBisect()
			},
		},
	}

	return gui.createMenu(
		gui.Tr.Bisect.BisectMenuTitle,
		menuItems,
		createMenuOptions{showCancel: true},
	)
}

func (gui *Gui) openStartBisectMenu(info *git_commands.BisectInfo, commit *models.Commit) error {
	return gui.createMenu(
		gui.Tr.Bisect.BisectMenuTitle,
		[]*menuItem{
			{
				displayString: fmt.Sprintf(gui.Tr.Bisect.MarkStart, commit.ShortSha(), info.NewTerm()),
				onPress: func() error {
					gui.logAction(gui.Tr.Actions.StartBisect)
					if err := gui.Git.Bisect.Start(); err != nil {
						return gui.surfaceError(err)
					}

					if err := gui.Git.Bisect.Mark(commit.Sha, info.NewTerm()); err != nil {
						return gui.surfaceError(err)
					}

					return gui.postBisectCommandRefresh()
				},
			},
			{
				displayString: fmt.Sprintf(gui.Tr.Bisect.MarkStart, commit.ShortSha(), info.OldTerm()),
				onPress: func() error {
					gui.logAction(gui.Tr.Actions.StartBisect)
					if err := gui.Git.Bisect.Start(); err != nil {
						return gui.surfaceError(err)
					}

					if err := gui.Git.Bisect.Mark(commit.Sha, info.OldTerm()); err != nil {
						return gui.surfaceError(err)
					}

					return gui.postBisectCommandRefresh()
				},
			},
		},
		createMenuOptions{showCancel: true},
	)
}

func (gui *Gui) resetBisect() error {
	return gui.ask(askOpts{
		title:  gui.Tr.Bisect.ResetTitle,
		prompt: gui.Tr.Bisect.ResetPrompt,
		handleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.ResetBisect)
			if err := gui.Git.Bisect.Reset(); err != nil {
				return gui.surfaceError(err)
			}

			return gui.postBisectCommandRefresh()
		},
	})
}

func (gui *Gui) showBisectCompleteMessage(candidateShas []string) error {
	prompt := gui.Tr.Bisect.CompletePrompt
	if len(candidateShas) > 1 {
		prompt = gui.Tr.Bisect.CompletePromptIndeterminate
	}

	formattedCommits, err := gui.Git.Commit.GetCommitsOneline(candidateShas)
	if err != nil {
		return gui.surfaceError(err)
	}

	return gui.ask(askOpts{
		title:  gui.Tr.Bisect.CompleteTitle,
		prompt: fmt.Sprintf(prompt, strings.TrimSpace(formattedCommits)),
		handleConfirm: func() error {
			gui.logAction(gui.Tr.Actions.ResetBisect)
			if err := gui.Git.Bisect.Reset(); err != nil {
				return gui.surfaceError(err)
			}

			return gui.postBisectCommandRefresh()
		},
	})
}

func (gui *Gui) afterMark(selectCurrent bool) error {
	done, candidateShas, err := gui.Git.Bisect.IsDone()
	if err != nil {
		return gui.surfaceError(err)
	}

	if err := gui.postBisectCommandRefresh(); err != nil {
		return gui.surfaceError(err)
	}

	if done {
		return gui.showBisectCompleteMessage(candidateShas)
	}

	if selectCurrent {
		gui.selectCurrentBisectCommit()
	}

	return nil
}

func (gui *Gui) selectCurrentBisectCommit() {
	info := gui.Git.Bisect.GetInfo()
	if info.GetCurrentSha() != "" {
		// find index of commit with that sha, move cursor to that.
		for i, commit := range gui.State.Commits {
			if commit.Sha == info.GetCurrentSha() {
				gui.State.Contexts.BranchCommits.GetPanelState().SetSelectedLineIdx(i)
				_ = gui.State.Contexts.BranchCommits.HandleFocus()
				break
			}
		}
	}
}

func (gui *Gui) postBisectCommandRefresh() error {
	return gui.refreshSidePanels(refreshOptions{mode: ASYNC, scope: []RefreshableView{}})
}
