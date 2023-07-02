package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/samber/lo"
)

// Read a git-rebase-todo file, change the action for the given sha to
// newAction, and write it back
func EditRebaseTodo(filePath string, sha string, oldAction todo.TodoCommand, newAction todo.TodoCommand, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(filePath, commentChar)
	if err != nil {
		return err
	}

	for i := range todos {
		t := &todos[i]
		// Comparing just the sha is not enough; we need to compare both the
		// action and the sha, as the sha could appear multiple times (e.g. in a
		// pick and later in a merge)
		if t.Command == oldAction && equalShas(t.Commit, sha) {
			t.Command = newAction
			return WriteRebaseTodoFile(filePath, todos, commentChar)
		}
	}

	// Should never get here
	return fmt.Errorf("Todo %s not found in git-rebase-todo", sha)
}

func equalShas(a, b string) bool {
	return strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
}

func ReadRebaseTodoFile(fileName string, commentChar byte) ([]todo.Todo, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	todos, err := todo.Parse(f, commentChar)
	err2 := f.Close()
	if err == nil {
		err = err2
	}
	return todos, err
}

func WriteRebaseTodoFile(fileName string, todos []todo.Todo, commentChar byte) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = todo.Write(f, todos, commentChar)
	err2 := f.Close()
	if err == nil {
		err = err2
	}
	return err
}

func PrependStrToTodoFile(filePath string, linesToPrepend []byte) error {
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	linesToPrepend = append(linesToPrepend, existingContent...)
	return os.WriteFile(filePath, linesToPrepend, 0o644)
}

func MoveTodoDown(fileName string, sha string, action todo.TodoCommand, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return err
	}
	rearrangedTodos, err := moveTodoDown(todos, sha, action)
	if err != nil {
		return err
	}
	return WriteRebaseTodoFile(fileName, rearrangedTodos, commentChar)
}

func MoveTodoUp(fileName string, sha string, action todo.TodoCommand, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return err
	}
	rearrangedTodos, err := moveTodoUp(todos, sha, action)
	if err != nil {
		return err
	}
	return WriteRebaseTodoFile(fileName, rearrangedTodos, commentChar)
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

func MoveFixupCommitDown(fileName string, originalSha string, fixupSha string, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return err
	}

	newTodos, err := moveFixupCommitDown(todos, originalSha, fixupSha)
	if err != nil {
		return err
	}

	return WriteRebaseTodoFile(fileName, newTodos, commentChar)
}

func moveFixupCommitDown(todos []todo.Todo, originalSha string, fixupSha string) ([]todo.Todo, error) {
	isOriginal := func(t todo.Todo) bool {
		return t.Command == todo.Pick && equalShas(t.Commit, originalSha)
	}

	isFixup := func(t todo.Todo) bool {
		return t.Command == todo.Pick && equalShas(t.Commit, fixupSha)
	}

	originalShaCount := lo.CountBy(todos, isOriginal)
	if originalShaCount != 1 {
		return nil, fmt.Errorf("Expected exactly one original SHA, found %d", originalShaCount)
	}

	fixupShaCount := lo.CountBy(todos, isFixup)
	if fixupShaCount != 1 {
		return nil, fmt.Errorf("Expected exactly one fixup SHA, found %d", fixupShaCount)
	}

	_, fixupIndex, _ := lo.FindIndexOf(todos, isFixup)
	_, originalIndex, _ := lo.FindIndexOf(todos, isOriginal)

	newTodos := MoveElement(todos, fixupIndex, originalIndex+1)

	newTodos[originalIndex+1].Command = todo.Fixup

	return newTodos, nil
}

// We render a todo in the commits view if it's a commit or if it's an
// update-ref. We don't render label, reset, or comment lines.
func isRenderedTodo(t todo.Todo) bool {
	return t.Commit != "" || t.Command == todo.UpdateRef
}
