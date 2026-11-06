package utils

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stefanhaller/git-todo-parser/todo"
	"github.com/stretchr/testify/assert"
)

func TestMoveTodos(t *testing.T) {
	todos := []todo.Todo{
		{Command: todo.Pick, Commit: "a"},
		{Command: todo.Pick, Commit: "b"},
		{Command: todo.Label, Label: "hidden"},
		{Command: todo.Pick, Commit: "c"},
		{Command: todo.Pick, Commit: "d"},
		{Command: todo.Pick, Commit: "e"},
		{Command: todo.Pick, Commit: "f"},
	}

	t.Run("moves a range up multiple rendered rows", func(t *testing.T) {
		actual, err := moveTodos(slices.Clone(todos), []Todo{{Hash: "d"}, {Hash: "c"}}, -2)

		assert.NoError(t, err)
		assert.Equal(t, []todo.Todo{
			{Command: todo.Pick, Commit: "a"},
			{Command: todo.Pick, Commit: "b"},
			{Command: todo.Label, Label: "hidden"},
			{Command: todo.Pick, Commit: "e"},
			{Command: todo.Pick, Commit: "f"},
			{Command: todo.Pick, Commit: "c"},
			{Command: todo.Pick, Commit: "d"},
		}, actual)
	})

	t.Run("moves a range down multiple rendered rows", func(t *testing.T) {
		actual, err := moveTodos(slices.Clone(todos), []Todo{{Hash: "e"}, {Hash: "d"}}, 2)

		assert.NoError(t, err)
		assert.Equal(t, []todo.Todo{
			{Command: todo.Pick, Commit: "a"},
			{Command: todo.Pick, Commit: "d"},
			{Command: todo.Pick, Commit: "e"},
			{Command: todo.Pick, Commit: "b"},
			{Command: todo.Label, Label: "hidden"},
			{Command: todo.Pick, Commit: "c"},
			{Command: todo.Pick, Commit: "f"},
		}, actual)
	})
}

