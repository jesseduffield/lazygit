package gui

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestStartCmdWithPipeWhenPipeCannotBeCreated(t *testing.T) {
	cmd := exec.Command("non-existent-command")
	// Assigning stdout up front makes cmd.StdoutPipe fail. This happens in
	// practice on the Unix pty fallback path: a failed pty start can leave
	// the tty assigned to the command's stdout.
	cmd.Stdout = &bytes.Buffer{}

	_, r := startCmdWithPipe(cmd, utils.NewDummyLog())

	// NewCmdTask's scanner panics on a nil reader, so startCmdWithPipe must
	// not return one even when it can't create the pipe.
	/* EXPECTED:
	assert.NotNil(t, r)
	ACTUAL: */
	assert.Nil(t, r)
}
