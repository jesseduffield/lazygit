package helpers

import (
	"strconv"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
		return commitSet.Includes(commit.Hash())
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

	self.getData().DidPaste = false

	self.rerender()
	return nil
}

// HandlePasteCommits begins a cherry-pick rebase with the commits the user has copied.
// Only to be called from the branch commits controller
func (self *CherryPickHelper) Paste() error {
	self.c.Confirm(types.ConfirmOpts{
		Title: self.c.Tr.CherryPick,
		Prompt: utils.ResolvePlaceholderString(
			self.c.Tr.SureCherryPick,
			map[string]string{
				"numCommits": strconv.Itoa(len(self.getData().CherryPickedCommits)),
			}),
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.CherryPickingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.CherryPick)
				result := self.c.Git().Rebase.CherryPickCommits(self.getData().CherryPickedCommits)
				err := self.rebaseHelper.CheckMergeOrRebase(result)
				if err != nil {
					return result
				}

				// If we're in the cherry-picking state at this point, it must
				// be because there were conflicts. Don't clear the copied
				// commits in this case, since we might want to abort and try
				// pasting them again.
				isInCherryPick, result := self.c.Git().Status.IsInCherryPick()
				if result != nil {
					return result
				}
				if !isInCherryPick {
					self.getData().DidPaste = true
					self.rerender()
				}
				return nil
			})
		},
	})

	return nil
}

func (self *CherryPickHelper) CanPaste() bool {
	return self.getData().CanPaste()
}

func (self *CherryPickHelper) Reset() error {
	self.getData().ContextKey = ""
	self.getData().CherryPickedCommits = nil

	self.rerender()
	return nil
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

func (self *CherryPickHelper) rerender() {
	for _, context := range []types.Context{
		self.c.Contexts().LocalCommits,
		self.c.Contexts().ReflogCommits,
		self.c.Contexts().SubCommits,
	} {
		self.c.PostRefreshUpdate(context)
	}
}
