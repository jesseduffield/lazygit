package context

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
)

// Returns the index of the commit that the "go to bottom" command selects in a
// list of commits. If the list has several sections, the command first goes to
// the beginning of the section that the selection is in, which is its oldest
// commit. From there it goes to the beginning of the next section, and from the
// last section to the bottom of the list.
//
// The sections are:
//   - the todos of a rebase, cherry-pick or revert
//   - each branch of a stack, given by branchesBelow. Among the todos of a
//     rebase, update-ref todos separate the branches.
//   - the commits that aren't merged into a main branch yet
//   - the remote commits of a divergence view, which are above the local ones
//
// During a bisect, it stops at the bottommost commit marked bad, at the commit
// to test, and at the topmost commit marked good instead. The commits list
// doesn't color the commits by whether they are merged during a bisect, so the
// beginnings of the sections would look like arbitrary commits.
func indexForGotoBottomInCommits(
	commits []*models.Commit,
	selectedIdx int,
	branchesBelow []*models.Branch,
	bisectInfo *git_commands.BisectInfo,
) int {
	if bisectInfo.Started() {
		return closestStopBelow(bisectStops(commits, bisectInfo), selectedIdx, len(commits)-1)
	}

	stops := []int{}

	if headIdx := models.HeadCommitIdx(commits); headIdx > 0 {
		stops = append(stops, headIdx-1)
	}

	tipsOfBranchesBelow := set.NewFromSlice(lo.Map(branchesBelow, func(branch *models.Branch, _ int) string {
		return branch.CommitHash
	}))
	for idx, commit := range commits {
		if commit.IsTODO() {
			// If several update-ref todos follow each other, the branches
			// between them have no todos, so only stop above the first one
			if commit.Action == todo.UpdateRef && idx > 0 && commits[idx-1].Action != todo.UpdateRef {
				stops = append(stops, idx-1)
			}
		} else if tipsOfBranchesBelow.Includes(commit.Hash()) {
			stops = append(stops, idx-1)
		}
	}

	if _, idx, found := lo.FindIndexOf(commits, func(c *models.Commit) bool {
		return c.Status == models.StatusMerged
	}); found {
		stops = append(stops, idx-1)
	}

	if _, idx, found := lo.FindIndexOf(commits, func(c *models.Commit) bool {
		return c.Divergence == models.DivergenceLeft
	}); found {
		stops = append(stops, idx-1)
	}

	return closestStopBelow(stops, selectedIdx, len(commits)-1)
}

func bisectStops(commits []*models.Commit, bisectInfo *git_commands.BisectInfo) []int {
	hasStatus := func(status git_commands.BisectStatus) func(*models.Commit) bool {
		return func(commit *models.Commit) bool {
			commitStatus, ok := bisectInfo.Status(commit.Hash())
			return ok && commitStatus == status
		}
	}
	isCommitToTest := func(commit *models.Commit) bool {
		return bisectInfo.HasCommitToTest() && commit.Hash() == bisectInfo.GetCurrentHash()
	}

	stops := []int{}
	if _, idx, found := lo.FindLastIndexOf(commits, hasStatus(git_commands.BisectStatusNew)); found {
		stops = append(stops, idx)
	}
	if _, idx, found := lo.FindIndexOf(commits, isCommitToTest); found {
		stops = append(stops, idx)
	}
	if _, idx, found := lo.FindIndexOf(commits, hasStatus(git_commands.BisectStatusOld)); found {
		stops = append(stops, idx)
	}
	return stops
}

// Returns the smallest of the given stops that is greater than the selected
// index, or the last index if there is none
func closestStopBelow(stops []int, selectedIdx int, lastIdx int) int {
	result := lastIdx
	for _, stop := range stops {
		if stop > selectedIdx && stop < result {
			result = stop
		}
	}
	return result
}
