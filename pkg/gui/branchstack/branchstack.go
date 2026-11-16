package branchstack

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
)

// IsStackBranch tells whether the given branch can be one of the stack of
// branches below the current branch. The commits list visualizes them with
// markers at their heads. The current branch and the main branches are not,
// because the goal is to show stacks of local branches, and they don't
// contribute to one.
func IsStackBranch(branch *models.Branch, currentBranchName string, mainBranches []string) bool {
	// Branches without a commit hash are excluded too. As far as I can see,
	// this happens for a detached head.
	return branch.CommitHash != "" &&
		branch.Name != currentBranchName &&
		!lo.Contains(mainBranches, branch.Name)
}

// Returns the branches stacked below the branch with the given name: those
// whose tip is one of the given commits of the branch that hasn't been
// merged into a main branch yet. The branch itself and the main branches
// are left out. The result is in the order in which the tips appear in
// commits, so the branch closest to the given one comes first.
//
// A commit counts as merged if the remote branch of a main branch contains it
// (this is what the commit's status says), but also if the local main branch
// contains it. So if a branch was merged into the local main branch, it isn't
// stacked below a branch on top of that main branch, and if the given branch
// is itself a main branch, nothing is stacked below it.
//
// Only the commits passed in are looked at, so a branch whose tip is further
// down the history than the loaded commits is not found.
func BranchesBelow(
	commits []*models.Commit,
	branches []*models.Branch,
	branchName string,
	mainBranches []string,
) []*models.Branch {
	branchesByTip := map[string][]*models.Branch{}
	for _, other := range branches {
		if IsStackBranch(other, branchName, mainBranches) {
			branchesByTip[other.CommitHash] = append(branchesByTip[other.CommitHash], other)
		}
	}

	mergedLocally := CommitsMergedIntoLocalMainBranch(commits, branches, mainBranches)

	result := []*models.Branch{}
	for _, commit := range commits {
		if commit.IsTODO() || commit.Status == models.StatusMerged || mergedLocally.Includes(commit.Hash()) {
			continue
		}
		result = append(result, branchesByTip[commit.Hash()]...)
	}

	return result
}

// Returns the hashes of those of the given commits that a local main branch
// contains, but whose status doesn't say that they are merged. The status
// looks at the remote branch of a main branch if it has one, so these are
// commits that were merged into a local main branch but not pushed yet.
func CommitsMergedIntoLocalMainBranch(
	commits []*models.Commit,
	branches []*models.Branch,
	mainBranches []string,
) *set.Set[string] {
	result := set.New[string]()

	mainBranchTips := lo.FilterMap(branches, func(branch *models.Branch, _ int) (string, bool) {
		return branch.CommitHash, lo.Contains(mainBranches, branch.Name)
	})
	if len(mainBranchTips) == 0 {
		return result
	}

	// The parents of a merged commit are merged too, so the walk never needs
	// to go past a merged commit. Usually the tip of a local main branch is
	// merged itself, so there is nothing to walk; this keeps it cheap enough to
	// do on every render of the commits list.
	unmergedCommitsByHash := map[string]*models.Commit{}
	for _, commit := range commits {
		if !commit.IsTODO() && commit.Status != models.StatusMerged {
			unmergedCommitsByHash[commit.Hash()] = commit
		}
	}

	pending := mainBranchTips
	for len(pending) > 0 {
		hash := pending[len(pending)-1]
		pending = pending[:len(pending)-1]

		commit, ok := unmergedCommitsByHash[hash]
		if !ok || result.Includes(hash) {
			continue
		}
		result.Add(hash)
		pending = append(pending, commit.Parents()...)
	}

	return result
}
