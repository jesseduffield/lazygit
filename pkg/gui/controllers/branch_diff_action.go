package controllers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// BranchDiffAction lists the files of a branch's diff since it diverged from its base
// branch, for the panels that offer it.
type BranchDiffAction struct {
	c *ControllerCommon
}

// View opens the commit files panel on top of parent, listing the files of the diff of
// branch, a local branch or a remote one.
func (self *BranchDiffAction) View(branch models.Ref, parent types.Context) error {
	mainBranches := self.c.Model().MainBranches
	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(branch.FullRefName(), mainBranches)
	if err != nil {
		return err
	}
	mergeBase := mainBranches.GetMergeBase(branch.FullRefName())
	if baseBranch == "" || mergeBase == "" {
		return errors.New(self.c.Tr.CouldNotDetermineBaseBranch)
	}
	tip, err := self.c.Git().Branch.GetCommitHash(branch.FullRefName())
	if err != nil {
		return err
	}

	ref := &models.BranchDiff{
		Branch:     branch,
		BaseBranch: helpers.ShortBranchName(baseBranch),
		MergeBase:  mergeBase,
		Tip:        tip,
	}
	return viewCommitFiles(self.c, ref, nil, false, parent)
}
