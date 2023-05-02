package todo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrUnexpectedCommand = errors.New("unexpected command")
	ErrMissingLabel      = errors.New("missing label")
	ErrMissingCommit     = errors.New("missing commit")
	ErrMissingExecCmd    = errors.New("missing command for exec")
	ErrMissingRef        = errors.New("missing ref")
)

func Parse(f io.Reader) ([]Todo, error) {
	var result []Todo

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		cmd, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %q: %w", line, err)
		}

		result = append(result, cmd)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	return result, nil
}

func parseLine(line string) (Todo, error) {
	var todo Todo

	if strings.HasPrefix(line, CommentChar) {
		todo.Command = Comment
		todo.Comment = strings.TrimLeft(line, CommentChar)
		return todo, nil
	}

	fields := strings.Fields(line)

	var commandLen int
	for i := Pick; i < Comment; i++ {
		if isCommand(i, fields[0]) {
			todo.Command = i
			commandLen = len(fields[0])
			fields = fields[1:]
			break
		}
	}

	if todo.Command == 0 {
		// unexpected command
		return todo, ErrUnexpectedCommand
	}

	if todo.Command == Break {
		return todo, nil
	}

	if todo.Command == Label || todo.Command == Reset {
		restOfLine := strings.TrimSpace(line[commandLen:])
		if todo.Command == Reset && restOfLine == "[new root]" {
			todo.Label = restOfLine
		} else if len(fields) == 0 {
			return todo, ErrMissingLabel
		} else {
			todo.Label = fields[0]
		}
		return todo, nil
	}

	if todo.Command == Exec {
		if len(fields) == 0 {
			return todo, ErrMissingExecCmd
		}
		todo.ExecCommand = strings.Join(fields, " ")
		return todo, nil
	}

	if todo.Command == Merge {
		if fields[0] == "-C" || fields[0] == "-c" {
			todo.Flag = fields[0]
			fields = fields[1:]
			if len(fields) == 0 {
				return todo, ErrMissingCommit
			}
			todo.Commit = fields[0]
			fields = fields[1:]
		}
		if len(fields) == 0 {
			return todo, ErrMissingLabel
		}
		todo.Label = fields[0]
		fields = fields[1:]
		if fields[0] == "#" {
			fields = fields[1:]
			todo.Msg = strings.Join(fields, " ")
		}
		return todo, nil
	}

	if todo.Command == Fixup {
		if len(fields) == 0 {
			return todo, ErrMissingCommit
		}
		// Skip flags
		if fields[0] == "-C" || fields[0] == "-c" {
			todo.Flag = fields[0]
			fields = fields[1:]
		}
	}

	if todo.Command == UpdateRef {
		if len(fields) == 0 {
			return todo, ErrMissingRef
		}
		todo.Ref = fields[0]
		return todo, nil
	}

	if len(fields) == 0 {
		return todo, ErrMissingCommit
	}

	todo.Commit = fields[0]
	fields = fields[1:]

	// Trim # and whitespace
	todo.Msg = strings.TrimPrefix(strings.Join(fields, " "), CommentChar+" ")

	return todo, nil
}

func isCommand(i TodoCommand, s string) bool {
	if i < 0 || i > Comment {
		return false
	}
	return len(s) > 0 &&
		(todoCommandInfo[i].cmd == s || todoCommandInfo[i].nickname == s)
}
