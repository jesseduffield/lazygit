package models

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stefanhaller/git-todo-parser/todo"
	"github.com/stretchr/testify/assert"
)

func TestHeadCommitIdx(t *testing.T) {
	testCases := []struct {
		name     string
		commits  []*Commit
		expected int
	}{
		{
			name:     "first commit without rebase todos",
			commits:  makeTestCommits("a", "b"),
			expected: 0,
		},
		{
			name: "first non-todo commit during an interactive rebase",
			commits: []*Commit{
				makeTestTodoCommit(todo.Pick),
				makeTestTodoCommit(todo.Reword),
				makeTestCommit("a"),
				makeTestCommit("b"),
			},
			expected: 2,
		},
		{
			name:     "no commits",
			commits:  nil,
			expected: -1,
		},
		{
			name: "only rebase todos",
			commits: []*Commit{
				makeTestTodoCommit(todo.Pick),
				makeTestTodoCommit(todo.Reword),
			},
			expected: -1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, HeadCommitIdx(testCase.commits))
		})
	}
}

func TestIsHeadCommit(t *testing.T) {
	commits := []*Commit{
		makeTestTodoCommit(todo.Pick),
		makeTestCommit("a"),
		makeTestCommit("b"),
	}

	assert.False(t, IsHeadCommit(commits, 0))
	assert.True(t, IsHeadCommit(commits, 1))
	assert.False(t, IsHeadCommit(commits, 2))
}

func makeTestCommits(hashes ...string) []*Commit {
	commits := make([]*Commit, 0, len(hashes))
	for _, hash := range hashes {
		commits = append(commits, makeTestCommit(hash))
	}

	return commits
}

func makeTestCommit(hash string) *Commit {
	return NewCommit(&utils.StringPool{}, NewCommitOpts{Hash: hash})
}

func makeTestTodoCommit(action todo.TodoCommand) *Commit {
	return NewCommit(&utils.StringPool{}, NewCommitOpts{Action: action})
}
