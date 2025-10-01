package helpers

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/modes/cherrypicking"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type CherryPickHelper struct {
	c *HelperCommon

	rebaseHelper *MergeAndRebaseHelper

	postPasteCleanup            func(markDidPaste bool) error
	postPasteShouldMarkDidPaste bool

	postPasteSelection *postPasteSelection

	prePasteHeadHash     string
	pasteProducedCommits bool

	deferPostPasteCleanup bool
}

type postPasteSelection struct {
	hash           string
	idx            int
	shouldReselect bool
}

// I'm using the analogy of copy+paste in the terminology here because it's intuitively what's going on,
// even if in truth we're running git cherry-pick

func NewCherryPickHelper(
	c *HelperCommon,
	rebaseHelper *MergeAndRebaseHelper,
) *CherryPickHelper {
	helper := &CherryPickHelper{
		c:            c,
		rebaseHelper: rebaseHelper,
	}

	rebaseHelper.SetCherryPickHelper(helper)

	return helper
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
			return self.c.WithWaitingStatusSync(self.c.Tr.CherryPickingStatus, func() error {
				mustStash := IsWorkingTreeDirty(self.c.Model().Files)

				self.c.LogAction(self.c.Tr.Actions.CherryPick)

				if mustStash {
					if err := self.c.Git().Stash.Push(self.c.Tr.AutoStashForCherryPicking); err != nil {
						return err
					}
				}

				cherryPickedCommits := self.getData().CherryPickedCommits
				self.preparePostPasteSelection(cherryPickedCommits)
				self.capturePrePasteHeadHash()

				self.setPostPasteCleanup(func(markDidPaste bool) error {
					if markDidPaste {
						self.getData().DidPaste = true
					}
					self.rerender()

					if mustStash {
						if err := self.c.Git().Stash.Pop(0); err != nil {
							return err
						}
						self.c.Refresh(types.RefreshOptions{
							Scope: []types.RefreshableView{types.STASH, types.FILES},
						})
					}

					self.restorePostPasteSelection()

					return nil
				})

				result := self.c.Git().Rebase.CherryPickCommits(cherryPickedCommits)
				err := self.rebaseHelper.CheckMergeOrRebaseWithRefreshOptions(result, types.RefreshOptions{Mode: types.SYNC})
				if err != nil {
					return err
				}

				self.markPasteProducedCommitsIfHeadChanged()

				// If we're in the cherry-picking state at this point, it must
				// be because there were conflicts. Don't clear the copied
				// commits in this case, since we might want to abort and try
				// pasting them again.
				isInCherryPick, result := self.c.Git().Status.IsInCherryPick()
				if result != nil {
					return result
				}
				if !isInCherryPick {
					if err := self.runPostPasteCleanup(true); err != nil {
						return err
					}
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
	self.postPasteCleanup = nil
	self.postPasteShouldMarkDidPaste = false
	self.postPasteSelection = nil
	self.prePasteHeadHash = ""
	self.pasteProducedCommits = false
	self.deferPostPasteCleanup = false

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

func (self *CherryPickHelper) preparePostPasteSelection(commits []*models.Commit) {
	selectedCommit := self.c.Contexts().LocalCommits.GetSelected()
	selectedIdx := self.c.Contexts().LocalCommits.GetSelectedLineIdx()

	self.postPasteSelection = nil

	if selectedCommit == nil {
		return
	}

	self.postPasteSelection = &postPasteSelection{
		hash:           selectedCommit.Hash(),
		idx:            selectedIdx,
		shouldReselect: !selectedCommit.IsTODO() && len(commits) > 0,
	}
}

func (self *CherryPickHelper) restorePostPasteSelection() {
	if self.postPasteSelection == nil || !self.postPasteSelection.shouldReselect {
		return
	}

	localCommits := self.c.Contexts().LocalCommits

	if self.postPasteSelection.hash != "" && localCommits.SelectCommitByHash(self.postPasteSelection.hash) {
		localCommits.FocusLine()
		return
	}

	if self.postPasteSelection.idx >= 0 {
		localCommits.SetSelectedLineIdx(self.postPasteSelection.idx)
		localCommits.FocusLine()
	}
}

func (self *CherryPickHelper) DisablePostPasteReselect() {
	if self.postPasteSelection == nil {
		return
	}

	self.postPasteSelection.shouldReselect = false
}

func (self *CherryPickHelper) ShouldRestorePostPasteSelection() bool {
	return self.postPasteSelection != nil && self.postPasteSelection.shouldReselect
}

func (self *CherryPickHelper) setPostPasteCleanup(cleanup func(markDidPaste bool) error) {
	self.postPasteCleanup = cleanup
	self.postPasteShouldMarkDidPaste = true
}

func (self *CherryPickHelper) runPostPasteCleanup(markDidPaste bool) error {
	if self.postPasteCleanup == nil {
		return nil
	}

	cleanup := self.postPasteCleanup
	self.postPasteCleanup = nil
	defer func() {
		self.postPasteShouldMarkDidPaste = false
		self.postPasteSelection = nil
		self.prePasteHeadHash = ""
		self.pasteProducedCommits = false
		self.deferPostPasteCleanup = false
	}()

	return cleanup(markDidPaste && self.postPasteShouldMarkDidPaste)
}

func (self *CherryPickHelper) setPostPasteShouldMarkDidPaste(mark bool) {
	if self.postPasteCleanup == nil {
		return
	}

	self.postPasteShouldMarkDidPaste = mark
}

func (self *CherryPickHelper) capturePrePasteHeadHash() {
	headHash, err := self.getHeadHash()
	if err != nil {
		self.prePasteHeadHash = ""
		return
	}

	self.prePasteHeadHash = headHash
	self.pasteProducedCommits = false
}

func (self *CherryPickHelper) PasteProducedCommits() bool {
	if self.pasteProducedCommits {
		return true
	}

	if self.prePasteHeadHash == "" {
		return true
	}

	headHash, err := self.getHeadHash()
	if err != nil {
		return true
	}

	if headHash != self.prePasteHeadHash {
		self.pasteProducedCommits = true
		return true
	}

	return false
}

func (self *CherryPickHelper) DeferPostPasteCleanup() {
	self.deferPostPasteCleanup = true
}

func (self *CherryPickHelper) ClearDeferredPostPasteCleanup() {
	self.deferPostPasteCleanup = false
}

func (self *CherryPickHelper) ShouldDeferPostPasteCleanup() bool {
	return self.deferPostPasteCleanup
}

func (self *CherryPickHelper) markPasteProducedCommitsIfHeadChanged() {
	if self.pasteProducedCommits {
		return
	}

	if self.prePasteHeadHash == "" {
		self.pasteProducedCommits = true
		return
	}

	headHash, err := self.getHeadHash()
	if err != nil {
		self.pasteProducedCommits = true
		return
	}

	if headHash != self.prePasteHeadHash {
		self.pasteProducedCommits = true
	}
}

func (self *CherryPickHelper) getHeadHash() (string, error) {
	output, err := self.c.OS().Cmd.New(
		git_commands.NewGitCmd("rev-parse").Arg("HEAD").ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}
