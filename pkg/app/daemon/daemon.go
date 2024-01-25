package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// Sometimes lazygit will be invoked in daemon mode from a parent lazygit process.
// We do this when git lets us supply a program to run within a git command.
// For example, if we want to ensure that a git command doesn't hang due to
// waiting for an editor to save a commit message, we can tell git to invoke lazygit
// as the editor via 'GIT_EDITOR=lazygit', and use the env var
// 'LAZYGIT_DAEMON_KIND=1' (exit immediately) to specify that we want to run lazygit
// as a daemon which simply exits immediately.
//
// 'Daemon' is not the best name for this, because it's not a persistent background
// process, but it's close enough.

type DaemonKind int

const (
	// for when we fail to parse the daemon kind
	DaemonKindUnknown DaemonKind = iota

	DaemonKindExitImmediately
	DaemonKindCherryPick
	DaemonKindMoveTodosUp
	DaemonKindMoveTodosDown
	DaemonKindInsertBreak
	DaemonKindChangeTodoActions
	DaemonKindMoveFixupCommitDown
)

const (
	DaemonKindEnvKey string = "LAZYGIT_DAEMON_KIND"

	// Contains json-encoded arguments to the daemon
	DaemonInstructionEnvKey string = "LAZYGIT_DAEMON_INSTRUCTION"
)

func getInstruction() Instruction {
	jsonData := os.Getenv(DaemonInstructionEnvKey)

	mapping := map[DaemonKind]func(string) Instruction{
		DaemonKindExitImmediately:     deserializeInstruction[*ExitImmediatelyInstruction],
		DaemonKindCherryPick:          deserializeInstruction[*CherryPickCommitsInstruction],
		DaemonKindChangeTodoActions:   deserializeInstruction[*ChangeTodoActionsInstruction],
		DaemonKindMoveFixupCommitDown: deserializeInstruction[*MoveFixupCommitDownInstruction],
		DaemonKindMoveTodosUp:         deserializeInstruction[*MoveTodosUpInstruction],
		DaemonKindMoveTodosDown:       deserializeInstruction[*MoveTodosDownInstruction],
		DaemonKindInsertBreak:         deserializeInstruction[*InsertBreakInstruction],
	}

	return mapping[getDaemonKind()](jsonData)
}