func TestMoveCommitTodos(t *testing.T) {
	pick := func(hash string) todo.Todo { return todo.Todo{Command: todo.Pick, Commit: hash} }
	updateRef := func(ref string) todo.Todo { return todo.Todo{Command: todo.UpdateRef, Ref: ref} }

	scenarios := []struct {
		testName      string
		todos         []todo.Todo
		hashes        []string
		destination   TodoMoveDestination
		expectedErr   string
		expectedTodos []todo.Todo
	}{
		{
			testName:      "move to the beginning",
			todos:         []todo.Todo{pick("a"), pick("b"), pick("c")},
			hashes:        []string{"c"},
			destination:   TodoMoveDestination{},
			expectedTodos: []todo.Todo{pick("c"), pick("a"), pick("b")},
		},
		{
			testName: "move to the beginning, after labels and resets",
			todos: []todo.Todo{
				{Command: todo.Label, Label: "onto"},
				{Command: todo.Reset, Label: "onto"},
				pick("a"),
				pick("b"),
			},
			hashes:      []string{"b"},
			destination: TodoMoveDestination{},
			expectedTodos: []todo.Todo{
				{Command: todo.Label, Label: "onto"},
				{Command: todo.Reset, Label: "onto"},
				pick("b"),
				pick("a"),
			},
		},
		{
			testName:      "move towards the beginning",
			todos:         []todo.Todo{pick("a"), pick("b"), pick("c"), pick("d")},
			hashes:        []string{"d"},
			destination:   TodoMoveDestination{AfterHash: "a"},
			expectedTodos: []todo.Todo{pick("a"), pick("d"), pick("b"), pick("c")},
		},
		{
			testName:      "move towards the end",
			todos:         []todo.Todo{pick("a"), pick("b"), pick("c"), pick("d")},
			hashes:        []string{"a"},
			destination:   TodoMoveDestination{AfterHash: "c"},
			expectedTodos: []todo.Todo{pick("b"), pick("c"), pick("a"), pick("d")},
		},
		{
			testName:      "move several commits",
			todos:         []todo.Todo{pick("a"), pick("b"), pick("c"), pick("d")},
			hashes:        []string{"c", "b"},
			destination:   TodoMoveDestination{AfterHash: "d"},
			expectedTodos: []todo.Todo{pick("a"), pick("d"), pick("b"), pick("c")},
		},
		{
			testName:      "update-refs between the moved commits move along",
			todos:         []todo.Todo{pick("a"), pick("b"), updateRef("refs/heads/x"), pick("c"), pick("d")},
			hashes:        []string{"c", "b"},
			destination:   TodoMoveDestination{AfterHash: "d"},
			expectedTodos: []todo.Todo{pick("a"), pick("d"), pick("b"), updateRef("refs/heads/x"), pick("c")},
		},
		{
			testName:      "update-refs after the moved commits stay",
			todos:         []todo.Todo{pick("a"), pick("b"), updateRef("refs/heads/x"), pick("c")},
			hashes:        []string{"b"},
			destination:   TodoMoveDestination{},
			expectedTodos: []todo.Todo{pick("b"), pick("a"), updateRef("refs/heads/x"), pick("c")},
		},
		{
			testName:      "update-refs at the beginning stay before the moved commits",
			todos:         []todo.Todo{pick("a"), updateRef("refs/heads/x"), pick("b")},
			hashes:        []string{"a"},
			destination:   TodoMoveDestination{},
			expectedTodos: []todo.Todo{updateRef("refs/heads/x"), pick("a"), pick("b")},
		},
		{
			testName:      "update-refs at the beginning that are listed go after the moved commits",
			todos:         []todo.Todo{pick("a"), updateRef("refs/heads/x"), pick("b")},
			hashes:        []string{"a"},
			destination:   TodoMoveDestination{RefsAfter: []string{"refs/heads/x"}},
			expectedTodos: []todo.Todo{pick("a"), updateRef("refs/heads/x"), pick("b")},
		},
		{
			testName:      "update-refs of the destination stay before the moved commits",
			todos:         []todo.Todo{pick("a"), updateRef("refs/heads/x"), pick("b"), pick("c")},
			hashes:        []string{"c"},
			destination:   TodoMoveDestination{AfterHash: "a"},
			expectedTodos: []todo.Todo{pick("a"), updateRef("refs/heads/x"), pick("c"), pick("b")},
		},
		{
			testName: "update-refs of the destination that are listed go after the moved commits",
			todos: []todo.Todo{
				pick("a"),
				updateRef("refs/heads/x"),
				updateRef("refs/heads/y"),
				updateRef("refs/heads/z"),
				pick("b"),
			},
			hashes:      []string{"b"},
			destination: TodoMoveDestination{AfterHash: "a", RefsAfter: []string{"refs/heads/x", "refs/heads/z"}},
			expectedTodos: []todo.Todo{
				pick("a"),
				updateRef("refs/heads/y"),
				pick("b"),
				updateRef("refs/heads/x"),
				updateRef("refs/heads/z"),
			},
		},
		{
			testName: "update-refs after the moved commits join those of the destination",
			todos: []todo.Todo{
				pick("a"),
				updateRef("refs/heads/x"),
				pick("b"),
				updateRef("refs/heads/y"),
				pick("c"),
			},
			hashes:      []string{"b"},
			destination: TodoMoveDestination{AfterHash: "a", RefsAfter: []string{"refs/heads/x", "refs/heads/y"}},
			expectedTodos: []todo.Todo{
				pick("a"),
				pick("b"),
				updateRef("refs/heads/x"),
				updateRef("refs/heads/y"),
				pick("c"),
			},
		},
		{
			testName:    "commit not found",
			todos:       []todo.Todo{pick("a"), pick("b")},
			hashes:      []string{"c"},
			destination: TodoMoveDestination{},
			expectedErr: "Todo c not found in git-rebase-todo",
		},
		{
			testName:    "destination not found",
			todos:       []todo.Todo{pick("a"), pick("b")},
			hashes:      []string{"b"},
			destination: TodoMoveDestination{AfterHash: "c"},
			expectedErr: "Todo c not found in git-rebase-todo",
		},
		{
			testName:    "other commit between the moved commits",
			todos:       []todo.Todo{pick("a"), pick("b"), pick("c"), pick("d")},
			hashes:      []string{"d", "b"},
			destination: TodoMoveDestination{},
			expectedErr: "The commits to move are not contiguous in git-rebase-todo",
		},
		{
			testName: "label between the moved commits",
			todos: []todo.Todo{
				pick("a"),
				pick("b"),
				{Command: todo.Label, Label: "x"},
				pick("c"),
			},
			hashes:      []string{"c", "b"},
			destination: TodoMoveDestination{},
			expectedErr: "The commits to move are not contiguous in git-rebase-todo",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			actual, err := moveCommitTodos(s.todos, s.hashes, s.destination)
			if s.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, s.expectedTodos, actual)
			} else {
				assert.EqualError(t, err, s.expectedErr)
			}
		})
	}
}

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
			todoToMoveDown: Todo{Hash: "5678"},
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
			todoToMoveDown: Todo{Hash: "abcd"},
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
			todoToMoveDown: Todo{Ref: "refs/heads/some_branch"},
			expectedErr:    "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "5678"},
			},
		},
		{
			testName: "move across update-ref todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "5678"},
			},
			todoToMoveDown: Todo{Hash: "5678"},
			expectedErr:    "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
			},
		},
		{
			testName: "move across exec todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Exec, ExecCommand: "make test"},
				{Command: todo.Pick, Commit: "5678"},
			},
			todoToMoveDown: Todo{Hash: "5678"},
			expectedErr:    "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.Exec, ExecCommand: "make test"},
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
			todoToMoveDown: Todo{Hash: "5678"},
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
			todoToMoveDown: Todo{Hash: "def0"},
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
			todoToMoveDown: Todo{Hash: "1234"},
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
			todoToMoveDown: Todo{Hash: "1234"},
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
			todoToMoveUp: Todo{Hash: "5678"},
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
			todoToMoveUp: Todo{Hash: "1234"},
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
			todoToMoveUp: Todo{Ref: "refs/heads/some_branch"},
			expectedErr:  "",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
			},
		},
		{
			testName: "move across update-ref todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "5678"},
			},
			todoToMoveUp: Todo{Hash: "1234"},
			expectedErr:  "",
			expectedTodos: []todo.Todo{
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
			},
		},
		{
			testName: "move across exec todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Exec, ExecCommand: "make test"},
				{Command: todo.Pick, Commit: "5678"},
			},
			todoToMoveUp: Todo{Hash: "1234"},
			expectedErr:  "",
			expectedTodos: []todo.Todo{
				{Command: todo.Exec, ExecCommand: "make test"},
				{Command: todo.Pick, Commit: "1234"},
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
			todoToMoveUp: Todo{Hash: "abcd"},
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
			todoToMoveUp:  Todo{Hash: "def0"},
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
			todoToMoveUp:  Todo{Hash: "abcd"},
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
			todoToMoveUp:  Todo{Hash: "5678"},
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
		originalHash  string
		fixupHash     string
		changeToFixup bool
		expectedTodos []todo.Todo
		expectedErr   error
	}{
		{
			name: "fixup commit is the last commit (change to fixup)",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: true,
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Fixup, Commit: "fixup"},
			},
			expectedErr: nil,
		},
		{
			name: "fixup commit is the last commit (don't change to fixup)",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: false,
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			expectedErr: nil,
		},
		{
			name: "fixup commit is separated from original commit",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "other"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: true,
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Fixup, Commit: "fixup"},
				{Command: todo.Pick, Commit: "other"},
			},
			expectedErr: nil,
		},
		{
			name: "fixup commit is separated from original merge commit",
			todos: []todo.Todo{
				{Command: todo.Merge, Commit: "original"},
				{Command: todo.Pick, Commit: "other"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: true,
			expectedTodos: []todo.Todo{
				{Command: todo.Merge, Commit: "original"},
				{Command: todo.Fixup, Commit: "fixup"},
				{Command: todo.Pick, Commit: "other"},
			},
			expectedErr: nil,
		},
		{
			name: "More original hashes than expected",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: true,
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one original hash, found 2"),
		},
		{
			name: "More fixup hashes than expected",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
				{Command: todo.Pick, Commit: "fixup"},
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: true,
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one fixup hash, found 2"),
		},
		{
			name: "No fixup hashes found",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "original"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: true,
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one fixup hash, found 0"),
		},
		{
			name: "No original hashes found",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "fixup"},
			},
			originalHash:  "original",
			fixupHash:     "fixup",
			changeToFixup: true,
			expectedTodos: nil,
			expectedErr:   errors.New("Expected exactly one original hash, found 0"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actualTodos, actualErr := moveFixupCommitDown(scenario.todos, scenario.originalHash, scenario.fixupHash, scenario.changeToFixup)

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
				{Ref: "refs/heads/some_branch"},
				{Hash: "abcd"},
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
				{Hash: "abcd"},
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

func TestRebaseCommands_changeTodoCommit(t *testing.T) {
	scenarios := []struct {
		name          string
		todos         []todo.Todo
		hash          string
		newHash       string
		expectedTodos []todo.Todo
		expectedErr   error
	}{
		{
			name: "success",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234", Msg: "first"},
				{Command: todo.Pick, Commit: "5678", Msg: "second"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "abcd", Msg: "third"},
			},
			hash:    "5678",
			newHash: "ef01",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234", Msg: "first"},
				{Command: todo.Pick, Commit: "ef01", Msg: "second"},
				{Command: todo.UpdateRef, Ref: "refs/heads/some_branch"},
				{Command: todo.Pick, Commit: "abcd", Msg: "third"},
			},
			expectedErr: nil,
		},
		{
			name: "abbreviated hash in todo",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
			},
			hash:    "56789abc",
			newHash: "ef01",
			expectedTodos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "ef01"},
			},
			expectedErr: nil,
		},
		{
			name: "failure",
			todos: []todo.Todo{
				{Command: todo.Pick, Commit: "1234"},
				{Command: todo.Pick, Commit: "5678"},
			},
			hash:          "abcd",
			newHash:       "ef01",
			expectedTodos: nil,
			expectedErr:   errors.New("Todo abcd not found in git-rebase-todo"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			actualTodos, actualErr := changeTodoCommit(scenario.todos, scenario.hash, scenario.newHash)

			if scenario.expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.EqualError(t, actualErr, scenario.expectedErr.Error())
			}

			assert.EqualValues(t, scenario.expectedTodos, actualTodos)
		})
	}
}

func Test_equalHash(t *testing.T) {
	scenarios := []struct {
		a        string
		b        string
		expected bool
	}{
		{"", "", true},
		{"", "123", false},
		{"123", "", false},
		{"123", "123", true},
		{"123", "123abc", true},
		{"123abc", "123", true},
		{"123", "a", false},
		{"1", "abc", false},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("'%s' vs. '%s'", scenario.a, scenario.b), func(t *testing.T) {
			assert.Equal(t, scenario.expected, equalHash(scenario.a, scenario.b))
		})
	}
}
