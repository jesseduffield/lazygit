package oscommands

import (
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTmuxTestOSCommand wires a dummy OSCommand to a fake command runner and a
// controllable environment lookup so we can assert on the exact tmux
// invocations without shelling out to a real tmux server.
func newTmuxTestOSCommand(runner *FakeCmdObjRunner, env map[string]string) *OSCommand {
	osCommand := NewDummyOSCommandWithDeps(OSCommandDeps{
		GetenvFn: func(key string) string { return env[key] },
	})
	// NewDummyOSCommandWithDeps does not wire up Cmd, so attach the fake runner
	// explicitly; every tmux invocation goes through it.
	osCommand.Cmd = NewDummyCmdObjBuilder(runner)
	return osCommand
}

// TestInTmuxRequiresEnvVar verifies the first half of the detection guard:
// without the TMUX environment variable we must report that we are not in tmux,
// regardless of whether the tmux binary happens to be installed.
func TestInTmuxRequiresEnvVar(t *testing.T) {
	osCommand := newTmuxTestOSCommand(NewFakeRunner(t), map[string]string{})

	assert.False(t, osCommand.inTmux())
}

// TestInTmuxWithEnvVarTracksBinary verifies the second half of the guard: with
// the TMUX environment variable set, the result must follow whether the tmux
// binary is resolvable on PATH. We compute the expectation from the same
// lookup the implementation uses so the test is deterministic in any
// environment, with or without tmux installed.
func TestInTmuxWithEnvVarTracksBinary(t *testing.T) {
	osCommand := newTmuxTestOSCommand(NewFakeRunner(t), map[string]string{tmuxEnvVar: "/tmp/tmux-1000/default,1234,0"})

	_, lookErr := exec.LookPath(tmuxBinaryName)
	tmuxBinaryPresent := lookErr == nil

	assert.Equal(t, tmuxBinaryPresent, osCommand.inTmux())
}

// TestCopyToClipboardTmux asserts that copying loads the text into a tmux
// buffer with the set-clipboard flag, streaming the payload in on stdin.
func TestCopyToClipboardTmux(t *testing.T) {
	const payload = "branch-name-to-copy"

	runner := NewFakeRunner(t)
	runner.ExpectFunc(
		"tmux load-buffer with set-clipboard flag and payload on stdin",
		func(cmdObj *CmdObj) bool {
			expectedArgs := []string{tmuxBinaryName, tmuxLoadBufferSubcommand, tmuxSetClipboardFlag, tmuxStdioArg}
			if !assert.ObjectsAreEqual(expectedArgs, cmdObj.GetCmd().Args) {
				return false
			}

			// The buffer contents must be piped in on stdin verbatim.
			stdin, err := io.ReadAll(cmdObj.GetCmd().Stdin)
			assert.NoError(t, err)
			return string(stdin) == payload
		},
		"",
		nil,
	)

	osCommand := newTmuxTestOSCommand(runner, map[string]string{tmuxEnvVar: "session"})

	assert.NoError(t, osCommand.copyToClipboardTmux(payload))
	runner.CheckForMissingCalls()
}

// TestPasteFromClipboardTmux asserts that pasting saves the tmux buffer to
// stdout and returns its contents.
func TestPasteFromClipboardTmux(t *testing.T) {
	const bufferContents = "previously-copied-text"

	runner := NewFakeRunner(t)
	runner.ExpectArgs(
		[]string{tmuxBinaryName, tmuxSaveBufferSubcommand, tmuxStdioArg},
		bufferContents,
		nil,
	)

	osCommand := newTmuxTestOSCommand(runner, map[string]string{tmuxEnvVar: "session"})

	pasted, err := osCommand.pasteFromClipboardTmux()
	assert.NoError(t, err)
	assert.Equal(t, bufferContents, pasted)
	runner.CheckForMissingCalls()
}
