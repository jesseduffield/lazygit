package helpers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

func (self *CherryPickHelper) Copy(commit *models.Commit, commitsList []*models.Commit, context types.Context) error {
	if err := self.resetIfNecessary(context); err != nil {
		return err
	}

	// we will un-copy it if it's already copied
	if self.getData().SelectedShaSet().Includes(commit.Sha) {
		self.getData().Remove(commit, commitsList)
	} else {
		self.getData().Add(commit, commitsList)
	}

	return self.rerender()
}

func (self *CherryPickHelper) CopyRange(selectedIndex int, commitsList []*models.Commit, context types.Context) error {
	if err := self.resetIfNecessary(context); err != nil {
		return err
	}

	commitSet := self.getData().SelectedShaSet()

	// find the last commit that is copied that's above our position
	// if there are none, startIndex = 0
	startIndex := 0
	for index, commit := range commitsList[0:selectedIndex] {
		if commitSet.Includes(commit.Sha) {
			startIndex = index
		}
	}

	for index := startIndex; index <= selectedIndex; index++ {
		commit := commitsList[index]
		self.getData().Add(commit, commitsList)
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
			return self.c.WithWaitingStatus(self.c.Tr.CherryPickingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.CherryPick)
				err := self.c.Git().Rebase.CherryPickCommits(self.getData().CherryPickedCommits)
				return self.rebaseHelper.CheckMergeOrRebase(err)
			})
		},
	})
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
