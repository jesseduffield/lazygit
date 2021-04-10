package commands

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
)

// DeletePatchesFromCommit applies a patch in reverse for a commit
func (c *GitCommand) DeletePatchesFromCommit(commits []*models.Commit, commitIndex int, p *patch.PatchManager) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// apply each patch in reverse
	if err := p.ApplyPatches(true); err != nil {
		if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
			return err
		}
		return err
	}

	// time to amend the selected commit
	if err := c.AmendHead(); err != nil {
		return err
	}

	c.onSuccessfulContinue = func() error {
		c.PatchManager.Reset()
		return nil
	}

	// continue
	return c.GenericMergeOrRebaseAction("rebase", "continue")
}

func (c *GitCommand) MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int, p *patch.PatchManager) error {
	if sourceCommitIdx < destinationCommitIdx {
		if err := c.BeginInteractiveRebaseForCommit(commits, destinationCommitIdx); err != nil {
			return err
		}

		// apply each patch forward
		if err := p.ApplyPatches(false); err != nil {
			if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
				return err
			}
			return err
		}

		// amend the destination commit
		if err := c.AmendHead(); err != nil {
			return err
		}

		c.onSuccessfulContinue = func() error {
			c.PatchManager.Reset()
			return nil
		}

		// continue
		return c.GenericMergeOrRebaseAction("rebase", "continue")
	}

	if len(commits)-1 < sourceCommitIdx {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if c.UsingGpg() {
		return errors.New(c.Tr.DisabledForGPG)
	}

	baseIndex := sourceCommitIdx + 1
	todo := ""
	for i, commit := range commits[0:baseIndex] {
		a := "pick"
		if i == sourceCommitIdx || i == destinationCommitIdx {
			a = "edit"
		}
		todo = a + " " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(commits[baseIndex].Sha, todo, true)
	if err != nil {
		return err
	}

	if err := c.OSCommand.RunPreparedCommand(cmd); err != nil {
		return err
	}

	// apply each patch in reverse
	if err := p.ApplyPatches(true); err != nil {
		if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
			return err
		}
		return err
	}

	// amend the source commit
	if err := c.AmendHead(); err != nil {
		return err
	}

	if c.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	c.onSuccessfulContinue = func() error {
		// now we should be up to the destination, so let's apply forward these patches to that.
		// ideally we would ensure we're on the right commit but I'm not sure if that check is necessary
		if err := p.ApplyPatches(false); err != nil {
			if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
				return err
			}
			return err
		}

		// amend the destination commit
		if err := c.AmendHead(); err != nil {
			return err
		}

		c.onSuccessfulContinue = func() error {
			c.PatchManager.Reset()
			return nil
		}

		return c.GenericMergeOrRebaseAction("rebase", "continue")
	}

	return c.GenericMergeOrRebaseAction("rebase", "continue")
}

func (c *GitCommand) PullPatchIntoIndex(commits []*models.Commit, commitIdx int, p *patch.PatchManager, stash bool) error {
	if stash {
		if err := c.StashSave(c.Tr.StashPrefix + commits[commitIdx].Sha); err != nil {
			return err
		}
	}

	if err := c.BeginInteractiveRebaseForCommit(commits, commitIdx); err != nil {
		return err
	}

	if err := p.ApplyPatches(true); err != nil {
		if c.WorkingTreeState() == REBASE_MODE_REBASING {
			if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
				return err
			}
		}
		return err
	}

	// amend the commit
	if err := c.AmendHead(); err != nil {
		return err
	}

	if c.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	c.onSuccessfulContinue = func() error {
		// add patches to index
		if err := p.ApplyPatches(false); err != nil {
			if c.WorkingTreeState() == REBASE_MODE_REBASING {
				if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
					return err
				}
			}
			return err
		}

		if stash {
			if err := c.StashDo(0, "apply"); err != nil {
				return err
			}
		}

		c.PatchManager.Reset()
		return nil
	}

	return c.GenericMergeOrRebaseAction("rebase", "continue")
}

func (c *GitCommand) PullPatchIntoNewCommit(commits []*models.Commit, commitIdx int, p *patch.PatchManager) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIdx); err != nil {
		return err
	}

	if err := p.ApplyPatches(true); err != nil {
		if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
			return err
		}
		return err
	}

	// amend the commit
	if err := c.AmendHead(); err != nil {
		return err
	}

	// add patches to index
	if err := p.ApplyPatches(false); err != nil {
		if err := c.GenericMergeOrRebaseAction("rebase", "abort"); err != nil {
			return err
		}
		return err
	}

	head_message, _ := c.GetHeadCommitMessage()
	new_message := fmt.Sprintf("Split from \"%s\"", head_message)
	err := c.OSCommand.RunCommand(c.CommitCmdStr(new_message, ""))
	if err != nil {
		return err
	}

	if c.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	c.PatchManager.Reset()
	return c.GenericMergeOrRebaseAction("rebase", "continue")
}
