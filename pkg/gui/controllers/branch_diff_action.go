package controllers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// BranchDiffAction lists the files of a branch's diff since it diverged from its base
// branch, for the panels that offer it.
type BranchDiffAction struct {
	c *ControllerCommon
}

// branchDiffBase is a branch that a branch diff can start from, along with the commit
// on it that the diff starts from.
type branchDiffBase struct {
	name      string
	mergeBase string
}

// View opens the commit files panel on top of parent, listing the files of the diff of
// branch, a local branch or a remote one, since the closest of its bases.
func (self *BranchDiffAction) View(branch models.Ref, parent types.Context) error {
	tip, err := self.c.Git().Branch.GetCommitHash(branch.FullRefName())
	if err != nil {
		return err
	}
	bases, err := self.bases(branch)
	if err != nil {
		return err
	}
	return self.view(branch, tip, bases[0], parent)
}

// MenuItems offers the diff of branch since each of its bases, closest first.
func (self *BranchDiffAction) MenuItems(branch models.Ref, parent types.Context) ([]*types.MenuItem, error) {
	tip, err := self.c.Git().Branch.GetCommitHash(branch.FullRefName())
	if err != nil {
		return nil, err
	}
	bases, err := self.bases(branch)
	if err != nil {
		return []*types.MenuItem{{
			LabelColumns:   []string{self.c.Tr.ViewBranchDiff},
			DisabledReason: &types.DisabledReason{Text: err.Error()},
		}}, nil
	}

	return lo.Map(bases, func(base branchDiffBase, _ int) *types.MenuItem {
		return &types.MenuItem{
			LabelColumns: []string{utils.ResolvePlaceholderString(
				self.c.Tr.ViewBranchDiffSince,
				map[string]string{"baseBranch": base.name},
			)},
			OnPress: func() error {
				return self.view(branch, tip, base, parent)
			},
		}
	}), nil
}

func (self *BranchDiffAction) view(branch models.Ref, tip string, base branchDiffBase, parent types.Context) error {
	ref := &models.BranchDiff{
		Branch:     branch,
		BaseBranch: base.name,
		MergeBase:  base.mergeBase,
		Tip:        tip,
	}
	return viewCommitFiles(self.c, ref, nil, false, parent)
}

// bases returns the branches that the diff of branch can start from.
func (self *BranchDiffAction) bases(branch models.Ref) ([]branchDiffBase, error) {
	mainBranches := self.c.Model().MainBranches
	baseBranch, err := self.c.Git().Loaders.BranchLoader.GetBaseBranch(branch.FullRefName(), mainBranches)
	if err != nil {
		return nil, err
	}
	mergeBase := mainBranches.GetMergeBase(branch.FullRefName())
	if baseBranch == "" || mergeBase == "" {
		return nil, errors.New(self.c.Tr.CouldNotDetermineBaseBranch)
	}
	return []branchDiffBase{{name: helpers.ShortBranchName(baseBranch), mergeBase: mergeBase}}, nil
}
