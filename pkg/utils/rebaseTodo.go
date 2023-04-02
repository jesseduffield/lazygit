package utils

import (
	"os"

	"github.com/fsmiamoto/git-todo-parser/todo"
)

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
