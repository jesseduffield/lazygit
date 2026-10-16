package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
)

// Returns the branches stacked below the given branch: those whose tip is one
// of the given commits of the branch that hasn't been merged into a main
// branch yet. The branch itself and the main branches are left out. The
// result is in the order in which the tips appear in commits, so the branch
// closest to the given one comes first.
//
// Only the commits passed in are looked at, so a branch whose tip is further
// down the history than the loaded commits is not found.
func BranchesBelowInStack(
	commits []*models.Commit,
	branches []*models.Branch,
	branch *models.Branch,
	mainBranches []string,
) []*models.Branch {
	branchesByTip := map[string][]*models.Branch{}
	for _, other := range branches {
		if other.Name == branch.Name || lo.Contains(mainBranches, other.Name) {
			continue
		}
		branchesByTip[other.CommitHash] = append(branchesByTip[other.CommitHash], other)
	}

	result := []*models.Branch{}
	for _, commit := range commits {
		if commit.IsTODO() || commit.Status == models.StatusMerged {
			continue
		}
		result = append(result, branchesByTip[commit.Hash()]...)
	}

	return result
}
