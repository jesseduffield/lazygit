package oscommands

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/samber/lo"
)

// Pty is the master side of a pseudo-terminal running a subprocess. The
// concrete implementation is platform-specific: creack/pty on Unix and
// ConPTY on Windows.
type Pty interface {
	io.ReadWriteCloser
	Resize(cols, rows uint16) error
}

// StartedPty is the result of StartPty.
type StartedPty struct {
	// Pty is the master side of the pseudo-terminal; read from it to get
	// the child's combined stdout/stderr and write to it to feed stdin.
	Pty Pty
	// Process is the spawned child. Useful for signalling; on Windows the
	// original *exec.Cmd was not Start()ed (ConPTY spawns via
	// CreateProcess, not os/exec) so cmd.Process is nil and this is the
	// only handle.
	Process *os.Process
	// Wait blocks until the child exits and returns a non-nil error on a
	// nonzero exit status, matching *exec.Cmd.Wait semantics.
	Wait func() error
}

// StartPty runs cmd in a pseudo-terminal with the given initial dimensions.
// Implemented per-platform in pty_unix.go / pty_windows.go.
//
// func StartPty(cmd *exec.Cmd, cols, rows uint16) (StartedPty, error)

// renderWithoutPtyEnvVar makes a render take the piped path on a platform that
// would otherwise use a pty, so that tests can exercise it anywhere.
const renderWithoutPtyEnvVar = "LAZYGIT_RENDER_WITHOUT_PTY"

// RendersThroughAPipe reports whether a render feeds the diff renderer the
// command's output through a pipe rather than running it in a pty.
//
// On Windows it has to. ConPTY doesn't pass a command's output through; it
// parses it into a screen buffer and re-encodes that for the terminal side,
// and it hands a sequence it can't represent there the moment it parses it,
// separately from the text around it. So what a renderer writes is not what
// lazygit reads. A pipe carries the bytes as the renderer wrote them.
//
// Everywhere else the pty is kept, since a renderer can read the width it
// should lay out to off it, and a configuration that doesn't name a width would
// otherwise render at whatever width the renderer falls back to.
func RendersThroughAPipe() bool {
	return runtime.GOOS == "windows" || os.Getenv(renderWithoutPtyEnvVar) != ""
}

// RunInPtyWithOutput runs cmd in a pseudo-terminal of the given size and returns
// everything it wrote. Commands that render for a terminal need one to render at
// all: git only pipes its output through a pager when it thinks it is talking to
// a terminal, and the pager itself commonly decides whether to use color the same
// way.
func RunInPtyWithOutput(cmd *exec.Cmd, cols uint16, rows uint16) (string, error) {
	startedPty, err := StartPty(cmd, cols, rows)
	if err != nil {
		return "", err
	}
	defer startedPty.Pty.Close()

	// Reading from the master side of a pty fails as soon as the child has closed
	// the other side, so an error here only tells us that the command is done.
	// What it wrote before that is what we came for.
	output, _ := io.ReadAll(startedPty.Pty)

	// A terminal ends each line with a carriage return, which is of no use to a
	// caller that treats the output as text.
	return strings.ReplaceAll(string(output), "\r\n", "\n"), startedPty.Wait()
}

// SetDumbTerminalEnv tells diff renderers that we're in a very simple terminal
// that they should not expect to have much capabilities.
// Moving the cursor, clearing the screen, or querying for colors are among such
// "advanced" capabilities.
// Context: https://github.com/jesseduffield/lazygit/issues/3419
func SetDumbTerminalEnv(cmd *exec.Cmd) {
	cmd.Env = append(removeExistingTermEnvVars(cmd.Env), "TERM=dumb")
}

func removeExistingTermEnvVars(env []string) []string {
	return lo.Filter(env, func(envVar string, _ int) bool {
		return !isTermEnvVar(envVar)
	})
}

// Terminals set a variety of different environment variables
// to identify themselves to processes. This list should catch the most common among them.
func isTermEnvVar(envVar string) bool {
	return strings.HasPrefix(envVar, "TERM=") ||
		strings.HasPrefix(envVar, "TERM_PROGRAM=") ||
		strings.HasPrefix(envVar, "TERM_PROGRAM_VERSION=") ||
		strings.HasPrefix(envVar, "TERMINAL_EMULATOR=") ||
		strings.HasPrefix(envVar, "TERMINAL_NAME=") ||
		strings.HasPrefix(envVar, "TERMINAL_VERSION_")
}
