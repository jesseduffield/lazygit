package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// BaseBranchHelper resolves the base branch for a given branch. The
// closeness rule (smallest ahead value) usually picks a single answer
// but can leave a tie when the branch's fork point is reachable from
// more than one main branch — in that case the helper surfaces the
// candidates so the caller can disambiguate.
type BaseBranchHelper struct {
	c *HelperCommon
}

func NewBaseBranchHelper(c *HelperCommon) *BaseBranchHelper {
	return &BaseBranchHelper{c: c}
}

// ResolveBaseBranch returns the base branch for the given branch, the
// full set of tied candidates (for any disambiguation UI), and whether
// the answer is genuinely ambiguous (more than one candidate tied at
// the closest position).
//
// An empty baseRef (with no error) means no configured main branch
// contains the branch — not an error condition.
func (self *BaseBranchHelper) ResolveBaseBranch(branch *models.Branch) (baseRef string, ambiguous bool, candidates []string, err error) {
	mainBranches := self.c.Model().MainBranches
	candidates, err = self.c.Git().Loaders.BranchLoader.GetBaseBranchCandidates(branch, mainBranches)
	if err != nil {
		return "", false, nil, err
	}
	if len(candidates) == 0 {
		return "", false, nil, nil
	}
	return candidates[0], len(candidates) > 1, candidates, nil
}

// ShowPicker presents a menu of candidate base branches and runs
// onPicked with the user's selection. Callers should only invoke this
// when ResolveBaseBranch reported ambiguous=true; for the
// single-candidate case there is nothing to pick.
func (self *BaseBranchHelper) ShowPicker(
	branch *models.Branch,
	candidates []string,
	onPicked func(baseRef string) error,
) error {
	items := lo.Map(candidates, func(ref string, _ int) *types.MenuItem {
		return &types.MenuItem{
			Label:   ShortBranchName(ref),
			OnPress: func() error { return onPicked(ref) },
		}
	})
	return self.c.Menu(types.CreateMenuOptions{
		Title: utils.ResolvePlaceholderString(self.c.Tr.PickBaseBranchTitle,
			map[string]string{"branchName": branch.Name}),
		Prompt: self.c.Tr.PickBaseBranchPrompt,
		Items:  items,
	})
}
