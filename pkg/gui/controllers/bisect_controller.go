package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	if !info.Started() {
		return self.openStartBisectMenu(info, commit)
	}
	if info.GetCurrentHash() == "" {
		return self.openMenuBeforeBisecting(info, commit)
	}
	return self.openMidBisectMenu(info, commit)
}

// The menu for when a bisect has started, but git doesn't know both a good and
// a bad commit yet. Its entries act on the selected commit.
func (self *BisectController) openMenuBeforeBisecting(info *git_commands.BisectInfo, commit *models.Commit) error {
	disabledReason := self.require(self.singleItemSelected())()
	ref := utils.ResolvePlaceholderString(self.c.Tr.Bisect.CommitRef,
		map[string]string{"hash": commit.ShortHash()})
	markLabel := func(term string) string {
		return utils.ResolvePlaceholderString(self.c.Tr.Bisect.MarkRef,
			map[string]string{"ref": ref, "term": term})
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Bisect.BisectMenuTitle,
		Items: []*types.MenuItem{
			self.markMenuItem(markLabel(info.NewTerm()), commit.Hash(), info.NewTerm(), disabledReason, 'b'),
			self.markMenuItem(markLabel(info.OldTerm()), commit.Hash(), info.OldTerm(), disabledReason, 'g'),
			self.skipMenuItem(
				utils.ResolvePlaceholderString(self.c.Tr.Bisect.SkipRef, map[string]string{"ref": ref}),
				commit.Hash(), types.KeepCommitSelectionByHash, disabledReason, 's'),
			self.resetMenuItem(),
		},
	})
}

// The menu for while bisecting. Its entries act on the commit that git wants the
// user to test, except for the one that skips the selected commit.
func (self *BisectController) openMidBisectMenu(info *git_commands.BisectInfo, commit *models.Commit) error {
	currentHash := info.GetCurrentHash()
	shortCurrentHash := utils.ShortHash(currentHash)

	menuItems := []*types.MenuItem{
		self.markMenuItem(
			fmt.Sprintf(self.c.Tr.Bisect.Mark, shortCurrentHash, info.NewTerm()),
			currentHash, info.NewTerm(), nil, 'b'),
		self.markMenuItem(
			fmt.Sprintf(self.c.Tr.Bisect.Mark, shortCurrentHash, info.OldTerm()),
			currentHash, info.OldTerm(), nil, 'g'),
		self.skipMenuItem(
			fmt.Sprintf(self.c.Tr.Bisect.SkipCurrent, shortCurrentHash),
			currentHash, types.KeepCommitSelectionByHash, nil, 's'),
	}
	if currentHash != commit.Hash() {
		menuItems = append(menuItems, self.skipMenuItem(
			fmt.Sprintf(self.c.Tr.Bisect.SkipSelected, commit.ShortHash()),
			commit.Hash(),
			// Git may pick a different commit to test now. Keep the skipped
			// commit selected anyway, so that the user can go on skipping
			// commits near it.
			types.KeepCommitSelectionIndex,
			self.require(self.singleItemSelected())(),
			'S'))
	}
	menuItems = append(menuItems, self.resetMenuItem())

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Bisect.BisectMenuTitle,
		Items: menuItems,
	})
}

func (self *BisectController) markMenuItem(label string, hash string, term string, disabledReason *types.DisabledReason, key rune) *types.MenuItem {
	return &types.MenuItem{
		Label: label,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.BisectMark)
			if err := self.c.Git().Bisect.Mark(hash, term); err != nil {
				return err
			}

			return self.afterMark(types.KeepCommitSelectionByHash)
		},
		DisabledReason: disabledReason,
		Keys:           menuKey(key),
	}
}

func (self *BisectController) skipMenuItem(label string, hash string, commitSelection types.CommitSelectionBehavior, disabledReason *types.DisabledReason, key rune) *types.MenuItem {
	return &types.MenuItem{
		Label: label,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.BisectSkip)
			if err := self.c.Git().Bisect.Skip(hash); err != nil {
				return err
			}

			return self.afterMark(commitSelection)
		},
		DisabledReason: disabledReason,
		Keys:           menuKey(key),
	}
}

func (self *BisectController) resetMenuItem() *types.MenuItem {
	return &types.MenuItem{
		Label: self.c.Tr.Bisect.ResetOption,
		OnPress: func() error {
			return self.c.Helpers().Bisect.Reset()
		},
		Keys: menuKey('r'),
	}
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
