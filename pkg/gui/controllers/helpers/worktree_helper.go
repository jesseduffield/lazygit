package helpers

import (
	"errors"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type IWorktreeHelper interface {
	GetMainWorktreeName() string
	GetCurrentWorktreeName() string
}

type WorktreeHelper struct {
	c                 *HelperCommon
	reposHelper       *ReposHelper
	refsHelper        *RefsHelper
	suggestionsHelper *SuggestionsHelper
}

func NewWorktreeHelper(c *HelperCommon, reposHelper *ReposHelper, refsHelper *RefsHelper, suggestionsHelper *SuggestionsHelper) *WorktreeHelper {
	return &WorktreeHelper{
		c:                 c,
		reposHelper:       reposHelper,
		refsHelper:        refsHelper,
		suggestionsHelper: suggestionsHelper,
	}
}

func (self *WorktreeHelper) GetMainWorktreeName() string {
	for _, worktree := range self.c.Model().Worktrees {
		if worktree.IsMain {
			return worktree.Name
		}
	}

	return ""
}

// If we're on the main worktree, we return an empty string
func (self *WorktreeHelper) GetLinkedWorktreeName() string {
	worktrees := self.c.Model().Worktrees
	if len(worktrees) == 0 {
		return ""
	}

	// worktrees always have the current worktree on top
	currentWorktree := worktrees[0]
	if currentWorktree.IsMain {
		return ""
	}

	return currentWorktree.Name
}

func (self *WorktreeHelper) NewWorktree() error {
	branch := self.refsHelper.GetCheckedOutRef()
	currentBranchName := branch.RefName()

	f := func(detached bool) error {
		return self.c.Prompt(types.PromptOpts{
			Title:               self.c.Tr.NewWorktreeBase,
			InitialContent:      currentBranchName,
			FindSuggestionsFunc: self.suggestionsHelper.GetRefsSuggestionsFunc(),
			HandleConfirm: func(base string) error {
				// we assume that the base can be checked out
				canCheckoutBase := true
				return self.NewWorktreeCheckout(base, canCheckoutBase, detached, context.WORKTREES_CONTEXT_KEY)
			},
		})
	}

	placeholders := map[string]string{"ref": "ref"}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.WorktreeTitle,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{utils.ResolvePlaceholderString(self.c.Tr.CreateWorktreeFrom, placeholders)},
				OnPress: func() error {
					return f(false)
				},
			},
			{
				LabelColumns: []string{utils.ResolvePlaceholderString(self.c.Tr.CreateWorktreeFromDetached, placeholders)},
				OnPress: func() error {
					return f(true)
				},
			},
		},
	})
}

func (self *WorktreeHelper) NewWorktreeCheckout(base string, canCheckoutBase bool, detached bool, contextKey types.ContextKey) error {
	opts := git_commands.NewWorktreeOpts{
		Base:   base,
		Detach: detached,
	}

	f := func() error {
		return self.c.WithWaitingStatus(self.c.Tr.AddingWorktree, func(gocui.Task) error {
			self.c.LogAction(self.c.Tr.Actions.AddWorktree)
			if err := self.c.Git().Worktree.New(opts); err != nil {
				return err
			}

			return self.reposHelper.DispatchSwitchTo(opts.Path, self.c.Tr.ErrWorktreeMovedOrRemoved, contextKey)
		})
	}

	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.NewWorktreePath,
		HandleConfirm: func(path string) error {
			opts.Path = path

			if detached {
				return f()
			}

			if canCheckoutBase {
				title := utils.ResolvePlaceholderString(self.c.Tr.NewBranchNameLeaveBlank, map[string]string{"default": base})
				// prompt for the new branch name where a blank means we just check out the branch
				return self.c.Prompt(types.PromptOpts{
					Title: title,
					HandleConfirm: func(branchName string) error {
						opts.Branch = branchName

						return f()
					},
				})
			} else {
				// prompt for the new branch name where a blank means we just check out the branch
				return self.c.Prompt(types.PromptOpts{
					Title: self.c.Tr.NewBranchName,
					HandleConfirm: func(branchName string) error {
						if branchName == "" {
							return errors.New(self.c.Tr.BranchNameCannotBeBlank)
						}

						opts.Branch = branchName

						return f()
					},
				})
			}
		},
	})
}

func (self *WorktreeHelper) Switch(worktree *models.Worktree, contextKey types.ContextKey) error {
	if worktree.IsCurrent {
		return errors.New(self.c.Tr.AlreadyInWorktree)
	}

	self.c.LogAction(self.c.Tr.SwitchToWorktree)

	return self.reposHelper.DispatchSwitchTo(worktree.Path, self.c.Tr.ErrWorktreeMovedOrRemoved, contextKey)
}

func (self *WorktreeHelper) Remove(worktree *models.Worktree, force bool) error {
	title := self.c.Tr.RemoveWorktreeTitle
	var templateStr string
	if force {
		templateStr = self.c.Tr.ForceRemoveWorktreePrompt
	} else {
		templateStr = self.c.Tr.RemoveWorktreePrompt
	}
	message := utils.ResolvePlaceholderString(
		templateStr,
		map[string]string{
			"worktreeName": worktree.Name,
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.RemovingWorktree, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.RemoveWorktree)
				if err := self.c.Git().Worktree.Delete(worktree.Path, force); err != nil {
					errMessage := err.Error()
					if !strings.Contains(errMessage, "--force") {
						return err
					}

					if !force {
						return self.Remove(worktree, true)
					}
					return err
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.WORKTREES, types.BRANCHES, types.FILES}})
			})
		},
	})
}

func (self *WorktreeHelper) Detach(worktree *models.Worktree) error {
	return self.c.WithWaitingStatus(self.c.Tr.DetachingWorktree, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.RemovingWorktree)

		err := self.c.Git().Worktree.Detach(worktree.Path)
		if err != nil {
			return err
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.WORKTREES, types.BRANCHES, types.FILES}})
	})
}

func (self *WorktreeHelper) ViewWorktreeOptions(context types.IListContext, ref string) error {
	currentBranch := self.refsHelper.GetCheckedOutRef()
	canCheckoutBase := context == self.c.Contexts().Branches && ref != currentBranch.RefName()

	return self.ViewBranchWorktreeOptions(ref, canCheckoutBase)
}

func (self *WorktreeHelper) ViewBranchWorktreeOptions(branchName string, canCheckoutBase bool) error {
	placeholders := map[string]string{"ref": branchName}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.WorktreeTitle,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{utils.ResolvePlaceholderString(self.c.Tr.CreateWorktreeFrom, placeholders)},
				OnPress: func() error {
					return self.NewWorktreeCheckout(branchName, canCheckoutBase, false, context.LOCAL_BRANCHES_CONTEXT_KEY)
				},
			},
			{
				LabelColumns: []string{utils.ResolvePlaceholderString(self.c.Tr.CreateWorktreeFromDetached, placeholders)},
				OnPress: func() error {
					return self.NewWorktreeCheckout(branchName, canCheckoutBase, true, context.LOCAL_BRANCHES_CONTEXT_KEY)
				},
			},
		},
	})
}
