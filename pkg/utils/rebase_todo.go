package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/samber/lo"
)

func equalShas(a, b string) bool {
	return strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
}

func ReadRebaseTodoFile(fileName string) ([]todo.Todo, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	todos, err := todo.Parse(f)
	err2 := f.Close()
	if err == nil {
		err = err2
	}
	return todos, err
}

func WriteRebaseTodoFile(fileName string, todos []todo.Todo) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = todo.Write(f, todos)
	err2 := f.Close()
	if err == nil {
		err = err2
	}
	return err
}

func MoveTodoDown(fileName string, sha string, action todo.TodoCommand) error {
	todos, err := ReadRebaseTodoFile(fileName)
	if err != nil {
		return err
	}
	rearrangedTodos, err := moveTodoDown(todos, sha, action)
	if err != nil {
		return err
	}
	return WriteRebaseTodoFile(fileName, rearrangedTodos)
}

func MoveTodoUp(fileName string, sha string, action todo.TodoCommand) error {
	todos, err := ReadRebaseTodoFile(fileName)
	if err != nil {
		return err
	}
	rearrangedTodos, err := moveTodoUp(todos, sha, action)
	if err != nil {
		return err
	}
	return WriteRebaseTodoFile(fileName, rearrangedTodos)
}

func moveTodoDown(todos []todo.Todo, sha string, action todo.TodoCommand) ([]todo.Todo, error) {
	rearrangedTodos, err := moveTodoUp(lo.Reverse(todos), sha, action)
	return lo.Reverse(rearrangedTodos), err
}

func moveTodoUp(todos []todo.Todo, sha string, action todo.TodoCommand) ([]todo.Todo, error) {
	_, sourceIdx, ok := lo.FindIndexOf(todos, func(t todo.Todo) bool {
		// Comparing just the sha is not enough; we need to compare both the
		// action and the sha, as the sha could appear multiple times (e.g. in a
		// pick and later in a merge)
		return t.Command == action && equalShas(t.Commit, sha)
	})

	if !ok {
		// Should never happen
		return []todo.Todo{}, fmt.Errorf("Todo %s not found in git-rebase-todo", sha)
	}

	// The todos are ordered backwards compared to our model commits, so
	// actually move the commit _down_ in the todos slice (i.e. towards
	// the end of the slice)

	// Find the next todo that we show in lazygit's commits view (skipping the rest)
	_, skip, ok := lo.FindIndexOf(todos[sourceIdx+1:], isRenderedTodo)

	if !ok {
		// We expect callers to guard against this
		return []todo.Todo{}, fmt.Errorf("Destination position for moving todo is out of range")
	}

	destinationIdx := sourceIdx + 1 + skip

	rearrangedTodos := MoveElement(todos, sourceIdx, destinationIdx)

	return rearrangedTodos, nil
}

// We render a todo in the commits view if it's a commit or if it's an
// update-ref. We don't render label, reset, or comment lines.
func isRenderedTodo(t todo.Todo) bool {
	return t.Commit != "" || t.Command == todo.UpdateRef
}
