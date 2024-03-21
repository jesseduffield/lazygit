package utils

import (
	"errors"
	"testing"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/stretchr/testify/assert"
)

func TestRebaseCommands_moveTodoDown(t *testing.T) {
	type scenario struct {
		testName       string
		todos          []todo.Todo
		todoToMoveDown Todo
		expectedErr    string
		expectedTodos  []todo.Todo
	}

	scenarios := []scenario{
		{
			testName: "simple case 1 - move to beginning",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveDown: Todo{Sha: "5678", Action: todo.Pick},
			expectedErr:    "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "abcd"},
			},
		},
		{
			testName: "simple case 2 - move from end",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveDown: Todo{Sha: "abcd", Action: todo.Pick},
			expectedErr:    "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "abcd"},
				{Command: todo.Pick, Commit: "5678"},
			},
		},
		{
			testName: "move update-ref todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
			},
			todoToMoveDown: Todo{Ref: "refs/heads/some_branch", Action: todo.UpdateRef},
			expectedErr:    "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "5678"},
			},
		},
		{
			testName: "skip an invisible todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "abcd"},
				{Command: todo.Label, Label: "myLabel"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "def0"},
			},
			todoToMoveDown: Todo{Sha: "5678", Action: todo.Pick},
			expectedErr:    "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
				{Command: todo.Label, Label: "myLabel"},
				{Command: todo.Pick, Commit: "def0"},
			},
		},

		// Error cases
		{
			testName: "commit not found",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveDown: Todo{Sha: "def0", Action: todo.Pick},
			expectedErr:    "Todo def0 not found in git-rebase-todo",
			expectedTodos:  []todo.Todo{},
		},
		{
			testName: "trying to move first commit down",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveDown: Todo{Sha: "1234", Action: todo.Pick},
			expectedErr:    "Destination position for moving todo is out of range",
			expectedTodos:  []todo.Todo{},
		},
		{
			testName: "trying to move commit down when all commits before are invisible",
			todos: []todo.Todo{
				{Command: todo.Label, Label: "myLabel"},
				{Command: todo.Reset, Label: "otherlabel"},
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
			},
			todoToMoveDown: Todo{Sha: "1234", Action: todo.Pick},
			expectedErr:    "Destination position for moving todo is out of range",
			expectedTodos:  []todo.Todo{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			rearrangedTodos, err := moveTodoDown(s.todos, s.todoToMoveDown)
			if s.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, s.expectedErr)
			}
			assert.Equal(t, s.expectedTodos, rearrangedTodos)
		},
		)
	}
}

func TestRebaseCommands_moveTodoUp(t *testing.T) {
	type scenario struct {
		testName      string
		todos         []todo.Todo
		todoToMoveUp  Todo
		expectedErr   string
		expectedTodos []todo.Todo
	}

	scenarios := []scenario{
		{
			testName: "simple case 1 - move to end",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveUp: Todo{Sha: "5678", Action: todo.Pick},
			expectedErr:  "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "abcd"},
				{Command: todo.Pick, Commit: "5678"},
			},
		},
		{
			testName: "simple case 2 - move from beginning",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveUp: Todo{Sha: "1234", Action: todo.Pick},
			expectedErr:  "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "abcd"},
			},
		},
		{
			testName: "move update-ref todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "5678"},
			},
			todoToMoveUp: Todo{Ref: "refs/heads/some_branch", Action: todo.UpdateRef},
			expectedErr:  "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
			},
		},
		{
			testName: "skip an invisible todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "abcd"},
				{Command: todo.Label, Label: "myLabel"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "def0"},
			},
			todoToMoveUp: Todo{Sha: "abcd", Action: todo.Pick},
			expectedErr:  "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Label, Label: "myLabel"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
				{Command: todo.Pick, Commit: "def0"},
			},
		},

		// Error cases
		{
			testName: "commit not found",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveUp:  Todo{Sha: "def0", Action: todo.Pick},
			expectedErr:   "Todo def0 not found in git-rebase-todo",
			expectedTodos: []todo.Todo{},
		},
		{
			testName: "trying to move last commit up",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todoToMoveUp:  Todo{Sha: "abcd", Action: todo.Pick},
			expectedErr:   "Destination position for moving todo is out of range",
			expectedTodos: []todo.Todo{},
		},
		{
			testName: "trying to move commit up when all commits after it are invisible",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Label, Label: "myLabel"},
				{Command: todo.Reset, Label: "otherlabel"},
			},
			todoToMoveUp:  Todo{Sha: "5678", Action: todo.Pick},
			expectedErr:   "Destination position for moving todo is out of range",
			expectedTodos: []todo.Todo{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			rearrangedTodos, err := moveTodoUp(s.todos, s.todoToMoveUp)
			if s.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, s.expectedErr)
			}
			assert.Equal(t, s.expectedTodos, rearrangedTodos)
		},
		)
	}
}

func TestRebaseCommands_moveFixupCommitDown(t *testing.T) {
	scenarios := []struct {
		name          string
		todos         []todo.Todo
		originalSha   string
		fixupSha      string
		expectedTodos []todo.Todo
		expectedErr   error
	}{
		{
			name: "fixup commit is the last commit",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalSha: "original",
			fixupSha:    "fixup",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Fixup, Commit: "fixup"},
			},
			expectedErr: nil,
		},
		{
			// TODO: is this something we actually want to support?
			name: "fixup commit is separated from original commit",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "other"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalSha: "original",
			fixupSha:    "fixup",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Fixup, Commit: "fixup"},
				{Command: todo.Pick, Commit: "other"},
			},
			expectedErr: nil,
		},
		{
			name: "More original SHAs than expected",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalSha:   "original",
			fixupSha:      "fixup",
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one original SHA, found 2"),
		},
		{
			name: "More fixup SHAs than expected",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalSha:   "original",
			fixupSha:      "fixup",
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one fixup SHA, found 2"),
		},
		{
			name: "No fixup SHAs found",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
			},
			originalSha:   "original",
			fixupSha:      "fixup",
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one fixup SHA, found 0"),
		},
		{
			name: "No original SHAs found",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalSha:   "original",
			fixupSha:      "fixup",
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one original SHA, found 0"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actualTodos, actualErr := moveFixupCommitDown(scenario.todos, scenario.originalSha, scenario.fixupSha)

			if scenario.expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.EqualError(t, actualErr, scenario.expectedErr.Error())
			}

			assert.EqualValues(t, scenario.expectedTodos, actualTodos)
		})
	}
}

func TestRebaseCommands_deleteTodos(t *testing.T) {
	scenarios := []struct {
		name          string
		todos         []todo.Todo
		todosToDelete []Todo
		expectedTodos []todo.Todo
		expectedErr   error
	}{
		{
			name: "success",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Pick, Commit: "abcd"},
			},
			todosToDelete: []Todo{
				{Ref: "refs/heads/some_branch", Action: todo.UpdateRef},
				{Sha: "abcd", Action: todo.Pick},
			},
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
			},
			expectedErr: nil,
		},
		{
			name: "failure",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
			},
			todosToDelete: []Todo{
				{Sha: "abcd", Action: todo.Pick},
			},
			expectedTodos: []todo.Todo{},
			expectedErr:   errors.New("Todo abcd not found in git-rebase-todo"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actualTodos, actualErr := deleteTodos(scenario.todos, scenario.todosToDelete)

			if scenario.expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.EqualError(t, actualErr, scenario.expectedErr.Error())
			}

			assert.EqualValues(t, scenario.expectedTodos, actualTodos)
		})
	}
}
