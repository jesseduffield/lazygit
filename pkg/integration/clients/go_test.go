//go:build !windows
// +build !windows

package clients

// This file allows you to use `go test` to run integration tests.
// See pkg/integration/README.md for more info.

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

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
		MaxAttempts: 2,
	})

	assert.NoError(t, err)
}

func runCmdHeadless(cmd *exec.Cmd) (int, error) {
	cmd.Env = append(
		cmd.Env,
		"HEADLESS=true",
		"TERM=xterm",
	)

	// not writing stderr to the pty because we want to capture a panic if
	// there is one. But some commands will not be in tty mode if stderr is
	// not a terminal. We'll need to keep an eye out for that.
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	// these rows and columns are ignored because internally we use tcell's
	// simulation screen. However we still need the pty for the sake of
	// running other commands in a pty.
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 300, Cols: 300})
	if err != nil {
		return -1, err
	}

	_, _ = io.Copy(io.Discard, f)

	if cmd.Wait() != nil {
		_ = f.Close()
		// return an error with the stderr output
		return cmd.Process.Pid, errors.New(stderr.String())
	}

	return cmd.Process.Pid, f.Close()
}
