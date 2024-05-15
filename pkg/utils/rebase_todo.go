package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type Todo struct {
	Hash   string // for todos that have one, e.g. pick, drop, fixup, etc.
	Ref    string // for update-ref todos
	Action todo.TodoCommand
}

// In order to change a TODO in git-rebase-todo, we need to specify the old action,
// because sometimes the same hash appears multiple times in the file (e.g. in a pick
// and later in a merge)
type TodoChange struct {
	Hash      string
	OldAction todo.TodoCommand
	NewAction todo.TodoCommand
}

// Read a git-rebase-todo file, change the actions for the given commits,
// and write it back
func EditRebaseTodo(filePath string, changes []TodoChange, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(filePath, commentChar)
	if err != nil {
		return err
	}

	matchCount := 0
	for i := range todos {
		t := &todos[i]
		// This is a nested loop, but it's ok because the number of todos should be small
		for _, change := range changes {
			if t.Command == change.OldAction && equalHash(t.Commit, change.Hash) {
				matchCount++
				t.Command = change.NewAction
			}
		}
	}

	if matchCount < len(changes) {
		// Should never get here
		return errors.New("Some todos not found in git-rebase-todo")
	}

	return WriteRebaseTodoFile(filePath, todos, commentChar)
}

func equalHash(a, b string) bool {
	return strings.HasPrefix(a, b) || strings.HasPrefix(b, a)
}

func findTodo(todos []todo.Todo, todoToFind Todo) (int, bool) {
	_, idx, ok := lo.FindIndexOf(todos, func(t todo.Todo) bool {
		// Comparing just the hash is not enough; we need to compare both the
		// action and the hash, as the hash could appear multiple times (e.g. in a
		// pick and later in a merge). For update-ref todos we also must compare
		// the Ref.
		return t.Command == todoToFind.Action &&
			equalHash(t.Commit, todoToFind.Hash) &&
			t.Ref == todoToFind.Ref
	})
	return idx, ok
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

func todosToString(todos []todo.Todo, commentChar byte) ([]byte, error) {
	buffer := bytes.Buffer{}
	err := todo.Write(&buffer, todos, commentChar)
	return buffer.Bytes(), err
}

func PrependStrToTodoFile(filePath string, linesToPrepend []byte) error {
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	linesToPrepend = append(linesToPrepend, existingContent...)
	return os.WriteFile(filePath, linesToPrepend, 0o644)
}

// Unlike the other functions in this file, which write the changed todos file
// back to disk, this one returns the new content as a byte slice. This is
// because when deleting update-ref todos, we must perform a "git rebase
// --edit-todo" command to pass the changed todos to git so that it can do some
// housekeeping around the deleted todos. This can only be done by our caller.
func DeleteTodos(fileName string, todosToDelete []Todo, commentChar byte) ([]byte, error) {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return nil, err
	}
	rearrangedTodos, err := deleteTodos(todos, todosToDelete)
	if err != nil {
		return nil, err
	}
	return todosToString(rearrangedTodos, commentChar)
}

func deleteTodos(todos []todo.Todo, todosToDelete []Todo) ([]todo.Todo, error) {
	for _, todoToDelete := range todosToDelete {
		idx, ok := findTodo(todos, todoToDelete)

		if !ok {
			// Should never happen
			return []todo.Todo{}, fmt.Errorf("Todo %s not found in git-rebase-todo", todoToDelete.Hash)
		}

		todos = Remove(todos, idx)
	}

	return todos, nil
}

func MoveTodosDown(fileName string, todosToMove []Todo, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return err
	}
	rearrangedTodos, err := moveTodosDown(todos, todosToMove)
	if err != nil {
		return err
	}
	return WriteRebaseTodoFile(fileName, rearrangedTodos, commentChar)
}

func MoveTodosUp(fileName string, todosToMove []Todo, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return err
	}
	rearrangedTodos, err := moveTodosUp(todos, todosToMove)
	if err != nil {
		return err
	}
	return WriteRebaseTodoFile(fileName, rearrangedTodos, commentChar)
}

func moveTodoDown(todos []todo.Todo, todoToMove Todo) ([]todo.Todo, error) {
	rearrangedTodos, err := moveTodoUp(lo.Reverse(todos), todoToMove)
	return lo.Reverse(rearrangedTodos), err
}

