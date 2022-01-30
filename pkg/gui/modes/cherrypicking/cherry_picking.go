package cherrypicking

import (
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

func (m *CherryPicking) Active() bool {
	return len(m.CherryPickedCommits) > 0
}
