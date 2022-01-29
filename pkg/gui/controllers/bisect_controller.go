package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type BisectController struct {
	c          *types.ControllerCommon
	getContext func() types.IListContext
	git        *commands.GitCommand

	getSelectedLocalCommit func() *models.Commit
	getCommits             func() []*models.Commit
}

var _ types.IController = &BisectController{}

func NewBisectController(
	c *types.ControllerCommon,
	getContext func() types.IListContext,
	git *commands.GitCommand,

	getSelectedLocalCommit func() *models.Commit,
	getCommits func() []*models.Commit,
) *BisectController {
	return &BisectController{
		c:          c,
		getContext: getContext,
		git:        git,

		getSelectedLocalCommit: getSelectedLocalCommit,
		getCommits:             getCommits,
	}
}

func (self *BisectController) Keybindings(getKey func(key string) interface{}, config config.KeybindingConfig, guards types.KeybindingGuards) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         getKey(config.Commits.ViewBisectOptions),
			Handler:     guards.OutsideFilterMode(self.checkSelected(self.openMenu)),
			Description: self.c.Tr.LcViewBisectOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *BisectController) openMenu(commit *models.Commit) error {
	// no shame in getting this directly rather than using the cached value
	// given how cheap it is to obtain
	info := self.git.Bisect.GetInfo()
	if info.Started() {
		return self.openMidBisectMenu(info, commit)
	} else {
		return self.openStartBisectMenu(info, commit)
	}
}

func (self *BisectController) openMidBisectMenu(info *git_commands.BisectInfo, commit *models.Commit) error {
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
	// we need to wait to reselect if our bisect commits aren't ancestors of our 'start'
	// ref, because we'll be reloading our commits in that case.
	waitToReselect := selectCurrentAfter && !self.git.Bisect.ReachableFromStart(info)

	menuItems := []*types.MenuItem{
		{
			DisplayString: fmt.Sprintf(self.c.Tr.Bisect.Mark, commit.ShortSha(), info.NewTerm()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectMark)
				if err := self.git.Bisect.Mark(commit.Sha, info.NewTerm()); err != nil {
					return self.c.Error(err)
				}

				return self.afterMark(selectCurrentAfter, waitToReselect)
			},
		},
		{
			DisplayString: fmt.Sprintf(self.c.Tr.Bisect.Mark, commit.ShortSha(), info.OldTerm()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectMark)
				if err := self.git.Bisect.Mark(commit.Sha, info.OldTerm()); err != nil {
					return self.c.Error(err)
				}

				return self.afterMark(selectCurrentAfter, waitToReselect)
			},
		},
		{
			DisplayString: fmt.Sprintf(self.c.Tr.Bisect.Skip, commit.ShortSha()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectSkip)
				if err := self.git.Bisect.Skip(commit.Sha); err != nil {
					return self.c.Error(err)
				}

				return self.afterMark(selectCurrentAfter, waitToReselect)
			},
		},
		{
			DisplayString: self.c.Tr.Bisect.ResetOption,
			OnPress: func() error {
				return self.Reset()
			},
		},
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Bisect.BisectMenuTitle,
		Items: menuItems,
	})
}

func (self *BisectController) openStartBisectMenu(info *git_commands.BisectInfo, commit *models.Commit) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Bisect.BisectMenuTitle,
		Items: []*types.MenuItem{
			{
				DisplayString: fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortSha(), info.NewTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.git.Bisect.Start(); err != nil {
						return self.c.Error(err)
					}

					if err := self.git.Bisect.Mark(commit.Sha, info.NewTerm()); err != nil {
						return self.c.Error(err)
					}

					return self.postBisectCommandRefresh()
				},
			},
			{
				DisplayString: fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortSha(), info.OldTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.git.Bisect.Start(); err != nil {
						return self.c.Error(err)
					}

					if err := self.git.Bisect.Mark(commit.Sha, info.OldTerm()); err != nil {
						return self.c.Error(err)
					}

					return self.postBisectCommandRefresh()
				},
			},
		},
	})
}

func (self *BisectController) Reset() error {
	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.Bisect.ResetTitle,
		Prompt: self.c.Tr.Bisect.ResetPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ResetBisect)
			if err := self.git.Bisect.Reset(); err != nil {
				return self.c.Error(err)
			}

			return self.postBisectCommandRefresh()
		},
	})
}

func (self *BisectController) showBisectCompleteMessage(candidateShas []string) error {
	prompt := self.c.Tr.Bisect.CompletePrompt
	if len(candidateShas) > 1 {
		prompt = self.c.Tr.Bisect.CompletePromptIndeterminate
	}

	formattedCommits, err := self.git.Commit.GetCommitsOneline(candidateShas)
	if err != nil {
		return self.c.Error(err)
	}

	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.Bisect.CompleteTitle,
		Prompt: fmt.Sprintf(prompt, strings.TrimSpace(formattedCommits)),
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ResetBisect)
			if err := self.git.Bisect.Reset(); err != nil {
				return self.c.Error(err)
			}

			return self.postBisectCommandRefresh()
		},
	})
}

func (self *BisectController) afterMark(selectCurrent bool, waitToReselect bool) error {
	done, candidateShas, err := self.git.Bisect.IsDone()
	if err != nil {
		return self.c.Error(err)
	}

	if err := self.afterBisectMarkRefresh(selectCurrent, waitToReselect); err != nil {
		return self.c.Error(err)
	}

	if done {
		return self.showBisectCompleteMessage(candidateShas)
	}

	return nil
}

func (self *BisectController) postBisectCommandRefresh() error {
	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{}})
}

func (self *BisectController) afterBisectMarkRefresh(selectCurrent bool, waitToReselect bool) error {
	selectFn := func() {
		if selectCurrent {
			self.selectCurrentBisectCommit()
		}
	}

	if waitToReselect {
		return self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{}, Then: selectFn})
	} else {
		selectFn()

		return self.postBisectCommandRefresh()
	}
}

func (self *BisectController) selectCurrentBisectCommit() {
	info := self.git.Bisect.GetInfo()
	if info.GetCurrentSha() != "" {
		// find index of commit with that sha, move cursor to that.
		for i, commit := range self.getCommits() {
			if commit.Sha == info.GetCurrentSha() {
				self.getContext().GetPanelState().SetSelectedLineIdx(i)
				_ = self.getContext().HandleFocus()
				break
			}
		}
	}
}

func (self *BisectController) checkSelected(callback func(*models.Commit) error) func() error {
	return func() error {
		commit := self.getSelectedLocalCommit()
		if commit == nil {
			return nil
		}

		return callback(commit)
	}
}

func (self *BisectController) Context() types.Context {
	return self.getContext()
}
