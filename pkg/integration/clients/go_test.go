//go:build !windows

package clients

// This file allows you to use `go test` to run integration tests.
// See pkg/integration/README.md for more info.

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	parallelTotal := tryConvert(os.Getenv("PARALLEL_TOTAL"), 1)
	parallelIndex := tryConvert(os.Getenv("PARALLEL_INDEX"), 0)
	raceDetector := os.Getenv("LAZYGIT_RACE_DETECTOR") != ""
	// LAZYGIT_GOCOVERDIR is the directory where we write coverage files to. If this directory
	// is defined, go binaries built with the -cover flag will write coverage files to
	// to it.
	codeCoverageDir := os.Getenv("LAZYGIT_GOCOVERDIR")
	testNumber := 0

	err := components.RunTests(components.RunTestArgs{
		Tests:  tests.GetTests(utils.GetLazyRootDirectory()),
		Logf:   t.Logf,
		RunCmd: runCmdHeadless,
		TestWrapper: func(test *components.IntegrationTest, f func() error) {
			defer func() { testNumber += 1 }()
			if testNumber%parallelTotal != parallelIndex {
				return
			}

			t.Run(test.Name(), func(t *testing.T) {
				t.Parallel()
				err := f()
				assert.NoError(t, err)
			})
		},
		Sandbox:         false,
		WaitForDebugger: false,
		RaceDetector:    raceDetector,
		CodeCoverageDir: codeCoverageDir,
		InputDelay:      0,
		// Allow two attempts at each test to get around flakiness
		MaxAttempts: 1,
	})

	assert.NoError(t, err)
}

func runCmdHeadless(cmd *exec.Cmd) (int, error) {
	cmd.Env = append(
		cmd.Env,
		"LAZYGIT_HEADLESS=true",
		"TERM=xterm",
	)

	// not writing stderr to the pty because we want to capture a panic if
	// there is one. But some commands will not be in tty mode if stderr is
	// not a terminal. We'll need to keep an eye out for that.
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	// If lazygit exits but leaves behind a subprocess that inherited its stderr
	// pipe, cmd.Wait blocks waiting for that pipe to reach EOF for as long as the
	// subprocess stays alive. Unbounded, that hangs the whole test binary until
	// its global timeout fires, and the timeout throws away whatever lazygit
	// wrote to stderr before exiting (a panic, a -race report) -- the very output
	// needed to diagnose the failure. WaitDelay caps the wait: once the process
	// has exited, Wait gives the stderr goroutine at most this long to drain,
	// then closes the pipe and returns ErrWaitDelay, so the captured stderr
	// surfaces as the test error instead of being lost.
	cmd.WaitDelay = 5 * time.Second

	// these rows and columns are ignored because internally we use tcell's
	// simulation screen. However we still need the pty for the sake of
	// running other commands in a pty.
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 300, Cols: 300})
	if err != nil {
		return -1, err
	}

	// pty.StartWithSize starts lazygit in its own process group, so we can signal
	// the whole group at once. Capture the id now, while the process is alive:
	// once Wait has reaped it we can no longer look it up.
	pgid, pgidErr := syscall.Getpgid(cmd.Process.Pid)

	_, _ = io.Copy(io.Discard, f)

	waitErr := cmd.Wait()

	// On any failure -- including a WaitDelay expiry caused by a leaked
	// subprocess -- kill the whole process group so a straggler can't linger and
	// wedge a later test or pile up across a CI run. Best effort: usually the
	// group is already gone (ESRCH), and a subprocess that called setsid to
	// detach into its own group is out of reach, but WaitDelay still unblocks us.
	if waitErr != nil && pgidErr == nil {
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
	}

	if waitErr != nil {
		_ = f.Close()
		// Prefer lazygit's own stderr as the error; fall back to the wait error
		// itself (e.g. ErrWaitDelay) when it exited without printing anything.
		if stderr.Len() > 0 {
			return cmd.Process.Pid, errors.New(stderr.String())
		}
		return cmd.Process.Pid, waitErr
	}

	return cmd.Process.Pid, f.Close()
}
