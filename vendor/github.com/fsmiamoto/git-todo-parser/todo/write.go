package todo

import (
	"io"
	"strings"
)

func Write(f io.Writer, todos []Todo, commentChar byte) error {
	for _, todo := range todos {
		if err := writeTodo(f, todo, commentChar); err != nil {
			return err
		}
	}

	return nil
}

func writeTodo(f io.Writer, todo Todo, commentChar byte) error {
	var sb strings.Builder
	if todo.Command != Comment {
		sb.WriteString(todo.Command.String())
	}

	switch todo.Command {
	case NoOp:
		return nil

	case Comment:
		sb.WriteByte(commentChar)
		sb.WriteString(todo.Comment)

	case Break:

	case Label:
		fallthrough
	case Reset:
		sb.WriteByte(' ')
		sb.WriteString(todo.Label)

	case Exec:
		sb.WriteByte(' ')
		sb.WriteString(todo.ExecCommand)

	case Merge:
		sb.WriteByte(' ')
		if todo.Commit != "" {
			sb.WriteString(todo.Flag)
			sb.WriteByte(' ')
			sb.WriteString(todo.Commit)
			sb.WriteByte(' ')
		}
		sb.WriteString(todo.Label)
		if todo.Msg != "" {
			sb.WriteString(" # ")
			sb.WriteString(todo.Msg)
		}

	case Fixup:
		sb.WriteByte(' ')
		if todo.Flag != "" {
			sb.WriteString(todo.Flag)
			sb.WriteByte(' ')
		}
		sb.WriteString(todo.Commit)

	case UpdateRef:
		sb.WriteByte(' ')
		sb.WriteString(todo.Ref)

	case Pick:
		fallthrough
	case Revert:
		fallthrough
	case Edit:
		fallthrough
	case Reword:
		fallthrough
	case Squash:
		fallthrough
	case Drop:
		sb.WriteByte(' ')
		sb.WriteString(todo.Commit)
		if todo.Msg != "" {
			sb.WriteByte(' ')
			sb.WriteString(todo.Msg)
		}
	}

	sb.WriteByte('\n')
	_, err := f.Write([]byte(sb.String()))
	return err
}
