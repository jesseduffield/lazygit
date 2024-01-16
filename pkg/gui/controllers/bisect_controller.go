package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type BisectController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c *ControllerCommon
}

var _ types.IController = &BisectController{}

func NewBisectController(
	c *ControllerCommon,
) *BisectController {
	return &BisectController{
		baseController: baseController{},
		c:              c,
		ListControllerTrait: NewListControllerTrait[*models.Commit](
			c,
			c.Contexts().LocalCommits,
			c.Contexts().LocalCommits.GetSelected,
			c.Contexts().LocalCommits.GetSelectedItems,
		),
	}
}

func (self *BisectController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewBisectOptions),
			Handler:     opts.Guards.OutsideFilterMode(self.withItem(self.openMenu)),
			Description: self.c.Tr.ViewBisectOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *BisectController) openMenu(commit *models.Commit) error {
	// no shame in getting this directly rather than using the cached value
	// given how cheap it is to obtain
	info := self.c.Git().Bisect.GetInfo()
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
	waitToReselect := selectCurrentAfter && !self.c.Git().Bisect.ReachableFromStart(info)

	// If we have a current sha already, then we always want to use that one. If
	// not, we're still picking the initial commits before we really start, so
	// use the selected commit in that case.

	bisecting := info.GetCurrentSha() != ""
	shaToMark := lo.Ternary(bisecting, info.GetCurrentSha(), commit.Sha)
	shortShaToMark := utils.ShortSha(shaToMark)

	// For marking a commit as bad, when we're not already bisecting, we require
	// a single item selected, but once we are bisecting, it doesn't matter because
	// the action applies to the HEAD commit rather than the selected commit.
	var singleItemIfNotBisecting *types.DisabledReason
	if !bisecting {
		singleItemIfNotBisecting = self.require(self.singleItemSelected())()
	}

	menuItems := []*types.MenuItem{
		{
			Label: fmt.Sprintf(self.c.Tr.Bisect.Mark, shortShaToMark, info.NewTerm()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectMark)
				if err := self.c.Git().Bisect.Mark(shaToMark, info.NewTerm()); err != nil {
					return self.c.Error(err)
				}

				return self.afterMark(selectCurrentAfter, waitToReselect)
			},
			DisabledReason: singleItemIfNotBisecting,
			Key:            'b',
		},
		{
			Label: fmt.Sprintf(self.c.Tr.Bisect.Mark, shortShaToMark, info.OldTerm()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectMark)
				if err := self.c.Git().Bisect.Mark(shaToMark, info.OldTerm()); err != nil {
					return self.c.Error(err)
				}

				return self.afterMark(selectCurrentAfter, waitToReselect)
			},
			DisabledReason: singleItemIfNotBisecting,
			Key:            'g',
		},
		{
			Label: fmt.Sprintf(self.c.Tr.Bisect.SkipCurrent, shortShaToMark),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectSkip)
				if err := self.c.Git().Bisect.Skip(shaToMark); err != nil {
					return self.c.Error(err)
				}

				return self.afterMark(selectCurrentAfter, waitToReselect)
			},
			DisabledReason: singleItemIfNotBisecting,
			Key:            's',
		},
	}
	if info.GetCurrentSha() != "" && info.GetCurrentSha() != commit.Sha {
		menuItems = append(menuItems, lo.ToPtr(types.MenuItem{
			Label: fmt.Sprintf(self.c.Tr.Bisect.SkipSelected, commit.ShortSha()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectSkip)
				if err := self.c.Git().Bisect.Skip(commit.Sha); err != nil {
					return self.c.Error(err)
				}

				return self.afterMark(selectCurrentAfter, waitToReselect)
			},
			DisabledReason: self.require(self.singleItemSelected())(),
			Key:            'S',
		}))
	}
	menuItems = append(menuItems, lo.ToPtr(types.MenuItem{
		Label: self.c.Tr.Bisect.ResetOption,
		OnPress: func() error {
			return self.c.Helpers().Bisect.Reset()
		},
		Key: 'r',
	}))

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
				Label: fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortSha(), info.NewTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.c.Git().Bisect.Start(); err != nil {
						return self.c.Error(err)
					}

					if err := self.c.Git().Bisect.Mark(commit.Sha, info.NewTerm()); err != nil {
						return self.c.Error(err)
					}

					return self.c.Helpers().Bisect.PostBisectCommandRefresh()
				},
				DisabledReason: self.require(self.singleItemSelected())(),
				Key:            'b',
			},
			{
				Label: fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortSha(), info.OldTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.c.Git().Bisect.Start(); err != nil {
						return self.c.Error(err)
					}

					if err := self.c.Git().Bisect.Mark(commit.Sha, info.OldTerm()); err != nil {
						return self.c.Error(err)
					}

					return self.c.Helpers().Bisect.PostBisectCommandRefresh()
				},
				DisabledReason: self.require(self.singleItemSelected())(),
				Key:            'g',
			},
			{
				Label: self.c.Tr.Bisect.ChooseTerms,
				OnPress: func() error {
					return self.c.Prompt(types.PromptOpts{
						Title: self.c.Tr.Bisect.OldTermPrompt,
						HandleConfirm: func(oldTerm string) error {
							return self.c.Prompt(types.PromptOpts{
								Title: self.c.Tr.Bisect.NewTermPrompt,
								HandleConfirm: func(newTerm string) error {
									self.c.LogAction(self.c.Tr.Actions.StartBisect)
									if err := self.c.Git().Bisect.StartWithTerms(oldTerm, newTerm); err != nil {
										return self.c.Error(err)
									}

									return self.c.Helpers().Bisect.PostBisectCommandRefresh()
								},
							})
						},
					})
				},
				Key: 't',
			},
		},
	})
}

func (self *BisectController) showBisectCompleteMessage(candidateShas []string) error {
	prompt := self.c.Tr.Bisect.CompletePrompt
	if len(candidateShas) > 1 {
		prompt = self.c.Tr.Bisect.CompletePromptIndeterminate
	}

	formattedCommits, err := self.c.Git().Commit.GetCommitsOneline(candidateShas)
	if err != nil {
		return self.c.Error(err)
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Bisect.CompleteTitle,
		Prompt: fmt.Sprintf(prompt, strings.TrimSpace(formattedCommits)),
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ResetBisect)
			if err := self.c.Git().Bisect.Reset(); err != nil {
				return self.c.Error(err)
			}

			return self.c.Helpers().Bisect.PostBisectCommandRefresh()
		},
	})
}

func (self *BisectController) afterMark(selectCurrent bool, waitToReselect bool) error {
	done, candidateShas, err := self.c.Git().Bisect.IsDone()
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

		return self.c.Helpers().Bisect.PostBisectCommandRefresh()
	}
}

func (self *BisectController) selectCurrentBisectCommit() {
	info := self.c.Git().Bisect.GetInfo()
	if info.GetCurrentSha() != "" {
		// find index of commit with that sha, move cursor to that.
		for i, commit := range self.c.Model().Commits {
			if commit.Sha == info.GetCurrentSha() {
				self.context().SetSelection(i)
				_ = self.context().HandleFocus(types.OnFocusOpts{})
				break
			}
		}
	}
}

func (self *BisectController) context() *context.LocalCommitsContext {
	return self.c.Contexts().LocalCommits
}
