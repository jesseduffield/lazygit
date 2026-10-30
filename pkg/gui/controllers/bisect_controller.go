package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
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
		ListControllerTrait: NewListControllerTrait(
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
			Keys:        opts.GetKeys(opts.Config.Commits.ViewBisectOptions),
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
	}
	return self.openStartBisectMenu(info, commit)
}

func (self *BisectController) openMidBisectMenu(info *git_commands.BisectInfo, commit *models.Commit) error {
	// If we have a current hash already, then we always want to use that one. If
	// not, we're still picking the initial commits before we really start, so
	// use the selected commit in that case.

	bisecting := info.GetCurrentHash() != ""
	hashToMark := lo.Ternary(bisecting, info.GetCurrentHash(), commit.Hash())
	shortHashToMark := utils.ShortHash(hashToMark)

	// For marking a commit as bad, when we're not already bisecting, we require
	// a single item selected, but once we are bisecting, it doesn't matter because
	// the action applies to the HEAD commit rather than the selected commit.
	var singleItemIfNotBisecting *types.DisabledReason
	if !bisecting {
		singleItemIfNotBisecting = self.require(self.singleItemSelected())()
	}

	menuItems := []*types.MenuItem{
		{
			Label: fmt.Sprintf(self.c.Tr.Bisect.Mark, shortHashToMark, info.NewTerm()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectMark)
				if err := self.c.Git().Bisect.Mark(hashToMark, info.NewTerm()); err != nil {
					return err
				}

				return self.afterMark(types.KeepCommitSelectionByHash)
			},
			DisabledReason: singleItemIfNotBisecting,
			Keys:           menuKey('b'),
		},
		{
			Label: fmt.Sprintf(self.c.Tr.Bisect.Mark, shortHashToMark, info.OldTerm()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectMark)
				if err := self.c.Git().Bisect.Mark(hashToMark, info.OldTerm()); err != nil {
					return err
				}

				return self.afterMark(types.KeepCommitSelectionByHash)
			},
			DisabledReason: singleItemIfNotBisecting,
			Keys:           menuKey('g'),
		},
		{
			Label: fmt.Sprintf(self.c.Tr.Bisect.SkipCurrent, shortHashToMark),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectSkip)
				if err := self.c.Git().Bisect.Skip(hashToMark); err != nil {
					return err
				}

				return self.afterMark(types.KeepCommitSelectionByHash)
			},
			DisabledReason: singleItemIfNotBisecting,
			Keys:           menuKey('s'),
		},
	}
	if info.GetCurrentHash() != "" && info.GetCurrentHash() != commit.Hash() {
		menuItems = append(menuItems, lo.ToPtr(types.MenuItem{
			Label: fmt.Sprintf(self.c.Tr.Bisect.SkipSelected, commit.ShortHash()),
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.BisectSkip)
				if err := self.c.Git().Bisect.Skip(commit.Hash()); err != nil {
					return err
				}

				// Git may pick a different commit to test now. Keep the skipped
				// commit selected anyway, so that the user can go on skipping
				// commits near it.
				return self.afterMark(types.KeepCommitSelectionIndex)
			},
			DisabledReason: self.require(self.singleItemSelected())(),
			Keys:           menuKey('S'),
		}))
	}
	menuItems = append(menuItems, lo.ToPtr(types.MenuItem{
		Label: self.c.Tr.Bisect.ResetOption,
		OnPress: func() error {
			return self.c.Helpers().Bisect.Reset()
		},
		Keys: menuKey('r'),
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
				Label: fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortHash(), info.NewTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.c.Git().Bisect.Start(); err != nil {
						return err
					}

					if err := self.c.Git().Bisect.Mark(commit.Hash(), info.NewTerm()); err != nil {
						return err
					}

					self.c.Helpers().Bisect.PostBisectCommandRefresh()
					return nil
				},
				DisabledReason: self.require(self.singleItemSelected())(),
				Keys:           menuKey('b'),
			},
			{
				Label: fmt.Sprintf(self.c.Tr.Bisect.MarkStart, commit.ShortHash(), info.OldTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.c.Git().Bisect.Start(); err != nil {
						return err
					}

					if err := self.c.Git().Bisect.Mark(commit.Hash(), info.OldTerm()); err != nil {
						return err
					}

					self.c.Helpers().Bisect.PostBisectCommandRefresh()
					return nil
				},
				DisabledReason: self.require(self.singleItemSelected())(),
				Keys:           menuKey('g'),
			},
			{
				Label: self.c.Tr.Bisect.ChooseTerms,
				OnPress: func() error {
					self.c.Prompt(types.PromptOpts{
						Title: self.c.Tr.Bisect.OldTermPrompt,
						HandleConfirm: func(oldTerm string) error {
							self.c.Prompt(types.PromptOpts{
								Title: self.c.Tr.Bisect.NewTermPrompt,
								HandleConfirm: func(newTerm string) error {
									self.c.LogAction(self.c.Tr.Actions.StartBisect)
									if err := self.c.Git().Bisect.StartWithTerms(oldTerm, newTerm); err != nil {
										return err
									}

									self.c.Helpers().Bisect.PostBisectCommandRefresh()
									return nil
								},
							})
							return nil
						},
					})
					return nil
				},
				Keys: menuKey('t'),
			},
		},
	})
}

func (self *BisectController) showBisectCompleteMessage(candidateHashes []string) error {
	prompt := self.c.Tr.Bisect.CompletePrompt
	if len(candidateHashes) > 1 {
		prompt = self.c.Tr.Bisect.CompletePromptIndeterminate
	}

	formattedCommits, err := self.c.Git().Commit.GetCommitsOneline(candidateHashes)
	if err != nil {
		return err
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.Bisect.CompleteTitle,
		Prompt: fmt.Sprintf(prompt, strings.TrimSpace(formattedCommits)),
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.ResetBisect)
			if err := self.c.Git().Bisect.Reset(); err != nil {
				return err
			}

			self.c.Helpers().Bisect.PostBisectCommandRefresh()
			return nil
		},
	})

	return nil
}

func (self *BisectController) afterMark(commitSelection types.CommitSelectionBehavior) error {
	done, candidateHashes, err := self.c.Git().Bisect.IsDone()
	if err != nil {
		return err
	}

	self.c.Refresh(types.RefreshOptions{CommitSelection: commitSelection})

	if done {
		return self.showBisectCompleteMessage(candidateHashes)
	}

	return nil
}
