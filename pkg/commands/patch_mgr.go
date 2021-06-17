package commands

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

//counterfeiter:generate . IPatchesMgr
type IPatchesMgr interface {
	NewPatchManager() *patch.PatchManager
	DeletePatchesFromCommit(commits []*models.Commit, commitIndex int, p *patch.PatchManager) error
	MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int, p *patch.PatchManager) error
	MovePatchIntoIndex(commits []*models.Commit, commitIdx int, p *patch.PatchManager, stash bool) error
	PullPatchIntoNewCommit(commits []*models.Commit, commitIdx int, p *patch.PatchManager) error
	ApplyPatch(patch string, flags ...string) error
}

type PatchesMgr struct {
	IRebasingMgr
	ICommander

	commitsMgr ICommitsMgr
	statusMgr  IStatusMgr
	stashMgr   IStashMgr
	diffMgr    IDiffMgr

	config IGitConfigMgr
	tr     *i18n.TranslationSet
	log    *logrus.Entry
	os     oscommands.IOS
}

func NewPatchesMgr(
	commander ICommander,
	config IGitConfigMgr,
	commitsMgr ICommitsMgr,
	rebasingMgr IRebasingMgr,
	statusMgr IStatusMgr,
	stashMgr IStashMgr,
	diffMgr IDiffMgr,
	tr *i18n.TranslationSet,
	log *logrus.Entry,
	os oscommands.IOS,
) *PatchesMgr {
	return &PatchesMgr{
		ICommander:   commander,
		config:       config,
		commitsMgr:   commitsMgr,
		IRebasingMgr: rebasingMgr,
		tr:           tr,
		log:          log,
		os:           os,
		statusMgr:    statusMgr,
		stashMgr:     stashMgr,
		diffMgr:      diffMgr,
	}
}

// DeletePatchesFromCommit applies a patch in reverse for a commit
func (c *PatchesMgr) DeletePatchesFromCommit(commits []*models.Commit, commitIndex int, p *patch.PatchManager) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// apply each patch in reverse
	if err := p.ApplyPatches(c.ApplyPatch, true); err != nil {
		if err := c.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	// time to amend the selected commit
	if err := c.commitsMgr.AmendHead(); err != nil {
		return err
	}

	c.getWorkflow().Start(func() error {
		p.Reset()
		return nil
	})

	// continue
	return c.ContinueRebase()
}

func (c *PatchesMgr) MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int, p *patch.PatchManager) error {
	if sourceCommitIdx < destinationCommitIdx {
		if err := c.BeginInteractiveRebaseForCommit(commits, destinationCommitIdx); err != nil {
			return err
		}

		// apply each patch forward
		if err := p.ApplyPatches(c.ApplyPatch, false); err != nil {
			if err := c.AbortRebase(); err != nil {
				return err
			}
			return err
		}

		// amend the destination commit
		if err := c.commitsMgr.AmendHead(); err != nil {
			return err
		}

		c.getWorkflow().Start(func() error {
			p.Reset()
			return nil
		})

		// continue
		return c.ContinueRebase()
	}

	if len(commits)-1 < sourceCommitIdx {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if c.config.UsingGpg() {
		return errors.New(c.tr.DisabledForGPG)
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

	if err := c.Run(
		c.InteractiveRebaseCmdObj(commits[baseIndex].Sha, todo, true),
	); err != nil {
		return err
	}

	// apply each patch in reverse
	if err := p.ApplyPatches(c.ApplyPatch, true); err != nil {
		if err := c.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	// amend the source commit
	if err := c.commitsMgr.AmendHead(); err != nil {
		return err
	}

	if c.getWorkflow().InProgress() {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	c.getWorkflow().Start(func() error {
		// now we should be up to the destination, so let's apply forward these patches to that.
		// ideally we would ensure we're on the right commit but I'm not sure if that check is necessary
		if err := p.ApplyPatches(c.ApplyPatch, false); err != nil {
			if err := c.AbortRebase(); err != nil {
				return err
			}
			return err
		}

		// amend the destination commit
		if err := c.commitsMgr.AmendHead(); err != nil {
			return err
		}

		c.getWorkflow().Start(func() error {
			p.Reset()
			return nil
		})

		return c.ContinueRebase()
	})

	return c.ContinueRebase()
}

func (c *PatchesMgr) MovePatchIntoIndex(commits []*models.Commit, commitIdx int, p *patch.PatchManager, stash bool) error {
	if stash {
		if err := c.stashMgr.Save(c.tr.StashPrefix + commits[commitIdx].Sha); err != nil {
			return err
		}
	}

	if err := c.BeginInteractiveRebaseForCommit(commits, commitIdx); err != nil {
		return err
	}

	if err := p.ApplyPatches(c.ApplyPatch, true); err != nil {
		if c.statusMgr.IsRebasing() {
			if err := c.AbortRebase(); err != nil {
				return err
			}
		}
		return err
	}

	// amend the commit
	if err := c.commitsMgr.AmendHead(); err != nil {
		return err
	}

	if c.getWorkflow().InProgress() {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	c.getWorkflow().Start(func() error {
		// add patches to index
		if err := p.ApplyPatches(c.ApplyPatch, false); err != nil {
			if c.statusMgr.IsRebasing() {
				if err := c.AbortRebase(); err != nil {
					return err
				}
			}
			return err
		}

		if stash {
			if err := c.stashMgr.Do(0, "apply"); err != nil {
				return err
			}
		}

		p.Reset()
		return nil
	})

	return c.ContinueRebase()
}

func (c *PatchesMgr) PullPatchIntoNewCommit(commits []*models.Commit, commitIdx int, p *patch.PatchManager) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIdx); err != nil {
		return err
	}

	if err := p.ApplyPatches(c.ApplyPatch, true); err != nil {
		if err := c.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	// amend the commit
	if err := c.commitsMgr.AmendHead(); err != nil {
		return err
	}

	// add patches to index
	if err := p.ApplyPatches(c.ApplyPatch, false); err != nil {
		if err := c.AbortRebase(); err != nil {
			return err
		}
		return err
	}

	head_message, _ := c.commitsMgr.GetHeadMessage()
	new_message := fmt.Sprintf("Split from \"%s\"", head_message)
	err := c.Run(c.commitsMgr.CommitCmdObj(new_message, ""))
	if err != nil {
		return err
	}

	if c.getWorkflow().InProgress() {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	p.Reset()
	return c.ContinueRebase()
}

func (c *PatchesMgr) ApplyPatch(patch string, flags ...string) error {
	filepath := filepath.Join(c.config.GetUserConfigDir(), utils.GetCurrentRepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	c.log.Infof("saving temporary patch to %s", filepath)
	if err := c.os.CreateFileWithContent(filepath, patch); err != nil {
		return err
	}

	flagStr := ""
	for _, flag := range flags {
		flagStr += " --" + flag
	}

	return c.RunGitCmdFromStr(fmt.Sprintf("apply %s %s", flagStr, c.Quote(filepath)))
}

func (c *PatchesMgr) NewPatchManager() *patch.PatchManager {
	return patch.NewPatchManager(c.log, c.diffMgr.ShowFileDiff)
}
