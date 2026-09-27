package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
	"github.com/stretchr/testify/assert"
)

func TestBranchesBelowInStack(t *testing.T) {
	hashPool := &utils.StringPool{}
	commit := func(hash string, status models.CommitStatus) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{Hash: hash, Status: status})
	}
	todoCommit := func(hash string) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{Hash: hash, Action: todo.Pick, Status: models.StatusRebasing})
	}
	branch := func(name string, tip string) *models.Branch {
		return &models.Branch{Name: name, CommitHash: tip}
	}

	current := branch("current", "c3")
	mainBranches := []string{"master", "main"}

	scenarios := []struct {
		testName string
		commits  []*models.Commit
		branches []*models.Branch
		expected []string
	}{
		{
			testName: "stack of branches, closest to the current one first",
			commits: []*models.Commit{
				commit("c3", models.StatusUnpushed),
				commit("c2", models.StatusPushed),
				commit("c1", models.StatusUnpushed),
				commit("m1", models.StatusMerged),
			},
			branches: []*models.Branch{current, branch("first", "c1"), branch("second", "c2"), branch("master", "m1")},
			expected: []string{"second", "first"},
		},
		{
			testName: "branch pointing at the same commit as the current one",
			commits: []*models.Commit{
				commit("c3", models.StatusUnpushed),
				commit("m1", models.StatusMerged),
			},
			branches: []*models.Branch{current, branch("twin", "c3")},
			expected: []string{"twin"},
		},
		{
			testName: "several branches pointing at one commit keep their order",
			commits: []*models.Commit{
				commit("c3", models.StatusUnpushed),
				commit("c1", models.StatusUnpushed),
			},
			branches: []*models.Branch{current, branch("one", "c1"), branch("two", "c1")},
			expected: []string{"one", "two"},
		},
		{
			testName: "main branches are left out even if their tip isn't merged",
			commits: []*models.Commit{
				commit("c3", models.StatusUnpushed),
				commit("c1", models.StatusUnpushed),
			},
			branches: []*models.Branch{current, branch("main", "c1")},
			expected: []string{},
		},
		{
			testName: "branches whose tip is merged are left out",
			commits: []*models.Commit{
				commit("c3", models.StatusUnpushed),
				commit("m1", models.StatusMerged),
				commit("m0", models.StatusMerged),
			},
			branches: []*models.Branch{current, branch("old", "m0"), branch("master", "m1")},
			expected: []string{},
		},
		{
			testName: "branches whose tip is not among the commits are left out",
			commits: []*models.Commit{
				commit("c3", models.StatusUnpushed),
			},
			branches: []*models.Branch{current, branch("sibling", "s1")},
			expected: []string{},
		},
		{
			testName: "todo commits of a rebase don't count",
			commits: []*models.Commit{
				todoCommit("t1"),
				commit("c3", models.StatusUnpushed),
			},
			branches: []*models.Branch{current, branch("other", "t1")},
			expected: []string{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := BranchesBelowInStack(s.commits, s.branches, current, mainBranches)
			names := lo.Map(result, func(b *models.Branch, _ int) string { return b.Name })
			assert.Equal(t, s.expected, names)
		})
	}
}
