package models

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stefanhaller/git-todo-parser/todo"
	"github.com/stretchr/testify/assert"
)

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

func makeTestCommit(hash string) *Commit {
	return NewCommit(&utils.StringPool{}, NewCommitOpts{Hash: hash})
}

func makeTestTodoCommit(action todo.TodoCommand) *Commit {
	return NewCommit(&utils.StringPool{}, NewCommitOpts{Action: action})
}
