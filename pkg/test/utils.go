package test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/mgutz/str"
	"github.com/stretchr/testify/assert"
)

// CommandSwapper takes a command, verifies that it is what it's expected to be
// and then returns a replacement command that will actually be called by the os
type CommandSwapper struct {
	Expect  string
	Replace string
}

// SwapCommand verifies the command is what we expected, and swaps it out for a different command
func (i *CommandSwapper) SwapCommand(t *testing.T, cmd string, args []string) *exec.Cmd {
	splitCmd := str.ToArgv(i.Expect)
	assert.EqualValues(t, splitCmd[0], cmd, fmt.Sprintf("received command: %s %s", cmd, strings.Join(args, " ")))
	if len(splitCmd) > 1 {
		assert.EqualValues(t, splitCmd[1:], args, fmt.Sprintf("received command: %s %s", cmd, strings.Join(args, " ")))
	}

	splitCmd = str.ToArgv(i.Replace)
	return exec.Command(splitCmd[0], splitCmd[1:]...)
}

// CreateMockCommand creates a command function that will verify its receiving the right sequence of commands from lazygit
func CreateMockCommand(t *testing.T, swappers []*CommandSwapper) func(cmd string, args ...string) *exec.Cmd {
	commandIndex := 0

	return func(cmd string, args ...string) *exec.Cmd {
		var command *exec.Cmd
		if commandIndex > len(swappers)-1 {
			assert.Fail(t, fmt.Sprintf("too many commands run. This command was (%s %s)", cmd, strings.Join(args, " ")))
		}
		command = swappers[commandIndex].SwapCommand(t, cmd, args)
		commandIndex++
		return command
	}
}
