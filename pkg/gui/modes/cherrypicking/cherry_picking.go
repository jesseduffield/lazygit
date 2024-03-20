package cherrypicking

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
)

type CherryPicking struct {
	CherryPickedCommits []*models.Commit

	// we only allow cherry picking from one context at a time, so you can't copy a commit from the local commits context and then also copy a commit in the reflog context
	ContextKey string
}

func New() *CherryPicking {
	return &CherryPicking{
		CherryPickedCommits: make([]*models.Commit, 0),
		ContextKey:          "",
	}
}

func (self *CherryPicking) Active() bool {
	return len(self.CherryPickedCommits) > 0
}

func (self *CherryPicking) SelectedHashSet() *set.Set[string] {
	hashes := lo.Map(self.CherryPickedCommits, func(commit *models.Commit, _ int) string {
		return commit.Hash
	})
	return set.NewFromSlice(hashes)
}

func (self *CherryPicking) Add(selectedCommit *models.Commit, commitsList []*models.Commit) {
	commitSet := self.SelectedHashSet()
	commitSet.Add(selectedCommit.Hash)

	self.update(commitSet, commitsList)
}

func (self *CherryPicking) Remove(selectedCommit *models.Commit, commitsList []*models.Commit) {
	commitSet := self.SelectedHashSet()
	commitSet.Remove(selectedCommit.Hash)

	self.update(commitSet, commitsList)
}

func (self *CherryPicking) update(selectedHashSet *set.Set[string], commitsList []*models.Commit) {
	cherryPickedCommits := lo.Filter(commitsList, func(commit *models.Commit, _ int) bool {
		return selectedHashSet.Includes(commit.Hash)
	})

	self.CherryPickedCommits = lo.Map(cherryPickedCommits, func(commit *models.Commit, _ int) *models.Commit {
		return &models.Commit{Name: commit.Name, Hash: commit.Hash}
	})
}
