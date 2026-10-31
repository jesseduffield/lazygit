package controllers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type CanBisect interface {
	types.IListContext
	GetSelectedRef() models.Ref
}

// Offers the bisect menu for the selected item. The commits panel is where the
// user bisects: it shows the good, bad and skipped commits, and the one to test
// next. The other views only offer marking their selected item as good or bad,
// until git knows both a good and a bad commit.
type BisectController struct {
	baseController
	*ListControllerTrait[models.Ref]
	c       *ControllerCommon
	context CanBisect
}

var _ types.IController = &BisectController{}

func NewBisectController(
	c *ControllerCommon,
	context CanBisect,
) *BisectController {
	return &BisectController{
		baseController: baseController{},
		c:              c,
		context:        context,
		ListControllerTrait: NewListControllerTrait(
			c,
			context,
			context.GetSelectedRef,
			func() ([]models.Ref, int, int) {
				panic("Not implemented")
			},
		),
	}
}

func (self *BisectController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:              opts.GetKeys(opts.Config.Commits.ViewBisectOptions),
			Handler:           opts.Guards.OutsideFilterMode(self.withItem(self.openMenu)),
			GetDisabledReason: self.notBisectingOutsideCommitsPanel,
			Description:       self.c.Tr.ViewBisectOptions,
			OpensMenu:         true,
		},
	}

	return bindings
}

func (self *BisectController) isCommitsPanel() bool {
	return self.context.GetKey() == context.LOCAL_COMMITS_CONTEXT_KEY
}

func (self *BisectController) notBisectingOutsideCommitsPanel() *types.DisabledReason {
	if !self.isCommitsPanel() && self.c.Model().BisectInfo.Bisecting() {
		return &types.DisabledReason{Text: self.c.Tr.Bisect.UseCommitsPanel}
	}

	return nil
}

func (self *BisectController) openMenu(ref models.Ref) error {
	// no shame in getting this directly rather than using the cached value
	// given how cheap it is to obtain
	info := self.c.Git().Bisect.GetInfo()
	if !info.Started() {
		return self.openStartBisectMenu(info, ref)
	}
	if !info.Bisecting() {
		return self.openMenuBeforeBisecting(info, ref)
	}
	if !info.HasCommitToTest() {
		return self.showBisectResult()
	}
	return self.openMidBisectMenu(info, ref)
}

// Git has nothing left for the user to test, so show the result, and offer to
// reset the bisect.
func (self *BisectController) showBisectResult() error {
	done, candidateHashes, err := self.c.Git().Bisect.IsDone()
	if err != nil {
		return err
	}

	if !done {
		// Git stopped without finding the first bad commit, for example
		// because the merge base of the good and the bad commit is bad. Git
		// said why when it stopped, so all that's left to do is to reset.
		return self.c.Helpers().Bisect.Reset()
	}

	return self.showBisectCompleteMessage(candidateHashes)
}

// The menu for when a bisect has started, but git doesn't know both a good and
// a bad commit yet. Its entries act on the selected item.
func (self *BisectController) openMenuBeforeBisecting(info *git_commands.BisectInfo, ref models.Ref) error {
	disabledReason := self.require(self.singleItemSelected())()
	refDescription := self.describeRef(ref)
	markLabel := func(term string) string {
		return utils.ResolvePlaceholderString(self.c.Tr.Bisect.MarkRef,
			map[string]string{"ref": refDescription, "term": term})
	}

	menuItems := []*types.MenuItem{
		self.markMenuItem(markLabel(info.NewTerm()), ref.FullRefName(), info.NewTerm(), disabledReason, 'b'),
		self.markMenuItem(markLabel(info.OldTerm()), ref.FullRefName(), info.OldTerm(), disabledReason, 'g'),
	}
	if self.isCommitsPanel() {
		menuItems = append(menuItems, self.skipMenuItem(
			utils.ResolvePlaceholderString(self.c.Tr.Bisect.SkipRef, map[string]string{"ref": refDescription}),
			ref.FullRefName(), types.KeepCommitSelectionByHash, disabledReason, 's'))
	}
	menuItems = append(menuItems, self.resetMenuItem())

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Bisect.BisectMenuTitle,
		Items: menuItems,
	})
}