func moveTodosDown(todos []todo.Todo, todosToMove []Todo) ([]todo.Todo, error) {
	rearrangedTodos, err := moveTodosUp(lo.Reverse(todos), lo.Reverse(todosToMove))
	return lo.Reverse(rearrangedTodos), err
}

func moveTodoUp(todos []todo.Todo, todoToMove Todo) ([]todo.Todo, error) {
	sourceIdx, ok := findTodo(todos, todoToMove)

	if !ok {
		// Should never happen
		return []todo.Todo{}, fmt.Errorf("Todo %s not found in git-rebase-todo", todoToMove.Hash)
	}

	// The todos are ordered backwards compared to our model commits, so
	// actually move the commit _down_ in the todos slice (i.e. towards
	// the end of the slice)

	// Find the next todo that we show in lazygit's commits view (skipping the rest)
	_, skip, ok := lo.FindIndexOf(todos[sourceIdx+1:], isRenderedTodo)

	if !ok {
		// We expect callers to guard against this
		return []todo.Todo{}, errors.New("Destination position for moving todo is out of range")
	}

	destinationIdx := sourceIdx + 1 + skip

	rearrangedTodos := MoveElement(todos, sourceIdx, destinationIdx)

	return rearrangedTodos, nil
}

func moveTodosUp(todos []todo.Todo, todosToMove []Todo) ([]todo.Todo, error) {
	for _, todoToMove := range todosToMove {
		var newTodos []todo.Todo
		newTodos, err := moveTodoUp(todos, todoToMove)
		if err != nil {
			return nil, err
		}
		todos = newTodos
	}

	return todos, nil
}

func MoveFixupCommitDown(fileName string, originalHash string, fixupHash string, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return err
	}

	newTodos, err := moveFixupCommitDown(todos, originalHash, fixupHash)
	if err != nil {
		return err
	}

	return WriteRebaseTodoFile(fileName, newTodos, commentChar)
}

func moveFixupCommitDown(todos []todo.Todo, originalHash string, fixupHash string) ([]todo.Todo, error) {
	isOriginal := func(t todo.Todo) bool {
		return (t.Command == todo.Pick || t.Command == todo.Merge) && equalHash(t.Commit, originalHash)
	}

	isFixup := func(t todo.Todo) bool {
		return t.Command == todo.Pick && equalHash(t.Commit, fixupHash)
	}

	originalHashCount := lo.CountBy(todos, isOriginal)
	if originalHashCount != 1 {
		return nil, fmt.Errorf("Expected exactly one original hash, found %d", originalHashCount)
	}

	fixupHashCount := lo.CountBy(todos, isFixup)
	if fixupHashCount != 1 {
		return nil, fmt.Errorf("Expected exactly one fixup hash, found %d", fixupHashCount)
	}

	_, fixupIndex, _ := lo.FindIndexOf(todos, isFixup)
	_, originalIndex, _ := lo.FindIndexOf(todos, isOriginal)

	newTodos := MoveElement(todos, fixupIndex, originalIndex+1)

	newTodos[originalIndex+1].Command = todo.Fixup

	return newTodos, nil
}

func RemoveUpdateRefsForCopiedBranch(fileName string, commentChar byte) error {
	todos, err := ReadRebaseTodoFile(fileName, commentChar)
	if err != nil {
		return err
	}

	// Filter out comments
	todos = lo.Filter(todos, func(t todo.Todo, _ int) bool {
		return t.Command != todo.Comment
	})

	// Delete any update-ref todos at the end of the todo list. These are not
	// part of a stack of branches, and so shouldn't be updated. This makes it
	// possible to create a copy of a branch and rebase the copy without
	// affecting the original branch.
	if _, i, found := lo.FindLastIndexOf(todos, func(t todo.Todo) bool {
		return t.Command != todo.UpdateRef
	}); found && i < len(todos)-1 {
		todos = slices.Delete(todos, i+1, len(todos))
		return WriteRebaseTodoFile(fileName, todos, commentChar)
	}

	return nil
}

// We render a todo in the commits view if it's a commit or if it's an
// update-ref. We don't render label, reset, or comment lines.
func isRenderedTodo(t todo.Todo) bool {
	return t.Commit != "" || t.Command == todo.UpdateRef
}
