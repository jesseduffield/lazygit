package helpers

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
)

// Returns the branches stacked below the given branch: those whose tip is one
// of the given commits of the branch that hasn't been merged into a main
// branch yet. The branch itself and the main branches are left out. The
// result is in the order in which the tips appear in commits, so the branch
// closest to the given one comes first.
//
// A commit counts as merged if the remote branch of a main branch contains it
// (this is what the commit's status says), but also if the local main branch
// contains it. So if a branch was merged into the local main branch, it isn't
// stacked below a branch on top of that main branch, and if the given branch
// is itself a main branch, nothing is stacked below it.
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
	mainBranchTips := []string{}
	for _, other := range branches {
		if lo.Contains(mainBranches, other.Name) {
			mainBranchTips = append(mainBranchTips, other.CommitHash)
		} else if other.Name != branch.Name {
			branchesByTip[other.CommitHash] = append(branchesByTip[other.CommitHash], other)
		}
	}

	mergedLocally := reachableCommits(commits, mainBranchTips)

	result := []*models.Branch{}
	for _, commit := range commits {
		if commit.IsTODO() || commit.Status == models.StatusMerged || mergedLocally.Includes(commit.Hash()) {
			continue
		}
		result = append(result, branchesByTip[commit.Hash()]...)
	}

	return result
}

// Returns the hashes of those of the given commits that can be reached from
// any of the given tips by following their parents
func reachableCommits(commits []*models.Commit, tips []string) *set.Set[string] {
	commitsByHash := map[string]*models.Commit{}
	for _, commit := range commits {
		if !commit.IsTODO() {
			commitsByHash[commit.Hash()] = commit
		}
	}

	reachable := set.New[string]()
	pending := append([]string{}, tips...)
	for len(pending) > 0 {
		hash := pending[len(pending)-1]
		pending = pending[:len(pending)-1]

		commit, ok := commitsByHash[hash]
		if !ok || reachable.Includes(hash) {
			continue
		}
		reachable.Add(hash)
		pending = append(pending, commit.Parents()...)
	}

	return reachable
}
