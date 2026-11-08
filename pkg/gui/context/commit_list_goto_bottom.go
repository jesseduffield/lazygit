package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
)

// Returns the index of the commit that the "go to bottom" command selects in a
// list of commits. If the list has several sections, the command first goes to
// the beginning of the section that the selection is in, which is its oldest
// commit. From there it goes to the beginning of the next section, and from the
// last section to the bottom of the list.
//
// The sections are:
//   - the todos of a rebase, cherry-pick or revert
//   - the commits of the branch that aren't merged into a main branch yet
func indexForGotoBottomInCommits(commits []*models.Commit, selectedIdx int) int {
	stops := []int{}

	if headIdx := models.HeadCommitIdx(commits); headIdx > 0 {
		stops = append(stops, headIdx-1)
	}

	if _, idx, found := lo.FindIndexOf(commits, func(c *models.Commit) bool {
		return c.Status == models.StatusMerged
	}); found {
		stops = append(stops, idx-1)
	}

	return closestStopBelow(stops, selectedIdx, len(commits)-1)
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