func Handle(common *common.Common) {
	if !InDaemonMode() {
		return
	}

	instruction := getInstruction()

	if err := instruction.run(common); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func InDaemonMode() bool {
	return getDaemonKind() != DaemonKindUnknown
}

func getDaemonKind() DaemonKind {
	intValue, err := strconv.Atoi(os.Getenv(DaemonKindEnvKey))
	if err != nil {
		return DaemonKindUnknown
	}

	return DaemonKind(intValue)
}

func getCommentChar() byte {
	cmd := exec.Command("git", "config", "--get", "--null", "core.commentChar")
	if output, err := cmd.Output(); err == nil && len(output) == 2 {
		return output[0]
	}

	return '#'
}

// An Instruction is a command to be run by lazygit in daemon mode.
// It is serialized to json and passed to lazygit via environment variables
type Instruction interface {
	Kind() DaemonKind
	SerializedInstructions() string

	// runs the instruction
	run(common *common.Common) error
}

func serializeInstruction[T any](instruction T) string {
	jsonData, err := json.Marshal(instruction)
	if err != nil {
		// this should never happen
		panic(err)
	}

	return string(jsonData)
}

func deserializeInstruction[T Instruction](jsonData string) Instruction {
	var instruction T
	err := json.Unmarshal([]byte(jsonData), &instruction)
	if err != nil {
		panic(err)
	}

	return instruction
}

func ToEnvVars(instruction Instruction) []string {
	return []string{
		fmt.Sprintf("%s=%d", DaemonKindEnvKey, instruction.Kind()),
		fmt.Sprintf("%s=%s", DaemonInstructionEnvKey, instruction.SerializedInstructions()),
	}
}

type ExitImmediatelyInstruction struct{}

func (self *ExitImmediatelyInstruction) Kind() DaemonKind {
	return DaemonKindExitImmediately
}

func (self *ExitImmediatelyInstruction) SerializedInstructions() string {
	return serializeInstruction(self)
}

func (self *ExitImmediatelyInstruction) run(common *common.Common) error {
	return nil
}

func NewExitImmediatelyInstruction() Instruction {
	return &ExitImmediatelyInstruction{}
}

type CherryPickCommitsInstruction struct {
	Todo string
}

func NewCherryPickCommitsInstruction(commits []*models.Commit) Instruction {
	todoLines := lo.Map(commits, func(commit *models.Commit, _ int) TodoLine {
		return TodoLine{
			Action: "pick",
			Commit: commit,
		}
	})

	todo := TodoLinesToString(todoLines)

	return &CherryPickCommitsInstruction{
		Todo: todo,
	}
}

func (self *CherryPickCommitsInstruction) Kind() DaemonKind {
	return DaemonKindCherryPick
}

func (self *CherryPickCommitsInstruction) SerializedInstructions() string {
	return serializeInstruction(self)
}

func (self *CherryPickCommitsInstruction) run(common *common.Common) error {
	return handleInteractiveRebase(common, func(path string) error {
		return utils.PrependStrToTodoFile(path, []byte(self.Todo))
	})
}

type ChangeTodoActionsInstruction struct {
	Changes []ChangeTodoAction
}

func NewChangeTodoActionsInstruction(changes []ChangeTodoAction) Instruction {
	return &ChangeTodoActionsInstruction{
		Changes: changes,
	}
}

func (self *ChangeTodoActionsInstruction) Kind() DaemonKind {
	return DaemonKindChangeTodoActions
}

func (self *ChangeTodoActionsInstruction) SerializedInstructions() string {
	return serializeInstruction(self)
}

func (self *ChangeTodoActionsInstruction) run(common *common.Common) error {
	return handleInteractiveRebase(common, func(path string) error {
		changes := lo.Map(self.Changes, func(c ChangeTodoAction, _ int) utils.TodoChange {
			return utils.TodoChange{
				Sha:       c.Sha,
				OldAction: todo.Pick,
				NewAction: c.NewAction,
			}
		})

		return utils.EditRebaseTodo(path, changes, getCommentChar())
	})
}

// Takes the sha of some commit, and the sha of a fixup commit that was created
// at the end of the branch, then moves the fixup commit down to right after the
// original commit, changing its type to "fixup"
type MoveFixupCommitDownInstruction struct {
	OriginalSha string
	FixupSha    string
}

func NewMoveFixupCommitDownInstruction(originalSha string, fixupSha string) Instruction {
	return &MoveFixupCommitDownInstruction{
		OriginalSha: originalSha,
		FixupSha:    fixupSha,
	}
}

func (self *MoveFixupCommitDownInstruction) Kind() DaemonKind {
	return DaemonKindMoveFixupCommitDown
}

func (self *MoveFixupCommitDownInstruction) SerializedInstructions() string {
	return serializeInstruction(self)
}

func (self *MoveFixupCommitDownInstruction) run(common *common.Common) error {
	return handleInteractiveRebase(common, func(path string) error {
		return utils.MoveFixupCommitDown(path, self.OriginalSha, self.FixupSha, getCommentChar())
	})
}

type MoveTodosUpInstruction struct {
	Shas []string
}

func NewMoveTodosUpInstruction(shas []string) Instruction {
	return &MoveTodosUpInstruction{
		Shas: shas,
	}
}

func (self *MoveTodosUpInstruction) Kind() DaemonKind {
	return DaemonKindMoveTodosUp
}

func (self *MoveTodosUpInstruction) SerializedInstructions() string {
	return serializeInstruction(self)
}

func (self *MoveTodosUpInstruction) run(common *common.Common) error {
	todosToMove := lo.Map(self.Shas, func(sha string, _ int) utils.Todo {
		return utils.Todo{
			Sha:    sha,
			Action: todo.Pick,
		}
	})

	return handleInteractiveRebase(common, func(path string) error {
		return utils.MoveTodosUp(path, todosToMove, getCommentChar())
	})
}

type MoveTodosDownInstruction struct {
	Shas []string
}

func NewMoveTodosDownInstruction(shas []string) Instruction {
	return &MoveTodosDownInstruction{
		Shas: shas,
	}
}

func (self *MoveTodosDownInstruction) Kind() DaemonKind {
	return DaemonKindMoveTodosDown
}

func (self *MoveTodosDownInstruction) SerializedInstructions() string {
	return serializeInstruction(self)
}

func (self *MoveTodosDownInstruction) run(common *common.Common) error {
	todosToMove := lo.Map(self.Shas, func(sha string, _ int) utils.Todo {
		return utils.Todo{
			Sha:    sha,
			Action: todo.Pick,
		}
	})

	return handleInteractiveRebase(common, func(path string) error {
		return utils.MoveTodosDown(path, todosToMove, getCommentChar())
	})
}

type InsertBreakInstruction struct{}

func NewInsertBreakInstruction() Instruction {
	return &InsertBreakInstruction{}
}

func (self *InsertBreakInstruction) Kind() DaemonKind {
	return DaemonKindInsertBreak
}

func (self *InsertBreakInstruction) SerializedInstructions() string {
	return serializeInstruction(self)
}

func (self *InsertBreakInstruction) run(common *common.Common) error {
	return handleInteractiveRebase(common, func(path string) error {
		return utils.PrependStrToTodoFile(path, []byte("break\n"))
	})
}
