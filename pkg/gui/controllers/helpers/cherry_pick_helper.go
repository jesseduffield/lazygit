package helpers

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CherryPickHelper struct {
	c *types.HelperCommon

	git *commands.GitCommand

	contexts *context.ContextTree
	getData  func() *cherrypicking.CherryPicking

	rebaseHelper *MergeAndRebaseHelper
}

// I'm using the analogy of copy+paste in the terminology here because it's intuitively what's going on,
// even if in truth we're running git cherry-pick

func NewCherryPickHelper(
	c *types.HelperCommon,
	git *commands.GitCommand,
	contexts *context.ContextTree,
	getData func() *cherrypicking.CherryPicking,
	rebaseHelper *MergeAndRebaseHelper,
) *CherryPickHelper {
	return &CherryPickHelper{
		c:            c,
		git:          git,
		contexts:     contexts,
		getData:      getData,
		rebaseHelper: rebaseHelper,
	}
}

func (self *CherryPickHelper) Copy(commit *models.Commit, commitsList []*models.Commit, context types.Context) error {
	if err := self.resetIfNecessary(context); err != nil {
		return err
	}

	// we will un-copy it if it's already copied
	for index, cherryPickedCommit := range self.getData().CherryPickedCommits {
		if commit.Sha == cherryPickedCommit.Sha {
			self.getData().CherryPickedCommits = append(
				self.getData().CherryPickedCommits[0:index],
				self.getData().CherryPickedCommits[index+1:]...,
			)
			return self.rerender()
		}
	}

	self.add(commit, commitsList)
	return self.rerender()
}

func (self *CherryPickHelper) CopyRange(selectedIndex int, commitsList []*models.Commit, context types.Context) error {
	if err := self.resetIfNecessary(context); err != nil {
		return err
	}

	commitSet := self.CherryPickedCommitShaSet()

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
		self.add(commit, commitsList)
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
			return self.c.WithWaitingStatus(self.c.Tr.CherryPickingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.CherryPick)
				err := self.git.Rebase.CherryPickCommits(self.getData().CherryPickedCommits)
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

func (self *CherryPickHelper) CherryPickedCommitShaSet() *set.Set[string] {
	shas := slices.Map(self.getData().CherryPickedCommits, func(commit *models.Commit) string {
		return commit.Sha
	})
	return set.NewFromSlice(shas)
}

func (self *CherryPickHelper) add(selectedCommit *models.Commit, commitsList []*models.Commit) {
	commitSet := self.CherryPickedCommitShaSet()
	commitSet.Add(selectedCommit.Sha)

	cherryPickedCommits := slices.Filter(commitsList, func(commit *models.Commit) bool {
		return commitSet.Includes(commit.Sha)
	})

	self.getData().CherryPickedCommits = slices.Map(cherryPickedCommits, func(commit *models.Commit) *models.Commit {
		return &models.Commit{Name: commit.Name, Sha: commit.Sha}
	})
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
		self.contexts.LocalCommits,
		self.contexts.ReflogCommits,
		self.contexts.SubCommits,
	} {
		if err := self.c.PostRefreshUpdate(context); err != nil {
			return err
		}
	}

	return nil
}
