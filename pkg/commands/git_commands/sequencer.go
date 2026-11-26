package git_commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
)

// sequencerTodoState tells what has happened so far with the todo that a
// cherry-pick or revert of several commits stopped at. This is the first todo
// in .git/sequencer/todo.
type sequencerTodoState int

const (
	// Git started to apply the todo, and stopped because of conflicts (or
	// because it couldn't commit the result). CHERRY_PICK_HEAD or REVERT_HEAD
	// exists, and continuing commits the changes.
	sequencerTodoStarted sequencerTodoState = iota
	// Git stopped before it started to apply the todo, for example because the
	// commit would overwrite an untracked file.
	sequencerTodoNotStarted
	// The user has committed the result of the todo themselves, for example
	// after resolving its conflicts. Continuing skips it.
	sequencerTodoCommitted
)

// getCurrentSequencerTodoState returns the state of the todo that the
// cherry-pick or revert of several commits that is in progress in gitDir
// stopped at.
//
// Without CHERRY_PICK_HEAD or REVERT_HEAD, we can tell whether the user has
// committed the todo by whether HEAD has moved since git last tried to apply a
// commit, which git records in .git/sequencer/abort-safety. Git does the same
// when it decides whether "--skip" is possible (see rollback_is_safe in git's
// sequencer.c).
func getCurrentSequencerTodoState(cmd oscommands.ICmdObjBuilder, gitDir string) sequencerTodoState {
	for _, file := range []string{"CHERRY_PICK_HEAD", "REVERT_HEAD"} {
		if _, err := os.Stat(filepath.Join(gitDir, file)); err == nil {
			return sequencerTodoStarted
		}
	}

	// If the file doesn't exist, git compares HEAD to the null hash; an empty
	// string does the same for us
	abortSafety, _ := os.ReadFile(filepath.Join(gitDir, "sequencer", "abort-safety"))

	// If HEAD is unborn, this fails and returns an empty string
	head, _ := cmd.New(NewGitCmd("rev-parse").Arg("--verify", "--quiet", "HEAD").ToArgv()).DontLog().RunWithOutput()

	if strings.TrimSpace(head) == strings.TrimSpace(string(abortSafety)) {
		return sequencerTodoNotStarted
	}
	return sequencerTodoCommitted
}

// repeatCurrentSequencerTodo prepares the cherry-pick or revert of several
// commits that is in progress in gitDir for "--continue". If git stopped it
// before it started to apply the current todo, "--continue" would skip that
// todo, because git can't tell it apart from a todo that the user has
// committed themselves. So we put a copy of the todo in front of it, and git
// skips the copy instead.
//
// The returned function removes the copy again if git hasn't got as far as
// skipping it, for example because it refused to continue.
func repeatCurrentSequencerTodo(cmd oscommands.ICmdObjBuilder, gitDir string, commentChar byte) (func(), error) {
	noop := func() {}

	todoPath := filepath.Join(gitDir, "sequencer", "todo")
	content, err := os.ReadFile(todoPath)
	if err != nil {
		// No cherry-pick or revert of several commits is in progress
		return noop, nil
	}

	if getCurrentSequencerTodoState(cmd, gitDir) != sequencerTodoNotStarted {
		return noop, nil
	}

	currentTodo, found := lo.Find(strings.Split(string(content), "\n"), func(line string) bool {
		line = strings.TrimSpace(line)
		return line != "" && line[0] != commentChar
	})
	if !found {
		return noop, nil
	}

	newContent := append([]byte(currentTodo+"\n"), content...)
	if err := os.WriteFile(todoPath, newContent, 0o644); err != nil {
		return nil, err
	}

	return func() {
		if currentContent, err := os.ReadFile(todoPath); err == nil && bytes.Equal(currentContent, newContent) {
			_ = os.WriteFile(todoPath, content, 0o644)
		}
	}, nil
}
