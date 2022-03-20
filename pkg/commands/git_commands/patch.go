package git_commands

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
)

type PatchCommands struct {
	*GitCommon
	rebase *RebaseCommands
	commit *CommitCommands
	status *StatusCommands
	stash  *StashCommands

	PatchManager *patch.PatchManager
}

func NewPatchCommands(
	gitCommon *GitCommon,
	rebase *RebaseCommands,
	commit *CommitCommands,
	status *StatusCommands,
	stash *StashCommands,
	patchManager *patch.PatchManager,
) *PatchCommands {
	return &PatchCommands{
		GitCommon:    gitCommon,
		rebase:       rebase,
		commit:       commit,
		status:       status,
		stash:        stash,
		PatchManager: patchManager,
	}
}

// DeletePatchesFromCommit applies a patch in reverse for a commit
func (self *PatchCommands) DeletePatchesFromCommit(commits []*models.Commit, commitIndex int) error {
	if err := self.rebase.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// apply each patch in reverse
	if err := self.PatchManager.ApplyPatches(true); err != nil {
		if err := self.rebase.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	// time to amend the selected commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	self.rebase.onSuccessfulContinue = func() error {
		self.PatchManager.Reset()
		return nil
	}

	// continue
	return self.rebase.ContinueRebase()
}

func (self *PatchCommands) MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int) error {
	if sourceCommitIdx < destinationCommitIdx {
		if err := self.rebase.BeginInteractiveRebaseForCommit(commits, destinationCommitIdx); err != nil {
			return err
		}

		// apply each patch forward
		if err := self.PatchManager.ApplyPatches(false); err != nil {
			if err := self.rebase.AbortRebase(); err != nil {
				return err
			}
			return err
		}

		// amend the destination commit
		if err := self.commit.AmendHead(); err != nil {
			return err
		}

		self.rebase.onSuccessfulContinue = func() error {
			self.PatchManager.Reset()
			return nil
		}

		// continue
		return self.rebase.ContinueRebase()
	}

	if len(commits)-1 < sourceCommitIdx {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if self.config.UsingGpg() {
		return errors.New(self.Tr.DisabledForGPG)
	}

	baseIndex := sourceCommitIdx + 1

	todoLines := self.rebase.BuildTodoLines(commits[0:baseIndex], func(commit *models.Commit, i int) string {
		if i == sourceCommitIdx || i == destinationCommitIdx {
			return "edit"
		} else {
			return "pick"
		}
	})

	err := self.rebase.PrepareInteractiveRebaseCommand(commits[baseIndex].Sha, todoLines, true).Run()
	if err != nil {
		return err
	}

	// apply each patch in reverse
	if err := self.PatchManager.ApplyPatches(true); err != nil {
		if err := self.rebase.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	// amend the source commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	if self.rebase.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	self.rebase.onSuccessfulContinue = func() error {
		// now we should be up to the destination, so let's apply forward these patches to that.
		// ideally we would ensure we're on the right commit but I'm not sure if that check is necessary
		if err := self.PatchManager.ApplyPatches(false); err != nil {
			if err := self.rebase.AbortRebase(); err != nil {
				return err
			}
			return err
		}

		// amend the destination commit
		if err := self.commit.AmendHead(); err != nil {
			return err
		}

		self.rebase.onSuccessfulContinue = func() error {
			self.PatchManager.Reset()
			return nil
		}

		return self.rebase.ContinueRebase()
	}

	return self.rebase.ContinueRebase()
}

func (self *PatchCommands) MovePatchIntoIndex(commits []*models.Commit, commitIdx int, stash bool) error {
	if stash {
		if err := self.stash.Save(self.Tr.StashPrefix + commits[commitIdx].Sha); err != nil {
			return err
		}
	}

	if err := self.rebase.BeginInteractiveRebaseForCommit(commits, commitIdx); err != nil {
		return err
	}

	if err := self.PatchManager.ApplyPatches(true); err != nil {
		if self.status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
			if err := self.rebase.AbortRebase(); err != nil {
				return err
			}
		}
		return err
	}

	// amend the commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	if self.rebase.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	self.rebase.onSuccessfulContinue = func() error {
		// add patches to index
		if err := self.PatchManager.ApplyPatches(false); err != nil {
			if self.status.WorkingTreeState() == enums.REBASE_MODE_REBASING {
				if err := self.rebase.AbortRebase(); err != nil {
					return err
				}
			}
			return err
		}

		if stash {
			if err := self.stash.Apply(0); err != nil {
				return err
			}
		}

		self.PatchManager.Reset()
		return nil
	}

	return self.rebase.ContinueRebase()
}

func (self *PatchCommands) PullPatchIntoNewCommit(commits []*models.Commit, commitIdx int) error {
	if err := self.rebase.BeginInteractiveRebaseForCommit(commits, commitIdx); err != nil {
		return err
	}

	if err := self.PatchManager.ApplyPatches(true); err != nil {
		if err := self.rebase.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	// amend the commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	// add patches to index
	if err := self.PatchManager.ApplyPatches(false); err != nil {
		if err := self.rebase.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	head_message, _ := self.commit.GetHeadCommitMessage()
	new_message := fmt.Sprintf("Split from \"%s\"", head_message)
	err := self.commit.CommitCmdObj(new_message).Run()
	if err != nil {
		return err
	}

	if self.rebase.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	self.PatchManager.Reset()
	return self.rebase.ContinueRebase()
}
