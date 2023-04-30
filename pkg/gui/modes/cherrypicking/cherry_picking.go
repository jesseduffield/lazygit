package cherrypicking

import (
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
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

func (self *CherryPicking) SelectedShaSet() *set.Set[string] {
	shas := slices.Map(self.CherryPickedCommits, func(commit *models.Commit) string {
		return commit.Sha
	})
	return set.NewFromSlice(shas)
}

func (self *CherryPicking) Add(selectedCommit *models.Commit, commitsList []*models.Commit) {
	commitSet := self.SelectedShaSet()
	commitSet.Add(selectedCommit.Sha)

	self.update(commitSet, commitsList)
}

func (self *CherryPicking) Remove(selectedCommit *models.Commit, commitsList []*models.Commit) {
	commitSet := self.SelectedShaSet()
	commitSet.Remove(selectedCommit.Sha)

	self.update(commitSet, commitsList)
}

func (self *CherryPicking) update(selectedShaSet *set.Set[string], commitsList []*models.Commit) {
	cherryPickedCommits := slices.Filter(commitsList, func(commit *models.Commit) bool {
		return selectedShaSet.Includes(commit.Sha)
	})

	self.CherryPickedCommits = slices.Map(cherryPickedCommits, func(commit *models.Commit) *models.Commit {
		return &models.Commit{Name: commit.Name, Sha: commit.Sha}
	})
}
