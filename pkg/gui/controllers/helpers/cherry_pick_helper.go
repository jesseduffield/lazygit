package helpers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type CherryPickHelper struct {
	c *HelperCommon

	rebaseHelper *MergeAndRebaseHelper
}

// I'm using the analogy of copy+paste in the terminology here because it's intuitively what's going on,
// even if in truth we're running git cherry-pick

func NewCherryPickHelper(
	c *HelperCommon,
	rebaseHelper *MergeAndRebaseHelper,
) *CherryPickHelper {
	return &CherryPickHelper{
		c:            c,
		rebaseHelper: rebaseHelper,
	}
}

func (self *CherryPickHelper) getData() *cherrypicking.CherryPicking {
	return self.c.Modes().CherryPicking
}

func (self *CherryPickHelper) CopyRange(commitsList []*models.Commit, context types.IListContext) error {
	startIdx, endIdx := context.GetList().GetSelectionRange()

	if err := self.resetIfNecessary(context); err != nil {
		return err
	}

	commitSet := self.getData().SelectedHashSet()

	allCommitsCopied := lo.EveryBy(commitsList[startIdx:endIdx+1], func(commit *models.Commit) bool {
		return commitSet.Includes(commit.Hash)
	})

	// if all selected commits are already copied, we'll uncopy them
	if allCommitsCopied {
		for index := startIdx; index <= endIdx; index++ {
			commit := commitsList[index]
			self.getData().Remove(commit, commitsList)
		}
	} else {
		for index := startIdx; index <= endIdx; index++ {
			commit := commitsList[index]
			self.getData().Add(commit, commitsList)
		}
	}

	return self.rerender()
}

// HandlePasteCommits begins a cherry-pick rebase with the commits the user has copied.
// Only to be called from the branch commits controller
func (self *CherryPickHelper) Paste() error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.CherryPick,
		Prompt: self.c.Tr.SureCherryPick,
		HandleConfirm: func() error {
			isInRebase, err := self.c.Git().Status.IsInInteractiveRebase()
			if err != nil {
				return err
			}
			if isInRebase {
				if err := self.c.Git().Rebase.CherryPickCommitsDuringRebase(self.getData().CherryPickedCommits); err != nil {
					return err
				}
				err = self.c.Refresh(types.RefreshOptions{
					Mode: types.SYNC, Scope: []types.RefreshableView{types.REBASE_COMMITS},
				})
				if err != nil {
					return err
				}

				return self.Reset()
			}

			return self.c.WithWaitingStatus(self.c.Tr.CherryPickingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.CherryPick)
				err := self.c.Git().Rebase.CherryPickCommits(self.getData().CherryPickedCommits)
				err = self.rebaseHelper.CheckMergeOrRebase(err)
				if err != nil {
					return err
				}

				// If we're in an interactive rebase at this point, it must
				// be because there were conflicts. Don't clear the copied
				// commits in this case, since we might want to abort and
				// try pasting them again.
				isInRebase, err = self.c.Git().Status.IsInInteractiveRebase()
				if err != nil {
					return err
				}
				if !isInRebase {
					return self.Reset()
				}
				return nil
			})
		},
	})
}

func (self *CherryPickHelper) CanPaste() bool {
	return self.getData().Active()
}

func (self *CherryPickHelper) Reset() error {
	self.getData().ContextKey = ""
	self.getData().CherryPickedCommits = nil

	return self.rerender()
}

// you can only copy from one context at a time, because the order and position of commits matter
func (self *CherryPickHelper) resetIfNecessary(context types.Context) error {
	oldContextKey := types.ContextKey(self.getData().ContextKey)

	if oldContextKey != context.GetKey() {
		// need to reset the cherry picking mode
		self.getData().ContextKey = string(context.GetKey())
		self.getData().CherryPickedCommits = make([]*models.Commit, 0)
	}

	return nil
}

func (self *CherryPickHelper) rerender() error {
	for _, context := range []types.Context{
		self.c.Contexts().LocalCommits,
		self.c.Contexts().ReflogCommits,
		self.c.Contexts().SubCommits,
	} {
		if err := self.c.PostRefreshUpdate(context); err != nil {
			return err
		}
	}

	return nil
}
