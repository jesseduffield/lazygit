package daemon

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Sometimes lazygit will be invoked in daemon mode from a parent lazygit process.
// We do this when git lets us supply a program to run within a git command.
// For example, if we want to ensure that a git command doesn't hang due to
// waiting for an editor to save a commit message, we can tell git to invoke lazygit
// as the editor via 'GIT_EDITOR=lazygit', and use the env var
// 'LAZYGIT_DAEMON_KIND=EXIT_IMMEDIATELY' to specify that we want to run lazygit
// as a daemon which simply exits immediately. Any additional arguments we want
// to pass to a daemon can be done via other env vars.

type DaemonKind string

const (
	InteractiveRebase DaemonKind = "INTERACTIVE_REBASE"
	ExitImmediately   DaemonKind = "EXIT_IMMEDIATELY"
)

const (
	DaemonKindEnvKey string = "LAZYGIT_DAEMON_KIND"

	// Contains a json-encoded instance of the InteractiveRebaseInstructions struct
	InteractiveRebaseInstructionsEnvKey string = "LAZYGIT_DAEMON_INSTRUCTIONS"
)

// Exactly one of the fields in this struct is expected to be non-empty
type InteractiveRebaseInstructions struct {
	// If this is non-empty, this string is prepended to the git-rebase-todo
	// file. The string is expected to have newlines at the end of each line.
	LinesToPrependToRebaseTODO string

	// If this is non-empty, it tells lazygit to read the original todo file, and
	// change the action for one or more entries in it.
	// The existing action of the todo to be changed is expected to be "pick".
	ChangeTodoActions []ChangeTodoAction

	// Can be set to the sha of a "pick" todo that will be moved down by one.
	ShaToMoveDown string

	// Can be set to the sha of a "pick" todo that will be moved up by one.
	ShaToMoveUp string
}

type ChangeTodoAction struct {
	Sha       string
	NewAction todo.TodoCommand
}

type Daemon interface {
	Run() error
}

func Handle(common *common.Common) {
	d := getDaemon(common)
	if d == nil {
		return
	}

	if err := d.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func InDaemonMode() bool {
	return getDaemonKind() != ""
}

func getDaemon(common *common.Common) Daemon {
	switch getDaemonKind() {
	case InteractiveRebase:
		return &rebaseDaemon{c: common}
	case ExitImmediately:
		return &exitImmediatelyDaemon{c: common}
	}

	return nil
}

func getDaemonKind() DaemonKind {
	return DaemonKind(os.Getenv(DaemonKindEnvKey))
}

type rebaseDaemon struct {
	c *common.Common
}

func (self *rebaseDaemon) Run() error {
	self.c.Log.Info("Lazygit invoked as interactive rebase demon")
	self.c.Log.Info("args: ", os.Args)
	path := os.Args[1]

	if strings.HasSuffix(path, "git-rebase-todo") {
		return self.writeTodoFile(path)
	} else if strings.HasSuffix(path, filepath.Join(gitDir(), "COMMIT_EDITMSG")) { // TODO: test
		// if we are rebasing and squashing, we'll see a COMMIT_EDITMSG
		// but in this case we don't need to edit it, so we'll just return
	} else {
		self.c.Log.Info("Lazygit demon did not match on any use cases")
	}

	return nil
}

func (self *rebaseDaemon) writeTodoFile(path string) error {
	jsonData := os.Getenv(InteractiveRebaseInstructionsEnvKey)
	instructions := InteractiveRebaseInstructions{}
	err := json.Unmarshal([]byte(jsonData), &instructions)
	if err != nil {
		return err
	}

	if instructions.LinesToPrependToRebaseTODO != "" {
		return utils.PrependStrToTodoFile(path, []byte(instructions.LinesToPrependToRebaseTODO))
	} else if len(instructions.ChangeTodoActions) != 0 {
		return self.changeTodoAction(path, instructions.ChangeTodoActions)
	} else if instructions.ShaToMoveDown != "" {
		return utils.MoveTodoDown(path, instructions.ShaToMoveDown, todo.Pick)
	} else if instructions.ShaToMoveUp != "" {
		return utils.MoveTodoUp(path, instructions.ShaToMoveUp, todo.Pick)
	}

	self.c.Log.Error("No instructions were given to daemon")
	return nil
}

func (self *rebaseDaemon) changeTodoAction(path string, changeTodoActions []ChangeTodoAction) error {
	for _, c := range changeTodoActions {
		if err := utils.EditRebaseTodo(path, c.Sha, todo.Pick, c.NewAction); err != nil {
			return err
		}
	}

	return nil
}

func gitDir() string {
	dir := env.GetGitDirEnv()
	if dir == "" {
		return ".git"
	}
	return dir
}

type exitImmediatelyDaemon struct {
	c *common.Common
}

func (self *exitImmediatelyDaemon) Run() error {
	return nil
}