// The menu for while bisecting. Its entries act on the commit that git wants the
// user to test, except for the one that skips the selected commit.
func (self *BisectController) openMidBisectMenu(info *git_commands.BisectInfo, ref models.Ref) error {
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
	if currentHash != ref.FullRefName() {
		menuItems = append(menuItems, self.skipMenuItem(
			fmt.Sprintf(self.c.Tr.Bisect.SkipSelected, utils.ShortHash(ref.FullRefName())),
			ref.FullRefName(),
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

// Returns how the bisect menu refers to the given item.
func (self *BisectController) describeRef(ref models.Ref) string {
	switch ref := ref.(type) {
	case *models.Branch:
		// The branches view shows a detached head as a branch, too
		if !ref.DetachedHead {
			return utils.ResolvePlaceholderString(self.c.Tr.Bisect.BranchRef,
				map[string]string{"name": ref.RefName()})
		}
	case *models.RemoteBranch:
		return utils.ResolvePlaceholderString(self.c.Tr.Bisect.BranchRef,
			map[string]string{"name": ref.RefName()})
	case *models.Tag:
		return utils.ResolvePlaceholderString(self.c.Tr.Bisect.TagRef,
			map[string]string{"name": ref.RefName()})
	}

	return utils.ResolvePlaceholderString(self.c.Tr.Bisect.CommitRef,
		map[string]string{"hash": utils.ShortHash(ref.FullRefName())})
}

func (self *BisectController) markMenuItem(label string, ref string, term string, disabledReason *types.DisabledReason, key rune) *types.MenuItem {
	return &types.MenuItem{
		Label: label,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.BisectMark)
			if err := self.c.Git().Bisect.Mark(ref, term); err != nil {
				return err
			}

			if !self.isCommitsPanel() && self.c.Git().Bisect.GetInfo().Bisecting() {
				// Git knows both a good and a bad commit now, so the user goes
				// on in the commits panel
				self.c.Context().Push(self.c.Contexts().LocalCommits, types.OnFocusOpts{})
			}

			return self.afterMark(types.KeepCommitSelectionByHash)
		},
		DisabledReason: disabledReason,
		Keys:           menuKey(key),
	}
}

func (self *BisectController) skipMenuItem(label string, ref string, commitSelection types.CommitSelectionBehavior, disabledReason *types.DisabledReason, key rune) *types.MenuItem {
	return &types.MenuItem{
		Label: label,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.BisectSkip)
			if err := self.c.Git().Bisect.Skip(ref); err != nil {
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

func (self *BisectController) openStartBisectMenu(info *git_commands.BisectInfo, ref models.Ref) error {
	refDescription := self.describeRef(ref)
	markLabel := func(term string) string {
		return utils.ResolvePlaceholderString(self.c.Tr.Bisect.MarkRefStart,
			map[string]string{"ref": refDescription, "term": term})
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Bisect.BisectMenuTitle,
		Items: []*types.MenuItem{
			{
				Label: markLabel(info.NewTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.c.Git().Bisect.Start(); err != nil {
						return err
					}

					if err := self.c.Git().Bisect.Mark(ref.FullRefName(), info.NewTerm()); err != nil {
						return err
					}

					self.c.Helpers().Bisect.PostBisectCommandRefresh()
					return nil
				},
				DisabledReason: self.require(self.singleItemSelected())(),
				Keys:           menuKey('b'),
			},
			{
				Label: markLabel(info.OldTerm()),
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StartBisect)
					if err := self.c.Git().Bisect.Start(); err != nil {
						return err
					}

					if err := self.c.Git().Bisect.Mark(ref.FullRefName(), info.OldTerm()); err != nil {
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
